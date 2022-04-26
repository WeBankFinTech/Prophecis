package dss

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
	"time"
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

	httpClient := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}

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
	data WorkspacesData `json:"data"`
}

type WorkflowsData struct {
	PersonalFlows    []PersonalFlow `json:"personalFlows"`
	ProjectVersionId int64          `json:"projectVersionId"`
}


type Workspace struct {
	CreateBy   string `json:"createBy"`
	Department string `json:"department"`
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Product    string `json:"product"`
}

type PersonalFlow struct {
	CreateBy   string `json:"createBy"`
	Department string `json:"department"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ProjectID  int    `json:"projectID"`
}

var log = logger.GetLogger()

func (c *DSSClient) GetWorkspaces(user string) ([]Workspace, error) {
	requestPath := "/api/rest_j/v1/dss/workspaces"
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

func (c *DSSClient) ListPersonalWorkflows(workspaceId int64, user string) (*WorkflowsData, error) {
	//requestPath := "/api/rest_j/v1/dss/listPersonalWorkflows"
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/listPersonalWorkflows?workspaceId=%v", workspaceId)

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
	var res WorkflowsData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type AddFlowRequest struct {
	Description      string `json:"description"`
	Name             string `json:"name"`
	ProjectVersionID int64  `json:"projectVersionID"`
	Uses             string `json:"uses"`
	TaxonomyID      int64  `json:"taxonomyID"`
}

func (c *DSSClient) AddFlow(params AddFlowRequest, user string) (*CreateWorkflowData, error) {
	requestPath := "/api/rest_j/v1/dss/addFlow"

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
	var res CreateWorkflowData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type SaveWorkflowParams struct {
	BmlVersion       string `json:"bmlVersion"`
	Comment          string `json:"comment"`
	FlowEditLock     string `json:"flowEditLock"`
	FlowVersion      string `json:"flowVersion"`
	Id               int64  `json:"id"`
	Json             string `json:"json"`
	ProjectID        int64  `json:"projectID"`
	ProjectVersionID int64  `json:"projectVersionID"`
	IsNotHaveLock bool `json:"isNotHaveLock"`
	//dssFlowId        int64  `json:"dss_flow_id"`
}

type SaveFlowResponse struct {
	Response
	data SaveFlowData `json:"flowEditLock"`
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
	requestPath := "/api/rest_j/v1/dss/saveFlow"
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

type FlowGetResponse struct {
	Response
	Data    FlowGetData `json:"data"`
}

type FlowGetData struct {
	Flow Flow `json:"flow"`
}

type Flow struct {
	CreateTime    int64          `json:"createTime"`
	Creator       string         `json:"creator"`
	Description   string         `json:"description"`
	FlowType      string         `json:"flowType"`
	FlowVersions  []*FlowVersion `json:"flowVersions"`
	LatestVersion *FlowVersion   `json:"latestVersion"`
	HasSaved      bool           `json:"hasSaved"`
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	ProjectID     int64          `json:"projectID"`
	Rank          int64          `json:"rank"`
	RootFlow      bool           `json:"rootFlow"`
	Source        string         `json:"source"`
	State         bool           `json:"state"`
	Uses          string         `json:"uses"`
}
type FlowVersion struct {
	BmlVersion       string `json:"bmlVersion"`
	FlowID           int    `json:"flowID"`
	Id               int    `json:"id"`
	Json             string `json:"json"`
	JsonPath         string `json:"jsonPath"`
	NotPublished     bool   `json:"notPublished"`
	Source           string `json:"source"`
	UpdateTime       int64  `json:"updateTime"`
	Updator          string `json:"updator"`
	Version          string  `json:"version"`
	ProjectVersionID int64    `json:"projectVersionID"`
}

func (c *DSSClient) GetFlow(id int64, version string, projectVersionID int64, user string) (*FlowGetData, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/get?id=%v&version=%v&projectVersionID=%v", id, version, projectVersionID)
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
	var res FlowGetData
	err = json.Unmarshal(result.Data, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
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

func (c * DSSClient) CheckCooperateProject(projectName, username string, workspaceId int64) (bool, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/listCooperateProjects?workspaceId=%v",workspaceId)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("dss get cooperate projects request failed, ", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return false, err
	}
	cooperateProjectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &cooperateProjectsObj)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return false, err
	}
	cooperateProjectList := []map[string]interface{}{}
	err = json.Unmarshal(cooperateProjectsObj["cooperateProjects"], &cooperateProjectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return false, err
	}
	checkFlag := false
	for _, cooperateProject := range cooperateProjectList {
		if cooperateProject["name"] == projectName{
			checkFlag = true
			break
		}
	}
	return checkFlag, nil
}

type CooperProjectData struct {
	ProjectName    string
	ProjectDesc    string
	ProjectProduct string
}

type CreateCooperProjectRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	TaxonomyID      int64    `json:"taxonomyID"`
	Business        string   `json:"business"`
	ApplicationArea int64    `json:"applicationArea"`
	Product         string   `json:"product"`
	EditUsers       []string `json:"editUsers"`
	AccessUsers     []string `json:"accessUsers"`
	WorkspaceId     int64    `json:"workspaceId"`
}

func (c * DSSClient) CreateCooperProject(cooperProjectData *CooperProjectData, username string, workspaceId int64) error {
	requestPath := "/api/rest_j/v1/dss/createCooperateProject"
	url := c.Address + requestPath
	createCooperProjectRequest := &CreateCooperProjectRequest{
		Name:            cooperProjectData.ProjectName,
		Description:     cooperProjectData.ProjectDesc,
		Product:         cooperProjectData.ProjectProduct,
		WorkspaceId:     workspaceId,
		ApplicationArea: 4,
		Business:        "",
		TaxonomyID:      1,
		EditUsers:       []string{},
		AccessUsers:     []string{},
	}
	requestJsonBytes, err := json.Marshal(&createCooperProjectRequest)
	if err != nil {
		log.Println("dss create cooperate projects request failed, ", err)
		return err
	}
	req, err := http.NewRequest("POST",url,bytes.NewBuffer(requestJsonBytes))
	if err != nil {
		log.Println("dss create cooperate projects request failed, ", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	_,err = c.DoLinkisHttpRequest(req, username)
	if err != nil {
		log.Println("dss create cooperate projects request failed, ", err)
		return err
	}
	return nil
}

func (c * DSSClient) GetProjectVersionId(projectName, username string, workspaceId int64) (int64, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/listCooperateProjects?workspaceId=%v",workspaceId)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("dss get cooperate projects request failed, ", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return 0, err
	}
	cooperateProjectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &cooperateProjectsObj)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return 0, err
	}
	cooperateProjectList := []map[string]interface{}{}
	err = json.Unmarshal(cooperateProjectsObj["cooperateProjects"], &cooperateProjectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return 0, err
	}
	projectVersionId := 0
	for _, cooperateProject := range cooperateProjectList {
		if cooperateProject["name"].(string) == projectName{
			latestVersionObj := cooperateProject["latestVersion"].(map[string]interface{})
			projectVersionId = int(latestVersionObj["id"].(float64))
			break
		}
	}
	return int64(projectVersionId), nil
}

func (c * DSSClient) GetProjectNameById(projectId int64, username string, workspaceId int64) (string, error) {
	requestPath := fmt.Sprintf("/api/rest_j/v1/dss/listCooperateProjects?workspaceId=%v",workspaceId)
	url := c.Address + requestPath
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		log.Println("dss get cooperate projects request failed, ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	c.setAuthHeader(username, req)
	result,err := c.DoLinkisHttpRequest(req, username)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return "", err
	}
	cooperateProjectsObj := map[string]json.RawMessage{}
	err = json.Unmarshal(result.Data, &cooperateProjectsObj)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return "", err
	}
	cooperateProjectList := []map[string]interface{}{}
	err = json.Unmarshal(cooperateProjectsObj["cooperateProjects"], &cooperateProjectList)
	if err != nil{
		log.Println("dss  get cooperate projects request failed, ",err)
		return "", err
	}
	var projectName string
	for _, cooperateProject := range cooperateProjectList {
		if cooperateProject["id"].(float64) == float64(projectId){
			projectName = cooperateProject["name"].(string)
			break
		}
	}
	return projectName, nil
}