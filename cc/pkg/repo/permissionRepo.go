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
	"mlss-controlcenter-go/pkg/models"
)

func GetPermissions() ([]*models.Permission, error) {
	var pms []*models.Permission
	err := datasource.GetDB().Find(&pms, "enable_flag = ?", 1).Error
	return pms, err
}

func GetPermissionIdsByRoleId(ids []int) ([]models.RolePermission, error) {
	var rps []models.RolePermission
	err := datasource.GetDB().Find(&rps, "enable_flag = ? AND role_id IN (?)", 1, ids).Error
	return rps, err
}

func GetPermissionByRoleId(id int) ([]models.RolePermission, error) {
	var rps []models.RolePermission
	err := datasource.GetDB().Find(&rps, "enable_flag = ? AND role_id = ?", 1, id).Error
	return rps, err
}

func GetPermissionsByLike(queryStr string) ([]*models.Permission, error) {
	var rps []*models.Permission
	err := datasource.GetDB().Find(&rps, "enable_flag = ? AND url LIKE ?", 1, "%"+queryStr+"%").Error
	return rps, err
}
