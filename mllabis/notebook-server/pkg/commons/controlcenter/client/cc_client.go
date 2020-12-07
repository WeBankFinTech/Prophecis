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
package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	cfg "webank/AIDE/notebook-server/pkg/commons/config"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"
)

type CcClient struct {
	ccUrl      string
	httpClient http.Client
}

type UserResponse struct {
	Name      string `json:"name"`
	GID       int    `json:"gid"`
	UID       int    `json:"uid"`
	Token     string  `json:"token"`
	Type      string `json:"type"`
	GUIDCheck bool   `json:"guidCheck"`
	Remarks   string `json:"remarks"`
}

const CcAuthToken = "MLSS-Token"

// FIXME MLSS Change: get models filter by username and namespace
const CcSuperadmin = "MLSS-Superadmin"

var once sync.Once
var client *CcClient

func GetCcClient(ccUrl string) *CcClient {
	once.Do(func() {
		client = &CcClient{}
		//client.httpClient = http.Client{}
		client.ccUrl = ccUrl

		//FIXME
		pool := x509.NewCertPool()

		caCertPath := cfg.GetValue(cfg.CaCertPath)
		logger.Logger().Infof("ccClient to get caCertPath: %v", caCertPath)
		//调用ca.crt文件
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(caCrt)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            pool,
			},
		}

		client.httpClient = http.Client{Transport: tr}
	})
	return client
}

func (cc *CcClient) AuthAccessCheck(authOptions *map[string]string) (int, string, string, error) {
	//apiUrl := cc.ccUrl + "/cc/v1/sample"
	apiUrl := cc.ccUrl + "/cc/v1/inter/user"

	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		logger.Logger().Debug("Create AuthAccessCheck Request Failed !")
		logger.Logger().Debug(err.Error())
		return 500, "", "", err
	}
	for key, value := range *authOptions {
		req.Header.Set(key, value)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		logger.Logger().Debug("Do AuthAccessCheck Request Failed !")
		logger.Logger().Debug(err.Error())
		return 500, "", "", err
	}
	defer resp.Body.Close()
	token := resp.Header.Get(CcAuthToken)
	isSA := resp.Header.Get(CcSuperadmin)
	statusCode := resp.StatusCode
	if statusCode == 401 || (statusCode >= 200 && statusCode < 300 && token == "") {
		logger.Logger().Debugf("Auth Access Check Request Failed, StatusCode: %v, Token: %v", statusCode, token)
		logger.Logger().Debugf("Message: %v", resp.Body)
		return 401, "", "", *new(error)
	}
	if statusCode > 300 || token == "" {
		logger.Logger().Debugf("Auth Access Check Request Failed from cc, StatusCode: %v, Token: %v", statusCode, token)
		logger.Logger().Debugf("Message: %v", resp.Body)
		return 500, "", "", *new(error)
	}
	return statusCode, token, isSA, nil
}

func (cc *CcClient) UserNamespaceCheck(token string, userName string, namespace string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/namespaces/%v", userName, namespace)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New("user cannot access the specified namespace")
	}
	return nil
}

func (cc *CcClient) PostNotebookRequestCheck(token string, namespace *models.NewNotebookRequest) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/check-notebook-request")

	bytes, err := json.Marshal(namespace)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(string(bytes)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New(string(body))
	}
	return nil
}

func (cc *CcClient) GetGUIDFromUserId(token string, userId string) (*string, *string, *string, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/users/name/%v", userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, nil, nil, err
	}
	var u = new(UserResponse)
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &u)
	if err != nil {
		logger.Logger().Errorf("parse resultMap Unmarshal failed. %v", err.Error())
		return nil, nil, nil,err
	}
	logger.Logger().Infof("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)
	if u.GUIDCheck && (u.GID == 0 || u.UID == 0) {
		return nil, nil, nil, errors.New("error GUIDCheck is on and GID or UID is 0")
	}
	gid := strconv.Itoa(u.GID)
	uid := strconv.Itoa(u.UID)

	return &gid, &uid, &u.Token, nil
}

func (cc *CcClient) AdminUserCheck(token string, adminUserId string, userId string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/users/%v", adminUserId, userId)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CcClient) UserStoragePathCheck(token string, username string, path string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/storages?path=%v", username, path)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CcClient) CheckNamespace(token string, namespace string) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/namespaces/%v", namespace)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}

	defer resp.Body.Close()
	return body, nil
}

func (cc *CcClient) CheckNamespaceUser(token string, namespace string, user string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/namespaces/%v/users/%v", namespace, user)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New(string(body))
	}
	return nil
}

func (cc *CcClient) CheckNamespacedNotebook(token string, namespace string, notebook string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/namespaces/%v/notebooks/%v", namespace, notebook)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New(string(body))
	}
	return nil
}

func (cc *CcClient) CheckNotebookPathCheck(token string, username string, namespace string, path string) (result []byte, error error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/namespaces/%v/storages?path=%v", username, namespace, path)
	res, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		error = err
		return res, err
	}
	result = res
	return res, err
}

func (cc *CcClient) CheckUserGetNamespace(token string, user string) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/users/%v", user)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}
	return body, nil
}

func (cc *CcClient) accessCheckFromCc(token string, apiUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	statusCode := resp.StatusCode
	if statusCode != 200 {
		return nil, errors.New("error from cc " + string(body))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}
