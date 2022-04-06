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

package lcm_training

// "strconv"

type sigleJobTraining struct {
	*training
}

func (t sigleJobTraining) Start() error {
	// serviceSpec := learner.CreateServiceSpec(t.learner.name, t.req.TrainingId, t.k8sClient)
	// jobUser := t.req.UserId
	// if t.req.ProxyUser != "" {
	// 	jobUser = t.req.ProxyUser
	// }
	// t.logr.Debugf("Service defined: %v ", serviceSpec.Name)
	// if true {
	// 	// serviceSpec = learner.CreateServiceSpec(t.learner.name, t.req.TrainingId, t.k8sClient)
	// 	t.learner.envVars = append(t.learner.envVars, v1core.EnvVar{
	// 		Name: "Driver_Host",
	// 		ValueFrom: &v1core.EnvVarSource{
	// 			FieldRef: &v1core.ObjectFieldSelector{
	// 				FieldPath: "spec.nodeName",
	// 			},
	// 		},
	// 	})

	// 	t.learner.envVars = append(t.learner.envVars, v1core.EnvVar{
	// 		Name:  "Spark_Driver_BlockManager_Port",
	// 		Value: strconv.Itoa(int(serviceSpec.Spec.Ports[0].NodePort)),
	// 	})

	// 	t.learner.envVars = append(t.learner.envVars, v1core.EnvVar{
	// 		Name:  "Spark_Driver_Port",
	// 		Value: strconv.Itoa(int(serviceSpec.Spec.Ports[1].NodePort)),
	// 	})

	// 	t.learner.envVars = append(t.learner.envVars, v1core.EnvVar{
	// 		Name:  "JOB_USER",
	// 		Value: jobUser,
	// 	})

	// }

	// numLearners := int(t.req.GetResources().Learners)

	// job, err := t.jobSpecForLearner(serviceSpec)
	// if err != nil {
	// 	t.logr.WithError(err).Errorf("Could not create job spec for %s", serviceSpec.Name)
	// 	return err
	// }

	// t.logr.Debugf("Split training job deploy namespace: %v ", t.req.JobNamespace)

	// if t.req.JobType != "dist-tf" {
	// 	return t.NewCreateFromBOM(&splitTrainingBOM{
	// 		t.learner.secrets,
	// 		serviceSpec,
	// 		numLearners,
	// 		t.req.JobNamespace,
	// 		job,
	// 	})
	// } else {
	// 	return t.CreateFromBOMForTFJob(&splitTrainingBOM{
	// 		t.learner.secrets,
	// 		serviceSpec,
	// 		numLearners,
	// 		t.req.JobNamespace,
	// 		job,
	// 	})
	// }
	return nil
}

// func (t *sigleJobTraining) NewCreateFromBOM(bom *splitTrainingBOM) error {
// 	logr := t.logr

// 	//namespace := config.GetLearnerNamespace()
// 	namespace := bom.namespace
// 	logr.Infof("namespace in bom is %s", namespace)

// 	for _, secret := range bom.secrets {
// 		//create the secrets
// 		if err := backoff.RetryNotify(func() error {
// 			_, err := t.k8sClient.CoreV1().Secrets(namespace).Create(secret)
// 			if k8serrors.IsAlreadyExists(err) {
// 				logr.WithError(err).Warnf("secret %s already exists", secret.Name)
// 				return nil
// 			}
// 			return err
// 		}, utils.k8sInteractionBackoff(), func(err error, window time.Duration) {
// 			logr.WithError(err).Errorf("Failed in creating secret %s while deploying for training ", secret.Name)
// 			// k8sFailureCounter.With(component, "secret").Add(1)
// 		}); err != nil {
// 			return err
// 		}
// 		logr.Infof("Created secret %s", secret.Name)
// 	}

// 	if bom.numLearners > 1 {
// 		return errors.New("numLearners must be 1")
// 	}

// 	//create the stateful set
// 	return backoff.RetryNotify(func() error {
// 		_, err := t.k8sClient.CoreV1().Services(namespace).Create(bom.service)
// 		logr.WithError(err).Warnf("Service Create err %s :", bom.service.Name)

// 		_, err = t.k8sClient.BatchV1().Jobs(namespace).Create(bom.jobBOM)

// 		if k8serrors.IsAlreadyExists(err) {
// 			logr.WithError(err).Warnf("Job %s already exists", bom.jobBOM.Name)
// 			return nil
// 		}
// 		return err
// 	}, utils.(), func(err error, window time.Duration) {
// 		logr.WithError(err).Errorf("failed in creating Job %s while deploying for training ", bom.jobBOM.Name)
// 		// k8sFailureCounter.With(component, "learner").Add(1)

// 		logr.Debugf("NewCreateFromBOM jobBOM, Volumes: %+v, VolumeMounts: %+v", bom.jobBOM.Spec.Template.Spec.Volumes, bom.jobBOM.Spec.Template.Spec.Containers[0].VolumeMounts)
// 	})

// }

// //CreateFromBOM ... eventually use with controller and make this transactional
// func (t *sigleJobTraining) CreateFromBOMForTFJob(bom *splitTrainingBOM) error {
// 	logr := t.logr
// 	logr.Infof("splitTraining.go, CreateFromBOMForTFJob")

// 	//namespace := config.GetLearnerNamespace()
// 	namespace := bom.namespace
// 	logr.Infof("namespace in bom is %s", namespace)

// 	for _, secret := range bom.secrets {
// 		//create the secrets
// 		if err := backoff.RetryNotify(func() error {
// 			_, err := t.k8sClient.CoreV1().Secrets(namespace).Create(secret)
// 			if k8serrors.IsAlreadyExists(err) {
// 				logr.WithError(err).Warnf("secret %s already exists", secret.Name)
// 				return nil
// 			}
// 			return err
// 		}, utils.k8sInteractionBackoff(), func(err error, window time.Duration) {
// 			logr.WithError(err).Errorf("Failed in creating secret %s while deploying for training ", secret.Name)
// 			// k8sFailureCounter.With(component, "secret").Add(1)
// 		}); err != nil {
// 			return err
// 		}
// 		logr.Infof("Created secret %s", secret.Name)
// 	}

// 	jobBOM := bom.jobBOM

// 	logr.Debugf("now starting to deploy learners for tfjob")
// 	logr.Debugf("CreateFromBOMForTFJob, req, TrainingId: %s, Framework: %s, JobType: %s, Learners: %s, Cpus: %s, Memory: %s, Gpus: %s, GpuType: %s, ", t.req.TrainingId, t.req.Framework, t.req.JobType, t.req.Resources.Learners, t.req.Resources.Cpus, t.req.Resources.Memory, t.req.Resources.Gpus, t.req.Resources.GpuType)
// 	err := tfjob.NewSubmitTFJob(t.req.TrainingId, t.req.JobNamespace, "wedatasphere/prophecis:tf-mnist-with-summaries-v1.0.0", t.req, jobBOM, logr)
// 	if err != nil {
// 		return err
// 	}
// 	return nil

// }
