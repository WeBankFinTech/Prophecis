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
package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubernetes-client/go/kubernetes/client"
	"github.com/kubernetes-client/go/kubernetes/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	stringsUtil "strings"
	"sync"
	"time"
	cfg "webank/AIDE/notebook-server/pkg/commons/config"
	"webank/AIDE/notebook-server/pkg/commons/constants"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"
)

type NotebookControllerClient struct {
	k8sClient *client.APIClient
}

var ncClientInitOnce sync.Once
var ncClient *NotebookControllerClient

// TODO move these configs to a configmap
const (
	MistakenTokenHeader = "Authentication"
	CorrectTokenHeader  = "Authorization"
	TokenValuePrefix    = "Bearer "
	SATokenPath         = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	NBApiGroup          = "kubeflow.org"
	NBApiVersion        = "v1alpha1"
	NBApiPlural         = "notebooks"
)

var PlatformNamespace = viper.GetString(cfg.PlatformNamespace)

func GetNBCClient() *NotebookControllerClient {
	ncClientInitOnce.Do(func() {
		cfg, err := config.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		tokenValue := cfg.DefaultHeader[MistakenTokenHeader]
		if len(tokenValue) == 0 {
			tv, err := ioutil.ReadFile(SATokenPath)
			if err != nil {
				panic(err.Error())
			}
			tokenValue = TokenValuePrefix + string(tv)
		}
		defaultHeader := map[string]string{CorrectTokenHeader: tokenValue}
		cfg.DefaultHeader = defaultHeader
		// creates the clientset
		ncClient = &NotebookControllerClient{}
		ncClient.k8sClient = client.NewAPIClient(cfg)
	})
	return ncClient
}

