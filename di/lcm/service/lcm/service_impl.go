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
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"webank/DI/commons/constants"
	"webank/DI/commons/tfjob"
	// "webank/DI/pkg/operators/tf-operator/client/clientset/versioned"

	// "webank/DI/linkisexecutor/linkisclientutils"
	// "webank/DI/linkisexecutor/models"
	"webank/DI/commons/operators/tf-operator/client/clientset/versioned"
	"webank/DI/storage/storage/grpc_storage"
	"webank/DI/trainer/storage"
	"webank/DI/trainer/trainer"
	"webank/DI/trainer/trainer/grpc_trainer_v2"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"webank/DI/lcm/coord"

	"google.golang.org/grpc"

	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/lcm/lcmconfig"
	"webank/DI/storage/client"

	"github.com/cenkalti/backoff"
	"github.com/coreos/etcd/clientv3"
	"github.com/go-kit/kit/metrics"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"webank/DI/commons/locker"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	tfclientset "webank/DI/commons/operators/tf-operator/client/clientset/versioned"
)

const (
	internalObjectStoreID = "dlaas_internal_os"

	LimitsCpu            = "limits.cpu"
	RequestsNvidiaComGpu = "requests.nvidia.com/gpu"
	LimitsMemory         = "limits.memory"
	JobmonitorCpu        = 500
	JobmonitorMemory     = 536870912

	HelperCpu    = 20
	HelperMemory = 104857600

	//modelsBucketKey        = "objectstore.bucket.models"
	trainedModelsBucketKey = "objectstore.bucket.trainedmodels"

	defaultModelsBucket        = "di-models"
	defaultTrainedModelsBucket = "di-trained-models"

	//collectionNameTrainingJobs = "training_jobs"
	collectionNameJobHistory = "job_history"

	//debugLogsMode = false

	//oldEndpointInternalPageSize = 10

	jobLeanerNumber  = 1
	jobHelperNumber  = 1
	jobMonitorNumber = 1

	MongoAddressKey             = "MONGO_ADDRESS"
	MongoDatabaseKey            = "MONGO_DATABASE"
	MongoUsernameKey            = "MONGO_USERNAME"
	MongoPasswordKey            = "MONGO_PASSWORD"
	MongoAuthenticationDatabase = "MONGO_Authentication_Database"

	//pollIntervalKey = "queue.poll.interval"
	queuePollInterval = "queuePollInterval"
	queueSize         = "queueSize"

	hardwareType = "hardware-type"
	nvidiagpu    = "NVIDIAGPU"
	cpucompute   = "CPUCOMPUTE"
)

// Confuse `go vet' to not check this `Errorf' call. :(
// See https://github.com/grpc/grpc-go/issues/90
var gerrf = grpc.Errorf

var (
	//NativeFrameworks which support native distribution
	NativeFrameworks = []string{"tensorflow", "caffe2", "mxnet", "horovod"}
	totalTrainingCounter, finishedTrainingCounter,
	failedToLaunchTrainingsCounter, k8sFailureCounter metrics.Counter
)

var lockForCreatingQueue = locker.NewNamedMutex()

//Service LCM manages the lifecycle of the entire distributed deep learning job
type Service interface {
	service.LifecycleManagerServer
	service.LifecycleHandler
	StopLCM()
}

type lcmService struct {
	service.Lifecycle
	k8sClient  kubernetes.Interface
	etcdClient coord.Coordinator
	// metrics             *lcmMetrics
	datastore           storage.DataStore
	queues              map[string]*queueHandler
	mtx                 sync.RWMutex
	repo                repository
	jobHistoryRepo      jobHistoryRepository
	queuesStarted       bool
	modelsBucket        string
	trainedModelsBucket string
	tfjobClient         *tfclientset.Clientset
}

type queueHandler struct {
	stopQueue chan struct{}
	*trainer.TrainingJobQueue
	isQueueStarted int64
}

// Entry represents a single training job in the queue
type Entry struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	TrainingID string        `bson:"training_id" json:"training_id"`
	Submitted  time.Time     `bson:"submitted" json:"submitted"`
}

//NewService is a constructor to initialize LCM
func NewService() (Service, error) {
	s, err := newService()
	return s, err
}

func (s *lcmService) StopLCM() {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyLcmService))
	logr.Debugf(" ###### shutting down lcm ###### ")
	s.etcdClient.Close(logr)
	s.Stop() // stop Service
}

func newService() (*lcmService, error) {
	logr := logger.LocLogger(logger.LogServiceBasic(logger.LogkeyLcmService))
	logr.Debugf("Instantiating LCM service...")
	logr.Debugf("LCM Deployment Mode is %s , pod namespace %s , learner namespace %s , databroker tag %s, learner tag %s",
		viper.GetString(config.LCMDeploymentKey), config.GetPodNamespace(), config.GetLearnerNamespace(),
		viper.GetString(config.DataBrokerTagKey), viper.GetString(config.LearnerTagKey))

	//config.FatalOnAbsentKey(mongoAddressKey)
	config.FatalOnAbsentKey(config.ETCDEndpoints)
	//config.SetDefault(gpuLimitsQuerySizeKey, 200)
	config.SetDefault(queueSize, 5000)
	config.SetDefault(queuePollInterval, 10) // in seconds

	defaultBackoff := backoff.NewExponentialBackOff()
	defaultBackoff.MaxElapsedTime = 1 * time.Minute
	// assert necessary config keys
	// FIXME MLSS Change: v_1.4.1 init dataStore
	ds, err := storage.CreateDataStore(config.GetDataStoreType(), config.GetDataStoreConfig())
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create datastore")
	}
	err = ds.Connect()
	if err != nil {
		logr.WithError(err).Fatalf("Cannot connect to object store")
	}
	// FIXME MLSS Change: v_1.4.1 init repo

	logr.Debugf("MongoAddressKey: %v, MongoDatabaseKey: %v, MongoUsernameKey: %v, MongoPasswordKey: %v",
		viper.GetString(MongoAddressKey), viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey),
		viper.GetString(MongoPasswordKey))

	repo, err := newTrainingsRepository(viper.GetString(MongoAddressKey),
		viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey),
		viper.GetString(MongoPasswordKey), viper.GetString(MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), "training_jobs")
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create repository with %s %s %s", viper.GetString(MongoAddressKey),
			viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey))
	}
	jobHistoryRepo, err := newJobHistoryRepository(viper.GetString(MongoAddressKey),
		viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey),
		viper.GetString(MongoPasswordKey), viper.GetString(MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), collectionNameJobHistory)
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create repository with %s %s %s %s", viper.GetString("mongo.address"),
			viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey), collectionNameJobHistory)
	}

	k8sClient, err := kubernetes.NewForConfig(lcmconfig.GetKubernetesConfig())

	tfjobClient := versioned.NewForConfigOrDie(lcmconfig.GetKubernetesConfig())

	if err != nil {
		logr.WithError(err).Errorf("Failed to create a kubernetes client: %v", lcmconfig.GetKubernetesConfig())
		return nil, err
	}

	client, connectivityErr := coordinator(logr)
	if connectivityErr != nil {
		logr.WithError(connectivityErr).Errorln("failed to connect to etcd when starting, this should trigger restart of lcm")
		return nil, connectivityErr
	}
	// FIXME MLSS Change: v_1.4.1 init queues
	queues, err := getQueuesFromDB(viper.GetString(MongoAddressKey),
		viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey),
		viper.GetString(MongoPasswordKey), viper.GetString(config.MongoAuthenticationDatabase),
		config.GetMongoCertLocation())

	if err != nil {
		logr.WithError(err).Errorf("Failed to getQueues from mongo when newService", lcmconfig.GetKubernetesConfig())

		return nil, err
	}

	s := &lcmService{
		k8sClient:  k8sClient,
		etcdClient: client,
		datastore:  ds,
		// FIXME MLSS Change: v_1.4.1 add field to lcmService
		repo:                repo,
		jobHistoryRepo:      jobHistoryRepo,
		modelsBucket:        trainer.GetModelsBucket(),
		trainedModelsBucket: getTrainedModelsBucket(),
		queues:              queues,
		queuesStarted:       false,
		tfjobClient:         tfjobClient,
	}

	s.RegisterService = func() {
		service.RegisterLifecycleManagerServer(s.Server, s)
	}
	s.startQueues()

	return s, nil
}

