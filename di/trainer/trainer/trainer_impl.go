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

package trainer

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"sync"
	"webank/DI/commons/util/random"

	// "webank/DI/linkisexecutor/linkisclientutils"
	// "webank/DI/linkisexecutor/models"

	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/commons/service/client"
	tdsClient "webank/DI/metrics/client"
	tdsService "webank/DI/metrics/service/grpc_training_data_v1"
	trainerClient "webank/DI/trainer/client"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"time"

	"regexp"
	"strconv"
	"strings"

	"errors"
	"webank/DI/trainer/instrumentation"
	"webank/DI/trainer/storage"

	storageClient "webank/DI/storage/client"
	"webank/DI/storage/storage/grpc_storage"
)

const internalObjectStoreID = "dlaas_internal_os"

const (
	modelsBucketKey        = "objectstore.bucket.models"
	trainedModelsBucketKey = "objectstore.bucket.trainedmodels"

	//defaultModelsBucket        = "dlaas-models"
	//defaultTrainedModelsBucket = "dlaas-trained-models"
	defaultModelsBucket        = "di-models"
	defaultTrainedModelsBucket = "di-trained-models"

	collectionNameTrainingJobs = "training_jobs"
	collectionNameJobHistory   = "job_history"

	debugLogsMode = false

	oldEndpointInternalPageSize = 10

	mongoAddressKey  = "MONGO_ADDRESS"
	mongoDatabaseKey = "MONGO_DATABASE"
	mongoUsernameKey = "MONGO_USERNAME"
	mongoPasswordKey = "MONGO_PASSWORD"
	MongoAuthenticationDatabase = "MONGO_Authentication_Database"

	gpuLimitsKey          = "gpu.limits"
	gpuLimitsQuerySizeKey = "gpu.limits.query.size"

	pollIntervalKey = "queue.poll.interval"

	bdp          = "bdp"
	bdap         = "bdap"
	bdapsafe     = "bdapsafe"
	safe         = "safe"
	defaultValue = "default"
)

const (
	// Used by a counter metric.
	dlaasStoreKind = "dlaas"
	userStoreKind  = "user"
)

// Confuse `go vet' to not check this `Errorf' call. :(
// See https://github.com/grpc/grpc-go/issues/90
var gerrf = status.Errorf

//var shortidGenHigh = shortid.Shortid{}
//var shortidGenLow = shortid.Shortid{}

// Service represents the functionality of the trainer service
type Service interface {
	grpc_trainer_v2.TrainerServer
	service.LifecycleHandler
	StopTrainer()
}


type queueHandler struct {
	stopQueue chan struct{}
	*TrainingJobQueue
}

type trainerService struct {
	mtx                 sync.RWMutex //this lock should only be used for instantiating job queue
	datastore           storage.DataStore
	lcm                 client.LcmClient
	repo                Repository
	jobHistoryRepo      jobHistoryRepository
	modelsBucket        string
	trainedModelsBucket string
	// metrics             *trainerMetrics
	tds           tdsClient.TrainingDataClient
	queues        map[string]*queueHandler
	queuesStarted bool
	service.Lifecycle
}

// NewService creates a new trainer service.
func NewService() Service {
	logr := logger.LogServiceBasic(logger.LogkeyTrainerService)

	config.FatalOnAbsentKey(mongoAddressKey)
	config.SetDefault(gpuLimitsQuerySizeKey, 200)
	config.SetDefault(pollIntervalKey, 60) // in seconds

	ds, err := storage.CreateDataStore(config.GetDataStoreType(), config.GetDataStoreConfig())
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create datastore")
	}
	err = ds.Connect()
	if err != nil {
		logr.WithError(err).Fatalf("Cannot connect to object store")
	}

	repo, err := newTrainingsRepository(viper.GetString(mongoAddressKey),
		viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey),
		viper.GetString(mongoPasswordKey), viper.GetString(config.MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "training_jobs")
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create Repository with %s %s %s %s", viper.GetString(mongoAddressKey),
			viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey), viper.Get(config.MongoAuthenticationDatabase))
	}
	jobHistoryRepo, err := newJobHistoryRepository(viper.GetString(mongoAddressKey),
		viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey),
		viper.GetString(mongoPasswordKey),  viper.GetString(config.MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), collectionNameJobHistory)
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create Repository with %s %s %s %s", viper.GetString("mongo.address"),
			viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey), collectionNameJobHistory)
	}

	queues := make(map[string]*queueHandler)

	s := &trainerService{
		datastore:           ds,
		repo:                repo,
		jobHistoryRepo:      jobHistoryRepo,
		modelsBucket:        getModelsBucket(),
		trainedModelsBucket: getTrainedModelsBucket(),
		queues:              queues,
		queuesStarted:       false,
	}
	logr.Infof("Bucket for model definitions: %s", s.modelsBucket)
	logr.Infof("Bucket for trained models: %s", s.trainedModelsBucket)

	s.RegisterService = func() {
		grpc_trainer_v2.RegisterTrainerServer(s.Server, s)
	}
	//s.StartQueues()
	return s
}

// NewTestService creates a new service instance for testing
func NewTestService(ds storage.DataStore, repo Repository, jobHistoryRepo jobHistoryRepository,
	lcm client.LcmClient, tds tdsClient.TrainingDataClient, queues map[string]*queueHandler) Service {

	config.SetDefault(gpuLimitsQuerySizeKey, 100)
	config.SetDefault(pollIntervalKey, 1) // set poll interval lower to run tests faster

	s := &trainerService{
		datastore:           ds,
		repo:                repo,
		jobHistoryRepo:      jobHistoryRepo,
		lcm:                 lcm,
		modelsBucket:        getModelsBucket(),
		trainedModelsBucket: getTrainedModelsBucket(),
		tds:                 tds,
		queues:              queues,
		queuesStarted:       false,
	}

	s.RegisterService = func() {
		grpc_trainer_v2.RegisterTrainerServer(s.Server, s)
	}
	s.StartQueues()
	return s
}

func debugLogger(logrr *logrus.Entry, isEnabled bool) *logger.LocLoggingEntry {
	logr := new(logger.LocLoggingEntry)
	logr.Logger = logrr
	logr.Enabled = isEnabled

	return logr
}

// Cover for deprecated grpc function.
func grpcErrorDesc(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	return err.Error()
}

func (s *trainerService) StartQueues() {
	logr := logger.LocLogger(logEntry())
	// ensure only one thread per trainer pulling jobs
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if !s.queuesStarted {
		s.queuesStarted = true
		for gpuType, qHandler := range s.queues {
			logr.Debugf("starting queue for %s", gpuType)
			tick := time.NewTicker(time.Duration(viper.GetInt(pollIntervalKey)) * time.Second).C
			go func(gpuType string, qHandler *queueHandler) {
				for {
					select {
					case <-tick:
						s.pullJobFromQueue(gpuType)
					case <-qHandler.stopQueue:
						logr.Debugf("%s queue stopped", gpuType)
						return
					}
				}
			}(gpuType, qHandler)
		}
	}
}

func (s *trainerService) StopTrainer() {
	logr := logger.LocLogger(logEntry())
	logr.Debugf("stopping trainer")

	// Close mongodb connections
	s.repo.Close()
	s.jobHistoryRepo.Close()
	for _, qHandler := range s.queues {
		qHandler.stopQueue <- struct{}{}
		close(qHandler.stopQueue)
		qHandler.TrainingJobQueue.session.Close()
	}
	s.Stop() // stop Service

}

func (s *trainerService) pullJobFromQueue(gpuType string) {
	logr := logger.LocLogger(logEntry())

	qHandler := s.queues[gpuType]
	if qHandler == nil {
		logr.Warnf("there is no queue for type %s", gpuType)
		return
	}

	locked := true
	qerr := qHandler.Lock()
	if qerr != nil {
		logr.WithError(qerr).Errorf("failed to lock %s queue", gpuType)
		return
	}
	defer func() {
		if locked {
			qHandler.Unlock()
		}
	}()

	qSize, err := qHandler.Size()
	logr.Debugf("queue %s has %d elements", gpuType, qSize)
	// s.metrics.queueSizeGauge.With("gpuType", gpuType).Set(float64(qSize))

	empty, err := qHandler.Empty()
	if err != nil {
		logr.WithError(err).Errorf("failed to check if %s queue is empty", gpuType)
		return
	}
	if empty {
		return
	}

	nextJobID, err := qHandler.Peek()
	if err != nil {
		logr.Errorf("failed to peek %s training job queue", gpuType)
		return
	}
	if nextJobID == "" {
		logr.Errorf("job pulled from %s queue is nil", gpuType)
		return
	}

	trainingRecord, err := s.repo.Find(nextJobID)
	logr.Debugf("Training ID is %v", nextJobID)
	if err != nil {
		if err == mgo.ErrNotFound {
			logr.Debugf("job %s not found in mongo, assuming job was deleted", nextJobID)
			qHandler.Delete(nextJobID)
			// s.metrics.deleteJobFromQueueCounter.Add(1)
			return
		}
		logr.WithError(err).Errorf("error retrieving training job")
		return
	}
	if trainingRecord.TrainingStatus == nil {
		logr.Debugf("PullJobFromQueue: Training Status is nil")
	}
	if trainingRecord.Deleted {
		logr.Debugf("job %s was deleted", nextJobID)
		qHandler.Delete(nextJobID)
		// s.metrics.deleteJobFromQueueCounter.Add(1)
		return
	}
	if trainingRecord.TrainingStatus.Status != grpc_trainer_v2.Status_QUEUED {
		logr.Warnf("job %s expected status QUEUED but found %s, removing job from queue", nextJobID, trainingRecord.TrainingStatus)
		qHandler.Delete(nextJobID)
		// s.metrics.deleteJobFromQueueCounter.Add(1)
		return
	}

	logr.Debugf("got training job %s from %s queue", nextJobID, gpuType)

	rateLimited := s.rateLimitTrainingJob(trainingRecord, logr)
	if rateLimited {
		logr.Debugf("training job %s is rate-limited, leaving in %s queue", trainingRecord.TrainingID, gpuType)
		return
	}

	err = s.submitJobToLCM(trainingRecord, logr)
	if err != nil {
		// submitting to LCM failed, don't update job history or dequeue
		logr.Error(" submit job to LCM error:" + err.Error())
		return
	}

	dequeuedJob, dequeueErr := qHandler.Dequeue()
	if dequeueErr != nil {
		logr.WithError(dequeueErr).Errorf("Failed to dequeue training job %s", trainingRecord.TrainingID)
	}
	if dequeueErr == nil && dequeuedJob != trainingRecord.TrainingID {
		logr.Errorf("expected to dequeue job %s, but got %s instead. This should never happen", trainingRecord.TrainingID, dequeuedJob)
		enqueueErr := qHandler.Enqueue(dequeuedJob)
		if enqueueErr != nil {
			logr.Errorf("job %s should not have been dequeued, and could not be re-enqueued. the record will stay in mongo but the job will never run", dequeuedJob)
			// find and update record with FAILED status
			if dequeuedTrainingRecord, err := s.repo.Find(dequeuedJob); err != nil {
				// this is only a problem if the dequeued job is still QUEUED, since it will stay QUEUED forever. if it has already been submitted to LCM, it should not be in the queue
				if dequeuedTrainingRecord.TrainingStatus.Status == grpc_trainer_v2.Status_QUEUED {
					_, err = updateTrainingJobPostLock(s, &grpc_trainer_v2.UpdateRequest{
						TrainingId:    dequeuedJob,
						UserId:        dequeuedTrainingRecord.UserID,
						Status:        grpc_trainer_v2.Status_FAILED,
						StatusMessage: "Job was dequeued without being submitted",
						ErrorCode:     trainerClient.ErrCodeFailDequeue,
					})
					if err != nil {
						logr.WithError(err).Errorln("Unable to update job status to FAILED")
					}
				}
			}
		}
	}
	// s.metrics.dequeueJobCounter.Add(1)

	qHandler.Unlock()
	locked = false

	// store job state transition
	timestamp := trainerClient.CurrentTimestampAsString()
	e := &JobHistoryEntry{
		TrainingID:    trainingRecord.TrainingID,
		Timestamp:     timestamp,
		Status:        trainingRecord.TrainingStatus.Status,
		StatusMessage: trainingRecord.TrainingStatus.StatusMessage,
		ErrorCode:     trainingRecord.TrainingStatus.ErrorCode,
	}
	s.jobHistoryRepo.RecordJobStatus(e)
}

