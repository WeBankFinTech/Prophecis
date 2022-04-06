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
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/models"
	"time"

	"github.com/jinzhu/gorm"
)

type ModelVersionDao struct {
}

type ModelVersionAndModelAndEventDao struct {
}

type ModelVersionAndEventDao struct {
}

func (m *ModelVersionAndModelAndEventDao) ListModelVersionAndModelAndEvent(limit, offset int64) ([]*gormmodels.ModelVersionAndModelAndEvent, error) {
	results := make([]*gormmodels.ModelVersionAndModelAndEvent, 0)
	err := common.GetDB().Table("t_modelversion").Where("t_modelversion.enable_flag = ? ", 1).
		Order("t_modelversion.id ASC").
		Offset(offset).
		Limit(limit).
		Preload("Model").
		Preload("Event").
		Find(&results).Error
	return results, err
}

func (m *ModelVersionAndEventDao) ListModelVersionAndEventByModelId(modelId int64) ([]*gormmodels.ModelVersionAndEvent, error) {
	results := make([]*gormmodels.ModelVersionAndEvent, 0)
	err := common.GetDB().Table("t_modelversion").Where("t_modelversion.enable_flag = ? and t_modelversion.model_id = ? ", 1, modelId).
		Joins("LEFT JOIN t_pushevent " +
			"ON t_pushevent.id = t_modelversion.push_id").
		Order("t_modelversion.id ASC").
		Preload("Event").
		Find(&results).Error
	return results, err
}

func (m *ModelVersionAndModelAndEventDao) CountModelVersionAndModelAndEvent() (int64, error) {
	var count int64
	err := common.GetDB().Table("t_modelversion").Where("t_modelversion.enable_flag = ? ", 1).
		Preload("Model").
		Preload("Event").
		Count(&count).Error
	return count, err
}

func (m *ModelVersionDao) AddModelVersion(modelVersion *gormmodels.ModelVersion) error {
	return common.GetDB().Table("t_modelversion").Create(&modelVersion).Error
}

func (m *ModelVersionDao) AddModelVersionBase(modelVersion *gormmodels.ModelVersionBase) error {
	return common.GetDB().Table("t_modelversion").Create(&modelVersion).Error
}

func (m *ModelVersionDao) GetModelVersion(modelVersionId int64) (*gormmodels.ModelVersion, error) {
	var modelVersion gormmodels.ModelVersion
	err := common.GetDB().Table("t_modelversion").Preload("Model").Find(&modelVersion, "id = ? AND enable_flag = ?", modelVersionId, 1).Error
	return &modelVersion, err
}

func (m *ModelVersionDao) GetModelVersionByNameAndGroupId(modelName string, version string, groupId int64) (*gormmodels.ModelVersion, error) {
	var model gormmodels.Model
	err := common.GetDB().Debug().Table("t_model").
		Find(&model, "model_name = ? and group_id = ? and enable_flag = ?", modelName, groupId, 1).Error
	if err != nil {
		return nil, err
	}
	var modelVersion gormmodels.ModelVersion
	err = common.GetDB().Debug().Table("t_modelversion").
		Preload("Model", "id = ? ", model.ID).
		Find(&modelVersion, "model_id = ? and version = ? and  enable_flag = ?",
			model.ID, version, 1).Error
	return &modelVersion, err
}

func (m *ModelVersionDao) GetModelVersionBaseByModelIdAndVersion(modelId int64, modelVersion string) (*gormmodels.ModelVersionBase, error) {
	var modelVersionObj gormmodels.ModelVersionBase
	err := common.GetDB().Table("t_modelversion").First(&modelVersionObj, "model_id = ? AND version = ?  AND enable_flag = ?",
		modelId, modelVersion, 1).Error
	return &modelVersionObj, err
}

func (m *ModelVersionDao) GetModelVersionBaseByID(id int64) (*gormmodels.ModelVersionBase, error) {
	var modelVersionObj gormmodels.ModelVersionBase
	err := common.GetDB().Table("t_modelversion").First(&modelVersionObj, "id = ? AND enable_flag = ?",
		id, 1).Error
	return &modelVersionObj, err
}

func (m *ModelVersionDao) GetModelVersionByModelIdAndVersion(modelId int64, version string) (*gormmodels.ModelVersionBase, error) {
	var modelVersion gormmodels.ModelVersionBase
	err := common.GetDB().Table("t_modelversion").
		Where("model_id = ? AND version = ? AND enable_flag = ?", modelId, version, 1).
		First(&modelVersion).
		Error
	return &modelVersion, err
}

func (m *ModelVersionDao) UpdateModelVersionParams(modelVersionId int64, params string) error {
	err := common.GetDB().Table("t_modelversion").Where("id = ? AND enable_flag = ?", modelVersionId, 1).UpdateColumn("params", params).Error
	return err
}

func (m *ModelVersionDao) GetModelVersionByServiceId(serviceId int64) ([]*gormmodels.ModelVersion, error) {
	var modelVersion []*gormmodels.ModelVersion
	err := common.GetDB().Table("t_modelversion").Preload("Model").
		Joins("LEFT JOIN t_service_modelversion ON t_modelversion.id = t_service_modelversion.modelversion_id").
		Find(&modelVersion, "t_service_modelversion.service_id = ? AND t_modelversion.enable_flag = ?", serviceId, 1).Error
	return modelVersion, err
}

