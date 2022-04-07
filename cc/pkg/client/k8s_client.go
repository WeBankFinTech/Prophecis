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

package client

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"mlss-controlcenter-go/pkg/logger"
	"sync"
)

type K8sClient struct {
	ClientSet *kubernetes.Clientset
}

var k8sClient *K8sClient
var restConfig *rest.Config
var k8sClientOnce sync.Once

func GetK8sClient() *K8sClient {
	k8sClientOnce.Do(func() {
		config, err := rest.InClusterConfig()
		if err != nil {
			logger.Logger().Error("Init K8s client err, ", err)
			panic(err)
		}
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			logger.Logger().Error("K8s NewForConfig err, ", err)
			panic(err)
		}
		/*flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", "D://kube/config")
		if err != nil {
			panic(err.Error())
		}
		restConfig = config
		client, err := kubernetes.NewForConfig(config)
		*/
		k8sClient = &K8sClient{
			ClientSet: client,
		}
	})
	return k8sClient
}

func (k8sClient K8sClient) ListNamespacedResourceQuota(ns string) ([]v1.ResourceQuota, error) {
	list, err := k8sClient.ClientSet.CoreV1().ResourceQuotas(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Logger().Error("Get resource quotas err, ", err)
		return nil, err
	}
	logger.Logger().Infof("ResourceQuotas,ns: %v, list: %v", ns, list)
	return list.Items, nil
}

func (k8sClient K8sClient) ListNamespaceNodes(ns string) ([]v1.Node, error) {
	v1Namespace, err := k8sClient.ClientSet.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Error("Get namespace err, ", err)
		return nil, err
	}

	namespaceAnnotations := v1Namespace.Annotations
	logger.Logger().Infof("listNamespaceNodes for namespaceAnnotations: %v", namespaceAnnotations)

	selector := namespaceAnnotations["scheduler.alpha.kubernetes.io/node-selector"]
	list, err := k8sClient.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		logger.Logger().Error("Get nodes err, ", err)
		return nil, err
	}
	logger.Logger().Infof("listNamespaceNodes by namespace: %v, list.items.length: %v", ns, len(list.Items))
	return list.Items, nil
}
