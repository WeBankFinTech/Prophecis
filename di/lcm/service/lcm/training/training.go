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
package lcm_training

import (
	"fmt"
	"path"
	"strconv"

	"webank/DI/lcm/service/lcm/learner"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"

	"golang.org/x/net/context"

	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	stringsUtil "strings"
)

// PodLevelLogDir represents the place to store the per-learner logs.
const PodLevelLogDir = "/job/logs"

//Training ...
type Training interface {
	Start() error
	//Halt() error
	//Stop() error
}

type training struct {
	ctx        context.Context
	k8sClient  kubernetes.Interface
	req        *service.JobDeploymentRequest
	trainingID string
	learner    trainingDefinition
	logging       *logger.LocLoggingEntry
}

type trainingResource struct {
	service   *v1core.Service
	job       *v1.Job
	namespace string
}

type trainingDefinition struct {
	secrets                                                     []*v1core.Secret
	volumes                                                     []v1core.Volume
	volumeMounts                                                []v1core.VolumeMount
	envVars                                                     []v1core.EnvVar
	name                                                        string
}


// func NewLCMTraining(ctx context.Context, k8sClient kubernetes.Interface, req *service.JobDeploymentRequest, log *logger.LocLoggingEntry) Training {
// 	const cosMountDriverName = "ibm/ibmc-s3fs"
// 	const cosMountType = "mount_cos"
// 	learnerName := fmt.Sprintf("learner-%s", req.Name)
// 	helperName := fmt.Sprintf("lhelper-%s", req.Name)
// 	numLearners := int(req.GetResources().Learners)
// 	if numLearners < 1 {
// 		numLearners = 1
// 	}

// 	dataStoreType := req.EnvVars["DATA_STORE_TYPE"]
// 	resultStoreStype := req.EnvVars["RESULT_STORE_TYPE"]
// 	// FIXME MLSS Change: v_1.5.1 added_de
// 	workStoreType := req.EnvVars["WORK_STORE_TYPE"]

// 	mountTrainingDataStoreInLearner := !(dataStoreType == learner.DataStoreTypeS3)
// 	mountResultsStoreInLearner := !(resultStoreStype == learner.DataStoreTypeS3)

// 	mountWorkStoreInLearner := !(workStoreType == learner.DataStoreTypeS3)

// 	logr := log.WithFields(logrus.Fields{
// 		"learner_name": learnerName,
// 		"helper_name":  helperName,
// 		"mounted_cos":  mountResultsStoreInLearner && mountTrainingDataStoreInLearner,
// 	})
// 	logr.Debugf("debug_for_param ParameterServer lcm newTraining pss: %v, psCpu: %v, psImage: %v, psMemory: %v", req.PSs, req.PSCPU, req.PSImage, req.PSMemory)
// 	logr.Debugf("extractEnvVarsFromDeploymentRequest(req): %+v", req.EnvVars)

// 	envVarsFromDeploymentRequest := extractEnvVarsFromDeploymentRequest(req) //shared across all containers of training
// 	logr.Debugf("envVarsFromDeploymentRequest: %+v", envVarsFromDeploymentRequest)

// 	envvarsForLearner := envVarsForDeployingLearner(envVarsFromDeploymentRequest, req.TrainingId,
// 		numLearners, learnerName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner) //only for learner
// 	logr.Debugf("NewLCMTraining envvarsForLearner: %+v", envvarsForLearner)

// 	learnerVolumes := volumesForLearner(req, envvarsForLearner, mountTrainingDataStoreInLearner, mountResultsStoreInLearner, mountWorkStoreInLearner, logr)
// 	learnerVolumeSpecs := learnerVolumes.CreateVolumeForLearner()
// 	learnerVolumeSpecs = extendLearnerVolumes(learnerVolumeSpecs, logr)
// 	learnerDefn := learnerDefinition{
// 		// secrets:                         secretsForDeployingLearner(req, mountTrainingDataStoreInLearner, mountResultsStoreInLearner),
// 		volumes:                         learnerVolumeSpecs,
// 		volumeMounts:                    learnerVolumes.CreateVolumeMountsForLearner(),
// 		envVars:                         envvarsForLearner,
// 		numberOfLearners:                numLearners,
// 		mountTrainingDataStoreInLearner: mountTrainingDataStoreInLearner,
// 		mountResultsStoreInLearner:      mountResultsStoreInLearner,
// 		name:                            learnerName,
// 	}
// 	logr.Debugf("NewLCMTraining learnerDefn: %+v", envvarsForLearner)

