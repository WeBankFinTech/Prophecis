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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/users"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func GetAllUsers(params users.GetAllUsersParams) middleware.Responder {
	page := *params.Page
	size := *params.Size
	error := common.CheckPageParams(page, size)
	if nil != error {
		return ResponderFunc(http.StatusBadRequest, "failed to check for page and size", error.Error())
	}

	logger.Logger().Debugf("v1/users GetAllUsers params page: %v, size: %v", page, size)

	pageUsers, error := service.GetAllUsers(page, size)
	if nil != error {
		return ResponderFunc(http.StatusBadRequest, "failed to GetAllUsers", error.Error())
	}
	marshal, marshalErr := json.Marshal(pageUsers)

	return GetResult(marshal, marshalErr)
}

func GetMyUsers(params users.GetMyUsersParams) middleware.Responder {
	//clusterName := *params.ClusterName
	sessionUser := GetSessionByContext(params.HTTPRequest)

	logger.Logger().Debugf("GetMyUsers cookies length: %v", len(params.HTTPRequest.Cookies()))
	logger.Logger().Debugf("GetMyUsers cookies: %v", params.HTTPRequest.Cookies())

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	var myUsers []*models.User

	if IsSuperAdmin(sessionUser) {
		saUserList, err := service.GetAllUser()
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetAllUser", err.Error())
		}
		myUsers = saUserList
	} else {
		username := sessionUser.UserName
		user, err := service.GetUserByUserName(username)
		if err != nil {
			return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
		}
		ids, err := service.GetGroupIdListByUserIdRoleId(user.ID, constants.GA_ID)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserIdRoleId", err.Error())
		}

		if len(ids) > 0 {
			myUsers, err = service.GetAllUserNameByGroupIds(ids)
			if err != nil {
				return ResponderFunc(http.StatusInternalServerError, "failed to GetAllUserNameByGroupIds", err.Error())
			}
		} else {
			myUsers = append(myUsers, user)
		}
	}
	marshal, marshalErr := json.Marshal(myUsers)

	return GetResult(marshal, marshalErr)
}

func GetSAByName(params users.GetSAByNameParams) middleware.Responder {
	adminName := params.Name
	saByName := service.GetSAByName(adminName)
	if adminName != saByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to get SA by adminName", "the SA is not exist in db")
	}

	marshal, marshalErr := json.Marshal(saByName)

	return GetResult(marshal, marshalErr)
}

func AddUser(params users.AddUserParams) middleware.Responder {
	//log := logger.Logger()
	userCheckError := common.CheckForAddUser(*params.User)
	if userCheckError != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to add user", userCheckError.Error())
	}
	userRequest := *params.User

	existed, err := service.IfUserIdExisted(userRequest.UID)
	if err != nil {
		logger.Logger().Error("Check UserID is Existed Error: ", err.Error())
	}
	if existed {
		userByUID, err := service.GetUserByUID(userRequest.UID)
		if err != nil {
			logger.Logger().Error("GetUserByUID Error: ", err.Error())
		}

		if userByUID.EnableFlag != 1 {
			return ResponderFunc(http.StatusBadRequest, "failed to add user", "the uid has been assigned")
		}
	}

	repoUser, addUserErr := service.AddUser(userRequest)
	if nil != addUserErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to add user", addUserErr.Error())
	}

	marshal, marshalErr := json.Marshal(repoUser)
	return GetResult(marshal, marshalErr)
}

func UpdateUser(params users.UpdateUserParams) middleware.Responder {
	userCheckError := common.CheckForAddingUser(*params.User)
	if userCheckError != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update user", userCheckError.Error())
	}

	userRequest := common.FromUpdateUserRequest(params)

	userByUserName, err := service.GetUserByUserId(userRequest.ID)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetUserByUserId:%v", err.Error())
	}

	if userRequest.ID != userByUserName.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update user", "user not exist in db")
	}

	if userRequest.Name != userByUserName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to update user", "The user id and name do not match")
	}

	repoUser, err := service.UpdateUser(userRequest)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update user", err.Error())
	}
	marshal, marshalErr := json.Marshal(repoUser)
	return GetResult(marshal, marshalErr)
}

