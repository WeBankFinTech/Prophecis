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

package rest_impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"webank/DI/commons/constants"
	"webank/DI/jobmonitor/jobmonitor"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/pkg/operators/tf-operator/apis/tensorflow/v1alpha1"
	"webank/DI/pkg/repo"
	"webank/DI/restapi/api_v1/server/operations"
	"webank/DI/restapi/service"
	"webank/DI/storage/storage/grpc_storage"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/mholt/archiver/v3"
	"github.com/modern-go/reflect2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"

	commonModels "webank/DI/commons/models"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/models"
	storageClient "webank/DI/storage/client"
	trainerClient "webank/DI/trainer/client"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"bufio"
	"bytes"
	"io"

	"webank/DI/commons/logger"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	newWebsocket "github.com/gorilla/websocket"
	"golang.org/x/net/websocket"

	stringsUtil "strings"
	"time"

	trainingDataClient "webank/DI/metrics/client"
	"webank/DI/restapi/api_v1/server/operations/training_data"

	"webank/DI/metrics/service/grpc_training_data_v1"

	cc "webank/DI/commons/controlcenter/client"

	"webank/DI/commons/config"

	"github.com/spf13/viper"
)

const (
	defaultLogPageSize = 10
	codeFile           = "codeFile"
	storagePath        = "storagePath"
	SA                 = "true"
	userSA             = "1"
	userGA             = "1"
	userGU             = "2"
)

var (
	PermissionDeniedError = errors.New("permission denied")
)

// postModel posts a model definition and starts the training
func PostModel(params models.PostModelParams) middleware.Responder {
	operation := "postModel"
	// Read MainFest File
	logr := logger.LocLogger(logWithPostModelParams(params))
	logr.Debugf("postModel invoked: %v", params.HTTPRequest.Header)
	manifestBytes, err := ioutil.ReadAll(params.Manifest)
	if err != nil {
		logr.Errorf("Cannot read 'manifest' parameter: %s", err.Error())
		return httpResponseHandle(http.StatusUnauthorized, err, operation, []byte("Incorrect parameters"))
	}
	logr.Debug("Loading Manifest")
	manifest, err := LoadManifestV1(manifestBytes)
	if err != nil {
		logr.WithError(err).Errorf("Parameter 'manifest' contains incorrect YAML")
		return httpResponseHandle(http.StatusNotFound, err, operation, []byte("Incorrect manifest"))
	}
	return postModelInclude(params, manifest, operation)
}

// PostModelInclude is postModel include func
func postModelInclude(params models.PostModelParams, manifest *ManifestV1, operation string) middleware.Responder {
	var err error
	modelFile := params.ModelDefinition
	var userID = getUserID(params.HTTPRequest)

	logr := logger.LocLogger(logWithPostModelParams(params))
	logr.Infof("PostModel manifest: %v", manifest)

	//TODO: Remove
	if manifest.Name == "" {
		manifest.Name = manifest.JobType
	}

	logr.Info("Loading Manifest: ", manifest.DssTaskId)
	var modelDefinition []byte
	if modelFile != nil {
		modelBytes, err := ioutil.ReadAll(modelFile)
		logger.GetLogger().Info("Upload code length: ", len(modelBytes), " to bucket")
		var experimentService = service.ExperimentService
		reader := bytes.NewReader(modelBytes)
		readerCloser := ioutil.NopCloser(reader)
		s3Path, err := experimentService.UploadCodeBucket(readerCloser, "di-model")
		if err != nil {
			logger.GetLogger().Error("UploadCode failed, ", err)
		}
		if len(s3Path) <= 0 {
			logger.GetLogger().Error("UpdateCode failed, s3 path's length is zero")
		}
		manifest.FileName = "codeFile.zip"
		manifest.CodePath = s3Path
		params.ModelDefinition = nil
	}
	// IF Codeath is Not Null, Download Code From Minio
	if manifest.CodePath != "" {
		modelDefinitionReader, err := codePathToModelDefinition(manifest.CodePath)
		params.ModelDefinition = modelDefinitionReader
		if err != nil {
			logr.Error("Transfer S3 Path to ModelDefinition Byte Error", err.Error())
			return httpResponseHandle(http.StatusBadRequest, err, operation, nil)
		}
		//log.Info("ModelDefintion is ", modelDefinition)
	}

	if "" != manifest.Framework.Command {
		manifest.Framework.Command, err = CommandTransform(manifest.Framework.Command)
		if err != nil {
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("Incorrect Manifest"))
		}
	}

	if "tfos" == manifest.JobType {
		err := handleTfosManifest(userID, manifest, logr)
		if err != nil {
			return httpResponseHandle(http.StatusUnauthorized, err, operation, []byte("Incorrect Manifest"))
		}
	}

	if "MLPipeline" == manifest.JobType {

	}

	// Check Manifest
	errForCheckManifest := checkManifest(manifest, logr)
	if errForCheckManifest != nil {
		logr.Errorf("checkManifest failed: %s", errForCheckManifest.Error())
		return httpResponseHandle(http.StatusNotFound, err, operation, []byte("checkManifest failed: "+errForCheckManifest.Error()))
	}

	// Check user match storages
	//get first dataStore,if exists
	DSs := manifest.DataStores
	namespaceFromManifest := manifest.Namespace
	dataPath, resultsPath := defineJobPath(DSs, manifest)
	logr.Debugf("namespace: %v, dataPath: %v, resultsPath: %v ", namespaceFromManifest, dataPath, resultsPath)

	//check alert params
	err = checkAlertParams(manifest, logr)
	if err != nil {
		logr.Errorf("checkAlertParams failed: %s", err.Error())
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("checkAlertParams failed: "+err.Error()))
	}

	//Permission Check
	ccAddress := viper.GetString(config.CCAddress)
	ccClient := cc.GetCcClient(ccAddress)
	//Check DataPath,TFOS Job Use HDFS, Not Need to Check Data Path
	if manifest.JobType != "tfos" && manifest.ProxyUser == "" {
		err = ccClient.UserStorageCheck(params.HTTPRequest.Header.Get(cc.CcAuthToken), userID, dataPath)
		if err != nil {
			logr.Errorf("check dataPath failed: %v", err.Error())
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("check dataPath failed: "))
		}
	}

	//Check ResultsPath
	if manifest.ProxyUser == "" {
		err = ccClient.UserStorageCheck(params.HTTPRequest.Header.Get(cc.CcAuthToken), userID, resultsPath)
		if err != nil {
			logr.Errorf("check resultsPath failed: %v", err.Error())
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("Check ResultsPath Failed "+err.Error()))
		}
	}

	if "tfos" != manifest.JobType {
		codeSelector := manifest.CodeSelector
		modelDefinitionRes, err := checkForJobParams(manifest, params, codeSelector, logr)
		if err != nil {
			logr.Errorf("checkForJobParams failed: %v", err.Error())
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("Check ResultsPath Failed "+err.Error()))
		}
		modelDefinition = modelDefinitionRes
	}

	// manifest we do not have a way to select the from a list of dataStores so we cap it to only one.
	if len(manifest.DataStores) > 1 {
		err = errors.New("Please only specify one data_store in the manifest. This constraint will go away with the once the new manifest is in place.")
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("Check ResultsPath Failed "+err.Error()))
	}

	commandStr, err := CommandTransform(manifest.Framework.Command)
	if err != nil {
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("Command Parse Error "+err.Error()))
	}
	manifest.Framework.Command = commandStr

	trainer, err := trainerClient.NewTrainer()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte("Cannot create client for trainer service"+err.Error()))
	}
	defer trainer.Close()

	// TODO do some basic manifest.yml validation to avoid a panic
	// Validate the Manifest File to Avoid panic
	createRequestParam, err := manifest2TrainingRequest(manifest, modelDefinition, params.HTTPRequest, logr, dataPath)
	if err != nil {
		logr.WithError(err).Errorf("CreateRequest parameters check or config failed")
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte(err.Error()))
	}

	//get client_token from db
	userFromDB, err := ccClient.GetUserByName(params.HTTPRequest.Header.Get(cc.CcAuthToken), userID)
	if err != nil {
		logr.Errorf("check dataPath failed: %v", err.Error())
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("GetUserByName failed: "+err.Error()))
	}
	logr.Debugf("userFromDB.Token: %v", userFromDB.Token)
	createRequestParam.Token = userFromDB.Token

	//Proxy User check
	if manifest.ProxyUser != "" {
		auth, err := ccClient.ProxyUserCheck(params.HTTPRequest.Header.Get(cc.CcAuthToken), userID, manifest.ProxyUser)
		if err != nil {
			logr.Errorf("check proxy user permission: %v", err.Error())
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("check proxy user permission failed: "+err.Error()))
		}
		if !auth {
			authStr := strconv.FormatBool(auth)
			logr.Errorf("check proxy user permission: %v", authStr)
			err = errors.New("check proxy user permission: " + authStr)
			return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("check proxy user permission: "+err.Error()))
		}
	}

	logr.Infof("createRequestParam.TFosRequest: %v", createRequestParam.TFosRequest)

	// Check GPU, CPU and Memory
	if createRequestParam.Training.Resources.Cpus <= 0 || createRequestParam.Training.Resources.Gpus < 0 || createRequestParam.Training.Resources.Memory <= 0 {
		err = errors.New("gpu, cpu and memory has to be greater than 0")
		return httpResponseHandle(http.StatusBadRequest, err, operation, []byte("failed to creatTrainingJob when check for resource"))
	}
	//Get data
	db := datasource.GetDB()
	if !reflect2.IsNil(manifest.ExpRunId) && manifest.ExpRunId > 0 {
		createRequestParam.ExpRunId = strconv.Itoa(int(manifest.ExpRunId))
		experimentRun, err := repo.ExperimentRunRepo.Get(manifest.ExpRunId, db)
		if err != nil {
			logger.GetLogger().Error("get experiment run failed, ", err)
		}
		experiment, err := repo.ExperimentRepo.Get(experimentRun.ExpID, db)
		if err != nil {
			logger.GetLogger().Error("get experiment failed, ", err)
		}
		if experiment.ExpName != nil {
			createRequestParam.ExpName = *experiment.ExpName
		}
	}
	if !reflect2.IsNil(manifest.FileName) && len(manifest.FileName) > 0 {
		createRequestParam.FileName = manifest.FileName
	}
	if !reflect2.IsNil(manifest.CodePath) && len(manifest.CodePath) > 0 {
		createRequestParam.FilePath = manifest.CodePath
	}

	// SERVICE METHOD
	tresp, err := trainer.Client().CreateTrainingJob(params.HTTPRequest.Context(), createRequestParam)
	if err != nil {
		logr.Println("== CreateTrainingJobError：", err.Error())
		logr.WithError(err).Errorf("Trainer service call failed")

		if status.Code(err) == codes.InvalidArgument || grpc.Code(err) == codes.NotFound {
			return httpResponseHandle(http.StatusBadRequest, err, operation, nil)
		}

		if status.Code(err) == codes.ResourceExhausted {
			return httpResponseHandle(http.StatusTooManyRequests, err, operation, nil)
		}
		logr.Debugf("model_impl CreateTrainingJob errMsg: %s", err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte("Trainer Interneal Error."))
	}

	loc := params.HTTPRequest.URL.Path + "/" + tresp.TrainingId
	return models.NewPostModelCreated().
		WithLocation(loc).
		WithPayload(&restmodels.BasicNewModel{
			BasicModel: restmodels.BasicModel{
				ModelID: tresp.TrainingId,
			},
			Location: loc,
		})
}

