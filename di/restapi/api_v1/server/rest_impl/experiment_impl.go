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
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"webank/DI/commons"
	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/experiments"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"

	//"webank/DI/restapi/service"
	service "webank/DI/restapi/newService"
)

//var experimentService = service.ExperimentService
var experimentService = service.ExperimentService
var log = logger.GetLogger()

func CreateExperiment(params experiments.CreateExperimentParams) middleware.Responder {
	operation := "CreateExperiment"
	currentUserId := getUserID(params.HTTPRequest)
	createType := ""
	if params.Experiment != nil && params.Experiment.CreateType != nil {
		createType = *params.Experiment.CreateType
	}
	if len(*params.Experiment.ExpName) <= 0 || len(*params.Experiment.ExpName) <= 0 {
		err := errors.New("ExpName's length <= 0 or ExpDesc's length <= 0")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	experiment, err := experimentService.CreateExperiment(params.Experiment, currentUserId, createType)
	if err != nil {
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	res := restmodels.ProphecisExperimentIDResponse{}
	err = copier.Copy(&res, experiment)
	if err != nil {
		log.Error("Copy Experiment To Response Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		log.WithError(err).Errorf(operation + "json.Marshal failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func CopyExperiment(params experiments.CopyExperimentParams) middleware.Responder {
	operation := "DeleteExperiment"
	isSA := isSuperAdmin(params.HTTPRequest)
	user := getUserID(params.HTTPRequest)
	exp, err := experimentService.CopyExperiment(params.ID, *params.CreateType, user, isSA)
	res := restmodels.ProphecisExperimentIDResponse{}
	err = copier.Copy(&res, exp)
	if err != nil {
		log.Error("Copy Experiment To Response Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	marshal, err := json.Marshal(res)
	if err != nil {
		log.WithError(err).Errorf(operation + "json.Marshal failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)

}

func UpdateExperiment(params experiments.UpdateExperimentParams) middleware.Responder {
	operation := "UpdateExperiment"
	currentUserId := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	_, err := experimentService.UpdateExperiment(params.Experiment, currentUserId, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.WithError(err).Errorf("experimentService.UpdateExperiment failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, nil)
}

func UpdateExperimentInfo(params experiments.UpdateExperimentInfoParams) middleware.Responder {
	operation := "UpdateExperiment"
	currentUserId := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	_, err := experimentService.UpdateExperimentInfo(params.ID, params.Experiment.ExpName, params.Experiment.ExpDesc, params.Experiment.GroupName,
		params.Experiment.TagList, currentUserId, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.WithError(err).Errorf("experimentService.UpdateExperiment failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, nil)
}

func GetExperiment(params experiments.GetExperimentParams) middleware.Responder {
	operation := "GetExperiment"
	currentUserId := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	experiment, flowJson, err := experimentService.GetExperiment(params.ID, currentUserId, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.WithError(err).Errorf("ExperimentService.GetExperiment failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	prophecisExperiment := restmodels.ProphecisExperiment{}
	err = copier.Copy(&prophecisExperiment, experiment)
	if err != nil {
		log.WithError(err).Errorf("Copier.Copy experiment failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	newResponse := restmodels.ProphecisExperimentGetResponse{
		ProphecisExperiment: &prophecisExperiment,
		FlowJSON:            flowJson,
	}
	marshal, err := json.Marshal(newResponse)
	if err != nil {
		log.WithError(err).Errorf("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}

func ListExperiments(params experiments.ListExperimentsParams) middleware.Responder {
	operation := "listExperiments"
	user := getUserID(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)

	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	//if experimentService == nil {
	//	log.Info("experimentSerivce is NIl")
	//	log.Info("experimentSerivce is NIl,%s,%s,%s",*params.Size, *params.Page, *params.QueryStr)
	//}

	expsModel, count, err := experimentService.ListExperiments(*params.Size, *params.Page, user, queryStr, isSA)

	var exps []*restmodels.ProphecisExperiment
	err = copier.Copy(&exps, &expsModel)
	if err != nil {
		log.Error("Copy Experiments Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	//log.Debug("experiments len is %s",len(exps))
	//log.Debug("DAO experiments len is %s",len(expsModel))
	response := restmodels.ProphecisExperiments{
		Experiments: exps,
		TotalPage:   int64(math.Ceil(float64(count) / float64(*params.Size))),
		Total:       count,
		PageNumber:  *params.Page,
		PageSize:    *params.Size,
	}
	marshal, err := json.Marshal(response)
	if err != nil {
		log.Error(operation + "json.Marshal failed: " + err.Error())
		return httpResponseHandle(http.StatusOK, nil, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, nil, operation, marshal)
}

func CreateExperimentTag(params experiments.CreateExperimentTagParams) middleware.Responder {
	operation := "CreateExperimentTag"
	user := getUserID(params.HTTPRequest)
	expTag, err := experimentService.CreateExperimentTag(params.ExperimentTag.ExpID, *params.ExperimentTag.ExpTag, user)
	if err != nil {
		log.Error("Delete Experiment Tag Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}

	expTagRest := restmodels.ProphecisExperimentTag{}
	err = copier.Copy(&expTagRest, expTag)
	if err != nil {
		log.Error("Copy ExperimentTag To Response Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	marshal, err := json.Marshal(expTagRest)
	if err != nil {
		log.Error(operation + "json.Marshal failed: " + err.Error())
		return httpResponseHandle(http.StatusOK, nil, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)

}

func DeleteExperimentTag(params experiments.DeleteExperimentTagParams) middleware.Responder {
	operation := "DeleteExperimentTag"
	isSA := isSuperAdmin(params.HTTPRequest)
	user := getUserID(params.HTTPRequest)
	err := experimentService.DeleteExperimentTag(params.ID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}

	if err != nil {
		log.Error("Delete Experiment Tag Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, nil)
}

func DeleteExperiment(params experiments.DeleteExperimentParams) middleware.Responder {
	operation := "DeleteExperiment"
	isSA := isSuperAdmin(params.HTTPRequest)
	user := getUserID(params.HTTPRequest)
	err := experimentService.DeleteExperiment(params.ID, user, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.Error("Delete Experiment Error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, nil)

}

func ExportExperiment(params experiments.ExportExperimentParams) middleware.Responder {
	operation := "ExportExperiment"
	currentUserId := getUserID(params.HTTPRequest)
	superAdmin := getSuperadmin(params.HTTPRequest)
	isSA := isSuperAdmin(params.HTTPRequest)
	log.Debugf(">>exportExperiment<< metaUserID: %v, superadmin: %v", currentUserId, superAdmin)

	expId := params.ID

	zipFile, err := service.ExperimentService.Export(expId, currentUserId, isSA)
	if err == commons.PermissionDeniedError {
		return httpResponseHandle(http.StatusUnauthorized, err, operation, nil)
	}
	if err != nil {
		log.WithError(err).Errorf("ExperimentService.GetExperiment failed")
		//return error500(log, err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	fileBytes, err := ioutil.ReadAll(zipFile)
	defer zipFile.Close()
	if err != nil {
		log.Error(err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	if len(fileBytes) <= 0 {
		log.Error("fileBytes is null")
		return httpResponseHandle(http.StatusInternalServerError, errors.New("fileBytes is null"), operation, []byte("fileBytes is null"))
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		/*defer storerClient.Close()

		//chunk, err := stream.Recv()
		//fileBytes, err := ioutil.ReadAll(zipFile)

		if err == io.EOF {
			log.Debugf("final write")

			w.WriteHeader(200)
			//if err = stream.CloseSend(); err != nil {
			//	getLogger.WithError(err).Error("Closing stream failed")
			//}
			break
		} else if err != nil {
			log.WithError(err).Error("Cannot read trained model")
			////FIXME MLSS Change: change Description to Msg
			w.Header().Set(runtime.HeaderContentType, runtime.JSONMime)
			w.WriteHeader(500)
			errResp := restmodels.Error{Code: http.StatusInternalServerError, Msg: err.Error()}
			payload, _ := json.Marshal(&errResp)
			w.Write(payload)
			break
		}
		//getLogger.Debugf("stream.Recv len, %v", len(chunk.ChucksData))

		//n, err := w.Write(chunk.ChucksData)

		//getLogger.Debugf("w.Write len, %v", n)
		fileBytes, err := ioutil.ReadAll(zipFile)
		writeLen, err := w.Write(fileBytes)
		if err != nil {
			log.Errorf("w.Write( failed, %v", err.Error())
		}*/
		_, err = w.Write(fileBytes)
		if err != nil {
			log.Errorf("w.Write( failed, %v", err.Error())
		}
	})
}

func ImportExperiment(params experiments.ImportExperimentParams) middleware.Responder {
	logr := logger.GetLogger()
	operation := "ImportExperiment"
	var currentUserId = getUserID(params.HTTPRequest)
	var superAdmin = getSuperadmin(params.HTTPRequest)
	logr.Debugf(">>importExperiment<< metaUserID: %v, superadmin: %v", currentUserId, superAdmin)
	if params.CreateType == nil || *params.CreateType == "" {
		params.CreateType = experiments.NewImportExperimentParams().CreateType
	}

	//import zip 文件

	// 创建flow(exp)

	// 解析json nodes id, codePath, codePaths 关系

	//上传code.zip, 更新node 中 codePath

	// update exp(flowJson)
	closer := params.File
	fileName := params.FileName
	expIdResp, err := service.ExperimentService.Import(closer, *params.ExperimentID, fileName,
		currentUserId, *params.CreateType, "NULL")
	// expIdResp, err := service.ExperimentService.Import(closer, *params.ExperimentID, fileName, currentUserId, "MLFLOW","NULL")
	if err != nil {
		logr.WithError(err).Errorf("json.Marshal(mp) failed")
		//return error500(logr, err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	response := restmodels.ProphecisExperimentIDResponse{
		ID:      expIdResp.ID,
		ExpName: expIdResp.ExpName,
	}

	marshal, err := json.Marshal(&response)
	if err != nil {
		logr.WithError(err).Errorf("json.Marshal(mp) failed")
		return handleErrorResponse(logr, err.Error())
	}
	return httpResponseHandle(http.StatusOK, err, "importExperiment", marshal)

	//var result = commonModels.Result{
	//	Code:   "200",
	//	Msg:    "success",
	//	Result: json.RawMessage(marshal),
	//}
	//return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
	//	payload, _ := json.Marshal(result)
	//	w.Write(payload)
	//	w.WriteHeader(200)
	//})
}

func CodeUpload(params experiments.CodeUploadParams) middleware.Responder {
	operation := "CodeUpload"
	log.Info("experimentService")
	log.Info(experimentService)
	log.Info("File")
	log.Info(params.File)

	s3Path, err := experimentService.UploadCodeBucket(params.File, "di-model")
	if err != nil {
		logger.GetLogger().Error("UploadCodeBucket err, ", err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	res := restmodels.CodeUploadResponse{S3Path: s3Path}

	marshal, err := json.Marshal(&res)
	if err != nil {
		logger.GetLogger().Error("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)

}

func ExportExperimentDss(params experiments.ExportExperimentDssParams) middleware.Responder {
	operation := "ExportExperimentDss"
	currentUserId := getUserID(params.HTTPRequest)
	superAdmin := getSuperadmin(params.HTTPRequest)
	log.Debugf(">>exportExperimentDss<< metaUserID: %v, superadmin: %v", currentUserId, superAdmin)
	expId := params.ID
	res, err := experimentService.ExportDSS(expId, currentUserId)
	marshal, err := json.Marshal(&res)
	if err != nil {
		log.Println("export experiment dss error, ", err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, marshal)
	}
	return httpResponseHandle(http.StatusOK, err, "exportExperimentDss", marshal)
}

func ImportExperimentDss(params experiments.ImportExperimentDssParams) middleware.Responder {
	operation := "ImportExperimentDss"
	currentUserId := getUserID(params.HTTPRequest)
	superAdmin := getSuperadmin(params.HTTPRequest)
	log.Debugf(">>importExperimentDss<< metaUserID: %v, superadmin: %v", currentUserId, superAdmin)
	resourceId := params.ResourceID
	version := params.Version
	desc := params.Desc
	res, err := experimentService.ImportDSS(*params.ExperimentID, resourceId, version, desc, currentUserId)
	marshal, err := json.Marshal(&res)
	if err != nil {
		log.Println("import experiment dss error, ", err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, marshal)
	}
	return httpResponseHandle(http.StatusOK, err, "importExperimentDss", marshal)
}

func isSuperAdmin(r *http.Request) bool {
	superAdminStr := r.Header.Get(config.CcAuthSuperadmin)
	if "true" == superAdminStr {
		return true
	}
	return false
}
