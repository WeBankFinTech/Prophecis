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
	"strings"
)

type ImageDao struct {
}

func (s *ImageDao) AddImage(image *gormmodels.Image) error {
	err := common.GetDB().Table("t_image").Create(&image).Error
	return err
}

func (s *ImageDao) DeleteImage(imageId int64) error {
	return common.GetDB().Table("t_image").Where("id = ?", imageId).UpdateColumn("enable_flag", "0").Error
}

func (s *ImageDao) UpdateImage(imageId int64, image *gormmodels.Image) error {
	return common.GetDB().Table("t_image").Where("id = ?", imageId).Updates(&image).Error
}

func (s *ImageDao) UpdateImageMsg(imageId int64, msg string) error {
	return common.GetDB().Table("t_image").Where("id = ?", imageId).UpdateColumn("msg", msg).Error
}

func (s *ImageDao) RecoverImage(imageName string) (*gormmodels.Image, error) {
	image := &gormmodels.Image{ImageName: imageName}
	err := common.GetDB().Model(&image).Where("image_name = ?", imageName).UpdateColumn("enable_flag", 1).Error

	return image, err
}

func (s *ImageDao) GetImage(imageId int64) (*gormmodels.Image, error) {
	image := &gormmodels.Image{}
	return image, common.GetDB().Table("t_image").Where("id = ?", imageId).
		Preload("User").
		Preload("Group").
		Preload("ModelVersion").
		Find(image).Error
}

func (s *ImageDao) ListImagesCountByUser(queryStr string, cluster string, currentUserId string) (int64, error) {
	var count int64
	var err error
	if len(cluster) > 0 {
		err = common.GetDB().Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ? AND t_group.cluster_name = ? AND t_user.name = ?", "%"+queryStr+"%", 1, cluster, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Table("t_image").
			Count(&count).Error
	} else {
		err = common.GetDB().Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ? AND t_user.name = ?", "%"+queryStr+"%", 1, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Table("t_image").
			Count(&count).Error
	}
	return count, err
}

func (s *ImageDao) ListImagesByUser(size, offSet int64, queryStr string, cluster string, currentUserId string) ([]*gormmodels.Image, error) {
	var err error
	images := []*gormmodels.Image{}
	if len(cluster) > 0 {
		err = common.GetDB().Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_model.enable_flag = ? AND t_group.cluster_name = ? AND t_user.name = ?", "%"+queryStr+"%", 1, cluster, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Find(&images).Error
	} else {
		err = common.GetDB().Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ? AND t_user.name = ?", "%"+queryStr+"%", 1, currentUserId).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Find(&images).Error
	}
	return images, err
}

func (s *ImageDao) ListImagesByModelVersionId(modelVersionId int64) ([]*gormmodels.Image, error) {
	var err error
	images := []*gormmodels.Image{}
	err = common.GetDB().Where("t_image.enable_flag = ? AND t_image.model_version_id = ?", 1, modelVersionId).
		Joins("LEFT JOIN t_user " +
			"ON t_image.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_image.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		Preload("User").
		Preload("Group").
		Preload("ModelVersion").
		Find(&images).Error
	return images, err
}

func (s *ImageDao) ListImagesByUserAndModelVersionId(modelVersionId int64, currentUserId string) ([]*gormmodels.Image, error) {
	var err error
	images := []*gormmodels.Image{}
	err = common.GetDB().Where("t_image.enable_flag = ? AND t_image.model_version_id = ? AND t_user.name = ?", 1, modelVersionId, currentUserId).
		Joins("LEFT JOIN t_user " +
			"ON t_image.user_id = t_user.id " +
			"LEFT JOIN t_group " +
			"ON t_image.group_id = t_group.id ").
		Order("creation_timestamp DESC").
		Preload("User").
		Preload("Group").
		Preload("ModelVersion").
		Find(&images).Error
	return images, err
}

func (s *ImageDao) HasImageByImageName(imageName string) (bool, bool, error) {
	var err error
	var flag = false
	var dataFlag = false
	images := []*gormmodels.Image{}
	err = common.GetDB().Find(&images).Error
	for _, v := range images {
		if strings.Split(v.ImageName, ":")[1] == strings.Split(imageName, ":")[1] {
			if v.EnableFlag {
				dataFlag = true
			}
			flag = true
			break
		}
	}
	return flag, dataFlag, err
}

func (s *ImageDao) ListImages(size, offSet int64, queryStr string, cluster string) ([]*gormmodels.Image, error) {
	var err error
	images := []*gormmodels.Image{}
	if len(cluster) > 0 {
		err = common.GetDB().Table("t_image").Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ? AND t_group.cluster_name = ?", "%"+queryStr+"%", 1, cluster).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Find(&images).Error
	} else {
		err = common.GetDB().Table("t_image").Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ?", "%"+queryStr+"%", 1).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Offset(int(offSet)).
			Limit(int(size)).
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Find(&images).Error
	}
	return images, err
}

func (s *ImageDao) ListImagesCount(queryStr, cluster string) (int64, error) {
	var err error
	var count int64 = 0
	if len(cluster) > 0 {
		err = common.GetDB().Table("t_image").Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ? AND t_group.cluster_name = ?", "%"+queryStr+"%", 1, cluster).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Count(&count).Error
	} else {
		err = common.GetDB().Table("t_image").Where("CONCAT(t_image.image_name,t_user.name,t_group.name,t_image.remarks,t_model.model_name,t_user.name,t_group.name) LIKE ?"+
			" AND t_image.enable_flag = ?", "%"+queryStr+"%", 1).
			Joins("LEFT JOIN t_user " +
				"ON t_image.user_id = t_user.id " +
				"LEFT JOIN t_group " +
				"ON t_image.group_id = t_group.id " +
				"LEFT JOIN t_modelversion " +
				"ON  t_modelversion.id = t_image.model_version_id " +
				"LEFT JOIN t_model " +
				"ON t_model.id = t_modelversion.model_id ").
			Order("creation_timestamp DESC").
			Preload("User").
			Preload("Group").
			Preload("ModelVersion").
			Count(&count).Error
	}
	return count, err
}
