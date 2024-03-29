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
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
	"webank/DI/storage/storage/grpc_storage"

	"github.com/cenkalti/backoff"
	"gopkg.in/yaml.v2"

	"webank/DI/commons/config"

	"webank/DI/commons/logger"
	"webank/DI/commons/service"
	"webank/DI/commons/util"
	"webank/DI/lcm/coord"

	"webank/DI/storage/client"

	"golang.org/x/net/context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//helper function to construct a job monitor name from job name
func constructJMName(jobName string) string {
	jmName := "jobmonitor-" + jobName
	return jmName
}

//helper function to construct a learner name from job name
func constructLearnerName(learnerID int, jobName string) string {
	return "learner-" + strconv.Itoa(learnerID) + "-" + jobName
}

//helper function to construct a learnerHelper name from job name
func constructLearnerHelperName(learnerID int, jobName string) string {
	return "lhelper-" + strconv.Itoa(learnerID) + "-" + jobName
}

//helper function to construct a learner service name from job name
func constructLearnerServiceName(learnerID int, jobName string) string {
	return constructLearnerName(learnerID, jobName)
}

//helper function to construct a learner service name from job name
func constructLearnerVolumeClaimName(learnerID int, jobName string) string {
	return constructLearnerName(learnerID, jobName)
}

//helper function to construct a parameter server name from job name
func constructPSName(jobName string) string {
	psName := "grpc-ps-" + jobName
	return psName
}

// Get the disk size (in bytes) requested for a job.
func getStorageSize(r *service.ResourceRequirements) int64 {
	// The default size for all jobs
	size := config.GetVolumeSize()

	// Use the requested volume size if it's specified
	if r.Storage > 0 {
		storageSizeInBytes := int64(calcStorage(r) * 1024 * 1024)
		size = storageSizeInBytes
	}

	return size
}

// Return the name of a volume to use for a job.
func getStaticVolume(logr *logger.LocLoggingEntry) string {
	type Items struct {
		Name   string `yaml:"name"`
		Label  string `yaml:"label"`
		Status string `yaml:"status"`
	}
	type Volumes struct {
		Volumes []Items `yaml:"static-volumes-v2"`
	}

	var staticVolumes Volumes
	pvcConfigMap := "/etc/static-volumes-v2/PVCs-v2.yaml"
	bytes, err := ioutil.ReadFile(pvcConfigMap)
	if err != nil {
		logr.WithError(err).Warnf("Unable to load %s: %s", pvcConfigMap, err)
		return ""
	}
	err = yaml.Unmarshal(bytes, &staticVolumes)
	if err != nil {
		logr.WithError(err).Errorf("Can not unmarshel from /etc/static-volumes-v2/PVCs-v2.yaml")
		return ""
	}

	if len(staticVolumes.Volumes) > 0 {
		logr.Debugf("Fetching a static volume")
		n := rand.Int() % len(staticVolumes.Volumes)
		return staticVolumes.Volumes[n].Name
	} else {
		logr.Debugf("len(staticVolumes.Volumes) is zero, so can't allocate a volume!")
	}
	return ""
}

func jobBasePath(trainingID string) string {
	return config.GetEtcdPrefix() + trainingID
}

// Return the etcd base path of learner znodes.
func learnerEtcdBasePath(trainingID string) string {
	return jobBasePath(trainingID) + "/learners"
}

// Return the etcd base path of status of learner znodes.
func learnerNodeEtcdStatusPath(trainingID string, learnerID int) string {
	return fmt.Sprintf("%s/learner_%d/status", learnerEtcdBasePath(trainingID), learnerID)
}

func learnerNodeEtcdStatusPathRelative(trainingID string, learnerID int) string {
	return fmt.Sprintf("%s/learner_%d/status", trainingID, learnerID)
}

// Return the etcd base path of learner znodes.
func learnerNodeEtcdBasePath(trainingID string, learnerID int) string {
	return fmt.Sprintf("%s/learner_%d/", learnerEtcdBasePath(trainingID), learnerID)
}

// calcMemory is a utility to convert the memory from DLaaS resource requirements
// to the default MB notation
func calcMemory(r *service.ResourceRequirements) float64 {
	return calcSize(r.Memory, r.MemoryUnit)
}

// calcStorage is a utility to convert the storage from DLaaS resource requirements
// to the default MB notation
func calcStorage(r *service.ResourceRequirements) float64 {
	return calcSize(r.Storage, r.StorageUnit)
}

