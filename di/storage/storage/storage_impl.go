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

package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
	"webank/DI/commons/util"
	"webank/DI/commons/util/random"
	// "webank/DI/linkisexecutor/linkisclientutils"
	// "webank/DI/linkisexecutor/models"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	uuid2 "github.com/google/uuid"
	"github.com/mholt/archiver/v3"

	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/commons/service/client"
	tdsClient "webank/DI/metrics/client"
	tdsService "webank/DI/metrics/service/grpc_training_data_v1"
	storageClient "webank/DI/storage/client"
	"webank/DI/storage/storage/grpc_storage"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"time"

	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"errors"

	"webank/DI/trainer/instrumentation"
	"webank/DI/trainer/storage"

	"github.com/jinzhu/copier"
	"github.com/minio/minio-go/v6"
	es "gopkg.in/olivere/elastic.v5"
)

const internalObjectStoreID = "dlaas_internal_os"

const (
	indexName = "tds_index" // New alias
)

const (
	modelsBucketKey        = "objectstore.bucket.models"
	trainedModelsBucketKey = "objectstore.bucket.trainedmodels"

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
	mongoAuthenticationDatabase = "MONGO_Authentication_Database"

	gpuLimitsKey          = "gpu.limits"
	gpuLimitsQuerySizeKey = "gpu.limits.query.size"

	pollIntervalKey = "queue.poll.interval"

	defaultValue = "default"
)

const (
	// Used by a counter metric.
	dlaasStoreKind  = "dlaas"
	userStoreKind   = "user"
	docTypeLog      = "logline"
	docTypeEmetrics = "emetrics"
	defaultPageSize = 10
	indexNameV1     = "dlaas_learner_data"
	indexNameV2     = "dlaas_learner_data_v2"

	metaSubRecord = `"meta" : {
		"properties" : {
			"trainer_id" : { "type" : "keyword", "index" : "not_analyzed" },
			"user_id" : { "type" : "keyword", "index" : "not_analyzed" },
			"time" : { "type" : "long" },
			"rindex" : { "type" : "integer" },
			"subid" : { "type" : "keyword", "null_value" : "NULL", "index" : "not_analyzed" }
		}
	}`
	indexMappingLogs = `{
                            "logline" : {
                                "properties" : {
									` + metaSubRecord + `,
                                    "line" : { "type" : "text", "index" : "not_analyzed" }
                                }
                            }
                        }`

	indexMappingEmetrics = `{
                            "emetrics" : {
                                "properties" : {
									` + metaSubRecord + `,
	                       			"grouplabel" : { "type" : "text", "index" : "not_analyzed" }
                                }
                            }
                        }`
)

// Confuse `go vet' to not check this `Errorf' call. :(
// See https://github.com/grpc/grpc-go/issues/90
var gerrf = status.Errorf

var (
	// TdsDebugMode = viper.GetBool(TdsDebug)
	TdsDebugMode = false

	// TdsDebugLogLineAdd outputs diagnostic loglines as they are added.  Should normally be false.
	TdsDebugLogLineAdd = false

	// TdsDebugEMetricAdd outputs diagnostic emetrics as they are added.  Should normally be false.
	TdsDebugEMetricAdd = false

	// TdsReportTimes if try report timings for Elastic Search operations
	TdsReportTimes = false
)

//var shortidGenHigh = shortid.Shortid{}
//var shortidGenLow = shortid.Shortid{}

// Service represents the functionality of the trainer service
type Service interface {
	grpc_storage.StorageServer
	service.LifecycleHandler
	StopTrainer()
}


type queueHandler struct {
	stopQueue chan struct{}
	*TrainingJobQueue
}

type storageService struct {
	mtx                 sync.RWMutex //this lock should only be used for instantiating job queue
	datastore           storage.DataStore
	lcm                 client.LcmClient
	repo                Repository
	jobHistoryRepo      jobHistoryRepository
	modelsBucket        string
	trainedModelsBucket string
	tds                 tdsClient.TrainingDataClient
	queues              map[string]*queueHandler
	queuesStarted       bool
	service.Lifecycle

	es              *es.Client
	esBulkProcessor *es.BulkProcessor
}

// NewService creates a new trainer service.
func NewService() Service {
	logr := logger.LogServiceBasic(logger.LogkeyStorageService)

	config.FatalOnAbsentKey(mongoAddressKey)
	config.SetDefault(gpuLimitsQuerySizeKey, 200)
	config.SetDefault(pollIntervalKey, 60) // in seconds

	//get s3 dataStore
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
		viper.GetString(mongoPasswordKey), viper.GetString(mongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "training_jobs")
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create Trainings Repository with %s %s %s", viper.GetString(mongoAddressKey), viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey))
	}

	jobHistoryRepo, err := newJobHistoryRepository(viper.GetString(mongoAddressKey),
		viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey),
		viper.GetString(mongoPasswordKey), viper.GetString(mongoAuthenticationDatabase),
		config.GetMongoCertLocation(), collectionNameJobHistory)
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create Job History Repository with %s %s %s %s", viper.GetString("mongo.address"),
			viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey), collectionNameJobHistory)
	}

	queues := make(map[string]*queueHandler)

	s := &storageService{
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
		grpc_storage.RegisterStorageServer(s.Server, s)
	}
	//s.StartQueues()
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

func (s *storageService) StartQueues() {
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

func (s *storageService) StopTrainer() {
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

func (s *storageService) pullJobFromQueue(gpuType string) {
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
	if err != nil {
		if err == mgo.ErrNotFound {
			logr.Debugf("job %s not found in mongo, assuming job was deleted", nextJobID)
			qHandler.Delete(nextJobID)
			return
		}
		logr.WithError(err).Errorf("error retrieving training job")
		return
	}

	if trainingRecord.Deleted {
		logr.Debugf("job %s was deleted", nextJobID)
		qHandler.Delete(nextJobID)
		return
	}
	if trainingRecord.TrainingStatus.Status != grpc_storage.Status_QUEUED {
		logr.Warnf("job %s expected status QUEUED but found %s, removing job from queue", nextJobID, trainingRecord.TrainingStatus)
		qHandler.Delete(nextJobID)
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
				if dequeuedTrainingRecord.TrainingStatus.Status == grpc_storage.Status_QUEUED {
					_, err = updateTrainingJobPostLock(s, &grpc_storage.UpdateRequest{
						TrainingId:    dequeuedJob,
						UserId:        dequeuedTrainingRecord.UserID,
						Status:        grpc_storage.Status_FAILED,
						StatusMessage: "Job was dequeued without being submitted",
						ErrorCode:     storageClient.ErrCodeFailDequeue,
					})
					if err != nil {
						logr.WithError(err).Errorln("Unable to update job status to FAILED")
					}
				}
			}
		}
	}

	qHandler.Unlock()
	locked = false

	// store job state transition
	timestamp := storageClient.CurrentTimestampAsString()
	e := &JobHistoryEntry{
		TrainingID:    trainingRecord.TrainingID,
		Timestamp:     timestamp,
		Status:        trainingRecord.TrainingStatus.Status,
		StatusMessage: trainingRecord.TrainingStatus.StatusMessage,
		ErrorCode:     trainingRecord.TrainingStatus.ErrorCode,
	}
	s.jobHistoryRepo.RecordJobStatus(e)
}

func (s *storageService) CreateTrainingJob(ctx context.Context, req *grpc_storage.CreateRequest) (*grpc_storage.CreateResponse, error) {

	sidHigh := random.String(9, random.Lowercase+random.Numeric+"-")
	sidLow := random.String(9, random.Lowercase+random.Numeric)

	id := fmt.Sprintf("training-%s%s", sidHigh, sidLow)

	logr := logger.LocLogger(logWith(id, req.UserId))

	logr.Infof("CreateTrainingJob, req.TFosRequest: %v", req.TfosRequest)

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

	outputDatastore := s.getOutputDatastore(req.Training.OutputData, req.Datastores)

	// upload model definition ZIP file to object store and set location
	logr.Debugf("debug_1.5.1 start to s3")
	if req.ModelDefinition.Content != nil {
		logr.Debugf("debug_1.5.1 start to s3 is nil")
		// MLSS Change: upload model to object store no matter what output datastore is used
		//if outputDatastore.Type != "mount_volume" {
		// Upload to DLaaS Object store.
		logr.Debugf("Uploading content to: %s filename: %s", s.modelsBucket, getModelZipFileName(id))
		err := s.datastore.UploadArchive(s.modelsBucket, getModelZipFileName(id), req.ModelDefinition.Content)
		if err != nil {
			logr.WithError(err).Errorf("Error uploading model to object store")
			return nil, err
		}
		req.ModelDefinition.Location = fmt.Sprintf("%s/%s.zip", s.modelsBucket, id)
		cl.Observe("uploaded model to dlaas object store: %s", req.ModelDefinition.Location)

		// Upload to user's result object store.
		logr.Debugf("Upload to user's result object store, type: %s bucket: %s", outputDatastore.Type, outputDatastore.Fields["bucket"])
		ds, err := storage.CreateDataStore(outputDatastore.Type, outputDatastore.Connection)
		if err != nil {
			logr.WithError(err).Fatalf("Cannot create datastore for output data store %s", outputDatastore.Id)
			return nil, err
		}
		err = ds.Connect()
		if err != nil {
			logr.WithError(err).Fatalf("Cannot connect to output object store %s", outputDatastore.Id)
			return nil, err
		}
		defer ds.Disconnect()
		bucket := outputDatastore.Fields["bucket"]
		object := fmt.Sprintf("%s/_submitted_code/model.zip", id)
		logr.Debugf("Writing to output object store: %s -> %s, length: %d", bucket, object, len(req.ModelDefinition.Content))
		err = ds.UploadArchive(bucket, object, req.ModelDefinition.Content)
		if err != nil {
			logr.WithError(err).Errorf("Error uploading model to output object store")
			return nil, err
		}
		cl.Observe("uploaded model to user's object store")
		//}
	}

	// create a copy of the model definition without the content field (do not store it to the database)
	modelWithoutContent := *req.ModelDefinition
	modelWithoutContent.Content = nil

	// get evaluation metrics from create request
	evaluationMetricsSpec := ""
	if req.EvaluationMetrics != nil {
		logr.Debugf("EMExtractionSpec ImageTag: %s", req.EvaluationMetrics.ImageTag)
		wrapper := make(map[string]interface{})
		wrapper["evaluation_metrics"] = req.EvaluationMetrics
		data, err := yaml.Marshal(wrapper)
		if err != nil {
			logr.WithError(err).Errorf("Can't re-marshal evaluation metrics specification")
		}
		evaluationMetricsSpec = string(data)
		logr.Debugf("Set evaluation_metrics to: %s<eof>", evaluationMetricsSpec)
	}

	tr := &TrainingRecord{
		TrainingID:      id,
		UserID:          req.UserId,
		ModelDefinition: &modelWithoutContent,
		Training:        req.Training,
		Datastores:      req.Datastores,
		TrainingStatus: &grpc_storage.TrainingStatus{
			Status:              grpc_storage.Status_QUEUED,
			SubmissionTimestamp: storageClient.CurrentTimestampAsString(),
		},
		Metrics:               nil,
		EvaluationMetricsSpec: evaluationMetricsSpec,
		Namespace:             req.Namespace,
		// FIXME MLSS Change: add GID & UID from request to training records
		GID: req.Gid,
		UID: req.Uid,

		// FIXME MLSS Change: v_1.4.1
		JobAlert:     req.JobAlert,
		CodeSelector: req.CodeSelector,
		PSs:          req.Pss,
		PSCPU:        req.PsCpu,
		PSImage:      req.PsImage,
		PSMemory:     req.PsMemory,
		DataPath:     req.DataPath,
		JobType:      req.JobType,
		TFosRequest:  req.TfosRequest,
		ExpRunId:     req.ExpRunId,
		ExpName:      req.ExpName,
		FileName:     req.FileName,
		FilePath:     req.FilePath,
		ProxyUser:    req.ProxyUser,
	}
	logr.Debugf("CreateTrainingJob, req.JobType: %v, reqCpu: %v, reqGpu: %v, reqMem: %v", req.JobType, req.Training.Resources.Cpus, req.Training.Resources.Gpus, req.Training.Resources.Memory)

	err := s.submitJobToLCM(tr, logr)
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

	return &grpc_storage.CreateResponse{TrainingId: id}, nil
}

