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
package dss1_0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	netUrl "net/url"
	"os"
	"strconv"
	"strings"
	"webank/DI/commons/logger"
	"webank/DI/pkg/client"
	"webank/DI/restapi/api_v1/restmodels"

	"github.com/spf13/viper"
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
	Status  int64  `json:"status"`
	Message string `json:"message"`
}

func GetDSSClient() DSSClient {

	httpClient := http.Client{}

	dssClient := DSSClient{
		//client.BaseClient{
		//	Address: "127.0.0.1:8088",
		//	Client:  httpClient,
		//},
		TokenCode: viper.GetString("linkis.tokenCode"),
		BaseClient: client.BaseClient{
			Address: "http://" + viper.GetString("linkis.address"),
			Client:  httpClient,
		},
	}
	return dssClient
}

func GetDSSProClient() DSSClient {

	httpClient := http.Client{}

	dssClient := DSSClient{
		//client.BaseClient{
		//	Address: "127.0.0.1:8088",
		//	Client:  httpClient,
		//},
		TokenCode: viper.GetString("linkispro.tokenCode"),
		BaseClient: client.BaseClient{
			Address: "http://" + viper.GetString("linkispro.address"),
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
	CreateBy       string `json:"createBy"`
	Department     string `json:"department"`
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Product        string `json:"product"`
	LastUpdateUser string `json:"lastUpdateUser"`
}

const (
	DefaultWorkspace     = "bdapWorkspace"
	actualFlowFolderName = "flow"
)

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

func (c *DSSClient) GetDefaultWorkspaces(user string) (*Workspace, error) {
	var res Workspace
	workspaceList, err := c.GetWorkspaces(user)
	if err != nil {
		log.Errorf(err.Error())
		return &res, err
	}
	for _, ws := range workspaceList {
		if ws.Name == DefaultWorkspace {
			res = ws
			break
		}
	}
	return &res, nil
}

func (c *DSSClient) setAuthHeader(user string, req *http.Request) {
	req.Header.Set(headerTokenUser, user)
	req.Header.Set(headerTokenCode, c.TokenCode)
}

type CreateOrchestratorRequest struct {
	Description      string   `json:"description"`
	OrchestratorName string   `json:"orchestratorName"`
	Uses             string   `json:"uses"`
	OrchestratorMode string   `json:"orchestratorMode"`
	OrchestratorWays []string `json:"orchestratorWays"`
	ProjectId        int64    `json:"projectId"`
	WorkspaceId      string   `json:"workspaceId"`
	Labels           DSSLabel `json:"labels"`
}

type CreateOrchestratorResponse struct {
	OrchestratorId int64 `json:"orchestratorId"`
}

type CookieHeaderParams struct {
	WorkspaceName string
	WorkspaceId string
}

func (c *DSSClient) setCookieHeader(req *http.Request, user string) error{
	workspace, err := c.GetDefaultWorkspaces(user)
	if err != nil {
		log.Error("GetDefaultWorkspaces:%s")
		return err
	}

	cookieWorkspaceName := &http.Cookie{Name: "workspaceName",Value: DefaultWorkspace, HttpOnly: true}
	cookieWorkspaceId := &http.Cookie{Name: "workspaceId", Value: strconv.Itoa(int(workspace.Id)), HttpOnly: true}
	req.AddCookie(cookieWorkspaceName)
	req.AddCookie(cookieWorkspaceId)
	return nil
}

func (c *DSSClient) SetCookieHeaderLink(req *http.Request, user string) error{
	workspace, err := c.GetDefaultWorkspaces(user)
	if err != nil {
		log.Error("GetDefaultWorkspaces:%s")
		return err
	}
	cookieWorkspaceName := &http.Cookie{Name: "workspaceName",Value: DefaultWorkspace, HttpOnly: true}
	cookieWorkspaceId := &http.Cookie{Name: "workspaceId", Value: strconv.Itoa(int(workspace.Id)), HttpOnly: true}
	req.AddCookie(cookieWorkspaceName)
	req.AddCookie(cookieWorkspaceId)
	return nil
}

func (c *DSSClient) CreateOrchestrator(params CreateOrchestratorRequest, user string) (*CreateOrchestratorResponse, error) {
	requestPath := "/api/rest_j/v1/dss/framework/orchestrator/createOrchestrator"

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

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, user); err != nil {
		return nil, err
	}

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
	Id            int64    `json:"id"`
	Json          string   `json:"json"`
	WorkspaceName string   `json:"workspaceName"`
	ProjectName   string   `json:"projectName"`
	Label         DSSLabel `json:"labels"`
	FlowEditLock  string   `json:"flowEditLock"`
	NotHaveLock   bool     `json:"notHaveLock"`
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
	log.Infof("SaveFlow1 url:%v", url)

	bytes, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	log.Infof("SaveFlow jsonParam: %v", string(bytes))
	//init request
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, user); err != nil {
		return nil, err
	}

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
	Data OrchestratorData `json:"data"`
}