// create a notebook in a namespace
func (*NotebookControllerClient) CreateNotebook(notebook *models.NewNotebookRequest, userId string, gid string, uid string, token string) error {
	/*** 0. get a nb template for k8s request && get cluster message from namespaces ***/
	nbForK8s := createK8sNBTemplate()

	nbContainer := nbForK8s.Spec.Template.Spec.Containers[0]

	/*** 1. volumes ***/
	var nbVolumes []client.V1Volume
	var localPath string
	if nil != notebook.WorkspaceVolume {
		localPath = *notebook.WorkspaceVolume.LocalPath
		if userId != stringsUtil.Split(localPath[1:], "/")[3] {
			gid = viper.GetString(cfg.MLSSGroupId)
		}
		subPath := notebook.WorkspaceVolume.SubPath
		nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
			Name:  "workdir",
			Value: subPath,
		})
		if subPath != "" && subPath != "/" {
			if !stringsUtil.HasPrefix(subPath, "/") && !stringsUtil.HasSuffix(localPath, "/") {
				localPath = localPath + "/" + subPath
			} else if stringsUtil.HasPrefix(subPath, "/") && stringsUtil.HasSuffix(localPath, "/") {
				localPath = localPath + subPath[1:]
			} else {
				localPath = localPath + subPath
			}
			logger.Logger().Debugf("k8s createNotebook workPath: %v", localPath+subPath)
		}
		nbVolumes = append(nbVolumes,
			client.V1Volume{
				Name: "workdir",
				HostPath: &client.V1HostPathVolumeSource{
					Path: localPath,
				},
			})

		// persist jupyter config
		//nbVolumes = append(nbVolumes,
		//	client.V1Volume{
		//		Name: "notebook-config",
		//		HostPath: &client.V1HostPathVolumeSource{
		//			Path: localPath + "/notebook-config/.jupyter",
		//		},
		//	})
		//
		//nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
		//	Name:      "notebook-config",
		//	MountPath: "/.jupyter",
		//})

	}
	if nil != notebook.DataVolume {
		//for i, dv := range notebook.DataVolume {
		localPath := *notebook.DataVolume.LocalPath
		if userId != stringsUtil.Split(localPath[1:], "/")[3] {
			gid = viper.GetString(cfg.MLSSGroupId)
		}
		subPath := notebook.DataVolume.SubPath
		nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
			Name:  fmt.Sprint("datadir"),
			Value: subPath,
		})
		if subPath != "" && subPath != "/" {
			if !stringsUtil.HasPrefix(subPath, "/") && !stringsUtil.HasSuffix(localPath, "/") {
				localPath = localPath + "/" + subPath
			} else if stringsUtil.HasPrefix(subPath, "/") && stringsUtil.HasSuffix(localPath, "/") {
				localPath = localPath + subPath[1:]
			} else {
				localPath = localPath + subPath
			}
			logger.Logger().Debugf("k8s createNotebook workPath: %v", localPath+subPath)
		}
		dvForNb := client.V1Volume{
			Name: fmt.Sprint("datadir"),
			HostPath: &client.V1HostPathVolumeSource{
				Path: localPath,
			},
		}
		nbVolumes = append(nbVolumes, dvForNb)
		//}
	}

	if viper.GetString(cfg.HadoopEnable) == "true" {
		nbVolumes = append(nbVolumes,
			client.V1Volume{
				Name: "bdap-appcom-config",
				HostPath: &client.V1HostPathVolumeSource{
					Path: viper.GetString(cfg.ConfigPath),
				},
			},
			client.V1Volume{
				Name: "bdap-appcom-commonlib",
				HostPath: &client.V1HostPathVolumeSource{
					Path: viper.GetString(cfg.CommonlibPath),
				},
			},
			client.V1Volume{
				Name: "bdap-appcom-install",
				HostPath: &client.V1HostPathVolumeSource{
					Path: viper.GetString(cfg.InstallPath),
				},
			},
			client.V1Volume{
				Name: "jdk",
				HostPath: &client.V1HostPathVolumeSource{
					Path: viper.GetString(cfg.JavaPath),
				},
			})

		// mount HADOOP Install File
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
			Name:      "bdap-appcom-config",
			MountPath: viper.GetString(cfg.ConfigPath),
		})
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
			Name:      "bdap-appcom-install",
			MountPath: viper.GetString(cfg.InstallPath),
		})
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
			Name:      "bdap-appcom-commonlib",
			MountPath: viper.GetString(cfg.CommonlibPath),
		})
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
			Name:      "jdk",
			MountPath: viper.GetString(cfg.JavaPath),
		})
	}

	nbVolumes = append(nbVolumes,
		client.V1Volume{
			Name: "config-json",
			ConfigMap: &client.V1ConfigMapVolumeSource{
				DefaultMode: 0755,
				Name:        *notebook.Name + "-" + "yarn-resource-setting",
			},
		}, client.V1Volume{
			Name: "linkismagic-json",
			ConfigMap: &client.V1ConfigMapVolumeSource{
				DefaultMode: 0755,
				Name:        *notebook.Name + "-" + "yarn-resource-setting",
			},
		})

	//}
	nbForK8s.Spec.Template.Spec.Volumes = append(nbForK8s.Spec.Template.Spec.Volumes, nbVolumes...)

	/*** 3. container ***/
	startWorkspace := ""
	if nil == notebook.DataVolume {
		startWorkspace = "/workspace"
	}
	nbContainer.Command = []string{"/etc/config/start-sh"}
	nbContainer.Args = []string{"jupyter lab --notebook-dir=/home/" + userId + startWorkspace +
		" --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' " +
		"--NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"}
	var extraResourceMap = make(map[string]string)
	if notebook.ExtraResources != "" {
		err := json.Unmarshal([]byte(notebook.ExtraResources), &extraResourceMap)
		if nil != err {
			logger.Logger().Errorf("Error parsing extra resources: %s", notebook.ExtraResources)
			logger.Logger().Errorf(err.Error())
			return err
		}
	}
	var limitMap = make(map[string]string)
	limitMap["cpu"] = fmt.Sprint(*(notebook.CPU))
	MemoryAmount := notebook.Memory.MemoryAmount
	float := fmt.Sprintf("%.1f", *MemoryAmount)
	s := notebook.Memory.MemoryUnit
	limitMap["memory"] = fmt.Sprint(float, s)
	if extraResourceMap["nvidia.com/gpu"] != "" && extraResourceMap["nvidia.com/gpu"] != "0" {
		limitMap["nvidia.com/gpu"] = extraResourceMap["nvidia.com/gpu"]
	} else {
		nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
			Name:  "NVIDIA_VISIBLE_DEVICES",
			Value: "",
		})
	}

	logger.Logger().Infof("CreateNotebook limitMap: %s", limitMap)

	nbContainer.Resources = &client.V1ResourceRequirements{
		Requests: limitMap,
		Limits:   limitMap,
	}
	nbContainer.Image = *(notebook.Image.ImageName)

	if nil != notebook.WorkspaceVolume {
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
			Name:      "workdir",
			MountPath: "/workspace",
		})
	}
	if nil != notebook.DataVolume {
		mountPath := *notebook.DataVolume.MountPath
		if !stringsUtil.HasPrefix(mountPath, "/") && !stringsUtil.HasSuffix(mountPath, "/") {
			if stringsUtil.Contains(mountPath, "/") {
				return errors.New("mountPath can only have level 1 directories")
			}
		} else if stringsUtil.HasPrefix(mountPath, "/") && !stringsUtil.HasSuffix(mountPath, "/") {
			if stringsUtil.Contains(mountPath[1:], "/") {
				return errors.New("mountPath can only have level 1 directories")
			}
		} else if !stringsUtil.HasPrefix(mountPath, "/") && stringsUtil.HasSuffix(mountPath, "/") {
			if stringsUtil.Index(mountPath, "/") < len(mountPath)-1 {
				return errors.New("mountPath can only have level 1 directories")
			}
		} else if stringsUtil.HasPrefix(mountPath, "/") && stringsUtil.HasSuffix(mountPath, "/") {
			if stringsUtil.Contains(mountPath[1:(len(mountPath) - 1)], "/") {
				return errors.New("mountPath can only have level 1 directories")
			}
		}
		if stringsUtil.HasPrefix(mountPath, "/") {
			mountPath = "data-" + mountPath[1:]
		} else {
			mountPath = "data-" + mountPath
		}
		dvmForNb := client.V1VolumeMount{
			Name:      fmt.Sprint("datadir"),
			MountPath: mountPath,
		}
		nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, dvmForNb)
	}
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name: "Driver_Host",
		ValueFrom: &client.V1EnvVarSource{
			FieldRef: &client.V1ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	})
	logger.Logger().Debugf("create notebook for gid: %v", gid)
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "NB_GID",
		Value: gid,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "NB_UID",
		Value: uid,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "NB_USER",
		Value: userId,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "GRANT_SUDO",
		Value: "no",
	})

	nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
		Name:      "config-json",
		MountPath: "/etc/sparkmagic/",
	})
	nbContainer.VolumeMounts = append(nbContainer.VolumeMounts, client.V1VolumeMount{
		Name:      "linkismagic-json",
		MountPath: "/etc/linkismagic/",
	})

	//create NodePort Service
	notebookService, nodePortArray, err := createK8sNBServiceTemplate(notebook)
	if nil != err {
		logger.Logger().Errorf(err.Error())
		return err
	}

	nbForK8s.Service = *notebookService
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "Spark_Driver_Port",
		Value: strconv.Itoa(nodePortArray[0]),
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "Spark_Driver_BlockManager_Port",
		Value: strconv.Itoa(nodePortArray[1]),
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "Linkis_Token",
		Value: viper.GetString(cfg.LinkisToken),
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "GUARDIAN_TOKEN",
		Value: token,
	})

	nbForK8s.Spec.Template.Spec.Containers[0] = nbContainer

	/*** 4. fill in other info ***/
	nbForK8s.Metadata.Name = *(notebook.Name)
	nbForK8s.Metadata.Namespace = *(notebook.Namespace)

	userIdLabel := PlatformNamespace + "-UserId"
	workDirLabel := PlatformNamespace + "-WorkdDir"
	//set notebook create time in label
	nowTimeStr := strconv.FormatInt(time.Now().Unix(), 10)

	//nbForK8s.Metadata.Labels = map[string]string{"MLSS-UserId": userId}
	var workDirPath string
	if nil != notebook.DataVolume {
		workDirPath = *notebook.DataVolume.LocalPath
	} else {
		workDirPath = *notebook.WorkspaceVolume.LocalPath
	}
	logger.Logger().Info("localPath is ", localPath)
	nbForK8s.Metadata.Labels = map[string]string{
		userIdLabel:          userId,
		workDirLabel:         stringsUtil.Replace(workDirPath[1:], "/", "_", -1),
		constants.CreateTime: nowTimeStr,
	}

	_, resp, err := ncClient.k8sClient.CustomObjectsApi.CreateNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, nbForK8s.Metadata.Namespace, NBApiPlural, nbForK8s, nil)
	if nil != err {
		logger.Logger().Errorf("Error creating notebook: %+v", err.Error())
		return err
	}
	if resp != nil {
		logger.Logger().Debugf("response from k8s: %+v", resp)
	}
	return err
}