func DeleteModel(params models.DeleteModelParams) middleware.Responder {
	//1、Recv Parmas Init
	logr := logger.LocLogger(logWithDeleteModelParams(params))
	logr.Debugf("deleteModel invoked: %v", params.HTTPRequest.Header)
	storageClient, err := storageClient.NewStorage()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		return handleErrorResponse(logr, err.Error())
	}
	defer storageClient.Close()
	//2、Params Check & Perminssion Check
	// FIXME MLSS Change: auth user in restapi
	_, errResp := checkJobVisibility(params.HTTPRequest.Context(), params.ModelID, getUserID(params.HTTPRequest), params.HTTPRequest.Header.Get(cc.CcAuthToken))
	if errResp != nil {
		return errResp
	}
	_, err = storageClient.Client().DeleteTrainingJob(params.HTTPRequest.Context(), &grpc_storage.DeleteRequest{
		TrainingId: params.ModelID,
		UserId:     getUserID(params.HTTPRequest),
	})
	if err != nil {
		logr.WithError(err).Errorf("storage.Client().DeleteTrainingJob call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		if grpc.Code(err) == codes.NotFound {
			//FIXME MLSS Change: change Description to Msg
			return models.NewDeleteModelNotFound().WithPayload(&restmodels.Error{
				Error: "Not found",
				Code:  http.StatusNotFound,
				Msg:   "",
			})
		}
		return handleErrorResponse(logr, err.Error())
	}
	return models.NewDeleteModelOK().WithPayload(
		&restmodels.BasicModel{
			ModelID: params.ModelID,
		})
}

func GetModel(params models.GetModelParams) middleware.Responder {
	logr := logger.LocLogger(logWithGetModelParams(params))
	logr.Debugf("getModel invoked: %v", params.HTTPRequest.Header)
	rresp, errFromTrainer := getTrainingJobFromStorage(params)
	if rresp == nil {
		return errFromTrainer
	}

	// FIXME MLSS Change: more properties for result
	m := createModel(params.HTTPRequest, rresp.Job, logr)

	logr.Debugf("m: %+v", m)
	return models.NewGetModelOK().WithPayload(m)
}

func RetryModel(params models.RetryModelParams) middleware.Responder {
	var (
		getModelParams  models.GetModelParams
		postModelParams models.PostModelParams
		manifest        *ManifestV1
		err             error
	)
	operation := "retryModel"
	getModelParams.HTTPRequest = params.HTTPRequest
	getModelParams.ModelID = params.ModelID
	getModelParams.Version = params.Version
	logr := logger.LocLogger(logWithRepeatModelParams(params))
	//logr.Debugf("retryModel invoked: %v", params.HTTPRequest.Header)
	rresp, errFromTrainer := getTrainingJobFromStorage(getModelParams)
	if rresp == nil {
		return errFromTrainer
	}

	// FIXME MLSS Change: more properties for result
	m := createModel(params.HTTPRequest, rresp.Job, logr)

	//logr.Debugf("m: %+v", m)
	manifest, err = trainingToManifest(m, nil, params.HTTPRequest, logr, "")
	if err != nil {
		return handleErrorResponse(logr, err.Error())
	}
	postModelParams.HTTPRequest = params.HTTPRequest
	return postModelInclude(postModelParams, manifest, operation)
}

func ListModels(params models.ListModelsParams) middleware.Responder {
	//1.Recv Parmas Init
	logr := logger.LocLogger(logWithGetListModelsParams(params))
	userID := getUserID(params.HTTPRequest)
	isSA := getSuperadmin(params.HTTPRequest)

	queryUser, namespace, page, size, clusterName, expRunId, err := getParamsForListModels(params)
	if err != nil {
		return handleErrorResponse(logr, err.Error())
	}
	//2.Params Check & Perminssion Check
	err = checkPageInfo(page, size)
	if err != nil {
		return handleErrorResponse(logr, err.Error())
	}
	//3.Variable Init
	//Init Model Storage Client
	sClient, err := storageClient.NewStorage()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for storage service")
		return handleErrorResponse(logr, "")
	}
	defer sClient.Close()

	//4.Get Model From Storage Client
	var resp *grpc_storage.GetAllResponse = nil
	logr.Debugf("get clusterName: %v", clusterName)
	var role string
	//if not sa
	if isSA == SA {
		role = userSA
		logr.Debugf("pass for SA.")
	} else {
		db := datasource.GetDB()
		// Permission Check
		if expRunId != "" {
			expRunIdInt, err := strconv.ParseInt(expRunId, 10, 64)
			if err != nil {
				logr.Error("Parse expId error, ", err.Error())
				return handleErrorResponse(logr, err.Error())
			}

			expRun, err := repo.ExperimentRunRepo.Get(expRunIdInt, db)
			if err != nil {
				logr.Error("Get ExperimentRun error, ", err.Error())
				return handleErrorResponse(logr, err.Error())

			}
			if *expRun.User.Name != userID {
				return httpResponseHandle(http.StatusUnauthorized, PermissionDeniedError, "List Models",
					[]byte("permission denied"))
			}
		}

		role = userGA
		if namespace != "" && queryUser != "" {
			role, err = getRoleNum(params.HTTPRequest.Header.Get(cc.CcAuthToken), getUserID(params.HTTPRequest), namespace)
			if err != nil {
				logr.Error("getRoleNum error, ", err)
				return handleErrorResponse(logr, err.Error())
			}
		}
	}
	nsList, err := getUserNSList(params.HTTPRequest.Header.Get(cc.CcAuthToken), role)
	if err != nil {
		logger.GetLogger().Info("getUserNSList error, ", err)
		return handleErrorResponse(logr, err.Error())
	}
	if queryUser == "" && namespace == "" {
		if isSA == SA {
			// FIXME MLSS Change: v_1.5.0 add logic when get nsList from cc
			if nsList != nil && len(nsList) > 0 {
				resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespaceList(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
					NamespaceList: nsList,
					Page:          page,
					Size:          size,
					ExpRunId:      expRunId,
				})
				if err != nil {
					logger.GetLogger().Error("GetAllTrainingsJobsByUserIdAndNamespaceList error, ", err.Error())
					return handleErrorResponse(logr, err.Error())
				}
			} else {
				resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespaceList(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
					NamespaceList: nsList,
					Page:          page,
					Size:          size,
					ExpRunId:      expRunId,
				})
			}
			if err != nil {
				logger.GetLogger().Error("GetAllTrainingsJobsByUserIdAndNamespaceList error, ", err.Error())
				return handleErrorResponse(logr, err.Error())
			}
		} else {
			logr.Debugf("pass nothing for GA/GU.")
			//other role
			//nsList,user
			//CC:getNamespaceWithRole
			logr.Debugf("debug_metaUserID: %v", userID)
			resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespace(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
				UserId:      getUserID(params.HTTPRequest),
				Username:    getUserID(params.HTTPRequest),
				ClusterName: clusterName,
				Page:        page,
				Size:        size,
				ExpRunId:    expRunId,
			})
			if err != nil {
				logger.GetLogger().Error("GetAllTrainingsJobsByUserIdAndNamespaceList error, ", err)
				return handleErrorResponse(logr, err.Error())
			}
			logr.Debugf("storage.Client().GetAllTrainingsJobsByUserIdAndNamespaceList, page: %v, total %v", resp.Pages, resp.Total)
		}
	} else if queryUser == "" && namespace != "" {
		//pass namespace only
		//CC:/access/admin/namespaces/{namespace}
		if isSA == SA {
			resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespace(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
				Namespace: namespace,
				Page:      page,
				Size:      size,
				ExpRunId:  expRunId,
			})
		} else {
			resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespace(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
				UserId:    getUserID(params.HTTPRequest),
				Username:  queryUser,
				Namespace: namespace,
				Page:      page,
				Size:      size,
				ExpRunId:  expRunId,
			})
		}
	} else if queryUser != "" && namespace != "" {
		//pass user and namespace
		//CC:/access/admin/namespaces/{namespace}/users/{user}
		if isSA == SA {
			resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespace(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
				Username:      queryUser,
				ClusterName:   clusterName,
				NamespaceList: nsList,
				Page:          page,
				Size:          size,
				ExpRunId:      expRunId,
			})
		} else {
			resp, err = sClient.Client().GetAllTrainingsJobsByUserIdAndNamespace(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
				UserId:        getUserID(params.HTTPRequest),
				Username:      queryUser,
				ClusterName:   clusterName,
				NamespaceList: nsList,
				Page:          page,
				Size:          size,
				ExpRunId:      expRunId,
			})
		}
	} else {
		return handleErrorResponse(logr, "namespace can't be null, when user exists")
	}

	//get trainingJob error handle
	if err != nil {
		logr.WithError(err).Error("Trainer readAll service call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		return handleErrorResponse(logr, err.Error())
	}

	// Build Model List Response
	pageInt := 0
	totalInt := 0
	marr := make([]*restmodels.Model, 0, len(resp.Jobs))

	if resp != nil {
		for _, job := range resp.Jobs {
			// use append(); we may have skipped some because they were gone by the time we got to them.
			// FIXME MLSS Change: more properties for result
			marr = append(marr, createModel(params.HTTPRequest, job, logr))
		}
		pageInt, err = strconv.Atoi(resp.Pages)
		if err != nil {
			logr.Error("resp.Pages parse int failed")
		}
		totalInt, err = strconv.Atoi(resp.Total)
		if err != nil {
			logr.Error("resp.Total parse int failed")
		}
	}
	return models.NewListModelsOK().WithPayload(&restmodels.ModelList{
		Models: marr,
		Pages:  int64(pageInt),
		Total:  int64(totalInt),
	})
}

func PatchModel(params models.PatchModelParams) middleware.Responder {
	logr := logger.LocLogger(logWithUpdateStatusParams(params))
	logr.Debugf("patchModel invoked: %v", params.HTTPRequest.Header)

	if params.Payload.Status != "halt" {
		//FIXME MLSS Change: change Description to Msg
		return models.NewPatchModelBadRequest().WithPayload(&restmodels.Error{
			Error: "Bad request",
			Code:  http.StatusBadRequest,
			Msg:   "status parameter has incorrect value",
		})
	}

	trainer, err := trainerClient.NewTrainer()

	if err != nil {
		logr.Errorf("Cannot create client for trainer service: %s", err.Error())
		return handleErrorResponse(logr, "")
	}
	defer trainer.Close()

	if err != nil {
		logr.Errorf("Trainer status update service call failed: %s", err.Error())
		if grpc.Code(err) == codes.NotFound {
			//FIXME MLSS Change: change Description to Msg
			return models.NewPatchModelNotFound().WithPayload(&restmodels.Error{
				Error: "Not found",
				Code:  http.StatusNotFound,
				Msg:   "model_id not found",
			})
		}
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
	}
	return models.NewPatchModelAccepted().WithPayload(&restmodels.BasicModel{
		ModelID: params.ModelID,
	})
}

func getUserNSList(token string, role string) ([]string, error) {
	ccAddress := viper.GetString(config.CCAddress)
	ccClient := cc.GetCcClient(ccAddress)
	bodyForNS, err := ccClient.GetCurrentUserNamespaceWithRole(token, role)
	if err != nil {
		logger.GetLogger().Error("Error when checking namespace from cc, ", err)
		return []string{}, errors.New("Error when checking namespace from cc, " + err.Error())
	}

	logger.GetLogger().Debugf("ns from cc: %v", string(bodyForNS))
	var nsVOList []commonModels.NamespaceVO
	err = commonModels.GetResultData(bodyForNS, &nsVOList)
	if err != nil {
		logger.GetLogger().Errorf("GetResultData failed")
		return []string{}, errors.New("GetResultData failed")
	}
	
	nsList := []string{}
	for _, v := range nsVOList {
		nsList = append(nsList, v.Namespace)
	}
	return nsList, nil
}

