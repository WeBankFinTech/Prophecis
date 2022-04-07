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

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/proxy_user"
	"mlss-controlcenter-go/pkg/service"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func AddProxyUser(params proxy_user.AddProxyUserParams) middleware.Responder {
	path := params.KeyPair.Path
	if len(path) <= 0 {
		path = fmt.Sprintf("/data/bdap-ss/mlss-data/%v", params.KeyPair.Name)
	}
	count, err := service.GetProxyUserByUserIdAndNameCount(params.KeyPair.UserID, params.KeyPair.Name)
	addProxyUser := &models.ProxyUser{
		Name:        params.KeyPair.Name,
		Path:        path,
		Remarks:     params.KeyPair.Remarks,
		Gid:         params.KeyPair.Gid,
		UserID:      params.KeyPair.UserID,
		UID:         params.KeyPair.UID,
		IsActivated: params.KeyPair.IsActivated,
		EnableFlag:  1,
	}
	if count <= 0 || err != nil {
		err = service.AddProxyUser(addProxyUser)
		if err != nil {
			logger.Logger().Error("Add proxy user err, ", err)
			return ResponderFunc(http.StatusInternalServerError, "failed to add proxy user", err.Error())
		}
	} else {
		return ResponderFunc(http.StatusBadRequest, "failed to add proxy user", errors.New("name is unqiue").Error())
	}
	return ResponderFunc(http.StatusOK, "", "ProxyUser Add success")
}

func DeleteProxyUser(params proxy_user.DeleteProxyUserParams) middleware.Responder {
	err := service.DeleteProxyUser(params.ID)
	if err != nil {
		logger.Logger().Error("Delete proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to delete proxy user", err.Error())
	}
	return ResponderFunc(http.StatusOK, "", "ProxyUser Delete success")
}

func UpdateProxyUser(params proxy_user.UpdateProxyUserParams) middleware.Responder {
	userId, err := service.GetProxyUserIdByID(params.ID)
	if err != nil {
		logger.Logger().Error("Update proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to update proxy user", err.Error())
	}
	updateProxyUser := &models.ProxyUser{
		ID:          params.ID,
		Path:        params.ProxyUser.Path,
		Remarks:     params.ProxyUser.Remarks,
		Gid:         params.ProxyUser.Gid,
		UID:         params.ProxyUser.UID,
		UserID:      userId,
		IsActivated: params.ProxyUser.IsActivated,
		EnableFlag:  1,
	}
	err = service.UpdateProxyUser(updateProxyUser)
	if err != nil {
		logger.Logger().Error("Update proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to update proxy user", err.Error())
	}
	return ResponderFunc(http.StatusOK, "", "ProxyUser Update success")
}

func GetProxyUser(params proxy_user.GetProxyUserParams) middleware.Responder {
	proxyUser, err := service.GetProxyUser(params.ID)
	if err != nil {
		logger.Logger().Error("Get proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to get proxy user", err.Error())
	}
	marshal, marshalErr := json.Marshal(&proxyUser)
	return GetResult(marshal, marshalErr)
}

func GetProxyUserByUser(params proxy_user.GetProxyUserByUserIDParams) middleware.Responder {
	proxyUser, err := service.GetProxyUserByUser(params.Name, params.UserName)
	if err != nil {
		logger.Logger().Error("Get proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to get proxy user", err.Error())
	}
	marshal, marshalErr := json.Marshal(&proxyUser)
	return GetResult(marshal, marshalErr)
}

func ProxyUserList(params proxy_user.ListProxyUserParams) middleware.Responder {
	isSA := isSuperAdmin(params.HTTPRequest)
	currentUserId := getUserID(params.HTTPRequest)
	user, err := service.GetUserByUserName(params.User)
	if err != nil {
		logger.Logger().Error("Get user err, ", err)
		return ResponderFunc(http.StatusUnauthorized, err.Error(), "permission denied")
		//return ResponderFunc(http.StatusInternalServerError, "failed to user,", err.Error())
	}
	proxyUserList, err := service.ProxyUserList(user.ID, currentUserId, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return ResponderFunc(http.StatusUnauthorized, err.Error(), "permission denied")
		}
		logger.Logger().Error("Get proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to proxy user list,", err.Error())
	}
	marshal, marshalErr := json.Marshal(&proxyUserList)
	return GetResult(marshal, marshalErr)
}

func ProxyUserCheck(params proxy_user.ProxyUserCheckParams) middleware.Responder {
	user, err := service.GetUserByUserName(params.UserName)
	if err != nil {
		logger.Logger().Error("Get user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to user,", err.Error())
	}
	check, err := service.ProxyUserCheck(user.ID, params.ProxyUserName)
	if err != nil {
		logger.Logger().Error("Get proxy user err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed to proxy user list,", err.Error())
	}
	marshal, marshalErr := json.Marshal(&check)
	return GetResult(marshal, marshalErr)
}

func isSuperAdmin(r *http.Request) bool {
	superAdminStr := r.Header.Get("MLSS-Superadmin")
	if "true" == superAdminStr {
		return true
	}
	return false
}

func getUserID(r *http.Request) string {
	return r.Header.Get("MLSS-UserID")
}
