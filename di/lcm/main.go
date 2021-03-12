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
	"os"

	"github.com/spf13/viper"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/service/client"
	"webank/DI/commons/util"
	"webank/DI/lcm/service/lcm"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.InitViper()
	logger.Config()

	log.Printf("lcm FLUENT_ELASTICSEARCH_HOST %s", os.Getenv("FLUENT_ELASTICSEARCH_HOST"))
	log.Printf("lcm FLUENT_ELASTICSEARCH_PORT %s", os.Getenv("FLUENT_ELASTICSEARCH_PORT"))
	log.Printf("lcm FLUENT_ELASTICSEARCH_USER %s", os.Getenv("FLUENT_ELASTICSEARCH_USER"))
	log.Printf("lcm FLUENT_ELASTICSEARCH_PASSWD %s", os.Getenv("FLUENT_ELASTICSEARCH_PASSWD"))

	log.Printf("DLAAS_PUSHGATEPWAY_SERVER %s", os.Getenv("DLAAS_PUSHGATEPWAY_SERVER"))
	log.Printf("DLAAS_PUSHGATEPWAY_PORT %s", os.Getenv("DLAAS_PUSHGATEPWAY_PORT"))

	port := viper.GetInt(config.PortKey)
	if port == 0 {
		port = client.LcmLocalPort
	}
	service, err := lcm.NewService()
	if err != nil {
		log.WithError(err).Errorf("Failed to start lcm since nil instance of lcm service")
		panic(err)
	}

	util.HandleOSSignals(func() {
		service.StopLCM()
	})
	service.Start(port, false)
}
