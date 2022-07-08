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
	"errors"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/manager"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"strings"
)

func UserStorageCheck(username string, path string) error {
	saByName, err := repo.GetSAByName(username)
	db := datasource.GetDB()
	var ids []int64
	if err != nil {
		userByUserName, err := repo.GetUserByName(username, db)
		if err != nil {
			logger.Logger().Errorf("GetUserByName error:%v", err.Error())
			return err
		}
		if userByUserName == nil || username != userByUserName.Name {
			return errors.New("user is not exist in db")
		}
		groupIds, err := repo.GetGroupIdsByUserId(userByUserName.ID)
		if err != nil {
			logger.Logger().Errorf("GetGroupIdsByUserId error:%v", err.Error())
			return err
		}
		ids = *groupIds
	} else {
		logger.Logger().Debugf("UserStorageCheck saByName: %v", saByName)
		if username == saByName.Name {
			groupIds, err := repo.GetAllGroupIds()
			if err != nil {
				logger.Logger().Errorf("GetAllGroupIds error:%v", err.Error())
				return err
			}
			ids = *groupIds
		}
	}
	logger.Logger().Debugf("UserStorageCheck ids: %v", ids)

	if nil == ids || len(ids) <= 0 {
		return errors.New("user is not exist in any group")
	}

	groupStorageByGroupIds, err := repo.GetAllGroupStorageByGroupIds(ids)
	if err != nil {
		logger.Logger().Errorf("GetAllGroupStorageByGroupIds error:%v", err.Error())
		return err
	}

	if len(groupStorageByGroupIds) < 1 {
		return errors.New("auth failed and the user can't access any storages")
	}

	logger.Logger().Debugf("UserStorageCheck groupStorageByGroupIds: %v", len(groupStorageByGroupIds))

	var flag bool = false
	for _, groupStorage := range groupStorageByGroupIds {
		if strings.HasPrefix(path, groupStorage.Path) {
			logger.Logger().Debugf("UserStorageCheck path: %v, gs path: %v", path, groupStorage.Path)
			flag = true
			break
		}
	}

	if !flag {
		return errors.New("auth failed for path")
	}

	return nil
}

func UserNamespaceCheck(username string, namespace string) (*models.UserNamespaceAccessResponse, error) {
	db := datasource.GetDB()
	if "" == username || "" == namespace {
		return nil, errors.New("username and namespace can't be empty")
	}

	logger.Logger().Debugf("UserNamespaceCheck username: %v and namespace: %v", username, namespace)
	userFromDB, err := repo.GetUserByName(username, db)
	if err == nil {
		return nil, err
	}
	if userFromDB == nil || username != userFromDB.Name {
		return nil, errors.New("user is not exist in db")
	}
	getNamespace, err := repo.GetNamespaceByName(namespace)
	if err != nil {
		logger.Logger().Error("Get namespace by name err, ", err)
		return nil, err
	}
	if namespace != getNamespace.Namespace {
		return nil, errors.New("namespace is not exist in db")
	}
	var accessible = false
	sa, err := repo.GetSAByName(username)
	if err != nil {
		logger.Logger().Errorf("GetSAByName error:%v", err.Error())
		return nil, err
	}
	if username == sa.Name {
		accessible = true
	} else {
		groupNamespaces, err := GetAllNamespaceByUserId(userFromDB.ID)
		if nil != err {
			return nil, err
		}
		if len(groupNamespaces) > 0 {
			for _, groupNamespace := range groupNamespaces {
				if namespace == groupNamespace.Namespace {
					accessible = true
					break
				}
			}
		}
	}

	var response = &models.UserNamespaceAccessResponse{
		Username:   username,
		Namespace:  namespace,
		Accessible: accessible,
	}

	if !response.Accessible {
		return nil, errors.New("the query of the namespace is not allowed")
	}

	return response, nil
}

