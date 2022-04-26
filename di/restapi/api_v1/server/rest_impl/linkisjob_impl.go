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
	"net/http"
	"webank/DI/restapi/api_v1/server/operations/linkis_job"
	"webank/DI/restapi/service"
)

var linkisJobService = service.LinkisJobService

func GetLinkisJobStatus(params linkis_job.GetLinkisJobStatusParams) middleware.Responder {
	operation := "GetLinkisJobStatus"
	user := getUserID(params.HTTPRequest)
	execData, err := linkisJobService.GetLinkisJobStatus(params.LinkisExecID, user)
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

func GetLinkisJobLog(params linkis_job.GetLinkisJobLogParams) middleware.Responder {
	operation := "GetLinkisJobStatus"
	user := getUserID(params.HTTPRequest)
	execData, err := linkisJobService.GetLinkisJobLog(params.LinkisTaskID, user)
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
