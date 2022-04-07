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
	"mlss-controlcenter-go/pkg/logger"
	"strings"
)

type MLFLOW_TYPE string

const (
	MLFLOW_EXP MLFLOW_TYPE = "experiment"
	MLFLOW_RUN MLFLOW_TYPE = "run"
	MLFLOW_ART MLFLOW_TYPE = "artifacts"
)

func CheckMLFlowResource(requestUri string, user string, isSA bool) (bool, error) {
	if isSA {
		return true, nil
	}

	mlflowType := getMLFlowResourceType(requestUri)
	isAccess := true
	var err interface{}
	if mlflowType == MLFLOW_EXP {
		isAccess, err = checkMLFlowExpRes(requestUri, user)
		if err != nil {
			logger.Logger().Error("Check MLFlow Resource error", err.(error).Error())
			return isAccess, err.(error)
		}
	} else if mlflowType == MLFLOW_RUN {
		isAccess, err = checkMLFlowRunRes(requestUri, user)
		if err != nil {
			logger.Logger().Error("Check MLFlow Resource error", err.(error).Error())
			return isAccess, err.(error)
		}
	} else if mlflowType == MLFLOW_ART {
		isAccess, err = checkMLFlowARTRes(requestUri, user)
		if err != nil {
			logger.Logger().Error("Check MLFlow Resource error", err.(error).Error())
			return isAccess, err.(error)
		}
	}
	return isAccess, nil
}

func getMLFlowResourceType(requestUri string) MLFLOW_TYPE {
	if strings.Contains(requestUri, "/mlflow/runs/") || strings.Contains(requestUri, "/mlflow/metrics/") {
		return MLFLOW_RUN
	} else if strings.Contains(requestUri, "/mlflow/artifacts/") {
		return MLFLOW_ART
	} else {
		return MLFLOW_EXP
	}
}

func checkMLFlowExpRes(requestUri string, user string) (bool, error) {
	return false, nil
}

func checkMLFlowRunRes(requestUri string, user string) (bool, error) {
	return false, nil
}

func checkMLFlowARTRes(requestUri string, user string) (bool, error) {
	return false, nil
}

func getRunIdByURI(requestUri string) (*string, error) {
	return nil, nil
}

func getExpIdByURI(requestUri string) (*string, error) {
	return nil, nil
}