type OrchestratorData struct {
	Flow Flow `json:"flow"`
}

type Flow struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	CreateTime        int64  `json:"createTime"`
	Creator           string `json:"creator"`
	Description       string `json:"description"`
	Rank              int64  `json:"rank"`
	ProjectID         int64  `json:"projectID"`
	LinkisAppConnName string `json:"projectID"`
	Labels            string `json:"dssLabels"`
	Uses              string `json:"uses"`
	ResourceId        string `json:"resourceId"`
	BmlVersion        string `json:"bmlVersion"`
	Version           string `json:"version"`
	FlowJson          string `json:"flowJson"`
	FlowType          string `json:"flowType"`
	HasSaved          bool   `json:"hasSaved"`
	RootFlow          bool   `json:"rootFlow"`
	Source            string `json:"source"`
	State             bool   `json:"state"`
	FlowEditLock      string `json:"flowEditLock"`
}

func (c *DSSClient) GetOrchestrator(flowId int64, labels string, user string) (*OrchestratorData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/workflow/get?flowId=%v&labels=%v&isNotHaveLock=%v",
		flowId, labels, true)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, user)
	if err != nil {
		return nil, err
	}
	var orchestratorData OrchestratorData
	err = json.Unmarshal(result.Data, &orchestratorData)
	if err != nil {
		log.Logger.Error("unmarshal error , response data:" + string(result.Data))
		return nil, err
	}

	return &orchestratorData, nil
}

type UploadData struct {
	UploadFile   *os.File `json:"file"`
	UploadSystem string   `json:"system"`
}

type UploadResponse struct {
	ResourceId string `json:"resourceId"`
	Version    string `json:"version"`
	TaskId     int64  `json:"taskId"`
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
	if err != nil {
		log.Println("upload to dss failed, ", err)
		return nil, err
	}
	request.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data;boundary=%v", boundary))
	request.ContentLength = fileInfo.Size() + int64(bodyBuffer.Len()) + int64(closeBuffer.Len())
	params := make(netUrl.Values)
	params.Set("system", "MLSS")
	request.PostForm = params
	c.setAuthHeader(user, request)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(request, user); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(request, user)
	if err != nil {
		log.Println("do http request(upload dss) failed, ", err)
		return nil, err
	}
	var res UploadResponse
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		log.Println("upload response un marshal failed, ", err)
		return nil, err
	}
	return &res, nil
}

type DowloadData struct {
	ResourceId string `json:"resourceId"`
	Version    string `json:"version"`
}

func (c *DSSClient) Download(downloadData *DowloadData, user string) ([]byte, error) {
	//requestPath := fmt.Sprintf("/api/rest_j/v1/filesystem/downloadShareResource?resourceId=%v&version=%v",downloadData.ResourceId,downloadData.Version)
	requestPath := fmt.Sprintf("/api/rest_j/v1/bml/downloadShareResource?resourceId=%v&version=%v", downloadData.ResourceId, downloadData.Version)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("dss download request failed, ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, user); err != nil {
		return nil, err
	}

	result, err := c.DoHttpDownloadRequest(req, user)
	if err != nil {
		log.Println("dss download request failed, ", err)
		return nil, err
	}
	return result, nil
}

func (c *DSSClient) GetUserInfo(cookies []*http.Cookie) (*restmodels.DssUserInfo, error) {
	requestPath := "/api/rest_j/v1/user/userInfo"
	url := c.Address + requestPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("dss  get user info request failed, ", err)
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := c.DoLinkisHttpRequest(req, "")
	if err != nil {
		log.Println("dss  get user info request failed, ", err)
		return nil, err
	}
	var dssUserInfo *restmodels.DssUserInfo
	err = json.Unmarshal(result.Data, &dssUserInfo)
	if err != nil {
		log.Println("dss  get user info request failed, ", err)
		return nil, err
	}
	return dssUserInfo, nil
}

type ProjectData struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Business             string   `json:"business"`
	ApplicationArea      int64    `json:"applicationArea"`
	Product              string   `json:"product"`
	CreateBy             string   `json:"createBy"`
	Source               string   `json:"source"`
	ReleaseUsers         []string `json:"releaseUsers"`
	EditUsers            []string `json:"editUsers"`
	AccessUsers          []string `json:"accessUsers"`
	DevProcessList       []string `json:"devProcessList"`
	OrchestratorModeList []string `json:"orchestratorModeList"`
	WorkspaceId          int64    `json:"workspaceId"`
	Id                   int64    `json:"Id"`
}

