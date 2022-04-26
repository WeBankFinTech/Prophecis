package dss1_0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"net/http"
	netUrl "net/url"
	"os"
	"strings"
	"webank/DI/commons/logger"
	"webank/DI/pkg/client"
	"webank/DI/restapi/api_v1/restmodels"
)

const (
	headerTokenUser = "Token-User"
	headerTokenCode = "Token-Code"
)

type DSSClient struct {
	client.BaseClient
	TokenCode string
}

type Response struct {
	Method  string `json:"method"`
	Status  int64 `json:"status"`
	Message string `json:"message"`
}

func GetDSSClient() DSSClient {

	httpClient := http.Client{}

	dssClient := DSSClient{
		TokenCode: viper.GetString("linkis.tokenCode"),
		BaseClient: client.BaseClient{
			Address: "http://"+viper.GetString("linkis.address"),
			Client:  httpClient,
		},
	}
	return dssClient
}

func GetDSSProClient() DSSClient {

	httpClient := http.Client{}

	dssClient := DSSClient{
		TokenCode: viper.GetString("linkispro.tokenCode"),
		BaseClient: client.BaseClient{
			Address: "http://"+viper.GetString("linkispro.address"),
			Client:  httpClient,
		},
	}
	return dssClient
}

type WorkspacesData struct {
	Workspaces []Workspace `json:"workspaces"`
}

type WorkspacesResponse struct {
	Response
	Data WorkspacesData `json:"data"`
}

type Workspace struct {
	CreateBy   string `json:"createBy"`
	Department string `json:"department"`
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Product    string `json:"product"`
	LastUpdateUser string `json:"lastUpdateUser"`
}

var log = logger.GetLogger()

func (c *DSSClient) GetWorkspaces(user string) ([]Workspace, error) {
	requestPath := "/api/rest_j/v1/dss/framework/workspace/workspaces"
	url := c.Address + requestPath
	// do http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var res WorkspacesData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return res.Workspaces, nil
}

func (c *DSSClient) setAuthHeader(user string, req *http.Request) {
	req.Header.Set(headerTokenUser, user)
	req.Header.Set(headerTokenCode, c.TokenCode)
}

type CreateOrchestratorRequest struct {
	Description      string `json:"description"`
	OrchestratorName string `json:"orchestratorName"`
	Uses             string `json:"uses"`
	OrchestratorMode string  `json:"orchestratorMode"`
	OrchestratorWays []string `json:"orchestratorWays"`
	ProjectId      int64  `json:"projectId"`
	WorkspaceId      string  `json:"workspaceId"`
	Labels      DSSLabel  `json:"labels"`
}

type CreateOrchestratorResponse struct {
	OrchestratorId int64 `json:"orchestratorId"`
}

func (c *DSSClient) CreateOrchestrator(params CreateOrchestratorRequest, user string) (*CreateOrchestratorResponse, error) {
	requestPath := "/api/rest_j/v1/dss/framework/project/createOrchestrator"

	url := c.Address + requestPath

	bytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	//init request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	log.Info("Debug Log address:%s" + url)
	if err != nil {
		log.Error("address:%s" + url)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		log.Error("Do Request Error: " + err.Error())
		return nil, err
	}
	var res CreateOrchestratorResponse
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}



type SaveWorkflowParams struct {
	Id               int64  `json:"id"`
	Json             string `json:"json"`
	WorkspaceName string  `json:"workspaceName"`
	ProjectName string  `json:"projectName"`
	Label DSSLabel `json:"labels"`
	//FlowEditLock string `json:"flowEditLock"`
	//IsNotHaveLock bool `json:"isNotHaveLock"`
}

type DSSLabel struct {
	Route string `json:"route"`
}

type SaveFlowResponse struct {
	Response
	Data SaveFlowData `json:"flowEditLock"`
}

type SaveFlowData struct {
	FlowEditLock string `json:"flowEditLock"`
	FlowVersion  string `json:"flowVersion"`
}

type CreateWorkflowResponse struct {
	Response
	Flow CreateWorkflowData `json:"data"`
}

type CreateWorkflowData struct {
	Flow Flow `json:"flow"`
}

