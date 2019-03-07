/*
 * test develop functions and client-go api
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

/*
 * testScheduler : test random scheduler function
 */
func testScheduler(sch *Scheduler, appPath string, weightPath string) {

	fmt.Println("test scheduler functions")
	podReq, weight := readApplication(appPath, weightPath)
	podLen := len(podReq)
	// for i := 0; i < appLen; i++ {
	// 	fmt.Println(*podReq[i].resReq)
	// 	fmt.Println(weight[i])
	// }
	fmt.Println(podLen)

	// test the random scheduler function
	fmt.Println("The random scheduler result")
	randList := sch.RandomSchedule(podReq)
	for i := 0; i < podLen; i++ {
		fmt.Printf("%d ", randList[i].nodeName)
	}
	fmt.Printf("\nThe FirstFit scheduler result")
	firstFitList := sch.FirstFitSchedule(podReq)
	for i := 0; i < podLen; i++ {
		fmt.Printf("%d ", firstFitList[i].nodeName)
	}

	fmt.Printf("\nThe kubernetes scheduler result")
	kubList := sch.KubernetesSchedule(podReq)
	for i := 0; i < podLen; i++ {
		fmt.Printf("%d ", kubList[i].nodeName)
	}

	fmt.Printf("\nThe mrws scheduler result")
	mrwsList := sch.MrwsSchedule(podReq, weight)
	for i := 0; i < podLen; i++ {
		fmt.Printf("%d ", mrwsList[i].nodeName)
	}
	fmt.Println()
}

/*
 * readApplication : read the application and application matrix
 */
func readApplication(appPath string, weightPath string) (podReq []PodRequest, weight [][DIMENSION + 1]float64) {

	appFile, err := os.Open(appPath)
	if err != nil {
		panic(err.Error())
	}
	defer appFile.Close()
	appBuf := bufio.NewReader(appFile)
	for {
		line, err := appBuf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil || err == io.EOF {
			break
		}
		var podRes [DIMENSION]float64
		str := strings.Split(line, " ")
		strLen := len(str)
		for i := 0; i < strLen; i++ {
			podRes[i], _ = strconv.ParseFloat(str[i], 64)
		}
		newPod := PodRequest{
			resReq:   &podRes,
			typePod:  "hadoop",
			nodeName: -1,
		}
		podReq = append(podReq, newPod)
	}
	// read weight matrix file
	weightFile, err := os.Open(weightPath)
	if err != nil {
		panic(err.Error())
	}
	defer weightFile.Close()
	weightBuf := bufio.NewReader(weightFile)
	for {
		line, err := weightBuf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil || err == io.EOF {
			break
		}
		var weightM [DIMENSION + 1]float64
		str := strings.Split(line, " ")
		strLen := len(str)
		for i := 0; i < strLen; i++ {
			weightM[i], _ = strconv.ParseFloat(str[i], 64)
		}
		weight = append(weight, weightM)
	}
	return podReq, weight
}

/*
 * test main
 */
func main() {
	var appPath, weightPah string
	appPath = "application.txt"
	weightPah = "weight.txt"

	var reTotal = [DIMENSION]float64{2400.0, 16000.0, 1000.0, 1000.0}
	var thold float64
	thold = 0.1
	sch := &Scheduler{
		reTotal: &reTotal,
		thold:   thold,
	}
	testScheduler(sch, appPath, weightPah)

}
