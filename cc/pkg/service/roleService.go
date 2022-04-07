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
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
)

func GetRoles() ([]*models.Role, error) {
	return repo.GetRoles()
}

func AddRole(role models.Role) (models.Role, error) {
	err := repo.AddRole(role)
	if err != nil {
		logger.Logger().Error("Add role err, ", err)
		return models.Role{}, err
	}
	return repo.GetRoleByName(role.Name)
}

func UpdateRole(role models.Role) (models.Role, error) {
	err := repo.UpdateRole(role)
	if err != nil {
		logger.Logger().Error("Update role err, ", err)
		return models.Role{}, err
	}
	return repo.GetRoleById(role.ID)
}

func GetRoleByName(name string) (models.Role, error) {
	return repo.GetRoleByName(name)
}

func GetDeleteRoleByName(name string) (models.Role, error) {
	return repo.GetDeleteRoleByName(name)
}

func GetRoleById(id int64) (models.Role, error) {
	return repo.GetRoleById(id)
}

func GetDeleteRoleById(id int64) (models.Role, error) {
	return repo.GetDeleteRoleById(id)
}