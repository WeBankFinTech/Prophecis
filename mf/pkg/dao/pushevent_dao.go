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

package dao

import (
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/gormmodels"
)

func GetEventById(id int64) (*gormmodels.PushEventBase, error) {
	event := gormmodels.PushEventBase{}
	err := common.GetDB().Table("t_pushevent").Where("id = ? ", id).First(&event).Error
	return &event, err
}

func AddEvent(event *gormmodels.PushEventBase) error {
	err := common.GetDB().Table("t_pushevent").Create(event).Error
	return err
}

func UpdateEventById(id int64, m map[string]interface{}) error {
	err := common.GetDB().Table("t_pushevent").Where("id = ? AND enable_flag = ?", id, 1).Update(m).Error
	return err
}

func UpdateEventByFileIdAndFileType(fileId int64, filtType string, m map[string]interface{}) error {
	err := common.GetDB().Table("t_pushevent").Where("file_id = ? AND file_type = ? AND enable_flag = ?",
		fileId, filtType, 1).Update(m).Error
	return err
}

const (
	ReportPushEventFileType = "DATA"
	ModelPushEventFileType  = "MODEL"
)

func ListReportPushEventByReportVersionId(reportVersionId, offSet, size int64) ([]*gormmodels.PushEvent, error) {
	results := make([]*gormmodels.PushEvent, 0)
	err := common.GetDB().Table("t_pushevent").Where("t_pushevent.file_id = ? AND t_pushevent.file_type = ?  AND t_pushevent.enable_flag = ?",
		reportVersionId, ReportPushEventFileType, 1).
		Joins(" LEFT JOIN t_user " +
			" ON t_user.id = t_pushevent.user_id").
		Order("id ASC").
		Offset(offSet).
		Limit(size).
		Preload("User").
		Find(&results).Error
	return results, err
}

func CountReportPushEventByReportVersionId(reportVersionId int64) (int64, error) {
	var count int64
	err := common.GetDB().Table("t_pushevent").Where("t_pushevent.file_id = ? AND t_pushevent.file_type = ?  AND t_pushevent.enable_flag = ?",
		reportVersionId, ReportPushEventFileType, 1).
		Joins(" LEFT JOIN t_user " +
			" ON t_user.id = t_pushevent.user_id").
		Preload("User").
		Count(&count).Error
	return count, err
}

func ListModelPushEventByModelVersionId(modelVersionId, offSet, size int64) ([]*gormmodels.PushEvent, error) {
	results := make([]*gormmodels.PushEvent, 0)
	err := common.GetDB().Table("t_pushevent").Where("t_pushevent.file_id = ? AND t_pushevent.file_type = ?  AND t_pushevent.enable_flag = ?",
		modelVersionId, ModelPushEventFileType, 1).
		Joins("LEFT JOIN t_user " +
			"ON t_user.id = t_pushevent.user_id").
		Order("id ASC").
		Offset(offSet).
		Limit(size).
		Preload("User").
		Find(&results).Error

	return results, err
}

func CountModelPushEventByModelVersionId(modelVersionId int64) (int64, error) {
	var count int64
	err := common.GetDB().Table("t_pushevent").Where("t_pushevent.file_id = ? AND t_pushevent.file_type = ?  AND t_pushevent.enable_flag = ?",
		modelVersionId, ModelPushEventFileType, 1).
		Joins(" LEFT JOIN t_user " +
			" ON t_user.id = t_pushevent.user_id").
		Preload("User").Count(&count).Error
	return count, err
}
