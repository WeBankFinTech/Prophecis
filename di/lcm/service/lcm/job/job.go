/*
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package job

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"webank/DI/commons/constants"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/lcm/service/lcm/learner"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	psPort int32 = 50051
	// FfDL Change: Specialized for FfDL
	caffeFrameworkName      string = "caffe"
	tfFrameworkName         string = "tensorflow"
	torchFrameworkName      string = "torch"
	caffe2FrameworkName     string = "caffe2"
	pytorchFrameworkName    string = "pytorch"
	customFrameworkName     string = "custom"
	h2o3FrameworkName       string = "h2o3"
	horovodFrameworkName    string = "horovod"
	pytorchMPIFrameworkName string = "pytorchmpi"
	numRetries                     = 5
	maxGPUsPerNode                 = 4

	// Not sure if these should stay or go, -sb 3/15/2018
	errCodeNormal                = "000"
	errCodeInsufficientResources = "100"
	errCodeFailedDeploy          = "101"
	errCodeFailedPS              = "102" // TODO: unused?
	errCodeImagePull             = "103"
	errFailedPodReasonUnknown    = "104"
	errCodeK8SConnection         = "200"
	errCodeEtcdConnection        = "201"
)

const learnerEntrypointFilesVolume = "learner-entrypoint-files"
const learnerEntrypointFilesPath = "/entrypoint-files"
const hostDataMountVolume = "mounted-host-data"

const hostDataMountPath = "/cosdata"
const imagePrefixKey = "image.prefix"
const mountHostDataKey = "mount.host.data"

//CreateServiceSpec ... this service will govern the statefulset
func CreateServiceSpec(name string, trainingID string, k8sClient kubernetes.Interface) *v1core.Service {
	startPort, _ := strconv.Atoi(os.Getenv(constants.DI_START_NODEPORT))
	endport, _ := strconv.Atoi(os.Getenv(constants.DI_END_NODEPORT))
	_, arr := learner.GenerateNodePort(startPort, endport, k8sClient)
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

//needs to happen before the volume creation since we muck around with the value and change the paths of data/result dir
func EnvVarsForDeployingJob(existingEnvVars []v1core.EnvVar, trainingID string, numLearners int, statefulsetName string, mountTrainingDataStoreInLearner, mountResultsStoreInLearner bool) []v1core.EnvVar {
	return learner.PopulateLearnerEnvVariablesAndLabels(existingEnvVars, trainingID, numLearners, statefulsetName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner)

}

func GenerateJobEnv(req *service.JobDeploymentRequest, service *v1core.Service) []v1core.EnvVar {
	envVars := extractEnvVarsFromDeploymentRequest(req)

	jobUser := req.UserId
	if req.ProxyUser != "" {
		jobUser = req.ProxyUser
	}
	envVars = append(envVars, v1core.EnvVar{
		Name: "Driver_Host",
		ValueFrom: &v1core.EnvVarSource{
			FieldRef: &v1core.ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	})
	envVars = append(envVars, v1core.EnvVar{
		Name:  "Spark_Driver_BlockManager_Port",
		Value: strconv.Itoa(int(service.Spec.Ports[0].NodePort)),
	})
	envVars = append(envVars, v1core.EnvVar{
		Name:  "Spark_Driver_Port",
		Value: strconv.Itoa(int(service.Spec.Ports[1].NodePort)),
	})

	envVars = append(envVars, v1core.EnvVar{
		Name:  "JOB_USER",
		Value: jobUser,
	})

	return envVars
}

func extractEnvVarsFromDeploymentRequest(req *service.JobDeploymentRequest) []v1core.EnvVar {
	var envVars []v1core.EnvVar
	for k, v := range req.EnvVars {
		envVars = append(envVars, v1core.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return envVars
}


// Wrap a sequence of commands with start and exit files.
func WrapCommands(commands []ContainerCommands) string {
	var allCommands string

	for _, command := range commands {
		allCommands += command.Cmd

	}

	allCommands += `
			echo "exit now"
			exit $TRAINING_RESULT`

	return allCommands
}

// Wrap a single command with start and exit files.
func WrapCommand(cmd string, containerName string, doCondExitWrite bool) string {

	vars := map[string]string{
		"Name": containerName,
		// "Dir":  controlFilesDirectory,
		"Cmd": cmd,
	}

	var exitWriteStr string
	if doCondExitWrite {
		exitWriteStr = `
		if [ ! -f {{.Dir}}/{{.Name}}.exit ]; then
			echo $main_cmd_status > {{.Dir}}/{{.Name}}.exit
        fi
		`
	} else {
		exitWriteStr = `
		echo $? > {{.Dir}}/{{.Name}}.exit
		`
	}

	var buf bytes.Buffer
	tmpl, _ := template.New("wrapped command").Parse(`
		# Don't repeat if already executed.
		if [ -f {{.Dir}}/{{.Name}}.exit ]; then
			while true; do sleep 1000; done
		fi
		# Wait for start signal.
		while [ ! -f {{.Dir}}/{{.Name}}.start ]; do sleep 2; done
		# Record the start time. Note: In distributed mode, this
		# file will get overwritten by each learner (this is intentional)
		date "+%s%N" | cut -b1-13 > {{.Dir}}/{{.Name}}.start_time
		{{.Cmd}} # do the actual work` + exitWriteStr +
		`while true; do sleep 2; done
	`)
	tmpl.Execute(&buf, vars)

	return buf.String()
}

func ExtendLearnerContainer(learner *v1core.Container, req *service.JobDeploymentRequest,
	logr *logger.LocLoggingEntry) {
	learnerImage := req.Framework + ":" + req.Version

	extCmd := ""
	extMount := v1core.VolumeMount{
		Name:      learnerEntrypointFilesVolume,
		MountPath: learnerEntrypointFilesPath,
	}
	learner.VolumeMounts = append(learner.VolumeMounts, extMount)

	if req.JobType == "MLPipeline"{
		mlpipelineMount := v1core.VolumeMount{
			Name:      "mlpipeline-script",
			MountPath: "/data/python/script",
		}
		learner.VolumeMounts = append(learner.VolumeMounts, mlpipelineMount)
	}

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

func doMountHostData() bool {
	mountPath := viper.GetString(mountHostDataKey)
	if mountPath == "1" || mountPath == "true" {
		return true
	}
	return false
}
