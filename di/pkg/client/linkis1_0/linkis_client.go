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
package linkis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"webank/DI/pkg/client"
	dss "webank/DI/pkg/client/dss1_0"
)

const (
	headerTokenUser = "Token-User"
	headerTokenCode = "Token-Code"
)

type LinkisClient struct {
	client.BaseClient
	TokenCode string
}

func GetLinkisClient() LinkisClient {

	httpClient := http.Client{}

	linkisClient := LinkisClient{
		TokenCode: viper.GetString("linkis.tokenCode"),
		BaseClient: client.BaseClient{
			//Address: "http://127.0.0.1:8088",
			Address: "http://" + viper.GetString("linkis.address"),
			Client:  httpClient,
		},
	}
	return linkisClient
}

type Response struct {
	Message string `json:"message"`
	Method  string `json:"method"`
	Status  int    `json:"status"`
}

type ExecuteData struct {
	ExecID string `json:"execID"`
	TaskID int64  `json:"taskID"`
}
type KillData struct {
	ExecID string `json:"execID"`
}

type LogResponse struct {
	Response
	data LogData `json:"data"`
}
type LogData struct {
	ExecID   string   `json:"execID"`
	Log      []string `json:"log"`
	FromLine int64    `json:"fromLine"`
}
type JobHistory struct {
	CostTime               int64  `json:"costTime"`
	CreatedTime            int64  `json:"createdTime"`
	EngineInstance         string `json:"engineInstance"`
	EngineStartTime        int64  `json:"engineStartTime"`
	EngineType             string `json:"engineType"`
	ErrCode                int64  `json:"errCode"`
	ErrDesc                string `json:"errDesc"`
	ExecId                 string `json:"execId"`
	ExecuteApplicationName string `json:"executeApplicationName"`
	ExecutionCode          string `json:"executionCode"`
	Instance               string `json:"instance"`
	LogPath                string `json:"logPath"`
	ParamsJson             string `json:"paramsJson"`
	Progress               string `json:"progress"`
	RequestApplicationName string `json:"requestApplicationName"`
	ResultLocation         string `json:"resultLocation"`
	RunType                string `json:"runType"`
	SourceJson             string `json:"sourceJson"`
	SourceTailor           string `json:"sourceTailor"`
	Status                 string `json:"status"`
	StrongerExecId         string `json:"strongerExecId"`
	TaskID                 int64  `json:"taskID"`
	UmUser                 string `json:"umUser"`
	UpdatedTime            int64  `json:"updatedTime"`
}

type JobHistoryTask struct {
	Task JobHistory `json:"task"`
}

type ExecuteRequest struct {
	ExecuteApplicationName string                 `json:"executeApplicationName"`
	ExecutionCode          string                 `json:"executionCode"`
	Params                 map[string]interface{} `json:"params"`
	RequestApplicationName string                 `json:"requestApplicationName"`
	Label                  dss.DSSLabel           `json:"labels"`
	RunType                string                 `json:"runType"`
	Source                 SourceOfExecuteRequest `json:"source"`
}

type SourceOfExecuteRequest struct {
	FlowName    string `json:"flowName"`
	ProjectName string `json:"projectName"`
}

func (c *LinkisClient) FlowExecute(executeRequest ExecuteRequest, user string) (*ExecuteData, error) {
	requestPath := "/api/rest_j/v1/dss/flow/entrance/execute"
	url := c.Address + requestPath

	bytes, err := json.Marshal(executeRequest)
	if err != nil {
		return nil, err
	}
	//init request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}

	var res ExecuteData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type StatusData struct {
	ExecID string `json:"execID"`
	Status string `json:"status"`
}

func (c *LinkisClient) FlowStatus(execID string, taskID int64, user string) (*StatusData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/flow/entrance/%v/status?taskID=%v&labels=dev", execID, taskID)
	url := c.Address + requestPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}
	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}

	var res StatusData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *LinkisClient) GetLogPath(taskId int64, user string) (string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/jobhistory/%v/get", taskId)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return "", err
	}

	if result.Status != 0 {
		return "", errors.New(result.Message)
	}

	log.Printf("GetLogPath result: %v", result)

	res := &JobHistoryTask{}
	err = json.Unmarshal(result.Data, res)
	log.Printf("GetLogPath Task: %v", res)

	if err != nil {
		return "", err
	}
	return res.Task.LogPath, nil
}

func (c *LinkisClient) GetOpenLog(logPath string, user string) (*LogData, error) {
	requestPath := "/api/rest_j/v1/filesystem/openLog?path=" + logPath
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	res := &LogData{}
	err = json.Unmarshal(result.Data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *LinkisClient) Log(execID string, fromLine int64, size int64, user string) (*LogData, error) {
	//requestPath := "/api/rest_j/v1/entrance//status?taskID=29842"
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%v/log?fromLine=%v&size=%v", execID, fromLine, size)
	url := c.Address + requestPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}
	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}

	var res LogData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *LinkisClient) Kill(execID string, taskID int64, user string) (*string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%v/kill?taskID=%v&labels=dev", execID, taskID)
	url := c.Address + requestPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var res KillData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res.ExecID, nil
}

func (c *LinkisClient) KillWorkflow(execID string, taskID int64, user string) (*string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/flow/entrance/%v/killWorkflow?taskID=%v&labels=dev",
		execID,taskID )
	url := c.Address + requestPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var res KillData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res.ExecID, nil
}

func (c *LinkisClient) setAuthHeader(user string, req *http.Request) {
	req.Header.Set(headerTokenUser, user)
	req.Header.Set(headerTokenCode, c.TokenCode)
}

type ExecutionData struct {
	SucceedJobs []Job `json:"succeedJobs"`
	RunningJobs []Job `json:"runningJobs"`
	PendingJobs []Job `json:"pendingJobs"`
	SkippedJobs []Job `json:"skippedJobs"`
	FailedJobs  []Job `json:"failedJobs"`
}

type Job struct {
	StartTime int64   `json:"startTime"`
	NodeID    string  `json:"nodeID"`
	TaskID    string  `json:"taskID"`
	ExecID    string  `json:"execID"`
	Info      *string `json:"info"`
	nowTime   int64   `json:"nowTime"`
}

func (c *LinkisClient) GetFlowExecution(execID, user string) (*ExecutionData, error) {
	//requestPath := "/api/rest_j/v1/dss/get?id=1011&version=&projectVersionID=602"
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/flow/entrance/%v/execution?labels=dev", execID)
	url := c.Address + requestPath

	//init request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var res ExecutionData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *LinkisClient) Status(execID, user string) (*StatusData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%s/status", execID)
	url := c.Address + requestPath

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}

	var res StatusData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type ExecuteReq struct {
	ExecuteApplicationName string                                  `json:"executeApplicationName"`
	ExecutionCode          string                                  `json:"executionCode"`
	Params                 map[string]map[string]map[string]string `json:"params"`
	Label                  map[string]string                       `json:"labels"`
	RunType                string                                  `json:"runType"`
}

func (c *LinkisClient) Execute(executeRequest ExecuteReq, user string) (*ExecuteData, error) {
	requestPath := "/api/rest_j/v1/entrance/execute"
	url := c.Address + requestPath
	bytes, err := json.Marshal(executeRequest)
	if err != nil {
		return nil, err
	}
	//init request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	cl := dss.GetDSSClient()
	if err := cl.SetCookieHeaderLink(req, user); err != nil {
		return nil, err
	}
	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}

	var res ExecuteData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
