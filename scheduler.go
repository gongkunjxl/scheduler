/*
 * scheduler pod调度策率的方法，根据pod的资源需求 决定pod调度到哪个节点上
 */
package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

/*
 * PHYNUM : the number of physical machines
 * define the cluster physical resources
 * 分别是物理机数量 资源量 维度 阈值

 * DIME : the dimension
 */
const PHYNUM int = 5

// DIMENSION : dimension
const DIMENSION int = 4

// Scheduler : the scheduler struct
type Scheduler struct {
	reTotal *[DIMENSION]float64
	thold   float64
}

/*
 * RandomSchedule : random scheduler method
 */
func (sch *Scheduler) RandomSchedule(podReq []PodRequest) []PodRequest {
	//the randUsed array
	var randUsed [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			randUsed[i][j] = 1.0
		}
	}
	// schedule all the pod
	podLen := len(podReq)
	for i := 0; i < podLen; i++ {
		fitInd := sch.RandomEvaluate(&randUsed, podReq[i])
		if fitInd != -1 {
			// add the used resource
			podReq[i].nodeName = fitInd
			for j := 0; j < DIMENSION; j++ {
				randUsed[fitInd][j] = randUsed[fitInd][j] + podReq[i].resReq[j]
			}

		}
	}
	// calculate the cluster resource rate
	sch.CalResourceRate(&randUsed)

	// calculate the balance value
	sch.CalClusterBalance(&randUsed, podReq)
	return podReq
}

/*
 * RandomEvaluate : calculate the physical machine idle resource
 */
func (sch *Scheduler) RandomEvaluate(randUsed *[PHYNUM][DIMENSION]float64, podReq PodRequest) int {
	var fitInd int
	fitInd = -1
	// get the physical resource idle rate
	var randIdle [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			randIdle[i][j] = (sch.reTotal[j] - randUsed[i][j] - podReq.resReq[j]) / sch.reTotal[j]
		}
	}

	// get the satisfy physical machine index
	saInd := sch.ResourceSatisfy(&randIdle)
	if saInd != nil {
		size := len(saInd)
		rd := rand.New(rand.NewSource(time.Now().UnixNano()))
		randInd := rd.Intn(size)
		fitInd = saInd[randInd]
		// fmt.Printf("%d, %d, %d\n", size, randInd, fitInd)
	}

	return fitInd
}

/*
 * FirstFitSchedule : first fit scheduler method
 */
func (sch *Scheduler) FirstFitSchedule(podReq []PodRequest) []PodRequest {
	//the randUsed array
	var firstFitUsed [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			firstFitUsed[i][j] = 1.0
		}
	}
	// schedule all the pod
	podLen := len(podReq)
	for i := 0; i < podLen; i++ {
		fitInd := sch.FirstFitEvaluate(&firstFitUsed, podReq[i])
		if fitInd != -1 {
			// add the used resource
			podReq[i].nodeName = fitInd
			for j := 0; j < DIMENSION; j++ {
				firstFitUsed[fitInd][j] = firstFitUsed[fitInd][j] + podReq[i].resReq[j]
			}

		}
	}
	// calculate the cluster resource rate
	sch.CalResourceRate(&firstFitUsed)

	// calculate the balance value
	sch.CalClusterBalance(&firstFitUsed, podReq)
	return podReq
}

/*
 * FirstFitEvaluate : calculate the physical machine idle resource
 */
func (sch *Scheduler) FirstFitEvaluate(firstFitUsed *[PHYNUM][DIMENSION]float64, podReq PodRequest) int {
	var fitInd int
	fitInd = -1
	// get the physical resource idle rate
	var firstFitIdle [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			firstFitIdle[i][j] = (sch.reTotal[j] - firstFitUsed[i][j] - podReq.resReq[j]) / sch.reTotal[j]
		}
	}

	// get the satisfy physical machine index
	saInd := sch.ResourceSatisfy(&firstFitIdle)
	if saInd != nil {
		fitInd = saInd[0]
	}
	return fitInd
}

/*
 * KubernetesSchedule : kubernetes default scheduler
 */
