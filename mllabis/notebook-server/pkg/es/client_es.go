package es

import (
	"fmt"
	"os"
	"webank/AIDE/notebook-server/pkg/commons/logger"

	elasticsearch "github.com/elastic/go-elasticsearch/v5"
)

const FLUENT_ELASTICSEARCH_HOST = "FLUENT_ELASTICSEARCH_HOST"
const FLUENT_ELASTICSEARCH_PORT = "FLUENT_ELASTICSEARCH_PORT"
const FLUENT_ELASTICSEARCH_USER = "FLUENT_ELASTICSEARCH_USER"
const FLUENT_ELASTICSEARCH_PASSWD = "FLUENT_ELASTICSEARCH_PASSWD"

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
		// EnableDebugLogger: true,
		// EnableMetrics:     true,
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