type ProjectParam struct {
	ProjectName    string
	ProjectDesc    string
	ProjectProduct string
}

type CreateProjectRequest struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Business             string   `json:"business"`
	ApplicationArea      int64    `json:"applicationArea"`
	Product              string   `json:"product"`
	EditUsers            []string `json:"editUsers"`
	AccessUsers          []string `json:"accessUsers"`
	DevProcessList       []string `json:"devProcessList"`
	OrchestratorModeList []string `json:"orchestratorModeList"`
	WorkspaceId          int64    `json:"workspaceId"`
}

func (c *DSSClient) CreateProject(projectParam *ProjectParam, username string, devProcessList []string,
	workspaceId int64, workspaceName string) (*ProjectData, error) {
	requestPath := "/api/rest_j/v1/dss/framework/project/createProject"
	url := c.Address + requestPath
	defaultOrchestratorModeList := []string{"pom_work_flow"}
	defaultDevProcessList := []string{"dev", "prod"}
	if devProcessList == nil {
		devProcessList = defaultDevProcessList
	}
	createProjectRequest := &CreateProjectRequest{
		Name:                 projectParam.ProjectName,
		Description:          projectParam.ProjectDesc,
		Product:              projectParam.ProjectProduct,
		WorkspaceId:          workspaceId,
		ApplicationArea:      4,
		Business:             "",
		OrchestratorModeList: defaultOrchestratorModeList,
		DevProcessList:       defaultDevProcessList,
		EditUsers:            []string{},
		AccessUsers:          []string{},
	}
	workspaceNameCookies := http.Cookie{
		Name: "workspaceName",
		Value: workspaceName,
	}
	workspaceIdCookies := http.Cookie{
		Name: "workspaceId",
		Value: strconv.FormatInt(workspaceId,10),
	}
	requestJsonBytes, err := json.Marshal(&createProjectRequest)
	if err != nil {
		log.Println("dss create  projects request failed, ", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJsonBytes))
	req.AddCookie(&workspaceNameCookies)
	req.AddCookie(&workspaceIdCookies)
	if err != nil {
		log.Println("dss create projects request failed, ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)


	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return nil, err
	}

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

	return &projectData, nil
}

type GetProjectParams struct {
	WorkspaceId int64 `json:"workspaceId"`
}

func (c *DSSClient) GetProjectNameById(projectId int64, username string, WorkspaceId int64) (string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	url := c.Address + requestPath
	getProjectParams := &GetProjectParams{
		WorkspaceId: WorkspaceId,
	}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", url, strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return "", err
	}

	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss  get projects request failed, ", err)
		return "", err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil {
		log.Println("dss  get cooperate projects request failed, ", err)
		return "", err
	}
	projectList := []map[string]interface{}{}
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil {
		log.Println("dss  get cooperate projects request failed, ", err)
		return "", err
	}
	var projectName string
	for _, project := range projectList {
		if project["id"].(float64) == float64(projectId) {
			projectName = project["name"].(string)
			break
		}
	}
	return projectName, nil
}

func (c *DSSClient) GetProjectIdByName(projectName string, username string, WorkspaceId int64) (int64, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	url := c.Address + requestPath
	getProjectParams := &GetProjectParams{
		WorkspaceId: WorkspaceId,
	}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return 0, err
	}

	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss  get projects request failed, ", err)
		return -1, err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil {
		log.Println("dss  get cooperate projects request failed, ", err)
		return -1, err
	}
	projectList := []map[string]interface{}{}
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil {
		log.Println("dss  get cooperate projects request failed, ", err)
		return -1, err
	}
	var projectID float64
	for _, project := range projectList {
		if project["name"].(string) == projectName {
			projectID = project["id"].(float64)
			break
		}
	}
	return int64(projectID), nil
}

type project struct {
	Name                 string   `json:"name"`
	Description          string   `json:"description"`
	Business             string   `json:"business"`
	ApplicationArea      int64    `json:"applicationArea"`
	Product              string   `json:"product"`
	EditUsers            []string `json:"editUsers"`
	AccessUsers          []string `json:"accessUsers"`
	DevProcessList       []string `json:"devProcessList"`
	OrchestratorModeList []string `json:"orchestratorModeList"`
	WorkspaceId          int64    `json:"workspaceId"`
}

func (c *DSSClient) CheckProjectExist(workspaceId int64, projectName string, username string) (bool, *ProjectData, error) {
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

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects request failed, ", err)
		return isExisted, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return isExisted, nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss  get projects request failed, ", err)
		return isExisted, nil, err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil {
		log.Println("dss  get  projects request failed, ", err)
		return false, nil, err
	}

	projectList := []ProjectData{}
	var project ProjectData
	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil {
		log.Println("dss  get cooperate projects request failed, ", err)
		return isExisted, nil, err
	}
	for _, p := range projectList {
		if p.Name == projectName {
			isExisted = true
			project = p
			break
		}
	}
	return isExisted, &project, nil

}

