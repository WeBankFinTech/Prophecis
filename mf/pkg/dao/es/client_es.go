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
package es

import (
	"fmt"
	"mlss-mf/pkg/logger"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v5"
)

const FLUENT_ELASTICSEARCH_HOST = "FLUENT_ELASTICSEARCH_HOST"
const FLUENT_ELASTICSEARCH_PORT = "FLUENT_ELASTICSEARCH_PORT"

// const FLUENT_ELASTICSEARCH_USER = "FLUENT_ELASTICSEARCH_USER"
// const FLUENT_ELASTICSEARCH_PASSWD = "FLUENT_ELASTICSEARCH_PASSWD"

var esClient *elasticsearch.Client

func init() {
	host := os.Getenv(FLUENT_ELASTICSEARCH_HOST)
	port := os.Getenv(FLUENT_ELASTICSEARCH_PORT)
	// username := os.Getenv(FLUENT_ELASTICSEARCH_USER)
	// passwd := os.Getenv(FLUENT_ELASTICSEARCH_PASSWD)
	logger.Logger().Infof("elasticsearch client host: %q, port: %q", host, port)
	if host == "" || port == "" {
		panic(fmt.Sprintf("elastic search host or port can't be empty string, host: %s, port: %s", host, port))
	}

	url := fmt.Sprintf("http://%s:%s", host, port)
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		// Username:          username,
		// Password:          passwd,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	esClient = client
}

func GetESClient() *elasticsearch.Client {
	return esClient
}
