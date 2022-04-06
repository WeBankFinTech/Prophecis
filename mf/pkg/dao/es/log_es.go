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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mlss-mf/pkg/logger"
)

type Log struct {
	EsIndex       string
	EsType        string
	PodName       string
	ContainerName string
	Ns            string
	From          int
	Size          int
	Asc           bool
}

const (
	K8S_LOG_CONTAINER_ES_INDEX = "kube_log_container_es_index"
	K8S_LOG_CONTAINER_ES_TYPE  = "kube_log_container_es_type"
	ES_MAX_SIZE                = 1000
)

func LogList(param Log) (map[string]interface{}, error) {
	esClient := GetESClient()

	order := "desc"
	if param.Asc {
		order = "asc"
	}
	if param.Size > ES_MAX_SIZE {
		param.Size = ES_MAX_SIZE
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"ns.keyword": param.Ns,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"pod_name.keyword": param.PodName,
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"container_name.keyword": param.ContainerName,
						},
					},
				},
			},
		},
		"size": param.Size,
		"from": param.From,
		"sort": map[string]interface{}{
			"time": map[string]interface{}{
				"order": order,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Logger().Errorf("fail to encode query: %s", err.Error())
		return nil, err
	}

	resp, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(param.EsIndex),
		esClient.Search.WithBody(&buf),
		esClient.Search.WithPretty(),
	)

	if err != nil {
		logger.Logger().Errorf("fail to get es data, error: %s", err.Error())
		return nil, errors.New("log service connection failed")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			logger.Logger().Errorf("Error parsing the response body: %s", err)
			return nil, err
		} else {
			logger.Logger().Errorf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			logger.Logger().Errorf("error map: %+v", e)
			err = fmt.Errorf("[%s] %s: %s",
				resp.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"])
			return nil, err
		}
	}

	r := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		logger.Logger().Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	return r, nil
}
