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
package newService

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"webank/DI/commons"
	"webank/DI/commons/logger"
	"webank/DI/commons/models"
	"webank/DI/pkg/client/mlflow"

	//"webank/DI/pkg/client/dss"
	dss "webank/DI/pkg/client/dss1_0"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/pkg/repo"
	"webank/DI/restapi/api_v1/restmodels"
	storageClient "webank/DI/storage/client"
	"webank/DI/storage/storage/grpc_storage"

	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"github.com/modern-go/reflect2"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

const (
	DefaultWorkspace     = "bdapWorkspace"
	actualFlowFolderName = "flow"
)

var ExperimentService ExperimentServiceIF = &ExperimentServiceImpl{}
var log = logger.GetLogger()

type ExperimentServiceIF interface {
	UpdateExperiment(*restmodels.ProphecisExperimentPutRequest, string, bool) (*models.Experiment, error)
	CreateExperiment(*restmodels.ProphecisExperimentRequest, string, string, string) (*models.Experiment, error)
	GetExperiment(int64, string, bool) (*models.Experiment, *string, error)
	DeleteExperiment(int64, string, bool) error
	DeleteExperimentTag(int64, string, bool) error
	CreateExperimentTag(int64, string, string) (*models.ExperimentTag, error)
	ListExperiments(int64, int64, string, string, bool) ([]*models.Experiment, int64, error)
	Export(int64, string, bool) (*os.File, error)
	ExportDSS(int64, string) (*restmodels.ProphecisExperimentDssResponse, error)
	Import(io.ReadCloser, int64, *string, string, string, string) (*restmodels.ProphecisExperimentIDResponse, error)
	UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error)
	UploadCode(closer io.ReadCloser) (string, error)
	UpdateExperimentInfo(int64, string, string,
		[]*restmodels.ProphecisExperimentTagPutBasicInfo, string, bool) (*models.Experiment, error)
	ImportDSS(experimentId int64, resourceId string, version *string, desc *string, user string) (*restmodels.ProphecisImportExperimentDssResponse, error)
}
type ExperimentServiceImpl struct {
}

func (expService *ExperimentServiceImpl) CreateExperimentTag(expId int64, expTagStr string, username string) (*models.ExperimentTag, error) {
	db := datasource.GetDB()
	user, err := repo.UserRepo.Get(&username, db)
	if err != nil {
		log.Errorf("DataBase Update Experiment Get User Info Error:" + err.Error())
		return nil, err
	}
	expTag := models.ExperimentTag{
		ExpID:              expId,
		ExpTag:             &expTagStr,
		ExpTagCreateTime:   time.Now(),
		ExpTagCreateUserID: user.ID,
		BaseModel: models.BaseModel{
			EnableFlag: true,
		},
	}
	err = repo.ExperimentTagRepo.Add(&expTag, db)
	if err != nil {
		log.Errorf("DataBase Create Experiment Tag Error:" + err.Error())
		return nil, err
	}
	return &expTag, err
}

func (expService *ExperimentServiceImpl) UpdateExperimentInfo(id int64, expName string, expDesc string,
	expTags []*restmodels.ProphecisExperimentTagPutBasicInfo, username string, isSA bool) (*models.Experiment, error) {
	db := datasource.GetDB()
	user, err := repo.UserRepo.Get(&username, db)
	if err != nil {
		log.Errorf("DataBase Update Experiment Get User Info Error:" + err.Error())
		return nil, err
	}

	exp,err := repo.ExperimentRepo.Get(id, db)
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return nil, err
	}
	err = PermissionCheck(username, exp.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	expTagList, err := repo.ExperimentTagRepo.GetByExpId(id, db)
	expTagMap := map[string]*models.ExperimentTag{}
	if err != nil {
		log.Errorf("DataBase Update Experiment Get Experiment Tag Info Error:" + err.Error())
		return nil, err
	}
	for _, tag := range expTagList {
		for _, v := range expTags {
			if *tag.ExpTag == v.ExpTag {
				expTagMap[*tag.ExpTag] = tag
			}
		}
	}
	for _, tag := range expTagList {
		if _, ok := expTagMap[*tag.ExpTag]; !ok {
			err := repo.ExperimentTagRepo.Delete(tag.TagID, db)
			if err != nil {
				log.Errorf("DataBase Delete Experiment Tag Error:" + err.Error())
				return nil, err
			}
		}
	}
	for _, v := range expTags {
		//Update
		if v.ExpID == -1 {
			err := repo.ExperimentTagRepo.Add(&models.ExperimentTag{
				ExpID:              id,
				ExpTag:             &v.ExpTag,
				ExpTagCreateTime:   time.Now(),
				ExpTagCreateUserID: user.ID,
				BaseModel: models.BaseModel{
					EnableFlag: true,
				},
			}, db)
			if err != nil {
				log.Errorf("DataBase Update Experiment Tag Error:" + err.Error())
				return nil, err
			}
		}
	}
	expUp := models.Experiment{
		ExpName:         &expName,
		ExpDesc:         &expDesc,
		ExpCreateUserID: user.ID,
		ExpModifyUserID: user.ID,
	}
	expUp.ID = id
	err = repo.ExperimentRepo.Update(&expUp, db)
	if err != nil {
		log.Errorf("DataBase Update Experiment Error:" + err.Error())
	}
	return &expUp, err
}

