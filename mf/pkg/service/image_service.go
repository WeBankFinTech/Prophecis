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
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"mlss-mf/storage/client"
	"mlss-mf/storage/storage/grpc_storage"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/restapi/operations/image"
	"strconv"

	"github.com/jinzhu/copier"

	"time"

	"github.com/google/uuid"
)

var (
	imageDao        *dao.ImageDao
	modelVersionDao *dao.ModelVersionDao
)

const (
	IMAGE_STATUS_CREATING = "CREATING"
	IMAGE_STATUS_BUILDING = "BUILDING"
	IMAGE_STATUS_PUSHING  = "PUSHING"
	IMAGE_STATUS_COMPLETE = "COMPLETE"
	IMAGE_STATUS_FAIL     = "FAIL"
	DOCKERFILE_PATH       = "/etc/mlss/Dockerfile"
)

type ImageService struct {
}

func (s ImageService) AddImage(currentUserId string, groupId int64, params image.CreateImageParams) (*gormmodels.Image, error) {
	newImage := &gormmodels.Image{}
	err := copier.Copy(&newImage, &params.Image)
	if err != nil {
		logger.Logger().Error("Add image copy params image error,", err)
		return nil, err
	}

	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		logger.Logger().Error("Get user error,", err)
		return nil, err
	}
	newImage.EnableFlag = true
	newImage.UserID = user.Id
	newImage.Status = "Creating"
	newImage.GroupID = groupId
	newImage.LastUpdatedTimestamp = time.Now()
	newImage.CreationTimestamp = time.Now()
	err = imageDao.AddImage(newImage)
	if err != nil {
		logger.Logger().Error("Add image error,", err.Error())
		return newImage, err
	}

	//Get Model Version & S3 Path
	modelVersion, err := modelVersionDao.GetModelVersion(newImage.ModelVersionID)
	if err != nil {
		logger.Logger().Error("Add Image error,", err.Error())
		return newImage, err
	}
	source := modelVersion.Source
	modelName := modelVersion.Model.ModelName
	go s.intergateImage(newImage.ID, newImage.ImageName, source, modelVersion.FileName, modelName)

	return newImage, err
}

func (s ImageService) DeleteImage(params image.DeleteImageParams) error {
	err := imageDao.DeleteImage(params.ImageID)
	if err != nil {
		logger.Logger().Error("Delete image error,", err)
		return err
	}
	return nil
}

func (s ImageService) RecoverImage(imageName string) (*gormmodels.Image, error) {
	image, err := imageDao.RecoverImage(imageName)
	if err != nil {
		logger.Logger().Error("Recover image error,", err)
		return image, err
	}
	return image, nil
}

func (s ImageService) UpdateImage(params image.UpdateImageParams) error {
	newImage := &gormmodels.Image{}
	err := copier.Copy(&newImage, &params.Image)
	if err != nil {
		logger.Logger().Error("Update image copy params image error,", err)
		return err
	}
	getImage, err := imageDao.GetImage(params.ImageID)
	if err != nil {
		logger.Logger().Error("Get image error,", err)
		return err
	}
	modelVersion, err := modelVersionDao.GetModelVersion(newImage.ModelVersionID)
	if err != nil {
		logger.Logger().Error("Get model version error,", err.Error())
		return err
	}
	source := modelVersion.Source
	go s.intergateImage(getImage.ID, getImage.ImageName, source, modelVersion.FileName, modelVersion.Model.ModelName)
	err = imageDao.UpdateImage(params.ImageID, newImage)
	if err != nil {
		logger.Logger().Error("Update image error,", err)
		return err
	}
	return nil
}

func (s ImageService) GetImage(params image.GetImageParams) (*gormmodels.Image, error) {
	getImage, err := imageDao.GetImage(params.ImageID)
	if err != nil {
		logger.Logger().Error("Get image error,", err)
		return &gormmodels.Image{}, err
	}
	model, err := modelDao.GetImageModel(strconv.Itoa(int(getImage.ModelVersion.ModelID)))
	if err != nil {
		logger.Logger().Error("Get model error,", err)
		return &gormmodels.Image{}, err
	}
	getImage.ModelVersion.Model = model
	return getImage, err
}

