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
/**
  Created by Luck
  2019/6/4 15:57
*/
package models

import (
	"github.com/kubernetes-client/go/kubernetes/client"
	"k8s.io/api/apps/v1"
)

type K8sNotebook struct {
	ApiVersion string              `json:"apiVersion,omitempty"`
	Kind       string              `json:"kind,omitempty"`
	Metadata   client.V1ObjectMeta `json:"metadata,omitempty"`
	Spec       NBSpec              `json:"spec,omitempty"`
	Status     NBStatus            `json:"status,omitempty"`
	Service    client.V1Service
}
type NBSpec struct {
	Template NBTemplate `json:"template,omitempty"`
}
type NBTemplate struct {
	Spec NBDetailedSpec `json:"spec,omitempty"`
}
type NBDetailedSpec struct {
	ServiceAccountName      string               `json:"serviceAccountName,omitempty"`
	Containers              []client.V1Container `json:"containers,omitempty"`
	TtlSecondsAfterFinished int                  `json:"ttlSecondsAfterFinished,omitempty"`
	Volumes                 []client.V1Volume    `json:"volumes,omitempty"`
	SecurityContext         client.V1SecurityContext `json:"securityContext,omitempty"`
}

type NBStatus struct {
	Conditions     []v1.StatefulSetCondition `json:"conditions,omitempty"`
	ContainerState client.V1ContainerState   `json:"containerState,omitempty"`
	ReadyReplicas  int32                     `json:"readyReplicas,omitempty"`
}