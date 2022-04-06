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
	"errors"
	"math"
	"mlss-mf/pkg/common"
	cc "mlss-mf/pkg/common/controlcenter"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/restapi/operations/image"
	"mlss-mf/pkg/service"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
)

var (
	imageService    = service.ImageService{}
	modelVersionDao = dao.ModelVersionDao{}
	imageDao        = dao.ImageDao{}
)

func CreateImage(params image.CreateImageParams) middleware.Responder {
	operation := "CreateImage"
	//Permission Check
	flag, dataFlag, err := imageDao.HasImageByImageName(params.Image.ImageName)
	if err != nil {
		logger.Logger().Error("Has image error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	if flag && dataFlag == false {
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(errors.New("Image name already exists").Error()))
	}
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := service.CheckGroupAuth(currentUserId, params.Image.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, operation, []byte("Check group auth failed"))
	}

	var image *gormmodels.Image
	if dataFlag == false {
		//Add image
		if err != nil {
			logger.Logger().Error("Create image error, ", err.Error())
			return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
		}
		image, err = imageService.AddImage(currentUserId, params.Image.GroupID, params)
	} else {
		//Recover image
		image, err = imageService.RecoverImage(params.Image.ImageName)
		if err != nil {
			logger.Logger().Error("Recover image error, ", err.Error())
			return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
		}
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(errors.New("The image has been created or created by someone else. Please try changing the image name.").Error()))
	}
	imageJson, err := json.Marshal(&image)

	return returnResponse(http.StatusOK, nil, operation, imageJson)
}

func DeleteImage(params image.DeleteImageParams) middleware.Responder {
	operation := "DeleteImage"
	//Params Check
	if params.ImageID <= 0 {
		err := errors.New("iamgeId is zero")
		logger.Logger().Error(err)
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	//Permission Check: image group id perimission check
	getImage, err := imageService.GetImageById(params.ImageID)
	if err != nil {
		logger.Logger().Error("Get image error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := service.CheckGroupAuth(currentUserId, getImage.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, operation, []byte("Check group auth failed"))
	}
	//Delete image
	err = imageService.DeleteImage(params)
	if err != nil {
		logger.Logger().Error("Delete image error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, operation, nil)
}

func ListImage(params image.ListImageParams) middleware.Responder {
	operation := "ListImage"
	//Params Check
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	//Permission Check: SA GA GU
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	//logger.Logger().Debugf("ListImage headers: %+v, isSa: %v\n", params.HTTPRequest.Header, isSA)
	resMap := map[string]interface{}{}
	images, count, err := imageService.ListImage(params, currentUserId, queryStr, "", isSA)
	resMap["images"] = images
	resMap["pageNumber"] = *params.Page
	resMap["total"] = count
	resMap["totalPage"] = int64(math.Ceil(float64(count) / float64(*params.Size)))
	imagesJson, err := json.Marshal(resMap)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, err, operation, imagesJson)
}

func ListImageByModelVersionId(params image.ListImageByModelVersionIDParams) middleware.Responder {
	operation := "ListImageByModelVersionId"
	//Permission Check: SA GA GU
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	//user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	images, err := imageService.ListImageByModelVersionId(params, currentUserId, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "ListImageByModelVersionId", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	imagesJson, err_ := json.Marshal(&images)
	if err_ != nil {
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, err, operation, imagesJson)
}

func UpdateImage(params image.UpdateImageParams) middleware.Responder {
	operation := "UpdateImage"
	//Permission Check:
	//Image Id & Group Id Check(GA)\UserID
	//ModelVersionID & GroupId Check
	getImage, err := imageService.GetImageById(params.ImageID)
	if err != nil {
		logger.Logger().Error("Get image error,", err)
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}

	auth := service.CheckGroupAuth(currentUserId, getImage.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, operation, []byte("Check group auth failed"))
	}
	modelVersionAuth := checkImageModelVersionAuth(getImage.GroupID, getImage.ModelVersionID, isSA)
	if !modelVersionAuth {
		return returnResponse(service.StatusAuthReject, nil, operation, []byte("Check model version auth failed"))
	}
	err = imageService.UpdateImage(params)
	if err != nil {
		logger.Logger().Error("Update image error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, operation, nil)
}

func GetImage(params image.GetImageParams) middleware.Responder {
	operation := "GetImage"
	//Params Check
	if params.ImageID <= 0 {
		err := errors.New("imageId is zero")
		logger.Logger().Error(err)
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	//Permission Check:
	//Image Id & Group Id Check(GA)\UserID
	getImage, err := imageService.GetImage(params)
	if err != nil {
		logger.Logger().Error("Get image error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := service.CheckGroupAuth(currentUserId, getImage.GroupID, isSA)
	if !auth {
		return returnResponse(service.StatusAuthReject, nil, operation, []byte("Check group auth failed"))
	}
	imageJson, err := json.Marshal(&getImage)
	if err != nil {
		logger.Logger().Error("Json marshal error, ", err.Error())
		return returnResponse(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	return returnResponse(http.StatusOK, nil, operation, imageJson)
}

//Check Image model version auth
func checkImageModelVersionAuth(groupId int64, modelVersionId int64, isSA bool) bool {
	if isSA {
		return true
	}
	modelVersion, err := modelVersionDao.GetModelVersion(modelVersionId)
	if err != nil {
		logger.Logger().Error("Get model version error:" + err.Error())
		return false
	}
	model, err := modelDao.GetModel(strconv.Itoa(int(modelVersion.ModelID)))
	if err != nil {
		logger.Logger().Error("Get model error:" + err.Error())
		return false
	}
	if model.GroupID == groupId {
		return true
	}
	return false
}
