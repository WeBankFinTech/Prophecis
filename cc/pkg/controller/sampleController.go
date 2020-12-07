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
package controller

import (
	"encoding/json"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/samples"
	"net/http"
)

func SamplePost(params samples.SamplePostParams) middleware.Responder {
	request := params.HTTPRequest

	log := logger.Logger()
	log.Debugf("SamplePost request header: %v", request.Header)
	log.Debugf("SamplePost request userId: %v", request.Header.Get(constants.SAMPLE_REQUEST_USER))

	username := request.Header.Get(constants.SAMPLE_REQUEST_USER)

	if "" == username {
		return ResponderFunc(http.StatusBadRequest, "failed to access sample", "failed to get username from header")
	}

	token := request.Header.Get(constants.AUTH_HEADER_TOKEN)
	cookieKey := request.Header.Get(constants.AUTH_HEADER_COOKIE)

	log.Debugf("SamplePost request cookieKey: %v and token: %v", cookieKey, token)

	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    token,
		HttpOnly: true,
		Path:     common.GetAppConfig().Core.Cookie.Path,
		//MaxAge:   common.GetAppConfig().Core.Cookie.DefaultTime,
	}
	log.Debugf("SamplePost request cookie.domain: %v", cookie.Domain)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage([]byte{}),
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		http.SetCookie(w, cookie)
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})
}

func SampleGet(params samples.SampleGetParams) middleware.Responder {
	request := params.HTTPRequest
	log := logger.Logger()
	log.Debugf("SampleGet request header: %v", request.Header)
	log.Debugf("SampleGet request userId: %v", request.Header.Get(constants.SAMPLE_REQUEST_USER))

	username := request.Header.Get(constants.SAMPLE_REQUEST_USER)

	if "" == username {
		return ResponderFunc(http.StatusBadRequest, "failed to access sample", "failed to get username from header")
	}

	return GetResult([]byte{}, nil)
}
