package service_v2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	api "github.com/kubeflow/pipelines/backend/api/v2beta1/go_client"
	"github.com/mholt/archiver/v3"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"webank/DI/commons/constants"
	"webank/DI/commons/logger"
	loggerv2 "webank/DI/commons/logger/v2"
	datasource "webank/DI/pkg/datasource/mysql"
	"webank/DI/pkg/v2/client"
	models "webank/DI/pkg/v2/model"
	"webank/DI/pkg/v2/repo"
	thirdRepo "webank/DI/pkg/v2/repo/third"
	"webank/DI/pkg/v2/util"
	"webank/DI/restapi/api_v2/restmodels"
	storageClient "webank/DI/storage/client"
	"webank/DI/storage/storage/grpc_storage"
)

var ExperimentService ExperimentServiceIF = &ExperimentServiceImpl{}
var log = logger.GetLogger()

type ExperimentServiceIF interface {
	CreateExperiment(currentUserName string, clusterType string, request *restmodels.CreateExperimentRequest) (*models.Experiment, error)
	CopyExperiment(currentUserName string, clusterType string, sourceExpId string, request *restmodels.PatchExperimentRequest) (*models.Experiment, error)
	ExportExperimentFlowJson(currentUserName string, clusterType string, expId string) (*os.File, error)
	DeleteExperiment(currentUserName string, clusterType string, expId string) error
	DeleteExperimentVersion(currentUserName string, clusterType string, expId string, versionName string) error
	ListExperiments(ctx context.Context, currentUserName string, clusterType string, name *string, tag *string, groupID *string, groupName *string, sourceSystem *string,
		createUser *string, createTimeSt *string, createTimeEd *string,
		updateTimeSt *string, updateTimeEd *string, dssProjectName, dssFlowName *string,
		size *int64, page *int64) ([]*models.Experiment, int64, error)
	ListExperimentVersionNames(currentUserName string, clusterType string, expId string, canDeploy *bool) ([]string, error)

	CreateExperimentVersion(currentUserName string, clusterType string, expId string, newDescription string) (archivedVersionName string, latestVersionName string, err error)
	GetExperimentVersion(currentUserName string, clusterType string, expId string, version_name string) (*models.Experiment, error)
	GetExperimentVersionFlowJson(currentUserName string, clusterType string, expId string, versionName string) (string, string, string, error)
	GetExperimentVersionGlobalVariablesStr(currentUserName string, clusterType string, expId string, versionName string) (globalVariablesStr string, err error)

	PatchExperimentVersion(currentUserName string, clusterType string, expId string, versionName string, request *restmodels.PatchExperimentRequest) error
	ListExperimentVersions(currentUserName string, clusterType string, expId *string, versionName *string,
		updateTimeSt *string, updateTimeEd *string, size *int64, page *int64) ([]*models.Experiment, int64, error)

	UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error)
}

type ExperimentServiceImpl struct {
}

func (e ExperimentServiceImpl) UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error) {
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

func (e ExperimentServiceImpl) GetExperimentVersionGlobalVariablesStr(currentUser string, clusterType string, expId string, versionName string) (globalVariablesStr string, err error) {
	flowJson, _, _, err := e.GetExperimentVersionFlowJson(currentUser, clusterType, expId, versionName)
	if err != nil {
		log.Errorf("GetExperimentVersionFlowJson failed: %+v", err)
		return "", err
	}
	if flowJson == "" {
		return "", errors.New("GetExperimentVersionFlowJson is empty")
	}
	var flowJsonMap map[string]json.RawMessage
	err = json.Unmarshal([]byte(flowJson), &flowJsonMap)
	if err != nil {
		log.Error("parse flow json to map failed: ", err.Error())
		return "", err
	}
	if _, check := flowJsonMap["global_variables"]; !check {
		return "", nil
	}
	return string(flowJsonMap["global_variables"]), nil
}

func (e ExperimentServiceImpl) GetExperimentVersionFlowJson(currentUserName string, clusterType string, expId string, versionName string) (flowjson string, kfpPipelineId, kfpPipelineVersion string, err error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return "", "", "", err
	}
	if !canAccess {
		return "", "", "", AccessError{}
	}

	// 获取实验的flowjosn, 如果为空就报错返回
	var exp *models.Experiment
	if versionName == "" {
		exp, err = repo.ExperimentRepo.GetLatestVersion(expId)
	} else {
		exp, err = repo.ExperimentRepo.GetByVersion(expId, versionName)
	}
	if err != nil {
		log.Errorf("ExperimentRepo.GetLatestVersion/GetByVersion failed: +%v", err)
		return "", "", "", err
	}
	if len(exp.KfpPipelineVersionId) <= 0 {
		log.Warnf("flowJson is empty(kfpPipelineVersionId is empty)")
		return "", "", "", errors.New("flowJson is empty(kfpPipelineVersionId is empty)")
	}
	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Errorf("client.GetKFPClient() failed: +%v", err)
		return "", "", "", err
	}
	getPipLineVersionRequest := &api.GetPipelineVersionRequest{PipelineId: exp.KfpPipelineId, PipelineVersionId: exp.KfpPipelineVersionId}
	pipelineVersion, err := kfpClient.PipelineServiceClient.GetPipelineVersion(context.Background(), getPipLineVersionRequest)
	if err != nil {
		log.Errorf("kfpClient.GetPipelineVersion() failed: +%v", err)
		return "", "", "", err
	}
	if len(pipelineVersion.FlowJson) <= 0 {
		return "", "", "", errors.New("flowJson is empty")
	}
	return pipelineVersion.FlowJson, exp.KfpPipelineId, exp.KfpPipelineVersionId, nil
}

