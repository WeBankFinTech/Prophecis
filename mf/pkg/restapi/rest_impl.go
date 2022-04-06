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

package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/config"
	cc "mlss-mf/pkg/common/controlcenter"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/models"
	"mlss-mf/pkg/restapi/operations/container"
	"mlss-mf/pkg/restapi/operations/model_deploy"
	"mlss-mf/pkg/restapi/operations/model_storage"
	"mlss-mf/pkg/service"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

var (
	mylogger *logrus.Logger
)

func init() {
	mylogger = logger.Logger()
}

type result struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

var (
	ccClient   = cc.GetCcClient(config.GetMFconfig().ServerConfig.CCAddress)
	serviceDao = dao.SerivceDao{}
	userDao    = dao.UserDao{}
	modelDao   = dao.ModelDao{}
)

func ServicePost(params model_deploy.PostServiceParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	//判断是否有NS权限 or 判断是否为GA或SA
	if params.Service.Namespace == nil {
		logger.Logger().Error("Namespace is nil")
		err := errors.New("Namespace is nil")
		return returnResponse(service.StatusError, err, "servicePost", []byte(err.Error()))
	}
	auth := service.CheckAuth(*params.Service.Namespace, currentUserId, token, isSA)
	if !auth {
		logger.Logger().Error("checkAuth failed, ", fmt.Sprintf("namespace: %v, currentUserId: %v, token: %v, isSA: %v", *params.Service.Namespace, currentUserId, token, isSA))
		return returnResponse(service.StatusAuthReject, nil, "servicePost", nil)
	}
	ContainerEngineCPU := config.GetMFconfig().ResourceConfig.ContainerEngineCpu
	ContainerEngineMemory := config.GetMFconfig().ResourceConfig.ContainerEngineMemory
	ContainerEngineGPU := config.GetMFconfig().ResourceConfig.ContainerEngineGpu
	IstioSideCarCPU := config.GetMFconfig().ResourceConfig.IstioSideCarCpu
	IstioSideCarMemory := config.GetMFconfig().ResourceConfig.IstioSideCarMemory
	IstioSideCarGPU := config.GetMFconfig().ResourceConfig.IstioSideCarGpu
	cpu := *params.Service.CPU + ContainerEngineCPU + IstioSideCarCPU
	gpu := "0"
	memory := *params.Service.Memory + ContainerEngineMemory + IstioSideCarMemory
	if params.Service.Gpu != nil {
		if *params.Service.Gpu != "" {
			gpu = *params.Service.Gpu
		}
	}
	gpuInt, err := strconv.Atoi(gpu)
	if err != nil {
		logger.Logger().Error("Gpu atoi err, ", err)
		return returnResponse(service.StatusError, err, "servicePost", []byte(err.Error()))
	}
	containerEngineGpuInt, err := strconv.Atoi(ContainerEngineGPU)
	if err != nil {
		logger.Logger().Error("EngineGpu atoi err, ", err)
		return returnResponse(service.StatusError, err, "servicePost", []byte(err.Error()))
	}
	istioSideCarGpuInt, err := strconv.Atoi(IstioSideCarGPU)
	if err != nil {
		logger.Logger().Error("IstioSideCarGpu atoi err, ", err)
		return returnResponse(service.StatusError, err, "servicePost", []byte(err.Error()))
	}
	gpu = strconv.Itoa(gpuInt + containerEngineGpuInt + istioSideCarGpuInt)
	if modelTypeFlag := service.CheckModelType(params.Service.ServicePostModels[0].ModelType); !modelTypeFlag {
		err = errors.New("Service create failed, model type cannot is " + params.Service.ServicePostModels[0].ModelType)
		return returnResponse(service.StatusBadRequest, err, "servicePost", []byte(err.Error()))
	}
	//resource quota
	flag, err := ccClient.CheckResource(*params.Service.Namespace, cpu, gpu, memory, token)
	if flag == false || err != nil {
		logger.Logger().Error("Add service err, ", err)
		return returnResponse(service.StatusError, err, "servicePost", []byte(errors.New("Add service err, "+err.Error()).Error()))
	}
	si, err := service.ServiceCreate(params, currentUserId, gpu)
	serviceJson, err := json.Marshal(&si)
	return returnResponse(http.StatusOK, err, "servicePost", serviceJson)
}

