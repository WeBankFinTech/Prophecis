package newService

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/models"
	dss "webank/DI/pkg/client/dss1_0"
	"webank/DI/pkg/client/linkis1_0"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/pkg/repo"
	"webank/DI/storage/storage"
)

var ExperimentRunService ExperimentRunServiceIF

func init() {
	ExperimentRunService = &ExperimentRunServiceImpl{}
}

type ExperimentRunServiceIF interface {
	CreateExperimentRun(expId int64, flowJson *string, runDate string, username string, isSA bool) (*models.ExperimentRun, error)
	ListExperimentRun(page int64, size int64, queryStr string, username string, isSA bool) ([]*models.ExperimentRun, int64, error)
	GetExperimentRun(execId int64, username string, isSA bool) (*models.ExperimentRun, *string, error)
	GetExperimentRunLog(expRunId int64, fromLine int64, size int64, user string, isSA bool) (*linkis.LogData, error)
	GetExperimentRunStatus(expRunId int64, username string, isSA bool) (*string, error)
	DeleteExperimentRun(expRunId int64, username string, isSA bool) error
	KillExperimentRun(expRunId int64, username string, isSA bool) error
	GetRunExecution(expRunId int64, user string, isSA bool) (*linkis.ExecutionData, error)
	GetLogPath(taskId int64, user string) (string, error)
	GetOpenLog(user, logPath string) (*linkis.LogData, error)
	GetRunHistory(expId int64, page int64, size int64, queryStr string, username string, isSA bool) ([]*models.ExperimentRun, int64, error)
}

type ExperimentRunServiceImpl struct {
}

func (impl *ExperimentRunServiceImpl) GetOpenLog(user, logPath string) (*linkis.LogData, error) {
	linkisClient := linkis.GetLinkisClient()
	return linkisClient.GetOpenLog(logPath, user)
}

func (impl *ExperimentRunServiceImpl) GetLogPath(taskId int64, user string) (string, error) {
	linkisClient := linkis.GetLinkisClient()
	return linkisClient.GetLogPath(taskId, user)
}

func (impl *ExperimentRunServiceImpl) GetExperimentRunLog(expRunId int64, fromLine int64, size int64, user string,
	isSA bool) (*linkis.LogData, error) {
	linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()

	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = PermissionCheck(user, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	logPath, err := linkisClient.GetLogPath(experimentRun.DssTaskID, user)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	logData, err := linkisClient.GetOpenLog(logPath, user)
	if err != nil {
		log.Error(err.Error())
		//Fix Return Error, TODO Control Linkis Error Log
		logData = &linkis.LogData{
			ExecID:   "",
			Log:      []string{err.Error()},
			FromLine: 0,
		}
		return logData, nil
	}
	return logData, err
}

func (impl *ExperimentRunServiceImpl) GetExperimentRunStatus(expRunId int64, username string, isSA bool) (*string, error) {
	//linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()

	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = PermissionCheck(username, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}
	//statusData, err := linkisClient.Status(experimentRun.DssExecId, username)
	//if err != nil {
	//	log.Error("")
	//	return nil, err
	//}

	return &experimentRun.ExpExecStatus, err
}

//TODO: ERROR CONTROL
func (impl *ExperimentRunServiceImpl) KillExperimentRun(expRunId int64, username string, isSA bool) error {
	linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()
	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		logger.GetLogger().Error("Get experiment run error, ", err)
		return err
	}
	err = PermissionCheck(username, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return err
	}

	execID, err := linkisClient.Kill(experimentRun.DssExecID, experimentRun.DssTaskID, username)
	if err != nil {
		logger.GetLogger().Error("Kill experiment run error, ", err)
		return err
	}
	if *execID != experimentRun.DssExecID {
		err = errors.New("Kill Flow Execution ExecId is Not Correct")
	}
	sts, err := linkisClient.FlowStatus(experimentRun.DssExecID, experimentRun.DssTaskID, username)
	if err != nil {
		logger.GetLogger().Error("Get experiment run's status err, ", err)
	}
	if sts == nil {
		log.Error("status is nil, ready retry...")
		ticker := time.NewTicker(time.Second * 3)
	End:
		for {
			select {
			case <-ticker.C:
				sts, err = linkisClient.FlowStatus(experimentRun.DssExecID, experimentRun.DssTaskID, username)
				if err != nil {
					logger.GetLogger().Error("Get experiment run's status err, ", err)
					return err
				}
				if sts != nil {
					break End
				}
			}
		}
	}
	if sts.Status == "" {
		sts.Status = "Failed"
	}
	err = repo.ExperimentRunRepo.UpdateStatus(expRunId, sts.Status, db)
	if err != nil {
		logger.GetLogger().Error("Update status error, ", err)
	}
	return err
}

