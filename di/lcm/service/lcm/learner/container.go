package learner

import (
	"webank/DI/lcm/lcmconfig"

	v1core "k8s.io/api/core/v1"
	v1resource "k8s.io/apimachinery/pkg/api/resource"
)

//Container ...
type Container struct {
	Image
	Resources
	VolumeMounts  []v1core.VolumeMount
	Name, Command string //FIXME eventually get rid of command as well
	EnvVars       []v1core.EnvVar
}

//Resources ...
type Resources struct {
	CPUs, Memory, GPUs v1resource.Quantity
}

//CreateContainerSpec ...
func CreateContainerSpec(container Container) v1core.Container {
	image := GetLearnerImageForFramework(container.Image)
	resources := generateResourceRequirements(container.CPUs, container.Memory, container.GPUs)
	mounts := container.VolumeMounts
	return generateContainerSpec(container.Name, image, container.Command, container.EnvVars, resources, mounts)
}

func generateContainerSpec(name, image, cmd string, vars []v1core.EnvVar, resourceRequirements v1core.ResourceRequirements, mounts []v1core.VolumeMount) v1core.Container {

	privileged := false
	allowPrivilegeEscalation := true

	securityContext := v1core.SecurityContext{
		//Capabilities: &caps,
		AllowPrivilegeEscalation: &allowPrivilegeEscalation,
		Privileged:               &privileged,
	}

	return v1core.Container{
		Name:            name,
		Image:           image,
		ImagePullPolicy: lcmconfig.GetImagePullPolicy(),
		Command:         []string{"bash", "-c", cmd},
		Env:             vars,
		Ports: []v1core.ContainerPort{
			v1core.ContainerPort{ContainerPort: int32(22), Protocol: v1core.ProtocolTCP},
			v1core.ContainerPort{ContainerPort: int32(2222), Protocol: v1core.ProtocolTCP},
		},
		Resources:       resourceRequirements,
		VolumeMounts:    mounts,
		SecurityContext: &securityContext,
	}
}

func generateResourceRequirements(cpus, memory, gpus v1resource.Quantity) v1core.ResourceRequirements {

	var resourceGPU v1core.ResourceName = "nvidia.com/gpu"
	//if !config.GetDevicePlugin() {
	//	resourceGPU = v1core. ResourceNvidiaGPU
	//}

	return v1core.ResourceRequirements{
		Requests: v1core.ResourceList{
			v1core.ResourceCPU:    cpus,
			v1core.ResourceMemory: memory,
			resourceGPU:           gpus,
		},
		Limits: v1core.ResourceList{
			v1core.ResourceCPU:    cpus,
			v1core.ResourceMemory: memory,
			resourceGPU:           gpus,
		},
	}
}