func (expService *ExperimentServiceImpl) UploadCode(closer io.ReadCloser) (string, error) {
	hostpath := "/data/oss-storage/"
	filename := uuid.New().String() + "file.zip"
	filepath := hostpath + filename
	//out, err := os.Create(filename)
	//wt := bufio.NewWriter(out)
	//_, err = io.Copy(wt, closer)
	fileBytes, err := ioutil.ReadAll(closer)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	err = ioutil.WriteFile(filepath, fileBytes, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	//上传 code.zip 后, 更新flowJson 中每个 node 的 codePath
	//创建 storer rpc client, 用于 上传code.zip
	//ctx := context.TODO()
	sClient, err := storageClient.NewStorage()
	if err != nil {
		log.Info("storerClient init Error " + err.Error())
		return "", err
	}

	req := &grpc_storage.CodeUploadRequest{
		FileName: filename,
		HostPath: hostpath,
	}
	ctx := context.TODO()
	if sClient == nil {
		log.Info("storerClient init Error,sclient is nil")
		return "", errors.New("rpc upload code zip error, sclient is nil")
	}
	response, err := sClient.Client().UploadCode(ctx, req)
	if err != nil {
		log.Error("sclient upload code error: " + err.Error())
		return "", err
	}
	err = os.Remove(filepath)
	if err != nil {
		log.Error("" + err.Error())
	}

	return response.S3Path, err
}

func (expService *ExperimentServiceImpl) UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error) {
	hostpath := "/data/oss-storage/"
	filename := uuid.New().String() + "file.zip"
	filepath := hostpath + filename
	fileBytes, err := ioutil.ReadAll(closer)
	if err != nil {
		logger.GetLogger().Error("UploadCodeBucket err, ", err.Error())
		return "", err
	}

	err = ioutil.WriteFile(filepath, fileBytes, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	//上传 code.zip 后, 更新flowJson 中每个 node 的 codePath
	//创建 storer rpc client, 用于 上传code.zip
	//ctx := context.TODO()
	sClient, err := storageClient.NewStorage()
	if err != nil {
		logger.GetLogger().Error("storerClient init Error " + err.Error())
		return "", err
	}
	req := &grpc_storage.CodeUploadRequest{
		FileName: filename,
		HostPath: hostpath,
		Bucket:   bucket,
	}
	ctx := context.TODO()
	if sClient == nil {
		logger.GetLogger().Error("storerClient init Error,sclient is nil")
		return "", errors.New("rpc upload code zip error, sclient is nil")
	}
	response, err := sClient.Client().UploadCode(ctx, req)
	if err != nil {
		logger.GetLogger().Error("sclient upload code error: " + err.Error())
		return "", err
	}
	logger.GetLogger().Info("==Response S3Path：", response.S3Path)
	err = os.Remove(filepath)
	if err != nil {
		logger.GetLogger().Error("os remove err, " + err.Error())
	}
	return response.S3Path, err
}

func (expService *ExperimentServiceImpl) DeleteExperiment(id int64, user string, isSuperAdmin bool) error {
	db := datasource.GetDB()

	// Permission Check
	exp,err := repo.ExperimentRepo.Get(id, db)
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return err
	}
	err = PermissionCheck(user, exp.ExpCreateUserID, nil, isSuperAdmin)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return  err
	}

	//Execute DeleteOperation
	err = db.Transaction(func(tx *gorm.DB) error {
		err := repo.ExperimentRepo.Delete(id, db)
		if err != nil {
			return err
		}

		err = repo.ExperimentTagRepo.BatchDelete(id, db)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Error("ExperimentService DeleteExperiment Error: " + err.Error())
	}
	return err
}

func (expService *ExperimentServiceImpl) DeleteExperimentTag(id int64, user string, isSA bool) error {
	db := datasource.GetDB()
	exp,err := repo.ExperimentRepo.Get(id, db)
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return  err
	}
	err = PermissionCheck(user, exp.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return  err
	}

	return repo.ExperimentTagRepo.Delete(id, db)
}

/**

 */
