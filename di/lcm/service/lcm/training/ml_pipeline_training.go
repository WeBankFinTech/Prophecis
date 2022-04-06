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
package lcm_training

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"time"

	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/lcm/service/lcm/job"
	"webank/DI/lcm/service/lcm/learner"
	"webank/DI/lcm/service/lcm/utils"

	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const loadDataContainerName = "load-data"
const loadModelContainerName = "load-model"
const learnerContainerName = "learner"
const storeResultsContainerName = "store-results"
const storeLogsContainerName = "store-logs"
const learnerConfigDir = "/etc/learner-config"

type mlPipelineTraining struct {
	*training
	tr *trainingResource
}

type commandParams struct {
	jobParams string
	dataset string
	model string
	resultDir  string
	factoryName string
	fitParams string
	APIType string
}

func (t mlPipelineTraining) Start() error {
	//Define Resource Spec
	err := t.transoformTrainingResourceSpec()
	if err != nil {
		return err
	}
	return t.createTrainingResourceInK8s()
}

//Create Job Kubernetes Resource
func (t *mlPipelineTraining) createTrainingResourceInK8s() error {
	logging := t.logging
	logging.Infof("Create Job in namespace: %s", t.tr.namespace)
	logging.Debugf("CreateJobResource in MLPipeline Training, req, TrainingId: %s, Framework: %s, JobType: %s, Learners: %s, Cpus: %s, Memory: %s, Gpus: %s, GpuType: %s, ", t.req.TrainingId, t.req.Framework, t.req.JobType, t.req.Resources.Learners, t.req.Resources.Cpus, t.req.Resources.Memory, t.req.Resources.Gpus, t.req.Resources.GpuType)

	//create the kubernetes resource
	return backoff.RetryNotify(func() error {
		_, err := t.k8sClient.CoreV1().Services(t.tr.namespace).Create(t.tr.service)
		logging.WithError(err).Warnf("Service Create err %s :", t.tr.service.Name)
		_, err = t.k8sClient.BatchV1().Jobs(t.tr.namespace).Create(t.tr.job)
		if k8serrors.IsAlreadyExists(err) {
			logging.WithError(err).Warnf("Job %s already exists", t.tr.job.Name)
			return nil
		}
		return err
	}, utils.K8sInteractionBackoff(), func(err error, window time.Duration) {
		logging.WithError(err).Errorf("Failed in Creating Job %s while deploying for training ", t.tr.job.Name)
		logging.Debugf("Create Job in K8s jobBOM, Volumes: %+v, VolumeMounts: %+v", t.tr.job.Spec.Template.Spec.Volumes, t.tr.job.Spec.Template.Spec.Containers[0].VolumeMounts)
	})

}

//Transform training to k8s Job and SVC  Spec
func (t *mlPipelineTraining) transoformTrainingResourceSpec() error {
	//SVC Definition
	serviceSpec := learner.CreateServiceSpec(t.learner.name, t.req.TrainingId, t.k8sClient)
	// if err != nil {
	// 	t.logr.WithError(err).Errorf("Could not create svc spec for %s", serviceSpec.Name)
	// 	return err
	// }

	//Job Definition
	jobSpec, err := t.createJobSpec(serviceSpec)
	if err != nil {
		t.logging.WithError(err).Errorf("Could not create job spec for %s", serviceSpec.Name)
		return err
	}

	//Set Training Resource
	t.tr = &trainingResource{
		serviceSpec,
		jobSpec,
		t.req.JobNamespace,
	}
	return nil
}

