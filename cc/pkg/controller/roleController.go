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
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/roles"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
)

func GetRoles(params roles.GetRolesParams) middleware.Responder {
	r, err := service.GetRoles()
	if err != nil {
		logger.Logger().Error("GetRoles err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(r)

	return GetResult(marshal, marshalErr)
}

func AddRole(params roles.AddRoleParams) middleware.Responder {
	roleRequest := common.FromAddRoleRequest(params)

	roleByName, err := service.GetRoleByName(roleRequest.Name)
	if err != nil {
		logger.Logger().Error("AddRole err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if roleRequest.Name == roleByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to add role to db", "role is exist in db")
	}

	deleteRoleByName, err := service.GetDeleteRoleByName(roleRequest.Name)
	if err != nil {
		logger.Logger().Error("DeleteRole err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	var role models.Role
	if roleRequest.Name == deleteRoleByName.Name {
		roleRequest.ID = deleteRoleByName.ID
		role, err = service.UpdateRole(roleRequest)
		if err != nil {
			logger.Logger().Error("UpdateRole err, ", err)
			return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
		}
	} else {
		role, err = service.AddRole(roleRequest)
		if err != nil {
			logger.Logger().Error("AddRole err, ", err)
			return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
		}
	}

	marshal, marshalErr := json.Marshal(role)

	return GetResult(marshal, marshalErr)
}

func UpdateRole(params roles.UpdateRoleParams) middleware.Responder {
	roleRequest := common.FromUpdateRoleRequest(params)

	roleById, err := service.GetRoleById(roleRequest.ID)
	if err != nil {
		logger.Logger().Error("GetRoles err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if roleRequest.ID != roleById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to update role to db", "role is not exist in db")
	}

	role, err := service.UpdateRole(roleRequest)
	if err != nil {
		logger.Logger().Error("UpdateRole err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	marshal, marshalErr := json.Marshal(role)

	return GetResult(marshal, marshalErr)
}

func GetRoleById(params roles.GetRoleByIDParams) middleware.Responder {
	id := params.ID

	roleById, err := service.GetRoleById(id)
	if err != nil {
		logger.Logger().Error("GetRoleById err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if id != roleById.ID {
		return ResponderFunc(http.StatusBadRequest, "failed to get role by id", "role is not exist in db")
	}

	marshal, marshalErr := json.Marshal(roleById)

	return GetResult(marshal, marshalErr)
}

func GetRoleByName(params roles.GetRoleByNameParams) middleware.Responder {
	name := params.Name

	roleByName, err := service.GetRoleByName(name)
	if err != nil {
		logger.Logger().Error("GetRoleByName err, ", err)
		return ResponderFunc(http.StatusBadRequest, err.Error(), err.Error())
	}
	if name != roleByName.Name {
		return ResponderFunc(http.StatusBadRequest, "failed to get role by name", "role is not exist in db")
	}

	marshal, marshalErr := json.Marshal(roleByName)

	return GetResult(marshal, marshalErr)
}
