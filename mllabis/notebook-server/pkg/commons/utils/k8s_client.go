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
	"io/ioutil"
	"net/http"
	"os"
	//"path"
	"strconv"
	stringsUtil "strings"
	"sync"
	"time"
	cfg "webank/AIDE/notebook-server/pkg/commons/config"
	"webank/AIDE/notebook-server/pkg/commons/constants"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"

	"github.com/kubernetes-client/go/kubernetes/client"
	"github.com/kubernetes-client/go/kubernetes/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PortMap struct {
	Mutex sync.Mutex
	M     map[int]struct{}
}

var Ports *PortMap = &PortMap{
	Mutex: sync.Mutex{},
	M:     make(map[int]struct{}),
}

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
func (*NotebookControllerClient) CreateNotebook(notebook *models.NewNotebookRequest,
	userId string, gid string, uid string, token string) ([]byte, error) {
	/*** 0. get a nb template for k8s request && get cluster message from namespaces ***/
	nbForK8s := createK8sNBTemplate()
	//Set Container User
	containerUser := userId
	if notebook.ProxyUser != "" {
		containerUser = notebook.ProxyUser
	}

	nbContainer := nbForK8s.Spec.Template.Spec.Containers[0]

	/*** 1. volumes ***/
	var nbVolumes []client.V1Volume
	var localPath string
	if nil != notebook.WorkspaceVolume {
		localPath = *notebook.WorkspaceVolume.LocalPath
		if len(userId) > len(localPath) {
			gid = viper.GetString(cfg.MLSSGroupId)
		} else if userId != localPath[len(localPath)-len(userId):] {
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
		localPath := *notebook.DataVolume.LocalPath
		// if local path directory is end with user Id, this directory will be considered as a common directory.
		if len(userId) > len(localPath) {
			gid = viper.GetString(cfg.MLSSGroupId)
		} else if userId != localPath[len(localPath)-len(userId):] {
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
		}, client.V1Volume{
			Name: "climagic-json",
			ConfigMap: &client.V1ConfigMapVolumeSource{
				DefaultMode: 0755,
				Name:        *notebook.Name + "-" + "yarn-resource-setting",
			},
		})

	nbForK8s.Spec.Template.Spec.Volumes = append(nbForK8s.Spec.Template.Spec.Volumes, nbVolumes...)

	/*** 3. container ***/
	nbContainer.Command = []string{"/etc/config/start-sh"}
	//nbContainer.Args = []string{" --notebook-dir=/home/jovyan --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"}
	nbContainer.Args = []string{"jupyter lab --notebook-dir=/home/" + containerUser + "/workspace --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.disable_check_xsrf=True  --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"}
	var extraResourceMap = make(map[string]string)
	if notebook.ExtraResources != "" {
		err := json.Unmarshal([]byte(notebook.ExtraResources), &extraResourceMap)
		if nil != err {
			logger.Logger().Errorf("Error parsing extra resources: %s", notebook.ExtraResources)
			logger.Logger().Errorf(err.Error())
			//return
			return nil, err
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
	if nil != notebook.DataVolume && nil != notebook.DataVolume.MountPath {
		mountPath := *notebook.DataVolume.MountPath
		if mountPath != "" {
			if !stringsUtil.HasPrefix(mountPath, "/") && !stringsUtil.HasSuffix(mountPath, "/") {
				if stringsUtil.Contains(mountPath, "/") {
					return nil, errors.New("mountPath can only have level 1 directories")
				}
			} else if stringsUtil.HasPrefix(mountPath, "/") && !stringsUtil.HasSuffix(mountPath, "/") {
				if stringsUtil.Contains(mountPath[1:], "/") {
					return nil, errors.New("mountPath can only have level 1 directories")
				}
			} else if !stringsUtil.HasPrefix(mountPath, "/") && stringsUtil.HasSuffix(mountPath, "/") {
				if stringsUtil.Index(mountPath, "/") < len(mountPath)-1 {
					return nil, errors.New("mountPath can only have level 1 directories")
				}
			} else if stringsUtil.HasPrefix(mountPath, "/") && stringsUtil.HasSuffix(mountPath, "/") {
				if stringsUtil.Contains(mountPath[1:(len(mountPath)-1)], "/") {
					return nil, errors.New("mountPath can only have level 1 directories")
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
		Value: containerUser,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "CREATE_USER",
		Value: userId,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "PROXY_USER",
		Value: notebook.ProxyUser,
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "GRANT_SUDO",
		Value: "no",
	})

	// add SPARK&HDFS Container Config if cluster is BDAP
	nodePortArrayInSparkSessionService := make([]int, 0)
	//create NodePort Service
	notebookService, nodePortArray, err := createK8sNBServiceTemplate(notebook) //todo 待更新
	defer func(nodePortArray []int) {
		if len(nodePortArray) > 0 {
			for _, port := range nodePortArray {
				if _, ok := Ports.M[port]; ok {
					delete(Ports.M, port)
					logger.Logger().Debugf("Port '%d' has been released\n", port)
				}
			}
			Ports.Mutex.Unlock()
		}

	}(nodePortArray)
	if nil != err {
		logger.Logger().Errorf(err.Error())
		return nil, err
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
		Value: "MLSS",
	})
	nbContainer.Env = append(nbContainer.Env, client.V1EnvVar{
		Name:  "GUARDIAN_TOKEN",
		Value: token,
	})
	nodePortArrayInSparkSessionService = nodePortArray[2:]
	//}

	// create notebook spark session nodePort
	sparkBytes := make([]byte, 0)
	logger.Logger().Debugf("CreateNotebook notebook.SparkSessionNum: %v, "+
		"nodePortArrayInSparkSessionService: %v\n", notebook.SparkSessionNum, nodePortArrayInSparkSessionService)
	if notebook.SparkSessionNum > 0 && len(nodePortArrayInSparkSessionService) > 0 {
		sparkSessionMap, sparkService, err := createK8sNBSparkSessionServiceTemplate(notebook, nodePortArrayInSparkSessionService)
		if err != nil {
			logger.Logger().Errorf("fail to create spark session service,"+
				" notebook param: %+v, nodePortArrayInSparkSessionService: %+v, err: %s\n",
				*notebook, nodePortArrayInSparkSessionService, err.Error())
			return nil, err
		}
		sparkBytes, err = json.Marshal(sparkSessionMap)
		if err != nil {
			logger.Logger().Errorf("fail to marshal spark session map, sparkSessionMap: %+v, err: %s\n",
				sparkSessionMap, err.Error())
			return nil, err
		}
		nbForK8s.SparkSessionService = *sparkService
	} else {
		if len(nodePortArrayInSparkSessionService) == 0 {
			logger.Logger().Infoln("cluster is not 'bdap', not create spark session service")
		}
	}

	nbForK8s.Spec.Template.Spec.Containers[0] = nbContainer

	/*** 4. fill in other info ***/
	nbForK8s.Metadata.Name = *(notebook.Name)
	nbForK8s.Metadata.Namespace = *(notebook.Namespace)

	userIdLabel := PlatformNamespace + "-UserId"
	proxyUserIdLabel := PlatformNamespace + "-ProxyUserId"
	workDirLabel := PlatformNamespace + "-WorkdDir"
	//set notebook create time in label
	nowTimeStr := strconv.FormatInt(time.Now().Unix(), 10)

	//nbForK8s.Metadata.Labels = map[string]string{"MLSS-UserId": userId}
	var workDirPath string
	if nil != notebook.DataVolume && notebook.DataVolume.LocalPath != nil && *notebook.DataVolume.LocalPath != "" {
		workDirPath = *notebook.DataVolume.LocalPath
	} else {
		workDirPath = *notebook.WorkspaceVolume.LocalPath
	}
	logger.Logger().Info("workDirPath is ", workDirPath)
	nbForK8s.Metadata.Labels = map[string]string{
		userIdLabel:          userId,
		proxyUserIdLabel:     notebook.ProxyUser,
		workDirLabel:         stringsUtil.Replace(workDirPath[1:], "/", "_", -1),
		constants.CreateTime: nowTimeStr,
	}

	_, resp, err := ncClient.k8sClient.CustomObjectsApi.CreateNamespacedCustomObject(context.Background(),
		NBApiGroup, NBApiVersion, nbForK8s.Metadata.Namespace, NBApiPlural, nbForK8s, nil)
	if nil != err {
		logger.Logger().Errorf("Error creating notebook: %+v", err.Error())
		return nil, err
	}
	if resp != nil {
		logger.Logger().Debugf("response from k8s: %+v", resp)
	}
	return sparkBytes, err
}

// delete notebook in a namespace
func (*NotebookControllerClient) DeleteNotebook(nb string, ns string) error {
	_, resp, err := ncClient.k8sClient.CustomObjectsApi.DeleteNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, ns, NBApiPlural, nb, client.V1DeleteOptions{}, nil)
	if err != nil {
		if !stringsUtil.Contains(err.Error(), "not found") {
			logger.Logger().Errorf("deleting notebook: %+v", err.Error())
			return err
		}
	}
	if resp == nil {
		logger.Logger().Error("empty response")
	}

	err = deleteK8sNBConfigMap(ns, nb+"-"+"yarn-resource-setting")
	if err != nil {
		if !stringsUtil.Contains(err.Error(), "not found") {
			logger.Logger().Errorf("deleting notebook configmap: %+v", err.Error())
			return err
		}
	}
	//delete spark service
	err = deleteK8sSparkSessionService(ns, nb+"-spark-session-service")
	if err != nil {
		if !stringsUtil.Contains(err.Error(), "not found") {
			logger.Logger().Errorf("fail to delete spark session service, name:%s, namespace: %s, "+
				"err: %s\n", nb+"-spark-session-service", ns, err.Error())
			return err
		}
	}
	return nil
}

func (*NotebookControllerClient) GetNotebookByNamespaceAndName(namespace, name string) (interface{}, error) {
	if namespace == "" || name == "" {
		return nil, fmt.Errorf("namespace %q or notebook name %q can't be empty", namespace, name)
	}
	nb, resp, err := ncClient.k8sClient.CustomObjectsApi.GetNamespacedCustomObject(context.Background(),
		NBApiGroup, NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}
	if resp == nil {
		logger.Logger().Error("emtpy response")
		return nil, nil
	}
	return nb, nil
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
		workDirLabel = fmt.Sprint(PlatformNamespace+"-WorkDir=", workDir)
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

func (*NotebookControllerClient) ListNotebooks(ns string) (interface{}, error) {
	list, resp, err := ncClient.k8sClient.CustomObjectsApi.ListNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, ns, NBApiPlural, nil)
	if err != nil {
		logger.Logger().Error(err.Error())
		return nil, err
	}
	logger.Logger().Debugf("ListNotebooks notebook list: %v", list)
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
	for idx, v := range list.Items {
		logger.Logger().Debugf("listNamespacedResourceQuota[%d] status: %+v, spec: %+v, "+
			"Metadata: %+v", idx, *v.Status, *v.Spec, *v.Metadata)
	}
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
			Name:              "",
			Namespace:         "",
			CreationTimestamp: time.Now(),
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
									Name:      "timezone",
									MountPath: "/etc/timezone",
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
							//Args:    []string{"jupyter notebook --notebook-dir=/home/jovyan --ip=0.0.0.0 --no-browser --allow-root --port=8888 --NotebookApp.token='' --NotebookApp.password='' --NotebookApp.allow_origin='*' --NotebookApp.base_url=${NB_PREFIX}"},
							// FIXME not working
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
							Name: "timezone",
							HostPath: &client.V1HostPathVolumeSource{
								Path: "/etc/timezone",
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

func createK8sNBSparkSessionServiceTemplate(notebook *models.NewNotebookRequest, portArray []int) (map[string]interface{}, *client.V1Service, error) {
	if len(portArray) == 0 || len(portArray) != int(notebook.SparkSessionNum*2) {
		logger.Logger().Errorf("port array number error, portArray: %v", portArray)
		return nil, nil, fmt.Errorf("port array number error, portArray: %v", portArray)
	}
	serviceName := *(notebook.Name) + "-spark-session-service"
	//ncClient.k8sClient.AppsV1Api.CreateNamespacedStatefulSet()
	selector := map[string]string{"statefulset": *notebook.Name}
	servicePortArr := []client.V1ServicePort{}

	sparkSessionMap := make(map[string]interface{})
	sparkPortArray := make([]map[string]string, 0)
	for i := int64(0); i < notebook.SparkSessionNum; i++ {
		sparkDriverServicePort := client.V1ServicePort{
			Name:       fmt.Sprintf("spark-driver-port-%d", i),
			NodePort:   int32(portArray[2*i]),
			Port:       int32(portArray[2*i]),
			Protocol:   "TCP",
			TargetPort: nil,
		}
		sparkDriverBlockManagerServicePort := client.V1ServicePort{
			Name:       fmt.Sprintf("spark-driver-block-manager-port-%d", i),
			NodePort:   int32(portArray[2*i+1]),
			Port:       int32(portArray[2*i+1]),
			Protocol:   "TCP",
			TargetPort: nil,
		}
		servicePortArr = append(servicePortArr, sparkDriverServicePort, sparkDriverBlockManagerServicePort)
		m := make(map[string]string)
		m["Spark_Driver_Port"] = strconv.FormatInt(int64(portArray[2*i]), 10)
		m["spark.driver.blockManager.port"] = strconv.FormatInt(int64(portArray[2*i]+1), 10)
		sparkPortArray = append(sparkPortArray, m)
	}
	sparkSessionMap["portArray"] = sparkPortArray
	sparkSessionMap["spark.driver.memory"] = notebook.DriverMemory + "g"
	sparkSessionMap["spark.executor.memory"] = notebook.ExecutorMemory + "g"
	sparkSessionMap["spark.executor.instances"] = notebook.Executors
	body := client.V1Service{
		ApiVersion: "",
		Kind:       "",
		Metadata: &client.V1ObjectMeta{
			Annotations:                nil,
			ClusterName:                "",
			CreationTimestamp:          time.Now(),
			DeletionGracePeriodSeconds: 0,
			DeletionTimestamp:          time.Time{},
			Finalizers:                 nil,
			GenerateName:               "",
			Generation:                 0,
			Initializers:               nil,
			Labels:                     nil,
			Name:                       serviceName,
			Namespace:                  *notebook.Namespace,
			OwnerReferences:            nil,
			ResourceVersion:            "",
			SelfLink:                   "",
			Uid:                        "",
		},
		Spec: &client.V1ServiceSpec{
			ClusterIP:                "",
			ExternalIPs:              nil,
			ExternalName:             "",
			ExternalTrafficPolicy:    "",
			HealthCheckNodePort:      0,
			LoadBalancerIP:           "",
			LoadBalancerSourceRanges: nil,
			Ports:                    servicePortArr,
			PublishNotReadyAddresses: false,
			Selector:                 selector,
			SessionAffinity:          "",
			SessionAffinityConfig:    nil,
			Type_:                    "NodePort",
		},
		Status: nil,
	}
	svc, resp, err := ncClient.k8sClient.CoreV1Api.CreateNamespacedService(context.TODO(), *notebook.Namespace, body, nil)
	logger.Logger().Debugf("createK8sNBSparkSessionServiceTemplate service: %+v, resp: %+v\n", svc, *resp)
	if err != nil {
		logger.Logger().Debugf("fail to create spark session service, err: %s\n", err.Error())
		return nil, nil, err
	}
	return sparkSessionMap, &svc, nil
}

func createK8sNBServiceTemplate(notebook *models.NewNotebookRequest) (*client.V1Service, []int, error) {

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
		return nil, nil, err
	}
	endPort, err := strconv.Atoi(os.Getenv("END_PORT"))
	if nil != err {
		logger.Logger().Errorf("parse END_PORT failed: %s", os.Getenv("END_PORT"))
		return nil, nil, err
	}

	portNum := int((notebook.SparkSessionNum * 2) + 2)
	_, nodePortArray := getNodePort(startPort, endPort, portNum)
	if len(nodePortArray) < portNum {
		logger.Logger().Errorf("nodePort number must be larger than 2, nodePortArray: %v", nodePortArray)
		return nil, nodePortArray, err
	}
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

// 该函数返回后需解锁全局变量
func getNodePort(startPort int, endPort int, portNum int) (error, []int) {

	logger.Logger().Info("ConfigMap")
	clusterService, _, err := ncClient.k8sClient.CoreV1Api.ListServiceForAllNamespaces(context.Background(), nil)

	portLen := (endPort - startPort) + 1
	portArray := make([]bool, portLen, portLen)
	//nodePortArray := [2]int{-1, -1}
	nodePortArray := make([]int, portNum, portNum)
	for i := range nodePortArray {
		nodePortArray[i] = -1
	}

	if err != nil {
		logrus.Error(err.Error())
		return err, nil
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
	Ports.Mutex.Lock()
	for i, value := range portArray {
		if value != true {
			nodePortArray[flag] = i + startPort
			Ports.M[i+startPort] = struct{}{}
			flag = flag + 1
			if flag >= portNum {
				break
			}
		}
	}
	return err, nodePortArray
}

func (*NotebookControllerClient) CreateYarnResourceConfigMap(notebook *models.NewNotebookRequest, userId string, sparkSessionBytes []byte) error {
	configmap, _, err := ncClient.k8sClient.CoreV1Api.ReadNamespacedConfigMap(context.Background(), "yarn-resource-setting", *notebook.Namespace, nil)
	if err != nil {
		logrus.Error(err.Error())
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
		logger.Logger().Info("json linkismagic configmap error")
		logrus.Error(err.Error())
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
	sparkSessionStr := ""
	if len(sparkSessionBytes) > 0 {
		sparkSessionStr = string(sparkSessionBytes)
	}
	nsConfigmap := client.V1ConfigMap{
		Metadata: &client.V1ObjectMeta{
			Name:      *notebook.Name + "-" + "yarn-resource-setting",
			Namespace: *notebook.Namespace,
		},
		Data: map[string]string{
			"config.json":       configJsonStr,
			"linkismagic.json":  linkismagicStr,
			"sparkSession.json": sparkSessionStr,
		},
	}

	_, _, err = ncClient.k8sClient.CoreV1Api.CreateNamespacedConfigMap(context.Background(), *notebook.Namespace, nsConfigmap, nil)
	logger.Logger().Info("Create Configmap yarn resource setting for notebook " + *notebook.Namespace)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	if err == nil {
		logger.Logger().Info("success")
		return err
	}
	return err
}

func (*NotebookControllerClient) PatchSparkSessionService(notebook *models.PatchNotebookRequest, name, namespace string) (
	map[string]interface{}, *client.V1Service, []int, error) {
	serviceName := name + "-spark-session-service"

	sparkPortArray := make([]map[string]string, 0)
	sparkSessionMap := make(map[string]interface{})
	servicePortArr := []client.V1ServicePort{}
	portArray := make([]int, 0)
	if notebook.SparkSessionNum > 0 {
		startPort, err := strconv.Atoi(os.Getenv("START_PORT"))
		if nil != err {
			logger.Logger().Errorf("parse START_PORT failed: %s", os.Getenv("START_PORT"))
			return nil, nil, portArray, err
		}
		endPort, err := strconv.Atoi(os.Getenv("END_PORT"))
		if nil != err {
			logger.Logger().Errorf("parse END_PORT failed: %s", os.Getenv("END_PORT"))
			return nil, nil, portArray, err
		}
		err, portArray = getNodePort(startPort, endPort, int(notebook.SparkSessionNum*2))
		if err != nil {
			return nil, nil, portArray, err
		}

		if len(portArray) == 0 || len(portArray) != int(notebook.SparkSessionNum*2) {
			logger.Logger().Errorf("port array number error, portArray: %v", portArray)
			return nil, nil, portArray, fmt.Errorf("port array number error, portArray: %v", portArray)
		}
		for i := int64(0); i < notebook.SparkSessionNum; i++ {
			sparkDriverServicePort := client.V1ServicePort{
				Name:       fmt.Sprintf("spark-driver-port-%d", i),
				NodePort:   int32(portArray[2*i]),
				Port:       int32(portArray[2*i]),
				Protocol:   "TCP",
				TargetPort: nil,
			}
			sparkDriverBlockManagerServicePort := client.V1ServicePort{
				Name:       fmt.Sprintf("spark-driver-block-manager-port-%d", i),
				NodePort:   int32(portArray[2*i+1]),
				Port:       int32(portArray[2*i+1]),
				Protocol:   "TCP",
				TargetPort: nil,
			}
			servicePortArr = append(servicePortArr, sparkDriverServicePort, sparkDriverBlockManagerServicePort)
			m := make(map[string]string)
			m["Spark_Driver_Port"] = strconv.FormatInt(int64(portArray[2*i]), 10)
			m["spark.driver.blockManager.port"] = strconv.FormatInt(int64(portArray[2*i]+1), 10)
			sparkPortArray = append(sparkPortArray, m)
		}
	}
	sparkSessionMap["portArray"] = sparkPortArray
	sparkSessionMap["spark.driver.memory"] = notebook.DriverMemory
	sparkSessionMap["spark.executor.memory"] = notebook.ExecutorMemory
	sparkSessionMap["spark.executor.instances"] = notebook.Executors

	logger.Logger().Debugf("spark ports has been allocated, ports: %v\n", servicePortArr)
	selector := map[string]string{"statefulset": name}

	service, _, err := ncClient.k8sClient.CoreV1Api.ReadNamespacedService(context.TODO(),
		serviceName, namespace, nil)
	if err != nil {
		if stringsUtil.Contains(err.Error(), "not found") || stringsUtil.Contains(err.Error(), "Not Found") {
			if len(servicePortArr) > 0 {
				serviceObj := client.V1Service{
					ApiVersion: "",
					Kind:       "",
					Metadata: &client.V1ObjectMeta{
						Annotations:                nil,
						ClusterName:                "",
						CreationTimestamp:          time.Now(),
						DeletionGracePeriodSeconds: 0,
						DeletionTimestamp:          time.Time{},
						Finalizers:                 nil,
						GenerateName:               "",
						Generation:                 0,
						Initializers:               nil,
						Labels:                     nil,
						Name:                       serviceName,
						Namespace:                  namespace,
						OwnerReferences:            nil,
						ResourceVersion:            "",
						SelfLink:                   "",
						Uid:                        "",
					},
					Spec: &client.V1ServiceSpec{
						ClusterIP:                "",
						ExternalIPs:              nil,
						ExternalName:             "",
						ExternalTrafficPolicy:    "",
						HealthCheckNodePort:      0,
						LoadBalancerIP:           "",
						LoadBalancerSourceRanges: nil,
						Ports:                    servicePortArr,
						PublishNotReadyAddresses: false,
						Selector:                 selector,
						SessionAffinity:          "",
						SessionAffinityConfig:    nil,
						Type_:                    "NodePort",
					},
					Status: nil,
				}
				createdService, _, err := ncClient.k8sClient.CoreV1Api.CreateNamespacedService(context.TODO(),
					namespace, serviceObj, nil)
				if err != nil {
					logger.Logger().Errorf("fail to create spark service in update spark service,"+
						" service: %+v, err: %v\n", serviceObj, err)
					return nil, nil, portArray, err
				}
				return sparkSessionMap, &createdService, portArray, nil
			} else {
				return sparkSessionMap, nil, portArray, nil
			}
		}
		logger.Logger().Errorf("fail to read namespaced service, service name: %s, "+
			"namespace: %s, err: %v\n", serviceName, namespace, err)
		return nil, nil, portArray, err
	}
	logger.Logger().Debugf("PatchSparkSessionService spark session service: %+v\n", service) //todo
	logger.Logger().Debugf("PatchSparkSessionService spark session service.spec: %+v\n", *service.Spec)
	if len(service.Spec.Ports) != int(notebook.SparkSessionNum*2) {
		if len(servicePortArr) == 0 {
			// delete spark session service
			logger.Logger().Infoln("spark session service ports is nil")
			_, _, err := ncClient.k8sClient.CoreV1Api.DeleteNamespacedService(context.TODO(),
				serviceName, namespace, client.V1DeleteOptions{}, nil)
			if err != nil {
				logger.Logger().Errorf("fail to delete spark session service in update spark, "+
					"service name: %s, namespace: %s, err: %v\n", serviceName, namespace, err)
				return nil, nil, portArray, err
			}
			logger.Logger().Debugln("PatchSparkSessionService servicePortArr is 0")
			return sparkSessionMap, nil, portArray, nil
		} else {
			// 使用 service update方法有问题，故改成先删除后新建的方式更新service
			_, _, err := ncClient.k8sClient.CoreV1Api.DeleteNamespacedService(context.TODO(),
				serviceName, namespace, client.V1DeleteOptions{}, nil)
			if err != nil {
				logger.Logger().Errorf("fail to delete spark session service in update spark, "+
					"service name: %s, namespace: %s, err: %v\n", serviceName, namespace, err)
				return nil, nil, portArray, err
			}
			logger.Logger().Infof("delete spark session service %q success.\n", serviceName)
			serviceObj := client.V1Service{
				ApiVersion: "",
				Kind:       "",
				Metadata: &client.V1ObjectMeta{
					Annotations:                nil,
					ClusterName:                "",
					CreationTimestamp:          time.Now(),
					DeletionGracePeriodSeconds: 0,
					DeletionTimestamp:          time.Time{},
					Finalizers:                 nil,
					GenerateName:               "",
					Generation:                 0,
					Initializers:               nil,
					Labels:                     nil,
					Name:                       serviceName,
					Namespace:                  namespace,
					OwnerReferences:            nil,
					ResourceVersion:            "",
					SelfLink:                   "",
					Uid:                        "",
				},
				Spec: &client.V1ServiceSpec{
					ClusterIP:                "",
					ExternalIPs:              nil,
					ExternalName:             "",
					ExternalTrafficPolicy:    "",
					HealthCheckNodePort:      0,
					LoadBalancerIP:           "",
					LoadBalancerSourceRanges: nil,
					Ports:                    servicePortArr,
					PublishNotReadyAddresses: false,
					Selector:                 selector,
					SessionAffinity:          "",
					SessionAffinityConfig:    nil,
					Type_:                    "NodePort",
				},
				Status: nil,
			}
			createdService, _, err := ncClient.k8sClient.CoreV1Api.CreateNamespacedService(context.TODO(),
				namespace, serviceObj, nil)
			if err != nil {
				logger.Logger().Errorf("fail to create spark service in update spark service,"+
					" service: %+v, err: %v\n", serviceObj, err)
				return nil, nil, portArray, err
			}
			logger.Logger().Infof("create spark session service %q success.\n", serviceName)
			return sparkSessionMap, &createdService, portArray, nil
		}
	}
	logger.Logger().Debugln("PatchSparkSessionService not update")
	logger.Logger().Debugf("PatchSparkSessionService not update, service.Spec.Ports: %v,"+
		"notebook.SparkSessionNum: %d\n", service.Spec.Ports, notebook.SparkSessionNum)
	return nil, &service, portArray, nil
}

func (*NotebookControllerClient) PatchNoteBookCRD(namespace, name string, sparkService *client.V1Service,
	nbReq *models.PatchNotebookRequest) error {
	body, _, err := ncClient.k8sClient.CustomObjectsApi.GetNamespacedCustomObject(context.TODO(), NBApiGroup,
		NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook CRD, namespace: %s, name: %s, "+
			"NBApiGroup: %s, NBApiVersion: %s, NBApiPlural: %s, err: %v\n", namespace, name,
			NBApiGroup, NBApiVersion, NBApiPlural, err)
		return err
	}
	logger.Logger().Infof("PatchNoteBookCRD body: %v\n", body)
	notebook := models.K8sNotebook{}
	nbBytes, err := json.Marshal(body)
	if err != nil {
		logger.Logger().Errorf("fail to marshal crd body: %v, err: %v\n", body, err)
		return err
	}
	err = json.Unmarshal(nbBytes, &notebook)
	if err != nil {
		logger.Logger().Errorf("fail to unmarshal crd body to notebook, bytes: %v, err: %v\n", body, err)
		return err
	}

	// check resource
	//todo
	limitMap := notebook.Spec.Template.Spec.Containers[0].Resources.Limits
	requestMap := notebook.Spec.Template.Spec.Containers[0].Resources.Requests
	if nbReq.CPU > 0 {
		limitMap["cpu"] = fmt.Sprint(nbReq.CPU)
		requestMap["cpu"] = fmt.Sprint(nbReq.CPU)
	}
	if nbReq.MemoryAmount > 0 {
		memoryAmount := nbReq.MemoryAmount
		floatStr := fmt.Sprintf("%.1f", memoryAmount)
		s := nbReq.MemoryUnit
		limitMap["memory"] = fmt.Sprint(floatStr, *s) //todo
		requestMap["memory"] = fmt.Sprint(floatStr, *s)
	}
	if nbReq.ExtraResources != "" {
		gpuMap := make(map[string]string)
		err := json.Unmarshal([]byte(nbReq.ExtraResources), &gpuMap)
		if err != nil {
			return err
		}
		if v, ok := gpuMap["nvidia.com/gpu"]; ok {
			limitMap["nvidia.com/gpu"] = v
			requestMap["nvidia.com/gpu"] = v
		}
	}
	notebook.Spec.Template.Spec.Containers[0].Resources.Limits = limitMap
	notebook.Spec.Template.Spec.Containers[0].Resources.Requests = requestMap
	if nbReq.ImageName != "" {
		notebook.Spec.Template.Spec.Containers[0].Image = nbReq.ImageName
	}

	logger.Logger().Debugf("PatchNoteBookCRD notebook: %+v\n", notebook)
	if sparkService == nil {
		notebook.SparkSessionService = client.V1Service{}
		logger.Logger().Debugln("PatchNoteBookCRD SparkSessionService is nil")
	} else {
		notebook.SparkSessionService = *sparkService
		logger.Logger().Debugf("PatchNoteBookCRD SparkSessionService: %+v\n", *sparkService)
	}
	notebookBytes, _ := json.Marshal(notebook)
	notebookMap := make(map[string]interface{})
	err = json.Unmarshal(notebookBytes, &notebookMap)
	if err != nil {
		logger.Logger().Errorf("fail to unmarshal notebook to map, notebookBytes: %s, err: %v\n", string(notebookBytes), err)
		return err
	}
	logger.Logger().Debugf("PatchNoteBookCRD notebookMap: %+v\n", notebookMap)
	replacedBody, _, err := ncClient.k8sClient.CustomObjectsApi.ReplaceNamespacedCustomObject(context.TODO(), NBApiGroup,
		NBApiVersion, namespace, NBApiPlural, name, notebookMap)
	//replacedBody, _, err := ncClient.k8sClient.CustomObjectsApi.ReplaceNamespacedCustomObject(context.TODO(), NBApiGroup,
	//	NBApiVersion, namespace, NBApiPlural, name, &notebook)
	//ncClient.k8sClient.CustomObjectsApi.PatchNamespacedCustomObject(context.TODO(), NBApiGroup,
	//	NBApiVersion, namespace, NBApiPlural, name, nil)
	if err != nil {
		logger.Logger().Errorf("fail to update notebook CRD, namespace: %s, NBApiGroup: %s, "+
			"NBApiVersion: %s, NBApiPlural: %s, err: %v\n", namespace, NBApiGroup, NBApiVersion, NBApiPlural, err)
		return err
	}
	logger.Logger().Infof("PatchNoteBookCRD, replacedBody: %v\n", replacedBody)

	return nil
}

func (*NotebookControllerClient) GetNoteBookCRD(namespace, name string) (v1ResourceRequirements client.V1ResourceRequirements, res *http.Response, err error) {
	var (
		body    interface{}
		nbBytes []byte
	)
	body, res, err = ncClient.k8sClient.CustomObjectsApi.GetNamespacedCustomObject(context.TODO(), NBApiGroup,
		NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook CRD, namespace: %s, name: %s, "+
			"NBApiGroup: %s, NBApiVersion: %s, NBApiPlural: %s, err: %v\n", namespace, name,
			NBApiGroup, NBApiVersion, NBApiPlural, err)
		return
	}
	logger.Logger().Infof("PatchNoteBookCRD body: %v\n", body)
	notebook := models.K8sNotebook{}
	nbBytes, err = json.Marshal(body)
	if err != nil {
		logger.Logger().Errorf("fail to marshal crd body: %v, err: %v\n", body, err)
		return
	}
	err = json.Unmarshal(nbBytes, &notebook)
	if err != nil {
		logger.Logger().Errorf("fail to unmarshal crd body to notebook, bytes: %v, err: %v\n", body, err)
		return
	}

	// check resource
	//todo
	v1ResourceRequirements.Limits = notebook.Spec.Template.Spec.Containers[0].Resources.Limits
	v1ResourceRequirements.Requests = notebook.Spec.Template.Spec.Containers[0].Resources.Requests
	return
}

func (*NotebookControllerClient) GetNoteBookStatus(namespace, name string) (status string, res *http.Response, err error) {
	var (
		body interface{}
		nbBytes []byte
	)
	body, res, err = ncClient.k8sClient.CustomObjectsApi.GetNamespacedCustomObject(context.TODO(), NBApiGroup,
		NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook CRD, namespace: %s, name: %s, "+
			"NBApiGroup: %s, NBApiVersion: %s, NBApiPlural: %s, err: %v\n", namespace, name,
			NBApiGroup, NBApiVersion, NBApiPlural, err)
		return
	}
	logger.Logger().Infof("PatchNoteBookCRD body: %v\n", body)
	notebook := models.K8sNotebook{}
	nbBytes, err = json.Marshal(body)
	if err != nil {
		logger.Logger().Errorf("fail to marshal crd body: %v, err: %v\n", body, err)
		return
	}
	err = json.Unmarshal(nbBytes, &notebook)
	if err != nil {
		logger.Logger().Errorf("fail to unmarshal crd body to notebook, bytes: %v, err: %v\n", body, err)
		return
	}

	// check resource
	//todo
	status = "NotReady"
	if nil != notebook.Status.ContainerState.Running {
		status = "Ready"
	} else if nil != notebook.Status.ContainerState.Waiting {
		logger.Logger().Debugf("GetNoteBookStatus for state Waiting: %v", notebook.Status.ContainerState.Waiting.Message)
		status = "Waiting"
	} else if nil != notebook.Status.ContainerState.Terminated {
		logger.Logger().Debugf("GetNoteBookStatus for state terminater: %v", notebook.Status.ContainerState.Terminated.Message)
		status = "Terminated"

	} else {
		logger.Logger().Debugf("GetNoteBookStatus for state terminater: %v", notebook.Status.ContainerState)
		status = "Waiting"
	}
	return
}

func (*NotebookControllerClient) PatchYarnSettingConfigMap(notebook *models.PatchNotebookRequest,
	sparkSessionMap map[string]interface{}, name, namespace string) error { //todo
	nbName := name
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
		logrus.Error("unmarshal yarn json error" + err.Error())
		return err
	}
	err = json.Unmarshal([]byte(configmap.Data["linkismagic.json"]), &linkismagicConfigmap)
	if err != nil {
		logrus.Error("unmarshal linkismagic json error" + err.Error())
		return err
	}
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

	configJsonStr := string(configByte)
	lmConfigJsonStr := string(linkismagicByte)

	valueMap := make(map[string]interface{})
	valueMap["config.json"] = configJsonStr
	valueMap["linkismagic.json"] = lmConfigJsonStr
	if len(sparkSessionMap) > 0 {
		sparkSessionMapBytes, err := json.Marshal(sparkSessionMap)
		if err != nil {
			logrus.Error("marshal spark session map error" + err.Error())
			return err
		}
		valueMap["sparkSession.json"] = string(sparkSessionMapBytes)
	}
	//opr_map := []map[string]interface{}{
	//	{"op": "replace",
	//		"path": "/data",
	//		"value": map[string]interface{}{
	//			"config.json":      configJsonStr,
	//			"linkismagic.json": lmConfigJsonStr,
	//		}},
	//}

	opr_map := []map[string]interface{}{
		{"op": "replace",
			"path":  "/data",
			"value": valueMap, //todo 测试sparkSession未更新时是否被替换
		},
	}
	_, _, err = ncClient.k8sClient.CoreV1Api.PatchNamespacedConfigMap(context.Background(),
		nbName+"-"+"yarn-resource-setting", namespace, opr_map, nil)
	if err != nil {
		logrus.Error("PatchNamespacedConfigMap" + err.Error())
		return err
	}
	if err == nil {
		logger.Logger().Info("success")
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

func deleteK8sSparkSessionService(namespace, name string) error {
	body := client.V1DeleteOptions{
		ApiVersion:         "",
		GracePeriodSeconds: 0,
		Kind:               "",
		OrphanDependents:   false,
		Preconditions:      nil,
		PropagationPolicy:  "",
	}
	_, _, err := ncClient.k8sClient.CoreV1Api.DeleteNamespacedService(context.TODO(), name,
		namespace, body, nil)
	return err
}

func (*NotebookControllerClient) GetNBConfigMaps(namespace string) (*map[string]interface{}, error) {
	var list client.V1ConfigMapList
	var response *http.Response
	var err error
	configs := map[string]interface{}{}
	if "" == namespace || "null" == namespace {
		list, response, err = ncClient.k8sClient.CoreV1Api.ListConfigMapForAllNamespaces(context.Background(), nil)
	} else {
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

	logger.Logger().Info("configmap length:" + string(len(list.Items)))
	for _, configmap := range list.Items {
		//if configmap.Metadata.Name != "yarn-resource-setting" {
		//	continue
		//}
		if !stringsUtil.HasSuffix(configmap.Metadata.Name, "yarn-resource-setting") {
			continue
		}
		var yarnConfigmap map[string]interface{}
		err = json.Unmarshal([]byte(configmap.Data["config.json"]), &yarnConfigmap)
		//logr.Info("configmap keys:" + configmap.Metadata.Namespace + "-" + configmap.Metadata.Name  )
		if err != nil {
			logger.Logger().Error("Json Parse Error: ", configmap.Metadata.Namespace+"-"+configmap.Metadata.Name, ",", err.Error())
		}
		configs[configmap.Metadata.Namespace+"-"+configmap.Metadata.Name] = yarnConfigmap
	}
	return &configs, err
}

func GetNotebook(namespace string, name string) (interface{}, *http.Response, error) {
	notebook, res, err := ncClient.k8sClient.CustomObjectsApi.
		GetNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Error("GetNotebook error, %v", err.Error())
		return nil, res, err
	}
	if res == nil {
		logger.Logger().Error("empty response")
		return nil, nil, nil
	}
	return notebook, res, nil
}

func (*NotebookControllerClient) CheckNotebookExisted(namespace string, name string) (bool, error) {
	nb, res, err := ncClient.k8sClient.CustomObjectsApi.
		GetNamespacedCustomObject(context.Background(), NBApiGroup, NBApiVersion, namespace, NBApiPlural, name)
	if err != nil {
		logger.Logger().Error("GetNotebook error res, %v", res.StatusCode)
		//Notebook Not Existed if status code is 404
		if res.StatusCode == 404 {
			logger.Logger().Error("GetNotebook error res, %v", res.StatusCode)
			return false, nil
		}
		logger.Logger().Error("GetNotebook error, %v", err.Error())
		return false, err
	}
	//Notebook is Existed
	if nb != nil {
		return true, nil
	}
	if res == nil {
		logger.Logger().Error("empty response")
		return false, errors.New("kubernetes client ressponse is empty.")
	}
	return false, nil
}

func GetPod(namespace, name string) (client.V1Pod, *http.Response, error) {
	return ncClient.k8sClient.CoreV1Api.ReadNamespacedPod(context.TODO(), name, namespace, nil)
}

func (*NotebookControllerClient) GetPodStatus(namespace, name string) (rets string, er error) {
	name = name + "-0"

	temp := make(map[string]interface{})
	temp["pretty"] = nil
	nb, res, err := ncClient.k8sClient.CoreV1Api.ReadNamespacedPodStatus(context.Background(), name, namespace, temp)
	if err != nil {
		logger.Logger().Error("GetNotebook error res, %v", res.StatusCode)
		//Notebook Not Existed if status code is 404
		if res.StatusCode == 404 {
			logger.Logger().Error("GetNotebook error res, %v", res.StatusCode)
			return
		}
		logger.Logger().Error("GetNotebook error, %v", err.Error())
		return
	}
	logger.Logger().Infof("GetPodStatus nb:%v", nb)

	rets = nb.Status.Phase
	logger.Logger().Infof("GetPodStatus rets:%v", rets)
	return
}