func uploadArchive(bucket string, file string, bytes []byte, onlyMinio bool) error {
	getLogger := logger.GetLogger()
	sClient, err := storageClient.NewStorage()
	defer sClient.Close()
	if err != nil {
		getLogger.Errorf("init storage client failed, %v", err.Error())
		return err
	}
	todo := context.TODO()
	ctx, cancel := context.WithCancel(todo)
	defer cancel()
	uploadClient, err := sClient.Client().Upload(ctx)
	if err != nil {
		getLogger.Errorf("storageClient.Client().Upload failed, %v", err.Error())
		return err
	}
	//先发送fileInfo
	request := grpc_storage.UploadRequest{Data: &grpc_storage.UploadRequest_Info{Info: &grpc_storage.FileInfo{Bucket: bucket, FileName: file, OnlyMinio: onlyMinio}}}
	err = uploadClient.Send(&request)
	if err != nil {
		getLogger.Errorf(" Upload UploadRequest_Info failed, %v", err.Error())
		return err
	}

	//开始发送文件
	if err == io.EOF {
		getLogger.Debugf("read model done, %v", err.Error())
		return nil
	}
	if err != nil {
		getLogger.Errorf("file.Read(buffer) failed, %v", err.Error())
		return err
	}
	request = grpc_storage.UploadRequest{Data: &grpc_storage.UploadRequest_ChucksData{ChucksData: bytes}}
	err = uploadClient.Send(&request)
	if err != nil {
		getLogger.Errorf("uploadClient.Send(&request) failed, %v", err.Error())
		return err
	}
	response, err := uploadClient.CloseAndRecv()
	if err != nil {
		getLogger.Errorf("CloseAndRecv failed, %v", err.Error())
		return err
	}
	getLogger.Debugf("upload response: %+v", response)
	return nil
}

func (s *trainerService) CreateTrainingJob(ctx context.Context, req *grpc_trainer_v2.CreateRequest) (*grpc_trainer_v2.CreateResponse, error) {
	// duration := metrics.NewTimer(s.metrics.createTrainingDuration)
	// defer duration.ObserveDuration()

	//sClient, err := storageClient.NewStorage()
	getLogger := logger.GetLogger()
	startTime := time.Now()
	getLogger.Debugf("CreateTrainingJob startTime: %v", startTime.String())

	// FIXME MLSS Change: The job name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character
	sidHigh := random.String(9, random.Lowercase+random.Numeric+"-")
	sidLow := random.String(9, random.Lowercase+random.Numeric)

	id := fmt.Sprintf("training-%s%s", sidHigh, sidLow)

	logr := logger.LocLogger(logWith(id, req.UserId))

	logr.Infof("CreateTrainingJob, req.TFosRequest: %v", req.TFosRequest)

	cl := instrumentation.NewCallLogger(ctx, "CreateTrainingJob", logr)
	defer cl.Returned()

	if err := s.validateRequest(logr.Logger, req); err != nil {
		return nil, err
	}

	setDefaultResourceRequirements(req.Training)

	//request is validated, now bump up the counter
	logFrameworkVersionValue := fmt.Sprintf("%s-%s", req.ModelDefinition.Framework.Name, req.ModelDefinition.Framework.Version)
	logGpuTypeUsagesValue := fmt.Sprintf("%s-%v", req.Training.Resources.GpuType, req.Training.Resources.Gpus)
	logr = logr.WithFields(logrus.Fields{
		logger.LogkeyFramework:        req.ModelDefinition.Framework.Name,
		logger.LogkeyFrameworkVersion: logFrameworkVersionValue,
		logger.LogkeyGpuType:          req.Training.Resources.GpuType,
		logger.LogkeyGpuUsage:         logGpuTypeUsagesValue,
	})

	logr.Debug(" metrics for total number of training jobs ")

	//TODO get ds
	outputDatastore := s.getOutputDatastore(req.Training.OutputData, req.Datastores)

	// upload model definition ZIP file to object store and set location
	logr.Debugf("debug_1.5.1 start to s3")
	if req.ModelDefinition.Content != nil {
		logr.Debugf("debug_1.5.1 start to s3 is nil")
		// MLSS Change: upload model to object store no matter what output datastore is used
		//if outputDatastore.Type != "mount_volume" {
		// Upload to DLaaS Object store.
		logr.Debugf("Uploading content to: %s filename: %s", s.modelsBucket, getModelZipFileName(id))
		//err := s.datastore.UploadArchive(s.modelsBucket, getModelZipFileName(id), req.ModelDefinition.Content)
		//todo mlss 1.10.0
		err := uploadArchive(s.modelsBucket, getModelZipFileName(id), req.ModelDefinition.Content, true)

		if err != nil {
			logr.WithError(err).Errorf("Error uploading model to object store")
			// s.metrics.uploadModelFailedCounter.With("kind", dlaasStoreKind).Add(1)
			return nil, err
		}
		req.ModelDefinition.Location = fmt.Sprintf("%s/%s.zip", s.modelsBucket, id)
		cl.Observe("uploaded model to minio object store: %s", req.ModelDefinition.Location)

		// Upload to user's result object store.
		logr.Debugf("Upload to user's result object store, type: %s bucket: %s", outputDatastore.Type, outputDatastore.Fields["bucket"])

		bucket := outputDatastore.Fields["bucket"]
		object := fmt.Sprintf("%s/_submitted_code/model.zip", id)
		logr.Debugf("Writing to output object store: %s -> %s, length: %d", bucket, object, len(req.ModelDefinition.Content))
		//err = ds.UploadArchive(bucket, object, req.ModelDefinition.Content)
		//todo
		err = uploadArchive(bucket, object, req.ModelDefinition.Content, true)
		if err != nil {
			// s.metrics.uploadModelFailedCounter.With("kind", userStoreKind).Add(1)
			logr.WithError(err).Errorf("Error uploading model to output object store")
			return nil, err
		}

		cl.Observe("uploaded model to user's object store")
		//}
	}
	elapsed := time.Since(startTime)
	getLogger.Debugf("uploadArchive endTime: %v", elapsed.String())

	// create a copy of the model definition without the content field (do not store it to the database)
	modelWithoutContent := *req.ModelDefinition
	modelWithoutContent.Content = nil

	trainStatus := grpc_trainer_v2.TrainingStatus{
		Status:              grpc_trainer_v2.Status_QUEUED,
		SubmissionTimestamp: trainerClient.CurrentTimestampAsString(),
	}
	tr := &TrainingRecord{
		TrainingID:            id,
		UserID:                req.UserId,
		ModelDefinition:       &modelWithoutContent,
		Training:              req.Training,
		Datastores:            req.Datastores,
		TrainingStatus:        &trainStatus,
		Metrics:               nil,
		Namespace:             req.Namespace,
		GID:                   req.Gid,
		UID:                   req.Uid,
		GuardianToken:         req.Token,
		JobAlert:              req.JobAlert,
		CodeSelector:          req.CodeSelector,
		PSs:                   req.PSs,
		PSCPU:                 req.PSCpu,
		PSImage:               req.PSImage,
		PSMemory:              req.PSMemory,
		DataPath:              req.DataPath,
		JobType:               req.JobType,
		TFosRequest:           req.TFosRequest,
		ExpRunId:              req.ExpRunId,
		ExpName:               req.ExpName,
		FileName:              req.FileName,
		FilePath:              req.FilePath,
		ProxyUser:             req.ProxyUser,
		DataSet: 			   req.DataSet,
		MFModel: 			   req.MfModel,
		Algorithm:             req.Algorithm,
		JobParams:             req.JobParams,
		FitParams: 			   req.FitParams,
		APIType: 			   req.APIType,
	}

	if req.ExpRunId == "0" {
		tr.ExpRunId = ""
	}
	logr.Debugf("CreateTrainingJob, req.JobType: %v, reqCpu: %v, reqGpu: %v, reqMem: %v, fileName: %v, filePath: %v", req.JobType, req.Training.Resources.Cpus, req.Training.Resources.Gpus, req.Training.Resources.Memory, req.FileName, req.FilePath)

	elapsed = time.Since(startTime)
	getLogger.Debugf("submitJobToLCM startTime: %v", elapsed.String())
	getLogger.Debugf("submitJobToLCM STATUS: %v", tr.TrainingStatus)
	err := s.submitJobToLCM(tr, logr)

	elapsed = time.Since(startTime)
	getLogger.Debugf("submitJobToLCM endTime: %v", elapsed.String())

	if err != nil {
		// err logged in submitJobToLCM
		return nil, err
	}
	cl.Observe("submitted job to lcm")
	//}

	//request is validated, now bump up the counter
	logUserGpuValue := fmt.Sprintf("%s-%s-%v", req.UserId, req.Training.Resources.GpuType, req.Training.Resources.Gpus)

	logUserFrameworkVersionGpuValue := fmt.Sprintf("%s-%s-%s-%s-%v", req.UserId, req.Training.Resources.GpuType, req.ModelDefinition.Framework.Name, req.ModelDefinition.Framework.Version, req.Training.Resources.Gpus)

	logr.WithFields(logrus.Fields{
		"userid_gputype_gpus":           logUserGpuValue,
		"userid_framework_gputype_gpus": logUserFrameworkVersionGpuValue,
	}).Debug(" incrementing userid log")

	return &grpc_trainer_v2.CreateResponse{TrainingId: id}, nil
}