func ServiceList(params model_deploy.ListServicesParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	pageModel, status, err := service.ServiceList(params, currentUserId, isSA)
	resultJson, err := json.Marshal(pageModel)

	return returnResponse(status, err, "serviceList", resultJson)
}

func ServiceGet(params model_deploy.GetServiceParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	var newService models.GetService
	getService, err := serviceDao.GetService(params.ID)
	if err != nil {
		logger.Logger().Error("Get Service error, ", err)
		return returnResponse(service.StatusError, err, "serviceGet", []byte(err.Error()))
	}
	auth := service.CheckUserAuth(currentUserId, getService.GroupID, getService.UserID, isSA)
	if !auth {
		err := errors.New("Check User Auth Error")
		return returnResponse(service.StatusAuthReject, err, "serviceGet", nil)
	}
	err = copier.Copy(&newService, &getService)
	if err != nil {
		logger.Logger().Error("copy service err, ", err)
		return returnResponse(service.StatusError, err, "serviceGet", []byte(err.Error()))
	}
	serviceRes, status, err := service.ServiceGet(params.ID, getService, &newService)
	serviceJson, err := json.Marshal(serviceRes)
	return returnResponse(status, err, "serviceGet", serviceJson)
}

func ServiceDashborad(params model_deploy.ServiceDashboradParams) middleware.Responder {
	isSA := false

	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}

	token := params.HTTPRequest.Header.Get(common.CcAuthToken)
	getDashborad, status, err := service.ServiceDashborad(token, isSA)
	dashBoradJson, err := json.Marshal(getDashborad)
	return returnResponse(status, err, "serviceDashborad", dashBoradJson)
}

func ServiceDelete(params model_deploy.DeleteServiceParams) middleware.Responder {
	// currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	// isSA := false
	// if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
	// 	isSA = true
	// }
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)

	getService, err := serviceDao.GetService(params.ID)
	if err != nil {
		logger.Logger().Error("Get Service error, ", err)
		return returnResponse(service.StatusError, err, "serviceGet", []byte(err.Error()))
	}
	// auth := service.CheckGroupAuth(currentUserId, getService.GroupID, isSA)
	// if !auth {
	// 	return returnResponse(service.StatusAuthReject, nil, "serviceGet", nil)
	// }
	err = service.PermissionCheck(user, getService.UserID, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return returnResponse(http.StatusForbidden, err, "ServiceDelete", nil)
	}

	status, err := service.ServiceDelete(params.ID, getService)
	return returnResponse(status, err, "serviceDelete", nil)
}

func ServiceUpdate(params model_deploy.UpdateServiceParams) middleware.Responder {
	// currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	// isSA := false
	// if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
	// isSA = true
	// }
	// getService, err := serviceDao.GetService(params.ID)
	// if err != nil {
	// logger.Logger().Error("Get Service error, ", err)
	// return returnResponse(service.StatusError, err, "serviceUpdate", []byte(err.Error()))
	// }
	// auth := service.CheckGroupAuth(currentUserId, getService.GroupID, isSA)
	// if !auth {
	// 	return returnResponse(service.StatusAuthReject, nil, "serviceUpdate", nil)
	// }
	//resource quota
	ContainerEngineCPU := config.GetMFconfig().ResourceConfig.ContainerEngineCpu
	ContainerEngineMemory := config.GetMFconfig().ResourceConfig.ContainerEngineMemory
	ContainerEngineGPU := config.GetMFconfig().ResourceConfig.ContainerEngineGpu
	IstioSideCarCPU := config.GetMFconfig().ResourceConfig.IstioSideCarCpu
	IstioSideCarMemory := config.GetMFconfig().ResourceConfig.IstioSideCarMemory
	IstioSideCarGPU := config.GetMFconfig().ResourceConfig.IstioSideCarGpu
	cpu := params.Service.CPU + ContainerEngineCPU + IstioSideCarCPU
	gpu := "0"
	memory := params.Service.Memory + ContainerEngineMemory + IstioSideCarMemory
	if params.Service.Gpu != "" {
		gpuInt, err := strconv.Atoi(params.Service.Gpu)
		if err != nil {
			logger.Logger().Error("Gpu atoi err, ", err)
			return returnResponse(service.StatusError, err, "serviceUpdate", []byte(err.Error()))
		}
		containerEngineGpuInt, err := strconv.Atoi(ContainerEngineGPU)
		if err != nil {
			logger.Logger().Error("ContainerEngineGPU atoi err, ", err)
			return returnResponse(service.StatusError, err, "serviceUpdate", []byte(err.Error()))
		}
		istioSideCarGpuInt, err := strconv.Atoi(IstioSideCarGPU)
		if err != nil {
			logger.Logger().Error("istioSideCarGpuInt atoi err, ", err)
			return returnResponse(service.StatusError, err, "serviceUpdate", []byte(err.Error()))
		}
		gpu = strconv.Itoa(gpuInt + containerEngineGpuInt + istioSideCarGpuInt)
	}
	//resource quota
	flag, err := ccClient.CheckResource(params.Service.Namespace, cpu, gpu, memory, token)
	if flag == false || err != nil {
		logger.Logger().Error("Update service err, ", err)
		return returnResponse(service.StatusError, err, "serviceDelete", []byte(err.Error()))
	}
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	status, err := service.ServiceUpdate(params, gpu, user, isSA)
	return returnResponse(status, err, "serviceUpdate", nil)
}