func (expService *ExperimentServiceImpl) CreateExperiment(params *restmodels.ProphecisExperimentRequest,
	username string, createType string, desc string) (*models.Experiment, error) {
	//TODO: Move To Controller, check params length
	if len(*params.ExpName) <= 0 || len(*params.ExpDesc) <= 0 {
		return nil, errors.New("ExpName's length <= 0 or ExpDesc's length <= 0")
	}

	//Get Default Workspace Id(BDAP Workspace)  From DSS
	wsID, projectId, err := getCreateInfo(username)
	if err != nil {
		log.Error("ExperimentService CreateExperiment GetCreateInfo Error: " + err.Error())
		return nil, err
	}

	//Build Flow In Dss BDAP Workspace
	orchestratorData, _, err := createFlow(*params.ExpName, projectId, strconv.Itoa(int(wsID)), createType, username)
	if err != nil {
		log.Error("ExperimentService CreateExperiment getFlowData Error: " + err.Error())
		return nil, err
	}


	//Build MLFlow Experiment
	mlflowExpId, err := createMLFlowExperiment(*params.ExpName, username)
	if err != nil {
		log.Error("ExperimentService CreateMLFLowExperiment Error: " + err.Error())
		return nil, err
	}

	//TODO
	if params.GroupName == ""{

	}

	//Save Experiment In DataBase
	//if createType is wtss
	if createType == "WTSS" {
		*params.ExpDesc = desc
	}
	//check create is mlflow or dss
	expDescSlice := strings.Split(*params.ExpDesc, ",")
	if len(expDescSlice) > 0 && createType != "WTSS" && createType != "MLFLOW" &&
		expDescSlice[0] == "dss-appjoint" {
		createType = "DSS"
	}
	var newExperiment *models.Experiment

	//TODO
	defaultVersion := "v000001"
	err = datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		user, err := repo.UserRepo.Get(&username, tx)
		if err != nil {
			return err
		}

		projectName := "MLSS" + "_" + username
		newExperiment = &models.Experiment{
			ExpName:             params.ExpName,
			ExpDesc:             params.ExpDesc,
			GroupName:             params.GroupName,
			MlflowExpID:		 mlflowExpId,
			DssBmlLastVersion:   &orchestratorData.Flow.BmlVersion,
			DssDssFlowName: &orchestratorData.Flow.Name,
			DssDssFlowVersion: &defaultVersion,
			DssFlowID:           orchestratorData.Flow.Id,
			DssFlowLastVersion:  &defaultVersion,
			DssDssProjectName: &projectName,
			DssProjectID:        orchestratorData.Flow.Id,
			DssWorkspaceID:      wsID,
			ExpCreateTime:       time.Now(),
			ExpCreateUserID:     user.ID,
			ExpModifyTime:       time.Now(),
			ExpModifyUserID:     user.ID,
			CreateType:          &createType,
			BaseModel: models.BaseModel{
				EnableFlag: true,
			},
		}
		if createType == "DSS" {
			if strings.ToUpper(expDescSlice[1]) != "NULL" {
				newExperiment.DssDssFlowName = &expDescSlice[1]
			}
			if strings.ToUpper(expDescSlice[2]) != "NULL" {
				newExperiment.DssDssFlowVersion = &expDescSlice[2]
			}
			if len(expDescSlice) >= 3 {
				projectName := expDescSlice[4]
				dssClient := dss.GetDSSClient()
				projectID, err := dssClient.GetProjectIdByName(projectName, username, wsID)
				if err != nil {
					log.Error("Get projectName failed,  ", err)
					return nil
				}
				newExperiment.DssDssProjectName = &projectName
				newExperiment.DssDssProjectId = int64(projectID)
				newExpDesc := "dss-appjoint," + expDescSlice[1] + "," + expDescSlice[2] + "," + expDescSlice[4]
				newExperiment.ExpDesc = &newExpDesc
			} else {
				log.Error("expDescSlice's length is ", len(expDescSlice), " get expDescSlice[6] failed.")
				return nil
			}
		}
		if createType == "WTSS" && params.ExpID <= 0 {
			log.Println("createType is WTSS, expId is 0")
			return errors.New("createType is WTSS, expId is 0")
		}
		if createType == "WTSS" && params.ExpID > 0 {
			dssClient := dss.GetDSSClient()
			if len(expDescSlice) >= 5 {
				experiment, _, err := ExperimentService.GetExperiment(params.ExpID, username, true)
				if err != nil {
					log.Error("CreateExperiment(WTSS) GetExperiment custom_error, ", err)
					return nil
				}
				flowData, err := dssClient.GetOrchestrator(experiment.DssFlowID, "dev", username)
				if err != nil {
					log.Error("CreateExperiment(WTSS) GetFlow custom_error, ", err)
					return nil
				}
				newExperiment.DssDssProjectName = experiment.DssDssProjectName
				newExperiment.DssDssProjectId = experiment.DssDssProjectId
				newExperiment.DssDssFlowName = experiment.DssDssFlowName
				newExperiment.DssDssFlowVersion = new(string)
				*newExperiment.DssDssFlowVersion = "v000001"
				dssDssFlowName := ""
				if experiment.DssDssFlowName != nil {
					dssDssFlowName = *experiment.DssDssFlowName
				}
				newExpDesc := "dss-appjoint," + dssDssFlowName + "," + flowData.Flow.Version + "," + expDescSlice[4]
				newExperiment.ExpDesc = &newExpDesc
			}
		}
		//newExperiment.DssDssFlowProjectVersionId = projectVersionId
		err = repo.ExperimentRepo.Add(newExperiment, tx)
		if err != nil {
			return err
		}

		if len(params.TagList) > 0 {
			tagsList := []string{}
			for _, tag := range params.TagList {
				tagsList = append(tagsList, tag.ExpTag)
			}
			err = addTag(newExperiment.ID, tagsList, user.ID, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Errorf("DataBase Save Experiment Error:" + err.Error())
		return nil, err
	}

	return newExperiment, err
}

func addTag(expID int64, tagList []string, userID int64, db *gorm.DB) error {
	var tagModelList []*models.ExperimentTag
	var tempStr *string
	for _, s := range tagList {
		tempStr = &s
		v := *tempStr
		expTag := &models.ExperimentTag{
			ExpID:              expID,
			ExpTag:             &v,
			ExpTagCreateTime:   time.Now(),
			ExpTagCreateUserID: userID,
			BaseModel: models.BaseModel{
				EnableFlag: true,
			},
		}
		tagModelList = append(tagModelList, expTag)
	}
	return repo.ExperimentTagRepo.BatchAdd(tagModelList, db)
}

func (expService *ExperimentServiceImpl) ListExperiments(size int64, page int64, user string, queryStr string, isSA bool) ([]*models.Experiment, int64, error) {
	//client := dss.GetDSSClient()
	db := datasource.GetDB()
	var count int64
	var expes []*models.Experiment
	var err error
	if isSA {
		expes, err = repo.ExperimentRepo.GetAllByOffset((page-1)*size, size, queryStr, db)
		if err != nil {
			log.Error("")
			return nil, -1, err
		}
		count, err = repo.ExperimentRepo.Count(queryStr, db)
		if err != nil {
			log.Error("")
			return expes, -1, err
		}
	} else {
		expes, err = repo.ExperimentRepo.GetAllByUserOffset((page-1)*size, size, queryStr, user, db)
		if err != nil {
			log.Error("")
			return expes, -1, err
		}
		count, err = repo.ExperimentRepo.CountByUser(queryStr, user, db)
		if err != nil {
			log.Error("")
			return expes, -1, err
		}
	}

	return expes, count, nil
}

func (expService *ExperimentServiceImpl) GetExperiment(id int64, username string, isSA bool) (*models.Experiment, *string, error) {
	client := dss.GetDSSClient()
	db := datasource.GetDB()
	//Get Experiment Basic Info From Database
	experiment, err := repo.ExperimentRepo.Get(id, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, nil, err
	}

	//Permission Check
	exp,err := repo.ExperimentRepo.Get(id, db)
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return nil, nil, err
	}
	err = PermissionCheck(username, exp.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, nil, err
	}

	//Get Flow Json From DSS
	log.Debugf("experiment.DssFlowId: %v, data.ProjectVersionId: %v", experiment.DssFlowID, experiment.DssProjectVersionID)
	//TODO: Check Default Label
	defaultLabel := "dev"
	flowData, err := client.GetOrchestrator(experiment.DssFlowID, defaultLabel, username)
	if err != nil {
		log.Errorf(err.Error())
		return nil, nil, err
	}
	return experiment, &flowData.Flow.FlowJson, nil
}

