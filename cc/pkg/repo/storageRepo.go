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
package repo

import (
	"github.com/jinzhu/gorm"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/models"
)

func GetAllStorageByOffSet(offset int64, size int64) ([]*models.Storage, error) {
	var storageList []*models.Storage
	err := datasource.GetDB().Offset(offset).Limit(size).Table("t_storage").Find(&storageList, "enable_flag = ?", 1).Error
	return storageList, err
}

func GetAllStorage() ([]*models.Storage, error) {
	var storageList []*models.Storage
	err := datasource.GetDB().Table("t_storage").Find(&storageList, "enable_flag = ?", 1).Error
	return storageList, err
}

func CountStorage() (int64, error) {
	var count int64
	err := datasource.GetDB().Table("t_storage").Where("enable_flag = ?", 1).Count(&count).Error
	return count, err
}

func AddStorage(storage models.Storage, db *gorm.DB) error {
	return db.Create(&storage).Error
}

func GetStorageById(storageId int64) (models.Storage, error) {
	var storage models.Storage
	err := datasource.GetDB().Find(&storage, "id = ? AND enable_flag = ?", storageId, 1).Error
	return storage, err
}

func DeleteStorageById(db *gorm.DB, storageId int64) error {
	return db.Table("t_storage").Where("id = ?", storageId).Update("enable_flag", 0).Error
}

func GetStorageByPath(path string, db *gorm.DB) (*models.Storage, error) {
	var storage models.Storage
	err := db.Find(&storage, "path = ? AND enable_flag = ?", path, 1).Error
	return &storage, err
}

func GetDeleteStorageByPath(path string, db *gorm.DB) (*models.Storage, error) {
	var storage models.Storage
	err := db.Find(&storage, "path = ? AND enable_flag = ?", path, 0).Error
	return &storage, err
}

func UpdateStorage(storage models.Storage, db *gorm.DB) error {
	return db.Model(&storage).Omit("path").Updates(&storage).Error
}