func (c *DSSClient) SaveFlow(params SaveWorkflowParams, user string) (*SaveFlowData, error) {
	requestPath := "/api/rest_j/v1/dss/workflow/saveFlow"
	url := c.Address + requestPath

	bytes, err := json.Marshal(params)
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
	var res SaveFlowData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type OrchestratorGetResponse struct {
	Response
	Data    OrchestratorData `json:"data"`
}

type OrchestratorData struct {
	Flow Flow `json:"flow"`
}

type Flow struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	CreateTime    int64          `json:"createTime"`
	Creator       string         `json:"creator"`
	Description   string         `json:"description"`
	Rank          int64          `json:"rank"`
	ProjectID     int64          `json:"projectID"`
	LinkisAppConnName    string   `json:"projectID"`
	Labels    string          `json:"dssLabels"`
	Uses          string         `json:"uses"`
	ResourceId string          `json:"resourceId"`
	BmlVersion string          `json:"bmlVersion"`
	Version string           `json:"version"`
	FlowJson string           `json:"flowJson"`
	FlowType      string         `json:"flowType"`
	HasSaved      bool           `json:"hasSaved"`
	RootFlow      bool           `json:"rootFlow"`
	Source        string         `json:"source"`
	State         bool           `json:"state"`
	//FlowEditLock string 		`json:"flowEditLock"`
}

func (c *DSSClient) GetOrchestrator( flowId int64, labels string, user string) (*OrchestratorData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/workflow/get?flowId=%v&labels=%v&isNotHaveLock=%v",
		flowId, labels, true)
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
	var orchestratorData OrchestratorData
	err = json.Unmarshal(result.Data, &orchestratorData)
	if err != nil {
		log.Logger.Error("unmarshal error , response data:"+string(result.Data))
		return nil, err
	}

	return &orchestratorData, nil
}


type UploadData struct {
	UploadFile 		*os.File `json:"file"`
	UploadSystem	string	`json:"system"`
}

type UploadResponse struct {
	ResourceId 		string  `json:"resourceId"`
	Version			string `json:"version"`
	TaskId			int64  `json:"taskId"`
}

func (c *DSSClient) Upload(uploadData *UploadData, user string) (*UploadResponse, error) {
	requestPath := "/api/rest_j/v1/bml/upload"
	url := c.Address + requestPath
	bodyBuffer := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuffer)
	_, err := bodyWriter.CreateFormFile("file", uploadData.UploadFile.Name())
	boundary := bodyWriter.Boundary()
	fileInfo, _ := uploadData.UploadFile.Stat()
	closeBuffer := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	requestReader := io.MultiReader(bodyBuffer, uploadData.UploadFile, closeBuffer)
	request, err := http.NewRequest("POST", url, requestReader)
	if err != nil{
		log.Println("upload to dss failed, ",err)
		return nil,err
	}
	request.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data;boundary=%v",boundary))
	request.ContentLength = fileInfo.Size() + int64(bodyBuffer.Len()) + int64(closeBuffer.Len())
	params := make(netUrl.Values)
	params.Set("system","MLSS")
	request.PostForm = params
	c.setAuthHeader(user, request)
	result, err := c.DoLinkisHttpRequest(request, user)
	if err != nil{
		log.Println("do http request(upload dss) failed, ",err)
		return nil,err
	}
	var res UploadResponse
	err = json.Unmarshal(result.Data, &res)
	if err != nil{
		log.Println("upload response un marshal failed, ",err)
		return nil,err
	}
	return &res, nil
}

type DowloadData struct {
	ResourceId string `json:"resourceId"`
	Version		  string `json:"version"`
}

func (c *DSSClient) Download(downloadData *DowloadData, user string) ([]byte, error) {
	//requestPath := fmt.Sprintf("/api/rest_j/v1/filesystem/downloadShareResource?resourceId=%v&version=%v",downloadData.ResourceId,downloadData.Version)
	requestPath := fmt.Sprintf("/api/rest_j/v1/bml/downloadShareResource?resourceId=%v&version=%v",downloadData.ResourceId,downloadData.Version)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("dss download request failed, ",err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)
	result, err := c.DoHttpDownloadRequest(req, user)
	if err != nil{
		log.Println("dss download request failed, ",err)
		return nil, err
	}
	return result, nil
}

