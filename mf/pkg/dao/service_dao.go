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

package dao

import (
	"github.com/jinzhu/gorm"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/models"
)

type SerivceDao struct {
}

func (s *SerivceDao) AddService(db *gorm.DB, service models.Service) (*models.Service, error) {
	err := db.Table("t_service").Create(&service).Error
	return &service, err
}

func (s *SerivceDao) GetService(serviceId int64) (*gormmodels.Service, error) {
	var service gormmodels.Service
	err := common.GetDB().Table("t_service").Where("t_service.id = ? AND t_service.enable_flag = ?", serviceId, 1).
		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").Find(&service).Error
	return &service, err
}

func (s *SerivceDao) GetServiceByNSAndName(namespace string, name string) (*gormmodels.Service, error) {
	var service gormmodels.Service
	err := common.GetDB().Table("t_service").Where("service_name = ? AND namespace = ? AND t_service.enable_flag"+
		" = ?", name, namespace, 1).Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").Find(&service).Error
	return &service, err
}

func (s *SerivceDao) ListServices(username string, namespace string, offset int64, size int64, queryStr string) ([]*models.GormService, error) {
	whereClause := ""
/*	if username != "" {
		whereClause = " AND t_user.name = " + "\"" + username + "\" "
	}*/
	if namespace != "" {
		whereClause = whereClause + " AND t_service.namespace = " + "\"" + namespace + "\" "
	}
	var serviceList []*models.GormService
	err := common.DB.Debug().Table("t_service").Offset(offset).Limit(size).Select("*").
		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
		Joins("LEFT JOIN t_service_modelversion ON t_service.id = t_service_modelversion.service_id ").
		Joins("LEFT JOIN t_modelversion ON t_service_modelversion.modelversion_id=t_modelversion.id ").
		Joins("LEFT JOIN t_model ON t_modelversion.model_id = t_model.id ").
		Joins("LEFT JOIN t_group ON t_group.id=t_service.group_id ").
		Order("t_service.creation_timestamp DESC").
		//Where("t_service.enable_flag = 1 AND t_modelversion.latest_flag=1 AND t_group.cluster_name = ? "+
		//	"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, clusterName, "%"+queryStr+"%").
		Where("t_service.enable_flag = 1 AND t_modelversion.latest_flag=1 "+
			"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, "%"+queryStr+"%").
		Find(&serviceList).Error
	return serviceList, err
}

func (s *SerivceDao) ListServicesByUserGroup(username string, namespace string, offset int64, size int64, queryStr string) ([]*models.GormService, error) {
	//var serviceList []*models.Service
	var serviceList []*models.GormService
	var user struct {
		Id string
	}

	err := common.DB.Table("t_user").Where("name = ?", username).First(&user).Error
	whereClause := ""
	if namespace != "" {
		whereClause = " AND t_service.namespace = " + "\"" + namespace + "\" "
	}

	if err != nil {
		return nil, err
	}
	subQuery := common.DB.Table("t_user_group").Select("t_user_group.group_id").Joins("LEFT JOIN t_user on "+
		"t_user_group.user_id=t_user.id AND t_user_group.role_id = 1").Where("t_user.name = ? AND t_user.enable_flag = 1", username).SubQuery()
	err = common.DB.Table("t_service").Offset(offset).Limit(size).Select("*").
		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
		Joins("LEFT JOIN t_service_modelversion ON t_service.id = t_service_modelversion.service_id ").
		Joins("LEFT JOIN t_modelversion ON t_service_modelversion.modelversion_id=t_modelversion.id ").
		Joins("LEFT JOIN t_model ON t_modelversion.model_id = t_model.id ").
		Joins("LEFT JOIN t_group ON t_group.id=t_service.group_id ").
		Order("t_service.creation_timestamp DESC").
		/*		Where(" (t_service.group_id in ? or t_service.user_id = ?) AND t_service.enable_flag = 1 "+
				"AND t_modelversion.latest_flag=1 AND t_group.cluster_name = ? "+
				"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, subQuery, user.Id, clusterName, "%"+queryStr+"%").*/
		Where(" (t_service.group_id in ? or t_service.user_id = ?) AND t_service.enable_flag = 1 "+
			"AND t_modelversion.latest_flag=1 "+
			"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, subQuery, user.Id, "%"+queryStr+"%").
		Find(&serviceList).Error
	return serviceList, err
}

func (s *SerivceDao) CountService(username string, namespace string, queryStr string) (int64, error) {
	var count int64
	var err error
	whereClause := ""
	if namespace != "" {
		whereClause = " AND t_service.namespace = " + "\"" + namespace + "\" "
	}
/*	subQuery := common.DB.Table("t_user_group").Select("t_user_group.group_id").Joins("LEFT JOIN t_user on "+
		"t_user_group.user_id=t_user.id AND t_user_group.role_id = 1").Where("t_user.name = ? AND t_user.enable_flag = 1", username).SubQuery()*/
	err = common.DB.Table("t_service").Select("*").
		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
		Joins("LEFT JOIN t_service_modelversion ON t_service.id = t_service_modelversion.service_id ").
		Joins("LEFT JOIN t_modelversion ON t_service_modelversion.modelversion_id=t_modelversion.id ").
		Joins("LEFT JOIN t_model ON t_modelversion.model_id = t_model.id ").
		Joins("LEFT JOIN t_group ON t_group.id=t_service.group_id ").
		Order("t_service.creation_timestamp DESC").
		/*		Where(" (t_service.group_id in ? or t_service.user_id = ?) AND t_service.enable_flag = 1 "+
				"AND t_modelversion.latest_flag=1 AND t_group.cluster_name = ? "+
				"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, subQuery, user.Id, clusterName, "%"+queryStr+"%").*/
		Where("t_service.enable_flag = 1 "+
			"AND t_modelversion.latest_flag=1 "+
			"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, "%"+queryStr+"%").
		Count(&count).Error
	return count, err
}

func (s *SerivceDao) CountServiceByUserGroup(username string, namespace string, queryStr string) (int64, error) {
	var count int64
	var user struct {
		Id string
	}
	err := common.DB.Table("t_user").Where("name = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	whereClause := ""
	if namespace != "" {
		whereClause = " AND t_service.namespace = " + "\"" + namespace + "\" "
	}
	subQuery := common.DB.Table("t_user_group").Select("t_user_group.group_id").Joins("LEFT JOIN t_user on "+
		"t_user_group.user_id=t_user.id AND t_user_group.role_id = 1").Where("t_user.name = ? AND t_user.enable_flag = 1", username).SubQuery()
	err = common.DB.Table("t_service").Select("*").
		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
		Joins("LEFT JOIN t_service_modelversion ON t_service.id = t_service_modelversion.service_id ").
		Joins("LEFT JOIN t_modelversion ON t_service_modelversion.modelversion_id=t_modelversion.id ").
		Joins("LEFT JOIN t_model ON t_modelversion.model_id = t_model.id ").
		Joins("LEFT JOIN t_group ON t_group.id=t_service.group_id ").
		Preload("User").
		Preload("Group").
		/*		Where(" (t_service.group_id in ? or t_service.user_id = ?) AND t_service.enable_flag = 1 "+
				"AND t_modelversion.latest_flag=1 AND t_group.cluster_name = ? "+
				"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, subQuery, user.Id, clusterName, "%"+queryStr+"%").*/
		Where(" (t_service.group_id in ? or t_service.user_id = ?) AND t_service.enable_flag = 1 "+
			"AND t_modelversion.latest_flag=1 "+
			"AND CONCAT(t_model.model_name,t_service.service_name, t_service.remark, t_service.type, t_user.name, t_group.name) LIKE ?"+whereClause, subQuery, user.Id, "%"+queryStr+"%").
		Count(&count).Error
	return count, err
}

//func(s *SerivceDao) CountServiceByUserName(username string,  clusterName string) (int64,error){
//	var count int64
//	err := common.DB.Table("t_service").Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
//		Where("t_user.name = ? and t_service.enable_flag = ? AND t_service.namespace like ?",
//			username , 1, "%"+clusterName+"%").Count(&count).Error
//	return count, err
//}

func (s *SerivceDao) UpdateModelsService(service *models.Service) error {
	err := common.GetDB().Model(&service).Where("id = ? AND enable_flag = ?", service.ID, 1).Update(&service).Error
	return err
}

func (s *SerivceDao) UpdateService(service *gormmodels.Service) error {
	err := common.GetDB().Model(&service).Where("id = ? AND enable_flag = ?", service.ID, 1).Update(&service).Error
	return err
}

func (s *SerivceDao) DeleteServiceById(id int64) error {
	err := common.GetDB().Table("t_service").Where("t_service.id = ?", id).UpdateColumn("enable_flag", 0).Error
	return err
}

//type SerivceDao interface {
//	addService(db *gorm.DB, service models.Service)
//	getService(serviceId string)  models.Service
//	listServices(offset int64, size int64) []*models.Service
//	countService() int64
//	updateService(db *gorm.DB, service models.Service)
//
//}

//func(s *SerivceDao) ListServicesByUserName( userName string, offset int64, size int64) ([]*models.GormService,error){
//	var serviceList []*models.GormService
//	var user struct{
//		Id string
//	}
//	err := common.DB.Table("t_user").Where("name = ?",userName).First(&user).Error
//	if err != nil {
//		return nil,err
//	}
//	err = common.DB.Debug().Table("t_service").Offset(offset).Limit(size).Select("*").
//		Joins("LEFT JOIN t_user on t_service.user_id = t_user.id").
//		Joins("LEFT JOIN t_model ON t_service.id = t_model.service_id ").
//		Joins("LEFT JOIN t_service_modelversion ON t_service.id = t_service_modelversion.service_id ").
//		Joins("LEFT JOIN t_modelversion ON t_service_modelversion.modelversion_id=t_modelversion.id ").
//		Where(" t_service.user_id = ? AND t_service.enable_flag = 1 AND t_modelversion.latest_flag=1", user.Id).
//		Find(&serviceList).Error
//	return serviceList,err
//}