//FIXME 1.9.0 lcm DeployTrainingJob
func (s *lcmService) DeployTrainingJob(ctx context.Context, req *service.JobDeploymentRequest) (*service.JobDeploymentResponse, error) {
	//extend the logger with required fields and this logr will be passed around
	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId).WithFields(logrus.Fields{
		"name":      req.Name,
		"framework": req.Framework,
		"gpus":      req.Resources.Gpus,
		"cpus":      req.Resources.Cpus,
		"memory":    req.Resources.Memory,
		"namespace": req.JobNamespace,
	}))

	tr := &trainer.TrainingRecord{
		TrainingID:            req.TrainingId,
		UserID:                req.UserId,
		ModelDefinition:       req.ModelDefinition,
		Training:              req.Training,
		Datastores:            req.DataStores,
		TrainingStatus:        req.TrainingStatus,
		Metrics:               nil,
		EvaluationMetricsSpec: req.EvaluationMetricsSpec,
		Namespace:             req.JobNamespace,
		GID:                   req.EnvVars["GID"],
		UID:                   req.EnvVars["UID"],
		GuardianToken:         req.EnvVars["GUARDIAN_TOKEN"],
		JobAlert:              req.JobAlert,
		PSs:                   req.PSs,
		PSCPU:                 req.PSCPU,
		PSImage:               req.PSImage,
		PSMemory:              req.PSMemory,
		CodeSelector:          req.CodeSelector,
		DataPath:              req.DataPath,
		JobType:               req.JobType,
		TFosRequest:           req.TFosRequest,
		ExpRunId:              req.ExpRunId,
		ExpName:               req.ExpName,
		FileName:              req.FileName,
		FilePath:              req.FilePath,
		ProxyUser:             req.ProxyUser,
		DataSet:               req.DataSet,
		MFModel:               req.MfModel,
		Algorithm:             req.Algorithm,
		JobParams:             req.JobParams,
		FitParams:             req.FitParams,
		APIType:               req.APIType,
	}
	// JobParmas:    req.Params,
	// DataSet:      req.DataSet}

	// FIXME MLSS Change: v_1.4.1 added
	ns := req.JobNamespace
	queueName := trainer.QueueName(ns)

	qHandler := s.queues[queueName]
	if qHandler == nil {
		//ensure get Handler
		handler, ensureQueueErr := s.ensureQueue(ns, logr)
		if ensureQueueErr != nil {
			logr.WithError(ensureQueueErr).Errorln("FAILED to ensureQueue")
			return &service.JobDeploymentResponse{Name: req.Name}, ensureQueueErr
		}
		qHandler = handler
	}

	err := qHandler.Lock()
	if err != nil {
		logr.WithError(err).Errorln("FAILED to qHandler lock")
		return &service.JobDeploymentResponse{Name: req.Name}, err
	}
	defer qHandler.Unlock()

	//check if rateLimited
	rateLimited, err := s.resourceChecking(req, tr, logr)
	if err != nil {
		logr.WithError(err).Errorln("resourceChecking failed")
		return &service.JobDeploymentResponse{Name: req.Name}, err
	}

	queueSizeFromConfig := viper.GetInt(queueSize)

	//rateLimited for the job
	if !rateLimited && qHandler == nil {
		logr.Debugf("deployJobToK8s with trainingID : s%", req.TrainingId)
		err = s.deployJobToK8s(tr, logr)
	} else {
		if qHandler == nil {
			logr.Debugf("deployTrainingJob qHandler not exists")
			//init and return handler for the ns
			qHandler, err = s.ensureQueue(ns, logr)
			if err != nil {
				logr.WithError(err).Errorln("FAILED to ensureQueue")
				return &service.JobDeploymentResponse{Name: req.Name}, err
			}
		}
		err = s.checkQueueAndEnqueue(qHandler, queueName, queueSizeFromConfig, tr, logr)
	}

	if err != nil {
		return &service.JobDeploymentResponse{Name: req.Name}, err
	}

	//unlock
	return &service.JobDeploymentResponse{Name: req.Name}, nil
}

// FIXME MLSS Change: v_1.4.1 added to create queue if queue is nil
func (s *lcmService) ensureQueue(namespace string, logr *logger.LocLoggingEntry) (*queueHandler, error) {
	queueName := trainer.QueueName(namespace)

	//named lock
	const timeoutForTryLock = 30
	timeout := lockForCreatingQueue.TryLockTimeout(namespace, time.Duration(timeoutForTryLock)*time.Second)
	if !timeout {
		timeoutMsg := fmt.Sprintf("try lock failed in %v sec, namespace: %s", viper.GetString(MongoAddressKey),
			viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey))
		logr.Fatalf(timeoutMsg)
		return nil, errors.New(timeoutMsg)
	}

	defer lockForCreatingQueue.Unlock(namespace)

	//if ns's queue Handler exists,return it
	handler := s.queues[queueName]
	if handler != nil {
		return handler, nil
	}

	logr.Debugf("DeployTrainingJob to newTrainingJobQueue with queueName: s%", queueName)
	queue, err := trainer.GetNewTrainingJobQueue(viper.GetString(MongoAddressKey),
		viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey),
		viper.GetString(MongoPasswordKey), viper.GetString(MongoAuthenticationDatabase),
		config.GetMongoCertLocation(), trainer.QueueName(namespace), trainer.LockName(namespace))
	if err != nil {
		logr.WithError(err).Fatalf("Cannot create queue with %s %s %s", viper.GetString(MongoAddressKey),
			viper.GetString(MongoDatabaseKey), viper.GetString(MongoUsernameKey))
		return nil, err
	}
	s.queues[queueName] = &queueHandler{make(chan struct{}), queue, constants.Fault}
	qHandler := s.queues[queueName]

	s.startQueue(queueName, logr)

	return qHandler, nil
}

// FIXME MLSS Change: v_1.4.1 added to startQueue
func (s *lcmService) startQueue(queueName string, logr *logger.LocLoggingEntry) error {
	//queueName := trainer.QueueName(namespace)
	qHandler := s.queues[queueName]

	swapped := atomic.CompareAndSwapInt64(&qHandler.isQueueStarted, constants.Fault, constants.True)

	if swapped {
		logr.Debugf("CompareAndSwapInt64 success, pull queue: %v in a goroutine", queueName)
		tick := time.NewTicker(time.Duration(viper.GetInt(queuePollInterval)) * time.Second).C
		go func(queueName string, qHandler *queueHandler) {
			for {
				select {
				case <-tick:
					s.pullJobFromQueue(queueName)
				case <-qHandler.stopQueue:
					logr.Debugf("%s queue stopped", queueName)
					return
				}
			}
		}(queueName, qHandler)
	}
	return nil
}

// FIXME MLSS Change: v_1.4.1 added method submit job to k8s
func (s *lcmService) deployJobToK8s(tr *trainer.TrainingRecord, logr *logger.LocLoggingEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// totalTrainingCounter.With("framework", tr.ModelDefinition.Framework.Name).Add(1)
	err := updateJobStatus(tr.TrainingID, grpc_storage.Status_PENDING, tr.UserID, service.StatusMessages_NORMAL_OPERATION.String(), client.ErrCodeNormal, logr)
	if err != nil {
		logr.WithError(err).Errorf("(deployDistributedTrainingJob) Before deploying job, error while calling Trainer service client update for trainingID %s , but still carrying on ", tr.TrainingID)
	}

	//tr.TrainingStatus.Status = grpc_trainer_v2.Status_PENDING
	//logr.Debugf("deployJobToK8s_de tr cpu: %v", tr.Training.Resources.Cpus)
	//logr.Debugf("deployJobToK8s, tr.TFosRequest: %v", tr.TFosRequest)
	//logr.Debugf("deployJobToK8s, tr.ProxyUser: %v", tr.ProxyUser)
	//repoErr := s.repo.Store(tr)
	//if repoErr != nil {
	//	return gerrf(codes.Internal, repoErr.Error())
	//}

	jobConfig, createJobConfigErr := s.createJobConfig(tr)

	if createJobConfigErr != nil {
		return gerrf(codes.Internal, createJobConfigErr.Error())
	}
	logr.Debugf("deployJobToK8s_de jobConfig cpu: %v", jobConfig.Resources.Cpus)
	logr.Debugf("deployJobToK8s_de jobConfig tf_request: %v", jobConfig.TFosRequest)

	// FIXME MLSS Change: v_1.4.1 changed name to PutIfKeyMissingError
	logr.Debugf("deployJobToK8s tr.JobType: %v", tr.JobType)

	go s.deployJob(ctx, jobConfig, logr)
	return nil
}

func (s *lcmService) deployJob(ctx context.Context, req *service.JobDeploymentRequest, logr *logger.LocLoggingEntry) {

	//logr.Debugf("now starting to deploy learners for training job")
	if err := NewMLFlowTraining(ctx, s.k8sClient, req, logr).Start(); err != nil {
		//Deploying learner helpers has failed. So update status
		handleDeploymentFailure(s, req.Name, req.TrainingId, req.UserId, req.JobNamespace, "learner deployment", logr)
		return
	}

}

func handleDeploymentFailure(s *lcmService, dlaasJobName string, tID string,
	userID string, namesapce string, component string, logr *logger.LocLoggingEntry) {

	logr.Errorf("updating status to FAILED")
	if errUpd := updateJobStatus(tID, grpc_storage.Status_FAILED, userID, service.StatusMessages_INTERNAL_ERROR.String(), client.ErrCodeFailedDeploy, logr); errUpd != nil {
		logr.WithError(errUpd).Errorf("after failed %s, error while calling Trainer service client update", component)
	}

	//Cleaning up resources out of an abundance of caution
	logr.Errorf("training FAILED so going ahead and cleaning up resources")
	if errKill := s.killDeployedJob(dlaasJobName, tID, userID, namesapce); errKill != nil {
		logr.WithError(errKill).Errorf("after failed %s, problem calling KillDeployedJob for job ", component)
	}

}