func (t *mlPipelineTraining) createJobSpec(service *v1core.Service) (*v1.Job, error) {
	t.logging.Debugf("debug_for_param ParameterServer lcm newTraining pss: %v, psCpu: %v, psImage: %v, psMemory: %v", t.req.PSs, t.req.PSCPU, t.req.PSImage, t.req.PSMemory)

	//Resource Definiton
	//1. Namespace & Basic Info
	nsObj, err := t.k8sClient.CoreV1().Namespaces().Get(t.req.JobNamespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	//2.Meta Setting
	///NodeSelector
	nodeSelectors := ""
	namespaceAnnotations := (nsObj.ObjectMeta).GetAnnotations()
	selector, ok := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]
	if ok {
		nodeSelectors = selector
	}
	nodeSelectorMap := make(map[string]string)
	if len(nodeSelectors) > 0 {
		nodeSelectorList := strings.Split(nodeSelectors, ",")
		for i := 0; i < len(nodeSelectorList); i++ {
			nodeSelector := nodeSelectorList[i]
			nodeSelectorKv := strings.Split(nodeSelector, "=")
			nodeSelectorMap[nodeSelectorKv[0]] = nodeSelectorKv[1]
		}
	}

	//3.Env
	envVarsFromDeploymentRequest := job.GenerateJobEnv(t.req,service) //shared across all containers of training
	mountTrainingDataStoreInLearner := true
	mountResultsStoreInLearner := true
	t.logging.Debugf("envVarsFromDeploymentRequest: %+v", envVarsFromDeploymentRequest)
	envvarsForJob := job.EnvVarsForDeployingJob(envVarsFromDeploymentRequest, t.req.TrainingId,
		1, t.trainingID, mountTrainingDataStoreInLearner, mountResultsStoreInLearner) //only for learner


	//4.Volume
	//Volumn Init
	t.logging.Debugf("NewLCMTraining envvarsForLearner: %+v", envvarsForJob)
	jobVolumes := job.VolumesForJob(t.req, envvarsForJob, t.logging)
	jobVolumeSpecs := jobVolumes.CreateVolumeForLearner()


	//5.Secret
	//Job Definition
	jobDefinition := trainingDefinition{
		// secrets:                         secretsForDeployingLearner(req, mountTrainingDataStoreInLearner, mountResultsStoreInLearner),
		volumes:                         jobVolumeSpecs,
		volumeMounts:                    jobVolumes.CreateVolumeMountsForLearner(),
		envVars:                         envvarsForJob,
		name:                            t.trainingID,
	}

	//Container Definition
	imagePullSecret, err := learner.GenerateImagePullSecret(t.k8sClient, t.req)
	if err != nil {
		return nil, err
	}


	//Container Build
	learnerDefn := jobDefinition
	learnerContainer := constructContainer(t.req, jobDefinition.envVars, jobDefinition.volumeMounts,  t.logging) // nil for mounting shared NFS volume since non split mode
	// FIXME MLSS Change: add nodeSelectors to deployment/sts pods
	splitLearnerPodSpec := learner.CreatePodSpecForJob([]v1core.Container{learnerContainer}, learnerDefn.volumes, map[string]string{"training_id": t.req.TrainingId, "user_id": t.req.UserId}, nodeSelectorMap, imagePullSecret)
	jobSpec := learner.CreateJobSpecForLearner(t.req.TrainingId, splitLearnerPodSpec)

	return jobSpec, nil
}

//construct learner container
func constructContainer(req *service.JobDeploymentRequest, envVars []v1core.EnvVar, learnerVolumeMounts []v1core.VolumeMount, logging *logger.LocLoggingEntry) v1core.Container {
	//Resource Init
	cpuCount := v1resource.NewMilliQuantity(int64(float64(req.Resources.Cpus)*1000.0), v1resource.DecimalSI)
	gpuCount := v1resource.NewQuantity(int64(req.Resources.Gpus), v1resource.DecimalSI)
	memInBytes := int64(utils.CalcMemory(req.Resources) * 1024 * 1024)
	memCount := v1resource.NewQuantity(memInBytes, v1resource.DecimalSI)

	//Param Definition
	//resultDir := ""
	//jobParams := req.JobParams
	//dataSet := req.Dataset
	//storageModel := req.StorageModel


	learnerBashCommand := `
            mkdir -p $JOB_RESULT_DIR/learner-$TRAINING_ID
            chown -R  $UID:$UID $JOB_RESULT_DIR
			# Source BDP ENV & Host
			if [ -d /appcom ]
			then
				sh /appcom/config/MLSS-config/MLSS_AIDE-config/HOST_ENV.sh
				source /appcom/config/MLSS-config/MLSS_AIDE-config/Notebook_ENV.sh
				if [ ! -d /workspace/tmp ]
				then
				   mkdir -p /workspace/tmp
				   chown $JOB_USER:users /workspace/tmp
				fi
				if [ ! -d /workspace/logs ]
				then
				   mkdir -p /workspace/logs
				   chown $JOB_USER:users /workspace/logs
				fi
				ln -s /workspace/tmp /appcom
				ln -s /workspace/logs /appcom
			fi ;`

	commandParams,err := buildCommandParamsStr(req)
	if err != nil{

	}
	fitParamsStr := ""
	if commandParams.fitParams != "'{}'"{
		fitParamsStr = ` --fit_params ` + commandParams.fitParams
	}

	if commandParams.APIType != "" {
		fitParamsStr = fitParamsStr + " --API_type " + commandParams.APIType
	}

	pythonCommand := `/opt/conda/bin/python /data/python/script/` + getScriptName(req.Algorithm) + ` ` +
		            ` --job_params `+ commandParams.jobParams +
		            ` --dataset `+  commandParams.dataset +
				    ` --model `+  commandParams.model +
					` --factory_name `+  commandParams.factoryName + fitParamsStr +
					` --result_dir `+  commandParams.resultDir + ` 2>&1 $@  | tee -a $JOB_RESULT_DIR/learner-$TRAINING_ID/training-log.txt /logs/$TRAINING_ID.log ; ` +
					`cmd_exit=${PIPESTATUS[0]};` +
				    `echo "python script status: ${cmd_exit}"`
					//`  >> /tmp/latest-log 2>&1 ;`
	learnerBashCommand += `
            echo "start command"
			export PYTHONPATH=$PWD ;
			echo "$(date): Starting training job";` + pythonCommand + ";" +
			`chown $UID:$UID $JOB_RESULT_DIR/learner-$TRAINING_ID/training-log.txt;
			echo ${cmd_exit};
			exit ${cmd_exit};
			echo "$(date): Training exit with exit code ${cmd_exit}." >> /tmp/latest-log`

			//eval "/opt/conda/bin/python /data/python/script/lr.py 2>&1" >> /tmp/latest-log 2>&1 ;




	//FIXME need to have the learner IDs start from 1 rather than 0
	var cmd string
	//FIXME we assume all images has python installed to unzip the model.zip
	var loadModelCommand string
	learnerCommand := learnerBashCommand

	// FIXME MLSS Change: just move the logs, and make sure the owner of the files is the UID user
	var storeLogsCommand string
	//sleep 40000
	storeLogsCommand = `
		echo ${cmd_exit}
		echo Calling copy logs.
		export LEARNER_ID=${HOSTNAME} ;
		chown -R $UID:$GID $JOB_RESULT_DIR/learner-$LEARNER_ID/
		bash -c 'exit $? '`

	cmd = job.WrapCommands([]job.ContainerCommands{
		//{cmd: getCreateTrainCMD(), container: ""},
		{Cmd: loadModelCommand, Container: loadModelContainerName},
		{Cmd: learnerCommand, Container: learnerContainerName},
		{Cmd: storeLogsCommand, Container: storeLogsContainerName}})

	container := learner.Container{
		Image: learner.Image{Framework: req.Framework, Version: req.Version, Tag: req.EnvVars["DLAAS_LEARNER_IMAGE_TAG"]},
		Resources: learner.Resources{
			CPUs: *cpuCount, Memory: *memCount, GPUs: *gpuCount,
		},
		VolumeMounts: learnerVolumeMounts,
		Name:         learnerContainerName,
		EnvVars:      envVars,
		Command:      cmd,
	}

	logging.Debugf("learner.CreateContainerSpec(container), container: %+v", container.EnvVars)
	learnerContainer := learner.CreateContainerSpec(container)

	job.ExtendLearnerContainer(&learnerContainer, req, logging)

	return learnerContainer
}

