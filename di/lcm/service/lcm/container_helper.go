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

package lcm

import (
	"fmt"
	"strings"
	"text/template"
	"webank/DI/commons/constants"

	"webank/DI/commons/config"
	"webank/DI/commons/service"
	"webank/DI/lcm/service/lcm/learner"

	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"

	"bytes"

	"webank/DI/commons/logger"

	v1core "k8s.io/api/core/v1"
	v1resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const logCollectorContainerName string = "log-collector" // the name of the learner container in the pod
const loadDataContainerName = "load-data"
const loadModelContainerName = "load-model"
const learnerContainerName = "learner"
const storeResultsContainerName = "store-results"
const storeLogsContainerName = "store-logs"
const learnerConfigDir = "/etc/learner-config"

const simpleLogCollectorName = "log_collector"

const codeFile = "codeFile"
const storagePath = "storagePath"

const logCollectorBadTagNoTDSFound = "dummy-tag-no-tds-found"

const FLUENTBIT = "fluent-bit"

// valid names of databroker types that map to "databroker_<type>" Docker image names
var validDatabrokerTypes = []string{"objectstorage", "s3"}

// default databroker type
var defaultDatabrokerType = "objectstorage"

const (
	workerPort int32 = 2222
	sshPort    int32 = 22
)

//helpers won't know about learner ID since that is only available to stateful set learners
// TODO used by controller (for tracking state of individual learners) and logs, do we really need this for logs
// need to use 1 and not 0 because job monitor tracks path starting with learner 1 and not 0
const masterLearnerID = 1

func fetchImageNameFromEvaluationMetrics(evalMetricsString string,
	learnerTag string,
	framework string,
	version string,
	logr *logger.LocLoggingEntry) (string, string) {

	logr.Debugf("evaluation_metrics: %v<end>", evalMetricsString)
	logCollectorImageShortName := simpleLogCollectorName

	learnerEMTag := learnerTag

	logr.Debugf("evalMetricsString: %s", evalMetricsString)
	if evalMetricsString != "" {
		em := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(evalMetricsString), &em)
		if err != nil {
			// Assuming pre-validation, this is unlikely to happen, so this is mostly a programmer assertion.
			logr.WithError(err).Error("evaluation_metrics was specified in manifest, but can't be parsed!")
		}

		m := em["evaluation_metrics"].(map[interface{}]interface{})

		if m != nil {
			val, ok := m["image_tag"]
			logr.Debugf("learner tag: %s %t", val, ok)
			if ok == false {
				// TODO: fix dropping underbar problem.  Somehow.
				// Having a hard time with, I think the yaml to string stuff, dropping underbars.
				val, ok = m["imagetag"]
			}
			if ok && val.(string) != "" {
				learnerEMTag = val.(string)
			}

			imageType, ok := m["type"]

			// Allow some synonyms for simple file extractor
			if ok && (imageType == "optivist" || imageType == "emetrics_file" || imageType == "file") {
				imageType = "emetrics_file_extractor"
			}

			if ok && imageType.(string) != "" {
				logr.Debugf("initial evaluation_metrics type: %s", imageType)
				// Assume the image name has been validated upstream
				logCollectorImageShortName = imageType.(string)
				if logCollectorImageShortName == "tensorboard" || logCollectorImageShortName == "tensorboard_extractor" {
					// For the moment we're just going to use TF 1.3, but the tag should change to be non-version
					// specific, and we should just use latest TF.
					logCollectorImageShortName = fmt.Sprintf("%s_extract", "tensorboard")
				}
				// Be flexible
				if logCollectorImageShortName == "null" || logCollectorImageShortName == "nil" ||
					logCollectorImageShortName == "logger" || logCollectorImageShortName == "none" {
					logCollectorImageShortName = simpleLogCollectorName
				}

			} else {
				logr.Error("evaluation_metrics type is empty")
				logCollectorImageShortName = simpleLogCollectorName
			}
		} else {
			logr.Debug("No evaluation metrics specified! (2)")
		}
	} else {
		logr.Debug("No evaluation metrics specified! (1)")
	}
	return logCollectorImageShortName, learnerEMTag
}

func findTrainingDataServiceTag(k8sClient kubernetes.Interface, logr *logger.LocLoggingEntry) string {
	selector := "service==ffdl-trainingdata"
	podInterface := k8sClient.CoreV1().Pods(config.GetPodNamespace())
	pods, err := podInterface.List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		logr.WithError(err).Debugf("Could not find service=ffdl-trainingdata")
		// bad fallback, ideally should never happen
		return logCollectorBadTagNoTDSFound
	}
	nPods := len(pods.Items)
	if nPods > 0 {
		for i := nPods - 1; i >= 0; i-- {
			containerStatuses := pods.Items[i].Status.ContainerStatuses
			for _, containerStatus := range containerStatuses {
				imageName := containerStatus.Image
				// No tag, don't build a log-collector
				splits := strings.SplitAfter(imageName, ":")
				if splits != nil && len(splits) > 1 {
					return splits[len(splits)-1]
				}
			}
		}
	}
	// bad fallback, ideally should never happen
	return logCollectorBadTagNoTDSFound
}

