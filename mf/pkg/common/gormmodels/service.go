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

type Service struct {
	BaseModel
	ServiceName          string    `json:"service_name"`
	Type                 string    `json:"type"`
	Cpu                  float64   `json:"cpu"`
	Memory               float64   `json:"memory"`
	Gpu                  string    `json:"gpu"`
	Namespace            string    `json:"namespace"`
	CreationTimestamp    time.Time `json:"creation_timestamp"`
	LastUpdatedTimestamp time.Time `json:"last_updated_timestamp"`
	UserID               int64     `json:"user_id"`
	GroupID              int64     `json:"group_id"`
	ModelCount           int64     `json:"model_count"`
	LogPath              string    `json:"log_path"`
	Remark               string    `json:"remark"`
	ImageName            string    `json:"image_name"`
	Group                Group     `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	User                 User      `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`
}

type GetService struct {
	BaseModel
	ServiceName          string    `json:"service_name"`
	Type                 string    `json:"type"`
	Cpu                  float64   `json:"cpu"`
	Memory               float64   `json:"memory"`
	Gpu                  string    `json:"gpu"`
	Namespace            string    `json:"namespace"`
	CreationTimestamp    time.Time `json:"creation_timestamp"`
	LastUpdatedTimestamp time.Time `json:"last_updated_timestamp"`
	UserID               int64     `json:"user_id"`
	GroupID              int64     `json:"group_id"`
	ModelCount           int64     `json:"model_count"`
	LogPath              string    `json:"log_path"`
	Remark               string    `json:"remark"`
	Status               string    `json:"status"`
	Group                Group     `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	User                 User      `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`
}
