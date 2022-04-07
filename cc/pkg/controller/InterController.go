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
	"log"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/inters"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gopkg.in/cas.v2"
)

func UserInterceptor(params inters.UserInterceptorParams) middleware.Responder {
	headers := params.HTTPRequest.Header
	token := headers.Get("MLSS-Token")
	userID := headers.Get(constants.AUTH_HEADER_USERID)
	if len(token) <= 0 {
		token = headers.Get("Mlss-Token")
	}
	logger.Logger().Debugf("UserInterceptor MLSS-Token header: %v", token)
	var sessionUser models.SessionUser
	if "" != token {
		get, b := authcache.TokenCache.Get(token)
		if b {
			//logger.Logger().Debugf("get sessionUser success: %v", b)
			log.Printf("get sessionUser success: %v ", b)
			if nil != get {
				//reflash token expire time
				authcache.TokenCache.Set(token, get, cache.DefaultExpiration)

				session, e := common.FromInterfaceToSessionUser(get)
				log.Printf("get sessionUser by token: %v", session)
				if nil != e {
					return ResponderFunc(http.StatusInternalServerError, "failed to auth by UserInterceptor", e.Error())
				}
				sessionUser = *session
				var result = models.Result{
					Code:    "200",
					Message: "success",
					Result:  json.RawMessage([]byte{}),
				}

				if sessionUser.UserName != userID {
					logger.Logger().Debugf("UserInterceptor SessionUser : %v, Header User : %v", sessionUser.UserID, userID)
					return ResponderFunc(http.StatusForbidden, "failed to auth by UserInterceptor", "token presented but invalid")
				}

				return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
					service.SetResponseHeader(w, sessionUser)
					w.Header().Set("MLSS-Token", token)
					payload, _ := json.Marshal(result)
					w.Write(payload)
				})
			} else {
				return ResponderFunc(http.StatusForbidden, "failed to auth by UserInterceptor", "token presented but invalid")
			}

		} else {
			return ResponderFunc(http.StatusUnauthorized, "failed to auth by UserInterceptor", "token was not found")
		}
	}

	authType := headers.Get("MLSS-Auth-Type")

	var user *models.SessionUser
	var cheToken string
	var err error

	if "SYSTEM" == authType {
		user, cheToken, err = checkSystemAuth(params.HTTPRequest)
	}

	if "LDAP" == authType {
		user, cheToken, err = checkLDAPAuth(params.HTTPRequest)
	}

	if nil != err {
		logger.Logger().Errorf("user: %v, token: %v, err: %v", user, cheToken, err)
		log.Printf("user: %v, token: %v, err: %v", user, cheToken, err)
		errMsg := fmt.Sprintf("failed to auth by: %v", authType)
		return ResponderFunc(http.StatusForbidden, errMsg, err.Error())
	}

	logger.Logger().Debugf("auth success for user: %v", user)

	return getSuccessRes(cheToken, user)
}

func AuthInterceptor(params inters.AuthInterceptorParams) middleware.Responder {
	flag, err := checkPermissionAuth(params.HTTPRequest)
	if nil != err || !flag {
		return ResponderFunc(http.StatusForbidden, "failed to auth by AuthInterceptor", err.Error())
	}

	return getSuccessRes(params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN), nil)
}

func IpInterceptor(params inters.IPInterceptorParams) middleware.Responder {
	flag, err := checkPermissionIP(params.HTTPRequest)
	if nil != err || !flag {
		return ResponderFunc(http.StatusForbidden, "failed to auth by IpInterceptor", err.Error())
	}

	return getSuccessRes(params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN), nil)
}

