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
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"os"
)

func GetAllGroups(page int64, size int64) (*models.PageGroupList, error) {
	log := logger.Logger()
	var allGroups []*models.Group
	offSet := common.GetOffSet(page, size)
	total, err := repo.CountGroups()
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}
	if err != nil {
		log.Errorf("GetAllGroups Count error: %v", err.Error())
		return nil, err
	}

	if page > 0 && size > 0 {
		allGroups, err = repo.GetAllGroupsByOffset(offSet, size)
	} else {
		allGroups, err = repo.GetAllGroups()
	}
	if err != nil {
		log.Errorf("GetAllGroups error: %v", err.Error())
		return nil, err
	}

	var pageGroup = models.PageGroupList{
		Models:     allGroups,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}
	return &pageGroup, err
}

func AddGroup(group models.Group) (*models.Group, error) {
	log := logger.Logger()
	var groupByName *models.Group
	log.Infof("groupName:%v", group.Name)
	err := datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		err := repo.AddGroup(group, tx)
		if err != nil {
			log.Errorf("AddGroup Error:%v", err.Error())
			return err
		}
		groupByName, err = repo.GetGroupByName(group.Name, tx)
		if err != nil {
			log.Errorf("GetGroupByName Error:%v", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("AddGroup transaction commit error:%v", err.Error())
		return nil, err
	}
	return groupByName, nil
}

func UpdateGroupByDB(group models.Group) (*models.Group, error) {
	log := logger.Logger()

	var groupByName *models.Group
	err := datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		//err := repo.UpdateGroup(group, tx)
		m := make(map[string]interface{})
		m["cluster_name"] = group.ClusterName
		m["department_id"] = group.DepartmentID
		m["department_name"] = group.DepartmentName
		m["enable_flag"] = group.EnableFlag
		m["remarks"] = group.Remarks
		m["subsystem_id"] = group.SubsystemID
		m["subsystem_name"] = group.SubsystemName
		err := repo.UpdateGroupById(group.ID, m, tx)
		if err != nil {
			log.Errorf("UpdateGroupByDB error:%v", err.Error())
			return err
		}
		groupByName, err = repo.GetGroupByName(group.Name, tx)
		if err != nil {
			log.Errorf("GetGroupByName error:%v", err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("UpdateGroupByDB transaction error:%v", err.Error())
		return nil, err
	}
	return groupByName, nil
}

func UpdateGroup(group models.Group, db *gorm.DB) (*models.Group, error) {
	log := logger.Logger()
	//err := repo.UpdateGroup(group, db)
	//if err != nil {
	//	log.Errorf("Update Group error：%v", err.Error())
	//	return nil, err
	//}

	m := make(map[string]interface{})
	m["cluster_name"] = group.ClusterName
	m["department_id"] = group.DepartmentID
	m["department_name"] = group.DepartmentName
	m["enable_flag"] = group.EnableFlag
	m["remarks"] = group.Remarks
	m["subsystem_id"] = group.SubsystemID
	m["subsystem_name"] = group.SubsystemName

	err := repo.UpdateGroupById(group.ID, m, db)
	if err != nil {
		log.Errorf("Update Group error：%v", err.Error())
		return nil, err
	}

	name, err := repo.GetGroupByName(group.Name, db)
	if err != nil {
		log.Errorf("Update Group error：%v", err.Error())
		return nil, err
	}
	return name, nil
}

func GetGroupByGroupId(groupId int64) (*models.Group, error) {
	log := logger.Logger()
	group, err := repo.GetGroupByGroupId(groupId)
	if err != nil {
		log.Error("GetGroupByGroupId error: %v", err.Error())
		return nil, err
	}
	return group, err
}

func GetGroupIdsByUserIdAndClusterName(userId int64, clusterName string) []int64 {
	return repo.GetGroupIdsByUserIdAndClusterName(userId, clusterName)
}

func DeleteGroupByIdDB(tx *gorm.DB, groupId int64) error {
	log := logger.Logger()
	userGroupByGroupId, err := repo.GetAllUserGroupByGroupId(tx, groupId)
	if err != nil {
		return err
	}

	if len(*userGroupByGroupId) > 0 {
		if err := repo.DeleteUserGroupByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}
	groupStorageByGroupId, err := repo.GetAllGroupStorageByGroupId(tx, groupId)
	if err != nil {
		log.Errorf("GetAllGroupStorageByGroupId error: %v", err.Error())
		return err
	}

	if len(groupStorageByGroupId) > 0 {
		if err := repo.DeleteGroupStorageByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}
	groupNamespaceByGroupId, err := repo.GetAllGroupNamespaceByGroupId(tx, groupId)
	if err != nil {
		log.Errorf("GetAllGroupNamespaceByGroupId error: %v", err.Error())
		return err
	}

	if len(groupNamespaceByGroupId) > 0 {
		if err := repo.DeleteGroupNamespaceByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}

	if err := repo.DeleteGroupById(tx, groupId); nil != err {
		tx.Rollback()
		return err
	}

	return nil
}

func DeleteGroupById(groupId int64) error {
	log := logger.Logger()
	tx := datasource.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	userGroupByGroupId, err := repo.GetAllUserGroupByGroupId(tx, groupId)
	if err != nil {
		log.Errorf("GetAllUserGroupByGroupId error:%v", err.Error())
		return err
	}

	if userGroupByGroupId != nil && len(*userGroupByGroupId) > 0 {
		if err := repo.DeleteUserGroupByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}

	groupStorageByGroupId, err := repo.GetAllGroupStorageByGroupId(tx, groupId)
	if err != nil {
		log.Errorf("GetAllGroupStorageByGroupId error:%v", err.Error())
		return err
	}

	if nil != groupStorageByGroupId && len(groupStorageByGroupId) > 0 {
		if err := repo.DeleteGroupStorageByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}

	groupNamespaceByGroupId, err := repo.GetAllGroupNamespaceByGroupId(tx, groupId)
	if err != nil {
		log.Errorf("GetAllGroupNamespaceByGroupId error:%v", err.Error())
		return err
	}

	if len(groupNamespaceByGroupId) > 0 {
		if err := repo.DeleteGroupNamespaceByGroupId(tx, groupId); nil != err {
			tx.Rollback()
			return err
		}
	}

	if err := repo.DeleteGroupById(tx, groupId); nil != err {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetGroupByName(groupName string) (*models.Group, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	group, err := repo.GetGroupByName(groupName, db)
	if err != nil {
		log.Errorf("GetGroupByName error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetGroupByUsername(username string) (*models.Group, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	groupName := "gp-private-" + username
	group, err := repo.GetGroupByName(groupName, db)
	if err != nil {
		log.Errorf("GetGroupByName error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func DeleteGroupByName(groupName string) (*models.Group, error) {
	log := logger.Logger()
	group, err := repo.DeleteGroupByName(groupName)
	if err != nil {
		log.Errorf("DeleteGroupByName error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetDeleteGroupByName(groupName string) (*models.Group, error) {
	log := logger.Logger()
	group, err := repo.GetDeleteGroupByName(groupName)
	if err != nil {
		log.Errorf("GetDeleteGroupByName error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetUserGroupByUserIdAndGroupId(userId int64, groupId int64) (*models.UserGroup, error) {
	log := logger.Logger()
	group, err := repo.GetUserGroupByUserIdAndGroupId(userId, groupId)
	if err != nil {
		log.Errorf("GetUserGroupByUserIdAndGroupId error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetDeleteUserGroupByUserIdAndGroupId(userId int64, groupId int64) (*models.UserGroup, error) {
	log := logger.Logger()
	group, err := repo.GetUserGroupByUserIdAndGroupId(userId, groupId)
	if err != nil {
		log.Errorf("GetDeleteUserGroupByUserIdAndGroupId error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetUserGroupById(id int64) (*models.UserGroup, error) {
	log := logger.Logger()
	group, err := repo.GetUserGroupById(id)
	if err != nil {
		log.Errorf("GetUserGroupById error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func AddUserToGroup(userGroup models.UserGroup) (*models.UserGroup, error) {
	log := logger.Logger()
	err := repo.AddUserToGroup(userGroup)
	if err != nil {
		log.Errorf("AddUserToGroup error:%v", err.Error())
		return nil, err
	}
	group, err := repo.GetUserGroupByUserIdAndGroupId(userGroup.UserID, userGroup.GroupID)
	if err != nil {
		log.Errorf("AddUserToGroup Repo GetUserGroupByUserIdAndGroupId error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func UpdateUserGroup(userGroup models.UserGroup) (*models.UserGroup, error) {
	log := logger.Logger()
	err := repo.UpdateUserGroup(userGroup)
	if err != nil {
		log.Errorf("UpdateUserGroup error:%v", err.Error())
		return nil, err
	}
	group, err := repo.GetUserGroupById(userGroup.ID)
	if err != nil {
		log.Errorf("UpdateUserGroup Repo GetUserGroupById error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func DeleteUserFromGroupById(id int64) (*models.UserGroup, error) {
	log := logger.Logger()
	err := repo.DeleteUserFromGroupById(id)
	if err != nil {
		log.Errorf("DeleteUserFromGroupById error:%v", err.Error())
		return nil, err
	}

	group, err := repo.GetDeleteUserGroupById(id)
	if err != nil {
		log.Errorf("DeleteUserFromGroupById Repo GetDeleteUserGroupById error:%v", err.Error())
		return nil, err
	}
	return group, err
}

func GetAllUserGroupByUserId(userId int64, page int64, size int64) (*models.PageUserGroupResList, error) {
	log := logger.Logger()
	var userGroups []models.UserGroup

	//Get Page Message
	total, err := repo.CountUserGroup(userId)
	if err != nil {
		log.Errorf("Service GetAllUserGroupByUserId CountUserGroup error: %v", err.Error())
		return nil, err
	}
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	//Get User Group
	if page > 0 && size > 0 {
		offSet := common.GetOffSet(page, size)
		userGroupsFromDB, err := repo.GetAllUserGroupByUserIdAndPageAndSize(userId, offSet, size)
		if err != nil {
			log.Errorf("Service GetAllUserGroupByUserIdAndPageAndSize error: %v", err.Error())
			return nil, err
		}
		userGroups = userGroupsFromDB
	} else {
		userGroupsFromDB, err := repo.GetAllUserGroupByUserId(userId)
		if err != nil {
			log.Errorf("Service GetAllUserGroupByUserId error: %v", err.Error())
			return nil, err
		}
		userGroups = userGroupsFromDB
	}

	var ugList []*models.UserGroupRes
	if nil != userGroups && len(userGroups) > 0 {
		for _, ug := range userGroups {
			group, err := repo.GetGroupByGroupId(ug.GroupID)
			if err != nil {
				log.Errorf("GetAllUserGroupByUserId for loop GetGroupByGroupId error: %v", err.Error())
				return nil, err
			}
			user, err := repo.GetUserByUserId(ug.UserID)
			if err != nil {
				log.Errorf("GetAllUserGroupByUserId for loop GetUserByUserId error: %v", err.Error())
				return nil, err
			}
			role, err := repo.GetRoleById(ug.RoleID)
			if err != nil {
				log.Errorf("GetAllUserGroupByUserId for loop GetRoleById error: %v", err.Error())
				return nil, err
			}

			var ugs = &models.UserGroupRes{
				ID:         ug.ID,
				GroupID:    ug.GroupID,
				GroupType:  group.GroupType,
				GroupName:  group.Name,
				UserID:     ug.UserID,
				Username:   user.Name,
				RoleID:     ug.RoleID,
				RoleName:   role.Name,
				Remarks:    ug.Remarks,
				EnableFlag: ug.EnableFlag,
			}
			ugList = append(ugList, ugs)
		}
	}

	var pageUserGroup = models.PageUserGroupResList{
		Models:     ugList,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}

	return &pageUserGroup, err
}

func GetAllUserNameByGroupIds(groupIds []int64) ([]*models.User, error) {
	log := logger.Logger()
	ugs, err := repo.GetAllUserNameByGroupIds(groupIds)
	if err != nil {
		log.Errorf("GetAllUserNameByGroupIds error:%v", err.Error())
		return nil, err
	}
	userIds := common.NewIntSet()
	if len(ugs) > 0 {
		for _, ug := range ugs {
			userIds.Add(int(ug.UserID))
		}
	}

	var users []*models.User
	if nil != userIds && userIds.Len() > 0 {
		users, err = repo.GetUserByIds(userIds.List())
	}

	return users, err
}

func GetGroupIdListByUserId(userId int64) ([]int64, error) {
	log := logger.Logger()
	userGroupByUserId, err := repo.GetAllUserGroupByUserId(userId)
	if err != nil {
		log.Errorf("GetGroupIdListByUserId error:%v", err.Error())
		return nil, err
	}

	var ids []int64
	if nil != userGroupByUserId && len(userGroupByUserId) > 0 {
		for _, userGroup := range userGroupByUserId {
			ids = append(ids, userGroup.GroupID)
		}
	}
	return ids, err
}

func GetRoleIdsByUserId(userId int64) ([]int, error) {
	log := logger.Logger()
	userGroups, err := repo.GetAllUserGroupByUserId(userId)
	if err != nil {
		log.Errorf("GetRoleIdsByUserId error:%v", err.Error())
		return nil, err
	}

	ids := common.NewIntSet()
	if nil != userGroups && len(userGroups) > 0 {
		for _, ug := range userGroups {
			ids.Add(int(ug.RoleID))
		}
	}
	return ids.List(), err
}

func GetGroupStorageById(id int64) (*models.GroupStorage, error) {
	log := logger.Logger()
	group, err := repo.GetGroupStorageById(id)
	if err != nil {
		log.Errorf("GetGroupStorageById error:%v", err.Error())
		return nil, err
	}
	return &group, err
}

func AllGroupStorage() ([]models.GroupStorage, error) {
	log := logger.Logger()
	groupStorage, err := repo.GetAllGroupStorage()
	if err != nil {
		log.Errorf("AllGroupStorage error:%v", err.Error())
	}
	return groupStorage, err
}

func GetGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64) (*models.GroupStorage, error) {
	log := logger.Logger()
	groupStorage, err := repo.GetGroupStorageByStorageIdAndGroupId(storageId, groupId)
	if err != nil {
		log.Errorf("GetGroupStorageByStorageIdAndGroupId error:", err.Error())
		return nil, err
	}
	return &groupStorage, err
}

func CountGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64) (int64, error) {
	log := logger.Logger()
	count, err := repo.CountGroupStorageByStorageIdAndGroupId(storageId, groupId)
	if err != nil {
		log.Errorf("CountGroupStorageByStorageIdAndGroupId error:", err.Error())
		return 0, err
	}
	return count, err
}

func GetDeleteGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64) (*models.GroupStorage, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	groupStorage, err := repo.GetDeleteGroupStorageByStorageIdAndGroupId(storageId, groupId, db)
	if err != nil {
		log.Errorf("GetDeleteGroupStorageByStorageIdAndGroupId error:%v", err.Error())
		return nil, err
	}
	return &groupStorage, err
}

func GetGroupStorageByGroupIds(groupIds []int64) ([]models.GroupStorage, error) {
	log := logger.Logger()
	groupStorage, err := repo.GetAllGroupStorageByGroupIds(groupIds)
	if err != nil {
		log.Errorf("GetGroupStorageByGroupIds error:%v", err.Error())
		return nil, err
	}
	return groupStorage, err

}

func AddStorageToGroup(groupStorage models.GroupStorage) (*models.GroupStorage, error) {
	log := logger.Logger()
	err := repo.AddStorageToGroup(groupStorage)
	if err != nil {
		log.Errorf("AddStorageToGroup error:%v", err.Error())
		return nil, err
	}

	groupStorages, err := GetGroupStorageByStorageIdAndGroupId(groupStorage.StorageID, groupStorage.GroupID)
	if err != nil {
		log.Errorf("GetGroupStorageByStorageIdAndGroupId error:%v", err.Error())
		return nil, err
	}
	return groupStorages, err
}

func UpdateGroupStorage(groupStorage models.GroupStorage) (*models.GroupStorage, error) {
	log := logger.Logger()
	err := repo.UpdateGroupStorage(groupStorage)
	if err != nil {
		log.Errorf("UpdateGroupStorage error:%v", err.Error())
		return nil, err
	}

	groupStorages, err := repo.GetGroupStorageById(groupStorage.ID)
	if err != nil {
		log.Errorf("GetGroupStorageById error:%v", err.Error())
		return nil, err
	}
	return &groupStorages, err
}

func DeleteStorageFromGroupById(id int64) (*models.GroupStorage, error) {
	log := logger.Logger()
	err := repo.DeleteStorageFromGroupById(id)
	if err != nil {
		log.Errorf("DeleteStorageFromGroupById error:%v", err.Error())
		return nil, err
	}

	groupStorages, err := repo.GetDeleteGroupStorageById(id)
	if err != nil {
		log.Errorf("GetDeleteGroupStorageById error:%v", err.Error())
		return nil, err
	}
	if isExist, _ := common.PathExists(groupStorages.Path); isExist {
		err := os.RemoveAll(groupStorages.Path)
		if err != nil {
			return nil, err
		}
	}

	return &groupStorages, err
}

func AddNamespaceToGroup(groupNamespace models.GroupNamespace) (*models.GroupNamespace, error) {
	log := logger.Logger()
	err := repo.AddNamespaceToGroup(groupNamespace)
	if err != nil {
		log.Errorf("AddNamespaceToGroup error:%v", err.Error())
		return nil, err
	}
	groupNamespaces, err := GetGroupNamespaceByGroupIdAndNamespaceId(groupNamespace.GroupID, groupNamespace.NamespaceID)
	return groupNamespaces, err
}

func GetGroupNamespaceByGroupIdAndNamespaceId(groupId int64, namespaceId int64) (*models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespace, err := repo.GetGroupNamespaceByGroupIdAndNamespaceId(groupId, namespaceId)
	if err != nil {
		log.Errorf("GetGroupNamespaceByGroupIdAndNamespaceId error:%v", err.Error())
		return nil, err
	}
	return &groupNamespace, err
}

func GetDeleteGroupNamespaceByGroupIdAndNamespaceId(groupId int64, namespaceId int64) (*models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespace, err := repo.GetDeleteGroupNamespaceByGroupIdAndNamespaceId(groupId, namespaceId)
	if err != nil {
		log.Errorf("GetDeleteGroupNamespaceByGroupIdAndNamespaceId error:%v", err.Error())
		return nil, err
	}
	return &groupNamespace, err
}

func UpdateGroupNamespace(groupNamespace models.GroupNamespace) (*models.GroupNamespace, error) {
	log := logger.Logger()
	err := repo.UpdateGroupNamespace(groupNamespace)
	if err != nil {
		log.Errorf("UpdateGroupNamespace error:%v", err.Error())
		return nil, err
	}

	groupNamespaces, err := repo.GetGroupNamespaceById(groupNamespace.ID)
	if err != nil {
		log.Errorf("GetGroupNamespaceById error:%v", err.Error())
		return nil, err
	}
	return &groupNamespaces, err
}

func GetGroupNamespaceById(id int64) (*models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespace, err := repo.GetGroupNamespaceById(id)
	if err != nil {
		log.Errorf("GetGroupNamespaceById error:%v", err.Error())
		return nil, err
	}
	return &groupNamespace, err
}

func DeleteGroupNamespaceById(id int64) (*models.GroupNamespace, error) {
	log := logger.Logger()
	err := repo.DeleteGroupNamespaceById(id)
	if err != nil {
		log.Errorf("DeleteGroupNamespaceById error:%v", err.Error())
		return nil, err
	}
	groupNamespace, err := repo.GetGroupNamespaceById(id)
	if err != nil {
		log.Errorf("GetGroupNamespaceById error:%v", err.Error())
		return nil, err
	}
	return &groupNamespace, err
}

func GetAllGroupNamespaceByNamespaceId(namespaceId int64, page int64, size int64) (*models.PageGroupNamespaceResList, error) {
	log := logger.Logger()
	var groupNamespace []models.GroupNamespace
	total, err := repo.CountGroupNamespace(namespaceId)
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}
	if err != nil {
		log.Errorf("CountGroupNamespace error:%v", err.Error())
		return nil, err
	}

	if page > 0 && size > 0 {
		offSet := common.GetOffSet(page, size)
		groupNamespace, err = repo.GetAllGroupNamespaceByNamespaceIdAndPageAndSize(namespaceId, offSet, size)
	} else {
		groupNamespace, err = repo.GetAllGroupNamespaceByNamespaceId(namespaceId)
	}
	if err != nil {
		log.Errorf("GetAllGroupNamespaceByNamespaceId error:%v", err.Error())
		return nil, err
	}

	var gnrList []*models.GroupNamespaceRes
	for _, gn := range groupNamespace {
		group, err := repo.GetGroupByGroupId(gn.GroupID)
		if err != nil {
			log.Errorf("GetAllGroupNamespaceByNamespaceId for loop GetGroupByGroupId error:%v", err.Error())
			return nil, err
		}
		var gnr = &models.GroupNamespaceRes{
			ID:          gn.ID,
			GroupID:     gn.GroupID,
			GroupName:   group.Name,
			GroupType:   group.GroupType,
			NamespaceID: gn.NamespaceID,
			Namespace:   gn.Namespace,
			Remarks:     gn.Remarks,
			EnableFlag:  gn.EnableFlag,
		}
		gnrList = append(gnrList, gnr)
	}

	var pageGroupNamespace = models.PageGroupNamespaceResList{
		Models:     gnrList,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}

	return &pageGroupNamespace, err
}

func GetAllGroupNamespaceByGroupId(groupId int64) ([]models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespaces, err := repo.GetAllGroupNamespaceByGroupId(datasource.GetDB(), groupId)
	if err != nil {
		log.Errorf("GetAllGroupNamespaceByGroupId error:%v", err.Error())
		return nil, err
	}
	return groupNamespaces, err
}

func GetAllGroupNamespaceByGroupIds(groupIds []int64) ([]models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespaces, err := repo.GetAllGroupNamespaceByGroupIds(groupIds)
	if err != nil {
		log.Errorf("GetAllGroupNamespaceByGroupId error:%v", err.Error())
		return nil, err
	}
	return groupNamespaces, err
}

func GetAllNamespaceByUserId(userId int64) ([]models.GroupNamespace, error) {
	log := logger.Logger()
	userGroupsList, err := GetAllUserGroupByUserId(userId, 0, 0)
	if err != nil {
		log.Errorf("GetAllUserGroupByUserId error:%v", err.Error())
		return nil, err
	}
	userGroups := userGroupsList.Models
	var groupNamespaceList []models.GroupNamespace
	if len(userGroups) > 0 {
		for _, userGroup := range userGroups {
			groupNamespaces, err := repo.GetAllGroupNamespaceByGroupId(datasource.GetDB(), userGroup.GroupID)
			if err != nil {
				log.Errorf("GetAllGroupNamespaceByGroupId error:%v", err.Error())
			}
			groupNamespaceList = append(groupNamespaceList, groupNamespaces...)
		}
	}
	return groupNamespaceList, nil
}

func CheckAdmin(namespace string, admin string, username string) error {
	log := logger.Logger()
	db := datasource.GetDB()
	if admin == username {
		return nil
	}
	_, err := repo.GetSAByName(admin)
	if err != nil {
		getNamespace, err := repo.GetNamespaceByName(namespace)
		if err != nil {
			logger.Logger().Error("Get namespace by name err, ", err)
			return err
		}
		if namespace != getNamespace.Namespace {
			return errors.New("namespace is not exists in db")
		}

		groupNamespaces, err := repo.GetAllGroupNamespaceByNamespaceId(getNamespace.ID)
		if err != nil {
			log.Errorf("GetAllGroupNamespaceByNamespaceId error:%v", err.Error())
			return nil
		}

		if len(groupNamespaces) < 1 {
			return errors.New("namespace is not exists in any group")
		}

		userByAdmin, err := repo.GetUserByName(admin, db)
		if err != nil {
			log.Errorf("repo.GetUserByName for userByAdmin, %v", err.Error())
			return err
		}
		userByName, err := repo.GetUserByName(username, db)
		if err != nil {
			log.Errorf("repo.GetUserByName for userByName, %v", err.Error())
			log.Errorf(err.Error())
			return err
		}
		if admin != userByAdmin.Name && username != userByName.Name {
			return errors.New("checkAdmin, user does not existed")
		}

		roleByName, err := GetRoleByName(constants.ADMIN)
		if err != nil {
			logger.Logger().Error("GetRoles err, ", err)
			return err
		}
		if constants.ADMIN != roleByName.Name {
			return errors.New("role:ADMIN does not existed in db")
		}

		var adminGroups []models.UserGroup
		var userGroups []models.UserGroup

		for _, gn := range groupNamespaces {
			adminGroup, err := repo.GetUserGroupByUserIdAndGroupId(userByAdmin.ID, gn.GroupID)
			if err != nil {
				log.Errorf("CheckAdmin for loop GetUserGroupByUserIdAndGroupId error: %v", err.Error())
				return err
			}
			if userByAdmin.ID == adminGroup.UserID && gn.GroupID == adminGroup.GroupID {
				adminGroups = append(adminGroups, *adminGroup)
			}

			userGroup, err := repo.GetUserGroupByUserIdAndGroupId(userByName.ID, gn.GroupID)
			if err != nil {
				log.Errorf("CheckAdmin for loop GetUserGroupByUserIdAndGroupId error: %v", err.Error())
				return err
			}
			if userByName.ID == userGroup.UserID && gn.GroupID == userGroup.GroupID {
				userGroups = append(userGroups, *userGroup)
			}
		}

		if nil == adminGroups || len(adminGroups) < 1 {
			return errors.New("checkAdmin, adminName is not ga in this group")
		}

		var isGA = false
		for _, ug := range adminGroups {
			if ug.RoleID == roleByName.ID {
				isGA = true
				break
			}
		}

		if !isGA {
			return errors.New("checkAdmin, adminName is not ga in this group")
		}

		if nil == userGroups || len(userGroups) < 1 {
			return errors.New("checkAdmin, userName not in this group")
		}
	}
	return nil
}

func GetAllGroupNamespace() ([]models.GroupNamespace, error) {
	log := logger.Logger()
	groupNamespaces, err := repo.GetAllGroupNamespace()
	if err != nil {
		log.Errorf("repo.GetAllGroupNamespace error: %v", err.Error())
		return nil, err
	}
	return groupNamespaces, err
}

func GetGroupIdListByUserIdRoleId(userId int64, roleId int64) ([]int64, error) {
	log := logger.Logger()
	groupIdList, err := repo.GetGroupIdListByUserIdRoleId(userId, roleId)
	if err != nil {
		log.Errorf("GetGroupIdListByUserIdRoleId error: %v", err.Error())
		return nil, err
	}
	return *groupIdList, err
}

func GetAllGroupStorageByStorageId(storageId int64, page int64, size int64) (*models.PageGroupStorageResList, error) {
	log := logger.Logger()
	tx := datasource.GetDB()
	total, err := repo.CountGroupStorage(storageId)
	if err != nil {
		log.Errorf("CountGroupStorage error:%v", err.Error())
		return nil, err
	}
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	var groupStorage []models.GroupStorage
	if page > 0 && size > 0 {
		offSet := common.GetOffSet(page, size)
		groupStorage, err = repo.GetAllGroupStorageByStorageIdAndPageAndSize(tx, storageId, offSet, size)
	} else {
		groupStorage, err = repo.GetAllGroupStorageByStorageId(tx, storageId)
	}
	if err != nil {
		log.Errorf("GetAllGroupStorageByStorageId error:%v", err.Error())
		return nil, err
	}

	var gsrList []*models.GroupStorageRes

	for _, gs := range groupStorage {
		group, err := repo.GetGroupByGroupId(gs.GroupID)
		if err != nil {
			log.Errorf("GetAllGroupStorageByStorageId for loop GetGroupByGroupId error:%v", err.Error())
			return nil, err
		}

		var gsr = &models.GroupStorageRes{
			ID:          gs.ID,
			GroupID:     gs.GroupID,
			GroupName:   group.Name,
			GroupType:   group.GroupType,
			Path:        gs.Path,
			StorageID:   gs.StorageID,
			Remarks:     gs.Remarks,
			Permissions: gs.Permissions,
			EnableFlag:  gs.EnableFlag,
			Type:        gs.Type,
		}
		gsrList = append(gsrList, gsr)
	}

	var pageGroupStorage = &models.PageGroupStorageResList{
		Models:     gsrList,
		Total:      total,
		TotalPage:  int64(totalPage),
		PageNumber: page,
		PageSize:   size,
	}

	return pageGroupStorage, nil
}

func GetGroups(userId int64) []*models.Group {
	userGroups, _ := repo.GetAllUserGroupByUserId(userId)
	var groups []*models.Group
	if userGroups != nil && len(userGroups) > 0 {
		for _, userGroup := range userGroups {
			getUserGroup, _ := repo.GetGroupByGroupId(userGroup.GroupID)
			groups = append(groups, getUserGroup)
		}
	}
	return groups
}

func GetAllGroupsNoPage() ([]*models.Group, error) {
	//var groups []*models.Group
	log := logger.Logger()
	groups, err := repo.GetAllGroups()
	if err != nil {
		log.Errorf("GetAllGroupsNoPage error:%v", err.Error())
		return nil, err
	}
	return groups, nil
}
