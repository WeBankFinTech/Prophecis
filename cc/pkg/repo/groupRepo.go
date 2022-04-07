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
package repo

import (
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"

	"github.com/jinzhu/gorm"
)

type GroupStorage struct {
	EnableFlag  int64  `json:"enableFlag,omitempty" gorm:"enable_flag"`
	GroupID     int64  `json:"groupId" gorm:"group_id"`
	ID          int64  `json:"id,omitempty" gorm:"id"`
	Path        string `json:"path,omitempty" gorm:"path"`
	Permissions string `json:"permissions,omitempty" gorm:"permissions"`
	Remarks     string `json:"remarks,omitempty" gorm:"remarks"`
	StorageID   int64  `json:"storageId,omitempty" gorm:"storage_id"`
	Type        string `json:"type,omitempty" gorm:"type"`
	//GroupType string `json:"groupType,omitempty"`
}

func GetAllGroupsByOffset(offSet int64, size int64) ([]*models.Group, error) {
	var groups []*models.Group
	err := datasource.GetDB().Offset(offSet).Limit(size).Table("t_group").Find(&groups, "enable_flag = ?", 1).Error
	return groups, err
}

func GetAllGroups() ([]*models.Group, error) {
	var groups []*models.Group
	err := datasource.GetDB().Table("t_group").Find(&groups, "enable_flag = ?", 1).Error
	return groups, err
}

func CountGroups() (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_group").Where("enable_flag = ?", 1).Count(&count).Error
	return count, err
}

func AddGroup(group models.Group, db *gorm.DB) error {
	return db.Create(&group).Error
}

func UpdateGroup(group models.Group, db *gorm.DB) error {
	return db.Model(&group).Omit("name", "group_type").Updates(&group).Error
}

func UpdateGroupById(groupId int64, m map[string]interface{}, db *gorm.DB) error {
	return datasource.GetDB().Table("t_group").Where("id = ? AND enable_flag = ?", groupId, 1).Update(m).Error
}

func GetGroupByName(groupName string, db *gorm.DB) (*models.Group, error) {
	var group models.Group
	err := db.Find(&group, "name = ? AND enable_flag = ?", groupName, 1).Error
	return &group, err
}

func DeleteGroupByName(groupName string) (*models.Group, error) {
	var group models.Group
	err := datasource.GetDB().Table("t_group").Where("name = ?", groupName).Update("enable_flag", 0).Error
	if err != nil {
		return nil, err
	}
	err = datasource.GetDB().Table("t_group").Find(&group, "name = ?", groupName).Error
	return &group, err
}

func GetDeleteGroupByName(groupName string) (*models.Group, error) {
	var group models.Group
	err := datasource.GetDB().Find(&group, "name = ? AND enable_flag = ?", groupName, 0).Error
	return &group, err
}

func GetGroupByGroupId(groupId int64) (*models.Group, error) {
	var group models.Group
	err := datasource.GetDB().Find(&group, "id = ? AND enable_flag = ?", groupId, 1).Error
	return &group, err
}

func DeleteGroupById(db *gorm.DB, groupId int64) error {
	return db.Table("t_group").Where("id = ?", groupId).Update("enable_flag", 0).Error
}

func AddUserToGroupDB(db *gorm.DB, userGroup models.UserGroup) error {
	return db.Create(&userGroup).Error
}

func AddUserToGroup(userGroup models.UserGroup) error {
	return datasource.GetDB().Create(&userGroup).Error
}

func UpdateUserGroupDB(db *gorm.DB, userGroup models.UserGroup) error {
	return db.Model(&userGroup).Omit("user_id", "group_id").Updates(&userGroup).Error
}

func UpdateUserGroup(userGroup models.UserGroup) error {
	return datasource.GetDB().Model(&userGroup).Omit("user_id", "group_id").Updates(&userGroup).Error
}

func GetUserGroupById(id int64) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	err := datasource.GetDB().Find(&userGroup, "id = ? and enable_flag = ?", id, 1).Error
	return &userGroup, err
}