func constructLearnerContainer(req *service.JobDeploymentRequest, envVars []v1core.EnvVar, learnerVolumeMounts []v1core.VolumeMount, mountTrainingDataStoreInLearner, mountResultsStoreInLearner bool, logr *logger.LocLoggingEntry) v1core.Container {

	cpuCount := v1resource.NewMilliQuantity(int64(float64(req.Resources.Cpus)*1000.0), v1resource.DecimalSI)
	gpuCount := v1resource.NewQuantity(int64(req.Resources.Gpus), v1resource.DecimalSI)
	memInBytes := int64(calcMemory(req.Resources) * 1024 * 1024)
	memCount := v1resource.NewQuantity(memInBytes, v1resource.DecimalSI)

	// FIXME MLSS Change: read platform namespace and download model from s3 svc in that namespace
	platformNamespace := viper.GetString(config.PlatformNamespace)

	//argh!!! this should be abstracted out as well
	command := "for i in ${!ALERTMANAGER*} ${!DLAAS*} ${!ETCD*} ${!GRAFANA*} ${!HOSTNAME*} ${!KUBERNETES*} ${!MONGO*} ${!PUSHGATEWAY*}; do unset $i; done;"
	//learnerBashCommand := `bash -c 'train.sh >> $JOB_STATE_DIR/latest-log 2>&1 ; exit ${PIPESTATUS[0]}'`

	learnerBashCommand := `
			# Source BDP ENV & Host
			if [ -d /appcom ]
			then
				sh /appcom/config/MLSS-config/MLSS_DI-config/HOST_ENV.sh
				source /appcom/config/MLSS-config/MLSS_DI-config/DI_ENV.sh
				if [ ! -d /workspace/tmp ]
				then
				   mkdir -p /workspace/tmp
				   chown $NB_USER:users /workspace/tmp
				fi
				if [ ! -d /workspace/logs ]
				then
				   mkdir -p /workspace/logs
				   chown $NB_USER:users /workspace/logs
				fi
                mkdir -p /workspace/logs/hadoop
                chmod -R 777 /workspace/logs
				chmod -R 777 /workspace/logs/hadoop
				ln -s /workspace/tmp /appcom
				ln -s /workspace/logs /appcom
			fi ;`

	if req.JobType == "dist-tf" {
		learnerBashCommand += `
			if [[ $GID != 0 || $UID != 0 ]];then
				echo "running as $USER_ID with GID $GID and UID $UID";
				if [[ $LEARNER_ID =~ "worker-0" ]];then
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt and $JOB_STATE_DIR/latest-log";
					su $USER_ID -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt $JOB_STATE_DIR/latest-log ; exit ${PIPESTATUS[0]}';
				else
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt";
					su $USER_ID -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt ; exit ${PIPESTATUS[0]}';
				fi
			else
				echo "running as root";
				if [[ $LEARNER_ID =~ "worker-0" ]];then
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt and $JOB_STATE_DIR/latest-log";
					bash -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt $JOB_STATE_DIR/latest-log ; exit ${PIPESTATUS[0]}';
				else
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt";
					bash -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt ; exit ${PIPESTATUS[0]}';
				fi
			fi`
	} else {
		learnerBashCommand += `
			# Source BDP ENV & Host
			if [ -d /appcom ]
			then
				sh /appcom/config/MLSS-config/MLSS_AIDE-config/HOST_ENV.sh
				source /appcom/config/MLSS-config/MLSS_AIDE-config/Notebook_ENV.sh
				if [ ! -d /workspace/tmp ]
				then
				   mkdir -p /workspace/tmp
				   chown $NB_USER:users /workspace/tmp
				fi
				if [ ! -d /workspace/logs ]
				then
				   mkdir -p /workspace/logs
				   chown $NB_USER:users /workspace/logs
				fi
				ln -s /workspace/tmp /appcom
				ln -s /workspace/logs /appcom
			fi ;
			if [[ $GID != 0 || $UID != 0 ]];then
				echo "running as $USER_ID with GID $GID and UID $UID";
				su $USER_ID -c 'train.sh 2>&1 | tee $JOB_STATE_DIR/latest-log ; exit ${PIPESTATUS[0]}';
			else
				echo "running as root";
				bash -c 'train.sh 2>&1 | tee $JOB_STATE_DIR/latest-log ; exit ${PIPESTATUS[0]}';
			fi`
	}

	image := learner.Image{
		Framework: req.Framework,
		Version:   req.Version,
		Tag:       req.ImageTag,
	}

	// special settings for custom image
	if req.ImageLocation != nil {
		image.Registry = req.ImageLocation.Registry
		image.Namespace = req.ImageLocation.Namespace
		learnerBashCommand = `
            echo "start command"
			cd "$MODEL_DIR" ;
			export PYTHONPATH=$PWD ;
			echo "$(date): Starting training job" > $JOB_STATE_DIR/latest-log ;
			eval "$TRAINING_COMMAND 2>&1" >> $JOB_STATE_DIR/latest-log 2>&1 ;
			cmd_exit=$? ;
			echo "$(date): Training exit with exit code ${cmd_exit}." >> $JOB_STATE_DIR/latest-log`
	}
	//FIXME need to have the learner IDs start from 1 rather than 0
	var cmd string
	var doCondExitWrite = true
	//FIXME we assume all images has python installed to unzip the model.zip
	if mountResultsStoreInLearner {
		codeSelector := req.CodeSelector
		var loadModelComand string
		if codeSelector == codeFile {
			loadModelComand = `
				echo "Starting Training $TRAINING_ID and $DLAAS_MLSSGID"
				echo "Checking $DATA_DIR and $RESULT_DIR ownership"
				DATA_DIR_UID=$(ls -nd $DATA_DIR | awk '{print $3}')
				DATA_DIR_GID=$(ls -nd $DATA_DIR | awk '{print $4}')
				RESULT_DIR_UID=$(ls -nd $RESULT_DIR | awk '{print $3}')
				RESULT_DIR_GID=$(ls -nd $RESULT_DIR | awk '{print $4}')
				echo "DATA_DIR belongs to $DATA_DIR_UID:$DATA_DIR_GID"
				echo "RESULT_DIR belongs to $RESULT_DIR_UID:$RESULT_DIR_GID"
				if [[ $GID == 0 && $UID == 0 ]];then
					echo "current user is root, ownership check passed"
			    elif [[ $GID == $DLAAS_MLSSGID && $DATA_DIR_GID == $DLAAS_MLSSGID && $RESULT_DIR_GID == $DLAAS_MLSSGID ]];then 
					echo "current user GID matches DATA_DIR and RESULT_DIR ownership, ownership check passed"
				elif [[ $UID == $DATA_DIR_UID && $GID == $DATA_DIR_GID && $UID == $RESULT_DIR_UID && $GID == $RESULT_DIR_GID ]];then
					echo "current user GID and UID matches DATA_DIR and RESULT_DIR ownership, ownership check passed"
				else
					echo "current user GID and UID does not match DATA_DIR and RESULT_DIR ownership, ownership check failed"
					exit 1
				fi;
				mkdir -p $JOB_RESULT_DIR/_submitted_code ;
				export LEARNER_ID=${HOSTNAME} ;
				mkdir -p $JOB_RESULT_DIR/learner-$LEARNER_ID ;
				mkdir -p $CHECKPOINT_DIR ;
				ls -l $JOB_RESULT_DIR;
				chown -R $UID:$GID $JOB_RESULT_DIR ;
				chown -R $UID:$GID $CHECKPOINT_DIR ;
				if [[ $GID != 0 || $UID != 0 ]];then
					echo "loading model as $USER_ID with GID $GID and UID $UID"
					groupadd -g $GID $USER_ID
					useradd -g $GID -u $UID $USER_ID
					su $USER_ID
				else
					echo "loading model as root"
				fi;
				DEFAULT_TRAINING_DATA_BUCKET=tf_trained_model;
				DEFAULT_TRAINING_MODEL_ZIP_BUCKET=dlaas-models;
				MODEL_URL_S3=http://s3.` + platformNamespace + `.svc.cluster.local/$DEFAULT_TRAINING_DATA_BUCKET/$TRAINING_ID/_submitted_code/model.zip;
				MODEL_URL_HOSTMOUNT=http://s3.` + platformNamespace + `.svc.cluster.local/$DEFAULT_TRAINING_MODEL_ZIP_BUCKET/$TRAINING_ID.zip;
				MODEL_URL=$MODEL_URL_S3;
				S3_CODE="curl -XGET -Li $MODEL_URL_S3 -s -w \"%{http_code}\" -o /dev/null -H \"AWS_ACCESS_KEY_ID: test\" -H \"AWS_DEFAULT_REGION: us-east-1\" -H \"AWS_SECRET_ACCESS_KEY: test\" -H \"Authorization: Basic Og==\"" ;
				if [ "$S3_CODE" != 200 ];then
					echo "Cannot get model zip from default training data bucket, assuming using hostmount."
					MODEL_URL=$MODEL_URL_HOSTMOUNT
				fi;
				curl -XGET $MODEL_URL \
					-H "AWS_ACCESS_KEY_ID: test" -H "AWS_DEFAULT_REGION: us-east-1" -H "AWS_SECRET_ACCESS_KEY: test" -H "Authorization: Basic Og==" > $JOB_RESULT_DIR/_submitted_code/model.zip ;
				python -m zipfile -e $JOB_RESULT_DIR/_submitted_code/model.zip $MODEL_DIR`
		} else if codeSelector == storagePath {
			loadModelComand = `
				echo "$CODE_SELECTOR";
				echo "Starting Training $TRAINING_ID and $DLAAS_MLSSGID";
				echo "Checking $DATA_DIR and $RESULT_DIR ownership";
				DATA_DIR_UID=$(ls -nd $DATA_DIR | awk '{print $3}');
				DATA_DIR_GID=$(ls -nd $DATA_DIR | awk '{print $4}');
				RESULT_DIR_UID=$(ls -nd $RESULT_DIR | awk '{print $3}');
				RESULT_DIR_GID=$(ls -nd $RESULT_DIR | awk '{print $4}');
				echo "DATA_DIR belongs to $DATA_DIR_UID:$DATA_DIR_GID";
				echo "RESULT_DIR belongs to $RESULT_DIR_UID:$RESULT_DIR_GID";
				if [[ $GID == 0 && $UID == 0 ]];then
					echo "current user is root, ownership check passed";
                elif [[ $GID == $DLAAS_MLSSGID && $DATA_DIR_GID == $DLAAS_MLSSGID && $RESULT_DIR_GID == $DLAAS_MLSSGID ]];then 
					echo "current user GID matches DATA_DIR and RESULT_DIR ownership, ownership check passed"
				elif [[ $UID == $DATA_DIR_UID && $GID == $DATA_DIR_GID && $UID == $RESULT_DIR_UID && $GID == $RESULT_DIR_GID ]];then
					echo "current user GID and UID matches DATA_DIR and RESULT_DIR ownership, ownership check passed";
				else
					echo "current user GID and UID does not match DATA_DIR and RESULT_DIR ownership, ownership check failed";
					exit 1;
				fi;
				mkdir -p $JOB_STATE_DIR/logs ;
				mkdir -p $JOB_RESULT_DIR/_submitted_code ;
				export LEARNER_ID=${HOSTNAME} ;
				mkdir -p $JOB_RESULT_DIR/learner-$LEARNER_ID ;
				mkdir -p $CHECKPOINT_DIR ;
				ls -l $JOB_RESULT_DIR;
				chown -R $UID:$GID $JOB_RESULT_DIR ;
				chown -R $UID:$GID $CHECKPOINT_DIR ;
				if [[ $GID != 0 || $UID != 0 ]];then
					echo "loading model as $USER_ID with GID $GID and UID $UID";
					groupadd -g $GID $USER_ID;
					useradd -g $GID -u $UID $USER_ID;
					su $USER_ID;
				else
					echo "loading model as root";
				fi`
		}

		learnerCommand := learnerBashCommand
		//var learnerCommand string
		//if req.JobType == "" {
		//	learnerCommand = learnerBashCommand
		//} else {
		//	learnerCommand = getLearnerBashCommandForTFJob()
		//}

		// FIXME MLSS Change: just move the logs, and make sure the owner of the files is the UID user
		var storeLogsCommand string
		if req.JobType == "dist-tf" {
			storeLogsCommand = `
			echo Calling copy logs.
			export LEARNER_ID=${HOSTNAME} ;
			chown -R $UID:$GID $JOB_RESULT_DIR/learner-$LEARNER_ID/
			bash -c 'exit 0'`
		} else {
			storeLogsCommand = `
			echo Calling copy logs.
			export LEARNER_ID=${HOSTNAME} ;
			mv -nf $LOG_DIR/* $JOB_RESULT_DIR/learner-$LEARNER_ID/ ;
			ERROR_CODE=$? ;
			echo $ERROR_CODE > $JOB_RESULT_DIR/learner-$LEARNER_ID/.log-copy-complete ;	
			chown -R $UID:$GID $JOB_RESULT_DIR/learner-$LEARNER_ID/
			bash -c 'exit $ERROR_CODE'`
		}
		logr.Debugf("wrapCommands", loadModelComand, learnerCommand)

		if req.JobType != "dist-tf" {
			cmd = wrapCommands([]containerCommands{
				//{cmd: getCreateTrainCMD(), container: ""},
				{cmd: loadModelComand, container: loadModelContainerName},
				{cmd: learnerCommand, container: learnerContainerName},
				{cmd: storeLogsCommand, container: storeLogsContainerName},
			})
		} else {
			cmd = wrapCommandsForTFJob([]containerCommands{
				{cmd: getCreateTrainCMD(), container: ""},
				{cmd: loadModelComand, container: loadModelContainerName},
				{cmd: learnerCommand, container: learnerContainerName},
				{cmd: storeLogsCommand, container: storeLogsContainerName},
			})
		}
	} else {
		command = fmt.Sprintf(`%s mkdir -p $JOB_RESULT_DIR ; bash -c ' train.sh 2>&1 | tee -a %s/latest-log; exit ${PIPESTATUS[0]}'`, command)
		doCondExitWrite = false
		cmd = wrapCommand(command, learnerContainerName, doCondExitWrite)
	}

	logr.Debugf("debug_learnerVolumeMounts: %v", learnerVolumeMounts)
	// FIXME MLSS Change: add timezone VolumeMount
	learnerVolumeMounts = append(learnerVolumeMounts, constants.TimezoneMount)

	//
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

	learnerContainer := learner.CreateContainerSpec(container)
	extendLearnerContainer(&learnerContainer, req, logr)
	//logr.Debugf("learnerContainer: %v+", learnerContainer)
	// dumpAsYaml("learnerContainer", learnerContainer, logr)

	return learnerContainer
}

