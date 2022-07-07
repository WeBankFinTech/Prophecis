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
package rest_impl

import (
	"testing"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/experiments"
)

func TestListExperimentRest(t *testing.T) {
	var size int64 = 10
	var page int64 = 1
	user := "alexwu"
	queryStr := ""
	datasource.InitDS(false)
	params := experiments.ListExperimentsParams{
		Page:     &page,
		Size:     &size,
		QueryStr: &queryStr,
		Username: &user,
	}
	result := ListExperiments(params)
	//result := 1
	println(result)
}

func TestCreateExperimentRest(t *testing.T) {
	//user := "alexwu"
	expName := "MLSS-T23WT"
	expDesc := "MLSS-TRSTETEWT"
	datasource.InitDS(false)
	req := restmodels.ProphecisExperimentRequest{
		ExpName: &expName,
		ExpDesc: &expDesc,
	}

	params := experiments.CreateExperimentParams{
		Experiment: &req,
	}
	result := CreateExperiment(params)
	//result := 1
	println(result)
}

func TestGetExperimentRest(t *testing.T) {

	datasource.InitDS(false)
	params := experiments.GetExperimentParams{
		ID: 34,
	}
	result := GetExperiment(params)
	//result := 1
	println(result)
}

func TestUpdateExperimentRest(t *testing.T) {
	datasource.InitDS(false)
	exp := restmodels.ProphecisExperimentPutRequest{

	}
	params := experiments.UpdateExperimentParams{
		Experiment: &exp,
	}

	result := UpdateExperiment(params)
	//result := 1
	println(result)
}

func TestUpdateExperimentInfoRest(t *testing.T) {
	datasource.InitDS(false)

	exp := restmodels.ProphecisExperimentPutBasicInfo{
		ExpName: "ipdateaname",
		ExpDesc: "testeestsets",
	}

	params := experiments.UpdateExperimentInfoParams{
		ID:         7,
		Experiment: &exp,
	}

	result := UpdateExperimentInfo(params)
	//result := 1
	println(result)
}

func TestDeleteExperimentRest(t *testing.T) {
	datasource.InitDS(false)
	params := experiments.DeleteExperimentParams{
		ID: -1,
	}
	result := DeleteExperiment(params)
	println(result)
}