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

package dao

import (
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/models"
)

type UserDao struct {
}

func (s *UserDao) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := common.GetDB().Table("t_user").Where("name = ?", name).Find(&user).Error
	return &user, err
}

func (s *UserDao) GetUserById(id int64) (*gormmodels.ImageUser, error) {
	var user gormmodels.ImageUser
	err := common.GetDB().Table("t_user").Where("id = ?", id).Find(&user).Error
	return &user, err
}

func (s *UserDao) GetUserByUserId(id int64) (*gormmodels.User, error) {
	var user gormmodels.User
	err := common.GetDB().Table("t_user").Where("id = ?", id).First(&user).Error
	return &user, err
}

func (s *UserDao) GetNamespaceByName(name string) (*models.Namespace, error) {
	var namespaceModel models.Namespace
	err := common.GetDB().Table("t_group_namespace").Select("id, group_id, namespace").Where("namespace = ?", name).Find(&namespaceModel).Error
	return &namespaceModel, err
}
