package learner

import (
	"os"
	"strconv"
	"sync"
	"webank/DI/commons/constants"

	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const portArrayLen = 10001

var mutex sync.Mutex
var portArray [portArrayLen]bool

//CreatePodSpec ...
// FIXME MLSS Change: add nodeSelectors to deployment/sts pods
func CreatePodSpec(containers []v1core.Container, volumes []v1core.Volume, labels map[string]string, nodeSelector map[string]string, imagePullSecret string) v1core.PodTemplateSpec {
	labels["service"] = "dlaas-learner" //label that denies ingress/egress
	automountSeviceToken := false
	return v1core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Annotations: map[string]string{
				"scheduler.alpha.kubernetes.io/tolerations": `[ { "key": "dedicated", "operator": "Equal", "value": "gpu-task" } ]`,
				"scheduler.alpha.kubernetes.io/nvidiaGPU":   `{ "AllocationPriority": "Dense" }`,
			},
		},
		Spec: v1core.PodSpec{
			Containers: containers,
			Volumes:    volumes,
			ImagePullSecrets: []v1core.LocalObjectReference{
				v1core.LocalObjectReference{
					Name: imagePullSecret,
				},
			},
			Tolerations: []v1core.Toleration{
				v1core.Toleration{
					Key:      "dedicated",
					Operator: v1core.TolerationOpEqual,
					Value:    "gpu-task",
					Effect:   v1core.TaintEffectNoSchedule,
				},
			},
			// For FfDL we don't need any constraint on GPU type.
			// NodeSelector:                 nodeSelector,
			AutomountServiceAccountToken: &automountSeviceToken,
			// FIXME MLSS Change: but for MLSS we need constraints on GPU type & Nodes zones you know?
			NodeSelector: nodeSelector,
		},
	}
}
func CreatePodSpecForJob(containers []v1core.Container, volumes []v1core.Volume, labels map[string]string, nodeSelector map[string]string, imagePullSecret string) v1core.PodTemplateSpec {
	labels["service"] = "dlaas-learner" //label that denies ingress/egress
	automountSeviceToken := false
	podSpec := v1core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Annotations: map[string]string{
				"scheduler.alpha.kubernetes.io/tolerations": `[ { "key": "dedicated", "operator": "Equal", "value": "gpu-task" } ]`,
				"scheduler.alpha.kubernetes.io/nvidiaGPU":   `{ "AllocationPriority": "Dense" }`,
			},
		},
		Spec: v1core.PodSpec{
			Containers: containers,
			Volumes:    volumes,
			ImagePullSecrets: []v1core.LocalObjectReference{
				v1core.LocalObjectReference{
					Name: imagePullSecret,
				},
			},
			Tolerations: []v1core.Toleration{
				v1core.Toleration{
					Key:      "dedicated",
					Operator: v1core.TolerationOpEqual,
					Value:    "gpu-task",
					Effect:   v1core.TaintEffectNoSchedule,
				},
			},
			// For FfDL we don't need any constraint on GPU type.
			// NodeSelector:                 nodeSelector,
			AutomountServiceAccountToken: &automountSeviceToken,
			// FIXME MLSS Change: but for MLSS we need constraints on GPU type & Nodes zones you know?
			NodeSelector:  nodeSelector,
			RestartPolicy: "Never",
		},
	}
	labels["service"] = "dlaas-learner"
	labels["service"] = "dlaas-learner"
	return podSpec
}

//CreateStatefulSetSpecForLearner ...
func CreateStatefulSetSpecForLearner(name, servicename string, replicas int, podTemplateSpec v1core.PodTemplateSpec) *v1beta1.StatefulSet {
	var replicaCount = int32(replicas)
	revisionHistoryLimit := int32(0) //https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#clean-up-policy

	return &v1beta1.StatefulSet{

		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: podTemplateSpec.Labels,
		},
		Spec: v1beta1.StatefulSetSpec{
			ServiceName:          servicename,
			Replicas:             &replicaCount,
			Template:             podTemplateSpec,
			RevisionHistoryLimit: &revisionHistoryLimit, //we never rollback these
			//PodManagementPolicy: v1beta1.ParallelPodManagement, //using parallel pod management in stateful sets to ignore the order. not sure if this will affect the helper pod since any pod in learner can come up now
		},
	}
}

//FIXME mlss 1.9.0 sts -> job
func CreateJobSpecForLearner(name string, podTemplateSpec v1core.PodTemplateSpec) *v1.Job {

	labels := podTemplateSpec.Labels
	labels[constants.ENVIR] = os.Getenv(constants.ENVIR_UPRER)
	labels[constants.JOBTYPE] = constants.SINGLE
	backoffLimit := int32(0)
	return &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1.JobSpec{
			Template:     podTemplateSpec,
			BackoffLimit: &backoffLimit,
		},
	}
}

//CreateServiceSpec ... this service will govern the statefulset
func CreateServiceSpec(name string, trainingID string, k8sClient kubernetes.Interface) *v1core.Service {
	startPort, _ := strconv.Atoi(os.Getenv(constants.DI_START_NODEPORT))
	endport, _ := strconv.Atoi(os.Getenv(constants.DI_END_NODEPORT))
	_, arr := GenerateNodePort(startPort, endport, k8sClient)
	return &v1core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: trainingID,
			Labels: map[string]string{
				"job-name": trainingID,
			},
		},
		Spec: v1core.ServiceSpec{
			Selector: map[string]string{"job-name": trainingID},
			Type:     v1core.ServiceTypeNodePort,
			Ports: []v1core.ServicePort{
				v1core.ServicePort{
					Name:     "sparkcontextport",
					Protocol: v1core.ProtocolTCP,
					Port:     int32(arr[0]),
					NodePort: int32(arr[0]),
				},
				v1core.ServicePort{
					Name:     "sparkblockmanagerport",
					Protocol: v1core.ProtocolTCP,
					Port:     int32(arr[1]),
					NodePort: int32(arr[1]),
				},
			},
		},
	}
}

func GenerateNodePort(startPort int, endPort int, k8sClient kubernetes.Interface) (error, [2]int) {
	mutex.Lock()
	defer mutex.Unlock()
	// clusterService, _, err := k8sClient.CoreV1().Services().List()
	clusterService, err := k8sClient.CoreV1().Services("").List(metav1.ListOptions{})

	// portLen := (endPort - startPort) + 1
	// portArray = make([]bool, portLen, portLen)
	nodePortArray := [2]int{-1, -1}

	if err != nil {
		return err, nodePortArray
	}

	for _, serviceItem := range clusterService.Items {
		if serviceItem.Spec.Type == "NodePort" {
			usedNodeportList := serviceItem.Spec.Ports
			for _, usedNodeportItem := range usedNodeportList {
				usedNodeport := int(usedNodeportItem.NodePort)
				if usedNodeport >= startPort && usedNodeport <= endPort {
					portArray[usedNodeport-startPort] = true
				}
			}
		}
	}
	flag := 0
	for i, value := range portArray {
		if value != true {
			portArray[i] = true
			nodePortArray[flag] = i + startPort + 1
			flag = flag + 1
			if flag >= 2 {
				break
			}
		}
	}



	// mutex.Unlock()
	return err, nodePortArray
}
