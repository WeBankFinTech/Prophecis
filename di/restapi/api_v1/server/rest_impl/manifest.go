/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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

package rest_impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	cc "webank/DI/commons/controlcenter/client"
	"webank/DI/commons/logger"
	"webank/DI/pkg/operators/tf-operator/apis/tensorflow/v1alpha1"
	"webank/DI/restapi/api_v1/restmodels"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	"webank/DI/commons/config"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ManifestV1 represents a manifest used to define the configurations for a training job
type ManifestV1 struct {
	Name         string                         `yaml:"name,omitempty"`
	Description  string                         `yaml:"description,omitempty"`
	Version      string                         `yaml:"version,omitempty"`
	Cpus         float64                        `yaml:"cpus,omitempty"`
	Gpus         float64                        `yaml:"gpus,omitempty"`
	GpuType      string                         `yaml:"gpu_type,omitempty"`
	Learners     int32                          `yaml:"learners,omitempty"`
	Memory       string                         `yaml:"memory,omitempty"`
	Storage      string                         `yaml:"storage,omitempty"`
	DataStores   []*dataStoreRef                `yaml:"data_stores,omitempty"`
	Framework    *frameworkV1                   `yaml:"framework,omitempty"`
	Namespace    string                         `yaml:"namespace,omitempty"`
	JobAlert     map[string][]map[string]string `yaml:"job_alert,omitempty"`
	CodeSelector string                         `yaml:"code_selector,omitempty"`
	CodePath     string                         `yaml:"code_path,omitempty"`
	JobType      string                         `yaml:"job_type,omitempty"`
	ExpName      string                         `yaml:"exec_type,omitempty"`
	ExpRunId     int64                          `yaml:"exp_run_id,omitempty"`
	FileName     string                         `yaml:"fileName,omitempty"`
	PSs          string                         `yaml:"pss,omitempty"`
	PSCPU        string                         `yaml:"ps_cpu,omitempty"`
	PSImage      string                         `yaml:"ps_image,omitempty"`
	PSMemory     string                         `yaml:"ps_memory,omitempty"`
	TFosRequest  *v1alpha1.TFosRequest          `yaml:"tfos_request,omitempty"`
	DssTaskId    int64                          `yaml:"dss_task_id,omitempty"`
	ProxyUser    string                         `yaml:"proxy_user,omitempty"`
	JobParams    string                         `yaml:"job_params,omitempty"`
	DataSet      *DataSet                       `yaml:"data_set,omitempty"`
	MFModel      *MFModel                       `yaml:"mf_model,omitempty"`
	Algorithm    string                         `yaml:"algorithm,omitempty"`
	FitParams    string                         `yaml:"fit_params,omitempty"`
	APIType      string                         `yaml:"API_type,omitempty"`
	SubmitId     string                         `yaml:"submit_id,omitempty"`
	UserId       string                         `yaml:"user_id,omitempty"`
}

type DataSet struct {
	TrainingDataPath    string `yaml:"training_data_path,omitempty"`
	TrainingLabelPath   string `yaml:"training_label_path,omitempty"`
	TestingDataPath     string `yaml:"testing_data_path,omitempty"`
	TestingLabelPath    string `yaml:"testing_label_path,omitempty"`
	ValidationDataPath  string `yaml:"validation_data_path,omitempty"`
	ValidationLabelPath string `yaml:"validation_label_path,omitempty"`
}

type MFModel struct {
	GroupId     string `yaml:"group_id,omitempty"`
	ModelName   string `yaml:"model_name,omitempty"`
	Version     string `yaml:"version,omitempty"`
	FactoryName string `yaml:"factory_name,omitempty"`
}

type dataStoreRef struct {
	ID                string              `yaml:"id,omitempty"`
	Type              string              `yaml:"type,omitempty"`
	TrainingData      *storageContainerV1 `yaml:"training_data,omitempty"`
	TrainingResults   *storageContainerV1 `yaml:"training_results,omitempty"`
	TrainingWorkspace *storageContainerV1 `yaml:"training_workspace,omitempty"`
	Connection        map[string]string   `yaml:"connection,omitempty"`
}

type frameworkV1 struct {
	Name     string `yaml:"name,omitempty"`
	Version  string `yaml:"version,omitempty"`
	ImageTag string `yaml:"image_tag,omitempty"`
	Command  string `yaml:"command,omitempty"`
}

