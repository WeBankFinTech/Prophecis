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
	"github.com/jinzhu/gorm"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/models"
	"mlss-mf/pkg/restapi/operations/model_result"
	"time"

	"github.com/go-openapi/strfmt"
)

func GetResultByModelName(params model_result.GetResultByModelNameParams) (*models.Results, error) {
	if params.ModelName == "" {
		logger.Logger().Error("model name can't be empty string")
		err := fmt.Errorf("model name can't be empty string")
		return nil, err
	}
	results, err := dao.GetResultByModelName(params.ModelName)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get result list by model name, error: %v", err)
			err := fmt.Errorf("fail to get result list by model name, error: %v", err)
			return nil, err
		}
		return nil, nil
	}
	resp := models.Results{}
	timeFormat := "2006-01-02 15:04:05"
	for _, v := range results {
		modelVersion := models.ModelVersionInfo{
			CreationTimestamp: v.ModelVersion.CreationTimestamp.Format(timeFormat),
			EnableFlag:        v.ModelVersion.EnableFlag,
			FileName:          v.ModelVersion.FileName,
			Filepath:          v.ModelVersion.Filepath,
			ID:                v.ModelVersion.ID,
			LatestFlag:        v.ModelVersion.LatestFlag,
			ModelID:           v.ModelVersion.ModelID,
			PushTimestamp:     v.ModelVersion.PushTimestamp.Format(timeFormat),
			Source:            v.ModelVersion.Source,
			TrainingFlag:      v.ModelVersion.TrainingFlag,
			TrainingID:        v.ModelVersion.TrainingId,
			Version:           v.ModelVersion.Version,
		}
		model := models.ModelInfo{
			CreationTimestamp:    v.Model.CreationTimestamp.Format(timeFormat),
			EnableFlag:           v.Model.EnableFlag,
			UserID:               v.Model.UserID,
			GroupID:              v.Model.GroupID,
			ModelName:            v.Model.ModelName,
			ModelType:            v.Model.ModelName,
			ModelVersion:         &modelVersion,
			Position:             v.Model.Position,
			Reamrk:               v.Model.Reamrk,
			ServiceID:            v.Model.ServiceID,
			UpdateTimestamp:      v.Model.UpdateTimestamp.Format(timeFormat),
			ModelLatestVersionID: v.Model.ModelLatestVersionID,
		}
		result := models.Result{
			EnableFlag:        v.EnableFlag,
			ID:                v.Id,
			ResultMsg:         v.ResultMsg,
			TrainingID:        v.TrainingId,
			CreationTimestamp: strfmt.DateTime(v.CreationTimestamp),
			UpdateTimestamp:   strfmt.DateTime(v.UpdateTimestamp),
			Model:             &model,
		}
		resp = append(resp, &result)
	}

	return &resp, nil
}

func DeleteResultByID(params model_result.DeleteResultByIDParams) error {
	if params.ResultID <= 0 {
		logger.Logger().Errorf("result id error, result id: %d", params.ResultID)
		err := fmt.Errorf("model name can't be empty string")
		return err
	}
	err := dao.DeleteResultByID(params.ResultID)
	if err != nil {
		logger.Logger().Errorf("fail to delete result  by id, id: %d, error: %v", params.ResultID, err)
		err := fmt.Errorf("fail to delete result  by id, id: %d, error: %v", params.ResultID, err)
		return err
	}
	return nil
}

func CreateResult(params model_result.CreateResultParams) error {
	if params.Result == nil {
		logger.Logger().Error("params result can't be nil")
		err := fmt.Errorf("params result can't be nil")
		return err
	}
	now := time.Now()
	result := gormmodels.Result{
		ModelId:           params.Result.ModelID,
		ModelVersionId:    params.Result.ModelVersionID,
		TrainingId:        *params.Result.TrainingID,
		ResultMsg:         *params.Result.ResultMsg,
		EnableFlag:        params.Result.EnableFlag,
		CreationTimestamp: now,
		UpdateTimestamp:   now,
	}
	err := dao.CreateResult(&result)
	if err != nil {
		logger.Logger().Errorf("fail to add result, result: %+v, error: %v", params.Result, err)
		err := fmt.Errorf("fail to add result, result: %+v, error: %v", params.Result, err)
		return err
	}
	return nil
}

func UpdateResultByID(params model_result.UpdateResultByIDParams) error {
	if params.ResultID <= 0 || params.Result == nil {
		logger.Logger().Errorf("params error, result Id: %d, result: %+v", params.ResultID, params.Result)
		err := fmt.Errorf("params error, result Id: %d, result: %+v", params.ResultID, params.Result)
		return err
	}
	now := time.Now()
	result := gormmodels.Result{
		ModelId:         params.Result.ModelID,
		ModelVersionId:  params.Result.ModelVersionID,
		TrainingId:      *params.Result.TrainingID,
		ResultMsg:       *params.Result.ResultMsg,
		EnableFlag:      params.Result.EnableFlag,
		UpdateTimestamp: now,
	}
	err := dao.UpdateResultByID(&result, params.ResultID)
	if err != nil {
		logger.Logger().Errorf("fail to update result, result Id: %d, result: %+v, error: %v", params.ResultID, params.Result, err)
		err := fmt.Errorf("fail to update result, result Id: %d, result: %+v, error: %v", params.ResultID, params.Result, err)
		return err
	}
	return nil
}
