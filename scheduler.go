/*
 * scheduler pod调度策率的方法，根据pod的资源需求 决定pod调度到哪个节点上
 */
package main

import (
	"math/rand"
	"strconv"
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
	reTotal []float64
	thold   float64
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
			//
			podReq[i].nodeName = "node" + strconv.Itoa(fitInd) + ".example.com"
			for j := 0; j < DIMENSION; j++ {
				randUsed[fitInd] = randUsed[fitInd] + podReq[i][j]
			}
		}
	}
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
		randInd := rand.Intn(size)
		fitInd = saInd[randInd]
	}

	return fitInd
}