//Stops a currently executing training job
func (s *lcmService) HaltTrainingJob(ctx context.Context, req *service.JobHaltRequest) (*service.JobHaltResponse, error) {

	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId))
	logr.Debugf("Halting training job: %s", req.TrainingId)

	path := req.TrainingId + "/halt"
	// FIXME MLSS Change: v_1.4.1 changed name to PutIfKeyMissingError
	success, PutIfKeyMissingError := s.etcdClient.PutIfKeyMissing(path, "", logr)
	if PutIfKeyMissingError != nil {
		logr.WithError(PutIfKeyMissingError).Errorf("Failed to update the halt training job status on path %s for training job %s", path, req.TrainingId)
		return nil, PutIfKeyMissingError
	}
	if !success {
		logr.Warnf("While updating halt for training job %s at path %s , the path already exists", req.TrainingId, path)
	}
	// counter.With(progress, "etcdKeysDeleted").Add(1)

	return &service.JobHaltResponse{}, nil
}

// FIXME MLSS Change: v_1.4.1 added method startQueues
func (s *lcmService) startQueues() {
	logr := logger.LocLogger(logEntry())
	logr.Debugf("v1.4.0_newService startQueues")
	// ensure only one thread per lcm pulling jobs
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for queueName := range s.queues {
		s.startQueue(queueName, logr)
	}
}