// delete notebook in a namespace
func (*NotebookControllerClient) DeleteNotebook(nb string, ns string) error {
	_, resp, err := ncClient.k8sClient.CustomObjectsApi.DeleteNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, ns, NBApiPlural, nb, client.V1DeleteOptions{}, nil)
	if err != nil {
		logger.Logger().Errorf(""+
			" deleting notebook: %+v", err.Error())
		return err
	}
	if resp == nil {
		logger.Logger().Error("empty response")
	}

	err = deleteK8sNBConfigMap(ns, nb+"-"+"yarn-resource-setting")
	if err != nil {
		logger.Logger().Errorf(""+
			" deleting notebook configmap: %+v", err.Error())
		return err
	}
	return nil
}

// get notebooks in a namespace, or of a user, or of a user in a namespace
func (*NotebookControllerClient) GetNotebooks(ns string, userId string, workDir string) (interface{}, error) {
	//createK8sNBConfigMap("ns-ns-bdap-common-common-dev-csf",logr)
	var optionals = make(map[string]interface{})
	userIdLabelKey := PlatformNamespace + "-UserId="
	workDirLabel := ""
	//optionals["labelSelector"] = fmt.Sprint("MLSS-UserId=", userId)
	if workDir != "" {
		workDir = stringsUtil.Replace(workDir[1:], "/", "_", -1)
		workDirLabel := fmt.Sprint(PlatformNamespace+"-WorkDir=", workDir)
		optionals["labelSelector"] = fmt.Sprint(userIdLabelKey, userId) + "," + workDirLabel
		logger.Logger().Info("labelselector:", optionals["labelSelector"])
	} else {
		optionals["labelSelector"] = fmt.Sprint(userIdLabelKey, userId)
		logger.Logger().Info("workdir is  nil label is %s", optionals["labelSelector"])
	}
	var list, resp interface{}
	var err error
	if len(ns) == 0 && len(userId) == 0 {
		optionals["labelSelector"] = workDirLabel
		list, resp, err = ncClient.k8sClient.CustomObjectsApi.ListClusterCustomObject(context.Background(), NBApiGroup, NBApiVersion, NBApiPlural, optionals)
	} else if len(ns) == 0 && len(userId) != 0 {
		// get user notebooks in all namespaces
		logger.Logger().Debugf("len(ns) == 0 && len(userId)")
		list, resp, err = ncClient.k8sClient.CustomObjectsApi.ListClusterCustomObject(context.Background(), NBApiGroup, NBApiVersion, NBApiPlural, optionals)
	} else if ns == "null" && userId == "null" {
		optionals["labelSelector"] = PlatformNamespace + "-UserId" + workDirLabel
		list, resp, err = ncClient.k8sClient.CustomObjectsApi.ListClusterCustomObject(context.Background(), NBApiGroup, NBApiVersion, NBApiPlural, optionals)
	} else {
		if len(userId) == 0 && len(ns) != 0 {
			optionals = nil
		}
		list, resp, err = ncClient.k8sClient.CustomObjectsApi.ListNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, ns, NBApiPlural, optionals)
	}

	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}
	if resp == nil {
		logger.Logger().Error("emtpy response")
		return nil, nil
	}
	return list, nil
}

