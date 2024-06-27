package service_v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/kubeflow/pipelines/backend/api/v2beta1/go_client"
	"github.com/kubeflow/pipelines/backend/src/apiserver/server/flowjson"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
	loggerv2 "webank/DI/commons/logger/v2"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/pkg/v2/client"
	models "webank/DI/pkg/v2/model"
	"webank/DI/pkg/v2/repo"
	"webank/DI/pkg/v2/util"
	"webank/DI/restapi/api_v2/restmodels"
)

var ExperimentRunService ExperimentRunServiceIF = &ExperimentRunServiceImpl{}

type ExperimentRunServiceIF interface {
	CreateExperimentRun(ctx context.Context, currentUserName string, clusterType string, request *restmodels.CreateExperimentRunRequest) (*models.ExperimentRun, error)
	GetExperimentRun(currentUserName string, clusterType string, expRunId string) (*models.ExperimentRun, error)
	KillExperimentRun(currentUserName string, clusterType string, expRunId string) error
	RetryExperimentRun(currentUserName string, clusterType string, expRunId string) error
	ListExperimentRuns(internal bool, currentUserName string, clusterType string, expId *string, expVersionName *string, status []string,
		createTimeSt *string, createTimeEd *string,
		size *int64, page *int64) ([]*models.ExperimentRun, int64, error)
	ListExperimentRunNodes(currentUserName string, clusterType string, expRunId string) ([]*models.ExperimentRunNode, error)

	DeleteExperimentRun(currentUserName string, clusterType string, expRunId string) error
	UpdateExperimentRunExecuteStatus() error
}

type ExperimentRunServiceImpl struct {
}