func ServiceStop(params model_deploy.StopNamespacedServiceParams) middleware.Responder {
	// currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	// token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	// isSA := false
	// if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
	// isSA = true
	// }
	// namespace := params.Namespace
	//auth := service.CheckAuth(namespace, currentUserId, token, isSA)
	// if !auth {
	// 	return returnResponse(service.StatusAuthReject, nil, "serviceStop", nil)
	// }
	isSA := isSuperAdmin(params.HTTPRequest)
	user := getUserID(params.HTTPRequest)
	status, err := service.ServiceStop(params, user, isSA)
	return returnResponse(status, err, "serviceStop", nil)
}

func ServiceRun(params model_deploy.CreateNamespacedServiceRunParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	// token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	// isSA := false
	// if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
	// isSA = true
	// }
	namespace := *params.Service.Namespace
	// auth := service.CheckAuth(namespace, currentUserId, token, isSA)
	// if !auth {
	// 	return returnResponse(service.StatusAuthReject, nil, "serviceRun", nil)
	// }
	status, err := service.ServiceRun(params, currentUserId, namespace, isSA)
	return returnResponse(status, err, "serviceRun", nil)
}

//func PostNamespacedModel(params model_deploy.PostNamespacedServiceParams) middleware.Responder {
//	//currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
//	//token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
//	//model, status, err := service.ModelAddByService(params, params.Model.GroupID, currentUserId, token, params.Model.)
//	//modelJson, err  := json.Marshal(model)
//	return returnResponse(http.StatusInternalServerError, errors.New("PostNamespacedModel is not implemented yet"), "postNamespacedModel",nil)
//}

func PostModel(params model_storage.PostModelParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := service.CheckGroupAuth(currentUserId, *params.Model.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, "addModel", nil)
	}
	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		return returnResponse(service.StatusError, err, "addModel", []byte(err.Error()))
	}
	modelFlag, err := modelDao.CheckModelByModelNameAndGroupId(*params.Model.ModelName, *params.Model.GroupID)
	if err != nil {
		return returnResponse(service.StatusError, err, "addModel", []byte(err.Error()))
	}

	if modelTypeFlag := service.CheckModelType(*params.Model.ModelType); !modelTypeFlag {
		err := errors.New("Model create failed, model type cannot is " + *params.Model.ModelType)
		return returnResponse(service.StatusError, err, "addModel", []byte(err.Error()))
	}

	// if modelFlag {
	// 	err = errors.New("The same group model name cannot be repeated")
	// 	return returnResponse(service.StatusBadRequest, err, "addModel", []byte(err.Error()))
	// }
	logger.Logger().Debugf("post model modelflag: %v, model name: %s, group id: %d",
		modelFlag, *params.Model.ModelName, *params.Model.GroupID)
	resp := models.PostModelResp{}
	switch modelFlag {
	case true:
		// update model verson
		modelResp, statusCode, err := service.AddModelVersion(params, currentUserId, token)
		if err != nil {
			return returnResponse(statusCode, err, "addModel", []byte(err.Error()))
		}
		resp.ModelID = modelResp.ModelID
		resp.ModelVersionID = modelResp.ModelVersionID
		resp.ModelVersion = modelResp.ModelVersion
	default:
		modelResp, status, err := service.AddModel(params, currentUserId, token, user)
		if err != nil {
			return returnResponse(status, err, "addModel", []byte(err.Error()))
		}
		resp.ModelID = modelResp.ModelID
		resp.ModelVersionID = modelResp.ModelVersionID
		resp.ModelVersion = modelResp.ModelVersion
	}

	bytes, err := json.Marshal(&resp)
	if err != nil {
		return returnResponse(service.StatusError, err, "addModel", []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, err, "addModel", bytes)
}

