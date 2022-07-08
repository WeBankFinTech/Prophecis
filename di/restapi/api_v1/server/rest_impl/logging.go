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

package rest_impl

import (
	"webank/DI/commons/logger"
	"webank/DI/restapi/api_v1/server/operations"
	"webank/DI/restapi/api_v1/server/operations/models"
	"webank/DI/restapi/api_v1/server/operations/training_data"

	logr "github.com/sirupsen/logrus"
)

func logWithPostModelParams(params models.PostModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	//data[logger.LogkeyModelFilename] = params.Manifest.Header.Filename
	//data[logger.LogkeyModelFilename] = params.Manifest.Header.Filename

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithDeleteModelParams(params models.DeleteModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithGetModelParams(params models.GetModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithRepeatModelParams(params models.RetryModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithDownloadModelDefinitionParams(params models.DownloadModelDefinitionParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithDownloadTrainedModelParams(params models.DownloadTrainedModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

//
//func logWithEMetricsParams(params training_data.GetEMetricsParams) *logr.Entry {
//	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
//
//	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
//	data[logger.LogkeyTrainingID] = params.ModelID
//
//	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
//}

func logWithLoglinesParams(params training_data.GetLoglinesParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)

	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data[logger.LogkeyTrainingID] = params.ModelID

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithGetLogsParams(params models.GetLogsParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
	data[logger.LogkeyTrainingID] = params.ModelID
	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithGetListModelsParams(params models.ListModelsParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithGetMetricsParams(params models.GetMetricsParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
	data[logger.LogkeyTrainingID] = params.ModelID
	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)

	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithUpdateStatusParams(params models.PatchModelParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
	data[logger.LogkeyTrainingID] = params.ModelID
	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	data["status"] = params.Payload.Status
	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}

func logWithGetDashboards(params operations.GetDashboardsParams) *logr.Entry {
	data := logger.NewDlaaSLogData(logger.LogkeyRestAPIService)
	data[logger.LogkeyUserID] = getUserID(params.HTTPRequest)
	return &logr.Entry{Logger: logr.StandardLogger(), Data: data}
}