func (sch *Scheduler) KubernetesSchedule(podReq []PodRequest) []PodRequest {
	//the kubUsed array
	var kubUsed [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			kubUsed[i][j] = 1.0
		}
	}
	// schedule all the pod
	podLen := len(podReq)
	for i := 0; i < podLen; i++ {
		fitInd := sch.KubernetesEvaluate(&kubUsed, podReq[i])
		if fitInd != -1 {
			// add the used resource
			podReq[i].nodeName = fitInd
			for j := 0; j < DIMENSION; j++ {
				kubUsed[fitInd][j] = kubUsed[fitInd][j] + podReq[i].resReq[j]
			}

		}
	}
	// calculate the cluster resource rate
	sch.CalResourceRate(&kubUsed)

	// calculate the balance value
	sch.CalClusterBalance(&kubUsed, podReq)

	return podReq
}

/*
 * KubernetesEvaluate : calculate the physical machine idle resource and determine the desitination
 */
func (sch *Scheduler) KubernetesEvaluate(kubUsed *[PHYNUM][DIMENSION]float64, podReq PodRequest) int {
	var fitInd int
	fitInd = -1
	// get the physical resource idle rate
	var kubIdle [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			kubIdle[i][j] = (sch.reTotal[j] - kubUsed[i][j] - podReq.resReq[j]) / sch.reTotal[j]
		}
	}

	// get the satisfy physical machine index
	saInd := sch.ResourceSatisfy(&kubIdle)
	if saInd != nil {
		fitInd = sch.KubernetesMaxScore(&kubIdle, saInd)
		// fmt.Printf("%d \n", fitInd)
	}

	return fitInd
}

/*
 * KubernetesMaxScore : get the max score physical machine
 */
func (sch *Scheduler) KubernetesMaxScore(kubIdle *[PHYNUM][DIMENSION]float64, saInd []int) (maxInd int) {
	saLen := len(saInd)
	var maxScore float64
	maxScore = -1.0
	maxInd = saInd[0]
	for i := 0; i < saLen; i++ {
		scoreVal := kubIdle[saInd[i]][0]*0.5 + kubIdle[saInd[i]][1]*0.5
		if scoreVal > maxScore {
			maxScore = scoreVal
			maxInd = saInd[i]
		}
		// fmt.Printf("%d %.3f %d \n", saInd[i], maxScore, maxInd)
	}
	return maxInd
}

/*
 * MrwsSchedule : the mrws scheduler
 */
func (sch *Scheduler) MrwsSchedule(podReq []PodRequest, weightPod [][DIMENSION + 1]float64) []PodRequest {
	//the mrwsUsed array
	var mrwsUsed [PHYNUM][DIMENSION + 1]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION+1; j++ {
			mrwsUsed[i][j] = 1.0
		}
	}
	// schedule all the pod
	podLen := len(podReq)
	for i := 0; i < podLen; i++ {
		fitInd := sch.MrwsEvaluate(&mrwsUsed, podReq[i], &weightPod[i])
		if fitInd != -1 {
			// add the used resource
			podReq[i].nodeName = fitInd
			for j := 0; j < DIMENSION; j++ {
				mrwsUsed[fitInd][j] = mrwsUsed[fitInd][j] + podReq[i].resReq[j]
			}
			mrwsUsed[fitInd][DIMENSION] = mrwsUsed[fitInd][DIMENSION] + 1.0
		}
	}
	// calculate the cluster resource rate
	var calMrwsUsed [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			calMrwsUsed[i][j] = mrwsUsed[i][j]
		}
	}
	sch.CalResourceRate(&calMrwsUsed)

	// calculate the balance value
	sch.CalClusterBalance(&calMrwsUsed, podReq)

	return podReq
}

/*
 * MrwsEvaluate : calculate the physical machine idle resource and determine the desitination
 * mrwsUsed: the resource used , weightPod: the pod weight, podReq: the pod resource
 */
