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

package repo

import (
	"fmt"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestGetDeleteStorageByPath(t *testing.T) {
	log := logger.Logger()
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}

	dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local") // common.GetAppConfig().Application.Datasource.Username,
	// common.GetAppConfig().Application.Datasource.Password,
	// common.GetAppConfig().Application.Datasource.Ip,
	// common.GetAppConfig().Application.Datasource.Port,
	// common.GetAppConfig().Application.Datasource.Db,

	log.Infof("mlss-db url: %v", dbUrl)

	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("failed to connect database, reason: %v", err.Error())
	}
	db.SingularTable(true)
	var storage models.Storage
	err = db.Find(&storage, "path = ? AND enable_flag = ?", "/data/bdap-ss/mlss-data/neiljianliu", 1).Error
	if err != nil {
		log.Fatalf("failed to connect database, reason: %v", err.Error())
	}

}
