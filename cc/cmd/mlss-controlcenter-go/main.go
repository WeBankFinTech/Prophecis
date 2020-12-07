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
package main

import (
	"github.com/IBM/FfDL/commons/config"
	"github.com/spf13/viper"
	"log"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/restapi/restapi"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	config.InitViper()

	viper.SetConfigName("application.yml") // name of config file (without extension)
	viper.AddConfigPath("/etc/config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper.ReadInConfig failed, %v \n", err.Error())
	}

	datasource.InitDS()
	authcache.Init()

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		logger.Logger().Fatalf(err.Error())
	}

	api := operations.NewMlssCcAPI(swaggerSpec)
	server := restapi.NewServer(api)
	//server.EnabledListeners = []string{"http"}
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "mlss-cc API!!"
	parser.LongDescription = "API description in Markdown."

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			logger.Logger().Fatalf(err.Error())
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		logger.Logger().Fatalf(err.Error())
	}

}