func (c *DSSClient) GetUserInfo(cookies []*http.Cookie) (*restmodels.DssUserInfo, error) {
	requestPath := "/api/rest_j/v1/user/userInfo"
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("dss  get user info request failed, ",err)
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := c.DoLinkisHttpRequest(req, "")
	if err != nil{
		log.Println("dss  get user info request failed, ",err)
		return nil, err
	}
	var dssUserInfo *restmodels.DssUserInfo
	err = json.Unmarshal(result.Data, &dssUserInfo)
	if err != nil{
		log.Println("dss  get user info request failed, ",err)
		return nil, err
	}
	return dssUserInfo, nil
}

type ProjectData struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Business        string   `json:"business"`
	ApplicationArea int64    `json:"applicationArea"`
	Product         string   `json:"product"`
	CreateBy         string   `json:"createBy"`
	Source         string   `json:"source"`
	ReleaseUsers       []string `json:"releaseUsers"`
	EditUsers       []string `json:"editUsers"`
	AccessUsers     []string `json:"accessUsers"`
	DevProcessList     []string `json:"devProcessList"`
	OrchestratorModeList     []string `json:"orchestratorModeList"`
	WorkspaceId     int64    `json:"workspaceId"`
	Id     int64    `json:"Id"`
}

type ProjectParam struct {
	ProjectName    string
	ProjectDesc    string
	ProjectProduct string
}

type CreateProjectRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Business        string   `json:"business"`
	ApplicationArea int64    `json:"applicationArea"`
	Product         string   `json:"product"`
	EditUsers       []string `json:"editUsers"`
	AccessUsers     []string `json:"accessUsers"`
	DevProcessList     []string `json:"devProcessList"`
	OrchestratorModeList     []string `json:"orchestratorModeList"`
	WorkspaceId     int64    `json:"workspaceId"`
}

func (c * DSSClient) CreateProject(projectParam *ProjectParam, username string, devProcessList []string ,workspaceId int64) (*ProjectData, error) {
	requestPath := "/api/rest_j/v1/dss/framework/project/createProject"
	url := c.Address + requestPath
	defaultOrchestratorModeList := []string{"pom_work_flow"}
	defaultDevProcessList := []string{"dev","prod"}
	if devProcessList == nil {
		devProcessList = defaultDevProcessList
	}
	createProjectRequest := &CreateProjectRequest{
		Name:            projectParam.ProjectName,
		Description:     projectParam.ProjectDesc,
		Product:         projectParam.ProjectProduct,
		WorkspaceId:     workspaceId,
		ApplicationArea: 4,
		Business:        "",
		OrchestratorModeList: defaultOrchestratorModeList,
		DevProcessList: defaultDevProcessList,
		EditUsers:       []string{},
		AccessUsers:     []string{},
	}
	requestJsonBytes, err := json.Marshal(&createProjectRequest)
	if err != nil {
		log.Println("dss create  projects request failed, ", err)
		return nil, err
	}
	req, err := http.NewRequest("POST",url,bytes.NewBuffer(requestJsonBytes))
	if err != nil {
		log.Println("dss create projects request failed, ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		return nil, err
	}

	var projectData ProjectData
	err = json.Unmarshal(result.Data, &projectData)
	if err != nil {
		log.Println("dss create cooperate projects request failed, ", err)
		return nil, err
	}

	return &projectData,nil
}

type GetProjectParams struct {
	WorkspaceId   int64  `json:"workspaceId"`
}

func (c * DSSClient) GetProjectNameById(projectId int64, username string, WorkspaceId int64) (string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	url := c.Address + requestPath
	getProjectParams := &GetProjectParams{
		WorkspaceId: WorkspaceId,
	}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET",url,strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get projects request failed, ",err)
		return "", err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return "", err
	}
	projectList := []map[string]interface{}{}
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return "", err
	}
	var projectName string
	for _, project := range projectList {
		if project["id"].(float64) == float64(projectId){
			projectName = project["name"].(string)
			break
		}
	}
	return projectName, nil
}


