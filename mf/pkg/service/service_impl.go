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

package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/config"
	cc "mlss-mf/pkg/common/controlcenter"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/models"
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	v1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"mlss-mf/pkg/restapi/operations"
	"mlss-mf/pkg/restapi/operations/model_deploy"
	"mlss-mf/pkg/restapi/operations/model_storage"
	"mlss-mf/storage/client"
	"mlss-mf/storage/storage/grpc_storage"
	"os"
	"strconv"
	"time"
)

const (
	StatusBadRequest = 400
	StatusAuthReject = 401
	StatusForbidden  = 403
	StatusSuccess    = 200
	StatusNotFound   = 404
	StatusError      = 500

	RoleGA   = "SA"
	RoleSA   = "GA"
	RoleUser = "User"
)

var (
	seldonClient           *versioned.Clientset
	transactionClient      *gorm.DB
	serviceDao             *dao.SerivceDao
	modelDao               *dao.ModelDao
	userDao                *dao.UserDao
	modelversionDao        *dao.ModelVersionDao
	groupDao               *dao.GroupDao
	serviceModelVersionDao *dao.ServiceModelVersionDao
	mylogger               *logrus.Logger
	platformNamespace      string
	ccClient               *cc.CcClient
	PermissionDeniedError  = errors.New("permission denied")
)

func init() {
	seldonClient = initSeldonClient()
	//ccClient = cc.GetCcClient(viper.GetString(config.CCAddress))
	ccClient = cc.GetCcClient(config.GetMFconfig().ServerConfig.CCAddress)
	transactionClient = common.GetDB()
	serviceDao = new(dao.SerivceDao)
	modelDao = new(dao.ModelDao)
	modelversionDao = new(dao.ModelVersionDao)
	serviceModelVersionDao = new(dao.ServiceModelVersionDao)
	mylogger = logger.Logger()
	platformNamespace = config.GetMFconfig().ServerConfig.Namespace
}

type service struct {
}

func initSeldonClient() (seldonClient *versioned.Clientset) {
	seldonClient = versioned.NewForConfigOrDie(common.Restconfig)
	return seldonClient
}

func ServiceCreate(params model_deploy.PostServiceParams, currentUserId string, gpu string) (*models.Service, error) {
	//拼成seldon deployment
	sldpInstance := DefineSingleDeployment(params.Service, currentUserId)

	//save到mysql
	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		mylogger.Error("Get User Error: " + err.Error())
		return nil, err
	}

	namespaceModel, err := userDao.GetNamespaceByName(*params.Service.Namespace)
	tx := transactionClient.Begin()
	if err != nil {
		mylogger.Error("Get Namespace Error: " + err.Error())
		return nil, err
	}
	// Add Service
	paramsJson, err := json.Marshal(&params.Service.ServicePostModels[0].Parameters)
	if err != nil {
		mylogger.Error("Json Marshal Error: " + err.Error())
		return nil, err
	}
	serviceInstance := models.Service{
		CPU:                  *params.Service.CPU,
		Gpu:                  gpu,
		GroupID:              namespaceModel.GroupId,
		LogPath:              params.Service.LogPath,
		Memory:               *params.Service.Memory,
		Namespace:            *params.Service.Namespace,
		Remark:               params.Service.Remark,
		ServiceName:          *params.Service.ServiceName,
		Type:                 *params.Service.Type,
		UserID:               user.Id,
		EndpointType:         params.Service.ServicePostModels[0].EndpointType,
		CreationTimestamp:    time.Now().Format("2006-01-02 15:04:05"),
		LastUpdatedTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	si, err := serviceDao.AddService(tx, serviceInstance)
	if err != nil {
		tx.Rollback()
		mylogger.Error("Add Service Error: " + err.Error())
		return nil, err
	}

	smv := models.GormServiceModelVersion{
		ModelVersionId: params.Service.ServicePostModels[0].ModelversionID,
		ServiceId:      si.ID,
		EnableFlag:     true,
	}
	err = modelversionDao.AddServiceModelVersion(tx, &smv)
	if err != nil {
		tx.Rollback()
		mylogger.Error("Create service model version Error: " + err.Error())
		return nil, err
	}
	err = modelversionDao.UpdateModelVersionParams(smv.ModelVersionId, string(paramsJson))
	if err != nil {
		tx.Rollback()
		mylogger.Error("Update model version params error: " + err.Error())
		return nil, err
	}
	s, err := seldonClient.MachinelearningV1().SeldonDeployments(*params.Service.Namespace).Create(sldpInstance)
	if err != nil {
		tx.Rollback()
		mylogger.Error("Create CRD Error: " + err.Error())
		return nil, err
	}
	logger.Logger().Info("s" + s.Namespace)
	tx.Commit()
	ticker := time.NewTicker(time.Second * 1)
End:
	for {
		select {
		case <-ticker.C:
			sldpInstance, err = seldonClient.MachinelearningV1().SeldonDeployments(sldpInstance.Namespace).Get(*params.Service.ServiceName, metav1.GetOptions{})
			if strings.ToUpper(string(sldpInstance.Status.State)) == "STOP" || strings.ToUpper(string(sldpInstance.Status.State)) == "CREATING" {
				ticker.Stop()
				break End
			}
			if err != nil {
				logger.Logger().Error("Get seldonDeployment err, ", err)
				return nil, err
			}
			if sldpInstance.Status.State != s.Status.State {
				ticker.Stop()
				break End
			}
		}
	}
	return si, nil
}

func ServiceList(params model_deploy.ListServicesParams, currentUserId string, isSA bool) (*models.PageModel, int, error) {
	Namespace := ""
	//获取deployment crd的状态
	ns := ""
	if Namespace != "None" {
		ns = Namespace
	}
	sldpList, err := seldonClient.MachinelearningV1().SeldonDeployments(ns).List(metav1.ListOptions{})
	if err != nil {
		mylogger.Error("List Service Error: " + err.Error())
		return nil, StatusAuthReject, err
	}
	mylogger.Info("sldplist len is " + strconv.Itoa(len(sldpList.Items)))

	sldpResponseMap := map[string]interface{}{}

	for _, sldpRes := range sldpList.Items {
		mylogger.Info("namespace is : " + Namespace)
		ns = Namespace
		if ns == "None" || ns == "" {
			mylogger.Info("add sldp to map " + sldpRes.Namespace + "-" + sldpRes.Name)
			sldpResponseMap[sldpRes.Namespace+"-"+sldpRes.Name] = string(sldpRes.Status.State)
			sldpResponseMap[sldpRes.Namespace+"-"+sldpRes.Name+"-"+"image"] = string(sldpRes.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image)
		} else {
			if sldpRes.Namespace == ns {
				sldpResponseMap[sldpRes.Namespace+"-"+sldpRes.Name] = string(sldpRes.Status.State)
				sldpResponseMap[sldpRes.Namespace+"-"+sldpRes.Name+"-"+"image"] = string(sldpRes.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image)
			}
		}
	}

	//获取DB service数量
	//serviceList := serviceDao.ListServicesByUserName(user,0,100)
	offset := common.GetOffSet(*params.Page, *params.Size)
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *params.QueryStr
	}
	var count int64
	var serviceList []*models.GormService
	username := ""
	if Namespace == "" {
		username = currentUserId
	} else {
		//TODO: FIX
		username = "2"
	}

	if isSA {
		serviceList, err = serviceDao.ListServices(username, ns, offset, *params.Size, queryStr)
		for _, v := range serviceList {
			var newModelVersion models.GormModelVersion
			modelVersion, err := modelversionDao.GetModelVersion(v.ModelVersionId)
			if err != nil {
				continue
			}
			err = copier.Copy(&newModelVersion, &modelVersion)
			v.ModelVersion = newModelVersion
		}
		count, err = serviceDao.CountService(username, ns, queryStr)
	} else {
		serviceList, err = serviceDao.ListServicesByUserGroup(username, ns, offset, *params.Size, queryStr)
		for _, v := range serviceList {
			var newModelVersion models.GormModelVersion
			modelVersion, err := modelversionDao.GetModelVersion(v.ModelVersionId)
			if err != nil {
				continue
			}
			err = copier.Copy(&newModelVersion, &modelVersion)
			v.ModelVersion = newModelVersion
		}
		count, err = serviceDao.CountServiceByUserGroup(username, ns, queryStr)
	}

	if err != nil {
		mylogger.Error("List Service Error: " + err.Error())
		return nil, StatusAuthReject, err
	}

	//删除NS无关Service

	mylogger.Info("serviceList len is " + strconv.Itoa(len(serviceList)))
	for _, service := range serviceList {
		service.Status = "Stop"
		key := service.Namespace + "-" + service.ServiceName
		imageKey := service.Namespace + "-" + service.ServiceName + "-" + "image"
		if _, ok := sldpResponseMap[key]; ok {
			sldpState := sldpResponseMap[key]
			sldpImage := sldpResponseMap[imageKey]
			service.Status = sldpState.(string)
			service.Image = sldpImage.(string)
		}
	}

	var pageModel = &models.PageModel{
		Models:     serviceList,
		TotalPage:  int64(math.Ceil(float64(count) / float64(*params.Size))),
		Total:      count,
		PageNumber: *params.Page,
		PageSize:   *params.Size,
	}
	return pageModel, StatusSuccess, nil
}

//Get Service
func ServiceGet(serviceId int64, service *gormmodels.Service, newService *models.GetService) (*models.GetService, int, error) {
	seldonDeployment, err := seldonClient.MachinelearningV1().SeldonDeployments(service.Namespace).Get(service.ServiceName, metav1.GetOptions{})
	if err != nil || seldonDeployment == nil {
		logger.Logger().Error("get seldonDeployment err, ", err)
		return nil, StatusError, err
	}
	parameters := []models.ModelParameters{}
	err = copier.Copy(&parameters, seldonDeployment.Spec.Predictors[0].Graph.Parameters)
	if err != nil {
		logger.Logger().Error("Copy parameters err, ", err)
		return nil, StatusError, err
	}
	newService.Image = fmt.Sprintf("%v", seldonDeployment.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image)
	newService.Status = fmt.Sprintf("%v", seldonDeployment.Status.State)

	var parameters2 []*models.ModelParameters
	for i := 0; i < len(parameters); i++ {
		parameters2 = append(parameters2, &parameters[i])
	}

	newService.Parameters = parameters2
	newService.EndpointType = fmt.Sprintf("%v", seldonDeployment.Spec.Predictors[0].Graph.Endpoint.Type)
	modelversions, err := modelversionDao.GetModelVersionByServiceId(serviceId)
	if err != nil {
		logger.Logger().Error("Get model version by service id err, ", err)
		return nil, StatusError, err
	}
	err = copier.Copy(&newService.Modelversion, &modelversions)
	if err != nil {
		logger.Logger().Error("copy models err, ", err)
		return nil, StatusError, err
	}
	return newService, StatusSuccess, nil
}

