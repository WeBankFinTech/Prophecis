/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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
	"io/ioutil"
	"net/http"
)

type BaseClient struct {
	Address string
	Client  http.Client
}

type BaseResult struct {
	Message string `json:"message"`
	Method  string `json:"method"`
	Status  int    `json:"status"`
}

type NormalResult struct {
	Data json.RawMessage `json:"data"`
	BaseResult
}

func (c *BaseClient) DoLinkisHttpRequest(req *http.Request, user string) (*NormalResult, error) {
	//c.setAuthHeader(user, req)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if statusCode != 200 {
		return nil, errors.New("error from dss " + string(body))
	}
	if err != nil {
		return nil, err
	}

	var result NormalResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *BaseClient) DoHttpRequest(req *http.Request) ([]byte, error) {
	//c.setAuthHeader(user, req)
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	statusCode := res.StatusCode
	if statusCode != 200 {
		return nil, errors.New("error from server " + string(body))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *BaseClient) DoHttpDownloadRequest(req *http.Request, user string) ([]byte, error) {
	//c.setAuthHeader(user, req)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if statusCode != 200 {
		return nil, errors.New("error from dss " + string(body))
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}
