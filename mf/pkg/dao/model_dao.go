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

	"github.com/jinzhu/copier"
)

type ModelDao struct {
}

func (m *ModelDao) AddModel(model *gormmodels.Model) error {
	err := common.GetDB().Table("t_model").Create(&model).Error
	return err
}

func (m *ModelDao) AddModelBase(model *gormmodels.ModelBase) error {
	err := common.GetDB().Table("t_model").Create(&model).Error
	return err
}

func (m *ModelDao) GetModel(modelId string) (gormmodels.Model, error) {
	var model gormmodels.Model
	err := common.GetDB().Table("t_model").
		Preload("Group").
		Preload("ModelLatestVersion").
		Find(&model, "id = ? AND enable_flag = ?", modelId, 1).
		Error
	return model, err
}

func (m *ModelDao) GetImageModel(modelId string) (gormmodels.ImageModel, error) {
	var model gormmodels.ImageModel
	err := common.GetDB().Table("t_model").
		Find(&model, "id = ? AND enable_flag = ?", modelId, 1).
		Error
	return model, err
}

func (m *ModelDao) CheckModelByModelNameAndGroupId(modelName string, groupId int64) (bool, error) {
	var count int64
	err := common.GetDB().Table("t_model").
		Where("model_name = ? AND group_id = ? AND enable_flag = ?", modelName, groupId, 1).
		Count(&count).
		Error
	return count >= 1, err
}

func (m *ModelDao) GetModelByModelNameAndGroupId(modelName string, groupId int64) (*gormmodels.ModelBase, error) {
	var result gormmodels.ModelBase
	err := common.GetDB().Table("t_model").
		Where("model_name = ? AND group_id = ? AND enable_flag = ?", modelName, groupId, 1).
		First(&result).
		Error
	return &result, err
}

func (m *ModelDao) GetModelByModelVersionId(modelVersionId int64) (*gormmodels.Model, error) {
	var model gormmodels.Model
	err := common.GetDB().Table("t_model").
		Joins("LEFT JOIN t_modelversion ON t_model.id = t_modelversion.model_id ").
		Preload("Group").
		Preload("ModelLatestVersion").
		Find(&model, "t_modelversion.id = ? AND t_modelversion.enable_flag = ?", modelVersionId, 1).
		Error
	return &model, err
}

func (m *ModelDao) GetModelBaseById(id int64) (*gormmodels.ModelBase, error) {
	var model gormmodels.ModelBase
	err := common.GetDB().Table("t_model").
		Where("t_model.id = ? AND t_model.enable_flag = ?", id, 1).
		First(&model).
		Error
	return &model, err
}

func (m *ModelDao) UpdateModel(model *gormmodels.Model) error {
	var updateModel gormmodels.UpdateModel
	err := copier.Copy(&updateModel, &model)
	if err == nil {
		err = common.GetDB().Table("t_model").Where("id = ? AND enable_flag = ?", model.BaseModel.ID, 1).Updates(&updateModel).Error
	}
	return err
}

func (m *ModelDao) ListModels(size, offSet int64, queryStr string, cluster string) ([]*gormmodels.Model, error) {
	//current database is reamrk
	var err error
	modelList := []*gormmodels.Model{}
	if len(cluster) > 0 {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_group.cluster_name = ?", "%"+queryStr+"%", 1, cluster).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Find(&modelList).Error
	} else {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ?", "%"+queryStr+"%", 1).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Find(&modelList).Error
	}
	return modelList, err
}

func (m *ModelDao) ListModelsCount(queryStr string, cluster string) (int64, error) {
	var count int64
	var err error
	if len(cluster) > 0 {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_group.cluster_name = ?", "%"+queryStr+"%", 1, cluster).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Table("t_model").
			Count(&count).Error
	} else {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ?", "%"+queryStr+"%", 1).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Table("t_model").
			Count(&count).Error
	}
	return count, err
}

func (m *ModelDao) ListModelsByUser(size, offSet int64, queryStr string, cluster string, currentUserId string) ([]*gormmodels.Model, error) {
	//current database is reamrk
	var err error
	modelList := []*gormmodels.Model{}
	if len(cluster) > 0 {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_group.cluster_name = ? AND t_user.name = ?", "%"+queryStr+"%", 1, cluster, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Find(&modelList).Error
	} else {
		err = common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_user.name = ?", "%"+queryStr+"%", 1, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Find(&modelList).Error
	}
	return modelList, err
}

