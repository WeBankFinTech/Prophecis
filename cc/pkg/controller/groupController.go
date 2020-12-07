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
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/groups"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func GetAllGroups(params groups.GetAllGroupsParams) middleware.Responder {
	page := *params.Page
	size := *params.Size
	err := common.CheckPageParams(page, size)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check for page and size", err.Error())
	}
	logger.Logger().Debugf("v1/groups GetAllGroups params page: %v, size: %v", page, size)
	allGroups, _ := service.GetAllGroups(page, size)
	marshal, marshalErr := json.Marshal(allGroups)
	return GetResult(marshal, marshalErr)
}

func AddGroup(params groups.AddGroupParams) middleware.Responder {
	checkForGroupErrMsg := common.CheckForGroup(*params.Group)
	if checkForGroupErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to add group", checkForGroupErrMsg.Error())
	}
	groupRequest := common.FromAddGroupRequest(params)
	logger.Logger().Debugf("AddGroup format to get groupRequest: %v", groupRequest)
	_, err := service.GetGroupByName(groupRequest.Name)
	var group *models.Group
	if err != nil {
		logger.Logger().Debugf("AddGroup format to add")
		group, err = service.AddGroup(groupRequest)
		if err != nil {
			logger.Logger().Error("AddGroup err,",err)
			return ResponderFunc(http.StatusBadRequest, "failed to add group", err.Error())
		}
	}else{
		deleteGroupByName, err := service.GetGroupByName(groupRequest.Name)
		if err == nil {
			if groupRequest.Name == deleteGroupByName.Name {
				logger.Logger().Debugf("AddGroup format to get deleteGroupByName: %v", deleteGroupByName)
				groupRequest.ID = deleteGroupByName.ID
				group, err = service.UpdateGroupByDB(groupRequest)
				if err != nil {
					return ResponderFunc(http.StatusBadRequest, "failed to update group", err.Error())
				}
			}
		}
	}
	marshal, err := json.Marshal(group)
	if err != nil {
		logger.Logger().Errorf("Json Marshal Error:%v",err.Error())
	}
	return GetResult(marshal, err)
}

func UpdateGroup(params groups.UpdateGroupParams) middleware.Responder {
	checkForGroupErrMsg := common.CheckForGroup(*params.Group)
	db := datasource.GetDB()

	if checkForGroupErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update group", checkForGroupErrMsg.Error())
	}

	if params.Group.ID == 0 {
		return ResponderFunc(http.StatusBadRequest, "failed to update group", "group id is invalid")
	}

	groupRequest := common.FromUpdateGroupRequest(params)

	groupByName, err := service.GetGroupByGroupId(groupRequest.ID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group", err.Error())
	}
	if groupRequest.ID != groupByName.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update group", "group not exist in db")
	}

	if groupRequest.Name != groupByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to update group", "The group id and name do not match")
	}

	group, err := service.UpdateGroup(groupRequest, db)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update group", err.Error())
	}

	marshal, marshalErr := json.Marshal(group)

	return GetResult(marshal, marshalErr)
}

func GetGroupByGroupId(params groups.GetGroupByGroupIDParams) middleware.Responder {
	groupId := params.GroupID

	groupByGroupId, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group:%v", err.Error())
	}

	if groupId != groupByGroupId.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get group", "group not exist in db")
	}

	marshal, marshalErr := json.Marshal(groupByGroupId)

	return GetResult(marshal, marshalErr)
}

func DeleteGroupById(params groups.DeleteGroupByIDParams) middleware.Responder {
	groupId := params.GroupID

	groupByGroupId, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by id:%v", err.Error())
	}
	if groupId != groupByGroupId.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete group", "group not exist in db")
	}
	err = service.DeleteGroupById(groupId)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to delete group", err.Error())
	}

	marshal, marshalErr := json.Marshal(groupByGroupId)

	return GetResult(marshal, marshalErr)
}

