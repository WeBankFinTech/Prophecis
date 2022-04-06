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
package controlcenter

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mlss-mf/pkg/common/config"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type CcClient struct {
	ccUrl      string
	httpClient http.Client
}
type UserResponse struct {
	Name      string `json:"name"`
	GID       int    `json:"gid"`
	UID       int    `json:"uid"`
	Type      string `json:"type"`
	GUIDCheck bool   `json:"guidCheck"`
	Remarks   string `json:"remarks"`
}

const CcAuthToken = "MLSS-Token"

// FIXME MLSS Change: get models filter by username and namespace
const CcSuperadmin = "MLSS-Superadmin"
const CcAuthUser = "MLSS-UserID"

var once sync.Once
var client *CcClient

func GetCcClient(ccUrl string) *CcClient {
	once.Do(func() {
		client = &CcClient{}
		//client.httpClient = http.Client{}
		client.ccUrl = ccUrl

		//FIXME
		pool := x509.NewCertPool()
		caCertPath := config.GetMFconfig().ServerConfig.CaCertPath
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

// FIXME MLSS Change: get models filter by username and namespace
func (cc *CcClient) AuthAccessCheck(authOptions *map[string]string) (int, string, string, error) {
	apiUrl := cc.ccUrl + "/cc/v1/inter/user"
	//apiUrl := cc.ccUrl + "/cc/v1/sample"

	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		logger.Logger().Error("Create AuthAccessCheck Request Failed !", err.Error())
		return 500, "", "", err
	}
	for key, value := range *authOptions {
		req.Header.Set(key, value)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		logger.Logger().Error("Do AuthAccessCheck Request Failed !", err.Error())
		return 500, "", "", err
	}
	defer resp.Body.Close()
	token := resp.Header.Get(CcAuthToken)
	isSA := resp.Header.Get(CcSuperadmin)
	statusCode := resp.StatusCode
	if statusCode == 401 || (statusCode >= 200 && statusCode < 300 && token == "") {
		logger.Logger().Error("Auth Access Check Request Failed, StatusCode: %v, Token: %v", statusCode, token)
		logger.Logger().Error("Message: %v", resp.Body)
		return 401, "", "", *new(error)
	}
	if statusCode > 300 || token == "" {
		logger.Logger().Error("Auth Access Check Request Failed from cc, StatusCode: %v, Token: %v", statusCode, token)
		logger.Logger().Error("Message: %v", resp.Body)
		return 500, "", "", *new(error)
	}
	return statusCode, token, isSA, nil
}

func (cc *CcClient) UserStorageCheck(token string, username string, path string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/storages?path=%v", username, path)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
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

func (cc *CcClient) GetGUIDFromUserId(token string, userId string) (*string, *string, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/users/name/%v", userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, nil, err
	}
	var u = new(UserResponse)
	//json.Unmarshal(body, &u)
	err = GetResultData(body, &u)
	if err != nil {
		return nil, nil, errors.New("GetResultData failed")
	}
	log.Debugf("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)
	if u.GUIDCheck && (u.GID == 0 || u.UID == 0) {
		return nil, nil, errors.New("error GUIDCheck is on and GID or UID is 0")
	}
	gid := strconv.Itoa(u.GID)
	uid := strconv.Itoa(u.UID)

	return &gid, &uid, nil
}

func (cc *CcClient) AdminUserCheck(token string, adminUserId string, userId string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/users/%v", adminUserId, userId)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CcClient) UserStoragePathCheck(token string, username string, namespace string, path string) (result []byte, error error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/namespaces/%v/storages?path=%v", username, namespace, path)
	res, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		error = err
		return res, err
	}
	result = res
	return res, err
}

func (cc *CcClient) accessCheckFromCc(token string, apiUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, doErr := client.httpClient.Do(req)
	statusCode := resp.StatusCode
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if doErr != nil {
		return body, doErr
	}
	defer resp.Body.Close()
	//statusCode := resp.StatusCode
	if statusCode != 200 {
		return body, errors.New("error from cc ")
	}
	//body, err := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		return body, bodyErr
	}
	return body, nil
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

func (cc *CcClient) GetUserGroups(token string) ([]*gormmodels.Group, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/usergroups")
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
	groups := []*gormmodels.Group{}
	statusCode := resp.StatusCode
	err = json.Unmarshal(body, &groups)
	if err != nil {
		return nil, err
	}
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}
	return groups, nil
}

func (cc *CcClient) CheckGroupByUser(token string, userId string) (int64, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/group/%v", userId)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return 0, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return 0, err
	}
	var group gormmodels.ResultGroup
	statusCode := resp.StatusCode
	err = json.Unmarshal(body, &group)
	if err != nil {
		return 0, err
	}
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return 0, errors.New(string(body))
	}
	return group.GroupID, nil
}

func (cc *CcClient) GetCurrentUserNamespaceWithRole(token string, roleId string) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/groups/users/roles/%v/namespaces", roleId)
	log.Debugf("get apiUrl: %v", apiUrl)
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

func (cc *CcClient) GetGroupsById(token string, groupId int64) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/groups/id/%v", groupId)
	log.Debugf("get apiUrl: %v", apiUrl)
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

func (cc *CcClient) CheckResource(namespace string, cpu float64, gpu string, memory float64, token string) (bool, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/%v/checkresource?cpu=%v&gpu=%v&memory=%v", namespace, cpu, gpu, memory)
	logger.Logger().Debugf("get apiUrl: %v", apiUrl)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		logger.Logger().Error("New check resource request err, ", err)
		return false, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		logger.Logger().Error("Check resource do http request err, ", err)
		return false, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logger.Logger().Error("Read response body err, ", err)
		return false, err
	}
	statusCode := resp.StatusCode
	defer resp.Body.Close()
	if statusCode == 200 {
		logger.Logger().Info("Check resource success")
		return true, nil
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		return false, err
	}
	return false, errors.New(jsonMap["error"].(string))
}
