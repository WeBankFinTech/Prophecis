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
package rest_impl

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"math"
	"net/http"
	"webank/DI/commons"
	"webank/DI/commons/logger"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/experiment_runs"

	//"webank/DI/restapi/service"

	service "webank/DI/restapi/newService"
)

var experimentRunService = service.ExperimentRunService

func CreateExperimentRun(params experiment_runs.CreateExperimentRunParams) middleware.Responder {
	logr := logger.GetLogger()
	operation := "CreateExperimentRun"
	var currentUserId = getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	run, err := service.ExperimentRunService.CreateExperimentRun(params.ExpID, &params.ExperimentRunRequest.FlowJSON,
		currentUserId, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		logr.WithError(err).Errorf("ExperimentService.CreateExperimentRun failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	newResponse := restmodels.ProphecisExperimentRun{}
	err = copier.Copy(&newResponse, run)
	if err != nil {
		logr.WithError(err).Errorf("copier.Copy failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	marshal, err := json.Marshal(newResponse)
	if err != nil {
		logr.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)

}

func GetExperimentRun(params experiment_runs.GetExperimentRunParams) middleware.Responder {
	operation := "GetExperimentRun"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	expRun, flowJson, err := experimentRunService.GetExperimentRun(params.ID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.Error(operation + " Error: " + err.Error())
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte(err.Error()))
	}

	expRunRes := restmodels.ProphecisExperimentRun{}
	err = copier.Copy(&expRunRes, expRun)
	if err != nil {
		log.Error(err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	response := restmodels.ProphecisExperimentRunsGetResponse{
		FlowJSON:               flowJson,
		ProphecisExperimentRun: &expRunRes,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func DeleteExperimentRun(params experiment_runs.DeleteExperimentRunParams) middleware.Responder {
	operation := "DeleteExperimentRun"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	err := experimentRunService.DeleteExperimentRun(params.ID, user, isSA)

	if err != nil {
		log.Error(operation + " Error: " + err.Error())
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, nil)
}

func KillExperimentRun(params experiment_runs.KillExperimentRunParams) middleware.Responder {
	operation := "KillExperimentRun"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	err := experimentRunService.KillExperimentRun(params.ID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}

	if err != nil {
		log.Error(operation + " Error: " + err.Error())
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, nil)
}

func GetExperimentRunLog(params experiment_runs.GetExperimentRunLogParams) middleware.Responder {
	operation := "GetExperimentRunLog"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	runLog, err :=
		experimentRunService.GetExperimentRunLog(params.ID, *params.FromLine, *params.Size, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.Error(operation + " Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	marshal, err := json.Marshal(runLog)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func StatusExperimentRun(params experiment_runs.GetExperimentRunStatusParams) middleware.Responder {
	operation := "StatusExperimentRun"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	status, err := experimentRunService.GetExperimentRunStatus(params.ID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}

	if err != nil {
		log.Error("")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	response := restmodels.ExperimentRunStatusResponse{RunID: params.ID, Status: *status}
	marshal, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func ListExperimentRuns(params experiment_runs.ListExperimentRunsParams) middleware.Responder {
	operation := "ListExperimentRuns"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	experimentRuns, count, err := experimentRunService.ListExperimentRun(*params.Page, *params.Size, queryStr, user, isSA)
	if err != nil {
		log.Error(err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	var expRuns []*restmodels.ProphecisExperimentRun
	err = copier.Copy(&expRuns, &experimentRuns)
	if err != nil {
		log.Error("Copy Experiment Runs Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	response := restmodels.ProphecisExperimentRuns{
		ExperimentRuns: expRuns,
		PageNumber:     *params.Page,
		PageSize:       *params.Size,
		Total:          count,
		TotalPage:      int64(math.Ceil(float64(count) / float64(*params.Size))),
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func GetExperimentRunExecution(params experiment_runs.GetExperimentRunExecutionParams) middleware.Responder {
	operation := "GetExperimentRunExecution"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	execData, err := experimentRunService.GetRunExecution(params.ExecID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.Error(err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	marshal, err := json.Marshal(execData)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func GetExperimentRunsHistory(params experiment_runs.GetExperimentRunsHistoryParams) middleware.Responder {
	operation := "GetExperimentRunsHistory"
	user := getUserID(params.HTTPRequest)
	isSA :=  isSuperAdmin(params.HTTPRequest)
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	experimentRuns, count, err := experimentRunService.GetRunHistory(params.ExpID, *params.Page, *params.Size, queryStr, user, isSA)
	if err != nil{
		log.Error("Get ExperimentRun's history failed, ",err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	var expRuns []*restmodels.ProphecisExperimentRun
	err = copier.Copy(&expRuns, &experimentRuns)
	if err != nil {
		log.Error("Copy Experiment Runs Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	response := restmodels.ProphecisExperimentRuns{
		ExperimentRuns: 	expRuns,
		PageNumber:     		*params.Page,
		PageSize:       			*params.Size,
		Total:          	 			count,
		TotalPage:      			int64(math.Ceil(float64(count) / float64(*params.Size))),
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		log.Error("json.Marshal failed, ",err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