//get role's number, GA is 1, GU is 2
func getRoleNum(token string, username string, namespace string) (string, error) {
	var role string
	var roleNum string
	ccAddress := viper.GetString(config.CCAddress)
	ccClient := cc.GetCcClient(ccAddress)
	err := ccClient.CheckNamespaceUser(token, namespace, username)
	if err != nil {
		logger.GetLogger().Errorf("CheckNamespaceUser failed, %v", err.Error())
		return "", err
	}
	bodyForRole, err := ccClient.CheckNamespace(token, namespace)
	err = commonModels.GetResultData(bodyForRole, &role)
	if err != nil {
		logger.GetLogger().Errorf("GetResultData failed")
		return "", err
	}
	if role == "SA" || role == "GA" {
		roleNum = "1"
	} else {
		roleNum = "2"
	}
	return roleNum, nil
}

func handleTfosManifest(userId string, manifest *ManifestV1, logr *logger.LocLoggingEntry) error {
	logr.Debugf("handleTfosManifest start")
	//fixme: skip loading model in learner
	manifest.CodeSelector = "storagePath"
	cpuInfloat64, err := strconv.ParseFloat(os.Getenv("LINKIS_EXECUTOR_CPU"), 64)
	if err != nil {
		logr.Debugf("ParseFloat LINKIS_EXECUTOR_CPU failed: %v", err.Error())
		return err
	}
	manifest.Cpus = cpuInfloat64
	manifest.Memory = os.Getenv("LINKIS_EXECUTOR_MEM")
	manifest.Gpus = 0

	clusterNameForDir, err := getClusterNameForDir(manifest.Namespace)
	if err != nil {
		logr.Debugf("handleTfosManifest, getClusterNameForDir: %v", err.Error())
		return err
	}

	var dss = make([]*dataStoreRef, 1)
	manifest.DataStores = dss
	var ds = &dataStoreRef{
		TrainingData: &storageContainerV1{
			Container: "data",
		},
		TrainingResults: &storageContainerV1{
			Container: "result",
		},
		ID:   "hostmount",
		Type: "mount_volume",
		Connection: map[string]string{
			"type": "host_mount",
			"name": "host-mount",
			"path": fmt.Sprintf("/data/%v/mlss-data/%v", clusterNameForDir, userId),
		},
	}
	dss[0] = ds
	manifest.Framework = &frameworkV1{
		Name:    os.Getenv("LINKIS_EXECUTOR_IMAGE"),
		Version: os.Getenv("LINKIS_EXECUTOR_TAG"),
	}
	manifest.Framework.Command = "/main"
	logr.Debugf("after handleTfosManifest: %+v", manifest)
	return nil
}

func getClusterNameForDir(namespace string) (string, error) {
	split := stringsUtil.Split(namespace, "-")

	if len(split) != 7 {
		return "", errors.New("ns must have 7 Field")
	}

	if split[2] == constants.BDP {
		return constants.BDP_SS, nil
	} else if split[2] == constants.BDAP {
		if split[6] == constants.SAFE {
			return constants.BDAPSAFE_SS, nil
		} else {
			return constants.BDAP_SS, nil
		}
	} else {
		return "", errors.New("check ns pls,there is only 3 clusterName bdapsafe/bdap/bdp now")
	}
}

func ExportModel(params models.ExportModelParams) middleware.Responder {
	currentUserId := getUserID(params.HTTPRequest)
	isSA := getSuperadminBool(params.HTTPRequest)

	//1、Recv Parmas Init
	folderUid := uuid.New().String()
	manifestPath := "./exportFolder/" + folderUid + "/model"
	manifestDefYml := "manifest.yaml"
	yamlFile, err := createYamlFile(manifestPath, manifestDefYml)
	if err != nil {
		logger.GetLogger().Error("CreateYamlFile error, ", err)
		return handleErrorResponse(logger.GetLogger(), err.Error())
	}
	storage, err := storageClient.NewStorage()
	defer storage.Close()
	if err != nil {
		logger.GetLogger().Error("Cannot create client for trainer service")
		return handleErrorResponse(logger.GetLogger(), "")
	}
	//2、Get TrainingJob
	trainingJob, manifest, err := getTrainingJob(params, storage)
	if err != nil {
		logger.GetLogger().Error("Get training job err, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	if !isSA && trainingJob.Job.UserId != currentUserId {
		return getErrorResultWithPayload("forbidden", http.StatusForbidden,
			errors.New("permission deniedKillTrainingJob"))
	}

	//3、Write yaml file
	yamlBytes, err := yaml.Marshal(&manifest)
	if err != nil {
		logger.GetLogger().Error("Yaml marshal faild err, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	yamlFile.Write(yamlBytes)
	err = yamlFile.Close()
	if err != nil {
		logger.GetLogger().Error("Yaml file close err, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	//4、Check codeFile or storagePath
	logger.GetLogger().Debugf("ExportModel trainingJob.Job.FilePath: %v\n", trainingJob.Job.FilePath)
	if trainingJob.Job.FilePath != "" {
		s3Client, err := storageClient.NewStorage()
		todo := context.TODO()
		ctx, _ := context.WithCancel(todo)
		downloadReq := grpc_storage.CodeDownloadRequest{
			S3Path: trainingJob.Job.FilePath,
		}
		res, err := s3Client.Client().DownloadCode(ctx, &downloadReq)
		if err != nil {
			logger.GetLogger().Error("s3Client download code failed, ", err)
			return handleErrorResponse(logger.GetLogger(), "")
		}
		//5、Write to code zip file
		err = writeExportCodeZipFile(manifestPath, trainingJob, res)
		if err != nil {
			logger.GetLogger().Error("Write export code zip file error, ", err)
			return handleErrorResponse(logger.GetLogger(), err.Error())
		}
	}
	//6、Archive and export
	uid := uuid.New().String()
	zipName := fmt.Sprintf("%v.zip", uid)
	err = archiver.Archive([]string{manifestPath}, zipName)
	if err != nil {
		logger.GetLogger().Error("zip archiver failed, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	//read code zip file's binary
	zipFileBytes, err := ioutil.ReadFile(fmt.Sprintf("./%v", zipName))
	if err != nil {
		logger.GetLogger().Error("ReadFile failed, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	//clear manifestPath's all file
	err = os.RemoveAll(manifestPath)
	if err != nil {
		logger.GetLogger().Error("clear folder failed, ", err)
		return handleErrorResponse(logger.GetLogger(), "")
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.Write(zipFileBytes)
		w.WriteHeader(200)
	})
}

func GetDashboards(params operations.GetDashboardsParams) middleware.Responder {
	//1、Recv Parmas Init
	logr := logger.LocLogger(logWithGetDashboards(params))
	var currentUserId = getUserID(params.HTTPRequest)
	var isSA = getSuperadmin(params.HTTPRequest)
	var resp *grpc_storage.GetAllResponse = nil
	var role string
	logr.Debugf(">>getDashboards<< metaUserID: %v, superadmin: %v", currentUserId, isSA)
	storage, err := storageClient.NewStorage()
	//defer func() {
	//	err := storage.Close()
	//	if err != nil {
	//		logr.WithError(err).Errorf("Cannot create client for trainer service")
	//	}else{
	//		logr.Debugf("Storage Client is close")
	//	}
	//}()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		return handleErrorResponse(logr, "")
	}

	//2、Get userId and namespace from params
	logr.Debugf("Calling trainer.Client().GetAllTrainingsJobs(...)")
	logr.Debugf("debug")

	//3、Check role
	if isSA == SA {
		role = userSA
	} else {
		role = userGA
	}
	nsList, err := getUserNSList(params.HTTPRequest.Header.Get(cc.CcAuthToken), role)
	if err != nil {
		logr.WithError(err).Error("getUserNSList error, ", err)
		return handleErrorResponse(logr, err.Error())
	}

	//4、Get nsList by role
	jobTotal := 0
	jobRunning := 0
	var gpuCount float32
	var currentPage int64 = 1
	var pageSize int64 = 10
	queryFlag := true
	pageSizeStr := strconv.FormatInt(pageSize, 10)
	if isSA == SA {
		logr.Infof("pass nothing for SA.")
		//if role == SA
		//nil,nil
		// FIXME MLSS Change: v_1.5.0 added param clusterName to filter models
		if nsList != nil && len(nsList) > 0 {
			for queryFlag {
				pageStr := strconv.FormatInt(currentPage, 10)
				resp, err = storage.Client().GetAllTrainingsJobsByUserIdAndNamespaceList(
					params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{
						NamespaceList: nsList,
						Page:          pageStr,
						Size:          pageSizeStr,
					})
				if err != nil {
					logr.WithError(err).Errorf("GetAllTrainingsJobsByUserIdAndNamespaceList failed")
					return handleErrorResponse(logr, err.Error())
				}
				totalInt64, err := strconv.ParseInt(resp.Total, 10, 64)
				if err != nil {
					logr.WithError(err).Errorf("GetAllTrainingsJobsByUserIdAndNamespaceList failed to"+
						" parse response total, total: %s, err: %v", resp.Total, err)
					return handleErrorResponse(logr, err.Error())
				}
				if (currentPage * pageSize) < totalInt64 {
					currentPage += 1
				} else {
					queryFlag = false
				}

				//5 make response
				if resp != nil {
					jobTotal += len(resp.Jobs)
					if len(resp.Jobs) > 0 {
						for _, j := range resp.Jobs {
							if grpc_storage.Status_COMPLETED != j.Status.Status && grpc_storage.Status_FAILED != j.Status.Status &&
								grpc_storage.Status_PENDING != j.Status.Status && grpc_storage.Status_QUEUED != j.Status.Status {
								jobRunning += 1
								if j.GetTraining().GetResources().GetGpus() > 0 {
									gpu := j.GetTraining().GetResources().GetGpus()
									gpuCount += gpu
								}
							}
						}
					}
				}
			}
		}
	} else {
		logr.Debugf("pass nothing for GA/GU.")
		if nsList != nil && len(nsList) > 0 {
			for queryFlag {
				pageStr := strconv.FormatInt(currentPage, 10)
				//resp, err = storage.Client().GetAllTrainingsJobsByUserIdAndNamespaceList(params.HTTPRequest.Context(), &grpc_storage.GetAllRequest{ //todo
				resp, err = storage.Client().GetAllTrainingsJobsByUserIdAndNamespaceList(context.TODO(), &grpc_storage.GetAllRequest{
					Username:      currentUserId,
					NamespaceList: nsList,
					Page:          pageStr,
					Size:          pageSizeStr,
				})
				if err != nil {
					logr.WithError(err).Errorf("GetAllTrainingsJobsByUserIdAndNamespaceList failed")
					return handleErrorResponse(logr, err.Error())
				}
				totalInt64, err := strconv.ParseInt(resp.Total, 10, 64)
				if err != nil {
					logr.WithError(err).Errorf("GetAllTrainingsJobsByUserIdAndNamespaceList failed to"+
						" parse response total, total: %s, err: %v", resp.Total, err)
					return handleErrorResponse(logr, err.Error())
				}
				if (currentPage * pageSize) < totalInt64 {
					currentPage += 1
				} else {
					queryFlag = false
				}

				//5 make response
				if resp != nil {
					jobTotal += len(resp.Jobs)
					if len(resp.Jobs) > 0 {
						for _, j := range resp.Jobs {
							if grpc_storage.Status_COMPLETED != j.Status.Status && grpc_storage.Status_FAILED != j.Status.Status &&
								grpc_storage.Status_PENDING != j.Status.Status && grpc_storage.Status_QUEUED != j.Status.Status {
								jobRunning += 1
								if j.GetTraining().GetResources().GetGpus() > 0 {
									gpu := j.GetTraining().GetResources().GetGpus()
									gpuCount += gpu
								}
							}
						}
					}
				}
			}
		}
	}

	//5、Make response
	mp := map[string]int{"jobTotal": jobTotal, "jobRunning": jobRunning, "gpuCount": int(gpuCount)}
	marshal, err := json.Marshal(mp)
	if err != nil {
		logr.WithError(err).Errorf("json.Marshal(mp) failed")
		return handleErrorResponse(logr, err.Error())
	}
	var result = commonModels.Result{
		Code:   "200",
		Msg:    "success",
		Result: json.RawMessage(marshal),
	}

	//6.Close Storage Client
	err = storage.Close()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		w.WriteHeader(200)
	})
}

func GetLogs(params models.GetLogsParams) middleware.Responder {
	return getLogsOrMetrics(params, false)
}

func DownloadModelDefinition(params models.DownloadModelDefinitionParams) middleware.Responder {
	logr := logger.LocLogger(logWithDownloadModelDefinitionParams(params))
	logr.Debugf("downloadModelDefinition invoked: %v", params.HTTPRequest.Header)

	//trainer, err := trainerClient.NewTrainer()
	storage, err := storageClient.NewStorage()
	defer storage.Close()
	if err != nil {
		logger.GetLogger().Error("Cannot create client for trainer service")
		return handleErrorResponse(logr, "")
	}

	stream, err := storage.Client().GetModelDefinition(params.HTTPRequest.Context(), &grpc_storage.ModelDefinitionRequest{
		TrainingId: params.ModelID,
		UserId:     getUserID(params.HTTPRequest),
	})
	if err != nil {
		logr.WithError(err).Error("Trainer GetModelDefinition service call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		if grpc.Code(err) == codes.NotFound {
			//FIXME MLSS Change: change Description to Msg
			return models.NewDownloadModelDefinitionNotFound().WithPayload(&restmodels.Error{
				Error: "Not found",
				Code:  http.StatusNotFound,
				Msg:   "",
			})
		}
		return handleErrorResponse(logr, "")
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				w.WriteHeader(200)
				if err = stream.CloseSend(); err != nil {
					logr.WithError(err).Error("Closing stream failed")
				}
				break
			} else if err != nil {
				logr.WithError(err).Errorf("Cannot read model definition")
				// this error handling is a bit of a hack with overriding the content-type
				// but the swagger-generated REST client chokes if we leave the content-type
				// as application-octet stream
				w.WriteHeader(500)
				w.Header().Set(runtime.HeaderContentType, runtime.JSONMime)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(&restmodels.Error{
					Error: "Internal server error",
					Msg:   "",
					Code:  500,
				})
				w.Write(payload)
				break
			}
			w.Write(chunk.Data)
		}
	})
}

func DownloadTrainedModel(params models.DownloadTrainedModelParams) middleware.Responder {
	logr := logger.LocLogger(logWithDownloadTrainedModelParams(params))
	logr.Debugf("downloadTrainedModel invoked: %v", params.HTTPRequest.Header)

	//trainer, err := trainerClient.NewTrainer()
	storage, err := storageClient.NewStorage()
	defer storage.Close()
	if err != nil {
		logger.GetLogger().Error("Cannot create client for trainer service")
		return handleErrorResponse(logr, "")
	}

	stream, err := storage.Client().GetTrainedModel(params.HTTPRequest.Context(), &grpc_storage.TrainedModelRequest{
		TrainingId: params.ModelID,
		UserId:     getUserID(params.HTTPRequest),
	})
	if err != nil {
		logr.WithError(err).Errorf("Trainer GetTrainedModel service call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		if grpc.Code(err) == codes.NotFound {
			//FIXME MLSS Change: change Description to Msg
			return models.NewDownloadTrainedModelNotFound().WithPayload(&restmodels.Error{
				Error: "Not found",
				Code:  http.StatusNotFound,
				Msg:   "",
			})
		}
		return handleErrorResponse(logr, "")
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				w.WriteHeader(200)
				if err = stream.CloseSend(); err != nil {
					logr.WithError(err).Error("Closing stream failed")
				}
				break
			} else if err != nil {
				logr.WithError(err).Error("Cannot read trained model")
				// this error handling is a bit of a hack with overriding the content-type
				// but the swagger-generated REST client chokes if we leave the content-type
				// as application-octet stream
				w.WriteHeader(500)
				w.Header().Set(runtime.HeaderContentType, runtime.JSONMime)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(&restmodels.Error{
					Error: "Internal server error",
					Msg:   "",
					Code:  500,
				})
				w.Write(payload)
				break
			}
			w.Write(chunk.Data)
		}
	})
}

