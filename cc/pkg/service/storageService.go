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
	"math"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/datasource"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
	"os"
	"path"
)

func GetAllStorage(page int64, size int64) (models.PageStorageList, error) {
	var allStorage []*models.Storage
	var err error
	offSet := common.GetOffSet(page, size)
	if page > 0 && size > 0 {
		allStorage, err = repo.GetAllStorageByOffSet(offSet, size)
		if err != nil {
			logger.Logger().Error("Get all storage by offset err, ", err)
			return models.PageStorageList{}, err
		}
	} else {
		allStorage, err = repo.GetAllStorage()
		if err != nil {
			logger.Logger().Error("Get all storage err, ", err)
			return models.PageStorageList{}, err
		}
	}
	total, err := repo.CountStorage()
	if err != nil {
		logger.Logger().Error("Get count storage err, ", err)
		return models.PageStorageList{}, err
	}
	var totalPage float64
	if size > 0 {
		totalPage = math.Ceil(float64(total) / float64(size))
	}

	var pageStorage = models.PageStorageList{
		Models:     allStorage,
		TotalPage:  int64(totalPage),
		Total:      total,
		PageNumber: page,
		PageSize:   size,
	}

	return pageStorage, nil
}

func AddStorage(storage models.Storage) (*models.Storage, error) {
	//mkdir -p $Path/data
	//mkdir -p $Path/result
	//chmod -R 700 $Path
	dataPath := path.Join(storage.Path, "data")
	if ok, _ := common.PathExists(dataPath); !ok {
		err := os.MkdirAll(dataPath, 0700)
		if err != nil {
			logger.Logger().Errorf("fail to make dir %q, err: %v\n", dataPath, err)
			return nil, err
		}
	}
	resultPath := path.Join(storage.Path, "result")
	if ok, _ := common.PathExists(resultPath); !ok {
		err := os.MkdirAll(resultPath, 0700)
		if err != nil {
			logger.Logger().Errorf("fail to make dir %q, err: %v\n", resultPath, err)
			return nil, err
		}
	}
	if err := os.Chmod(storage.Path, 0700); err != nil {
		logger.Logger().Debugf("fail to chmod path %q, err :%v\n", storage.Path, err)
		return nil, err
	}

	db := datasource.GetDB()
	err := repo.AddStorage(storage, db)
	if err != nil {
		logger.Logger().Errorf("Get storage by path err, ", err)
		return nil, err
	}
	path, err := repo.GetStorageByPath(storage.Path, db)
	if err != nil {
		logger.Logger().Errorf("Get storage by path err, ", err)
		return nil, err
	}
	return path, nil
}

func GetStorageById(storageId int64) (models.Storage, error) {
	return repo.GetStorageById(storageId)
}

func DeleteStorageById(storageId int64) error {
	tx := datasource.GetDB().Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	groupStorageByStorageId, err := repo.GetAllGroupStorageByStorageId(tx, storageId)
	if err != nil {
		logger.Logger().Errorf("DeleteStorageById Repo GetAllGroupStorageByStorageId error:%v", err.Error())
		return err
	}

	if len(groupStorageByStorageId) > 0 {
		err := repo.DeleteGroupStorageByStorageId(tx, storageId)
		if nil != err {
			tx.Rollback()
			return err
		}
	}
	storageObj, err := repo.GetStorageById(storageId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := repo.DeleteStorageById(tx, storageId); nil != err {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	if ok, _ := common.PathExists(storageObj.Path); ok {
		if err := os.RemoveAll(storageObj.Path); err != nil {
			return err
		}
	}

	return nil
}

func GetStorageByPath(path string) (*models.Storage, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	storage, err := repo.GetStorageByPath(path, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return storage, err
}

func GetDeleteStorageByPath(path string) (*models.Storage, error) {
	log := logger.Logger()
	db := datasource.GetDB()
	storage, err := repo.GetDeleteStorageByPath(path, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return storage, nil
}

func UpdateStorage(storage models.Storage) (*models.Storage, error) {
	dataPath := path.Join(storage.Path, "data")
	if ok, _ := common.PathExists(dataPath); !ok {
		err := os.MkdirAll(dataPath, 0700)
		if err != nil {
			logger.Logger().Errorf("fail to make dir %q, err: %v\n", dataPath, err)
			return nil, err
		}
	}
	resultPath := path.Join(storage.Path, "result")
	if ok, _ := common.PathExists(resultPath); !ok {
		err := os.MkdirAll(resultPath, 0700)
		if err != nil {
			logger.Logger().Errorf("fail to make dir %q, err: %v\n", resultPath, err)
			return nil, err
		}
	}
	if err := os.Chmod(storage.Path, 0700); err != nil {
		logger.Logger().Debugf("fail to chmod path %q, err :%v\n", storage.Path, err)
		return nil, err
	}
	log := logger.Logger()
	db := datasource.GetDB()
	err := repo.UpdateStorage(storage, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	path, err := repo.GetStorageByPath(storage.Path, db)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return path, err
}
