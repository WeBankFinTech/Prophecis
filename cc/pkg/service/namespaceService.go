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
	"github.com/jinzhu/gorm"
	"math"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/manager"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"strings"
)

func GetAllNamespaces(page int64, size int64) (*models.PageNamespaceList, error) {
	var allNamespaces []*models.Namespace
	var err error
	offSet := common.GetOffSet(page, size)
	if page > 0 && size > 0 {
		allNamespaces, err = repo.GetAllNamespacesByOffset(offSet, size)
		if err != nil {
			logger.Logger().Error("Get all namespace by offset err, ", err)
			return &models.PageNamespaceList{}, err
		}
	} else {
		allNamespaces, err = repo.GetAllNamespaces()
		if err != nil {
			logger.Logger().Error("Get all namespaces err, ", err)
			return &models.PageNamespaceList{}, err
		}
	}

	total, err := repo.CountNamespaces()
	if err != nil {
		logger.Logger().Error("Count namespaces err, ", err)
		return &models.PageNamespaceList{}, err
	}

	namespaceResList := common.FromNamespaceToNamespaceResponse(allNamespaces)

	kac := manager.GetK8sApiClient()

	namespacesFromK8s, listNSErr := kac.GetAllNamespacesFromK8s()

	if nil != listNSErr {
		return nil, listNSErr
	}

	namespaceForK8s := namespacesFromK8s.Items

	logger.Logger().Debugf("GetAllNamespaces len(namespaceForK8s): %v", len(namespaceForK8s))
	//var nsList []*models.NamespaceResponse

	if nil != namespaceResList && len(namespaceResList) > 0 {
		for _, ns := range namespaceResList {
			if nil != namespaceForK8s {
				for _, v1NS := range namespaceForK8s {
					if ns.Namespace == v1NS.ObjectMeta.Name {
						if nil != v1NS.ObjectMeta.Annotations {
							nodeSelector := v1NS.ObjectMeta.Annotations[constants.NodeSelectorKey]
							if "" != nodeSelector {
								customSelectorArr := strings.Split(nodeSelector, ",")
								ns.Annotations = map[string]string{}
								for _, customSelector := range customSelectorArr {
									customSelectorKey := strings.Split(customSelector, "=")[0]
									if "lb-gpu-model" == customSelectorKey && len(strings.Split(customSelector, "=")) > 1 {
										ns.Annotations["lb-gpu-model"] = strings.Split(customSelector, "=")[1]
									}
									if "lb-bus-type" == customSelectorKey && len(strings.Split(customSelector, "=")) > 1 {
										ns.Annotations["lb-bus-type"] = strings.Split(customSelector, "=")[1]
									}
								}
							}
						}
						quotaList, reErr := kac.ListRqOfNamespace(ns.Namespace)
						if nil != reErr {
							return nil, reErr
						}
						quota := quotaList.Items[0]
						hard := quota.Spec.Hard
						//ns.Hard = hard

						hardMap := make(map[string]interface{})
						for k, v := range hard {
							hardMap[k.String()] = v.String()
						}
						ns.Hard = hardMap
						//nsList = append(nsList, ns)
					}
				}
			}
		}
	}

	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	var pageNamespace = &models.PageNamespaceList{
		Models:     namespaceResList,
		TotalPage:  int64(totalPage),
		Total:      total,
		PageNumber: page,
		PageSize:   size,
	}

	return pageNamespace, nil
}

func GetAllNamespace() ([]models.GroupNamespace, error) {
	groupNamespace, err := repo.GetAllGroupNamespace()
	if err != nil {
		logger.Logger().Errorf("GetAllGroupNamespace error:%v", err.Error())
		return nil, err
	}
	return groupNamespace, err
}

func GetNamespaceByID(id int64) (models.Namespace, error) {
	return repo.GetNamespaceByID(id)
}

func DeleteNamespaceByID(id int64) (*models.Namespace, error) {
	groupNamespaceByNamespaceId, err := repo.GetAllGroupNamespaceByNamespaceId(id)
	if err != nil {
		logger.Logger().Errorf("GetAllGroupNamespaceByNamespaceId err:%v", err.Error())
		return nil, err
	}

	if len(groupNamespaceByNamespaceId) > 0 {
		err := repo.DeleteGroupNamespaceByNamespaceId(id)
		if err != nil {
			logger.Logger().Errorf("GetAllGroupNamespaceByNamespaceId err:%v", err.Error())
			return nil, err
		}
	}
	err = repo.DeleteNamespaceByID(id)
	if err != nil {
		logger.Logger().Errorf("DeleteNamespaceByID err:%v", err.Error())
		return nil, err
	}
	namepsace, err := repo.GetDeleteNamespaceByID(id)
	return &namepsace, err
}