//Get Service Dashborad
func ServiceDashborad(token string, isSA bool) (*gormmodels.Dashborad, int, error) {
	var seldonDeploymentList *v1.SeldonDeploymentList
	dashborad := gormmodels.Dashborad{
		RunningCount:   0,
		ExceptionCount: 0,
		CardCount:      0,
	}
	if isSA {
		var err error
		seldonDeploymentList, err = seldonClient.MachinelearningV1().SeldonDeployments("").List(metav1.ListOptions{})
		if err != nil {
			logger.Logger().Error("Get seldonDeployment err, ", err)
			return nil, StatusError, err
		}
	} else {
		namespace, err := ccClient.GetCurrentUserNamespaceWithRole(token, "2")
		if err != nil {
			logger.Logger().Error("Get current user namespace with role err, ", err)
			return nil, StatusError, err
		}
		seldonDeploymentList, err = seldonClient.MachinelearningV1().SeldonDeployments(string(namespace)).List(metav1.ListOptions{})
		if err != nil {
			logger.Logger().Error("Get seldonDeployment err, ", err)
			return nil, StatusError, err
		}
	}
	for _, v := range seldonDeploymentList.Items {
		//if len(strings.Split(v.Namespace, "-")) >= 3
		if strings.HasPrefix(v.Namespace,"ns-"){
			if strings.ToUpper(string(v.Status.State)) != "AVAILABLE" {
				dashborad.ExceptionCount += 1
				dashborad.CardCount = getPredictorsGPUResource(v.Spec.Predictors)
			} else {
				dashborad.RunningCount += 1
			}
		}
	}
	return &dashborad, StatusSuccess, nil
}

func getPredictorsGPUResource(predictors []v1.PredictorSpec)  int64{
	totalResource := int64(0)
	for i := 0; i < len(predictors); i++ {
		predictor := predictors[i]
		if len(predictor.ComponentSpecs) == 0 {
			break
		}
		componentSpecs := predictor.ComponentSpecs[0]
		for j := 0; j < len(componentSpecs.Spec.Containers); j++ {
			container := predictor.ComponentSpecs[0].Spec.Containers[j]
			resourceQuantity := container.Resources.Requests["nvidia.com/gpu"]
			gpuCount, _ := resourceQuantity.AsInt64()
			//gpu, _ := container.Resources.Requests["nvidia.com/gpu"].AsInt64()
			totalResource = totalResource + gpuCount
		}

	}
	return totalResource
}

//Delete Service
func ServiceDelete(serviceId int64, service *gormmodels.Service) (int, error) {
	tx := common.GetDB().Begin()
	err := serviceDao.DeleteServiceById(serviceId)
	if err != nil {
		tx.Rollback()
		logger.Logger().Error("Service Delete Error: " + err.Error())
		return StatusError, err
	}
	err = serviceModelVersionDao.DeleteServiceModelVersion(serviceId)
	if err != nil {
		tx.Rollback()
		logger.Logger().Error("Delete Service ModelVersion error, ", err)
		return StatusNotFound, err
	}
	err = seldonClient.MachinelearningV1().SeldonDeployments(service.Namespace).Delete(service.ServiceName, &metav1.DeleteOptions{})
	if err != nil {
		logger.Logger().Error("Delete Service SeldonDeployment error, ", err)
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Logger().Error("transaction commit Error: " + err.Error())
		return StatusError, err
	}
	return StatusSuccess, nil
}

func ServiceStop(params model_deploy.StopNamespacedServiceParams, user string, isSA bool) (int, error) {
	service, err := serviceDao.GetService(params.ID)
	if err != nil {
		logger.Logger().Error("Service stop err, get service from db error ", err.Error())
		return StatusError, err
	}
	if !isSA && service.User.Name != user {
		return StatusForbidden, PermissionDeniedError
	}

	sldpInstance, err := seldonClient.MachinelearningV1().SeldonDeployments(params.Namespace).Get(params.Name, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Error("Service stop err, get service from k8s error ", err.Error())
		return StatusError, err
	}
	service.ImageName = sldpInstance.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image
	err = serviceDao.UpdateService(service)
	if err != nil {
		logger.Logger().Error("Service stop err, update service error ", err.Error())
		return StatusError, err
	}
	err = seldonClient.MachinelearningV1().SeldonDeployments(params.Namespace).Delete(params.Name, &metav1.DeleteOptions{})
	if err != nil {
		logger.Logger().Error("Service stop err, delete service from k8s error  ", err.Error())
		return StatusError, err
	}
	return StatusSuccess, nil
}

func ServiceRun(params model_deploy.CreateNamespacedServiceRunParams, currentUserId string, namespace string, isSA bool) (int, error) {
	service, err := serviceDao.GetService(params.Service.ServiceID)
	if err != nil {
		return StatusError, err
	}
	if !isSA && service.User.Name != currentUserId {
		return StatusForbidden, PermissionDeniedError
	}
	params.Service.ServicePostModels[0].Image = service.ImageName
	sldpInstance := DefineSingleDeployment(params.Service, currentUserId)
	sldpInstance, err = seldonClient.MachinelearningV1().SeldonDeployments(namespace).Create(sldpInstance)
	if err != nil {
		return StatusError, err
	}
	return StatusSuccess, nil
}

//func serviceRunToServicePost(params model_deploy.CreateNamespacedServiceRunParams) (*models.ServicePost, error) {
//	servicePost := models.ServicePost{}
//
//	service := serviceDao.GetService(params.ID)
//	err := copier.Copy(servicePost, service)
//	if err != nil {
//		logger.Logger().Error("Service Delete CRD Error: " + err.Error())
//		return nil, err
//	}
//	modelVersions, err := modelversionDao.GetModelVersionByServiceId(service.ID)
//	modelVersions[0]
//
//
//	return &servicePost,err
//}