func isExperimentRunStop(status string) bool {
	isStop := false
	if status == models.RUN_STATUS_CANCELLED || status == models.RUN_STATUS_FAILED ||
		status == models.RUN_STATUS_TIMEOUT || status == models.RUN_STATUS_SUCCEED {
		isStop = true
	}
	return isStop
}

//TODO: ERROR CONTROL
func (impl *ExperimentRunServiceImpl) DeleteExperimentRun(expRunId int64, username string, isSA bool) error {
	linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()
	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = PermissionCheck(username, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return err
	}

	if isExperimentRunStop(experimentRun.ExpExecStatus) == false {
		_, err = linkisClient.Kill(experimentRun.DssExecID, experimentRun.DssTaskID, username)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}

	err = repo.ExperimentRunRepo.Delete(expRunId, db)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return err
}

func (impl *ExperimentRunServiceImpl) ListExperimentRun(page int64, size int64, queryStr string, username string,
	isSA bool) ([]*models.ExperimentRun, int64, error) {
	db := datasource.GetDB()
	var count int64
	var experimentRuns []*models.ExperimentRun
	var err error
	if isSA {
		count, err = repo.ExperimentRunRepo.Count(queryStr, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}

		experimentRuns, err = repo.ExperimentRunRepo.GetAllByOffset((page-1)*size, size, queryStr, db)
		if err != nil {
			log.Error(err.Error())
			return nil, count, err
		}

	} else {
		count, err = repo.ExperimentRunRepo.CountByUser(queryStr, username, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}
		experimentRuns, err = repo.ExperimentRunRepo.GetAllByUserOffset((page-1)*size, size, queryStr, username, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}
	}
	return experimentRuns, count, err
}

//TODO: ERROR CONTROL
func (impl *ExperimentRunServiceImpl) GetExperimentRun(expRunId int64, username string, isSA bool) (*models.ExperimentRun, *string, error) {

	db := datasource.GetDB()

	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		log.Error("Get Experiment Run Error:", err.Error())
		return nil, nil, err
	}
	err = PermissionCheck(username, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, nil, err
	}
	//10.107.127.19:37017
	flowRepo, err := storage.NewFlowRepository(viper.GetString(config.MongoAddressKey),
		viper.GetString(config.MongoDatabaseKey), viper.GetString(config.MongoUsernameKey),
		viper.GetString(config.MongoPasswordKey), viper.GetString(config.MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "flow_json")
	if err != nil {
		log.Error(err.Error())
		return experimentRun, nil, err
	}

	flowData, err := flowRepo.GetFlowJson(expRunId)
	if err != nil {
		log.Error(err.Error())
		return experimentRun, nil, err
	}

	return experimentRun, &flowData.ExpJson, err
}