// Get namespaces
func (*NotebookControllerClient) GetNamespaces() []string {
	list, response, err := ncClient.k8sClient.CoreV1Api.ListNamespace(context.Background(), nil)
	if err != nil {
		logrus.Error(err.Error())
		return nil
	}
	if response == nil {
		logrus.Error("empty response")
		return nil
	}
	var namespaceList []string
	for _, namespace := range list.Items {
		namespaceList = append(namespaceList, namespace.Metadata.Name)
	}
	return namespaceList
}

func (*NotebookControllerClient) listNamespacedResourceQuota(ns string) ([]client.V1ResourceQuota, error) {
	list, response, err := ncClient.k8sClient.CoreV1Api.ListNamespacedResourceQuota(context.Background(), ns, nil)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if response == nil {
		logrus.Error("empty response")
		return nil, err
	}
	logrus.Infof("listNamespacedResourceQuota,ns: %v, list: %v", ns, list)
	return list.Items, nil
}

func (*NotebookControllerClient) listNamespaceNodes(ns string) ([]client.V1Node, error) {
	v1Namespace, _, readErr := ncClient.k8sClient.CoreV1Api.ReadNamespace(context.Background(), ns, nil)
	if readErr != nil {
		logrus.Error(readErr.Error())
		return nil, readErr
	}

	namespaceAnnotations := v1Namespace.Metadata.Annotations
	logrus.Debugf("listNamespaceNodes for namespaceAnnotations: %v", namespaceAnnotations)

	selector := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]

	var localVarOptionals = make(map[string]interface{})
	localVarOptionals["labelSelector"] = selector

	list, response, err := ncClient.k8sClient.CoreV1Api.ListNode(context.Background(), localVarOptionals)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if response == nil {
		logrus.Error("empty response")
		return nil, err
	}
	logrus.Debugf("listNamespaceNodes by namespace: %v, list.items.length: %v", ns, len(list.Items))
	return list.Items, nil
}

