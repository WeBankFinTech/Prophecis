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
package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
	"testing"
	"webank/DI/commons/logger"
	"webank/DI/pkg/datasource/mysql"
	"webank/DI/restapi/api_v1/restmodels"
)

func TestCreateExperiment(t *testing.T) {
	username := "alexwu"
	expeName := "t513er344531" + uuid.New().String()
	expDesc := "test3414213386444"
	datasource.InitDS()
	params := &restmodels.ProphecisExperimentRequest{
		ExpName: &expeName,
		ExpDesc: &expDesc,
	}
	experiment, err := ExperimentService.CreateExperiment(params, username)
	if err != nil {
		println(err.Error())
	}
	println(experiment.ID)
}

func TestListCreateExperiment(t *testing.T) {
	var size int64 = 10
	var page int64 = 1
	user := "alexwu"
	queryStr := ""
	datasource.InitDS()
	experiment, count,err := ExperimentService.ListExperiments(size, page, user, queryStr, true)
	if err != nil {
		println(err.Error())
	}
	println(len(experiment))
	println(count)
}

func TestCreateFlow(t *testing.T) {
	flowName := "buoy_test_service_15"
	username := "hduser05"

	datasource.InitDS()

	params := &restmodels.ProphecisExperimentRequest{
		ExpName: &flowName,
		ExpDesc: nil,
	}

	experiment, err := ExperimentService.CreateExperiment(params, username)
	if err != nil {
		log.Printf("%v\n", err.Error())
	} else {
		log.Printf("%+v\n", experiment)
	}
}

func TestUpdateFlow(t *testing.T) {
	var flowName int64 = 1
	username := "hduser05"
	//username := "alexwu"

	datasource.InitDS()

	json := `{"edges":[],"nodes":[{"key":"346cc47a-578a-4d94-9e3a-c9348349d5e4","title":"qualitis_9411","desc":"","layout":{"height":40,"width":150,"x":204,"y":245},"resources":[],"selected":false,"id":"346cc47a-578a-4d94-9e3a-c9348349d5e4","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","title":"qualitis_202","desc":"","layout":{"height":40,"width":150,"x":286,"y":336},"resources":[],"selected":false,"id":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"62e616e6-8002-469c-9a7f-4a2c354a7302","title":"subFlow_2611","desc":"","layout":{"height":40,"width":150,"x":681,"y":303},"params":{"configuration":{"special":{},"runtime":{},"startup":{}}},"resources":[],"createTime":1598495890907,"selected":false,"creator":"","jobContent":{"embeddedFlowId":1013},"bindViewKey":"","id":"62e616e6-8002-469c-9a7f-4a2c354a7302","jobType":"workflow.subflow"}],"comment":"paramsSave","type":"flow","updateTime":1598495936594,"props":[{"user.to.proxy":"hduser05"}],"resources":[],"scheduleParams":{"proxyuser":"hduser05"},"contextID":"{\"type\":\"HAWorkFlowContextID\",\"value\":\"{\\\"instance\\\":null,\\\"backupInstance\\\":null,\\\"user\\\":\\\"hduser05\\\",\\\"workspace\\\":\\\"bdapWorkspace\\\",\\\"project\\\":\\\"bdapWorkspace_hduser05\\\",\\\"flow\\\":\\\"buoytest2\\\",\\\"contextId\\\":\\\"24-24--YmRwZHdzMTEwMDAxOjkxMTY\\\\u003dYmRwZHdzMTEwMDAxOjkxMTY\\\\u003d95974\\\",\\\"version\\\":\\\"v000001\\\",\\\"env\\\":\\\"BDP_DEV\\\"}\"}"}`
	//json := `{}`
	params := &restmodels.ProphecisExperimentPutRequest{
		ExperimentID: &flowName,
		FlowJSON:     &json,
	}

	experiment, err := ExperimentService.UpdateExperiment(params, username)
	if err != nil {
		log.Printf("%v\n", err.Error())
	} else {
		log.Printf("%+v\n", experiment)
	}
}

