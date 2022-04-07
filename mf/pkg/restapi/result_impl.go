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
	"fmt"
	"mlss-mf/pkg/restapi/operations/model_result"
	"mlss-mf/pkg/service"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func GetResultByModelName(params model_result.GetResultByModelNameParams) middleware.Responder {
	results, err := service.GetResultByModelName(params)
	if err != nil {
		err := fmt.Errorf("fail to get result list by model name, error: %v", err)
		return returnResponse(http.StatusInternalServerError, err, "GetResultByModelName", []byte(err.Error()))
	}
	resultsJson, err := json.Marshal(results)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetResultByModelName", []byte(err.Error()))
	}

	return returnResponse(http.StatusOK, err, "GetResultByModelName", resultsJson)
}

func DeleteResultByID(params model_result.DeleteResultByIDParams) middleware.Responder {
	err := service.DeleteResultByID(params)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "DeleteResultByID", []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, "DeleteResultByID", nil)
}

func CreateResult(params model_result.CreateResultParams) middleware.Responder {
	err := service.CreateResult(params)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "CreateResult", []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, "CreateResult", nil)
}

func UpdateResultByID(params model_result.UpdateResultByIDParams) middleware.Responder {
	err := service.UpdateResultByID(params)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "UpdateResultByID", []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, "UpdateResultByID", nil)
}
