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
	"strings"
	"webank/DI/commons/constants"
	"webank/DI/commons/tfjob"

	"github.com/cenkalti/backoff"
	v1 "k8s.io/api/batch/v1"

	//"webank/DI/commons/metricsmon"
	"time"
	"webank/DI/lcm/service/lcm/learner"

	v1core "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (t splitTraining) Start() error {

	serviceSpec := learner.CreateServiceSpec(t.learner.name, t.req.TrainingId)

	numLearners := int(t.req.GetResources().Learners)

	job, err := t.jobSpecForLearner(serviceSpec.Name)
	if err != nil {
		t.logr.WithError(err).Errorf("Could not create job spec for %s", serviceSpec.Name)
		return err
	}

	t.logr.Debugf("Split training job deploy namespace: %v ", t.req.JobNamespace)

	if t.req.JobType != "dist-tf" {
		return t.NewCreateFromBOM(&splitTrainingBOM{
			t.learner.secrets,
			serviceSpec,
			numLearners,
			t.req.JobNamespace,
			job,
		})
	} else {
		return t.CreateFromBOMForTFJob(&splitTrainingBOM{
			t.learner.secrets,
			serviceSpec,
			numLearners,
			t.req.JobNamespace,
			job,
		})
	}
}

func (t splitTraining) jobSpecForLearner(serviceName string) (*v1.Job, error) {

	gpus := make(map[string]string)
	if t.req.Resources.Gpus > 0 {
		gpus["ibm-cloud.kubernetes.io/gpu-type"] = t.req.Resources.GpuType
	}

	learnerDefn := t.learner
	helperDefn := t.helper

	helperAndLearnerVolumes := append(learnerDefn.volumes, helperDefn.sharedVolume)

	// FIXME MLSS Change: define timezone volume
	helperAndLearnerVolumes = append(helperAndLearnerVolumes, constants.TimezoneVolume)
	// FIXME MLSS Change: add joblogs Volume for fluent-bit
	helperAndLearnerVolumes = append(helperAndLearnerVolumes, constants.JobLogsVolume)

	imagePullSecret, err := learner.GenerateImagePullSecret(t.k8sClient, t.req)
	if err != nil {
		return nil, err
	}

	// FIXME MLSS Change: get namespace annotations
	namespaceObj, err := t.k8sClient.CoreV1().Namespaces().Get(t.req.JobNamespace, metav1.GetOptions{})
	var nodeSelectors string = ""
	if nil == err {
		namespaceAnnotations := (namespaceObj.ObjectMeta).GetAnnotations()
		selector, ok := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]
		if ok {
			nodeSelectors = selector
		}
	}
	//FIXME MLSS Change: parse nodeSelectors and add to pods
	nodeSelectorMap := make(map[string]string)
	if len(nodeSelectors) > 0 {
		nodeSelectorList := strings.Split(nodeSelectors, ",")
		for i := 0; i < len(nodeSelectorList); i++ {
			nodeSelector := nodeSelectorList[i]
			nodeSelectorKv := strings.Split(nodeSelector, "=")
			nodeSelectorMap[nodeSelectorKv[0]] = nodeSelectorKv[1]
		}
	}
	// add training-related selectors(gpus) to nodeSelector map

	logr := t.logr
	//now create the learner container
	//logr.Debugf("learnerDefn.envVars: %+v", learnerDefn.envVars)
	learnerContainer := newConstructLearnerContainer(t.req, learnerDefn.envVars, learnerDefn.volumeMounts, helperDefn.sharedVolumeMount, learnerDefn.mountTrainingDataStoreInLearner, learnerDefn.mountResultsStoreInLearner, t.logr) // nil for mounting shared NFS volume since non split mode
	//logr.Debugf("learnerContainer envVars: %+v", learnerContainer.Env)

	// FIXME MLSS Change: add nodeSelectors to deployment/sts pods
	logr.Debugf("debug_imagePullSecret: %v", imagePullSecret)
	splitLearnerPodSpec := learner.CreatePodSpecForJob([]v1core.Container{learnerContainer}, helperAndLearnerVolumes, map[string]string{"training_id": t.req.TrainingId, "user_id": t.req.UserId}, nodeSelectorMap, imagePullSecret)
	forLearner := learner.CreateJobSpecForLearner(t.req.TrainingId, splitLearnerPodSpec)

	return forLearner, nil
}

