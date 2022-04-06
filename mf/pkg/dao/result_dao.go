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
	"time"
)

func GetResultByModelName(name string) ([]*gormmodels.Result, error) {
	var results []*gormmodels.Result
	err := common.GetDB().Table("t_result").Where("t_result.enable_flag = ? ", 1).
		Preload("Model", "model_name = ?", name).
		Preload("ModelVersion").
		Joins("LEFT JOIN t_model " +
			"ON t_model.id = t_result.model_id ").
		Find(&results).Error
	return results, err
}

func DeleteResultByID(id int64) error {
	m := make(map[string]interface{})
	m["enable_flag"] = 0
	m["update_timestamp"] = time.Now()
	return common.GetDB().Table("t_result").Where("id = ?", id).Updates(m).Error
}

func CreateResult(result *gormmodels.Result) error {
	return common.GetDB().Table("t_result").Omit("Model").Create(&result).Error
}

func UpdateResultByID(result *gormmodels.Result, id int64) error {
	return common.GetDB().Table("t_result").Where("id = ?", id).Omit("Model").UpdateColumns(&result).Error
}
