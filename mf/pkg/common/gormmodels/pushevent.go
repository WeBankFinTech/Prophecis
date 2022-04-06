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

type PushEventBase struct {
	Id                int64     `gorm:"primary_key;column:id" json:"id"`
	FileType          string    `gorm:"column:file_type;type:varchar(100)" json:"file_type"`
	FileId            int64     `gorm:"column:file_id;type:bigint(20)" json:"file_id"` // 根据FileType不同，分别对应report_version_id, model_version_id
	FileName          string    `gorm:"column:file_name;type:varchar(100)" json:"file_name"`
	FpsFileId         string    `gorm:"column:fps_file_id" json:"fps_file_id"` //FPS返回的56位文件ID
	FpsHashValue      string    `gorm:"column:fps_hash_value;type:varchar(100)" json:"fps_hash_value"`
	Params            string    `gorm:"column:params;type:varchar(200)" json:"params"`
	Idc               string    `gorm:"column:idc;type:varchar(200)" json:"idc"`
	Dcn               string    `gorm:"column:dcn;type:varchar(200)" json:"dcn"`
	Status            string    `gorm:"column:status;type:varchar(100)" json:"status"`
	CreationTimestamp time.Time `gorm:"column:creation_timestamp;type:datetime" json:"creation_timestamp"`
	UpdateTimestamp   time.Time `gorm:"column:update_timestamp;type:datetime" json:"update_timestamp"`
	EnableFlag        int8      `gorm:"column:enable_flag;type:tinyint(4)" json:"enable_flag"`
	RmbS3path         string    `gorm:"column:rmb_s3path;type:varchar(100)" json:"rmb_s3path"`
	RmbRespFileName   string    `gorm:"column:rmb_resp_file_name;type:varchar(100)" json:"rmb_resp_file_name"`
	UserId            int64     `gorm:"column:user_id;type:bigint(20)" json:"user_id"`
	Version           string    `gorm:"column:version;type:varchar(100)" json:"version"`
}

type PushEvent struct {
	PushEventBase
	//ReportVersion ReportVersionAndGroupAndReport `gorm:"ForeignKey:FileId;AssociationForeignKey:id" json:"report_version"`
	User User `gorm:"ForeignKey:UserId;AssociationForeignKey:id" json:"user"`
}

func (PushEventBase) TableName() string {
	return "pushevent"
}

func (PushEvent) TableName() string {
	return "pushevent"
}