func GetDeleteUserGroupById(id int64) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	err := datasource.GetDB().Find(&userGroup, "id = ? AND enable_flag = ?", id, 0).Error
	return &userGroup, err
}

func GetUserGroupByUserIdAndGroupId(userId int64, groupId int64) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	err := datasource.GetDB().Find(&userGroup, "user_id = ? AND group_id = ?", userId, groupId).Error
	return &userGroup, err
}

func GetDeleteUserGroupByUserIdAndGroupId(userId int64, groupId int64) (*models.UserGroup, error) {
	var userGroup models.UserGroup
	err := datasource.GetDB().Find(&userGroup, "user_id = ? AND group_id = ? AND enable_flag = ?", userId, groupId, 0).Error
	return &userGroup, err
}

func GetAllUserGroupByUserIdAndPageAndSize(userId int64, offSet int64, size int64) ([]models.UserGroup, error) {
	var userGroups []models.UserGroup
	err := datasource.GetDB().Offset(offSet).Limit(size).Find(&userGroups, "user_id = ? AND enable_flag = ?", userId, 1).Error
	return userGroups, err
}

func GetAllUserGroupByUserId(userId int64) ([]models.UserGroup, error) {
	var userGroups []models.UserGroup
	err := datasource.GetDB().Find(&userGroups, "user_id = ? AND enable_flag = ?", userId, 1).Error
	return userGroups, err
}

func CountUserGroup(userId int64) (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_user_group").Where("user_id = ? AND enable_flag = ?", userId, 1).Count(&count).Error
	return count, err
}

func GetAllUserNameByGroupIds(groupIds []int64) ([]models.UserGroup, error) {
	var ug []models.UserGroup
	err := datasource.GetDB().Find(&ug, "enable_flag = ? AND group_id IN (?)", 1, groupIds).Error
	return ug, err
}

func GetAllUserGroupByGroupId(db *gorm.DB, groupId int64) (*[]models.UserGroup, error) {
	var userGroups []models.UserGroup
	err := db.Find(&userGroups, "group_id = ? AND enable_flag = ?", groupId, 1).Error
	if err != nil {
		return nil, err
	}
	return &userGroups, err
}

func DeleteUserFromGroupById(id int64) error {
	return datasource.GetDB().Table("t_user_group").Where("id = ?", id).Update("enable_flag", 0).Error
}

func DeleteUserGroupByUserId(db *gorm.DB, userId int64) error {
	return db.Table("t_user_group").Where("user_id = ?", userId).Update("enable_flag", 0).Error
}

func DeleteUserGroupByGroupId(db *gorm.DB, groupId int64) error {
	return db.Table("t_user_group").Where("group_id = ?", groupId).Update("enable_flag", 0).Error
}

func GetGroupStorageById(id int64) (models.GroupStorage, error) {
	var groupStorage models.GroupStorage
	err := datasource.GetDB().Find(&groupStorage, "id = ? AND enable_flag = ?", id, 1).Error
	return groupStorage, err
}

func GetAllGroupStorage() ([]models.GroupStorage, error) {
	var groupStorage []models.GroupStorage
	err := datasource.GetDB().Find(&groupStorage, "enable_flag = ?", 1).Error
	return groupStorage, err
}

func GetDeleteGroupStorageById(id int64) (models.GroupStorage, error) {
	var groupStorage models.GroupStorage
	err := datasource.GetDB().Find(&groupStorage, "id = ? AND enable_flag = ?", id, 0).Error
	return groupStorage, err
}

func GetGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64) (models.GroupStorage, error) {
	var groupStorage models.GroupStorage
	err := datasource.GetDB().Find(&groupStorage, "storage_id = ? AND group_id = ?", storageId, groupId).Error
	return groupStorage, err
}

func CountGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64) (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_group_storage").Where("storage_id = ? AND group_id = ? AND enable_flag = ?",
		storageId, groupId, 1).Count(&count).Error
	return count, err
}

func GetDeleteGroupStorageByStorageIdAndGroupId(storageId int64, groupId int64, db *gorm.DB) (models.GroupStorage, error) {
	var groupStorage models.GroupStorage
	err := db.Find(&groupStorage, "storage_id = ? AND group_id = ? AND enable_flag = ?", storageId, groupId, 0).Error
	return groupStorage, err
}

