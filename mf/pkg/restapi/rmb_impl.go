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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/restapi/operations/rmb"
	"mlss-mf/pkg/service"
	"net/http"
)

func DownloadRmbLogByEventID(params rmb.DownloadRmbLogByEventIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	eventId := params.EventID
	bytes, fileName, err := service.DownloadRmbLogByEventID(eventId, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "DownloadRmbLogByEventID", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "DownloadRmbLogByEventID", nil)
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.Header().Add("File-Name", fileName)
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	})
}