func (s *storageService) GetTrainingJob(ctx context.Context, req *grpc_storage.GetRequest) (*grpc_storage.GetResponse, error) {
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
	jobb := &grpc_storage.Job{
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
	return &grpc_storage.GetResponse{
		Job: jobb,
	}, nil
}

func (s *storageService) GetTrainingStatusID(ctx context.Context, req *grpc_storage.GetRequest) (*grpc_storage.GetStatusIDResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	statusID, err := s.repo.FindTrainingStatusID(req.TrainingId)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
		}
		logr.WithError(err).Errorf("Cannot retrieve record for training %s", req.TrainingId)
		return nil, err
	}
	return &grpc_storage.GetStatusIDResponse{
		Status: statusID,
	}, nil
}

func (s *storageService) UpdateTrainingJob(ctx context.Context, req *grpc_storage.UpdateRequest) (*grpc_storage.UpdateResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("UpdateTrainingJob called for training %s", req.TrainingId)

	return updateTrainingJobPostLock(s, req)
}

// This method contains all the functionality of UpdateTrainingJob, minus the lock on the database.  This enables it to be called
// from within another function, which already has the lock itself (Halt)
func updateTrainingJobPostLock(s *storageService, req *grpc_storage.UpdateRequest) (*grpc_storage.UpdateResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	training, err := s.repo.Find(req.TrainingId)
	if err != nil {
		logr.WithError(err).Errorf("Cannot retrieve training '%s'", req.TrainingId)
		return nil, err
	}
	if training == nil {
		// training does not exist
		logr.Errorf("training does not exist")
		return nil, gerrf(codes.NotFound, "Training with id %s not found.", req.TrainingId)
	}

	if training.UserID != req.UserId {
		msg := fmt.Sprintf("User %s does not have permission to update training data with id %s.",
			req.UserId, req.TrainingId)
		logr.Error(msg)
		return nil, gerrf(codes.PermissionDenied, "permission denied")
	}

	ts := training.TrainingStatus
	originalStatus := ts.Status

	// If status is completed/failed/halted and the update is requesting a halt, then do nothing and return error
	if originalStatus == grpc_storage.Status_COMPLETED || originalStatus == grpc_storage.Status_FAILED {
		logr.Error("Cannot update training status: %s", originalStatus)
		//Change to Response
		return &grpc_storage.UpdateResponse{TrainingId: training.TrainingID}, err
	}

	ts.Status = req.Status
	ts.StatusMessage = req.StatusMessage
	ts.ErrorCode = req.ErrorCode

	nowMillis := storageClient.CurrentTimestampAsString()

	if req.Status == grpc_storage.Status_COMPLETED || req.Status == grpc_storage.Status_FAILED {
		// FIXME MLSS Change: keep datastores even job finished
		//training.Datastores = nil
		ts.CompletionTimestamp = nowMillis
		if req.Timestamp != "" {
			ts.CompletionTimestamp = req.Timestamp
		}
		// erase sensitive data from the db
		training.ModelDefinition.Framework.ImageLocation = nil
	}
	if req.Status == grpc_storage.Status_PROCESSING {
		ts.ProcessStartTimestamp = nowMillis
		if req.Timestamp != "" {
			ts.ProcessStartTimestamp = req.Timestamp
		}
	}

	// send monitoring metrics for failed/succeeded jobs
	if req.Status == grpc_storage.Status_COMPLETED || req.Status == grpc_storage.Status_FAILED {
		if req.Status == grpc_storage.Status_FAILED {
			errorType := "server"
			if strings.HasPrefix(req.ErrorCode, "C") {
				errorType = "client"
			}

			logFrameworkErrorsValue := fmt.Sprintf("%s-%s-%s", training.ModelDefinition.Framework.Name, errorType, req.ErrorCode)
			logr.WithFields(logrus.Fields{
				"framework_errors": logFrameworkErrorsValue,
			}).Debug(" metrics for failed training jobs framework")

		}
		gpusUsed := training.Training.Resources.Gpus
		if training.Training.Resources.Learners > 1 {
			gpusUsed = training.Training.Resources.Gpus * float32(training.Training.Resources.Learners)
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
	logr.Debugf("CHECKING Stored training %s, Status %s Code %s Message %s", req.TrainingId, ts.Status, ts.ErrorCode, ts.StatusMessage)

	// Additionally, store any job state transitions in the job_history DB collection
	// We store a history record if either (1) the status is different, or (2) if this is
	// a PROCESSING->PROCESSING transition, to record the full picture for distributed jobs.
	if req.Status != originalStatus || req.Status == grpc_storage.Status_PROCESSING {
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

	return &grpc_storage.UpdateResponse{TrainingId: training.TrainingID}, nil
}

func (s *storageService) GetTrainingJobByTrainingJobId(ctx context.Context, req *grpc_storage.GetTrainingJobRequest) (*grpc_storage.GetResponse, error) {
	logger.GetLogger().Info("Get training job, id:", req.TrainingId)
	tr, err := s.repo.Find(req.TrainingId)
	if err != nil {
		logger.GetLogger().Error("Find training err, ", err)
		return nil, errors.New(fmt.Sprintf("Find training err, %v", err))
	}

	newJob := &grpc_storage.Job{
		UserId:          tr.UserID,
		JobId:           tr.JobID,
		ModelDefinition: tr.ModelDefinition,
		TrainingId:      tr.TrainingID,
		Training:        tr.Training,
		Status:          tr.TrainingStatus,
		Datastores:      tr.Datastores,
		JobNamespace:    tr.Namespace,
		JobAlert:        tr.JobAlert,
		Pss:             tr.PSs,
		PsCpu:           tr.PSCPU,
		PsImage:         tr.PSImage,
		PsMemory:        tr.PSMemory,
		JobType:         tr.JobType,
		ExpRunId:        tr.ExpRunId,
		ExpName:         tr.ExpName,
		FileName:        tr.FileName,
		FilePath:        tr.FilePath,
		Metrics:         tr.Metrics,
		CodeSelector:    tr.CodeSelector,
		TfosRequest:     tr.TFosRequest,
		ProxyUser:       tr.ProxyUser,
		DataSet:         tr.DataSet,
		MfModel:         tr.MFModel,
	}
	return &grpc_storage.GetResponse{
		Job: newJob,
	}, nil
}

func (s *storageService) GetAllTrainingsJobs(ctx context.Context, req *grpc_storage.GetAllRequest) (*grpc_storage.GetAllResponse, error) {
	logr := logger.LocLogger(logEntry().WithField(logger.LogkeyUserID, req.UserId))
	logr.Debugf("GetAllTrainingsJobs called")

	cl := instrumentation.NewCallLogger(ctx, "GetAllTrainingsJobs", logr)
	defer cl.Returned()

	jobPaged, err := s.repo.FindAll(req.UserId, req.Page, req.Size, req.ExpRunId)
	if err != nil {
		msg := "Failed to retrieve all training jobPaged"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_storage.GetAllResponse{
		Jobs: make([]*grpc_storage.Job, len(jobPaged.TrainingRecords)),
	}
	for i, job := range jobPaged.TrainingRecords {
		resp.Jobs[i] = &grpc_storage.Job{
			UserId:          job.UserID,
			JobId:           job.JobID,
			ModelDefinition: job.ModelDefinition,
			TrainingId:      job.TrainingID,
			Training:        job.Training,
			Status:          job.TrainingStatus,
			Datastores:      job.Datastores,
			ExpRunId:        job.ExpRunId,
			ExpName:         job.ExpName,
			FileName:        job.FileName,
			FilePath:        job.FilePath,
			CodeSelector:    job.CodeSelector,
			TfosRequest:     job.TFosRequest,
			ProxyUser:       job.ProxyUser,
		}
	}

	pageStr := strconv.Itoa(jobPaged.Pages)
	resp.Pages = pageStr
	totalStr := strconv.Itoa(jobPaged.Total)
	resp.Total = totalStr

	return resp, nil
}

func (s *storageService) GetAllTrainingsJobsByUserIdAndNamespace(ctx context.Context, req *grpc_storage.GetAllRequest) (*grpc_storage.GetAllResponse, error) {
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
	expRunId := req.ExpRunId
	logr.Debugf("storageService GetAllTrainingsJobsByUserIdAndNamespace, Username: %v, Namespace: %v, Page: %v, Size: %v, expRunId: %v", username, namespace, page, size, expRunId)

	//jobs, err := s.repo.FindAll(req.UserId)
	jobs, err := s.repo.FindAllByUserIdAndNamespace(&username, &namespace, page, size, expRunId)
	if err != nil {
		msg := "Failed to retrieve all training jobs"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_storage.GetAllResponse{
		Jobs: make([]*grpc_storage.Job, len(jobs.TrainingRecords)),
	}
	for i, job := range jobs.TrainingRecords {
		logr.Debugf("trainer_impl GetAllTrainingsJobsByUserIdAndNamespace job.JobAlert: %v", job.JobAlert)
		resp.Jobs[i] = &grpc_storage.Job{
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
			FileName:     job.FileName,
			FilePath:     job.FilePath,
			CodeSelector: job.CodeSelector,
			TfosRequest:  job.TFosRequest,
			ProxyUser:    job.ProxyUser,
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
func (s *storageService) GetAllTrainingsJobsByUserIdAndNamespaceList(ctx context.Context, req *grpc_storage.GetAllRequest) (*grpc_storage.GetAllResponse, error) {
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
	expRunId := req.ExpRunId
	clusterName := req.ClusterName
	logr.Debugf("storageService GetAllTrainingsJobsByUserIdAndNamespaceList, Username: %v, NamespaceList: %v, Page: %v, Size: %v, ExpRunId: %v", username, namespaceList, page, size, expRunId)

	//jobs, err := s.repo.FindAll(req.UserId)
	jobs, err := s.repo.FindAllByUserIdAndNamespaceList(&username, &namespaceList, page, size, clusterName, expRunId)
	if err != nil {
		msg := "Failed to retrieve all training jobs"
		logr.WithError(err).Errorf(msg)
		return nil, gerrf(codes.Internal, msg)
	}
	resp := &grpc_storage.GetAllResponse{
		Jobs: make([]*grpc_storage.Job, len(jobs.TrainingRecords)),
	}
	for i, job := range jobs.TrainingRecords {
		//logr.Debugf("debug_for_param ParameterServer trainer pss: %v, psCpu: %v, psImage: %v, psMemory: %v", job.PSs, job.PSCPU, job.PSImage, job.PSMemory)
		resp.Jobs[i] = &grpc_storage.Job{
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
			FileName:     job.FileName,
			FilePath:     job.FilePath,
			CodeSelector: job.CodeSelector,
			TfosRequest:  job.TFosRequest,
			ProxyUser:    job.ProxyUser,
		}
		logr.Debugf("storageService GetAllTrainingsJobsByUserIdAndNamespaceList job[%d]: %+v", i, resp.Jobs[i])
	}
	logr.Debugf("storageService GetAllTrainingsJobsByUserIdAndNamespaceList length: %d", len(resp.Jobs))


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

func (s *storageService) deleteJobFromTDS(query *tdsService.Query, logr *logger.LocLoggingEntry) error {
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

func (s *storageService) deleteJobFromQueue(trainingID string, namespaceName string, logr *logger.LocLoggingEntry) error {
	qHandler := s.queues[QueueName(namespaceName)]
	if qHandler == nil {
		// FIXME MLSS Change: v_1.4.1 to create queue  when queue is nil
		queue, err := newTrainingJobQueue(viper.GetString(mongoAddressKey),
			viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey),
			viper.GetString(mongoPasswordKey), viper.GetString(mongoAuthenticationDatabase),
			config.GetMongoCertLocation(), QueueName(namespaceName), LockName(namespaceName))
		if err != nil {
			logr.WithError(err).Fatalf("Cannot create queue with %s %s %s", viper.GetString(mongoAddressKey), viper.GetString(mongoDatabaseKey), viper.GetString(mongoUsernameKey))
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
	}
	return nil
}

func (s *storageService) DeleteTrainingJob(ctx context.Context,
	req *grpc_storage.DeleteRequest) (*grpc_storage.DeleteResponse, error) {

	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))

	cl := instrumentation.NewCallLogger(ctx, "DeleteTrainingJob", logr)
	defer cl.Returned()

	readResp, err := s.GetTrainingJob(ctx, &grpc_storage.GetRequest{
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

	var job *grpc_storage.Job
	if readResp != nil {
		job = readResp.Job

		// delete from queue
		if job.Status.Status == grpc_storage.Status_QUEUED {
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

		// delete from DB
		err = s.repo.Delete(job.TrainingId)
		if err != nil {
			logr.WithError(err).Errorf("Failed to delete training job '%s' from database", job.TrainingId)
			return nil, err
		}
		cl.Observe("deleted model from mongo")

		//return &grpc_storage.DeleteResponse{TrainingId: job.JobId}, nil
		return &grpc_storage.DeleteResponse{Success: true}, nil
	}
	return nil, gerrf(codes.NotFound, "Training with id '%s' not found.", req.TrainingId)
}

func (s *storageService) ResumeTrainingJob(ctx context.Context, req *grpc_storage.ResumeRequest) (*grpc_storage.ResumeResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("HaltTrainingJob called")
	return nil, gerrf(codes.Unimplemented, "ResumeTrainingJob not implemented yet")
}

func (s *storageService) GetModelDefinition(req *grpc_storage.ModelDefinitionRequest, stream grpc_storage.Storage_GetModelDefinitionServer) error {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("GetModelDefinition")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.GetTrainingJob(ctx, &grpc_storage.GetRequest{
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
	err = stream.Send(&grpc_storage.ZippedDataChunk{
		Data: payload,
	})
	if err != nil {
		logr.WithError(err).Errorf("Failed to send zipped chunk.")
		return err
	}
	return nil
}

func (s *storageService) GetTrainedModel(req *grpc_storage.TrainedModelRequest, stream grpc_storage.Storage_GetTrainedModelServer) error {
	//s.mtx.Lock()
	//defer s.mtx.Unlock()

	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("GetTrainedModel")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.GetTrainingJob(ctx, &grpc_storage.GetRequest{
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
		err = stream.Send(&grpc_storage.ZippedDataChunk{
			Data: buf,
		})
		if err != nil {
			logr.WithError(err).Error("Failed to send zipped data chunk")
			return err
		}
	}
	return nil
}

func trainedModelLogRequestToTrainerQuery(req *grpc_storage.TrainedModelLogRequest, rindex int64, pageSize int32) *tdsService.Query {
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

func trainedGetModelLogRequestToTrainerQuery(req *grpc_storage.GetTrainedModelLogRequest, size int64, from int64) *tdsService.Query {
	query := &tdsService.Query{
		Meta: &tdsService.MetaInfo{
			TrainingId: req.TrainingId,
			UserId:     req.UserId,
		},
		StartToEndLine: strconv.Itoa(int(size)) + ":" + strconv.Itoa(int(from)),
		Pagesize:       1,
	}
	return query
}

func (s *storageService) isLearningFinished(req *grpc_storage.TrainedModelLogRequest) (bool, error) {
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
	if statusID == grpc_storage.Status_COMPLETED ||
		statusID == grpc_storage.Status_FAILED {
		jobCompleted = true
	}

	return jobCompleted, nil
}

func (s *storageService) waitUntilJobStart(req *grpc_storage.TrainedModelLogRequest,
	outStream grpc_storage.Storage_GetTrainedModelLogsServer,
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
			if !(statusID == grpc_storage.Status_NOT_STARTED ||
				statusID == grpc_storage.Status_QUEUED ||
				statusID == grpc_storage.Status_PENDING) {
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
				errSend := outStream.Send(&grpc_storage.ByteStreamResponse{Data: []byte(msg)})
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

func (s *storageService) GetTrainedModelLogs(req *grpc_storage.TrainedModelLogRequest,
	outStream grpc_storage.Storage_GetTrainedModelLogsServer) error {

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
			errSend := outStream.Send(&grpc_storage.ByteStreamResponse{Data: []byte(chunk.Line)})
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

func (s *storageService) GetTrainedModelLog(ctx context.Context, req *grpc_storage.GetTrainedModelLogRequest) (*grpc_storage.GetTrainedModelLogResponse, error) {
	tds, err := s.tdsClient()
	if err != nil {
		logger.GetLogger().Error("Cannot create LCM service client")
		return nil, err
	}
	query := trainedGetModelLogRequestToTrainerQuery(req, req.Size, req.From)
	res, err := tds.Client().GetLog(ctx, query)
	if err != nil {
		logger.GetLogger().Error("get logs error, ", err)
		return nil, err
	}
	return &grpc_storage.GetTrainedModelLogResponse{
		Logs: res.Log,
	}, nil
}

func marshalQuerySearchType(st grpc_storage.Query_SearchType) tdsService.Query_SearchType {
	searchType := tdsService.Query_TERM

	switch st {
	case grpc_storage.Query_TERM:
		searchType = tdsService.Query_TERM
		break
	case grpc_storage.Query_NESTED:
		searchType = tdsService.Query_NESTED
		break
	case grpc_storage.Query_MATCH:
		searchType = tdsService.Query_MATCH
		break
	case grpc_storage.Query_ALL:
		searchType = tdsService.Query_ALL
		break
	}
	return searchType
}

func marshalTDSQueryToTrainerQuery(in *grpc_storage.Query) *tdsService.Query {
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

func (s *storageService) GetTrainingLogs(in *grpc_storage.Query,
	outStream grpc_storage.Storage_GetTrainingLogsServer) error {

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
		errSend := outStream.Send(&grpc_storage.LogLine{
			Meta: &grpc_storage.MetaInfo{
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

func (s *storageService) GetVersions(ctx context.Context, req *grpc_storage.GetVersionsRequest) (*grpc_storage.Frameworks, error) {
	//call the frameworks.go and then getAllVersions for the frameworks
	//Return response from getAll Versions
	frameworks, err := getExternalVersions()
	if err != nil {
		return nil, err
	}
	return &frameworks, nil
}

func (s *storageService) validateRequest(log *logrus.Entry, req *grpc_storage.CreateRequest) error {
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

	return nil
}

func findDatastore(id string, ds []*grpc_storage.Datastore) *grpc_storage.Datastore {
	for _, v := range ds {
		if v.Id == id {
			return v
		}
	}
	return nil
}


func (s *storageService) failCreateRequest(msg string, req *grpc_storage.CreateRequest, log *logrus.Entry) error {
	return s.failCreateRequestWithCode(storageClient.ErrInvalidManifestFile, msg, req, log)
}

func (s *storageService) failCreateRequestWithCode(errorCode string, msg string, req *grpc_storage.CreateRequest, log *logrus.Entry) error {
	log.Errorf("Failed to validate CreateRequest: %s", msg)

	// send error event as monitoring metric

	if req.ModelDefinition != nil && req.ModelDefinition.Framework != nil {

		logFrameworkErrorsValue := fmt.Sprintf("%s-%s-%s", req.ModelDefinition.Framework.Name, "client", errorCode)
		log.WithFields(logrus.Fields{
			"framework_errors": logFrameworkErrorsValue,
		})
	}

	log.Debug("Metrics for failed training jobs framework")

	return gerrf(codes.InvalidArgument, msg)
}

func (s *storageService) validateDatastore(ds *grpc_storage.Datastore, req *grpc_storage.CreateRequest, log *logrus.Entry) error {

	if ds == nil {
		return s.failCreateRequest("Data store is not set", req, log)
	}
	log.Debugf("Datastore type: %s", ds.Type)
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

		ostore, err := storage.CreateDataStore(ds.Type, ds.Connection)
		if err != nil {
			log.Errorf("Validation failed: %s", err.Error())
			return s.failCreateRequestWithCode(storageClient.ErrInvalidCredentials,
				fmt.Sprintf("Data store authentication information for id '%s' incorrect or there is a connection problem", ds.Id), req, log)
		}

		if err := ostore.Connect(); err != nil {
			log.Errorf("Validation failed: %s", err.Error())
			return s.failCreateRequestWithCode(storageClient.ErrInvalidCredentials,
				fmt.Sprintf("Data store authentication information for id '%s' incorrect or there is a connection problem", ds.Id), req, log)
		}

		// validate bucket (or container as it is called in Swift)
		bucket := ds.Fields["bucket"]
		if bucket != "" {
			exists, err := ostore.ContainerExists(bucket)
			if !exists || err != nil {
				return s.failCreateRequestWithCode(storageClient.ErrInvalidCredentials,
					fmt.Sprintf("Data store bucket '%s' for data store id '%s' incorrect, there may be a connection problem or credentials do not allow access to the bucket", bucket, ds.Id), req, log)
			}
		}
	}
	return nil
}

// lcmClient established a connection if the storageService has nothing existing cached
func (s *storageService) lcmClient() (client.LcmClient, error) {
	if s.lcm == nil {
		return client.NewLcm(nil)
	}
	return s.lcm, nil
}

func (s *storageService) tdsClient() (tdsClient.TrainingDataClient, error) {
	if s.tds == nil {
		logger.GetLogger().Debugln("storageService tdsClient is nil", )
		address := fmt.Sprintf("ffdl-trainingdata.%s.svc.cluster.local:80", config.GetPodNamespace())
		tds, err := tdsClient.NewTrainingDataClientWithAddress(address)
		if err != nil {
			return nil, err
		}
		s.tds = tds
	}
	logger.GetLogger().Debugf("storageService create tdsClient pod namespace: %v\n", config.GetPodNamespace())
	return s.tds, nil
}

func (s *storageService) createJobConfig(tr *TrainingRecord) (*service.JobDeploymentRequest, error) {
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
	envvars["USER_ID"] = tr.UserID

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

	var dataStores []*grpc_trainer_v2.Datastore
	var modelDefinition grpc_trainer_v2.ModelDefinition
	var training grpc_trainer_v2.Training
	var trainingStatus grpc_trainer_v2.TrainingStatus
	var mfModel grpc_trainer_v2.MFModel
	var dataSet grpc_trainer_v2.DataSet
	err = copier.Copy(&dataStores, &tr.Datastores)
	if err != nil {
		logr.Debugf("Copy Datastore failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy Datastore: %+v", dataStores)

	err = copier.Copy(&modelDefinition, &tr.ModelDefinition)
	if err != nil {
		logr.Debugf("Copy ModelDefinition failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy ModelDefinition: %+v", dataStores)

	err = copier.Copy(&training, &tr.Training)
	if err != nil {
		logr.Debugf("Copy Training failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy Training: %+v", dataStores)

	err = copier.Copy(&trainingStatus, &tr.TrainingStatus)
	if err != nil {
		logr.Debugf("Copy TrainingStatus failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy TrainingStatus: %+v", dataStores)

	err = copier.Copy(&mfModel, &tr.MFModel)
	if err != nil {
		logr.Debugf("Copy MFModel failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy MFModel: %+v", mfModel)

	err = copier.Copy(&dataSet, &tr.DataSet)
	if err != nil {
		logr.Debugf("Copy MFModel failed: %v", err.Error())
		return nil, err
	}
	logr.Debugf("Copy MFModel: %+v", dataSet)

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
		DataStores:      dataStores,
		ModelDefinition: &modelDefinition,
		Training:        &training,
		TrainingStatus:  &trainingStatus,
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
		ExpName:      tr.ExpName,
		ExpRunId:     tr.ExpRunId,
		FileName:     tr.FileName,
		FilePath:     tr.FilePath,
		DataSet:      &dataSet,
		MfModel:      &mfModel,
		Algorithm:    tr.Algorithm,
		JobParams:    tr.JobParams,
		FitParams:    tr.FitParams,
		APIType:      tr.APIType,
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

func setDefaultResourceRequirements(t *grpc_storage.Training) {
	if t == nil || t.Resources == nil {
		t.Resources = &grpc_storage.ResourceRequirements{ // set sensible defaults
			Cpus:        5.0,
			Gpus:        1.0,
			Memory:      12.0,
			MemoryUnit:  grpc_storage.SizeUnit_GiB,
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
		t.Resources.MemoryUnit = grpc_storage.SizeUnit_GiB
	}
	if t.Resources.Schedpolicy == "" || strings.ToLower(t.Resources.Schedpolicy) != "spread" {
		t.Resources.Schedpolicy = "dense"
	}
	if t.Resources.GpuType == "" {
		t.Resources.GpuType = "nvidia-TeslaK80"
	}
}

func getResourceRequirements(t *grpc_storage.Training) *service.ResourceRequirements {
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

// getOutputDatastore retrieves the output data store or return the internal datastore if none has been defined
func (s *storageService) getOutputDatastore(outputData []string, datastores []*grpc_storage.Datastore) *grpc_storage.Datastore {
	var ds *grpc_storage.Datastore
	if len(outputData) > 0 {
		ds = findDatastore(outputData[0], datastores) // we assume there is only one output data at this point b/c the underlying system does not support more
	}
	if ds == nil {
		ds = &grpc_storage.Datastore{
			Id:         internalObjectStoreID,
			Type:       config.GetDataStoreType(),
			Connection: config.GetDataStoreConfig(),
			Fields:     map[string]string{"bucket": s.trainedModelsBucket},
		}
	}
	return ds
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

func (s *storageService) submitJobToLCM(tr *TrainingRecord, logr *logger.LocLoggingEntry) error {
	logr.Printf("submitJobToLCM, tr.JobType: %v", tr.JobType)

	jobConfig, err := s.createJobConfig(tr)
	if err != nil {
		logr.WithError(err).Errorf("Failed to create job config")
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}

	lcm, err := s.lcmClient()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create LCM service client")
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}
	defer lcm.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = lcm.Client().DeployTrainingJob(ctx, jobConfig)
	if err != nil {
		logr.WithError(err).Errorf("Cannot deploy training job with id %s", tr.TrainingID)
		return gerrf(codes.Internal, grpcErrorDesc(err))
	}

	logr.Printf("training job %s submitted to lcm", tr.TrainingID)

	// capture the gpu usage when the job is submitted to LCM
	gpusUsed := tr.Training.Resources.Gpus
	if tr.Training.Resources.Learners > 1 {
		gpusUsed = tr.Training.Resources.Gpus * float32(tr.Training.Resources.Learners)
	}
	// log the gpu usages requested by user
	logGpuTypeIncrementValue := fmt.Sprintf("%s-%v", tr.Training.Resources.GpuType, gpusUsed)
	logr.WithFields(logrus.Fields{
		"gputype_increment": logGpuTypeIncrementValue,
	}).Debug(" incrementing the gpus")

	// increment the counter

	return nil
}

// determine if this job should be rate-limited by checking if the total number of GPUs
// would exceed the limit set for the GPU type
func (s *storageService) rateLimitTrainingJob(trainingRecord *TrainingRecord, logr *logger.LocLoggingEntry) bool {
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
		gpusRequested = trainingRecord.Training.Resources.Gpus * float32(trainingRecord.Training.Resources.Learners)
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
	var totalGPUsUsedCount float32
	var matchingGPUConsumingRecords []*TrainingRecord
	for _, record := range records {
		trainingStatus := record.TrainingStatus.Status
		if trainingStatus == grpc_storage.Status_COMPLETED || trainingStatus == grpc_storage.Status_FAILED || trainingStatus == grpc_storage.Status_QUEUED {
			//ignore these since they don't count towards active resource usage
		} else if TransformResourceName(record.Training.Resources.GpuType) == TransformResourceName(gpuType) {
			//only count matching gpu type
			matchingGPUConsumingRecords = append(matchingGPUConsumingRecords, record)
			gpusUsed := record.Training.Resources.Gpus
			if record.Training.Resources.Learners > 1 {
				gpusUsed = record.Training.Resources.Gpus * float32(record.Training.Resources.Learners)
			}
			totalGPUsUsedCount = totalGPUsUsedCount + gpusUsed
		}
	}

	if int64(totalGPUsUsedCount+gpusRequested) > limit {
		rateLimit = true
		if logr.Logger.Level >= logrus.DebugLevel {
			for _, record := range matchingGPUConsumingRecords {
				logr.Debugf("Found a gpu consuming training %v has a status %s and using gpus %v with submission time as %v and process start time as %v and code %v",
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
func (s *storageService) DeleteSubmittedCode(ctx context.Context, req *grpc_storage.DeleteRequest) (*grpc_storage.DeleteResponse, error) {
	logr := logger.LocLogger(logWith(req.TrainingId, req.UserId))
	logr.Debugf("DeleteSubmittedCode enter, TrainingId: %s,UserId: %s", req.TrainingId, req.UserId)

	bucket := GetModelsBucket()
	TrainingId := req.TrainingId
	zipName := TrainingId + ".zip"

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()
	//resp, err := s.GetTrainingJob(ctx, &grpc_storage.GetRequest{
	//	TrainingId: req.TrainingId,
	//	UserId:     req.UserId,
	//})
	//if err != nil {
	//	logr.WithError(err).Errorf("Error reading training with id: %s", req.TrainingId)
	//	return nil, err
	//}
	//if resp == nil || resp.Job == nil {
	//	//return gerrf(codes.NotFound, "Training with id '%s' not found.", req.TrainingId)url
	//}

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
	//return &grpc_storage.DeleteResponse{TrainingId: TrainingId}, nil
	return &grpc_storage.DeleteResponse{Success: true}, nil
}

// AddEMetrics adds the passed evaluation metrics record to storage.
func (s *storageService) AddEMetrics(ctx context.Context, in *grpc_storage.EMetrics) (*grpc_storage.AddResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)
	s.reportOnCluster("AddEMetrics", logr)

	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)

	start := time.Now()

	jsonBytes, err := json.Marshal(in)
	//noinspection GoBoolExpressions
	if TdsDebugMode {
		dlogr.Debugf("emetrics_in_json: %d: %d: %s",
			in.Meta.Rindex, in.Meta.Time, makeSnippetForDebug(string(jsonBytes), 20))
	}

	_, err = s.es.Index().
		Index(indexName).
		Type(docTypeEmetrics).
		BodyString(string(jsonBytes)).
		Do(ctx)

	doneQuery := time.Now()

	out := new(grpc_storage.AddResponse)
	if err != nil {
		logr.WithError(err).Errorf("Failed to add log line to elasticsearch")
		out.Success = false
		return out, err
	}

	out.Success = true

	s.reportTime(logr, "AddEMetrics", in.Meta.TrainingId, in.Meta.Rindex, in.Meta.Time, start, doneQuery)

	return out, nil
}

// AddLogLine adds the line line record to storage.
func (s *storageService) AddLogLine(ctx context.Context, in *grpc_storage.LogLine) (*grpc_storage.AddResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	s.reportOnCluster("AddLogLine", logr)

	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)
	dlogr.Debugf("AddLogLine: %d (%s): %s", in.Meta.Rindex, in.Meta.Subid, makeSnippetForDebug(in.Line, 7))

	start := time.Now()

	out := new(grpc_storage.AddResponse)

	_, err := s.es.Index().
		Index(indexName).
		Type(docTypeLog).
		BodyJson(in).
		Do(ctx)

	doneQuery := time.Now()

	if err != nil {
		logr.WithError(err).Errorf("Failed to add log line to elasticsearch")
		out.Success = false
		return out, err
	}

	out.Success = true

	s.reportTime(logr, "AddLogLine", in.Meta.TrainingId, in.Meta.Rindex, in.Meta.Time, start, doneQuery)
	dlogr.Debugf("exit")
	return out, nil
}

// AddEMetricsBatch adds the passed evaluation metrics record to storage.
//noinspection GoBoolExpressions
func (s *storageService) AddEMetricsBatch(ctx context.Context,
	inBatch *grpc_storage.EMetricsBatch) (*grpc_storage.AddResponse, error) {

	out := new(grpc_storage.AddResponse)
	if len(inBatch.Emetrics) == 0 {
		out.Success = true
		return out, nil
	}
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, inBatch.Emetrics[0].Meta.TrainingId).
		WithField(logger.LogkeyUserID, inBatch.Emetrics[0].Meta.UserId)
	s.reportOnCluster("AddEMetrics", logr)

	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)

	if TdsDebugEMetricAdd {
		fmt.Printf("------------------\n")
	}
	start := time.Now()
	//bulkRequest := c.es.Bulk()
	bulkRequest := s.es.Bulk().Index(indexName).Type(docTypeEmetrics)
	for _, in := range inBatch.Emetrics {

		jsonBytes, err := json.Marshal(in)
		if err != nil {
			logr.WithError(err).Errorf("Could not marshal request to string: %+v", *in)
			out.Success = false
			return out, err
		}
		if TdsDebugEMetricAdd {
			var jsonBytes []byte
			//noinspection GoBoolExpressions
			fmt.Printf("emetrics_in_json: %d: %d: %s\n",
				in.Meta.Rindex, in.Meta.Time, makeSnippetForDebug(string(jsonBytes), 20))
		}

		r := es.NewBulkIndexRequest().
			Index(indexName).
			Type(docTypeEmetrics).
			Doc(string(jsonBytes))

		bulkRequest = bulkRequest.Add(r)

	}
	bulkResponse, err := bulkRequest.Refresh("wait_for").Do(ctx)
	if err != nil {
		logr.WithError(err).Error("bulkRequest.Refresh returned error")
	}
	if bulkResponse == nil {
		logr.Warning("expected bulkResponse to be != nil; got nil")
	}
	s.refreshIndexes(ctx, dlogr)

	doneQuery := time.Now()
	if TdsDebugEMetricAdd {
		fmt.Printf("------------------\n")
	}

	//if inBatch.Force {
	//	c.esBulkProcessor.Flush()
	//}

	out.Success = true

	s.reportTime(logr, "AddEMetrics", inBatch.Emetrics[0].Meta.TrainingId,
		inBatch.Emetrics[0].Meta.Rindex, inBatch.Emetrics[0].Meta.Time, start, doneQuery)

	return out, nil
}

// AddLogLineBatch adds the line line record to storage.
//noinspection GoBoolExpressions
func (s *storageService) AddLogLineBatch(ctx context.Context,
	inBatch *grpc_storage.LogLineBatch) (*grpc_storage.AddResponse, error) {

	out := new(grpc_storage.AddResponse)
	if len(inBatch.LogLine) == 0 {
		out.Success = true
		return out, nil
	}
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	s.reportOnCluster("AddLogLine", logr)

	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)

	var err error
	if TdsDebugLogLineAdd {
		fmt.Printf("------------------\n")
	}
	start := time.Now()
	prevTimeStamp := int64(0)
	bulkRequest := s.es.Bulk().Index(indexName).Type(docTypeEmetrics)
	for _, in := range inBatch.LogLine {
		if TdsDebugLogLineAdd {
			if prevTimeStamp == 0 {
				prevTimeStamp = in.Meta.Time
			}

			timeSincePrev := in.Meta.Time - prevTimeStamp
			fmt.Printf("AddLogLineBatch: %d (tmSncPrev: %d)(%s): %s\n", in.Meta.Rindex, int(timeSincePrev),
				in.Meta.Subid, makeSnippetForDebug(in.Line, 7))

			prevTimeStamp = in.Meta.Time
		}
		jsonBytes, err := json.Marshal(in)
		if err != nil {
			logr.WithError(err).Errorf("Could not marshal request to string: %+v", *in)
			out.Success = false
			return out, err
		}

		r := es.NewBulkIndexRequest().
			Index(indexName).
			Type(docTypeLog).
			Doc(string(jsonBytes))

		bulkRequest.Add(r)
	}
	bulkResponse, err := bulkRequest.Refresh("wait_for").Do(ctx)
	if err != nil {
		logr.WithError(err).Error("bulkRequest.Refresh returned error")
	}
	if bulkResponse == nil {
		logr.Warning("expected bulkResponse to be != nil; got nil")
	}
	s.refreshIndexes(ctx, dlogr)

	doneQuery := time.Now()
	if TdsDebugLogLineAdd {
		fmt.Printf("------------------\n")
	}
	out.Success = true

	s.reportTime(logr, "AddLogLine", inBatch.LogLine[0].Meta.TrainingId,
		inBatch.LogLine[0].Meta.Rindex, inBatch.LogLine[0].Meta.Time, start, doneQuery)

	return out, err
}

// DeleteEMetrics deletes the queried evaluation metrics from storage.
func (s *storageService) DeleteEMetrics(ctx context.Context, in *grpc_storage.Query) (*grpc_storage.DeleteResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)
	logr.Debugf("function entry")
	out := new(grpc_storage.DeleteResponse)

	query, _, err := makeESQueryFromDlaasQuery(in)
	if err != nil {
		logr.WithError(err).Errorf("Could not make elasticsearch query from submitted query: %+v", *in)
		out.Success = false
		return out, err
	}

	_, err = s.es.DeleteByQuery(indexName).
		Index(indexName).
		Type(docTypeEmetrics).
		Query(query).
		Do(ctx)

	if err != nil {
		logr.WithError(err).Errorf("Search failed failed")
		out.Success = false
		return out, err
	}

	out.Success = true
	logr.Debugf("function exit")
	return out, err
}

// DeleteLogLines deletes the queried log lines from storage.
func (s *storageService) DeleteLogLines(ctx context.Context, in *grpc_storage.Query) (*grpc_storage.DeleteResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)
	logr.Debugf("function entry")
	out := new(grpc_storage.DeleteResponse)

	query, _, err := makeESQueryFromDlaasQuery(in)
	if err != nil {
		logr.WithError(err).Errorf("Could not make elasticsearch query from submitted query: %+v", *in)
		out.Success = false
		return out, err
	}

	_, err = s.es.DeleteByQuery(indexName).
		Index(indexName).
		Type(docTypeLog).
		Query(query).
		Do(ctx)

	if err != nil {
		logr.WithError(err).Errorf("Search failed failed")
		out.Success = false
		return out, err
	}

	out.Success = true
	logr.Debugf("function exit")
	return out, err
}

// DeleteJob deletes both queried evaluation metrics and log lines.
func (s *storageService) DeleteJob(ctx context.Context, in *grpc_storage.Query) (*grpc_storage.DeleteResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)
	logr.Debugf("function entry")
	out := new(grpc_storage.DeleteResponse)

	query, _, err := makeESQueryFromDlaasQuery(in)
	if err != nil {
		logr.WithError(err).Errorf("Could not make elasticsearch query from submitted query: %+v", *in)
		out.Success = false
		return out, err
	}

	_, err = s.es.DeleteByQuery(indexName).
		Index(indexName).
		Query(query).
		Do(ctx)

	if err != nil {
		logr.WithError(err).Errorf("Search failed failed")
		out.Success = false
		return out, err
	}

	out.Success = true
	logr.Debugf("function exit")
	return out, err
}

func createIndexWithLogsIfDoesNotExist(ctx context.Context, client *es.Client) error {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	logr.Debugf("function entry")

	mainIndex := indexNameV1

	logr.Infof("calling IndexExists for %s", mainIndex)

	exists, err := client.IndexExists(mainIndex).Do(ctx)
	if err != nil {
		logr.WithError(err).Errorf("IndexExists for %s failed", mainIndex)
	}

	if exists {
		// ignore error if already exist
		client.Alias().Add(mainIndex, indexName).Do(ctx)
		if err != nil {
			logr.WithError(err).Infof("alias alias for %s failed", mainIndex)
		}

		logr.Infof("Maintaining index: %s", mainIndex)
		return nil
	}

	logr.Debugf("calling CreateIndex")
	ires, err := client.CreateIndex(mainIndex).Do(ctx)
	if err != nil {
		logr.WithError(err).Debug("CreateIndex failed")
		return err
	}
	if !ires.Acknowledged {
		return errors.New("the put mapping was not acknowledged")
	}
	res, err := client.PutMapping().Index(mainIndex).Type(docTypeLog).BodyString(indexMappingLogs).Do(ctx)
	if err != nil {
		logr.WithError(err).Debug("PutMapping logs failed")
		return err
	}
	if !res.Acknowledged {
		return errors.New("the put mapping was not acknowledged")
	}
	res, err = client.PutMapping().Index(mainIndex).Type(docTypeEmetrics).BodyString(indexMappingEmetrics).Do(ctx)
	if err != nil {
		logr.WithError(err).Debug("PutMapping logs failed")
		return err
	}
	if !res.Acknowledged {
		return errors.New("the put mapping was not acknowledged")
	}
	client.Alias().Remove(indexNameV1, indexName).Do(ctx)
	client.Alias().Add(mainIndex, indexName).Do(ctx)

	exists, err = client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		logr.WithError(err).Errorf("IndexExists for %s failed", indexName)
	}
	logr.Infof("after creation of index %s, exists: %t", indexName, exists)

	if false {
		existsV1, err := client.IndexExists(indexNameV1).Do(ctx)
		//existsV1, err := client.IndexExists(indexNameV1).Do(context.Background())
		if err != nil {
			logr.WithError(err).Errorf("IndexExists for %s failed", indexNameV1)
		}
		if existsV1 {
			logr.Infof("reindex from %s to %s", indexNameV1, indexNameV2)

			src := es.NewReindexSource().Index(indexNameV1)
			dst := es.NewReindexDestination().Index(indexNameV2)
			res, err := client.Reindex().Source(src).Destination(dst).Refresh("true").Do(context.Background())

			if err != nil {
				logr.WithError(err).Debugf("Reindex from %s to %s failed", indexNameV1, indexNameV2)
			}
			logr.Infof("Reindex of %v documents took %v",
				res.Total, res.Took)

			logr.Infof("deleting index %s", indexNameV1)
			client.DeleteIndex(indexNameV1).Do(ctx)
		}
	}

	logr.Debugf("function exit")
	return err
}

func (s *storageService) GetLogs(in *grpc_storage.Query, stream grpc_storage.Storage_GetLogsServer) error {
	start := time.Now()

	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)

	s.reportOnCluster("GetLogs", logr)

	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)
	dlogr.Debugf("function entry: %+v", in)

	query, shouldPostSortProcess, err := makeESQueryFromDlaasQuery(in)
	if err != nil {
		logr.WithError(err).Errorf("Could not make elasticsearch query from submitted query: %+v", *in)
		return err
	}
	dlogr.Debugf("pos: %d, since: %s, shouldPostSortProcess: %t", in.Pos, in.Since, shouldPostSortProcess)

	pagesize := int(in.Pagesize)
	if pagesize == 0 {
		pagesize = defaultPageSize
	}

	pos, isBackward := adjustOffsetPos(int(in.Pos))

	res, err := s.executeQuery(indexName, query, isBackward, pos,
		pagesize, docTypeLog, shouldPostSortProcess, dlogr)

	doneQuery := time.Now()

	if err != nil {
		logr.WithError(err).Errorf("Search failed")
		return err
	}

	dlogr.Debugf("Found: %d out of %d", len(res.Hits.Hits), res.TotalHits())

	if !(res.Hits == nil || res.Hits.Hits == nil || len(res.Hits.Hits) == 0) {
		logLineRecord := new(grpc_storage.LogLine)
		start := 0
		if isBackward {
			start = len(res.Hits.Hits) - 1
		}
		count := 0

		for i := start; ; {
			err := json.Unmarshal(*res.Hits.Hits[i].Source, &logLineRecord)
			if err != nil {
				logr.WithError(err).Errorf("Unmarshal from ES failed!")
				return err
			}
			dlogr.Debugf("GetLogs: %d (%s): %s",
				logLineRecord.Meta.Rindex, logLineRecord.Meta.Subid, makeSnippetForDebug(logLineRecord.Line, 7))
			err = stream.Send(logLineRecord)
			if err != nil {
				logr.WithError(err).Errorf("stream.Send failed")
				return err
			}

			if isBackward {
				i--
				if i < 0 {
					break
				}
			} else {
				i++
				if i >= len(res.Hits.Hits) {
					break
				}
			}
			count++
			if count >= pagesize {
				break
			}
		}
	}
	s.reportTime(logr, "GetLogs", in.Meta.TrainingId, in.Meta.Rindex, in.Meta.Time, start, doneQuery)

	dlogr.Debugf("function exit")
	return nil
}

// Hello is simple a gRPC test endpoint.
func (s *storageService) Hello(ctx context.Context, in *grpc_storage.Empty) (*grpc_storage.HelloResponse, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	logr.Debugf("function entry")
	logr.Debugf("Hello!")
	out := new(grpc_storage.HelloResponse)
	out.Msg = "Hello from ffdl-trainingdata-service!"
	logr.Debugf("function exit")
	return out, nil
}

func makeESQueryFromDlaasQuery(in *grpc_storage.Query) (es.Query, bool, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)
	var query es.Query

	shouldPostSortProcess := true

	trainingIDFieldName := "meta.training_id.keyword"
	timeFieldName := "meta.time"
	rindexFieldName := "meta.rindex"
	subidFieldName := "meta.subid"

	if in.SearchType == grpc_storage.Query_TERM {
		// TODO: This should be a time-based string
		var since int64
		var err error
		if in.Since != "" {
			dlogr.Debugf("Query_ since: %s", in.Since)
			since, err = humanStringToUnixTime(in.Since)
			if err != nil {
				logr.WithError(err).Errorf(
					"For now the since argument must be an integer representing " +
						"the number of milliseconds since midnight January 1, 1970")
				return nil, false, err
			}
		} else if in.Meta.Time != 0 {
			since = in.Meta.Time
		} else {
			since = 0
		}
		var idQuery es.Query
		if in.Meta.Subid != "" {
			idQuery = es.NewBoolQuery().Filter(
				es.NewTermQuery(trainingIDFieldName, in.Meta.TrainingId),
				es.NewTermQuery(subidFieldName, in.Meta.Subid),
			)
		} else {
			idQuery = es.NewTermQuery(trainingIDFieldName, in.Meta.TrainingId)
		}

		if since == 0 {
			if in.Pos > 0 {
				dlogr.Debugf("Query_ NewRangeQuery (pos)")
				query = es.NewBoolQuery().Filter(
					idQuery,
					es.NewBoolQuery().Filter(
						es.NewRangeQuery(rindexFieldName).Gte(in.Pos),
						es.NewRangeQuery(rindexFieldName).Lt(in.Pos+int64(in.Pagesize)),
					),
				)
				shouldPostSortProcess = false
			} else {
				dlogr.Debugf("Query_ NewTermQuery")
				query = idQuery
			}
		} else {
			dlogr.Debugf("Query_ NewRangeQuery")
			query = es.NewBoolQuery().Filter(
				idQuery,
				es.NewRangeQuery(timeFieldName).Gte(since),
			)
		}
	} else {
		err := fmt.Errorf("search type not supported: %s", grpc_storage.Query_SearchType_name[int32(in.SearchType)])
		logr.WithError(err).Error("Can't perform query")
		return nil, false, err
	}
	return query, shouldPostSortProcess, nil
}

func makeDebugLogger(logrr *logrus.Entry, isEnabled bool) *logger.LocLoggingEntry {
	logr := new(logger.LocLoggingEntry)
	logr.Logger = logrr
	logr.Enabled = isEnabled

	return logr
}

func (s *storageService) reportOnCluster(method string, logr *logger.LocLoggingEntry) {
	//ctx := context.Background()
	//
	//res, err := c.es.ClusterHealth().Do(ctx)
	//if err != nil {
	//	logr.WithError(err).Errorf("can't get cluster health!")
	//	return
	//}
	//logr.Debugf("cluster health (%s): %s", method, res.Status)
}

func (s *storageService) KillTrainingJob(ctx context.Context, request *grpc_storage.KillRequest) (*grpc_storage.KillResponse, error) {
	//Init Response & Logger & LCM Client
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService))
	killResponse := &grpc_storage.KillResponse{
		Success: false, IsCompleted:false }

	//Get Training Job From Mongo
	trainingJob, err := s.repo.Find(request.TrainingId)
	if err != nil {
		logr.Errorf("Kill Training Job Error:" , err.Error())
		return killResponse, err
	}

	if request.IsSa == false{
		if trainingJob.UserID != request.UserId {
			msg := fmt.Sprintf("User %s does not have permission to KILL training job with id %s.",
				request.UserId , request.TrainingId)
			logr.Error(msg)
			return nil, gerrf(codes.PermissionDenied, msg)
		}
	}

	jobStatus := trainingJob.TrainingStatus.Status
	//IF COMPLETED OR FAILED RETURN
	if jobStatus == grpc_storage.Status_COMPLETED || jobStatus == grpc_storage.Status_FAILED{
		killResponse.IsCompleted = true
		return killResponse, nil
	}else if jobStatus == grpc_storage.Status_KILLING || jobStatus == grpc_storage.Status_CANCELLED {
		killResponse.Success = true
		return killResponse, nil
	}

	if jobStatus == grpc_storage.Status_QUEUED {
		trainingJob.TrainingStatus.Status = grpc_storage.Status_CANCELLED
	}else{
		trainingJob.TrainingStatus.Status = grpc_storage.Status_KILLING
	}
	//TODO: TRANSACTION
	err = s.repo.Store(trainingJob)
	if err != nil {
		logr.Errorf("Kill/Cancel Training Job Store Error:" , err.Error())
		return killResponse, err
	}

	//IF Job is Running, Start LCM Kill
	if jobStatus == grpc_storage.Status_RUNNING ||
		jobStatus == grpc_storage.Status_PENDING{
		go func() {
			lcmClient, err := s.lcmClient()
			if err != nil {
				logr.Errorf("Kill Training Job Error:" , err.Error())
				return
			}
			defer lcmClient.Close()

			//STOP Job in LCM QUEUE
			jobKillRequest := service.JobKillRequest{
				TrainingId:   trainingJob.TrainingID,
				UserId:       trainingJob.UserID,
				JobNamespace: trainingJob.Namespace,
				Name: trainingJob.TrainingID,
				JobType: trainingJob.JobType,
			}
			_, err =  lcmClient.Client().KillMLFlowJob(context.Background(), &jobKillRequest)
			if err != nil{
				logr.Errorf("Kill Training Job Error:" , err.Error())
				return
			}

			//UPDATE TO CANCEL
			trainingJob.TrainingStatus.Status = grpc_storage.Status_CANCELLED
			trainingJob.TrainingStatus.CompletionTimestamp = util.CurrentTimestampAsString()
			//TODO: Transaction
			err = s.repo.Store(trainingJob)
			if err != nil {
				logr.Errorf("Kill Training Job Error:" , err.Error())
				return
			}
		}()
	}
	killResponse.Success = true
	return killResponse, nil
}

func makeSnippetForDebug(str string, maxLen int) string {
	//noinspection GoBoolExpressions
	if TdsDebugMode {
		str = strings.TrimSpace(str)
		if str != "" {
			endIndex := maxLen
			if endIndex >= len(str) {
				endIndex = len(str) - 1
			}
			return str[:endIndex]
		}
	}
	return ""
}

//noinspection GoBoolExpressions
func (s *storageService) reportTime(logr *logger.LocLoggingEntry,
	method string,
	trainingID string,
	rIndex int64,
	logReportTimeInMilliseconds int64,
	endPointEntry time.Time, doneQuery time.Time,
) {
	if TdsReportTimes {
		logReportTime := time.Unix(0, logReportTimeInMilliseconds*int64(time.Millisecond))

		elapsedFromReportTimeToEndPointEntryTime := endPointEntry.Sub(logReportTime)

		elapsedFromEndPointEntryToQueryDone := doneQuery.Sub(endPointEntry)

		elapsedDoneQueryToExit := time.Since(doneQuery)

		//logr.Debugf("%s query took %f", method, elapsed.Seconds())
		fmt.Printf("%7d\t%12.3f\t%12.3f\t%12.3f\t%s\t%s\n",
			int(rIndex),
			elapsedFromReportTimeToEndPointEntryTime.Seconds(),
			elapsedFromEndPointEntryToQueryDone.Seconds(),
			elapsedDoneQueryToExit.Seconds(),
			method, trainingID)
	}
}

func (s *storageService) refreshIndexes(ctx context.Context, logr *logger.LocLoggingEntry) {

	logr.Debugf("Refresh")

	res, err := s.es.Refresh(indexName).Do(ctx)
	if err != nil {
		logr.WithError(err).Debugf("Refresh failed")
	}
	if res == nil {
		logr.Debugf("Refresh expected result; got nil")
	}
	s.es.Update()
}

func (s *storageService) executeQuery(index string, query es.Query,
	isBackward bool, pos int, pagesize int, typ string, postSortProcess bool,
	dlogr *logger.LocLoggingEntry) (*es.SearchResult, error) {
	var res *es.SearchResult
	var err error
	ctx := context.Background()

	//c.refreshIndexes(dlogr)

	if postSortProcess {
		dlogr.Debugf("executing query with post-sort processing")
		res, err = s.es.Search(index).
			Index(index).
			Type(typ).
			Query(query).
			Sort("meta.time", !isBackward).
			From(pos).
			Size(pagesize).
			Do(ctx)
	} else {
		dlogr.Debugf("executing query with no post-sort processing")
		res, err = s.es.Search(index).
			Index(index).
			Type(typ).
			Query(query).
			Sort("meta.time", !isBackward).
			Do(ctx)
	}

	return res, err
}

func adjustOffsetPos(pos int) (int, bool) {
	isBackward := false
	if pos < 0 {
		isBackward = true
		pos = -pos
		if pos >= 1 {
			pos = pos - 1
		}
	}
	return pos, isBackward
}

// GetEMetrics returns a stream of evaluation metrics records.
func (s *storageService) GetEMetrics(in *grpc_storage.Query, stream grpc_storage.Storage_GetEMetricsServer) error {
	start := time.Now()
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyStorageService)).
		WithField(logger.LogkeyTrainingID, in.Meta.TrainingId).
		WithField(logger.LogkeyUserID, in.Meta.UserId)
	s.reportOnCluster("GetEMetrics", logr)
	//noinspection GoBoolExpressions
	dlogr := makeDebugLogger(logr.Logger, TdsDebugMode)
	dlogr.Debugf("function entry: %+v", in)

	query, shouldPostSortProcess, err := makeESQueryFromDlaasQuery(in)
	if err != nil {
		logr.WithError(err).Errorf("Could not make elasticsearch query from submitted query: %+v", *in)
		return err
	}

	pagesize := int(in.Pagesize)
	if pagesize == 0 {
		pagesize = defaultPageSize
	}

	pos, isBackward := adjustOffsetPos(int(in.Pos))

	dlogr.Debugf("pos: %d, pagesize: %d isBackward: %t", pos, pagesize, isBackward)

	res, err := s.executeQuery(indexName, query, isBackward, pos, pagesize,
		docTypeEmetrics, shouldPostSortProcess, dlogr)

	doneQuery := time.Now()

	if err != nil {
		logr.WithError(err).Errorf("Search failed")
		return err
	}

	dlogr.Debugf("Found: %d out of %d", len(res.Hits.Hits), res.TotalHits())

	if !(res.Hits == nil || res.Hits.Hits == nil || len(res.Hits.Hits) == 0) {
		start := 0
		if isBackward {
			start = len(res.Hits.Hits) - 1
		}
		count := 0
		for i := start; ; {
			emetricsRecord := new(grpc_storage.EMetrics)
			err := json.Unmarshal(*res.Hits.Hits[i].Source, &emetricsRecord)
			if err != nil {
				logr.WithError(err).Errorf("Unmarshal from ES failed!")
				return err
			}
			dlogr.Debugf("Sending record with rindex %v, time %v",
				emetricsRecord.Meta.Rindex,
				emetricsRecord.Meta.Time)

			err = stream.Send(emetricsRecord)
			if err != nil {
				logr.WithError(err).Errorf("stream.Send failed")
				return err
			}

			// logr.Debugf("EMetrics record: %+v\n", emetricsRecord)
			if isBackward {
				i--
				if i < 0 {
					break
				}
			} else {
				i++
				if i >= len(res.Hits.Hits) {
					break
				}
			}
			count++
			if count >= pagesize {
				break
			}
		}
	}
	s.reportTime(logr, "GetEMetrics", in.Meta.TrainingId, in.Meta.Rindex, in.Meta.Time, start, doneQuery)

	dlogr.Debugf("function exit")
	return nil
}

func (s *storageService) Upload(stream grpc_storage.Storage_UploadServer) error {
	getLogger := logger.GetLogger()
	//count upload time
	//now := time.Now()
	recv, err := stream.Recv()
	if err != nil {
		getLogger.Errorf("stream.Recv() err: %v", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	info := recv.GetInfo()

	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	useSSL := false
	bucketName := info.Bucket
	fileName := info.FileName
	hostPath := info.HostPath
	onlyMinio := info.OnlyMinio
	location := os.Getenv("AWS_TIMEZONE")
	getLogger.Debugf("upload, endpoint: %v, accessKeyID: %v, secretAccessKey: %v", endpoint, accessKeyID, secretAccessKey)
	getLogger.Debugf("upload, bucket: %v, fileName: %v, hostPath: %v, onlyMinio: %v, location: %v", bucketName, fileName, hostPath, onlyMinio, location)

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		getLogger.Errorf("minio.New err: %v", err.Error())
		return err
	}

	//--------------------------------------
	if hostPath != "" {
		//获取container中ceph的绝对路径
		//absHostPath, err := getAbsoluteHostPath(hostPath)

		//获取文件名
		splits := strings.Split(hostPath, "/")
		fileName := splits[len(splits)-1]
		getLogger.Debugf("filepath, fileName: %s", fileName)

		ctx, _ := context.WithCancel(context.TODO())
		//当为hostPath时
		//创建 bucket
		err = minioClient.MakeBucket(bucketName, location)
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := minioClient.BucketExists(bucketName)
			if errBucketExists == nil && exists {
				getLogger.Debugf("We already own %s", bucketName)
			} else {
				getLogger.Errorf("minioClient.MakeBucket failed, %v", err.Error())
				return err
			}
		} else {
			getLogger.Debugf("Successfully created %s", bucketName)
		}
		//put obj
		_, err = minioClient.FPutObjectWithContext(ctx, bucketName, fileName, hostPath, minio.PutObjectOptions{})
		if err != nil {
			getLogger.Errorf("minioClient.FPutObjectWithContext failed, err: %v\n", err.Error())
			return status.Error(codes.Internal, err.Error())
		}
		getLogger.Debugf("successful uploaded %v", fileName)
		response := grpc_storage.UploadResponse{
			Msg:      "success",
			Code:     "200",
			S3Path:   fmt.Sprintf("s3://%v/", bucketName),
			HostPath: hostPath,
		}
		err = stream.SendAndClose(&response)
		if err != nil {
			getLogger.Errorf("stream.SendAndClose failed, err: %v", err.Error())
			return err
		}
		return nil
	}
	//--------------------------------------

	//make bucket if it not exists
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			getLogger.Debugf("We already own %s\n", bucketName)
		} else {
			getLogger.Errorf("minioClient.MakeBucket failed, %v", err.Error())
			return err
		}
	} else {
		getLogger.Debugf("Successfully created %s", bucketName)
	}

	reader, writerToMinio := io.Pipe()
	defer reader.Close()

	var finalWriter io.Writer
	var filePath string
	if !onlyMinio {
		//create folder
		folderPath := fmt.Sprintf("%v/%v", viper.GetString(config.UploadContainerPath), bucketName)
		err = os.Mkdir(folderPath, 0775)

		if os.IsExist(err) {
			getLogger.Debugf("uploading folder is exist, err: %v", err.Error())
		} else if err != nil {
			getLogger.Errorf("Mkdir for uploading failed, err: %v", err.Error())
			return err
		}
		//upload to ceph at same time
		filePath = fmt.Sprintf("%v/%v/%v", viper.GetString(config.UploadContainerPath), bucketName, fileName)
		cephFile, err := os.Create(filePath)
		if err != nil {
			getLogger.Errorf("create file in ceph failed, err: %v", err.Error())
			return err
		}
		defer cephFile.Close()
		//upload to ceph and minio
		finalWriter = io.MultiWriter(writerToMinio, cephFile)
	} else {
		finalWriter = writerToMinio
	}

	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		defer writerToMinio.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				recv, err := stream.Recv()
				if err == io.EOF {
					getLogger.Debugf("no more data of %v/%v", bucketName, fileName)
					return
				}
				if err != nil {
					getLogger.Errorf("stream.Recv, err: %v", err.Error())
					cancel()
					return
				}
				thunk := recv.GetChucksData()
				_, err = finalWriter.Write(thunk)
				if err != nil {
					getLogger.Errorf("file.Write, err: %v", err.Error())
					cancel()
					return
				}
			}
		}
	}()
	_, err = minioClient.PutObjectWithContext(ctx, bucketName, fileName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		getLogger.Errorf("minioClient.PutObjectWithContext failed, err: %v\n", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	getLogger.Debugf("successful uploaded %v", fileName)
	//since := time.Since(now)
	//getLogger.Debugf("upload elapsed: ", since)

	response := grpc_storage.UploadResponse{
		Msg:    "success",
		Code:   "200",
		S3Path: fmt.Sprintf("s3://%v/", bucketName),
		//HostPath: fmt.Sprintf("%v//%v", viper.GetString(config.UploadContainerPath), fileName),
		HostPath: filePath,
		FileName: info.FileName,
	}
	err = stream.SendAndClose(&response)
	if err != nil {
		getLogger.Errorf("stream.SendAndClose failed, err: %v", err.Error())
		return err
	}

	return nil
}

func (s *storageService) Download(req *grpc_storage.DownloadRequest, stream grpc_storage.Storage_DownloadServer) error {
	getLogger := logger.GetLogger()
	//count Download duration
	now := time.Now()
	fileInfo := req.FileInfo

	//---------------
	//TODO need to get from ENV
	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	//endpoint := "mlss-minio.mlss-dev.svc.cluster.local:9000"
	//accessKeyID := "AKIAIOSFODNN7EXAMPLE"
	//secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	getLogger.Debugf("Download, endpoint: %v, accessKeyID: %v, secretAccessKey: %v, bucket: %v, fileName: %v", endpoint, accessKeyID, secretAccessKey, req.FileInfo.Bucket, req.FileInfo.FileName)

	//useSSL := true
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		getLogger.Errorf("minio.New err: %v\n", err.Error())
		return err
	}

	ctx, _ := context.WithCancel(context.TODO())
	obj, err := minioClient.GetObjectWithContext(ctx, fileInfo.Bucket, fileInfo.FileName, minio.GetObjectOptions{})
	if err != nil {
		getLogger.Errorf("minioClient.GetObjectWithContext failed, err: %v\n", err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	//defer obj.Close()
	buffer := make([]byte, 4096)
	var flen = 0
	for {
		readLen, err := obj.Read(buffer)
		if err == io.EOF {
			getLogger.Debugf("no more data of %v/%v", fileInfo.Bucket, fileInfo.FileName)

			getLogger.Debugf("stream.Send one more time! so fucking strange!")
			bytes := buffer[:readLen]
			response := grpc_storage.DownloadResponse{ChucksData: bytes}
			err = stream.Send(&response)
			if err != nil {
				getLogger.Errorf("stream.Send, err: %v", err.Error())
				return status.Error(codes.Internal, err.Error())
			}
			flen += len(bytes)

			getLogger.Debugf("final stream.Send( len, %v", flen)
			break
		}
		if err != nil {
			getLogger.Errorf("obj.Read(buffer), err: %v", err.Error())
			return status.Error(codes.Internal, err.Error())
		}
		//getLogger.Debugf("obj.Read len, %v", readLen)

		bytes := buffer[:readLen]
		response := grpc_storage.DownloadResponse{ChucksData: bytes}
		err = stream.Send(&response)
		if err != nil {
			getLogger.Errorf("stream.Send, err: %v", err.Error())
			return status.Error(codes.Internal, err.Error())
		}
		//getLogger.Debugf("stream.Send( len, %v", len(bytes))
		flen += len(bytes)
	}

	getLogger.Debugf("successful downloaded %v", fileInfo.FileName)
	since := time.Since(now)
	getLogger.Debugf("download elapsed: %v", since)

	return nil
}

//func (s *storageService) BucketExists(req *grpc_storage.DownloadRequest, stream grpc_storage.Storage_DownloadServer) error {
func (s *storageService) BucketExists(context.Context, *grpc_storage.BucketExistsRequest) (*grpc_storage.BucketExistsResponse, error) {
	getLogger := logger.GetLogger()
	getLogger.Debugf("BucketExists unimplemented")

	return nil, fmt.Errorf("BucketExists unimplemented")
}

func (s *storageService) UploadModelZip(ctx context.Context, req *grpc_storage.CodeUploadRequest) (*grpc_storage.CodeUploadResponse, error) {
	res := grpc_storage.CodeUploadResponse{
		Code:   "500",
		Msg:    "",
		S3Path: "",
	}
	logger.GetLogger().Debugf("UploadModelZip req: %+v\n", *req)
	filePath := req.HostPath + req.FileName
	//useSSL := true
	useSSL := false
	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		logger.GetLogger().Error("new minio client error," + err.Error())
		return &res, err
	}
	uid := uuid2.New()
	bucketName := "mlss-mf"
	unFilePath := uid.String()
	dirName := "nil"
	dirCount := 0
	err = archiver.Unarchive(filePath, unFilePath)
	if err != nil {
		logger.GetLogger().Error("unzip err," + err.Error())
		return nil, err
	}
	err = filepath.Walk(unFilePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logger.GetLogger().Debugf("lk-test downloadZip path: %s, info.Name: %s", path, info.Name())
			_, err = minioClient.FPutObject(bucketName, path, path, minio.PutObjectOptions{ContentType: "application/zip"})
			if err != nil {
				logger.GetLogger().Error("minioClient putObject err," + err.Error())
				return err
			}
		} else {
			if dirCount == 1 {
				logger.GetLogger().Debugf("lk-test dirCount = 1 downloadZip path: %s, info.Name: %s", path, info.Name())
				dirName = info.Name()
			}
			dirCount++
		}
		return err
	})
	if err != nil {
		logger.GetLogger().Error("filepath walk err," + err.Error())
		return nil, err
	}
	res.Code = "200"
	res.S3Path = "s3://" + bucketName + "/" + uid.String() + "/" + dirName
	logger.GetLogger().Debugf("lk-test S3Path: %s", res.S3Path)
	return &res, err
}

func (s *storageService) UploadCode(ctx context.Context, req *grpc_storage.CodeUploadRequest) (*grpc_storage.CodeUploadResponse, error) {
	res := grpc_storage.CodeUploadResponse{
		Code:   "500",
		Msg:    "",
		S3Path: "",
	}
	// filePath := req.HostPath + req.FileName
	filePath := path.Join(req.HostPath, req.FileName)
	getLogger := logger.GetLogger()
	//useSSL := true
	useSSL := false
	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		getLogger.Error("new minio client error," + err.Error())
		return &res, err
	}
	bucketName := "mlss-di"

	getLogger.Info("Request's bucketName：", req.Bucket)
	if len(req.Bucket) > 0 {
		uid := uuid2.New()
		bucketName = req.Bucket
		getLogger.Info("bucketName：", bucketName)
		if req.Bucket == "mlss-mf" {
			req.FileName = fmt.Sprintf("%v/%v", uid, req.FileName)
		}
	}
	location := os.Getenv("AWS_TIMEZONE")
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		getLogger.Printf("MakeBucket Error: " + err.Error())
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			getLogger.Printf("We already own %s\n", bucketName)
		} else {
			getLogger.Fatalln(err)
		}
	} else {
		getLogger.Printf("Successfully created %s\n", bucketName)
	}
	getLogger.Debugf("Upload, endpoint: %v, accessKeyID: %v, secretAccessKey: %v, bucket: %v, fileName: %v, filePath: %s", endpoint, accessKeyID, secretAccessKey,
		bucketName, req.FileName, filePath)

	num, err := minioClient.FPutObject(bucketName, req.FileName, filePath, minio.PutObjectOptions{ContentType: "application/zip"})

	if err != nil {
		getLogger.Errorf("minio.put object err: %v\n", err.Error())
		return &res, err
	}
	if num == 0 {
		file, err := os.Open(filePath)
		if err != nil {
			getLogger.Errorf("upload file num is zero, bucketName if : open file error")
		}
		fs, err := file.Stat()
		if err != nil {
			getLogger.Errorf("upload file num is zero, bucketName if : stat file error")
		}
		defer file.Close()
		getLogger.Errorf("upload file num is zero, bucketName if : %v\n", fs.Size())
		getLogger.Errorf("upload file num is zero, bucketName if : %v\n", bucketName)
		getLogger.Errorf("upload file num is zero, filepath if : %v\n", filePath)
		getLogger.Errorf("upload file num is zero, filename if : %v\n", req.FileName)
	}
	if len(req.Bucket) > 0 {
		if req.Bucket == "mlss-mf" {
			res.Code = "200"
			res.S3Path = "s3://" + bucketName + "/" + strings.Split(req.FileName, "/")[0]
		} else {
			res.Code = "200"
			res.S3Path = "s3://" + bucketName + "/" + req.FileName
		}
	} else {
		res.Code = "200"
		res.S3Path = "s3://" + bucketName + "/" + req.FileName
	}
	//return nil, status.Errorf(codes.Unimplemented, "method UploadCode not implemented")
	getLogger.Debugf("Upload response: %+v\n", res)
	return &res, err
}

func (s *storageService) DownloadCode(CTX context.Context, request *grpc_storage.CodeDownloadRequest) (*grpc_storage.CodeDownloadResponse, error) {
	s3Path := request.S3Path

	getLogger := logger.GetLogger()
	res := grpc_storage.CodeDownloadResponse{
		Code:     "500",
		Msg:      "",
		FileName: "",
		HostPath: "",
		Bucket:   "",
	}

	if !strings.Contains(s3Path, "s3://") {
		res.Msg = "s3 Path is not correct, did not contains s3:// "
		return &res, errors.New("s3 Path is not correct")
	}

	s3Array := strings.Split(s3Path[5:], "/")
	if len(s3Array) < 2 {
		res.Msg = "s3 Path is not correct"
		getLogger.Error("new minio client error," + s3Path)
		return &res, errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]
	filename := ""
	filename = s3Path[(5 + len(bucketName)):]
	// if bucketName == "mlss-mf" {
	// 	// s3Path 应该包含upload时的uid+文件名， 例： s3://mlss-mf/09380357-a9c9-4c38-a4e7-39b04a8b886f/main.go
	// 	// filename = "09380357-a9c9-4c38-a4e7-39b04a8b886f/main.go"
	// 	filename = s3Path[(5 + len(bucketName)):]

	// 	// todo not deal zip file,
	// 	// because type zip file is uploaded after unzip, cannot download  zip file
	// } else {
	// 	filename = s3Path[(5 + len(bucketName)):]
	// }

	getLogger.Infof("DownloadCode file name %q from s3Path %q", filename, s3Path)

	useSSL := false

	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		getLogger.Error("new minio client error," + err.Error())
		return &res, err
	}

	filedir, err := uuid.NewV4()
	if err != nil {
		getLogger.Error("download code uuid new Error Path:", err.Error())
		return &res, err
	}
	hostpath := "/data/oss-storage/download/" + filedir.String()
	// test zip file by lk
	// ok
	// filename = "b9b15716-65c5-443a-b31a-d8c9f1192f36/fluentbit-lk/fluent-bit-configmap-alexwu.yaml"
	err = minioClient.FGetObject(bucketName, filename, hostpath+"/"+filename, minio.GetObjectOptions{})
	if err != nil {
		getLogger.Error("FGETObject Error Path:", bucketName+"/"+filename)
		getLogger.Error("FGETObject Error: " + err.Error())
		return &res, err
	}

	// objPrefix := "b9b15716-65c5-443a-b31a-d8c9f1192f36/fluentbit-lk"
	// doneCh := make(chan struct{})
	// defer close(doneCh)
	// getLogger.Infoln("lk-test strat to list object from minio ......")
	// for objInfo := range minioClient.ListObjects(bucketName, objPrefix, true, doneCh) {
	// 	getLogger.Infof("lk-test listObjects %+v\n", objInfo)
	// }
	// getLogger.Infoln("lk-test finish to list object from minio ......")

	// filename = "b9b15716-65c5-443a-b31a-d8c9f1192f36/fluentbit-lk.zip"
	// err = minioClient.FGetObject(bucketName, filename, hostpath+"/"+filename, minio.GetObjectOptions{})
	// if err != nil {
	// 	getLogger.Error("[2] FGETObject Error Path:", bucketName+"/"+filename)
	// 	getLogger.Error("[2]FGETObject Error: " + err.Error())
	// 	return &res, err
	// }

	// filename = "b9b15716-65c5-443a-b31a-d8c9f1192f36/fluentbit-lk"
	// err = minioClient.FGetObject(bucketName, filename, hostpath+"/"+filename, minio.GetObjectOptions{})
	// if err != nil {
	// 	getLogger.Error("[3] FGETObject Error Path:", bucketName+"/"+filename)
	// 	getLogger.Error("[3]FGETObject Error: " + err.Error())
	// 	return &res, err
	// }

	res.Bucket = bucketName
	res.FileName = filename
	res.HostPath = hostpath
	res.Code = "200"
	return &res, err
}

func (s *storageService) DownloadCodeByDir(CTX context.Context, request *grpc_storage.CodeDownloadByDirRequest) (*grpc_storage.CodeDownloadByDirResponse, error) {
	s3Path := request.S3Path

	getLogger := logger.GetLogger()
	res := grpc_storage.CodeDownloadByDirResponse{
		RespList: make([]*grpc_storage.CodeDownloadResponse, 0),
	}

	if !strings.Contains(s3Path, "s3://") {
		resSub := grpc_storage.CodeDownloadResponse{
			Code: "500",
		}
		resSub.Msg = "s3 Path is not correct, did not contains s3:// "
		res.RespList = append(res.RespList, &resSub)
		return &res, errors.New("s3 Path is not correct, did not contains s3://")
	}

	s3Array := strings.Split(s3Path[5:], "/")
	logger.GetLogger().Debugf("lk-test s3Array: %+v, length: %d\n", s3Array, len(s3Array))
	if len(s3Array) < 2 {
		resSub := grpc_storage.CodeDownloadResponse{
			Code: "500",
		}
		resSub.Msg = "s3 Path is not correct"
		res.RespList = append(res.RespList, &resSub)
		getLogger.Error("new minio client error," + s3Path)
		return &res, fmt.Errorf("s3 Path is not correct, s3Array: %+v, length: %d", s3Array, len(s3Array))
	}
	bucketName := s3Array[0]
	filename := s3Path[(5 + len(bucketName) + 1):]

	getLogger.Infof("DownloadCode file name %q from s3Path %q", filename, s3Path)

	useSSL := false
	storeConfig := config.GetDataStoreConfig()
	endpoint := storeConfig[config.AuthURLKey]
	accessKeyID := storeConfig[config.UsernameKey]
	secretAccessKey := storeConfig[config.PasswordKey]
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		resSub := grpc_storage.CodeDownloadResponse{
			Code: "500",
		}
		resSub.Msg = fmt.Sprintf("new minio client error, %s", err.Error())
		res.RespList = append(res.RespList, &resSub)
		getLogger.Error("new minio client error," + err.Error())
		return &res, err
	}

	filedir, err := uuid.NewV4()
	if err != nil {
		resSub := grpc_storage.CodeDownloadResponse{
			Code: "500",
		}
		resSub.Msg = fmt.Sprintf("download code uuid new Error Path, %s", err.Error())
		res.RespList = append(res.RespList, &resSub)
		getLogger.Error("download code uuid new Error Path:", err.Error())
		return &res, err
	}
	hostpath := "/data/oss-storage/download/" + filedir.String()

	if strings.Contains(filename, ".zip") {
		//eg: objPrefix = "b9b15716-65c5-443a-b31a-d8c9f1192f36"
		objPrefix := s3Array[1]
		doneCh := make(chan struct{})
		defer close(doneCh)
		getLogger.Infoln("lk-test start to list object from minio ......")
		for objInfo := range minioClient.ListObjects(bucketName, objPrefix, true, doneCh) {
			getLogger.Infof("lk-test listObjects %+v\n", objInfo)
			resSub := grpc_storage.CodeDownloadResponse{
				Code:     "500",
				Msg:      "",
				FileName: "",
				HostPath: "",
				Bucket:   "",
			}

			//todo
			// Key:b9311acd-6ae0-4624-a9a7-87c134827a61/fluentbit-lk/fluent-bit-role.yaml
			err = minioClient.FGetObject(bucketName, objInfo.Key, hostpath+"/"+objInfo.Key, minio.GetObjectOptions{})
			if err != nil {
				getLogger.Error("FGETObject Error Path:", bucketName+"/"+filename)
				getLogger.Error("FGETObject Error: " + err.Error())
				resSub.Msg = fmt.Sprintf("fail to get object, %s", err.Error())
				res.RespList = append(res.RespList, &resSub)
				return &res, err
			}
			resSub.Bucket = bucketName
			resSub.FileName = objInfo.Key
			resSub.HostPath = hostpath
			resSub.Code = "200"
			res.RespList = append(res.RespList, &resSub)
			getLogger.Debugf("lk-test sub object FileName: %s, HostPath: %s\n", resSub.FileName, resSub.HostPath)
		}
	} else {
		// eg: filename = "b9b15716-65c5-443a-b31a-d8c9f1192f36/fluentbit-lk/fluent-bit-configmap-alexwu.yaml"
		resSub := grpc_storage.CodeDownloadResponse{
			Code:     "500",
			Msg:      "",
			FileName: "",
			HostPath: "",
			Bucket:   "",
		}
		err = minioClient.FGetObject(bucketName, filename, hostpath+"/"+filename, minio.GetObjectOptions{})
		if err != nil {
			getLogger.Error("FGETObject Error Path:", bucketName+"/"+filename)
			getLogger.Error("FGETObject Error: " + err.Error())
			resSub.Msg = fmt.Sprintf("fail to get object, %s", err.Error())
			res.RespList = append(res.RespList, &resSub)
			return &res, err
		}
		resSub.Bucket = bucketName
		resSub.FileName = filename
		resSub.HostPath = hostpath
		resSub.Code = "200"
		res.RespList = append(res.RespList, &resSub)
	}

	return &res, err
}