//Kills a currently executing training job and cleans up its zookeeper entries
func (s *lcmService) KillTrainingJob(ctx context.Context, req *service.JobKillRequest) (*service.JobKillResponse, error) {
	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId))

	logr.Debugf("Killing training job: %s, namespace: %s", req.Name, req.JobNamespace)

	selector := "training_id==" + req.TrainingId
	backgroundPropagation := metav1.DeletePropagationBackground
	backgroundDeleteOpts := &metav1.DeleteOptions{
		PropagationPolicy: &backgroundPropagation,
	}

	logr.Debugf(" Checking if there are kubernetes services associated with training job %s", req.TrainingId)
	svcs, err := s.k8sClient.CoreV1().Services(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Services for job with name '%s' found by querying kubernetes.", req.Name)
		for _, svc := range svcs.Items {
			logr.Debugf(" Deleting service '%s'", svc.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().Services(req.JobNamespace).Delete(svc.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes service '%s' failed", svc.ObjectMeta.Name)
			}
		}
	}
	// counter.With(progress, servicesDeletedPhaseComplete).Add(1)

	logr.Debugf(" Checking if there are kubernetes statefulsets associated with training job %s", req.TrainingId)
	sets, err := s.k8sClient.AppsV1beta1().StatefulSets(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Stateful for job with name '%s' found by querying kubernetes.", req.Name)
		for _, set := range sets.Items {
			logr.Debugf(" Deleting stateful '%s'", set.ObjectMeta.Name)
			err := s.k8sClient.AppsV1beta1().StatefulSets(req.JobNamespace).Delete(set.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes stateful '%s' failed", set.ObjectMeta.Name)
			}
		}
	}

	logr.Debugf(" Checking if there are kubernetes Job associated with training job %s", req.TrainingId)
	//sets, err := s.k8sClient.AppsV1beta1().StatefulSets(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	//err := s.k8sClient.BatchV1().RESTClient().Delete().Namespace(j.Namespace).Name(j.Name).Do().Error()
	if err == nil {
		logr.Debugf(" Job for job with name '%s' found by querying kubernetes.", req.TrainingId)
		go func() {
			//todo if job failed, job will be deleted after 20s
			time.Sleep(20 * time.Second)
			logr.Debugf(" Deleting Job '%s'", req.TrainingId)
			err := s.k8sClient.BatchV1().Jobs(req.JobNamespace).Delete(req.TrainingId, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes Job '%s' failed", req.TrainingId)
			}
		}()

	}

	logr.Debugf(" Checking if there are kubernetes learner persistent volume claims associated with training job %s", req.TrainingId)
	claims, err := s.k8sClient.CoreV1().PersistentVolumeClaims(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		for _, claim := range claims.Items {
			logr.Debugf(" Deleting persistent volume claim '%s'", claim.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().PersistentVolumeClaims(req.JobNamespace).Delete(claim.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes persistent volume '%s' failed", claim.ObjectMeta.Name)
			}
		}
	}
	// counter.With(progress, pvsDeletedPhaseComplete).Add(1)

	logr.Debugf(" Checking if there are kubernetes learner COS mount secrets associated with training job %s", req.TrainingId)
	secrets, err := s.k8sClient.CoreV1().Secrets(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		for _, secret := range secrets.Items {
			logr.Debugf(" Deleting Secret '%s'", secret.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().Secrets(req.JobNamespace).Delete(secret.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes Secret '%s' failed", secret.ObjectMeta.Name)
			}
		}
	}
	// counter.With(progress, secretsDeletedPhaseComplete).Add(1)

	logr.Debugf(" Checking if there are kubernetes deployments associated with training job %s", req.TrainingId)
	deploys, err := s.k8sClient.AppsV1beta1().Deployments(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Deployments for job with name '%s' found by querying kubernetes.", req.Name)
		for _, deploy := range deploys.Items {
			logr.Debugf(" Deleting deployment '%s'", deploy.ObjectMeta.Name)
			err := s.k8sClient.AppsV1beta1().Deployments(req.JobNamespace).Delete(deploy.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes deployment '%s' failed", deploy.ObjectMeta.Name)
			}
		}
	}

	// counter.With(progress, deploymentsDeletedPhaseComplete).Add(1)

	//After Deleting the application, delete the etcd directory.
	deleteKeyWithOptsErr := s.etcdClient.DeleteKeyWithOpts(req.TrainingId, logr, clientv3.WithPrefix())
	if deleteKeyWithOptsErr != nil {
		logr.WithError(err).Errorf("failed to s.etcdClient.DeleteKeyWithOpts(req.TrainingId, logr, clientv3.WithPrefix())")
	}
	// counter.With(progress, etcdKeysDeletedPhaseComplete).Add(1)

	// FIXME MLSS Change: delete model.zip from S3 after the job is killed
	//trainerClient, err := client.NewTrainer()
	storage, err := client.NewStorage()

	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
		//return error500(logr, "")
	}
	//defer trainerClient.Close()
	defer storage.Close()
	//err = trainer.Client().DeleteSubmittedCode(req, logr)

	logr.Debugf(" DeleteSubmittedCode in service_impl start,TrainingId:  %s", req.TrainingId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err = storage.Client().DeleteSubmittedCode(ctx, &grpc_storage.DeleteRequest{
		TrainingId: req.TrainingId,
		UserId:     req.UserId,
	})

	if err != nil {
		logr.WithError(err).Errorf("DeleteSubmittedCode failed")
	}

	// FIXME MLSS Change: TFJob, Kills a currently executing training job through arena
	// todo delete TFJob
	logr.Debugf(" tfjob.DeleteTrainingJob:  %s", req.TrainingId)
	err = tfjob.DeleteTrainingJob(req.TrainingId, req.JobNamespace, tfjob.TFJob)
	if err != nil {
		logr.WithError(err).Errorf("kills a currently executing training job through arena failed")
	}

	return &service.JobKillResponse{}, nil
}

// FIXME MLSS Change: TFJob, Kills a currently executing training job and cleans up
func (s *lcmService) KillTrainingTFJob(ctx context.Context, req *service.JobKillRequest) (*service.JobKillResponse, error) {
	// counter := finishedTrainingCounter.With(outcome, killed)
	// counter.With(progress, started).Add(1)
	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId))

	logr.Debugf("Killing training job: %s, namespace: %s", req.Name, req.JobNamespace)

	selector := "training_id==" + req.TrainingId
	backgroundPropagation := metav1.DeletePropagationBackground
	backgroundDeleteOpts := &metav1.DeleteOptions{
		PropagationPolicy: &backgroundPropagation,
	}

	logr.Debugf(" Checking if there are kubernetes services associated with training job %s", req.TrainingId)
	svcs, err := s.k8sClient.CoreV1().Services(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Services for job with name '%s' found by querying kubernetes.", req.Name)
		for _, svc := range svcs.Items {
			logr.Debugf(" Deleting service '%s'", svc.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().Services(req.JobNamespace).Delete(svc.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes service '%s' failed", svc.ObjectMeta.Name)
			}
		}
	}
	// counter.With(progress, servicesDeletedPhaseComplete).Add(1)

	logr.Debugf(" Checking if there are kubernetes statefulsets associated with training job %s", req.TrainingId)
	sets, err := s.k8sClient.AppsV1beta1().StatefulSets(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Stateful for job with name '%s' found by querying kubernetes.", req.Name)
		for _, set := range sets.Items {
			logr.Debugf(" Deleting stateful '%s'", set.ObjectMeta.Name)
			err := s.k8sClient.AppsV1beta1().StatefulSets(req.JobNamespace).Delete(set.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes stateful '%s' failed", set.ObjectMeta.Name)
			}
		}
	}

	logr.Debugf(" Checking if there are kubernetes learner persistent volume claims associated with training job %s", req.TrainingId)
	claims, err := s.k8sClient.CoreV1().PersistentVolumeClaims(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		for _, claim := range claims.Items {
			logr.Debugf(" Deleting persistent volume claim '%s'", claim.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().PersistentVolumeClaims(req.JobNamespace).Delete(claim.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes persistent volume '%s' failed", claim.ObjectMeta.Name)
			}
		}
	}

	logr.Debugf(" Checking if there are kubernetes learner COS mount secrets associated with training job %s", req.TrainingId)
	secrets, err := s.k8sClient.CoreV1().Secrets(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		for _, secret := range secrets.Items {
			logr.Debugf(" Deleting Secret '%s'", secret.ObjectMeta.Name)
			err := s.k8sClient.CoreV1().Secrets(req.JobNamespace).Delete(secret.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes Secret '%s' failed", secret.ObjectMeta.Name)
			}
		}
	}

	logr.Debugf(" Checking if there are kubernetes deployments associated with training job %s", req.TrainingId)
	deploys, err := s.k8sClient.AppsV1beta1().Deployments(req.JobNamespace).List(metav1.ListOptions{LabelSelector: selector})
	if err == nil {
		logr.Debugf(" Deployments for job with name '%s' found by querying kubernetes.", req.Name)
		for _, deploy := range deploys.Items {
			logr.Debugf(" Deleting deployment '%s'", deploy.ObjectMeta.Name)
			err := s.k8sClient.AppsV1beta1().Deployments(req.JobNamespace).Delete(deploy.ObjectMeta.Name, backgroundDeleteOpts)
			if err != nil {
				logr.WithError(err).Errorf(" Deleting kubernetes deployment '%s' failed", deploy.ObjectMeta.Name)
			}
		}
	}

	//After Deleting the application, delete the etcd directory.
	deleteKeyWithOptsErr := s.etcdClient.DeleteKeyWithOpts(req.TrainingId, logr, clientv3.WithPrefix())
	if deleteKeyWithOptsErr != nil {
		logr.WithError(err).Errorf("failed to s.etcdClient.DeleteKeyWithOpts(req.TrainingId, logr, clientv3.WithPrefix())")
	}

	// FIXME MLSS Change: delete model.zip from S3 after the job is killed
	//trainerClient, err := client.NewTrainer()
	storageCient, err := client.NewStorage()
	if err != nil {
		logr.WithError(err).Errorf("Cannot create client for trainer service")
	}
	defer storageCient.Close()

	logr.Debugf(" DeleteSubmittedCode in service_impl start,TrainingId:  %s", req.TrainingId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err = storageCient.Client().DeleteSubmittedCode(ctx, &grpc_storage.DeleteRequest{
		TrainingId: req.TrainingId,
		UserId:     req.UserId,
	})

	if err != nil {
		logr.WithError(err).Errorf("DeleteSubmittedCode failed")
	}

	// FIXME MLSS Change: TFJob, Kills a currently executing training job through arena
	// todo delete TFJob
	logr.Debugf(" tfjob.DeleteTrainingJob:  %s", req.TrainingId)
	err = tfjob.DeleteTrainingJob(req.TrainingId, req.JobNamespace, tfjob.TFJob)
	if err != nil {
		logr.WithError(err).Errorf("kills a currently executing training job through arena failed")
	}

	return &service.JobKillResponse{}, nil
}

//Wrapper function for LCM's KillTrainingJob
func (s *lcmService) killDeployedJob(jobName string, trainingID string, userID string, jobNamespace string) error {
	job := &service.JobKillRequest{Name: string(jobName), TrainingId: trainingID, UserId: userID, JobNamespace: jobNamespace}
	logr := logger.LocLogger(InitLogger(trainingID, userID))
	logr.Debugf("Killing job request: %+v", job)
	_, err := s.KillTrainingJob(context.Background(), job)
	return err
}

//manages a DLaaS training job
func deployJobMonitor(s *lcmService, req *service.JobDeploymentRequest, trainingID string, numLearners int, jobName string, userID string, useNativeDistribution bool, logr *logger.LocLoggingEntry) error {

	// FIXME MLSS Change: get namespace annotations
	namespaceObj, err := s.k8sClient.CoreV1().Namespaces().Get(req.JobNamespace, metav1.GetOptions{})
	var nodeSelectors string = ""
	if nil == err {
		namespaceAnnotations := (namespaceObj.ObjectMeta).GetAnnotations()
		selector, ok := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]
		if ok {
			logr.Debugf("namespace %s annotations selector: %+v", req.JobNamespace, selector)
			nodeSelectors = selector
		}
	}

	logr.Debugf("deployJobMonitor_jobReq: %v", req)
	envVars, jmLabels := populateJobMonitorEnvVariablesAndLabels(req, trainingID, jobName, userID, numLearners, useNativeDistribution)
	//deploySpec := defineJobMonitorDeployment(req, envVars, jmLabels, logr)
	// FIXME MLSS Change: add nodeSelectors from namespace annotations to jm pods
	deploySpec := defineJobMonitorDeployment(req, envVars, jmLabels, nodeSelectors, logr)

	return backoff.RetryNotify(func() error {
		_, err := s.k8sClient.AppsV1beta1().Deployments(req.JobNamespace).Create(deploySpec)
		if k8serrors.IsAlreadyExists(err) {
			logr.WithError(err).Warnf("deployment %s already exists", deploySpec.ObjectMeta.Name)
			return nil
		}
		return err
	}, k8sInteractionBackoff(), func(err error, window time.Duration) {
		logr.WithError(err).Errorf("Could not connect to learner kube cluster and/or deploy job monitor for Training Job %s. Problem may either be in reaching the kubernetes API server or in creating a deployment", req.TrainingId)
		// k8sFailureCounter.With(component, "jobmonitor").Add(1)
	})
}

// FIXME MLSS Change: v_1.4.1 added
func getTrainedModelsBucket() string {
	if viper.IsSet(trainedModelsBucketKey) {
		return viper.GetString(trainedModelsBucketKey)
	}
	return defaultTrainedModelsBucket
}

// FIXME MLSS Change: v_1.4.1 added
func (s *lcmService) rateLimitTrainingJob(tr *trainer.TrainingRecord, quotaList *v1.ResourceQuotaList, nodeList *v1.NodeList, logr *logger.LocLoggingEntry) bool {
	logr.Debugf("rateLimitTrainingJob start to check resources limit TrainingID: %v", tr.TrainingID)
	cpu := tr.Training.Resources.Cpus
	gpu := tr.Training.Resources.Gpus
	memory := tr.Training.Resources.Memory
	requestMemoryUnit := tr.Training.Resources.MemoryUnit

	var requestMEMInByte float64
	if requestMemoryUnit.String() == "Mib" {
		requestMEMInByte = (float64(memory)) * 1024 * 1024
	} else if requestMemoryUnit.String() == "GB" {
		requestMEMInByte = (float64(memory)) * 1024 * 1024 * 1024
	} else {
		///default unit is MB, check out func convertMemoryFromManifest in manifest.go
		requestMEMInByte = (float64(memory)) * 1000 * 1000
	}

	//FIXME MLSS Change: v_1.5.2 added to get ps param
	pss, b, pssErr := s.stringToFloat64(tr.PSs, logr)
	if nil != pssErr || b == false {
		logr.WithError(pssErr).Warnf("failed to parse pss to faloat64")
		return false
	}
	psCPU, c, psCPUErr := s.stringToFloat64(tr.PSCPU, logr)
	if nil != psCPUErr || c == false {
		logr.WithError(psCPUErr).Warnf("failed to parse psCPU to faloat64")
		return false
	}
	psMemory, m, psMemoryErr := s.stringToFloat64(tr.PSMemory, logr)
	if nil != psMemoryErr || m == false {
		logr.WithError(psMemoryErr).Warnf("failed to parse psMemory to faloat64")
		return false
	}

	learners := float64(tr.Training.Resources.Learners)
	if learners <= 0 {
		learners = 1
	}
	//learners := tr.Training.Resources.Learners
	logr.Debugf("rateLimitTrainingJob debug for parse to get pss: %v, psCPU: %v, psMemory: %v, learners: %v", pss, psCPU, psMemory, learners)

	isLimit := true
	var countsGPU float64 = 0
	var countsCPU float64 = 0
	var countsMemory float64 = 0
	//nodePodsLimit := true
	var availableNodePodsCount float64 = 0
	for _, node := range nodeList.Items {

		labels := node.ObjectMeta.Labels
		//logr.Debugf("methodName_rateLimitTrainingJob node.ObjectMeta.Labels: %v", labels)
		hardwareType := labels[hardwareType]

		var nodeCpusRequested, nodeMemRequested, nodeGpuRequested float64 = 0.0, 0.0, 0.0
		//gets the node's allocatable resources
		availableResources := node.Status.Allocatable
		cpuQty := availableResources["cpu"]
		nodeCpu, _ := strconv.ParseFloat(cpuQty.AsDec().String(), 64)
		memQty := availableResources["memory"]
		nodeMem, _ := strconv.ParseFloat(memQty.AsDec().String(), 64)
		var gpuQty resource.Quantity
		var nodeGpu float64
		if nvidiagpu == hardwareType {
			gpuQty = availableResources["nvidia.com/gpu"]
			nodeParseFloatGpu, _ := strconv.ParseFloat(gpuQty.AsDec().String(), 64)
			nodeGpu = nodeParseFloatGpu
		}
		podsQty := availableResources["pods"]
		nodePods, _ := strconv.ParseFloat(podsQty.AsDec().String(), 64)
		logr.Debugf("rateLimitTrainingJob debug for node availableResources nodeCpu: %v, nodeMem: %v, nodeGpu: %v, nodePods: %v", nodeCpu, nodeMem, nodeGpu, nodePods)

		nodeName := node.Name
		pods, podListErr := s.k8sClient.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + nodeName,
		})
		if podListErr != nil {
			logr.WithError(podListErr).Warnf("failed to get podList")
			return true
		}
		//determine the pods ceiling in node
		//join the queue if the number of pods is not sufficient for a job
		//if nodePods-(learners+jobHelperNumber+jobMonitorNumber) >= float64(len(pods.Items)) {
		//	nodePodsLimit = false
		//}

		availableNodePodsCount = availableNodePodsCount + (nodePods - float64(len(pods.Items)))
		logr.Debugf("rateLimitTrainingJob debug for rateLimitTrainingJob nodeName: %v,  podList.Items.length: %v, nodePods: %v, availableNodePodsCount: %v", nodeName, len(pods.Items), nodePods, availableNodePodsCount)

		//determine the remaining resources in node
		for _, pod := range pods.Items {
			logr.Debugf("Inspecting pod %s", pod.ObjectMeta.Name)
			containers := pod.Spec.Containers
			for j := 0; j < len(containers); j++ {
				resourcesRequested := containers[j].Resources.Requests
				cpuQty := resourcesRequested["cpu"]
				pCpu, _ := strconv.ParseFloat(cpuQty.AsDec().String(), 64)
				memQty := resourcesRequested["memory"]
				pMem, _ := strconv.ParseFloat(memQty.AsDec().String(), 64)
				gpuQty := resourcesRequested["nvidia.com/gpu"]
				pGpu, _ := strconv.ParseFloat(gpuQty.AsDec().String(), 64)
				nodeCpusRequested += pCpu
				nodeMemRequested += pMem
				nodeGpuRequested += pGpu
			}
		}
		logr.Debugf("rateLimitTrainingJob debug for nod's pods available resources , nodeName: %v, cpus: %v, "+
			"memory: %v, gpu: %v", nodeName, nodeCpu-nodeCpusRequested, nodeMem-nodeMemRequested, nodeGpu-nodeGpuRequested)
		logr.Debugf("rateLimitTrainingJob debug for request resources memory: %v, cpu: %v, gpu: %v, learners:"+
			" %v", nodeCpusRequested, nodeMemRequested, nodeGpuRequested, learners)
		logr.Debugf("rateLimitTrainingJob debug for request resources memory: %v, cpu: %v, gpu: %v, learners: %v", memory, cpu, gpu, learners)
		logr.Debugf("rateLimitTrainingJob debug for nodeCpu-nodeCpusRequested >= cpu: %v,"+
			" nodeMem-nodeMemRequested >= float64(memory*1024*1024*1024): %v, "+
			"nodeGpu-nodeGpuRequested >= float64(gpu): %v, learners: %v", nodeCpu-nodeCpusRequested >= float64(cpu),
			nodeMem-nodeMemRequested >= float64(memory)*1024*1024*1024, nodeGpu-nodeGpuRequested >= float64(gpu), learners)

		if nodeCpu-nodeCpusRequested >= float64(cpu) && nodeMem-nodeMemRequested >= requestMEMInByte &&
			nodeGpu-nodeGpuRequested >= float64(gpu) && nodeCpu-nodeCpusRequested >= psCPU &&
			nodeMem-nodeMemRequested >= psMemory*1024*1024*1024 {
			logr.Debugf("node has sufficient cpu, gpu, memory resources")
			isLimit = false
		} else {
			logr.Debugf("node has insufficient cpu, gpu, memory resources")
			logr.Debugf("insufficient cpu: %v", nodeCpu-nodeCpusRequested >= float64(cpu))
			logr.Debugf("insufficient cpu: %v", nodeCpu-nodeCpusRequested)
			logr.Debugf("insufficient cpu: %v", float64(cpu))
			logr.Debugf("insufficient mem: %v", nodeMem-nodeMemRequested >= requestMEMInByte)
			logr.Debugf("insufficient mem: %v", nodeMem-nodeMemRequested)
			logr.Debugf("insufficient mem: %v", requestMEMInByte)
			logr.Debugf("insufficient gpu: %v", nodeGpu-nodeGpuRequested >= float64(gpu))
			logr.Debugf("insufficient gpu: %v", float64(gpu))
			logr.Debugf("insufficient psCPU: %v", nodeCpu-nodeCpusRequested >= psCPU)
			logr.Debugf("insufficient psCPU: %v", psCPU)
			logr.Debugf("insufficient memory: %v", nodeMem-nodeMemRequested >= psMemory*1024*1024*1024)
			logr.Debugf("insufficient memory: %v", psMemory*1024*1024*1024)
		}

		if gpu > 0 {
			countsGPU = countsGPU + math.Floor((nodeGpu-nodeGpuRequested)/float64(gpu))
		}
		if cpu > 0 {
			countsCPU = countsCPU + math.Floor((nodeCpu-nodeCpusRequested)/float64(cpu))
		}
		if memory > 0 {
			countsMemory = countsMemory + math.Floor((nodeMem-nodeMemRequested)/requestMEMInByte)
		}
	}

	logr.Debugf("rateLimitTrainingJob debug for isLimit: %v", isLimit)
	logr.Debugf("rateLimitTrainingJob debug for countsGPU: %v, countsCPU: %v, countsMemory: %v, learners: %v", countsGPU, countsCPU, countsMemory, learners)

	if availableNodePodsCount < learners+jobHelperNumber+jobMonitorNumber {
		logr.Debugf("node has not enough pod resources for availableNodePodsCount: %v, reqPods: %v", availableNodePodsCount, learners+jobHelperNumber+jobMonitorNumber)
		return true
	}

	if isLimit {
		return true
	}
	if gpu > 0 {
		if countsGPU < learners || countsCPU < learners || countsMemory < learners {
			return true
		}
	} else {
		if countsCPU < learners || countsMemory < learners {
			return true
		}
	}

	rq := quotaList.Items[0]
	hard := rq.Status.Hard
	used := rq.Status.Used
	hardLimitsCPU := hard[LimitsCpu]
	usedLimitsCPU := used[LimitsCpu]
	hardLimitsGPU := hard[RequestsNvidiaComGpu]
	usedLimitsGPU := used[RequestsNvidiaComGpu]
	hardLimitsMemory := hard[LimitsMemory]
	usedLimitsMemory := used[LimitsMemory]

	//isEnough, e := CheckResource(cpu, gpu, memory, learners, hardLimitsCPU, hardLimitsGPU, hardLimitsMemory, logr)
	//if !isEnough {
	//	logr.Debugf("cpu, gpu or memory is limit for need cpu: v% with hardLimitsCPU: v%, gpu: v% with hardLimitsGPU: v%, memory: v% with hardLimitsMemory: v%", cpu, hardLimitsCPU, gpu, hardLimitsGPU, memory, hardLimitsMemory)
	//	return true
	//}

	logr.Debugf("lcmService rateLimitTrainingJob cpu: %v hardLimitsCPU: %v gpu: %v hardLimitsGPU: %v usedLimitsCPU: %v memory: %v(requestMEMInByte: %v)", cpu, hardLimitsCPU, gpu, hardLimitsGPU, usedLimitsCPU, memory, requestMEMInByte)
	//check cpu is limit
	availableCPU, cpuErr := CheckAvailableCPU(float64(cpu)*learners+(JobmonitorCpu+HelperCpu*2)/1000+pss*psCPU, hardLimitsCPU.String(), usedLimitsCPU.String())
	if cpuErr != nil {
		logr.WithError(cpuErr).Warnf("failed to CheckAvailableCPU: %v", tr.Namespace)
		return true
	}
	if *availableCPU < 0 {
		return true
	}
	//check gpu is limit
	//FIXME MLSS Change: v_1.5.2 added
	availableGPU, gpuErr := CheckAvailableGPU(float64(gpu)*learners, hardLimitsGPU.String(), usedLimitsGPU.String())
	if gpuErr != nil {
		logr.WithError(gpuErr).Warnf("failed to CheckAvailableGPU: %v", tr.Namespace)
		return true
	}
	if *availableGPU < 0 {
		return true
	}
	//check memory is limit
	mem, memoryErr := CheckAvailableMemory(requestMEMInByte*learners+float64((JobmonitorMemory+HelperMemory*2)/(1024*1024*1024))+pss*psMemory, hardLimitsMemory.String(), usedLimitsMemory.String())
	if memoryErr != nil {
		logr.WithError(memoryErr).Warnf("memory can not greater than the limit: %v", tr.Namespace)
		return true
	}
	if *mem < 0 {
		return true
	}
	return false
}