func (s *trainerService) GetTrainingJob(ctx context.Context, req *grpc_trainer_v2.GetRequest) (*grpc_trainer_v2.GetResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	cl := instrumentation.NewCallLogger(ctx, "GetTrainingJob", logr)
	defer cl.Returned()

	tr, err := s.repo.Find(req.TrainingId)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
		}
		logr.WithError(err).Errorf("Cannot retrieve training record")
		return nil, err
	}

	cl.Observe("got training job record")
	jobb := &grpc_trainer_v2.Job{
		UserId:          tr.UserID,
		JobId:           tr.JobID,
		ModelDefinition: tr.ModelDefinition,
		TrainingId:      tr.TrainingID,
		Training:        tr.Training,
		Status:          tr.TrainingStatus,
		Datastores:      tr.Datastores,
		Metrics:         tr.Metrics,
		JobNamespace:    tr.Namespace,
	}
	return &grpc_trainer_v2.GetResponse{
		Job: jobb,
	}, nil
}

func (s *trainerService) GetTrainingStatusID(ctx context.Context, req *grpc_trainer_v2.GetRequest) (*grpc_trainer_v2.GetStatusIDResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	statusID, err := s.repo.FindTrainingStatusID(req.TrainingId)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
		}
		logr.WithError(err).Errorf("Cannot retrieve record for training %s", req.TrainingId)
		return nil, err
	}
	return &grpc_trainer_v2.GetStatusIDResponse{
		Status: statusID,
	}, nil
}

func (s *trainerService) UpdateTrainingJob(ctx context.Context, req *grpc_trainer_v2.UpdateRequest) (*grpc_trainer_v2.UpdateResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("UpdateTrainingJob called for training %s", req.TrainingId)

	return updateTrainingJobPostLock(s, req)
}

// This method contains all the functionality of UpdateTrainingJob, minus the lock on the database.  This enables it to be called
// from within another function, which already has the lock itself (Halt)
func updateTrainingJobPostLock(s *trainerService, req *grpc_trainer_v2.UpdateRequest) (*grpc_trainer_v2.UpdateResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	training, err := s.repo.Find(req.TrainingId)
	if err != nil {
		logr.WithError(err).Errorf("Cannot retrieve training '%s'", req.TrainingId)
		return nil, err
	}
	if training == nil {
		// training does not exist
		return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
	}

	if training.UserID != req.UserId {
		msg := fmt.Sprintf("User %s does not have permission to update training data with id %s.", req.UserId, req.TrainingId)
		logr.Error(msg)
		return nil, gerrf(codes.PermissionDenied, msg)
	}

	ts := training.TrainingStatus
	originalStatus := ts.Status

	// If status is completed/failed/halted and the update is requesting a halt, then do nothing and return error
	if originalStatus == grpc_trainer_v2.Status_COMPLETED || originalStatus == grpc_trainer_v2.Status_FAILED {
		return nil, err
	}

	ts.Status = req.Status
	ts.StatusMessage = req.StatusMessage
	ts.ErrorCode = req.ErrorCode

	nowMillis := trainerClient.CurrentTimestampAsString()

	if req.Status == grpc_trainer_v2.Status_COMPLETED || req.Status == grpc_trainer_v2.Status_FAILED {
		// FIXME MLSS Change: keep datastores even job finished
		//training.Datastores = nil
		ts.CompletionTimestamp = nowMillis
		if req.Timestamp != "" {
			ts.CompletionTimestamp = req.Timestamp
		}
		// erase sensitive data from the db
		training.ModelDefinition.Framework.ImageLocation = nil
	}

	if req.Status == grpc_trainer_v2.Status_PROCESSING {
		ts.ProcessStartTimestamp = nowMillis
		if req.Timestamp != "" {
			ts.ProcessStartTimestamp = req.Timestamp
		}
	}

	// send monitoring metrics for failed/succeeded jobs
	if req.Status == grpc_trainer_v2.Status_COMPLETED || req.Status == grpc_trainer_v2.Status_FAILED {
		// counter := s.metrics.trainingJobSucceededCounter
		if req.Status == grpc_trainer_v2.Status_FAILED {
			errorType := "server"
			if strings.HasPrefix(req.ErrorCode, "C") {
				errorType = "client"
			}

			logFrameworkErrorsValue := fmt.Sprintf("%s-%s-%s", training.ModelDefinition.Framework.Name, errorType, req.ErrorCode)
			logr.WithFields(logrus.Fields{
				"framework_errors": logFrameworkErrorsValue,
			}).Debug(" metrics for failed training jobs framework")

			// counter = s.metrics.trainingJobFailedCounter.With("type", errorType, "errorcode", req.ErrorCode)
		}
		gpusUsed := training.Training.Resources.Gpus
		if training.Training.Resources.Learners > 1 {
			gpusUsed = training.Training.Resources.Gpus * float64(training.Training.Resources.Learners)
		}

		logGpuTypeDecrementValue := fmt.Sprintf("%s-%v", training.Training.GetResources().GpuType, gpusUsed)
		logr.WithFields(logrus.Fields{
			"gputype_decrement": logGpuTypeDecrementValue,
		}).Debug(" decrementing the gpus")

		logUserGpuValue := fmt.Sprintf("%s-%s-%v", req.UserId, training.Training.GetResources().GpuType, gpusUsed)

		logUserFrameworkVersionGpuValue := fmt.Sprintf("%s-%s-%s-%s-%v", req.UserId, training.Training.GetResources().GpuType, training.ModelDefinition.Framework.Name, training.ModelDefinition.Framework.Version, gpusUsed)

		logr.WithFields(logrus.Fields{
			"userid_gputype_gpus":           logUserGpuValue,
			"userid_framework_gputype_gpus": logUserFrameworkVersionGpuValue,
		}).Debug(" decrementing userid log")

	}
	err = s.repo.Store(training)
	if err != nil {
		logr.WithError(err).Errorf("Failed updating status of training %s in DB", req.TrainingId)
		return nil, err
	}

	// verify that the training job details have been updated properly
	training, err = s.repo.Find(req.TrainingId)
	if err != nil {
		logr.WithError(err).Errorf("Cannot retrieve training '%s'", req.TrainingId)
		return nil, err
	}
	if training == nil {
		// training does not exist
		return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
	}
	ts = training.TrainingStatus
	logr.Debugf("CHECKING Stored training %s, Status %s Error Code %s Message %s", req.TrainingId, ts.Status, ts.ErrorCode, ts.StatusMessage)

	// Additionally, store any job state transitions in the job_history DB collection
	// We store a history record if either (1) the status is different, or (2) if this is
	// a PROCESSING->PROCESSING transition, to record the full picture for distributed jobs.
	if req.Status != originalStatus || req.Status == grpc_trainer_v2.Status_PROCESSING {
		timestamp := req.Timestamp
		if req.Timestamp == "" {
			// If timestamp is missing, we may end up storing duplicate events (with different timestamps)
			// in the job_history DB collection. Technically, that shouldn't happen (as we always add a
			// timestamp in controller/jobmonitor when calling UpdateTrainingJob(..))
			logr.Warnf("Timestamp missing in UpdateTrainingJob(..) request, adding current time.")
			timestamp = nowMillis
		}
		e := &JobHistoryEntry{
			TrainingID:    req.TrainingId,
			Timestamp:     timestamp,
			Status:        req.Status,
			StatusMessage: req.StatusMessage,
			ErrorCode:     req.ErrorCode,
		}
		s.jobHistoryRepo.RecordJobStatus(e)
	}

	return &grpc_trainer_v2.UpdateResponse{TrainingId: training.TrainingID}, nil
}

