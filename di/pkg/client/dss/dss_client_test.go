package dss

import (
	"testing"
)

func TestGetWorkspaces(t *testing.T) {
	client := GetDSSClient()
	workspaces, err := client.GetWorkspaces("hduser05")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("result: %+v", workspaces)

}

func TestListPersonalWorkflows(t *testing.T) {
	client := GetDSSClient()
	workspaces, err := client.ListPersonalWorkflows(104, "hduser05")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("err: %+v", workspaces)

}

func TestCreatePersonalWorkflow(t *testing.T) {
	params := CreatePersonalWorkflowParams{
		Name:             "buoy-test-dss-cli4",
		Description:      "test",
		ProjectVersionID: 602,
		WorkspaceId:      104,
		Uses:             "",
	}

	client := GetDSSClient()
	workspaces, err := client.CreatePersonalWorkflow(params, "hduser05")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("err: %+v", workspaces)

}

func TestSaveWorkflow(t *testing.T) {
	params := SaveWorkflowParams{
		BmlVersion:       "v000009",
		Comment:          "paramsSave",
		FlowEditLock:     "",
		FlowVersion:      "v000001",
		Id:               1011,
		Json:             `{"edges":[],"nodes":[{"key":"346cc47a-578a-4d94-9e3a-c9348349d5e4","title":"qualitis_9411","desc":"","layout":{"height":40,"width":150,"x":204,"y":245},"resources":[],"selected":false,"id":"346cc47a-578a-4d94-9e3a-c9348349d5e4","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","title":"qualitis_202","desc":"","layout":{"height":40,"width":150,"x":286,"y":336},"resources":[],"selected":false,"id":"b4f7aa77-53fb-4f12-b45c-fdcb45df5576","businessTag":"","appTag":null,"jobType":"linkis.appjoint.qualitis"},{"key":"62e616e6-8002-469c-9a7f-4a2c354a7302","title":"subFlow_2611","desc":"","layout":{"height":40,"width":150,"x":681,"y":303},"params":{"configuration":{"special":{},"runtime":{},"startup":{}}},"resources":[],"createTime":1598495890907,"selected":false,"creator":"","jobContent":{"embeddedFlowId":1013},"bindViewKey":"","id":"62e616e6-8002-469c-9a7f-4a2c354a7302","jobType":"workflow.subflow"}],"comment":"paramsSave","type":"flow","updateTime":1598495936594,"props":[{"user.to.proxy":"hduser05"}],"resources":[],"scheduleParams":{"proxyuser":"hduser05"},"contextID":"{\"type\":\"HAWorkFlowContextID\",\"value\":\"{\\\"instance\\\":null,\\\"backupInstance\\\":null,\\\"user\\\":\\\"hduser05\\\",\\\"workspace\\\":\\\"bdapWorkspace\\\",\\\"project\\\":\\\"bdapWorkspace_hduser05\\\",\\\"flow\\\":\\\"buoytest2\\\",\\\"contextId\\\":\\\"24-24--YmRwZHdzMTEwMDAxOjkxMTY\\\\u003dYmRwZHdzMTEwMDAxOjkxMTY\\\\u003d95974\\\",\\\"version\\\":\\\"v000001\\\",\\\"env\\\":\\\"BDP_DEV\\\"}\"}"}`,
		ProjectID:        10,
		ProjectVersionID: 602,
	}

	client := GetDSSClient()
	workspaces, err := client.SaveFlow(params, "hduser05")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("err: %+v", workspaces)

}
func TestGetWorkflow(t *testing.T) {

	client := GetDSSClient()
	workspaces, err := client.GetFlow(1011, "", 602, "hduser05")
	if err != nil {
		log.Printf("err: %v", err.Error())
	}
	log.Printf("err: %+v", workspaces)

}