// func PostModel(params model_storage.PostModelParams) middleware.Responder {
// 	modelID, err := service.CreateModel(params)
// 	if err != nil {
// 		return returnResponse(http.StatusInternalServerError, err, "PostModel", nil)
// 	}
// 	resp := models.PostModelResp{
// 		ModelID: modelID,
// 	}
// 	respBytes, err := json.Marshal(&resp)
// 	if err != nil {
// 		return returnResponse(http.StatusInternalServerError, err, "PostModel", nil)
// 	}
// 	return returnResponse(http.StatusOK, nil, "PostModel", respBytes)
// }

func DeleteModel(params model_storage.DeleteModelParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	model, err := modelDao.GetModel(strconv.Itoa(int(params.ModelID)))
	if err != nil {
		logger.Logger().Error("GetModel err, ", err)
		return returnResponse(service.StatusError, nil, "deleteModel", []byte(err.Error()))
	}
	auth := service.CheckGroupAuth(currentUserId, model.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, "deleteModel", nil)
	}
	if flag := service.CheckServiceRunState(params.ModelID); flag == false {
		err = errors.New("The service bound to this model is running")
		return returnResponse(service.StatusError, err, "deleteModel", []byte(err.Error()))
	}
	status, err := service.ModelDelete(&model)
	if err != nil {
		return returnResponse(status, err, "deleteModel", []byte(err.Error()))
	}
	return returnResponse(status, err, "deleteModel", nil)
}

func GetModel(params model_storage.GetModelParams) middleware.Responder {
	model, status, err := service.ModelGet(params.ModelID)
	modelJson, err := json.Marshal(model)
	if err != nil {
		return returnResponse(status, err, "getModel", []byte(err.Error()))
	}
	return returnResponse(status, err, "getModel", modelJson)
}

func UpdateModel(params model_storage.UpdateModelParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := service.CheckGroupAuth(currentUserId, params.Model.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, "updateModel", nil)
	}
	status, err := service.ModelUpdate(params, currentUserId, token)
	if err != nil {
		return returnResponse(status, err, "updateModel", []byte(err.Error()))
	}
	return returnResponse(status, err, "updateModel", nil)
}

func ListModels(params model_storage.GetModelsParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	resMap := map[string]interface{}{}
	modelList, count, status, err := service.ModelList(currentUserId, *params.Size, *params.Page, queryStr, isSA)
	resMap["models"] = modelList
	resMap["pageNumber"] = *params.Page
	resMap["total"] = count
	resMap["totalPage"] = int64(math.Ceil(float64(count) / float64(*params.Size)))
	modelsJson, err := json.Marshal(resMap)
	if err != nil {
		return returnResponse(status, err, "listModels", []byte(err.Error()))
	}
	return returnResponse(status, err, "listModels", modelsJson)
}

func ListModelsByCluster(params model_storage.GetModelsByClusterParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	resMap := map[string]interface{}{}
	modelList, count, status, err := service.ModelListByCluster(currentUserId, *params.Size, *params.Page, queryStr, params.Cluster, isSA)
	resMap["models"] = modelList
	resMap["pageNumber"] = *params.Page
	resMap["total"] = count
	resMap["totalPage"] = int64(math.Ceil(float64(count) / float64(*params.Size)))
	modelsJson, err := json.Marshal(resMap)
	if err != nil {
		return returnResponse(status, err, "listModels", []byte(err.Error()))
	}
	return returnResponse(status, err, "listModels", modelsJson)
}

