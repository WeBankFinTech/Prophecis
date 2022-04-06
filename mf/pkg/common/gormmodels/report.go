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

type ReportBase struct {
	Id                int64     `gorm:"primary_key;column:id" json:"id"`
	ReportName        string    `gorm:"column:report_name" json:"report_name"`
	ModelId           int64     `gorm:"column:model_id" json:"model_id"`
	ModelVersionId    int64     `gorm:"column:model_version_id" json:"model_version_id"`
	UserId            int64     `gorm:"column:user_id" json:"user_id"`
	GroupId           int64     `gorm:"column:group_id" json:"group_id"`
	CreationTimestamp time.Time `gorm:"column:creation_timestamp;type:datetime" json:"creation_timestamp"`
	UpdateTimestamp   time.Time `gorm:"column:update_timestamp;type:datetime" json:"update_timestamp"`
	//PushStatus            string    `gorm:"column:push_status;type:varchar(100)" json:"push_status"`
	EnableFlag            int8  `gorm:"column:enable_flag;type:tinyint(4)" json:"enable_flag"`
	ReportLatestVersionId int64 `gorm:"column:report_latest_version_id" json:"report_latest_version_id"`
}

type ReportAndReportVersionBase struct {
	ReportBase
	ReportVersion ReportVersionBase `gorm:"ForeignKey:ReportLatestVersionId;AssociationForeignKey:id" json:"report_version"`
}

type Report struct {
	ReportBase
	Model         ModelBase             `gorm:"ForeignKey:ModelId;AssociationForeignKey:id" json:"model"`
	ModelVersion  ModelVersionBase      `gorm:"ForeignKey:ModelVersionId;AssociationForeignKey:id" json:"model_version"`
	Group         Group                 `gorm:"ForeignKey:GroupId;AssociationForeignKey:id" json:"group"`
	ReportVersion ReportVersionAndEvent `gorm:"ForeignKey:ReportLatestVersionId;AssociationForeignKey:id" json:"report_version"`
	User          User                  `gorm:"ForeignKey:UserId;AssociationForeignKey:id" json:"user"`
}

func (ReportBase) TableName() string {
	return "report"
}
func (ReportAndReportVersionBase) TableName() string {
	return "report"
}
func (Report) TableName() string {
	return "report"
}