// 	logr.Infof("starting deploying learner infra for split learning")
// 	return splitTraining{&training{ctx, k8sClient, req, req.TrainingId, learnerDefn, helperDefn, logr}}
// }

func volumesForLearner(req *service.JobDeploymentRequest, learnerEnvVars []v1core.EnvVar, mountTrainingDataStoreInLearner, mountResultsStoreInLearner, mountWorkStoreInLearner bool, logr *logger.LocLoggingEntry) learner.Volumes {

	volumesStruct := learner.Volumes{}

	if mountTrainingDataStoreInLearner {
		logr.Debugf("debug_volumesForLearner mountTrainingDataStoreInLearner: %v", mountTrainingDataStoreInLearner)
		dataStoreType := req.EnvVars["DATA_STORE_TYPE"]
		logr.Debugf("dataStoreType: %s", dataStoreType)
		if dataStoreType == learner.DataStoreTypeMountCOSS3 {
			logr.Debugf("debug_volumesForLearner DataStoreTypeMountCOSS3_dataStoreType: %v", dataStoreType)
			region := req.EnvVars["DATA_STORE_REGION"]
			if region == "" {
				region = "us-standard"
			}
			configValStr := config.GetString("MOUNTCOS_GB_CACHE_PER_GPU")
			cacheSize, err := strconv.Atoi(configValStr)
			if err != nil {
				cacheSize = 6
				logr.Warnf("DLAAS_MOUNTCOS_GB_CACHE_PER_GPU value %s is not an integer.  Defaulting to %dGB/GPU", configValStr, cacheSize)
			}
			cacheSize = cacheSize * int(req.Resources.Gpus)
			// reserve 1/3 of cache for prefetching, up to a limit (diskFree is specified in MB, cache in GB)
			diskFree := (cacheSize * 1024) / 3
			if diskFree > 10000 {
				diskFree = 10000
			}

			volumesStruct.TrainingData = &learner.COSVolume{
				VolumeType: dataStoreType,
				ID:         "cosinputmount-" + req.Name,
				Region:     region,
				Bucket:     req.EnvVars["DATA_STORE_OBJECTID"],
				Endpoint:   req.EnvVars["DATA_STORE_AUTHURL"],
				SecretRef:  "cossecretdata-" + req.Name,
				MountSpec: learner.VolumeMountSpec{
					MountPath: getValue(learnerEnvVars, "DATA_DIR"),
					SubPath:   "",
				},
				CacheSize: strconv.Itoa(cacheSize),
				DiskFree:  strconv.Itoa(diskFree),
			}
		} else if dataStoreType == learner.DataStoreHostMountVolume {
			logr.Debugf("debug_volumesForLearner DataStoreHostMountVolume_dataStoreType: %v", dataStoreType)
			hostPath := req.EnvVars["DATA_STORE_PATH"]
			// The variable the learner will get is the concatenated mount path from envvars.go
			mountPath := getValue(learnerEnvVars, "DATA_DIR")
			// While this is the unadulterated value
			subPath := req.EnvVars["DATA_DIR"]
			fmt.Printf("(data) hostPath=%s\nmountPath=%s\nsubPath=%s\n", hostPath, mountPath, subPath)
			// FIXME MLSS Change:
			// 1. remove subpath in mount spec
			// 2. add subpath in host path
			volumesStruct.TrainingData = &learner.COSVolume{
				VolumeType: dataStoreType,
				ID:         "inputmount-" + req.Name,
				//HostPath: hostPath,
				HostPath: hostPath + "/" + subPath,

				MountSpec: learner.VolumeMountSpec{
					Name:      "inputmount-" + req.Name,
					MountPath: mountPath,
					//SubPath: subPath,
					SubPath: "",
				},
			}
			logr.Debugf("TrainingData volume request: %v+", volumesStruct.TrainingData)
		}
	}
	if mountResultsStoreInLearner {
		resultStoreType := req.EnvVars["RESULT_STORE_TYPE"]
		logr.Debugf("resultStoreType: %s", resultStoreType)
		if resultStoreType == learner.DataStoreTypeMountCOSS3 {
			logr.Debugf("debug_volumesForLearner DataStoreTypeMountCOSS3_resultStoreType: %v", resultStoreType)
			region := req.EnvVars["RESULT_STORE_REGION"]
			if region == "" {
				region = "us-standard"
			}
			resultBucketDir := getValue(learnerEnvVars, "RESULT_BUCKET_DIR")
			// drop the "/mnt/results" part of the path and only keep the bucket name
			_, resultBucketName := path.Split(resultBucketDir)
			volumesStruct.ResultsDir = &learner.COSVolume{
				VolumeType: resultStoreType,
				ID:         "cosoutputmount-" + req.Name,
				Region:     region,
				Bucket:     resultBucketName,
				Endpoint:   req.EnvVars["RESULT_STORE_AUTHURL"],
				SecretRef:  "cossecretresults-" + req.Name,
				MountSpec: learner.VolumeMountSpec{
					MountPath: resultBucketDir,
					SubPath:   "",
				},
				CacheSize: "0",
				DiskFree:  "2048",
			}
		} else if resultStoreType == learner.DataStoreHostMountVolume {
			logr.Debugf("debug_volumesForLearner DataStoreHostMountVolume_resultStoreType: %v", resultStoreType)
			hostPath := req.EnvVars["RESULT_STORE_PATH"]
			// The variable the learner will get is the concatenated mount path from envvars.go
			mountPath := getValue(learnerEnvVars, "RESULT_DIR")
			// While this is the unadulterated value
			subPath := req.EnvVars["RESULT_DIR"]

			fmt.Printf("(result) hostPath=%s\nmountPath=%s\nsubPath=%s\n", hostPath, mountPath, subPath)
			volumesStruct.ResultsDir = &learner.COSVolume{
				VolumeType: resultStoreType,
				ID:         "outputmount-" + req.Name,
				// deprecated: MLSS Change: specify sub path in host
				// FIXME MLSS Change:
				// 1. remove training ID in hostpath, create the training ID sub directory in pod initialization
				// 2. remove subpath in mount spec
				//HostPath: hostPath + "/" + subPath + "/" + req.TrainingId,
				HostPath: hostPath + "/" + subPath,
				MountSpec: learner.VolumeMountSpec{
					Name:      "outputmount-" + req.Name,
					MountPath: mountPath,
					//SubPath: subPath,
					SubPath: "",
				},
			}
			logr.Debugf("ResultsDir volume request: %v+", volumesStruct.ResultsDir)
		}
	}
	// FIXME MLSS Change: v_1.5.1 added mount the storage shared directory
	//if mountWorkStoreInLearner {
	//	logr.Debugf("debug_volumesForLearner mountWorkStoreInLearner: %v", mountWorkStoreInLearner)
	//	workStoreType := req.EnvVars["WORK_STORE_TYPE"]
	//	logr.Debugf("debug_volumesForLearner workStoreType: %v", workStoreType)
	//	if workStoreType == learner.DataStoreHostMountVolume {
	//		logr.Debugf("debug_volumesForLearner DataStoreHostMountVolume_workStoreType: %v", workStoreType)
	//		hostPath := req.EnvVars["WORK_STORE_PATH"]
	//		// The variable the learner will get is the concatenated mount path from envvars.go
	//		mountPath := getValue(learnerEnvVars, "MODEL_DIR")
	//		// get WORK_DIR as the container directory
	//		subPath := req.EnvVars["WORK_DIR"]
	//		logr.Debugf("debug_volumesForLearner hostPath: %v, mountPath: %v", hostPath, mountPath)
	//		fmt.Printf("(data) hostPath=%s\nmountPath=%s\nsubPath=%s\n", hostPath, mountPath, subPath)
	//		// 1. remove subpath in mount spec
	//		// 2. add subpath in host path
	//		volumesStruct.WorkDir = &learner.COSVolume{
	//			VolumeType: workStoreType,
	//			ID:         "work-" + req.Name,
	//			//HostPath: hostPath,
	//			HostPath: hostPath + "/" + subPath,
	//
	//			MountSpec: learner.VolumeMountSpec{
	//				Name:      "work-" + req.Name,
	//				MountPath: mountPath,
	//				//SubPath: subPath,
	//				SubPath: "",
	//			},
	//		}
	//		logr.Debugf("WorkDir volume request: %v+", volumesStruct.WorkDir)
	//	}
	//}

	// FIXME MLSS Change: v_1.6.0 added HDFS config
	logr.Debugf("debug_volumesForLearner jobNamespace: %v", req.JobNamespace)
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

	logr.Debugf("debug_volumesForLearner: %v", volumesStruct)
	return volumesStruct
}