func (expService *ExperimentServiceImpl) UpdateExperiment(req *restmodels.ProphecisExperimentPutRequest,
	username string, isSA bool) (*models.Experiment, error) {
	getLogger := logger.GetLogger()
	client := dss.GetDSSClient()

	// Get experiment From Database
	expe, err := repo.ExperimentRepo.Get(*req.ExperimentID, datasource.GetDB())
	if err != nil {
		getLogger.Errorf(err.Error())
		return expe, err
	}
	if expe == nil {
		log.Error("Update Service: experiment can not be found, check the experiment id")
		return nil, errors.New(" experiment can not be found, check the experiment id ")
	}

	//Permission Check
	err = PermissionCheck(username, expe.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	//Get Flow in DSS For Update Experiment Basic Info
	defaultLabel := "dev"
	flowData, err := client.GetOrchestrator(expe.DssFlowID, defaultLabel, username)
	if err != nil {
		log.Errorf("Error Get Flow in DSS: " + err.Error())
		return expe, err
	}

	//Save Flow in DSS
	data, err := saveFlow(expe, req.FlowJSON, username)
	//data, err := saveFlow(expe, req.FlowJSON, flowData.Flow.FlowEditLock, username)
	if err != nil {
		log.Errorf("Error Save Flow in DSS: " + err.Error())
		return expe, err
	}


	//Update Experiment in Transaction
	err = datasource.GetDB().Transaction(func(tx *gorm.DB) error {
		expe.DssFlowLastVersion = &data.FlowVersion
		expe.DssBmlLastVersion = &flowData.Flow.BmlVersion
		expe.ExpModifyTime = time.Now()
		err = repo.ExperimentRepo.Update(expe, tx)
		return err
	})
	if err != nil {
		getLogger.Errorf("Update Experiment Info Error:" + err.Error())
		return expe, err
	}
	return expe, err
}

func (expService *ExperimentServiceImpl) Import(data io.ReadCloser, experimentId int64, fileName *string,
	username string, createType string, desc string) (*restmodels.ProphecisExperimentIDResponse, error) {
	if experimentId < 0 {
		log.Printf("experiment id must be larger than or equal to 0 in experiment import, experimentId: %d\n", experimentId)
		return nil, fmt.Errorf("experiment id must be larger than or equal to 0 in experiment import, experimentId: %d", experimentId)
	}
	getLogger := logger.GetLogger()
	folderUid := uuid.New().String()
	flowFolder := "./importFolder/" + folderUid
	ossStorageFolder := "/data/oss-storage"
	flowFileFolder := "flow"
	flowDefYaml := "flow-definition.yaml"
	var flowFileName string
	uid := uuid.New().String()
	if fileName == nil || *fileName == "" {
		flowFileName = fmt.Sprintf("%v.zip", uid)
	} else {
		flowFileName = *fileName
	}
	flowFolderFile, err := os.Open(flowFolder)
	defer flowFolderFile.Close()
	if err != nil && os.IsNotExist(err) {
		//create file
		err := os.MkdirAll(flowFolder, os.ModePerm)
		if err != nil {
			log.Println("create folder flowFolder failed, ", err)
			return nil, err
		}
	}
	//save zip file
	importZipFilePath := fmt.Sprintf("%v/%v", flowFolder, flowFileName)
	log.Println("import zip file, path:", importZipFilePath)
	//check file is not zip
	if !strings.HasSuffix(importZipFilePath, ".zip") {
		log.Println("import zip file is not zip file, path:", importZipFilePath)
		return nil, errors.New("import zip file is not zip file, path:" + importZipFilePath)
	}
	//un's zip file
	unZipPathOfFlowZip := strings.TrimRight(importZipFilePath, ".zip")
	log.Println("un zip, path:", unZipPathOfFlowZip)
	//create flow.zip file and write
	flowZipFile, err := os.Create(importZipFilePath)
	if err != nil {
		log.Println("create flow.zip failed, err:", err)
		return nil, err
	}
	//read data
	fileBytes, err := ioutil.ReadAll(data)
	_, err = flowZipFile.Write(fileBytes)
	if err != nil {
		log.Println("flow.zip write failed, ", err)
		return nil, err
	}
	//un zip file
	err = archiver.Unarchive(importZipFilePath, unZipPathOfFlowZip)
	if err != nil {
		log.Println("un archive file failed, ", err)
		return nil, err
	}
	//code.zip, flow-definition.yaml's dir
	flowFolderPath := fmt.Sprintf("%v/%v", unZipPathOfFlowZip, flowFileFolder)
	log.Println("flowFolderPath: ", flowFolderPath)
	yamlPath := fmt.Sprintf("%v/%v", flowFolderPath, flowDefYaml)
	log.Println("yamlPath: ", yamlPath)
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Println("read yaml file failed, ", err)
		return nil, err
	}
	flowDefYamlModel := &models.ExportFlowDefinitionYamlModel{}
	err = yaml.Unmarshal(yamlBytes, flowDefYamlModel)
	if err != nil {
		log.Println("yaml un marshal failed, ", err)
		return nil, err
	}

	//查询experiment 是否存在，不存在就创建，存在就更新
	createFlag := false
	if experimentId != 0 {
		_, _, err = expService.GetExperiment(experimentId, username, true) //todo check
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Printf("fail to get experiment, experimentId: %d, username: %s, err: %v\n", experimentId, username, err)
				return nil, err
			}
			createFlag = true
		} else {
			m := make(map[string]interface{})
			m["exp_name"] = flowDefYamlModel.ExportExp.ExpName
			m["exp_desc"] = flowDefYamlModel.ExportExp.ExpDesc
			m["exp_modify_time"] = time.Now()
			if err := repo.ExperimentRepo.UpdateByMap(experimentId, m, datasource.GetDB()); err != nil {
				log.Printf("fail to update experiment, experimentId: %d, m: %s, err: %v\n", experimentId, m, err)
				return nil, err
			}
		}
	}else{
		createFlag = true
	}


	if createFlag {
		// 构建 createExp request
		var createFlowReq restmodels.ProphecisExperimentRequest
		if createType == "WTSS" {
			createFlowReq = restmodels.ProphecisExperimentRequest{
				ExpID:   flowDefYamlModel.ExportExp.ExpId,
				ExpName: flowDefYamlModel.ExportExp.ExpName,
				ExpDesc: flowDefYamlModel.ExportExp.ExpDesc,
			}
		} else {
			createFlowReq = restmodels.ProphecisExperimentRequest{
				ExpName: flowDefYamlModel.ExportExp.ExpName,
				ExpDesc: flowDefYamlModel.ExportExp.ExpDesc,
			}
		}
		// 调用di 创建 exp 接口
		createFlowResp, err := expService.CreateExperiment(&createFlowReq, username, createType, desc)
		if err != nil {
			getLogger.Debugf("expService.CreateFlow failed")
			return nil, err
		}
		experimentId = createFlowResp.ID
	}

	//为什么这种方式解析, 因为flowJson 结构不确定, 所以只修改对应的node 的 codePath 属性.
	flowJson := flowDefYamlModel.ExportDSSFlow.FlowJson
	//解析出flowJson第一层,
	var flowObj map[string]json.RawMessage
	err = json.Unmarshal([]byte(flowJson), &flowObj)
	if err != nil {
		getLogger.WithError(err).Errorf("parser flowObj failed")
		return nil, err
	}
	//check flow json
	err = checkFlowJson(flowObj)
	if err != nil {
		log.Println("check flow json err, ", err)
		return nil, err
	}
	//解析flowJson 的 nodes 属性
	var nodesObj []map[string]json.RawMessage
	err = json.Unmarshal(flowObj["nodes"], &nodesObj)
	if err != nil {
		getLogger.WithError(err).Errorf("parser flowObj failed")
		return nil, err
	}
	err = json.Unmarshal(flowObj["nodes"], &nodesObj)
	if err != nil {
		log.Println("nodes json unmarshal err,", err)
		return nil, err
	}
	for _, nodeObj := range nodesObj {

		jobContentObj := map[string]json.RawMessage{}
		err := json.Unmarshal(nodeObj["jobContent"], &jobContentObj)
		if err != nil {
			log.Error("job content un marshal failed,", err)
			return nil, err
		}
		if string(jobContentObj["mlflowJobType"]) != "DistributedModel" {
			break
		}

		//check manifest
		maniFestObj, err := checkManifest(jobContentObj)
		if err != nil {
			log.Println("check manifest err, ", err)
			return nil, err
		}
		//check code_selector
		if _, codeSelectCheck := maniFestObj["code_selector"]; codeSelectCheck {
			codeSelect := maniFestObj["code_selector"]
			//if code_selector is codeFile
			if strings.Trim(codeSelect.(string), `"`) == "codeFile" {
				ctx := context.TODO()
				sClient, err := storageClient.NewStorage()
				codeZipFile, err := os.Open(fmt.Sprintf("%v/%v", flowFolderPath, flowDefYamlModel.CodePaths[strings.Trim(string(nodeObj["id"]), `"`)]))
				if err != nil {
					log.Println("open codeZipFile failed, ", err)
					return nil, err
				}
				_, err = os.Open(ossStorageFolder)
				if err != nil && os.IsNotExist(err) {
					//create file
					err := os.MkdirAll(ossStorageFolder, os.ModePerm)
					if err != nil {
						log.Println("create oss-storage folder failed, ", err)
						return nil, err
					}
				}
				codeZipInfo, err := codeZipFile.Stat()
				if err != nil {
					log.Println("open codeZipFile failed, ", err)
					return nil, err
				}
				ossStorageFile, err := os.Create(fmt.Sprintf("%v/%v", ossStorageFolder, codeZipInfo.Name()))
				if err != nil {
					log.Println("create ossStorageFile failed, ", err)
					return nil, err
				}
				codeZipFileBytes, err := ioutil.ReadAll(codeZipFile)
				if err != nil {
					log.Println("codeZipFile read failed, ", err)
					return nil, err
				}
				_, err = ossStorageFile.Write(codeZipFileBytes)
				if err != nil {
					log.Println("ossStorageFile write failed, ", err)
					return nil, err
				}
				err = ossStorageFile.Sync()
				if err != nil {
					log.Println("ossStorageFile sync failed, ", err)
					return nil, err
				}
				req := grpc_storage.CodeUploadRequest{
					FileName: codeZipInfo.Name(),
					HostPath: "/data/oss-storage/",
				}
				res, err := sClient.Client().UploadCode(ctx, &req)
				if err != nil {
					getLogger.Errorf("sClient.Client().UploadCode(ctx, &req) failed, %v", err.Error())
					return nil, err
				}
				delete(flowDefYamlModel.CodePaths, flowDefYamlModel.CodePaths[string(nodeObj["id"])])
				flowDefYamlModel.CodePaths[string(nodeObj["id"])] = res.S3Path
				err = sClient.Close()
				if err != nil {
					log.Println("sClient close err,", err)
					return nil, err
				}
				err = ossStorageFile.Close()
				if err != nil {
					log.Println("ossStorageFile close err,", err)
					return nil, err
				}
			}
		}
	}
	err = flowZipFile.Close()
	if err != nil {
		log.Println("flowZipFile close err,", err)
		return nil, err
	}
	flowJsonBytes, err := json.Marshal(&flowObj)
	latestJson := string(flowJsonBytes)
	flowJson = latestJson
	// update exp(flowJson)
	updateExpReq := restmodels.ProphecisExperimentPutRequest{
		ExperimentID: &experimentId,
		FlowJSON:     &latestJson,
	}
	expResp, err := expService.UpdateExperiment(&updateExpReq, username, true)
	if err != nil {
		getLogger.Errorf("expService.UpdateExperiment failed, %v", err.Error())
		return nil, err
	}
	getLogger.Debugf("expResp result: %+v", expResp)

	//创建 expRun, 执行flow
	var runResp *models.ExperimentRun
	runResp, err = ExperimentRunService.CreateExperimentRun(experimentId, &flowJson, username, true)
	if err != nil {
		getLogger.Errorf("ExperimentRunService.Execute failed, %v", err.Error())
		return nil, err
	}
	err = os.RemoveAll(fmt.Sprintf("%v/%v", flowFolderPath, uid))
	if err != nil {
		log.Println("remove all failed, ", err)
		return nil, err
	}
	getLogger.Debugf("runResp result: %+v", runResp)
	return &restmodels.ProphecisExperimentIDResponse{
		ID:      &experimentId,
		ExpName: *flowDefYamlModel.ExportExp.ExpName,
	}, nil
}

