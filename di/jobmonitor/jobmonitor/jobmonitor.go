package jobmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"webank/DI/commons/config"
	"webank/DI/commons/constants"
	cc "webank/DI/commons/controlcenter/client"
	"webank/DI/commons/logger"
	commonmodels "webank/DI/commons/models"
	"webank/DI/commons/service"
	lcmClient "webank/DI/commons/service/client"
	"webank/DI/commons/tfjob"
	"webank/DI/lcm/lcmconfig"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	// "webank/DI/linkisexecutor/linkisclientutils"
	// "webank/DI/linkisexecutor/models"
	"webank/DI/commons/operators/tf-operator/client/clientset/versioned"
	"webank/DI/storage/storage/grpc_storage"
	"webank/DI/trainer/trainer"

	"github.com/cenkalti/backoff"
	"github.com/go-kit/kit/metrics"
	"github.com/spf13/viper"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	//"webank/DI/trainer/client"
	//"webank/DI/trainer/trainer/grpc_trainer_v2"

	"webank/DI/storage/client"
)

const (
	FAILED            = "FAILED"
	COMPLETED         = "COMPLETED"

	CRITICAL = "critical"
	MAJOR    = "major"
	MINOR    = "minor"
	WARNING  = "warning"
	INFO     = "info"
)

type JobMonitor struct {
	ctx                 context.Context
	K8sClient           kubernetes.Interface
	TFJobTrainer        tfjob.Trainer
	repo                trainer.Repository
	numTerminalLearners uint64
	trMap               map[string][]string
	Cli                 *kubernetes.Clientset
	TFClient            *versioned.Clientset
	Alertings           map[string]*Alerting
	//requestBase         models.LinkisRequestBase
}

var failedTrainerConnectivityCounter metrics.Counter

type void struct {}

