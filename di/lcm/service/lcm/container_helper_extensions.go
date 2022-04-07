package lcm

import (
	"webank/DI/commons/service"

	"github.com/spf13/viper"

	"fmt"
	"os"

	v1core "k8s.io/api/core/v1"

	//"github.com/sirupsen/logrus"
	"webank/DI/commons/logger"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const learnerEntrypointFilesVolume = "learner-entrypoint-files"
const learnerEntrypointFilesPath = "/entrypoint-files"
const hostDataMountVolume = "mounted-host-data"
const learnerLogFilesVolume = "log-dir"

//const hostDataMountPath = "/host-data"
const hostDataMountPath = "/cosdata"
const imagePrefixKey = "image.prefix"
const mountHostDataKey = "mount.host.data"

func extendLearnerContainer(learner *v1core.Container, req *service.JobDeploymentRequest,
	logr *logger.LocLoggingEntry) {

	learnerImage := "__unknown__"

	switch req.Framework {
	case caffeFrameworkName:
		learnerImage = "bvlc/caffe:" + req.Version
	case tfFrameworkName:
		learnerImage = "tensorflow/tensorflow:" + req.Version
	case caffe2FrameworkName:
		learnerImage = "caffe2ai/caffe2:" + req.Version
	case pytorchFrameworkName:
		learnerImage = "pytorch/pytorch:" + req.Version
	case h2o3FrameworkName:
		learnerImage = "opsh2oai/h2o3-ffdl:" + req.Version
	case horovodFrameworkName:
		learnerImage = "uber/horovod:" + req.Version
	case pytorchMPIFrameworkName:
		learnerImage = "tomcli/pytorch:" + req.Version
	case customFrameworkName:
		learnerImage = req.Version
	default:
		// TODO!
		// FIXME MLSS Change: user specified docker image
		learnerImage = req.Framework + ":" + req.Version
	}

	extCmd := ""
	// FIXME MLSS Change: tfjob
	if req.JobType == "dist-tf" {
		extCmd = "echo 'tfjob start';"
	} else if req.Framework != horovodFrameworkName && req.Framework != pytorchMPIFrameworkName {
		extCmd = "export PATH=/usr/local/bin/:$PATH; cp " + learnerEntrypointFilesPath + "/*.sh /usr/local/bin/; chmod +x /usr/local/bin/*.sh;"
	} else {
		extCmd = "export PATH=/usr/local/bin/:$PATH; cp " + learnerEntrypointFilesPath + "/*.sh /usr/local/bin/; chmod +x /usr/local/bin/*.sh; mv /usr/local/bin/train-horovod.sh /usr/local/bin/train.sh;"
	}
	extMount := v1core.VolumeMount{
		Name:      learnerEntrypointFilesVolume,
		MountPath: learnerEntrypointFilesPath,
	}
	learner.VolumeMounts = append(learner.VolumeMounts, extMount)

	//LogMount
	//logMount := v1core.VolumeMount{
	//	Name:  learnerLogFilesVolume,
	//	MountPath: "job-logs"}
	//learner.VolumeMounts = append(learner.VolumeMounts, logMount)

	if doMountHostData() {
		hostMount := v1core.VolumeMount{
			Name:      hostDataMountVolume,
			MountPath: hostDataMountPath,
		}
		learner.VolumeMounts = append(learner.VolumeMounts, hostMount)

		yamlbytes, err := yaml.Marshal(learner.VolumeMounts)
		if err != nil {
			logr.WithError(err).Errorf("Could not marshal volume mounts for diagnosis")
		}
		fmt.Printf("-------------------------\n")
		fmt.Printf("learner.VolumeMounts: %s", string(yamlbytes))
	}

	learner.Image = learnerImage
	learner.Command[2] = extCmd + learner.Command[2]
	if req.Framework == horovodFrameworkName || req.Framework == pytorchMPIFrameworkName {
		RSApub := v1core.EnvVar{Name: "RSA_PUB_KEY", ValueFrom: &v1core.EnvVarSource{
			SecretKeyRef: &v1core.SecretKeySelector{Key: "RSA_PUB_KEY", LocalObjectReference: v1core.LocalObjectReference{Name: "rsa-keys"}}}}
		RSApri := v1core.EnvVar{Name: "RSA_PRI_KEY", ValueFrom: &v1core.EnvVarSource{
			SecretKeyRef: &v1core.SecretKeySelector{Key: "RSA_PRI_KEY", LocalObjectReference: v1core.LocalObjectReference{Name: "rsa-keys"}}}}
		learner.Env = append(learner.Env, RSApub)
		learner.Env = append(learner.Env, RSApri)
	}
}

func extendLearnerVolumes(volumeSpecs []v1core.Volume, logr *logger.LocLoggingEntry) []v1core.Volume {
	logr.Debugf("extendLearnerVolumes entry")

	// learner entrypoint files volume
	//learnerEntrypointFilesVolume := v1core.Volume{
	//	Name: learnerEntrypointFilesVolume,
	//	VolumeSource: v1core.VolumeSource{
	//		ConfigMap: &v1core.ConfigMapVolumeSource{
	//			LocalObjectReference: v1core.LocalObjectReference{
	//				Name: learnerEntrypointFilesVolume,
	//			},
	//		},
	//	},
	//}
	//volumeSpecs = append(volumeSpecs, learnerEntrypointFilesVolume)

	shouldMountHostData := doMountHostData()

	logrus.Debugf("shouldMountHostData: %t", shouldMountHostData)

	if shouldMountHostData {
		hostDataVolume := v1core.Volume{
			Name: hostDataMountVolume,
			VolumeSource: v1core.VolumeSource{
				HostPath: &v1core.HostPathVolumeSource{
					Path: hostDataMountPath,
				},
			},
		}
		logrus.Debugf("created host mount volume spec: %v+", hostDataVolume)
		volumeSpecs = append(volumeSpecs, hostDataVolume)


	}
	logr.Debugf("extendLearnerVolumes exit")
	return volumeSpecs
}

func controllerImageNameExtended(dockerRegistry string, servicesTag string) string {
	imagePrefix := getImagePrefix()
	//	return fmt.Sprintf("%s/%scontroller:%s", dockerRegistry, imagePrefix, servicesTag)
	// FIXME MLSS Temporary Change: use fixed lhelper image
	if len(imagePrefix) > 0 {
		return fmt.Sprintf("%s/%s:%s", dockerRegistry, imagePrefix, servicesTag)
	} else {
		return fmt.Sprintf("%s:%s", dockerRegistry, servicesTag)
	}
}

func jobmonitorImageNameExtended(dockerRegistry string, jmTag string) string {
	imagePrefix := getImagePrefix()
	//	return fmt.Sprintf("%s/%sjobmonitor:%s", dockerRegistry, imagePrefix, jmTag)
	// FIXME MLSS Temporary Change: use fixed jobmonitor image
	if len(imagePrefix) > 0 {
		return fmt.Sprintf("%s/%s:%s", dockerRegistry, imagePrefix, jmTag)
	} else {
		return fmt.Sprintf("%s:%s", dockerRegistry, jmTag)
	}
}

func getImagePrefix() string {
	return viper.GetString(imagePrefixKey)
}

func doMountHostData() bool {
	mountPath := viper.GetString(mountHostDataKey)
	if mountPath == "1" || mountPath == "true" {
		return true
	}
	return false
}

func exists(path string) bool {
	_, err := os.Stat(hostDataMountPath)
	return err == nil
}