func UserStoragePathCheck(username string, namespace string, path string) error {
	db := datasource.GetDB()
	if "" == username || "" == namespace || "" == path {
		return errors.New("username, namespace and path can't be empty")
	}

	sa, err := repo.GetSAByName(username)
	if err != nil {
		logger.Logger().Errorf("GetSAByName error:%v", err.Error())
	}
	if username == sa.Name {
		return nil
	}

	userFromDB, err := repo.GetUserByName(username, db)
	if err == nil {
		return err
	}
	if userFromDB == nil || username != userFromDB.Name {
		return errors.New("user is not exist in db")
	}
	getNamespace, err := repo.GetNamespaceByName(namespace)
	if err != nil {
		logger.Logger().Error("Get namespace by name err, ", err)
		return err
	}
	if namespace != getNamespace.Namespace {
		return errors.New("namespace is not exist in db")
	}
	if username != userFromDB.Name {
		return errors.New("user is not exists in db")
	}

	groupIds, err := GetGroupIdListByUserId(userFromDB.ID)
	if err != nil {
		logger.Logger().Errorf("GetGroupIdListByUserId error:%v", err.Error())
		return err
	}
	if len(groupIds) <= 0 {
		return errors.New("user has not joined any groups")
	}

	groupStorageByGroupIds, err := repo.GetAllGroupStorageByGroupIds(groupIds)
	if err != nil {
		logger.Logger().Errorf("GetAllGroupStorageByGroupIds error:%v", err.Error())
		return err
	}

	var storageExistOfGroupId []int64

	if len(groupStorageByGroupIds) > 0 {
		for _, groupStorage := range groupStorageByGroupIds {
			if strings.HasPrefix(path, groupStorage.Path) {
				storageExistOfGroupId = append(storageExistOfGroupId, groupStorage.GroupID)
			}
		}
	}

	if len(storageExistOfGroupId) <= 0 {
		return errors.New("user not joined a group that have accessing of this path")
	}

	for _, groupId := range storageExistOfGroupId {
		groupNamespaceByGroupId, err := repo.GetAllGroupNamespaceByGroupId(db, groupId)
		if err != nil {
			logger.Logger().Errorf("UserStoragePathCheck for loop GetAllGroupNamespaceByGroupId error:%v", err.Error())
			return err
		}
		for _, groupNS := range groupNamespaceByGroupId {
			if namespace == groupNS.Namespace {
				return nil
			}
		}
	}

	return errors.New("no group namespace matching at the same time")
}

func AdminUserCheck(admin string, username string) error {
	log := logger.Logger()
	db := datasource.GetDB()
	saByAdmin, err := repo.GetSAByName(admin)
	if err != nil {
		logger.Logger().Error("GetSAByName error:%v", err.Error())
		return err
	}
	userByAdmin, err := repo.GetUserByName(admin, db)
	if err != nil || userByAdmin == nil {
		log.Errorf(err.Error())
		return err
	}
	saByName, err := repo.GetSAByName(username)
	if err != nil {
		logger.Logger().Error("GetSAByName error:%v", err.Error())
		return err
	}
	userByName, err := repo.GetUserByName(username, db)
	if err != nil || userByName == nil {
		log.Errorf(err.Error())
		return err
	}

	if (admin != saByAdmin.Name && (admin != userByAdmin.Name || userByAdmin == nil)) || (username != saByName.Name && (username != userByName.Name || userByName == nil)) {
		return errors.New("找不到adminUsername或username")
	}

	if admin == saByAdmin.Name || username == saByName.Name {
		return nil
	}

	if userByAdmin.Name == userByName.Name {
		return nil
	}

	//if adminUsername is GA
	//adminUser and user in the same group, return
	num, err := repo.GetMatchGroupNum(userByAdmin.ID, userByName.ID)
	if err != nil {
		logger.Logger().Error("GetMatchGroupNum error:%v", err.Error())
		return err
	}
	if num > 0 {
		return nil
	}

	return errors.New("auth failed")
}

func CheckNamespace(namespace string, admin string) (string, error) {
	logger.Logger().Debugf("CheckNamespace with namespace: %v and admin: %v", namespace, admin)
	db := datasource.GetDB()
	saByName, err := repo.GetSAByName(admin)
	if err != nil {
		getNamespace, err := repo.GetNamespaceByName(namespace)
		if err != nil {
			logger.Logger().Error("Get namespace by name err, ", err)
			return "", err
		}
		if namespace != getNamespace.Namespace {
			return "", errors.New("namespace is not exist in db")
		}
		userByName, err := repo.GetUserByName(admin, db)
		if err != nil {
			return "", err
		}
		if userByName == nil || admin != userByName.Name {
			return "", errors.New("admin is not exists in db")
		}

		roleByName, err := GetRoleByName(constants.ADMIN)
		if err != nil {
			logger.Logger().Error("GetRoles err, ", err)
			return "", err
		}
		if constants.ADMIN != roleByName.Name {
			return "", errors.New("role:ADMIN does not existed in db")
		}

		ns, err := GetAllGroupNamespaceByNamespaceId(getNamespace.ID, 0, 0)
		if err != nil {
			logger.Logger().Errorf("GetAllGroupNamespaceByNamespaceId error:%v", err.Error())
			return "", err
		}
		groupNamespaces := ns.Models

		if len(groupNamespaces) < 0 {
			return "", errors.New("namespace is not exists in any group")
		}
		for _, gn := range groupNamespaces {
			logger.Logger().Debugf("CheckNamespace to get gn: %v", gn)
			userGroup, err := GetUserGroupByUserIdAndGroupId(userByName.ID, gn.GroupID)
			if err != nil {
				logger.Logger().Errorf("CheckNamespace for loop GetUserGroupByUserIdAndGroupId error:%v", err.Error())
				return "", err
			}
			if userGroup.RoleID == roleByName.ID {
				return constants.GA, nil
			}
			if userGroup.UserID == userByName.ID && userGroup.GroupID == gn.GroupID {
				return constants.GU, nil
			}
		}
	} else {
		if admin == saByName.Name {
			return "SA", nil
		}
	}
	return "", errors.New("user not in this group")
}