func (m *ModelVersionDao) GetModelVersionByModelId(modelId int64) (*gormmodels.ModelVersion, error) {
	var modelVersion gormmodels.ModelVersion
	err := common.GetDB().Table("t_modelversion").Where("model_id = ? AND enable_flag = ?", modelId, 1).Find(&modelVersion).Error
	return &modelVersion, err
}

func (m *ModelVersionDao) ListModelVersionBaseByModelId(modelId int64) ([]*gormmodels.ModelVersionBase, error) {
	modelVersion := make([]*gormmodels.ModelVersionBase, 0)
	err := common.GetDB().Table("t_modelversion").Where("model_id = ? AND enable_flag = ?", modelId, 1).Find(&modelVersion).Error
	return modelVersion, err
}

func (m *ModelVersionDao) CountByModelId(modelId string) (int64, error) {
	var count int64
	err := common.GetDB().Table("t_modelversion").Where("model_id = ? AND enable_flag = ?", modelId, 1).Count(&count).Error
	return count, err
}

func (m *ModelVersionDao) ListModelVersion(modelId, size, offSet int64, queryStr string) ([]*gormmodels.ModelVersionAndModel, error) {
	var modelVersionList []*gormmodels.ModelVersionAndModel
	err := common.GetDB().Table("t_modelversion").
		Where("CONCAT(t_model.model_name, t_modelversion.version) LIKE ? "+
			"AND t_modelversion.model_id = ? AND t_modelversion.enable_flag = ?", "%"+queryStr+"%", modelId, 1).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_modelversion.model_id ").
		Order("t_modelversion.id DESC").
		Offset(offSet).
		Limit(size).
		Preload("Model").
		Find(&modelVersionList).Error
	return modelVersionList, err
}

func (m *ModelVersionDao) ListModelVersionCount(modelId, queryStr string) (int64, error) {
	var count int64
	err := common.GetDB().Table("t_modelversion").
		Where("CONCAT(t_model.model_name, t_modelversion.version) LIKE ? "+
			"AND t_modelversion.model_id = ? AND t_modelversion.enable_flag = ?", "%"+queryStr+"%", modelId, 1).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_modelversion.model_id ").
		Preload("Model").
		Count(&count).Error
	return count, err
}

func (m *ModelVersionDao) DeleteModelVersion(db *gorm.DB, modelVersion models.ModelVersion) error {
	if db != nil {
		return db.Model(&modelVersion).UpdateColumn("enable_flag", "0").Error
	} else {
		return common.GetDB().Model(&modelVersion).UpdateColumn("enable_flag", "0").Error
	}
}

func (m *ModelVersionDao) UpdateModelVersionByID(id int64, u map[string]interface{}) error {
	err := common.GetDB().Table("t_modelversion").Where("id = ?", id).Update(u).Error
	return err
}

func (m *ModelVersionDao) DeleteModelVersionByModelId(modelId int64) error {
	return common.GetDB().Table("t_modelversion").Where("model_id = ?", modelId).UpdateColumn("enable_flag", "0").Error
}

func (m *ModelVersionDao) DeleteModelVersionById(modelVersionId int64) error {
	return common.GetDB().Table("t_modelversion").Where("id = ?", modelVersionId).UpdateColumn("enable_flag", "0").Error
}

func (m *ModelVersionDao) DeleteModelVersionByServiceId(serviceId int64) error {
	// Find Related Service From relation table
	subQuery := common.GetDB().Table("t_service_modelversion").Select("modelversion_id").Where("service_id = ?", serviceId).SubQuery()
	err := common.GetDB().Table("t_modelversion").Where("enable_flag = 1 and id in ?", subQuery).UpdateColumn("enable_flag", 0).Error
	return err
}

func (m *ModelVersionDao) UpdateModelVersion(db *gorm.DB, source string, filepath string, namespace string, name string, serviceId int64, modelId int64) error {
	if db == nil {
		db = common.GetDB()
	}
	//set modelversion not latest
	subQuery := db.Table("t_service_modelversion").Select("modelversion_id").Where(" service_id = ?", serviceId).SubQuery()
	err := db.Table("t_modelversion").Where("t_modelversion.id in ?", subQuery).Debug().Update("t_modelversion.latest_flag", 0).Error
	if err != nil {
		return err
	}

	// find service

	//add modelversion
	version := "v2"
	mvi := models.ModelVersion{
		Filepath:          &filepath,
		LatestFlag:        1,
		ModelID:           modelId,
		Source:            &source,
		Version:           &version,
		CreationTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		EnableFlag:        1,
	}
	err = db.Table("t_modelversion").Create(&mvi).Error
	if err != nil {
		return err
	}
	msvi := models.GormServiceModelVersion{
		ServiceId:      serviceId,
		ModelVersionId: mvi.ID,
	}
	return db.Table("t_service_modelversion").Create(&msvi).Error
}

func (m *ModelVersionDao) AddServiceModelVersion(db *gorm.DB, version *models.GormServiceModelVersion) error {
	if db != nil {
		//db.Model(&modelVersion).UpdateColumn("enable_flag","0")
		return db.Table("t_service_modelversion").Create(version).Error
	} else {
		return common.GetDB().Table("t_service_modelversion").Create(version).Error
	}
}