func AddStorageToGroupDB(db *gorm.DB, groupStorage models.GroupStorage) error {
	logger.Logger().Debugf("AddStorageToGroupDB groupStorage: %+v", groupStorage)
	obj := GroupStorage{
		EnableFlag:  groupStorage.EnableFlag,
		GroupID:     groupStorage.GroupID,
		Path:        groupStorage.Path,
		Permissions: groupStorage.Permissions,
		Remarks:     groupStorage.Remarks,
		StorageID:   groupStorage.StorageID,
		Type:        groupStorage.Type,
	}
	return db.Create(&obj).Error
}

func AddStorageToGroup(groupStorage models.GroupStorage) error {
	obj := GroupStorage{
		EnableFlag:  groupStorage.EnableFlag,
		GroupID:     groupStorage.GroupID,
		Path:        groupStorage.Path,
		Permissions: groupStorage.Permissions,
		Remarks:     groupStorage.Remarks,
		StorageID:   groupStorage.StorageID,
		Type:        groupStorage.Type,
	}
	logger.Logger().Debugf("AddStorageToGroup obj: %+v", obj)
	return datasource.GetDB().Table("t_group_storage").Create(&obj).Error
}

func UpdateGroupStorageDB(db *gorm.DB, groupStorage models.GroupStorage) error {
	return db.Model(&groupStorage).Omit("storage_id", "group_id").Updates(&groupStorage).Error
}

func UpdateGroupStorage(groupStorage models.GroupStorage) error {
	return datasource.GetDB().Model(&groupStorage).Omit("storage_id", "group_id").Updates(&groupStorage).Error
}

func DeleteStorageFromGroupById(id int64) error {
	return datasource.GetDB().Table("t_group_storage").Where("id = ?", id).Update("enable_flag", 0).Error
}

func DeleteGroupStorageByStorageId(db *gorm.DB, storageId int64) error {
	return db.Table("t_group_storage").Where("storage_id = ?", storageId).Update("enable_flag", 0).Error
}

func DeleteGroupStorageByGroupId(db *gorm.DB, groupId int64) error {
	return db.Table("t_group_storage").Where("group_id = ?", groupId).Update("enable_flag", 0).Error
}

func GetAllGroupStorageByStorageIdAndPageAndSize(db *gorm.DB, storageId int64, offSet int64, size int64) ([]models.GroupStorage, error) {
	var groupStorageList []models.GroupStorage
	err := db.Offset(offSet).Limit(size).Find(&groupStorageList, "storage_id = ? AND enable_flag = ?", storageId, 1).Error
	return groupStorageList, err
}

func GetAllGroupStorageByStorageId(db *gorm.DB, storageId int64) ([]models.GroupStorage, error) {
	var groupStorageList []models.GroupStorage
	err := db.Find(&groupStorageList, "storage_id = ? AND enable_flag = ?", storageId, 1).Error
	return groupStorageList, err
}

func GetAllGroupStorageByGroupId(db *gorm.DB, groupId int64) ([]models.GroupStorage, error) {
	var groupStorageList []models.GroupStorage
	err := db.Find(&groupStorageList, "group_id = ? AND enable_flag = ?", groupId, 1).Error
	return groupStorageList, err
}

func CountGroupStorage(storageId int64) (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_group_storage").Where("storage_id = ? AND enable_flag = ?", storageId, 1).Count(&count).Error
	return count, err
}

func GetAllGroupStorageByGroupIds(groupIds []int64) ([]models.GroupStorage, error) {
	var groupStorageList []models.GroupStorage
	err := datasource.GetDB().Find(&groupStorageList, "enable_flag = ? AND group_id IN (?)", 1, groupIds).Error
	return groupStorageList, err
}

func AddNamespaceToGroup(groupNamespace models.GroupNamespace) error {
	return datasource.GetDB().Create(&groupNamespace).Error
}

