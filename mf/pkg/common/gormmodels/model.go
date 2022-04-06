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

import (
	"time"
)

//model table struct
type ModelBase struct {
	ID                   int64     `gorm:"column:id; PRIMARY_KEY" json:"id"`
	EnableFlag           int8      `json:"enable_flag"`
	CreationTimestamp    time.Time `json:"creation_timestamp"`
	UpdateTimestamp      time.Time `json:"update_timestamp"`
	ModelLatestVersionID int64     `json:"model_latest_version_id"`
	ModelName            string    `json:"model_name"`
	ModelType            string    `json:"model_type"`
	Position             int64     `json:"position"`
	Reamrk               string    `json:"reamrk"`
	GroupID              int64     `json:"group_id"`
	UserID               int64     `json:"user_id"`
	ServiceID            int64     `json:"service_id"`
}

type Model struct {
	BaseModel
	CreationTimestamp    time.Time          `json:"creation_timestamp"`
	UpdateTimestamp      time.Time          `json:"update_timestamp"`
	ModelLatestVersionID int64              `json:"model_latest_version_id"`
	ModelName            string             `json:"model_name"`
	ModelType            string             `json:"model_type"`
	Position             int64              `json:"position"`
	Reamrk               string             `json:"reamrk"`
	GroupID              int64              `json:"group_id"`
	UserID               int64              `json:"user_id"`
	ModelLatestVersion   ModelLatestVersion `gorm:"ForeignKey:ModelLatestVersionID;AssociationForeignKey:id" json:"model_latest_version"`
	Group                Group              `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	User                 User               `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`

	ServiceID int64 `json:"service_id"`
}

type ModelVersionsModel struct {
	BaseModel
	ModelName            string               `json:"model_name"`
	ModelType            string               `json:"model_type"`
	GroupID              int64                `json:"group_id"`
	ServiceID            int64                `json:"service_id"`
	UserID               int64                `json:"user_id"`
	ModelLatestVersionID int64                `json:"model_latest_version_id"`
	Group                Group                `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	ModelVersion         ModelVersionsByGroup `gorm:"ForeignKey:ModelLatestVersionID;AssociationForeignKey:id" json:"model_version"`
}

type ImageModel struct {
	BaseModel
	ModelName         string     `json:"model_name"`
	ModelType         string     `json:"model_type"`
	GroupID           int64      `json:"group_id"`
	UserID            int64      `json:"user_id"`
	CreationTimestamp time.Time  `json:"creation_timestamp"`
	User              ImageUser  `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`
	Group             ImageGroup `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
}

type UpdateModel struct {
	BaseModel
	CreationTimestamp    time.Time `json:"creation_timestamp"`
	UpdateTimestamp      time.Time `json:"update_timestamp"`
	ModelLatestVersionID int64     `json:"model_latest_version_id"`
	ModelName            string    `json:"model_name"`
	ModelType            string    `json:"model_type"`
	Position             int64     `json:"position"`
	Reamrk               string    `json:"reamrk"`
	GroupID              int64     `json:"group_id"`
	ServiceID            int64     `json:"service_id"`
	UserID               int64     `json:"user_id"`
}

func (ImageModel) TableName() string {
	return "model"
}

func (ModelBase) TableName() string {
	return "model"
}
