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
	"k8s.io/api/core/v1"
	"math"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/manager"
	"mlss-controlcenter-go/pkg/models"
	"strings"
)

func SetNamespaceRQ(rqRequest models.ResourcesQuota) error {
	namespace := rqRequest.Namespace
	cpu := rqRequest.CPU
	gpu := rqRequest.Gpu
	memory := rqRequest.Memory

	logger.Logger().Debugf("update namespace rq with namespace: %v cpu: %v, gpu: %v, memory: %v", namespace, cpu, gpu, memory)
	kac := manager.GetK8sApiClient()
	nodes, getNodeErr := kac.GetNodesByNamespace(namespace)
	if nil != getNodeErr {
		return getNodeErr
	}
	logger.Logger().Debugf("get nodes by namespace rq with namespace: %v nodes: %v", namespace, len(nodes))

	checkForReqErr := kac.CheckForRQRequest(nodes, cpu, gpu, memory["memoryAmount"])
	if nil != checkForReqErr {
		return checkForReqErr
	}

	memReq := memory["memoryAmount"] + memory["memoryUnit"]

	updateRQErr := kac.UpdateRQ(namespace, cpu, gpu, memReq)

	if updateRQErr != nil {
		return updateRQErr
	}

	return nil
}

func GetLabelsOfNode(page int64, size int64) (*models.PageLabelsResponseList, error) {
	var labels []*models.LabelsResponse

	kac := manager.GetK8sApiClient()

	allNodes, getNodeErr := kac.GetAllNodes()

	if nil != getNodeErr {
		return nil, getNodeErr
	}

	if page > 0 && size > 0 {
		offSet := common.GetOffSet(page, size)
		for index, node := range allNodes.Items {
			if int64(index) >= offSet && int64(index) < offSet+size {
				var label = &models.LabelsResponse{
					IPAddress: node.Status.Addresses[0].Address,
					Labels:    node.ObjectMeta.Labels,
				}
				labels = append(labels, label)
			}
		}
	} else {
		for _, node := range allNodes.Items {
			var label = &models.LabelsResponse{
				IPAddress: node.Status.Addresses[0].Address,
				Labels:    node.ObjectMeta.Labels,
			}
			labels = append(labels, label)
		}
	}

	total := int64(len(allNodes.Items))
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	var pageLabelsRes = &models.PageLabelsResponseList{
		Models:     labels,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}

	return pageLabelsRes, nil
}

func AddLabels(labels models.LabelsRequest) (map[string]string, error) {
	var ipAddress = labels.IPAddress
	err, nodeName := GetNodeNameByIp(ipAddress)
	if nil != err {
		return nil, err
	}

	//isMasterNode, isMasterNodeErr := IsMasterNode(nodeName)
	//if nil != isMasterNodeErr {
	//	return nil, isMasterNodeErr
	//}

	//if isMasterNode {
	//	logger.Logger().Errorf("failed to add labels for node: %v", nodeName)
	//	return nil, errors.New("the labels of master node is not allowed to changed")
	//}

	var namespace = labels.Namespace
	var lbBussType = labels.LbBusType
	var lbGpuModel = labels.LbGpuModel

	kac := manager.GetK8sApiClient()

	nsFromK8s, namespaceErr := kac.GetNSFromK8s(namespace)
	if nil != namespaceErr {
		return nil, namespaceErr
	}

	annotations := nsFromK8s.ObjectMeta.Annotations
	newLabels := map[string]string{}
	nodeSelector := annotations[constants.NodeSelectorKey]

	split := strings.Split(nodeSelector, ",")
	logger.Logger().Debugf("")
	for _, nsLabel := range split {
		logger.Logger().Debugf("nsLabel: %v", nsLabel)
		if len(strings.Split(nsLabel, "=")) > 1 {
			newLabels[strings.Split(nsLabel, "=")[0]] = strings.Split(nsLabel, "=")[1]
		}
	}
	newLabels[constants.LB_BUS_TYPE] = lbBussType
	newLabels[constants.LB_GPU_MODEL] = lbGpuModel

	replaceErr := kac.ReplaceNode(nodeName, newLabels)

	if nil != replaceErr {
		return nil, replaceErr
	}

	return newLabels, nil
}

