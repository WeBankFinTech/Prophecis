/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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
	"gorm.io/gorm"
	"webank/DI/commons/models"
	datasource "webank/DI/pkg/datasource/mysql"
)

const TableNameUser = "t_user"

var UserRepo UserRepoIF

func init() {
	UserRepo = &UserRepoImpl{
		BaseRepoImpl{TableNameUser},
	}
}

type UserRepoIF interface {
	GetByTransaction(username *string, db *gorm.DB) (*models.User, error)
	Get(username *string) (*models.User, error)
}

type UserRepoImpl struct {
	BaseRepoImpl
}

func (i UserRepoImpl) GetByTransaction(username *string, db *gorm.DB) (*models.User, error) {
	var user models.User
	err := db.Table(i.TableName).Find(&user, "name = ? AND enable_flag = ?", username, 1).Error
	return &user, err
}

func (i UserRepoImpl) Get(username *string) (*models.User, error) {
	var user models.User
	err := datasource.GetDB().Table(i.TableName).Find(&user, "name = ? AND enable_flag = ?", username, 1).Error
	return &user, err
}
