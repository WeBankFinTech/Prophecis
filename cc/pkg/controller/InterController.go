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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gopkg.in/cas.v2"
	"log"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/inters"
	"mlss-controlcenter-go/pkg/service"
	"mlss-controlcenter-go/pkg/umclient"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func UserInterceptor(params inters.UserInterceptorParams) middleware.Responder {
	headers := params.HTTPRequest.Header

	token := headers.Get("MLSS-Token")
	hdserUser := headers.Get("Mlss-Userid")
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
				if hdserUser != sessionUser.UserName{
					return ResponderFunc(http.StatusForbidden, "failed to auth by UserInterceptor", "userid does match cookies")
				}

				var result = models.Result{
					Code:    "200",
					Message: "success",
					Result:  json.RawMessage([]byte{}),
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
			return ResponderFunc(http.StatusForbidden, "failed to auth by UserInterceptor", "token was not found")
		}
	}

	authType := headers.Get("MLSS-Auth-Type")

	var user *models.SessionUser
	var cheToken string
	var err error

	if "SSO" == authType {
		user, cheToken, err = checkSSOAuth(params.HTTPRequest)
	}

	if "UM" == authType {
		user, cheToken, err = checkUMAuth(params.HTTPRequest)
	}

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
	permissions, err := service.GetPermissions()
	if err != nil {
		return false, err
	}
	likePermissions, err := service.GetPermissionsByLike(path)
	if err != nil {
		return false, err
	}
	errMsg := fmt.Sprintf("user can't access the api: %s", path)
	isSA := service.HasSA(username)
	pathSplit := strings.Split(path,"/")
	//Get url's permission
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
				if len(likePermissions) >= 2 && pathSplit[len(pathSplit) - 1] != "*" && !isSA{
					for _, likePm := range likePermissions {
						urlSplit := strings.Split(likePm.URL, "/")
						if urlSplit[len(urlSplit) - 1] == pathSplit[len(pathSplit) - 1] {
							permission = pm
							break
						}
					}
				}else{
					permission = pm
					break
				}
			}
		}
	}
	if permission == nil {
		return false, errors.New(fmt.Sprintf("user can't access the api: %s", path))
	}
	if isSA {
		userRolePermissions, err := service.GetRolePermissionsByRoleId(1)
		if err != nil {
			return false, err
		}
		//check permission
		for _, rolePer := range userRolePermissions {
			if permission.ID == rolePer.PermissionID {
				return true, nil
			}
		}
	}else{
		userRolePermissions, err := service.GetRolePermissionsByRoleId(2)
		if err != nil {
			return false, err
		}
		//check permission
		for _, rolePer := range userRolePermissions {
			if permission.ID == rolePer.PermissionID {
				return true, nil
			}
		}
	}
	logger.Logger().Debugf("auth interceptor with permission: %v", permission)
	return false, errors.New(errMsg)
}

func checkSSOAuth(r *http.Request) (*models.SessionUser, string, error) {
	// 1. 验证头部完整性
	// 1. check if the header hss anything the verification need
	ticket := r.Header.Get(constants.AUTH_HEADER_TICKET)
	uiUrl := r.Header.Get(constants.AUTH_HEADER_UIURL)

	if "" == ticket || "" == uiUrl {
		return nil, "", errors.New("ticket and url can't be empty")
	}

	logger.Logger().Debugf("checkSSOAuth sso auth for param url: %v ticket: %v", uiUrl, ticket)

	var scheme, host, path string

	if strings.HasPrefix(uiUrl, "https") {
		scheme = "https"
		indexH := strings.Index(uiUrl[8:], "/")
		indexP := strings.Index(uiUrl[8:], "?")
		host = uiUrl[8:][0:indexH]
		path = uiUrl[8:][indexH:indexP]
	} else {
		scheme = "http"
		indexH := strings.Index(uiUrl[7:], "/")
		indexP := strings.Index(uiUrl[7:], "?")
		host = uiUrl[7:][0:indexH]
		path = uiUrl[7:][indexH:indexP]
	}

	// 2. 验证ticket，并获取用户信息
	var url = &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: ticket,
		Fragment: "/",
	}

	sessionUser, authErr := auth(url, ticket)

	if nil != authErr {
		return nil, "", authErr
	}

	logger.Logger().Debugf("get sessionUser by ticket: %v", sessionUser)

	// 3. check if user exists in db
	username := sessionUser.UserName
	if "" == username {
		return nil, "", errors.New("username is empty")
	}

	saByName := service.GetSAByName(username)
	userByName, err := service.GetUserByUserName(username)
	if err != nil {
		return nil, "", err
	}

	logger.Logger().Debugf("get sessionUser to get userByName: %v", userByName)

	if username != saByName.Name && username != userByName.Name {
		return nil, "", errors.New("user not exists in system")
	}

	// 3. 放置响应头部及缓存token
	// 3. put the auth token in the response header
	token := uuid.New().String()
	authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)

	return sessionUser, token, nil
}