func RemoveNodeLabel(label string, nodeName string) (map[string]string, error) {
	//isMasterNode, isMasterNodeErr := IsMasterNode(nodeName)
	//if nil != isMasterNodeErr {
	//	return nil, isMasterNodeErr
	//}
	//
	//if isMasterNode {
	//	logger.Logger().Errorf("failed to update labels for node: %v", nodeName)
	//	return nil, errors.New("the labels of master node is not allowed to changed")
	//
	//}

	split := strings.Split(label, ",")

	return manager.GetK8sApiClient().RemoveNodeLabel(nodeName, split)

}

func UpdateLabels(labels models.LabelsRequest) (map[string]string, error) {
	var ipAddress = labels.IPAddress
	err, nodeName := GetNodeNameByIp(ipAddress)
	if nil != err {
		return nil, err
	}

	//isMasterNode, isMasterNodeErr := IsMasterNode(nodeName)
	//if nil != isMasterNodeErr {
	//	return nil, isMasterNodeErr
	//}
	//
	//if isMasterNode {
	//	logger.Logger().Errorf("failed to update labels for node: %v", nodeName)
	//	return nil, errors.New("the labels of master node is not allowed to changed")
	//}

	var namespace = labels.Namespace
	var lbBussType = labels.LbBusType
	var lbGpuModel = labels.LbGpuModel

	kac := manager.GetK8sApiClient()

	nsFromK8s, namespaceErr := kac.GetNSFromK8s(namespace)
	if nil != namespaceErr {
		return nil, namespaceErr
	}

	annotations := nsFromK8s.ObjectMeta.Annotations
	newLabels := map[string]string{}
	nodeSelector := annotations[constants.NodeSelectorKey]

	split := strings.Split(nodeSelector, ",")
	logger.Logger().Debugf("")
	for _, nsLabel := range split {
		logger.Logger().Debugf("nsLabel: %v", nsLabel)
		if len(strings.Split(nsLabel, "=")) > 1 {
			newLabels[strings.Split(nsLabel, "=")[0]] = strings.Split(nsLabel, "=")[1]
		}
	}
	newLabels[constants.LB_BUS_TYPE] = lbBussType
	newLabels[constants.LB_GPU_MODEL] = lbGpuModel

	replaceErr := kac.ReplaceNode(nodeName, newLabels)

	if nil != replaceErr {
		return nil, replaceErr
	}

	return newLabels, nil
}

func GetNodeNameByIp(ipAddress string) (error, string) {
	var nodeName string

	kac := manager.GetK8sApiClient()
	nodeList, nodeErr := kac.GetAllNodes()
	if nil != nodeErr {
		return nodeErr, ""
	}
Loop:
	for _, node := range nodeList.Items {
		addresses := node.Status.Addresses
		for _, nodeAddress := range addresses {
			if ipAddress == nodeAddress.Address {
				nodeName = node.ObjectMeta.Name
				break Loop
			}
		}
	}
	return nil, nodeName
}

func IsMasterNode(nodeName string) (bool, error) {
	labels, getErr := GetLabelsByNodeName(nodeName)

	if nil != getErr {
		return false, getErr
	}

	if _, ok := labels["node-role.kubernetes.io/master"]; nil != labels && ok {
		return true, nil
	}

	return false, nil
}

func GetLabelsByNodeName(nodeName string) (map[string]string, error) {
	kac := manager.GetK8sApiClient()
	node, getErr := kac.GetNodeByName(nodeName)

	if nil != getErr {
		return nil, getErr
	}

	return node.ObjectMeta.Labels, nil
}

func GetNodeByName(nodeName string) (*v1.Node, error) {
	return manager.GetK8sApiClient().GetNodeByName(nodeName)
}
