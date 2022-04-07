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
package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	seldonV1 "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/seldonio/seldon-core/operator/client/machinelearning.seldon.io/v1/clientset/versioned"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/config"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"strconv"
	"strings"
	"sync"
)

func initSeldonClient() (seldonClient *versioned.Clientset) {
	// create the mpijobClientset
	seldonClient = versioned.NewForConfigOrDie(config.GetKubernetesConfig())
	return seldonClient
}

var (
	seldonClient *versioned.Clientset
)

func init() {
	seldonClient = initSeldonClient()
}

type K8sApiClient struct {
}

var once sync.Once
var k8sApiClient *K8sApiClient

func GetK8sApiClient() *K8sApiClient {
	once.Do(func() {
		k8sApiClient = &K8sApiClient{}
	})
	return k8sApiClient
}

func (kac *K8sApiClient) GetNSFromK8sForAdd(namespace string) (*v1.Namespace, error) {
	k8sClient := getClient()

	options := metav1.GetOptions{}
	nsFromK8s, e := k8sClient.CoreV1().Namespaces().Get(context.TODO(), namespace, options)
	if e != nil {
		logger.Logger().Debugf("AddNamespace getNS err: %v", e.Error())
		if strings.Contains(e.Error(), "not found") {
			return nil, nil
		}
		return nil, e
	}

	return nsFromK8s, nil
}

func (kac *K8sApiClient) GetNSFromK8s(namespace string) (*v1.Namespace, error) {
	k8sClient := getClient()

	options := metav1.GetOptions{}
	nsFromK8s, e := k8sClient.CoreV1().Namespaces().Get(context.TODO(), namespace, options)
	if e != nil {
		return nil, e
	}

	return nsFromK8s, nil
}

func (kac *K8sApiClient) DeleteNSFromK8s(namespace string) error {
	k8sClient := getClient()

	return k8sClient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
}

