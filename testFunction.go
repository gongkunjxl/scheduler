/*
 * test develop functions and client-go api
 */
package main

import (
	"fmt"
)

/*
 * testScheduler : test random scheduler function
 */
func testScheduler(sch *Scheduler) {
	// build the test sample for random scheduler
	var podRes = [DIMENSION]float64{300.0, 2000.0, 200.0, 100.0}
	var podNum int
	podNum = 13
	var podReq []PodRequest
	for i := 0; i < podNum; i++ {
		newPod := PodRequest{
			resReq:   &podRes,
			typePod:  "hadoop",
			nodeName: -1,
		}
		podReq = append(podReq, newPod)
	}

	fmt.Println("test scheduler functions")
	// test the random scheduler function
	//podList := sch.RandomSchedule(podReq)
	podList := sch.KubernetesSchedule(podReq)

	for i := 0; i < podNum; i++ {
		fmt.Println(podList[i].nodeName)
	}
}

/*
 * test main
 */
func main() {
	var reTotal = [DIMENSION]float64{2400.0, 16000.0, 1000.0, 1000.0}
	var thold float64
	thold = 0.1
	sch := &Scheduler{
		reTotal: &reTotal,
		thold:   thold,
	}
	testScheduler(sch)

}