type storageContainerV1 struct {
	Container string `yaml:"container,omitempty"`
}

// LoadManifestV1 constructs a manifest object from a byte array
func LoadManifestV1(data []byte) (*ManifestV1, error) {
	t := &ManifestV1{}
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// WriteManifestV1 writes a manifest object to a byte array
func WriteManifestV1(t *ManifestV1) ([]byte, error) {
	data, err := yaml.Marshal(t)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// manifest2TrainingRequest converts the existing manifest to a training request for the trainer microservice.
func manifest2TrainingRequest(m *ManifestV1, modelDefinition []byte, http *http.Request,
	logr *logger.LocLoggingEntry, dataPath string) (*grpc_trainer_v2.CreateRequest, error) {

	// 1. Params Init
	authToken := http.Header.Get(cc.CcAuthToken)
	userId := getUserID(http)

	if m.UserId == "" {
		m.UserId = userId
	}

	execUserId := userId
	if m.Namespace == "" {
		return nil, *new(error)
	}
	logr.Debugf("debug for ccAddress restApi: %v", viper.GetString(config.CCAddress))
	ccClient := cc.GetCcClient(viper.GetString(config.CCAddress))
	if m.ProxyUser != "" {
		logr.Debugf("manifest2TrainingRequest to ProxyUser to: %v", m.ProxyUser)
		execUserId = m.ProxyUser
	}

	// 2.User Namespace Check
	err := ccClient.UserNamespaceCheck(authToken, userId, m.Namespace)
	if err != nil {
		log.Error("UserNamespaceCheck error:", err.Error())
		return nil, err
	}

	//3. Get GUID From CC
	var gid, uid *string
	if m.ProxyUser == "" {
		gid, uid, err = ccClient.GetGUIDFromUserId(authToken, execUserId)
		if err != nil {
			log.Error("GetGUIDFromUserId error:", err.Error())
			return nil, err
		}
	} else {
		proxyUser, err := ccClient.GetProxyUserFromCC(authToken, userId, execUserId)
		if err != nil {
			log.Error("GetProxyUserFromCC error:", err.Error())
			return nil, err
		}
		gidStr := strconv.Itoa(proxyUser.GID)
		uidStr := strconv.Itoa(proxyUser.UID)
		gid = &gidStr
		uid = &uidStr
	}

	//4. DataStore Init
	if nil != m.DataStores && m.ProxyUser == "" {
		path := m.DataStores[0].Connection["path"]
		userDir := strings.Split(path[1:], "/")[3]
		if userId != userDir && execUserId != userDir {
			groupId := viper.GetString(config.MLSSGroupId)
			gid = &groupId
		}
		logr.Debugf("manifest2TrainingRequest to get gid: %v", gid)
	}

	// 5. JobAlert Transform
	jobAlert := m.JobAlert
	jobString := ""
	if nil != jobAlert {
		marshal, marErr := json.Marshal(m.JobAlert)
		if marErr != nil {
			log.Error("JobAlert Marshal error:", marErr.Error())
			return nil, marErr
		}
		jobString = string(marshal)
	}

	//6. TFOS Transform(Option)
	var tFRequest string
	if nil != m.TFosRequest {
		marshal, marErr := json.Marshal(m.TFosRequest)
		if marErr != nil {
			log.Error("TFosRequest Marshal error:", marErr.Error())
			return nil, marErr
		}
		tFRequest = string(marshal)
	}

	//7. Build Trainer RPC Request
	r := &grpc_trainer_v2.CreateRequest{
		UserId: m.UserId,
		ModelDefinition: &grpc_trainer_v2.ModelDefinition{
			Name:        m.Name,
			Description: m.Description,
			Framework: &grpc_trainer_v2.Framework{
				Name:     m.Framework.Name,
				Version:  m.Framework.Version,
				ImageTag: m.Framework.ImageTag,
			},
			Content: modelDefinition,
		},
		Training: &grpc_trainer_v2.Training{
			Command:   m.Framework.Command,
			InputData: []string{m.DataStores[0].ID + "-input"},
			Profiling: false,
		},
		Datastores: []*grpc_trainer_v2.Datastore{
			{
				Id:         m.DataStores[0].ID + "-input",
				Type:       m.DataStores[0].Type,
				Connection: m.DataStores[0].Connection,
				Fields: map[string]string{
					"bucket": m.DataStores[0].TrainingData.Container,
				},
			},
		},
		Namespace:    m.Namespace,
		Gid:          *gid,
		Uid:          *uid,
		JobAlert:     jobString,
		CodeSelector: m.CodeSelector,
		DataPath:     dataPath,
		PSs:          m.PSs,
		PSCpu:        m.PSCPU,
		PSImage:      m.PSImage,
		PSMemory:     m.PSMemory,
		JobType:      m.JobType,
		TFosRequest:  tFRequest,
		ProxyUser:    m.ProxyUser,
		JobParams:    m.JobParams,
		Algorithm:    m.Algorithm,
		FitParams:    m.FitParams,
		APIType:      m.APIType,
		SubmitId:     userId,
	}

	// FIXME MLSS Change: get alert parms
	logr.Infof("manifest, JobAlert connection: %v, type: %v", m.DataStores[0].Connection, m.DataStores[0].Connection["type"])
	logr.Infof("manifest2TrainingRequest, JobType: %v", m.JobType)
	logr.Infof("manifest2TrainingRequest, TfosRequest: %v", m.TFosRequest)

	if m.DataStores[0].TrainingResults != nil {
		r.Training.OutputData = []string{m.DataStores[0].ID + "-output"}
		r.Datastores = append(r.Datastores, &grpc_trainer_v2.Datastore{
			Id:         m.DataStores[0].ID + "-output",
			Type:       m.DataStores[0].Type,
			Connection: m.DataStores[0].Connection,
			Fields: map[string]string{
				"bucket": m.DataStores[0].TrainingResults.Container,
			},
		})
	}
	// FIXME MLSS Change: v_1.5.1 get workData when TrainingWork is not null
	if m.DataStores[0].TrainingWorkspace != nil {
		r.Training.WorkData = []string{m.DataStores[0].ID + "-work"}
		r.Datastores = append(r.Datastores, &grpc_trainer_v2.Datastore{
			Id:         m.DataStores[0].ID + "-work",
			Type:       m.DataStores[0].Type,
			Connection: m.DataStores[0].Connection,
			Fields: map[string]string{
				"codeWork": m.DataStores[0].TrainingWorkspace.Container,
			},
		})
	}

	mem, memUnit, err := convertMemoryFromManifest(m.Memory)
	if err != nil {
		logr.WithError(err).Errorf("Incorrect memory specification in manifest")
	}

	storage, storageUnit := float32(0), grpc_trainer_v2.SizeUnit_MB // default to 0
	if m.Storage != "" {
		storage, storageUnit, err = convertMemoryFromManifest(m.Storage)
		if err != nil {
			logr.WithError(err).Errorf("Incorrect storage specification in manifest")
		}
	}
	logr.Debugf("convertMemoryFromManifest reqMem: %v", mem)
	r.Training.Resources = &grpc_trainer_v2.ResourceRequirements{
		Gpus:        m.Gpus,
		GpuType:     m.GpuType,
		Cpus:        m.Cpus,
		Memory:      mem,
		MemoryUnit:  memUnit,
		Storage:     storage,
		StorageUnit: storageUnit,
		Learners:    m.Learners,
		// TODO add storage support
	}

	if nil != m.DataSet {
		r.DataSet = &grpc_trainer_v2.DataSet{
			TrainingDataPath:    m.DataSet.TrainingDataPath,
			TrainingLabelPath:   m.DataSet.TrainingLabelPath,
			TestingDataPath:     m.DataSet.TestingDataPath,
			TestingLabelPath:    m.DataSet.TestingLabelPath,
			ValidationDataPath:  m.DataSet.ValidationDataPath,
			ValidationLabelPath: m.DataSet.ValidationLabelPath,
		}
	}

	if nil != m.MFModel {
		r.MfModel = &grpc_trainer_v2.MFModel{
			GroupId:     m.MFModel.GroupId,
			ModelName:   m.MFModel.ModelName,
			Version:     m.MFModel.Version,
			FactoryName: m.MFModel.FactoryName,
		}
	}

	return r, nil
}

// trainingToManifest converts the existing a training to manifest for the trainer microservice
func trainingToManifest(m *restmodels.Model, modelDefinition []byte, http *http.Request,
	logr *logger.LocLoggingEntry, dataPath string) (*ManifestV1, error) {

	var err error

	mem, err := convertMemoryToStringFromManifest(m.Training.Memory)
	if err != nil {
		logr.WithError(err).Errorf("Incorrect memory specification in manifest")
	}

	//7. Build Trainer RPC Request
	r := &ManifestV1{
		APIType:      m.APIType,
		Name:         m.Name,
		Description:  m.Description,
		PSs:          m.Pss,
		PSImage:      m.PsImage,
		JobType:      m.JobType,
		Namespace:    m.JobNamespace,
		Cpus:         m.Training.Cpus,
		Gpus:         m.Training.Gpus,
		PSCPU:        m.PsCPU,
		PSMemory:     m.PsMemory,
		Memory:       mem,
		ProxyUser:    m.ProxyUser,
		JobParams:    m.JobParams,
		CodeSelector: m.CodeSelector,
		UserId:       m.UserID,
		FileName:     m.FileName,
		CodePath:     m.FilePath,
	}

	// FIXME MLSS Change: get alert parms
	logr.Infof("manifest, JobAlert connection: %v, type: %v", m.DataStores[0].Connection, m.DataStores[0].Connection["type"])
	logr.Infof("manifest2TrainingRequest, JobType: %v", m.JobType)
	logr.Infof("manifest2TrainingRequest, TfosRequest: %v", m.TFosRequest)

	if nil != m.DataSet {
		r.DataSet = &DataSet{
			TrainingDataPath:    m.DataSet.TrainingDataPath,
			TrainingLabelPath:   m.DataSet.TrainingLabelPath,
			TestingDataPath:     m.DataSet.TestingDataPath,
			TestingLabelPath:    m.DataSet.TestingLabelPath,
			ValidationDataPath:  m.DataSet.ValidationDataPath,
			ValidationLabelPath: m.DataSet.ValidationLabelPath,
		}
	}

	r.MFModel = &MFModel{
		Version: m.Framework.Version,
	}

	r.Framework = &frameworkV1{
		Version: m.Framework.Version,
		Command: m.Training.Command,
		Name:    m.Framework.Name,
	}
	datas1 := &dataStoreRef{
		Connection:   m.DataStores[0].Connection,
		TrainingData: &storageContainerV1{Container: m.DataStores[0].Connection["bucket"]},
	}
	for _, v := range m.DataStores {

		if val, ok := v.Connection["bucket"]; ok {
			datas1.TrainingResults = &storageContainerV1{Container: val}
		}

		if val, ok := v.Connection["codeWork"]; ok {
			datas1.TrainingWorkspace = &storageContainerV1{Container: val}
		}
	}
	r.DataStores = append(r.DataStores, datas1)
	logr.Infof("DataStores:%v", r.DataStores)

	return r, err
}

// Converts an incoming memory spec like 4GB or 1024mb into a memory (float64)
// and a unit.
// Rules for conversion:
// 1GB (gigabytes) = 1000MB (megabytes)
// 1Gi (gibibytes) = 1024 Mi (mebibytes)
func convertMemoryFromManifest(memStr string) (float32, grpc_trainer_v2.SizeUnit, error) {
	regex, _ := regexp.Compile(`(?i)(?P<memory>(\d+)?)\s*(?P<unit>(mb|mib|gb|gib|tb|tib)?)`)
	match := regex.FindStringSubmatch(memStr)

	// convert to a map
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && i <= len(match) {
			result[name] = match[i]
		}
	}
	mem, err := strconv.ParseFloat(result["memory"], 64)
	if err != nil {
		return 512, grpc_trainer_v2.SizeUnit_MB, fmt.Errorf("Memory requirements not correctly specified: %s", memStr)
	}
	unit := strings.ToLower(result["unit"])

	for k, v := range grpc_trainer_v2.SizeUnit_value {
		if strings.ToLower(k) == unit {
			return float32(mem), grpc_trainer_v2.SizeUnit(v), nil
		}
	}
	return 512, grpc_trainer_v2.SizeUnit_MB, fmt.Errorf("Memory requirements not correctly specified: %s", memStr)
}

// convertMemoryToStringFromManifest + Gb
func convertMemoryToStringFromManifest(momory float64) (mom string, err error) {
	mom = strconv.FormatFloat(momory, 'f', 0, 64)
	mom = mom + "Gb"
	return
}