func ServiceUpdate(params model_deploy.UpdateServiceParams, gpu string, user string, isSA bool) (int, error) {
	if !isSA {
		service, err := serviceDao.GetService(params.ID)
		if err != nil {
			logger.Logger().Error("Service stop err, get service from db error ", err.Error())
			return StatusError, err
		}

		err = PermissionCheck(user, service.UserID, nil, isSA)
		if err != nil {
			logger.Logger().Errorf("Permission Check Error:" + err.Error())
			return StatusForbidden, err
		}
	}

	db := transactionClient.Begin()
	serviceModelVersion, err := serviceModelVersionDao.GetServiceModelVersionByServiceId(params.Service.ServiceID)
	if serviceModelVersion != nil && err == nil {
		serviceModelVersion.ModelVersionId = params.Service.ModelversionID
		err = serviceModelVersionDao.UpdateServiceModelVersion(params.Service.ServiceID, serviceModelVersion)
		if err != nil {
			db.Rollback()
			logger.Logger().Error("ServiceUpdate Error: " + err.Error())
			return StatusError, err
		}
	}
	serviceInstance := &models.Service{
		ID:                   params.Service.ServiceID,
		CPU:                  params.Service.CPU,
		Remark:               params.Service.Remark,
		Memory:               params.Service.Memory,
		Gpu:                  gpu,
		Namespace:            params.Service.Namespace,
		LastUpdatedTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	err = serviceDao.UpdateModelsService(serviceInstance)
	if err != nil {
		db.Rollback()
		logger.Logger().Error("ServiceUpdate Error: " + err.Error())
		return StatusError, err
	}
	sldpInstance, err := seldonClient.MachinelearningV1().SeldonDeployments(params.Service.Namespace).Get(params.Service.ServiceName, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Error("Get seldonDeployment err, ", err)
		return StatusError, err
	}
	if sldpInstance == nil {
		logger.Logger().Error("Get seldonDeployment err, sldpInstance is nil")
		return StatusError, errors.New("Get seldonDeployment err, sldpInstance is nil")
	}
	modelVersion, err := modelversionDao.GetModelVersion(params.Service.ModelversionID)
	if err != nil {
		logger.Logger().Error("Get modelVersion err,", err)
		return StatusError, err
	}
	db.Commit()
	checkResourceModelFlag := CheckResourceModel(params, sldpInstance)
	if checkResourceModelFlag == false {
		updateSldpInstance := UpdateResourceModel(params, modelVersion.Source, sldpInstance)
		_, err = seldonClient.MachinelearningV1().SeldonDeployments(sldpInstance.Namespace).Update(updateSldpInstance)
		if err != nil {
			db.Rollback()
			logger.Logger().Error("SeldonDeployments Update Error: " + err.Error())
			return StatusError, err
		}
		state := sldpInstance.Status.State
		if strings.ToUpper(string(sldpInstance.Status.State)) != "CREATING" {
			ticker := time.NewTicker(time.Second * 1)
		End:
			for {
				select {
				case <-ticker.C:
					sldpInstance, err = seldonClient.MachinelearningV1().SeldonDeployments(sldpInstance.Namespace).Get(params.Service.ServiceName, metav1.GetOptions{})
					if strings.ToUpper(string(sldpInstance.Status.State)) == "STOP" || strings.ToUpper(string(sldpInstance.Status.State)) == "CREATING" {
						ticker.Stop()
						break End
					}
					if err != nil {
						logger.Logger().Error("Get seldonDeployment err, ", err)
						return StatusError, err
					}
					if sldpInstance.Status.State != state {
						ticker.Stop()
						break End
					}
				}
			}
		}
	}
	return StatusSuccess, nil
}

func ModelAddByService(params model_deploy.PostServiceParams, groupId int64, userId int64, serviceId int64) (*models.Model, error) {
	var gormmodels gormmodels.Model
	modelInstance := models.Model{
		GroupID:   groupId,
		ModelName: *params.Service.ServiceName + "-Model",
		//ModelType:         params.Model.Models[0].ModelType,
		//TODO: FIX CHange
		ModelType:         "SKLEARN",
		Position:          0,
		Reamrk:            "",
		ServiceID:         serviceId,
		UserID:            userId,
		CreationTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		EnableFlag:        1,
	}
	err := copier.Copy(&gormmodels, &modelInstance)
	if err != nil {
		return &models.Model{}, err
	}
	err = modelDao.AddModel(&gormmodels)
	return &modelInstance, err
}

func ModelVersionAddByService(db *gorm.DB, params model_deploy.PostServiceParams, serviceId int64, modelId int64) (gormmodels.ModelVersion, error) {
	version := "v1"
	now := time.Now()
	modelversion := gormmodels.ModelVersion{
		//Filepath:           *params.Model.Models[0].Modelversions[0].Filepath,
		//Source:             *params.Model.Models[0].Modelversions[0].Source,
		Version:           version,
		ModelID:           modelId,
		LatestFlag:        1,
		CreationTimestamp: now,
		PushTimestamp:     now,
	}
	err := modelversionDao.AddModelVersion(&modelversion)

	//ServiceId & ModelVersionId Bind
	smv := models.GormServiceModelVersion{
		ServiceId:      serviceId,
		ModelVersionId: modelversion.ID,
		EnableFlag:     true,
	}
	err = modelversionDao.AddServiceModelVersion(db, &smv)
	return modelversion, err
}

//Add model
func AddModel(params model_storage.PostModelParams, currentUserId string, token string, user *models.User) (*models.PostModelResp, int, error) {
	if params.Model.S3Path != "" {
		if err := checkS3path(params.Model.S3Path); err != nil {
			return nil, 0, err
		}
	}

	var err error
	s3Path := params.Model.S3Path
	if len(params.Model.RootPath) > 0 {
		err = ccClient.UserStorageCheck(token, currentUserId, params.Model.RootPath)
		if err != nil {
			return nil, StatusAuthReject, err
		}
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%v%v", params.Model.RootPath, params.Model.ChildPath))
		if err != nil {
			return nil, StatusError, err
		}
		reader := bytes.NewReader(fileBytes)
		closer := ioutil.NopCloser(reader)
		pathSplit := strings.Split(params.Model.ChildPath, "/")
		res, status, err := UploadModel(pathSplit[len(pathSplit)-1], *params.Model.ModelType, closer)
		if err != nil || status != 200 {
			return nil, StatusError, err
		}
		s3Path = res
	}
	model := gormmodels.Model{
		GroupID:           *params.Model.GroupID,
		ModelName:         *params.Model.ModelName,
		ModelType:         *params.Model.ModelType,
		UserID:            user.Id,
		Position:          0,
		Reamrk:            "",
		CreationTimestamp: time.Now(),
		UpdateTimestamp:   time.Now(),
		BaseModel: gormmodels.BaseModel{
			EnableFlag: true,
		},
	}
	err = modelDao.AddModel(&model)
	if err != nil {
		logger.Logger().Error("Add model err, ", err)
		return nil, StatusError, err
	}
	now := time.Now()
	modelVersion := &gormmodels.ModelVersion{
		ModelID:           model.ID,
		CreationTimestamp: now,
		LatestFlag:        1,
		Version:           "v1",
		BaseModel: gormmodels.BaseModel{
			EnableFlag: true,
		},
		PushTimestamp: now,
	}
	if len(params.Model.RootPath) > 0 {
		modelVersion.Filepath = params.Model.RootPath + params.Model.ChildPath
	}
	if len(params.Model.FileName) > 0 {
		modelVersion.FileName = params.Model.FileName
	}
	if len(s3Path) > 0 {
		modelVersion.Source = params.Model.S3Path
	}
	err = modelversionDao.AddModelVersion(modelVersion)
	if err != nil {
		logger.Logger().Error("Add modelversion err, ", err)
		return nil, StatusError, err
	}
	model.ModelLatestVersionID = modelVersion.ID
	err = modelDao.UpdateModel(&model)
	if err != nil {
		logger.Logger().Error("Update model err, ", err)
		return nil, StatusError, err
	}
	return &models.PostModelResp{
		ModelID:        model.ID,
		ModelVersion:   "v1",
		ModelVersionID: modelVersion.ID,
	}, StatusSuccess, err
}

func checkS3path(s3Path string) error {
	subPath := strings.Split(strings.TrimPrefix(s3Path, "s3://"), "/")
	if subPath == nil || len(subPath) != 2 {
		return fmt.Errorf("s3path error, s3Path: %s", s3Path)
	}
	return nil
}

func AddModelVersion(params model_storage.PostModelParams, currentUserId, token string) (*models.PostModelResp, int, error) {
	if params.Model.S3Path != "" {
		if err := checkS3path(params.Model.S3Path); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	modelInfo, err := modelDao.GetModelByModelNameAndGroupId(*params.Model.ModelName, *params.Model.GroupID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if *params.Model.ModelType != modelInfo.ModelType {
		logger.Logger().Errorf("model has existed, modelType should be %q, not %q\n", modelInfo.ModelType, *params.Model.ModelType)
		err = fmt.Errorf("model has existed, modelType should be %q, not %q", modelInfo.ModelType, *params.Model.ModelType)
		return nil, http.StatusInternalServerError, err
	}

	if params.Model.ModelVersion == "" {
		params.Model.ModelVersion = "v1"
	}
	matched, err := checkVersion(params.Model.ModelVersion)
	if err != nil {
		logger.Logger().Errorf("fail to check model version, model version: %s, err: %s\n", params.Model.ModelVersion, err.Error())
		return nil, http.StatusInternalServerError, err
	}
	if !matched {
		logger.Logger().Errorf("model version not matched, model verison: %s\n", params.Model.ModelVersion)
		err = fmt.Errorf("model version not match by regexp '^[Vv][1-9][0-9]*$'")
		return nil, http.StatusInternalServerError, err
	}

	source := params.Model.S3Path
	if len(params.Model.RootPath) > 0 {
		err = ccClient.UserStorageCheck(token, currentUserId, params.Model.RootPath)
		if err != nil {
			return nil, http.StatusUnauthorized, err
		}
		fileBytes, err := ioutil.ReadFile(path.Join(params.Model.RootPath, params.Model.ChildPath))
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		reader := bytes.NewReader(fileBytes)
		closer := ioutil.NopCloser(reader)
		pathSplit := strings.Split(params.Model.ChildPath, "/")
		s3Path, status, err := UploadModel(pathSplit[len(pathSplit)-1], *params.Model.ModelType, closer)
		if err != nil || status != 200 {
			return nil, http.StatusInternalServerError, err
		}
		source = s3Path
	}

	modelVersion := params.Model.ModelVersion
	_, err = modelversionDao.GetModelVersionBaseByModelIdAndVersion(modelInfo.ID, modelVersion)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model latest version information, model id: %d, model version: %s,  err:%v\n",
				modelInfo.ID, modelVersion, err)
			return nil, http.StatusInternalServerError, err
		}
	}
	modelLatestVersionInfo, err := modelversionDao.GetModelVersionBaseByID(modelInfo.ModelLatestVersionID)
	if err != nil {
		logger.Logger().Errorf("fail to model latest version information, model latest version id: %d, err: %s\n",
			modelInfo.ModelLatestVersionID, err.Error())
		return nil, http.StatusInternalServerError, err
	}
	versionLatestBytes := []byte(modelLatestVersionInfo.Version)
	versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
	if err != nil {
		logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", modelLatestVersionInfo.Version, err.Error())
		return nil, http.StatusInternalServerError, err
	}

	versionParamBytes := []byte(params.Model.ModelVersion)
	versionParamInt64, err := strconv.ParseInt(string(versionParamBytes[1:]), 10, 64)
	if err != nil {
		logger.Logger().Errorf("fail to parse params version, version: %s, err: %s\n", params.Model.ModelVersion, err.Error())
		return nil, http.StatusInternalServerError, err
	}
	if versionLatestInt64 < versionParamInt64 {
		modelVersion = fmt.Sprintf("v%d", versionParamInt64)
	} else {
		modelVersion = fmt.Sprintf("v%d", versionLatestInt64+1)
	}
	if source == "" { // when not repeated to upload file，user latest model verison's s3Path
		source = modelLatestVersionInfo.Source
	}
	if source == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("s3Path can't be empty string")
	}

	now := time.Now()
	modelVersionBase := gormmodels.ModelVersionBase{
		EnableFlag:        1,
		Version:           modelVersion,
		CreationTimestamp: now,
		LatestFlag:        1,
		Source:            source, //TODO
		Filepath:          path.Join(params.Model.RootPath, params.Model.ChildPath),
		Params:            "",
		FileName:          params.Model.FileName,
		ModelID:           modelInfo.ID,
		TrainingId:        params.Model.TrainingID,
		TrainingFlag:      params.Model.TrainingFlag,
	}
	if err := modelVersionDao.AddModelVersionBase(&modelVersionBase); err != nil {
		logger.Logger().Errorf("fail to add model version to db, modelVersion: %+v, err:%v\n",
			modelVersion, err)
		return nil, http.StatusInternalServerError, err
	}

	m := make(map[string]interface{})
	m["model_latest_version_id"] = modelVersionBase.ID
	if err := modelDao.UpdateModelById(modelInfo.ID, m); err != nil {
		logger.Logger().Errorf("fail to update model latest version id, model version id: %d, m: %+v, err:%v\n",
			modelVersionBase.ID, m, err)
		return nil, http.StatusInternalServerError, err
	}

	return &models.PostModelResp{
		ModelID:        modelInfo.ID,
		ModelVersion:   modelVersion,
		ModelVersionID: modelVersionBase.ID,
	}, http.StatusInternalServerError, nil
}

// func CreateModel(params model_storage.PostModelParams) (int64, error) {
// 	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
// 	//token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
// 	isSA := false
// 	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
// 		isSA = true
// 	}
// 	auth := CheckGroupAuth(currentUserId, *params.Model.GroupID, isSA)
// 	if !auth {
// 		return 0, errors.New("fail to check group auth")
// 	}
// 	user, err := userDao.GetUserByName(currentUserId)
// 	if err != nil {
// 		return 0, err
// 	}

// 	version := params.Model.ModelVersion
// 	if version == "" {
// 		version = "v1"
// 	}
// 	matched, err := checkVersion(version)
// 	if err != nil {
// 		logger.Logger().Errorf("fail to check model version, model version: %s, err: %s\n", version, err.Error())
// 		return 0, err
// 	}
// 	if !matched {
// 		logger.Logger().Errorf("model version not matched, model verison: %s\n", version)
// 		err = fmt.Errorf("model version not match by regexp '^[Vv][1-9][0-9]*$'")
// 		return 0, err
// 	}

// 	modelInfo, err := modelDao.GetModelByModelNameAndGroupId(*params.Model.ModelName, *params.Model.GroupID)
// 	if err != nil {
// 		if err.Error() != gorm.ErrRecordNotFound.Error() {
// 			logger.Logger().Errorf("fail to get model information, model name: %s, group id: %d, err:%v\n",
// 				*params.Model.ModelName, *params.Model.GroupID, err)
// 			return 0, err
// 		}
// 		now := time.Now()
// 		modelVersion := gormmodels.ModelVersionBase{
// 			EnableFlag:        1,
// 			Version:           version,
// 			CreationTimestamp: now,
// 			LatestFlag:        1,
// 			Source:            *params.Model.S3Path,
// 			Filepath:          params.Model.Filepath,
// 			Params:            "",
// 			FileName:          *params.Model.FileName,
// 			ModelID:           modelInfo.ID,
// 			TrainingId:        params.Model.TrainingID,
// 			TrainingFlag:      params.Model.TrainingFlag,
// 		}
// 		if err := modelVersionDao.AddModelVersionBase(&modelVersion); err != nil {
// 			logger.Logger().Errorf("fail to add model version to db, modelVersion: %+v, err:%v\n",
// 				modelVersion, err)
// 			return 0, err
// 		}
// 		return modelInfo.ID, nil
// 	}