func (e ExperimentServiceImpl) ListExperimentVersionNames(currentUserName string, clusterType string, expId string, canDeploy *bool) ([]string, error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return []string{}, err
	}
	if !canAccess {
		return []string{}, AccessError{}
	}

	result, err := repo.ExperimentRepo.ListVersionNames(expId)
	if err != nil {
		log.Errorf("repo.ExperimentRepo.ListVersionNames failed: %+v", err)
		return []string{}, err
	}
	if canDeploy != nil && *canDeploy {
		lenResult := len(result)
		if lenResult >= 1 {
			return result[1:], nil
		}
	}
	return result, nil
}

func (e ExperimentServiceImpl) ListExperiments(ctx context.Context, currentUserName string, clusterType string, name *string, tag *string,
	groupID *string, groupName *string, sourceSystem *string,
	createUser *string, createTimeSt *string, createTimeEd *string,
	updateTimeSt *string, updateTimeEd *string, dssProjectName, dssFlowName *string,
	size *int64, page *int64) ([]*models.Experiment, int64, error) {
	log := loggerv2.GetLogger(ctx)

	// (1) 获取用户权限相关的数据
	userDao, err := thirdRepo.UserRepo.GetUser(currentUserName, clusterType)
	if err != nil {
		log.Errorf("thirdRepo.UserRepo.GetUser(%s) failed: %+v", currentUserName, err)
		return nil, 0, err
	}
	var canAccessGroupIDs []string
	isSA := thirdRepo.UserRoleRepo.IsSA(userDao.UserID)
	log.Infof("the currentUser(%s) isSA %v", currentUserName, isSA)
	if !isSA {
		groupIDs, err := thirdRepo.UserRoleRepo.GetGroupIdsByUserId(userDao.UserID)
		if err != nil {
			log.Errorf("thirdRepo.UserRoleRepo.GetGroupIdsByUserId(%s) failed: %+v", currentUserName, err)
			return nil, 0, err
		}
		canAccessGroupIDs = groupIDs
	}
	if !isSA && len(canAccessGroupIDs) == 0 {
		return nil, 0, nil
	}

	if groupID != nil {
		if !isSA {
			canAccess := false
			for _, item := range canAccessGroupIDs {
				if item == *groupID {
					canAccess = true
					break
				}
			}
			if !canAccess {
				dirtyData := "not_exist_group_id"
				groupID = &dirtyData
			}
		}
	}

	rOffset, rSize := GetOffSetAndSize(page, size)
	exps, total, err := repo.ExperimentRepo.GetAllByOffset(currentUserName, nil, name, tag, []string{string(models.ExperimentStatusActive)},
		canAccessGroupIDs, groupID, groupName, nil, sourceSystem, clusterType, createUser, createTimeSt, createTimeEd,
		updateTimeSt, updateTimeEd, dssProjectName, dssFlowName, rSize, rOffset)
	if err != nil {
		log.Errorf("ExperimentRepo.GetAllByOffset failed: %+v", err)
		return nil, 0, err
	}
	if exps == nil || len(exps) == 0 {
		log.Errorf("the repo.ExperimentRepo.GetAllByOffset result experiments is nil or empty")
	}
	return exps, total, nil
}