func GetLoglines(params training_data.GetLoglinesParams) middleware.Responder {
	logr := logger.LocLogger(logWithLoglinesParams(params))
	logr.Debug("function entry")

	trainingData, err := trainingDataClient.NewTrainingDataClient()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		return handleErrorResponse(logr, "")
	}
	defer trainingData.Close()

	// FIXME MLSS Change: get user id of the training from trainer
	//getModelParams := models.NewGetModelParams().WithModelID(params.ModelID).WithTimeout(10*time.Second)
	getModelParams := models.GetModelParams{
		HTTPRequest: params.HTTPRequest,
		ModelID:     params.ModelID,
	}
	//rresp, errFromTrainer := getTrainingJobFromTrainer(getModelParams)
	rresp, errFromTrainer := getTrainingJobFromStorage(getModelParams)

	if rresp == nil {
		return errFromTrainer
	}
	trainingUserID := rresp.Job.UserId

	var metaUserID = getUserID(params.HTTPRequest)

	// FIXME MLSS Change: check if current request user can access the records of the training creator
	//ccClient := cc.GetCcClient("http://controlcenter:7777")
	// FIXME MLSS Change: read cc address from configmap
	ccAddress := viper.GetString(config.CCAddress)
	ccClient := cc.GetCcClient(ccAddress)
	err = ccClient.AdminUserCheck(params.HTTPRequest.Header.Get(cc.CcAuthToken), metaUserID, trainingUserID)
	if err != nil {
		logr.WithError(err).Errorf("Error when checking user access from cc")
		return handleErrorResponse(logr, "")
	}

	var sinceQuery = ""
	if params.SinceTime != nil {
		sinceQuery = *params.SinceTime
	}
	var pagesize int32 = defaultLogPageSize
	if params.Pagesize != nil {
		pagesize = *params.Pagesize
	}
	var pos int64
	if params.Pos != nil {
		pos = *params.Pos
	}
	var searchType grpc_training_data_v1.Query_SearchType = grpc_training_data_v1.Query_TERM
	if params.SearchType != nil {
		searchType = makeGrpcSearchTypeFromRestSearchType(*params.SearchType)
	}

	query := &grpc_training_data_v1.Query{
		Meta: &grpc_training_data_v1.MetaInfo{
			TrainingId: params.ModelID,
			UserId:     metaUserID,
			Time:       0,
		},
		Pagesize:   pagesize,
		Pos:        pos,
		Since:      sinceQuery,
		SearchType: searchType,
	}

	// The marshal from the grpc record to the rest record should probably be just a byte stream transfer.
	// But for now, I prefer a structural copy, I guess.

	getLogsClient, err := trainingData.Client().GetLogs(params.HTTPRequest.Context(), query)
	if err != nil {
		trainingData.Close()
		logr.WithError(err).Error("GetLogs call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				//FIXME MLSS Change: change Description to Msg
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		return handleErrorResponse(logr, "")
	}

	marr := make([]*restmodels.V1LogLine, 0, pagesize)

	err = nil
	nRecordsActual := 0
	for ; nRecordsActual < int(pagesize); nRecordsActual++ {
		logsRecord, err := getLogsClient.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			logr.WithError(err).Errorf("Cannot read model definition")
			break
		}

		marr = append(marr, &restmodels.V1LogLine{
			Meta: &restmodels.V1MetaInfo{
				TrainingID: logsRecord.Meta.TrainingId,
				UserID:     logsRecord.Meta.UserId,
				Time:       logsRecord.Meta.Time,
				Rindex:     logsRecord.Meta.Rindex,
			},
			Line: logsRecord.Line,
		})
	}
	trimmedList := marr[0:nRecordsActual]

	response := training_data.NewGetLoglinesOK().WithPayload(&restmodels.V1LogLinesList{
		Models: trimmedList,
	})
	logr.Debug("function exit")
	return response
}

func createYamlFile(manifestPath string, manifestDefYml string) (*os.File, error) {
	//create yaml's dir
	err := os.MkdirAll(manifestPath, os.ModePerm)
	if err != nil {
		logger.GetLogger().Info("create flowPath failed:", err)
		return nil, err
	}
	//create yaml's file
	yamlFile, err := os.Create(fmt.Sprintf("%v/%v", manifestPath, manifestDefYml))
	if err != nil {
		logger.GetLogger().Info("yamlFile create failed, ", err)
		return nil, err
	}
	return yamlFile, err
}

func getTrainingJob(params models.ExportModelParams, storage storageClient.StorageClient) (*grpc_storage.GetResponse, *ManifestV1, error) {
	var tfReq v1alpha1.TFosRequest
	var framework frameworkV1
	var jobAlert map[string][]map[string]string
	dataStores := []*dataStoreRef{}
	req := grpc_storage.GetTrainingJobRequest{
		TrainingId: params.ModelID,
	}
	//get training job
	trainingJob, err := storage.Client().GetTrainingJobByTrainingJobId(context.TODO(), &req)
	if err != nil {
		logger.GetLogger().Error("GetTrainingJobByTrainingJobId faild err, ", err)
		return nil, nil, err
	}
	expRunId, _ := strconv.Atoi(trainingJob.Job.ExpRunId)
	copier.Copy(&tfReq, &trainingJob.Job.TfosRequest)
	copier.Copy(&framework, &trainingJob.Job.ModelDefinition.Framework)
	json.Unmarshal([]byte(trainingJob.Job.JobAlert), &jobAlert)
	framework.Command = trainingJob.Job.Training.Command
	//make data store
	dataStore := &dataStoreRef{
		ID:         "hostmount ",
		Type:       "mount_volume",
		Connection: trainingJob.Job.Datastores[0].Connection,
		TrainingData: &storageContainerV1{
			Container: trainingJob.Job.Datastores[0].Fields["bucket"],
		},
		TrainingResults: &storageContainerV1{
			Container: trainingJob.Job.Datastores[1].Fields["bucket"],
		},
	}

	mfModel := MFModel{}
	if trainingJob.Job.MfModel != nil {
		mfModel = MFModel{
			GroupId:     trainingJob.Job.MfModel.GroupId,
			ModelName:   trainingJob.Job.MfModel.ModelName,
			Version:     trainingJob.Job.MfModel.Version,
			FactoryName: trainingJob.Job.MfModel.FactoryName,
		}
	}

	//apped data store to training job's data stores
	dataStores = append(dataStores, dataStore)
	return trainingJob, &ManifestV1{
		Name:         trainingJob.Job.ModelDefinition.Name,
		Description:  trainingJob.Job.ModelDefinition.Description,
		Version:      "1.0",
		Cpus:         float64(trainingJob.Job.Training.Resources.Cpus),
		Gpus:         float64(trainingJob.Job.Training.Resources.Gpus),
		Memory:       fmt.Sprintf("%vGb", trainingJob.Job.Training.Resources.Memory),
		GpuType:      trainingJob.Job.Training.Resources.GpuType,
		Learners:     trainingJob.Job.Training.Resources.Learners,
		Storage:      fmt.Sprintf("%v", trainingJob.Job.Training.Resources.Storage),
		Namespace:    trainingJob.Job.JobNamespace,
		CodeSelector: trainingJob.Job.CodeSelector,
		JobType:      trainingJob.Job.JobType,
		JobAlert:     jobAlert,
		ExpName:      trainingJob.Job.ExpName,
		ExpRunId:     int64(expRunId),
		CodePath:     "",
		FileName:     "",
		PSs:          trainingJob.Job.Pss,
		PSCPU:        trainingJob.Job.PsCpu,
		PSImage:      trainingJob.Job.PsImage,
		PSMemory:     trainingJob.Job.PsMemory,
		TFosRequest:  &tfReq,
		DataStores:   dataStores,
		Framework:    &framework,
		MFModel:      &mfModel,
		ProxyUser:    trainingJob.Job.ProxyUser,
		Algorithm:    trainingJob.Job.Algorithm,
		FitParams:    trainingJob.Job.FitParams,
		APIType:      trainingJob.Job.APIType,
	}, nil
}