func GetGroupByName(params groups.GetGroupByNameParams) middleware.Responder {
	groupName := params.GroupName

	repoGroup, err := service.GetGroupByName(groupName)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by name:%v", err.Error())
	}
	if groupName != repoGroup.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to get group by name", "group is not exist in db")
	}

	marshal, marshalErr := json.Marshal(repoGroup)

	return GetResult(marshal, marshalErr)
}

func DeleteGroupByName(params groups.DeleteGroupByNameParams) middleware.Responder {
	groupName := params.GroupName

	groupByName, err := service.GetGroupByName(groupName)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by name:%v", err.Error())
	}
	if groupName != groupByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to delete group by name", "group is not exist in db")
	}

	repoGroup, err := service.DeleteGroupByName(groupName)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to delete group by name:%v", err.Error())
	}
	marshal, marshalErr := json.Marshal(repoGroup)

	return GetResult(marshal, marshalErr)
}

func AddUserToGroup(params groups.AddUserToGroupParams) middleware.Responder {
	//log := logger.Logger()

	userGroupRequest := common.FromAddUserToGroupRequest(params)
	userId := userGroupRequest.UserID
	groupId := userGroupRequest.GroupID

	byUserId,err := service.GetUserByUserId(userId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetUserByUserId:%v", err.Error())
	}

	byGroupId, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by group id:%v", err.Error())
	}
	if byUserId.ID != userId || byGroupId.ID != groupId {
		return ResponderFunc(http.StatusBadRequest, "failed to add user to group", "user or group is not exist in db")
	}
	var repoUserGroup *models.UserGroup
	deleteUserGroup, err := service.GetDeleteUserGroupByUserIdAndGroupId(userId, groupId)
	if err != nil {
		repoUserGroup, err = service.AddUserToGroup(userGroupRequest)
	}else{
		if deleteUserGroup.GroupID == groupId && deleteUserGroup.UserID == userId {
			userGroupRequest.ID = deleteUserGroup.ID
			repoUserGroup, err = service.UpdateUserGroup(userGroupRequest)
		}
	}
	marshal, marshalErr := json.Marshal(repoUserGroup)
	return GetResult(marshal, marshalErr)
}


func UpdateUserGroup(params groups.UpdateUserGroupParams) middleware.Responder {
	userGroupRequest := common.FromUpdateUserGroupRequest(params)

	userGroupByUserIdAndGroupId, err := service.GetUserGroupByUserIdAndGroupId(userGroupRequest.UserID, userGroupRequest.GroupID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by user id and group id:%v", err.Error())
	}

	if userGroupByUserIdAndGroupId.ID != userGroupRequest.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update userGroup", "user is not exist in group")
	}

	repoUserGroup, err := service.UpdateUserGroup(userGroupRequest)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to update user group", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoUserGroup)

	return GetResult(marshal, marshalErr)
}

func DeleteUserFromGroup(params groups.DeleteUserFromGroupParams) middleware.Responder {
	id := params.ID

	userGroupById, err := service.GetUserGroupById(id)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by id", err.Error())
	}
	if id != userGroupById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete user from", "user is not exist in group")
	}

	repoUserGroup, err := service.DeleteUserFromGroupById(id)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to delete user from group by id", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoUserGroup)

	return GetResult(marshal, marshalErr)
}

func GetAllUserGroupByUserId(params groups.GetAllUserGroupByUserIDParams) middleware.Responder {
	userId := params.UserID

	page := *params.Page
	size := *params.Size
	err := common.CheckPageParams(page, size)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check for page and size", err.Error())
	}

	repoUserGroup, err := service.GetAllUserGroupByUserId(userId, page, size)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get all user group by user id", err.Error())
	}
	logger.Logger().Debugf("GetAllUserGroupByUserId repoUserGroup: %v", repoUserGroup)

	marshal, marshalErr := json.Marshal(repoUserGroup)

	return GetResult(marshal, marshalErr)
}

