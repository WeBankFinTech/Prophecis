package mlflow

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
	"webank/DI/pkg/client"
)

const  (
	CREATE_EXPERIMENT_URL = "/ajax-api/2.0/preview/mlflow/experiments/create"
)

type CreateExperimentRequest struct {
	Name string `json:"name"`
	ArtifactLocation string `json:"artifact_location"`
}

type CreateExperimentResponse struct {
	ExperimentId string `json:"experiment_id"`
}

type MLFlowClient struct {
	client.BaseClient
}

func GetMLFlowClient() MLFlowClient {

	httpClient := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}

	mlflowClient := MLFlowClient{
		BaseClient: client.BaseClient{
			Address: "http://"+viper.GetString("mlflow.address"),
			Client:  httpClient,
		},
	}
	return mlflowClient
}

func (c *MLFlowClient) CreateExperiment(request *CreateExperimentRequest) (*CreateExperimentResponse, error) {
	requestURL := c.BaseClient.Address + CREATE_EXPERIMENT_URL

	bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//init request
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	//Do Request
	body, err := c.DoHttpRequest(req)
	if err != nil {
		return nil, err
	}

	var createExperimentResponse CreateExperimentResponse
	err = json.Unmarshal(body, &createExperimentResponse)
	if err != nil{
		return nil, err
	}

	return &createExperimentResponse, nil
}