func (s ImageService) GetImageById(imageId int64) (*gormmodels.Image, error) {
	getImage, err := imageDao.GetImage(imageId)
	if err != nil {
		logger.Logger().Error("Get image error,", err)
		return &gormmodels.Image{}, err
	}
	model, err := modelDao.GetImageModel(strconv.Itoa(int(getImage.ModelVersion.ModelID)))
	if err != nil {
		logger.Logger().Error("Get model error,", err)
		return &gormmodels.Image{}, err
	}
	getImage.ModelVersion.Model = model
	return getImage, err
}

func (s ImageService) ListImage(params image.ListImageParams, currentUserId, queryStr, clusterName string, isSA bool) ([]*gormmodels.Image, int64, error) {
	var images []*gormmodels.Image
	var count int64
	var err error
	if isSA {
		images, err = imageDao.ListImages(*params.Size, (*params.Page-1)*(*params.Size), queryStr, clusterName)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Error("List images error, ", err)
				return nil, 0, err
			}
			return nil, 0, nil
		}
		for _, image := range images {
			model, err := modelDao.GetImageModel(strconv.Itoa(int(image.ModelVersion.ModelID)))
			if err != nil {
				if err.Error() != gorm.ErrRecordNotFound.Error() {
					logger.Logger().Error("Get model error,", err)
					return nil, 0, err
				}
			}
			image.ModelVersion.Model = model
		}
		count, err = imageDao.ListImagesCount(queryStr, clusterName)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Error("List images error, ", err)
				return nil, 0, err
			}
			return nil, 0, nil
		}
	} else {
		images, err = imageDao.ListImagesByUser(*params.Size, (*params.Page-1)*(*params.Size), queryStr, clusterName, currentUserId)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Error("List images error, ", err)
				return nil, 0, err
			}
			return nil, 0, nil
		}
		for _, image := range images {
			model, err := modelDao.GetImageModel(strconv.Itoa(int(image.ModelVersion.ModelID)))
			if err != nil {
				if err.Error() != gorm.ErrRecordNotFound.Error() {
					logger.Logger().Error("Get model error,", err)
					return nil, 0, err
				}
			}
			image.ModelVersion.Model = model
		}
		count, err = imageDao.ListImagesCountByUser(queryStr, clusterName, currentUserId)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Error("List images error, ", err)
				return nil, 0, err
			}
			return nil, 0, nil
		}
	}
	return images, count, err
}

func (s ImageService) ListImageByModelVersionId(params image.ListImageByModelVersionIDParams, currentUserId string, isSA bool) ([]*gormmodels.Image, error) {
	var images []*gormmodels.Image
	var err error
	logger.Logger().Infof("ListImageByModelVersionId isSA:%v", isSA)
	if !isSA {
		model, err := modelDao.GetModelByModelVersionId(params.ModelVersionID)
		if err != nil {
			logger.Logger().Error("ModelVersionList error", err.Error())
			return nil, err
		}
		logger.Logger().Infof("ListImageByModelVersionId:%d", model.UserID)
		err = PermissionCheck(currentUserId, model.UserID, nil, isSA)
		if err != nil {
			logger.Logger().Errorf("Permission Check Error:" + err.Error())
			return nil, err
		}
	}

	if isSA {
		images, err = imageDao.ListImagesByModelVersionId(params.ModelVersionID)
		if err != nil {
			logger.Logger().Error("List images error, ", err)
			return nil, err
		}
		for _, v := range images {
			model, err := modelDao.GetImageModel(strconv.Itoa(int(v.ModelVersion.ModelID)))
			if err != nil {
				logger.Logger().Error("Get model error,", err)
				return nil, err
			}
			v.ModelVersion.Model = model
		}
	} else {
		images, err = imageDao.ListImagesByUserAndModelVersionId(params.ModelVersionID, currentUserId)
		if err != nil {
			logger.Logger().Error("List images error, ", err)
			return nil, err
		}
		for _, v := range images {
			model, err := modelDao.GetImageModel(strconv.Itoa(int(v.ModelVersion.ModelID)))
			if err != nil {
				logger.Logger().Error("Get model error,", err)
				return nil, err
			}
			v.ModelVersion.Model = model
		}
	}
	return images, err
}

