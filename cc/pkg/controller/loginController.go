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
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-ldap/ldap/v3"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/logins"
	"mlss-controlcenter-go/pkg/service"
	"mlss-controlcenter-go/pkg/umclient"
	"net/http"
	"strings"
)

func UMLogin(params logins.UMLoginParams) middleware.Responder {
	username := params.Username
	password := params.Password
	pwdEncoded := false
	if nil != params.PwdEncoded {
		pwdEncoded = *params.PwdEncoded
	}

	if pwdEncoded {
		pwBytes, deErr := base64.StdEncoding.DecodeString(password)
		if nil != deErr {
			ResponderFunc(http.StatusInternalServerError, "failed to login", deErr.Error())
		}
		password = strings.TrimSpace(string(pwBytes))
	}

	ipAddr := common.RemoteIp(params.HTTPRequest)
	logger.Logger().Debugf("login for username: %v, password: %v and request ip: %v", username, password, ipAddr)

	password = base64.StdEncoding.EncodeToString([]byte(password))
	password = strings.TrimSpace(password)

	loginResult, err := umclient.CheckUserAndPassword(username, password)

	if nil != err {
		return ResponderFunc(http.StatusForbidden, "failed to login", err.Error())
	}

	if nil != loginResult && loginResult.Code == 0 {
		logger.Logger().Debugf("The checking of user succeed. user: %v", username)
		sessionUser := models.SessionUser{
			UserName:     loginResult.UserName,
			UserID:       loginResult.ID,
			Email:        loginResult.Email,
			OrgCode:      loginResult.Org,
			DeptCode:     loginResult.Dept,
			AccountType:  loginResult.Actype,
			IsSuperadmin: service.GetSAByName(loginResult.UserName).Name == loginResult.UserName,
		}

		logger.Logger().Debugf("sessionUser: %v", sessionUser)

		marshal, _ := json.Marshal(sessionUser)
		var result = models.Result{
			Code:    "200",
			Message: "success",
			Result:  json.RawMessage(marshal),
		}

		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			service.UMLogin(w, sessionUser)
			payload, _ := json.Marshal(result)
			w.Write(payload)
		})
	}

	return ResponderFunc(http.StatusForbidden, "failed to login", "failed to login")
}

func LDAPLogin(params logins.LDAPLoginParams) middleware.Responder {
	username := params.LoginRequest.Username
	password := params.LoginRequest.Password

	isAccess, err := LDAPAuth(username, password)
	if err != nil || isAccess == false {
		return ResponderFunc(http.StatusBadRequest, "failed to login", "failed to login")
	}

	p, err := service.CheckUserPermission(username)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to login", err.Error())
	}
	if !p {
		return ResponderFunc(http.StatusUnauthorized, "failed to login, ", "User does not have system permissions")
	}

	sessionUser := models.SessionUser{
		UserName:     username,
		IsSuperadmin: service.GetSAByName(username).Name == username,
	}

	logger.Logger().Debugf("sessionUser: %v", sessionUser)

	marshal, _ := json.Marshal(sessionUser)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}

	//authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		cookie := service.LDAPLogin(w, common.GetAppConfig().Core.Cookie.Path, sessionUser)
		http.SetCookie(w, &cookie)
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func LDAPAuth(username string, password string) (bool, error) {
	l, err := ldap.DialURL(common.GetAppConfig().Application.LDAP)
	if err != nil {
		logger.Logger().Errorf("ldap server dial error" + err.Error())
		return false, err
	}
	if l == nil {
		logger.Logger().Errorf("ldap server dial error, connection is nil")
		return false, errors.New("ldap server dial error,connection is nil")
	}

	passwordDecode, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		logger.Logger().Errorf("Password decode Error" + err.Error())
	}

	_, err = l.SimpleBind(&ldap.SimpleBindRequest{
		Username: username,
		Password: string(passwordDecode),
	})
	if err != nil {
		logger.Logger().Errorf("LDAP Server Auth Error: %s\n", err)
		return false, err
	}
	defer l.Close()
	return true, err
}