func (m *ModelDao) ListModelsCountByUser(queryStr string, cluster string, currentUserId string) (int64, error) {
	var count int64
	var err error
	if len(cluster) > 0 {
		err = common.GetDB().Table("t_model").Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_group.cluster_name = ? AND t_user.name = ?", "%"+queryStr+"%", 1, cluster, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Count(&count).Error
	} else {
		err = common.GetDB().Table("t_model").Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_user.name = ?", "%"+queryStr+"%", 1, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_model.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_model.group_id = t_group.id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelLatestVersion").
			Count(&count).Error
	}
	return count, err
}

func (m *ModelDao) ListModelsByGroup(size, offSet int64, queryStr string, groupId, currentUserId string) ([]*gormmodels.ModelVersionsModel, error) {
	//current database is reamrk
	modelList := []*gormmodels.ModelVersionsModel{}
	err := common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
		" AND t_model.enable_flag = ?"+
		" AND t_model.group_id = ? AND t_user.name = ?", "%"+queryStr+"%", 1, groupId, currentUserId).
		Joins("LEFT JOIN t_user " +
			"ON t_model.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_model.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		Preload("Group").
		Preload("ModelVersion").
		Table("t_model").
		Find(&modelList).Error
	return modelList, err
}

func (m *ModelDao) ListModelsByGroupIdAndModelName(size, offSet int64, groupId, modelName, currentUserId string) ([]*gormmodels.ModelVersionsModel, error) {
	//current database is reamrk
	modelList := []*gormmodels.ModelVersionsModel{}
	err := common.GetDB().Where("t_model.enable_flag = ?"+
		" AND t_model.group_id = ? AND t_model.model_name = ? AND t_user.name = ?", 1, groupId, modelName, currentUserId).
		Joins("LEFT JOIN t_user " +
			"ON t_model.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_model.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		// Preload("User").
		Preload("Group").
		Preload("ModelVersion").
		Table("t_model").
		Find(&modelList).Error

	return modelList, err
}

func (m *ModelDao) ListModelsByGroupIdAndModelNameCount(groupId, modelName, currentUserId string) (int64, error) {
	//current database is reamrk
	var count int64 = 0
	err := common.GetDB().Where("t_model.enable_flag = ?"+
		" AND t_model.group_id = ? AND t_model.model_name = ? AND t_user.name = ?", 1, groupId, modelName, currentUserId).
		Joins("LEFT JOIN t_user " +
			"ON t_model.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_model.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		// Preload("User").
		Preload("Group").
		Preload("ModelVersion").
		Table("t_model").
		Count(&count).Error
	return count, err
}

func (m *ModelDao) ListModelsByGroupCount(queryStr string, groupId, currentUserId string) (int64, error) {
	var count int64
	err := common.GetDB().Where("CONCAT(t_model.model_name,t_model.model_type,t_model.reamrk,t_user.name,t_group.name) LIKE ?"+
		" AND t_model.enable_flag = ?"+
		" AND t_model.group_id = ? AND t_user.name = ?", "%"+queryStr+"%", 1, groupId, currentUserId).
		Joins("LEFT JOIN t_user " +
			"ON t_model.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_model.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		Preload("Group").
		Preload("ModelVersion").
		Table("t_model").
		Count(&count).Error
	return count, err
}

func (m *ModelDao) DeleteModel(model *gormmodels.Model) error {
	return common.GetDB().Table("t_model").Where("id = ?", model.ID).UpdateColumn("enable_flag", "0").Error
}

func (m *ModelDao) DeleteModelByModelVersionId(modelVersionId int64) error {
	subQuery := common.GetDB().Table("t_modelversion").Select("model_id").Where("id = ?", modelVersionId).SubQuery()
	err := common.GetDB().Table("t_model").Where("enable_flag = 1 AND id in ?", subQuery).UpdateColumn("enable_flag", "0").Error
	return err
}

func (m *ModelDao) UpdateModelById(id int64, vm map[string]interface{}) error {
	return common.GetDB().Table("t_model").Where("id = ?", id).Update(vm).Error
}
