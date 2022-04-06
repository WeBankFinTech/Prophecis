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

type Image struct {
	BaseModel
	ImageName            string            `json:"image_name"`
	ModelVersionID       int64             `json:"model_version_id"`
	UserID               int64             `json:"user_id"`
	GroupID              int64             `json:"group_id"`
	CreationTimestamp    time.Time         `json:"creation_timestamp"`
	LastUpdatedTimestamp time.Time         `json:"last_updated_timestamp"`
	Status               string            `json:"status"`
	Remarks              string            `json:"remarks"`
	Msg                  string            `json:"msg"`
	Group                ImageGroup        `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	ModelVersion         ImageModelVersion `gorm:"ForeignKey:ModelVersionID;AssociationForeignKey:id" json:"model_version"`
	User                 ImageUser         `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`
}