// NewMLPipelineTraining: New MLPipeline job type training, init training
func NewMLPipelineTraining(ctx context.Context, k8sClient kubernetes.Interface, req *service.JobDeploymentRequest, log *logger.LocLoggingEntry) Training {
	//Basic Info Init
	jobName := fmt.Sprintf("learner-%s", req.Name)
	logr := log.WithFields(logrus.Fields{
		"learner_name": jobName,
	})
	jobDefinition := trainingDefinition{}
	return mlPipelineTraining{&training{ctx, k8sClient, req, req.TrainingId, jobDefinition, logr}, nil}
}


func buildCommandParamsStr(req *service.JobDeploymentRequest) (*commandParams,error) {
	params := commandParams{}

	//Job Params Setting
	params.jobParams = req.JobParams
	if params.jobParams == "" {
		params.jobParams = "'{}'"
	}else{
		params.jobParams = "'" + req.JobParams + "'"
	}


	//Job Fit Params Setting
	params.fitParams = req.FitParams
	if params.fitParams == "" {
		params.fitParams = "'{}'"
	}else{
		params.fitParams = "'" + req.FitParams + "'"
	}

	//DataSet Setting
	dataSetStr, err := json.Marshal(req.DataSet)
	if  err != nil {
		return nil,err
	}
	params.dataset = "'" + string(dataSetStr) + "'"

	//resultDir
	params.resultDir = "/workspace/result/"+req.TrainingId

	//model
	mfModelStr, err := json.Marshal(req.MfModel)
	if  err != nil {
		return nil,err
	}
	params.model = "'" + string(mfModelStr) + "'"


	params.factoryName = req.MfModel.FactoryName
	if params.factoryName == ""{
		params.factoryName = "'None'"
	}

	params.APIType = req.APIType

	return &params, nil
}


func getScriptName(algorithm string) string{
	//TODO: decision tree
	if algorithm == "DecisionTree" {
		return "decision_tree.py"
	}else if algorithm == "LogisticRegression"{
		return "lr.py"
	}else if algorithm == "XGBoost"{
		return "xgb.py"
	}else if algorithm == "LightGBM"{
		return "lgb.py"
	}else if algorithm == "RandomForest"{
		return "random_forest.py"
	}else{
		return ""
	}
}

