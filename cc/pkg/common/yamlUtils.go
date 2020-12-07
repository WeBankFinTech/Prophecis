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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
)

var appConfig models.AppConfig

func init() {
	appConfig = loadYaml()
}

func GetAppConfig() models.AppConfig {
	return appConfig
}

func loadYaml() models.AppConfig {
	log := logger.Logger()
	conf := new(models.AppConfig)
	yamlFile, err := ioutil.ReadFile("/etc/config/application.yml")

	//log.Println("yamlFile:", yamlFile)
	if err != nil {
		log.Errorf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, conf)
	// err = yaml.Unmarshal(yamlFile, &resultMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	log.Debugf("conf: %+v", conf)
	return *conf
}