func GetGroupNamespaceById(id int64) (models.GroupNamespace, error) {
	var groupNamespace models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespace, "id = ? AND enable_flag = ?", id, 1).Error
	return groupNamespace, err
}

func DeleteGroupNamespaceById(id int64) error {
	return datasource.GetDB().Table("t_group_namespace").Where("id = ?", id).Update("enable_flag", 0).Error
}

func DeleteGroupNamespaceByNamespaceId(id int64) error {
	return datasource.GetDB().Table("t_group_namespace").Where("namespace_id = ?", id).Update("enable_flag", 0).Error
}

func DeleteGroupNamespaceByGroupId(db *gorm.DB, id int64) error {
	return db.Table("t_group_namespace").Where("group_id = ?", id).Update("enable_flag", 0).Error
}

func GetGroupNamespaceByGroupIdAndNamespaceId(groupId int64, namespaceId int64) (models.GroupNamespace, error) {
	var groupNamespace models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespace, "group_id = ? AND namespace_id = ? AND enable_flag = ?", groupId, namespaceId, 1).Error
	return groupNamespace, err
}

func GetDeleteGroupNamespaceByGroupIdAndNamespaceId(groupId int64, namespaceId int64) (models.GroupNamespace, error) {
	var groupNamespace models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespace, "group_id = ? AND namespace_id = ? AND enable_flag = ?", groupId, namespaceId, 0).Error
	return groupNamespace, err
}

func UpdateGroupNamespace(groupNamespace models.GroupNamespace) error {
	return datasource.GetDB().Model(&groupNamespace).Omit("group_id", "namespace_id").Update(&groupNamespace).Error
}

func GetAllGroupNamespaceByNamespaceIdAndPageAndSize(namespaceId int64, offSet int64, size int64) ([]models.GroupNamespace, error) {
	var groupNamespaceList []models.GroupNamespace
	err := datasource.GetDB().Offset(offSet).Limit(size).Find(&groupNamespaceList, "namespace_id = ? AND enable_flag = ?", namespaceId, 1).Error
	return groupNamespaceList, err
}

func GetAllGroupNamespaceByNamespaceId(namespaceId int64) ([]models.GroupNamespace, error) {
	var groupNamespaceList []models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespaceList, "namespace_id = ? AND enable_flag = ?", namespaceId, 1).Error
	return groupNamespaceList, err
}

func CountGroupNamespace(namespaceId int64) (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_group_namespace").Where("namespace_id = ? AND enable_flag = ?", namespaceId, 1).Count(&count).Error
	return count, err
}

func GetAllGroupNamespaceByGroupId(db *gorm.DB, groupId int64) ([]models.GroupNamespace, error) {
	var groupNamespaceList []models.GroupNamespace
	err := db.Find(&groupNamespaceList, "group_id = ? AND enable_flag = ?", groupId, 1).Error
	return groupNamespaceList, err
}

func GetAllGroupNamespaceByGroupIds(groupIds []int64) ([]models.GroupNamespace, error) {
	var groupNamespaceList []models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespaceList, "enable_flag = ? AND group_id IN (?)", 1, groupIds).Error
	return groupNamespaceList, err
}

func GetAllGroupIds() (*[]int64, error) {
	var groups []models.Group
	err := datasource.GetDB().Table("t_group").Find(&groups, "enable_flag = ?", 1).Error
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, group := range groups {
		ids = append(ids, group.ID)
	}
	return &ids, err
}

func GetGroupIdsByUserIds(userIds []int64) (*[]int64, error) {
	var userGroups []models.UserGroup
	err := datasource.GetDB().Find(&userGroups, "enable_flag = ? AND user_id IN (?)", 1, userIds).Error
	if err != nil {
		return nil, err
	}
	var groupIds []int64
	for _, ug := range userGroups {
		groupIds = append(groupIds, ug.GroupID)
	}
	return &groupIds, err
}

func GetGroupIdsByUserId(userId int64) (*[]int64, error) {
	var userGroups []models.UserGroup
	err := datasource.GetDB().Table("t_user_group").Find(&userGroups, "user_id = ? AND enable_flag = ?",
		userId, 1).Error
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, userGroup := range userGroups {
		ids = append(ids, userGroup.GroupID)
	}
	return &ids, err
}