func GetUserGroupByUserIdAndGroupId(params groups.GetUserGroupByUserIDAndGroupIDParams) middleware.Responder {
	userId := params.UserID
	groupId := params.GroupID
	logger.Logger().Debugf("GetUserGroupByUserIdAndGroupId with userId: %v and groupId: %v", userId, groupId)
	userGroup, err := service.GetUserGroupByUserIdAndGroupId(userId, groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get user group by user id and group id", err.Error())
	}

	if userId != userGroup.UserID || groupId != userGroup.GroupID {
		ResponderFunc(http.StatusBadRequest, "failed to get userGroup by userId and groupId", "userGroup is not exists in db")
	}

	marshal, marshalErr := json.Marshal(userGroup)

	return GetResult(marshal, marshalErr)
}

func DeleteUserGroupByUserIdAndGroupId(params groups.DeleteUserGroupByUserIDAndGroupIDParams) middleware.Responder {
	userId := params.UserID
	groupId := params.GroupID
	logger.Logger().Debugf("DeleteUserGroupByUserIdAndGroupId with userId: %v and groupId: %v", userId, groupId)

	userGroup, err := service.GetUserGroupByUserIdAndGroupId(userId, groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get user group by user id and group id", err.Error())
	}

	if userId != userGroup.UserID || groupId != userGroup.GroupID {
		ResponderFunc(http.StatusBadRequest, "failed to delete userGroup by userId and groupId", "userGroup is not exists in db")
	}

	userGroupById, err := service.DeleteUserFromGroupById(userGroup.ID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to delete user from group by id", err.Error())
	}
	if userId != userGroup.UserID || groupId != userGroup.GroupID {
		ResponderFunc(http.StatusBadRequest, "failed to delete userGroup by userId and groupId", "userGroup is not exists in db")
	}
	marshal, marshalErr := json.Marshal(userGroupById)

	return GetResult(marshal, marshalErr)
}

func AddStorageToGroup(params groups.AddStorageToGroupParams) middleware.Responder {
	groupStorageRequest := common.FromAddStorageToGroupRequest(params)

	deletedGroupStorage, err := service.GetDeleteGroupStorageByStorageIdAndGroupId(groupStorageRequest.StorageID, groupStorageRequest.GroupID)
	var repoGroupStorage *models.GroupStorage
	if err != nil {
		repoGroupStorage, err = service.AddStorageToGroup(groupStorageRequest)
	}else{
		if groupStorageRequest.GroupID == deletedGroupStorage.GroupID && groupStorageRequest.StorageID == deletedGroupStorage.StorageID {
			groupStorageRequest.ID = deletedGroupStorage.ID
			repoGroupStorage, err = service.UpdateGroupStorage(groupStorageRequest)
		}
	}
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to add or update group storage", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoGroupStorage)

	return GetResult(marshal, marshalErr)
}

func UpdateGroupStorage(params groups.UpdateGroupStorageParams) middleware.Responder {
	groupStorageRequest := common.FromUpdateGroupStorageRequest(params)

	groupStorage, err := service.GetGroupStorageByStorageIdAndGroupId(groupStorageRequest.StorageID, groupStorageRequest.GroupID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group storage by storageId and groupId", err.Error())
	}
	if groupStorageRequest.ID != groupStorage.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update groupStorage", "storage is not in group")
	}

	repoGroupStorage, err := service.UpdateGroupStorage(groupStorageRequest)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to update group storage", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoGroupStorage)

	return GetResult(marshal, marshalErr)
}

func DeleteStorageFromGroup(params groups.DeleteStorageFromGroupParams) middleware.Responder {
	id := params.ID

	groupStorageById, err := service.GetGroupStorageById(id)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group storage by id", err.Error())
	}
	if id != groupStorageById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete storage from group", "storage is not exist in group")
	}

	repoGroupStorage, err := service.DeleteStorageFromGroupById(id)

	marshal, marshalErr := json.Marshal(repoGroupStorage)

	return GetResult(marshal, marshalErr)
}