func (s *trainerService) GetAllTrainingsJobs(ctx context.Context, req *grpc_trainer_v2.GetAllRequest) (*grpc_trainer_v2.GetAllResponse, error) {
	logr := logger.LocLogger(logEntry().WithField(logger.LogkeyUserID, req.UserId))
	logr.Debugf("GetAllTrainingsJobs called")

	cl := instrumentation.NewCallLogger(ctx, "GetAllTrainingsJobs", logr)
	defer cl.Returned()

	jobPaged, err := s.repo.FindAll(req.UserId, req.Page, req.Size)
	if err != nil {
		msg := "Failed to retrieve all training jobPaged"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_trainer_v2.GetAllResponse{
		Jobs: make([]*grpc_trainer_v2.Job, len(jobPaged.TrainingRecords)),
	}
	for i, job := range jobPaged.TrainingRecords {
		resp.Jobs[i] = &grpc_trainer_v2.Job{
			UserId:          job.UserID,
			JobId:           job.JobID,
			ModelDefinition: job.ModelDefinition,
			TrainingId:      job.TrainingID,
			Training:        job.Training,
			Status:          job.TrainingStatus,
			Datastores:      job.Datastores,
			TfosRequest:     job.TFosRequest,
		}
	}

	pageStr := strconv.Itoa(jobPaged.Pages)
	resp.Pages = pageStr
	totalStr := strconv.Itoa(jobPaged.Total)
	resp.Total = totalStr

	return resp, nil
}

func (s *trainerService) GetAllTrainingsJobsByUserIdAndNamespace(ctx context.Context, req *grpc_trainer_v2.GetAllRequest) (*grpc_trainer_v2.GetAllResponse, error) {
	logr := logger.LocLogger(logEntry().WithField(logger.LogkeyUserID, req.UserId))
	logr.Debugf("GetAllTrainingsJobsByUserIdAndNamespace called")

	cl := instrumentation.NewCallLogger(ctx, "GetAllTrainingsJobsByUserIdAndNamespace", logr)
	defer cl.Returned()

	// FIXME MLSS Change: get models filter by username and namespace
	//userId := req.UserId
	username := req.Username
	namespace := req.Namespace
	page := req.Page
	size := req.Size
	logr.Debugf("trainerService GetAllTrainingsJobsByUserIdAndNamespace, Username: %v, Namespace: %v, Page: %v, Size: %v", username, namespace, page, size)

	//jobs, err := s.repo.FindAll(req.UserId)
	jobs, err := s.repo.FindAllByUserIdAndNamespace(&username, &namespace, page, size)
	if err != nil {
		msg := "Failed to retrieve all training jobs"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_trainer_v2.GetAllResponse{
		Jobs: make([]*grpc_trainer_v2.Job, len(jobs.TrainingRecords)),
	}
	for i, job := range jobs.TrainingRecords {
		logr.Debugf("trainer_impl GetAllTrainingsJobsByUserIdAndNamespace job.JobAlert: %v", job.JobAlert)
		resp.Jobs[i] = &grpc_trainer_v2.Job{
			UserId:          job.UserID,
			JobId:           job.JobID,
			ModelDefinition: job.ModelDefinition,
			TrainingId:      job.TrainingID,
			Training:        job.Training,
			Status:          job.TrainingStatus,
			Datastores:      job.Datastores,
			// FIXME MLSS Change: more properties for result
			JobNamespace: job.Namespace,
			JobAlert:     job.JobAlert,
			Pss:          job.PSs,
			PsCpu:        job.PSCPU,
			PsImage:      job.PSImage,
			PsMemory:     job.PSMemory,
			ExpRunId:     job.ExpRunId,
			ExpName:      job.ExpName,
			JobType:      job.JobType,
			TfosRequest:  job.TFosRequest,
		}
	}

	pageStr := strconv.Itoa(jobs.Pages)
	resp.Pages = pageStr
	totalStr := strconv.Itoa(jobs.Total)
	resp.Total = totalStr
	logr.Debugf("jobs.Pages: %v,jobs.Total: %v, debug_jobs: %v", jobs.Pages, jobs.Total, jobs)
	return resp, nil
}

// FIXME MLSS Change: get models filter by username and namespace
func (s *trainerService) GetAllTrainingsJobsByUserIdAndNamespaceList(ctx context.Context, req *grpc_trainer_v2.GetAllRequest) (*grpc_trainer_v2.GetAllResponse, error) {
	logr := logger.LocLogger(logEntry().WithField(logger.LogkeyUserID, req.UserId))
	logr.Debugf("GetAllTrainingsJobsByUserIdAndNamespace called")

	cl := instrumentation.NewCallLogger(ctx, "GetAllTrainingsJobsByUserIdAndNamespace", logr)
	defer cl.Returned()

	// FIXME MLSS Change: get models filter by username and namespace
	//userId := req.UserId
	username := req.Username
	namespaceList := req.NamespaceList
	page := req.Page
	size := req.Size
	clusterName := req.ClusterName
	logr.Debugf("trainerService GetAllTrainingsJobsByUserIdAndNamespaceList, Username: %v, NamespaceList: %v, Page: %v, Size: %v", username, namespaceList, page, size)

	//jobs, err := s.repo.FindAll(req.UserId)
	jobs, err := s.repo.FindAllByUserIdAndNamespaceList(&username, &namespaceList, page, size, clusterName)
	if err != nil {
		msg := "Failed to retrieve all training jobs"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_trainer_v2.GetAllResponse{
		Jobs: make([]*grpc_trainer_v2.Job, len(jobs.TrainingRecords)),
	}
	for i, job := range jobs.TrainingRecords {
		logr.Debugf("debug_for_param ParameterServer trainer pss: %v, psCpu: %v, psImage: %v, psMemory: %v", job.PSs, job.PSCPU, job.PSImage, job.PSMemory)
		resp.Jobs[i] = &grpc_trainer_v2.Job{
			UserId:          job.UserID,
			JobId:           job.JobID,
			ModelDefinition: job.ModelDefinition,
			TrainingId:      job.TrainingID,
			Training:        job.Training,
			Status:          job.TrainingStatus,
			Datastores:      job.Datastores,
			// FIXME MLSS Change: more properties for result
			JobNamespace: job.Namespace,
			JobAlert:     job.JobAlert,
			Pss:          job.PSs,
			PsCpu:        job.PSCPU,
			PsImage:      job.PSImage,
			PsMemory:     job.PSMemory,
			JobType:      job.JobType,
			ExpRunId:     job.ExpRunId,
			ExpName:      job.ExpName,
			TfosRequest:  job.TFosRequest,
		}
	}

	pageStr := strconv.Itoa(jobs.Pages)
	resp.Pages = pageStr
	totalStr := strconv.Itoa(jobs.Total)
	resp.Total = totalStr

	return resp, nil
}

// cover for depreciated grpc method
func grpcCode(err error) codes.Code {
	if s, ok := status.FromError(err); ok {
		return s.Code()
	}
	return codes.Unknown
}

func (s *trainerService) deleteJobFromTDS(query *tdsService.Query, logr *logger.LocLoggingEntry) error {
	tds, err := s.tdsClient()
	if err != nil {
		logr.WithError(err).Error("Cannot create TDS client")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*4)
	defer cancel()

	delResponse, err := tds.Client().DeleteJob(ctx, query)
	if err != nil {
		logr.WithError(err).Error("tds DeleteJob returned error")
		return err
	}
	if !delResponse.Success {
		logr.Warn("tds DeleteJob reported false for success")
	}
	return nil
}

func (s *trainerService) deleteJobFromQueue(trainingID string, namespaceName string, logr *logger.LocLoggingEntry) error {
	qHandler := s.queues[QueueName(namespaceName)]
	if qHandler == nil {
		// FIXME MLSS Change: v_1.4.1 to create queue  when queue is nil
		queue, err := newTrainingJobQueue(viper.GetString(mongoAddressKey),
			viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey),
			viper.GetString(mongoPasswordKey),  viper.GetString(config.MongoAuthenticationDatabase),
			config.GetMongoCertLocation(), QueueName(namespaceName), LockName(namespaceName))
		if err != nil {
			logr.WithError(err).Fatalf("Cannot create queue with %s %s %s", viper.GetString(mongoAddressKey), viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey))
			// s.metrics.trainerServiceRestartCounter.With("reason", "createqueue").Add(1)
		}
		s.queues[QueueName(namespaceName)] = &queueHandler{make(chan struct{}), queue}
		qHandler = s.queues[QueueName(namespaceName)]
	}

	qerr := qHandler.Lock()
	if qerr != nil {
		logr.WithError(qerr).Errorf("failed to lock %s queue", QueueName(namespaceName))
		return qerr
	}
	defer func() {
		qHandler.Unlock()
	}()

	deleted, err := qHandler.Delete(trainingID)
	if err != nil {
		logr.WithError(err).Errorf("failed to delete job %s from queue %s", trainingID, QueueName(namespaceName))
		return err
	}
	if !deleted {
		logr.Debugf("job %s not found in queue %s", trainingID, QueueName(namespaceName))
	} else {
		logr.Debugf("job %s deleted from queue %s", trainingID, QueueName(namespaceName))
		// s.metrics.deleteJobFromQueueCounter.Add(1)
	}
	return nil
}

func (s *trainerService) DeleteTrainingJob(ctx context.Context,
	req *grpc_trainer_v2.DeleteRequest) (*grpc_trainer_v2.DeleteResponse, error) {

	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	cl := instrumentation.NewCallLogger(ctx, "DeleteTrainingJob", logr)
	defer cl.Returned()

	// s.metrics.deleteTrainingJobCounter.Add(1)

	readResp, err := s.GetTrainingJob(ctx, &grpc_trainer_v2.GetRequest{
		TrainingId: req.TrainingId,
		UserId:     req.UserId,
	})

	if readResp == nil {
		logr.Debugf("readResp from trainer_impl in DeleteTrainingJob is nil!!!")
	}

	if err != nil {
		logr.WithError(err).Errorf("Failing querying training job")
		return nil, err
	}

	cl.Observe("got training job record")

	// We've noticed that deleting from the TDS can take several minutes, and we don't want to delay this
	// call due to that. This is a temporary workaround until we find out root cause of the TDS slowdowns.
	go func() {
		err = s.deleteJobFromTDS(&tdsService.Query{
			Meta: &tdsService.MetaInfo{
				TrainingId: req.TrainingId,
				UserId:     req.UserId,
			},
		}, logr)
		if err != nil {
			logr.WithError(err).Warn("deleteJobFromTDS returned error")
		}

		cl.Observe("cleaned up job in TDS")
	}()

	var job *grpc_trainer_v2.Job
	if readResp != nil {
		job = readResp.Job

		// delete from queue
		if job.Status.Status == grpc_trainer_v2.Status_QUEUED {
			// if this fails, the queue entry will be cleaned up when the job is pulled
			s.deleteJobFromQueue(job.TrainingId, job.JobNamespace, logr)
		}

		// Do the LCM cleanup in the background. We noticed this step can take a long time and cause context deadline
		// exceeded errors where there were many concurrent calls to delete a job.
		// As long as we can delete the record from mongo, object store, and the training data service, the user gets a
		// successful status back.
		// LCM failures are silently ignored, so we need alerts when the LCM cleanup fails, and be more proactive
		// in cleaning up stale learners.
		go func() {
			// delete the job if it exists
			lcm, err := s.lcmClient()
			if err != nil {
				logr.WithError(err).Errorln("Cannot create lcm service client")
				return
			}
			defer lcm.Close()

			lcmCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			_, err = lcm.Client().KillTrainingJob(lcmCtx, &service.JobKillRequest{
				Name:         job.JobId,
				TrainingId:   job.TrainingId,
				UserId:       job.UserId,
				JobNamespace: job.JobNamespace,
			})

			// tolerate "not found" because it just means the job is no longer running
			if err != nil && grpcCode(err) != codes.NotFound {
				logr.WithError(err).Errorf("Failed to kill job '%s'", job.JobId)
				return
			}
			logr.Debugf("Kubernetes job '%s' does not longer exist.", job.JobId)

			cl.Observe("killed job in LCM")
		}()

		// delete model content from data store

		// FIXME MLSS Change: delete model.zip from s3
		logr.Debugf("delete model zip from trainer_impl in DeleteTrainingJob, modelsBucket: %v, modelZipFileName: %v", s.modelsBucket, getModelZipFileName(job.JobId))
		err = s.datastore.DeleteArchive(s.modelsBucket, getModelZipFileName(job.JobId))
		if err != nil {
			logr.WithError(err).Errorf("Error deleting model from object store")
			// log this error, but continue with deleting the training record anyway
		}
		cl.Observe("deleted model from object store")

		// FIXME MLSS Change: delete tfos-job in linkis first
		err = killTfos(job, s, logr)
		if err != nil {
			logr.WithError(err).Errorf("kill Tfos in linkis failed: %v", err.Error())
			return nil, err
		}

		// delete from DB
		err = s.repo.Delete(job.TrainingId)
		if err != nil {
			logr.WithError(err).Errorf("Failed to delete training job '%s' from database", job.TrainingId)
			return nil, err
		}
		cl.Observe("deleted model from mongo")

		return &grpc_trainer_v2.DeleteResponse{TrainingId: job.JobId}, nil
	}
	return nil, gerrf(codes.NotFound, "Training with id '%s' not found.", req.TrainingId)
}

func killTfos(job *grpc_trainer_v2.Job, s *trainerService, logr *logger.LocLoggingEntry) error {
	training, err := s.repo.Find(job.TrainingId)
	if err != nil {
		logr.WithError(err).Errorf("Cannot retrieve training '%s'", job.TrainingId)
		return err
	}
	if training == nil {
		// training does not exist
		logr.WithError(err).Errorf("Training with id %s not found.", job.TrainingId)
		return err
	}

	logr.Infof(" killTfos,req.UserId:  %s, req.TrainingId:  %s, training.LinkisExecId:  %s", job.UserId, job.TrainingId, training.LinkisExecId)

	// if training.LinkisExecId != "" {
	// 	var linkisAddress = os.Getenv("LINKIS_ADDRESS")
	// 	var linkisTokenCode = os.Getenv("LINKIS_TOKEN_CODE")
	// 	logr.Infof(" killTfos,LINKIS_ADDRESS: %s, LINKIS_TOKEN_CODE: %v", linkisAddress, linkisTokenCode)

	// 	var requestBase = models.LinkisRequestBase{
	// 		LinkisAddress:   linkisAddress,
	// 		LinkisTokenCode: linkisTokenCode,
	// 		UserId:          job.UserId,
	// 	}
	// 	err = linkisclientutils.KillingJob(requestBase, training.LinkisExecId)
	// 	if err != nil {
	// 		logr.WithError(err).Errorf("Cannot kill training '%s': %v", job.TrainingId, err.Error())
	// 		return err
	// 	}
	// } else {
	// 	logr.Infof(" killTfos,training.LinkisExecId == '',no need to delete job in linkis ")
	// }

	return nil
}

func (s *trainerService) ResumeTrainingJob(ctx context.Context, req *grpc_trainer_v2.ResumeRequest) (*grpc_trainer_v2.ResumeResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("HaltTrainingJob called")
	return nil, gerrf(codes.Unimplemented, "ResumeTrainingJob not implemented yet")
}

func (s *trainerService) GetModelDefinition(req *grpc_trainer_v2.ModelDefinitionRequest, stream grpc_trainer_v2.Trainer_GetModelDefinitionServer) error {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("GetModelDefinition")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.GetTrainingJob(ctx, &grpc_trainer_v2.GetRequest{
		TrainingId: req.TrainingId,
		UserId:     req.UserId,
	})
	if err != nil {
		logr.WithError(err).Errorf("Failed to read training with id: %s", req.TrainingId)
		return gerrf(codes.Internal, "Failed to read training with id: %s", req.TrainingId)
	}
	if resp == nil || resp.Job == nil {
		return gerrf(codes.NotFound, "Training with id '%s' not found.", req.TrainingId)
	}

	// TODO we need to change this to accept a writer to be more efficient
	payload, err := s.datastore.DownloadArchive(s.modelsBucket, getModelZipFileName(req.TrainingId))
	if err != nil {
		logr.WithError(err).Errorf("Downloading model definition archive failed")
	}
	err = stream.Send(&grpc_trainer_v2.ZippedDataChunk{
		Data: payload,
	})
	if err != nil {
		logr.WithError(err).Errorf("Failed to send zipped chunk.")
		return err
	}
	return nil
}

func (s *trainerService) GetTrainedModel(req *grpc_trainer_v2.TrainedModelRequest, stream grpc_trainer_v2.Trainer_GetTrainedModelServer) error {
	//s.mtx.Lock()
	//defer s.mtx.Unlock()

	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("GetTrainedModel")

	// s.metrics.downloadTrainedModelJobCounter.Add(1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.GetTrainingJob(ctx, &grpc_trainer_v2.GetRequest{
		TrainingId: req.TrainingId,
		UserId:     req.UserId,
	})
	if err != nil {
		logr.WithError(err).Errorf("Error reading training with id: %s", req.TrainingId)
		return err
	}
	if resp == nil || resp.Job == nil {
		return gerrf(codes.NotFound, "Training with id '%s' not found.", req.TrainingId)
	}

	var ostore storage.DataStore
	ds := s.getOutputDatastore(resp.Job.Training.OutputData, resp.Job.Datastores)
	ostore, err = storage.CreateDataStore(ds.Type, ds.Connection)
	if err != nil {
		logr.WithError(err).Errorf("Error creating datastore: %v", ds)
		return err
	}
	if err := ostore.Connect(); err != nil {
		logr.WithError(err).Error("Error connect to datastore")
		return err
	}
	defer ostore.Disconnect()

	trainedModelSize, err := ostore.GetTrainedModelSize(fmt.Sprintf("%s/%s", ds.Fields["bucket"], resp.Job.TrainingId),
		resp.Job.Training.Resources.Learners)

	if err != nil {
		logr.WithError(err).Error("Error retrieving trained model size")
		return err
	}
	logr.Debugf("The size of the trained model is %d", trainedModelSize)

	// DP only allows downloads of sizes less than 200MBs
	if trainedModelSize > 200000000 {
		logr.Debugf("Trained model for '%s' exceeded download limit size.", req.TrainingId)
		return gerrf(codes.FailedPrecondition,
			"Trained model exceeded download limit. Download from your cloud storage directly")
	}

	r, w := io.Pipe() // connect I/O without temp space.

	go func() {
		// write to pipe by downloading
		err := ostore.DownloadTrainedModelAsZipStream(fmt.Sprintf("%s/%s", ds.Fields["bucket"], resp.Job.TrainingId),
			resp.Job.Training.Resources.Learners, w)

		if err != nil {
			logr.WithError(err).Error("Downloading trained model failed")
			w.CloseWithError(err)
		}
		if err := w.Close(); err != nil {
			logr.WithError(err).Error("Closing writer failed")
		}
	}()

	reader := bufio.NewReader(r)
	buf := make([]byte, 0, 10*1024)
	for {
		n, err := reader.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				//logr.Errorf("Downloading trained model failed: %s", err.Error())
				break
			}
			return err
		}
		// process buf
		if err != nil && err != io.EOF {
			logr.WithError(err).Errorf("Downloading trained model failed")
			return err
		}
		err = stream.Send(&grpc_trainer_v2.ZippedDataChunk{
			Data: buf,
		})
		if err != nil {
			logr.WithError(err).Error("Failed to send zipped data chunk")
			return err
		}
	}
	return nil
}