func createK8sNBTemplate() *models.K8sNotebook {
	nb := &models.K8sNotebook{
		ApiVersion: NBApiGroup + "/" + NBApiVersion,
		Kind:       "Notebook",
		Metadata: client.V1ObjectMeta{
			Name:      "",
			Namespace: "",
			//Labels: nil,
		},
		Spec: models.NBSpec{
			Template: models.NBTemplate{
				Spec: models.NBDetailedSpec{
					ServiceAccountName: "jupyter-notebook",
					Containers: []client.V1Container{
						{
							Name: "notebook-container",
							VolumeMounts: []client.V1VolumeMount{
								{
									Name:      "dshm",
									MountPath: "/dev/shm",
								},
								{
									Name:      "start-sh",
									MountPath: "/etc/config",
								}, {
									Name:      "timezone-volume",
									MountPath: "/etc/localtime",
								}},
							Env:     nil,
							Command: []string{"/etc/config/start-sh"},
							SecurityContext: &client.V1SecurityContext{
								RunAsGroup:               int64(0),
								RunAsUser:                int64(0),
								AllowPrivilegeEscalation: true,
								Privileged:               false,
							},
						},
					},
					TtlSecondsAfterFinished: 300,
					Volumes: []client.V1Volume{
						{
							Name: "dshm",
							EmptyDir: &client.V1EmptyDirVolumeSource{
								Medium: "Memory",
							}},
						{
							Name: "start-sh",
							ConfigMap: &client.V1ConfigMapVolumeSource{
								DefaultMode: 0777,
								Name:        "notebook-entrypoint-files",
							},
						},

						{
							Name: "timezone-volume",
							HostPath: &client.V1HostPathVolumeSource{
								Path: "/usr/share/zoneinfo/Asia/Shanghai",
							}},
					},
					// FIXME not working
					SecurityContext: client.V1SecurityContext{
						RunAsGroup:               int64(0),
						RunAsUser:                int64(0),
						AllowPrivilegeEscalation: true,
						Privileged:               false,
					},
				},
			},
		},
	}
	return nb
}