func (kac *K8sApiClient) CreateNamespaceWithResources(namespace string, annotations map[string]string, platformNamespace string) error {
	k8sClient := getClient()

	objectMeta := &metav1.ObjectMeta{
		Name: namespace,
	}
	if "" == platformNamespace {
		platformNamespace = "default"
	}

	var v1Namespace = &v1.Namespace{
		ObjectMeta: *objectMeta,
	}

	_, createErr := k8sClient.CoreV1().Namespaces().Create(context.TODO(), v1Namespace, metav1.CreateOptions{})
	if nil != createErr {
		return createErr
	}

	logger.Logger().Debugf("start creating configMaps & secrets for namespace: %v, platformNamespace: %v", namespace, platformNamespace)
	leaFiErr := kac.CopyCMFromPlatformNamespace("learner-entrypoint-files", namespace, platformNamespace)
	if leaFiErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: learner-entrypoint-files", namespace)
		return leaFiErr
	}

	leaConErr := kac.CopyCMFromPlatformNamespace("learner-config", namespace, platformNamespace)
	if leaConErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: learner-config", namespace)
		return leaConErr
	}

	diConErr := kac.CopyCMFromPlatformNamespace("di-config", namespace, platformNamespace)
	if diConErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: di-config", namespace)
		return diConErr
	}

	noteErr := kac.CopyCMFromPlatformNamespace("notebook-entrypoint-files", namespace, platformNamespace)
	if noteErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: notebook-entrypoint-files", namespace)
		return noteErr
	}

	yarnErr := kac.CopyCMFromPlatformNamespace("yarn-resource-setting", namespace, platformNamespace)
	if yarnErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: yarn-resource-setting", namespace)
		return yarnErr
	}

	//fluBitErr := kac.CopyCMFromPlatformNamespace("fluent-bit-log-collector-config", namespace, platformNamespace)
	//if fluBitErr != nil {
	//	logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: fluent-bit-log-collector-config", namespace)
	//	return fluBitErr
	//}

	//lcnSecErr := kac.CopySecretFromNamespace("lcm-secrets", namespace, platformNamespace)
	//if lcnSecErr != nil {
	//	logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: lcm-secrets", namespace)
	//	return lcnSecErr
	//}

	//hubSecErr := kac.CopySecretFromNamespace("hubsecret-go", namespace, platformNamespace)
	//if hubSecErr != nil {
	//	logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configName: %v", namespace)
	//	return hubSecErr
	//}

	//logger.Logger().Debugf("start creating pvc & its configMaps for  namespace: %v, platformNamespace: %v", namespace, platformNamespace)

	//createPVCErr := kac.CopyPVCFromPlatformNamespace("static-volume-1", namespace, platformNamespace)
	//if createPVCErr != nil {
	//	logger.Logger().Errorf("failed to create pvc to k8s for namespace: %v, pvcName: static-volume-1", namespace)
	//	return createPVCErr
	//}

	//generateCMFromPVCErr := kac.GenerateCMFromPVC("static-volume-1", "PVCs.yaml", "static-volumes", namespace)
	//if generateCMFromPVCErr != nil {
	//	logger.Logger().Errorf("failed to GenerateCMFromPVC for namespace: %v, pvcName: static-volume-1, configKeyName: PVCs.yaml", namespace)
	//	return createPVCErr
	//}

	//generateV2CMFromPVCErr := kac.GenerateV2CMFromPVC("static-volumes-v2", "static-volumes-v2", "PVCs-v2.yaml", namespace)
	//if generateV2CMFromPVCErr != nil {
	//	logger.Logger().Errorf("failed to GenerateV2CMFromPVC for namespace: %v, pvcName: static-volumes-v2, configKeyName: static-volumes-v2", namespace)
	//	return createPVCErr
	//}

	rqConfig := common.GetAppConfig().Core.Kube.NamespacedResourceConfig
	logger.Logger().Debugf("create namespace resourcesQuota with config: %v", rqConfig)

	createRQErr := kac.CreateRQ(rqConfig.DefaultRQName, namespace, rqConfig.DefaultRQCpu, rqConfig.DefaultRQCpu, rqConfig.DefaultRQMem, rqConfig.DefaultRQMem, rqConfig.DefaultRQGpu)
	if createRQErr != nil {
		logger.Logger().Errorf("failed to create resourceQuota for namespace: %v", namespace)
		return createRQErr
	}

	split := strings.Split(namespace, "-")

	//var netWork string

	//if len(split) == 6 {
	//	netWork = split[5]
	//} else {
	//	netWork = split[6]
	//}

	var nodeSelectors bytes.Buffer

	//nodeSelectors.WriteString("lb-idc=")
	//nodeSelectors.WriteString(split[1])
	//nodeSelectors.WriteString(",lb-department=")
	//nodeSelectors.WriteString(split[3])
	//nodeSelectors.WriteString(",lb-app=")
	//nodeSelectors.WriteString(split[4])
	//nodeSelectors.WriteString(",lb-network=")
	//nodeSelectors.WriteString(netWork)

	nodeSelectors.WriteString("lb-department=")
	nodeSelectors.WriteString(split[1])
	nodeSelectors.WriteString(",lb-app=")
	nodeSelectors.WriteString(split[2])

	if annotations != nil {
		customSelectors := annotations["scheduler.alpha.kubernetes.io/node-selector"]
		if "" == customSelectors {
			if "" != annotations["lb-bus-type"] && "" != annotations["lb-gpu-model"] {
				nodeSelectors.WriteString(",lb-bus-type=")
				nodeSelectors.WriteString(annotations["lb-bus-type"])
				nodeSelectors.WriteString(",lb-gpu-model=")
				nodeSelectors.WriteString(annotations["lb-gpu-model"])
			} else if "" != annotations["lb-bus-type"] {
				nodeSelectors.WriteString(",lb-bus-type=")
				nodeSelectors.WriteString(annotations["lb-bus-type"])
			} else if "" != annotations["lb-gpu-model"] {
				nodeSelectors.WriteString(",lb-gpu-model=")
				nodeSelectors.WriteString(annotations["lb-gpu-model"])
			}
		} else {
			customSelectorArr := strings.Split(customSelectors, ",")
			for _, customSelector := range customSelectorArr {
				cusSelKey := strings.Split(customSelector, "=")[0]
				if cusSelKey != "lb-department" && cusSelKey != "lb-app"  {
					nodeSelectors.WriteString(",")
					nodeSelectors.WriteString(customSelector)
				}
			}
			annotations["scheduler.alpha.kubernetes.io/node-selector"] = nodeSelectors.String()
		}
	} else {
		annotations = map[string]string{}
	}

	annotations["scheduler.alpha.kubernetes.io/node-selector"] = nodeSelectors.String()
	annotations["MLSS-PlatformNamespace"] = platformNamespace

	replaceNamespaceErr := kac.ReplaceNamespace(namespace, annotations)
	if nil != replaceNamespaceErr {
		logger.Logger().Errorf("failed to replaceNamespace from k8s for namespace: %v", namespace)
		return replaceNamespaceErr
	}

	createSAErr := kac.CreateSAForNotebook(namespace)
	if nil != createSAErr {
		return createSAErr
	}

	return nil
}