func (e ExperimentServiceImpl) DeleteExperiment(currentUserName string, clusterType string, expId string) error {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	exp, err := repo.ExperimentRepo.GetLatestVersion(expId)
	if err != nil {
		log.Errorf("get GetLatestVersion(%s) failed: %+v", expId, err)
		return err
	}
	updateMap := make(map[string]interface{})
	updateMap["update_user"] = currentUserName
	updateMap["update_time"] = time.Now()
	updateMap["status"] = string(models.ExperimentStatusDeleted)
	err = repo.ExperimentRepo.UpdateByMap(expId, exp.VersionName, updateMap)
	if err != nil {
		log.Errorf("DeleteExperiment Error:" + err.Error())
		return err
	}
	return nil
}

func (e ExperimentServiceImpl) DeleteExperimentVersion(currentUserName string, clusterType string, expId string, versionName string) error {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	exp, err := repo.ExperimentRepo.GetLatestVersion(expId)
	if err != nil {
		log.Errorf("get GetLatestVersion(%s) failed: %+v", expId, err)
		return err
	}
	log.Infof("get exp is: %+v", exp)
	if exp.VersionName == versionName {
		return errors.New("cant delete latest experiment version")
	}

	/* todo(gaouanhe): 暂时注释掉，因为发现将v1版本删除掉后，在用v3版本去创建一个新版本v4的时候，kfp报找不到v3版本的kfpPipelineVersion
	if exp.KfpPipelineVersionId != "" {
		kfpClient, err := client.GetKFPClient()
		if err != nil {
			log.Error("get kfpClient error : " + err.Error())
			return err
		}
		deletePipelineVersionRequest := &api.DeletePipelineVersionRequest{
			PipelineId:        exp.KfpPipelineId,
			PipelineVersionId: exp.KfpPipelineVersionId,
		}
		_, err = kfpClient.PipelineServiceClient.DeletePipelineVersion(context.Background(), deletePipelineVersionRequest)
		if err != nil {
			log.Errorf("kfp DeletePipelineVersion failed: %+v", err)
			return err
		}
	}
	*/

	updateMap := make(map[string]interface{})
	updateMap["update_user"] = currentUserName
	updateMap["update_time"] = time.Now()
	updateMap["status"] = string(models.ExperimentStatusDeleted)
	err = repo.ExperimentRepo.UpdateByMap(expId, versionName, updateMap)
	if err != nil {
		log.Errorf("ExperimentRepo.UpdateByMap failed: %+v", err)
		return err
	}
	return nil
}

func (e ExperimentServiceImpl) ExportExperimentFlowJson(currentUserName string, clusterType string, expId string) (*os.File, error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}

	flowJson, kfpPipelineId, kfpPipelineVersionId, err := e.GetExperimentVersionFlowJson(currentUserName, clusterType, expId, "")
	if err != nil {
		log.Errorf("e.GetExperimentFlowJson failed: +%v", err)
		return nil, err
	}
	log.Infof("GetExperimentFlowJson result is len(flowJosn) is: %d; kfpPipelineId is: %s; kfpPipelineVersionId is: %s", len(flowJson), kfpPipelineId, kfpPipelineVersionId)
	// 创建临时的导出文件目录
	timeNowStr := time.Now().Format("2006-01-02-15-04-05")
	parentFolderPath := filepath.Join("./exportExperimentFlowJson", expId, kfpPipelineId, kfpPipelineVersionId)
	folderPath := filepath.Join(parentFolderPath, timeNowStr, "flowjson")
	log.Infof("the parentFolderPath is: %s", parentFolderPath)
	log.Infof("the folderPath is: %s", folderPath)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Errorf("create folderPath failed: %+v", err)
		return nil, err
	}
	filePath := filepath.Join(folderPath, "flowjson.json")
	flowjsonFile, err := os.Create(filePath)
	if err != nil {
		log.Errorf("flowjsonFile create failed: %+v", err)
		return nil, err
	}
	_, err = flowjsonFile.Write([]byte(flowJson))
	if err != nil {
		log.Errorf("flowjsonFile write failed: %+v", err)
		return nil, err
	}
	defer flowjsonFile.Close()
	zipPath := filepath.Join(parentFolderPath, fmt.Sprintf("%s.zip", timeNowStr))
	err = archiver.Archive([]string{folderPath}, zipPath)
	if err != nil {
		log.Println("zip archiver failed, ", err)
		return nil, err
	}

	return os.Open(zipPath)
}