// FIXME mlss 1.9.0 new learner
func newConstructLearnerContainer(req *service.JobDeploymentRequest, envVars []v1core.EnvVar, learnerVolumeMounts []v1core.VolumeMount, mountTrainingDataStoreInLearner, mountResultsStoreInLearner bool, logr *logger.LocLoggingEntry) v1core.Container {
	//getLogger := logger.GetLogger()
	cpuCount := v1resource.NewMilliQuantity(int64(float64(req.Resources.Cpus)*1000.0), v1resource.DecimalSI)
	gpuCount := v1resource.NewQuantity(int64(req.Resources.Gpus), v1resource.DecimalSI)
	memInBytes := int64(calcMemory(req.Resources) * 1024 * 1024)
	memCount := v1resource.NewQuantity(memInBytes, v1resource.DecimalSI)

	// FIXME MLSS Change: read platform namespace and download model from s3 svc in that namespace
	platformNamespace := viper.GetString(config.PlatformNamespace)

	////get env from platformNamespace
	//split := strings.Split(platformNamespace, "-")
	//if len(split) != 2 {
	//	getLogger.Errorf("get env failed, len(split) != 2")
	//	return v1core.Container{}
	//}
	////ENV
	//env := split[1]

	//argh!!! this should be abstracted out as well
	command := "for i in ${!ALERTMANAGER*} ${!DLAAS*} ${!ETCD*} ${!GRAFANA*} ${!HOSTNAME*} ${!KUBERNETES*} ${!MONGO*} ${!PUSHGATEWAY*}; do unset $i; done;"
	//learnerBashCommand := `bash -c 'train.sh >> $JOB_STATE_DIR/latest-log 2>&1 ; exit ${PIPESTATUS[0]}'`

	learnerBashCommand := `
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

	learnerBashCommand += ` chmod -R 777 "$JOB_STATE_DIR";
							chmod -R 777 /job-logs;
							ls -al $JOB_STATE_DIR;`

	if req.JobType == "dist-tf" {
		learnerBashCommand += `
			if [[ $GID != 0 || $UID != 0 ]];then
				echo "running as $USER_ID with GID $GID and UID $UID";
				if [[ $LEARNER_ID =~ "worker-0" ]];then
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt and $JOB_STATE_DIR/latest-log";
					chmod -R 777 $JOB_STATE_DIR;
					ls -al $JOB_STATE_DIR;
					su $USER_ID -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt $JOB_STATE_DIR/latest-log /job-logs/${TRAINING_ID}.log ; exit ${PIPESTATUS[0]}';
				else
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt";
					su $USER_ID -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt /job-logs/${TRAINING_ID}.log ; exit ${PIPESTATUS[0]}';
				fi
			else
				echo "running as root";
				if [[ $LEARNER_ID =~ "worker-0" ]];then
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt and $JOB_STATE_DIR/latest-log";
					bash -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt $JOB_STATE_DIR/latest-log /job-logs/${TRAINING_ID}.log ; exit ${PIPESTATUS[0]}';
				else
                    echo "output: $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt";
					bash -c 'train.sh 2>&1 | tee $JOB_RESULT_DIR/learner-$LEARNER_ID/training-log.txt /job-logs/${TRAINING_ID}.log ; exit ${PIPESTATUS[0]}';
				fi
			fi
			TRAINING_RESULT=$? ;
			echo "TRAINING_RESULT: $TRAINING_RESULT" ;`
	} else {
		learnerBashCommand += `
			if [[ $GID != 0 || $UID != 0 ]];then
				echo "running as $USER_ID with GID $GID and UID $UID";
				su $USER_ID -c 'train.sh 2>&1 | tee $JOB_STATE_DIR/latest-log /job-logs/${TRAINING_ID}.log; exit ${PIPESTATUS[0]}';
			else
				echo "running as root";
				bash -c 'train.sh 2>&1 | tee $JOB_STATE_DIR/latest-log /job-logs/${TRAINING_ID}.log; exit ${PIPESTATUS[0]}';
			fi
			TRAINING_RESULT=$? ;
			echo "TRAINING_RESULT: $TRAINING_RESULT" ;`
	}

	//FIXME need to have the learner IDs start from 1 rather than 0
	var cmd string
	var doCondExitWrite = true
	//FIXME we assume all images has python installed to unzip the model.zip
	if mountResultsStoreInLearner {
		codeSelector := req.CodeSelector
		var loadModelComand string

		loadModelComand = `
				echo "CODE_SELECTOR: $CODE_SELECTOR";
				echo "Starting Training $TRAINING_ID and $DLAAS_MLSSGID"
				echo "Checking $DATA_DIR and $RESULT_DIR ownership"
				DATA_DIR_UID=$(ls -nd $DATA_DIR | awk '{print $3}')
				DATA_DIR_GID=$(ls -nd $DATA_DIR | awk '{print $4}')
				RESULT_DIR_UID=$(ls -nd $RESULT_DIR | awk '{print $3}')
				RESULT_DIR_GID=$(ls -nd $RESULT_DIR | awk '{print $4}')
				echo "DATA_DIR belongs to $DATA_DIR_UID:$DATA_DIR_GID"
				echo "RESULT_DIR belongs to $RESULT_DIR_UID:$RESULT_DIR_GID"
				if [[ $GID == 0 && $UID == 0 ]];then
					echo "current user is root, ownership check passed"
			    elif [[ $GID == $DLAAS_MLSSGID && $DATA_DIR_GID == $DLAAS_MLSSGID && $RESULT_DIR_GID == $DLAAS_MLSSGID ]];then 
					echo "current user GID matches DATA_DIR and RESULT_DIR ownership, ownership check passed"
				elif [[ $UID == $DATA_DIR_UID && $GID == $DATA_DIR_GID && $UID == $RESULT_DIR_UID && $GID == $RESULT_DIR_GID ]];then
					echo "current user GID and UID matches DATA_DIR and RESULT_DIR ownership, ownership check passed"
				else
					echo "current user GID and UID does not match DATA_DIR and RESULT_DIR ownership, ownership check failed"
					exit 1
				fi;

				echo "Create file for training logs."
				mkdir -p $JOB_STATE_DIR/logs ;
				touch "$JOB_STATE_DIR/logs/training-log.txt"
    			ln -sf "$JOB_STATE_DIR/logs/training-log.txt" "$JOB_STATE_DIR/latest-log"
				# Allow anyone to write to job state dir.
    			# This allows learner processes (which don't run as root) to write to the directory.
    			chmod -R 777 "$JOB_STATE_DIR"

				mkdir -p $JOB_RESULT_DIR/_submitted_code ;
				export LEARNER_ID=${HOSTNAME} ;
				mkdir -p $JOB_RESULT_DIR/learner-$LEARNER_ID ;
				mkdir -p $CHECKPOINT_DIR ;
				ls -l $JOB_RESULT_DIR;
				chown -R $UID:$GID $JOB_RESULT_DIR ;
				chown -R $UID:$GID $CHECKPOINT_DIR ;
				if [[ $GID != 0 || $UID != 0 ]];then
					echo "loading model as $USER_ID with GID $GID and UID $UID"
					groupadd -g $GID $USER_ID
					useradd -g $GID -u $UID $USER_ID
					su $USER_ID
				else
					echo "loading model as root"
				fi;`
		if codeSelector == codeFile {
			loadModelComand += `
				export DEFAULT_TRAINING_DATA_BUCKET=` + defaultTrainedModelsBucket + `;
				export DEFAULT_TRAINING_MODEL_ZIP_BUCKET=` + defaultModelsBucket + `;
				export S3_ENDPOINT=http://minio-prophecis`  + `.` + platformNamespace + `.svc.cluster.local:9000
				python /entrypoint-files/download_model.py $DEFAULT_TRAINING_MODEL_ZIP_BUCKET $TRAINING_ID.zip $JOB_RESULT_DIR/_submitted_code/model.zip ;

				python -m zipfile -e $JOB_RESULT_DIR/_submitted_code/model.zip $MODEL_DIR`
		}

		learnerCommand := learnerBashCommand

		// FIXME MLSS Change: just move the logs, and make sure the owner of the files is the UID user
		var storeLogsCommand string
		if req.JobType == "dist-tf" {
			storeLogsCommand = `
			echo Calling copy logs.
			export LEARNER_ID=${HOSTNAME} ;
			chown -R $UID:$GID $JOB_RESULT_DIR/learner-$LEARNER_ID/
			bash -c 'exit 0'`
		} else {
			storeLogsCommand = `
			echo Calling copy logs.
			export LEARNER_ID=${HOSTNAME} ;
			mv -nf $LOG_DIR/* $JOB_RESULT_DIR/learner-$LEARNER_ID/ ;
			ERROR_CODE=$? ;
			chown -R $UID:$GID $JOB_RESULT_DIR/learner-$LEARNER_ID/
			bash -c 'exit $ERROR_CODE'`
		}
		// logr.Debugf("wrapCommands", loadModelComand, learnerCommand)

		if req.JobType != "dist-tf" {
			cmd = wrapCommands([]containerCommands{
				//{cmd: getCreateTrainCMD(), container: ""},
				{cmd: loadModelComand, container: loadModelContainerName},
				{cmd: learnerCommand, container: learnerContainerName},
				{cmd: storeLogsCommand, container: storeLogsContainerName}})
		} else {
			cmd = wrapCommandsForTFJob([]containerCommands{
				{cmd: getCreateTrainCMD(), container: ""},
				{cmd: loadModelComand, container: loadModelContainerName},
				{cmd: learnerCommand, container: learnerContainerName},
				{cmd: storeLogsCommand, container: storeLogsContainerName}})
		}
	} else {
		command = fmt.Sprintf(`%s mkdir -p $JOB_RESULT_DIR ; bash -c ' train.sh 2>&1 | tee -a %s/latest-log; exit ${PIPESTATUS[0]}'`, command)
		doCondExitWrite = false
		cmd = wrapCommand(command, learnerContainerName, doCondExitWrite)
	}

	// FIXME MLSS Change: add timezone VolumeMount
	//TIMEZONE TODO
	// learnerVolumeMounts = append(learnerVolumeMounts, constants.TimezoneMount)
	// FIXME MLSS Change: add joblogs VolumeMount for fluent-bit
	// learnerVolumeMounts = append(learnerVolumeMounts, constants.JobLogsMount)

	//
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

	logr.Debugf("learner.CreateContainerSpec(container), container: %+v", container.EnvVars)
	learnerContainer := learner.CreateContainerSpec(container)

	extendLearnerContainer(&learnerContainer, req, logr)
	//logr.Debugf("learnerContainer: %v+", learnerContainer)
	// dumpAsYaml("learnerContainer", learnerContainer, logr)

	return learnerContainer
}

func getCreateTrainCMD() string {

	cmd := `
if [ ! -f "/entrypoint-files/train.sh" ];then
echo "now creating file /entrypoint-files/train.sh";
mkdir -p /entrypoint-files;
cat > /entrypoint-files/train.sh <<EOF
#!/bin/bash
echo "Training with training/test data at:"
echo "  DATA_DIR: \$DATA_DIR"
echo "  MODEL_DIR: \$MODEL_DIR"
echo "  TRAINING_JOB: \$TRAINING_JOB"
echo "  TRAINING_COMMAND: \$TRAINING_COMMAND"
echo "Storing trained model at:"
echo "  RESULT_DIR: \$RESULT_DIR"
echo 'Contents of \$MODEL_DIR'
if [-d "/$MODEL_DIR" ];then
ls -la \$MODEL_DIR;
echo 'Contents of \$DATA_DIR'
ls -la \$DATA_DIR
export LEARNER_ID=\${HOSTNAME}
echo "export LEARNER_ID=\${HOSTNAME}" >> ~/.bashrc
touch /job/\$TRAINING_ID
# Switch to model dir
cd \$MODEL_DIR
export PYTHONPATH=\$PYTHONPATH:\$PWD
# env | sort
echo "\$(date): Running training job"
eval "\$TRAINING_COMMAND 2>&1"
cmd_exit=\$?
echo "Training process finished. Exit code: \$cmd_exit"
if [ \${cmd_exit} -ne 0 ];
then
  echo "Job exited with error code \${cmd_exit}"
  exit \${cmd_exit}
fi;
EOF
fi;

if [ ! -f "/entrypoint-files/download_model.py" ];then
echo "now creating file /entrypoint-files/download_model.py";
mkdir -p /entrypoint-files;
cat > /entrypoint-files/download_model.py <<EOF
import os
import sys

import boto3


def main(argv):
    if len(argv) < 4:
        sys.exit("not enough args")
    endpoint = os.environ['S3_ENDPOINT']
    print("now endpoint: ", endpoint)

    s3 = boto3.resource('s3',
                        endpoint_url=endpoint,
                        aws_access_key_id='AKIAIOSFODNN7EXAMPLE',
                        aws_secret_access_key='wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY')
	print(str(argv[2]))
	print(str(argv[3]))
    s3.Bucket(str(argv[1])).download_file(str(argv[2]), str(argv[3]))

    print("Downloaded.")


if __name__ == '__main__':
    main(sys.argv)
EOF
fi;


export PATH=/usr/local/bin/:$PATH; cp /entrypoint-files/*.sh /usr/local/bin/; chmod +x /usr/local/bin/*.sh;
`

	return cmd
}

func getLearnerBashCommandForTFJob() string {
	//su $USER_ID -c 'train.sh >> $JOB_STATE_DIR/latest-log 2>&1 ; exit ${PIPESTATUS[0]}';

	learnerBashCommand := `
			echo "Training with training/test data at, " &&
    		echo "  DATA_DIR, $DATA_DIR" &&
    		echo "  MODEL_DIR, $MODEL_DIR" &&
			echo "  TRAINING_JOB, $TRAINING_JOB" &&
			echo "  TRAINING_COMMAND, $TRAINING_COMMAND" &&
			echo "Storing trained model at," &&
    		echo "  RESULT_DIR, $RESULT_DIR" &&
    		echo 'Contents of $MODEL_DIR' &&
    		echo 'Contents of $DATA_DIR' ;
			ls -la $MODEL_DIR  ; 
    		ls -la $DATA_DIR  ; 
    		export LEARNER_ID=${HOSTNAME}
    		touch /job/$TRAINING_ID ;

			bash -c 'touch /job/$TRAINING_ID' ; 

    		cd $MODEL_DIR;
    		export PYTHONPATH=$PYTHONPATH:$PWD;
    		env | sort 
    		echo "$(date), Running training job";

			bash -c 'eval "$TRAINING_COMMAND " ' ; 

    		cmd_exit=$?
    		echo "Training process finished. Exit code, $cmd_exit";
    		if [ ${cmd_exit} -ne 0 ];
    		then
      		  echo "Job exited with error code ${cmd_exit}";
			  exit ${cmd_exit};
    		fi`

	return learnerBashCommand
}

// Store relationship between a command and the "container" it's associated with.
// The "container" determines the control files used to communicate with the controller.
type containerCommands struct {
	cmd       string
	container string
}

// Wrap a sequence of commands with start and exit files.
func wrapCommands(commands []containerCommands) string {
	var allCommands string

	// FIXME MLSS Temporary Change: TFJOB
	for _, command := range commands {
		allCommands += command.cmd

	}

	//allCommands += `
	//	            while true; do sleep 2; done ;`
	allCommands += `
			echo "exit now"
			exit $TRAINING_RESULT`

	return allCommands
}

// Wrap a sequence of commands with start and exit files.
func wrapCommandsForTFJob(commands []containerCommands) string {
	var allCommands string

	// FIXME MLSS Temporary Change: TFJOB
	for i, command := range commands {

		if i == 0 {
			allCommands += command.cmd
			continue
		}
		//allCommands += buf.String()
		allCommands += command.cmd

	}

	//allCommands += `
	//	            while true; do sleep 2; done ;`
	//allCommands += `
	//	            sleep 6;`

	allCommands += `
			echo "exit now"
			exit $TRAINING_RESULT`

	return allCommands
}

// Wrap a single command with start and exit files.
func wrapCommand(cmd string, containerName string, doCondExitWrite bool) string {

	vars := map[string]string{
		"Name": containerName,
		// "Dir":  controlFilesDirectory,
		"Cmd": cmd,
	}

	var exitWriteStr string
	if doCondExitWrite {
		exitWriteStr = `
		if [ ! -f {{.Dir}}/{{.Name}}.exit ]; then
			echo $main_cmd_status > {{.Dir}}/{{.Name}}.exit
        fi
		`
	} else {
		exitWriteStr = `
		echo $? > {{.Dir}}/{{.Name}}.exit
		`
	}

	var buf bytes.Buffer
	tmpl, _ := template.New("wrapped command").Parse(`
		# Don't repeat if already executed.
		if [ -f {{.Dir}}/{{.Name}}.exit ]; then
			while true; do sleep 1000; done
		fi
		# Wait for start signal.
		while [ ! -f {{.Dir}}/{{.Name}}.start ]; do sleep 2; done
		# Record the start time. Note: In distributed mode, this
		# file will get overwritten by each learner (this is intentional)
		date "+%s%N" | cut -b1-13 > {{.Dir}}/{{.Name}}.start_time
		{{.Cmd}} # do the actual work` + exitWriteStr +
		`while true; do sleep 2; done
	`)
	tmpl.Execute(&buf, vars)

	return buf.String()
}

func constructVolumeClaim(name string, namespace string, volumeSize int64, labels map[string]string) *v1core.PersistentVolumeClaim {
	claim, err := GetVolumeClaim(volumeSize)
	if err != nil {
		logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyLcmService))
		logr.Errorf("constructVolumeClaim return nil, volumeSize == %d", volumeSize)
		return nil
	}
	claim.Name = name
	claim.Namespace = namespace
	claim.Labels = labels
	claim.Spec.AccessModes = []v1core.PersistentVolumeAccessMode{v1core.ReadWriteMany}
	return claim
}

// Return the value of the named environment variable.
func getValue(envVars []v1core.EnvVar, name string) string {
	value := ""
	for _, ev := range envVars {
		if ev.Name == name {
			value = ev.Value
			break
		}
	}
	return value
}

// Checks whether a value is contained in an array
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