func writeExportCodeZipFile(manifestPath string, trainingJob *grpc_storage.GetResponse, codeDownloadRes *grpc_storage.CodeDownloadResponse) error {
	codeZipPath := fmt.Sprintf("%v/%v", manifestPath, trainingJob.Job.FileName)
	///read oss's code zip file
	codeFileBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", codeDownloadRes.HostPath, codeDownloadRes.FileName))
	if err != nil {
		logger.GetLogger().Error("codeFile read failed, ", err)
		return err
	}
	//create code zip file
	codeZipFile, err := os.Create(codeZipPath)
	if err != nil {
		logger.GetLogger().Error("codeZipFile create failed, ", err)
		return err
	}
	//write to code zip file
	codeZipFile.Write(codeFileBytes)
	codeZipFile.Close()
	return nil
}

func defineJobPath(DSs []*dataStoreRef, manifest *ManifestV1) (string, string) {
	var dataPath string
	var resultsPath string
	for _, ds := range DSs {
		conn := ds.Connection
		parentPath := conn["path"]
		// Add logic for path when create notebook
		if manifest.JobType != "tfos" {
			subDataPath := ds.TrainingData.Container
			if !stringsUtil.HasPrefix(subDataPath, "/") && !stringsUtil.HasSuffix(parentPath, "/") {
				dataPath = parentPath + "/" + subDataPath
			} else if stringsUtil.HasPrefix(subDataPath, "/") && stringsUtil.HasSuffix(parentPath, "/") {
				dataPath = parentPath + subDataPath[1:]
			} else {
				dataPath = parentPath + subDataPath
			}
		}

		subResultsPath := ds.TrainingResults.Container
		if !stringsUtil.HasPrefix(subResultsPath, "/") && !stringsUtil.HasSuffix(parentPath, "/") {
			resultsPath = parentPath + "/" + subResultsPath
		} else if stringsUtil.HasPrefix(subResultsPath, "/") && stringsUtil.HasSuffix(parentPath, "/") {
			resultsPath = parentPath + subResultsPath[1:]
		} else {
			resultsPath = parentPath + subResultsPath
		}
		break
	}
	return dataPath, resultsPath
}

func codePathToModelDefinition(codePath string) (io.ReadCloser, error) {
	sClient, err := storageClient.NewStorage()
	defer sClient.Close()
	if err != nil {
		logger.GetLogger().Error("NewStorage Error: ", err.Error())
		return nil, err
	}
	req := &grpc_storage.CodeDownloadRequest{S3Path: codePath}
	ctx := context.TODO()
	res, err := sClient.Client().DownloadCode(ctx, req)
	if err != nil {
		logger.GetLogger().Error("Storage Download Code Error: ", err.Error())
		return nil, err
	}
	//fileByte, err := ioutil.ReadFile(res.HostPath+"/"+res.FileName)

	file, err := os.Open(res.HostPath + "/" + res.FileName)
	fileReader := bufio.NewReader(file)

	return ioutil.NopCloser(fileReader), err
}

// func handleTfosManifest(userId string, manifest *ManifestV1, logr *logger.LocLoggingEntry) error {
// 	logr.Debugf("handleTfosManifest start")

// 	//fixme: skip loading model in learner
// 	manifest.CodeSelector = "storagePath"

// 	cpuInfloat64, err := strconv.ParseFloat(os.Getenv("LINKIS_EXECUTOR_CPU"), 64)
// 	if err != nil {
// 		logr.Debugf("ParseFloat LINKIS_EXECUTOR_CPU failed: %v", err.Error())
// 		return err
// 	}
// 	manifest.Cpus = cpuInfloat64
// 	manifest.Memory = os.Getenv("LINKIS_EXECUTOR_MEM")
// 	manifest.Gpus = 0

// 	clusterNameForDir, err := getClusterNameForDir(manifest.Namespace)
// 	if err != nil {
// 		logr.Debugf("handleTfosManifest, getClusterNameForDir: %v", err.Error())
// 		return err
// 	}

// 	var dss = make([]*dataStoreRef, 1)
// 	manifest.DataStores = dss
// 	var ds = &dataStoreRef{
// 		TrainingData: &storageContainerV1{
// 			Container: "data",
// 		},
// 		TrainingResults: &storageContainerV1{
// 			Container: "result",
// 		},
// 		ID:   "hostmount",
// 		Type: "mount_volume",
// 		Connection: map[string]string{
// 			"type": "host_mount",
// 			"name": "host-mount",
// 			"path": fmt.Sprintf("/data/%v/mlss-data/%v", clusterNameForDir, userId),
// 		},
// 	}
// 	dss[0] = ds

// 	manifest.Framework = &frameworkV1{
// 		Name:    os.Getenv("LINKIS_EXECUTOR_IMAGE"),
// 		Version: os.Getenv("LINKIS_EXECUTOR_TAG"),
// 	}

// 	manifest.Framework.Command = "/main"

// 	logr.Debugf("after handleTfosManifest: %+v", manifest)

// 	return nil
// }

// FIXME MLSS Change: check manifest.yaml
func checkManifest(manifest *ManifestV1, logr *logger.LocLoggingEntry) error {
	if math.IsNaN(manifest.Cpus) {
		return errors.New("cpus is NaN")
	}
	if math.IsNaN(manifest.Gpus) {
		return errors.New("gpus is NaN")
	}
	logr.Debugf("checkManifest for manifest.JobType: %v, manifest.Learners: %v, cpu: %v, gpu: %v, memory: %v", manifest.JobType, manifest.Learners, manifest.Cpus, manifest.Gpus, manifest.Memory)
	//FIXME MLSS Change: v_1.5.2 added
	if "dist-tf" == manifest.JobType {
		if math.IsNaN(float64(manifest.Learners)) {
			return errors.New("learners is NaN")
		}
		if manifest.Learners <= 0 {
			return errors.New("learners must be integer greater than 0")
		}
	} else {
		if manifest.Learners > 1 {
			return errors.New("if jobType not dist-tf learners can't greater than 1")
		}
	}

	memory := manifest.Memory
	lower := stringsUtil.ToLower(memory)
	if stringsUtil.HasSuffix(lower, "mib") {
		if _, e := strconv.ParseFloat(stringsUtil.TrimSuffix(lower, "mib"), 64); e != nil {
			return errors.New("memory not a number")
		}
	} else if stringsUtil.HasSuffix(lower, "gb") {
		if _, e := strconv.ParseFloat(stringsUtil.TrimSuffix(lower, "gb"), 64); e != nil {
			return errors.New("memory not a number")
		}
	} else {
		return errors.New("now only support unit Mib, GB. current memory: " + manifest.Memory)
	}

	return nil
}

func checkPageInfo(page string, size string) error {
	if page != "" || size != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return errors.New("page parse to int failed")
		}
		if pageInt <= 0 {
			return errors.New("page must be > 0")
		}
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			return errors.New("size parse to int failed")
		}
		if sizeInt <= 0 {
			return errors.New("size must be > 0")
		}
	}
	return nil
}

func getParamsForListModels(params models.ListModelsParams) (string, string, string, string, string, string, error) {
	logr := logger.LocLogger(logWithGetListModelsParams(params))
	query := params.HTTPRequest.URL.Query()
	userValue := query["userid"]

	var username string
	if len(userValue) != 0 {
		logr.Debugf("listModels param user: %v", userValue[0])
		username = userValue[0]
	}
	namespaceValue := query["namespace"]

	var namespace string
	if len(namespaceValue) != 0 {
		logr.Debugf("listModels param namespace: %v", namespaceValue[0])
		namespace = namespaceValue[0]
	}
	pageValue := query["page"]

	var page string
	if len(pageValue) != 0 {
		logr.Debugf("listModels param page: %v", pageValue[0])
		page = pageValue[0]
	}
	sizeValue := query["size"]

	var size string
	if len(sizeValue) != 0 {
		logr.Debugf("listModels param size: %v", sizeValue[0])
		size = sizeValue[0]
	}
	expRunIdValue := query["expRunId"]

	var expRunId string
	if len(expRunIdValue) != 0 {
		logr.Debugf("listModels param expRunId: %v", expRunIdValue[0])
		expRunId = expRunIdValue[0]
	}

	// FIXME MLSS Change: v_1.5.0 add logic to get clusterName for request
	clusterNameValue := query["clusterName"]
	var clusterName string
	if len(clusterNameValue) != 0 {
		logr.Debugf("listModels param clusterName: %v", clusterNameValue[0])
		clusterName = clusterNameValue[0]
	} else {
		clusterName = "default"
	}

	return username, namespace, page, size, clusterName, expRunId, nil
}

// FIXME MLSS Change: v_1.5.0 add logic to get clusterName for request
func getParamsForDashboards(params operations.GetDashboardsParams) (string, error) {
	logr := logger.LocLogger(logWithGetDashboards(params))
	query := params.HTTPRequest.URL.Query()
	clusterNameValue := query["clusterName"]
	var clusterName string
	if len(clusterNameValue) != 0 {
		logr.Debugf("getDashboards param clusterName: %v", clusterNameValue[0])
		clusterName = clusterNameValue[0]
	} else {
		clusterName = "default"
	}
	return clusterName, nil
}

