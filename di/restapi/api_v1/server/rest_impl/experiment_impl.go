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
	"net/http"
	"webank/DI/commons/logger"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/experiments"
	"webank/DI/restapi/service"

	"github.com/go-openapi/runtime/middleware"
)

var experimentService = service.ExperimentService
var log = logger.GetLogger()

func CodeUpload(params experiments.CodeUploadParams) middleware.Responder {
	operation := "CodeUpload"
	log.Info("experimentService")
	log.Info(experimentService)
	log.Info("File")
	log.Info(params.File)

	s3Path, err := experimentService.UploadCode(params.File)

	if err != nil {
		log.Error(err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	res := restmodels.CodeUploadResponse{S3Path: s3Path}

	marshal, err := json.Marshal(&res)
	if err != nil {
		log.Error("json.Marshal(mp) failed")
		return httpResponseHandle(http.StatusInternalServerError, err, operation, nil)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)

}