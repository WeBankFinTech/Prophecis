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
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
)

func GetPermissions() ([]*models.Permission, error) {
	return repo.GetPermissions()
}

func GetPermissionsByLike(queryStr string) ([]*models.Permission, error) {
	return repo.GetPermissionsByLike(queryStr)
}

func GetPermissionIdsByRoleId(ids []int) (common.IntSet, error) {
	rolePermissions, err := repo.GetPermissionIdsByRoleId(ids)
	if err != nil {
		logger.Logger().Error("Get permission ids by roleId err, ", err)
		return common.IntSet{}, err
	}
	perIds := *common.NewIntSet()
	for _, rp := range rolePermissions {
		perIds.Add(int(rp.PermissionID))
	}
	return perIds, nil
}

func GetPermissionsByRoleId(id int) (common.IntSet, error) {
	rolePermissions, err := repo.GetPermissionByRoleId(id)
	if err != nil {
		logger.Logger().Error("Get permission ids by roleId err, ", err)
		return common.IntSet{}, err
	}
	perIds := *common.NewIntSet()
	for _, rp := range rolePermissions {
		perIds.Add(int(rp.PermissionID))
	}
	return perIds, nil
}
