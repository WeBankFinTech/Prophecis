package rest_impl

import (
	"encoding/base64"
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	loggerv2 "webank/DI/commons/logger/v2"
	"webank/DI/pkg/client/dss"
	models "webank/DI/pkg/v2/model"
	"webank/DI/pkg/v2/repo"
	"webank/DI/restapi/api_v2/restmodels"
	"webank/DI/restapi/api_v2/server/operations/experiment"
	"webank/DI/restapi/service_v2"
)

var experimentService = service_v2.ExperimentService

func GetDSSJumpMLSSUser(params experiment.GetDSSJumpMLSSUserParams) middleware.Responder {
	operation := "GetDSSJumpMLSSUser"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	dssClient := dss.GetDSSClient()
	cookies := []*http.Cookie{}
	ticket, err := url.QueryUnescape(params.DssUserTicketID)
	if err != nil {
		log.WithError(err).Errorf(operation + "QueryUnescape failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	ticketIdCookies := http.Cookie{
		Name:  "linkis_user_session_ticket_id_v1",
		Value: ticket,
	}
	workspaceIdCookies := http.Cookie{
		Name:  "workspaceId",
		Value: params.DssWorkspaceID,
	}
	cookies = append(cookies, &ticketIdCookies, &workspaceIdCookies)
	log.Infof("DssGetUserInfo cookies: %+v", cookies)
	res, err := dssClient.GetUserInfo(cookies)
	if err != nil {
		log.WithError(err).Errorf(operation + "get dss user info failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, &restmodels.GetDSSJumpMLSSUserResponse{
		UserName: res.UserName,
	})
}

func ListExperimentVersionNames(params experiment.ListExperimentVersionNamesParams) middleware.Responder {
	operation := "ListExperimentVersionNames"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)

	versions, err := experimentService.ListExperimentVersionNames(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.CanDeploy)
	if err != nil {
		log.WithError(err).Errorf(operation + "ListExperimentVersionNames failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, versions)
}

func ListExperimentSourceSystems(params experiment.ListExperimentSourceSystemsParams) middleware.Responder {
	//operation := "ListExperimentSourceSystems"
	ctx := GenerateContext(params.HTTPRequest)
	returnResult := []string{"MLSS", "DSS", "WTSS"}
	return HttpResponseHandle2(ctx, nil, returnResult)
}

func ListExperimentDssProjectNames(params experiment.ListExperimentDSSProjectNamesParams) middleware.Responder {
	operation := "ListExperimentDssProjectNames"
	ctx := GenerateContext(params.HTTPRequest)
	returnResult, err := repo.ExperimentRepo.GetDSSFieldValues("dss_project_name")
	if err != nil {
		log.WithError(err).Errorf(operation + "ExperimentRepo GetFieldValues")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, returnResult)
}

func ListExperimentDssFlowNames(params experiment.ListExperimentDSSFlowNamesParams) middleware.Responder {
	operation := "ListExperimentDssProjectNames"
	ctx := GenerateContext(params.HTTPRequest)
	returnResult, err := repo.ExperimentRepo.GetDSSFieldValues("dss_flow_name")
	if err != nil {
		log.WithError(err).Errorf(operation + "ExperimentRepo GetFieldValues")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, returnResult)
}

func CreateExperiment(params experiment.CreateExperimentParams) middleware.Responder {
	operation := "CreateExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)

	if err := validateCreateExperiment(params); err != nil {
		return HttpResponseHandle2(ctx, err, nil)
	}
	if ok, reason := validateFlowJson(params.Experiment.FlowJSON); !ok {
		log.Errorf("输入flowjson数据有问题")
		return HttpResponseHandle2(ctx, service_v2.InputParamsError{Field: "FlowJson", Reason: reason}, nil)
	}
	experiment, err := experimentService.CreateExperiment(currentUserName, strings.ToUpper(clusterType), params.Experiment)
	if err != nil {
		log.Errorf("experimentService.CreateExperiment failed: %+v", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	returnResult := &restmodels.CreateExperimentResponse{
		ID: &experiment.ID,
	}
	return HttpResponseHandle2(ctx, nil, returnResult)
}

func CreateExperimentByUpload(params experiment.CreateExperimentByUploadParams) middleware.Responder {
	operation := "CreateExperimentByUpload"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)

	// 1. 将zip文件存放到 ./importExperimentFlowJson/expName-2006-01-02-15-04-05.zip zipFilePath
	// 2. 解压该文件到 ./importExperimentFlowJson/expName-2006-01-02-15-04-05/ unArchivePath
	// 3. 去读第2步的目录的flowjson文件 ./importExperimentFlowJson/expName-2006-01-02-15-04-05/flowjson/flowjson.json flowjsonPath
	fileBytes, err := ioutil.ReadAll(params.File)
	timeNowStr := time.Now().Format("2006-01-02-15-04-05")
	expNameTimeNowFileName := fmt.Sprintf("%s-%s.zip", params.Name, timeNowStr)
	zipFilePath := filepath.Join("./importExperimentFlowJson", expNameTimeNowFileName)
	err = os.MkdirAll("./importExperimentFlowJson", os.ModePerm)
	if err != nil {
		log.Errorf("create folderPath failed: %+v", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	zipFilePathFile, err := os.Create(zipFilePath)
	if err != nil {
		log.Println("create flow.zip failed, err:", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	_, err = zipFilePathFile.Write(fileBytes)
	if err != nil {
		log.Println("flow.zip write failed, ", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	//un zip file
	unArchivePath := filepath.Join("./importExperimentFlowJson", fmt.Sprintf("%s-%s", params.Name, timeNowStr))
	err = archiver.Unarchive(zipFilePath, unArchivePath)
	if err != nil {
		log.Println("un archive file failed, ", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	flowjsonPath := filepath.Join(unArchivePath, "flowjson", "flowjson.json")
	jsonBytes, err := ioutil.ReadFile(flowjsonPath)
	if err != nil || len(jsonBytes) == 0 {
		return HttpResponseHandle2(ctx, err, nil)
	}
	createExperimentRequest := &restmodels.CreateExperimentRequest{
		ClusterType: params.ClusterType,
		//Description: *params.Description,
		//DssInfo: &restmodels.DSSInfo{
		//	FlowID:        params.DssFlowID,
		//	FlowName:      params.DssFlowName,
		//	FlowVersion:   params.DssFlowVersion,
		//	WorkspaceID:   params.DssWorkspaceID,
		//	WorkspaceName: params.DssWorkspaceName,
		//	ProjectID:     params.DssProjectID,
		//	ProjectName:   params.DssProjectName,
		//},
		GroupName:    params.GroupName,
		FlowJSON:     string(jsonBytes),
		Name:         &params.Name,
		SourceSystem: params.OriginSystem,
		Tags:         params.Tags,
	}
	exp, err := experimentService.CreateExperiment(currentUserName, strings.ToUpper(clusterType), createExperimentRequest)
	if err != nil {
		log.WithError(err).Errorf("experimentService.CreateExperiment failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	returnResult := modelExperimentToRestModel(exp)
	return HttpResponseHandle2(ctx, nil, returnResult)
}

func GetExperiment(params experiment.GetExperimentParams) middleware.Responder {
	operation := "GetExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	exp, err := experimentService.GetExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID, "")
	if err != nil {
		log.WithError(err).Errorf("experimentService.GetExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	returnResult := modelExperimentToRestModel(exp)
	return HttpResponseHandle2(ctx, nil, returnResult)
	//return ProcessResponse(operation, err, exp, &restmodels.ExperimentWithoutPipeline{})
}

func ListExperiments(params experiment.ListExperimentsParams) middleware.Responder {
	operation := "ListExperiments"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	exps, total, err := experimentService.ListExperiments(ctx, currentUserName, strings.ToUpper(clusterType), params.Name, params.Tag, params.GroupID, params.GroupName,
		params.SourceSystem, params.CreateUser, params.CreateTimeSt, params.CreateTimeEd,
		params.UpdateTimeSt, params.UpdateTimeEd, params.DssProjectName, params.DssFlowName, params.Size, params.Page)
	if err != nil {
		log.WithError(err).Errorf("experimentService.ListExperiments failed")
		return HttpResponseHandle2(ctx, err, nil)
	}

	var returnExps []*restmodels.ExperimentWithoutPipeline
	for _, exp := range exps {
		returnExp := modelExperimentToRestModel(exp)
		returnExps = append(returnExps, returnExp)
	}

	response := restmodels.ListExperimentsResponse{
		Experiments: returnExps,
		TotalPage:   int64(math.Ceil(float64(total) / float64(*params.Size))),
		Total:       total,
		PageNumber:  *params.Page,
		PageSize:    *params.Size,
	}
	return HttpResponseHandle2(ctx, err, response)
}

func DeleteExperiment(params experiment.DeleteExperimentParams) middleware.Responder {
	operation := "DeleteExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentService.DeleteExperiment(currentUserName, strings.ToUpper(clusterType), params.ExpID)
	if err != nil {
		log.WithError(err).Errorf("experimentService.DeleteExperiment failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func BatchDeleteExperiments(params experiment.BatchDeleteExperimentsParams) middleware.Responder {
	operation := "BatchDeleteExperiments"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	for _, expId := range params.BatchDeleteExperimentsRequest.Ids {
		err := experimentService.DeleteExperiment(currentUserName, strings.ToUpper(clusterType), expId)
		if err != nil {
			log.WithError(err).Errorf("experimentService.DeleteExperiment failed")
			return HttpResponseHandle2(ctx, err, nil)
		}
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func DeleteExperimentVersion(params experiment.DeleteExperimentVersionParams) middleware.Responder {
	operation := "DeleteExperimentVersion"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentService.DeleteExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.VersionName)
	if err != nil {
		log.WithError(err).Errorf("experimentService.DeleteExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func PatchExperiment(params experiment.PatchExperimentParams) middleware.Responder {
	operation := "PatchExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	if ok, reason := validateFlowJson(params.Experiment.FlowJSON); !ok {
		log.Errorf("输入flowjson数据有问题")
		return HttpResponseHandle2(ctx, service_v2.InputParamsError{Field: "FlowJson", Reason: reason}, nil)
	}
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentService.PatchExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID, "", params.Experiment)
	if err != nil {
		log.WithError(err).Errorf("experimentService.PatchExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func PatchExperimentVersion(params experiment.PatchExperimentVersionParams) middleware.Responder {
	operation := "PatchExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	request := &restmodels.PatchExperimentRequest{
		Description: params.ExperimentVersion.Description,
	}
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentService.PatchExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.VersionName, request)
	if err != nil {
		log.WithError(err).Errorf("experimentService.PatchExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func CreateExperimentVersion(params experiment.CreateExperimentVersionParams) middleware.Responder {
	operation := "CreateExperimentVersion"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	// CreateExperimentVersion 的作用是固化当前版本，所以如果有描述，是将当前版本的描述覆盖，新版本的描述应该为空
	var newExpDescription string
	if params.Experiment != nil {
		newExpDescription = params.Experiment.Description
	}
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	archivedVersionName, latestVersionName, err := experimentService.CreateExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID,
		newExpDescription)
	if err != nil {
		log.WithError(err).Errorf("experimentService.CreateExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := restmodels.CreateExperimentVersionResponse{
		archivedVersionName, latestVersionName,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func GetExperimentVersionFlowJson(params experiment.GetExperimentVersionFlowJSONParams) middleware.Responder {
	operation := "GetExperimentVersionFlowJson"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	flowJson, _, _, err := experimentService.GetExperimentVersionFlowJson(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.VersionName)
	if err != nil {
		// 如果是 flowJson is empty 的错误，则忽略该错误
		if !strings.Contains(err.Error(), "flowJson is empty") {
			log.WithError(err).Errorf("experimentService.GetExperimentVersionFlowJson failed")
			return HttpResponseHandle2(ctx, err, nil)
		}
	}
	response := restmodels.FlowJSONResponse{
		flowJson,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func GetExperimentVersionGlobalVariablesStr(params experiment.GetExperimentVersionGlobalVariablesStrParams) middleware.Responder {
	operation := "GetExperimentVersionGlobalVariables"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	globalVariablesStr, _ := experimentService.GetExperimentVersionGlobalVariablesStr(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.VersionName)
	response := restmodels.GetExperimentGlobalVariablesStrResponse{
		globalVariablesStr,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func GetExperimentVersion(params experiment.GetExperimentVersionParams) middleware.Responder {
	operation := "GetExperimentVersion"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	exp, err := experimentService.GetExperimentVersion(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.VersionName)
	if err != nil {
		log.WithError(err).Errorf("experimentService.GetExperimentVersion failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := modelExperimentToRestModel(exp)
	return HttpResponseHandle2(ctx, nil, response)
}

func ListExperimentVersions(params experiment.ListExperimentVersionsParams) middleware.Responder {
	operation := "ListExperimentVersions"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	exps, total, err := experimentService.ListExperimentVersions(currentUserName, strings.ToUpper(clusterType), &params.ExpID, params.VersionName,
		params.UpdateTimeSt, params.UpdateTimeEd, params.Size, params.Page)
	if err != nil {
		log.WithError(err).Errorf("experimentService.ListExperimentVersions failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	var returnExps []*restmodels.ExperimentWithoutPipeline
	for _, exp := range exps {
		returnExp := modelExperimentToRestModel(exp)
		returnExps = append(returnExps, returnExp)
	}
	response := restmodels.ListExperimentsResponse{
		Experiments: returnExps,
		TotalPage:   int64(math.Ceil(float64(total) / float64(*params.Size))),
		Total:       total,
		PageNumber:  *params.Page,
		PageSize:    *params.Size,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func CopyExperiment(params experiment.CopyExperimentParams) middleware.Responder {
	operation := "CopyExperiment"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	exp, err := experimentService.CopyExperiment(currentUserName, strings.ToUpper(clusterType), params.ExpID, params.Experiment)
	if err != nil {
		log.WithError(err).Errorf("experimentService.CopyExperiment failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := modelExperimentToRestModel(exp)
	return HttpResponseHandle2(ctx, nil, response)
}

func UploadExperimentFlowJson(params experiment.UploadExperimentFlowJSONParams) middleware.Responder {
	//operation := "UploadExperimentFlowJson"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	// 1. 将zip文件存放到 ./importExperimentFlowJson/expName-2006-01-02-15-04-05.zip zipFilePath
	// 2. 解压该文件到 ./importExperimentFlowJson/expName-2006-01-02-15-04-05/ unArchivePath
	// 3. 去读第2步的目录的flowjson文件 ./importExperimentFlowJson/expName-2006-01-02-15-04-05/flowjson/flowjson.json flowjsonPath
	fileBytes, err := ioutil.ReadAll(params.File)
	timeNowStr := time.Now().Format("2006-01-02-15-04-05")
	expNameTimeNowFileName := fmt.Sprintf("%s-%s.zip", "upload_flow_json", timeNowStr)
	zipFilePath := filepath.Join("./importExperimentFlowJson", expNameTimeNowFileName)
	log.Infof("the zipFilePath is: %s", zipFilePath)
	err = os.MkdirAll("./importExperimentFlowJson", os.ModePerm)
	if err != nil {
		log.Errorf("create folderPath failed: %+v", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	zipFilePathFile, err := os.Create(zipFilePath)
	if err != nil {
		log.Println("create flow.zip failed, err:", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	_, err = zipFilePathFile.Write(fileBytes)
	if err != nil {
		log.Println("flow.zip write failed, ", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	//un zip file
	unArchivePath := filepath.Join("./importExperimentFlowJson", fmt.Sprintf("%s-%s", "upload_flow_json", timeNowStr))
	err = archiver.Unarchive(zipFilePath, unArchivePath)
	if err != nil {
		log.Println("un archive file failed, ", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	log.Infof("the zip unArchivePath is: %s", unArchivePath)
	flowjsonPath := filepath.Join(unArchivePath, "flowjson", "flowjson.json")
	log.Infof("the flowjsonPath is: %s", flowjsonPath)
	jsonBytes, err := ioutil.ReadFile(flowjsonPath)
	if err != nil || len(jsonBytes) == 0 {
		return HttpResponseHandle2(ctx, err, nil)
	}
	flowjsonPathBase64 := base64.StdEncoding.EncodeToString([]byte(flowjsonPath))
	log.Infof("the flowjsonPathBase64 is: %s", flowjsonPathBase64)
	response := restmodels.UploadExperimentFlowJSONResponse{
		UploadID: flowjsonPathBase64,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func CodeUpload(params experiment.CodeUploadParams) middleware.Responder {
	//operation := "CodeUpload"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	s3Path, err := experimentService.UploadCodeBucket(params.File, "di-model")
	log.Infof("the s3Path is: %s", s3Path)
	if err != nil {
		log.Errorf("UploadCodeBucket faild: %+v ", err.Error())
		return HttpResponseHandle2(ctx, err, nil)
	}
	res := restmodels.CodeUploadResponse{S3path: s3Path}
	return HttpResponseHandle2(ctx, nil, res)
}

func ExportExperimentFlowJson(params experiment.ExportExperimentFlowJSONParams) middleware.Responder {
	operation := "ExportExperimentFlowJson"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	zipFile, err := experimentService.ExportExperimentFlowJson(currentUserName, strings.ToUpper(clusterType), params.ExpID)
	if err != nil {
		log.Errorf("experimentService.ExportExperimentFlowJson faild: %+v ", err.Error())
		return HttpResponseHandle2(ctx, err, nil)
	}

	fileBytes, err := ioutil.ReadAll(zipFile)
	defer zipFile.Close()
	if err != nil || len(fileBytes) == 0 {
		log.Errorf("read zipFile faild or zipFile is empty: %+v ", err.Error())
		return HttpResponseHandle2(ctx, err, nil)
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		fileName := "flowjson.zip"
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(fileBytes)
		if err != nil {
			// Handle error during file streaming
			log.Printf("Error streaming file to response: %v", err)
			// Note: Since streaming has started, setting the status code here might not have the intended effect.
		}
	})
}

func modelExperimentToRestModel(exp *models.Experiment) *restmodels.ExperimentWithoutPipeline {
	var returnExp restmodels.ExperimentWithoutPipeline
	log.Infof("exp.LatestExecuteTime is %+v", exp.LatestExecuteTime)
	err := copier.Copy(&returnExp, exp)
	if err != nil {
		log.Errorf("copier.Copy failed: %+v", err)
		return nil
	}
	if len(exp.Tags) != 0 {
		returnExp.Tags = strings.Split(exp.Tags, ",")
	}
	if exp.SourceSystem == restmodels.ExperimentWithoutPipelineSourceSystemDSS {
		returnExp.DssInfo = &restmodels.DSSInfo{
			exp.DSSFlowID,
			exp.DSSFlowName,
			exp.DSSFlowVersion,
			exp.DSSProjectID,
			exp.DSSProjectName,
			exp.DSSWorkspaceID,
			exp.DSSWorkspaceName,
		}
	}
	returnExp.DeploySetting = &restmodels.DeploySetting{}
	returnExp.DeploySetting.CanDeployAsOfflineService = &exp.CanDeployAsOfflineService
	log.Infof("before the returnExp.LatestExecuteTime is: %+v", returnExp.LatestExecuteTime)
	if exp.LatestExecuteTime.IsZero() {
		returnExp.LatestExecuteTime = ""
	} else {
		returnExp.LatestExecuteTime = exp.LatestExecuteTime.Format(time.RFC3339)
	}
	log.Infof("after the returnExp.LatestExecuteTime is: %+v", returnExp.LatestExecuteTime)
	return &returnExp
}