// FIXME MLSS Change: v_1.4.1 added
func (s *lcmService) pullJobFromQueue(gpuType string) {
	logr := logger.LocLogger(logEntry())

	qHandler := s.queues[gpuType]
	if qHandler == nil {
		logr.Warnf("there is no queue for type %v", gpuType)
		return
	}

	locked := true
	qerr := qHandler.Lock()
	if qerr != nil {
		logr.WithError(qerr).Errorf("failed to lock %v queue", gpuType)
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

	if trainingRecord.Deleted ||
		trainingRecord.TrainingStatus.Status == grpc_trainer_v2.Status_KILLING ||
		trainingRecord.TrainingStatus.Status == grpc_trainer_v2.Status_CANCELLED {
		logr.Debugf("job %s was deleted or cancel", nextJobID)
		qHandler.Delete(nextJobID)
		return
	}
	if trainingRecord.TrainingStatus.Status != grpc_trainer_v2.Status_QUEUED {
		logr.Warnf("job %s expected status QUEUED but found %s, removing job from queue", nextJobID, trainingRecord.TrainingStatus)
		qHandler.Delete(nextJobID)
		return
	}

	logr.Debugf("got training job %s from %s queue", nextJobID, gpuType)

	quotaList, listE := s.getRQList(logr, trainingRecord.Namespace)
	if listE != nil {
		logr.WithError(listE).Errorln("failed to getRQList")
		return
	}

	nodeList, getNodeListErr := s.getNodeList(trainingRecord, logr)
	if getNodeListErr != nil {
		logr.WithError(listE).Errorln("failed to getNodeList")
		return
	}
	//logr.Debugf("getResources for nodeList.Items: %v", nodeList.Items)
	if nodeList == nil {
		logr.WithError(listE).Errorln("namespace is not yet tied to any node")
		return
	}

	rateLimited := s.rateLimitTrainingJob(trainingRecord, quotaList, nodeList, logr)
	if rateLimited {
		logr.Debugf("training job %s is rate-limited, leaving in %s queue", trainingRecord.TrainingID, gpuType)
		return
	}
	// FIXME MLSS Change: v_1.4.1 added
	err = s.deployJobToK8s(trainingRecord, logr)
	// FIXME 1.9.0 lcm
	//err = s.newDeployJobToK8s(trainingRecord, logr)
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
				if dequeuedTrainingRecord.TrainingStatus.Status == grpc_trainer_v2.Status_QUEUED {
					_, err = updateTrainingJobPostLock(s, &grpc_trainer_v2.UpdateRequest{
						TrainingId:    dequeuedJob,
						UserId:        dequeuedTrainingRecord.UserID,
						Status:        grpc_trainer_v2.Status_FAILED,
						StatusMessage: "Job was dequeued without being submitted",
						ErrorCode:     client.ErrCodeFailDequeue,
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
	timestamp := client.CurrentTimestampAsString()
	e := &trainer.JobHistoryEntry{
		TrainingID:    trainingRecord.TrainingID,
		Timestamp:     timestamp,
		Status:        trainingRecord.TrainingStatus.Status,
		StatusMessage: trainingRecord.TrainingStatus.StatusMessage,
		ErrorCode:     trainingRecord.TrainingStatus.ErrorCode,
	}
	s.jobHistoryRepo.RecordJobStatus(e)
}

// This method contains all the functionality of UpdateTrainingJob, minus the lock on the database.  This enables it to be called
// from within another function, which already has the lock itself (Halt)
// FIXME MLSS Change: v_1.4.1 added
func updateTrainingJobPostLock(s *lcmService, req *grpc_trainer_v2.UpdateRequest) (*grpc_trainer_v2.UpdateResponse, error) {
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

	nowMillis := client.CurrentTimestampAsString()

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
		if req.Status == grpc_trainer_v2.Status_FAILED {
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
		e := &trainer.JobHistoryEntry{
			TrainingID:    req.TrainingId,
			Timestamp:     timestamp,
			Status:        req.Status,
			StatusMessage: req.StatusMessage,
			ErrorCode:     req.ErrorCode,
		}
		recordJobStatusErr := s.jobHistoryRepo.RecordJobStatus(e)
		if recordJobStatusErr != nil {
			return nil, recordJobStatusErr
		}
	}

	return &grpc_trainer_v2.UpdateResponse{TrainingId: training.TrainingID}, nil
}

// FIXME MLSS Change: v_1.4.1 added
func (s *lcmService) createJobConfig(tr *trainer.TrainingRecord) (*service.JobDeploymentRequest, error) {
	logr := logger.LocLogger(logWith(tr.TrainingID, tr.UserID))

	user := tr.UserID
	if tr.ProxyUser != "" {
		user = tr.ProxyUser
	}

	// training data/results - assume only one training input and output data at this point
	trainingData := trainer.FindDatastore(tr.Training.InputData[0], tr.Datastores)
	trainingResults := s.getOutputDatastore(tr.Training.OutputData, tr.Datastores)
	// FIXME MLSS Change: v_1.5.1 added
	trainingWorkspace := s.getCodeWorkDatastore(tr.Training.WorkData, tr.Datastores)
	logr.Debugf("debug_volumesForLearner trainingWorkspace: %v", trainingWorkspace)
	// Environment variables
	envvars := make(map[string]string)
	// FIXME MLSS Change: v_1.5.1 added workData to envvars
	if nil != trainingWorkspace {
		envvars["WORK_DIR"] = trainingWorkspace.Fields["codeWork"]
		envvars["WORK_STORE_TYPE"] = trainingWorkspace.Connection["type"]
		logr.Debugf("debug_volumesForLearner workStoreType: %v", trainingWorkspace.Connection["type"])
		if trainingWorkspace.Connection["region"] != "" {
			envvars["WORK_STORE_REGION"] = trainingWorkspace.Connection["region"]
		}
		if trainingData.Connection["path"] != "" {
			envvars["WORK_STORE_PATH"] = trainingData.Connection["path"]
		}
	}

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

	// Storing model in container at
	envvars["MODEL_DIR"] = "/model-code"

	// Storing trained model at
	envvars["RESULT_DIR"] = trainingResults.Fields["bucket"]

	// FIXME MLSS Change: v_1.4.1
	envvars["JOB_ALERT"] = tr.JobAlert

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
	envvars["USER_ID"] = user
	envvars["DLAAS_MLSSGID"] = viper.GetString(config.MLSSGroupId)
	//envvars["EVENT_CHECKER"] = tr.EventChecker
	//envvars["DEADLINE_CHECKER"] = tr.DeadlineChecker
	//envvars["OVERTIME_CHECKER"] = tr.OvertimeChecker
	//envvars["ALERT_LEVEL"] = tr.AlertLevel
	//envvars["RECEIVER"] = tr.Receiver
	//envvars["INTERVAL"] = tr.Interval

	envvars["TFOS_REQUEST"] = tr.TFosRequest

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

	job := &service.JobDeploymentRequest{
		Name:                  u4.String(),
		Resources:             trainer.GetResourceRequirements(tr.Training),
		EnvVars:               envvars,
		Labels:                labels,
		UserId:                user,
		TrainingId:            tr.TrainingID,
		Framework:             tr.ModelDefinition.Framework.Name,
		Version:               tr.ModelDefinition.Framework.Version,
		ImageLocation:         trainer.ParseImageLocation(tr),
		EvaluationMetricsSpec: tr.EvaluationMetricsSpec,
		JobNamespace:          tr.Namespace,
		// FIXME MLSS Change: add DataStore & ModelDefinition & Training from tr to JobDeploymentRequest
		DataStores:   tr.Datastores,
		JobAlert:     tr.JobAlert,
		PSs:          tr.PSs,
		PSCPU:        tr.PSCPU,
		PSImage:      tr.PSImage,
		PSMemory:     tr.PSMemory,
		CodeSelector: tr.CodeSelector,
		DataPath:     tr.DataPath,
		JobType:      tr.JobType,
		TFosRequest:  tr.TFosRequest,
		ProxyUser:    tr.ProxyUser,
		DataSet:      tr.DataSet,
		MfModel:      tr.MFModel,
		Algorithm:    tr.Algorithm,
		JobParams:    tr.JobParams,
		FitParams:    tr.FitParams,
		APIType:      tr.APIType,
	}

	return job, nil
}

// FIXME MLSS Change: v_1.4.1 added
func (s *lcmService) getOutputDatastore(outputData []string, datastores []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	var ds *grpc_trainer_v2.Datastore
	if len(outputData) > 0 {
		ds = trainer.FindDatastore(outputData[0], datastores) // we assume there is only one output data at this point b/c the underlying system does not support more
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

// FIXME MLSS Change: v_1.5.1 added
func (s *lcmService) getCodeWorkDatastore(workPath []string, datastores []*grpc_trainer_v2.Datastore) *grpc_trainer_v2.Datastore {
	var ds *grpc_trainer_v2.Datastore
	if len(workPath) > 0 {
		ds = trainer.FindDatastore(workPath[0], datastores) // we assume there is only one output data at this point b/c the underlying system does not support more
	}
	return ds
}

// FIXME MLSS Change: v_1.4.1 added
func CheckAvailableGPU(gpu float64, hardLimitsGPU string, usedLimitsGPU string) (*float64, error) {
	hard, e := strconv.ParseFloat(hardLimitsGPU, 64)
	if e != nil {
		return nil, e
	}
	used, e := strconv.ParseFloat(usedLimitsGPU, 64)
	if e != nil {
		return nil, e
	}
	logrus.Debugf("CheckAvailableGPU, hard: %v, used: %v, requestGpu: %v, ", hard, used, gpu)
	re := hard - used - gpu
	return &re, nil
}

// FIXME MLSS Change: v_1.4.1 added
func CheckAvailableCPU(cpu float64, hardLimitsCPU string, usedLimitsCPU string) (*float64, error) {
	hard, e := GetCPUValueWithUnits(hardLimitsCPU)
	if e != nil {
		return nil, e
	}
	used, e := GetCPUValueWithUnits(usedLimitsCPU)
	if e != nil {
		return nil, e
	}
	logrus.Debugf("CheckAvailableCPU, hard: %v, used: %v, requestCpu: %v, ", *hard, *used, cpu*1000)
	f := *hard - *used - cpu*1000
	return &f, nil
}

// FIXME MLSS Change: v_1.4.1 added
func GetCPUValueWithUnits(withUnits string) (*float64, error) {
	f, e := strconv.ParseFloat(withUnits, 64)
	if e == nil {
		i := f * 1000
		return &i, nil
	}
	if strings.HasSuffix(withUnits, "m") {
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-1]), 64)
		if e != nil {
			return nil, e
		} else {
			i := float
			return &i, nil
		}
	}
	return nil, e
}

// FIXME MLSS Change: v_1.4.1 added
func CheckAvailableMemory(memory float64, hardLimitsMemory string, usedLimitsMemory string) (*float64, error) {
	//mem := memory * 1024 * 1024 * 1024
	mem := memory
	hard, e := GetMemoryValueWithUnits(hardLimitsMemory)
	if e != nil {
		return nil, e
	}
	used, e := GetMemoryValueWithUnits(usedLimitsMemory)
	if e != nil {
		return nil, e
	}
	logrus.Debugf("CheckAvailableMemory, hard: %v, used: %v, requestMem: %v, ", *hard, *used, mem)
	re := *hard - *used - mem
	return &re, nil
}

// FIXME MLSS Change: v_1.4.1 added
func GetMemoryValueWithUnits(withUnits string) (*float64, error) {
	f, e := strconv.ParseFloat(withUnits, 64)
	if e == nil {
		return &f, nil
	}
	if strings.HasSuffix(withUnits, "k") {
		//return strconv.ParseFloat(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024;
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-1]), 64)
		if e != nil {
			return nil, e
		} else {
			i := float * 1024
			return &i, nil
		}
	}
	if strings.HasSuffix(withUnits, "Mi") {
		//return strconv.ParseFloat(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024;
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-2]), 64)
		if e != nil {
			return nil, e
		} else {
			i := float * 1024 * 1024
			return &i, nil
		}
	}
	if strings.HasSuffix(withUnits, "Gi") {
		//return Double.parseDouble(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024 * 1024;
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-2]), 64)
		if e != nil {
			return nil, e
		} else {
			i := float * 1024 * 1024 * 1024
			return &i, nil
		}
	}
	return nil, e
}

