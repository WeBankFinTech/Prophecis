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

type ReportVersionBase struct {
	Id                int64     `gorm:"primary_key;column:id" json:"id"`
	ReportName        string    `gorm:"column:report_name;type:varchar(100)" json:"report_name"`
	Version           string    `gorm:"column:version;type:varchar(100)" json:"version"`
	FileName          string    `gorm:"column:file_name;type:varchar(100)" json:"file_name"`
	FilePath          string    `gorm:"column:file_path;type:varchar(100)" json:"file_path"`
	Source            string    `gorm:"column:source;type:varchar(100)" json:"source"`
	PushId            int64     `gorm:"column:push_id;type:varchar(100)" json:"push_id"`
	CreationTimestamp time.Time `gorm:"column:creation_timestamp;type:datetime" json:"creation_timestamp"`
	UpdateTimestamp   time.Time `gorm:"column:update_timestamp;type:datetime" json:"update_timestamp"`
	EnableFlag        int8      `gorm:"column:enable_flag;type:tinyint(4)" json:"enable_flag"`
	ReportId          int64     `gorm:"column:report_id;type:bigint(20)" json:"report_id"`
}

type ReportVersionAndEvent struct {
	ReportVersionBase
	Event  PushEventBase `gorm:"ForeignKey:PushId;AssociationForeignKey:id" json:"event"`
	User   User
	Group  Group
	Report ReportBase `gorm:"ForeignKey:ReportId;AssociationForeignKey:id" json:"report"`
}

type ReportVersionAndGroupAndReport struct {
	ReportVersionBase
	Group  Group
	Report ReportBase
}

func (ReportVersionAndEvent) TableName() string {
	return "reportversion"
}

func (ReportVersionBase) TableName() string {
	return "reportversion"
}

func (ReportVersionAndGroupAndReport) TableName() string {
	return "reportversion"
}
