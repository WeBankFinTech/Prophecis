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
package mlflow

import (
	"github.com/spf13/viper"
	"strconv"
	"testing"
	"time"
)

func TestCreateExperimemnt(t *testing.T) {
	viper.Set("mlflow.address","")
	client := GetMLFlowClient()
	username := "alexwu"
	currentTime := time.Now().Unix()
	currentTimeStr := strconv.FormatInt(currentTime,10)

	request := CreateExperimentRequest{
		Name: "test"+"_"+ username +"_"+ currentTimeStr,
	}

	response, err := client.CreateExperiment(&request)
	if err !=nil{
		print(err.Error())
		return
	}
	print(response.ExperimentId)

}