//list of shared env vars shared by all containers in helper and learner pod
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

//needs to happen before the volume creation since we muck around with the value and change the paths of data/result dir
func envVarsForDeployingLearner(existingEnvVars []v1core.EnvVar, trainingID string, numLearners int, statefulsetName string, mountTrainingDataStoreInLearner, mountResultsStoreInLearner bool) []v1core.EnvVar {
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

// func NewMLPipeTraining(ctx context.Context, k8sClient kubernetes.Interface, req *service.JobDeploymentRequest, log *logger.LocLoggingEntry) Training {
// 	const cosMountDriverName = "ibm/ibmc-s3fs"
// 	const cosMountType = "mount_cos"
// 	learnerName := fmt.Sprintf("learner-%s", req.Name)
// 	helperName := fmt.Sprintf("lhelper-%s", req.Name)
// 	numLearners := int(req.GetResources().Learners)
// 	if numLearners < 1 {
// 		numLearners = 1
// 	}

// 	dataStoreType := req.EnvVars["DATA_STORE_TYPE"]
// 	resultStoreStype := req.EnvVars["RESULT_STORE_TYPE"]
// 	// FIXME MLSS Change: v_1.5.1 added_de
// 	workStoreType := req.EnvVars["WORK_STORE_TYPE"]

// 	mountTrainingDataStoreInLearner := !(dataStoreType == learner.DataStoreTypeS3)
// 	mountResultsStoreInLearner := !(resultStoreStype == learner.DataStoreTypeS3)

// 	mountWorkStoreInLearner := !(workStoreType == learner.DataStoreTypeS3)

// 	logr := log.WithFields(logrus.Fields{
// 		"learner_name": learnerName,
// 		"helper_name":  helperName,
// 		"mounted_cos":  mountResultsStoreInLearner && mountTrainingDataStoreInLearner,
// 	})
// 	logr.Debugf("debug_for_param ParameterServer lcm newTraining pss: %v, psCpu: %v, psImage: %v, psMemory: %v", req.PSs, req.PSCPU, req.PSImage, req.PSMemory)
// 	logr.Debugf("extractEnvVarsFromDeploymentRequest(req): %+v", req.EnvVars)

// 	envVarsFromDeploymentRequest := extractEnvVarsFromDeploymentRequest(req) //shared across all containers of training
// 	logr.Debugf("envVarsFromDeploymentRequest: %+v", envVarsFromDeploymentRequest)

// 	envvarsForLearner := envVarsForDeployingLearner(envVarsFromDeploymentRequest, req.TrainingId,
// 		numLearners, learnerName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner) //only for learner
// 	logr.Debugf("NewLCMTraining envvarsForLearner: %+v", envvarsForLearner)

// 	learnerVolumes := volumesForLearner(req, envvarsForLearner, mountTrainingDataStoreInLearner, mountResultsStoreInLearner, mountWorkStoreInLearner, logr)
// 	learnerVolumeSpecs := learnerVolumes.CreateVolumeForLearner()
// 	learnerVolumeSpecs = extendLearnerVolumes(learnerVolumeSpecs, logr)
// 	learnerDefn := learnerDefinition{
// 		secrets:                         secretsForDeployingLearner(req, mountTrainingDataStoreInLearner, mountResultsStoreInLearner),
// 		volumes:                         learnerVolumeSpecs,
// 		volumeMounts:                    learnerVolumes.CreateVolumeMountsForLearner(),
// 		envVars:                         envvarsForLearner,
// 		numberOfLearners:                numLearners,
// 		mountTrainingDataStoreInLearner: mountTrainingDataStoreInLearner,
// 		mountResultsStoreInLearner:      mountResultsStoreInLearner,
// 		name:                            learnerName,
// 	}
// 	logr.Debugf("MLPipelineTraining learnerDefn: %+v", envvarsForLearner)

// 	logr.Infof("starting deploying learner infra for split learning")
// 	return splitTraining{&training{ctx, k8sClient, req, req.TrainingId, learnerDefn, helperDefn, logr}}
// }

// func NewSigleJobTraining(ctx context.Context, k8sClient kubernetes.Interface, req *service.JobDeploymentRequest, log *logger.LocLoggingEntry) Training {
// 	const cosMountDriverName = "ibm/ibmc-s3fs"
// 	const cosMountType = "mount_cos"
// 	learnerName := fmt.Sprintf("learner-%s", req.Name)
// 	helperName := fmt.Sprintf("lhelper-%s", req.Name)
// 	numLearners := int(req.GetResources().Learners)
// 	if numLearners < 1 {
// 		numLearners = 1
// 	}

// 	dataStoreType := req.EnvVars["DATA_STORE_TYPE"]
// 	resultStoreStype := req.EnvVars["RESULT_STORE_TYPE"]
// 	// FIXME MLSS Change: v_1.5.1 added_de
// 	workStoreType := req.EnvVars["WORK_STORE_TYPE"]

// 	mountTrainingDataStoreInLearner := !(dataStoreType == learner.DataStoreTypeS3)
// 	mountResultsStoreInLearner := !(resultStoreStype == learner.DataStoreTypeS3)

// 	mountWorkStoreInLearner := !(workStoreType == learner.DataStoreTypeS3)

// 	logr := log.WithFields(logrus.Fields{
// 		"learner_name": learnerName,
// 		"helper_name":  helperName,
// 		"mounted_cos":  mountResultsStoreInLearner && mountTrainingDataStoreInLearner,
// 	})
// 	logr.Debugf("debug_for_param ParameterServer lcm newTraining pss: %v, psCpu: %v, psImage: %v, psMemory: %v", req.PSs, req.PSCPU, req.PSImage, req.PSMemory)
// 	logr.Debugf("extractEnvVarsFromDeploymentRequest(req): %+v", req.EnvVars)

// 	envVarsFromDeploymentRequest := extractEnvVarsFromDeploymentRequest(req) //shared across all containers of training
// 	logr.Debugf("envVarsFromDeploymentRequest: %+v", envVarsFromDeploymentRequest)

// 	envvarsForLearner := envVarsForDeployingLearner(envVarsFromDeploymentRequest, req.TrainingId,
// 		numLearners, learnerName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner) //only for learner
// 	logr.Debugf("NewLCMTraining envvarsForLearner: %+v", envvarsForLearner)

// 	learnerVolumes := volumesForLearner(req, envvarsForLearner, mountTrainingDataStoreInLearner, mountResultsStoreInLearner, mountWorkStoreInLearner, logr)
// 	learnerVolumeSpecs := learnerVolumes.CreateVolumeForLearner()
// 	learnerVolumeSpecs = extendLearnerVolumes(learnerVolumeSpecs, logr)
// 	learnerDefn := learnerDefinition{
// 		secrets:                         secretsForDeployingLearner(req, mountTrainingDataStoreInLearner, mountResultsStoreInLearner),
// 		volumes:                         learnerVolumeSpecs,
// 		volumeMounts:                    learnerVolumes.CreateVolumeMountsForLearner(),
// 		envVars:                         envvarsForLearner,
// 		numberOfLearners:                numLearners,
// 		mountTrainingDataStoreInLearner: mountTrainingDataStoreInLearner,
// 		mountResultsStoreInLearner:      mountResultsStoreInLearner,
// 		name:                            learnerName,
// 	}
// 	logr.Debugf("MLPipelineTraining learnerDefn: %+v", envvarsForLearner)

// 	logr.Infof("starting deploying learner infra for split learning")
// 	return splitTraining{&training{ctx, k8sClient, req, req.TrainingId, learnerDefn, helperDefn, logr}}
// }

// func NewOperatorTraining(ctx context.Context, k8sClient kubernetes.Interface, req *service.JobDeploymentRequest, log *logger.LocLoggingEntry) Training {
// 	const cosMountDriverName = "ibm/ibmc-s3fs"
// 	const cosMountType = "mount_cos"
// 	learnerName := fmt.Sprintf("learner-%s", req.Name)
// 	helperName := fmt.Sprintf("lhelper-%s", req.Name)
// 	numLearners := int(req.GetResources().Learners)
// 	if numLearners < 1 {
// 		numLearners = 1
// 	}

// 	dataStoreType := req.EnvVars["DATA_STORE_TYPE"]
// 	resultStoreStype := req.EnvVars["RESULT_STORE_TYPE"]
// 	// FIXME MLSS Change: v_1.5.1 added_de
// 	workStoreType := req.EnvVars["WORK_STORE_TYPE"]

// 	mountTrainingDataStoreInLearner := !(dataStoreType == learner.DataStoreTypeS3)
// 	mountResultsStoreInLearner := !(resultStoreStype == learner.DataStoreTypeS3)

// 	mountWorkStoreInLearner := !(workStoreType == learner.DataStoreTypeS3)

// 	logr := log.WithFields(logrus.Fields{
// 		"learner_name": learnerName,
// 		"helper_name":  helperName,
// 		"mounted_cos":  mountResultsStoreInLearner && mountTrainingDataStoreInLearner,
// 	})
// 	logr.Debugf("debug_for_param ParameterServer lcm newTraining pss: %v, psCpu: %v, psImage: %v, psMemory: %v", req.PSs, req.PSCPU, req.PSImage, req.PSMemory)
// 	logr.Debugf("extractEnvVarsFromDeploymentRequest(req): %+v", req.EnvVars)

// 	envVarsFromDeploymentRequest := extractEnvVarsFromDeploymentRequest(req) //shared across all containers of training
// 	logr.Debugf("envVarsFromDeploymentRequest: %+v", envVarsFromDeploymentRequest)

// 	envvarsForLearner := envVarsForDeployingLearner(envVarsFromDeploymentRequest, req.TrainingId,
// 		numLearners, learnerName, mountTrainingDataStoreInLearner, mountResultsStoreInLearner) //only for learner
// 	logr.Debugf("NewLCMTraining envvarsForLearner: %+v", envvarsForLearner)

// 	learnerVolumes := volumesForLearner(req, envvarsForLearner, mountTrainingDataStoreInLearner, mountResultsStoreInLearner, mountWorkStoreInLearner, logr)
// 	learnerVolumeSpecs := learnerVolumes.CreateVolumeForLearner()
// 	learnerVolumeSpecs = extendLearnerVolumes(learnerVolumeSpecs, logr)
// 	learnerDefn := learnerDefinition{
// 		secrets:                         secretsForDeployingLearner(req, mountTrainingDataStoreInLearner, mountResultsStoreInLearner),
// 		volumes:                         learnerVolumeSpecs,
// 		volumeMounts:                    learnerVolumes.CreateVolumeMountsForLearner(),
// 		envVars:                         envvarsForLearner,
// 		numberOfLearners:                numLearners,
// 		mountTrainingDataStoreInLearner: mountTrainingDataStoreInLearner,
// 		mountResultsStoreInLearner:      mountResultsStoreInLearner,
// 		name:                            learnerName,
// 	}
// 	logr.Debugf("MLPipelineTraining learnerDefn: %+v", envvarsForLearner)

// 	logr.Infof("starting deploying learner infra for split learning")
// 	return splitTraining{&training{ctx, k8sClient, req, req.TrainingId, learnerDefn, helperDefn, logr}}
// }
