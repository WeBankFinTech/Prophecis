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
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/models"
)

func AddKey(key models.Keypair) error {
	err := datasource.GetDB().Create(&key).Error
	return err
}

func UpdateKey(key models.Keypair) error {
	err := datasource.GetDB().Model(&key).Omit("name").Update(&key).Error
	return err
}

func DeleteByName(name string) error {
	err := datasource.GetDB().Table("t_keypair").Where("name = ?", name).Update("enable_flag", 0).Error
	return err
}

func GetKeyById(id int64) (models.Keypair, error) {
	var key models.Keypair
	err := datasource.GetDB().Find(&key, "id = ? AND enable_flag = ?", id, 1).Error
	return key, err
}

func GetKeyByName(name string) (models.Keypair, error) {
	var key models.Keypair
	err := datasource.GetDB().Find(&key, "name = ? AND enable_flag = ?", name, 1).Error
	return key, err
}

func GetDeleteKeyByName(name string) (models.Keypair, error) {
	var key models.Keypair
	err := datasource.GetDB().Find(&key, "name = ? AND enable_flag = ?", name, 0).Error
	return key, err
}

func GetKeyByApiKey(apiKey string) (models.Keypair, error) {
	var key models.Keypair
	err := datasource.GetDB().Find(&key, "api_key = ? AND enable_flag = ?", apiKey, 1).Error
	return key, err
}