// FIXME MLSS Change: v_1.4.1 getNodeList By labelSelector
func (s *lcmService) getNodeList(tr *trainer.TrainingRecord, log *logger.LocLoggingEntry) (*v1.NodeList, error) {
	namespace, nsErr := s.k8sClient.CoreV1().Namespaces().Get(tr.Namespace, metav1.GetOptions{})
	if nsErr != nil {
		return nil, nsErr
	}
	namespaceAnnotations := namespace.ObjectMeta.GetAnnotations()
	log.Debugf("getResources for namespaceAnnotations: %v", namespaceAnnotations)

	selector := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]
	nodeList, nodeErr := s.k8sClient.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if nodeErr != nil {
		return nil, nodeErr
	}
	return nodeList, nil
}

//FIXME MLSS Change: v_1.5.2 added
func (s *lcmService) stringToFloat64(str string, logr *logger.LocLoggingEntry) (float64, bool, error) {
	if "" != str {
		if strings.Contains(str, "Gi") {
			str = str[0:(len(str) - 2)]
		}
		s, e := strconv.ParseFloat(str, 64)
		if nil != e {
			logr.Debugf("strconv.ParseFloat param: %v", str)
			return 0, false, e
		}
		return s, true, nil
	}
	return 0, true, nil
}

// FIXME MLSS Change: v_1.4.1
func (s *lcmService) checkLimitResources(tr *trainer.TrainingRecord, quotaList *v1.ResourceQuotaList, nodeList *v1.NodeList, logr *logger.LocLoggingEntry) (bool, error) {
	logr.Debugf("rateLimitTrainingJob start to check resources limit TrainingID: %v", tr.TrainingID)
	requestCpu := tr.Training.Resources.Cpus
	requestGpu := tr.Training.Resources.Gpus

	requestPSCPU, c, psCPUErr := s.stringToFloat64(tr.PSCPU, logr)
	if nil != psCPUErr || c == false {
		logr.WithError(psCPUErr).Warnf("failed to parse requestPSCPU to faloat64")
		return false, psCPUErr
	}
	requestPSMemory, m, psMemoryErr := s.stringToFloat64(tr.PSMemory, logr)
	if nil != psMemoryErr || m == false {
		logr.WithError(psMemoryErr).Warnf("failed to parse requestPSMemory to faloat64")
		return false, psCPUErr
	}

	//FIXME MLSS Change: v_1.5.2 added
	requestLearners := float64(tr.Training.Resources.Learners)
	if requestLearners <= 0 {
		requestLearners = 1
	}
	logr.Debugf("checkLimitResources to get requestLearners: %v, get requestPSCPU: %v, requestPSMemory: %v", requestLearners, requestPSCPU, requestPSMemory)

	requestMemory := tr.Training.Resources.Memory
	requestMemoryUnit := tr.Training.Resources.MemoryUnit

	// FIXME MLSS Change: v_1.4.1
	//var cpusAllocatable, memAllocatable, gpuAllocatable float64 = 0, 0, 0

	var countsGPU float64 = 0
	var countsCPU float64 = 0
	var countsMemory float64 = 0

	var requestMEMInByte float64
	if requestMemoryUnit.String() == "Mib" {
		requestMEMInByte = (float64(requestMemory)) * 1024 * 1024
	} else if requestMemoryUnit.String() == "GB" {
		requestMEMInByte = (float64(requestMemory)) * 1024 * 1024 * 1024
	} else {
		///default unit is MB, check out func convertMemoryFromManifest in manifest.go
		requestMEMInByte = (float64(requestMemory)) * 1000 * 1000
	}

	cpuIsEnough := false
	gpuIsEnough := false
	memIsEnough := false
	errCPUMsg := ""
	errGPUMsg := ""
	errMemMsg := ""

	for _, node := range nodeList.Items {
		//FIXME MLSS Change: v_1.5.1 added logic to get GPU resources by hardwareType
		labels := node.ObjectMeta.Labels
		hardwareType := labels[hardwareType]

		availableResources := node.Status.Allocatable

		cpuQty := availableResources[v1.ResourceCPU]
		nodeCpu, _ := strconv.ParseFloat(cpuQty.AsDec().String(), 64)
		memQty := availableResources[v1.ResourceMemory]
		nodeMem, _ := strconv.ParseFloat(memQty.AsDec().String(), 64)
		var gpuQty resource.Quantity
		var nodeGpu float64
		if nvidiagpu == hardwareType {
			gpuQty = availableResources["nvidia.com/gpu"]
			nodeParseFloatGpu, _ := strconv.ParseFloat(gpuQty.AsDec().String(), 64)
			nodeGpu = nodeParseFloatGpu
			logr.Debugf("methodName_checkLimitResources hardwareType: %v, nodeGpu: %v", hardwareType, nodeGpu)
		}
		//cpusAllocatable += nodeCpu
		//memAllocatable += nodeMem
		//gpuAllocatable += nodeGpu

		logr.Debugf("checkLimitResources nodeMem: %v, requestMEMInByte: %v, psMem: %v, nodeCpu: %v, reqCpu: %v, requestPSCPU: %v, nodeGpu: %v, reqGpu: %v", nodeMem, requestMEMInByte, requestPSMemory, nodeCpu, requestCpu, requestPSCPU, nodeGpu, requestGpu)

		if nodeCpu >= float64(requestCpu) && nodeCpu >= float64(requestPSCPU) {
			cpuIsEnough = true
		}
		errCPUMsg = fmt.Sprintf("there is not enough cpu to meet you request nodeCPU: %v, reqCPU: %v, requestPSCPU: %v", nodeCpu, requestCpu, requestPSCPU)

		if nodeMem >= requestMEMInByte && nodeMem >= requestPSMemory {
			memIsEnough = true
		}
		errMemMsg = fmt.Sprintf("there is not enough memory to meet you request nodeMem: %v, requestMEMInByte: %v, requestPSMemory: %v", nodeMem, requestMEMInByte, requestPSMemory)

		if nodeGpu >= float64(requestGpu) {
			gpuIsEnough = true
		}
		errGPUMsg = fmt.Sprintf("there is not enough GPU to meet you request nodeGpu: %v, requestGpu: %v", nodeGpu, requestGpu)

		//FIXME MLSS Change: v_1.5.2 added
		if requestGpu > 0 {
			countsGPU = countsGPU + math.Floor(nodeGpu/float64(requestGpu))
		}
		if requestCpu > 0 {
			countsCPU = countsCPU + math.Floor(nodeCpu/float64(requestCpu))
		}
		if requestMemory > 0 {
			countsMemory = countsMemory + math.Floor(nodeMem/requestMEMInByte)
		}
	}

	logr.Debugf("checkLimitResources countsGPU: %v, countsCPU: %v, countsMemory: %v, requestLearners: %v, nodes of namespace len(nodeList.Items): %v)", countsGPU, countsCPU, countsMemory, requestLearners, len(nodeList.Items))

	if !cpuIsEnough {
		return false, errors.New(errCPUMsg)
	}

	if !gpuIsEnough {
		return false, errors.New(errGPUMsg)
	}

	if !memIsEnough {
		return false, errors.New(errMemMsg)
	}
	//FIXME MLSS Change: v_1.5.2 added
	if requestGpu > 0 {
		if countsGPU < requestLearners {
			return false, errors.New("there are not enough GPU resources of node to meet you resources needs, please reset the resources")
		}
	}
	if countsCPU < requestLearners {
		return false, errors.New("there are not enough CPU resources of node to meet you resources needs, please reset the resources")
	}
	if countsMemory < requestLearners {
		return false, errors.New("there are not enough Memory resources of node to meet you resources needs, please reset the resources")
	}
	//if !isEnough {
	//	return false, errors.New("there is no node that meets your resource requirements")
	//}

	//requestLearners := tr.Training.Resources.Learners
	rq := quotaList.Items[0]
	hard := rq.Status.Hard
	hardLimitsCPU := hard[LimitsCpu]
	hardLimitsGPU := hard[RequestsNvidiaComGpu]
	hardLimitsMemory := hard[LimitsMemory]

	hardCPU, e := GetCPUValueWithUnits(hardLimitsCPU.String())
	if e != nil {
		return false, e
	}
	if requestLearners*float64(requestCpu)*1000+JobmonitorCpu+HelperCpu*2 > *hardCPU {
		logr.Debugf("requestCpu is limit for need hardCPU: %v, needCpu: %v", hardCPU, requestLearners*float64(requestCpu)*1000+JobmonitorCpu+HelperCpu*2)
		return false, nil
	}
	hardGPU, e := strconv.ParseFloat(hardLimitsGPU.String(), 64)
	if e != nil {
		return false, e
	}
	if float64(requestGpu)*requestLearners > hardGPU {
		logr.Debugf("requestGpu is limit for need hardGPU: %v, needGpu: %v", float64(requestGpu)*requestLearners, hardGPU)
		return false, nil
	}
	units, e := GetMemoryValueWithUnits(hardLimitsMemory.String())
	if e != nil {
		return false, e
	}
	if float64(requestLearners*requestMEMInByte+JobmonitorMemory+HelperMemory*2) > *units {
		logr.Debugf("requestMemory is limit for need hardMemory: %v, reqMemory: %v", units, float64(requestLearners*requestMEMInByte+JobmonitorMemory+HelperMemory*2))
		return false, nil
	}
	return true, nil
}

