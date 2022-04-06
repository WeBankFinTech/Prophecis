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
	"fmt"
	"mlss-mf/pkg/common/config"
	"mlss-mf/pkg/restapi/operations/model_deploy"
	"strings"

	seldon "github.com/seldonio/seldon-core/operator/apis/machinelearning.seldon.io/v1"
	"github.com/seldonio/seldon-core/operator/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	model "mlss-mf/pkg/models"
)

const (
	ModelSKLEARN    = "SKLEARN"
	ModelXGBoost    = "XGBOOST"
	ModelTensorflow = "TENSORFLOW"
	CUSTOM          = "CUSTOM"
)

func DefineSingleDeployment(service *model.ServicePost, user string) *seldon.SeldonDeployment {
	containerEngineCPU := config.GetMFconfig().ResourceConfig.ContainerEngineCpu
	containerEngineMemory := config.GetMFconfig().ResourceConfig.ContainerEngineMemory
	containerEngineGPU := config.GetMFconfig().ResourceConfig.ContainerEngineGpu
	predictUnitType := seldon.PredictiveUnitType(seldon.MODEL)
	modelType := service.ServicePostModels[0].ModelType
	image := service.ServicePostModels[0].Image
	resourceList := defineRL(*service.CPU, *service.Memory, *service.Gpu)
	sldpLabel := map[string]string{"MLSS-USER": user}
	specName := "classifier"
	if strings.ToUpper(modelType) == CUSTOM {
		specName = "custom-predictor"
	}
	sldp := seldon.SeldonDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "SeldonDeployment",
			//APIVersion: "machinelearning.seldon.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      *service.ServiceName,
			Namespace: *service.Namespace,
			Labels:    sldpLabel,
		},
		Spec: seldon.SeldonDeploymentSpec{
			Name: *service.ServiceName,
			Predictors: []seldon.PredictorSpec{
				{
					Name: "default",
					Graph: seldon.PredictiveUnit{
						Name:             specName,
						Endpoint:         &seldon.Endpoint{Type: "REST"},
						Parameters:       []seldon.Parameter{},
						EnvSecretRefName: "seldon-init-container-secret",
						Type:             &predictUnitType,
					},
					SvcOrchSpec: seldon.SvcOrchSpec{
						Resources: &v1.ResourceRequirements{
							Limits:   defineRL(containerEngineCPU, containerEngineMemory, containerEngineGPU),
							Requests: defineRL(containerEngineCPU, containerEngineMemory, containerEngineGPU),
						},
					},
					EngineResources: v1.ResourceRequirements{
						Limits:   defineRL(containerEngineCPU, containerEngineMemory, containerEngineGPU),
						Requests: defineRL(containerEngineCPU, containerEngineMemory, containerEngineGPU),
					},
					ComponentSpecs: []*seldon.SeldonPodSpec{
						{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: specName,
										Resources: v1.ResourceRequirements{
											Requests: defineRL(*service.CPU, *service.Memory, *service.Gpu),
											Limits:   defineRL(*service.CPU, *service.Memory, *service.Gpu),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	setContainer(modelType, image, service.ServicePostModels, resourceList, &sldp)

	return &sldp
}

func setContainer(modelType string, image string, models model.ServicePostModels, resourceList v1.ResourceList, sldp *seldon.SeldonDeployment) {
	endpointType := models[0].EndpointType
	modelName := models[0].ModelName
	if strings.ToUpper(modelType) != CUSTOM {
		sldp.Spec.Predictors[0].Graph.ModelURI = models[0].Source
		sldp.Spec.Predictors[0].Graph.Implementation = definePredictUnitType(models[0].ModelType)
	}
	if strings.ToUpper(modelType) == ModelTensorflow {
		if endpointType == "GRPC" {
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = seldon.EndpointType("GRPC")
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/tfserving-proxy_grpc:1.3.0"
		} else {
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = seldon.EndpointType("REST")
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/tfserving-proxy_rest:1.3.0"
		}
		parameters := models[0].Parameters
		sldp.Spec.Predictors[0].Name = modelName
		sldp.Spec.Predictors[0].Graph.Name = modelName
		sldp.Spec.Predictors[0].Graph.Endpoint.Type = seldon.EndpointType("REST")
		sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = modelName
		//sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[1].Name = "tfserving"
		//sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[1].Image = "tensorflow/serving:2.1.0"
		engine := v1.Container{
			Name:  "tfserving",
			Image: "tensorflow/serving:2.1.0",
			Resources: v1.ResourceRequirements{
				Requests: resourceList,
				Limits:   resourceList,
			},
		}
		sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers = append(sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers, engine)
		if len(parameters) > 0 {
			for _, v := range parameters {
				sldp.Spec.Predictors[0].Graph.Parameters = append(sldp.Spec.Predictors[0].Graph.Parameters, seldon.Parameter{
					Type:  seldon.ParmeterType(v.Type),
					Name:  v.Name,
					Value: v.Value,
				})
			}
		}
	} else if strings.ToUpper(modelType) == ModelSKLEARN {
		if endpointType == "GRPC" {
			sldp.Spec.Transport = "grpc"
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = ""
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "classifier"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/sklearnserver_grpc:0.3"
		} else {
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "classifier"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/sklearnserver_rest:0.3"
		}
	} else if strings.ToUpper(modelType) == ModelXGBoost {
		if endpointType == "GRPC" {
			sldp.Spec.Transport = "grpc"
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = ""
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "classifier"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/xgboostserver_grpc:0.4"
		} else {
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "classifier"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = "seldonio/xgboostserver_rest:0.4"
		}
	} else if strings.ToUpper(modelType) == CUSTOM {
		if endpointType == "GRPC" {
			sldp.Spec.Transport = "grpc"
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = seldon.EndpointType("GRPC")
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "custom-predictor"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = image
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].ImagePullPolicy = "IfNotPresent"
		} else {
			sldp.Spec.Predictors[0].Graph.Endpoint.Type = seldon.EndpointType("REST")
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Name = "custom-predictor"
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Image = image
			sldp.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].ImagePullPolicy = "IfNotPresent"
		}
	}
}

func definePredictUnitType(modelType string) *seldon.PredictiveUnitImplementation {
	var impl seldon.PredictiveUnitImplementation
	if strings.ToUpper(modelType) == ModelSKLEARN {
		impl = seldon.PredictiveUnitImplementation(constants.PrePackedServerSklearn)
	}
	if strings.ToUpper(modelType) == ModelTensorflow {
		impl = seldon.PredictiveUnitImplementation(constants.PrePackedServerTensorflow)
	}
	if strings.ToUpper(modelType) == ModelXGBoost {
		impl = seldon.PredictiveUnitImplementation("XGBOOST_SERVER")
	}
	if strings.ToUpper(modelType) != ModelXGBoost && strings.ToUpper(modelType) != ModelTensorflow && strings.ToUpper(modelType) != ModelSKLEARN {
		impl = seldon.PredictiveUnitImplementation(seldon.MODEL)
	}
	return &impl
}

//Set Resource List
func defineRL(cpu float64, memory float64, gpu string) v1.ResourceList {
	rl := v1.ResourceList{}
	rl[v1.ResourceCPU] = resource.MustParse(fmt.Sprintf("%.1f", cpu))
	rl[v1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%.1fGi", memory))
	if len(gpu) > 0 {
		rl["nvidia.com/gpu"] = resource.MustParse(gpu)
	}
	return rl
}

func UpdateResourceModel(params model_deploy.UpdateServiceParams, source string, sldpInstance *seldon.SeldonDeployment) *seldon.SeldonDeployment {
	sldpInstance.Spec.Predictors[0].ComponentSpecs[0].Spec.Containers[0].Resources = v1.ResourceRequirements{
		Limits:   defineRL(params.Service.CPU, params.Service.Memory, params.Service.Gpu),
		Requests: defineRL(params.Service.CPU, params.Service.Memory, params.Service.Gpu),
	}
	sldpInstance.Spec.Predictors[0].Graph.ModelURI = source
	return sldpInstance
}