func (kac *K8sApiClient) CopyCMFromPlatformNamespace(name string, namespace string, platformNamespace string) error {
	k8sClient := getClient()

	options := metav1.GetOptions{}
	configMap, configErr := k8sClient.CoreV1().ConfigMaps(platformNamespace).Get(context.TODO(), name, options)
	if configErr != nil {
		return configErr
	}

	var meta = &metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}

	var typeMeta = &metav1.TypeMeta{
		APIVersion: configMap.APIVersion,
		Kind:       configMap.Kind,
	}

	var newConfig = &v1.ConfigMap{
		TypeMeta:   *typeMeta,
		Data:       configMap.Data,
		ObjectMeta: *meta,
	}

	create, createErr := k8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), newConfig, metav1.CreateOptions{})
	if configErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, config: %v", namespace, create)
		return createErr
	}

	logger.Logger().Debugf("create configMap to k8s for namespace: %v, config: %v", namespace, create)

	return nil
}

func (kac *K8sApiClient) CopySecretFromNamespace(name string, namespace string, platformNamespace string) error {
	k8sClient := getClient()

	options := metav1.GetOptions{}

	secret, secretErr := k8sClient.CoreV1().Secrets(platformNamespace).Get(context.TODO(), name, options)

	if secretErr != nil {
		return secretErr
	}

	var meta = &metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}

	var typeMeta = &metav1.TypeMeta{
		APIVersion: secret.TypeMeta.APIVersion,
		Kind:       secret.TypeMeta.Kind,
	}

	var newSecret = &v1.Secret{
		Type:       secret.Type,
		TypeMeta:   *typeMeta,
		Data:       secret.Data,
		ObjectMeta: *meta,
	}

	create, createErr := k8sClient.CoreV1().Secrets(namespace).Create(context.TODO(), newSecret, metav1.CreateOptions{})
	if createErr != nil {
		logger.Logger().Errorf("failed to create secret to k8s for namespace: %v, secret: %v", namespace, create)
		return createErr
	}

	logger.Logger().Debugf("create secret to k8s for namespace: %v, secret: %v", namespace, create)

	return nil

}

func (kac *K8sApiClient) CopyPVCFromPlatformNamespace(name string, namespace string, platformNamespace string) error {
	k8sClient := getClient()

	pvc, perVolumeErr := k8sClient.CoreV1().PersistentVolumeClaims(platformNamespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if perVolumeErr != nil {
		return perVolumeErr
	}

	newAnnotations := map[string]string{}
	newAnnotations["volume.beta.kubernetes.io/storage-class"] = pvc.ObjectMeta.Annotations["volume.beta.kubernetes.io/storage-class"]

	var meta = &metav1.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Annotations: newAnnotations,
		Labels:      pvc.ObjectMeta.Labels,
	}

	var typeMeta = &metav1.TypeMeta{
		APIVersion: pvc.APIVersion,
		Kind:       pvc.Kind,
	}

	var spec = &v1.PersistentVolumeClaimSpec{
		AccessModes: pvc.Spec.AccessModes,
		Resources:   pvc.Spec.Resources,
	}

	var newPVC = &v1.PersistentVolumeClaim{
		TypeMeta:   *typeMeta,
		ObjectMeta: *meta,
		Spec:       *spec,
	}

	claim, createPerVolumeErr := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), newPVC, metav1.CreateOptions{})
	if createPerVolumeErr != nil {
		logger.Logger().Errorf("failed to create pvc to k8s for namespace: %v, pvc: %v", namespace, claim)
		return createPerVolumeErr
	}

	logger.Logger().Debugf("create pvc to k8s for namespace: %v, pvc: %v", namespace, claim)

	return nil
}

