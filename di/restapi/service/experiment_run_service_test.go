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
	"fmt"
	"testing"
)

func TestSprintf(t *testing.T) {

	sprintf := fmt.Sprintf(`{"projectVersionId":%v,"flowId":%v,"version":"%v"}`, 1, 2, 3)
	fmt.Println(sprintf)
}
func TestExecute(t *testing.T) {
	//var expId int64 = 1
	//username := "hduser05"
	//datasource.InitDS()

	//run, err := ExperimentRunService.Execute(expId, username)
	//if err != nil {
	//	log.Printf("err: %v\n", err.Error())
	//} else {
	//	log.Printf("result: %+v\n", run)
	//}
}
func TestGetLogPath(t *testing.T) {
	logPath, err1 := ExperimentRunService.GetLogPath(35763, "alexwu")
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("Path：" + logPath)
}

func TestGetOpenLog(t *testing.T) {
	logPath, err1 := ExperimentRunService.GetLogPath(35763, "alexwu")
	if err1 != nil {
		log.Fatal(err1)
	}
	openLog, err2 := ExperimentRunService.GetOpenLog("alexwu", logPath)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(openLog.Log[3])
}