func GetMatchGroupNum(adminId int64, userId int64) (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_user_group").
		Joins("join t_user_group ug on t_user_group.group_id = ug.group_id").
		Where("t_user_group.role_id = ? AND t_user_group.user_id = ? AND ug.user_id = ?",
			1, adminId, userId).Count(&count).Error
	return count, err
}

func GetAllGroupNamespace() ([]models.GroupNamespace, error) {
	var groupNamespaces []models.GroupNamespace
	err := datasource.GetDB().Find(&groupNamespaces, "enable_flag = ?", 1).Error
	return groupNamespaces, err
}

func GetGroupIdListByUserIdRoleId(userId int64, roleId int64) (*[]int64, error) {
	var userGroups []models.UserGroup
	err := datasource.GetDB().Table("t_user_group").Find(&userGroups,
		"user_id = ? AND role_id = ? AND enable_flag = ?", userId, roleId, 1).Error
	if err != nil {
		return nil, err
	}
	var ids []int64
	for _, userGroup := range userGroups {
		ids = append(ids, userGroup.GroupID)
	}
	return &ids, err
}

func GetUserGroupByUserId(userId string) models.UserGroup {
	var userGroup models.UserGroup
	datasource.GetDB().Find(&userGroup, "user_id = ? AND enable_flag = ?", userId, 1)
	logger.Logger().Info("userGroup is ", userGroup)
	return userGroup
}

func GetGroupIdListByUserIdRoleIdAndClusterName(userId int64, roleId int64, clusterName string) []int64 {
	var userGroups []models.UserGroup
	db := datasource.GetDB()
	if "" == clusterName {
		db.Raw("select tug.* from t_user_group tug LEFT JOIN t_group tg ON tug.group_id = tg.id where tug.user_id = ? and tug.role_id = ? and tug.enable_flag = ? and tg.enable_flag = ?", userId, roleId, 1, 1).Scan(&userGroups)
	} else {
		db.Raw("select tug.* from t_user_group tug LEFT JOIN t_group tg ON tug.group_id = tg.id where tug.user_id = ? and tug.role_id = ? and tg.cluster_name = ? and tug.enable_flag = ? and tg.enable_flag = ?", userId, roleId, clusterName, 1, 1).Scan(&userGroups)
	}
	var ids []int64
	for _, userGroup := range userGroups {
		ids = append(ids, userGroup.GroupID)
	}
	return ids
}

func GetAllUsersByClusterName(clusterName string) []*models.User {
	var users []*models.User
	db := datasource.GetDB()
	if "" != clusterName {
		db.Raw("select * from t_user where enable_flag = 1 and id IN (select DISTINCT(user_id) from t_group tg right JOIN t_user_group tug on tg.id = tug.group_id where tg.cluster_name = ? and tg.enable_flag = 1 and tug.enable_flag = 1)", clusterName).Scan(&users)
	} else {
		//db.Table("t_user tu").Select([]string{"tu.id"}).Joins("left join t_user_group tug on tu.id = tug.user_id ANG enable_flag = ?", 1).Joins("left join t_group tg on")
		db.Raw("select * from t_user where enable_flag = 1 and id IN (select DISTINCT(user_id) from t_group tg right JOIN t_user_group tug on tg.id = tug.group_id where tg.enable_flag = 1 and tug.enable_flag = 1)").Scan(&users)
	}
	return users
}

func GetGroupIdsByUserIdAndClusterName(userId int64, clusterName string) []int64 {
	var groups []models.Group
	var ids []int64
	datasource.GetDB().Raw("select tg.id from t_group tg LEFT JOIN t_user_group tug on tg.id = tug.group_id where tg.cluster_name = ? AND tg.enable_flag = 1 AND tug.user_id = ? AND tug.enable_flag = 1", clusterName, userId).Scan(&groups)
	if len(groups) > 0 {
		for _, group := range groups {
			ids = append(ids, group.ID)
		}
	}
	return ids
}

