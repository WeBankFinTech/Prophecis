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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/auths"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/logins"
	"mlss-controlcenter-go/pkg/service"
	"mlss-controlcenter-go/pkg/umclient"
	"net/http"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
)

const (
	ldapPubKey  = "ldapPubKey"
	ldapPrivKey = "ldapPrivKey"
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
	decodeBytes, err := base64.StdEncoding.DecodeString(params.LoginRequest.Password)
	if err != nil {
		logger.Logger().Error("failed to login, base64 decode failed:%v", err.Error())
		return ResponderFunc(http.StatusInternalServerError, "failed to login, base64 decode failed:", err.Error())
	}
	decryptPassword, err := RsaDecrypt(decodeBytes)
	if err != nil {
		logger.Logger().Error("failed to login, rsa decrypt failed:%v", err.Error())
		return ResponderFunc(http.StatusInternalServerError, "failed to login, rsa decrypt failed:", err.Error())
	}

	isAccess := false
	if username == common.GetAppConfig().Application.Admin.User {
		if string(decryptPassword) == common.GetAppConfig().Application.Admin.Password {
			isAccess = true
		}
	} else {
		isAccess, err = LDAPAuth(username, string(decryptPassword))
		if err != nil {
			logger.Logger().Error("Failed to login, LDAP Auth Error:", err.Error())
			return ResponderFunc(http.StatusBadRequest, "LDAP auth failed.", "用户名或密码错误。")
		}
	}
	if isAccess == false {
		return ResponderFunc(http.StatusBadRequest, "LDAP auth failed.", "用户名或密码错误。")
	}

	// Check system permission
	p, err := service.CheckUserPermission(username)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to login", err.Error())
	}
	if !p {
		return ResponderFunc(http.StatusUnauthorized, "failed to login, ", "User does not have system permissions")
	}

	//Set Session User for Return
	isSA := service.GetSAByName(username).Name == username
	sessionUser := models.SessionUser{
		UserName:     username,
		IsSuperadmin: isSA,
	}
	logger.Logger().Debugf("sessionUser: %v", sessionUser)
	marshal, _ := json.Marshal(sessionUser)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		cookie := service.LDAPLogin(w, common.GetAppConfig().Core.Cookie.Path, sessionUser)
		http.SetCookie(w, &cookie)
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func LDAPAuth(username string, password string) (bool, error) {
	address := common.GetAppConfig().Application.LDAP.Address
	baseDN := common.GetAppConfig().Application.LDAP.BaseDN

	//Dial LDAP Server
	l, err := ldap.DialURL(address)
	defer l.Close()
	if err != nil {
		logger.Logger().Errorf("LDAP Dial Fail:%v", err.Error())
		return false, err
	}
	if l == nil {
		return false, errors.New("ldap server dial error,connection is nil")
	}

	//Search User in LDAP Server
	nsr := ldap.NewSearchRequest(baseDN, ldap.ScopeBaseObject, ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username), []string{"dn"}, nil)
	sr, err := l.Search(nsr)
	if err != nil {
		logger.Logger().Errorf("LDAP Search Fail:%v", err.Error())
		return false, err
	}

	//Auth User Password
	userDN := sr.Entries[0].DN
	err = l.Bind(userDN, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetRsaPubKey(params logins.GetRsaPubKeyParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		var result = models.Result{
			Code:    "200",
			Message: "success",
			Result:  viper.GetString(ldapPubKey),
		}
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func RsaDecrypt(context []byte) ([]byte, error) {
	privateKey := viper.GetString(ldapPrivKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, context)
}

func CheckGroupByUser(params auths.CheckGroupByUserParams) middleware.Responder {
	userId := params.UserID
	group := repo.GetUserGroupByUserId(userId)
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(group)
		w.Write(payload)
	})
}