func checkUMAuth(r *http.Request) (*models.SessionUser, string, error) {
	username := r.Header.Get(constants.AUTH_HEADER_USERID)
	password := r.Header.Get(constants.AUTH_HEADER_PWD)
	if "" == username || "" == password {
		return nil, "", errors.New("username and password can't be empty")
	}

	loginResult, checkErr := umclient.CheckUserAndPassword(username, password)
	if nil != checkErr {
		return nil, "", checkErr
	}
	if nil != loginResult && loginResult.Code == 0 {
		var username string
		if "" == loginResult.UserName && "" != loginResult.ID {
			username = loginResult.ID
		} else if "" != loginResult.UserName && "" == loginResult.ID {
			username = loginResult.UserName
		} else if "" != loginResult.UserName && "" != loginResult.ID {
			username = loginResult.UserName
		}
		var sessionUser = &models.SessionUser{
			OrgCode:      loginResult.Org,
			DeptCode:     loginResult.Dept,
			UserName:     username,
			UserID:       username,
			IsSuperadmin: service.GetSAByName(username).Name == username,
		}
		token := uuid.New().String()
		authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)
		return sessionUser, token, nil
	}

	return nil, "", errors.New("failed to check user and password")
}

func checkSystemAuth(r *http.Request) (*models.SessionUser, string, error) {
	// 1. 检查头部完整
	// 1. check if the header hss anything the verification need
	reqApiKey := r.Header.Get(constants.AUTH_HEADER_APIKEY)
	reqTimestamp := r.Header.Get(constants.AUTH_HEADER_TIMESTAMP)
	reqSign := r.Header.Get(constants.AUTH_HEADER_SIGNATURE)
	userId := r.Header.Get(constants.AUTH_HEADER_USERID)
	if "" == reqApiKey || "" == reqTimestamp || "" == reqSign || "" == userId {
		return nil, "", errors.New("missing apiKey or timestamp or signature or userId header")
	}

	timestamp, parErr := strconv.ParseInt(reqTimestamp, 10, 64)
	if nil != parErr {
		return nil, "", parErr
	}

	// 2. 检查timestamp超时
	// 2. whether the timestamp is expired
	//TODO
	current := time.Now().UnixNano() / 1e6
	//defaultTimestampTimeout := common.GetAppConfig().Core.Cache.DefaultTimestampTimeout
	logger.Logger().Debugf("current: %v", current)
	logger.Logger().Debugf("timestamp: %v", timestamp)
	logger.Logger().Debugf("current-timestamp: %v", current-timestamp)
	//if current-timestamp > defaultTimestampTimeout {
	//	return nil, "", errors.New("timestamp in request timed out")
	//}

	// 3. 根据apikey找到对应的secretkey
	// 3. find the secret key according the api key
	keyPair,err := service.GetKeyByApiKey(reqApiKey)
	if err != nil {
		logger.Logger().Error("Get key by api key err, ", err)
		return nil, "", err
	}
	if reqApiKey != keyPair.APIKey {
		return nil, "", errors.New("cannot find secretKey from apiKey")
	}
	secretKey := keyPair.SecretKey

	// 4. 将apikey及timestamp头部生成有序map
	// 4. generate the SortedMap with apikey and timestamp
	reqForm := map[string]string{}
	reqForm[constants.AUTH_HEADER_APIKEY] = reqApiKey
	reqForm[constants.AUTH_HEADER_TIMESTAMP] = reqTimestamp
	reqForm[constants.AUTH_HEADER_AUTH_TYPE] = constants.TypeSYSTEM
	reqForm[constants.AUTH_HEADER_USERID] = userId

	// 5. 计算出签名并与传来的签名比对
	// 5. calculate Sign and check whether it is equal with whe sign input

	logger.Logger().Debugf("secretKey: %v", secretKey)

	//验证用户
	userByUserName, err := service.GetUserByUserName(userId)
	if err != nil {
		return nil, "", err
	}
	token := uuid.New().String()

	var sessionUser = &models.SessionUser{
		UserName: userByUserName.Name,
		UserID:   userByUserName.Name,
	}
	authcache.TokenCache.Set(token, sessionUser, cache.DefaultExpiration)

	return sessionUser, token, nil
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

	loginResult, checkErr := LDAPAuth(username, password)
	if nil != checkErr {
		return nil, "", checkErr
	}
	if loginResult {
		var username string
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
