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
package workflow

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"webank/DI/commons/service"
	"webank/DI/pkg/util/helm"
	"webank/DI/pkg/util/kubectl"
)

/**
*	delete training job with the job name
**/

func DeleteJob(name, namespace, trainingType string) error {
	jobName := fmt.Sprintf("%s-%s", name, trainingType)

	appInfoFileName, err := kubectl.SaveAppConfigMapToFile(jobName, "app", namespace)
	if err != nil {
		log.Debugf("Failed to SaveAppConfigMapToFile due to %v", err)
		return err
	}

	result, err := kubectl.UninstallAppsWithAppInfoFile(appInfoFileName, namespace)
	if err != nil {
		log.Warnf("Failed to UninstallAppsWithAppInfoFile due to %v", err)
		log.Warnln("manually delete the following resource:")
	}
	fmt.Printf("%s", result)

	err = kubectl.DeleteAppConfigMap(jobName, namespace)
	if err != nil {
		log.Warningf("Delete configmap %s failed, please clean it manually due to %v.", jobName, err)
		log.Warningf("Please run `kubectl delete -n %s cm %s`", namespace, jobName)
	}

	return nil
}

/**
*	Submit training job
**/

func SubmitJob(name string, trainingType string, namespace string, values interface{}, chart string, cmd string, req *service.JobDeploymentRequest) error {
	found := kubectl.CheckAppConfigMap(fmt.Sprintf("%s-%s", name, trainingType), namespace)
	if found {
		return fmt.Errorf("the job %s is already exist, please delete it first. use 'arena delete %s'", name, name)
	}

	// 1. Generate value file
	valueFileName, err := helm.GenerateValueFile(values)
	if err != nil {
		return err
	}

	// 2. Generate Template file
	template, err := helm.GenerateHelmTemplate(name, namespace, valueFileName, chart, cmd, req)
	if err != nil {
		return err
	}

	// 3. Generate AppInfo file
	appInfoFileName, err := kubectl.SaveAppInfo(template, namespace)
	if err != nil {
		return err
	}

	// 4. Keep value file in configmap
	chartName := helm.GetChartName(chart)
	chartVersion, err := helm.GetChartVersion(chart)
	if err != nil {
		return err
	}

	err = kubectl.CreateAppConfigmap(name,
		trainingType,
		namespace,
		valueFileName,
		appInfoFileName,
		chartName,
		chartVersion)
	if err != nil {
		return err
	}
	err = kubectl.LabelAppConfigmap(name, trainingType, namespace, kubectl.JOB_CONFIG_LABEL)
	if err != nil {
		return err
	}

	// 5. Create Application
	_, err = kubectl.UninstallAppsWithAppInfoFile(appInfoFileName, namespace)
	if err != nil {
		log.Debugf("Failed to UninstallAppsWithAppInfoFile due to %v", err)
	}

	result, err := kubectl.InstallApps(template, namespace)
	fmt.Printf("%s", result)
	if err != nil {
		// clean configmap
		log.Infof("clean up the config map %s because creating application failed.", name)
		log.Warnf("Please clean up the training job by using `arena delete %s --type %s`", name, trainingType)
		return err
	}

	// 6. Clean up the template file
	if log.GetLevel() != log.DebugLevel {
		err = os.Remove(valueFileName)
		if err != nil {
			log.Warnf("Failed to delete %s due to %v", valueFileName, err)
		}

		err = os.Remove(template)
		if err != nil {
			log.Warnf("Failed to delete %s due to %v", template, err)
		}

		err = os.Remove(appInfoFileName)
		if err != nil {
			log.Warnf("Failed to delete %s due to %v", appInfoFileName, err)
		}
	}

	return nil
}