func (e ExperimentServiceImpl) CopyExperiment(currentUserName string, clusterType string, sourceExpId string, request *restmodels.PatchExperimentRequest) (*models.Experiment, error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, sourceExpId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, sourceExpId, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}

	exp, err := repo.ExperimentRepo.GetLatestVersion(sourceExpId)
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return nil, err
	}
	log.Infof("the input request is: %+v", request)

	copyName := fmt.Sprintf("%s_copy", exp.Name)
	createExperimentRequest := &restmodels.CreateExperimentRequest{
		ClusterType:  &exp.ClusterType,
		Description:  exp.Description,
		GroupID:      exp.GroupID,
		GroupName:    exp.GroupName,
		Name:         &copyName,
		SourceSystem: &exp.SourceSystem,
		Tags:         strings.Split(exp.Tags, ","),
	}
	if exp.SourceSystem == restmodels.ExperimentWithoutPipelineSourceSystemDSS &&
		exp.DSSWorkspaceName != "" && exp.DSSProjectName != "" {
		createExperimentRequest.DssInfo = &restmodels.DSSInfo{
			exp.DSSFlowID,
			exp.DSSFlowName,
			exp.DSSFlowVersion,
			exp.DSSProjectID,
			exp.DSSProjectName,
			exp.DSSWorkspaceID,
			exp.DSSWorkspaceName,
		}
	}
	flowjson, _, _, err := e.GetExperimentVersionFlowJson(currentUserName, clusterType, sourceExpId, "")
	if err != nil {
		log.Errorf("GetExperimentVersionFlowJson failed: %+v", err)
	}
	if flowjson != "" {
		createExperimentRequest.FlowJSON = flowjson
	}
	if request != nil {
		if request.Description != "" {
			createExperimentRequest.Description = request.Description
		}
		if request.GroupName != "" {
			createExperimentRequest.GroupName = request.GroupName
		}
		if request.Name != "" {
			createExperimentRequest.Name = &request.Name
		}
		if len(request.Tags) != 0 {
			createExperimentRequest.Tags = request.Tags
		}
		if request.FlowJSON != "" {
			createExperimentRequest.FlowJSON = request.FlowJSON
		}
	}
	return e.CreateExperiment(currentUserName, clusterType, createExperimentRequest)
}