func trainedModelLogRequestToTrainerQuery(req *grpc_trainer_v2.TrainedModelLogRequest, rindex int64, pageSize int32) *tdsService.Query {
	query := &tdsService.Query{
		Meta: &tdsService.MetaInfo{
			TrainingId: req.TrainingId,
			UserId:     req.UserId,
		},
		Pos:      rindex,
		Pagesize: pageSize,
	}
	return query
}

func (s *trainerService) isLearningFinished(req *grpc_trainer_v2.TrainedModelLogRequest) (bool, error) {
	tr, err := s.repo.Find(req.TrainingId)
	if err != nil {
		if err == mgo.ErrNotFound {
			// Maybe it was deleted.  Call it a day without reporting an error
			return true, nil
		}
		return true, err
	}
	statusID := tr.TrainingStatus.Status

	jobCompleted := false
	if statusID == grpc_trainer_v2.Status_COMPLETED ||
		statusID == grpc_trainer_v2.Status_FAILED {
		jobCompleted = true
	}

	return jobCompleted, nil
}

func (s *trainerService) waitUntilJobStart(req *grpc_trainer_v2.TrainedModelLogRequest,
	outStream grpc_trainer_v2.Trainer_GetTrainedModelLogsServer,
	logr *logger.LocLoggingEntry) error {

	startTime := time.Now()
	lastReportTime := time.Now()
	if req.Follow == true {
		for {
			tr, err := s.repo.Find(req.TrainingId)
			if err != nil {
				return err
			}
			statusID := tr.TrainingStatus.Status
			if !(statusID == grpc_trainer_v2.Status_NOT_STARTED ||
				statusID == grpc_trainer_v2.Status_QUEUED ||
				statusID == grpc_trainer_v2.Status_PENDING) {
				break
			}
			duration := time.Now().Sub(startTime)
			if duration > time.Minute*10 {
				err := errors.New(
					"gave up waiting for job to start when attempting to retrieve learner logs")
				logr.WithError(err).Debugf("gave up waiting")
				return err
			}
			durationSinceLastReport := time.Now().Sub(lastReportTime)
			if durationSinceLastReport.Seconds() == 15 {
				msg := fmt.Sprintf(
					"Waiting for training to start for log follow: %f minutes",
					duration.Minutes())
				logr.Debugf("%s", msg)
				errSend := outStream.Send(&grpc_trainer_v2.ByteStreamResponse{Data: []byte(msg)})
				if errSend != nil {
					logr.WithError(errSend).Errorf("cannot report status to user")
				}
				lastReportTime = time.Now()
			}

			time.Sleep(time.Second * 2)
		}
	}

	return nil
}

func (s *trainerService) GetTrainedModelLogs(req *grpc_trainer_v2.TrainedModelLogRequest,
	outStream grpc_trainer_v2.Trainer_GetTrainedModelLogsServer) error {

	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	//noinspection GoBoolExpressions
	dlogr := debugLogger(logr.Logger, debugLogsMode)

	dlogr.Debug("entry")

	err := s.waitUntilJobStart(req, outStream, logr)
	if err != nil {
		return err
	}

	tds, err := s.tdsClient()
	if err != nil {
		logr.WithError(err).Error("Cannot create LCM service client")
		return err
	}

	var rindex int64 = 1

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*4)
		defer cancel()

		// TODO: Create query from old request
		query := trainedModelLogRequestToTrainerQuery(req, rindex, oldEndpointInternalPageSize)

		dlogr.Debugf("Query to send to training-data: %+v", query)

		inStream, err := tds.Client().GetLogs(ctx, query)
		if err != nil {
			logr.WithError(err).Error("training data service GetLogs seems to have failed")
			return err
		}

		nRecordsFound := 0
		for {
			dlogr.Debugf("inStream.Recv()")
			chunk, err := inStream.Recv()
			if err == io.EOF {
				dlogr.Debug("eof")
				break
			}
			if err != nil {
				logr.WithError(err).Errorf("cannot read trained model log")
				return fmt.Errorf("cannot read trained model log: %v", err)
			}
			dlogr.Debugf("sending line: %d", chunk.Meta.Rindex)
			errSend := outStream.Send(&grpc_trainer_v2.ByteStreamResponse{Data: []byte(chunk.Line)})
			if errSend != nil {
				logr.WithError(errSend).Errorf("cannot send trained model log")
				return fmt.Errorf("cannot send trained model log: %v", err)
			}
			rindex++
			nRecordsFound++
			dlogr.Debugf("sent without error")
		}
		if nRecordsFound == 0 {
			if req.Follow == false {
				break
			}
			isDone, err := s.isLearningFinished(req)
			if err != nil {
				logr.WithError(err).Errorf("Can not get trainer status")
				return err
			}
			if isDone {
				break
			}

			time.Sleep(time.Second * 2)
		}
	}
	dlogr.Debug("exit with nil return")
	return nil
}

func marshalQuerySearchType(st grpc_trainer_v2.Query_SearchType) tdsService.Query_SearchType {
	searchType := tdsService.Query_TERM

	switch st {
	case grpc_trainer_v2.Query_TERM:
		searchType = tdsService.Query_TERM
		break
	case grpc_trainer_v2.Query_NESTED:
		searchType = tdsService.Query_NESTED
		break
	case grpc_trainer_v2.Query_MATCH:
		searchType = tdsService.Query_MATCH
		break
	case grpc_trainer_v2.Query_ALL:
		searchType = tdsService.Query_ALL
		break
	}
	return searchType
}

func marshalTDSQueryToTrainerQuery(in *grpc_trainer_v2.Query) *tdsService.Query {
	query := &tdsService.Query{
		Meta: &tdsService.MetaInfo{
			TrainingId: in.Meta.TrainingId,
			UserId:     in.Meta.UserId,
			Time:       in.Meta.Time,
			Rindex:     in.Meta.Rindex,
			Subid:      in.Meta.Subid,
		},
		Pos:        in.Pos,
		Pagesize:   in.Pagesize,
		Since:      in.Since,
		SearchType: marshalQuerySearchType(in.SearchType),
	}
	return query
}

