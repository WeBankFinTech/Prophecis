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
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/models"
)

type ServiceModelVersionDao struct {
}

func (s ServiceModelVersionDao) AddServiceModelVersion(serviceModelVersion *models.GormServiceModelVersion) error {
	err := common.GetDB().Table("t_service_modelversion").Create(&serviceModelVersion).Error
	return err
}

func (s ServiceModelVersionDao) GetServiceModelVersionByServiceId(serviceId int64) (*models.GormServiceModelVersion, error) {
	var serviceModelVersions models.GormServiceModelVersion
	err := common.GetDB().Table("t_service_modelversion").Where("service_id = ? AND enable_flag = ?", serviceId, 1).Find(&serviceModelVersions).Error
	return &serviceModelVersions, err
}

func (s ServiceModelVersionDao) GetServiceModelVersionByModelVersionId(modelVersionId int64) ([]*models.GormServiceModelVersion, error) {
	var serviceModelVersions []*models.GormServiceModelVersion
	err := common.GetDB().Table("t_service_modelversion").Where("modelversion_id = ? AND enable_flag = ?", modelVersionId, 1).Find(&serviceModelVersions).Error
	return serviceModelVersions, err
}

func (s ServiceModelVersionDao) DeleteServiceModelVersion(serviceId int64) error {
	err := common.GetDB().Table("t_service_modelversion").Where("service_id = ? AND enable_flag = ?", serviceId, 1).UpdateColumn("enable_flag", "0").Error
	return err
}

func (s ServiceModelVersionDao) UpdateServiceModelVersion(serviceId int64, serviceModelVersion *models.GormServiceModelVersion) error {
	err := common.GetDB().Table("t_service_modelversion").Where("service_id = ? AND enable_flag = ?", serviceId, 1).Updates(&serviceModelVersion).Error
	return err
}