// 	if modelTypeFlag := CheckModelType(params.Model.ModelType); !modelTypeFlag {
// 		logger.Logger().Errorf("Model create failed, model type cannot is %s\n", params.Model.ModelType)
// 		err := errors.New("Model create failed, model type cannot is " + params.Model.ModelType)
// 		return 0, err
// 	}

// 	// model and modelverison both not exist
// 	now := time.Now()
// 	model := gormmodels.ModelBase{
// 		EnableFlag:           1,
// 		CreationTimestamp:    now,
// 		UpdateTimestamp:      now,
// 		ModelLatestVersionID: 0,
// 		ModelName:            *params.Model.ModelName,
// 		ModelType:            params.Model.ModelType,
// 		Position:             0,
// 		Reamrk:               "",
// 		GroupID:              *params.Model.GroupID,
// 		UserID:               user.Id,
// 		ServiceID:            0,
// 	}
// 	if err := modelDao.AddModelBase(&model); err != nil {
// 		logger.Logger().Errorf("fail to add model to db, model:%+v, err: %v\n", model, err)
// 		return 0, err
// 	}

// 	modelVersion := gormmodels.ModelVersionBase{
// 		EnableFlag:        1,
// 		Version:           version,
// 		CreationTimestamp: now,
// 		LatestFlag:        1,
// 		Source:            *params.Model.S3Path,
// 		Filepath:          params.Model.Filepath,
// 		Params:            "",
// 		FileName:          *params.Model.FileName,
// 		ModelID:           model.ID,
// 		TrainingId:        params.Model.TrainingID,
// 		TrainingFlag:      params.Model.TrainingFlag,
// 	}
// 	if err := modelVersionDao.AddModelVersionBase(&modelVersion); err != nil {
// 		logger.Logger().Errorf("fail to add model version to db, modelVersion: %+v, err:%v\n",
// 			modelVersion, err)
// 		return 0, err
// 	}

// 	// update model latest version id
// 	m := make(map[string]interface{})
// 	m["model_latest_version_id"] = modelVersion.ID
// 	if err := modelDao.UpdateModelById(model.ID, m); err != nil {
// 		logger.Logger().Errorf("fail to update model latest version id, model version id: %d, m: %+v, err:%v\n",
// 			modelVersion.ID, m, err)
// 		return 0, err
// 	}
// 	return model.ID, nil
// }

//Delete model
func ModelDelete(model *gormmodels.Model) (int, error) {
	err := modelDao.DeleteModel(model)
	if err != nil {
		logger.Logger().Error("Delete Model err, ", err)
		return StatusError, err
	}

	modelVersionList, err := modelversionDao.ListModelVersionBaseByModelId(model.ID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Error("fail to get model version list, model_id: %d, err: %v\n", model.ID, err)
			return StatusError, err
		}
	}
	if len(modelVersionList) > 0 {
		err = modelversionDao.DeleteModelVersionByModelId(model.ID)
		if err != nil {
			logger.Logger().Error("Delete ModelVersion err, ", err)
			return StatusError, err
		}
		m := make(map[string]interface{})
		m["enable_flag"] = 0
		for idx := range modelVersionList {
			if modelVersionList[idx] != nil {
				err = dao.UpdateEventByFileIdAndFileType(modelVersionList[idx].ID, config.Push_Event_File_Type_Model, m)
				if err != nil {
					logger.Logger().Error("fail to delete model version push event , file_id: %d,"+
						" file_type: %s, err: %v\n", modelVersionList[idx].ID, config.Push_Event_File_Type_Model, err)
					return StatusError, err
				}
			}
		}
	}

	return StatusSuccess, nil
}

//CheckServiceRunState
func CheckServiceRunState(modelId int64) bool {
	modelVersion, err := modelversionDao.GetModelVersionByModelId(modelId)
	if err != nil {
		logger.Logger().Error("Get Model error, ", err)
		return false
	}
	serviceModelVersions, err := serviceModelVersionDao.GetServiceModelVersionByModelVersionId(modelVersion.ID)
	if err != nil || len(serviceModelVersions) <= 0 {
		logger.Logger().Error("Get Model Version error, ", err)
		return true
	}
	for _, v := range serviceModelVersions {
		service, err := serviceDao.GetService(v.ServiceId)
		if err != nil {
			logger.Logger().Error("Get Service error, ", err)
			return false
		}
		sldp, err := seldonClient.MachinelearningV1().SeldonDeployments(service.Namespace).Get(service.ServiceName, metav1.GetOptions{})
		if err != nil {
			logger.Logger().Error("Get seldonDeployment error, ", err)
			return false
		}
		if len(sldp.Status.State) <= 0 || strings.ToUpper(string(sldp.Status.State)) == "STOP" {
			return false
		}
	}
	return false
}

//Get model
func ModelGet(modelId int64) (*models.Model, int, error) {
	var newModel models.Model
	model, err := modelDao.GetModel(strconv.Itoa(int(modelId)))
	if err != nil {
		logger.Logger().Error("Get model error, ", err)
		return &models.Model{}, StatusError, err
	}
	err = copier.Copy(&newModel, &model)
	if err != nil {
		logger.Logger().Error("Get model error, ", err)
		return &models.Model{}, StatusError, err
	}
	return &newModel, StatusSuccess, nil
}

//Update model
//func ModelUpdate(params model_storage.UpdateModelParams, currentUserId string, token string) (int, error) {
//	db := transactionClient.Begin()
//	var updateModel gormmodels.Model
//	getModel, err := modelDao.GetModel(strconv.Itoa(int(params.ModelID)))
//	if err != nil {
//		logger.Logger().Error("Get model err, ", err)
//		return StatusError, err
//	}
//	count, err := modelversionDao.CountByModelId(strconv.Itoa(int(params.ModelID)))
//	if err != nil {
//		logger.Logger().Error("Count by model id err, ", err)
//		db.Rollback()
//		return StatusError, err
//	}
//	now := time.Now()
//	modelVersion := &gormmodels.ModelVersion{
//		ModelID:           params.ModelID,
//		CreationTimestamp: now,
//		LatestFlag:        1,
//		Version:           fmt.Sprintf("v%v", count+1),
//		Source:            getModel.ModelLatestVersion.Source,
//		BaseModel: gormmodels.BaseModel{
//			EnableFlag: true,
//		},
//		PushTimestamp: now,
//	}
//	if len(params.Model.RootPath) > 0 {
//		err = ccClient.UserStorageCheck(token, currentUserId, params.Model.RootPath)
//		if err != nil {
//			return StatusAuthReject, err
//		}
//		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%v%v", params.Model.RootPath, params.Model.ChildPath))
//		if err != nil {
//			return StatusError, err
//		}
//		reader := bytes.NewReader(fileBytes)
//		closer := ioutil.NopCloser(reader)
//		pathSplit := strings.Split(params.Model.ChildPath, "/")
//		s3Path, status, err := UploadModel(pathSplit[len(pathSplit)-1], params.Model.ModelType, closer)
//		modelVersion.Source = s3Path
//		if err != nil || status != 200 {
//			db.Rollback()
//			return StatusError, err
//		}
//	} else {
//		if len(params.Model.S3Path) > 0 && getModel.ModelLatestVersion.Source != params.Model.S3Path {
//			modelVersion.FileName = params.Model.FileName
//			modelVersion.Source = params.Model.S3Path
//		}
//	}
//	updateModel.ID = params.ModelID
//	updateModel.UpdateTimestamp = time.Now()
//	err = modelversionDao.AddModelVersion(modelVersion)
//	if err != nil {
//		logger.Logger().Error("Add modelversion err, ", err)
//		db.Rollback()
//		return StatusError, err
//	}
//	updateModel.ModelLatestVersionID = modelVersion.ID
//	err = modelDao.UpdateModel(&updateModel)
//	if err != nil {
//		logger.Logger().Error("Update model err, ", err)
//		db.Rollback()
//		return StatusError, err
//	}
//	db.Commit()
//	return StatusSuccess, nil
//}

func ModelUpdate(params model_storage.UpdateModelParams, currentUserId string, token string) (int, error) {
	getModel, err := modelDao.GetModel(strconv.Itoa(int(params.ModelID)))
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Error("Get model err, ", err)
			return StatusError, err
		}
		logger.Logger().Debugf("model not exist, modelId: %d, err: %v\n", params.ModelID, err)
		return StatusSuccess, err
	}
	version := ""
	modelLatestVersion, err := modelversionDao.GetModelVersionBaseByID(getModel.ModelLatestVersionID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Error("Get model version err, ", err)
			return StatusError, err
		}
		version = "v1"
	} else {
		versionLatestBytes := []byte(modelLatestVersion.Version)
		versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
		if err != nil {
			logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", modelLatestVersion.Version, err.Error())
			return StatusError, err
		}
		version = fmt.Sprintf("v%d", versionLatestInt64+1)
	}

	source := params.Model.S3Path
	fileName := params.Model.FileName
	if len(params.Model.RootPath) > 0 {
		err = ccClient.UserStorageCheck(token, currentUserId, params.Model.RootPath)
		if err != nil {
			return StatusAuthReject, err
		}
		fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%v%v", params.Model.RootPath, params.Model.ChildPath))
		if err != nil {
			return StatusError, err
		}
		reader := bytes.NewReader(fileBytes)
		closer := ioutil.NopCloser(reader)
		pathSplit := strings.Split(params.Model.ChildPath, "/")
		s3Path, status, err := UploadModel(pathSplit[len(pathSplit)-1], params.Model.ModelType, closer)
		if err != nil || status != 200 {
			logger.Logger().Errorf("fail to upload model, fileName: %s, err: %v\n", pathSplit[len(pathSplit)-1], err)
			return StatusError, err
		}
		source = s3Path
		fileName = pathSplit[len(pathSplit)-1]
	}
	if source == "" {
		logger.Logger().Errorln("s3path can't be empty string")
		return StatusError, errors.New("s3path can't be empty string")
	}
	if fileName == "" {
		logger.Logger().Errorln("fileName can't be empty string")
		return StatusError, errors.New("fileName can't be empty string")
	}

	now := time.Now()
	modelVersion := gormmodels.ModelVersionBase{
		EnableFlag:        1,
		Version:           version,
		CreationTimestamp: now,
		LatestFlag:        1,
		Source:            source,
		Filepath:          path.Join(params.Model.RootPath, params.Model.ChildPath),
		Params:            "",
		FileName:          fileName,
		ModelID:           params.ModelID,
		TrainingId:        "",
		TrainingFlag:      0,
		PushId:            0,
	}
	if err := modelVersionDao.AddModelVersionBase(&modelVersion); err != nil {
		logger.Logger().Errorf("fail to add modelversion , modelversion: %+v,  err: %v\n", modelVersion, err)
		return StatusError, errors.New("fail to add modelversion, " + err.Error())
	}

	m := make(map[string]interface{})
	m["model_latest_version_id"] = modelVersion.ID
	m["update_timestamp"] = now
	if err := modelDao.UpdateModelById(params.ModelID, m); err != nil {
		logger.Logger().Error("Update model err, ", err)
		return StatusError, err
	}
	return StatusSuccess, nil
}