type openOrchestratorRequest struct {
	OrchestratorId int64    `json:"orchestratorId"`
	WorkspaceName  string   `json:"workspaceName"`
	Labels         DSSLabel `json:"labels"`
}

type openOrchestratorResponse struct {
	OrchestratorOpenUrl string         `json:"OrchestratorOpenUrl"`
	OrchestratorVo      OrchestratorVo `json:"OrchestratorVo"`
}
type OrchestratorVo struct {
	DssOrchestratorInfo    dssOrchestratorInfo    `json:"dssOrchestratorInfo"`
	DssOrchestratorVersion dssOrchestratorVersion `json:"dssOrchestratorVersion"`
}

type dssOrchestratorVersion struct {
	Id             int64  `json:"id"`
	OrchestratorId int64  `json:"orchestratorId"`
	ProjectId      int64  `json:"projectId"`
	AppId          int64  `json:"appId"`
	Source         string `json:"source"`
	Version        string `json:"Version"`
	Comment        string `json:"Version"`
	Creator        string `json:"comment"`
	Content        string `json:"content"`
	ContextId      string `json:"contextId"`
}

type dssOrchestratorInfo struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	CreateTime    int64  `json:"createTime"`
	Creator       string `json:"creator"`
	Type          string `json:"type"`
	Uses          string `json:"uses"`
	AppConnName   string `json:"appConnName"`
	ProjectId     int64  `json:"projectId"`
	UUID          string `json:"uuid"`
	SecondaryType string `json:"secondaryType"`
	Description   string `json:"description"`
}

func (c *DSSClient) OpenOrchestrator(orchestratorId int64, labels DSSLabel, workspaceName string, user string) (*OrchestratorVo, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/orchestrator/openOrchestrator")
	url := c.Address + requestPath

	or := openOrchestratorRequest{
		OrchestratorId: orchestratorId,
		Labels:         labels,
		WorkspaceName:  workspaceName,
	}

	bytes, err := json.Marshal(or)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(user, req)
	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, user); err != nil {
		return nil, err
	}

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

func (c *DSSClient) GetAllProjectByUser(username string, workspaceId int64) ([]map[string]interface{}, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/project/getAllProjects")
	url := c.Address + requestPath
	getProjectParams := &GetProjectParams{
		WorkspaceId: workspaceId,
	}
	var projectList []map[string]interface{}
	bytes, err := json.Marshal(getProjectParams)
	if err != nil {
		return projectList, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get projects by user request failed, ", err)
		return projectList, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss  get projects by user request failed, ", err)
		return projectList, err
	}
	projectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &projectsObj)
	if err != nil {
		log.Println("dss  get cooperate projects by user request failed, ", err)
		return projectList, err
	}

	err = json.Unmarshal(projectsObj["projects"], &projectList)
	if err != nil {
		log.Println("dss  get cooperate projects by user request failed, ", err)
		return projectList, err
	}
	return projectList, nil
}

type getAllOrchestratorParams struct {
	WorkspaceId      string `json:"workspaceId"`
	ProjectId        int64  `json:"projectId"`
	OrchestratorMode string `json:"orchestratorMode"`
}

func (c *DSSClient) GetAllOrchestrator(username string, workspaceId int64, projectId int64, orchestratorMode string) ([]map[string]interface{}, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/framework/orchestrator/getAllOrchestrator")
	url := c.Address + requestPath
	logger.GetLogger().Printf("GetAllOrchestrator url ============= %v", url)
	getAllOrchestratorParams := &getAllOrchestratorParams{
		WorkspaceId:      strconv.FormatInt(workspaceId, 10),
		ProjectId:        projectId,
		OrchestratorMode: orchestratorMode,
	}
	var workFlowList []map[string]interface{}
	bytes, err := json.Marshal(getAllOrchestratorParams)
	if err != nil {
		return workFlowList, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(bytes)))
	if err != nil {
		log.Println("dss get workFlow list by user request failed, ", err)
		return workFlowList, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)

	//设置dss1.1.4 cookie header
	if err := c.setCookieHeader(req, username); err != nil {
		return nil, err
	}

	result, err := c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss  get workFlow list by user request failed, ", err)
		return workFlowList, err
	}
	FlowObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &FlowObj)
	if err != nil {
		log.Println("dss  get cooperate workFlow list by user request failed, ", err)
		return workFlowList, err
	}

	err = json.Unmarshal(FlowObj["page"], &workFlowList)
	if err != nil {
		log.Println("dss  get cooperate workFlow list by user request failed, ", err)
		return workFlowList, err
	}
	return workFlowList, nil
}