func GetAllGroupStorageByStorageId(params groups.GetAllGroupStorageByStorageIDParams) middleware.Responder {
	storageId := params.StorageID

	page := *params.Page
	size := *params.Size
	cheErr := common.CheckPageParams(page, size)
	if nil != cheErr {
		return ResponderFunc(http.StatusBadRequest, "failed to check for page and size", cheErr.Error())
	}

	repoGroupStorage, dbErr := service.GetAllGroupStorageByStorageId(storageId, page, size)
	if nil != dbErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to get groupStorage by storageId", dbErr.Error())
	}

	marshal, marshalErr := json.Marshal(repoGroupStorage)

	return GetResult(marshal, marshalErr)
}

func AddNamespaceToGroup(params groups.AddNamespaceToGroupParams) middleware.Responder {
	gNSRequest := common.FromAddNamespaceToGroupRequest(params)

	groupId := gNSRequest.GroupID
	group, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get group by group id%v", err.Error())
	}
	if groupId != group.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to add namespace to group", "group is not exist in db")
	}

	namespaceId := gNSRequest.NamespaceID
	namespace, err := service.GetNamespaceByID(namespaceId)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetNamespaceByID", err.Error())
	}
	if namespaceId != namespace.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to add namespace to group", "namespace is not exist in db")
	}
	delGroupNSDB, err := service.GetDeleteGroupNamespaceByGroupIdAndNamespaceId(gNSRequest.GroupID, gNSRequest.NamespaceID)
	var repoGroupNamespace *models.GroupNamespace
	if err != nil {
		repoGroupNamespace, err = service.AddNamespaceToGroup(gNSRequest)
	}else{
		if gNSRequest.GroupID == delGroupNSDB.GroupID && gNSRequest.NamespaceID == delGroupNSDB.NamespaceID {
			gNSRequest.ID = delGroupNSDB.ID
			repoGroupNamespace, err = service.UpdateGroupNamespace(gNSRequest)
		}
	}
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to add or update namespace: %v", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoGroupNamespace)

	return GetResult(marshal, marshalErr)
}

func DeleteGroupNamespace(params groups.DeleteGroupNamespaceParams) middleware.Responder {
	id := params.ID

	groupNamespaceById, err := service.GetGroupNamespaceById(id)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupNamespaceById: %v", err.Error())
	}

	if id != groupNamespaceById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete namespace from group", "namespace is not exist in group")
	}
	repoGroupNamespace, _ := service.DeleteGroupNamespaceById(id)
	marshal, marshalErr := json.Marshal(repoGroupNamespace)
	return GetResult(marshal, marshalErr)
}

func GetAllGroupNamespaceByNamespaceId(params groups.GetAllGroupNamespaceByNamespaceIDParams) middleware.Responder {
	namespaceId := params.NamespaceID

	page := *params.Page
	size := *params.Size
	err := common.CheckPageParams(page, size)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check for page and size", err.Error())
	}
	namespace, err := service.GetNamespaceByID(namespaceId)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetNamespaceByID", err.Error())
	}
	if namespaceId != namespace.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get groupNamespaces by namespaceId", "namespace is not exist in db")
	}

	groupNamespaces, err := service.GetAllGroupNamespaceByNamespaceId(namespaceId, page, size)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetAllGroupNamespaceByNamespaceId: %v", err.Error())
	}
	marshal, marshalErr := json.Marshal(groupNamespaces)

	return GetResult(marshal, marshalErr)
}

