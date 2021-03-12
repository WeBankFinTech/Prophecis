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

package main

import (
	"github.com/spf13/viper"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/util"
	"webank/DI/trainer/trainer"
)

func main() {
	config.InitViper()
	logger.Config()

	port := viper.GetInt(config.PortKey)
	if port == 0 {
		port = 30005 // TODO don't hardcode
	}
	service := trainer.NewService()

	// var stopSendingMetricsChannel chan struct{}

	util.HandleOSSignals(func() {
		service.StopTrainer()
		// if config.CheckPushGatewayEnabled() {
		// 	stopSendingMetricsChannel <- struct{}{}
		// }
	})
	service.Start(port, false)
}