// FIXME MLSS Change: v_1.4.1
func (s *lcmService) enQueue(qHandler *queueHandler, tr *trainer.TrainingRecord, logr *logger.LocLoggingEntry) error {
	//save trainingRecord to mongo
	enqueueErr := qHandler.Enqueue(tr.TrainingID)
	logr.Debugf("DeployTrainingJob tr enqueue with trainingID : s%", tr.TrainingID)

	if enqueueErr != nil {
		tr.TrainingStatus.Status = grpc_trainer_v2.Status_FAILED
		tr.TrainingStatus.StatusMessage = fmt.Sprintf("Job was rate-limited and could not be enqueued")
		tr.TrainingStatus.ErrorCode = client.ErrCodeFailEnqueue
		err := s.repo.Store(tr)
		if err != nil {
			logr.WithError(err).Errorln("Unable to store job with status FAILED")
			return err
		}
		return enqueueErr
	}
	//logr.Debugf("DeployTrainingJob_de tr enqueue with cpu : s%", tr.Training.Resources.Cpus)
	//logr.Debugf("DeployTrainingJob_de tr enqueue with tr.GuardianToken : s%", tr.GuardianToken)
	err := s.repo.Store(tr)
	if err != nil {
		logr.WithError(err).Errorf("Failed to resolve output datastore")
		return err
	}

	//if queue not start yet, create a goroutine
	//s.startQueue(trainer.QueueName(tr.Namespace), logr)
	return nil
}