func (expService *ExperimentServiceImpl) Export(expId int64, username string, isSA bool) (*os.File, error) {
	folderUid := uuid.New().String()
	fileZipFolder := "./exportFolder/" + folderUid
	flowFolderName := "flow"
	flowDefYml := "flow-definition.yaml"
	expTags := map[string]string{}
	experiment, flowJson, err := ExperimentService.GetExperiment(expId, username, true)
	if err != nil {
		log.Println("GetExperiment err, ", err)
		return nil, err
	}
	if experiment == nil {
		log.Println("Export err, experiment is nil")
		return nil, errors.New("Export err, experiment is nil")
	}
	//Permission Check
	err = PermissionCheck(username, experiment.ExpCreateUserID, nil, isSA)
	if err != nil {
		log.Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	//get expriment's tag
	for _, expTag := range experiment.TagList {
		expTags[strconv.Itoa(int(expTag.ID))] = *expTag.ExpTag
	}
	if len(*flowJson) <= 0 {
		return nil, errors.New("flowJson is null")
	}
	uid := uuid.New().String()
	flowZipId := fmt.Sprintf("%v-%v", experiment.DssFlowID, uid)
	flowFolderPath := fmt.Sprintf("%v/%v", fileZipFolder, flowZipId)
	flowPath := fmt.Sprintf("%v/%v", flowFolderPath, flowFolderName)
	err = os.MkdirAll(flowPath, os.ModePerm)
	if err != nil {
		log.Println("create flowPath failed:", err)
		return nil, err
	}
	flowObj := map[string]json.RawMessage{}
	nodesObj := []map[string]json.RawMessage{}
	codesPath := map[string]string{}
	err = json.Unmarshal([]byte(*flowJson), &flowObj)
	if err != nil {
		log.Println("flow json un marshal failed:", err)
		return nil, err
	}
	//check flow json
	err = checkFlowJson(flowObj)
	if err != nil {
		log.Println("check flow json err, ", err)
		return nil, err
	}
	err = json.Unmarshal(flowObj["nodes"], &nodesObj)
	if err != nil {
		log.Println("nodes un marshal err, ", err)
		return nil, err
	}
	for _, nodeObj := range nodesObj {

		//Parse node content to job conente
		jobContentObj := map[string]json.RawMessage{}
		err := json.Unmarshal(nodeObj["jobContent"], &jobContentObj)
		if err != nil {
			log.Error("job content un marshal failed,", err)
			return nil, err
		}
		if string(jobContentObj["mlflowJobType"]) != "DistributedModel" {
			break
		}

		//check manifest
		maniFestObj, err := checkManifest(jobContentObj)
		if err != nil {
			log.Println("check manifest err, ", err)
			return nil, err
		}
		//check code_selector
		if _, codeSelectCheck := maniFestObj["code_selector"]; !codeSelectCheck {
			return nil, errors.New("codeSelector is null")
		}
		codeSelect := maniFestObj["code_selector"]
		//if code_selector is codeFile
		if strings.Trim(codeSelect.(string), `"`) == "codeFile" {
			//download from s3 to codeZipPath
			s3Client, err := storageClient.NewStorage()
			todo := context.TODO()
			ctx, _ := context.WithCancel(todo)
			codePath := strings.Trim(maniFestObj["code_path"].(string), `""`)
			if err != nil {
				log.Println("create storerClient err, ", err)
				return nil, err
			}
			// if codePath is null or ""
			if codePath != "" {
				downloadReq := grpc_storage.CodeDownloadRequest{
					S3Path: codePath,
				}
				res, err := s3Client.Client().DownloadCode(ctx, &downloadReq)
				if err != nil {
					log.Println("s3Client download code failed, ", err)
					return nil, err
				}
				if res == nil {
					log.Println("s3Client download code res is null")
					return nil, errors.New("s3Client download code res is null")
				}
				nodeId := strings.Trim(string(nodeObj["id"]), `"`)
				codeZipPath := fmt.Sprintf("%v/%v", flowPath, res.FileName)
				codeZipFile, err := os.Create(codeZipPath)
				if err != nil {
					log.Println("codeZipFile create failed, ", err)
					return nil, err
				}
				codeBytes, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", res.HostPath, res.FileName))
				if err != nil {
					log.Println("code zip file read failed, ", err)
					return nil, err
				}
				count, err := codeZipFile.Write(codeBytes)
				if err != nil {
					log.Println("code zip file write failed:", err)
					return nil, err
				}
				if count == 0 {
					log.Printf(fmt.Sprintf("code zip file write byte is zero error, codebytes:%v", len(codeBytes)))
					log.Println("code zip file write byte is zero error, filepath" + res.HostPath + res.FileName)
				}
				err = codeZipFile.Sync()
				if err != nil {
					log.Println("codeZipFile sync failed, ", err)
					return nil, err
				}
				codeZipInfo, _ := codeZipFile.Stat()
				codesPath[nodeId] = codeZipInfo.Name()
				err = s3Client.Close()
				if err != nil {
					log.Println("s3Client close failed, ", err)
					return nil, err
				}
				err = codeZipFile.Close()
				if err != nil {
					log.Println("codeZipFile close failed, ", err)
					return nil, err
				}
			}
		}
	}
	yamlFile, err := os.Create(fmt.Sprintf("%v/%v", flowPath, flowDefYml))
	if err != nil {
		log.Println("yamlFile create failed, ", err)
		return nil, err
	}
	yamlModel := models.ExportFlowDefinitionYamlModel{
		ExportExp: models.ExportExp{
			ExpId:   expId,
			ExpName: experiment.ExpName,
			ExpDesc: experiment.ExpDesc,
			ExpTag:  expTags,
		},
		ExportDSSFlow: models.ExportDSSFlow{
			FlowJson: *flowJson,
		},
		CodePaths: codesPath,
	}
	yamlBytes, err := yaml.Marshal(yamlModel)
	if err != nil {
		log.Println("yaml marshal failed, ", err)
		return nil, err
	}
	_, err = yamlFile.Write(yamlBytes)
	if err != nil {
		log.Println("yamlFile write failed, ", err)
		return nil, err
	}
	err = yamlFile.Close()
	if err != nil {
		log.Println("yamlFile close failed, ", err)
		return nil, err
	}
	zipPath := fmt.Sprintf("%v", flowPath)
	zipName := fmt.Sprintf("%v.zip", flowFolderPath)
	err = archiver.Archive([]string{zipPath}, zipName)
	if err != nil {
		log.Println("zip archiver failed, ", err)
		return nil, err
	}
	err = os.RemoveAll(flowFolderPath)
	if err != nil {
		log.Println("clear folder failed, ", err)
		return nil, err
	}
	zipPath = fmt.Sprintf("%v/%v.zip", fileZipFolder, flowZipId)
	return os.Open(zipPath)
}

//func Get() ExperimentServiceImpl {
//	return ExperimentService
//}

//
//返回
func updateFlowJsonCodePath(flowJson string) error {
	return nil
}
func getCreateInfo(username string) (int64, int64, error) {
	var wsID int64
	var projectID int64
	client := dss.GetDSSClient()

	workspaceList, err := client.GetWorkspaces(username)
	if err != nil {
		log.Errorf(err.Error())
		return wsID, projectID, err
	}

	var workspace *dss.Workspace
	for _, ws := range workspaceList {
		if ws.Name == "bdapWorkspace" {
			workspace = &ws
			break
		}
	}
	if workspace == nil {
		log.Errorf("workspace named bdapWorkspace not exists")
		return wsID, projectID, errors.New("workspace named bdapWorkspace not exists")
	}

	wsID = workspace.Id
	//check cooperate project
	projectName := "MLSS" + "_" + username
	isExisted, project, err := client.CheckProjectExist(wsID ,projectName, username)
	if err != nil {
		return wsID, projectID, errors.New(fmt.Sprintf("check cooperate project failed, %v", err))
	}
	// create mlss cooperate project
	if !isExisted {
		projectParams := dss.ProjectParam{
			ProjectName:    "MLSS" + "_" + username ,
			ProjectDesc:    "MLSS",
			ProjectProduct: "MLSS",
		}
		project, err := client.CreateProject(&projectParams, username, nil ,workspace.Id)
		if err != nil {
			return wsID, -1, err
		}
		return wsID, project.Id, nil
	}

	//get projectID

	if err != nil {
		return wsID, projectID, errors.New(fmt.Sprintf("get cooperate project version id failed, %v", err))
	}
	wsID = workspace.Id
	return wsID, project.Id, err
}

func createFlow(expName string, projectID int64, workspaceID string, expType string ,username string) (*dss.OrchestratorData, *dss.CreateOrchestratorResponse, error) {
	// Init
	client := dss.GetDSSClient()
	orchestratorWays := []string{"pom_work_flow_DAG"}
	orchestratorMode := "pom_work_flow"
	label := dss.DSSLabel{
		Route: "dev",
	}
	//TODO: If type is WTSS, change label/way/mode
	if expType == "WTSS" {

	}

	//Create Flow By DSS Client
	createPersonalWorkflowParams := dss.CreateOrchestratorRequest{
		OrchestratorName: expName + "-" + uuid.New().String(),
		Uses:             "",
		ProjectId:        projectID,
		WorkspaceId:      workspaceID,
		OrchestratorWays: orchestratorWays,
		OrchestratorMode: orchestratorMode,
		Labels:           label,
		Description:      "MLSS experiment flow",
	}
	orchestratorResponse, err := client.CreateOrchestrator(createPersonalWorkflowParams, username)
	if err != nil {
		log.Errorf("GetFlowData, CreaatePersonalWorkFlow Error: %s" + err.Error())
		return nil, orchestratorResponse, err
	}

	//Get Flow ID
	orchestratorVO, err  := client.OpenOrchestrator(orchestratorResponse.OrchestratorId,label,"bdapWorkspace",username)
	if err != nil {
		log.Errorf("GetFlowId Error: %s" + err.Error())
		return nil, orchestratorResponse, err
	}


	// GET Flow Data From DSS
	orchestratorData, err := client.GetOrchestrator(orchestratorVO.DssOrchestratorVersion.AppId, "dev",  username)
	if err != nil {
		log.Errorf("GetFlowData, GetFlow Error: %s" + err.Error())
		return nil, orchestratorResponse, err
	}
	return orchestratorData, orchestratorResponse, err
}

//func saveFlow(experiment *models.Experiment, flowJson *string, flowEditLock string, username string) (*dss.SaveFlowData, error) {
func saveFlow(experiment *models.Experiment, flowJson *string, username string) (*dss.SaveFlowData, error) {
	client := dss.GetDSSClient()
	//TODO: Check Label
	label := dss.DSSLabel{
		Route: "dev",
	}
	//Save Flow In DSS
	params := dss.SaveWorkflowParams{
		Id:               experiment.DssFlowID,
		Json:             *flowJson,
		//TODO: ADD FLOW Name
		WorkspaceName: *experiment.DssDssFlowName,
		ProjectName: *experiment.DssDssProjectName,
		Label: label,
		//FlowEditLock: flowEditLock,
		//IsNotHaveLock: true,
	}
	data, err := client.SaveFlow(params, username)
	return data, err
}

func checkFlowJson(flowObj map[string]json.RawMessage) error {
	//check flowObj's nodes
	if _, check := flowObj["nodes"]; !check {
		log.Println("flowJson nodes is null")
		return errors.New("flowJson nodes is null")
	}
	//check nodes
	nodesObj := []map[string]json.RawMessage{}
	err := json.Unmarshal(flowObj["nodes"], &nodesObj)
	if err != nil {
		log.Println("nodes un marshal err, ", err)
		return err
	}
	for _, nodeObj := range nodesObj {
		//check node's jobContent
		if _, check := nodeObj["jobContent"]; !check {
			return errors.New("jobContent is null")
		}
		//get jobcontent
		jobContentObj := map[string]json.RawMessage{}
		err = json.Unmarshal(nodeObj["jobContent"], &jobContentObj)
		if err != nil {
			log.Println("job content un marshal failed,", err)
			return err
		}
		//get jobContent's ManiFest
		if _, check := jobContentObj["ManiFest"]; !check {
			log.Println("manifest's code_selector is null")
			return err
		}
	}
	return nil
}

func checkManifest(jobContentObj map[string]json.RawMessage) (map[string]interface{}, error) {
	//jobContentObj := map[string]json.RawMessage{}
	//err := json.Unmarshal(nodeObj["jobContent"], &jobContentObj)
	//if err != nil {
	//	log.Error("job content un marshal failed,", err)
	//	return nil, err
	//}
	maniFestObj := map[string]interface{}{}
	err := json.Unmarshal(jobContentObj["ManiFest"], &maniFestObj)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return maniFestObj, nil
}

func (expService *ExperimentServiceImpl) ExportDSS(expId int64, username string) (*restmodels.ProphecisExperimentDssResponse, error) {
	//init dss client
	dssClient := dss.GetDSSClient()
	exportFile, err := ExperimentService.Export(expId, username, true)
	if err != nil {
		log.Println("export file failed, ", err)
		return nil, err
	}
	if exportFile == nil {
		log.Println("export file failed, exportFile is nil")
		return nil, errors.New("export file failed, exportFile is nil")
	}
	uploadData := dss.UploadData{
		UploadFile:   exportFile,
		UploadSystem: "MLSS",
	}
	res, err := dssClient.Upload(&uploadData, username)
	if err != nil {
		log.Println("upload failed, ", err)
		return nil, err
	}
	data := restmodels.ProphecisExperimentDssResponse{
		ResourceID: res.ResourceId,
		Version:    res.Version,
	}
	return &data, nil
}

func (expService *ExperimentServiceImpl) ImportDSS(experimentId int64, resourceId string, version *string, desc *string, user string) (*restmodels.ProphecisImportExperimentDssResponse, error) {
	//init dss client
	dssClient := dss.GetDSSProClient()
	var download *dss.DowloadData
	if reflect2.IsNil(version) == false {
		download = &dss.DowloadData{
			ResourceId: resourceId,
			Version:    *version,
		}
	} else {
		download = &dss.DowloadData{
			ResourceId: resourceId,
		}
	}
	reply, err := dssClient.Download(download, user)
	if err != nil || reply == nil {
		log.Println("dssClient download failed, ", err)
		return nil, nil
	}
	reader := bytes.NewReader(reply)
	readerCloser := ioutil.NopCloser(reader)
	if desc == nil {
		nilDesc := "NULL"
		desc = &nilDesc
	}
	res, err := expService.Import(readerCloser, experimentId, nil, user, "WTSS", *desc)
	if err != nil {
		log.Println("read download err, ", err)
		return nil, nil
	}
	return &restmodels.ProphecisImportExperimentDssResponse{
		ExpID: *res.ID,
	}, nil
}

func PermissionCheck(execUserName string, entityUserId int64, groupId *string, isSuperAdmin bool) error{
	if isSuperAdmin {
		return nil
	}
	db := datasource.GetDB()
	user, err := repo.UserRepo.Get(&execUserName, db)
	if err != nil {
		return errors.New("Permission Check Error, Get user err: " + err.Error())
	}
	if user.ID != entityUserId{
		return commons.PermissionDeniedError
	}

	return nil
}

func createMLFlowExperiment(expName string, username string) (*string, error) {
	client := mlflow.GetMLFlowClient()
	currentTime := time.Now().Unix()
	currentTimeStr := strconv.FormatInt(currentTime,10)

	request := mlflow.CreateExperimentRequest{
		Name: expName+"_"+ username +"_"+ currentTimeStr,
	}

	result, err := client.CreateExperiment(&request)
	if err !=nil{
		return nil,err
	}
	return &(result.ExperimentId),nil
}