func TestParseJson(t *testing.T) {
	log.Printf("str: %v", "123")

	var r map[string]json.RawMessage

	jsonStr := `{"edges":[],"nodes":[{"key":"346cc47a-578a-4d94-9e3a-c9348349d5e4","title":"qualitis_9411","desc":"","layout":{"height":40,"width":150,"x":204,"y":245},"resources":[],"selected":false,"id":"346cc47a-578a-4d94-9e3a-c9348349d5e4","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","title":"qualitis_202","desc":"","layout":{"height":40,"width":150,"x":286,"y":336},"resources":[],"selected":false,"id":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"62e616e6-8002-469c-9a7f-4a2c354a7302","title":"subFlow_2611","desc":"","layout":{"height":40,"width":150,"x":681,"y":303},"params":{"configuration":{"special":{},"runtime":{},"startup":{}}},"resources":[],"createTime":1598495890907,"selected":false,"creator":"","jobContent":{"embeddedFlowId":1013},"bindViewKey":"","id":"62e616e6-8002-469c-9a7f-4a2c354a7302","jobType":"workflow.subflow"}],"comment":"paramsSave","type":"flow","updateTime":1598495936594,"props":[{"user.to.proxy":"hduser05"}],"resources":[],"scheduleParams":{"proxyuser":"hduser05"},"contextID":"{\"type\":\"HAWorkFlowContextID\",\"value\":\"{\\\"instance\\\":null,\\\"backupInstance\\\":null,\\\"user\\\":\\\"hduser05\\\",\\\"workspace\\\":\\\"bdapWorkspace\\\",\\\"project\\\":\\\"bdapWorkspace_hduser05\\\",\\\"flow\\\":\\\"buoytest2\\\",\\\"contextId\\\":\\\"24-24--YmRwZHdzMTEwMDAxOjkxMTY\\\\u003dYmRwZHdzMTEwMDAxOjkxMTY\\\\u003d95974\\\",\\\"version\\\":\\\"v000001\\\",\\\"env\\\":\\\"BDP_DEV\\\"}\"}"}`

	err := json.Unmarshal([]byte(jsonStr), &r)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	} else {
		log.Printf("%+v\n", r)
	}
	var rl []map[string]json.RawMessage

	err = json.Unmarshal([]byte(r["nodes"]), &rl)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	}
	for _, rr := range rl {
		//log.Printf("id: %v,codePath: %v\n", string(rr["id"]),rr["codePath"])
		rr["codePath"] = json.RawMessage(`"s3://asd.zip"`)
		log.Printf("id: %v,codePath: %v\n", string(rr["id"]), string(rr["codePath"]))
		break
	}

	//revert
	bytes, err := json.Marshal(&rl)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	}
	r["nodes"] = bytes

	bytes, err = json.Marshal(r)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	}

	//show result
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	} else {
		log.Printf("%+v\n", r)
	}
	var rl2 []map[string]json.RawMessage

	err = json.Unmarshal([]byte(r["nodes"]), &rl2)
	if err != nil {
		log.Printf("%v\n", err.Error())
		return
	}
	for _, rr := range rl2 {
		//log.Printf("id: %v,codePath: %v\n", string(rr["id"]),rr["codePath"])
		//rr["codePath"] = []byte("s3://asd.zip")
		log.Printf("id: %v,codePath: %v\n", string(rr["id"]), string(rr["codePath"]))

	}

}

//type node struct {
//	Id string `json:"id"`
//
//}

func TestGetExperiment(t *testing.T) {
	logger.Config()
	val, _ := logrus.ParseLevel("debug")
	logrus.SetLevel(val)

	var flowName int64 = 1
	username := "hduser05"

	datasource.InitDS()

	experiment, json, err := ExperimentService.GetExperiment(flowName, username)
	if err != nil {
		log.Printf("%v\n", err.Error())
	} else {
		log.Printf("%+v, %v\n", experiment, json)
	}
}

func TestExport(t *testing.T) {
	logger.Config()
	val, _ := logrus.ParseLevel("debug")
	logrus.SetLevel(val)

	var flowName int64 = 95
	username := "hduser05"
	datasource.InitDS()

	model, err := ExperimentService.Export(flowName, username)
	if err != nil {
		log.Printf("%v\n", err.Error())
	} else {
		log.Printf("%+v\n", model)
	}
}

func TestUnzip(t *testing.T) {
	err := archiver.Unarchive(`E:\gitrepo\FfDL\test_export.zip`, `E:\gitrepo\FfDL\test`)
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
}