func GetModelsByGroupIdAndModelName(params model_storage.GetModelsByGroupIDAndModelNameParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	respMap := map[string]interface{}{}
	modelList, count, err := service.ModelListByGroupIdAndModelName(*params.PageSize,
		*params.CurrentPage, params.GroupID, params.ModelName, currentUserId)
	respMap["models"] = modelList
	respMap["total"] = count
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetModelsByGroupIdAndModelName", []byte(err.Error()))
	}
	modelsJson, err := json.Marshal(respMap)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetModelsByGroupIdAndModelName", []byte(err.Error()))
	}

	return returnResponse(http.StatusOK, err, "GetModelsByGroupIdAndModelName", modelsJson)
}

func ListModelsByGroup(params model_storage.ListModelsByGroupIDParams) middleware.Responder {
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	resMap := map[string]interface{}{}
	modelList, count, status, err := service.ModelListByGroups(*params.Size, *params.Page, params.GroupID, queryStr, currentUserId)
	resMap["models"] = modelList
	resMap["pageNumber"] = *params.Page
	resMap["total"] = count
	resMap["totalPage"] = int64(math.Ceil(float64(count) / float64(*params.Size)))
	modelsJson, err := json.Marshal(resMap)
	if err != nil {
		return returnResponse(status, err, "listModelsByGroup", []byte(err.Error()))
	}
	return returnResponse(status, err, "listModelsByGroup", modelsJson)
}

func UploadModel(params model_storage.UploadModelParams) middleware.Responder {
	s3Path, status, err := service.UploadModel(params.FileName, params.ModelType, params.File)
	res := models.UploadModelResponse{
		S3Path:   s3Path,
		FileName: params.FileName,
	}
	resJson, err := json.Marshal(&res)
	if err != nil {
		return returnResponse(status, err, "uploadModel", []byte(err.Error()))
	}
	return returnResponse(status, err, "uploadModel", resJson)
}

func returnResponse(status int, err error, operation string, resultMsg []byte) middleware.Responder {
	result := result{Code: 200, Message: "success"}
	if resultMsg != nil {
		result.Result = json.RawMessage(resultMsg)
	} else {
		resultBytes, _ := json.Marshal((operation + " " + "success"))
		result.Result = resultBytes
	}
	if status != service.StatusSuccess {
		result.Code = status
		jsonStr, _ := json.Marshal(operation + " " + " error")
		result.Result = jsonStr
		if len(string(resultMsg)) > 0 {
			result.Message = string(resultMsg)
		} else {
			result.Message = err.Error()
		}
		if err != nil {
			mylogger.Error("Operation: " + operation + ";" + "Error: " + err.Error())
			jsonStr, _ := json.Marshal(operation + " " + " error " + err.Error())
			result.Result = jsonStr
		}
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(status)
		response, err := json.Marshal(&result)
		// logger.Logger().Debugf("result in %q , result: %+v, response: %s\n", operation, result, string(response))
		if err != nil {
			logger.Logger().Errorf("fail to marshal response in %q, err: %s\n", operation, err.Error())
			return
		}
		_, err = w.Write(response)
		if err != nil {
			logger.Logger().Errorf("fail to write response in %q, err: %s\n", operation, err.Error())
			return
		}
	})
}

func ListModelVersions(params model_storage.GetModelVersionParams) middleware.Responder {
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	resMap := map[string]interface{}{}
	modelVersionList, count, status, err := service.ModelVersionList(params.ModelID, *params.Size, *params.Page,
		queryStr, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "ListModelVersions", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "ListModelVersions", nil)
	}
	resMap["modelVersions"] = modelVersionList
	resMap["pageNumber"] = *params.Page
	resMap["total"] = count
	resMap["totalPage"] = int64(math.Ceil(float64(count) / float64(*params.Size)))
	modelVersionListJson, err := json.Marshal(resMap)
	if err != nil {
		return returnResponse(status, err, "listModelVersions", []byte(err.Error()))
	}
	return returnResponse(status, err, "listModelVersions", modelVersionListJson)
}

