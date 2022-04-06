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
	"fmt"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/gormmodels"
)

func CreateReport(report *gormmodels.ReportBase) error {
	return common.GetDB().Table("t_report").Create(&report).Error
}

func IsExistReportOfID(id int64) bool {
	ok := common.GetDB().Table("t_report").Where("id = ? ", id).RecordNotFound()
	if ok {
		return false
	}
	return true
}

func IsExistReportByReportNameAndGroupId(reportName string, groupId int64) bool {
	ok := common.GetDB().Table("t_report").Where("report_name = ? AND group_id = ? AND enable_flag = ?", reportName, groupId, 1).RecordNotFound()
	if ok {
		return false
	}
	return true
}

func UpdateReportById(id int64, m map[string]interface{}) error {
	return common.GetDB().Table("t_report").Where("id = ?", id).Updates(m).Error
}

func GetReportBaseByID(id int64) (*gormmodels.ReportBase, error) {
	report := gormmodels.ReportBase{}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? and t_report.id = ? ", 1, id).
		First(&report).Error
	return &report, err
}

func GetReportByID(id int64) (*gormmodels.Report, error) {
	report := gormmodels.Report{}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? and t_report.id = ? ", 1, id).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_report.model_id " +
			"LEFT JOIN t_modelversion " +
			"ON t_modelversion.id = t_report.model_version_id ").
		// Order("creation_timestamp DESC").
		Preload("Model").
		Preload("ModelVersion").
		Preload("Group").
		Preload("ReportVersion").
		Preload("User").
		First(&report).Error

	return &report, err
}

func GetReportBaseByReportNameAndGroupId(reportName string, groupId int64) (*gormmodels.ReportBase, error) {
	result := gormmodels.ReportBase{}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? and t_report.report_name = ? and t_report.group_id = ? ",
		1, reportName, groupId).
		First(&result).Error
	return &result, err
}

func ListReports(limit, offset int64, queryStr string) ([]*gormmodels.Report, error) {
	reports := []*gormmodels.Report{}
	subQuery := ""
	if queryStr != "" {
		subQuery = fmt.Sprint("  AND CONCAT(t_report.report_name) LIKE  '" +
			"%" + queryStr + "%'")
	}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? "+subQuery, 1).
		Offset(offset).Limit(limit).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_report.model_id " +
			"LEFT JOIN t_modelversion " +
			"ON t_modelversion.id = t_report.model_version_id ").
		Order("t_report.id DESC").
		Preload("Model").
		Preload("ModelVersion").
		Preload("Group").
		Preload("ReportVersion").
		Preload("User").
		Find(&reports).Error

	return reports, err
}

func CountReports(queryStr string) (int64, error) {
	var count int64 = 0
	subQuery := ""
	if queryStr != "" {
		subQuery = fmt.Sprint("  AND CONCAT(t_report.report_name) LIKE  '" +
			"%" + queryStr + "%'")
	}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? "+subQuery,
		1).
		Count(&count).Error
	return count, err
}

func ListReportsByUserId(userId, limit, offset int64, queryStr string) ([]*gormmodels.Report, error) {
	reports := []*gormmodels.Report{}
	subQuery := ""
	if queryStr != "" {
		subQuery = fmt.Sprint("  AND CONCAT(t_report.report_name) LIKE  '" +
			"%" + queryStr + "%'")
	}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ? AND t_report.user_id = ? "+subQuery,
		1, userId).
		Offset(offset).Limit(limit).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_report.model_id " +
			"LEFT JOIN t_modelversion " +
			"ON t_modelversion.id = t_report.model_version_id ").
		Order("t_report.id DESC").
		Preload("Model").
		Preload("ModelVersion").
		Preload("Group").
		Preload("ReportVersion").
		Preload("User").
		Find(&reports).Error
	return reports, err
}

func CountReportsByUserId(userId int64, queryStr string) (int64, error) {
	var count int64 = 0
	subQuery := ""
	if queryStr != "" {
		subQuery = fmt.Sprint("  AND CONCAT(t_report.report_name) LIKE  '" +
			"%" + queryStr + "%'")
	}
	err := common.GetDB().Table("t_report").Where("t_report.enable_flag = ?  AND t_report.user_id = ? "+subQuery, 1, userId).
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_report.model_id " +
			"LEFT JOIN t_modelversion " +
			"ON t_modelversion.id = t_report.model_version_id ").
		Count(&count).Error
	return count, err
}