func createK8sNBServiceTemplate(notebook *models.NewNotebookRequest) (*client.V1Service, [2]int, error) {

	ServiceName := *(notebook.Name) + "-service"
	nbService := &client.V1Service{
		Metadata: &client.V1ObjectMeta{
			Name:      ServiceName,
			Namespace: *(notebook.Namespace),
		},
		Spec: &client.V1ServiceSpec{},
	}
	nbServiceSelector := make(map[string]string)
	nbServiceSelector["notebook-name"] = *notebook.Name
	nbService.Metadata.Namespace = *notebook.Namespace
	nbService.Spec.Type_ = "NodePort"

	startPort, err := strconv.Atoi(os.Getenv("START_PORT"))
	if nil != err {
		logger.Logger().Errorf("parse START_PORT failed: %s", os.Getenv("START_PORT"))
		return nil, [2]int{0, 0}, err
	}
	endPort, err := strconv.Atoi(os.Getenv("END_PORT"))
	if nil != err {
		logger.Logger().Errorf("parse END_PORT failed: %s", os.Getenv("END_PORT"))
		return nil, [2]int{0, 0}, err
	}

	_, nodePortArray := getNodePort(startPort, endPort)
	logger.Logger().Info("Service NodePort:" + strconv.Itoa(nodePortArray[0]) + " and " + strconv.Itoa(nodePortArray[1]))
	nbService.Spec.Ports = []client.V1ServicePort{{
		Name:     "sparkcontextport",
		NodePort: int32(nodePortArray[0]),
		Port:     int32(nodePortArray[0]),
		Protocol: "TCP",
	}, {
		Name:     "sparkblockmanagerport",
		NodePort: int32(nodePortArray[1]),
		Port:     int32(nodePortArray[1]),
		Protocol: "TCP",
	},
	}

	nbService.Spec.Selector = nbServiceSelector
	return nbService, nodePortArray, nil
}

func getNodePort(startPort int, endPort int) (error, [2]int) {

	logger.Logger().Info("ConfigMap")
	clusterService, _, err := ncClient.k8sClient.CoreV1Api.ListServiceForAllNamespaces(context.Background(), nil)

	portLen := (endPort - startPort) + 1
	portArray := make([]bool, portLen, portLen)
	nodePortArray := [2]int{-1, -1}

	if err != nil {
		logrus.Error(err.Error())
		return err, nodePortArray
	}

	for _, serviceItem := range clusterService.Items {
		if serviceItem.Spec.Type_ == "NodePort" {
			usedNodeportList := serviceItem.Spec.Ports
			for _, usedNodeportItem := range usedNodeportList {
				usedNodeport := int(usedNodeportItem.NodePort)
				if usedNodeport >= startPort && usedNodeport <= endPort {
					logger.Logger().Info("NodePort:" + strconv.Itoa(usedNodeport) + " is allocated")
					portArray[usedNodeport-startPort] = true
				}
			}
		}
	}
	flag := 0
	for i, value := range portArray {
		if value != true {
			nodePortArray[flag] = i + startPort
			flag = flag + 1
			if flag >= 2 {
				break
			}
		}
	}

	return err, nodePortArray
}

