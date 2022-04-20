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
	"errors"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/models"

	"github.com/jinzhu/gorm"
)

func AddProxyUser(proxyUser *models.ProxyUser) error {
	return datasource.GetDB().Create(&proxyUser).Error
}

func DeleteProxyUser(proxyUserId int64) error {
	return datasource.GetDB().Table("t_proxy_user").Where("id = ?", proxyUserId).UpdateColumn("enable_flag", 0).Error
}

func UpdateProxyUser(proxyUser *models.ProxyUser) error {
	return datasource.GetDB().Table("t_proxy_user").Where("id = ?", proxyUser.ID).Updates(&proxyUser).Error
}

func GetProxyUser(proxyUserId int64) (*models.ProxyUser, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("id = ? AND enable_flag = ?", proxyUserId, 1).Find(&getProxyUser).Error
	return &getProxyUser, err
}

func GetProxyUserIdByID(proxyUserId int64) (int64, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("id = ? AND enable_flag = ?", proxyUserId, 1).Find(&getProxyUser).Error
	return getProxyUser.UserID, err
}

func ProxyUserList(userId int64) ([]*models.ProxyUser, error) {
	var proxyUserList []*models.ProxyUser
	err := datasource.GetDB().Where("user_id = ? AND enable_flag = ?", userId, 1).Find(&proxyUserList).Error
	return proxyUserList, err
}

func GetProxyUserByUserIdAndName(userId int64, name string) (*models.ProxyUser, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("user_id = ? AND name = ? AND enable_flag = ?", userId, name, 1).
		Find(&getProxyUser).Error
	return &getProxyUser, err
}

func GetProxyUserByUserName(name string) (*models.ProxyUser, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("name = ? AND enable_flag = ?", name, 1).
		Find(&getProxyUser).Error
	return &getProxyUser, err
}

func GetProxyUserByUserNameCount(name string) (int64, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("name = ? AND enable_flag = ?", name, 1).
		Find(&getProxyUser).Error
	return getProxyUser.ID, err
}

func GetProxyUserByUserIdAndNameCount(userId int64, name string) (int64, error) {
	var getProxyUser models.ProxyUser
	err := datasource.GetDB().Where("user_id = ? AND name = ? AND enable_flag = ?", userId, name, 1).
		Find(&getProxyUser).Error
	return getProxyUser.ID, err
}

func CountProxyUserByUserId(userId int64, name string) (int64, error) {
	// var getProxyUser models.ProxyUser
	var proxyUserList []*models.ProxyUser
	err := datasource.GetDB().Table("t_proxy_user").Where("user_id = ? AND name = ? AND enable_flag = ?", userId, name, 1).Find(&proxyUserList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return int64(len(proxyUserList)), err
}
