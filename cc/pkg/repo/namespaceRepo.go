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
	"fmt"
	"github.com/jinzhu/gorm"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
)

func GetAllNamespacesByOffset(offSet int64, size int64) ([]*models.Namespace, error) {
	var namespaces []*models.Namespace
	err := datasource.GetDB().Offset(offSet).Limit(size).Table("t_namespace").Find(&namespaces, "enable_flag = ?", 1).Error
	return namespaces, err
}

func CountNamespaces() (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_namespace").Where("enable_flag = ?", 1).Count(&count).Error
	return count, err
}

func GetAllNamespaces() ([]*models.Namespace, error) {
	var namespaces []*models.Namespace
	err := datasource.GetDB().Table("t_namespace").Find(&namespaces, "enable_flag = ?", 1).Error
	return namespaces, err
}

//func CountNamespaces() int64 {
//	var count int64
//	datasource.GetDB().Table("t_namespace").Where("enable_flag = ?", 1).Count(&count)
//	return count
//}

func GetNamespaceByID(id int64) (models.Namespace, error) {
	var namespace models.Namespace
	err := datasource.GetDB().Find(&namespace, "id = ? AND enable_flag = ?", id, 1).Error
	return namespace, err
}

func GetDeleteNamespaceByID(id int64) (models.Namespace, error) {
	var namespace models.Namespace
	err := datasource.GetDB().Find(&namespace, "id = ?", id).Error
	return namespace, err
}

func DeleteNamespaceByID(id int64) error {
	return datasource.GetDB().Table("t_namespace").Where("id = ?", id).Update("enable_flag", 0).Error
}

func GetNamespaceByName(name string) (models.Namespace, error) {
	var namespace models.Namespace
	err := datasource.GetDB().Find(&namespace, "namespace = ? AND enable_flag = ?", name, 1).Error
	return namespace, err
}

func GetDeleteNamespaceByName(name string) (models.Namespace, error) {
	var namespace models.Namespace
	err := datasource.GetDB().Find(&namespace, "namespace = ? AND enable_flag = ?", name, 0).Error
	return namespace, err
}

func AddNamespaceDB(db *gorm.DB, namespace models.Namespace) error {
	return db.Create(&namespace).Error
}

func UpdateNamespace(namespace models.Namespace) error {
	return datasource.GetDB().Model(&namespace).Omit("namespace").Updates(&namespace).Error
}

func UpdateNamespaceDB(db *gorm.DB, namespace models.Namespace) error {
	return db.Model(&namespace).Omit("namespace").Updates(&namespace).Error
}

func GetNamespaceByGroupIdList(groupIdList []int64) ([]*models.GroupNamespace, error) {
	var namespaces []*models.GroupNamespace
	subQuery := ""
	if len(groupIdList) > 0 {
		for _, id := range groupIdList {
			subQuery += fmt.Sprintf("%d,", id)
		}
	}
	if subQuery != "" {
		subQuery = subQuery[0 : len(subQuery)-1]
		subQuery = " AND t_group_namespace.group_id IN (" + subQuery + ")"
	}
	logger.Logger().Infof("GetNamespaceByGroupIdList subQuery: %v\n", subQuery)
	err := datasource.GetDB().Table("t_group_namespace").
		Where("t_group_namespace.enable_flag = ?"+subQuery, 1).
		Find(&namespaces).Error
	return namespaces, err
}