func (*NotebookControllerClient) CreateYarnResourceConfigMap(notebook *models.NewNotebookRequest, userId string) error {
	configmap, _, err := ncClient.k8sClient.
		CoreV1Api.ReadNamespacedConfigMap(context.Background(), "yarn-resource-setting", *notebook.Namespace, nil)

	if err != nil {
		logger.Logger().Errorf(err.Error())
		return err
	}
	var yarnConfigmap map[string]interface{}
	err = json.Unmarshal([]byte(configmap.Data["config.json"]), &yarnConfigmap)
	if err != nil {
		logger.Logger().Info("json yarn configmap error")
		logrus.Error(err.Error())
		return err
	}
	var linkismagicCM map[string]interface{}
	err = json.Unmarshal([]byte(configmap.Data["linkismagic.json"]), &linkismagicCM)
	if err != nil {
		logger.Logger().Error("json linkismagic configmap error:%v", err.Error())
		return err
	}
	//logr.Info("testtsetststest1:" + yarnConfigmap["fatal_error_suggestion"].(string))
	yarnConfigmap["session_configs"].(map[string]interface{})["proxyUser"] = userId
	if "" != notebook.Queue {
		yarnConfigmap["session_configs"].(map[string]interface{})["queue"] = notebook.Queue
		linkismagicCM["session_configs"].(map[string]interface{})["wds.linkis.yarnqueue"] = notebook.Queue
		// default setting
		yarnConfigmap["session_configs"].(map[string]interface{})["driverMemory"] = "2g"
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.cores"] = "2"
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.instances"] = "2"
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.memory"] = "2g"

		linkismagicCM["session_configs"].(map[string]interface{})["spark.driver.memory"] = "2"
		linkismagicCM["session_configs"].(map[string]interface{})["spark.executor.instances"] = "2"
		linkismagicCM["session_configs"].(map[string]interface{})["spark.executor.memory"] = "2"
	}
	if "" != notebook.DriverMemory {
		yarnConfigmap["session_configs"].(map[string]interface{})["driverMemory"] = notebook.DriverMemory + "g"
		linkismagicCM["session_configs"].(map[string]interface{})["spark.driver.memory"] = notebook.DriverMemory
	}
	if "" != notebook.ExecutorCores {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.cores"] = notebook.ExecutorCores
	}
	if "" != notebook.Executors {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.instances"] = notebook.Executors
		linkismagicCM["session_configs"].(map[string]interface{})["spark.executor.instances"] = notebook.Executors
	}
	if "" != notebook.ExecutorMemory {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.memory"] = notebook.ExecutorMemory + "g"
		linkismagicCM["session_configs"].(map[string]interface{})["spark.executor.memory"] = notebook.ExecutorMemory
	}
	configByte, _ := json.Marshal(yarnConfigmap)
	linkismagicCMByte, _ := json.Marshal(linkismagicCM)
	configJsonStr := string(configByte)
	linkismagicStr := string(linkismagicCMByte)
	nsConfigmap := client.V1ConfigMap{
		Metadata: &client.V1ObjectMeta{
			Name:      *notebook.Name + "-" + "yarn-resource-setting",
			Namespace: *notebook.Namespace,
		},
		Data: map[string]string{
			"config.json":      configJsonStr,
			"linkismagic.json": linkismagicStr,
		},
	}
	_, _, err = ncClient.k8sClient.CoreV1Api.CreateNamespacedConfigMap(context.Background(), *notebook.Namespace, nsConfigmap, nil)
	logger.Logger().Info("Create Configmap yarn resource setting for notebook " + *notebook.Namespace)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	return err
}

