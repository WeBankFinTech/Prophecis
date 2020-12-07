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
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/logouts"
	"net/http"
)

func Logout(params logouts.LogoutParams) middleware.Responder {
	logger.Logger().Debugf("logout params: %v", params)
	token := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)
	logger.Logger().Debugf("LogoutController token: %v", token)
	tokenErr := removeTokenFromCache(token)
	if nil != tokenErr {
		return ResponderFunc(http.StatusBadRequest, "failed to logout", tokenErr.Error())
	}

	get, b := authcache.TokenCache.Get(token)
	if b {
		sessionUser, e := common.FromInterfaceToSessionUser(get)
		if e != nil {
			logger.Logger().Errorf("LogoutController err :%v", e.Error())
		}
		logger.Logger().Debugf("LogoutController tokenCache sessionUser: %v", sessionUser)
	}

	//logger.Logger().Debugf("User IP: %v logged out", params.HTTPRequest.Header.Get(constants.AUTH_HEADER_REAL_IP))
	//ip := params.HTTPRequest.Header.Get(constants.X_REAL_IP)
	//logger.Logger().Debugf("LogoutController xRealIP: %v", ip)
	//ipErr := removeFromCache(ip)
	//if nil != ipErr {
	//	return ResponderFunc(http.StatusBadRequest, "failed to logout", ipErr.Error())
	//}

	//ipGet, ipB := authcache.IPCache.Get(ip)
	//if ipB {
	//	sessionUser, e := common.FromInterfaceToSessionUser(ipGet)
	//	if e != nil {
	//		logger.Logger().Errorf("LogoutController err :%v", e.Error())
	//	}
	//	logger.Logger().Debugf("LogoutController ipCache sessionUser: %v", sessionUser)
	//}

	//logout sso
	//ssoLogout()

	return GetResult([]byte{}, nil)
}

func ssoLogout() {
	logger.Logger().Debugf("ssoLogout() called.")
	resp, err := http.Get(common.GetAppConfig().Core.Sso.CasLogout)
	if err != nil {
		// handle error
		logger.Logger().Errorf("ssoLogout GET failed, %v", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		logger.Logger().Errorf("ssoLogout ReadAll body failed, %v", err.Error())
	}

	logger.Logger().Debugf(string(body))
}

func removeTokenFromCache(token string) error {
	if "" != token {
		get, b := authcache.TokenCache.Get(token)
		if b {
			sessionUser, e := common.FromInterfaceToSessionUser(get)
			if e != nil {
				return e
			}
			logger.Logger().Debugf("removeTokenFromCache sessionUser: %v", sessionUser)
			authcache.TokenCache.Delete(token)
		}
	}
	return nil
}

//func removeFromCache(ip string) error {
//	if "" != ip {
//		get, b := authcache.IPCache.Get(ip)
//		if b {
//			sessionUser, e := common.FromInterfaceToSessionUser(get)
//			if e != nil {
//				return e
//			}
//			logger.Logger().Debugf("removeTokenFromCache sessionUser: %v", sessionUser)
//			authcache.IPCache.Delete(ip)
//		}
//	}
//	return nil
//}
