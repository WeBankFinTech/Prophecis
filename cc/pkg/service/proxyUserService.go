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
	"errors"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
)

var (
	PermissionDeniedError = errors.New("permission denied")
)

func GetProxyUserByUser(userName string, name string) (*models.ProxyUser, error) {
	user, err := repo.GetUserByName(userName, datasource.GetDB())
	if err != nil {
		logger.Logger().Error("Get user err, ", err)
		return nil, err
	}
	return repo.GetProxyUserByUserIdAndName(user.ID, name)
}

func AddProxyUser(proxyUser *models.ProxyUser) error {
	err := repo.AddProxyUser(proxyUser)
	if err != nil {
		logger.Logger().Error("Add proxy user err, ", err)
		return err
	}
	return nil
}

func DeleteProxyUser(proxyUserId int64) error {
	err := repo.DeleteProxyUser(proxyUserId)
	if err != nil {
		logger.Logger().Error("Delete proxy user err, ", err)
		return err
	}
	return nil
}

func UpdateProxyUser(proxyUser *models.ProxyUser) error {
	err := repo.UpdateProxyUser(proxyUser)
	if err != nil {
		logger.Logger().Error("Update proxy user err, ", err)
		return err
	}
	return nil
}

func GetProxyUser(proxyUserId int64) (*models.ProxyUser, error) {
	proxyUser, err := repo.GetProxyUser(proxyUserId)
	if err != nil {
		logger.Logger().Error("Get proxy user err, ", err)
		return nil, err
	}
	return proxyUser, err
}

func GetProxyUserIdByID(proxyUserId int64) (int64, error) {
	userId, err := repo.GetProxyUserIdByID(proxyUserId)
	if err != nil {
		logger.Logger().Error("Update proxy user err, ", err)
		return 0, err
	}
	return userId, nil
}

func ProxyUserList(userId int64, username string, isSA bool) ([]*models.ProxyUser, error) {
	proxyUserList, err := repo.ProxyUserList(userId)
	if err != nil {
		logger.Logger().Error("Get proxy user list err, ", err)
		return nil, err
	}
	err = PermissionCheck(username, userId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	return proxyUserList, nil
}

func GetProxyUserByUserIdAndNameCount(userId int64, name string) (int64, error) {
	count, err := repo.GetProxyUserByUserIdAndNameCount(userId, name)
	if err != nil {
		logger.Logger().Error("Get proxy user count err, ", err)
		return 0, err
	}
	return count, nil
}

func ProxyUserCheck(userId int64, name string) (bool, error) {
	count, err := repo.CountProxyUserByUserId(userId, name)
	if err != nil {
		logger.Logger().Error("Get proxy user count err, ", err)
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	return true, nil
}

func PermissionCheck(execUserName string, entityUserId int64, groupId *string, isSuperAdmin bool) error {
	if isSuperAdmin {
		return nil
	}
	db := datasource.GetDB()
	user, err := repo.GetUserByName(execUserName, db)
	if err != nil {
		return errors.New("Permission Check Error, Get user err: " + err.Error())
	}
	if user.ID != entityUserId {
		return PermissionDeniedError
	}
	return nil
}
