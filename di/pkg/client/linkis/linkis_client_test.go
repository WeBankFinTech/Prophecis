package linkis

import (
	"log"
	"testing"
)

func TestExecute(t *testing.T) {
	username := "hduser05"
	request := ExecuteRequest{
		ExecuteApplicationName: "flowexecution",
		ExecutionCode:          `"{"projectVersionId":602,"flowId":1011,"version":"v000001"}"`,
		Params:                 make(map[string]interface{}),
		RequestApplicationName: "flowexecution",
		RunType:                "json",
		Source: SourceOfExecuteRequest{
			FlowName:    "buoytest2",
			ProjectName: "bdapWorkspace_hduser05",
		},
	}
	client := GetLinkisClient()
	data, err := client.Execute(request, username)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("result: %+v", data)
}

func TestStatus(t *testing.T) {
	username := "hduser05"
	execID := "131317flowexecutionflowexecutionbdpdws110001:9006flowexecution_hduser05_2"
	client := GetLinkisClient()
	data, err := client.Status(execID, username)
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("result: %+v", data)
}

func TestGetExecute(t *testing.T){
	username := "alexwu"
	execID := "131317flowexecutionflowexecutionbdpdws110001:9006flowexecution_alexwu_12"
	client := GetLinkisClient()
	execData, err := client.GetExecution(execID,username)
	if err != nil{
		println(err.Error())
	}
	println(len(execData.PendingJobs))
}