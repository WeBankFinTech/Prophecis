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
	"errors"
	"io"
	"io/ioutil"
	"os"
	"webank/DI/commons/logger"
	storageClient "webank/DI/storage/client"
	"webank/DI/storage/storage/grpc_storage"

	"github.com/google/uuid"
)

const (
	actualFlowFolderName = "flow"
)

var ExperimentService ExperimentServiceIF = &ExperimentServiceImpl{}
var log = logger.GetLogger()

type ExperimentServiceIF interface {
	UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error)
	UploadCode(closer io.ReadCloser) (string, error)
}
type ExperimentServiceImpl struct {
}

func (expService *ExperimentServiceImpl) UploadCode(closer io.ReadCloser) (string, error) {
	hostpath := "/data/oss-storage/"
	filename := uuid.New().String() + "file.zip"
	filepath := hostpath + filename
	//out, err := os.Create(filename)
	//wt := bufio.NewWriter(out)
	//_, err = io.Copy(wt, closer)
	fileBytes, err := ioutil.ReadAll(closer)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	err = ioutil.WriteFile(filepath, fileBytes, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	//上传 code.zip 后, 更新flowJson 中每个 node 的 codePath
	//创建 storage rpc client, 用于 上传code.zip
	//ctx := context.TODO()
	sClient, err := storageClient.NewStorage()
	if err != nil {
		log.Info("storageClient init Error " + err.Error())
		return "", err
	}

	req := &grpc_storage.CodeUploadRequest{
		FileName: filename,
		HostPath: hostpath,
	}
	ctx := context.TODO()
	if sClient == nil {
		log.Info("storageClient init Error,sclient is nil")
		return "", errors.New("rpc upload code zip error, sclient is nil")
	}
	response, err := sClient.Client().UploadCode(ctx, req)
	if err != nil {
		log.Error("sclient upload code error: " + err.Error())
		return "", err
	}
	err = os.Remove(filepath)
	if err != nil {
		log.Error("" + err.Error())
	}

	return response.S3Path, err

}

func (expService *ExperimentServiceImpl) UploadCodeBucket(closer io.ReadCloser, bucket string) (string, error) {
	hostpath := "/data/oss-storage/"
	filename := uuid.New().String() + "file.zip"
	filepath := hostpath + filename
	fileBytes, err := ioutil.ReadAll(closer)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	err = ioutil.WriteFile(filepath, fileBytes, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	//上传 code.zip 后, 更新flowJson 中每个 node 的 codePath
	//创建 storage rpc client, 用于 上传code.zip
	//ctx := context.TODO()
	sClient, err := storageClient.NewStorage()
	if err != nil {
		log.Info("storageClient init Error " + err.Error())
		return "", err
	}

	req := &grpc_storage.CodeUploadRequest{
		FileName: filename,
		HostPath: hostpath,
		Bucket:   bucket,
	}
	ctx := context.TODO()
	if sClient == nil {
		log.Info("storageClient init Error,sclient is nil")
		return "", errors.New("rpc upload code zip error, sclient is nil")
	}
	response, err := sClient.Client().UploadCode(ctx, req)
	if err != nil {
		log.Error("sclient upload code error: " + err.Error())
		return "", err
	}
	err = os.Remove(filepath)
	if err != nil {
		log.Error("" + err.Error())
	}
	return response.S3Path, err
}
