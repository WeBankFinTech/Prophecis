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
package models

import (
	"github.com/kubernetes-client/go/kubernetes/client"
	"github.com/spf13/viper"
	"strconv"
	stringsUtil "strings"
	"time"
	"webank/AIDE/notebook-server/pkg/commons/config"
	"webank/AIDE/notebook-server/pkg/commons/constants"
	"webank/AIDE/notebook-server/pkg/commons/logger"
)

type NotebookFromK8s struct {
	ApiVersion string            `json:"apiVersion,omitempty"`
	Items      []K8sNotebook     `json:"items,omitempty"`
	Kind       string            `json:"kind,omitempty"`
	Metadata   client.V1ListMeta `json:"metadata,omitempty"`
}

var PlatformNamespace = viper.GetString(config.PlatformNamespace)
var userIdLabel = PlatformNamespace + "-UserId"

func ParseToNotebookRes(m NotebookFromK8s, configs map[string]interface{}) []*Notebook {

	var notebookResList []*Notebook
	yarnConfig := map[string]string{}

	notebooks := m.Items
	if notebooks == nil || len(notebooks) <= 0 {
		return notebookResList
	}

	for _, nb := range notebooks {

		status := "NotReady"
		logger.Logger().Debugf("ParseToNotebookRes for state: %v", nb.Status.ContainerState.Running)
		if nil != nb.Status.ContainerState.Running {
			status = "Ready"
		}else if nil != nb.Status.ContainerState.Waiting {//error
			logger.Logger().Debugf("ParseToNotebookRes for state Waiting: %v", nb.Status.ContainerState.Waiting.Message)
			status = "Waiting"
		}else if nil != nb.Status.ContainerState.Terminated {
			logger.Logger().Debugf("ParseToNotebookRes for state terminater: %v", nb.Status.ContainerState.Terminated.Message)
			status = "Terminated"
		}else{
			logger.Logger().Debugf("ParseToNotebookRes for state terminater: %v", nb.Status.ContainerState)
			status = "Waiting"
		}

		var timeLayoutStr = "2006-01-02 15:04:05"

		//get creating time in label first,if it exists
		var updateTime string
		createTimeFromLabel := nb.Metadata.Labels[constants.CreateTime]
		if createTimeFromLabel != "" {
			timeFromLabel, err := strconv.ParseInt(createTimeFromLabel, 10, 64)
			if err != nil {
				logger.Logger().Errorf("ParseInt failed, time: %v", createTimeFromLabel)
			} else {
				format := time.Unix(timeFromLabel, 0).Format(timeLayoutStr)
				updateTime = format
			}
		}
		if updateTime == "" {
			updateTime = nb.Metadata.CreationTimestamp.Format(timeLayoutStr)
		}

		volumeMounts := nb.Spec.Template.Spec.Containers[0].VolumeMounts
		envVars := nb.Spec.Template.Spec.Containers[0].Env
		volumes := nb.Spec.Template.Spec.Volumes

		var dataVolume []*MountInfo
		var workspaceVolume MountInfo

		if nil != volumeMounts && nil != volumes {
			for _, v1VolumeMount := range volumeMounts {
				var mountInfo MountInfo
				var size int64 = 0
				mountType := "New"
				accessMode := "ReadWriteMany"
				mountInfo.MountType = &mountType
				mountInfo.AccessMode = &accessMode
				mountInfo.Size = &size

				if stringsUtil.Contains(v1VolumeMount.Name, "datadir") {
					mountPath := "/" + v1VolumeMount.MountPath[5:]
					mountInfo.MountPath = &mountPath
					mountInfo = GetMountInfo(v1VolumeMount, volumes, mountInfo, envVars)
					dataVolume = append(dataVolume, &mountInfo)
				} else if v1VolumeMount.Name == "workdir" {
					mountPath := v1VolumeMount.MountPath
					mountInfo.MountPath = &mountPath
					mountInfo = GetMountInfo(v1VolumeMount, volumes, mountInfo, envVars)
					workspaceVolume = mountInfo
				}
			}
		}

		config := configs[nb.Metadata.Namespace+"-"+nb.Metadata.Name+"-yarn-resource-setting"]
		//configs.
		if config != nil {
			sessionConfig := config.(map[string]interface{})["session_configs"]
			logger.Logger().Infof("config json 1 ")
			yarnConfig["queue"] = sessionConfig.(map[string]interface{})["queue"].(string)
			logger.Logger().Infof("config json 2 ")
			yarnConfig["executors"] = sessionConfig.(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.instances"].(string)
			logger.Logger().Infof("config json 3 ")
			yarnConfig["executorCores"] = sessionConfig.(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.cores"].(string)
			yarnConfig["executorMemory"] = sessionConfig.(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.memory"].(string)
			yarnConfig["driverMemory"] = sessionConfig.(map[string]interface{})["driverMemory"].(string)
		} else {
			yarnConfig["queue"] = ""
			yarnConfig["executors"] = ""
			yarnConfig["executorCores"] = ""
			yarnConfig["executorMemory"] = ""
			yarnConfig["driverMemory"] = ""
		}

		newNotebookRes := Notebook{
			CPU:             nb.Spec.Template.Spec.Containers[0].Resources.Requests["cpu"],
			Image:           nb.Spec.Template.Spec.Containers[0].Image,
			Gpu:             nb.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"],
			Memory:          nb.Spec.Template.Spec.Containers[0].Resources.Requests["memory"],
			Queue:           yarnConfig["queue"],
			Executors:       yarnConfig["executors"],
			ExecutorCores:   yarnConfig["executorCores"],
			ExecutorMemory:  yarnConfig["executorMemory"],
			DriverMemory:    yarnConfig["driverMemory"],
			Name:            nb.Metadata.Name,
			Namespace:       nb.Metadata.Namespace,
			Pods:            "",
			Service:         "",
			SrtImage:        "",
			Status:          status,
			Uptime:          updateTime,
			User:            nb.Metadata.Labels[userIdLabel],
			Volumns:         "",
			DataVolume:      dataVolume,
			WorkspaceVolume: &workspaceVolume,
		}

		notebookResList = append(notebookResList, &newNotebookRes)
	}

	return notebookResList
}

func GetMountInfo(v1VolumeMount client.V1VolumeMount, volumes []client.V1Volume, mountInfo MountInfo, vars []client.V1EnvVar) MountInfo {
	for _, v1Volume := range volumes {
		if v1VolumeMount.Name == v1Volume.Name {
			if nil != vars {
				for _, env := range vars {
					if env.Name == v1Volume.Name {
						mountInfo.SubPath = env.Value
						hostPath := v1Volume.HostPath.Path
						if env.Value == "" || env.Value == "/" {
							hostPath = hostPath
						} else if stringsUtil.HasPrefix(env.Value, "/") {
							hostPath = hostPath[0 : len(hostPath)-len(env.Value)]
						} else {
							hostPath = hostPath[0 : len(hostPath)-len(env.Value)-1]
						}
						mountInfo.LocalPath = &hostPath
					}
				}
			}
		}
	}
	return mountInfo
}

func GetQueue(envs []client.V1EnvVar) string {
	for _, env := range envs {
		if env.Name == "QUEUE" {
			return env.Value
		}
	}
	return ""
}

func GetEnv(envs []client.V1EnvVar) map[string]string {
	// Init Env Map
	var envMap = make(map[string]string)
	envMap["QUEUE"] = ""
	envMap["EXECUTOR_CORES"] = ""
	envMap["EXECUTOR_MEMORY"] = ""
	envMap["DRIVER_MEMORY"] = ""
	envMap["EXECUTORS"] = ""

	for _, env := range envs {
		if env.Name == "QUEUE" {
			envMap["QUEUE"] = env.Value
		} else if env.Name == "EXECUTOR_CORES" {
			envMap["EXECUTOR_CORES"] = env.Value
		} else if env.Name == "EXECUTOR_MEMORY" {
			envMap["EXECUTOR_MEMORY"] = stringsUtil.Trim(env.Value, "g")
		} else if env.Name == "DRIVER_MEMORY" {
			envMap["DRIVER_MEMORY"] = stringsUtil.Trim(env.Value, "g")
		} else if env.Name == "EXECUTORS" {
			envMap["EXECUTORS"] = env.Value
		}
	}
	return envMap
}