func CheckNamespaceUser(namespace string, admin string, username string) error {
	err := CheckAdmin(namespace, admin, username)

	if nil != err {
		return err
	}

	return nil
}

func CheckUserGetNamespace(admin string, username string) (*models.UserNoteBook, error) {
	db := datasource.GetDB()
	var userNotebook *models.UserNoteBook
	//if admin == username {
	//	userNotebook = &models.UserNoteBook{
	//		Role: "SELF",
	//	}
	//	return userNotebook, nil
	//}

	saByAdmin, err := repo.GetSAByName(username)
	if err != nil {
		logger.Logger().Errorf("GetSAByName error:%v", err.Error())
	}
	if admin == saByAdmin.Name {
		userNotebook = &models.UserNoteBook{
			Role: "SA",
		}
		return userNotebook, nil
	}

	user, err := repo.GetUserByName(username, db)
	if err != nil {
		logger.Logger().Errorf("GetUserByName error:%v", err.Error())
		return nil, err
	}
	groupIds, err := repo.GetGroupIdListByUserIdRoleId(user.ID, constants.GA_ID)
	if err != nil {
		logger.Logger().Errorf("GetGroupIdsByUserIds error:", err.Error())
		return nil, err
	}

	groupNamespaceByGroupIds, err := repo.GetAllGroupNamespaceByGroupIds(*groupIds)
	if err != nil {
		logger.Logger().Errorf("groupNamespaceByGroupIds error:", err.Error())
		return nil, err
	}
	var namespaces = common.NewStringSet()
	for _, gn := range groupNamespaceByGroupIds {
		namespaces.Add(gn.Namespace)
	}
	namespaceList := namespaces.SortList()
	if len(*groupIds) == 1 && len(namespaceList) == 0 {
		userNotebook := models.UserNoteBook{
			Role: "GU",
		}
		return &userNotebook, nil
	}
	userNotebook = &models.UserNoteBook{
		Role:       "GA",
		Namespaces: namespaceList,
	}

	return userNotebook, nil
}

func CheckUserNotebook(namespace string, username string, notebookUser string, notebook string) error {
	podList, err := manager.GetK8sApiClient().GetPodsByNS(namespace)

	if nil != err {
		return err
	}

	if nil == podList || len(podList.Items) <= 0 {
		return errors.New("auth failed, get no pod by namespace")
	}

	nbPrefix := "/notebook/" + namespace + "/" + notebook

	for _, pod := range podList.Items {
		vars := pod.Spec.Containers[0].Env
		var nbUser, nbPre string
		for _, env := range vars {
			if env.Name == "NB_PREFIX" {
				nbPre = env.Value
			}
			if env.Name == "NB_USER" {
				nbUser = env.Value
			}
		}
		if nbPrefix == nbPre && username == nbUser {
			logger.Logger().Debugf("check notebook: %v for user: %v by namespace: %v success", notebook, username, namespace)
			return nil
		}
	}

	return errors.New("auth failed")
}

func GetNotebookProxyUserAndUser(namespace string, name string) (string, string, error) {
	envs, err := manager.GetK8sApiClient().GetPodEnvlByNSAndName(namespace, name+"-0")
	var user string
	var proxyUser string
	var nbUser string
	if err != nil {
		return "", "", err
	}
	for _, env := range envs {
		if env.Name == "CREATE_USER" {
			user = env.Value
		}
		if env.Name == "PROXY_USER" {
			proxyUser = env.Value
		}
		if env.Name == "NB_USER" {
			nbUser = env.Value
		}
	}

	//Fix Old Version Notebook
	if user == "" {
		user = nbUser
	}

	return user, proxyUser, nil
}

func GetMFServiceUser(namespace string, mfService string) (*string, error) {
	sldp, err := manager.GetK8sApiClient().GetSeldonDeployment(namespace, mfService)
	if err != nil {
		return nil, err
	}
	user := sldp.Labels["MLSS-USER"]
	return &user, nil
}