func GetNamespaceByName(name string) (models.Namespace, error) {
	return repo.GetNamespaceByName(name)
}

func AddNamespace(namespace models.Namespace, namespaceRequest models.NamespaceRequest) error {
	tx := datasource.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	nsByName, err := repo.GetDeleteNamespaceByName(namespace.Namespace)
	if err != nil {
		if err := repo.AddNamespaceDB(tx, namespace); nil != err {
			tx.Rollback()
			return err
		}
	} else {
		if nsByName.Namespace == namespace.Namespace {
			namespace.ID = nsByName.ID
			if err := repo.UpdateNamespaceDB(tx, namespace); nil != err {
				tx.Rollback()
				return err
			}
		}
	}
	err = AddNSToK8s(namespace.Namespace, namespaceRequest)
	logger.Logger().Debugf("AddNamespace add ns to k8s")

	if nil != err {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateNamespace(namespace models.Namespace) (models.Namespace, error) {
	err := repo.UpdateNamespace(namespace)
	if err != nil {
		logger.Logger().Error("Update namespace err, ", err)
		return models.Namespace{}, err
	}
	return repo.GetNamespaceByName(namespace.Namespace)
}

func AddNSToK8s(namespace string, namespaceRequest models.NamespaceRequest) error {
	kac := manager.GetK8sApiClient()

	nsFromK8s, getNSErr := kac.GetNSFromK8sForAdd(namespace)

	if nil != getNSErr {
		return getNSErr
	}
	logger.Logger().Debugf("AddNamespace AddNSToK8s get ns from k8s")

	if nsFromK8s != nil {
		deleteNSFromK8sErr := kac.DeleteNSFromK8s(namespace)
		if deleteNSFromK8sErr != nil {
			return deleteNSFromK8sErr
		}
	}

	CreateNamespaceWithResourcesErrMsg := kac.CreateNamespaceWithResources(namespaceRequest.Namespace, namespaceRequest.Annotations, namespaceRequest.PlatformNamespace)

	if nil != CreateNamespaceWithResourcesErrMsg {
		return CreateNamespaceWithResourcesErrMsg
	}
	return nil
}

func ListNamespaceByRoleNameAndUserName(roleName, userName string) ([]string, error) {
	if roleName == "" || userName == "" {
		logger.Logger().Errorf("ListNamespaceByRoleNameAndUserName params can't be empty string,"+
			" roleName:%s, userName: %s\n", roleName, userName)
		return nil, errors.New("params error")
	}
	role, err := repo.GetRoleByName(roleName)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get role info by rolename, roleName: %s, err: %v\n", roleName, err)
			return nil, err
		}
		logger.Logger().Infof("role info not exist by rolename, roleName: %s, err: %v\n", roleName, err)
		return nil, nil
	}
	user, err := repo.GetUserByName(userName, datasource.GetDB())
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get user info by userName, userName: %s, err: %v\n", userName, err)
			return nil, err
		}
		logger.Logger().Errorf("user info not exist by userName, userName: %s, err: %v\n", userName, err)
		return nil, nil
	}
	groupIdList, err := repo.GetGroupIdListByUserIdRoleId(user.ID, role.ID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get group id list info by userId and roleId, "+
				"userId: %d, roleId: %d, err: %v\n", user.ID, role.ID, err)
			return nil, err
		}
		logger.Logger().Errorf("group id list info not exist by userId and roleId, "+
			"userId: %d, roleId: %d, err: %v\n", user.ID, role.ID, err)
		return nil, nil
	}
	groupNamespaceList, err := repo.GetNamespaceByGroupIdList(*groupIdList)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get group namespace list by group id list,"+
				" groupIdList: %v, err: %v\n", *groupIdList, err)
			return nil, err
		}
		logger.Logger().Errorf("group namespace list not found by group id list,"+
			" groupIdList: %v, err: %v\n", *groupIdList, err)
		return nil, nil
	}
	nsStr := make([]string, 0)
	if len(groupNamespaceList) > 0 {
		for _, ns := range groupNamespaceList {
			nsStr = append(nsStr, ns.Namespace)
		}
	}
	return nsStr, nil
}