func (c * DSSClient) GetProjectIdByName(projectName string, username string, WorkspaceId int64) (int64, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	url := c.Address + requestPath
	getProjectParams := &GetProjectParams{
		WorkspaceId: WorkspaceId,
	}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST",url,strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get projects request failed, ",err)
		return -1, err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return -1, err
	}
	projectList := []map[string]interface{}{}
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return -1, err
	}
	var projectID float64
	for _, project := range projectList {
		if project["name"].(string) == projectName{
			projectID = project["id"].(float64)
			break
		}
	}
	return int64(projectID), nil
}

type project struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Business        string   `json:"business"`
	ApplicationArea int64    `json:"applicationArea"`
	Product         string   `json:"product"`
	EditUsers       []string `json:"editUsers"`
	AccessUsers     []string `json:"accessUsers"`
	DevProcessList     []string `json:"devProcessList"`
	OrchestratorModeList     []string `json:"orchestratorModeList"`
	WorkspaceId     int64    `json:"workspaceId"`
}








func (c * DSSClient) CheckProjectExist(workspaceId int64, projectName string, username string) (bool,*ProjectData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	isExisted := false
	url := c.Address + requestPath


	getProjectParams := &GetProjectParams{
		WorkspaceId: workspaceId,
	}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return isExisted, nil, err
	}

	req, err := http.NewRequest("POST",url,strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return isExisted, nil , err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get projects request failed, ",err)
		return isExisted, nil, err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil{
		log.Println("dss  get  projects request failed, ",err)
		return false, nil , err
	}

	projectList := []ProjectData{}
	var project ProjectData
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return isExisted, nil , err
	}
	for _, p := range projectList {
		if p.Name == projectName{
			isExisted = true
			project = p
			break
		}
	}
	return isExisted, &project, nil

}


type openOrchestratorRequest struct {
	OrchestratorId   int64 `json:"orchestratorId"`
	WorkspaceName      string `json:"workspaceName"`
	Labels      DSSLabel  `json:"labels"`
}


type openOrchestratorResponse struct {
	OrchestratorOpenUrl string `json:"OrchestratorOpenUrl"`
	OrchestratorVo OrchestratorVo `json:"OrchestratorVo"`
}
type OrchestratorVo struct {
	DssOrchestratorInfo dssOrchestratorInfo  `json:"dssOrchestratorInfo"`
	DssOrchestratorVersion dssOrchestratorVersion  `json:"dssOrchestratorVersion"`
}


type dssOrchestratorVersion struct {
	Id            int64          `json:"id"`
	OrchestratorId          int64         `json:"orchestratorId"`
	ProjectId     int64          `json:"projectId"`
	AppId         int64          `json:"appId"`
	Source         string          `json:"source"`
	Version    string          `json:"Version"`
	Comment    string          `json:"Version"`
	Creator       string         `json:"comment"`
	Content       string         `json:"content"`
	ContextId     string          `json:"contextId"`
}


type dssOrchestratorInfo struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	CreateTime    int64          `json:"createTime"`
	Creator       string         `json:"creator"`
	Type          string         `json:"type"`
	Uses          string         `json:"uses"`
	AppConnName   string         `json:"appConnName"`
	ProjectId     int64          `json:"projectId"`
	UUID          string           `json:"uuid"`
	SecondaryType string  `json:"secondaryType"`
	Description   string         `json:"description"`
}


func (c *DSSClient) OpenOrchestrator( orchestratorId int64, labels DSSLabel, workspaceName string	, user string) (*OrchestratorVo, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/orchestrator/openOrchestrator")
	url := c.Address + requestPath

	or := openOrchestratorRequest{
		OrchestratorId: orchestratorId,
		Labels: labels,
		WorkspaceName: workspaceName,
	}

	bytes, err := json.Marshal(or)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST",url,strings.NewReader(string(bytes)))
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var orchestratorResponse openOrchestratorResponse
	err = json.Unmarshal(result.Data, &orchestratorResponse)
	if err != nil {
		return nil, err
	}

	return &orchestratorResponse.OrchestratorVo, nil
}