func (s ImageService) downloadModel(source string) (string, error) {
	//Download Model ZIP From MinIO
	storageClient, err := client.NewStorage()
	req := &grpc_storage.CodeDownloadRequest{
		S3Path: source,
	}
	res, err := storageClient.Client().DownloadCode(context.TODO(), req)
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", errors.New("storage download model response is nil.")
	}
	logger.Logger().Info("Read path：", fmt.Sprintf("%v/%v", res.HostPath, res.FileName))
	filepath := fmt.Sprintf("%v/%v", res.HostPath, res.FileName)
	return filepath, nil
}

func (s ImageService) intergateImage(imageID int64, imageName string, modelSource string, filename string, modelName string) {
	dockerBuildPath := "/data/ImageBuild/" + uuid.New().String()
	dest := dockerBuildPath + "/model"
	fmt.Print(dest) //MKDIR Temp Dir
	err := os.MkdirAll(dest, os.ModePerm)

	if err != nil {
		logger.Logger().Error("MKDir error:" + err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}

	//COPY Docker File to Dir
	err = s.copyDockerFile(dockerBuildPath)
	if err != nil {
		logger.Logger().Error("copy dockerfile failed:", err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}

	// Unzip Model Matrial
	src, err := s.downloadModel(modelSource + "/" + filename)
	if err != nil {
		logger.Logger().Error(err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}
	// src := ""
	err = Unzip(src, dest)
	if err != nil {
		logger.Logger().Error("Unzip error:" + err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}

	err = s.buildImage(imageID, imageName, modelName, dockerBuildPath)
	if err != nil {
		logger.Logger().Error("build image:" + err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}

	err = s.pushImage(imageID, imageName)
	if err != nil {
		logger.Logger().Error("push image:" + err.Error())
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		// return err
	}
	// return nil
}

func (s ImageService) buildImage(imageID int64, imageName string, modelName string, path string) error {
	s.updateImageIntergateStatus(imageID, IMAGE_STATUS_BUILDING)
	command := "docker build -t " + imageName + " --build-arg MODEL_NAME=" + modelName + " " + path
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger().Error("Build Image Error:", err.Error())
		logger.Logger().Error("Build Image Error:", string(out))
		s.updateImageIntergateErrorMessage(imageID, "Build Image Error:"+string(out))
		return err
	}
	return nil
}

func (s ImageService) pushImage(imageID int64, imageName string) error {
	s.updateImageIntergateStatus(imageID, IMAGE_STATUS_PUSHING)
	command := "docker push " + imageName
	// _ = exec.Command("/bin/sh", "-c", command)
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger().Error("Pushing Image Error:", err.Error())
		logger.Logger().Error("Pushing Image Error:", string(out))
		s.updateImageIntergateErrorMessage(imageID, "Pushing Image Error:"+string(out))
		return err
	}
	s.updateImageIntergateStatus(imageID, IMAGE_STATUS_COMPLETE)
	return nil
}

func (s ImageService) tagImage(imageID int64, imageName string) error {
	s.updateImageIntergateStatus(imageID, IMAGE_STATUS_BUILDING)
	path := ""
	command := "docker build -t " + imageName + " " + path
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.Output()
	if err != nil {
		s.updateImageIntergateStatus(imageID, IMAGE_STATUS_FAIL)
		logger.Logger().Error("Build Image Error:", err.Error())
		return err
	}
	fmt.Println(out)
	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ImageService) updateImageIntergateStatus(imageID int64, status string) {
	image := gormmodels.Image{
		Status:               status,
		LastUpdatedTimestamp: time.Now(),
	}
	err := imageDao.UpdateImage(imageID, &image)
	if err != nil {
		logger.Logger().Error("Update Image Status error: ", status)
	}
}

func (s ImageService) updateImageIntergateErrorMessage(imageID int64, errorMessage string) {
	err := imageDao.UpdateImageMsg(imageID, errorMessage)
	if err != nil {
		logger.Logger().Error("Update Image msg error: ", err)
	}
}

func (s ImageService) copyDockerFile(dstFloder string) error {
	dstName := dstFloder + "/" + "Dockerfile"
	src, err := os.Open(DOCKERFILE_PATH)
	if err != nil {
		logger.Logger().Error("open src file error:", err.Error())
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Logger().Error("open dst file error:", err.Error())
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		logger.Logger().Error("copy file error:", err.Error())
		return err
	}
	return err
}
