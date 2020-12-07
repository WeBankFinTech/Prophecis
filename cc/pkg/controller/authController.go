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
	"bytes"
	"encoding/json"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/auths"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func UserStorageCheck(params auths.UserStorageCheckParams) middleware.Responder {
	username := params.Username
	path := params.Path

	logger.Logger().Debugf("UserStorageCheck username: %v, path: %v", username, path)

	userStorageCheckErr := service.UserStorageCheck(username, path)

	logger.Logger().Debugf("UserStorageCheck userStorageCheckErr: %v", userStorageCheckErr)

	if nil != userStorageCheckErr {
		return ResponderFunc(http.StatusBadRequest, "failed to check user and path", userStorageCheckErr.Error())
	}

	return GetResult([]byte{}, nil)
}

func UserNamespaceCheck(params auths.UserNamespaceCheckParams) middleware.Responder {
	response, err := service.UserNamespaceCheck(params.Username, params.Namespace)

	if nil != err {
		logger.Logger().Debugf("UserNamespaceCheck with error: %v", err.Error())
		return ResponderFunc(http.StatusBadRequest, "failed to check user and namespace", err.Error())
	}

	marshal, marshalErr := json.Marshal(response)

	return GetResult(marshal, marshalErr)
}

func UserStoragePathCheck(params auths.UserStoragePathCheckParams) middleware.Responder {
	err := service.UserStoragePathCheck(params.Username, params.Namespace, params.Path)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check user and namespace", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func AdminUserCheck(params auths.AdminUserCheckParams) middleware.Responder {
	err := service.AdminUserCheck(params.AdminUsername, params.Username)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check admin and user", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func CheckNamespace(params auths.CheckNamespaceParams) middleware.Responder {
	namespace := params.Namespace

	sessionUser := GetSessionByContext(params.HTTPRequest)

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	role, err := service.CheckNamespace(namespace, sessionUser.UserName)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check admin and user", err.Error())
	}

	marshal, marshalErr := json.Marshal(role)

	return GetResult(marshal, marshalErr)
}

func CheckNamespaceUser(params auths.CheckNamespaceUserParams) middleware.Responder {
	username := params.Username
	namespace := params.Namespace

	sessionUser := GetSessionByContext(params.HTTPRequest)

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	err := service.CheckNamespaceUser(namespace, sessionUser.UserName, username)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check namespace and user", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func CheckUserGetNamespace(params auths.CheckUserGetNamespaceParams) middleware.Responder {
	username := params.Username
	logger.Logger().Debugf("CheckUserGetNamespace username: %v", username)

	sessionUser := GetSessionByContext(params.HTTPRequest)

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	userNoteBook, err := service.CheckUserGetNamespace(sessionUser.UserName, username)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to get userNoteBook by username", err.Error())
	}

	marshal, marshalErr := json.Marshal(userNoteBook)

	return GetResult(marshal, marshalErr)
}

func CheckCurrentUserNamespacedNotebook(params auths.CheckCurrentUserNamespacedNotebookParams) middleware.Responder {
	namespace := params.Namespace
	notebook := params.Notebook
	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook namespace: %v, notebook: %v", namespace , notebook)
	token := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)
	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook, checking token: %v", token)
	authType := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_AUTH_TYPE)


	var sessionUser *models.SessionUser
	if authType == constants.TypeSYSTEM{
		user := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_USERID)
		logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook get user from header: %v", user)
		sessionUser = &models.SessionUser{UserName:user}
	}else{
		//get, b := ipCache.Get(realIp)
		tokenCache := authcache.TokenCache
		get, b := tokenCache.Get(token)

		logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser b: %v", b)
		if b {
			session, err := common.FromInterfaceToSessionUser(get)
			if nil != err {
				return ResponderFunc(http.StatusInternalServerError, "failed to get sessionUser", err.Error())
			}
			sessionUser = session
		}

	}
	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser %+v", sessionUser)

	if nil == sessionUser {
		return ResponderFunc(http.StatusForbidden, "failed to CheckCurrentUserNamespacedNotebook", "token is empty or invalid")
	}

	username := sessionUser.UserName
	user, err := service.GetUserByUserName(username)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
	}

	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser.IsSuperAdmin: %v", sessionUser.IsSuperadmin)

	namespaceByName,err := service.GetNamespaceByName(namespace)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetNamespaceByName", err.Error())
	}
	if namespace != namespaceByName.Namespace {
		return ResponderFunc(http.StatusForbidden, "failed to access notebook", "namespace is not exist in db")
	}

	var checkErr error
	if !sessionUser.IsSuperadmin {
		var flag = false
		groupIds,err := service.GetGroupIdListByUserIdRoleId(user.ID, constants.GA_ID)
		if err != nil{
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserIdRoleId", err.Error())
		}
		for _, groupId := range groupIds {
			gn,err := service.GetGroupNamespaceByGroupIdAndNamespaceId(groupId, namespaceByName.ID)
			if err == nil{
				if gn.GroupID == groupId && gn.NamespaceID == namespaceByName.ID {
					flag = true
					break
				}
			}
		}
		if !flag {
			checkErr = service.CheckUserNotebook(namespace, username, notebook)
		}
	}

	if nil != checkErr {
		return ResponderFunc(http.StatusForbidden, "access notebook failed", checkErr.Error())
	}

	var nbAddress string

	var notebookBuffer bytes.Buffer
	notebookBuffer.WriteString("/notebook/")
	notebookBuffer.WriteString(namespace)
	notebookBuffer.WriteString("/")
	notebookBuffer.WriteString(notebook)

	//if constants.BDP == strings.Split(namespace, "-")[2] {
	//	logger.Logger().Debugf("gatewayProperties.getBdpAddress(): %v", common.GetAppConfig().Core.Gateway.BdpAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdpAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//} else if constants.BDAP == strings.Split(namespace, "-")[2] && !(constants.SAFE == strings.Split(namespace, "-")[6]) {
	//	logger.Logger().Debugf("gatewayProperties.BdapAddress(): %v", common.GetAppConfig().Core.Gateway.BdapAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdapAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//} else if constants.BDAP == strings.Split(namespace, "-")[2] && constants.SAFE == strings.Split(namespace, "-")[6] {
	//	logger.Logger().Debugf("gatewayProperties.BdapsafeAddress(): %v", common.GetAppConfig().Core.Gateway.BdapsafeAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdapsafeAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//}

	nbAddress = notebookBuffer.String()

	var userNotebookAddressResponse = models.UserNotebookAddressResponse{
		NotebookAddress: nbAddress,
	}

	marshal, marshalErr := json.Marshal(userNotebookAddressResponse)

	return GetResult(marshal, marshalErr)
}

func ResponderFunc(code int, error string, msg string) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(code)
		payload, _ := json.Marshal(models.Error{
			Error: error,
			Code:  int32(code),
			Msg:   msg,
		})
		w.Write(payload)
	})
}