func (e ExperimentServiceImpl) PatchExperimentVersion(currentUserName string, clusterType string, expId string,
	versionName string, request *restmodels.PatchExperimentRequest) error {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	var exp *models.Experiment
	if versionName == "" {
		exp, err = repo.ExperimentRepo.GetLatestVersion(expId)
	} else {
		exp, err = repo.ExperimentRepo.GetByVersion(expId, versionName)
	}
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return err
	} else if exp == nil {
		log.Errorf("Experiment(%s) does not exist", expId)
		return fmt.Errorf("experiment(%s) does not exist", expId)
	}
	updateMap := make(map[string]interface{})
	// updateMapAll是所有的实验版本都会修改
	updateMapAll := make(map[string]interface{})
	if request.Name != "" {
		updateMap["name"] = request.Name
	}
	if request.Description != "" {
		updateMap["description"] = request.Description
	}
	if request.GroupID != "" {
		updateMapAll["group_id"] = request.GroupID
	}
	if request.GroupName != "" {
		updateMapAll["group_name"] = request.GroupName
	}
	if len(request.Tags) != 0 {
		updateMap["tags"] = strings.Join(request.Tags, ",")
	}
	if request.DeploySetting != nil {
		if exp.SourceSystem == restmodels.ExperimentWithoutPipelineSourceSystemMLSS {
			return InputParamsError{Field: "deploy_setting", Reason: "MLSS类型的实验不能进行发布设置"}
		}
		if exp.SourceSystem == restmodels.ExperimentWithoutPipelineSourceSystemDSS {
			if !exp.DSSPublished {
				return InputParamsError{Field: "deploy_setting", Reason: "未在BDAP DSS开发中心中进行发布的实验不能进行发布设置"}
			}
		}
		if exp.GroupID == "" {
			return InputParamsError{Field: "deploy_setting", Reason: "没有设置项目组的实验不能进行发布设置"}
		}
		updateMap["can_deploy_as_offline_service"] = request.DeploySetting.CanDeployAsOfflineService
		updateMap["can_deploy_as_offline_service_operator"] = currentUserName
	}
	// 如果flowjson不为空，还要额外对kfp的pipelineVersion进行操作
	if request.FlowJSON != "" {
		kfpClient, err := client.GetKFPClient()
		if err != nil {
			log.Error("get kfpClient error : " + err.Error())
			return err
		}
		// 如果exp的KfpPipelineVersionId为空，则表示还没有创建 kfp PipelineVersion
		if exp.KfpPipelineVersionId == "" {
			createPipelineVersionRequest := &api.CreatePipelineVersionRequest{
				PipelineId: exp.KfpPipelineId,
				PipelineVersion: &api.PipelineVersion{
					PipelineId:  exp.KfpPipelineId,
					DisplayName: fmt.Sprintf("%s-%s", exp.ID, exp.VersionName),
					FlowJson:    request.FlowJSON,
				},
			}
			createPipelineVersionResponse, err := kfpClient.PipelineServiceClient.CreatePipelineVersionV3(context.Background(), createPipelineVersionRequest)
			if err != nil {
				log.Error("kfp CreatePipelineVersionV3 error : " + err.Error())
				return err
			}
			updateMap["kfp_pipeline_version_id"] = createPipelineVersionResponse.PipelineVersionId
		} else {
			updatePipelineVersionRequest := &api.UpdatePipelineVersionRequest{
				PipelineId: exp.KfpPipelineId,
				PipelineVersion: &api.PipelineVersion{
					PipelineId:        exp.KfpPipelineId,
					PipelineVersionId: exp.KfpPipelineVersionId,
					DisplayName:       fmt.Sprintf("%s-%s", exp.ID, exp.VersionName),
					FlowJson:          request.FlowJSON,
				},
			}
			_, err = kfpClient.PipelineServiceClient.UpdatePipelineVersion(context.Background(), updatePipelineVersionRequest)
			if err != nil {
				log.Error("kfp UpdatePipelineVersion error : " + err.Error())
				return err
			}
		}
	}
	// 沒有要修改的
	if len(updateMap) == 0 && request.FlowJSON == "" {
		return nil
	}
	updateMap["update_user"] = currentUserName
	updateMap["update_time"] = time.Now().Format(constants.GoDefaultTimeFormat)
	if len(updateMapAll) != 0 {
		updateMapAll["update_user"] = currentUserName
		updateMapAll["update_time"] = time.Now().Format(constants.GoDefaultTimeFormat)
	}
	err = datasource.GetDB().Transaction(func(db *gorm.DB) error {
		err = repo.ExperimentRepo.UpdateByMap(expId, exp.VersionName, updateMap)
		if err != nil {
			log.Errorf("Update Experiment Error:" + err.Error())
			return err
		}
		if len(updateMapAll) != 0 {
			err = repo.ExperimentRepo.UpdateByMapIfEmpty(expId, "", updateMapAll)
			if err != nil {
				log.Errorf("Update All Experiment Error:" + err.Error())
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("Patch Experiment Error:" + err.Error())
		return err
	}
	return nil
}

func (e ExperimentServiceImpl) UpdateExperimentDeploySetting(currentUserName string, clusterType string, expId string, request *restmodels.DeploySetting) error {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return err
	}
	if !canAccess {
		return AccessError{}
	}

	exp, err := repo.ExperimentRepo.GetLatestVersion(expId)
	if err != nil {
		log.Errorf("GetLatestVersion(%s) failed: %+v", expId, err)
	}
	updateMap := make(map[string]interface{})
	updateMap["update_user"] = currentUserName
	updateMap["update_time"] = time.Now().Format(constants.GoDefaultTimeFormat)
	updateMap["can_deploy_as_offline_service"] = request.CanDeployAsOfflineService
	return repo.ExperimentRepo.UpdateByMap(expId, exp.VersionName, updateMap)
}

func (e ExperimentServiceImpl) ListExperimentVersions(currentUserName string, clusterType string, expId *string, versionName *string,
	updateTimeSt *string, updateTimeEd *string, size *int64, page *int64) ([]*models.Experiment, int64, error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, *expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, *expId, err)
		return nil, 0, err
	}
	if !canAccess {
		return nil, 0, AccessError{}
	}

	// 如果实验(最新版本的实验）已经被删除了，应该返回空数组
	exp, err := repo.ExperimentRepo.GetLatestVersion(*expId)
	if err != nil {
		log.Errorf("ExperimentRepo.GetLatestVersion(%s) failed: %+v", *expId, err)
	}
	if exp.ID == "" {
		log.Warnf("实验的最新版本已经被删除了，返回空数据")
		return nil, 0, err
	}
	rOffset, rSize := GetOffSetAndSize(page, size)
	exps, total, err := repo.ExperimentRepo.GetAllByOffset("", expId, nil, nil,
		[]string{string(models.ExperimentStatusActive), string(models.ExperimentStatusArchive)},
		nil, nil, nil, versionName, nil, "", nil, nil, nil,
		updateTimeSt, updateTimeEd, nil, nil, rSize, rOffset)
	if err != nil {
		log.Error("ExperimentRepo.GetAllByOffset failed: %+v", err)
		return nil, 0, err
	}

	if exps == nil || len(exps) == 0 {
		log.Errorf("the repo.ExperimentRepo.GetAllByOffset result experiments is nil or empty")
	}
	return exps, total, nil
}

