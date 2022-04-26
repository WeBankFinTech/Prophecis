package linkis

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"webank/DI/pkg/client"
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
			Address: "http://"+viper.GetString("linkis.address"),
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
	CostTime               int64   `json:"costTime"`
	CreatedTime            int64   `json:"createdTime"`
	EngineInstance         string  `json:"engineInstance"`
	EngineStartTime        int64   `json:"engineStartTime"`
	EngineType             int64   `json:"engineType"`
	ErrCode                int64   `json:"errCode"`
	ErrDesc                string  `json:"errDesc"`
	ExecId                 string  `json:"execId"`
	ExecuteApplicationName string  `json:"executeApplicationName"`
	ExecutionCode          string  `json:"executionCode"`
	Instance               string  `json:"instance"`
	LogPath                string  `json:"logPath"`
	ParamsJson             string  `json:"paramsJson"`
	Progress               float64 `json:"progress"`
	RequestApplicationName string  `json:"requestApplicationName"`
	ResultLocation         string  `json:"resultLocation"`
	RunType                string  `json:"runType"`
	SourceJson             string  `json:"sourceJson"`
	SourceTailor           string  `json:"sourceTailor"`
	Status                 string  `json:"status"`
	StrongerExecId         string  `json:"strongerExecId"`
	TaskID                 int64   `json:"taskID"`
	UmUser                 string  `json:"umUser"`
	UpdatedTime            int64   `json:"updatedTime"`
}

type ExecuteRequest struct {
	ExecuteApplicationName string                 `json:"executeApplicationName"`
	ExecutionCode          string                 `json:"executionCode"`
	Params                 map[string]interface{} `json:"params"`
	RequestApplicationName string                 `json:"requestApplicationName"`
	RunType                string                 `json:"runType"`
	Source                 SourceOfExecuteRequest `json:"source"`
}

type SourceOfExecuteRequest struct {
	FlowName    string `json:"flowName"`
	ProjectName string `json:"projectName"`
}

func (c *LinkisClient) Execute(executeRequest ExecuteRequest, user string) (*ExecuteData, error) {
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

func (c *LinkisClient) Status(execID, user string) (*StatusData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%v/status", execID)
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
	res := &JobHistory{}
	err = json.Unmarshal(result.Data[8:len(result.Data)-1], res)
	if err != nil {
		return "", err
	}
	return res.LogPath, nil
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

func (c *LinkisClient) Kill(execID, user string) (*string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%v/kill", execID)
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

func (c *LinkisClient) GetExecution(execID, user string) (*ExecutionData, error) {
	//requestPath := "/api/rest_j/v1/dss/get?id=1011&version=&projectVersionID=602"
	requestPath := fmt.Sprintf("/api/rest_j/v1/entrance/%v/execution", execID)
	url := c.Address + requestPath

	//init request
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
	var res ExecutionData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
