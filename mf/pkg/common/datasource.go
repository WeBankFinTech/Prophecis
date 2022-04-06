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

package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mlss-mf/pkg/common/config"
	"mlss-mf/pkg/logger"
)

var (
	DB *gorm.DB
)

func init() {
	logger := logger.Logger()
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
	//dbConfig := config.GetMFconfig().DBConfig
	//logger.Infof("reason: %v", dbConfig)
	//dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Server, dbConfig.Port, dbConfig.Database)
	//logger.Info("echo ",dbUrl)
	dbUrl := fmt.Sprintf("%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.GetMFconfig().DBConfig.Username,
		config.GetMFconfig().DBConfig.Password,
		config.GetMFconfig().DBConfig.Server,
		config.GetMFconfig().DBConfig.Port,
		config.GetMFconfig().DBConfig.Database,
	)
	db, err := gorm.Open("mysql", dbUrl)
	if err != nil {
		logger.Infof("reason: %v", err.Error())
		panic("failed to connect database")
	}
	db.SingularTable(true)
	db.LogMode(config.GetMFconfig().DBConfig.SqlLog)
	DB = db
}

func GetDB() *gorm.DB {
	return DB.New()
}
