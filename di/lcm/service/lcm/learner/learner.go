package learner

import (
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"webank/DI/commons/constants"
)

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
			NodeSelector:  nodeSelector,
			RestartPolicy: "Never",
		},
	}
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
func CreateServiceSpec(name string, trainingID string) *v1core.Service {

	return &v1core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"training_id": trainingID,
			},
		},
		Spec: v1core.ServiceSpec{
			Selector: map[string]string{"training_id": trainingID},
			Ports: []v1core.ServicePort{
				v1core.ServicePort{
					Name:     "ssh",
					Protocol: v1core.ProtocolTCP,
					Port:     22,
				},
				v1core.ServicePort{
					Name:     "tf-distributed",
					Protocol: v1core.ProtocolTCP,
					Port:     2222,
				},
			},
			ClusterIP: "None",
		},
	}
}
