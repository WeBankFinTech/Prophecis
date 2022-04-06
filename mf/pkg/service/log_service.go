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

package service

import (
	"fmt"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/dao/es"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/models"
	"mlss-mf/pkg/restapi/operations/container"
	"reflect"
	"sort"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespacedServiceContainerLog(params container.GetNamespacedServiceContainerLogParams, user string, isSA bool) (
	models.GetServiceContainerLogResponse, error) {
	rst := models.GetServiceContainerLogResponse{}
	if !isSA {
		serverInfo, err := serviceDao.GetServiceByNSAndName(params.Namespace, params.ServiceName)
		if err != nil {
			logger.Logger().Error("GetNamespacedServiceContainerLog error", err.Error())
			return rst, err
		}
		err = PermissionCheck(user, serverInfo.UserID, nil, isSA)
		if err != nil {
			logger.Logger().Errorf("Permission Check Error:" + err.Error())
			return rst, err
		}
	}
	if err := checkGetNamespacedNotebookLogParams(&params); err != nil {
		logger.Logger().Errorf("fail to check param, error: %s", err.Error())
		return rst, err
	}

	// get container information from SeldonDeployments
	podName, containerName, isExist, err := checkContainerFromSeldonDeployments(params.Namespace, params.ServiceName, params.ContainerName)
	if err != nil {
		logger.Logger().Errorf("fail to check container info from seldDeployment, namesapce: %s, "+
			"service name: %s, container name: %s", params.Namespace, params.ServiceName, params.ContainerName)
		return rst, err
	}
	logger.Logger().Debugf("GetNamespacedServiceContainerLog pod name: %s, container name: %s\n", podName, containerName)
	if !isExist {
		logger.Logger().Infof("service not exist in seldDeployment, namesapce: %s, "+
			"service name: %s, container name: %s", params.Namespace, params.ServiceName, params.ContainerName)
		return rst, nil
	}
	if containerName == "" {
		logger.Logger().Infof("fail to get container info from seldDeployment, namesapce: %s, "+
			"service name: %s,", params.Namespace, params.ServiceName)
		return rst, nil
	}
	if podName == "" {
		logger.Logger().Infof("fail to get pod info from seldDeployment, namesapce: %s, "+
			"service name: %s,", params.Namespace, params.ServiceName)
		return rst, nil
	}

	esClient := es.GetESClient()
	if esClient == nil {
		logger.Logger().Errorf("fail to get es client, es client: %v", esClient)
		err := fmt.Errorf("fail to get es client, es client: %v", esClient)
		return rst, err
	}

	logParam := es.Log{
		EsIndex:       es.K8S_LOG_CONTAINER_ES_INDEX,
		EsType:        es.K8S_LOG_CONTAINER_ES_TYPE,
		PodName:       podName,
		ContainerName: containerName,
		Ns:            params.Namespace,
		From:          int((*params.CurrentPage - 1) * (*params.PageSize)),
		Size:          int(*params.PageSize),
		Asc:           *params.Asc,
	}
	r, err := es.LogList(logParam)
	if err != nil {
		logger.Logger().Errorf("fail to get service container log information from es, err: %v", err)
		return rst, err
	}

	rst, err = generateLogList(r)
	if err != nil {
		logger.Logger().Errorf("fail to format log information from map, err: %v, map: %+v", err, r)
		return rst, err
	}
	return rst, nil
}

func checkGetNamespacedNotebookLogParams(param *container.GetNamespacedServiceContainerLogParams) error {
	if param == nil || param.Namespace == "" || param.ServiceName == "" || param.ContainerName == "" {
		logger.Logger().Errorf("params can't be empty string or nil, param: %v", *param)
		return fmt.Errorf("params can't be empty string or nil, param: %v", *param)
	}
	if param.CurrentPage == nil || *param.CurrentPage == 0 {
		var currentPage int64 = 1
		param.CurrentPage = &currentPage
	}
	if param.PageSize == nil || *param.PageSize == 0 {
		var pageSize int64 = 10
		param.CurrentPage = &pageSize
	}

	asc := false
	if *(param.Asc) {
		asc = true
	}
	param.Asc = &asc

	return nil
}

func generateLogList(r map[string]interface{}) (models.GetServiceContainerLogResponse, error) {
	rst := models.GetServiceContainerLogResponse{}
	if len(r) == 0 || r["hits"] == nil {
		logger.Logger().Errorf("es hits is nil, r: %+v", r)
		err := fmt.Errorf("es hits is nil, r: %+v", r)
		return rst, err
	}

	switch r["hits"].(type) {
	case map[string]interface{}:
		// logger.Logger().Info("hits type is map[string]interface{}")
		if len(r["hits"].(map[string]interface{})) == 0 {
			logger.Logger().Info("es hits type is map[string]interface{}, but map is empty")
			return rst, nil
		}
	default:
		logger.Logger().Errorf("hits type is undefined, type by  reflect is %v", reflect.TypeOf(r["hits"]).Kind().String())
		err := fmt.Errorf("hits type is undefined, type by  reflect is %v", reflect.TypeOf(r["hits"]).Kind().String())
		return rst, err
	}

	serviceContainerLogList := make([]*models.ServiceContainerLog, 0)
	var total int64 = 0
	if hitsTotalInter, ok := r["hits"].(map[string]interface{})["total"]; ok {
		var hitsTotal int64 = 0
		switch hitsTotalVal := hitsTotalInter.(type) {
		case int:
			hitsTotal = int64(hitsTotalVal)
		case int64:
			hitsTotal = hitsTotalVal
		case float32:
			hitsTotal = int64(hitsTotalVal)
		case float64:
			hitsTotal = int64(hitsTotalVal)
		case string:
			hitsTotalInt64, err := strconv.ParseInt(hitsTotalVal, 10, 64)
			if err != nil {
				logger.Logger().Infof("hits total type is string, but parse error: %v", err)
			}
			hitsTotal = hitsTotalInt64
		case map[string]interface{}:
			logger.Logger().Debugf("test-lk hits total type is map[string]interface{}, hitsTotal: %+v", hitsTotal)
			if valueInter, exist := hitsTotalVal["value"]; exist {
				if totalFloat64, existFloat64 := valueInter.(float64); existFloat64 {
					hitsTotal = int64(totalFloat64)
				}
			}
		default:
			logger.Logger().Infof("hits total type undefined, %v", reflect.TypeOf(hitsTotalInter).Kind().String())
			// total = hitsTotal
		}
		total = hitsTotal
	}
	if _, ok := r["hits"].(map[string]interface{})["hits"]; ok {
		if hits, ok2 := r["hits"].(map[string]interface{})["hits"].([]interface{}); ok2 {
			for _, hit := range hits {
				var log interface{} = ""
				switch logV := hit.(type) {
				case map[string]interface{}:
					if source, ok := logV["_source"]; ok {
						sourceSub, sourceOK := source.(map[string]interface{})
						if !sourceOK && len(sourceSub) == 0 {
							continue
						}

						if v, ok2 := sourceSub["log"]; ok2 {
							log = v
						}
					}
					serviceContainerLog := models.ServiceContainerLog{
						Log: log.(string),
					}
					serviceContainerLogList = append(serviceContainerLogList, &serviceContainerLog)
				default:
					logger.Logger().Infof("hit type undefined, hit type: %v", reflect.TypeOf(logV).Kind().String())
					continue
				}
			}
		}
	}

	rst.LogList = serviceContainerLogList
	rst.Total = &total

	return rst, nil
}

func checkContainerFromSeldonDeployments(ns, name, containerName string) (string, string, bool, error) {
	seldonDeployment, err := seldonClient.MachinelearningV1().SeldonDeployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Errorf("fail to get seldonDeployment, namespace: %s, name: %s, error: %v\n",
			ns, name, err)
		return "", "", false, err
	}

	logger.Logger().Infof("test-lk seldonDeployment.Status.DeploymentStatus: %+v, length: %d",
		seldonDeployment.Status.DeploymentStatus, len(seldonDeployment.Status.DeploymentStatus))
	for deploymentName, _ := range seldonDeployment.Status.DeploymentStatus {
		logger.Logger().Debugf("deployment name: %s\n", deploymentName)
		deployment, err := common.Clientset.AppsV1().Deployments(ns).Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			logger.Logger().Infof("fail to get deployment, namespace: %s, name: %s, error: %v\n",
				ns, deploymentName, err)
			continue
		}

		label := ""
		for key, value := range deployment.Spec.Selector.MatchLabels {
			label += fmt.Sprintf("%s=%s,", key, value)
		}
		if len(label) > 0 {
			label = string([]byte(label)[:len(label)-1])
		}
		logger.Logger().Debugf("rs label: %s", label)
		rsList, err := common.Clientset.AppsV1().ReplicaSets(ns).List(metav1.ListOptions{LabelSelector: label})
		if err != nil {
			logger.Logger().Infof("ListContainer: fail to get replicaSet list, namespace: %s"+
				", LabelSelector: %s, error: %v\n", ns, label, err)
			continue
		}

		if rsList == nil || rsList.Items == nil || len(rsList.Items) == 0 {
			continue
		}
		sort.Sort(ReplicaSetItem(rsList.Items))

		rsName := rsList.Items[0].Name
		logger.Logger().Debugf("rs name: %s", rsName)
		rs, err := common.Clientset.AppsV1().ReplicaSets(ns).Get(rsName, metav1.GetOptions{})
		if err != nil {
			logger.Logger().Infof("fail to get replicaSet, namespace: %s"+
				", name: %s, error: %v\n", ns, rsName, err)
			continue
		}

		podTemplateHash, ok := rs.Labels["pod-template-hash"]
		if !ok {
			logger.Logger().Infof("fail to get pod-template-hash, labels: %+v\n", rs.Labels)
			continue
		}

		podList, err := common.Clientset.CoreV1().Pods(ns).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("pod-template-hash=%s", podTemplateHash),
		})
		if err != nil {
			logger.Logger().Infof("fail to get pod list by namespace %q and label pod-template-hash %q\n",
				ns, podTemplateHash)
			continue
		}

		for _, pod := range podList.Items {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.Name == containerName {
					if ok := strings.HasPrefix(containerStatus.ContainerID, "docker://"); ok {
						containerID := strings.TrimPrefix(containerStatus.ContainerID, "docker://")
						containerNameWithId := fmt.Sprintf("%s-%s", containerName, containerID)
						return pod.Name, containerNameWithId, true, nil
					}
					return pod.Name, containerName, true, nil
				}
			}
		}
	}

	logger.Logger().Errorf("fail to get container info from seldonDeployment, namespace: %s, service name: %s, "+
		"container name: %s", ns, name, containerName)
	return "", "", false, err
}
