package main

import (
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

type podParameters struct {
	spaceName  string
	podName    string
	podImage   string
	podCommand []string
	cpuLimit   string
	memLimit   string
}

/*
 * create hadoop-master namespace
 */
func createNamespace(clientset *kubernetes.Clientset, spaceName string) (*apiv1.Namespace, error) {
	nc := new(apiv1.Namespace)
	ncTypeMeta := metav1.TypeMeta{Kind: "NameSpace", APIVersion: "v1"}
	nc.TypeMeta = ncTypeMeta

	nc.ObjectMeta = metav1.ObjectMeta{
		Name: spaceName,
	}

	nc.Spec = apiv1.NamespaceSpec{}
	return clientset.CoreV1().Namespaces().Create(nc)
}

/*
 * get specify namespace by name
 */
func getNamespace(clientset *kubernetes.Clientset, spaceName string) (*apiv1.Namespace, error) {
	return clientset.CoreV1().Namespaces().Get(spaceName, metav1.GetOptions{})
}

/*
 * create hadoop-master service
 */
func createHadoopMasterService(clientset *kubernetes.Clientset, spaceName string, svcName string) (*apiv1.Service, error) {
	masterSvc := new(apiv1.Service)
	svcTypeMeta := metav1.TypeMeta{Kind: "Service", APIVersion: "V1"}
	masterSvc.TypeMeta = svcTypeMeta

	svcObjectMeta := metav1.ObjectMeta{Name: svcName, Namespace: spaceName, Labels: map[string]string{"name": svcName}}
	masterSvc.ObjectMeta = svcObjectMeta

	svcServiceSpec := apiv1.ServiceSpec{
		Ports: []apiv1.ServicePort{
			apiv1.ServicePort{
				Name:       "app1",
				Port:       22,
				TargetPort: intstr.FromInt(22),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app2",
				Port:       7373,
				TargetPort: intstr.FromInt(7373),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app3",
				Port:       7946,
				TargetPort: intstr.FromInt(7946),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app4",
				Port:       9000,
				TargetPort: intstr.FromInt(9000),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app5",
				Port:       50010,
				TargetPort: intstr.FromInt(50010),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app6",
				Port:       50020,
				TargetPort: intstr.FromInt(50020),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app7",
				Port:       50070,
				TargetPort: intstr.FromInt(50070),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app8",
				Port:       50075,
				TargetPort: intstr.FromInt(50075),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app9",
				Port:       50475,
				TargetPort: intstr.FromInt(50475),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app10",
				Port:       8030,
				TargetPort: intstr.FromInt(8030),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app11",
				Port:       8031,
				TargetPort: intstr.FromInt(8031),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app12",
				Port:       8032,
				TargetPort: intstr.FromInt(8032),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app13",
				Port:       8033,
				TargetPort: intstr.FromInt(8033),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app14",
				Port:       8040,
				TargetPort: intstr.FromInt(8040),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app15",
				Port:       8042,
				TargetPort: intstr.FromInt(8042),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app16",
				Port:       8060,
				TargetPort: intstr.FromInt(8060),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app17",
				Port:       8088,
				TargetPort: intstr.FromInt(8088),
				Protocol:   "TCP",
			},
			apiv1.ServicePort{
				Name:       "app18",
				Port:       50060,
				TargetPort: intstr.FromInt(50060),
				Protocol:   "TCP",
			},
		},
		Selector:        map[string]string{"name": svcName},
		ClusterIP:       "",
		Type:            "ClusterIP",
		SessionAffinity: "None",
	}
	masterSvc.Spec = svcServiceSpec
	return clientset.CoreV1().Services(spaceName).Create(masterSvc)
}

/*
 * create hadoop master and slave pods
 */
func createHadoopPods(clientset *kubernetes.Clientset, podPara podParameters) (*apiv1.Pod, error) {
	masterPriviged := true
	masterPod := new(apiv1.Pod)
	podTypeNeta := metav1.TypeMeta{Kind: "Pod", APIVersion: "V1"}
	masterPod.TypeMeta = podTypeNeta

	podObjectMeta := metav1.ObjectMeta{Name: masterName, Namespace: spaceName, Labels: map[string]string{"name": masterName}}
	masterPod.ObjectMeta = podObjectMeta

	podSpec := apiv1.PodSpec{
		Containers: []apiv1.Container{
			apiv1.Container{
				Name:    masterName,
				Image:   masterImage,
				Command: masterCommand,
				Ports: []apiv1.ContainerPort{
					apiv1.ContainerPort{
						ContainerPort: 22,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 7373,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 7946,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 9000,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50010,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50020,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50070,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50075,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50090,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50475,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8030,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8031,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8032,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8033,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8040,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8042,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8060,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 8088,
						Protocol:      apiv1.ProtocolTCP,
					},
					apiv1.ContainerPort{
						ContainerPort: 50060,
						Protocol:      apiv1.ProtocolTCP,
					},
				},
				ImagePullPolicy: "IfNotPresent",
				SecurityContext: &apiv1.SecurityContext{
					Privileged: &masterPriviged,
				},
				Resources: apiv1.ResourceRequirements{
					Limits: apiv1.ResourceList{
						apiv1.ResourceCPU:    resource.MustParse("500m"),
						apiv1.ResourceMemory: resource.MustParse("2Gi"),
						// apiv1.ResourceStorage: resource.MustParse("50Gi"),
					},
				},
			},
		},
		RestartPolicy: apiv1.RestartPolicyAlways,
		DNSPolicy:     "ClusterFirst",
	}
	masterPod.Spec = podSpec
	return clientset.CoreV1().Pods(podPara.spaceName).Create(masterPod)
}

func main() {
	fmt.Println("come main")
	flag.Parse()
	// uses the current context in kubeconfig"k8s.io/apimachinery/pkg/api/resource"
	config, err := clientcmd.BuildConfigFromFlags("https://master.example.com:8443", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//create namesapce
	//newNamespace, err := createNamespace(clientset,"k8s-test")

	svcName := "hadoop-master"
	spaceName := "hadoop-test"

	//get namespace
	hadoopNc, err := getNamespace(clientset, spaceName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(hadoopNc.Name)

	//create servie in namespace hadoopNc
	newSvc, err := createHadoopMasterService(clientset, spaceName, svcName)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Print(newSvc.Name)
		fmt.Println("New master service create successful")
	}

	//create master pod
	masterName := "hadoop-master"
	masterImage := "172.30.7.23:5000/openshift/hadoop-master:0.1.0"
	masterCommand := []string{"bash", "-c", "/root/start-ssh-serf.sh && sleep 365d"}

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Print("New master pod create successful")
		fmt.Println(newSvcPod.Name)
	}

	//create N slave pods
	// N := 4
	// var i int
	// podName := "hadoop-slave-"
	// slaveImage := "172.30.7.23:5000/openshift/hadoop-slave:0.1.0"
	// slaveCommand := []string{"bash", "-c", "export JOIN_IP=$HADOOP_MASTER_SERVICE_HOST;/root/start-ssh-serf.sh ; ssh -o StrictHostKeyChecking=no $JOIN_IP \"/root/config.sh;/root/restart.sh\"; sleep 365d"}
	// slavePriviged := true
	// for i = 1; i <= N; i++ {
	// 	slaveName := podName + strconv.Itoa(i)

	// }
}
