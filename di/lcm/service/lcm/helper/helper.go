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

package helper

import (
	"webank/DI/commons/config"
	"webank/DI/commons/constants"

	"github.com/spf13/viper"
	v1beta1 "k8s.io/api/apps/v1beta1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const FLUENTBIT = "fluent-bit"

//CreatePodSpec ...
// FIXME MLSS Change: add nodeSelectors to deployment/sts pods
//func CreatePodSpec(containers []v1core.Container, volumes []v1core.Volume, labels map[string]string) v1core.PodTemplateSpec {
func CreatePodSpec(containers []v1core.Container, volumes []v1core.Volume, labels map[string]string, nodeSelectorMap map[string]string) v1core.PodTemplateSpec {
	labels["service"] = "dlaas-lhelper" //controls ingress/egress
	imagePullSecret := viper.GetString(config.LearnerImagePullSecretKey)
	automountSeviceToken := false

	// FIXME MLSS Change: define timezone volume
	volumes = append(volumes, constants.TimezoneVolume)

	// FIXME MLSS Temporary Change: use fluent-bit
	logCollectorImageShortName := labels["log_collector_image_short_name"]
	if logCollectorImageShortName == FLUENTBIT {
		volumes = append(volumes, constants.FluentbitConfigVolume)
	}

	return v1core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1core.PodSpec{
			Containers: containers,
			Volumes:    volumes,
			ImagePullSecrets: []v1core.LocalObjectReference{
				v1core.LocalObjectReference{
					Name: imagePullSecret,
				},
			},
			AutomountServiceAccountToken: &automountSeviceToken,
			// FIXME MLSS Change: add nodeSelectors
			NodeSelector: nodeSelectorMap,
		},
	}
}

//CreateDeploymentForHelper ...
func CreateDeploymentForHelper(name string, podTemplateSpec v1core.PodTemplateSpec) *v1beta1.Deployment {

	revisionHistoryLimit := int32(0) //https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#clean-up-policy

	//TODO consider this as well https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#progress-deadline-seconds
	//but not sure if we can nicely revert back
	return &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1beta1.DeploymentSpec{
			Template:             podTemplateSpec,
			RevisionHistoryLimit: &revisionHistoryLimit, //we never rollback these
		},
	}
}