func (kac *K8sApiClient) GenerateCMFromPVC(pvcName string, configKeyName string, configMapName string, namespace string) error {
	k8sClient := getClient()

	pvc, pvcErr := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if pvcErr != nil {
		return pvcErr
	}

	marshal, _ := yaml.Marshal(pvc)
	pvcString := string(marshal)
	dataMap := map[string]string{}
	dataMap[configKeyName] = pvcString
	var typeMeta = &metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	}

	var meta = &metav1.ObjectMeta{
		Name:      configMapName,
		Namespace: namespace,
	}

	var newConfig = &v1.ConfigMap{
		TypeMeta:   *typeMeta,
		Data:       dataMap,
		ObjectMeta: *meta,
	}

	configMap, cfgErr := k8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), newConfig, metav1.CreateOptions{})
	if cfgErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configMap: %v", namespace, newConfig)
		return cfgErr
	}

	logger.Logger().Debugf("create configMap to k8s for namespace: %v, configMap: %v", namespace, configMap)

	return nil
}

func (kac *K8sApiClient) GenerateV2CMFromPVC(pvcName string, configKeyName string, configMapName string, namespace string) error {
	k8sClient := getClient()

	objectMap := map[string]string{}
	objectMap["name"] = pvcName
	objectMap["zlabel"] = pvcName
	objectMap["status"] = "active"

	internalDataMap := map[string]interface{}{}
	objectMapList := []map[string]string{objectMap}
	internalDataMap[configKeyName] = objectMapList
	marshal, _ := yaml.Marshal(internalDataMap)
	pvcString := string(marshal)
	dataMap := map[string]string{}
	dataMap[configMapName] = pvcString

	var typeMeta = &metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ConfigMap",
	}

	var meta = &metav1.ObjectMeta{
		Name:      pvcName,
		Namespace: namespace,
	}

	var newConfig = &v1.ConfigMap{
		TypeMeta:   *typeMeta,
		Data:       dataMap,
		ObjectMeta: *meta,
	}

	configMap, cfgErr := k8sClient.CoreV1().ConfigMaps(namespace).Create(context.TODO(), newConfig, metav1.CreateOptions{})
	if cfgErr != nil {
		logger.Logger().Errorf("failed to create configMap to k8s for namespace: %v, configMap: %v", namespace, newConfig)
		return cfgErr
	}

	logger.Logger().Debugf("create configMap to k8s for namespace: %v, configMap: %v", namespace, configMap)

	return nil
}

func (kac *K8sApiClient) CreateRQ(name string, namespace string, cpuRequests string, cpuLimits string, memRequests string, memLimits string, gpu string) error {
	k8sClient := getClient()

	var reqList = v1.ResourceList{}
	hardRestraint := v1.ResourceList{}
	parCPUReq, parCPUReqErr := resource.ParseQuantity(cpuRequests)
	if nil != parCPUReqErr {
		return parCPUReqErr
	}

	parCPULit, parCPULitErr := resource.ParseQuantity(cpuLimits)
	if nil != parCPULitErr {
		return parCPULitErr
	}

	parGPU, parGPUErr := resource.ParseQuantity(gpu)
	if nil != parGPUErr {
		return parGPUErr
	}

	parMemReq, parMemReqErr := resource.ParseQuantity(memRequests)
	if nil != parMemReqErr {
		return parMemReqErr
	}

	parMemLit, parMemLitErr := resource.ParseQuantity(memLimits)
	if nil != parMemLitErr {
		return parMemLitErr
	}

	hardRestraint["requests.cpu"] = parCPUReq
	hardRestraint["limits.cpu"] = parCPULit
	hardRestraint["requests.memory"] = parMemReq
	hardRestraint["limits.memory"] = parMemLit
	hardRestraint["requests.nvidia.com/gpu"] = parGPU

	var typeMeta = &metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ResourceQuota",
	}

	var meta = &metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}

	spec := &v1.ResourceQuotaSpec{
		Hard: reqList,
	}

	rq := &v1.ResourceQuota{
		TypeMeta:   *typeMeta,
		ObjectMeta: *meta,
		Spec:       *spec,
	}

	quota, quotaErr := k8sClient.CoreV1().ResourceQuotas(namespace).Create(context.TODO(), rq, metav1.CreateOptions{})

	if quotaErr != nil {
		logger.Logger().Errorf("failed to create quota to k8s for namespace: %v, quota: %v", namespace, rq)
		return quotaErr
	}

	logger.Logger().Debugf("create quota to k8s for namespace: %v, quota: %v", namespace, quota)

	return nil
}