// GetExperimentVersion 获取实验的某个具体的版本，如果version_name为空，则表示获取最新的版本
func (e ExperimentServiceImpl) GetExperimentVersion(currentUserName string, clusterType string, expId string, version_name string) (*models.Experiment, error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return nil, err
	}
	if !canAccess {
		return nil, AccessError{}
	}

	var exp *models.Experiment
	if version_name == "" {
		exp, err = repo.ExperimentRepo.GetLatestVersion(expId)
	} else {
		exp, err = repo.ExperimentRepo.GetByVersion(expId, version_name)
	}
	if err != nil {
		log.Errorf("Get Experiment Error:" + err.Error())
		return nil, err
	}
	return exp, nil
}

func (e ExperimentServiceImpl) CreateExperimentVersion(currentUserName string, clusterType string, expId string, newDescription string) (archivedVersionName string, latestVersionName string, err error) {
	canAccess, err := CanAccessExp(currentUserName, clusterType, expId)
	if err != nil {
		log.Errorf("CanAccessExp(%s, %s) failed: +%v", currentUserName, expId, err)
		return "", "", err
	}
	if !canAccess {
		return "", "", AccessError{}
	}

	exp, err := repo.ExperimentRepo.GetLatestVersion(expId)
	if err != nil {
		log.Errorf("ExperimentRepo.GetLatestVersion(%s) failed: %+v", expId, err)
		return "", "", err
	} else if exp.KfpPipelineVersionId == "" {
		log.Warnf("该实验（%s）没有任何工作流信息，另存为新版本没有多大意义", expId)
	}

	newExpVersionName := nextVersion(exp.VersionName)
	newExpKfpPipelineVersionId := ""
	if exp.KfpPipelineVersionId != "" {
		kfpClient, err := client.GetKFPClient()
		if err != nil {
			log.Error("get kfpClient error : " + err.Error())
			return "", "", err
		}
		getPipelineVersionRequest := &api.GetPipelineVersionRequest{
			PipelineId:        exp.KfpPipelineId,
			PipelineVersionId: exp.KfpPipelineVersionId,
		}
		getPipelineVersionResponse, err := kfpClient.PipelineServiceClient.GetPipelineVersion(context.Background(), getPipelineVersionRequest)
		if err != nil {
			log.Errorf("kfpClient.PipelineServiceClient.CreatePipeline failed: %+v", err)
			return "", "", err
		}
		createPipelineVersionRequestV3 := &api.CreatePipelineVersionRequest{
			PipelineId: exp.KfpPipelineId,
			PipelineVersion: &api.PipelineVersion{
				PipelineId:  exp.KfpPipelineId,
				DisplayName: fmt.Sprintf("%s-%s", exp.ID, newExpVersionName),
				FlowJson:    getPipelineVersionResponse.FlowJson,
			},
		}

		createPipelineVersionResponseV3, err := kfpClient.PipelineServiceClient.CreatePipelineVersionV3(context.Background(), createPipelineVersionRequestV3)
		if err != nil {
			log.Errorf("kfpClient.PipelineServiceClient.CreatePipeline failed: %+v", err)
			return "", "", err
		}
		newExpKfpPipelineVersionId = createPipelineVersionResponseV3.PipelineVersionId
	}

	// todo(gaoyuanhe): 该方法不鲁棒，后面考虑优化
	dss_published := exp.DSSPublished
	if strings.Contains(newDescription, "该版本是由dss开发中心发布而触发的") {
		dss_published = true
	}
	timeNow := time.Now()
	newExp := models.Experiment{
		ID:                        exp.ID,
		Name:                      exp.Name,
		Description:               "",
		GroupID:                   exp.GroupID,
		GroupName:                 exp.GroupName,
		VersionName:               newExpVersionName,
		SourceSystem:              exp.SourceSystem,
		ClusterType:               exp.ClusterType,
		CreateUser:                currentUserName,
		CreateTime:                timeNow,
		UpdateUser:                currentUserName,
		UpdateTime:                timeNow,
		Status:                    string(models.ExperimentStatusActive),
		Tags:                      exp.Tags,
		KfpPipelineId:             exp.KfpPipelineId,
		KfpPipelineVersionId:      newExpKfpPipelineVersionId,
		CanDeployAsOfflineService: exp.CanDeployAsOfflineService,
		LatestExecuteTime:         exp.LatestExecuteTime,
		DSSFlowID:                 exp.DSSFlowID,
		DSSFlowName:               exp.DSSFlowName,
		DSSFlowVersion:            exp.DSSFlowVersion,
		DSSProjectID:              exp.DSSProjectID,
		DSSProjectName:            exp.DSSProjectName,
		DSSWorkspaceID:            exp.DSSWorkspaceID,
		DSSWorkspaceName:          exp.DSSWorkspaceName,
		DSSPublished:              dss_published,
	}

	err = datasource.GetDB().Transaction(func(db *gorm.DB) error {
		// 保存新创建的实验版本
		err = repo.ExperimentRepo.Add(&newExp)
		if err != nil {
			return err
		}
		// 更新旧的实验版本信息
		expUpdateMap := make(map[string]interface{})
		expUpdateMap["update_user"] = currentUserName
		expUpdateMap["update_time"] = time.Now()
		expUpdateMap["dss_published"] = dss_published
		expUpdateMap["status"] = string(models.ExperimentStatusArchive)
		if newDescription != "" {
			expUpdateMap["description"] = newDescription
		}
		err = repo.ExperimentRepo.UpdateByMap(expId, exp.VersionName, expUpdateMap)
		if err != nil {
			return err
		}
		return nil
	})
	return exp.VersionName, newExpVersionName, nil
}