/**
Send Experiment Run Flow Engine Execute to Linkis Flow Engine
*/
func flowExecute(dssFlowId int64, version string, flowName string, label dss.DSSLabel,
	username string, isSA bool) (*linkis.ExecuteData, error) {
	client := linkis.GetLinkisClient()
	//TODO Check Default Label
	request := linkis.ExecuteRequest{
		ExecuteApplicationName: "flowexecution",
		ExecutionCode:          fmt.Sprintf(`{"flowId":%v,"version":"%v"}`, dssFlowId, version),
		RequestApplicationName: "flowexecution",
		Params:                 make(map[string]interface{}),
		RunType:                "json",
		Label:                  label,
		Source: linkis.SourceOfExecuteRequest{
			ProjectName: fmt.Sprintf("%v_%v", DefaultWorkspace, username),
			FlowName:    flowName,
		},
	}
	execData, err := client.FlowExecute(request, username)
	return execData, err
}
func flowJsonSave(experimentRun *models.ExperimentRun, flowJson string) error {
	//Save Current Flow Json In MongoDB
	flowRepo, err := storage.NewFlowRepository(viper.GetString(config.MongoAddressKey),
		viper.GetString(config.MongoDatabaseKey), viper.GetString(config.MongoUsernameKey),
		viper.GetString(config.MongoPasswordKey), viper.GetString(config.MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "flow_json")
	if err != nil {
		log.WithError(err).Fatalf("Cannot create Repository with %s %s %s %s", "mongo",
			"test", "mlss", "mlss")
	}
	strconv.FormatInt(experimentRun.ID, 10)
	flowEntry := storage.FlowJsonEntry{
		ExecId:           experimentRun.ID,
		DssFlowVersionID: experimentRun.DssFlowLastVersion,
		DssFlowId:        experimentRun.DssFlowID,
		ExpJson:          flowJson,
	}
	err = flowRepo.AddFlowJson(flowEntry)
	return err
}

func (*ExperimentRunServiceImpl) GetRunExecution(expRunId int64, user string, isSA bool) (*linkis.ExecutionData, error) {
	linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()

	experimentRun, err := repo.ExperimentRunRepo.Get(expRunId, db)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = PermissionCheck(user, experimentRun.User.ID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	execData, err := linkisClient.GetFlowExecution(experimentRun.DssExecID, user)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return execData, err
}

func (*ExperimentRunServiceImpl) CreateExperimentRun(expId int64, flowJson *string, runDate string, username string,
	isSA bool) (*models.ExperimentRun, error) {
	dssClient := dss.GetDSSClient()
	db := datasource.GetDB()
	user, err := repo.UserRepo.GetByTransaction(&username, db)
	experiment, err := repo.ExperimentRepo.Get(expId, db)
	if err != nil {
		log.Errorf("CreateExperimentRun method Get Experiment From Database Error: " + err.Error())
		return nil, err
	}
	err = PermissionCheck(username, experiment.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}
	if experiment.DssFlowLastVersion == nil {
		log.Errorf("CreateExperimentRun method experiment DssFlowLastVersion is nil")
		return nil, err
	}
	expRun := models.ExperimentRun{
		DssFlowID:          experiment.DssFlowID,
		DssFlowLastVersion: *experiment.DssFlowLastVersion,
		DssExecID:          "",
		DssTaskID:          0,
		ExpExecStatus:      models.RUN_STATUS_INITED,
		ExpExecType:        *experiment.CreateType,
		ExpID:              experiment.ID,
		ExpRunCreateTime:   time.Now(),
		ExpRunModifyUserID: user.ID,
		ExpRunCreateUserID: user.ID,
		BaseModel: models.BaseModel{
			EnableFlag: true,
		},
	}
	//TODO ADD Transaction
	err = repo.ExperimentRunRepo.Add(&expRun, db)
	if err != nil {
		log.Errorf("CreateExperimentRun method Save ExperimentRun in DataBase Error: " + err.Error())
		return &expRun, err
	}
	//var flowData *dss.FlowGetData
	defaultLabel := ""
	flowData, err := dssClient.GetOrchestrator(experiment.DssFlowID, defaultLabel, username)
	if err != nil {
		log.Errorf("CreateExperimentRun method Get Flow From Dss Error: " + err.Error())
		return nil, err
	}
	if flowJson == nil {
		flowJson = &flowData.Flow.FlowJson
	}
	//retry get flow josn
	if len(*flowJson) <= 0 {
		latestFlowJson := flowData.Flow.FlowJson
		flowJson = &latestFlowJson
		if len(latestFlowJson) <= 0 {
			ticker := time.NewTicker(time.Second * 3)
		End:
			for {
				select {
				case <-ticker.C:
					flowData, err := dssClient.GetOrchestrator(experiment.DssFlowID, defaultLabel, username)
					if err != nil {
						log.Errorf("CreateExperimentRun method Get Flow From Dss Error: " + err.Error())
						return nil, err
					}
					log.Println("GetFlow retry...")
					latestFlowJson = flowData.Flow.FlowJson
					if len(latestFlowJson) > 0 {
						ticker.Stop()
						break End
					}
				}
			}
		}
	}

	//1st get flowObj
	var flowObj map[string]json.RawMessage
	err = json.Unmarshal([]byte(*flowJson), &flowObj)
	if err != nil {
		log.Errorf("parser flowObj failed,", err)
		return nil, err
	}
	//check flowObj's nodes
	if _, check := flowObj["nodes"]; !check {
		log.Errorf("flowJson nodes is null: ", *flowJson)
		return nil, err
	}
	//2nd get nodes
	var nodesObj []map[string]json.RawMessage
	err = json.Unmarshal(flowObj["nodes"], &nodesObj)
	if err != nil {
		log.Println("parser flowObj failed ", err)
		return nil, err
	}
	nodesBytes := []byte{}
	for _, nodeObj := range nodesObj {

		//check node's jobContent
		if _, check := nodeObj["jobContent"]; !check {
			log.Println("jobContent is null")
			return nil, err
		}
		jobContentObj := map[string]json.RawMessage{}
		err = json.Unmarshal(nodeObj["jobContent"], &jobContentObj)
		if err != nil {
			log.Println("job content un marshal failed,", err)
			return nil, err
		}

		//update run date
		nodeType := fmt.Sprint(jobContentObj["mlflowJobType"])
		if strings.Compare(nodeType, "ModelMonitor") == 1 && runDate != "" {
			content := map[string]interface{}{}
			//jobContentObj["content"] = json.RawMessage("\"" + runDate + "\"")
			err = json.Unmarshal(jobContentObj["content"], &content)
			if err != nil {
				log.Errorf("monitor job parse error:", err.Error())
				return nil, err
			}
			nodeRunDate := fmt.Sprint(content["run_date"])
			if nodeRunDate == "${run_date}" {
				content["wtss_run_date"] = json.RawMessage("\"" + runDate + "\"")
			}
			contentBytes, _ := json.Marshal(content)
			_ = json.Unmarshal(contentBytes, &content)
			jobContentObj["content"], _ = json.Marshal(&content)
			jobContentBytes, _ := json.Marshal(jobContentObj)
			_ = json.Unmarshal(jobContentBytes, &jobContentObj)
			nodeObj["jobContent"], _ = json.Marshal(&jobContentObj)
			nodeObjBytes, _ := json.Marshal(nodeObj)
			_ = json.Unmarshal(nodeObjBytes, &nodeObj)
			nodesBytes = append(nodesBytes, nodeObjBytes...)
		}

		//GPU Job Check
		if _, check := jobContentObj["ManiFest"]; check {
			maniFestObj := map[string]interface{}{}
			err = json.Unmarshal(jobContentObj["ManiFest"], &maniFestObj)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			maniFestObj["exp_run_id"] = expRun.ID
			maniFestBytes, _ := json.Marshal(maniFestObj)
			_ = json.Unmarshal(maniFestBytes, &maniFestObj)
			jobContentObj["ManiFest"], _ = json.Marshal(&maniFestObj)
			jobContentBytes, _ := json.Marshal(jobContentObj)
			_ = json.Unmarshal(jobContentBytes, &jobContentObj)
			nodeObj["jobContent"], _ = json.Marshal(&jobContentObj)
			nodeObjBytes, _ := json.Marshal(nodeObj)
			_ = json.Unmarshal(nodeObjBytes, &nodeObj)
			nodesBytes = append(nodesBytes, nodeObjBytes...)
		}
	}
	_ = json.Unmarshal(nodesBytes, &nodesObj)
	newFlowObj := map[string]interface{}{}
	flowObjBytes, _ := json.Marshal(&flowObj)
	_ = json.Unmarshal(flowObjBytes, &newFlowObj)
	newFlowObj["nodes"] = nodesObj
	flowJsonBytes, err := json.Marshal(&newFlowObj)
	if err != nil {
		log.Errorf("CreateExperimentRun marshal Error: " + err.Error())
		return nil, err
	}
	newFlowJson := string(flowJsonBytes)
	flowJson = &newFlowJson
	data, err := saveFlow(experiment, flowJson, flowData.Flow.FlowEditLock, username)
	if err != nil {
		log.Errorf("CreateExperimentRun method Save Flow From Dss Error: " + err.Error())
		return nil, err
	}

	//Update Experiment in Transaction
	err = datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		experiment.DssFlowLastVersion = &data.FlowVersion
		experiment.DssBmlLastVersion = &flowData.Flow.BmlVersion
		experiment.ExpModifyTime = time.Now()
		err = repo.ExperimentRepo.Update(experiment, tx)
		return err
	})
	if err != nil {
		log.Errorf("Update Experiment Info Error:" + err.Error())
		return nil, err
	}

	devLabel := dss.DSSLabel{
		Route: "dev",
	}

	execData, err := flowExecute(experiment.DssFlowID, data.FlowVersion, flowData.Flow.Name, devLabel, username, true)
	if err != nil {
		log.Errorf("CreateExperimentRun method Send Run Flow Request to Linkis Error: " + err.Error())
		return nil, err
	}
	expRun.DssExecID = execData.ExecID
	expRun.DssTaskID = execData.TaskID
	err = repo.ExperimentRunRepo.Update(&expRun, db)
	if err != nil {
		log.Errorf("CreateExperimentRun update experiment run Error: " + err.Error())
		return nil, err
	}
	// Save Latest Flow Json in MongoDB
	err = flowJsonSave(&expRun, *flowJson)
	if err != nil {
		log.Errorf("CreateExperimentRun method Save flow Error: " + err.Error())
		return nil, err
	}
	go updateFlowStatus(expRun.ID, expRun.DssExecID, expRun.DssTaskID, expRun.ExpExecStatus, username)
	return &expRun, nil
}

