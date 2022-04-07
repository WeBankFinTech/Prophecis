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
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"net/http"
	"strconv"
)

func UMLogin(w http.ResponseWriter, user models.SessionUser) {
	token := uuid.New().String()
	logger.Logger().Debugf("login success to set token: %v", token)
	authcache.TokenCache.Set(token, user, cache.DefaultExpiration)
	w.Header().Add(constants.AUTH_HEADER_TOKEN, token)
	SetResponseHeader(w, user)
}

func SetResponseHeader(res http.ResponseWriter, user models.SessionUser) {
	if "" == user.UserID && "" != user.UserName {
		user.UserID = user.UserName
	} else if "" != user.UserID && "" == user.UserName {
		user.UserName = user.UserID
	}
	res.Header().Add(constants.AUTH_RSP_USERID, user.UserID)
	res.Header().Add(constants.AUTH_RSP_ORGCODE, user.OrgCode)
	res.Header().Add(constants.AUTH_RSP_DEPTCODE, user.DeptCode)
	res.Header().Add(constants.AUTH_RSP_LOGINID, user.UserID)
	res.Header().Add(constants.AUTH_RSP_SUPERADMIN, strconv.FormatBool(user.IsSuperadmin))
}

func LDAPLogin(w http.ResponseWriter, cookiesPath string, user models.SessionUser) http.Cookie {
	token := uuid.New().String()
	logger.Logger().Debugf("login success to set token: %v", token)
	logger.Logger().Infof("Login success to set token, user:  %v , token: %v", user, token)
	authcache.TokenCache.Set(token, user, cache.DefaultExpiration)
	w.Header().Add(constants.AUTH_HEADER_TOKEN, token)
	SetResponseHeader(w, user)
	cookie := &http.Cookie{
		Name:     "PROPHECIS",
		Value:    token,
		HttpOnly: true,
		Path:     cookiesPath,
	}
	return *cookie
}