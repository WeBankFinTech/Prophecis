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

type ModelVersionBase struct {
	ID                int64     `gorm:"column:id; PRIMARY_KEY" json:"id"`
	EnableFlag        int8      `json:"enable_flag"`
	Version           string    `json:"version"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LatestFlag        int64     `json:"latest_flag"`
	Source            string    `json:"source"`
	Filepath          string    `json:"filepath"`
	Params            string    `json:"params"`
	FileName          string    `json:"file_name"`
	ModelID           int64     `json:"model_id"`
	TrainingId        string    `json:"training_id"`
	TrainingFlag      int8      `json:"training_flag"`
	PushTimestamp     time.Time `json:"push_timestamp"`
	PushId            int64     `json:"push_id"`
}

type ModelVersionAndModelAndEvent struct {
	ModelVersionBase
	Model Model         `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`
	Event PushEventBase `gorm:"ForeignKey:PushId;AssociationForeignKey:id" json:"event"`
	User  User
}

type ModelVersionAndModel struct {
	ModelVersionBase
	Model Model `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`
}

type ModelVersionAndEvent struct {
	ModelVersionBase
	Event PushEventBase `gorm:"ForeignKey:PushId;AssociationForeignKey:id" json:"event"`
	User  User
}

type ModelVersion struct {
	BaseModel
	Version           string    `json:"version"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	LatestFlag        int64     `json:"latest_flag"`
	Source            string    `json:"source"`
	Filepath          string    `json:"filepath"`
	Params            string    `json:"params"`
	FileName          string    `json:"file_name"`
	ModelID           int64     `json:"model_id"`
	Model             Model     `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`

	PushTimestamp time.Time `json:"push_timestamp"`
}

type ImageModelVersion struct {
	BaseModel
	Version string     `json:"version"`
	Source  string     `json:"source"`
	ModelID int64      `json:"model_id"`
	Model   ImageModel `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`
}

type ModelVersions struct {
	BaseModel
	Version string     `json:"version"`
	Source  string     `json:"source"`
	ModelID int64      `json:"model_id"`
	Model   ImageModel `gorm:"ForeignKey:ModelID;AssociationForeignKey:id" json:"model"`
}

type ModelVersionsByGroup struct {
	BaseModel
	Version string `json:"version"`
	Source  string `json:"source"`
}

type ModelLatestVersion struct {
	BaseModel
	Version           string    `json:"version"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	ModelId           int64     `json:"model_id"`
	LatestFlag        int64     `json:"latest_flag"`
	Source            string    `json:"source"`
	Filepath          string    `json:"filepath"`
	FileName          string    `json:"file_name"`

	//add by version v1.16.0
	TrainingId    string    `json:"training_id"`
	TrainingFlag  int8      `json:"training_flag"`
	PushTimestamp time.Time `json:"push_timestamp"`
	PushId        int64     `json:"push_id"`
}

func (ImageModelVersion) TableName() string {
	return "modelversion"
}

func (ModelVersions) TableName() string {
	return "modelversion"
}

func (ModelVersionsByGroup) TableName() string {
	return "modelversion"
}

func (ModelVersion) TableName() string {
	return "modelversion"
}

func (ModelLatestVersion) TableName() string {
	return "modelversion"
}

func (ModelVersionBase) TableName() string {
	return "modelversion"
}

func (ModelVersionAndModelAndEvent) TableName() string {
	return "modelversion"
}

func (ModelVersionAndEvent) TableName() string {
	return "modelversion"
}

func (ModelVersionAndModel) TableName() string {
	return "modelversion"
}