func (sch *Scheduler) MrwsEvaluate(mrwsUsed *[PHYNUM][DIMENSION + 1]float64, podReq PodRequest, weightPod *[DIMENSION + 1]float64) int {
	var fitInd int
	fitInd = -1
	// get the physical resource and pod idle rate
	var mrwsIdle [PHYNUM][DIMENSION]float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			mrwsIdle[i][j] = (sch.reTotal[j] - mrwsUsed[i][j] - podReq.resReq[j]) / sch.reTotal[j]
		}
	}
	var podSum float64
	var podIdle [PHYNUM]float64
	for i := 0; i < PHYNUM; i++ {
		podSum = podSum + mrwsUsed[i][DIMENSION]
	}
	for i := 0; i < PHYNUM; i++ {
		podIdle[i] = 1.0 - mrwsUsed[i][DIMENSION]/podSum

	}

	// get the satisfy physical machine index and calculate the max value node
	saInd := sch.ResourceSatisfy(&mrwsIdle)
	if saInd != nil {
		saLen := len(saInd)
		//calculate the satisfy index physical machine podMean and resMean

		var resVal [DIMENSION]float64 // cal the sum and mean value
		var resMean [DIMENSION]float64
		var podVal float64
		for i := 0; i < saLen; i++ {
			podVal = podVal + podIdle[saInd[i]]
			for j := 0; j < DIMENSION; j++ {
				resVal[j] = resVal[j] + mrwsIdle[saInd[i]][j]
			}
		}
		podMean := podVal / podSum
		for j := 0; j < DIMENSION; j++ {
			resMean[j] = resVal[j] / (float64)(saLen)
		}

		var maxScore float64
		maxScore = -1.0
		fitInd = saInd[0]
		var bi, vi float64
		for i := 0; i < saLen; i++ {
			vi = 0.0
			bi = 0.0
			for j := 0; j < DIMENSION; j++ {
				vi = vi + mrwsIdle[saInd[i]][j]*weightPod[j]
				bi = bi + (mrwsIdle[saInd[i]][j]/resMean[j])*weightPod[DIMENSION]
			}
			vi = vi + podIdle[i]*weightPod[DIMENSION]
			bi = bi + (podIdle[i]/podMean)*weightPod[DIMENSION]
			// fmt.Printf("vi and bi %.3f %.3f \n", vi, bi)
			// bi = 0.0
			scoreVi := vi + bi
			if scoreVi > maxScore {
				fitInd = saInd[i]
				maxScore = scoreVi
			}
		}
		// fmt.Printf("%.3f  %d \n", maxScore, fitInd)
	}
	return fitInd
}

/*
 * ResourceSatisfy : judge the physical machine satisfy resource request
 */
func (sch *Scheduler) ResourceSatisfy(phyIdle *[PHYNUM][DIMENSION]float64) []int {
	var saInd []int
	for i := 0; i < PHYNUM; i++ {
		flag := true
		for j := 0; j < DIMENSION; j++ {
			if phyIdle[i][j] < sch.thold {
				flag = false
				break
			}
		}
		if flag {
			saInd = append(saInd, i)
		}
	}
	return saInd
}

/*
 * CalResourceRate : calculate the physical machine resource rate just print
 */
func (sch *Scheduler) CalResourceRate(podUsed *[PHYNUM][DIMENSION]float64) {
	for i := 0; i < PHYNUM; i++ {
		fmt.Printf("%s", "node"+strconv.Itoa(i)+": ")
		for j := 0; j < DIMENSION; j++ {
			fmt.Printf("%.3f ", podUsed[i][j]/sch.reTotal[j])
		}
		fmt.Println()
	}
}

/*
 * CalClusterBalance : calculate the balance value just print
 */
func (sch *Scheduler) CalClusterBalance(podUsed *[PHYNUM][DIMENSION]float64, podReq []PodRequest) {
	//cal the pod sum and used rate
	podLen := len(podReq)
	var podNum [PHYNUM]int
	var podSum int
	for i := 0; i < podLen; i++ {
		if podReq[i].nodeName != -1 {
			podSum++
			podNum[podReq[i].nodeName]++
		}
	}

	var podIdle [PHYNUM]float64
	var resIdle [PHYNUM][DIMENSION]float64
	var podVal float64
	var resVal [DIMENSION]float64 // cal the sum and mean value

	for i := 0; i < PHYNUM; i++ {
		podIdle[i] = 1.0 - (float64)(podNum[i])/(float64)(podSum)
		podVal = podVal + podIdle[i]
		for j := 0; j < DIMENSION; j++ {
			resIdle[i][j] = (sch.reTotal[j] - podUsed[i][j]) / sch.reTotal[j]
			resVal[j] = resVal[j] + resIdle[i][j]
		}
	}
	// cal the balance value
	podMean := podVal / (float64)(podSum)
	var resMean [DIMENSION]float64
	for j := 0; j < DIMENSION; j++ {
		resMean[j] = resVal[j] / (float64)(PHYNUM)
	}
	var baIdle float64
	for i := 0; i < PHYNUM; i++ {
		for j := 0; j < DIMENSION; j++ {
			baIdle = baIdle + math.Pow((resIdle[i][j]-resMean[j]), 2)
		}
		baIdle = baIdle + math.Pow((podIdle[i]-podMean), 2)
	}
	baIdle = math.Sqrt(baIdle)
	fmt.Printf("The balance value is %.3f \n", baIdle)
}
