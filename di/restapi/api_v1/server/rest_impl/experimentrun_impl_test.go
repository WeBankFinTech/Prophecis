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
	"io/ioutil"
	"testing"
	"time"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/restapi/api_v1/server/operations/experiment_runs"
)

func TestCreateExperimentRunRest(t *testing.T) {
	str, err := ioutil.ReadFile("D:\\flowjson.txt")
	if err != nil {
		println("read file error:" + err.Error())
		return
	}

	flowJson := string(str)
	expExecType := "MLFLOW"

	datasource.InitDS(false)
	expRes := restmodels.ProphecisExperimentRunRequest{
		ExpExecType: &expExecType,
		FlowJSON:    flowJson,
	}
	params := experiment_runs.CreateExperimentRunParams{
		HTTPRequest:          nil,
		ExpID:                9,
		ExperimentRunRequest: &expRes,
	}
	result := CreateExperimentRun(params)
	//result := 1
	time.Sleep(time.Second * 300)
	println(result)
}

func TestDeleteExperimentRunRest(t *testing.T) {
	datasource.InitDS(false)
	params := experiment_runs.DeleteExperimentRunParams{
		HTTPRequest: nil,
		ID:          47,
	}
	result := DeleteExperimentRun(params)
	//result := 1
	println(result)
}

func TestKillExperimentRunRest(t *testing.T) {

	datasource.InitDS(false)

	params := experiment_runs.KillExperimentRunParams{
		HTTPRequest: nil,
		ID:          48,
	}

	result := KillExperimentRun(params)
	//result := 1
	println(result)
}

func TestStatusExperimentRun(t *testing.T) {

	datasource.InitDS(false)

	params := experiment_runs.GetExperimentRunStatusParams{
		HTTPRequest: nil,
		ID:          48,
	}

	result := StatusExperimentRun(params)
	print(result)
}

func TestListExperimentRuns(t *testing.T) {
	datasource.InitDS(false)
	var page int64 = 1
	var size int64 = 10
	var queryStr = ""
	var username  = "alexwu"
	params := experiment_runs.ListExperimentRunsParams{
		HTTPRequest: nil,
		Page:        &page,
		QueryStr:    &queryStr,
		Size:        &size,
		Username:   &username,
	}

	exp := ListExperimentRuns(params)
	println(exp)
}

//func TestUpdateFlowStatus(t *testing.T){
//	exp := service.ExperimentRunServiceImpl{}
//	datasource.InitDS()
//	_,_ := exp.CreateExperimentRun(9,"MLFLOW","alexwu")
//	//exp.UpdateFlowStatus(expRun.DssExecId,expRun.ExpExecStatus,"alexwu")
//}

func TestGetFlowExecution(t *testing.T){
	datasource.InitDS(false)
	params := experiment_runs.GetExperimentRunExecutionParams{ExecID:182}
	GetExperimentRunExecution(params)
	println(1)
}

func TestGetExpermientRun(t *testing.T){
	datasource.InitDS(false)
	params := experiment_runs.GetExperimentRunParams{
		ID:          281,
	}
	GetExperimentRun(params)
	println(1)
}