func (s *trainerService) GetTrainingLogs(in *grpc_trainer_v2.Query,
	outStream grpc_trainer_v2.Trainer_GetTrainingLogsServer) error {

	logr := logger.LocLogger(logWith(in.Meta.TrainingId, in.Meta.UserId))

	//noinspection GoBoolExpressions
	dlogr := debugLogger(logr.Logger, debugLogsMode)

	dlogr.Debug("entry")

	tds, err := s.tdsClient()
	if err != nil {
		logr.WithError(err).Error("Cannot create LCM service client")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*4)
	defer cancel()

	dlogr.Debugf("Query to send from client: %+v", in)

	query := marshalTDSQueryToTrainerQuery(in)

	dlogr.Debugf("Query to send to training-data: %+v", query)

	inStream, err := tds.Client().GetLogs(ctx, query)

	if err != nil {
		logr.WithError(err).Error("training data service GetLogs seems to have failed")
		return err
	}

	for {
		dlogr.Debugf("inStream.Recv()")
		chunk, err := inStream.Recv()
		if err == io.EOF {
			dlogr.Debug("eof")
			break
		}
		if err != nil {
			logr.WithError(err).Errorf("cannot read trained model log")
			return fmt.Errorf("cannot read trained model log: %v", err)
		}
		dlogr.Debugf("sending line: %d", chunk.Meta.Rindex)
		errSend := outStream.Send(&grpc_trainer_v2.LogLine{
			Meta: &grpc_trainer_v2.MetaInfo{
				TrainingId: chunk.Meta.TrainingId,
				UserId:     chunk.Meta.UserId,
				Time:       chunk.Meta.Time,
				Rindex:     chunk.Meta.Rindex,
				Subid:      chunk.Meta.Subid,
			},
			Line: chunk.Line,
		})
		if errSend != nil {
			logr.WithError(errSend).Errorf("cannot send trained model log")
			return fmt.Errorf("cannot send trained model log: %v", err)
		}
		dlogr.Debugf("sent without error")
	}
	dlogr.Debug("exit with nil return")
	return nil
}

func marshalTDSDataType2TrainerDataType(dt tdsService.Any_DataType) grpc_trainer_v2.Any_DataType {
	dataType := grpc_trainer_v2.Any_STRING

	switch dt {
	case tdsService.Any_STRING:
		dataType = grpc_trainer_v2.Any_STRING
		break
	case tdsService.Any_JSONSTRING:
		dataType = grpc_trainer_v2.Any_JSONSTRING
		break
	case tdsService.Any_INT:
		dataType = grpc_trainer_v2.Any_INT
		break
	case tdsService.Any_FLOAT:
		dataType = grpc_trainer_v2.Any_FLOAT
		break
	}
	return dataType
}

func marshalTDSMapToTrainerMap(tdsMap map[string]*tdsService.Any) map[string]*grpc_trainer_v2.Any {
	grpcMap := make(map[string]*grpc_trainer_v2.Any)
	for k, v := range tdsMap {
		trainerDT := marshalTDSDataType2TrainerDataType(v.Type)
		grpcMap[k] = &grpc_trainer_v2.Any{Type: trainerDT, Value: v.Value}
	}
	return grpcMap
}

func (s *trainerService) GetVersions(ctx context.Context, req *grpc_trainer_v2.GetVersionsRequest) (*grpc_trainer_v2.Frameworks, error) {
	//call the frameworks.go and then getAllVersions for the frameworks
	//Return response from getAll Versions
	frameworks, err := getExternalVersions()
	if err != nil {
		return nil, err
	}
	return &frameworks, nil
}

func (s *trainerService) validateRequest(log *logrus.Entry, req *grpc_trainer_v2.CreateRequest) error {
	if req.UserId == "" {
		return s.failCreateRequest("UserId is nil", req, log)
	}

	// validate model definition object

	m := req.ModelDefinition
	if m == nil {
		return s.failCreateRequest("Model definition is not set", req, log)
	}
	if m.Name == "" {
		return s.failCreateRequest("Model definition name is not set", req, log)
	}
	if m.Framework == nil {
		return s.failCreateRequest("Framework is not set", req, log)
	}
	if m.Framework.Name == "" {
		return s.failCreateRequest("Framework name is not set", req, log)
	}

	if m.Framework.Version == "" {
		return s.failCreateRequest("Framework version is not set", req, log)
	}

	// FfDL Change: FIXME: temporarily disable framework validation
	// custom image check
	//if m.Framework.ImageLocation == nil {
	//	if ok, msg := validateFrameworks(m.Framework); !ok {
	//		return s.failCreateRequest(msg, req, log)
	//	}
	//}
	//
	//FIXME MLSS Change: v_1.5.1 remove logic for model definition content
	//if len(m.Content) == 0 {
	//	return s.failCreateRequest("Model definition content is not set", req, log)
	//}

	// validate Training object

	// t := req.Training
	// if t == nil {
	// 	return s.failCreateRequest("Training is not set", req, log)
	// }
	// if t.Command == "" {
	// 	return s.failCreateRequest("Training command is not set", req, log)
	// }
	// if t.InputData == nil || len(t.InputData) == 0 {
	// 	return s.failCreateRequest("Training input data is not set", req, log)
	// }
	// if len(t.InputData) > 1 {
	// 	return s.failCreateRequest("Training input data can only contain one id", req, log)
	// }
	// if t.OutputData != nil && len(t.OutputData) > 1 {
	// 	return s.failCreateRequest("Training output data can only contain one id", req, log)
	// }

	// // validate datastores

	// ds := req.Datastores
	// if ds == nil {
	// 	return s.failCreateRequest("Data stores is not set", req, log)
	// }
	// if len(ds) == 0 {
	// 	return s.failCreateRequest("Data stores is empty", req, log)
	// }

	// for _, name := range t.InputData {
	// 	ds := findDatastore(name, req.Datastores)
	// 	if ds == nil {
	// 		return s.failCreateRequest(fmt.Sprintf("Training input data reference '%s' does not reference an existing datastore id.", name), req, log)
	// 	}
	// 	if err := s.validateDatastore(ds, req, log); err != nil {
	// 		return err
	// 	}
	// }

	// if len(t.OutputData) > 0 {
	// 	for _, name := range t.OutputData {
	// 		ds := findDatastore(name, req.Datastores)
	// 		if ds == nil {
	// 			return s.failCreateRequest(fmt.Sprintf("Training output data reference '%s' does not reference an existing datastore id.", name), req, log)
	// 		}
	// 		if err := s.validateDatastore(ds, req, log); err != nil {
	// 			return err
	// 		}
	// 	}
	// }
	return nil
}

func findDatastore(id string, ds []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	for _, v := range ds {
		if v.Id == id {
			return v
		}
	}
	return nil
}

func FindDatastore(id string, ds []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	return findDatastore(id, ds)
}

func (s *trainerService) failCreateRequest(msg string, req *grpc_trainer_v2.CreateRequest, log *logrus.Entry) error {
	return s.failCreateRequestWithCode(trainerClient.ErrInvalidManifestFile, msg, req, log)
}

func (s *trainerService) failCreateRequestWithCode(errorCode string, msg string, req *grpc_trainer_v2.CreateRequest, log *logrus.Entry) error {
	log.Errorf("Failed to validate CreateRequest: %s", msg)

	// send error event as monitoring metric

	// counter := s.metrics.trainingJobFailedCounter.With("type", "client", "errorcode", errorCode)
	if req.ModelDefinition != nil && req.ModelDefinition.Framework != nil {
		// counter = counter.With("framework", req.ModelDefinition.Framework.Name, "version", req.ModelDefinition.Framework.Version)

		logFrameworkErrorsValue := fmt.Sprintf("%s-%s-%s", req.ModelDefinition.Framework.Name, "client", errorCode)
		log.WithFields(logrus.Fields{
			"framework_errors": logFrameworkErrorsValue,
		})
	}

	// if req.Training != nil && req.Training.Resources != nil {
	// 	counter = counter.With("gpus", strconv.Itoa(int(req.Training.Resources.Gpus)),
	// 		"cpus", strconv.Itoa(int(req.Training.Resources.Cpus)),
	// 		"memory", strconv.Itoa(int(req.Training.Resources.Memory)))
	// }

	log.Debug("Metrics for failed training jobs framework")
	// counter.Add(1)

	return gerrf(codes.InvalidArgument, msg)
}

func (s *trainerService) validateDatastore(ds *grpc_trainer_v2.Datastore, req *grpc_trainer_v2.CreateRequest, log *logrus.Entry) error {

	if ds == nil {
		return s.failCreateRequest("Data store is not set", req, log)
	}
	log.Debugf("validateDatastore, Datastore type: %s", ds.Type)
	if ds.Type == "mount_cos" {
		if ds.Id == "" {
			return s.failCreateRequest("Data store id is not set", req, log)
		}
		if ds.Connection == nil || len(ds.Connection) == 0 {
			return s.failCreateRequest("Data store connection info not set", req, log)
		}
		if ds.Fields == nil || len(ds.Fields) == 0 || ds.Fields["bucket"] == "" {
			return s.failCreateRequest("Data store bucket is not set", req, log)
		}

		//ostore, err := storage.CreateDataStore(ds.Type, ds.Connection)
		//if err != nil {
		//	log.Errorf("Validation failed: %s", err.Error())
		//	return s.failCreateRequestWithCode(trainerClient.ErrInvalidCredentials,
		//		fmt.Sprintf("Data store authentication information for id '%s' incorrect or there is a connection problem", ds.Id), req, log)
		//}

		//if err := ostore.Connect(); err != nil {
		//	log.Errorf("Validation failed: %s", err.Error())
		//	return s.failCreateRequestWithCode(trainerClient.ErrInvalidCredentials,
		//		fmt.Sprintf("Data store authentication information for id '%s' incorrect or there is a connection problem", ds.Id), req, log)
		//}

		// validate bucket (or container as it is called in Swift)
		bucket := ds.Fields["bucket"]
		if bucket != "" {
			//exists, err := ostore.ContainerExists(bucket)
			exists, err := bucketIsExists(bucket)
			if !exists || err != nil {
				return s.failCreateRequestWithCode(trainerClient.ErrInvalidCredentials,
					fmt.Sprintf("Data store bucket '%s' for data store id '%s' incorrect, there may be a connection problem or credentials do not allow access to the bucket", bucket, ds.Id), req, log)
			}
		}
	}
	log.Debugf("validate datastore successfully")
	return nil
}

func bucketIsExists(bucket string) (bool, error) {
	todo := context.TODO()
	ctx, cancel := context.WithCancel(todo)
	getLogger := logger.GetLogger()
	sClient, err := storageClient.NewStorage()
	if err != nil {
		getLogger.Errorf("storageClient.NewStorage failed, %v", err.Error())
		defer cancel()
		return false, err
	}
	defer sClient.Close()

	request := grpc_storage.BucketExistsRequest{
		Bucket: bucket,
	}
	response, err := sClient.Client().BucketExists(ctx, &request)
	if err != nil {
		getLogger.Errorf("check BucketExists failed, %v", err.Error())
		defer cancel()
		return false, err
	}
	getLogger.Debugf("BucketExists, %+v", response)
	return true, nil
}