func (kac *K8sApiClient) ReplaceNamespace(namespace string, annotations map[string]string) error {
	k8sClient := getClient()

	v1Namespace, getNSErr := k8sClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})

	if getNSErr != nil {
		logger.Logger().Errorf("failed to get ns from k8s for namespace: %v", namespace)
		return getNSErr
	}

	nsAnnotations := v1Namespace.ObjectMeta.Annotations

	if nil != nsAnnotations {
		for k, v := range annotations {
			nsAnnotations[k] = v
		}
	} else {
		nsAnnotations = annotations
	}

	v1Namespace.ObjectMeta.Annotations = nsAnnotations

	updateNS, updateNSErr := k8sClient.CoreV1().Namespaces().Update(context.TODO(), v1Namespace, metav1.UpdateOptions{})

	if nil != updateNSErr {
		logger.Logger().Errorf("failed to update namespace from k8s for namespace: %v with annotations: %v", namespace, nsAnnotations)
		return updateNSErr
	}

	logger.Logger().Debugf("update namespace to k8s for namespace: %v", updateNS)

	return nil
}

func (kac *K8sApiClient) CreateSAForNotebook(namespace string) error {
	k8sClient := getClient()

	var typeMeta = &metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ServiceAccount",
	}

	labelMap := map[string]string{"ksonnet.io/component": "jupyter-web-app"}

	var meta = &metav1.ObjectMeta{
		Name:      "jupyter-notebook",
		Namespace: namespace,
		SelfLink:  "/api/v1/namespaces/kubeflow/serviceaccounts/jupyter-notebook",
		Labels:    labelMap,
	}

	serviceAccount := &v1.ServiceAccount{
		TypeMeta:   *typeMeta,
		ObjectMeta: *meta,
	}

	account, SAErr := k8sClient.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if nil != SAErr {
		logger.Logger().Errorf("failed to create serviceAccount from k8s for namespace: %v with serviceAccount: %v", namespace, serviceAccount)
		return SAErr
	}

	logger.Logger().Debugf("update namespace to k8s for namespace: %v, serviceAccount: %v", namespace, account)

	return nil
}