func GetUserByUserId(params users.GetUserByUserIDParams) middleware.Responder {
	userId := params.UserID

	user, err := service.GetUserByUserId(userId)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to get user by id", "user is not exist in db"+err.Error())
	}
	if userId != user.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get user by id", "user is not exist in db")
	}

	marshal, marshalErr := json.Marshal(user)

	return GetResult(marshal, marshalErr)
}

func DeleteUserById(params users.DeleteUserByIDParams) middleware.Responder {
	userId := params.UserID

	userByUserId, err := service.GetUserByUserId(userId)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to GetUserByUserId:%v", err.Error())
	}

	if userId != userByUserId.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete user by id", "user is not exist in db")
	}

	repoUser, err := service.DeleteUserById(userId)

	if nil != err {
		ResponderFunc(http.StatusBadRequest, "failed to delete user by id", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoUser)

	return GetResult(marshal, marshalErr)
}

func GetUserByUserName(params users.GetUserByUserNameParams) middleware.Responder {
	userName := params.UserName

	if "" == userName {
		return ResponderFunc(http.StatusBadRequest, "failed to format user", "username can't be empty")
	}

	repoUser, err := service.GetUserByUserName(userName)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
	}
	if userName != repoUser.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to get user by name", "user is not exist in db")
	}

	userRes := common.FromUserToUserRes(*repoUser)

	marshal, marshalErr := json.Marshal(userRes)

	return GetResult(marshal, marshalErr)
}

func DeleteUserByName(params users.DeleteUserByNameParams) middleware.Responder {
	userName := params.UserName

	if "" == userName {
		return ResponderFunc(http.StatusBadRequest, "failed to format user", "username can't be empty")
	}

	userByUserName, err := service.GetUserByUserName(userName)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
	}
	if userName != userByUserName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to delete user by name", "user is not exist in db")
	}

	repoUser, err := service.DeleteUserById(userByUserName.ID)

	if nil != err {
		ResponderFunc(http.StatusBadRequest, "failed to delete user by name", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoUser)

	return GetResult(marshal, marshalErr)
}

func GetResult(marshal []byte, marshalErr error) middleware.Responder {
	if nil != marshalErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to format", marshalErr.Error())
	}

	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func IsSuperAdmin(sessionUser *models.SessionUser) bool {
	if sessionUser != nil && sessionUser.UserName != "" {
		return service.GetSAByName(sessionUser.UserName).Name == sessionUser.UserName
	}
	return false
}

func GetSessionByContext(request *http.Request) *models.SessionUser {
	logger.Logger().Debugf("context: %v", request.Context())
	if request.Header.Get(constants.AUTH_HEADER_AUTH_TYPE) == "SYSTEM" {
		logger.Logger().Debugf("SYSTEM AUTH Toekn is Nil")
		userName := request.Header.Get(constants.AUTH_HEADER_USERID)
		user := models.SessionUser{UserName: userName}
		return &user
	}
	logger.Logger().Info("Token：", request.Header.Get(constants.AUTH_HEADER_TOKEN))
	logger.Logger().Debugf("AUTH TYPE: %v", request.Context())
	sessionUser := request.Context().Value(request.Header.Get(constants.AUTH_HEADER_TOKEN))
	logger.Logger().Debugf("get sessionUser from context: %v\n", request.Header.Get(constants.AUTH_HEADER_AUTH_TYPE))
	if nil != sessionUser {
		return sessionUser.(*models.SessionUser)
	}
	logger.Logger().Error("GET SESSION BY CONTEXT ERROR: Session User is nill.")
	return nil
}

func FailedToGetUserByToken() middleware.Responder {
	return ResponderFunc(http.StatusUnauthorized, "failed to access", "token is empty or invalid")
}

func GetUserToken(params users.GetUserTokenParams) middleware.Responder {
	headerToken := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)
	user := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_USERID)
	logger.Logger().Infof("CheckURLAccess token: %v, url:%v ", headerToken, user)
	result := models.TokenMsg{User: user, Token: headerToken}
	marshal, marshalErr := json.Marshal(result)
	return GetResult(marshal, marshalErr)
}