func (s *lcmService) getRQList(logR *logger.LocLoggingEntry, namespace string) (*v1.ResourceQuotaList, error) {
	var opts = metav1.ListOptions{}
	quotaList, optsE := s.k8sClient.CoreV1().ResourceQuotas(namespace).List(opts)
	logR.Debugf("lcmService rateLimitTrainingJob listNamespaceResources v%", quotaList)
	if optsE != nil {
		logR.WithError(optsE).Warnf("failed to get quotaList of namespace v%", namespace)
		return nil, optsE
	}
	if len(quotaList.Items) < 1 {
		logR.WithError(optsE).Warnf("failed to get quotaList of namespace v%", namespace)
		return nil, errors.New("resourceQuota is not existed, please contact the administrator")
	}
	return quotaList, nil
}

func (s *lcmService) resourceChecking(req *service.JobDeploymentRequest, tr *trainer.TrainingRecord, logr *logger.LocLoggingEntry) (bool, error) {
	ns := req.JobNamespace

	quotaList, listE := s.getRQList(logr, ns)
	if listE != nil {
		logr.WithError(listE).Errorln("failed to getRQList")
		return false, listE
	}
	//get nodes by labelSelector
	nodeList, getNodeListErr := s.getNodeList(tr, logr)
	if getNodeListErr != nil {
		return false, getNodeListErr
	}
	if nodeList == nil || len(nodeList.Items) < 1 {
		return false, errors.New("namespace is not yet tied to any node, please contact the administrator")
	}

	//check the limit of ns's total resource
	isEnough, checkErr := s.checkLimitResources(tr, quotaList, nodeList, logr)
	if checkErr != nil {
		logr.WithError(checkErr).Errorln("failed to checkLimitResources" + checkErr.Error())
		return false, checkErr
	}
	if !isEnough {
		return false, errors.New("committed resources exceed the maximum limit")
	}

	//check if ns has enough resource for the job now
	rateLimited := s.rateLimitTrainingJob(tr, quotaList, nodeList, logr)
	return rateLimited, nil
}

func (s *lcmService) checkQueueAndEnqueue(handler *queueHandler, queueName string, maxQueueSize int, tr *trainer.TrainingRecord, logr *logger.LocLoggingEntry) error {
	size, err := handler.Size()
	if err != nil {
		logr.WithError(err).Errorln("FAILED to get queue size")
		return err
	}

	logGpuTypeQueueSize := fmt.Sprintf("%s_%s", queueName, "queue_size")
	logr.WithFields(logrus.Fields{
		logGpuTypeQueueSize: size,
	})
	logr.Debugf("queue %s has %d elements", queueName, size)

	logr.Infof("handler.Size: %v, viper.GetInt(queueSize): %v", size, maxQueueSize)
	if size >= maxQueueSize {
		return errors.New(fmt.Sprintf("queue is full, queue size: %v", maxQueueSize))
	}

	enErr := s.enQueue(handler, tr, logr)
	if enErr != nil {
		logr.WithError(enErr).Errorf("Failed to enQueue,when qHandler exists")
		return enErr
	}
	return nil
}

//Kills a currently executing training job and cleans up its zookeeper entries
func (s *lcmService) KillMLFlowJob(ctx context.Context, req *service.JobKillRequest) (*service.JobKillResponse, error) {
	if req.JobType == "MLPipeline" {
		return s.killMLPipelineJob(req)
	} else if req.JobType == "dist-tf" {
		return s.KillTrainingTFJob(ctx, req)
	} else { //Local
		return s.killSingleJob(req)
	}
}

func (s *lcmService) killSingleJob(req *service.JobKillRequest) (*service.JobKillResponse, error) {
	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId))
	propagation := metav1.DeletePropagationBackground
	options := metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	}

	logr.Debugf("Killing training job: %s, namespace: %s", req.Name, req.JobNamespace)
	logr.Debugf(" Deleting Job '%s'", req.TrainingId)

	err := s.k8sClient.CoreV1().Services(req.JobNamespace).Delete(req.Name, &options)
	if err != nil {
		logr.Errorf("deleting service: %v  failed, %v", req.Name, err.Error())
	}
	err = s.k8sClient.BatchV1().Jobs(req.JobNamespace).Delete(req.Name, &options)
	if err != nil {
		logr.Errorf("deleting job: %v in l8s failed, %v", req.Name, err.Error())
	}
	return &service.JobKillResponse{}, nil
}

func (s *lcmService) killDistJob(req *service.JobKillRequest) (*service.JobKillResponse, error) {

	return &service.JobKillResponse{}, nil

}

func (s *lcmService) killMLPipelineJob(req *service.JobKillRequest) (*service.JobKillResponse, error) {
	logr := logger.LocLogger(InitLogger(req.TrainingId, req.UserId))
	propagation := metav1.DeletePropagationBackground
	options := metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	}

	logr.Debugf("Killing training job: %s, namespace: %s", req.TrainingId, req.JobNamespace)
	logr.Debugf(" Deleting Job '%s'", req.TrainingId)

	err := s.k8sClient.CoreV1().Services(req.JobNamespace).Delete(req.TrainingId, &options)
	if err != nil {
		logr.Errorf("deleting service: %v  failed, %v", req.TrainingId, err.Error())
	}
	err = s.k8sClient.BatchV1().Jobs(req.JobNamespace).Delete(req.TrainingId, &options)
	if err != nil {
		logr.Errorf("deleting job: %v in l8s failed, %v", req.TrainingId, err.Error())
	}
	return &service.JobKillResponse{}, nil

}
