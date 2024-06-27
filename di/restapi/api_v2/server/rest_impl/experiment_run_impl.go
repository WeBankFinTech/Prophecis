package rest_impl

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
	loggerv2 "webank/DI/commons/logger/v2"
	models "webank/DI/pkg/v2/model"
	"webank/DI/restapi/api_v2/restmodels"
	"webank/DI/restapi/api_v2/server/operations/experiment_run"
	"webank/DI/restapi/service_v2"
)

var experimentRunService = service_v2.ExperimentRunService

func CreateExperimentRun(params experiment_run.CreateExperimentRunParams) middleware.Responder {
	operation := "CreateExperimentRun"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	experimentRun, err := experimentRunService.CreateExperimentRun(ctx, currentUserName, strings.ToUpper(clusterType), params.ExperimentRunRequest)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.CreateExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	returnExpRun := restmodels.CreateExperimentRunResponse{}
	err = copier.Copy(&returnExpRun, experimentRun)
	if err != nil {
		log.Errorf("copier.Copy failed: %+v", err)
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, returnExpRun)
}

func DeleteExperimentRun(params experiment_run.DeleteExperimentRunParams) middleware.Responder {
	operation := "DeleteExperimentRun"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentRunService.DeleteExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.DeleteExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func GetExperimentRun(params experiment_run.GetExperimentRunParams) middleware.Responder {
	operation := "GetExperimentRun"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	expRun, err := experimentRunService.GetExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.GetExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := modelExpRunToRestModel(expRun)
	return HttpResponseHandle2(ctx, nil, response)
}

// todo(gaoyuanhe): 待实现
func GetExperimentRunNodeLogs(params experiment_run.GetExperimentRunNodeLogsParams) middleware.Responder {
	operation := "GetExperimentRunNodeLogs"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	expRun, err := experimentRunService.GetExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.GetExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := modelExpRunToRestModel(expRun)
	return HttpResponseHandle2(ctx, nil, response)
}

// todo(gaoyuanhe): 待实现
func GetExperimentRunNodeExecutions(params experiment_run.GetExperimentRunNodeExecutionsParams) middleware.Responder {
	operation := "GetExperimentRunNodeExecutions"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	expRunNodes, err := experimentRunService.ListExperimentRunNodes(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.GetExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := modelExpRunNodesToRestModel(expRunNodes)
	return HttpResponseHandle2(ctx, nil, response)
}

func modelExpRunNodesToRestModel(expRunNodes []*models.ExperimentRunNode) *restmodels.NodeExecutionsResponse {
	var response restmodels.NodeExecutionsResponse
	for _, expRunNode := range expRunNodes {
		if expRunNode.State == "RUNNIng" {
			response.RunningNodes = append(response.RunningNodes, &restmodels.NodeExecutionInfo{
				ExecRunID:      expRunNode.ExperimentRunID,
				NodeID:         expRunNode.ID,
				NodeStatus:     expRunNode.State,
				TrainingTaskID: expRunNode.ProxyID,
			})
		}
	}
	return &response
}

func modelExpRunToRestModel(expRun *models.ExperimentRun) *restmodels.ExperimentRun {
	var returnExpRun restmodels.ExperimentRun
	err := copier.Copy(&returnExpRun, expRun)
	if err != nil {
		log.Errorf("copier.Copy failed: %+v", err)
		return nil
	}
	return &returnExpRun
}

func GetExperimentRunStatus(params experiment_run.GetExperimentRunStatusParams) middleware.Responder {
	operation := "GetExperimentRunStatus"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	expRun, err := experimentRunService.GetExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.GetExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := restmodels.ExperimentRunStatus{
		Status: expRun.ExecuteStatus,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func GetExperimentRunFlowJson(params experiment_run.GetExperimentRunFlowJSONParams) middleware.Responder {
	operation := "GetExperimentRunFlowJson"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	expRun, err := experimentRunService.GetExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.GetExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := restmodels.FlowJSONResponse{
		FlowJSON: expRun.FlowJson,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func UpdateExperimentRunExecuteStatus() {
	go experimentRunService.UpdateExperimentRunExecuteStatus()
}

func ListExperimentRuns(params experiment_run.ListExperimentRunsParams) middleware.Responder {
	operation := "ListExperimentRuns"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	var requestStatus []string
	if params.Status != nil {
		requestStatus = append(requestStatus, *params.Status)
	}
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)

	expRuns, total, err := experimentRunService.ListExperimentRuns(false, currentUserName, strings.ToUpper(clusterType), &params.ExpID, params.ExpVersionName, requestStatus,
		params.CreateTimeSt, params.CreateTimeEd, params.Size, params.Page)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.ListExperimentRuns failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	var returnExps []*restmodels.ExperimentRun
	err = copier.Copy(&returnExps, &expRuns)
	if err != nil {
		log.Error("Copy Experiments Error: " + err.Error())
		return HttpResponseHandle2(ctx, err, nil)
	}
	response := restmodels.ListExperimentRunsResponse{
		ExperimentRuns: returnExps,
		TotalPage:      int64(math.Ceil(float64(total) / float64(*params.Size))),
		Total:          total,
		PageNumber:     *params.Page,
		PageSize:       *params.Size,
	}
	return HttpResponseHandle2(ctx, nil, response)
}

func KillExperimentRun(params experiment_run.KillExperimentRunParams) middleware.Responder {
	operation := "KillExperimentRun"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentRunService.KillExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.KillExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}

func RetryExperimentRun(params experiment_run.RetryExperimentRunParams) middleware.Responder {
	operation := "RetryExperimentRun"
	ctx := GenerateContext(params.HTTPRequest)
	log := loggerv2.GetLogger(ctx)
	currentUserName := GetUserID(params.HTTPRequest)
	clusterType := GetClusterType(params.HTTPRequest)
	log.Infof("%s(%s, %s)", operation, currentUserName, clusterType)
	err := experimentRunService.RetryExperimentRun(currentUserName, strings.ToUpper(clusterType), params.RunID)
	if err != nil {
		log.WithError(err).Errorf("experimentRunService.KillExperimentRun failed")
		return HttpResponseHandle2(ctx, err, nil)
	}
	return HttpResponseHandle2(ctx, nil, nil)
}