func (e ExperimentRunServiceImpl) ListExperimentRunNodes(currentUserName string, clusterType string, expRunId string) ([]*models.ExperimentRunNode, error) {
	canAccess, err := CanAccessExpRun(currentUserName, clusterType, expRunId)
	if err != nil {
		log.Errorf("CanAccessExpRun(%s, %s) failed: +%v", currentUserName, expRunId, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}
	expRunNodes, err := repo.ExperimentRunNodeRepo.GetAllByExpRunID(expRunId)
	if err != nil {
		log.Errorf("ExperimentRunNodeRepo.GetAllByExpRunID failed: %+v", err)
		return nil, err
	}
	return expRunNodes, nil
}

func (e ExperimentRunServiceImpl) GetExperimentRun(currentUserName string, clusterType string, expRunId string) (*models.ExperimentRun, error) {
	canAccess, err := CanAccessExpRun(currentUserName, clusterType, expRunId)
	if err != nil {
		log.Errorf("CanAccessExpRun(%s, %s) failed: +%v", currentUserName, expRunId, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}

	expRun, err := repo.ExperimentRunRepo.Get(expRunId)
	if err != nil {
		log.Errorf("ExperimentRunRepo.Get failed: %+v", err)
		return nil, err
	} else if expRun.ExpRunID == "" {
		log.Errorf("do not find record of expRunId(%s)", expRunId)
		return nil, fmt.Errorf("do not find record of expRunId(%s)", expRunId)
	}
	return expRun, nil
}

/*
time="2024-06-12T15:40:44+08:00" level=info msg="thre are 1 expreimentRun is not final state" func=webank/DI/restapi/service_v2.ExperimentRunServiceImpl.UpdateExperimentRunExecuteStatus file="/var/jenkins_home/workspace/MLSS-DI-GPU-Develop-gaoyuan/restapi/service_v2/experiment_run.go:86" Operation=UpdateExperimentRunExecuteStatus X-Request-Id=ce5892cd2ed1408797062b339614ee05
time="2024-06-12T15:40:44+08:00" level=info msg="the runInfo.RunDetails of e630229e-3a56-4676-a3ab-ebf72a38aa28 is
task_details:{
  run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\"
  task_id:\"42438e86-b786-471c-9f65-2c6e2b3f5b6b\"
  display_name:\"executor\"
  create_time:{seconds:1718164385}
  start_time:{seconds:1718164416}
  end_time:{seconds:-62135596800}
  state:RUNNING
  state_history:{update_time:{seconds:1718164417} state:PENDING}
  state_history:{update_time:{seconds:1718164427} state:RUNNING}
}
task_details:{
  run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\"
  task_id:\"4aa6e16f-9ce9-4660-b831-7f69f896048b\"
  display_name:\"task-s14imtifm0xylt\"
  create_time:{seconds:1718164385}
  start_time:{seconds:1718164416}
  end_time:{seconds:-62135596800}
  state:RUNNING
  state_history:{update_time:{seconds:1718164417} state:RUNNING}
  child_tasks:{pod_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv-2269131434\"}
}
task_details:{
  run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\"
  task_id:\"4b8b448e-b262-4ea1-9f58-2293115a09bd\"
  display_name:\"task-s14imtifm0xylt-driver\" create_time:{seconds:1718164385} start_time:{seconds:1718164406}
  end_time:{seconds:1718164411}
  state:SUCCEEDED
  state_history:{update_time:{seconds:1718164407} state:PENDING} state_history:{update_time:{seconds:1718164417} state:SUCCEEDED}
  child_tasks:{pod_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv-2310288339\"}
}
task_details:{run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\" task_id:\"8718a0f2-8356-4c10-a201-c882f865082c\"
  display_name:\"root-driver\" create_time:{seconds:1718164385} start_time:{seconds:1718164386} end_time:{seconds:1718164397}
  state:SUCCEEDED
  state_history:{update_time:{seconds:1718164387} state:PENDING}
  state_history:{update_time:{seconds:1718164397} state:RUNNING}
  state_history:{update_time:{seconds:1718164407} state:SUCCEEDED}
  child_tasks:{pod_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv-3272604890\"}}
task_details:{run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\" task_id:\"dc653fec-f6be-4297-8abf-8da37d0abffa\"
  display_name:\"root\" create_time:{seconds:1718164385} start_time:{seconds:1718164406} end_time:{seconds:-62135596800}
  state:RUNNING state_history:{update_time:{seconds:1718164407} state:RUNNING}
  child_tasks:{pod_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv-4136973796\"}}
task_details:{run_id:\"e630229e-3a56-4676-a3ab-ebf72a38aa28\" task_id:\"e121edb6-2fa8-4b76-a61d-d932d766b43b\"
  display_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv\"
  create_time:{seconds:1718164385} start_time:{seconds:1718164386} end_time:{seconds:-62135596800}
  state:RUNNING
  state_history:{update_time:{seconds:1718164387} state:RUNNING}
  child_tasks:{pod_name:\"07fda893b2174070b4a9bb2374075024-v2-d3daafc15364487eb01bb4dkwtv-131277647\"}}"

func=webank/DI/restapi/service_v2.ExperimentRunServiceImpl.UpdateExperimentRunExecuteStatus file="/var/jenkins_home/workspace/MLSS-DI-GPU-Develop-gaoyuan/restapi/service_v2/experiment_run.go:101" Operation=UpdateExperimentRunExecuteStatus X-Request-Id=ce5892cd2ed1408797062b339614ee05
*/

func (e ExperimentRunServiceImpl) UpdateExperimentRunExecuteStatus() error {
	ctx := context.WithValue(context.Background(), loggerv2.OperationLogFieldKey, "UpdateExperimentRunExecuteStatus")
	ctx = context.WithValue(ctx, loggerv2.RequestIdLogFieldKey, util.CreateUUID())
	var log = loggerv2.GetLogger(ctx)

	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Errorf("client.GetKFPClient() failed: +%v", err)
		return err
	}
	taskSleepTime := 10 * time.Second
	for {
		ctx = context.WithValue(ctx, loggerv2.RequestIdLogFieldKey, util.CreateUUID())
		var log = loggerv2.GetLogger(ctx)
		time.Sleep(taskSleepTime)
		size := int64(math.MaxInt64)
		page := int64(1)
		experimentRuns, total, err := e.ListExperimentRuns(true, "", "", nil, nil,
			[]string{restmodels.ExperimentRunExecuteStatusInitializing, restmodels.ExperimentRunExecuteStatusRunning},
			nil, nil, &size, &page)
		if err != nil {
			log.Errorf("ListExperimentRuns failed: %+v", err)
			continue
		}
		log.Infof("thre are %d expreimentRun is not final state", total)
		for _, experimentRun := range experimentRuns {
			runRequest := &api.GetRunRequest{
				RunId: experimentRun.KfpRunId,
			}
			runInfo, err := kfpClient.RunServiceClient.GetRun(context.Background(), runRequest)
			if err != nil {
				log.Errorf("kfpClient.RunServiceClient.GetRunState failed: %+v", err)
				continue
			}
			experimentRunExecuteStatus := kfpRunStateToExperimentRunState(int(runInfo.State))
			if experimentRunExecuteStatus == "" {
				log.Warnf("the current runState.State(%s) will not be processed.", runInfo.State)
				continue
			}
			log.Infof("the runInfo.RunDetails of %s is %+v", experimentRun.KfpRunId, runInfo.RunDetails)
			for _, taskDetail := range runInfo.RunDetails.TaskDetails {
				// 找到displayname为 task-{nodeKey} 的taskDetail
				if strings.HasPrefix(taskDetail.DisplayName, "task-") && !strings.HasSuffix(taskDetail.DisplayName, "driver") {
					taskKey := strings.Split(taskDetail.DisplayName, "-")[1]
					expRunID := experimentRun.ExpRunID
					taskState := kfpRunStateToExperimentRunState(int(taskDetail.State))
					updateMap := make(map[string]interface{})
					updateMap["state"] = taskState
					err = repo.ExperimentRunNodeRepo.UpdateByMap1(expRunID, taskKey, updateMap)
					if err != nil {
						log.Errorf("ExperimentRunNodeRepo.UpdateByMap failed: %+v", err)
						continue
					}
				}
			}

			// todo(gaoyuanhe): 还有结束时间没有修改
			updateMap := make(map[string]interface{})
			updateMap["execute_status"] = experimentRunExecuteStatus
			err = repo.ExperimentRunRepo.UpdateByMap(experimentRun.ExpRunID, updateMap)
			if err != nil {
				log.Errorf("ExperimentRunRepo.UpdateByMap failed: %+v", err)
				continue
			}
		}

	}
}

func kfpRunStateToExperimentRunState(input int) string {
	switch input {
	case 0, 1:
		return restmodels.ExperimentRunExecuteStatusInitializing
	case 2:
		return restmodels.ExperimentRunExecuteStatusRunning
	case 3:
		return restmodels.ExperimentRunExecuteStatusSucceed
	case 5:
		return restmodels.ExperimentRunExecuteStatusFailed
	case 7:
		return restmodels.ExperimentRunExecuteStatusCancelled
	default:
		return ""
	}
}

func (e ExperimentRunServiceImpl) DeleteExperimentRun(currentUserName string, clusterType string, expRunId string) error {
	canAccess, err := CanAccessExpRun(currentUserName, clusterType, expRunId)
	if err != nil {
		log.Errorf("CanAccessExpRun(%s, %s) failed: +%v", currentUserName, expRunId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	updateMap := make(map[string]interface{})
	updateMap["deleted"] = 1
	err = repo.ExperimentRunRepo.UpdateByMap(expRunId, updateMap)
	if err != nil {
		log.Errorf("ExperimentRunRepo.UpdateByMap failed: %+v", err)
		return err
	}
	return nil
}

func (e ExperimentRunServiceImpl) ListExperimentRuns(internal bool, currentUserName string, clusterType string, expId *string, expVersionName *string, status []string,
	createTimeSt *string, createTimeEd *string,
	size *int64, page *int64) ([]*models.ExperimentRun, int64, error) {
	if !internal {
		canAccess, err := CanAccessExp(currentUserName, clusterType, *expId)
		if err != nil {
			log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, *expId, err)
			return nil, 0, err
		}
		if !canAccess {
			return nil, 0, AccessError{}
		}
	}

	rOffset, rSize := GetOffSetAndSize(page, size)
	expRuns, total, err := repo.ExperimentRunRepo.GetAllByOffset(expId, expVersionName, status,
		createTimeSt, createTimeEd, rSize, rOffset)
	if err != nil {
		log.Errorf("ExperimentRunRepo.GetAllByOffset failed: %+v", err)
		return nil, 0, err
	}
	if expRuns == nil || len(expRuns) == 0 {
		log.Errorf("the repo.ExperimentRepo.GetAllByOffset result experiments is nil or empty")
	}
	return expRuns, total, nil
}

func (e ExperimentRunServiceImpl) CreateExperimentRun(ctx context.Context, currentUserName string, clusterType string, request *restmodels.CreateExperimentRunRequest) (*models.ExperimentRun, error) {
	log := loggerv2.GetLogger(ctx)
	canAccess, err := CanAccessExp(currentUserName, clusterType, request.ExpID)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, request.ExpID, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}

	// 通过experimentId从db中获取experiment的详细信息

	// 通过experiment关联的kfp pipelineId和pipelineVersionId(比如v7)信息获取flowjson

	// 最挫的方法就是利用输入的globalVariable来新建一个kfp的flowjson
	// 然后更新kfp新建一个 v7_runTime 的临时版本，然后运行这个版本的kfp
	// 这样也没有修改原来的v7版本的kfp工作流
	// 最好的方法要和kfp的输入参数结合在一起，将全局变量直接传给kfp的输入参数

	// 下面的实现暂就是用runConfig实现的全局变量的使用
	var exp *models.Experiment
	if request.ExpVersionName == "" {
		exp, err = repo.ExperimentRepo.GetLatestVersion(request.ExpID)
	} else {
		exp, err = repo.ExperimentRepo.GetByVersion(request.ExpID, request.ExpVersionName)
	}
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return nil, err
	} else if exp == nil {
		log.Errorf("Experiment(%s, %s) does not exist", request.ExpID, request.ExpVersionName)
		return nil, fmt.Errorf("experiment(%s, %s) does not exist", request.ExpID, request.ExpVersionName)
	}
	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Error("get kfpClient error : " + err.Error())
		return nil, err
	}
	expRunId := util.CreateUUID()
	displayName := fmt.Sprintf("%s-%s-%s", exp.ID, exp.VersionName, expRunId)
	createRunRequest := &api.CreateRunRequest{
		Run: &api.Run{
			DisplayName: displayName,
			PipelineSource: &api.Run_PipelineVersionReference{
				PipelineVersionReference: &api.PipelineVersionReference{
					PipelineId:        exp.KfpPipelineId,
					PipelineVersionId: exp.KfpPipelineVersionId,
				},
			},
			RuntimeConfig: &api.RuntimeConfig{
				Parameters: make(map[string]*structpb.Value),
			},
		},
	}
	// todo(gaoyuanhe): 在mlss执行的时候，前端可能强制了用户填入runConfig.GlobalVariables
	// 但是在dss执行的时候，发布调度的时候就将globalVariables固化了，
	// 所以这里也要支持如果request.RunConfig为空，但是flowjson里的globalVariables是完整
	// 那就用flowjson里的globalVariables作为 createRunRequest.Run.RuntimeConfig.Parameters["global_variables"]
	if request.RunConfig != nil && len(request.RunConfig.GlobalVariables) != 0 {
		globalVariables, err := json.Marshal(request.RunConfig.GlobalVariables)
		log.Infof("the request.RunConfig.GlobalVariables is: %+v", globalVariables)
		if err != nil {
			log.Errorf("json.Marshal failed: %+v", err)
		} else {
			createRunRequest.Run.RuntimeConfig.Parameters["global_variables"] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: string(globalVariables)}}
		}
	}

	if request.RunConfig != nil && request.RunConfig.WtssVariables != nil {
		wtssVariablesJson, err := json.Marshal(request.RunConfig.WtssVariables)
		log.Infof("the request.RunConfig.WtssVariables is: %+v", wtssVariablesJson)
		if err != nil {
			log.Errorf("json.Marshal failed: %+v", err)
		} else {
			createRunRequest.Run.RuntimeConfig.Parameters["wtss_variables"] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: string(wtssVariablesJson)}}
		}
	}

	// 获取在kfp保存的flowjson, 存储在 t_experiment_run_v2 表中，因为kfp中的flowjson可能会被update
	getPipLineVersionRequest := &api.GetPipelineVersionRequest{PipelineId: exp.KfpPipelineId, PipelineVersionId: exp.KfpPipelineVersionId}
	pipelineVersion, err := kfpClient.PipelineServiceClient.GetPipelineVersion(context.Background(), getPipLineVersionRequest)
	if err != nil {
		log.Errorf("kfpClient.GetPipelineVersion() failed: +%v", err)
		return nil, err
	}
	if len(pipelineVersion.FlowJson) <= 0 {
		return nil, errors.New("flowJson is empty")
	}

	// 解析flowjson，创建对应的 ExperimentRunNode
	var flowJson flowjson.FlowJson
	err = json.Unmarshal([]byte(pipelineVersion.FlowJson), &flowJson)
	if err != nil {
		log.Errorf("json.Unmarshal(pipelineVersion.FlowJson) failed: +%v", err)
		return nil, err
	}

	createRunResponse, err := kfpClient.RunServiceClient.CreateRun(context.Background(), createRunRequest)
	if err != nil {
		log.Error("kfp CreateRun error : " + err.Error())
		return nil, err
	}

	currentTime := time.Now()
	newExperimentRun := &models.ExperimentRun{
		ExpRunID:       expRunId,
		Name:           displayName,
		CreateUser:     currentUserName,
		CreateTime:     currentTime,
		ExpID:          exp.ID,
		ExpVersionName: exp.VersionName,
		ExecuteStatus:  restmodels.ExperimentRunExecuteStatusInitializing,
		KfpRunId:       createRunResponse.RunId,
		FlowJson:       pipelineVersion.FlowJson,
	}

	newExperimentRunNodes := make([]*models.ExperimentRunNode, 0)
	for _, node := range flowJson.Nodes {
		newExperimentRunNode := &models.ExperimentRunNode{
			ExperimentRunID: expRunId,
			Key:             node.Key,
			Title:           node.Title,
			Type:            node.JobType,
			State:           "",
			ProxyID:         "",
		}
		newExperimentRunNodes = append(newExperimentRunNodes, newExperimentRunNode)
	}

	err = datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		err = repo.ExperimentRunRepo.Add(newExperimentRun)
		if err != nil {
			return err
		}
		/* todo(gaoyuanhe): 暂时注释
		for _, newExperimentRunNode := range newExperimentRunNodes {
			err = repo.ExperimentRunNodeRepo.Add(newExperimentRunNode)
			if err != nil {
				return err
			}
		}
		*/

		// 修改experiment的最后执行时间
		updateMap := make(map[string]interface{})
		updateMap["latest_execute_time"] = currentTime
		err = repo.ExperimentRepo.UpdateByMap(exp.ID, exp.VersionName, updateMap)
		if err != nil {
			log.Errorf("ExperimentRepo.UpdateByMap failed: %+v", err)
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Create ExperimentRun Error:" + err.Error())
		return nil, err
	}
	return newExperimentRun, nil
}

