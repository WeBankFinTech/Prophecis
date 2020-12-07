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
package controller

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/resources"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func SetNamespaceRQ(params resources.SetNamespaceRQParams) middleware.Responder {
	checkForRQErr := common.CheckForRQParams(*params.ResourcesQuota)

	if nil != checkForRQErr {
		return ResponderFunc(http.StatusBadRequest, "failed to list namespace from", checkForRQErr.Error())
	}

	rqRequest := common.FromSetNamespaceRQRequest(params)

	setNamespaceRQErr := service.SetNamespaceRQ(rqRequest)

	if nil != setNamespaceRQErr {
		return ResponderFunc(http.StatusBadRequest, "failed to update resourcesQuota for namespace", setNamespaceRQErr.Error())
	}

	return GetResult([]byte{}, nil)
}

func GetLabelsOfNode(params resources.GetLabelsOfNodeParams) middleware.Responder {
	page := *params.Page
	size := *params.Size
	logger.Logger().Debugf("v1/resources/labels GetLabelsOfNode params page: %v, size: %v", page, size)

	labelsResponses, getErr := service.GetLabelsOfNode(page, size)

	if nil != getErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to get labels of nodes", getErr.Error())
	}

	marshal, marshalErr := json.Marshal(labelsResponses)

	return GetResult(marshal, marshalErr)
}

func AddLabels(params resources.AddLabelsParams) middleware.Responder {
	checkErr := common.CheckForAddLabels(params)

	if nil != checkErr {
		return ResponderFunc(http.StatusBadRequest, "failed to add labels for node", checkErr.Error())
	}

	newLabels, addErr := service.AddLabels(*params.LabelsRequest)

	if nil != addErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to add labels for node", addErr.Error())
	}

	marshal, marshalErr := json.Marshal(newLabels)

	return GetResult(marshal, marshalErr)
}

func RemoveNodeLabel(params resources.RemoveNodeLabelParams) middleware.Responder {
	label := params.Label
	nodeName := params.NodeName

	newLabels, removeErr := service.RemoveNodeLabel(label, nodeName)

	if nil != removeErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to remove labels for node", removeErr.Error())
	}

	marshal, marshalErr := json.Marshal(newLabels)

	return GetResult(marshal, marshalErr)
}

func UpdateLabels(params resources.UpdateLabelsParams) middleware.Responder {

	checkErr := common.CheckForLabelsRequest(params)

	if nil != checkErr {
		return ResponderFunc(http.StatusBadRequest, "failed to update labels for node", checkErr.Error())
	}

	newLabels, updateErr := service.UpdateLabels(*params.ResourcesQuota)

	if nil != updateErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to update labels for node", updateErr.Error())
	}

	marshal, marshalErr := json.Marshal(newLabels)

	return GetResult(marshal, marshalErr)
}

func GetNodeByName(params resources.GetNodeByNameParams) middleware.Responder {

	nodeByName, getErr := service.GetNodeByName(params.NodeName)

	if nil != getErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to get node by name", getErr.Error())
	}

	marshal, marshalErr := json.Marshal(nodeByName)

	return GetResult(marshal, marshalErr)
}
