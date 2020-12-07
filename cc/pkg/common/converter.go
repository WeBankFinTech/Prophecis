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
package common

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/alerts"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/groups"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/namespaces"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/resources"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/roles"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/storages"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/users"
	"strconv"
	"strings"
)

func FromAddUserRequest(userRequest models.UserRequest) models.User {
	logger.Logger().Debugf("v1/users FromAddUserRequest params: %v", userRequest)
	var user = models.User{
		Name:       userRequest.Name,
		UID:        userRequest.UID,
		Gid:        userRequest.Gid,
		Type:       userRequest.Type,
		GUIDCheck:  userRequest.GUIDCheck,
		Remarks:    userRequest.Remarks,
		EnableFlag: 1,
		Token:      userRequest.Token,
	}
	return user
}

func FromUpdateUserRequest(params users.UpdateUserParams) models.User {
	logger.Logger().Debugf("v1/users FromUpdateUserRequest params: %v", params)
	var remarks = " "
	if params.User.Remarks != "" {
		remarks = params.User.Remarks
	}
	var user = models.User{
		ID:         params.User.ID,
		Name:       params.User.Name,
		UID:        params.User.UID,
		Gid:        params.User.Gid,
		Type:       params.User.Type,
		GUIDCheck:  params.User.GUIDCheck,
		Remarks:    remarks,
		EnableFlag: 1,
		Token:      params.User.Token,
	}
	return user
}

func FromAddGroupRequest(params groups.AddGroupParams) models.Group {
	logger.Logger().Debugf("v1/groups FromAddGroupRequest params: %v", params)
	//var departmentID = "0"
	//if "" != params.Group.DepartmentID {
	//	departmentID = params.Group.DepartmentID
	//}
	var group = models.Group{
		Name:           params.Group.Name,
		GroupType:      params.Group.GroupType,
		SubsystemID:    params.Group.SubsystemID,
		SubsystemName:  params.Group.SubsystemName,
		Remarks:        params.Group.Remarks,
		EnableFlag:     1,
		DepartmentID:   params.Group.DepartmentID,
		DepartmentName: params.Group.DepartmentName,
		ClusterName:    params.Group.ClusterName,
	}
	return group
}

func FromUpdateGroupRequest(params groups.UpdateGroupParams) models.Group {
	logger.Logger().Debugf("v1/groups FromUpdateGroupRequest params: %v", params)
	var remarks = " "
	if params.Group.Remarks != "" {
		remarks = params.Group.Remarks
	}
	var group = models.Group{
		ID:             params.Group.ID,
		Name:           params.Group.Name,
		GroupType:      params.Group.GroupType,
		SubsystemID:    params.Group.SubsystemID,
		SubsystemName:  params.Group.SubsystemName,
		Remarks:        remarks,
		EnableFlag:     1,
		DepartmentID:   params.Group.DepartmentID,
		DepartmentName: params.Group.DepartmentName,
		ClusterName:    params.Group.ClusterName,
	}
	return group
}

func FromAddStorageRequest(params storages.AddStorageParams) models.Storage {
	logger.Logger().Debugf("v1/storage FromAddStorageRequest params: %v", params)
	var storage = models.Storage{
		Path:       params.Storage.Path,
		Remarks:    params.Storage.Remarks,
		EnableFlag: 1,
		Type:       params.Storage.Type,
	}
	return storage
}

func FromUpdateStorageRequest(params storages.UpdateStorageParams) models.Storage {
	logger.Logger().Debugf("v1/storage FromUpdateStorageRequest params: %v", params)
	var remarks = " "
	if params.Storage.Remarks != "" {
		remarks = params.Storage.Remarks
	}
	var storage = models.Storage{
		ID:         params.Storage.ID,
		Path:       params.Storage.Path,
		Remarks:    remarks,
		EnableFlag: 1,
		Type:       params.Storage.Type,
	}
	return storage
}

func FromAddStorageToGroupRequest(params groups.AddStorageToGroupParams) models.GroupStorage {
	logger.Logger().Debugf("v1/groups FromAddStorageToGroupRequest params: %v", params)
	var groupStorage = models.GroupStorage{
		GroupID:     params.GroupStorage.GroupID,
		StorageID:   params.GroupStorage.StorageID,
		Path:        params.GroupStorage.Path,
		Permissions: params.GroupStorage.Permissions,
		Remarks:     params.GroupStorage.Remarks,
		EnableFlag:  1,
		Type:        params.GroupStorage.Type,
	}
	return groupStorage
}

