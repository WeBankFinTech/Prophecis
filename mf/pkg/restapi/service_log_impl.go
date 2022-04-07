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
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/restapi/operations/container"
	"mlss-mf/pkg/service"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func GetNamespacedServiceContainerLog(params container.GetNamespacedServiceContainerLogParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	resp, err := service.GetNamespacedServiceContainerLog(params, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "GetNamespacedServiceContainerLog", nil)
		}
		logger.Logger().Errorf("fail to get container log information, error: %s", err.Error())
		return returnResponse(http.StatusInternalServerError, err, "GetNamespacedServiceContainerLog", []byte(err.Error()))
	}
	respBytes, err := json.Marshal(&resp)
	if err != nil {
		logger.Logger().Errorf("fail to json marshal container log information, error: %s", err.Error())
		return returnResponse(http.StatusInternalServerError, err, "GetNamespacedServiceContainerLog", []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, "GetNamespacedServiceContainerLog", respBytes)
}
