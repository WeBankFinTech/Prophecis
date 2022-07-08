/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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
	"github.com/prometheus/common/log"
	"testing"
	"webank/DI/pkg/datasource/mysql"
)

func TestCreateExp(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	//s := "1"
	//exp := models.Experiment{
	//	//ExpID: 1,
	//	ExpDesc: &s,
	//	ExpName: &s,
	//}
	//err := Add(&exp, db)
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//s := "12355577"
	//name := "asd555555555555555"
	//bool := false
	//exp := models.Experiment{
	//
	//	ExpID:   1,
	//	ExpDesc: &s,
	//	ExpName: &name,
	//	BaseModel: models.BaseModel{
	//		ID:         4,
	//		EnableFlag: &bool,
	//	},
	//}
	//err := UpdateExperiment(&exp, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//err := DeleteExperimentById(1, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//experiment, err := GetExperimentByExperimentId(3, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//} else {
	//	fmt.Printf("experiment: %v", experiment)
	//}

	//offset := GetAllExperimentsByOffset(0, 99, db.Debug())
	//fmt.Printf("experiment: %v", offset)

	//experiment := GetAllExperiments(db.Debug())
	//fmt.Printf("experiment: %v", experiment)
	user := "alexwu"

	experiment,err := ExperimentRepo.CountByUser("",user, db.Debug())
	if err != nil{
		log.Error(err.Error())
	}
	fmt.Printf("experiment: %v", experiment)
}


func TestGetAllExperiments(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()
	experiments, err := ExperimentRepo.GetAllByOffset(0, 100, "", db.Debug())
	if err != nil{
		println(err.Error())
	}
	println(len(experiments))
}

func TestGetExperiments(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()
	experiments, err := ExperimentRepo.Get(9,db)
	if err != nil{
		println(err.Error())
	}
	println(experiments)
}