// getTrainingLogs establishes a long running http connection and streams the logs of the training container.
func getLogsOrMetrics(params models.GetLogsParams, isMetrics bool) middleware.Responder {
	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Infof("getLogsOrMetrics, params.ModelID: %v", params.ModelID)

	// FIXME MLSS Change: auth user in restapi
	_, errResp := checkJobVisibility(params.HTTPRequest.Context(), params.ModelID, getUserID(params.HTTPRequest), params.HTTPRequest.Header.Get(cc.CcAuthToken))
	if errResp != nil && params.ModelID != "wstest" {
		return errResp
	}

	ctx := params.HTTPRequest.Context()

	isFollow := (params.HTTPRequest.Header.Get("Sec-Websocket-Key") != "")

	if !isFollow {
		isFollow = *params.Follow
	}

	logr.Debugf("isFollow: %v, isMetrics: %v", isFollow, isMetrics)

	var timeout time.Duration

	if params.ModelID == "wstest" {
		logr.Debug("is wstest")
		timeout = 2 * time.Minute

	} else if isFollow {
		// Make this a *very* long time out.  In the longer run, if we
		// push to a message queue, we can hopefully just subscribe to a web
		// socket.
		timeout = 80 * time.Hour
	} else {
		timeout = 3 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	// don't cancel here as we are passing the cancel function to others.

	// HACK FOR WS TEST
	if params.ModelID == "wstest" {
		//return getTrainingLogsWSTest(ctx, params, cancel)
		return getNewTrainingLogsWSTest(ctx, params, cancel)
	}

	//trainer, err := trainerClient.NewTrainer()
	storage, err := storageClient.NewStorage()

	if err != nil {
		logr.WithError(err).Error("Cannot create client for lcm service")
		defer cancel()
		return handleErrorResponse(logr, "")
	}

	//var stream grpc_trainer_v2.Trainer_GetTrainedModelLogsClient
	var stream grpc_storage.Storage_GetTrainedModelLogsClient

	if isMetrics {
		logr.WithError(err).Errorf("GetTrainedModelMetrics has been removed")
		defer cancel()
		return handleErrorResponse(logr, "")
	} else {
		//stream, err = trainer.Client().GetTrainedModelLogs(ctx, &grpc_trainer_v2.TrainedModelLogRequest{
		stream, err = storage.Client().GetTrainedModelLogs(ctx, &grpc_storage.TrainedModelLogRequest{
			TrainingId: params.ModelID,
			UserId:     getUserID(params.HTTPRequest),
			Follow:     isFollow,
		})
	}
	if err != nil {
		logr.WithError(err).Errorf("GetTrainedModelLogs failed")
		defer cancel()
		return handleErrorResponse(logr, "")
	}
	// _, err = *stream.Header()
	// We seem to need this pause, which is somewhat disconcerting. -sb
	//time.Sleep(time.Second * 2)

	if params.HTTPRequest.Header.Get("Sec-Websocket-Key") != "" {
		//return getTrainingLogsWS(trainer, params, stream, cancel, isMetrics)
		//return newGetTrainingLogsWS(trainer, params, stream, cancel, isMetrics)
		return newGetTrainingLogsWSForStorage(storage, params, stream, cancel, isMetrics)
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, prod runtime.Producer) {
		// The close of the LcmClient should also close the stream, as far as i can tell.
		//defer trainer.Close()
		defer cancel()

		//logr.Debugln("w.WriteHeader(200)")
		w.WriteHeader(200)
		// w.Header().Set("Transfer-Encoding", "chunked")
		var onelinebytes []byte
		for {
			var logFrame *grpc_storage.ByteStreamResponse
			logFrame, err := stream.Recv()
			// time.Sleep(time.Second * 2)
			if logFrame == nil {
				if err != io.EOF && err != nil {
					logr.WithError(err).Errorf("stream.Recv() returned error")
				}
				break
			}
			if isMetrics {
				var bytesBuf []byte = logFrame.GetData()
				if bytesBuf != nil {

					byteReader := bytes.NewReader(bytesBuf)
					bufioReader := bufio.NewReader(byteReader)
					for {
						var lineBytes []byte
						lineBytes, err := bufioReader.ReadBytes('\n')
						if lineBytes != nil {
							lenRead := len(lineBytes)

							if err == nil || (err == io.EOF && lenRead > 0 && lineBytes[lenRead-1] == '}') {

								if onelinebytes != nil {
									lineBytes = append(onelinebytes[:], lineBytes[:]...)
									onelinebytes = nil
								}

								_, err := w.Write(lineBytes)
								//logr.Debugf("w.Write(bytes) says %d bytes written", n)
								if err != nil && err != io.EOF {
									logr.Errorf("getTrainingLogs(2) Write returned error: %s", err.Error())
								}
								// logr.Debugln("if f, ok := w.(http.Flusher); ok {")
								if f, ok := w.(http.Flusher); ok {
									logr.Debugln("f.Flush()")
									f.Flush()
								}
							} else {
								if onelinebytes == nil {
									onelinebytes = lineBytes
								} else {
									onelinebytes = append(onelinebytes[:], lineBytes[:]...)
								}
							}
						}
						if err == io.EOF {
							break
						}
					}
				}
			} else {
				var bytes []byte = logFrame.GetData()
				if bytes != nil {
					//logr.Debugln("w.Write(bytes) len = %d", len(bytes))
					_, err := w.Write(bytes)
					//logr.Debugf("w.Write(bytes) says %d bytes written", n)
					if err != nil && err != io.EOF {
						logr.Errorf("getTrainingLogs(2) Write returned error: %s", err.Error())
					}
					//logr.Debugln("if f, ok := w.(http.Flusher); ok {")
					if f, ok := w.(http.Flusher); ok {
						logr.Debugln("f.Flush()")
						f.Flush()
					}
				}
			}
			//logr.Debugln("bottom of for")
		}
	})
}

func makeGrpcSearchTypeFromRestSearchType(st string) grpc_training_data_v1.Query_SearchType {
	var searchType grpc_training_data_v1.Query_SearchType
	switch restmodels.QuerySearchType(st) {
	case restmodels.QuerySearchTypeTERM:
		searchType = grpc_training_data_v1.Query_TERM
		break
	case restmodels.QuerySearchTypeMATCH:
		searchType = grpc_training_data_v1.Query_MATCH
		break
	case restmodels.QuerySearchTypeNESTED:
		searchType = grpc_training_data_v1.Query_NESTED
		break
	case restmodels.QuerySearchTypeALL:
		searchType = grpc_training_data_v1.Query_ALL
		break
	}
	return searchType
}

func handleErrorResponse(log *logger.LocLoggingEntry, description string) middleware.Responder {
	logger.GetLogger().Errorf("Returning 500 error: %s", description)
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(http.StatusInternalServerError)
		//FIXME MLSS Change: change Description to Msg
		payload, _ := json.Marshal(&restmodels.Error{
			Error: "Internal server error",
			Msg:   description,
			Code:  http.StatusInternalServerError,
		})
		w.Write(payload)
	})
}

func getUserID(r *http.Request) string {
	return r.Header.Get(config.CcAuthUser)
}

func getSuperadmin(r *http.Request) string {
	return r.Header.Get(config.CcAuthSuperadmin)
}

func getSuperadminBool(r *http.Request) bool {
	superAdminStr := r.Header.Get(config.CcAuthSuperadmin)
	if "true" == superAdminStr {
		return true
	}
	return false
}

// Echo the data received on the WebSocket.
func serveLogHandler(trainer trainerClient.TrainerClient, stream grpc_trainer_v2.Trainer_GetTrainedModelLogsClient,
	logr *logger.LocLoggingEntry, cancel context.CancelFunc, isMetrics bool) websocket.Handler {

	return func(ws *websocket.Conn) {
		defer ws.Close()
		defer trainer.Close()
		defer cancel()

		// TODO: The second param should be log.LogCategoryServeLogHandler, but, for the
		// moment, just use the hard coded string, until the code is committed in dlaas-commons.

		logr.Debugf("Going into Recv() loop")
		var onelinebytes []byte
		for {
			var logFrame *grpc_trainer_v2.ByteStreamResponse
			logFrame, err := stream.Recv()
			if err == io.EOF {
				logr.Infof("serveLogHandler stream.Recv() is EOF")
			}
			if logFrame != nil && len(logFrame.Data) > 0 {
				if isMetrics {

					byteReader := bytes.NewReader(logFrame.Data)
					bufioReader := bufio.NewReader(byteReader)
					for {
						lineBytes, err := bufioReader.ReadBytes('\n')
						if lineBytes != nil {
							lenRead := len(lineBytes)

							if err == nil || (err == io.EOF && lenRead > 0 && lineBytes[lenRead-1] == '}') {
								if onelinebytes != nil {
									lineBytes = append(onelinebytes[:], lineBytes[:]...)
									onelinebytes = nil
								}
								// We should just scan for the first non-white space.
								if len(bytes.TrimSpace(lineBytes)) == 0 {
									continue
								}

								ws.Write(lineBytes)
								// Take a short snooze, just to not take over CPU, etc.
								// time.Sleep(time.Millisecond * 250)
							} else {
								if onelinebytes == nil {
									onelinebytes = lineBytes
								} else {
									onelinebytes = append(onelinebytes[:], lineBytes[:]...)
								}
							}
						}
						if err == io.EOF {
							break
						}
					}
					logr.Debug("==== done processing logFrame.Data ====")

				} else {
					var bytes []byte
					bytes = logFrame.Data
					n, errWrite := ws.Write(bytes)
					if errWrite != nil && errWrite != io.EOF {
						logr.WithError(errWrite).Errorf("serveLogHandler Write returned error")
						break
					}
					logr.Debugf("wrote %d bytes", n)
				}
			}

			// either EOF or error reading from trainer
			if err != nil {
				logr.WithError(err).Debugf("Breaking from Recv() loop")
				break
			}
			time.Sleep(time.Millisecond * 2)
		}
	}
}

// Echo the data received on the WebSocket.
func newServeLogHandler(params models.GetLogsParams, w http.ResponseWriter, trainer trainerClient.TrainerClient, stream grpc_trainer_v2.Trainer_GetTrainedModelLogsClient,
	logr *logger.LocLoggingEntry, cancel context.CancelFunc, isMetrics bool) {
	r := params.HTTPRequest
	upgrader := newWebsocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLogger().Error("upgrade failed: ", err.Error())
		return
	}
	defer c.Close()
	defer trainer.Close()
	defer cancel()

	// TODO: The second param should be log.LogCategoryServeLogHandler, but, for the
	// moment, just use the hard coded string, until the code is committed in dlaas-commons.

	logr.Debugf("Going into Recv() loop")
	var onelinebytes []byte
	for {
		var logFrame *grpc_trainer_v2.ByteStreamResponse
		logFrame, err := stream.Recv()
		if err == io.EOF {
			logr.Infof("serveLogHandler stream.Recv() is EOF")
		}
		if logFrame != nil && len(logFrame.Data) > 0 {
			if isMetrics {

				byteReader := bytes.NewReader(logFrame.Data)
				bufioReader := bufio.NewReader(byteReader)
				for {
					lineBytes, err := bufioReader.ReadBytes('\n')
					if lineBytes != nil {
						lenRead := len(lineBytes)

						if err == nil || (err == io.EOF && lenRead > 0 && lineBytes[lenRead-1] == '}') {
							if onelinebytes != nil {
								lineBytes = append(onelinebytes[:], lineBytes[:]...)
								onelinebytes = nil
							}
							// We should just scan for the first non-white space.
							if len(bytes.TrimSpace(lineBytes)) == 0 {
								continue
							}

							//ws.Write(lineBytes)
							err = c.WriteMessage(newWebsocket.TextMessage, lineBytes)
							// Take a short snooze, just to not take over CPU, etc.
							// time.Sleep(time.Millisecond * 250)
						} else {
							if onelinebytes == nil {
								onelinebytes = lineBytes
							} else {
								onelinebytes = append(onelinebytes[:], lineBytes[:]...)
							}
						}
					}
					if err == io.EOF {
						break
					}
				}
				logr.Debug("==== done processing logFrame.Data ====")

			} else {
				var bytes []byte
				bytes = logFrame.Data
				//n, errWrite := ws.Write(bytes)
				errWrite := c.WriteMessage(newWebsocket.TextMessage, bytes)

				if errWrite != nil && errWrite != io.EOF {
					logr.WithError(errWrite).Errorf("serveLogHandler Write returned error")
					break
				}
				//logr.Debugf("wrote %d bytes", n)
				logr.Debugf("wrote bytes")
			}
		}

		// either EOF or error reading from trainer
		if err != nil {
			logr.WithError(err).Debugf("Breaking from Recv() loop")
			break
		}
		time.Sleep(time.Millisecond * 2)
	}
	//}
}

func newServeLogHandlerForStorage(params models.GetLogsParams, w http.ResponseWriter, storage storageClient.StorageClient, stream grpc_storage.Storage_GetTrainedModelLogsClient,
	logr *logger.LocLoggingEntry, cancel context.CancelFunc, isMetrics bool) {
	r := params.HTTPRequest
	upgrader := newWebsocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLogger().Error("upgrade failed: ", err.Error())
		return
	}
	defer c.Close()
	defer cancel()

	// TODO: The second param should be log.LogCategoryServeLogHandler, but, for the
	// moment, just use the hard coded string, until the code is committed in dlaas-commons.

	logr.Debugf("Going into Recv() loop")
	var onelinebytes []byte
	for {
		logFrame, err := stream.Recv()
		if err != nil {
			logr.Error(err.Error())
		}
		if err == io.EOF {
			logr.Infof("serveLogHandler stream.Recv() is EOF")
			break
		}
		if logFrame != nil && len(logFrame.Data) > 0 {
			if isMetrics {

				byteReader := bytes.NewReader(logFrame.Data)
				bufioReader := bufio.NewReader(byteReader)
				for {
					lineBytes, err := bufioReader.ReadBytes('\n')
					if lineBytes != nil {
						lenRead := len(lineBytes)

						if err == nil || (err == io.EOF && lenRead > 0 && lineBytes[lenRead-1] == '}') {
							if onelinebytes != nil {
								lineBytes = append(onelinebytes[:], lineBytes[:]...)
								onelinebytes = nil
							}
							// We should just scan for the first non-white space.
							if len(bytes.TrimSpace(lineBytes)) == 0 {
								continue
							}

							//ws.Write(lineBytes)
							err = c.WriteMessage(newWebsocket.TextMessage, lineBytes)
							// Take a short snooze, just to not take over CPU, etc.
							// time.Sleep(time.Millisecond * 250)
						} else {
							if onelinebytes == nil {
								onelinebytes = lineBytes
							} else {
								onelinebytes = append(onelinebytes[:], lineBytes[:]...)
							}
						}
					}
					if err == io.EOF {
						break
					}
				}
				logr.Debug("==== done processing logFrame.Data ====")

			} else {
				var bytes []byte
				bytes = logFrame.Data
				//n, errWrite := ws.Write(bytes)
				errWrite := c.WriteMessage(newWebsocket.TextMessage, bytes)

				if errWrite != nil && errWrite != io.EOF {
					logr.WithError(errWrite).Errorf("serveLogHandler Write returned error")
					break
				}
				//logr.Debugf("wrote %d bytes", n)
				logr.Debugf("wrote bytes")
			}
		}

		// either EOF or error reading from trainer
		if err != nil {
			logr.WithError(err).Debugf("Breaking from Recv() loop")
			break
		}
		time.Sleep(time.Millisecond * 2)
	}
	//}
}