func GetCurrentUserNamespaceWithRole(params groups.GetCurrentUserNamespaceWithRoleParams) middleware.Responder {
	roleId := params.RoleID
	clusterName := params.ClusterName

	token := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)

	logger.Logger().Debugf("GetCurrentUserNamespaceWithRole by clusterName: %v and roleId: %v, token: %v", clusterName, roleId, token)
	if "" == token {
		return ResponderFunc(http.StatusForbidden, "failed to GetCurrentUserNamespaceWithRole", "token is empty")
	}
	sessionUser := GetSessionByContext(params.HTTPRequest)
	if nil == sessionUser || sessionUser.UserName == "" {
		return FailedToGetUserByToken()
	}

	var groupNamespaces []models.GroupNamespace
	if IsSuperAdmin(sessionUser) {
		groupNamespaces, err := service.GetAllGroupNamespace()
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetAllGroupNamespace:%v", err.Error())
		}
		logger.Logger().Debugf("GetCurrentUserNamespaceWithRole get groupNamespaces by SA: %v", groupNamespaces)
	} else {
		userName := sessionUser.UserName
		userByName, err := service.GetUserByUserName(userName)
		if err != nil {
			return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
		}
		groupIds, err := service.GetGroupIdListByUserIdRoleId(userByName.ID, roleId)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserIdRoleId:%v", err.Error())
		}
		logger.Logger().Debugf("GetCurrentUserNamespaceWithRole get userByName: %v and groupIds: %v", userByName, groupIds)
		if len(groupIds) > 0 {
			groupNamespaces, err = service.GetAllGroupNamespaceByGroupIds(groupIds)
			if err != nil {
				return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserIdRoleId:%v", err.Error())
			}
		}
	}

	logger.Logger().Debugf("GetCurrentUserNamespaceWithRole get groupNamespaces: %v", groupNamespaces)

	//gns := NamespaceFilter(groupNamespaces, clusterName)
	gns := groupNamespaces

	logger.Logger().Debugf("GetCurrentUserNamespaceWithRole to res ns: %v", gns)

	marshal, marshalErr := json.Marshal(gns)

	return GetResult(marshal, marshalErr)
}

func GetNamespaceByGroupIDAndNamespace(params groups.GetNamespaceByGroupIDAndNamespaceParams) middleware.Responder {
	groupId := params.GroupID
	group, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupByGroupId:%v", err.Error())
	}
	if groupId != group.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get groupNamespaces", "group is not exist in db")
	}

	namespaceId := params.NamespaceID
	namespace, err := service.GetNamespaceByID(namespaceId)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetNamespaceByID", err.Error())
	}
	if namespaceId != namespace.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get groupNamespaces", "namespace is not exist in db")
	}

	repo, err := service.GetGroupNamespaceByGroupIdAndNamespaceId(groupId, namespaceId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupNamespaceByGroupIdAndNamespaceId:%v", err.Error())
	}

	marshal, marshalErr := json.Marshal(repo)

	return GetResult(marshal, marshalErr)
}

func GetNamespacesByGroupId(params groups.GetNamespacesByGroupIDParams) middleware.Responder {
	groupId := params.GroupID
	group, err := service.GetGroupByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetNamespacesByGroupId:%v", err.Error())
	}
	if groupId != group.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get groupNamespace bu groupId", "group is not exist in db")
	}

	repo, err := service.GetAllGroupNamespaceByGroupId(groupId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetAllGroupNamespaceByGroupId:%v", err.Error())
	}
	marshal, marshalErr := json.Marshal(repo)

	return GetResult(marshal, marshalErr)
}

func GetCurrentUserStoragePath(params groups.GetCurrentUserStoragePathParams) middleware.Responder {
	sessionUser := GetSessionByContext(params.HTTPRequest)

	logger.Logger().Debugf("GetCurrentUserStoragePath by sessionUser: %v", sessionUser)

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	var myStorage = common.NewStringSet()
	var gss []models.GroupStorage
	if IsSuperAdmin(sessionUser) {
		groupStorage, err := service.AllGroupStorage()
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetNamespacesByGroupId:%v", err.Error())
		}
		gss = groupStorage
	} else {
		username := sessionUser.UserName
		user, err := service.GetUserByUserName(username)
		if err != nil {
			return ResponderFunc(http.StatusBadRequest, "failed to GetUserByUserName", err.Error())
		}
		groupIds, err := service.GetGroupIdListByUserId(user.ID)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserId:%v", err.Error())
		}
		gss, err = service.GetGroupStorageByGroupIds(groupIds)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupStorageByGroupIds:%v", err.Error())
		}
	}
	logger.Logger().Debugf("GetCurrentUserStoragePath gss: %v", len(gss))
	myStorage = common.GetMyStoragePath(gss, myStorage, "")

	marshal, marshalErr := json.Marshal(myStorage.List())

	return GetResult(marshal, marshalErr)
}