func FromUpdateGroupStorageRequest(params groups.UpdateGroupStorageParams) models.GroupStorage {
	logger.Logger().Debugf("v1/groups FromUpdateGroupStorageRequest params: %v", params)
	var remarks = " "
	if params.GroupStorage.Remarks != "" {
		remarks = params.GroupStorage.Remarks
	}
	var groupStorage = models.GroupStorage{
		ID:          params.GroupStorage.ID,
		GroupID:     params.GroupStorage.GroupID,
		StorageID:   params.GroupStorage.StorageID,
		Path:        params.GroupStorage.Path,
		Permissions: params.GroupStorage.Permissions,
		Remarks:     remarks,
		EnableFlag:  1,
		Type:        params.GroupStorage.Type,
	}
	return groupStorage
}

func FromAddUserToGroupRequest(params groups.AddUserToGroupParams) models.UserGroup {
	logger.Logger().Debugf("v1/groups FromAddUserToGroupRequest params: %v", params)
	var UserGroup = models.UserGroup{
		UserID:     params.UserGroup.UserID,
		RoleID:     params.UserGroup.RoleID,
		GroupID:    params.UserGroup.GroupID,
		Remarks:    params.UserGroup.Remarks,
		EnableFlag: 1,
	}
	return UserGroup
}

func FromUpdateUserGroupRequest(params groups.UpdateUserGroupParams) models.UserGroup {
	logger.Logger().Debugf("v1/groups FromUpdateUserGroupRequest params: %v", params)
	var remarks = " "
	if params.UserGroup.Remarks != "" {
		remarks = params.UserGroup.Remarks
	}
	var UserGroup = models.UserGroup{
		ID:         params.UserGroup.ID,
		UserID:     params.UserGroup.UserID,
		RoleID:     params.UserGroup.RoleID,
		GroupID:    params.UserGroup.GroupID,
		Remarks:    remarks,
		EnableFlag: 1,
	}
	return UserGroup
}

func FromAddNamespaceRequest(params namespaces.AddNamespaceParams) models.NamespaceRequest {
	var namespaceRequest = models.NamespaceRequest{
		Namespace:         params.Namespace.Namespace,
		Annotations:       params.Namespace.Annotations,
		PlatformNamespace: params.Namespace.PlatformNamespace,
		Remarks:           params.Namespace.Remarks,
		EnableFlag:        1,
	}
	return namespaceRequest
}

func FromUpdateNamespaceRequest(params namespaces.UpdateNamespaceParams) models.NamespaceRequest {
	var remarks = " "
	if params.Namespace.Remarks != "" {
		remarks = params.Namespace.Remarks
	}
	var namespaceRequest = models.NamespaceRequest{
		ID:                params.Namespace.ID,
		Namespace:         params.Namespace.Namespace,
		Annotations:       params.Namespace.Annotations,
		PlatformNamespace: params.Namespace.PlatformNamespace,
		Remarks:           remarks,
		EnableFlag:        1,
	}
	return namespaceRequest
}

func FromNamespaceRequestToNamespace(namespaceRequest models.NamespaceRequest) models.Namespace {
	var namespace = models.Namespace{
		ID:         namespaceRequest.ID,
		Namespace:  namespaceRequest.Namespace,
		Remarks:    namespaceRequest.Remarks,
		EnableFlag: 1,
	}
	return namespace
}

func FromAddNamespaceToGroupRequest(params groups.AddNamespaceToGroupParams) models.GroupNamespace {
	var groupNamespace = models.GroupNamespace{
		GroupID:     params.GroupNamespace.GroupID,
		NamespaceID: params.GroupNamespace.NamespaceID,
		Namespace:   params.GroupNamespace.Namespace,
		Remarks:     params.GroupNamespace.Remarks,
		EnableFlag:  1,
	}
	return groupNamespace
}

func FromNamespaceToNamespaceResponse(namespaces []*models.Namespace) []*models.NamespaceResponse {
	var namespaceResList []*models.NamespaceResponse

	for _, ns := range namespaces {
		var namespaceRes = &models.NamespaceResponse{
			ID:         ns.ID,
			Namespace:  ns.Namespace,
			EnableFlag: ns.EnableFlag,
			Remarks:    ns.Remarks,
		}
		namespaceResList = append(namespaceResList, namespaceRes)
	}
	return namespaceResList
}

