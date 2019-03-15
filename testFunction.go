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
	podReq, _ := readApplication(appPath, weightPath)
	podLen := len(podReq)
	// for i := 0; i < appLen; i++ {
	// 	fmt.Println(*podReq[i].resReq)
	// 	fmt.Println(weight[i])
	// }
	fmt.Println(podLen)

	// test the random scheduler function
	fmt.Println("The random scheduler result")
	// podList := sch.RandomSchedule(podReq)
	podList := sch.KubernetesSchedule(podReq)
	// podList := sch.MrwsSchedule(podReq, weight)

	for i := 0; i < podLen; i++ {
		fmt.Printf("%d ", podList[i].nodeName)
	}
	fmt.Println()
	var spec = [3]int{-1, -1, -1} // specify master or slave podReq
	var masterReq []PodRequest
	var slaveReq []PodRequest
	for i := 0; i < podLen; i++ {
		var j int
		j = 0
		for j = 0; j < 3; j++ {
			if podList[i].typePod == spec[j] {
				break
			}
		}
		if j == 3 {
			masterReq = append(masterReq, podList[i])
			spec[podList[i].typePod-1] = podList[i].typePod
		} else {
			slaveReq = append(slaveReq, podList[i])
		}
	}
	fmt.Printf("Master pod request length %d \n", len(masterReq))
	fmt.Printf("Slave pod request length %d \n", len(slaveReq))

	// podbyNamenode create master pod
	var typePod = []string{"mpi", "spark", "hadoop"}
	var nodeName [PHYNUM]string
	var typeMod string
	typeMod = "master"
	nodeName[0] = "master.example.com"
	for i := 1; i < PHYNUM; i++ {
		nodeName[i] = "node" + strconv.Itoa(i) + ".example.com"
	}
	masterCommand := []string{"bash", "-c", "/root/start-ssh-serf.sh && sleep 365d"}
	pyn := PodByName{
		typePod:  typePod,
		nodeName: nodeName,
		command:  masterCommand,
	}
	pyn.CreatePodByRequest(masterReq, typeMod)

	// podByNamenode create slave pod
	typeMod = "slave"
	pyn.CreatePodByRequest(slaveReq, typeMod)

}

/*
 * readApplication : read the application and application matrix
 */
func readApplication(appPath string, weightPath string) (podReq []PodRequest, weight [][DIMENSION + 1]float64) {

	var typeCommand = [][]string{{"bash", "-c", "export JOIN_IP=$MPI_MASTER_SERVICE_HOST && /root/start-ssh-serf.sh && sleep 365d"}, {"bash", "-c", "export JOIN_IP=$SPARK_MASTER_SERVICE_HOST && /root/start-ssh-serf.sh && sleep 365d"}, {"bash", "-c", "export JOIN_IP=$HADOOP_MASTER_SERVICE_HOST && /root/start-ssh-serf.sh && sleep 365d"}}
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
		strLen := len(str) - 1
		for i := 0; i < strLen; i++ {
			podRes[i], _ = strconv.ParseFloat(str[i], 64)
		}
		podType, _ := strconv.Atoi(str[strLen])
		// podType = 1
		newPod := PodRequest{
			resReq:   &podRes,
			typePod:  podType,
			nodeName: -1,
		}
		newPod.command = typeCommand[newPod.typePod-1]
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

	var reTotal = [DIMENSION]float64{4000.0, 16000.0, 1000.0, 100.0}
	var thold float64
	thold = 0.15
	sch := &Scheduler{
		reTotal: &reTotal,
		thold:   thold,
	}
	testScheduler(sch, appPath, weightPah)

}
