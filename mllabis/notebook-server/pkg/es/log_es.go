package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"webank/AIDE/notebook-server/pkg/commons/logger"
)

type Log struct {
	EsIndex string
	EsType  string
	PodName string
	Ns      string
	From    int
	Size    int
	Asc     bool
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
		return nil, errors.New("log service is not connected")
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
			logger.Logger().Errorf("test-lk error map: %+v", e)
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
