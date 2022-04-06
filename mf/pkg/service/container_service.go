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

package service

import (
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/models"
	"mlss-mf/pkg/restapi/operations/container"
	"sort"
)

func ListContainer(params container.ListContainerParams, user string, isSA bool) (models.Pods, int, error) {
	ns := params.Namespace
	name := params.ServiceName

	//Permission Check OK
	if !isSA {
		service, err := serviceDao.GetServiceByNSAndName(params.Namespace, params.ServiceName)
		if err != nil {
			mylogger.Errorf("ListContainer get service error(Check Permission Error), namespace:"+
				" %s, name: %s, error: %v\n",
				ns, name, err)
			return nil, StatusError, err
		}
		if service.User.Name != user {
			return nil, StatusForbidden, PermissionDeniedError
		}
	}

	seldonDeployment, err := seldonClient.MachinelearningV1().SeldonDeployments(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		mylogger.Errorf("ListContainer: fail to get seldonDeployment, namespace: %s, name: %s, error: %v\n",
			ns, name, err)
		return nil, StatusError, err
	}
	pods := models.Pods{}
	for deploymentName, _ := range seldonDeployment.Status.DeploymentStatus {
		mylogger.Debugf("ListContainer deployment name: %s\n", deploymentName)
		deployment, err := common.Clientset.AppsV1().Deployments(ns).Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			mylogger.Errorf("ListContainer: fail to get deployment, namespace: %s, name: %s, error: %v\n",
				ns, deploymentName, err)
			return nil, 500, err
		}
		mylogger.Debugf("ListContainer deployment name: %s, selector: %+v\n",
			deployment, deployment.Spec.Selector.MatchLabels)

		//todo
		//label := deployment.Spec.Selector.String()
		label := ""
		for key, value := range deployment.Spec.Selector.MatchLabels {
			label += fmt.Sprintf("%s=%s,", key, value)
		}
		if len(label) > 0 {
			label = string([]byte(label)[:len(label)-1])
		}
		mylogger.Debugf("ListContainer rs label: %s", label)
		rsList, err := common.Clientset.AppsV1().ReplicaSets(ns).List(metav1.ListOptions{LabelSelector: label})
		if err != nil {
			mylogger.Errorf("ListContainer: fail to get replicaSet list, namespace: %s"+
				", LabelSelector: %s, error: %v\n", ns, label, err)
			return nil, 500, err
		}

		mylogger.Debugf("ListContainer rs length: %d\n", len(rsList.Items))
		if rsList == nil || rsList.Items == nil || len(rsList.Items) == 0 {
			continue
		}
		sort.Sort(ReplicaSetItem(rsList.Items))

		//debug
		//for _, rs := range rsList.Items {
		//	mylogger.Debugf("afeter sort replicaset by createtime, rs name: %s, createtime: %s\n", rs.Name, rs.CreationTimestamp.String())
		//}

		rsName := rsList.Items[0].Name
		rs, err := common.Clientset.AppsV1().ReplicaSets(ns).Get(rsName, metav1.GetOptions{})
		if err != nil {
			mylogger.Errorf("ListContainer: fail to get replicaSet, namespace: %s"+
				", name: %s, error: %v\n", ns, rsName, err)
			return nil, 500, err
		}

		podTemplateHash, ok := rs.Labels["pod-template-hash"]
		if !ok {
			mylogger.Errorf("ListContainer: fail to get pod-template-hash, labels: %+v\n", rs.Labels)
			return nil, 500, errors.New("fail to get pod-template-hash")
		}
		//get pod info by labels
		podList, err := common.Clientset.CoreV1().Pods(ns).List(metav1.ListOptions{
			LabelSelector: fmt.Sprintf("pod-template-hash=%s", podTemplateHash),
		})
		if err != nil {
			mylogger.Errorf("ListContainer: fail to get pod list by namespace %q and label pod-template-hash %q\n",
				ns, podTemplateHash)
			return nil, 500, err
		}

		for _, pod := range podList.Items {
			mylogger.Debugf("ListContainer: k8s pod.status : %+v\n", pod.Status)
			podResult := models.Pod{
				Containers: nil,
				Message:    pod.Status.Message,
				Name:       pod.Name,
				Namespace:  pod.Namespace,
				Reason:     pod.Status.Reason,
				Status:     string(pod.Status.Phase),
			}
			containers := models.Containers{}
			for idx, containerInK8s := range pod.Spec.Containers {
				var gpu int64 = 0
				if gpuQuantity, ok := containerInK8s.Resources.Requests["nvidia.com/gpu"]; ok {
					gpu, _ = gpuQuantity.AsInt64()
				}
				mylogger.Debugf("ListContainer: k8s container info: %+v\n", containerInK8s)
				containerStatus := ""
				if pod.Status.ContainerStatuses[idx].State.Waiting != nil {
					containerStatus = "waiting"
				} else if pod.Status.ContainerStatuses[idx].State.Running != nil {
					containerStatus = "running"
				} else if pod.Status.ContainerStatuses[idx].State.Terminated != nil {
					containerStatus = "treminated"
				}
				container := models.Container{
					ContainerName: containerInK8s.Name,
					CPU:           containerInK8s.Resources.Requests.Cpu().Value(),
					FinishedTime:  "",
					Gpu:           gpu,
					Image:         containerInK8s.Image,
					ImageID:       pod.Status.ContainerStatuses[idx].ImageID,
					Memory:        containerInK8s.Resources.Requests.Memory().String(),
					Namespace:     pod.Namespace,
					PodName:       pod.Name,
					StartedTime:   pod.CreationTimestamp.String(),
					Status:        containerStatus,
				}
				containers = append(containers, &container)
			}
			podResult.Containers = containers
			pods = append(pods, &podResult)
		}
	}

	return pods, 200, nil
}

type ReplicaSetItem []appsv1.ReplicaSet

func (r ReplicaSetItem) Len() int {
	return len(r)

}
func (r ReplicaSetItem) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r ReplicaSetItem) Less(i, j int) bool {
	return r[i].CreationTimestamp.After(r[j].CreationTimestamp.Time)
}
