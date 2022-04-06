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
package gormmodels

import "time"

type Result struct {
	Id                int64            `gorm:"primary_key;column:id" json:"id"`
	ModelId           int64            `gorm:"column:model_id" json:"model_id"`
	ModelVersionId    int64            `gorm:"column:model_version_id" json:"model_version_id"`
	TrainingId        string           `gorm:"column:training_id;type:varchar(100)" json:"training_id"`
	ResultMsg         string           `gorm:"column:result_msg;type:varchar(200)" json:"result_msg"`
	EnableFlag        int8             `gorm:"column:enable_flag;type:tinyint(4)" json:"enable_flag"`
	CreationTimestamp time.Time        `gorm:"column:creation_timestamp;type:datetime" json:"creation_timestamp"`
	UpdateTimestamp   time.Time        `gorm:"column:update_timestamp;type:datetime" json:"update_timestamp"`
	Model             ModelBase        `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`
	ModelVersion      ModelVersionBase `gorm:"ForeignKey:ModelVersionId;AssociationForeignKey:id" json:"model_varsion"`
}
