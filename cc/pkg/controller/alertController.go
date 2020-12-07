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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/alerts"
	"net/http"
	"strings"
)

func ReceiveAlert(c *gin.Context) {
	addr := c.Request.RemoteAddr

	logger.Logger().Infof("ReceiveAlert addr: %v", addr)

	request := models.ReceiverAlertRequest{}
	c.ShouldBind(&request)

	marshal, err := json.Marshal(request)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:8080/cc/v1/user", strings.NewReader(string(marshal)))
	if err != nil {
		// handle error
		fmt.Println(err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("user", "buoy")
	req.Header.Set("pwd", "123")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

func ReceiveTaskAlert(params alerts.ReceiveTaskAlertParams) middleware.Responder {
	logger.Logger().Debugf("receiveTaskAlert request: %v", params.JobRequest)
	ipAddr := params.HTTPRequest.RemoteAddr
	ims := common.GetAppConfig().Core.Ims

	alertRequest, covertErr := common.FromAlertRequest(params, ipAddr, ims.AlertReceiver)
	if nil != covertErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to format alertRequest", covertErr.Error())
	}

	var alert = map[string][]models.IMSAlert{}
	alert["alertList"] = alertRequest
	marshal, marErr := json.Marshal(alert)
	if nil != marErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to format alertRequest", marErr.Error())
	}

	var imsRequestJson = string(marshal)
	logger.Logger().Debugf("receiveTaskAlert imsRequestJson: %v", imsRequestJson)

	imsUrl := ims.AccessUrl
	logger.Logger().Debugf("receiveTaskAlert imsUrl: %v", imsUrl)

	resp, postErr := http.Post(imsUrl, "application/json", strings.NewReader(imsRequestJson))
	if nil != postErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to post send alert info to ims", postErr.Error())
	}

	code := resp.StatusCode

	body, readErr := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if nil != readErr {
		return ResponderFunc(http.StatusInternalServerError, "failed to readBody", readErr.Error())
	}

	logger.Logger().Debugf("receiveTaskAlert resp code: %v, body: %v", code, string(body))

	marshal, marshalErr := json.Marshal(body)

	return GetResult(marshal, marshalErr)
}
