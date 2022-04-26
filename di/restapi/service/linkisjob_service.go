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
	linkis "webank/DI/pkg/client/linkis1_0"
)

var LinkisJobService LinkisJobServiceIF

func init() {
	LinkisJobService = &LinkisJobServiceImpl{}
}

type LinkisJobServiceIF interface {
	GetLinkisJobStatus(exeId string, user string) (*linkis.StatusData, error)
	GetLinkisJobLog(taskID int64, user string) (*linkis.LogData, error)
}

type LinkisJobServiceImpl struct {
}

func (*LinkisJobServiceImpl) GetLinkisJobStatus(exeId string, user string) (*linkis.StatusData, error) {
	linkisClient := linkis.GetLinkisClient()

	status, err := linkisClient.Status(exeId, user)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return status, err
}

func (*LinkisJobServiceImpl) GetLinkisJobLog(taskID int64, user string) (*linkis.LogData, error) {
	linkisClient := linkis.GetLinkisClient()

	logPath, err := linkisClient.GetLogPath(taskID, user)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	logData, err := linkisClient.GetOpenLog(logPath, user)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return logData, err
}