func (e ExperimentRunServiceImpl) KillExperimentRun(currentUserName string, clusterType string, expRunId string) error {
	canAccess, err := CanAccessExpRun(currentUserName, clusterType, expRunId)
	if err != nil {
		log.Errorf("CanAccessExpRun(%s, %s) failed: +%v", currentUserName, expRunId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	expRun, err := repo.ExperimentRunRepo.Get(expRunId)
	if err != nil {
		log.Errorf("ExperimentRunRepo.Get failed: %+v", err)
		return err
	} else if expRun.ExpRunID == "" {
		log.Errorf("do not find record of expRunId(%s)", expRunId)
		return fmt.Errorf("do not find record of expRunId(%s)", expRunId)
	}

	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Error("get kfpClient error : " + err.Error())
		return err
	}
	terminateRunRequest := &api.TerminateRunRequest{
		RunId: expRun.KfpRunId,
	}
	_, err = kfpClient.RunServiceClient.TerminateRun(context.Background(), terminateRunRequest)
	if err != nil {
		log.Errorf("TerminateRun failed: %+v", err)
		// rpc error: code = InvalidArgument desc = Failed to terminate a run: Failed to terminate run 6c12ef56-5840-493e-9b55-4c07c5a85b98:
		// Invalid input error: Failed to terminate a run 6c12ef56-5840-493e-9b55-4c07c5a85b98. Row not found
		// 这个错误一般是这个运行已经不存在了，已经被终止了，所以不用返回错误；应该让前端控制已终止的不能点终止
		if strings.Contains(strings.ToLower(err.Error()), "row not found") {
			log.Warnf("这个错误一般是这个运行已经不存在了，已经被终止了，所以不用返回错误")
			return nil
		}
		return err
	}
	return nil
}

func (e ExperimentRunServiceImpl) RetryExperimentRun(currentUserName string, clusterType string, expRunId string) error {
	canAccess, err := CanAccessExpRun(currentUserName, clusterType, expRunId)
	if err != nil {
		log.Errorf("CanAccessExpRun(%s, %s) failed: +%v", currentUserName, expRunId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	expRun, err := repo.ExperimentRunRepo.Get(expRunId)
	if err != nil {
		log.Errorf("ExperimentRunRepo.Get failed: %+v", err)
		return err
	} else if expRun.ExpRunID == "" {
		log.Errorf("do not find record of expRunId(%s)", expRunId)
		return fmt.Errorf("do not find record of expRunId(%s)", expRunId)
	}

	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Error("get kfpClient error : " + err.Error())
		return err
	}
	retryRunRequest := &api.RetryRunRequest{
		RunId: expRun.KfpRunId,
	}
	_, err = kfpClient.RunServiceClient.RetryRun(context.Background(), retryRunRequest)
	if err != nil {
		log.Error("RetryRun failed: %+v", err)
		return err
	}

	// todo(gaoyuanhe): 这里是否需要修改状态？还是由独立的线程去修改状态？
	updateMap := make(map[string]interface{})
	updateMap["execute_status"] = restmodels.ExperimentRunExecuteStatusInitializing
	err = repo.ExperimentRunRepo.UpdateByMap(expRunId, updateMap)
	if err != nil {
		log.Errorf("ExperimentRunRepo.UpdateByMap failed: %+v", err)
		return err
	}
	return nil
}
