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
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
)

func AddKey(key models.Keypair) (models.Keypair, error) {
	err := repo.AddKey(key)
	if err != nil {
		logger.Logger().Error("Add key err, ", err)
		return models.Keypair{}, err
	}
	return repo.GetKeyByName(key.Name)
}

func UpdateKey(key models.Keypair) (models.Keypair, error) {
	err := repo.UpdateKey(key)
	if err != nil {
		logger.Logger().Error("Update key err, ", err)
		return models.Keypair{}, err
	}
	return repo.GetKeyById(key.ID)
}

func GetKeyByName(name string) (models.Keypair, error) {
	return repo.GetKeyByName(name)
}

func GetKeyByApiKey(apiKey string) (models.Keypair, error) {
	return repo.GetKeyByApiKey(apiKey)
}

func GetDeleteKeyByName(name string) (models.Keypair, error) {
	return repo.GetDeleteKeyByName(name)
}

func DeleteByName(name string) (models.Keypair, error) {
	err := repo.DeleteByName(name)
	if err != nil {
		logger.Logger().Error("DeleteByName err, ", err)
		return models.Keypair{}, err
	}
	return repo.GetDeleteKeyByName(name)
}
