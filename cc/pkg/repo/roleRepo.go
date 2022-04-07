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

func GetRoles() ([]*models.Role, error) {
	var roles []*models.Role
	err := datasource.GetDB().Find(&roles, "enable_flag = ?", 1).Error
	return roles, err
}

func AddRole(role models.Role) error {
	err := datasource.GetDB().Create(&role).Error
	return err
}

func UpdateRole(role models.Role) error {
	err := datasource.GetDB().Model(&role).Omit("name").Updates(&role).Error
	return err
}

func GetRoleById(id int64) (models.Role, error) {
	var role models.Role
	err := datasource.GetDB().Find(&role, "id = ? AND enable_flag = ?", id, 1).Error
	return role, err
}

func GetDeleteRoleById(id int64) (models.Role, error) {
	var role models.Role
	err := datasource.GetDB().Find(&role, "id = ? AND enable_flag = ?", id, 0).Error
	return role, err
}

func GetRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := datasource.GetDB().Find(&role, "name = ? AND enable_flag = ?", name, 1).Error
	return role, err
}

func GetDeleteRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := datasource.GetDB().Find(&role, "name = ? AND enable_flag = ?", name, 0).Error
	return role, err
}