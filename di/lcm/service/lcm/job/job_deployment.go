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
	"fmt"
	"os"

	"webank/DI/lcm/service/lcm/learner"

	"webank/DI/commons/constants"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"

	v1core "k8s.io/api/core/v1"

	stringsUtil "strings"
)

type ContainerCommands struct {
	Cmd       string
	Container string
}

func VolumesForJob(req *service.JobDeploymentRequest, learnerEnvVars []v1core.EnvVar, logging *logger.LocLoggingEntry) learner.Volumes {

	//HostPath
	hostPath := req.DataStores[2].Connection["path"]


	volumesStruct := learner.Volumes{}
	//hostPath := req.EnvVars["DATA_STORE_PATH"]
	mountPath := getValue(learnerEnvVars, "DATA_DIR")
	// While this is the unadulterated value
	subPath := req.EnvVars["DATA_DIR"]
	fmt.Printf("(data) hostPath=%s\nmountPath=%s\nsubPath=%s\n", hostPath, mountPath, subPath)

	//req.DataPath


	// 1. Data Dir
	volumesStruct.TrainingData = &learner.COSVolume{
		VolumeType: "HostPath",
		ID:         "inputmount-" + req.Name,
		//HostPath: hostPath,
		HostPath: hostPath + "/" + subPath,

		MountSpec: learner.VolumeMountSpec{
			Name:      "workspace",
			MountPath: mountPath,
			SubPath: "",
		},
	}
	resultStoreType := "HostPath"
	logging.Debugf("TrainingData volume request: %v+", volumesStruct.TrainingData)
	logging.Debugf("debug_volumesForLearner DataStoreHostMountVolume_resultStoreType: %v", resultStoreType)

	// 2. Result Dir
	mountPath = getValue(learnerEnvVars, "RESULT_DIR")
	subPath = req.EnvVars["RESULT_DIR"]
	fmt.Printf("(result) hostPath=%s\nmountPath=%s\nsubPath=%s\n", hostPath, mountPath, subPath)
	volumesStruct.ResultsDir = &learner.COSVolume{
		VolumeType: "HostPath",
		ID:         "outputmount-" + req.Name,
		HostPath: hostPath + "/" + subPath,
		MountSpec: learner.VolumeMountSpec{
			Name:      "outputmount-" + req.Name,
			MountPath: mountPath,
			SubPath: "",
		},
	}
	logging.Debugf("ResultsDir volume request: %v+", volumesStruct.ResultsDir)


	// 4. HADOOP CONFIG DiR
	logging.Debugf("debug_volumesForLearner jobNamespace: %v", req.JobNamespace)
	if stringsUtil.Split(req.JobNamespace, "-")[2] == "bdap" {
		volumesStruct.AppComConfig = &learner.COSVolume{
			VolumeType: "HostPath",
			ID:         "bdap-appcom-config",
			HostPath:   "/appcom/config",
			MountSpec: learner.VolumeMountSpec{
				Name:      "bdap-appcom-config",
				MountPath: "/appcom/config",
			},
		}
		volumesStruct.AppComInstall = &learner.COSVolume{
			VolumeType: "HostPath",
			ID:         "bdap-appcom-install",
			HostPath:   "/appcom/Install",
			MountSpec: learner.VolumeMountSpec{
				Name:      "bdap-appcom-install",
				MountPath: "/appcom/Install",
			},
		}
		volumesStruct.JDK = &learner.COSVolume{
			VolumeType: "HostPath",
			ID:         "jdk",
			HostPath:   "/nemo/jdk1.8.0_141",
			MountSpec: learner.VolumeMountSpec{
				Name:      "jdk",
				MountPath: "/nemo/jdk1.8.0_141",
			},
		}
	}


	//5. Mount MLPipeline Script
	if req.JobType == "MLPipeline" {
		volumesStruct.Script = &learner.ConfigMapVolume{
			VolumeType: "ConfigMap",
			Name:       "mlpipeline-script",
			MountSpec: learner.VolumeMountSpec{
				Name:      "mlpipeline-script",
				MountPath: "/data/python/script",
				SubPath:   "",
			},
		}
	}
	///job-logs/${TRAINING_ID}.log
	//6. Log Dir
	volumesStruct.LogDir = &learner.HostPathVolume{
		VolumeType: resultStoreType,
		ID:         "log-dir",
		HostPath: os.Getenv(constants.FLUENT_BIT_LOG_PATH) + os.Getenv(constants.ENVIR_UPRER),
		MountSpec: learner.VolumeMountSpec{
			Name:      "log-dir",
			MountPath: "/logs/",
		},
	}


	//workspace
	volumesStruct.WorkDir = &learner.HostPathVolume{
		VolumeType: "HostPath",
		ID:         "work-" + req.Name,
		HostPath: hostPath ,
		MountSpec: learner.VolumeMountSpec{
			Name:      "work-" + req.Name,
			MountPath: "/workspace",
			SubPath: "",
		},
	}

	return volumesStruct
}

//needs to happen before the volume creation since we muck around with the value and change the paths of data/result dir
func EnvVarsForJob(existingEnvVars []v1core.EnvVar, trainingID string, numLearners int, statefulsetName string, mountTrainingDataStoreInLearner, mountResultsStoreInLearner bool) []v1core.EnvVar {
	return learner.PopulateLearnerEnvVariablesAndLabels(existingEnvVars, trainingID, numLearners, statefulsetName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner)
}

// Return the value of the named environment variable.
func getValue(envVars []v1core.EnvVar, name string) string {
	value := ""
	for _, ev := range envVars {
		if ev.Name == name {
			value = ev.Value
			break
		}
	}
	return value
}