func NewJobMonitor(cli kubernetes.Interface) (*JobMonitor, error) {
	getLogger := logger.GetLogger()

	//var K8sClient kubernetes.Interface
	var err error
	//if Cli != nil {
	//	K8sClient = Cli
	//} else {
	k8sClient, err := kubernetes.NewForConfig(lcmconfig.GetKubernetesConfig())
	if err != nil {
		getLogger.WithError(err).Errorf("Failed to create a kubernetes client: %v", lcmconfig.GetKubernetesConfig())
		return nil, err
	}
	//}

	tfjobClient := versioned.NewForConfigOrDie(lcmconfig.GetKubernetesConfig())

	getLogger.Debugf("MongoAddressKey: %v, MongoDatabaseKey: %v, MongoUsernameKey: %v, MongoPasswordKey: %v", viper.GetString(constants.MongoAddressKey), viper.GetString(constants.MongoDatabaseKey), viper.GetString(constants.MongoUsernameKey), viper.GetString(constants.MongoAddressKey))

	repo, err := trainer.NewTrainingsRepository(viper.GetString(constants.MongoAddressKey),
		viper.GetString(constants.MongoDatabaseKey), viper.GetString(constants.MongoUsernameKey),
		viper.GetString(constants.MongoPasswordKey), viper.GetString(constants.MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "training_jobs")
	if err != nil {
		getLogger.WithError(err).Fatalf("Cannot create repository with %s %s %s", viper.GetString(constants.MongoAddressKey), viper.GetString(constants.MongoDatabaseKey), viper.GetString(constants.MongoUsernameKey))
	}

	ctx := context.TODO()

	TFJobTrainer := tfjob.NewTensorFlowJobTrainer(k8sClient, "")

	alertings := make(map[string]*Alerting, 0)

	monitor := JobMonitor{
		K8sClient:    k8sClient,
		repo:         repo,
		ctx:          ctx,
		trMap:        initTransitionMap(),
		TFJobTrainer: TFJobTrainer,
		Cli:          k8sClient,
		TFClient:     tfjobClient,
		Alertings:    alertings,
	}

	return &monitor, nil
}

func (jm *JobMonitor) Start() {
	go jm.monitoring()
}

func (jm *JobMonitor) monitoring() {
	//go jm.monitoringSingleLearnerJob()
	//MLFlow Job & Single Job
	go jm.monitoringMLFlowJob()
	go jm.monitoringTFJob()
	// go jm.monitoringLinkisJob()
}

func (jm *JobMonitor) handleErrorJob() {
	//Init
	logger := logger.GetLogger()
	var labels = make(map[string]string)
	labels[constants.ENVIR] = os.Getenv(constants.ENVIR_UPRER)
	labels[constants.JOBTYPE] = constants.SINGLE

	for {
		select {
		case <-jm.ctx.Done():
			return
		default:
			//1. Get Running Job from Mongo: Set A
			jobListFromMongo,err := jm.repo.FindAllByStatus(grpc_trainer_v2.Status_RUNNING,"","")
			if err != nil {
				logger.Errorf("GetNormalJobListFromK8s failed, %v", err)
			}

			//2. Get Running Job From Kubernetes: Set B
			jobListFromK8s, err := jm.GetNormalJobListFromK8s(labels)
			if err != nil {
				logger.Errorf("GetNormalJobListFromK8s failed, %v", err)
			}

			//3. Cal: A-B
			var trainingIdSet map[string]void
			var errorJobList []string
			var setMember void
			for i := 0; i < len(jobListFromMongo.TrainingRecords); i++ {
				trainingIdSet[jobListFromMongo.TrainingRecords[i].TrainingID] = setMember
			}
			for i := 0; i < len(jobListFromK8s); i++ {
				if _,ok := trainingIdSet[jobListFromK8s[i].Name]; ok{
					errorJobList = append(errorJobList, jobListFromK8s[i].Name)
				}
			}

			//4. Handle Error Set A-B
			for i := 0; i < len(errorJobList); i++ {
				record, err := jm.repo.Find(errorJobList[i])
				if err != nil {
					logger.Errorf("GetNormalJobListFromK8s failed, %v", err)
				}

				updateStatus := &client.TrainingStatusUpdate{
					Status: grpc_storage.Status_FAILED,
					StatusMessage: "Job is Not Found in Kubernetes platform",
				}

				//Update Mongo Job to Failed
				err = updateJobStatusInTrainer(record.TrainingID, record.UserID, updateStatus, logger)
				if err != nil{
					logger.Errorf("GetNormalJobListFromK8s failed, %v", err)
				}

			}
			time.Sleep(10 * time.Second)
		}
	}



}


func (jm *JobMonitor) monitoringMLFlowJob() {
	go jm.monitoringNormalMLFlowJob()
	go jm.handleErrorJob()

}
func (jm *JobMonitor) monitoringNormalMLFlowJob() {
	getLogger := logger.GetLogger()
	getLogger.Debugf("monitoringSingleLearnerJob start")
	var labels = make(map[string]string)
	labels[constants.ENVIR] = os.Getenv(constants.ENVIR_UPRER)
	labels[constants.JOBTYPE] = constants.SINGLE
	getLogger.Debugf("labels: %+v", labels)
	for {
		select {
		case <-jm.ctx.Done():
			return
		default:
			//get job by envir and jobType from k8s
			jobList, err := jm.GetNormalJobListFromK8s(labels)
			if err != nil {
				getLogger.Errorf("GetNormalJobListFromK8s failed, %v", err)
			}
			err = jm.updateStatusInDB(jobList)
			if err != nil {
				getLogger.Errorf("updateStatusInDB failed, %v", err)
			}
			time.Sleep(10 * time.Second)
		}
	}
}





func (jm *JobMonitor) monitoringTFJob() {
	getLogger := logger.GetLogger()
	getLogger.Debugf("monitoringTFJob start")

	for {
		select {
		case <-jm.ctx.Done():
			return
		default:
			tfJobs, err := jm.getTFJobListFromK8s()
			if err != nil {
				//return nil
				getLogger.Debugf("getTFJobListFromK8s failed, %v", err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			err = jm.updateStatusInDBForTFJob(tfJobs)
			if err != nil {
				//return nil
				getLogger.Debugf("updateStatusInDBForTFJob failed, %v", err.Error())
			}
			time.Sleep(5 * time.Second)
		}
	}

}

func (jm *JobMonitor) GetNormalJobListFromK8s(labels map[string]string) ([]v1.Job, error) {
	getLogger := logger.GetLogger()

	labelStr := parseLabelString(labels)
	options := metav1.ListOptions{
		LabelSelector: labelStr,
	}
	jobList := v1.JobList{}
	err := jm.K8sClient.BatchV1().RESTClient().Get().Resource("jobs").VersionedParams(&options, scheme.ParameterCodec).Do().Into(&jobList)
	if err != nil {
		getLogger.Debugf(err.Error())
		return nil, err
	}
	getLogger.Debugf("jobList size: %v", len(jobList.Items))
	return jobList.Items, nil
}

func parseLabelString(labels map[string]string) string {
	var labelList []string
	for k, v := range labels {
		s := fmt.Sprintf("%v=%v", k, v)
		labelList = append(labelList, s)
	}
	result := strings.Join(labelList, ",")
	return result
}

func (jm *JobMonitor) getTFJobListFromK8s() ([]tfjob.TrainingJob, error) {
	getLogger := logger.GetLogger()
	//jobs, e := jm.TFJobTrainer.GetAllTrainingJob()
	getLogger.Debugf("getTFJobListFromK8s start")

	//jobs, err := tfjob.ListAllTFJob(jm.Cli, jm.TFJobTrainer)

	err := tfjob.GetTFJobFromk8sToOriginTFJob(jm.TFClient)
	if err != nil {
		getLogger.Errorf("GetTFJobFromk8sToOriginTFJob failed, %v", err.Error())
	}

	jobs, err := jm.TFJobTrainer.ListTrainingJobs()
	if err != nil {
		getLogger.Errorf("ListAllTFJob failed, %v", err.Error())
	}

	for i, j := range jobs {
		getLogger.Debugf("NewTensorFlowJobTrainer tfjobList[%v], namespaceL %v, name: %v", i, j.Namespace(), j.Name())
	}

	//notFinishedTFJobs := getNotFinishedTFJob(jobs)

	return jobs, nil
	//return nil, nil
}

func getNotFinishedTFJob(jobs []tfjob.TrainingJob) []tfjob.TrainingJob {
	var result []tfjob.TrainingJob
	for _, j := range jobs {
		if !strings.EqualFold(j.GetStatus(), "FAILED") && !strings.EqualFold(j.GetStatus(), "SUCCEEDED") {
			result = append(result, j)
		}
	}

	return result
}

// func (jm *JobMonitor) getUnfinishedLinkisJobFromDB() ([]*trainer.TrainingRecord, error) {
// 	getLogger := logger.GetLogger()

// 	records, err := jm.repo.FindAllNotCompletedLinkisJobs(true, false)
// 	if err != nil {
// 		getLogger.Errorf("jm.repo.Find(JobId) failed, %v", err.Error())
// 		return nil, err
// 	}
// 	return records, nil
// }

func (jm *JobMonitor) updateStatusInDB(list []v1.Job) error {
	getLogger := logger.GetLogger()

	for _, j := range list {
		startTime := j.Status.StartTime
		JobId := j.Labels["training_id"]

		currentJobStatus := j.Status.String()
		actualStatus, err := GetStatusFromJob(j, getLogger)
		getLogger.Debugf("updateStatusInDB, jobId Status: %v , %v", JobId, currentJobStatus)
		if err != nil {
			getLogger.Errorf("GetStatusFromJob(j, getLogger) failed, %v", err.Error())
			continue
		}
		getLogger.Debugf("updateStatusInDB, currentJobStatus: %v, actualStatus: %v", currentJobStatus, actualStatus.Status)

		record, err := jm.repo.Find(JobId)
		if err != nil {
			getLogger.Errorf("jm.repo.Find(JobId) failed, %v", err.Error())
			continue
		}
		getLogger.Debugf("jm.repo.Find(JobId): %v, record.TrainingStatus.Status: %v", record.TrainingID, record.TrainingStatus.Status.String())
		getLogger.Debugf("jm.repo.Find(JobId): %v, record.JobAlert: %+v", record.TrainingID, record.JobAlert)

		err = jm.processUpdateLearnerStatus(record, startTime.Time, actualStatus, getLogger)
		if err != nil {
			getLogger.Errorf("jm processUpdateLearnerStatus failed, %v", err.Error())
			return err
		}
		getLogger.Debugf("uptdated jobId: %v, record.TrainingStatus.Status: %v", record.TrainingID, record.TrainingStatus.Status.String())

		update, err := GetStatusFromJob(j, getLogger)
		propagation := metav1.DeletePropagationBackground
		options := metav1.DeleteOptions{
			PropagationPolicy: &propagation,
		}
		if strings.EqualFold(update.Status.String(), grpc_storage.Status_COMPLETED.String()) ||
			strings.EqualFold(update.Status.String(), grpc_storage.Status_FAILED.String()) {
			getLogger.Infof("job: %v done, cleaning up now", record.TrainingID)

			//删除对应 job alerting对象
			delete(jm.Alertings, record.TrainingID)

			go func() {
				//IF Job Finish
				time.Sleep(10 * time.Second)
				getLogger.Debugf(" Deleting Job '%s'", record.TrainingID)
				error := updateJobStatusInTrainer(record.TrainingID, record.UserID, actualStatus, getLogger)
				if error != nil {
					getLogger.WithError(error).Errorf("Failed to write the status %s for training %s to trainer", actualStatus.Status.String(), record.TrainingID)
				}
				err := jm.K8sClient.BatchV1().Jobs(j.Namespace).Delete(j.Name, &options)
				if err != nil {
					getLogger.Errorf("deleting job: %v in l8s failed, %v", j.Name, err.Error())
				}
			}() //todo error marked by lk
		}

		jm.HandlerAlertForDeadlineAndOvertimeJobStatus(actualStatus.Status.String(), startTime.Time, record, getLogger)

	}
	//submit
	return nil
}

func (jm *JobMonitor) processUpdateLearnerStatus(tr *trainer.TrainingRecord, startTime time.Time, update *client.TrainingStatusUpdate, logr *logger.LocLoggingEntry) error {

	learnerStatus := update.Status

	//learnerStatus = GetStatus(learnerStatusValue, logr).Status
	logr.Infof("got triggered with the current value %s (beforeStatus %s)", tr.TrainingStatus.Status.String(), learnerStatus)

	beforeStatus := tr.TrainingStatus.Status

	if jm.isTransitionAllowed(beforeStatus.String(), learnerStatus.String()) {
		logr.Infof("Transition was allowed, changing overall beforeStatus of job from %s to learners beforeStatus %s", beforeStatus, learnerStatus)
		//_, _ = jm.EtcdClient.CompareAndSwap(overallJobStatusPath(jm.TrainingID), learnerStatusValue, currentOverallJobStatus, logr)
		//FIXME MLSS Change: trigger alert by CC
		logr.Debugf("processUpdateLearnerStatus start to HandlerAlertForJobStatus time: %v", time.Now())
		_ = jm.HandlerAlertForJobStatus(tr, startTime, learnerStatus.String(), logr)
		updated := jm.processUpdateJobStatus(tr, update, logr)
		logr.Infof("jm.processUpdateJobStatus result: %v", updated)

	} else {
		logr.Warnf("Transition not allowed job from overall job beforeStatus %s to learner beforeStatus %s", beforeStatus, learnerStatus)
	}
	//keep an eye on idividual learners as well, if they terminate then check if all of them are done then check if job can be terminated
	if learnerStatus == grpc_storage.Status_COMPLETED || learnerStatus == grpc_storage.Status_FAILED {
		atomic.AddUint64(&jm.numTerminalLearners, 1)
	}
	return nil
}

func (jm *JobMonitor) HandlerAlertForJobStatus(tr *trainer.TrainingRecord, startTime time.Time, currStatus string, logr *logger.LocLoggingEntry) error {
	getLogger := logger.GetLogger()
	getLogger.Debugf("HandlerAlertForJobStatus job: %v, startTime: %v", tr.TrainingID, startTime.String())
	userId := tr.UserID
	jobNamespace := tr.Namespace
	jobName := tr.JobID
	trainingId := tr.TrainingID

	alerting, err := jm.getJobAlerting(tr, startTime)
	if err != nil {
		getLogger.Debugf("jm.getJobAlerting failed, %v", err.Error())
		return err
	}
	jobAlert := alerting.AlertRequest.JobAlert
	logr.Debugf("alerting.AlertRequest.JobAlert: %+v", jobAlert)

	logr.Infof("userId: %v, jobNamespace: %v, jobName: %v, trainingId: %v", userId, jobNamespace, jobName, trainingId)
	//alertRequest := jm.alertRequest
	statusUpdate := client.GetStatus(currStatus, logr)
	status := statusUpdate.Status

	dto := commonmodels.AlertDTO{
		UserID:       userId,
		JobNamespace: jobNamespace,
		JobName:      jobName,
		TrainingID:   trainingId,
		StartTime:    startTime.Format("2006-01-02 15:04:05"),
	}

	//jobAlert := alertRequest.JobAlert
	event := jobAlert["event"]
	if nil != event && len(event) > 0 {
		logr.Infof("HandlerAlertForJobStatus rangeFor eventMap: %v", event)
		for _, eMap := range event {
			alertLevel := eMap["alert_level"]
			receiver := eMap["receiver"]
			eventChecker := eMap["event_checker"]
			if (eventChecker == FAILED && status.String() == FAILED) || (eventChecker == COMPLETED && status.String() == COMPLETED) {
				dto.EndTime = time.Now().Format("2006-01-02 15:04:05")
				dto.AlertLevel = alertLevel
				dto.Receiver = receiver
				logr.Infof("HandlerAlertForJobStatus to alert")
				eventErr := alert(dto, status.String(), "job"+eventChecker, logr)
				if eventErr != nil {
					return eventErr
				}
			}
		}
	}
	return nil
}

func (jm *JobMonitor) processUpdateJobStatus(tr *trainer.TrainingRecord, actualStatus *client.TrainingStatusUpdate, logr *logger.LocLoggingEntry) bool {
	logr.Infof("(processUpdateJobStatus) got triggered with the current status %s", tr.TrainingStatus.Status.String())
	//Variable to notify whether the job needs further status monitoring
	markComplete := false

	error := updateJobStatusInTrainer(tr.TrainingID, tr.UserID, actualStatus, logr)
	if error != nil {
		logr.WithError(error).Errorf("Failed to write the status %s for training %s to trainer", actualStatus.Status.String(), tr.TrainingID)
	}

	return markComplete
}

func updateJobStatusInTrainer(trainingID string, userID string, statusUpdate *client.TrainingStatusUpdate, logr *logger.LocLoggingEntry) error {
	updStatus := statusUpdate.Status
	logr.Infof("(updateJobStatus) Updating status of %s to %s", trainingID, updStatus.String())
	//updateRequest := &grpc_trainer_v2.UpdateRequest{TrainingId: trainingID, Status: updStatus, Timestamp: statusUpdate.Timestamp,
	//	UserId: userID, StatusMessage: statusUpdate.StatusMessage, ErrorCode: statusUpdate.ErrorCode}
	updateRequest := &grpc_storage.UpdateRequest{TrainingId: trainingID, Status: updStatus, Timestamp: statusUpdate.Timestamp,
		UserId: userID, StatusMessage: statusUpdate.StatusMessage, ErrorCode: statusUpdate.ErrorCode}
	//trainer, err := client.NewTrainer()
	storageClient, err := client.NewStorage()

	if err != nil {
		logr.WithError(err).Errorf("(updateJobStatus) Creating training client for status update failed. Training ID %s New Status %s", trainingID, updStatus.String())
	}
	defer storageClient.Close()

	defaultBackoff := backoff.NewExponentialBackOff()
	defaultBackoff.MaxElapsedTime = 1 * time.Minute
	defaultBackoff.MaxInterval = 5 * time.Second

	err = backoff.RetryNotify(func() error {
		_, err = storageClient.Client().UpdateTrainingJob(context.Background(), updateRequest)
		return err
	}, defaultBackoff, func(err error, t time.Duration) {
		logr.WithError(err).Errorf("Failed to update status to the trainer. Retrying WARNING: Status updates for %s may be temporarily inconsistent due to failure to communicate with Trainer.", trainingID)
	})

	if err != nil {
		failedTrainerConnectivityCounter.Add(1)
		logr.WithError(err).Errorf("Failed to update status to the trainer. Already retried several times.WARNING : Status of job %s will likely be incorrect", trainingID)
		return err
	}

	return err
}

func (jm *JobMonitor) isTransitionAllowed(fromStatus string, toStatus string) bool {
	validFroms := jm.trMap[toStatus]
	for _, allowed := range validFroms {
		if fromStatus == allowed {
			return true
		}
	}
	return false
}

func (jm *JobMonitor) updateStatusInDBForTFJob(jobs []tfjob.TrainingJob) error {
	getLogger := logger.GetLogger()
	for _, j := range jobs {
		startTime := j.StartTime()

		JobId := j.Name()
		getLogger.Debugf("updateStatusInDB, jobId: %v", JobId)
		currentJobStatus := j.GetStatus()
		getLogger.Debugf("updateStatusInDB, currentJobStatus: %v", currentJobStatus)

		statusUpdate, err := GetStatusFromJobForTFjob(j, getLogger)
		if err != nil {
			getLogger.Errorf("jm.repo.Find(JobId) failed, %v", err.Error())
			//return nil, err
			continue
		}

		record, err := jm.repo.Find(JobId)
		if err != nil {
			getLogger.Errorf("jm.repo.Find(JobId) failed, %v", err.Error())
			//return nil, err
			continue
		}
		getLogger.Debugf("jm.repo.Find(JobId): %v, record.TrainingStatus.Status: %v", record.TrainingID, record.TrainingStatus.Status.String())

		err = jm.processUpdateLearnerStatus(record, startTime.Time, statusUpdate, getLogger)
		if err != nil {
			getLogger.Errorf("jm processUpdateLearnerStatus failed, %v", err.Error())
			continue
		}
		getLogger.Debugf("uptdated jobId: %v, record.TrainingStatus.Status: %v", record.TrainingID, record.TrainingStatus.Status.String())

		update, err := GetStatusFromJobForTFjob(j, getLogger)

		getLogger.Debugf("update.Status.String(): %v", update.Status.String())
		if strings.EqualFold(update.Status.String(), grpc_storage.Status_COMPLETED.String()) || strings.EqualFold(update.Status.String(), grpc_storage.Status_FAILED.String()) {
			getLogger.Infof("tfjob: %v done, cleaning up now", record.TrainingID)
			err := KillDeployedJob(record.TrainingID, record.UserID, record.JobID, record.Namespace, getLogger)
			if err != nil {
				getLogger.Errorf("KillDeployedJob tfjob: %v in k8s failed, %v", j.Name, err.Error())
				continue
			}
		}

	}

	//submit
	return nil
}

//KillDeployedJob ... Contact the LCM and kill training job
func KillDeployedJob(trainingID string, userID string, jobName string, jobNamespace string, logr *logger.LocLoggingEntry) error {
	time.Sleep(10 * time.Second)
	logr.Infof("(killDeployedJob) Sending job kill request to LCM for %s", trainingID)
	jobKillReq := &service.JobKillRequest{Name: jobName, TrainingId: trainingID, UserId: userID, JobNamespace: jobNamespace}
	lcm, err := lcmClient.NewLcm(nil)
	if err != nil {
		logr.Errorln("(KillDeployedJob) Cannot create lcm service client: ", err.Error())
		return err
	}
	defer lcm.Close()

	defaultBackoff := backoff.NewExponentialBackOff()
	defaultBackoff.MaxElapsedTime = 1 * time.Minute
	defaultBackoff.MaxInterval = 5 * time.Second

	err = backoff.Retry(func() error {
		_, err = lcm.Client().KillTrainingJob(context.Background(), jobKillReq)
		if err != nil {
			logr.WithError(err).Errorf("Failed to send request to LCM to garbage collect Training Job %s. Retrying", trainingID)
		}
		return err
	}, defaultBackoff)

	if err != nil {
		logr.WithError(err).Errorf("(killDeployedJob) Successfully sent request to LCM to garbage collect Failed to send request to LCM to garbage collect Training Job %s. Already retried several times.", trainingID)
		return err
	}

	return err
}

func tranformStatus(status v1.JobStatus) string {
	if status.Failed > 0 {
		return "FAILED"
	} else if status.Succeeded > 0 && status.Active == 0 {
		return "COMPLETED"
	} else if status.Active > 0 {
		return "RUNNING"
	} else {
		return "PENDING"
	}

}

func GetStatusFromJob(job v1.Job, logr *logger.LocLoggingEntry) (*client.TrainingStatusUpdate, error) {
	getLogger := logger.GetLogger()

	statusFromK8s := tranformStatus(job.Status)
	//unify PROCESSING status in single job
	if strings.EqualFold(statusFromK8s, "running") {
		statusFromK8s = "RUNNING"
	}
	//unify COMPLETED status in single job
	if strings.EqualFold(statusFromK8s, "SUCCEEDED") || strings.EqualFold(statusFromK8s, "SUCCEED") {
		statusFromK8s = "COMPLETED"
	}

	upper := strings.ToUpper(statusFromK8s)
	getLogger.Debugf("statusFromK8s: %v,upper: %v", statusFromK8s, upper)

	var updStatus grpc_storage.Status
	switch upper {
	case grpc_storage.Status_PENDING.String():
		updStatus = grpc_storage.Status_PENDING
	case grpc_storage.Status_FAILED.String():
		updStatus = grpc_storage.Status_FAILED
	case grpc_storage.Status_PROCESSING.String():
		updStatus = grpc_storage.Status_PROCESSING
	case grpc_storage.Status_RUNNING.String():
		updStatus = grpc_storage.Status_RUNNING
	case grpc_storage.Status_QUEUED.String():
		updStatus = grpc_storage.Status_QUEUED
	case grpc_storage.Status_COMPLETED.String():
		updStatus = grpc_storage.Status_COMPLETED
	}
	getLogger.Debugf("after updStatus: %v", updStatus.String())

	result := client.TrainingStatusUpdate{}
	result.ErrorCode = ""
	result.Status = updStatus
	result.StatusMessage = ""
	result.Timestamp = client.CurrentTimestampAsString()
	return &result, nil
}

type revision struct {
	Status string `json:"status"`
}

func GetStatusFromJobForTFjob(job tfjob.TrainingJob, logr *logger.LocLoggingEntry) (*client.TrainingStatusUpdate, error) {
	getLogger := logger.GetLogger()
	statusFromK8s := job.GetStatus()

	//unify PROCESSING status in single job
	if strings.EqualFold(statusFromK8s, "running") {
		statusFromK8s = "RUNNING"
	}
	//unify COMPLETED status in single job
	if strings.EqualFold(statusFromK8s, "SUCCEEDED") {
		statusFromK8s = "COMPLETED"
	}

	upper := strings.ToUpper(statusFromK8s)
	getLogger.Debugf("statusFromK8s: %v,upper: %v", statusFromK8s, upper)

	var updStatus grpc_storage.Status
	switch upper {
	case grpc_storage.Status_PENDING.String():
		updStatus = grpc_storage.Status_PENDING
	case grpc_storage.Status_FAILED.String():
		updStatus = grpc_storage.Status_FAILED
	case grpc_storage.Status_PROCESSING.String():
		updStatus = grpc_storage.Status_PROCESSING
	case grpc_storage.Status_RUNNING.String():
		updStatus = grpc_storage.Status_RUNNING
	case grpc_storage.Status_QUEUED.String():
		updStatus = grpc_storage.Status_QUEUED
	case grpc_storage.Status_COMPLETED.String():
		updStatus = grpc_storage.Status_COMPLETED
	}
	getLogger.Debugf("after updStatus: %v", updStatus.String())

	result := client.TrainingStatusUpdate{}
	result.ErrorCode = ""
	result.Status = updStatus
	result.StatusMessage = ""
	result.Timestamp = client.CurrentTimestampAsString()
	return &result, nil
}

func initTransitionMap() map[string][]string {
	transistionMap := make(map[string][]string)
	allowQUEUED := []string{grpc_storage.Status_QUEUED.String()}
	allowPROCESSING := []string{grpc_storage.Status_QUEUED.String(), grpc_storage.Status_PROCESSING.String(), grpc_storage.Status_PENDING.String(), grpc_storage.Status_NOT_STARTED.String()}
	allowRUNNING := []string{grpc_storage.Status_QUEUED.String(), grpc_storage.Status_PROCESSING.String(), grpc_storage.Status_PENDING.String(), grpc_storage.Status_NOT_STARTED.String()}
	allowCOMPLETED := []string{grpc_storage.Status_RUNNING.String(), grpc_storage.Status_QUEUED.String(), grpc_storage.Status_PROCESSING.String(), grpc_storage.Status_PENDING.String(), grpc_storage.Status_NOT_STARTED.String()}
	allowFAILED := []string{grpc_storage.Status_QUEUED.String(), grpc_storage.Status_PROCESSING.String(), grpc_storage.Status_PENDING.String(), grpc_storage.Status_NOT_STARTED.String()}

	transistionMap[grpc_storage.Status_QUEUED.String()] = allowQUEUED
	transistionMap[grpc_storage.Status_PROCESSING.String()] = allowPROCESSING
	transistionMap[grpc_storage.Status_COMPLETED.String()] = allowCOMPLETED
	transistionMap[grpc_storage.Status_RUNNING.String()] = allowRUNNING
	transistionMap[grpc_storage.Status_FAILED.String()] = allowFAILED

	return transistionMap
}

func alert(dto commonmodels.AlertDTO, jobStatus string, alertReason string, logr *logger.LocLoggingEntry) error {
	dto.JobStatus = jobStatus
	dto.AlertReason = alertReason
	ccAddress := viper.GetString(config.CCAddress)
	ccClient := cc.GetCcClient(ccAddress)
	logr.Debugf("start client to cc with ccAddress:" + ccAddress)
	body, e := ccClient.AlertTraining(dto)
	if e != nil {
		logr.Error("AlertTraining failed. %v", e.Error())
		return e
	}
	logr.Infof("AlertTraining to CC success, AlertTraining body: %v", body)
	return nil
}

func (jm *JobMonitor) HandlerAlertForDeadlineAndOvertimeJobStatus(currStatus string, startTime time.Time, tr *trainer.TrainingRecord, logr *logger.LocLoggingEntry) error {
	getLogger := logger.GetLogger()
	getLogger.Debugf("HandlerAlertForDeadlineAndOvertimeJobStatus job: %v, startTime: %v", tr.TrainingID, startTime.String())

	userId := tr.UserID
	jobNamespace := tr.Namespace
	jobName := tr.JobID
	trainingId := tr.TrainingID

	//startTime := tr.alertRequest.StartTime
	//alerting := jm.Alertings[tr.TrainingID]
	//alertRequest := alerting.AlertRequest
	alerting, err := jm.getJobAlerting(tr, startTime)
	if err != nil {
		return err
	}
	alertRequest := alerting.AlertRequest

	logr.Infof("userId: %v, jobNamespace: %v, jobName: %v, trainingId: %v", userId, jobNamespace, jobName, trainingId)
	//alertRequest := jm.alertRequest
	statusUpdate := client.GetStatus(currStatus, logr)
	status := statusUpdate.Status

	dto := commonmodels.AlertDTO{
		UserID:       userId,
		JobNamespace: jobNamespace,
		JobName:      jobName,
		TrainingID:   trainingId,
		StartTime:    startTime.Format("2006-01-02 15:04:05"),
	}

	jobAlert := alertRequest.JobAlert
	deadline := jobAlert["deadline"]
	overtime := jobAlert["overtime"]

	now := time.Now()
	if nil != deadline && len(deadline) > 0 {
		logr.Infof("HandlerAlertForJobStatus rangeFor deadlineMap: %v", deadline)
		for _, dMap := range deadline {
			alertLevel := dMap["alert_level"]
			interval, atoErr := strconv.Atoi(dMap["interval"])
			if atoErr != nil {
				return errors.New("dMap interval parse int failed")
			}
			receiver := dMap["receiver"]
			deadlineChecker := dMap["deadline_checker"]
			if deadlineChecker != "" {
				deadlineTime, parseErr := time.ParseInLocation("2006-01-02 15:04:05", deadlineChecker, time.Local)
				if parseErr != nil {
					logr.Error("deadlineChecker time Parse failed.")
					return parseErr
				}
				logr.Infof("HandlerAlertForJobStatus rangeFor deadline: %v, now: %v, deadline.Before: %v, lastAlertTimeForDeadlineMap[deadlineChecker]: %v", deadlineTime, now, deadlineTime.Before(now), alerting.LastAlertTimeForDeadlineMap[deadlineChecker])
				if deadlineTime.Before(now) && status.String() != FAILED && status.String() != COMPLETED {
					if alerting.LastAlertTimeForDeadlineMap[deadlineChecker].Before(now) {
						logr.Infof("HandlerAlertForJobStatus rangeFor lastAlertTimeForDeadline: %v, lastAlertTimeForDeadlineMap[deadlineChecker].Before: %v", alerting.LastAlertTimeForDeadlineMap[deadlineChecker], alerting.LastAlertTimeForDeadlineMap[deadlineChecker].Before(now))
						dto.AlertLevel = alertLevel
						dto.Receiver = receiver
						err := alert(dto, status.String(), "The training job "+trainingId+" timed out for deadlineChecker: "+deadlineChecker+" and nowTime: "+now.Format("2006-01-02 15:04:05"), logr)
						if err != nil {
							return err
						}
						alerting.LastAlertTimeForDeadlineMap[deadlineChecker] = now.Add(time.Duration(interval) * time.Minute)
					}
				}
			}
		}
	}

	if nil != overtime && len(overtime) > 0 {
		logr.Infof("HandlerAlertForJobStatus rangeFor overtimeMap: %v", overtime)
		for _, oMap := range overtime {
			alertLevel := oMap["alert_level"]
			interval, atoErr := strconv.Atoi(oMap["interval"])
			if atoErr != nil {
				return errors.New("overtime interval parse int failed")
			}
			receiver := oMap["receiver"]
			overtimeChecker := oMap["overtime_checker"]

			if overtimeChecker != "" {
				f, err := strconv.ParseFloat(overtimeChecker, 64)
				if err != nil {
					logr.Error("overtimeChecker float Parse failed.")
					return err
				}
				duration := int64(math.Floor(f * 60))
				over := startTime.Add(time.Duration(duration) * time.Minute)
				logr.Debugf("HandlerAlertForJobStatus for over: %v, now: %v, duration, lastAlertTimeForOvertimeMap[overtimeChecker]: %v", over, now, duration, alerting.LastAlertTimeForOvertimeMap[overtimeChecker])
				if over.Before(now) && status.String() != FAILED && status.String() != COMPLETED {
					if alerting.LastAlertTimeForDeadlineMap[overtimeChecker].Before(now) {
						dto.AlertLevel = alertLevel
						dto.Receiver = receiver
						logr.Debugf("HandlerAlertForJobStatus overtime to alert lastAlertTimeForOvertime: %v, now: %v", alerting.LastAlertTimeForOvertimeMap[overtimeChecker], now)
						err := alert(dto, status.String(), "The training job "+trainingId+" timed out for overtime."+overtimeChecker+"and nowTime: "+now.Format("2006-01-02 15:04:05"), logr)
						if err != nil {
							return err
						}
						alerting.LastAlertTimeForOvertimeMap[overtimeChecker] = now.Add(time.Duration(interval) * time.Minute)
					}
				}
			}
		}
	}
	return nil
}

func (jm *JobMonitor) getJobAlerting(tr *trainer.TrainingRecord, startTime time.Time) (*Alerting, error) {
	getLogger := logger.GetLogger()
	var alerting Alerting
	getLogger.Debugf("getJobAlerting tr.TrainingID, %v", tr.TrainingID)
	getLogger.Debugf("getJobAlerting jm.Alertings[tr.TrainingID], %+v", jm.Alertings[tr.TrainingID])

	if jm.Alertings[tr.TrainingID] == nil {
		getLogger.Debugf("jm.Alertings[tr.TrainingID] == nil")
		jobAlertStr := tr.JobAlert
		bytes := []byte(jobAlertStr)
		var jobAlert = map[string][]map[string]string{}
		getLogger.Debugf("main JOB_ALERT: %v", jobAlertStr)
		err := json.Unmarshal(bytes, &jobAlert)
		if err != nil {
			return nil, err
		}
		getLogger.Debugf("getJobAlerting json.Unmarshal, &jobAlert: %+v", jobAlert)

		alerting = Alerting{
			LastAlertTimeForOvertimeMap: make(map[string]time.Time),
			LastAlertTimeForDeadlineMap: make(map[string]time.Time),
			AlertRequest: commonmodels.AlertRequest{
				StartTime: startTime,
				JobAlert:  jobAlert,
			},
		}
		getLogger.Debugf("getJobAlerting alerting.AlertRequest.JobAlert: %+v", alerting.AlertRequest.JobAlert)

		jm.Alertings[tr.TrainingID] = &alerting
		getLogger.Debugf("getJobAlerting jm.Alertings[tr.TrainingID].AlertRequest.JobAlert: %+v", jm.Alertings[tr.TrainingID].AlertRequest.JobAlert)

	} else {
		getLogger.Debugf("jm.Alertings[tr.TrainingID] excited")
		jm.Alertings[tr.TrainingID].AlertRequest.StartTime = startTime
		getLogger.Debugf("getJobAlertingm.Alertings[tr.TrainingID].AlertRequest.StartTime: %+v", jm.Alertings[tr.TrainingID].AlertRequest.StartTime)
		alerting = *jm.Alertings[tr.TrainingID]
	}

	return &alerting, nil
}

type Alerting struct {
	AlertRequest                commonmodels.AlertRequest
	LastAlertTimeForDeadlineMap map[string]time.Time
	LastAlertTimeForOvertimeMap map[string]time.Time
}