//Get model list
func ModelList(currentUserId string, size, page int64, queryStr string, isSA bool) ([]*gormmodels.Model, int64, int, error) {
	var modelList []*gormmodels.Model
	var count int64
	var err error
	if isSA {
		modelList, err = modelDao.ListModels(size, (page-1)*size, queryStr, "")
		count, err = modelDao.ListModelsCount(queryStr, "")
	} else {
		modelList, err = modelDao.ListModelsByUser(size, (page-1)*size, queryStr, "", currentUserId)
		count, err = modelDao.ListModelsCountByUser(queryStr, "", currentUserId)
	}
	if err != nil {
		return []*gormmodels.Model{}, 0, StatusError, err
	}
	return modelList, count, StatusSuccess, err
}

//Get model list by cluster
func ModelListByCluster(currentUserId string, size, page int64, queryStr string, cluster string, isSA bool) ([]*gormmodels.Model, int64, int, error) {
	var modelList []*gormmodels.Model
	var count int64
	var err error
	if isSA {
		modelList, err = modelDao.ListModels(size, (page-1)*size, queryStr, cluster)
		count, err = modelDao.ListModelsCount(queryStr, cluster)
	} else {
		modelList, err = modelDao.ListModelsByUser(size, (page-1)*size, queryStr, cluster, currentUserId)
		count, err = modelDao.ListModelsCountByUser(queryStr, cluster, currentUserId)
	}
	if err != nil {
		return []*gormmodels.Model{}, 0, StatusError, err
	}
	return modelList, count, StatusSuccess, err
}

//Get model version list
func ModelVersionList(modelId, size, page int64, queryStr string, user string, isSA bool) (models.GetModelVersionResp, int64, int, error) {

	if !isSA {
		model, err := modelDao.GetModel(strconv.FormatInt(modelId, 10))
		if err != nil {
			logger.Logger().Error("ModelVersionList error", err.Error())
			return nil, 0, StatusError, err
		}
		err = PermissionCheck(user, model.UserID, nil, isSA)
		if err != nil {
			logger.Logger().Errorf("Permission Check Error:" + err.Error())
			return nil, 0, StatusForbidden, PermissionDeniedError
		}
	}

	modelversionList, err := modelversionDao.ListModelVersion(modelId, size, (page-1)*size, queryStr)
	for idx := range modelversionList {
		user, err := userDao.GetUserByUserId(modelversionList[idx].Model.UserID)
		if err != nil {
			logger.Logger().Error("Get model user failed, ", err)
			return nil, 0, StatusError, err
		}
		group, err := groupDao.GetGroupByGroupId(modelversionList[idx].Model.GroupID)
		if err != nil {
			logger.Logger().Error("Get image group failed, ", err)
			return nil, 0, StatusError, err
		}
		modelversionList[idx].Model.User = *user
		modelversionList[idx].Model.Group = *group
	}
	if err != nil {
		logger.Logger().Error("Get list model version failed, ", err)
		return nil, 0, StatusError, err
	}
	count, err := modelversionDao.ListModelVersionCount(strconv.Itoa(int(modelId)), queryStr)
	if err != nil {
		logger.Logger().Error("Get list model version count failed, ", err)
		return nil, 0, StatusError, err
	}
	respData := generateGetModelversionResp(modelversionList)
	return respData, count, StatusSuccess, err
}

func generateGetModelversionResp(versions []*gormmodels.ModelVersionAndModel) models.GetModelVersionResp {
	if len(versions) == 0 {
		return nil
	}
	resp := make(models.GetModelVersionResp, 0)
	for k := range versions {
		group := models.Group{
			ClusterName:    versions[k].Model.Group.ClusterName,
			DepartmentID:   versions[k].Model.Group.DepartmentId,
			DepartmentName: versions[k].Model.Group.DepartmentName,
			GroupType:      versions[k].Model.Group.GroupType,
			ID:             versions[k].Model.Group.ID,
			Name:           versions[k].Model.Group.Name,
			Remarks:        versions[k].Model.Group.Remarks,
			RmbDcn:         versions[k].Model.Group.RmbDcn,
			RmbIdc:         versions[k].Model.Group.RmbIdc,
			ServiceID:      versions[k].Model.Group.ServiceId,
			SubsystemID:    versions[k].Model.Group.SubsystemId,
			SubsystemName:  versions[k].Model.Group.SubsystemName,
		}
		var enableFlag int8 = 0
		if versions[k].Model.EnableFlag {
			enableFlag = 1
		}
		model := models.ModelBase{
			CreationTimestamp:    versions[k].Model.CreationTimestamp.Format(config.TimeFormat),
			EnableFlag:           enableFlag,
			GroupID:              versions[k].Model.GroupID,
			ID:                   versions[k].Model.ID,
			ModelLatestVersionID: versions[k].Model.ModelLatestVersionID,
			ModelName:            versions[k].Model.ModelName,
			ModelType:            versions[k].Model.ModelType,
			Position:             versions[k].Model.Position,
			Reamrk:               versions[k].Model.Reamrk,
			ServiceID:            versions[k].Model.ServiceID,
			UpdateTimestamp:      versions[k].Model.UpdateTimestamp.Format(config.TimeFormat),
			UserID:               versions[k].Model.UserID,
		}
		enableFlag = 0
		if versions[k].Model.User.EnableFlag {
			enableFlag = 1
		}
		user := models.UserInfo{
			EnableFlag: enableFlag,
			Gid:        versions[k].Model.User.Gid,
			GUIDCheck:  int8(versions[k].Model.User.GuidCheck),
			ID:         versions[k].Model.User.ID,
			Name:       versions[k].Model.User.Name,
			Remarks:    versions[k].Model.User.Remarks,
			UID:        versions[k].Model.User.UID,
			UserType:   versions[k].Model.User.Type,
		}
		modelVersion := models.GetModelVersionRespBase{
			CreationTimestamp: versions[k].CreationTimestamp.Format(config.TimeFormat),
			EnableFlag:        versions[k].EnableFlag,
			FileName:          versions[k].FileName,
			Filepath:          versions[k].Filepath,
			Group:             &group,
			ID:                versions[k].ID,
			LatestFlag:        versions[k].LatestFlag,
			Model:             &model,
			ModelID:           versions[k].ModelID,
			Params:            versions[k].Params,
			PushID:            versions[k].PushId,
			PushTimestamp:     versions[k].PushTimestamp.Format(config.TimeFormat),
			Source:            versions[k].Source,
			TrainingFlag:      versions[k].TrainingFlag,
			TrainingID:        versions[k].TrainingId,
			User:              &user,
			Version:           versions[k].Version,
		}
		resp = append(resp, &modelVersion)
	}
	return resp
}

//Bind service and model version
func ModelVersionServiceCreate(s *models.Service, m *models.ModelVersion) bool {
	serviceModelVersion := &models.GormServiceModelVersion{
		ServiceId:      s.ID,
		ModelVersionId: m.ID,
		EnableFlag:     true,
	}
	err := serviceModelVersionDao.AddServiceModelVersion(serviceModelVersion)
	if err != nil {
		return false
	}
	return true
}

//Get model list by group
func ModelListByGroups(size, page int64, groupId, queryStr, currentUserId string) ([]*gormmodels.ModelVersionsModel, int64, int, error) {
	modelList, err := modelDao.ListModelsByGroup(size, (page-1)*size, queryStr, groupId, currentUserId)
	if err != nil {
		logger.Logger().Error("Get list models by group failed, ", err)
		return []*gormmodels.ModelVersionsModel{}, 0, StatusError, err
	}
	count, err := modelDao.ListModelsByGroupCount(queryStr, groupId, currentUserId)
	if err != nil {
		logger.Logger().Error("Get list models by group count failed, ", err)
		return []*gormmodels.ModelVersionsModel{}, 0, StatusError, err
	}
	return modelList, count, StatusSuccess, err
}

//Get model list by group id and model name
func ModelListByGroupIdAndModelName(size, currentPage int64, groupId, modelName, currentUserId string) ([]*gormmodels.ModelVersionsModel, int64, error) {
	modelList, err := modelDao.ListModelsByGroupIdAndModelName(size, (currentPage-1)*size,
		groupId, modelName, currentUserId)
	if err != nil {
		logger.Logger().Errorf("fail to get model by group id and model name, %s", err.Error())
		return modelList, 0, err
	}

	count, err := modelDao.ListModelsByGroupIdAndModelNameCount(groupId, modelName, currentUserId)
	if err != nil {
		logger.Logger().Errorf("fail to get total of model by group id and model name, %s", err.Error())
		return modelList, 0, err
	}
	return modelList, count, nil
}