// lcmClient established a connection if the trainerService has nothing existing cached
func (s *trainerService) lcmClient() (client.LcmClient, error) {
	if s.lcm == nil {
		return client.NewLcm(nil)
	}
	return s.lcm, nil
}

func (s *trainerService) tdsClient() (tdsClient.TrainingDataClient, error) {
	if s.tds == nil {
		address := fmt.Sprintf("ffdl-trainingdata.%s.svc.cluster.local:80", config.GetPodNamespace())
		tds, err := tdsClient.NewTrainingDataClientWithAddress(address)
		if err != nil {
			return nil, err
		}
		s.tds = tds
	}
	return s.tds, nil
}

func (s *trainerService) createJobConfig(tr *TrainingRecord) (*service.JobDeploymentRequest, error) {
	logr := logger.LocLogger(logWith(tr.TrainingID, tr.UserID))

	// training data/results - assume only one training input and output data at this point
	trainingData := findDatastore(tr.Training.InputData[0], tr.Datastores)
	trainingResults := s.getOutputDatastore(tr.Training.OutputData, tr.Datastores)

	// Environment variables
	envvars := make(map[string]string)

	// Fetching data from user's Object Store with following info
	envvars["DATA_STORE_TYPE"] = trainingData.Type
	envvars["DATA_STORE_AUTHURL"] = trainingData.Connection["auth_url"]
	if trainingData.Connection["path"] != "" {
		envvars["DATA_STORE_PATH"] = trainingData.Connection["path"]
	}
	if trainingData.Connection["project_id"] != "" {
		envvars["DATA_STORE_PROJECTID"] = trainingData.Connection["project_id"]
	}
	if trainingData.Connection["type"] != "" {
		envvars["DATA_STORE_TYPE"] = trainingData.Connection["type"]
	}
	if trainingData.Connection["user_name"] != "" {
		envvars["DATA_STORE_USERNAME"] = trainingData.Connection["user_name"]
	}
	if trainingData.Connection["password"] != "" {
		envvars["DATA_STORE_APIKEY"] = trainingData.Connection["password"]
	}
	if trainingData.Connection["domain_name"] != "" {
		envvars["DATA_STORE_DOMAINNAME"] = trainingData.Connection["domain_name"]
	}
	if trainingData.Connection["region"] != "" {
		envvars["DATA_STORE_REGION"] = trainingData.Connection["region"]
	}
	envvars["DATA_STORE_OBJECTID"] = trainingData.Fields["bucket"]

	// Allow to fetch model from DLaaS's Object Store to the container
	osConf := config.GetDataStoreConfig()
	envvars["MODEL_STORE_USERNAME"] = osConf[storage.UsernameKey]
	envvars["MODEL_STORE_APIKEY"] = osConf[storage.PasswordKey]
	envvars["MODEL_STORE_AUTHURL"] = osConf[storage.AuthURLKey] // this will inside SL so we need the internal one
	if osConf[storage.StorageType] != "" {
		envvars["MODEL_STORE_TYPE"] = osConf[storage.StorageType]
	}
	// only needed for Bluemix objectstore
	if val, ok := osConf[storage.DomainKey]; ok {
		envvars["MODEL_STORE_DOMAINNAME"] = val
	}
	if val, ok := osConf[storage.RegionKey]; ok {
		envvars["MODEL_STORE_REGION"] = val
	}
	if val, ok := osConf[storage.DomainKey]; ok {
		envvars["MODEL_STORE_PROJECTID"] = val
	}
	envvars["MODEL_STORE_OBJECTID"] = tr.ModelDefinition.Location

	// "Storing trained model in DLaaS Object Store with following info:"
	envvars["RESULT_STORE_TYPE"] = trainingResults.Type
	if trainingData.Connection["path"] != "" {
		envvars["RESULT_STORE_PATH"] = trainingData.Connection["path"]
	}
	envvars["RESULT_STORE_USERNAME"] = trainingResults.Connection["user_name"]
	envvars["RESULT_STORE_APIKEY"] = trainingResults.Connection["password"]
	envvars["RESULT_STORE_AUTHURL"] = trainingResults.Connection["auth_url"]
	if trainingResults.Connection[storage.StorageType] != "" {
		envvars["RESULT_STORE_TYPE"] = trainingResults.Connection[storage.StorageType]
	}
	// only needed for Bluemix objectstore
	if trainingResults.Connection["domain_name"] != "" {
		envvars["RESULT_STORE_DOMAINNAME"] = trainingResults.Connection["domain_name"]
	}
	if trainingResults.Connection["region"] != "" {
		envvars["RESULT_STORE_REGION"] = trainingResults.Connection["region"]
	}
	if trainingResults.Connection["project_id"] != "" {
		envvars["RESULT_STORE_PROJECTID"] = trainingResults.Connection["project_id"]
	}
	envvars["RESULT_STORE_OBJECTID"] = fmt.Sprintf("%s/%s", trainingResults.Fields["bucket"], tr.TrainingID)

	// Storing data in container at
	envvars["DATA_DIR"] = trainingData.Fields["bucket"]

	envvars["CODE_WORK_DIR"] = trainingData.Fields["work"]

	// Storing model in container at
	envvars["MODEL_DIR"] = "/model-code"

	// Storing trained model at
	envvars["RESULT_DIR"] = trainingResults.Fields["bucket"]

	// TODO: This is pointing to currently where the logs are put, but should be redefined per nfs log mount proposal.
	// (by the time it gets to the learners/log-collectors, it will be "/job/logs", at the time of this writing.)
	envvars["LOG_DIR"] = "/logs"

	re := regexp.MustCompile(`\r?\n`)
	input := re.ReplaceAllString(fmt.Sprint(tr.Training.Command), " ")

	envvars["TRAINING_COMMAND"] = input

	envvars["TRAINING_ID"] = tr.TrainingID

	envvars["CODE_DATA_PATH"] = tr.DataPath

	envvars["CODE_SELECTOR"] = tr.CodeSelector

	envvars["GPU_COUNT"] = strconv.FormatFloat(float64(tr.Training.Resources.Gpus), 'f', 6, 64)

	envvars["SCHED_POLICY"] = strings.ToLower(tr.Training.Resources.Schedpolicy)

	// tag to use to lookup learner image to use; this is a Docker image tag
	if tr.ModelDefinition.Framework.ImageTag != "" {
		envvars["DLAAS_LEARNER_IMAGE_TAG"] = tr.ModelDefinition.Framework.ImageTag
	}

	// envvar for profile
	if tr.Training.Profiling {
		envvars["DLAAS_PROFILING"] = "true"
	}

	// FIXME MLSS Change: add GID & UID & USER_ID to envvar
	envvars["GID"] = tr.GID
	envvars["UID"] = tr.UID
	envvars["GUARDIAN_TOKEN"] = tr.GuardianToken
	envvars["USER_ID"] = tr.UserID
	// FIXME MLSS Change: v_1.4.1
	//envvars["EVENT_CHECKER"] = tr.EventChecker
	//envvars["DEADLINE_CHECKER"] = tr.DeadlineChecker
	//envvars["OVERTIME_CHECKER"] = tr.OvertimeChecker
	//envvars["ALERT_LEVEL"] = tr.AlertLevel
	//envvars["RECEIVER"] = tr.Receiver
	//envvars["INTERVAL"] = tr.Interval

	// labels
	labels := make(map[string]string)
	labels["training_id"] = tr.TrainingID
	labels["user_id"] = tr.UserID
	labels["gpu_type"] = tr.Training.Resources.GpuType

	u4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	logr.Debugf("Training job env vars: %v", envvars)
	logr.Debugf(" tr.Training: %v", tr.Training)
	logr.Debugf(" tr.ProxyUser: %v", tr.ProxyUser)
	job := &service.JobDeploymentRequest{
		Name:                  u4.String(),
		Resources:             getResourceRequirements(tr.Training),
		EnvVars:               envvars,
		Labels:                labels,
		UserId:                tr.UserID,
		TrainingId:            tr.TrainingID,
		Framework:             tr.ModelDefinition.Framework.Name,
		Version:               tr.ModelDefinition.Framework.Version,
		ImageLocation:         parseImageLocation(tr),
		EvaluationMetricsSpec: tr.EvaluationMetricsSpec,
		JobNamespace:          tr.Namespace,
		// FIXME MLSS Change: v_1.4.1 add DataStore & ModelDefinition & Training from tr to JobDeploymentRequest JobAlert
		DataStores:      tr.Datastores,
		ModelDefinition: tr.ModelDefinition,
		Training:        tr.Training,
		TrainingStatus:  tr.TrainingStatus,
		JobAlert:        tr.JobAlert,
		// FIXME MLSS Change: v_1.5.2 added ParameterServer
		PSs:      tr.PSs,
		PSCPU:    tr.PSCPU,
		PSImage:  tr.PSImage,
		PSMemory: tr.PSMemory,
		// FIXME MLSS Change: v_1.5.1 added CodeSelector and DataPath
		CodeSelector: tr.CodeSelector,
		DataPath:     tr.DataPath,
		JobType:      tr.JobType,
		TFosRequest:  tr.TFosRequest,
		ExpRunId:     tr.ExpRunId,
		ExpName:      tr.ExpName,
		FileName:     tr.FileName,
		FilePath:     tr.FilePath,
		ProxyUser:    tr.ProxyUser,
		DataSet:      tr.DataSet,
		MfModel: 	  tr.MFModel,
		Algorithm:     tr.Algorithm,
		JobParams:  tr.JobParams,
		FitParams: tr.FitParams,
		APIType: tr.APIType,
	}

	return job, nil
}

func parseImageLocation(tr *TrainingRecord) *service.ImageLocation {
	tril := tr.ModelDefinition.Framework.ImageLocation
	var il (*service.ImageLocation)
	if tril != nil {
		il = &service.ImageLocation{
			Registry:    tril.Registry,
			Namespace:   tril.Namespace,
			AccessToken: tril.AccessToken,
			Email:       tril.Email,
		}
	}
	return il
}

func ParseImageLocation(tr *TrainingRecord) *service.ImageLocation {
	return parseImageLocation(tr)
}

func setDefaultResourceRequirements(t *grpc_trainer_v2.Training) {
	if t == nil || t.Resources == nil {
		t.Resources = &grpc_trainer_v2.ResourceRequirements{ // set sensible defaults
			Cpus:        5.0,
			Gpus:        1.0,
			Memory:      12.0,
			MemoryUnit:  grpc_trainer_v2.SizeUnit_GiB,
			Learners:    1,
			Schedpolicy: "dense",
		}
		return
	}
	if t.Resources.Cpus == 0 {
		t.Resources.Cpus = 5.0
	}
	if t.Resources.Memory == 0 {
		t.Resources.Memory = 12
		t.Resources.MemoryUnit = grpc_trainer_v2.SizeUnit_GiB
	}
	if t.Resources.Schedpolicy == "" || strings.ToLower(t.Resources.Schedpolicy) != "spread" {
		t.Resources.Schedpolicy = "dense"
	}
	if t.Resources.GpuType == "" {
		t.Resources.GpuType = "nvidia-TeslaK80"
	}
}

