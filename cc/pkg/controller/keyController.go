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
	"github.com/go-openapi/runtime/middleware"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/keys"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func AddKey(params keys.AddKeyParams) middleware.Responder {
	name := params.KeyPair.Name
	var err error
	keyByName, err := service.GetKeyByName(name)
	if err !=nil {
		logger.Logger().Error("Add key err, ",err)
		return ResponderFunc(http.StatusBadRequest, err.Error(),err.Error())
	}
	if name == keyByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to add key", "key is exist in db")
	}
	keyPair := common.FromKeyRequestToKeyPair(params.KeyPair)
	deleteKeyByName, err := service.GetDeleteKeyByName(name)
	if err !=nil {
		logger.Logger().Error("Get delete key by name err, ",err)
		return ResponderFunc(http.StatusBadRequest, err.Error(),err.Error())
	}
	var repoKey models.Keypair
	if name == deleteKeyByName.Name {
		keyPair.ID = deleteKeyByName.ID
		repoKey, err = service.UpdateKey(keyPair)
		if err != nil {
			logger.Logger().Error("Update key err, ", err)
			return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
		}
	} else {
		repoKey, err = service.AddKey(keyPair)
		if err != nil {
			logger.Logger().Error("Add key err, ", err)
			return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
		}
	}

	marshal, marshalErr := json.Marshal(repoKey)

	return GetResult(marshal, marshalErr)
}

func GetByName(params keys.GetByNameParams) middleware.Responder {
	name := params.Name
	keyByName, err := service.GetKeyByName(name)
	if err != nil {
		logger.Logger().Error("Update err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if name != keyByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to get key by name", "key is not exist in db")
	}

	marshal, marshalErr := json.Marshal(keyByName)

	return GetResult(marshal, marshalErr)
}

func DeleteByName(params keys.DeleteByNameParams) middleware.Responder {
	name := params.Name
	keyByName, err := service.GetKeyByName(name)
	if err != nil {
		logger.Logger().Error("Get key by name err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if name != keyByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to delete key by name", "key is not exist in db")
	}

	deleteByName, err := service.DeleteByName(name)
	if err != nil {
		logger.Logger().Error("Delete by name err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(deleteByName)

	return GetResult(marshal, marshalErr)
}