// calcSize converts from memory resource requirements to the default MB notation
func calcSize(size float64, unit service.ResourceRequirements_MemoryUnit) float64 {
	// according to google unit converter :)
	switch unit {
	case service.ResourceRequirements_MiB:
		return util.RoundPlus(size*1.048576, 2)
	case service.ResourceRequirements_GB:
		return util.RoundPlus(size*1000, 2)
	case service.ResourceRequirements_TB:
		return util.RoundPlus(size*1000*1000, 2)
	case service.ResourceRequirements_GiB:
		return util.RoundPlus(size*1073.741824, 2)
	case service.ResourceRequirements_TiB:
		return util.RoundPlus(size*1073.741824*1073.741824, 2)
	default:
		return size // assume MB
	}
}

//update job status in the database
//update job status in cassandra
func updateJobStatus(trainingID string, updStatus grpc_storage.Status, userID string, statusMessage string, errorCode string, logr *logger.LocLoggingEntry) error {
	logr.Debugf("(updateJobStatus) Updating status of %s to %s", trainingID, updStatus.String())
	//updateRequest := &grpc_trainer_v2.UpdateRequest{TrainingId: trainingID, Status: updStatus, UserId: userID, StatusMessage: statusMessage, ErrorCode: errorCode}
	updateRequest := &grpc_storage.UpdateRequest{TrainingId: trainingID, Status: updStatus, UserId: userID, StatusMessage: statusMessage, ErrorCode: errorCode}

	//trainer, err := client.NewTrainer()
	storageClient, err := client.NewStorage()

	if err != nil {
		logr.WithError(err).Errorf("(updateJobStatus) Creating training client for status update failed. Training ID %s New Status %s", trainingID, updStatus.String())
		logr.Errorf("(updateJobStatus) Error while creating training client is %s", err.Error())
	}
	defer storageClient.Close()
	err = util.Retry(10, 100*time.Millisecond, "UpdateTrainingJob", logr, func() error {
		//ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		//defer cancel()
		_, err = storageClient.Client().UpdateTrainingJob(context.Background(), updateRequest)
		if err != nil {
			logr.WithError(err).Error("Failed to update status to the trainer. Retrying")
			logr.Infof("WARNING: Status updates for %s may be temporarily inconsistent due to failure to communicate with Trainer.", trainingID)
		}
		return err
	})
	if err != nil {
		logr.WithError(err).Errorf("Failed to update status to the trainer. Already retried several times.")
		logr.Infof("WARNING : Status of job %s will likely be incorrect", trainingID)
		return err
	}

	logr.Debugf("(updateJobStatus) Status update request for %s sent to trainer", trainingID)
	return nil
}

func isJobDone(jobStatus string, logr *logger.LocLoggingEntry) bool {
	statusUpdate := client.GetStatus(jobStatus, logr)
	status := statusUpdate.Status
	return status == grpc_storage.Status_COMPLETED || status == grpc_storage.Status_FAILED
}

// Set the DLaaS service type label to an object.
// This label is used to configure Calico network policy rules for the pod.
func setServiceTypeLabel(spec *metav1.ObjectMeta, value string) {
	spec.Labels["service"] = value
}

func k8sInteractionBackoff() *backoff.ExponentialBackOff {
	back := backoff.NewExponentialBackOff()
	back.MaxElapsedTime = 3 * time.Minute
	back.MaxInterval = 1 * time.Minute
	return back
}

func etdInteractionBackoff(maxElapsedTime, maxInterval time.Duration) *backoff.ExponentialBackOff {
	back := backoff.NewExponentialBackOff()
	back.MaxElapsedTime = maxElapsedTime
	back.MaxInterval = maxInterval
	return back
}

//onError function on how to deal with the scenario if connecting to coordinator failed. the error is still returned in case
func coordinator(logr *logger.LocLoggingEntry) (coord.Coordinator, error) {

	var instance coord.Coordinator
	var err error
	err = backoff.
		RetryNotify(func() error {
			instance, err = coord.NewCoordinator(coord.Config{Endpoints: config.GetEtcdEndpoints(), Prefix: config.GetEtcdPrefix(),
				Cert: config.GetEtcdCertLocation(), Username: config.GetEtcdUsername(), Password: config.GetEtcdPassword()}, logr)
			return err
		}, etdInteractionBackoff(1*time.Minute, 30*time.Second), func(err error, t time.Duration) {
			logr.WithError(err).Errorf("failed to establish connection with etcd")
		})

	return instance, err
}
