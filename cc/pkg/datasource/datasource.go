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
package datasource

import (
	"fmt"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/logger"
)
import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/mysql"

var DB *gorm.DB
//var MlssDS DataSource
var TokenDS DataSource

const (
	WbUserToken = "WB_USER_TOKEN"
)



func InitDS() {
	log := logger.Logger()
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}

	dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		common.GetAppConfig().Application.Datasource.Username,
		common.GetAppConfig().Application.Datasource.Password,
		common.GetAppConfig().Application.Datasource.Ip,
		common.GetAppConfig().Application.Datasource.Port,
		common.GetAppConfig().Application.Datasource.Db,
	)
	log.Infof("mlss-db url: %v", dbUrl)

	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("failed to connect database, reason: %v", err.Error())
	}
	db.SingularTable(true)
	DB = db

	//TokenDS = initTokenDB()
}

func GetDB() *gorm.DB {
	return DB
}



type DataSource struct {
	DB *gorm.DB
}

func (ds *DataSource) GetDB() *gorm.DB {
	return ds.DB
}

//func (ds *DataSource) GetDBByDBS(txs []*gorm.DB) *gorm.DB {
//	var db *gorm.DB
//	if len(txs) > 0 {
//		db = txs[0]
//	} else {
//		db = ds.GetDB().New()
//		defer db.Close()
//
//	}
//	return db
//}