func (kac *K8sApiClient) GetAllNamespacesFromK8s() (*v1.NamespaceList, error) {
	k8sClient := getClient()
	return k8sClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (kac *K8sApiClient) ListRqOfNamespace(namespace string) (*v1.ResourceQuotaList, error) {
	k8sClient := getClient()

	return k8sClient.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (kac *K8sApiClient) GetNodesByNamespace(namespace string) ([]v1.Node, error) {
	nsFromK8s, nsErr := kac.GetNSFromK8s(namespace)

	if nil != nsErr {
		return nil, nsErr
	}

	nodeList, nodeErr := kac.GetAllNodes()
	if nil != nodeErr {
		return nil, nodeErr
	}
	logger.Logger().Debugf("GetNodesByNamespace with namespace: %v nodeNum: %v", nsFromK8s.Namespace, len(nodeList.Items))

	annotations := nsFromK8s.ObjectMeta.Annotations
	var nodeSelector string
	if nil != annotations {
		nodeSelector = annotations[constants.NodeSelectorKey]
	}

	var nodes []v1.Node
	var flag = false
	if "" != nodeSelector {
		nsLabels := strings.Split(nodeSelector, ",")
		for _, v1Node := range nodeList.Items {
			noLabels := v1Node.ObjectMeta.Labels
			for _, nsLabel := range nsLabels {
				if len(strings.Split(nsLabel, "=")) > 1 {
					if _, ok := noLabels[strings.Split(nsLabel, "=")[0]]; ok {
						flag = noLabels[strings.Split(nsLabel, "=")[0]] == strings.Split(nsLabel, "=")[1]
					}
				} else if len(strings.Split(nsLabel, "=")) == 1 {
					if _, ok := noLabels[strings.Split(nsLabel, "=")[0]]; ok {
						flag = "" == noLabels[strings.Split(nsLabel, "=")[0]] && "" == strings.Split(nsLabel, "=")[1]
					}
				}
				if !flag {
					break
				}
			}
			if flag {
				nodes = append(nodes, v1Node)
			}
		}
	}

	return nodes, nil
}

func (kac *K8sApiClient) GetAllNodes() (*v1.NodeList, error) {
	k8sClient := getClient()
	return k8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (kac *K8sApiClient) CheckForRQRequest(nodes []v1.Node, cpu string, gpu string, memory string) error {
	if nil != nodes && len(nodes) > 0 {
		var countGPU float64 = 0
		var countCPU float64 = 0
		var countMemory float64 = 0
		for _, node := range nodes {
			var hardwareType string
			labels := node.ObjectMeta.Labels
			if nil != labels {
				hardwareType = labels["hardware-type"]
			}
			availableResources := node.Status.Allocatable
			cpuQty := availableResources[v1.ResourceCPU]
			nodeCpu, _ := strconv.ParseFloat(cpuQty.AsDec().String(), 64)
			memQty := availableResources[v1.ResourceMemory]
			nodeMem, _ := strconv.ParseFloat(memQty.AsDec().String(), 64)

			var gpuQty resource.Quantity
			var nodeGpu float64
			if constants.NVIDIAGPU == hardwareType {
				gpuQty = availableResources["nvidia.com/gpu"]
				nodeParseFloatGpu, _ := strconv.ParseFloat(gpuQty.AsDec().String(), 64)
				nodeGpu = nodeParseFloatGpu
			}
			countCPU += nodeCpu
			countMemory += nodeMem
			countGPU += nodeGpu
		}
		logger.Logger().Debugf("for resources with reqCpu: %v, avaCpu: %v, reqGpu: %v, vavGpu: %v, reqMem: %v, avaMem: %v, nodes: %v", cpu, countCPU, gpu, countGPU, memory, countMemory, len(nodes))

		flCPU, cpuErr := common.StringToFloat64(cpu)
		if nil != cpuErr {
			return cpuErr
		}

		flGPU, gpuErr := common.StringToFloat64(gpu)
		if nil != gpuErr {
			return gpuErr
		}

		flMem, memErr := common.StringToFloat64(memory)
		if nil != memErr {
			return memErr
		}

		logger.Logger().Debugf("request cpu: %v, gpu: %v, mem: %v", flCPU, flGPU, flMem)
		logger.Logger().Debugf("node countCPU: %v, countGPU: %v, mem: %v", countCPU, countGPU, countMemory)

		if flCPU > countCPU || flGPU > countGPU || flMem*1024*1024 > countMemory {
			errMsg := fmt.Sprintf("the resources set exceed the maximum limit for node cpu: %v, gpu: %v, mem: %v Mi", countCPU, countGPU, countMemory/1024/1024)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("namespace is not yet tied to any node")
	}

	return nil
}

func (kac *K8sApiClient) UpdateRQ(namespace string, cpu string, gpu string, memory string) error {
	k8sClient := getClient()

	hardRestraint := v1.ResourceList{}
	parCPU, parCPUErr := resource.ParseQuantity(cpu)
	if nil != parCPUErr {
		return parCPUErr
	}

	parGPU, parGPUErr := resource.ParseQuantity(gpu)
	if nil != parGPUErr {
		return parGPUErr
	}

	parMem, parMemErr := resource.ParseQuantity(memory)
	if nil != parMemErr {
		return parMemErr
	}

	hardRestraint["requests.cpu"] = parCPU
	hardRestraint["limits.cpu"] = parCPU
	hardRestraint["requests.memory"] = parMem
	hardRestraint["limits.memory"] = parMem
	hardRestraint["requests.nvidia.com/gpu"] = parGPU

	var typeMeta = &metav1.TypeMeta{
		APIVersion: "v1",
		Kind:       "ResourceQuota",
	}

	var meta = &metav1.ObjectMeta{
		Name:      "mlss-default-rq",
		Namespace: namespace,
	}

	var spec = &v1.ResourceQuotaSpec{
		Hard: hardRestraint,
	}

	var quota = &v1.ResourceQuota{
		TypeMeta:   *typeMeta,
		ObjectMeta: *meta,
		Spec:       *spec,
	}

	resourceQuota, updateRQErr := k8sClient.CoreV1().ResourceQuotas(namespace).
		Update(context.TODO(), quota, metav1.UpdateOptions{})

	if updateRQErr != nil {
		logger.Logger().Errorf("failed to update resourcesQuota for namespace: %v", namespace)
		return updateRQErr
	}

	logger.Logger().Errorf("update resourcesQuota for namespace: %v with resourcesQuota: %v", namespace, resourceQuota)

	return nil
}

func (kac *K8sApiClient) GetNodeByName(nodeName string) (*v1.Node, error) {
	k8sClient := getClient()
	return k8sClient.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
}

func (kac *K8sApiClient) ReplaceNode(nodeName string, labels map[string]string) error {
	k8sClient := getClient()

	node, getErr := kac.GetNodeByName(nodeName)
	if nil != getErr {
		return getErr
	}

	noLabels := node.ObjectMeta.Labels
	for k, v := range labels {
		noLabels[k] = v
	}
	node.ObjectMeta.Labels = noLabels

	updateNode, updateErr := k8sClient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if nil != updateErr {
		return updateErr
	}

	logger.Logger().Debugf("update node labels for node: %v", updateNode.Name)

	return nil
}

func getClient() *kubernetes.Clientset {
	k8sClient, _ := kubernetes.NewForConfig(config.GetKubernetesConfig())

	return k8sClient
}

func (kac *K8sApiClient) RemoveNodeLabel(nodeName string, labels []string) (map[string]string, error) {
	k8sClient := getClient()

	node, getErr := kac.GetNodeByName(nodeName)

	if nil != getErr {
		return nil, getErr
	}

	nodeLabels := node.ObjectMeta.Labels
	for k, _ := range nodeLabels {
		for _, nl := range labels {
			if k == nl {
				delete(nodeLabels, k)
			}
		}
	}

	node.ObjectMeta.Labels = nodeLabels

	_, updateErr := k8sClient.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if nil != updateErr {
		return nil, updateErr
	}

	return nodeLabels, nil
}

func (kac *K8sApiClient) GetPodsByNS(namespace string) (*v1.PodList, error) {
	k8sClient := getClient()
	return k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (kac *K8sApiClient) GetPodLabelByNSAndName(namespace string, name string) (map[string]string, error) {
	k8sClient := getClient()
	pod, err := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Errorf("Pod Get Error: %v", err.Error())
		return nil, err
	}

	return pod.Labels, nil
}

func (kac *K8sApiClient) GetPodEnvlByNSAndName(namespace string, name string) ([]v1.EnvVar, error) {
	k8sClient := getClient()
	pod, err := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Logger().Errorf("Pod Get Error: %v", err.Error())
		return nil, err
	}

	return pod.Spec.Containers[0].Env, nil
}

func (kac *K8sApiClient) GetSeldonDeployment(namespace string, sldpName string) (*seldonV1.SeldonDeployment, error) {
	sldp, err := seldonClient.MachinelearningV1().SeldonDeployments(namespace).Get(context.TODO(), sldpName, metav1.GetOptions{})
	if err != nil || sldp == nil {
		logger.Logger().Error("get seldonDeployment err, ", err)
		return nil, err
	}
	return sldp, nil
}