func ExportModel(params model_storage.ExportModelParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	fileBytes, status, err := service.ExportModel(currentUserId, params.ModelID, isSA)
	if err != nil {
		return returnResponse(status, err, "exportModel", []byte(err.Error()))
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(status)
		w.Write(fileBytes)
	})
}

func ListContainer(params container.ListContainerParams) middleware.Responder {
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	//TODO
	containers, status, err := service.ListContainer(params, user, isSA)
	if err != nil {
		return returnResponse(status, err, "ListContainer", nil)
	}
	resultJson, err := json.Marshal(containers)

	return returnResponse(status, err, "ListContainer", resultJson)
}

func DownloadModelByID(params model_storage.DownloadModelByIDParams) middleware.Responder {
	currentUserID := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	modelID := params.ModelID
	bytes, statusCode, fileName, err := service.DownloadModel(modelID, currentUserID, isSA)
	if err != nil {
		return returnResponse(statusCode, err, "DownloadModelByID", nil)
	}
	// return returnResponse(http.StatusOK, err, "DownloadModelByID", bytes)
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.Header().Add("File-Name", fileName)
		w.WriteHeader(statusCode)
		w.Write(bytes)
	})
}

func DownloadModelVersionByID(params model_storage.DownloadModelVersionByIDParams) middleware.Responder {
	currentUserID := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	modelVersionID := params.ModelVersionID
	bytes, statusCode, fileName, err := service.DownloadModelVersion(modelVersionID, currentUserID, isSA)
	if err != nil {
		return returnResponse(statusCode, err, "DownloadModelVersionByID", nil)
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.Header().Add("File-Name", fileName)
		w.WriteHeader(statusCode)
		w.Write(bytes)
	})
}

func PushModelByModelId(params model_storage.PushModelByModelIDParams) middleware.Responder {
	event, err := service.PushModelByModelID(params.ModelID, params.ModelPushEvent)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "PushModelByModelId", nil)
	}
	resultJson, err := json.Marshal(event)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "PushModelByModelId", nil)
	}
	return returnResponse(http.StatusOK, err, "PushModelByModelId", resultJson)
}

func PushModelByModelVersionId(params model_storage.PushModelByModelVersionIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	event, err := service.PushModelByModelVersionID(params.ModelVersionID, params.ModelPushEvent, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "PushModelByModelVersionId", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "PushModelByModelVersionId", nil)
	}
	resultJson, err := json.Marshal(event)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "PushModelByModelVersionId", nil)
	}
	return returnResponse(http.StatusOK, err, "PushModelByModelVersionId", resultJson)
}

func ListModelVersionPushEventsByModelVersionID(params model_storage.ListModelVersionPushEventsByModelVersionIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	resp, err := service.ListModelVersionPushEventsByModelVersionID(params.ModelVersionID, *params.CurrentPage, *params.PageSize, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "ListModelVersionPushEventsByModelVersionID", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "ListModelVersionPushEventsByModelVersionID", nil)
	}
	resultJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "ListModelVersionPushEventsByModelVersionID", nil)
	}
	return returnResponse(http.StatusOK, err, "ListModelVersionPushEventsByModelVersionID", resultJson)
}

func GetModelVersionByNameAndGroupID(params model_storage.GetModelVersionByNameAndVersionParams) middleware.Responder {
	res, _, err := service.GetModelVersionByNameAndGroupID(params)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetModelVersionByNameAndGroupID", nil)
	}
	resultJson, err := json.Marshal(res)
	if err != nil {
		logger.Logger().Error("GetModelVersionByNameAndGroupID resultJson marshal err: ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, "GetModelVersionByNameAndGroupID", nil)
	}
	return returnResponse(http.StatusOK, err, "GetModelVersionByNameAndGroupID", resultJson)
}

func isSuperAdmin(r *http.Request) bool {
	superAdminStr := r.Header.Get(cc.CcSuperadmin)
	if "true" == superAdminStr {
		return true
	}
	return false
}

func getUserID(r *http.Request) string {
	return r.Header.Get(cc.CcAuthUser)
}