func (t *splitTraining) NewCreateFromBOM(bom *splitTrainingBOM) error {
	logr := t.logr

	//namespace := config.GetLearnerNamespace()
	namespace := bom.namespace
	logr.Infof("namespace in bom is %s", namespace)

	for _, secret := range bom.secrets {
		//create the secrets
		if err := backoff.RetryNotify(func() error {
			_, err := t.k8sClient.CoreV1().Secrets(namespace).Create(secret)
			if k8serrors.IsAlreadyExists(err) {
				logr.WithError(err).Warnf("secret %s already exists", secret.Name)
				return nil
			}
			return err
		}, k8sInteractionBackoff(), func(err error, window time.Duration) {
			logr.WithError(err).Errorf("Failed in creating secret %s while deploying for training ", secret.Name)
			k8sFailureCounter.With(component, "secret").Add(1)
		}); err != nil {
			return err
		}
		logr.Infof("Created secret %s", secret.Name)
	}

	if bom.numLearners > 1 {
		return errors.New("numLearners must be 1")
	}

	//create the stateful set
	return backoff.RetryNotify(func() error {
		//_, err := t.k8sClient.AppsV1beta1().StatefulSets(namespace).Create(bom.learnerBOM)
		//logr.Debugf("deployJobToK8s, bom.jobBOM.Spec.Template.Spec.Containers[0].Env: %+v", bom.jobBOM.Spec.Template.Spec.Containers[0].Env)

		_, err := t.k8sClient.BatchV1().Jobs(namespace).Create(bom.jobBOM)

		if k8serrors.IsAlreadyExists(err) {
			logr.WithError(err).Warnf("Job %s already exists", bom.jobBOM.Name)
			return nil
		}
		return err
	}, k8sInteractionBackoff(), func(err error, window time.Duration) {
		logr.WithError(err).Errorf("failed in creating Job %s while deploying for training ", bom.jobBOM.Name)
		k8sFailureCounter.With(component, "learner").Add(1)

		logr.Debugf("NewCreateFromBOM jobBOM, Volumes: %+v, VolumeMounts: %+v", bom.jobBOM.Spec.Template.Spec.Volumes, bom.jobBOM.Spec.Template.Spec.Containers[0].VolumeMounts)
	})

}

//CreateFromBOM ... eventually use with controller and make this transactional
func (t *splitTraining) CreateFromBOMForTFJob(bom *splitTrainingBOM) error {
	logr := t.logr
	logr.Infof("splitTraining.go, CreateFromBOMForTFJob")

	//namespace := config.GetLearnerNamespace()
	namespace := bom.namespace
	logr.Infof("namespace in bom is %s", namespace)

	for _, secret := range bom.secrets {
		//create the secrets
		if err := backoff.RetryNotify(func() error {
			_, err := t.k8sClient.CoreV1().Secrets(namespace).Create(secret)
			if k8serrors.IsAlreadyExists(err) {
				logr.WithError(err).Warnf("secret %s already exists", secret.Name)
				return nil
			}
			return err
		}, k8sInteractionBackoff(), func(err error, window time.Duration) {
			logr.WithError(err).Errorf("Failed in creating secret %s while deploying for training ", secret.Name)
			k8sFailureCounter.With(component, "secret").Add(1)
		}); err != nil {
			return err
		}
		logr.Infof("Created secret %s", secret.Name)
	}

	jobBOM := bom.jobBOM

	logr.Debugf("now starting to deploy learners for tfjob")
	logr.Debugf("CreateFromBOMForTFJob, req, TrainingId: %s, Framework: %s, JobType: %s, Learners: %s, Cpus: %s, Memory: %s, Gpus: %s, GpuType: %s, ", t.req.TrainingId, t.req.Framework, t.req.JobType, t.req.Resources.Learners, t.req.Resources.Cpus, t.req.Resources.Memory, t.req.Resources.Gpus, t.req.Resources.GpuType)
	err := tfjob.NewSubmitTFJob(t.req.TrainingId, t.req.JobNamespace, "wedatasphere/prophecis:tf-mnist-with-summaries-v1.0.0", t.req, jobBOM, logr)
	if err != nil {
		return err
	}
	return nil

}
