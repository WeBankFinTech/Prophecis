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
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/storages"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func GetAllStorage(params storages.GetAllStorageParams) middleware.Responder {
	page := *params.Page
	size := *params.Size
	logger.Logger().Debugf("v1/groups GetAllGroups params page: %v, size: %v", page, size)

	repoStorage, err := service.GetAllStorage(page, size)
	if err != nil {
		logger.Logger().Error("Get all storage err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(repoStorage)

	return GetResult(marshal, marshalErr)
}

func AddStorage(params storages.AddStorageParams) middleware.Responder {
	checkForStorageErrMsg := common.CheckForStorage(*params.Storage)
	if checkForStorageErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to add storage", checkForStorageErrMsg.Error())
	}
	storageRequest := common.FromAddStorageRequest(params)
	deleteStorage, err := service.GetDeleteStorageByPath(storageRequest.Path)
	var repoStorage *models.Storage
	if err != nil {
		repoStorage, err = service.AddStorage(storageRequest)
		if err != nil {
			return ResponderFunc(http.StatusBadRequest, "failed to add storage", err.Error())
		}
	} else {
		if deleteStorage != nil {
			storageRequest.ID = deleteStorage.ID
			repoStorage, err = service.UpdateStorage(storageRequest)
			if err != nil {
				return ResponderFunc(http.StatusBadRequest, "failed to update storage", err.Error())
			}
		}
	}
	marshal, marshalErr := json.Marshal(repoStorage)
	return GetResult(marshal, marshalErr)
}

func UpdateStorage(params storages.UpdateStorageParams) middleware.Responder {
	checkForStorageErrMsg := common.CheckForStorage(*params.Storage)

	if checkForStorageErrMsg != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update storage", checkForStorageErrMsg.Error())
	}

	storageRequest := common.FromUpdateStorageRequest(params)

	storageById, err := service.GetStorageById(storageRequest.ID)
	if err != nil {
		logger.Logger().Error("Get storage by id err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if storageRequest.ID != storageById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update storage", "storage is not exist in db")
	}

	repoStorage, err := service.UpdateStorage(storageRequest)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to update storage", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoStorage)

	return GetResult(marshal, marshalErr)
}

func DeleteStorageByID(params storages.DeleteStorageByIDParams) middleware.Responder {
	storageId := params.StorageID

	storageById, err := service.GetStorageById(storageId)
	if err != nil {
		logger.Logger().Error("Get all storage by id err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if storageId != storageById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to delete storage", "storage is not exist in db")
	}

	err = service.DeleteStorageById(storageId)
	if nil != err {
		ResponderFunc(http.StatusBadRequest, "failed to delete storage by id", err.Error())
	}

	marshal, marshalErr := json.Marshal(storageById)

	return GetResult(marshal, marshalErr)
}

func GetStorageByPath(params storages.GetStorageByPathParams) middleware.Responder {
	path := params.Path

	repoStorage, err := service.GetStorageByPath(path)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to get storage by path", err.Error())
	}

	marshal, marshalErr := json.Marshal(repoStorage)

	return GetResult(marshal, marshalErr)

}

func DeleteStorageByPath(params storages.DeleteStorageByPathParams) middleware.Responder {
	path := params.Path

	storageByPath, err := service.GetStorageByPath(path)
	if err != nil {
		return ResponderFunc(http.StatusBadRequest, "failed to delete storage by path", err.Error())
	}

	err = service.DeleteStorageById(storageByPath.ID)
	if nil != err {
		ResponderFunc(http.StatusBadRequest, "failed to delete storage by id", err.Error())
	}

	marshal, marshalErr := json.Marshal(storageByPath)

	return GetResult(marshal, marshalErr)
}