func getSuccessRes(token string, sessionUser *models.SessionUser) middleware.Responder {
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage([]byte{}),
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		if sessionUser != nil && sessionUser.UserName != "" {
			sessionUser.IsSuperadmin = service.GetSAByName(sessionUser.UserName).Name == sessionUser.UserName
			service.SetResponseHeader(w, *sessionUser)
		}
		w.Header().Set("MLSS-Token", token)
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func checkPermissionIP(r *http.Request) (bool, error) {
	headers := r.Header
	for k, v := range headers {
		logger.Logger().Debugf("%v = %v", k, v)
	}
	logger.Logger().Debugf("checkPermissionIP receiving from gateway, REAL_IP: %v", headers.Get(constants.AUTH_HEADER_REAL_IP))

	ip := headers.Get(constants.X_REAL_IP)
	logger.Logger().Debugf("checkPermissionIP xRealIP: %v", ip)
	if "" == ip {
		ip = r.RemoteAddr
		logger.Logger().Debugf("checkPermissionIP RemoteAddr: %v", ip)
	}

	token := headers.Get(constants.AUTH_HEADER_TOKEN)
	if "" != token {
		get, b := authcache.TokenCache.Get(token)
		if b {
			sessionUser, e := common.FromInterfaceToSessionUser(get)
			if nil != e {
				return false, e
			}
			logger.Logger().Debugf("checkPermissionIP set user: %v to ip: %v cache", sessionUser, ip)
			//authcache.IPCache.Set(ip, sessionUser, cache.DefaultExpiration)
			return true, nil
		} else {
			return false, errors.New("invalid token")
		}
	}

	return false, errors.New("unauthorized IP")
}

func checkPermissionAuth(r *http.Request) (bool, error) {
	headers := r.Header
	if "" == headers.Get(constants.AUTH_HEADER_TOKEN) && "" == headers.Get(constants.AUTH_HEADER_USERID) {
		return false, errors.New("token or userId is not exist in header")
	}

	var username string
	if "" != headers.Get(constants.AUTH_HEADER_TOKEN) {
		get, b := authcache.TokenCache.Get(headers.Get(constants.AUTH_HEADER_TOKEN))
		if b {
			sessionUser, e := common.FromInterfaceToSessionUser(get)
			if nil != e {
				return false, e
			}
			username = sessionUser.UserName
		}
	}

	if "" != headers.Get(constants.AUTH_HEADER_USERID) {
		username = headers.Get(constants.AUTH_HEADER_USERID)
	}

	if username == service.GetSAByName(username).Name {
		return true, nil
	}

	path := headers.Get(constants.AUTH_REAL_PATH)
	method := headers.Get(constants.AUTH_REAL_METHOD)
	logger.Logger().Debugf("request for url: %v and method: %v", path, method)
	log.Printf("request for url: %v and method: %v", path, method)

	var permission *models.Permission

	permissions, _ := service.GetPermissions()
	if nil != permissions && len(permissions) > 0 {
		for _, pm := range permissions {
			pmMethod := strings.TrimSpace(pm.Method)
			pmUrlRegx := strings.ReplaceAll(strings.TrimSpace(pm.URL), "*", "[^/]+")
			matched, err := regexp.MatchString(pmUrlRegx, path)
			logger.Logger().Debugf("auth interceptor with pmUrlRegx: %v, pmMethod: %v, matched: %v", pmUrlRegx, pmMethod, matched)
			if nil != err {
				return false, err
			}
			if method == pmMethod && matched {
				permission = pm
				break
			}
		}
	}
	logger.Logger().Debugf("auth interceptor with permission: %v", permission)
	if nil == permission || "" == permission.URL {
		return false, errors.New("the url is not exist in db of you access")
	}

	if permission.EnableFlag == 0 {
		return false, errors.New("the url is not access of GA ")
	}

	roleId := getRoleIdByRequest(r)
	perSet, _ := service.GetPermissionsByRoleId(roleId)
	if len(perSet.List()) > 0 && perSet.Has(int(permission.ID)) {
		return true, nil
	}

	errMsg := fmt.Sprintf("user can't access the api: %s", path)

	return false, errors.New(errMsg)
}

func checkMLFlowPermissionAuth(r *http.Request) (bool, error) {
	headers := r.Header
	username := headers.Get(constants.AUTH_HEADER_USERID)

	if "" != headers.Get(constants.AUTH_HEADER_USERID) {
		username = headers.Get(constants.AUTH_HEADER_USERID)
	}

	if username == service.GetSAByName(username).Name {
		return true, nil
	}

	path := headers.Get(constants.AUTH_REAL_PATH)
	method := headers.Get(constants.AUTH_REAL_METHOD)
	logger.Logger().Debugf("request for url: %v and method: %v", path, method)
	log.Printf("request for url: %v and method: %v", path, method)

	var permission *models.Permission

	permissions, _ := service.GetPermissions()
	if nil != permissions && len(permissions) > 0 {
		for _, pm := range permissions {
			pmMethod := strings.TrimSpace(pm.Method)
			pmUrlRegx := strings.ReplaceAll(strings.TrimSpace(pm.URL), "*", "[^/]+")
			matched, err := regexp.MatchString(pmUrlRegx, path)
			logger.Logger().Debugf("auth interceptor with pmUrlRegx: %v, pmMethod: %v, matched: %v", pmUrlRegx, pmMethod, matched)
			if nil != err {
				return false, err
			}
			if method == pmMethod && matched {
				permission = pm
				break
			}
		}
	}
	logger.Logger().Debugf("auth interceptor with permission: %v", permission)
	if nil == permission || "" == permission.URL {
		return false, errors.New("the url is not exist in db of you access")
	}

	if permission.EnableFlag == 0 {
		return false, errors.New("the url is not access of GA ")
	}

	roleId := getRoleIdByRequest(r)
	perSet, _ := service.GetPermissionsByRoleId(roleId)
	if len(perSet.List()) > 0 && perSet.Has(int(permission.ID)) {
		return true, nil
	}

	errMsg := fmt.Sprintf("user can't access the api: %s", path)

	return false, errors.New(errMsg)
}

func getRoleIdByRequest(r *http.Request) int {
	roleId := 2
	if isSuperAdmin(r) {
		roleId = 1
	}
	return roleId
}

func checkSystemAuth(r *http.Request) (*models.SessionUser, string, error) {
	// 1. 检查头部完整
	// 1. check if the header hss anything the verification need
	reqApiKey := r.Header.Get(constants.AUTH_HEADER_APIKEY)
	// reqTimestamp := r.Header.Get(constants.AUTH_HEADER_TIMESTAMP)
	reqSign := r.Header.Get(constants.AUTH_HEADER_SIGNATURE)
	userId := r.Header.Get(constants.AUTH_HEADER_USERID)
	if "" == reqApiKey || "" == reqSign || "" == userId {
		return nil, "", errors.New("missing apiKey or timestamp or signature or userId header")
	}

	// 2. 根据apikey找到对应的secretkey
	// 2. find the secret key according the api key
	keyPair, err := service.GetKeyByApiKey(reqApiKey)
	if err != nil {
		logger.Logger().Error("Get key by api key err, ", err)
		return nil, "", err
	}
	if reqApiKey != keyPair.APIKey {
		return nil, "", errors.New("cannot find secretKey from apiKey")
	}
	secretKey := keyPair.SecretKey
	if reqSign != secretKey && reqApiKey != "MLFLOW" {
		return nil, "", errors.New("App signature is not matched, APP KEY is" + reqApiKey)
	}

	//3.检查用户是否匹配
	var sessionUser models.SessionUser
	user, err := service.GetUserByUserName(userId)
	if err != nil {
		return nil, "", err
	}
	if keyPair.SuperAdmin == 1 {
		sessionUser.UserName = userId
		sessionUser.UserID = strconv.Itoa(int(user.ID))
	} else {
		if keyPair.Name != userId {
			return nil, "", errors.New("User Name is not matched the app key")
		}
	}
	token := uuid.New().String()
	authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)

	return &sessionUser, token, nil

}