//Upload code
func UploadModel(fileName, modelType string, closer io.ReadCloser) (string, int, error) {
	uid := uuid.New().String()
	hostPath := fmt.Sprintf("/data/oss-storage/model/%v/", uid)
	err := os.MkdirAll(hostPath, os.ModePerm)
	if err != nil {
		logger.Logger().Error("Mkdir err, ", err.Error())
		return "", StatusError, errors.New("Mkdir err, " + err.Error())
	}

	filePath := hostPath + fileName
	logger.Logger().Infof("lk-test UploadModel filepath: %s", filePath)
	fileBytes, err := ioutil.ReadAll(closer)
	if err != nil {
		logger.Logger().Error("File read err, ", err.Error())
		return "", StatusError, errors.New("File read err," + err.Error())
	}

	//err = ioutil.WriteFile(filePath, fileBytes, os.ModePerm)
	err = WriteToFile(filePath, config.Max_RW, fileBytes)
	if err != nil {
		logger.Logger().Error("File write err, ", err.Error())
		return "", StatusError, errors.New("File write err," + err.Error())
	}
	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return "", StatusError, errors.New("Storer client create err," + err.Error())
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		return "", StatusError, errors.New("Storer client is nil")
	}
	req := &grpc_storage.CodeUploadRequest{
		Bucket:   "mlss-mf",
		FileName: fileName,
		HostPath: hostPath,
	}
	//var res *grpc_storage.CodeUploadResponse
	//if strings.HasSuffix(fileName, ".zip") && modelType != CUSTOM {
	//	res, err = storageClient.Client().UploadModelZip(context.TODO(), req)
	//	if err != nil {
	//		logger.Logger().Error("Storer client upload model err, ", err.Error())
	//		return "", StatusError, errors.New("Storer client upload model err," + err.Error())
	//	}
	//} else {
	//	res, err = storageClient.Client().UploadCode(context.TODO(), req)
	//	if err != nil {
	//		logger.Logger().Error("Storer client upload model err, ", err.Error())
	//		return "", StatusError, errors.New("Storer client upload model err," + err.Error())
	//	}
	//}

	res, err := storageClient.Client().UploadCode(context.TODO(), req)
	if err != nil {
		logger.Logger().Error("Storer client upload model err, ", err.Error())
		return "", StatusError, errors.New("Storer client upload model err," + err.Error())
	}
	//TODO clear folder
	err = os.RemoveAll(hostPath)
	if err != nil {
		logger.Logger().Error("clear folder err, " + err.Error())
		return "", StatusError, err
	}
	return res.S3Path, StatusSuccess, nil
}

func WriteToFile(fileName string, maxWriteNum int, fileBytes []byte) error {
	if fileName == "" {
		return errors.New("file name cam't be empty")
	}
	if maxWriteNum <= 0 {
		maxWriteNum = config.Max_RW
	}
	if len(fileBytes) == 0 {
		return nil
	}
	f, err := os.Create(fileName)
	if err != nil {
		f.Close()
		logger.Logger().Error("fail to create file, file: %s, err: %s\n", fileName, err.Error())
		return err
	}
	if len(fileBytes) <= maxWriteNum {
		_, err = f.Write(fileBytes)
		if err != nil {
			f.Close()
			logger.Logger().Error("fail to write file, file: %s, err: %s\n", fileName, err.Error())
			return err
		}
	} else {
		for len(fileBytes) > 0 {
			b := fileBytes
			if len(b) > maxWriteNum {
				b = fileBytes[:maxWriteNum]
			}
			n, err := f.Write(b)
			if err != nil {
				f.Close()
				logger.Logger().Error("fail to write file, file: %s, err: %s\n", fileName, err.Error())
				return err
			}
			fileBytes = fileBytes[n:]
		}
	}
	f.Close()
	return nil
}

//Export model
func ExportModel(currentUserId string, modelId int64, isSA bool) ([]byte, int, error) {
	model, err := modelDao.GetModel(strconv.Itoa(int(modelId)))
	if err != nil {
		return nil, StatusError, err
	}
	auth := CheckGroupAuth(currentUserId, model.GroupID, isSA)
	if !auth {
		return nil, StatusAuthReject, nil
	}
	storageClient, err := client.NewStorage()
	if err != nil {
		return nil, StatusError, err
	}
	req := &grpc_storage.CodeDownloadRequest{
		S3Path: model.ModelLatestVersion.Source,
	}
	res, err := storageClient.Client().DownloadCode(context.TODO(), req)
	if err != nil {
		return nil, StatusError, err
	}
	if res == nil {
		return nil, StatusError, errors.New("res is nil")
	}
	logger.Logger().Info("Read path：", fmt.Sprintf("%v/%v", res.HostPath, res.FileName))
	fileBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", res.HostPath, res.FileName))
	if err != nil {
		logger.Logger().Errorf("fail to read file %q, err: %s\n", path.Join(res.HostPath, res.FileName), err.Error())
		return nil, StatusError, err
	}
	if err := os.Remove(res.HostPath); err != nil {
		logger.Logger().Errorf("fail to remove folder %q, err: %s\n", res.HostPath, err.Error())
		return nil, StatusError, err
	}
	return fileBytes, StatusSuccess, nil
}

//Check auth
func CheckAuth(namespace string, currentUserId string, token string, isSA bool) bool {
	if isSA {
		return true
	}
	err := ccClient.CheckNamespaceUser(token, namespace, currentUserId)
	if err != nil {
		mylogger.Error("Check Auth-CheckNamespaceUser error:" + err.Error())
		return false
	} else {
		return true
	}
}

//Check User Group auth
func CheckUserAuth(currentUserId string, groupId int64, userId int64, isSA bool) bool {
	if isSA {
		return true
	}
	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		mylogger.Error("Get userId error:" + err.Error())
		return false
	}

	if userId == user.Id {
		return true
	}
	groups, err := groupDao.GetGroupsByUserId(user.Id)
	for _, v := range groups {
		if v.GroupID == groupId {
			return true
		}
	}
	return false
}

//Check Group auth
func CheckGroupAuth(currentUserId string, groupId int64, isSA bool) bool {
	if isSA {
		return true
	}
	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		mylogger.Error("Get userId error:" + err.Error())
		return false
	}
	groups, err := groupDao.GetGroupsByUserId(user.Id)
	for _, v := range groups {
		if v.GroupID == groupId {
			return true
		}
	}
	return false
}

//Check Model type
func CheckModelType(modelType string) bool {
	if strings.ToUpper(modelType) == "SKLEARN" ||
		strings.ToUpper(modelType) == "TENSORFLOW" ||
		strings.ToUpper(modelType) == "XGBOOST" ||
		strings.ToUpper(modelType) == "CUSTOM" || strings.ToUpper(modelType) == "MLPIPELINE" {
		return true
	}
	return false
}

//Check Resource model
func CheckResourceModel(params model_deploy.UpdateServiceParams, sldpInstance *v1.SeldonDeployment) bool {
	cpuQuantity := resource.MustParse(fmt.Sprintf("%.1f", params.Service.CPU))
	memoryString := fmt.Sprintf("%1.f", params.Service.Memory)
	memoryQuantity := resource.MustParse(fmt.Sprintf("%vGi", memoryString))
	cpuFlag := cpuQuantity.String() == sldpInstance.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Requests.Cpu().String()
	gpuFlag := params.Service.Gpu == "0"
	memoryFlag := memoryQuantity.String() == sldpInstance.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources.Requests.Memory().String()
	sourceFlag := params.Service.Source == sldpInstance.Spec.Predictors[0].Graph.ModelURI
	logger.Logger().Info(fmt.Sprintf("cpuFlag: %v, gpuFlag: %v, memoryFlag: %v, sourceFlag: %v", cpuFlag, gpuFlag, memoryFlag, sourceFlag))
	if cpuFlag && gpuFlag && memoryFlag && sourceFlag {
		return true
	}
	return false
}

