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
)

type GroupDao struct {
}

func (s *GroupDao) GetGroupsByUserId(userId int64) ([]*gormmodels.UserGroup, error) {
	var groups []*gormmodels.UserGroup
	err := common.GetDB().Table("t_user_group").Where("user_id = ?", userId).
		Preload("User").
		Preload("Role").
		Preload("Group").
		Find(&groups).
		Error
	return groups, err
}

func (s *GroupDao) GetGroupByGroupId(id int64) (*gormmodels.Group, error) {
	var group gormmodels.Group
	err := common.GetDB().Table("t_group").Where("id = ? and enable_flag = ?", id, 1).
		First(&group).
		Error
	return &group, err
}