func auth(su *url.URL, ticket string) (*models.SessionUser, error) {
	client := &http.Client{}

	config := common.GetAppConfig()
	casUrl := config.Core.Sso.CasServiceUrlPrefix
	logger.Logger().Debugf("casUrl: %v", casUrl)

	u, err := url.Parse(casUrl)
	if err != nil {
		fmt.Printf("url.Parse failed in cas auth: %v", err.Error())
		return nil, err
	}

	validator := cas.NewServiceTicketValidator(client, u)
	response, errVa := validator.ValidateTicket(su, ticket)
	logger.Logger().Debugf("get sessionUser from response: %v", response)
	if errVa != nil {
		fmt.Printf("%v", errVa.Error())
		return nil, errVa
	}
	return &models.SessionUser{
		OrgCode:  response.Attributes.Get("orgCode"),
		DeptCode: response.Attributes.Get("deptCode"),
		UserName: response.Attributes.Get("loginid"),
	}, nil
}

func checkLDAPAuth(r *http.Request) (*models.SessionUser, string, error) {
	username := r.Header.Get(constants.AUTH_HEADER_USERID)
	password := r.Header.Get(constants.AUTH_HEADER_PWD)
	if "" == username || "" == password {
		return nil, "", errors.New("username and password can't be empty")
	}

	var loginResult bool
	var checkErr error
	if username == common.GetAppConfig().Application.Admin.User {
		if password == common.GetAppConfig().Application.Admin.Password {
			loginResult = true
		}
	} else {
		loginResult, checkErr = LDAPAuth(username, password)
	}
	if nil != checkErr {
		return nil, "", checkErr
	}
	if loginResult == false {
		return nil, "admin login error", errors.New("admin password is wrong")
	}
	if loginResult {
		sa := service.GetSAByName(username)
		var sessionUser = &models.SessionUser{
			UserName:     username,
			UserID:       username,
			IsSuperadmin: sa.Name == username,
		}
		token := uuid.New().String()
		authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)
		return sessionUser, token, nil
	}

	return nil, "", errors.New("failed to check user and password")
}
