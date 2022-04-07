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
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/manager"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/namespaces"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
	"strings"
)

func GetAllNamespaces(params namespaces.GetAllNamespacesParams) middleware.Responder {
	page := *params.Page
	size := *params.Size
	err := common.CheckPageParams(page, size)
	if nil != err {
		return ResponderFunc(http.StatusForbidden, "failed to check for page and size", err.Error())
	}
	logger.Logger().Debugf("v1/users GetAllUsers params page: %v, size: %v", page, size)

	namespaceResList, getErr := service.GetAllNamespaces(page, size)

	if nil != getErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to list namespace from", getErr.Error())
	}

	marshal, marshalErr := json.Marshal(namespaceResList)

	return GetResult(marshal, marshalErr)
}

func AddNamespace(params namespaces.AddNamespaceParams) middleware.Responder {
	checkForNamespaceErrMsg := common.CheckForNamespace(*params.Namespace)

	if checkForNamespaceErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to add namespace", checkForNamespaceErrMsg.Error())
	}

	namespaceRequest := common.FromAddNamespaceRequest(params)

	logger.Logger().Debugf("AddNamespace namespaceRequest: %v", namespaceRequest)

	namespace := common.FromNamespaceRequestToNamespace(namespaceRequest)

	logger.Logger().Debugf("AddNamespace namespace: %v", namespace)

	namespaceByName, err := service.GetNamespaceByName(namespace.Namespace)
	logger.Logger().Debugf("AddNamespace namespaceByName: %v", namespaceByName)
	if err != nil {
		addErr := service.AddNamespace(namespace, namespaceRequest)
		if nil != addErr {
			return ResponderFunc(http.StatusBadRequest, "failed to add namespace", addErr.Error())
		}
	} else {
		if namespace.Namespace == namespaceByName.Namespace {
			return ResponderFunc(http.StatusBadRequest, "failed to add namespace", "namespace already in db")
		}
	}
	repoNamespace, err := service.GetNamespaceByName(namespace.Namespace)
	if err != nil {
		logger.Logger().Error("Get namespace by name err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(repoNamespace)
	return GetResult(marshal, marshalErr)
}

func UpdateNamespace(params namespaces.UpdateNamespaceParams) middleware.Responder {
	checkForNamespaceErrMsg := common.CheckForNamespace(*params.Namespace)

	if checkForNamespaceErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update namespace", checkForNamespaceErrMsg.Error())
	}

	if params.Namespace.ID == 0 {
		return ResponderFunc(http.StatusBadRequest, "failed to update namespace", "namespace id is invalid")
	}

	namespaceRequest := common.FromUpdateNamespaceRequest(params)

	logger.Logger().Debugf("UpdateNamespace namespaceRequest: %v", namespaceRequest)

	namespace := common.FromNamespaceRequestToNamespace(namespaceRequest)

	logger.Logger().Debugf("AddNamespace namespace: %v", namespace)

	namespaceByID, err := service.GetNamespaceByID(namespace.ID)
	if err != nil {
		logger.Logger().Error("Get namespace by id err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if namespace.ID != namespaceByID.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update namespace", "namespace not exist in db")
	}

	if namespace.Namespace != namespaceByID.Namespace {
		return ResponderFunc(http.StatusBadRequest, "failed to update namespace", "The namespace id and namespace do not match")
	}

	kac := manager.GetK8sApiClient()

	nsFromK8s, getNSErr := kac.GetNSFromK8s(namespace.Namespace)

	if nil != getNSErr {
		return ResponderFunc(http.StatusBadRequest, "failed to get namespace from k8s", getNSErr.Error())
	}

	if nil == nsFromK8s {
		return ResponderFunc(http.StatusBadRequest, "failed to get namespace from k8s", "namespace not exist in k8s")
	}

	repoNamespace, err := service.UpdateNamespace(namespace)
	if err != nil {
		logger.Logger().Error("Update namespace err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(repoNamespace)

	annotationsRequest := namespaceRequest.Annotations

	logger.Logger().Debugf("updateNamespace with annotations : %v", annotationsRequest)

	//annotation named "scheduler.alpha.kubernetes.io/node-selector" fetched from k8s
	var nodeSelector string
	//updated annotation named "scheduler.alpha.kubernetes.io/node-selector"
	//var selector bytes.Buffer

	//updating annotation named "scheduler.alpha.kubernetes.io/node-selector"
	//"lb-idc", "lb-department", "lb-app",  "lb-network" keys included by annotation named "scheduler.alpha.kubernetes.io/node-selector", will not be updated.
	var updatingNodeSelectorAnnoMap = make(map[string]string)
	if nil != nsFromK8s.ObjectMeta.Annotations {
		nodeSelector = nsFromK8s.ObjectMeta.Annotations[constants.NodeSelectorKey]
		if "" != nodeSelector {
			customSelectorArr := strings.Split(nodeSelector, ",")
			for _, customSelector := range customSelectorArr {
				split := strings.Split(customSelector, "=")
				cusSelKey := split[0]
				if "lb-idc" == cusSelKey || "lb-department" == cusSelKey || "lb-app" == cusSelKey || "lb-network" == cusSelKey {
					//if index == len(customSelectorArr)-1 {
					//	selector.WriteString(customSelector)
					//} else {
					//	selector.WriteString(customSelector)
					//	selector.WriteString(",")
					//}
					updatingNodeSelectorAnnoMap[split[0]] = split[1]
				}
			}
		}
	}

	if nil != annotationsRequest {
		//if "" != annotations["lb-bus-type"] && "" != annotations["lb-gpu-model"] {
		//	selector.WriteString(",lb-bus-type=")
		//	selector.WriteString(annotations["lb-bus-type"])
		//	selector.WriteString(",lb-gpu-model=")
		//	selector.WriteString(annotations["lb-gpu-model"])
		//	delete(annotations, "lb-bus-type")
		//	delete(annotations, "lb-gpu-model")
		//} else if "" != annotations["lb-bus-type"] {
		//	selector.WriteString(",lb-bus-type=")
		//	selector.WriteString(annotations["lb-bus-type"])
		//	delete(annotations, "lb-bus-type")
		//} else if "" != annotations["lb-gpu-model"] {
		//	selector.WriteString(",lb-gpu-model=")
		//	selector.WriteString(annotations["lb-gpu-model"])
		//	delete(annotations, "lb-gpu-model")
		//}
		if "" != annotationsRequest["lb-bus-type"] {
			updatingNodeSelectorAnnoMap["lb-bus-type"] = annotationsRequest["lb-bus-type"]
		}

		if "" != annotationsRequest["lb-gpu-model"] {
			updatingNodeSelectorAnnoMap["lb-gpu-model"] = annotationsRequest["lb-gpu-model"]
		}

	} else {
		annotationsRequest = map[string]string{}
	}

	//annotations[constants.NodeSelectorKey] = selector.String()
	var annoList []string
	for key, value := range updatingNodeSelectorAnnoMap {
		kvPair := fmt.Sprintf("%s=%s", key, value)
		annoList = append(annoList, kvPair)
	}
	updatingNodeSelectorAnnoString := strings.Join(annoList, ",")
	annotationsRequest[constants.NodeSelectorKey] = updatingNodeSelectorAnnoString

	replaceNamespaceErr := kac.ReplaceNamespace(namespace.Namespace, annotationsRequest)

	if nil != replaceNamespaceErr {
		return ResponderFunc(http.StatusBadRequest, "failed to update namespace to k8s", replaceNamespaceErr.Error())
	}

	return GetResult(marshal, marshalErr)
}

func GetNamespaceByName(params namespaces.GetNamespaceByNameParams) middleware.Responder {
	namespace := params.Namespace

	repoNamespace, err := service.GetNamespaceByName(namespace)
	if err != nil {
		logger.Logger().Error("Get namespace by name err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if namespace != repoNamespace.Namespace {
		return ResponderFunc(http.StatusBadRequest, "failed to get namespace", "namespace not exist in db")
	}

	kac := manager.GetK8sApiClient()

	nsFromK8s, err := kac.GetNSFromK8s(namespace)

	if nil != err {
		return ResponderFunc(http.StatusInternalServerError, "failed to get namespace", err.Error())
	}

	marshal, marshalErr := json.Marshal(nsFromK8s)

	return GetResult(marshal, marshalErr)
}

func DeleteNamespaceByName(params namespaces.DeleteNamespaceByNameParams) middleware.Responder {
	name := params.Namespace

	namespaceByName, err := service.GetNamespaceByName(name)
	if err != nil {
		logger.Logger().Error("Get namespace by name err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if name != namespaceByName.Namespace {
		return ResponderFunc(http.StatusBadRequest, "failed to delete namespace from db", "namespace is not exist in db")
	}

	repoNamespace, err := service.DeleteNamespaceByID(namespaceByName.ID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to DeleteNamespaceByID:%v", err.Error())
	}

	kac := manager.GetK8sApiClient()

	nsFromK8s, getNSErr := kac.GetNSFromK8s(name)

	if nil != getNSErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to get namespace from k8s", getNSErr.Error())
	}

	if nil != nsFromK8s {
		deleteNSErr := kac.DeleteNSFromK8s(name)
		if nil != deleteNSErr {
			return ResponderFunc(http.StatusInternalServerError, "failed to delete namespace from k8s", deleteNSErr.Error())
		}
	}

	marshal, marshalErr := json.Marshal(repoNamespace)

	return GetResult(marshal, marshalErr)
}

func GetMyNamespace(params namespaces.GetMyNamespaceParams) middleware.Responder {
	sessionUser := GetSessionByContext(params.HTTPRequest)

	logger.Logger().Debugf("getMyNamespace by sessionUser: %v", sessionUser)

	if sessionUser == nil || "" == sessionUser.UserName {
		logger.Logger().Error("SessionUser is nil or sessionUser.UserName's length <= 0")
		logger.Logger().Info("sessionUser == nil：", sessionUser == nil)
		if sessionUser != nil {
			logger.Logger().Info("sessionUser.UserName <= 0：", len(sessionUser.UserName) <= 0)
		}
		return FailedToGetUserByToken()
	}

	var myNamespace = common.NewStringSet()

	var groupNamespaces []models.GroupNamespace

	if IsSuperAdmin(sessionUser) {
		groupNamespaces, err := service.GetAllNamespace()
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetAllNamespace", err.Error())
		}
		myNamespace = common.GetMyNamespace(groupNamespaces, myNamespace)
		logger.Logger().Debugf("myNamespace for SA: %v", myNamespace)
	} else {
		username := sessionUser.UserName
		user, err := service.GetUserByUserName(username)
		if err != nil {
			return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
		}
		groupIds, err := service.GetGroupIdListByUserId(user.ID)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserId", err.Error())
		}
		logger.Logger().Debugf("myNamespace for groupIds: %v", groupIds)
		if nil != groupIds && len(groupIds) > 0 {
			groupNamespaces, err = service.GetAllGroupNamespaceByGroupIds(groupIds)
			if err != nil {
				return ResponderFunc(http.StatusInternalServerError, "failed to GetAllGroupNamespaceByGroupIds", err.Error())
			}
			logger.Logger().Debugf("myNamespace for GU groupNamespaces: %v, length: %v", groupNamespaces, len(groupNamespaces))
			myNamespace = common.GetMyNamespace(groupNamespaces, myNamespace)
			logger.Logger().Debugf("myNamespace for GA GU: %v", myNamespace)
		}
	}

	marshal, marshalErr := json.Marshal(myNamespace.List())

	return GetResult(marshal, marshalErr)
}

func ListNamespaceByRoleAndUserName(params namespaces.ListNamespaceByRoleNameAndUserNameParams) middleware.Responder {
	roleName := params.RoleName
	userName := params.UserName
	resp, err := service.ListNamespaceByRoleNameAndUserName(roleName, userName)
	if nil != err {
		logger.Logger().Errorf("fail to list namespace by role name and user name, "+
			"roleName: %s, userName:%s, err: %v\n", roleName, userName, err)
		return ResponderFunc(http.StatusInternalServerError, "failed to list namespace", err.Error())
	}

	marshal, marshalErr := json.Marshal(resp)
	return GetResult(marshal, marshalErr)
}