func getResourceRequirements(t *grpc_trainer_v2.Training) *service.ResourceRequirements {
	return &service.ResourceRequirements{
		Cpus:        float64(t.Resources.Cpus),
		Gpus:        float64(t.Resources.Gpus),
		Memory:      float64(t.Resources.Memory),
		MemoryUnit:  service.ResourceRequirements_MemoryUnit(service.ResourceRequirements_MemoryUnit_value[t.Resources.MemoryUnit.String()]),
		Storage:     float64(t.Resources.Storage),
		StorageUnit: service.ResourceRequirements_MemoryUnit(service.ResourceRequirements_MemoryUnit_value[t.Resources.StorageUnit.String()]),
		Learners:    t.Resources.Learners,
		GpuType:     t.Resources.GpuType,
	}
}

func GetResourceRequirements(t *grpc_trainer_v2.Training) *service.ResourceRequirements {
	return getResourceRequirements(t)
}

// getOutputDatastore retrieves the output data store or return the internal datastore if none has been defined
func (s *trainerService) getOutputDatastore(outputData []string, datastores []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	var ds *grpc_trainer_v2.Datastore
	if len(outputData) > 0 {
		ds = findDatastore(outputData[0], datastores) // we assume there is only one output data at this point b/c the underlying system does not support more
	}
	if ds == nil {
		ds = &grpc_trainer_v2.Datastore{
			Id:         internalObjectStoreID,
			Type:       config.GetDataStoreType(),
			Connection: config.GetDataStoreConfig(),
			Fields:     map[string]string{"bucket": s.trainedModelsBucket},
		}
	}
	return ds
}

// getOutputStoreForService is a wrapper function to make the logic in trainerService.getOutputDatastore available for testing
func getOutputStoreForService(s *trainerService, outputData []string, datastores []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	return s.getOutputDatastore(outputData, datastores)
}

func getModelsBucket() string {
	if viper.IsSet(modelsBucketKey) {
		return viper.GetString(modelsBucketKey)
	}
	return defaultModelsBucket
}

// FIXME MLSS Change: v_1.4.1 expose the func
func GetModelsBucket() string {
	return getModelsBucket()
}

func getTrainedModelsBucket() string {
	if viper.IsSet(trainedModelsBucketKey) {
		return viper.GetString(trainedModelsBucketKey)
	}
	return defaultTrainedModelsBucket
}

func getModelZipFileName(trainingID string) string {
	return fmt.Sprintf("%s.zip", trainingID)
}

func (s *trainerService) submitJobToLCM(tr *TrainingRecord, logr *logger.LocLoggingEntry) error {
	logr.Println("submitJobToLCM, tr.JobType: ", tr.JobType)

	jobConfig, err := s.createJobConfig(tr)
	if err != nil {
		logr.WithError(err).Errorf("Failed to create job config")
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}
	logr.Debugf("submitJobToLCM jobConfig.EnvVars TOKEN: %v", jobConfig.EnvVars["GUARDIAN_TOKEN"])

	lcm, err := s.lcmClient()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create LCM service client")
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}
	defer lcm.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	logr.Debug("DeployTrainingJob1-----------------1")
	logr.Debug("DeployTrainingJob1-----------------1" + reflect.TypeOf(jobConfig.EnvVars).String())
	logr.Debug("JOBCONFIG Trainging : ", jobConfig.Training)
	logr.Debug("JOBCONFIG Trainging Status: ", jobConfig.TrainingStatus)
	_, err = lcm.Client().DeployTrainingJob(ctx, jobConfig)
	logr.Debug("DeployTrainingJob12-----------------12")
	if err != nil {
		logr.WithError(err).Errorf("Cannot deploy training job with id %s", tr.TrainingID)
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}

	logr.Printf("training job %s submitted to lcm", tr.TrainingID)

	// capture the gpu usage when the job is submitted to LCM
	gpusUsed := tr.Training.Resources.Gpus
	if tr.Training.Resources.Learners > 1 {
		gpusUsed = tr.Training.Resources.Gpus * float64(tr.Training.Resources.Learners)
	}
	// log the gpu usages requested by user
	logGpuTypeIncrementValue := fmt.Sprintf("%s-%v", tr.Training.Resources.GpuType, gpusUsed)
	logr.WithFields(logrus.Fields{
		"gputype_increment": logGpuTypeIncrementValue,
	}).Debug(" incrementing the gpus")

	// increment the counter
	// s.metrics.clusterWideGPUUsageCounter.With("gpuType", tr.Training.Resources.GpuType, "gpus", strconv.Itoa(int(gpusUsed))).Add(1)

	return nil
}

// determine if this job should be rate-limited by checking if the total number of GPUs
// would exceed the limit set for the GPU type
func (s *trainerService) rateLimitTrainingJob(trainingRecord *TrainingRecord, logr *logger.LocLoggingEntry) bool {
	var rateLimit = false

	gpuType := trainingRecord.Training.Resources.GpuType
	limit := getGpuLimitByType(gpuType)
	if limit == 0 {
		return false
	}

	// get total number of GPUs requested for this job
	gpusRequested := trainingRecord.Training.Resources.Gpus
	if trainingRecord.Training.Resources.Learners > 1 {
		// trainingRecord.Training.Resources.Gpus is GPUs used per learner
		gpusRequested = trainingRecord.Training.Resources.Gpus * float64(trainingRecord.Training.Resources.Learners)
	}
	if gpusRequested == 0 {
		return false
	}

	// find the GPUs used that count toward this limit
	records, err := s.repo.FindCurrentlyRunningTrainings(getGpuLimitQuerySize())
	logr.Debugf("running records (%d)", len(records))
	if err != nil || len(records) == 0 {
		logr.WithError(err).Warnf("did not execute rate limiting correctly, returned number of records count is %d", len(records))
		return false
	}
	var totalGPUsUsedCount float64
	var matchingGPUConsumingRecords []*TrainingRecord
	for _, record := range records {
		trainingStatus := record.TrainingStatus.Status
		if trainingStatus == grpc_trainer_v2.Status_COMPLETED || trainingStatus == grpc_trainer_v2.Status_FAILED || trainingStatus == grpc_trainer_v2.Status_QUEUED {
			//ignore these since they don't count towards active resource usage
		} else if TransformResourceName(record.Training.Resources.GpuType) == TransformResourceName(gpuType) {
			//only count matching gpu type
			matchingGPUConsumingRecords = append(matchingGPUConsumingRecords, record)
			gpusUsed := record.Training.Resources.Gpus
			if record.Training.Resources.Learners > 1 {
				gpusUsed = record.Training.Resources.Gpus * float64(record.Training.Resources.Learners)
			}
			totalGPUsUsedCount = totalGPUsUsedCount + float64(gpusUsed)
		}
	}
	// s.metrics.clusterWideGPUUsageGauge.With("gpuType", gpuType).Set(float64(totalGPUsUsedCount))

	if int64(totalGPUsUsedCount+float64(gpusRequested)) > limit {
		rateLimit = true
		if logr.Logger.Level >= logrus.DebugLevel {
			for _, record := range matchingGPUConsumingRecords {
				logr.Debugf("Found a gpu consuming training %v has a status %s and using gpus %v with submission time as %v and process start time as %v and error code %v",
					record.TrainingID, record.TrainingStatus.Status, record.Training.Resources.Gpus, record.TrainingStatus.SubmissionTimestamp,
					record.TrainingStatus.ProcessStartTimestamp, record.TrainingStatus.ErrorCode)
			}
		}
	}

	logr.Debugf("result of rate-limiting for job %s is %t; gpu type %s has limit %d, total used %v, requested %v",
		trainingRecord.TrainingID, rateLimit, TransformResourceName(gpuType), limit, totalGPUsUsedCount, gpusRequested)
	// if rateLimit {
	// 	s.metrics.rateLimitTrainingJobCounter.Add(1)
	// }

	return rateLimit
}

func getGpuLimitQuerySize() int {
	return viper.GetInt(gpuLimitsQuerySizeKey)
}

func GetGpuLimitQuerySize() int {
	return getGpuLimitQuerySize()
}

//getAllResourceTypes returns all GPU types defined in resource limits
func getAllResourceTypes() []string {
	types := []string{}
	allLimits := strings.Split(viper.GetString(gpuLimitsKey), ",")
	for _, limit := range allLimits {
		if strings.Contains(limit, "=") {
			types = append(types, TransformResourceName(strings.Split(limit, "=")[0]))
		}
	}
	return types
}

//getGpuLimitByType returns the resource limit if it is defined, or returns 0 if not
func getGpuLimitByType(gpuType string) int64 {
	limit := int64(0)
	allLimits := strings.Split(viper.GetString(gpuLimitsKey), ",")
	for _, l := range allLimits {
		if TransformResourceName(strings.Split(l, "=")[0]) == TransformResourceName(gpuType) {
			lim, err := strconv.ParseInt(strings.Split(l, "=")[1], 10, 0)
			if err == nil {
				limit = lim
			}
			break
		}
	}
	return limit
}

// FIXME MLSS Change: v_1.4.1 added
func GetGpuLimitByType(gpuType string) int64 {
	return getGpuLimitByType(gpuType)
}

// FIXME MLSS Change: delete model.zip from s3
func (s *trainerService) DeleteSubmittedCode(ctx context.Context, req *grpc_trainer_v2.DeleteRequest) (*grpc_trainer_v2.DeleteResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("DeleteSubmittedCode enter, TrainingId: %s,UserId: %s", req.TrainingId, req.UserId)

	bucket := GetModelsBucket()
	TrainingId := req.TrainingId
	zipName := TrainingId + ".zip"

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	//ds := s.getOutputDatastore(resp.Job.Training.OutputData, resp.Job.Datastores)

	ds := s.getOutputDatastore(nil, nil)
	logr.Debugf("DeleteSubmittedCode, type: %s, auth_url: %s, user_name: %s, password: %s", ds.Type, ds.Connection["auth_url"], ds.Connection["user_name"], ds.Connection["password"])
	logr.Debugf(" Deleting model.zip, bucket: %s, zipName: %s", bucket, zipName)

	err := s.datastore.DeleteArchive(s.modelsBucket, zipName)
	if err != nil {
		logr.WithError(err).Errorf("Error deleting model from object store")
		// log this error, but continue with deleting the training record anyway
		return nil, err
	}

	if err != nil {
		logr.WithError(err).Errorf("DeleteArchive error: %s", err.Error())
		return nil, err
	}
	return &grpc_trainer_v2.DeleteResponse{TrainingId: TrainingId}, nil
}