func getTrainingLogsWS(trainer trainerClient.TrainerClient, params models.GetLogsParams,
	stream grpc_trainer_v2.Trainer_GetTrainedModelLogsClient,
	cancel context.CancelFunc, isMetrics bool) middleware.Responder {

	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Debugf("Setting up web socket: %v", params.HTTPRequest.Header)

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		logr.Debugf("In responderFunc")
		handler := serveLogHandler(trainer, stream, logr, cancel, isMetrics)
		handler.ServeHTTP(w, params.HTTPRequest)

	})
}

func newGetTrainingLogsWS(trainer trainerClient.TrainerClient, params models.GetLogsParams,
	stream grpc_trainer_v2.Trainer_GetTrainedModelLogsClient,
	cancel context.CancelFunc, isMetrics bool) middleware.Responder {

	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Debugf("Setting up web socket: %v", params.HTTPRequest.Header)

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		logr.Debugf("In responderFunc")
		newServeLogHandler(params, w, trainer, stream, logr, cancel, isMetrics)
	})
}

func newGetTrainingLogsWSForStorage(storage storageClient.StorageClient, params models.GetLogsParams,
	stream grpc_storage.Storage_GetTrainedModelLogsClient,
	cancel context.CancelFunc, isMetrics bool) middleware.Responder {

	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Debugf("Setting up web socket: %v", params.HTTPRequest.Header)

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		logr.Debugf("In responderFunc")
		newServeLogHandlerForStorage(params, w, storage, stream, logr, cancel, isMetrics)
	})
}

func getNewTrainingLogsWSTest(ctx context.Context, params models.GetLogsParams,
	cancel context.CancelFunc) middleware.Responder {

	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Debugf("Setting up web socket test: %v", params.HTTPRequest.Header)

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		//serveLogHandlerTest(logr, cancel).ServeHTTP(w, params.HTTPRequest)
		newWebsocketHandler(w, params.HTTPRequest)
	})

}

func newWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := newWebsocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.GetLogger().Error("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		s := time.Now().String()
		err = c.WriteMessage(newWebsocket.TextMessage, []byte(s))
		if err != nil {
			logger.GetLogger().Error("write:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}

}

func serveLogHandlerTest(logr *logger.LocLoggingEntry, cancel context.CancelFunc) websocket.Handler {
	logr.Debugf("In serveLogHandlerTest")

	return func(ws *websocket.Conn) {
		defer ws.Close()
		defer cancel()
		for i := 0; i < 30; i++ {
			currentTime := time.Now().Local()
			currentTimeString := currentTime.Format("  Sat Mar 7 11:06:39 EST 2015\n")
			logr.Debugf("In websocket test function, writing %s", currentTimeString)
			_, writeStringError := io.WriteString(ws, currentTimeString)
			if writeStringError != nil {
				logr.WithError(writeStringError).Debug("WriteString failed")
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func getTrainingLogsWSTest(ctx context.Context, params models.GetLogsParams,
	cancel context.CancelFunc) middleware.Responder {

	logr := logger.LocLogger(logWithGetLogsParams(params))
	logr.Debugf("Setting up web socket test: %v", params.HTTPRequest.Header)

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		serveLogHandlerTest(logr, cancel).ServeHTTP(w, params.HTTPRequest)
	})

}

// FIXME MLSS Change: more properties for result
func createModel(req *http.Request, job *grpc_storage.Job, logr *logger.LocLoggingEntry) *restmodels.Model {
	memUnit := job.Training.Resources.MemoryUnit.String()

	// FIXME MLSS Change: more properties for result
	//logr.Infof("createModel job.TrainingId: %v", job.TrainingId)
	cpus := float64(job.Training.Resources.Cpus)
	gpus := float64(job.Training.Resources.Gpus)
	memory := float64(job.Training.Resources.Memory)
	if math.IsNaN(cpus) {
		logr.Debugf("warning, cpus IsNaN: %v", cpus)
		cpus = 0
	}
	if math.IsNaN(gpus) {
		logr.Debugf("warning, gpus IsNaN: %v", gpus)
		gpus = 0
	}
	if math.IsNaN(memory) {
		logr.Debugf("warning, memory IsNaN: %v", memory)
		memory = 0
	}

	psM := ""
	if "" != job.PsMemory {
		if stringsUtil.Contains(job.PsMemory, "Gi") {
			psM = job.PsMemory[0 : len(job.PsMemory)-2]
		}
	}

	//var tFosRequest *v1alpha1.TFosRequest
	var tFosRequest *restmodels.TFosRequest
	str := job.TfosRequest
	unmarshal := json.Unmarshal([]byte(str), &tFosRequest)
	if job.JobType == "tfos" {
		if nil != unmarshal {
			logr.Errorf("failed format TFosRequest with err: %v", unmarshal.Error())
		}
	}

	m := &restmodels.Model{
		BasicNewModel: restmodels.BasicNewModel{
			BasicModel: restmodels.BasicModel{
				ModelID: job.TrainingId,
			},
			Location: req.URL.Path + "/" + job.TrainingId,
		},
		Pss:         job.Pss,
		PsCPU:       job.PsCpu,
		PsImage:     job.PsImage,
		PsMemory:    psM,
		JobType:     job.JobType,
		Name:        job.ModelDefinition.Name,
		Description: job.ModelDefinition.Description,
		Framework: &restmodels.Framework{
			Name:    job.ModelDefinition.Framework.Name,
			Version: job.ModelDefinition.Framework.Version,
		},
		Training: &restmodels.Training{
			Command:    job.Training.Command,
			Cpus:       cpus,
			Gpus:       gpus,
			Memory:     memory,
			MemoryUnit: &memUnit,
			Learners:   job.Training.Resources.Learners,
			InputData:  job.Training.InputData,
			OutputData: job.Training.OutputData,
			TrainingStatus: &restmodels.TrainingStatus{
				Status:            job.Status.Status.String(),
				StatusDescription: job.Status.Status.String(),
				StatusMessage:     job.Status.StatusMessage,
				ErrorCode:         job.Status.ErrorCode,
				Submitted:         job.Status.SubmissionTimestamp,
				Completed:         job.Status.CompletionTimestamp,
			},
		},

		// FIXME MLSS Change: more properties for result
		JobNamespace:        job.JobNamespace,
		UserID:              job.UserId,
		SubmissionTimestamp: job.Status.SubmissionTimestamp,
		CompletedTimestamp:  job.Status.CompletionTimestamp,
		// FIXME MLSS Change: v_1.4.1 added field  JobAlert
		JobAlert:     job.JobAlert,
		ExpRunID:     job.ExpRunId,
		ExpName:      job.ExpName,
		FileName:     job.FileName,
		FilePath:     job.FilePath,
		TFosRequest:  tFosRequest,
		ProxyUser:    job.ProxyUser,
		CodeSelector: job.CodeSelector,
	}

	// add data stores
	for i, v := range job.Datastores {
		m.DataStores = append(m.DataStores, &restmodels.Datastore{
			DataStoreID: v.Id,
			Type:        v.Type,
			Connection:  v.Connection,
		})
		for k, v := range job.Datastores[i].Fields {
			m.DataStores[i].Connection[k] = v
		}
	}
	return m
}

func getTrainingJobFromStorage(params models.GetModelParams) (*grpc_storage.GetResponse, middleware.Responder) {

	getLogger := logger.GetLogger()
	// FIXME MLSS Change: auth user in restapi
	rresp, errResp := checkJobVisibility(params.HTTPRequest.Context(), params.ModelID, getUserID(params.HTTPRequest), params.HTTPRequest.Header.Get(cc.CcAuthToken))
	if errResp != nil {
		getLogger.Errorf("getTrainingJobFromSorage, checkJobVisibility failed.")
		return nil, errResp
	}

	return rresp, middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(http.StatusOK)
		//FIXME MLSS Change: change Description to Msg
		payload, _ := json.Marshal(&restmodels.Error{
			Error: "OK",
			Msg:   "",
			Code:  http.StatusOK,
		})
		w.Write(payload)
	})
}

func getErrorResultWithPayload(error string, status int32, err error) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(http.StatusForbidden)
		//FIXME MLSS Change: change Description to Msg
		payload, _ := json.Marshal(restmodels.Error{
			Error: error,
			Code:  status,
			Msg:   err.Error(),
		})
		w.Write(payload)
	})
}

func checkJobVisibility(ctx context.Context, modelID string, userId string, ccAuthToken string) (*grpc_storage.GetResponse, middleware.Responder) {
	operation := "checkJobVisibility"
	//logr := logger.LocLogger(log.StandardLogger().WithField("module", "restapi,models_impl.go,checkJobVisibility"))
	//trainer, err := trainerClient.NewTrainer()
	storage, err := storageClient.NewStorage()
	defer storage.Close()
	if err != nil {
		logger.GetLogger().Error("New Storage Client Error: " + err.Error())
		return nil, getErrorResultWithPayload("New Client Error", http.StatusInternalServerError, err)
	}

	//rresp, err := trainer.Client().GetTrainingJob(ctx, &grpc_trainer_v2.GetRequest{
	rresp, err := storage.Client().GetTrainingJob(ctx, &grpc_storage.GetRequest{
		TrainingId: modelID,
		UserId:     userId,
	})
	if err != nil {
		logger.GetLogger().Errorf("Trainer GetTrainingJob service call failed")
		if grpc.Code(err) == codes.PermissionDenied {
			//return nil, getErrorResultWithPayload("forbidden", http.StatusForbidden, err)
			return nil, getErrorResultWithPayload("forbidden", http.StatusForbidden, err)
		}
		if grpc.Code(err) == codes.NotFound {
			//FIXME MLSS Change: change Description to Msg
			return nil, models.NewGetModelNotFound().WithPayload(&restmodels.Error{
				Error: "Not found",
				Code:  http.StatusNotFound,
				Msg:   "",
			})
		}
		return nil, httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}

	if rresp.Job == nil {
		//FIXME MLSS Change: change Description to Msg
		return nil, models.NewGetModelNotFound().WithPayload(&restmodels.Error{
			Error: "Not found",
			Code:  http.StatusNotFound,
			Msg:   "",
		})
	}

	// FIXME MLSS Change: auth user in restapi
	userFromJob := rresp.Job.UserId
	logger.GetLogger().Debugf("userFromJob: %v and userId: %v", userFromJob, userId)
	if userFromJob != userId {
		nsFromJob := rresp.Job.JobNamespace
		ccAddress := viper.GetString(config.CCAddress)
		ccClient := cc.GetCcClient(ccAddress)
		bodyForRole, err := ccClient.CheckNamespace(ccAuthToken, nsFromJob)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return nil, getErrorResultWithPayload("", http.StatusInternalServerError, err)
		}
		var role string
		err = commonModels.GetResultData(bodyForRole, &role)
		if err != nil {
			logger.GetLogger().Error("GetResultData failed")
			return nil, getErrorResultWithPayload("", http.StatusInternalServerError, err)
		}
		if role == "GU" {
			logger.GetLogger().Error("GU only access his own notebook")
			return nil, getErrorResultWithPayload("StatusForbidden", http.StatusForbidden, err)
		}
	}
	return rresp, nil
}