func DownloadModel(modelID, currentUserId string, isSA bool) ([]byte, int, string, error) {
	modelDao := &dao.ModelDao{}
	modelInfo, err := modelDao.GetModel(modelID)
	if err != nil {
		logger.Logger().Errorf("fail to get model information, error: %s", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	auth := CheckGroupAuth(currentUserId, modelInfo.GroupID, isSA)
	if !auth {
		return nil, StatusAuthReject, "", nil
	}

	logger.Logger().Infof("lk-test modelInfo, modelInfo: %+v", modelInfo)
	//tmpDir := "/data/oss-storage/model/" + uuid.New().String()
	//if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
	//	logger.Logger().Errorf("fail to create tmp dir, error: %s", err.Error())
	//	return nil, http.StatusInternalServerError, "", err
	//}

	s3Path := modelInfo.ModelLatestVersion.Source
	if len(s3Path) < 5 {
		logger.Logger().Errorf("s3Path length not correct, s3Path: %s", s3Path)
		return nil, http.StatusInternalServerError, "", fmt.Errorf("s3Path length not correct, s3Path: %s", s3Path)
	}
	s3Array := strings.Split(s3Path[5:], "/")
	logger.Logger().Debugf("lk-test s3Path: %v", s3Path)
	if len(s3Array) < 2 {
		logger.Logger().Errorf("s3 Path is not correct, s3Array: %+v", s3Array)
		return nil, http.StatusInternalServerError, "", errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]

	if bucketName != "mlss-mf" {
		err := fmt.Errorf("bucket name is not 'mlss-mf', bucket name: %s", bucketName)
		return nil, http.StatusInternalServerError, "", err
	}

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		err = fmt.Errorf("Storer client is nil ")
		return nil, http.StatusInternalServerError, "", err
	}

	//if is single file, download
	//bytes := make([]byte, 0)
	//switch strings.Contains(modelInfo.ModelLatestVersion.FileName, ".zip") {
	//case false:
	//	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	//	s3Path = s3Path + "/" + modelInfo.ModelLatestVersion.FileName //单个文件
	//	// 1. download file
	//	in := grpc_storage.CodeDownloadRequest{
	//		S3Path: s3Path,
	//	}
	//	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
	//		return nil, http.StatusInternalServerError, err
	//	}
	//
	//	// 2. read file body
	//	bytes, err = ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
	//		return nil, http.StatusInternalServerError, err
	//	}
	//
	//	err = os.RemoveAll(out.HostPath)
	//	if err != nil {
	//		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
	//		return nil, http.StatusInternalServerError, err
	//	}
	//default:
	//	in := grpc_storage.CodeDownloadByDirRequest{
	//		S3Path: s3Path + "/" + modelInfo.ModelLatestVersion.FileName,
	//	}
	//	outs, err := storageClient.Client().DownloadCodeByDir(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code by dir, error: %s", err.Error())
	//		return nil, http.StatusInternalServerError, err
	//	}
	//	if outs != nil && len(outs.RespList) > 0 {
	//		for _, resp := range outs.RespList {
	//			logger.Logger().Debugf("download code by dir, hostpath: %s, filename: %s\n", resp.HostPath, resp.FileName)
	//			bytes, err := ioutil.ReadFile(resp.HostPath + "/" + resp.FileName)
	//			if err != nil {
	//				logger.Logger().Errorf("fail to read file, filename: %s", resp.HostPath+"/"+resp.FileName)
	//				return nil, http.StatusInternalServerError, err
	//			}
	//
	//			// write to tmp file
	//			// todo
	//			logger.Logger().Debugf("read file %s, content-length: %d\n", resp.HostPath+"/"+resp.FileName, len(bytes))
	//			subPath := strings.Split(resp.FileName, "/")
	//			filename := path.Join(subPath[1:]...)
	//			tmpFile := path.Join(tmpDir, filename)
	//			if err := createDir(tmpFile); err != nil {
	//				logger.Logger().Errorf("fail to create tmp directory of file %q, error: %s", tmpFile, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//			err = ioutil.WriteFile(tmpFile, bytes, os.ModePerm)
	//			// logger.Logger().Debugf("file %q content: %s\n", resp.FileName, string(bytes))
	//			if err != nil {
	//				logger.Logger().Errorf("fail to write data to file %q, error: %s", tmpFile, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//		}
	//		for _, resp := range outs.RespList {
	//			err = os.RemoveAll(resp.HostPath)
	//			if err != nil {
	//				logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", resp.HostPath, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//		}
	//
	//		// gzip file
	//		dest := ""
	//		if modelInfo.ModelLatestVersion.FileName != "" {
	//			filenameOnly := getFileNameOnly(modelInfo.ModelLatestVersion.FileName)
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", filenameOnly))
	//		} else {
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", s3Array[len(s3Array)-1]))
	//		}
	//		source := []string{tmpDir}
	//
	//		if err := archiver.Archive(source, dest); err != nil {
	//			logger.Logger().Errorf("fail to archive file, error: %s, source: %s, dest: %s\n", err.Error(), source, dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//
	//		// read from zip file
	//		bytes, err = ioutil.ReadFile(dest)
	//		if err != nil {
	//			logger.Logger().Errorf("fail to read zip file , file: %s", dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//		logger.Logger().Debugf("zip file %q content: %s\n", dest, string(bytes))
	//		if err := os.Remove(dest); err != nil {
	//			logger.Logger().Errorf("fail to remove zip file , file: %s", dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//	}
	//}

	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	s3Path = s3Path + "/" + modelInfo.ModelLatestVersion.FileName //单个文件
	// 1. download file
	in := grpc_storage.CodeDownloadRequest{
		S3Path: s3Path,
	}
	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	if err != nil {
		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	// 2. read file body
	bytes, err := ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	if err != nil {
		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
		return nil, http.StatusInternalServerError, "", err
	}

	err = os.RemoveAll(out.HostPath)
	if err != nil {
		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	//// remove host tmp file
	//err = os.RemoveAll(tmpDir)
	//if err != nil {
	//	logger.Logger().Errorf("clear folder failed, dir name: %s, error: %s", tmpDir, err)
	//	return nil, http.StatusInternalServerError, "", err
	//}

	return bytes, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(modelInfo.ModelLatestVersion.FileName)), nil
}

func DownloadModelVersion(modelVersionId int64, currentUserId string, isSA bool) ([]byte, int, string, error) {
	modelVersionInDb, err := modelversionDao.GetModelVersion(modelVersionId)
	if err != nil {
		logger.Logger().Errorf("fail to get model version information, modelVersionId: %d, error: %s",
			modelVersionId, err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	auth := CheckGroupAuth(currentUserId, modelVersionInDb.Model.GroupID, isSA)
	if !auth {
		return nil, StatusAuthReject, "", nil
	}

	//tmpDir := "/data/oss-storage/modelversion/" + uuid.New().String()
	//if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
	//	logger.Logger().Errorf("fail to create tmp dir, error: %s", err.Error())
	//	return nil, http.StatusInternalServerError, "", err
	//}

	s3Path := modelVersionInDb.Source
	if len(s3Path) < 5 {
		logger.Logger().Errorf("s3Path length not correct, s3Path: %s", s3Path)
		return nil, http.StatusInternalServerError, "", fmt.Errorf("s3Path length not correct, s3Path: %s", s3Path)
	}
	s3Array := strings.Split(s3Path[5:], "/")
	logger.Logger().Debugf("lk-test s3Path: %v", s3Path)
	if len(s3Array) < 2 {
		logger.Logger().Errorf("s3 Path is not correct, s3Array: %+v", s3Array)
		return nil, http.StatusInternalServerError, "", errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]

	if bucketName != "mlss-mf" {
		err := fmt.Errorf("bucket name is not 'mlss-mf', bucket name: %s", bucketName)
		return nil, http.StatusInternalServerError, "", err
	}

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		err = fmt.Errorf("Storer client is nil ")
		return nil, http.StatusInternalServerError, "", err
	}

	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	s3Path = s3Path + "/" + modelVersionInDb.FileName
	// 1. download file
	in := grpc_storage.CodeDownloadRequest{
		S3Path: s3Path,
	}
	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	if err != nil {
		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	// 2. read file body
	bytes, err := ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	if err != nil {
		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
		return nil, http.StatusInternalServerError, "", err
	}

	err = os.RemoveAll(out.HostPath)
	if err != nil {
		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	//// remove host tmp file
	//err = os.RemoveAll(tmpDir)
	//if err != nil {
	//	logger.Logger().Errorf("clear folder failed, dir name: %s, error: %s", tmpDir, err)
	//	return nil, http.StatusInternalServerError, "", err
	//}

	return bytes, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(modelVersionInDb.FileName)), nil
}

func getFileNameOnly(filename string) string {
	filenameWithSuffix := path.Base(filename)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly
}

func createDir(pathStr string) error {
	dir := path.Dir(pathStr)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsExist(err) {
			// fmt.Printf("dir exist")
			return nil
		} else {
			// fmt.Printf("file not exist")
			os.IsNotExist(err)

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkModelPushEventParams(params *models.PushModelRequest) error {
	if params == nil {
		return errors.New("params can't be nil")
	}
	if params.FactoryName == nil || *params.FactoryName == "" {
		return errors.New("params factory name can't be empty")
	}

	if params.ModelType == nil {
		return errors.New("params model type can't be nil")
	}
	if *params.ModelType != config.RMB_ModelType_Logistic_Regression &&
		*params.ModelType != config.RMB_ModelType_Decision_Tree &&
		*params.ModelType != config.RMB_ModelType_Random_Forest &&
		*params.ModelType != config.RMB_ModelType_XGBoost &&
		*params.ModelType != config.RMB_ModelType_LightGBM {
		return fmt.Errorf("params model type not exist, model_type: %s", *params.ModelType)
	}

	if params.ModelUsage == nil {
		return errors.New("params model usage can't be nil")
	}
	if *params.ModelUsage != config.RMB_ModelUsage_Classification &&
		*params.ModelUsage != config.RMB_ModelUsage_Regression {
		return fmt.Errorf("params model usage not exist, model_usage: %s", *params.ModelUsage)
	}
	return nil
}

func PushModelByModelID(modelId int64, modelPushEventParam *models.PushModelRequest) (*models.Event, error) {
	if err := checkModelPushEventParams(modelPushEventParam); err != nil {
		return nil, err
	}

	modelIdStr := strconv.FormatInt(modelId, 10)
	modelInDb, err := modelDao.GetModel(modelIdStr)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model information, modelId: %d, err: %s\n", modelId, err.Error())
			return nil, err
		}
		return nil, nil
	}

	params, err := json.Marshal(modelPushEventParam)
	if err != nil {
		logger.Logger().Errorf("fail to marshal model event params, params: %+v, err: %s\n", *modelPushEventParam, err.Error())
		return nil, err
	}

	eventVersion := ""
	latestPushEvent, err := dao.GetEventById(modelInDb.ModelLatestVersion.PushId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get latest push event information, modelVersionID: %d, eventId: %d,  err: %s\n",
				modelInDb.ModelLatestVersionID, modelInDb.ModelLatestVersion.PushId, err.Error())
			return nil, err
		}
		eventVersion = "v1"
	} else {
		if latestPushEvent.Version != "" {
			versionLatestBytes := []byte(latestPushEvent.Version)
			versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
			if err != nil {
				logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", latestPushEvent.Version, err.Error())
				return nil, err
			}
			eventVersion = fmt.Sprintf("v%d", versionLatestInt64+1)
		} else {
			eventVersion = "v1"
		}
	}

	now := time.Now()
	event := gormmodels.PushEventBase{
		FileType:          "MODEL",
		FileId:            modelInDb.ModelLatestVersion.ID,
		FileName:          modelInDb.ModelLatestVersion.FileName,
		FpsFileId:         "",
		FpsHashValue:      "",
		Params:            string(params),
		Idc:               modelInDb.Group.RmbIdc,
		Dcn:               modelInDb.Group.RmbDcn,
		Status:            config.RMB_PushEvent_Status_Creating,
		CreationTimestamp: now,
		UpdateTimestamp:   now,
		EnableFlag:        1,
		RmbRespFileName:   "",
		RmbS3path:         "",
		UserId:            modelInDb.UserID,
		Version:           eventVersion,
	}
	if err := dao.AddEvent(&event); err != nil {
		logger.Logger().Errorf("fail to add event, event: %+v, err: %s\n", event, err.Error())
		return nil, err
	}

	go pushModel(modelInDb.ModelLatestVersion.Source, modelInDb.ModelLatestVersion.FileName, modelInDb.Group.RmbDcn, modelInDb.Group.ServiceId,
		event.Id, modelPushEventParam)

	//update model version push id
	m := make(map[string]interface{})
	m["push_id"] = event.Id
	err = modelVersionDao.UpdateModelVersionByID(modelInDb.ModelLatestVersion.ID, m)
	if err != nil {
		logger.Logger().Errorf("fail to update model version's push_id, modelVersion id: %d, m: %+v, err: %s\n",
			modelInDb.ModelLatestVersion.ID, m, err.Error())
		return nil, err
	}

	eventInfo := models.Event{
		CreationTimestamp: (*strfmt.DateTime)(&event.CreationTimestamp),
		Dcn:               &event.Dcn,
		EnableFlag:        &event.EnableFlag,
		FileID:            &event.FileId,
		FileName:          &event.FileName,
		FileType:          &event.FileType,
		FpsFileID:         &event.FpsFileId,
		HashValue:         &event.FpsHashValue,
		ID:                &event.Id,
		Idc:               &event.Idc,
		Params:            &event.Params,
		RmbS3path:         event.RmbS3path,
		RmbRespFileName:   event.RmbRespFileName,
		Status:            &event.Status,
		UpdateTimestamp:   (*strfmt.DateTime)(&event.UpdateTimestamp),
		Version:           eventVersion,
	}

	return &eventInfo, nil
}

func PushModelByModelVersionID(modelVersionId int64, modelPushEventParam *models.PushModelRequest, user string, isSA bool) (*models.Event, error) {
	if err := checkModelPushEventParams(modelPushEventParam); err != nil {
		return nil, err
	}

	modelVersionInDb, err := modelversionDao.GetModelVersionBaseByID(modelVersionId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model version information, modelVersionId: %d, err: %s\n", modelVersionId, err.Error())
			return nil, err
		}
		logger.Logger().Infof("record not exist, modelVersionId: %d, err: %s\n", modelVersionId, err.Error())
		return nil, nil
	}
	modelInDb, err := modelDao.GetModelBaseById(modelVersionInDb.ModelID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model information, modelId: %d, err: %s\n", modelVersionInDb.ModelID, err.Error())
			return nil, err
		}
		logger.Logger().Infof("record not exist, modelId: %d, err: %s\n", modelVersionInDb.ModelID, err.Error())
		return nil, nil
	}

	err = PermissionCheck(user, modelInDb.UserID, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	groupInDb, err := groupDao.GetGroupByGroupId(modelInDb.GroupID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get group information, groupId: %d, err: %s\n", modelInDb.GroupID, err.Error())
			return nil, err
		}
		logger.Logger().Infof("record not exist, groupId: %d, err: %s\n", modelInDb.GroupID, err.Error())
		return nil, nil
	}

	params, err := json.Marshal(modelPushEventParam)
	if err != nil {
		logger.Logger().Errorf("fail to marshal model event params, params: %+v, err: %s\n", *modelPushEventParam, err.Error())
		return nil, err
	}

	eventVersion := ""
	latestPushEvent, err := dao.GetEventById(modelVersionInDb.PushId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get latest push event information, modelVersionId: %d, pushEventId: %d,  err: %s\n",
				modelVersionId, modelVersionInDb.PushId, err.Error())
			return nil, err
		}
		eventVersion = "v1"
	} else {
		if latestPushEvent.Version != "" {
			versionLatestBytes := []byte(latestPushEvent.Version)
			versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
			if err != nil {
				logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", latestPushEvent.Version, err.Error())
				return nil, err
			}
			eventVersion = fmt.Sprintf("v%d", versionLatestInt64+1)
		} else {
			eventVersion = "v1"
		}
	}

	now := time.Now()
	event := gormmodels.PushEventBase{
		FileType:          "MODEL",
		FileId:            modelVersionInDb.ID,
		FileName:          modelVersionInDb.FileName,
		FpsFileId:         "",
		FpsHashValue:      "",
		Params:            string(params),
		Idc:               groupInDb.RmbIdc,
		Dcn:               groupInDb.RmbDcn,
		Status:            config.RMB_PushEvent_Status_Creating,
		CreationTimestamp: now,
		UpdateTimestamp:   now,
		EnableFlag:        1,
		RmbRespFileName:   "",
		RmbS3path:         "",
		UserId:            modelInDb.UserID,
		Version:           eventVersion,
	}
	if err := dao.AddEvent(&event); err != nil {
		logger.Logger().Errorf("fail to add event, event: %+v, err: %s\n", event, err.Error())
		return nil, err
	}
	go pushModel(modelVersionInDb.Source, modelVersionInDb.FileName, groupInDb.RmbDcn, groupInDb.ServiceId, event.Id, modelPushEventParam)

	//update model version push id
	m := make(map[string]interface{})
	m["push_id"] = event.Id
	err = modelVersionDao.UpdateModelVersionByID(modelVersionInDb.ID, m)
	if err != nil {
		logger.Logger().Errorf("fail to update model version's push_id, modelVersion id: %d, m: %+v, err: %s\n",
			modelVersionInDb.ID, m, err.Error())
		return nil, err
	}

	eventInfo := models.Event{
		CreationTimestamp: (*strfmt.DateTime)(&event.CreationTimestamp),
		Dcn:               &event.Dcn,
		EnableFlag:        &event.EnableFlag,
		FileID:            &event.FileId,
		FileName:          &event.FileName,
		FileType:          &event.FileType,
		FpsFileID:         &event.FpsFileId,
		HashValue:         &event.FpsHashValue,
		ID:                &event.Id,
		Idc:               &event.Idc,
		Params:            &event.Params,
		RmbS3path:         event.RmbS3path,
		RmbRespFileName:   event.RmbRespFileName,
		Status:            &event.Status,
		UpdateTimestamp:   (*strfmt.DateTime)(&event.UpdateTimestamp),
		Version:           eventVersion,
	}

	return &eventInfo, nil
}

