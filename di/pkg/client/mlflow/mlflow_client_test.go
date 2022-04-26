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