func (*NotebookControllerClient) PatchYarnSettingConfigMap(notebook *models.PatchNotebookRequest) error {
	namespace := *notebook.Namespace
	nbName := *notebook.Name
	configmap, _, err := ncClient.k8sClient.CoreV1Api.ReadNamespacedConfigMap(context.Background(), nbName+"-"+"yarn-resource-setting", namespace, nil)

	logger.Logger().Info("read configmap from notebook" + nbName)
	logger.Logger().Info("update yarn resource 1、 queue: " + notebook.Queue)
	if err != nil {
		logrus.Error("get configmap errror:" + err.Error())
		return err
	}
	var yarnConfigmap map[string]interface{}
	var linkismagicConfigmap map[string]interface{}
	err = json.Unmarshal([]byte(configmap.Data["config.json"]), &yarnConfigmap)
	if err != nil {
		logrus.Error("unmarshal yarn json erro" + err.Error())
		return err
	}
	err = json.Unmarshal([]byte(configmap.Data["linkismagic.json"]), &linkismagicConfigmap)
	if err != nil {
		logrus.Error("unmarshal linkismagic json error" + err.Error())
		return err
	}
	configJsonStr := ""
	lmConfigJsonStr := ""

	if "" != notebook.Queue {
		yarnConfigmap["session_configs"].(map[string]interface{})["queue"] = notebook.Queue
		linkismagicConfigmap["session_configs"].(map[string]interface{})["wds.linkis.yarnqueue"] = notebook.Queue
	}
	if "" != notebook.DriverMemory {
		yarnConfigmap["session_configs"].(map[string]interface{})["driverMemory"] = notebook.DriverMemory
		linkismagicConfigmap["session_configs"].(map[string]interface{})["spark.driver.memory"] = notebook.DriverMemory[0 : len(notebook.ExecutorMemory)-1]
	}
	if "" != notebook.ExecutorCores {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.cores"] = notebook.ExecutorCores
	}
	if "" != notebook.Executors {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.instances"] = notebook.Executors
		linkismagicConfigmap["session_configs"].(map[string]interface{})["spark.executor.instances"] = notebook.Executors
	}
	if "" != notebook.ExecutorMemory {
		yarnConfigmap["session_configs"].(map[string]interface{})["conf"].(map[string]interface{})["spark.executor.memory"] = notebook.ExecutorMemory
		linkismagicConfigmap["session_configs"].(map[string]interface{})["spark.executor.memory"] = notebook.ExecutorMemory[0 : len(notebook.ExecutorMemory)-1]
	}

	configByte, err := json.Marshal(yarnConfigmap)
	if err != nil {
		logrus.Error("marshal yarnConfigmap error" + err.Error())
		return err
	}
	linkismagicByte, err := json.Marshal(linkismagicConfigmap)
	if err != nil {
		logrus.Error("marshal linkismagic error" + err.Error())
		return err
	}
	configJsonStr = string(configByte)
	lmConfigJsonStr = string(linkismagicByte)
	logger.Logger().Info("read configmap from configJsonStr" + configJsonStr)
	opr_map := []map[string]interface{}{
		{"op": "replace",
			"path": "/data",
			"value": map[string]interface{}{
				"config.json":      configJsonStr,
				"linkismagic.json": lmConfigJsonStr,
			}},
	}

	_, _, err = ncClient.k8sClient.CoreV1Api.PatchNamespacedConfigMap(context.Background(), nbName+"-"+"yarn-resource-setting", namespace, opr_map, nil)
	if err != nil {
		logrus.Error("PatchNamespacedConfigMap" + err.Error())
		return err
	}
	return err
}

func deleteK8sNBConfigMap(namespace string, name string) error {
	logger.Logger().Info("Delete configmap " + name + " in namespace " + namespace)
	_, resp, err := ncClient.k8sClient.CoreV1Api.DeleteNamespacedConfigMap(context.Background(), name, namespace, client.V1DeleteOptions{}, nil)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	if resp == nil {
		logger.Logger().Error("empty response")
	}
	return err
}

func (*NotebookControllerClient) GetNBConfigMaps(namespace string) (*map[string]interface{}, error) {
	var list client.V1ConfigMapList
	var response *http.Response
	var err error
	configs := map[string]interface{}{}
	if "" == namespace || "null" == namespace {
		logger.Logger().Infof("ListConfigMapForAllNamespaces:%v", namespace)
		list, response, err = ncClient.k8sClient.CoreV1Api.ListConfigMapForAllNamespaces(context.Background(), nil)
	} else {
		logger.Logger().Infof("ListNamespacedConfigMap namespsace:%v", namespace)
		list, response, err = ncClient.k8sClient.CoreV1Api.ListNamespacedConfigMap(context.Background(), namespace, nil)
	}
	if err != nil {
		logger.Logger().Error(err.Error())
		return &configs, err
	}
	if response == nil {
		logger.Logger().Error("emtpy response")
		return nil, nil
	}

	logger.Logger().Infof("configmap length:%v", len(list.Items))
	for _, configmap := range list.Items {
		logger.Logger().Infof("configmap name:%v", configmap.Metadata.Name)
		if !stringsUtil.Contains(configmap.Metadata.Name, "yarn-resource-setting") {
			continue
		}
		var yarnConfigmap map[string]interface{}
		err = json.Unmarshal([]byte(configmap.Data["config.json"]), &yarnConfigmap)
		if err != nil {
			logger.Logger().Errorf("configmap "+configmap.Metadata.Name+" unmarshal error:%v", err.Error())
		}
		configs[configmap.Metadata.Namespace+"-"+configmap.Metadata.Name] = yarnConfigmap
	}
	return &configs, err
}