func (e ExperimentServiceImpl) CreateExperiment(currentUser string, clusterType string, request *restmodels.CreateExperimentRequest) (*models.Experiment, error) {
	expId := util.CreateUUID()
	// Experiment, Pipeline, PipelineVersion的关联关系
	// generate UUID as ExperimentID ---- ExperimentID ----- PipelineName ------ PipelineVersionName(PipelineName-VersionName)
	// 为什么不用experimentName作为PipelineName? 因为experimentName目前的设计是可以修改, 而PipelineName不可修改

	// 在kfp中创建 pipeline 和 pipelineVersion
	kfpPipelineName := expId
	initVersion := "v1"
	kfpPipelineVersionName := fmt.Sprintf("%s-%s", kfpPipelineName, initVersion)

	kfpClient, err := client.GetKFPClient()
	if err != nil {
		log.Error("get kfpClient error : " + err.Error())
		return nil, err
	}
	createPipelineRequest := &api.CreatePipelineRequest{
		Pipeline: &api.Pipeline{
			DisplayName: kfpPipelineName,
		},
	}
	createPipelineResponse, err := kfpClient.PipelineServiceClient.CreatePipeline(context.Background(), createPipelineRequest)
	if err != nil {
		log.Errorf("kfpClient.PipelineServiceClient.CreatePipeline failed: %+v", err)
		return nil, err
	}
	// 如果传入的flowjson不为空就创建 kfp的pipelineversion
	// 如果为空就暂时不创建，等到调用PatchExperiment接口的时候再创建
	// 因为kfp的CreatePipelineVersion接口暂时不支持传入空的flowjson
	var pipelineVersionId string
	var inputFlowJson string
	if request.FlowJSON != "" {
		inputFlowJson = request.FlowJSON
	} else if request.FlowJSONUploadID != "" {
		flowjsonPath, err := base64.StdEncoding.DecodeString(request.FlowJSONUploadID)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := ioutil.ReadFile(string(flowjsonPath))
		if err != nil || len(jsonBytes) == 0 {
			return nil, err
		}
		inputFlowJson = string(jsonBytes)
	}
	if inputFlowJson != "" {
		createPipelineVerionRequestV3 := &api.CreatePipelineVersionRequest{
			PipelineId: createPipelineResponse.PipelineId,
			PipelineVersion: &api.PipelineVersion{
				PipelineId:  createPipelineResponse.PipelineId,
				DisplayName: kfpPipelineVersionName,
				FlowJson:    inputFlowJson,
			},
		}
		createPipelineVerionResponseV3, err := kfpClient.PipelineServiceClient.CreatePipelineVersionV3(context.Background(), createPipelineVerionRequestV3)
		if err != nil {
			log.Errorf("kfpClient.PipelineServiceClient.CreatePipeline failed: %+v", err)
			return nil, err
		}
		pipelineVersionId = createPipelineVerionResponseV3.PipelineVersionId
	}
	if request.SourceSystem == nil {
		request.SourceSystem = StringPtr(restmodels.CreateExperimentRequestSourceSystemMLSS)
	}
	if request.ClusterType == nil {
		request.ClusterType = StringPtr(restmodels.CreateExperimentRequestClusterTypeBDAP)
	}

	newExperiment := &models.Experiment{
		ID:          expId,
		Name:        *request.Name,
		Description: request.Description,
		// todo(gaoyuanehe): 如果用户传入的GroupID和GroupName不匹配
		GroupID:              request.GroupID,
		GroupName:            request.GroupName,
		VersionName:          "v1",
		SourceSystem:         *request.SourceSystem,
		ClusterType:          *request.ClusterType,
		CreateUser:           currentUser,
		CreateTime:           time.Now(),
		UpdateUser:           currentUser,
		UpdateTime:           time.Now(),
		Status:               string(models.ExperimentStatusActive),
		Tags:                 strings.Join(request.Tags, ","),
		KfpPipelineId:        createPipelineResponse.PipelineId,
		KfpPipelineVersionId: pipelineVersionId,
	}
	if *request.SourceSystem == restmodels.CreateExperimentRequestSourceSystemDSS {
		newExperiment.DSSWorkspaceID = request.DssInfo.WorkspaceID
		newExperiment.DSSWorkspaceName = request.DssInfo.WorkspaceName
		newExperiment.DSSProjectID = request.DssInfo.ProjectID
		newExperiment.DSSProjectName = request.DssInfo.ProjectName
		newExperiment.DSSFlowID = request.DssInfo.FlowID
		newExperiment.DSSFlowName = request.DssInfo.FlowName
		newExperiment.DSSFlowVersion = request.DssInfo.FlowVersion
		newExperiment.DSSPublished = false
	}
	err = datasource.GetDB().Transaction(func(db *gorm.DB) error {
		err = repo.ExperimentRepo.Add(newExperiment)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("Create Experiment Error:" + err.Error())
		return nil, err
	}
	return newExperiment, nil
}

func GetOffSetAndSize(page *int64, size *int64) (int, int) {
	//sizeEnt := 0
	var rSize int
	var rOffSet int
	if size == nil {
		size = new(int64)
		*size = 10
	}
	if page == nil {
		page = new(int64)
		*page = 1
	}
	// 获取全部数据
	if *size == -1 {
		return 0, 0
	}
	if *size <= 0 {
		*size = 10
	}
	if *page <= 0 {
		*page = 1
	}
	rOffSet = int((*page - 1) * *size)
	rSize = int(*size)
	return rOffSet, rSize
}

// nextVersion : v1 => v2; v1.3 => v1.4; v0.3.6 => v0.3.7
// generated by chatGPT 3.5
func nextVersion(inputVersion string) string {
	// Split the input version string into components
	components := strings.Split(inputVersion, ".")

	// Extract the numeric part of the version if available
	var numericVersion int
	if len(components) > 1 {
		var err error
		numericVersion, err = strconv.Atoi(components[len(components)-1])
		if err != nil {
			// Handle error if conversion fails
			fmt.Println("Error:", err)
			return ""
		}
	} else {
		var err error
		numericVersion, err = strconv.Atoi(inputVersion[1:])
		if err != nil {
			// Handle error if conversion fails
			fmt.Println("Error:", err)
			return ""
		}
	}

	// Increment the numeric version
	nextNumericVersion := numericVersion + 1

	if len(components) > 1 {
		// Update the last component with the incremented version
		components[len(components)-1] = strconv.Itoa(nextNumericVersion)

		// Join the components back into a string
		nextVersionStr := strings.Join(components, ".")

		return nextVersionStr
	} else {
		nextVersionStr := fmt.Sprintf("v%d", nextNumericVersion)
		return nextVersionStr
	}

}

func StringPtr(input string) *string {
	return &input
}