func FromSetNamespaceRQRequest(params resources.SetNamespaceRQParams) models.ResourcesQuota {
	var rq = models.ResourcesQuota{
		Namespace: params.ResourcesQuota.Namespace,
		CPU:       params.ResourcesQuota.CPU,
		Gpu:       params.ResourcesQuota.Gpu,
		Memory:    params.ResourcesQuota.Memory,
	}
	return rq
}

func FromAddRoleRequest(params roles.AddRoleParams) models.Role {
	var role = models.Role{
		Name:       params.Role.Name,
		Remarks:    params.Role.Remarks,
		EnableFlag: 1,
	}

	return role
}

func FromUpdateRoleRequest(params roles.UpdateRoleParams) models.Role {
	var role = models.Role{
		ID:         params.Role.ID,
		Name:       params.Role.Name,
		Remarks:    params.Role.Remarks,
		EnableFlag: 1,
	}

	return role
}

func FromAlertRequest(params alerts.ReceiveTaskAlertParams, alertIP string, alertReceiver string) ([]models.IMSAlert, error) {

	if strings.Contains(alertIP, ":") {
		alertIP = alertIP[0:strings.Index(alertIP, ":")]
	}

	req := params.JobRequest
	var imsAlerts []models.IMSAlert
	alertLevel := req.AlertLevel
	logger.Logger().Debugf("alertLevel: %v", alertLevel)

	var title = "MLSS-DI训练任务告警"
	receivers := strings.Split(alertReceiver, ",")
	receiver := req.Receiver
	userId := req.UserID
	//format level to int
	if "" != receiver {
		split := strings.Split(receiver, ",")
		for _, re := range split {
			receivers = append(receivers, re)
		}
	}
	subSystemId, parErr := strconv.ParseInt(GetAppConfig().Core.Ims.SubSystemId, 10, 64)

	if nil != parErr {
		return nil, parErr
	}

	receivers = append(receivers, userId)
	var content = fmt.Sprintf("任务ID：%s\n任务提交人：%s\n任务状态：%s\n告警原因：%s\n开始时间：%s\n结束时间：%s\n", req.TrainingID, req.UserID, req.JobStatus, req.AlertReason, req.StartTime, req.EndTime)
	//TODO set alert level
	for _, rec := range receivers {
		var imsAlert = models.IMSAlert{
			AlertIP:       alertIP,
			AlertLevel:    1,
			AlertObj:      req.TrainingID,
			AlertTitle:    title,
			AlertInfo:     content,
			SubSystemID:   subSystemId,
			AlertWay:      "1,2,3,4",
			AlertReceiver: rec,
		}
		imsAlerts = append(imsAlerts, imsAlert)
	}

	return imsAlerts, nil
}

func FromKeyRequestToKeyPair(request *models.KeyPairRequest) models.Keypair {
	var key = models.Keypair{
		EnableFlag: 1,
		APIKey:     uuid.New().String(),
		SecretKey:  uuid.New().String(),
		Name:       request.Name,
		Remarks:    request.Remarks,
	}
	return key
}

func FromInterfaceToSessionUser(int interface{}) (*models.SessionUser, error) {
	var sessionUser *models.SessionUser
	bytes, e := json.Marshal(int)
	if nil != e {
		return nil, e
	}
	e = json.Unmarshal(bytes, &sessionUser)
	if nil != e {
		return nil, e
	}
	return sessionUser, nil
}

func GetMyNamespace(gns []models.GroupNamespace, set *StringSet) *StringSet {
	if nil != gns && len(gns) > 0 {
		for _, gn := range gns {
			set.Add(gn.Namespace)
		}
	}
	return set
}

func GetMyStoragePath(gss []models.GroupStorage, set *StringSet, cluster string) *StringSet {
	if nil != gss && len(gss) > 0 {
		for _, gn := range gss {
			path := strings.TrimSpace(gn.Path)
			set.Add(path)
		}
	}
	return set
}

func FromUserToUserRes(user models.User) models.UserRes {
	return models.UserRes{
		ID:         user.ID,
		Name:       user.Name,
		GUIDCheck:  "1" == user.GUIDCheck,
		Gid:        user.Gid,
		UID:        user.UID,
		Remarks:    user.Remarks,
		EnableFlag: user.EnableFlag,
		Type:       user.Type,
		Token:      user.Token,
	}
}
