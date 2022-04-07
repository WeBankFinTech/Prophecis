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
package umclient

import (
	bytesUtil "bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UmClient struct {
	UmAuth  string `json:"UmAuth,omitempty" yaml:"UmAuth" description:"UmAuth"`
	UmToken string `json:"UmToken,omitempty" yaml:"UmToken" description:"UmToken"`
}

func MD5(message string) string {
	h := md5.New()
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func GETSHA256hashCode(password []byte, username []byte) string {
	hash := sha256.New()
	hash.Write(password)
	hash.Write(username)
	bytes := hash.Sum(nil)
	var hashCode bytesUtil.Buffer
	for _, b := range bytes {
		hashCode.WriteString(fmt.Sprintf("%02x", string([]byte{b})))
	}

	return hashCode.String()
}

func CheckUserAndPassword(username string, password string) (*models.LoginResult, error) {
	appConfig := common.GetAppConfig()
	umHost := appConfig.Core.Um.Host
	umPort := appConfig.Core.Um.Port
	umAppId := appConfig.Core.Um.AppId
	umAppToken := appConfig.Core.Um.AppToken

	umURL := fmt.Sprintf("http://%s:%s/um_service", umHost, umPort)
	authResult, authErr := UMAuth(umURL, umAppId, umAppToken)
	if nil != authErr {
		return nil, authErr
	}

	logger.Logger().Debugf("login -> CheckUserAndPassword -> UMAuth result: %v", authResult)
	if 0 != authResult.RetCode || "" == authResult.Tok {
		return nil, errors.New("auth failed from um system")
	}

	if len(password) > 200 || strings.HasPrefix(password, "{RSA}") {
		if strings.HasPrefix(password, "{RSA}") {
			password = strings.TrimSpace(password[5:])
		}
	}

	pwBytes, deErr := base64.StdEncoding.DecodeString(password)

	if nil != deErr {
		return nil, deErr
	}
	password = strings.TrimSpace(string(pwBytes))

	result, authErr := Login(umURL, authResult, username, password)

	if nil != authErr {
		return nil, authErr
	}

	return result, nil
}

func UMAuth(umStyle2URL string, appId string, appToken string) (*models.AuthResult, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	currentTime := time.Now().UnixNano()
	nonce := strconv.FormatInt(currentTime%9000+1000, 10)

	tmp := MD5(appId + nonce + timeStamp)
	sign := MD5(tmp + appToken)

	data := fmt.Sprintf("style=2&appid=%s&timeStamp=%s&nonce=%s&sign=%s", appId, timeStamp, nonce, sign)
	fmt.Println("umStyle2URL", umStyle2URL)
	fmt.Println("data", data)

	resp, postErr := http.Post(umStyle2URL, "application/x-www-form-urlencoded", strings.NewReader(data))

	if postErr != nil {
		return nil, postErr
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	str := string(body)
	//
	//if readErr != readErr {
	//	return nil, readErr
	//}

	authRequest := &models.AuthResult{}

	unmarshalErr := json.Unmarshal([]byte(str), &authRequest)
	if nil != unmarshalErr {
		return nil, unmarshalErr
	}

	return authRequest, nil
}

func Login(umURL string, result *models.AuthResult, username string, password string) (*models.LoginResult, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	tmp := GETSHA256hashCode([]byte(password), []byte(fmt.Sprintf("{%s}", username)))
	sign := MD5(username + tmp + timeStamp)

	url := fmt.Sprintf("%s?id=%s&timeStamp=%s&sign=%s&style=6&appid=%s&token=%s&auth=%s", umURL, username, timeStamp, sign, result.ID, result.Tok, result.Auth)

	resp, postErr := http.Get(url)
	if nil != postErr {
		return nil, postErr
	}

	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if nil != readErr {
		return nil, readErr
	}

	str := string(body)

	loginResult := &models.LoginResult{}

	unmarshalErr := json.Unmarshal([]byte(str), &loginResult)

	fmt.Printf("str: %v", str)

	if nil != unmarshalErr {
		return nil, unmarshalErr
	}

	return loginResult, nil

}
