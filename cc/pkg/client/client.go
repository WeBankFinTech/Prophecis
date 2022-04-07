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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"net/http"
	"strings"
	"sync"
)

type AideClient struct {
	aideUrl    string
	httpClient http.Client
}

var once sync.Once
var client *AideClient

func GetAideClient(aideUrl string) *AideClient {
	once.Do(func() {
		client = &AideClient{}
		client.aideUrl = aideUrl
		client.httpClient = http.Client{}
	})
	return client
}

func (aidec *AideClient) GetAideProxyUser(headerUserId, headerToken string, namespace string, notebookName string) (string, string, error) {
	apiUrl := aidec.aideUrl + fmt.Sprintf("/aide/v1/notebook/user/%v/%v", namespace, notebookName)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		logger.Logger().Error("GetAideProxyUser request failed, ", err.Error())
		return "", "", err
	}
	req.Header.Set(constants.AUTH_HEADER_USERID, headerUserId)
	req.Header.Set("token", headerToken)
	//req.Header.Set(constants.AUTH_HEADER_TOKEN, headerToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		logger.Logger().Error("Do GetAideProxyUser Request Failed !")
		return "", "", err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if err != nil {
		logger.Logger().Error("Do GetAideProxyUser Request Failed !")
		return "", "", err
	}
	if statusCode != 200 {
		logger.Logger().Error("Message: ", string(resBody))
		return "", "", errors.New(string(resBody))
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(resBody, &jsonMap)
	if err != nil {
		logger.Logger().Error("Json unmarshal err: ", err)
		return "", "", err
	}
	return jsonMap["ProxyUser"].(string), jsonMap["User"].(string), nil
}
