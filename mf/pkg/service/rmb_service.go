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
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/storage/client"
	"mlss-mf/storage/storage/grpc_storage"
	"os"
	"strings"
)

func DownloadRmbLogByEventID(eventId int64, user string, isSA bool) ([]byte, string, error) {
	eventInDb, err := dao.GetEventById(eventId)
	if err != nil {
		logger.Logger().Errorf("fail to get event  information, eventId: %d, error: %s",
			eventId, err.Error())
		return nil, "", err
	}

	err = PermissionCheck(user, eventInDb.UserId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, "", err
	}

	s3Path := eventInDb.RmbS3path
	if len(s3Path) < 5 {
		logger.Logger().Errorf("s3Path length not correct, s3Path: %s", s3Path)
		return nil, "", fmt.Errorf("s3Path length not correct, s3Path: %s", s3Path)
	}
	s3Array := strings.Split(s3Path[5:], "/")
	logger.Logger().Debugf("lk-test s3Path: %v", s3Path)
	if len(s3Array) < 2 {
		logger.Logger().Errorf("s3 Path is not correct, s3Array: %+v", s3Array)
		return nil, "", errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]

	if bucketName != "mlss-mf" {
		err := fmt.Errorf("bucket name is not 'mlss-mf', bucket name: %s", bucketName)
		return nil, "", err
	}

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return nil, "", err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		err = fmt.Errorf("Storer client is nil ")
		return nil, "", err
	}

	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	s3Path = s3Path + "/" + eventInDb.RmbRespFileName
	// 1. download file
	in := grpc_storage.CodeDownloadRequest{
		S3Path: s3Path,
	}
	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	if err != nil {
		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
		return nil, "", err
	}

	// 2. read file body
	bytes, err := ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	if err != nil {
		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
		return nil, "", err
	}

	err = os.RemoveAll(out.HostPath)
	if err != nil {
		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
		return nil, "", err
	}

	return bytes, base64.StdEncoding.EncodeToString([]byte(eventInDb.RmbRespFileName)), nil
}
