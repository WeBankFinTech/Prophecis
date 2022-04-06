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

// func GetReportVersionByReportNameAndReportVersion(reportName, reportVersion string) (*gormmodels.ReportVersionAndReportBase, error) {
// 	result := gormmodels.ReportVersionAndReportBase{}
// 	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ?  AND t_reportversion.report_name = ? "+
// 		" AND t_reportversion.version = ? ", 1, reportName, reportVersion).
// 		Preload("Report").First(&result).Error

// 	return &result, err
// }

func GetReportVersionByReportIdAndReportVersion(reportId int64, reportVersion string) (*gormmodels.ReportVersionBase, error) {
	result := gormmodels.ReportVersionBase{}
	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ?  AND t_reportversion.report_id = ? "+
		" AND t_reportversion.version = ? ", 1, reportId, reportVersion).First(&result).Error
	return &result, err
}

func GetReportVersionByReportVersionId(reportVersionId int64) (*gormmodels.ReportVersionBase, error) {
	result := gormmodels.ReportVersionBase{}
	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ?  AND t_reportversion.id = ? ",
		1, reportVersionId).
		First(&result).Error
	return &result, err
}

func CreateReportVersion(reportVersion *gormmodels.ReportVersionBase) error {
	return common.GetDB().Table("t_reportversion").Create(reportVersion).Error
}

func UpdateReportVersionById(id int64, m map[string]interface{}) error {
	return common.GetDB().Table("t_reportversion").Where("t_reportversion.id = ?", id).Update(m).Error
}

func UpdateReportVersionByReportId(reportId int64, m map[string]interface{}) error {
	return common.GetDB().Table("t_reportversion").Where("t_reportversion.report_id = ?", reportId).Update(m).Error
}

func ListReportVersionsByReportID(reportId int64) ([]gormmodels.ReportVersionAndEvent, error) {
	results := make([]gormmodels.ReportVersionAndEvent, 0)

	//reportVersionBaseSlice := make([]gormmodels.ReportVersionBase, 0)
	//err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ?  "+
	//	" AND t_reportversion.report_id = ? ", 1, reportId).
	//	Find(&reportVersionBaseSlice).Error
	//if err != nil {
	//	if err.Error() == gorm.ErrRecordNotFound.Error() {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//
	//reportInfo := gormmodels.ReportBase{}
	//err = common.GetDB().Table("t_report").Where("t_report.enable_flag = ?  "+
	//	" AND t_report.id = ? ", 1, reportId).First(&reportInfo).Error
	//if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
	//	return nil, err
	//}
	//
	//if len(reportVersionBaseSlice) > 0 {
	//	for idx := range reportVersionBaseSlice {
	//		eventInfo := gormmodels.PushEventBase{}
	//		if reportVersionBaseSlice[idx].PushId > 0 {
	//			err = common.GetDB().Table("t_pushevent").Where("t_pushevent.enable_flag = ? "+
	//				" AND t_pushevent.id = ? ", 1, reportVersionBaseSlice[idx].PushId).First(&eventInfo).Error
	//			if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
	//				return nil, err
	//			}
	//		}
	//
	//		result := gormmodels.ReportVersionAndEvent{
	//			ReportVersionBase: reportVersionBaseSlice[idx],
	//			Event:             eventInfo,
	//			Report:            reportInfo,
	//		}
	//		results = append(results, result)
	//	}
	//}

	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ? AND "+
		" t_reportversion.report_id = ?", 1, reportId).
		Joins("LEFT JOIN t_report " +
			"ON t_report.id = t_reportversion.report_id " +
			"LEFT JOIN t_pushevent " +
			"ON t_pushevent.id = t_reportversion.push_id ").
		Preload("Event").
		Preload("Report").
		Find(&results).Error

	return results, err
}

//TODO update
func ListReportVersions(limit, offset int64) ([]gormmodels.ReportVersionAndEvent, error) {
	result := make([]gormmodels.ReportVersionAndEvent, 0)
	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ? ", 1).
		Joins("LEFT JOIN t_report " +
			"ON t_report.id = t_reportversion.report_id " +
			"LEFT JOIN t_pushevent " +
			"ON t_pushevent.id = t_reportversion.push_id ").
		Order("t_reportversion.id DESC").
		Offset(offset).
		Limit(limit).
		Preload("Event").
		Preload("Report").
		Find(&result).Error
	return result, err
}

func CountReportVersions() (int64, error) {
	var count int64 = 0
	err := common.GetDB().Table("t_reportversion").Where("t_reportversion.enable_flag = ? ", 1).
		Joins("LEFT JOIN t_report " +
			"ON t_report.id = t_reportversion.report_id " +
			"LEFT JOIN t_pushevent " +
			"ON t_pushevent.id = t_reportversion.push_id ").
		Preload("Event").
		Preload("Report").
		Count(&count).Error
	return count, err
}