func checkForJobParams(m *ManifestV1, params models.PostModelParams, codeSelector string, logr *logger.LocLoggingEntry) ([]byte, error) {
	var modelDefinition []byte = nil
	if codeSelector == codeFile && nil == params.ModelDefinition && m.ExpName != "" {
		return nil, errors.New("ModelDefinition can't be empty when codeSelector is codeFile")
	}
	if codeSelector == codeFile && nil != m.DataStores[0].TrainingWorkspace && "" != m.DataStores[0].TrainingWorkspace.Container {
		return nil, errors.New("training_work must be empty when codeSelector is codeFile")
	}
	if codeSelector == codeFile && nil != params.ModelDefinition {
		modelDefinitionData, err := ioutil.ReadAll(params.ModelDefinition)
		if err != nil {
			logr.Errorf("Cannot read 'model_definition' parameter: %s", err.Error())
			return nil, err

		}
		modelDefinition = modelDefinitionData
	}
	if codeSelector == storagePath && nil != params.ModelDefinition {
		return nil, errors.New("ModelDefinition must be empty when codeSelector is storagePath")
	}
	if codeSelector == storagePath && (nil == m.DataStores[0].TrainingWorkspace || (nil != m.DataStores[0].TrainingWorkspace && "" == m.DataStores[0].TrainingWorkspace.Container)) {
		return nil, errors.New("training_work can't be empty when codeSelector is storagePath")
	}
	return modelDefinition, nil
}

func checkAlertParams(mannifest *ManifestV1, logr *logger.LocLoggingEntry) error {
	//receiver := mannifest.Receiver
	jobAlert := mannifest.JobAlert
	codeSelector := mannifest.CodeSelector
	event := jobAlert["event"]
	deadline := jobAlert["deadline"]
	overtime := jobAlert["overtime"]
	if nil != event && len(event) > 0 {
		for _, eMap := range event {
			alertLevel := eMap["alert_level"]
			receiver := eMap["receiver"]
			eventChecker := eMap["event_checker"]
			if receiver == "" || eventChecker == "" {
				return errors.New("event receiver can't be nil")
			}
			if eventChecker != jobmonitor.FAILED && eventChecker != jobmonitor.COMPLETED {
				return errors.New("event_checker must be FAILED or COMPLETED")
			}
			if alertLevel != "" && alertLevel != jobmonitor.CRITICAL && alertLevel != jobmonitor.MAJOR && alertLevel != jobmonitor.MINOR && alertLevel != jobmonitor.WARNING && alertLevel != jobmonitor.INFO {
				return errors.New("event alert_level must be critical, major, minor, warning, info")
			}
		}
	}

	if nil != deadline && len(deadline) > 0 {
		for _, dMap := range deadline {
			alertLevel := dMap["alert_level"]
			interval, atoErr := strconv.Atoi(dMap["interval"])
			if atoErr != nil {
				return errors.New("deadline interval parse int failed")
			}
			receiver := dMap["receiver"]
			deadlineChecker := dMap["deadline_checker"]
			if receiver == "" || deadlineChecker == "" {
				return errors.New("alertLevel, receiver and deadlineChecker can't be nil")
			}
			_, err := time.Parse("2006-01-02 15:04:05", deadlineChecker)
			if err != nil {
				return errors.New("deadline_checker pattern must be 2006-01-02 15:04:05 ")
			}
			if alertLevel != "" && alertLevel != jobmonitor.CRITICAL && alertLevel != jobmonitor.MAJOR && alertLevel != jobmonitor.MINOR && alertLevel != jobmonitor.WARNING && alertLevel != jobmonitor.INFO {
				return errors.New("deadline alert_level must be critical, major, minor, warning, info")
			}
			if interval < 1 {
				return errors.New("deadline interval has to be gretter than 0")
			}
		}
	}

	if nil != overtime && len(overtime) > 0 {
		for _, oMap := range overtime {
			alertLevel := oMap["alert_level"]
			interval, atoErr := strconv.Atoi(oMap["interval"])
			if atoErr != nil {
				return errors.New("overtime interval parse int failed")
			}
			receiver := oMap["receiver"]
			overtimeChecker := oMap["overtime_checker"]

			if receiver == "" || overtimeChecker == "" {
				return errors.New("alertLevel, receiver and overtimeChecker can't be nil")
			}
			f, parseErr := strconv.ParseFloat(overtimeChecker, 64)
			if parseErr != nil {
				return errors.New("overtime_checker must be number")
			}
			if f <= 0 {
				return errors.New("overtime_checker must be greater than 0")
			}
			if alertLevel != "" && alertLevel != jobmonitor.CRITICAL && alertLevel != jobmonitor.MAJOR && alertLevel != jobmonitor.MINOR && alertLevel != jobmonitor.WARNING && alertLevel != jobmonitor.INFO {
				return errors.New("overtime alert_level must be critical, major, minor, warning, info")
			}
			if interval < 1 {
				return errors.New("overtime interval has to be gretter than 0")
			}
		}
	}
	if mannifest.JobType != "tfos" {
		if codeSelector == "" {
			return errors.New("codeSelector is request")
		}
		if codeSelector != codeFile && codeSelector != storagePath {
			return errors.New("the value of codeSelector can only be codeFile or storagePath")
		}
	}
	return nil
}

func httpResponseHandle(status int, err error, operation string, resultMsg []byte) middleware.Responder {
	result := commonModels.Result{
		Code:   "200",
		Msg:    "success",
		Result: json.RawMessage(resultMsg),
	}

	if resultMsg == nil {
		jsonStr, err := json.Marshal(operation + " " + "success")
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
		result.Result = jsonStr
	}

	if status != http.StatusOK {
		result.Code = strconv.Itoa(status)
		result.Msg = "Error"
		jsonStr, _ := json.Marshal(operation + " " + " Error")
		result.Result = jsonStr
		if err != nil {
			logger.GetLogger().Error("Operation: " + operation + ";" + "Error: " + err.Error())
			jsonStr, _ := json.Marshal(operation + " " + " error: " + err.Error())
			result.Result = jsonStr
		}
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(status)
		response, err := json.Marshal(result)
		_, err = w.Write(response)
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	})
}

func GetJobLogByLine(params models.GetJobLogByLineParams) middleware.Responder {
	user := getUserID(params.HTTPRequest)
	//client, err := storageClient.NewStorage()
	logger.GetLogger().Debugf("GetJobLogByLine params training id: %v\n", params.TrainingID)
	address := fmt.Sprintf("di-storage-rpc.%s.svc.cluster.local:80", viper.GetString(config.PlatformNamespace))
	client, err := storageClient.NewStorageWithAddress(address)
	if err != nil {
		logger.GetLogger().Error("New storage client error, ", err.Error())
		return handleErrorResponse(log, err.Error())
	}
	logger.GetLogger().Debugln("GetJobLogByLine create client success.")
	req := &grpc_storage.GetTrainedModelLogRequest{
		TrainingId: params.TrainingID,
		Size:       params.Size,
		From:       params.From,
		UserId:     user,
	}
	req2 := &grpc_storage.GetTrainedModelLogRequest{
		TrainingId: params.TrainingID,
		Size:       10000,
		From:       0,
		UserId:     user,
	}
	res, err := client.Client().GetTrainedModelLog(context.Background(), req)
	if err != nil {
		logger.GetLogger().Error("Get training log error, ", err.Error())
		return handleErrorResponse(log, err.Error())
	}
	logger.GetLogger().Debugln("GetJobLogByLine get log success.")
	if res == nil {
		logger.GetLogger().Error("Response is nil")
		return handleErrorResponse(log, "Response is nil")
	}
	res2, err := client.Client().GetTrainedModelLog(context.Background(), req2)
	if err != nil {
		logger.GetLogger().Error("Get training log error, ", err.Error())
		return handleErrorResponse(log, err.Error())
	}
	if res2 == nil {
		logger.GetLogger().Error("Response2 is nil")
		return handleErrorResponse(log, "Response2 is nil")
	}
	isLastLine := "false"
	if int64(len(res2.Logs)) <= req.Size {
		isLastLine = "true"
	}
	return models.NewGetJobLogByLineOK().WithPayload(&restmodels.LogResponse{
		Logs:       res.Logs,
		IsLastLine: isLastLine,
	})
}

func KillTrainingModel(params models.KillTrainingModelParams) middleware.Responder {
	logr := logger.GetLogger()
	logr.Debugf("deleteModel invoked: %v", params.HTTPRequest.Header)
	storageClient, err := storageClient.NewStorage()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		return handleErrorResponse(logr, err.Error())
	}
	defer storageClient.Close()

	res, err := storageClient.Client().KillTrainingJob(params.HTTPRequest.Context(), &grpc_storage.KillRequest{
		TrainingId: params.ModelID,
		UserId:     getUserID(params.HTTPRequest),
		IsSa:       getSuperadminBool(params.HTTPRequest),
	})
	if err != nil {
		logr.Error("Killing Model Error: ", err.Error())
		if grpc.Code(err) == codes.PermissionDenied {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				payload, _ := json.Marshal(restmodels.Error{
					Error: "forbidden",
					Code:  http.StatusForbidden,
					Msg:   err.Error(),
				})
				w.Write(payload)
			})
		}
		return handleErrorResponse(logr, res.String())
	}
	if !res.Success {
		if res.IsCompleted == true {
			return handleErrorResponse(logr, "任务已完成，无法停止任务，请刷新页面。")
		}
		return handleErrorResponse(logr, res.String())
	}
	return models.NewKillTrainingModelOK().WithPayload(
		&restmodels.Model{
			Name: params.ModelID,
		})

}

func CommandTransform(command string) (string, error) {

	// Get ${} Expression  (\${).*?(\})
	re := regexp.MustCompile("(\\${).*?(\\})")
	expArr := re.FindAllString(command, -1)
	expArrIndex := re.FindAllStringIndex(command, -1)

	if expArrIndex == nil {
		return command, nil
	}

	// For Loop Pre-Processing
	var originStr []string
	reCheck := regexp.MustCompile("run_date([+*/-]|\\d+)?")
	for i := 0; i < len(expArr); i++ {
		s := stringsUtil.TrimSpace(expArr[i])
		if reCheck.MatchString(s) {
			originStr = append(originStr, expArr[i])
		}
	}

	// Cal Expression
	var transString []string
	for i := 0; i < len(originStr); i++ {
		runDate, err := getRunDateFromExp(originStr[i])
		transString = append(transString, runDate)
		if err != nil {
			return "", err
		}
	}

	if transString == nil {
		return command, nil
	}

	// Replace ${run_date}
	nCommand := command
	for i := 0; i < len(originStr); i++ {
		nCommand = stringsUtil.ReplaceAll(nCommand, originStr[i], transString[i])
	}

	return nCommand, nil
}

func getRunDateFromExp(exp string) (string, error) {
	nExp := exp[2 : len(exp)-1]
	re := regexp.MustCompile("\\d+|[+/*-]")
	expArr := re.FindAllString(nExp, len("run_date")-1)
	runDate := time.Now().AddDate(0, 0, -1)

	for i := 0; i < len(expArr)-1; i++ {
		dayInt, err := strconv.Atoi(expArr[i+1])
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return "", err
		}
		if expArr[i] == "-" {
			runDate = runDate.AddDate(0, 0, -dayInt)
		} else if expArr[i] == "+" {
			runDate = runDate.AddDate(0, 0, dayInt)
		}
	}
	runDateStr := runDate.Format("20060102")
	return runDateStr, nil
}
