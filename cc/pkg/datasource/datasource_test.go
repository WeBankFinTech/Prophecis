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
	"github.com/jinzhu/gorm"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"testing"
)

func TestGetTokenDB(t *testing.T) {

	//db := GetTokenDB()
	db := TokenDS.GetDB()
	var token Token
	err := db.Table("WB_USER_TOKEN").Find(&token, "USER = ?", "hadoop").Error

	if err != nil {
		fmt.Printf("err: %v", err.Error())
	}
	fmt.Printf("token: %+v", token)
}

type Token struct {
	User        string `gorm:"column:USER"`
	Token       string `gorm:"column:TOKEN"`
	ClientToken string `gorm:"column:CLIENT_TOKEN"`
}

//type Token struct {
//	User        string
//	Token       string
//	ClientToken string
//}

func TestTX(t *testing.T) {
	log := logger.Logger()

	//db := GetDB()
	db := GetDB().New()
	defer db.Close()
	//TokenDS.

	var result models.Group

	err := db.Transaction(func(tx *gorm.DB) error {

		//add user-------------------
		user := models.User{
			Name:      "aha-buoy",
			Type:      "SYSTEM",
			Gid:       80001,
			UID:       80001,
			GUIDCheck: "1",
		}

		err := tx.Create(&user).Error
		if err != nil {
			log.Errorf(err.Error())
			return err
		}

		//add group-------------------
		group := models.Group{
			Name:         "aha-buoy",
			ClusterName:  "BDAP",
			DepartmentID: "89999",
		}

		err = tx.Create(&group).Error
		if err != nil {
			log.Errorf(err.Error())
			return err
		}
		result = group

		return nil
	})

	if err != nil {
		log.Errorf("tx err: %v", err.Error())
	}
	log.Infof("result: %v", result)
}