func pushModel(s3path, filename, rmbDcn string, serviceId, eventId int64, modelPushEventParam *models.PushModelRequest) error {
	return nil
}

func ListModelVersionPushEventsByModelVersionID(modelVersinId, currentPage, pageSize int64, user string, isSA bool) (*models.ModelVersionPushEventResp, error) {
	defaultParams := model_storage.NewListModelVersionPushEventsByModelVersionIDParams()
	if currentPage <= 0 {
		currentPage = *defaultParams.CurrentPage
	}
	if pageSize <= 0 {
		pageSize = *defaultParams.PageSize
	}
	modelVersionInfo, err := modelversionDao.GetModelVersionBaseByID(modelVersinId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model verison information, modelVersinId: %d, error: %s\n",
				modelVersinId, err.Error())
			return nil, err
		}
		return nil, nil
	}
	modelInfo, err := modelDao.GetModelBaseById(modelVersionInfo.ModelID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model information, modelId: %d, error: %s\n",
				modelVersionInfo.ModelID, err.Error())
			return nil, err
		}
		return nil, nil
	}
	err = PermissionCheck(user, modelInfo.UserID, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}
	eventAndUser, err := dao.ListModelPushEventByModelVersionId(modelVersinId, (currentPage-1)*pageSize, pageSize)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get model verison event and user information, modelVersinId: %d, error: %s\n",
				modelVersinId, err.Error())
			return nil, err
		}
		return nil, nil
	}
	count, err := dao.CountModelPushEventByModelVersionId(modelVersinId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to count model verison event, modelVersinId: %d, error: %s\n",
				modelVersinId, err.Error())
			return nil, err
		}
		return nil, nil
	}

	resp := generateModelVersionPushEventResp(count, modelVersionInfo, modelInfo, eventAndUser)
	return resp, nil
}

func generateModelVersionPushEventResp(count int64, modelVersion *gormmodels.ModelVersionBase,
	model *gormmodels.ModelBase, pushEvents []*gormmodels.PushEvent) *models.ModelVersionPushEventResp {
	events := make([]*models.ModelVersionPushEventRespBase, 0)
	if len(pushEvents) > 0 {
		for idx := range pushEvents {
			logger.Logger().Debugf("generateModelVersionPushEventResp eventsAndUser[%d]: %+v\n", idx, *pushEvents[idx])
			if pushEvents[idx] != nil {
				eventInfo := models.Event{
					CreationTimestamp: (*strfmt.DateTime)(&pushEvents[idx].CreationTimestamp),
					Dcn:               &pushEvents[idx].Dcn,
					EnableFlag:        &pushEvents[idx].EnableFlag,
					FileID:            &pushEvents[idx].FileId,
					FileName:          &pushEvents[idx].FileName,
					FileType:          &pushEvents[idx].FileType,
					FpsFileID:         &pushEvents[idx].FpsFileId,
					HashValue:         &pushEvents[idx].FpsHashValue,
					ID:                &pushEvents[idx].Id,
					Idc:               &pushEvents[idx].Idc,
					Params:            &pushEvents[idx].Params,
					RmbRespFileName:   pushEvents[idx].RmbRespFileName,
					RmbS3path:         pushEvents[idx].RmbS3path,
					Status:            &pushEvents[idx].Status,
					UpdateTimestamp:   (*strfmt.DateTime)(&pushEvents[idx].UpdateTimestamp),
					Version:           pushEvents[idx].Version,
				}
				var enableFlag int8 = 0
				if pushEvents[idx].User.EnableFlag {
					enableFlag = 1
				}
				userInfo := models.UserInfo{
					EnableFlag: enableFlag,
					Gid:        pushEvents[idx].User.Gid,
					GUIDCheck:  int8(pushEvents[idx].User.GuidCheck),
					ID:         pushEvents[idx].User.ID,
					Name:       pushEvents[idx].User.Name,
					Remarks:    pushEvents[idx].User.Remarks,
					UID:        pushEvents[idx].User.UID,
					UserType:   pushEvents[idx].User.Type,
				}
				event := models.ModelVersionPushEventRespBase{
					Event: &eventInfo,
					User:  &userInfo,
				}
				events = append(events, &event)
			}
		}
	} else {
		logger.Logger().Debugln("generateModelVersionPushEventResp pushEvents is nil")
	}
	modelInfo := models.ModelBase{
		CreationTimestamp:    model.CreationTimestamp.Format(config.TimeFormat),
		EnableFlag:           model.EnableFlag,
		GroupID:              model.GroupID,
		ID:                   model.ID,
		ModelLatestVersionID: model.ModelLatestVersionID,
		ModelName:            model.ModelName,
		ModelType:            model.ModelType,
		Position:             model.Position,
		Reamrk:               model.Reamrk,
		ServiceID:            model.ServiceID,
		UpdateTimestamp:      model.UpdateTimestamp.Format(config.TimeFormat),
		UserID:               model.UserID,
	}
	modelVersionInfo := models.ModelVersionBase{
		CreationTimestamp: modelVersion.CreationTimestamp.Format(config.TimeFormat),
		EnableFlag:        modelVersion.EnableFlag,
		FileName:          modelVersion.FileName,
		Filepath:          modelVersion.Filepath,
		ID:                modelVersion.ID,
		LatestFlag:        modelVersion.LatestFlag,
		ModelID:           modelVersion.ModelID,
		PushID:            modelVersion.PushId,
		PushTimestamp:     modelVersion.PushTimestamp.Format(config.TimeFormat),
		Source:            modelVersion.Source,
		TrainingFlag:      modelVersion.TrainingFlag,
		TrainingID:        modelVersion.TrainingId,
		Version:           modelVersion.Version,
	}
	resp := models.ModelVersionPushEventResp{
		Events:       events,
		Model:        &modelInfo,
		ModelVersion: &modelVersionInfo,
		Total:        count,
	}
	return &resp
}

func GetModelVersionByNameAndGroupID(params model_storage.GetModelVersionByNameAndVersionParams) (*models.ModelVersion, int, error) {
	var modelVersion models.ModelVersion
	gormModelVersion, err :=
		modelVersionDao.GetModelVersionByNameAndGroupId(params.ModelName, params.Version, params.GroupID)
	if err != nil {
		logger.Logger().Error("Get model error, ", err.Error())
		return nil, StatusError, err
	}
	err = copier.Copy(&modelVersion, &gormModelVersion)
	if err != nil {
		logger.Logger().Error("Get model error, ", err)
		return &models.ModelVersion{}, StatusError, err
	}
	return &modelVersion, StatusSuccess, nil
}

func PermissionCheck(execUserName string, entityUserId int64, groupId *string, isSuperAdmin bool) error {
	if isSuperAdmin {
		return nil
	}
	user, err := userDao.GetUserByName(execUserName)
	if err != nil {
		return errors.New("Permission Check Error, Get user err: " + err.Error())
	}
	if user.Id != entityUserId {
		return PermissionDeniedError
	}
	return nil
}