func updateFlowStatus(runId int64, dssExecId string, taskID int64, status string, user string) {
	linkisClient := linkis.GetLinkisClient()
	db := datasource.GetDB()
	dataChan := make(chan *linkis.StatusData)
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ticker := time.NewTicker(time.Second * 25)
End:
	for {
		select {
		case sts := <-dataChan:
			if sts != nil {
				if isExperimentRunStop(status) != true {
					if status != sts.Status {
						err := repo.ExperimentRunRepo.UpdateStatus(runId, sts.Status, db)
						if err != nil {
							log.Error(err.Error())
							ticker.Stop()
							break End
						}
					}
					status = sts.Status
				} else {
					experimentRun := &models.ExperimentRun{}
					updateErr := db.Model(&experimentRun).Where("id = ?", runId).Update("exp_run_end_time", time.Now()).Error
					if updateErr != nil {
						log.Error(errors.New(" update experiment run end time error: " + updateErr.Error()))
					}
					ticker.Stop()
					break End
				}
			}
		case <-ticker.C:
			go func() {
				sts, err := linkisClient.FlowStatus(dssExecId, taskID, user)
				if sts == nil {
					log.Error(errors.New(" update status error: status is nil "))
				}
				if err != nil {
					log.Error(err.Error())
				}
				dataChan <- sts
			}()
		case <-signalChan:
			done <- true
			break End
		}
	}
}

func (impl *ExperimentRunServiceImpl) GetRunHistory(expId int64, page int64, size int64, queryStr string, username string, isSA bool) ([]*models.ExperimentRun, int64, error) {
	db := datasource.GetDB()
	var count int64
	var experimentRuns []*models.ExperimentRun
	var err error
	if isSA {
		count, err = repo.ExperimentRunRepo.CountExpId(expId, queryStr, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}
		experimentRuns, err = repo.ExperimentRunRepo.GetByExpIdAndOffset(expId, page, size, queryStr, db)
		if err != nil {
			log.Error(err.Error())
			return nil, count, err
		}
	} else {
		count, err = repo.ExperimentRunRepo.CountByExpIdAndUser(expId, queryStr, username, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}
		experimentRuns, err = repo.ExperimentRunRepo.GetByUserAndExpIdAndOffset(expId, (page-1)*size, size, queryStr, username, db)
		if err != nil {
			log.Error(err.Error())
			return nil, -1, err
		}
	}
	return experimentRuns, count, err
}
