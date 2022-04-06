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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"io"
	"io/ioutil"
	"mlss-mf/pkg/common"
	"mlss-mf/pkg/common/config"
	cc "mlss-mf/pkg/common/controlcenter"
	"mlss-mf/pkg/common/gormmodels"
	"mlss-mf/pkg/dao"
	"mlss-mf/pkg/logger"
	"mlss-mf/pkg/models"
	"mlss-mf/pkg/restapi/operations/report"
	"mlss-mf/storage/client"
	"mlss-mf/storage/storage/grpc_storage"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func checkVersion(version string) (bool, error) {
	ok, err := regexp.MatchString("^[Vv][1-9][0-9]*$", version)
	return ok, err
}

func checkQueryStr(str string) (bool, error) {
	ok, err := regexp.MatchString("^[0-9a-zA-Z-]*$", str)
	return ok, err
}

func CreateReport(params report.CreateReportParams) (*models.PostReportResp, error) {
	token := params.HTTPRequest.Header.Get(cc.CcAuthToken)
	currentUserId := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	auth := CheckGroupAuth(currentUserId, *params.Report.GroupID, isSA)
	if !auth {
		err := errors.New("current user no authorization")
		return nil, err
	}
	userDao := &dao.UserDao{}
	user, err := userDao.GetUserByName(currentUserId)
	if err != nil {
		logger.Logger().Errorf("fail to get user information, error: %s", err.Error())
		return nil, err
	}
	reportName := params.Report.ReportName
	reportVersion := params.Report.ReportVersion
	modelName := params.Report.ModelName
	modelVersion := params.Report.ModelVersion
	groupId := params.Report.GroupID
	if reportVersion == "" {
		reportVersion = "v1"
	}

	// check report version
	matched, err := checkVersion(reportVersion)
	if err != nil {
		logger.Logger().Errorf("fail to check report version, report version: %s, err: %s\n", reportVersion, err.Error())
		return nil, err
	}
	if !matched {
		logger.Logger().Errorf("report version not matched, report verison: %s\n", reportVersion)
		err = fmt.Errorf("report version not match by regexp '^[Vv][1-9][0-9]*$'")
		return nil, err
	}

	modelDao := &dao.ModelDao{}
	modelInfo, err := modelDao.GetModelByModelNameAndGroupId(*modelName, *groupId)
	if err != nil {
		logger.Logger().Errorf("fail to get model info, modelName: %s, groupId: %d, error: %s\n", *modelName, *groupId, err.Error())
		return nil, err
	}
	modelVerisonDao := &dao.ModelVersionDao{}
	modelVersionInfo, err := modelVerisonDao.GetModelVersionByModelIdAndVersion(modelInfo.ID, *modelVersion)
	if err != nil {
		logger.Logger().Errorf("fail to get model version info, modelId: %d, modelVersion: %s, error: %s\n", modelInfo.ID, *modelVersion, err.Error())
		return nil, err
	}

	// report version not exist, add report and report version
	s3Path := params.Report.S3Path
	if len(params.Report.RootPath) > 0 {
		err = ccClient.UserStorageCheck(token, currentUserId, params.Report.RootPath)
		if err != nil {
			return nil, err
		}
		fileBytes, err := ioutil.ReadFile(path.Join(params.Report.RootPath, params.Report.ChildPath))
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(fileBytes)
		closer := ioutil.NopCloser(reader)
		pathSplit := strings.Split(params.Report.ChildPath, "/")
		uploadReportResp, _, err := UploadReport(pathSplit[len(pathSplit)-1], closer)
		if err != nil {
			logger.Logger().Errorf("fail to upload report, filename: %s, error: %s", pathSplit[len(pathSplit)-1], err.Error())
			return nil, err
		}
		s3Path = uploadReportResp.S3Path
	}

	var reportId int64 = 0
	reportInfo, err := dao.GetReportBaseByReportNameAndGroupId(*reportName, *groupId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report info report name: %s, group id: %d, err: %s\n", *reportName, *groupId, err.Error())
			return nil, err
		} else {
			logger.Logger().Infof("report record not found , reportName: %s, groupId: %d\n",
				*reportName, *groupId)
			//create report
			now := time.Now()
			reportObj := gormmodels.ReportBase{
				ReportName:            *reportName,
				ModelId:               modelInfo.ID,
				ModelVersionId:        modelVersionInfo.ID,
				UserId:                user.Id,
				GroupId:               *groupId,
				CreationTimestamp:     now,
				UpdateTimestamp:       now,
				EnableFlag:            1,
				ReportLatestVersionId: 0,
			}
			if err := dao.CreateReport(&reportObj); err != nil {
				logger.Logger().Errorf("fail to add report, reportObj: %+v, error: %s", reportObj, err.Error())
				return nil, err
			}
			reportId = reportObj.Id
		}
	} else {
		reportId = reportInfo.Id
	}

	var reportVersionId int64 = 0
	reportVersionResp := ""
	if reportInfo.ReportLatestVersionId == 0 {
		// report version not exist, add report and report version
		if s3Path == "" {
			err = errors.New("when not repeat upload file, s3Path can not be empty string, ")
			return nil, err
		}
		now := time.Now()
		reportVersionObj := gormmodels.ReportVersionBase{
			ReportName:        *reportName,
			Version:           reportVersion, //todo
			FileName:          params.Report.FileName,
			FilePath:          path.Join(params.Report.RootPath, params.Report.ChildPath),
			Source:            s3Path,
			PushId:            0,
			CreationTimestamp: now,
			UpdateTimestamp:   now,
			EnableFlag:        1,
			ReportId:          reportId,
		}
		if err := dao.CreateReportVersion(&reportVersionObj); err != nil {
			logger.Logger().Errorf("fail to add report version, report version object: %+v, err: %s\n", reportVersionObj, err.Error())
			return nil, err
		}
		reportVersionId = reportVersionObj.Id
		reportVersionResp = reportVersion
		reportInfo.ReportLatestVersionId = reportVersionObj.Id
	} else {
		//reportVersionInfo, err := dao.GetReportVersionByReportIdAndReportVersion(reportId, reportVersion)
		reportVersionInfo, err := dao.GetReportVersionByReportVersionId(reportInfo.ReportLatestVersionId)
		if err != nil {
			logger.Logger().Errorf("fail to get report version information, reportName: %s, reportVersion: %s, error: %s\n",
				*reportName, reportVersion, err.Error())
			return nil, err

		} else {
			version := ""
			versionLatestBytes := []byte(reportVersionInfo.Version)
			versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
			if err != nil {
				logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", reportVersionInfo.Version, err.Error())
				return nil, err
			}

			versionParamBytes := []byte(reportVersion)
			versionParamInt64, err := strconv.ParseInt(string(versionParamBytes[1:]), 10, 64)
			if err != nil {
				logger.Logger().Errorf("fail to parse params version, version: %s, err: %s\n", reportVersion, err.Error())
				return nil, err
			}
			if versionLatestInt64 < versionParamInt64 {
				version = fmt.Sprintf("v%d", versionParamInt64)
			} else {
				version = fmt.Sprintf("v%d", versionLatestInt64+1)
			}

			now := time.Now()
			if s3Path == "" {
				s3Path = reportVersionInfo.Source
			}
			reportVersionObj := gormmodels.ReportVersionBase{
				ReportName:        *reportName,
				Version:           version,
				FileName:          params.Report.FileName,
				FilePath:          reportVersionInfo.FilePath,
				Source:            s3Path,
				PushId:            0,
				CreationTimestamp: now,
				UpdateTimestamp:   now,
				EnableFlag:        1,
				ReportId:          reportId,
			}
			if err := dao.CreateReportVersion(&reportVersionObj); err != nil {
				logger.Logger().Errorf("fail to add report version, report version object: %+v, err: %s\n", reportVersionObj, err.Error())
				return nil, err
			}
			reportVersionId = reportVersionObj.Id
			reportVersionResp = version
			reportInfo.ReportLatestVersionId = reportVersionObj.Id
		}
	}

	//update report latest version id
	m := make(map[string]interface{})
	m["report_latest_version_id"] = reportVersionId
	logger.Logger().Debugf("update report version id in table report, report id: %d, m: %+v\n", reportId, m)
	if err := dao.UpdateReportById(reportId, m); err != nil {
		logger.Logger().Errorf("fail to update report latest version id, report id: %d, m: %+v, err: %s\n", reportId, m, err.Error())
		return nil, err
	}

	return &models.PostReportResp{ReportID: reportId, ReportVersion: reportVersionResp, ReportVersionID: reportVersionId}, nil
}

func ListReportVersionsByReportID(params report.ListReportVersionsByReportIDParams, user string, isSA bool) (*models.ListReportVersionResp, error) {

	if !isSA {
		reportInDb, err := dao.GetReportByID(params.ReportID)
		logger.Logger().Infof("ListReportVersionsByReportIDParams report: %+v\n", reportInDb)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
				return nil, err
			}
			return nil, nil
		}
		err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
		if err != nil {
			logger.Logger().Errorf("Permission Check Error:" + err.Error())
			return nil, err
		}
	}

	reportVersionListDao, err := dao.ListReportVersionsByReportID(params.ReportID)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to list report versions by report id, report id: %d, err: %s\n", params.ReportID, err.Error())
			return nil, err
		}
		return nil, nil
	}

	//logger.Logger().Infof("ListReportVersionsByReportID reportVersionListDao[0]: %+v", reportVersionListDao[0])
	result := models.ListReportVersionResp{}
	if len(reportVersionListDao) > 0 {
		for idx := range reportVersionListDao {
			userInfo, err := userDao.GetUserByUserId(reportVersionListDao[idx].Event.UserId)
			if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Errorf("fail to get user information by user id, user id: %d, err: %s\n", reportVersionListDao[idx].Event.UserId, err.Error())
				return nil, err
			}
			goupInfo, err := groupDao.GetGroupByGroupId(reportVersionListDao[idx].Report.GroupId)
			if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Errorf("fail to get group information by user id, user id: %d, err: %s\n", reportVersionListDao[idx].Report.GroupId, err.Error())
				return nil, err
			}
			createTime := strfmt.DateTime(reportVersionListDao[idx].Event.CreationTimestamp)
			updateTime := strfmt.DateTime(reportVersionListDao[idx].Event.UpdateTimestamp)
			event := models.Event{
				Dcn:               &reportVersionListDao[idx].Event.Dcn,
				EnableFlag:        &reportVersionListDao[idx].Event.EnableFlag,
				FileID:            &reportVersionListDao[idx].Event.FileId,
				FileName:          &reportVersionListDao[idx].Event.FileName,
				FileType:          &reportVersionListDao[idx].Event.FileType,
				FpsFileID:         &reportVersionListDao[idx].Event.FpsFileId,
				HashValue:         &reportVersionListDao[idx].Event.FpsHashValue,
				ID:                &reportVersionListDao[idx].Event.Id,
				Idc:               &reportVersionListDao[idx].Event.Idc,
				Params:            &reportVersionListDao[idx].Event.Params,
				Status:            &reportVersionListDao[idx].Event.Status,
				CreationTimestamp: &createTime,
				UpdateTimestamp:   &updateTime,
			}
			var enableFlag int8 = 0
			if userInfo.EnableFlag {
				enableFlag = 1
			}
			user := models.UserInfo{
				EnableFlag: enableFlag,
				Gid:        userInfo.Gid,
				GUIDCheck:  int8(userInfo.GuidCheck),
				ID:         userInfo.ID,
				Name:       userInfo.Name,
				Remarks:    userInfo.Remarks,
				UID:        userInfo.UID,
				UserType:   userInfo.Type,
			}
			group := models.Group{
				ClusterName:    goupInfo.ClusterName,
				DepartmentID:   goupInfo.DepartmentId,
				DepartmentName: goupInfo.DepartmentName,
				GroupType:      goupInfo.GroupType,
				ID:             goupInfo.ID,
				Name:           goupInfo.Name,
				Remarks:        goupInfo.Remarks,
				RmbDcn:         goupInfo.RmbDcn,
				RmbIdc:         goupInfo.RmbIdc,
				ServiceID:      goupInfo.ServiceId,
				SubsystemID:    goupInfo.SubsystemId,
				SubsystemName:  goupInfo.SubsystemName,
			}
			report := models.ReportBase{
				CreationTimestamp:     strfmt.DateTime(reportVersionListDao[idx].Report.CreationTimestamp),
				EnableFlag:            &reportVersionListDao[idx].Report.EnableFlag,
				GroupID:               reportVersionListDao[idx].Report.GroupId,
				ID:                    reportVersionListDao[idx].Report.Id,
				ModelID:               reportVersionListDao[idx].Report.ModelId,
				ModelVersionID:        reportVersionListDao[idx].Report.ModelVersionId,
				ReportLatestVersionID: reportVersionListDao[idx].Report.ReportLatestVersionId,
				ReportName:            reportVersionListDao[idx].Report.ReportName,
				UpdateTimestamp:       strfmt.DateTime(reportVersionListDao[idx].Report.UpdateTimestamp),
				UserID:                reportVersionListDao[idx].Report.UserId,
			}
			reportVersion := models.ListReportVersionBase{
				CreationTimestamp: strfmt.DateTime(reportVersionListDao[idx].CreationTimestamp),
				EnableFlag:        &reportVersionListDao[idx].EnableFlag,
				Event:             &event,
				FileName:          reportVersionListDao[idx].FileName,
				Filepath:          reportVersionListDao[idx].FilePath,
				ID:                reportVersionListDao[idx].Id,
				PushID:            reportVersionListDao[idx].PushId,
				ReportID:          reportVersionListDao[idx].ReportId,
				ReportName:        reportVersionListDao[idx].ReportName,
				Source:            reportVersionListDao[idx].Source,
				UpdateTimestamp:   strfmt.DateTime(reportVersionListDao[idx].UpdateTimestamp),
				Version:           reportVersionListDao[idx].Version,
				User:              &user,
				Report:            &report,
				Group:             &group,
			}
			result = append(result, &reportVersion)
		}
	}
	return &result, nil
}

// func GetReportByModelNameAndModelVersion(modelName, modelVersion string) (*models.Report, error) {
// 	reportInDb, err := dao.GetReportByModelNameAndModelVersion(modelName, modelVersion)
// 	if err != nil {
// 		logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
// 		return nil, err
// 	}
// 	resp := generateReport(reportInDb)

// 	return resp, nil
// }

func generateReport(r *gormmodels.Report) *models.Report {
	if r == nil {
		return nil
	}
	timeFormat := "2006-01-02 15:04:05"
	createTime := strfmt.DateTime(r.ReportVersion.Event.UpdateTimestamp)
	updateTime := strfmt.DateTime(r.ReportVersion.Event.CreationTimestamp)
	event := models.Event{
		Dcn:               &r.ReportVersion.Event.Dcn,
		EnableFlag:        &r.ReportVersion.Event.EnableFlag,
		FileID:            &r.ReportVersion.Event.FileId,
		FileName:          &r.ReportVersion.Event.FileName,
		FileType:          &r.ReportVersion.Event.FileType,
		FpsFileID:         &r.ReportVersion.Event.FpsFileId,
		HashValue:         &r.ReportVersion.Event.FpsHashValue,
		ID:                &r.ReportVersion.Event.Id,
		Idc:               &r.ReportVersion.Event.Idc,
		Params:            &r.ReportVersion.Event.Params,
		RmbS3path:         r.ReportVersion.Event.RmbS3path,
		RmbRespFileName:   r.ReportVersion.Event.RmbRespFileName,
		Status:            &r.ReportVersion.Event.Status,
		UpdateTimestamp:   &createTime,
		CreationTimestamp: &updateTime,
	}
	modelVersion := models.ModelVersionBase{
		ID:                r.ModelVersion.ID,
		FileName:          r.ModelVersion.FileName,
		Filepath:          r.ModelVersion.Filepath,
		LatestFlag:        r.ModelVersion.LatestFlag,
		ModelID:           r.ModelVersion.ModelID,
		PushTimestamp:     r.ModelVersion.PushTimestamp.Format(timeFormat),
		Source:            r.ModelVersion.Source,
		TrainingFlag:      r.ModelVersion.TrainingFlag,
		TrainingID:        r.ModelVersion.TrainingId,
		Version:           r.ModelVersion.Version,
		EnableFlag:        r.ModelVersion.EnableFlag,
		CreationTimestamp: r.ModelVersion.CreationTimestamp.Format(timeFormat),
		PushID:            r.ModelVersion.PushId,
	}

	model := models.ModelBase{
		ID:                   r.Model.ID,
		ModelLatestVersionID: r.Model.ModelLatestVersionID,
		ModelName:            r.Model.ModelName,
		ModelType:            r.Model.ModelType,
		Position:             r.Model.Position,
		Reamrk:               r.Model.Reamrk,
		ServiceID:            r.Model.ServiceID,
		UserID:               r.Model.UserID,
		EnableFlag:           r.Model.EnableFlag,
		GroupID:              r.Model.GroupID,
		CreationTimestamp:    r.Model.CreationTimestamp.Format(timeFormat),
		UpdateTimestamp:      r.Model.UpdateTimestamp.Format(timeFormat),
	}
	group := models.Group{
		ClusterName:    r.Group.ClusterName,
		DepartmentID:   r.Group.DepartmentId,
		DepartmentName: r.Group.DepartmentName,
		GroupType:      r.Group.GroupType,
		ID:             r.Group.ID,
		Name:           r.Group.Name,
		Remarks:        r.Group.Remarks,
		SubsystemID:    r.Group.SubsystemId,
		SubsystemName:  r.Group.SubsystemName,
		ServiceID:      r.Group.ServiceId,
	}
	reportVersion := models.ReportVersionBase{
		CreationTimestamp: strfmt.DateTime(r.ReportVersion.CreationTimestamp),
		EnableFlag:        &r.ReportVersion.EnableFlag,
		FileName:          r.ReportVersion.FileName,
		Filepath:          r.ReportVersion.FilePath,
		ID:                r.ReportVersion.Id,
		PushID:            r.ReportVersion.PushId,
		ReportID:          r.ReportVersion.ReportId,
		ReportName:        r.ReportVersion.ReportName,
		Source:            r.ReportVersion.Source,
		UpdateTimestamp:   strfmt.DateTime(r.ReportVersion.UpdateTimestamp),
		Version:           r.ReportVersion.Version,
	}
	var enable_flag int8 = 1
	if !r.User.EnableFlag {
		enable_flag = 0
	}
	user := models.UserInfo{
		EnableFlag: enable_flag,
		Gid:        r.User.Gid,
		GUIDCheck:  int8(r.User.GuidCheck),
		ID:         r.User.ID,
		Name:       r.User.Name,
		Remarks:    r.User.Remarks,
		UID:        r.User.UID,
		UserType:   r.User.Type,
	}
	resp := models.Report{
		EnableFlag:            &r.EnableFlag,
		Group:                 &group,
		ID:                    r.Id,
		Model:                 &model,
		UserID:                r.UserId,
		CreationTimestamp:     strfmt.DateTime(r.CreationTimestamp),
		UpdateTimestamp:       strfmt.DateTime(r.UpdateTimestamp),
		LatestReportVersion:   &reportVersion,
		ModelVersion:          &modelVersion,
		ReportLatestVersionID: r.ReportLatestVersionId,
		Event:                 &event,
		User:                  &user,
	}

	return &resp
}

func GetReportByID(id int64, user string, isSA bool) (*models.Report, error) {
	reportInDb, err := dao.GetReportByID(id)
	logger.Logger().Infof("GetReportByID report: %+v\n", reportInDb)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
			return nil, err
		}
		return nil, nil
	}
	err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	if reportInDb.ReportVersion.PushId > 0 {
		event, err := dao.GetEventById(reportInDb.ReportVersion.PushId)
		if err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				logger.Logger().Debugf("push event record not found, id: %d, err: %s\n", reportInDb.ReportVersion.PushId, err.Error())
			} else {
				logger.Logger().Debugf("fail to get push event information, id: %d, err: %s\n", reportInDb.ReportVersion.PushId, err.Error())
				return nil, err
			}
		} else {
			reportInDb.ReportVersion.Event = *event
		}
	}

	resp := generateReport(reportInDb)

	return resp, nil
}

func ListReports(currentPage, pageSize int64, currentUserName, queryStr string, isSA bool) (*models.ListReportsResp, error) {
	defaultParams := report.NewListReportsParams()
	if currentPage <= 0 {
		currentPage = *defaultParams.CurrentPage
	}
	if pageSize <= 0 {
		pageSize = *defaultParams.PageSize
	}
	limit := pageSize
	offset := (currentPage - 1) * pageSize

	if queryStr != "" {
		matched, err := checkQueryStr(queryStr)
		if err != nil {
			logger.Logger().Errorf("fail to check query string, queryStr: %s, error: %s", queryStr, err.Error())
			return nil, err
		}
		if !matched {
			logger.Logger().Errorf("queryStr not matched, queryStr: %s\n", queryStr)
			err = fmt.Errorf("queryStr not match by regexp '^[0-9a-zA-Z-]*$'")
			return nil, err
		}
	}

	if isSA {
		total, err := dao.CountReports(queryStr)
		if err != nil {
			logger.Logger().Errorf("fail to get reports total, error: %s", err.Error())
			return nil, err
		}
		reportSInDb, err := dao.ListReports(limit, offset, queryStr)
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				logger.Logger().Errorf("fail to get reports list, error: %s", err.Error())
				return nil, err
			}
			return nil, nil
		}
		for k := range reportSInDb {
			//get event by id
			if reportSInDb[k].ReportVersion.PushId <= 0 {
				continue
			}
			event, err := dao.GetEventById(reportSInDb[k].ReportVersion.PushId)
			if err != nil {
				if err.Error() == gorm.ErrRecordNotFound.Error() {
					logger.Logger().Debugf("push event record not found, id: %d, err: %s\n", reportSInDb[k].ReportVersion.PushId, err.Error())
					continue
				} else {
					logger.Logger().Debugf("fail to get push event information, id: %d, err: %s\n", reportSInDb[k].ReportVersion.PushId, err.Error())
					return nil, err
				}
			}
			reportSInDb[k].ReportVersion.Event = *event
		}
		reportList := make([]*models.Report, 0)
		for _, reportInDb := range reportSInDb {
			report := generateReport(reportInDb)
			reportList = append(reportList, report)
		}
		resp := models.ListReportsResp{
			Reports: reportList,
			Total:   total,
		}
		return &resp, nil
	}

	userDao := &dao.UserDao{}
	user, err := userDao.GetUserByName(currentUserName)
	if err != nil {
		logger.Logger().Errorf("fail to get user information, error: %s", err.Error())
		return nil, err
	}

	total, err := dao.CountReportsByUserId(user.Id, queryStr)
	if err != nil {
		logger.Logger().Errorf("fail to get reports total, error: %s", err.Error())
		return nil, err
	}
	reportSInDb, err := dao.ListReportsByUserId(user.Id, limit, offset, queryStr)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get reports list, error: %s", err.Error())
			return nil, err
		}
		return nil, nil
	}
	reportList := make([]*models.Report, 0)
	for _, reportInDb := range reportSInDb {
		report := generateReport(reportInDb)
		reportList = append(reportList, report)
	}
	resp := models.ListReportsResp{
		Reports: reportList,
		Total:   total,
	}

	return &resp, nil
}

func DeleteReportByID(id int64, user string, isSA bool) error {
	if id <= 0 {
		err := fmt.Errorf("id can't be 0 or negatice number, id: %d", id)
		logger.Logger().Errorf("id can't be 0 or negatice number, id: %d", id)
		return err
	}

	//Permission Check
	reportInDb, err := dao.GetReportByID(id)
	logger.Logger().Infof("GetReportByID report: %+v\n", reportInDb)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
			return err
		}
		return nil
	}
	err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return err
	}

	m := make(map[string]interface{})
	m["enable_flag"] = 0
	err = dao.UpdateReportById(id, m)
	if err != nil {
		logger.Logger().Errorf("fail to delete report by id, id: %d, err: %v\n", id, err)
		return err
	}

	reportVersionInfo, err := dao.ListReportVersionsByReportID(id)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report list by report id, report_id: %d, err: %v\n", id, err)
			return err
		}
	}

	if len(reportVersionInfo) > 0 {
		err = dao.UpdateReportVersionByReportId(id, m)
		if err != nil {
			logger.Logger().Errorf("fail to delete report version by  report id, report id: %d, err: %v\n", id, err)
			return err
		}
		for idx := range reportVersionInfo {
			err = dao.UpdateEventByFileIdAndFileType(reportVersionInfo[idx].Id, config.Push_Event_File_Type_Report, m)
			if err != nil {
				logger.Logger().Errorf("fail to delete report version push event by file_id, file_id: %d, "+
					"file_type: %s, err: %v\n", reportVersionInfo[idx].Id, config.Push_Event_File_Type_Report, err)
				return err
			}
		}
	}

	return nil
}

func DeleteReportVersionByReportID(id int64, user string, isSA bool) error {
	if id <= 0 {
		err := fmt.Errorf("id can't be 0 or negatice number, id: %d", id)
		logger.Logger().Errorf("id can't be 0 or negatice number, id: %d", id)
		return err
	}

	//Permission Check
	reportInDb, err := dao.GetReportByID(id)
	logger.Logger().Infof("GetReportByID report: %+v\n", reportInDb)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
			return err
		}
		return nil
	}
	err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return err
	}

	m := make(map[string]interface{})
	m["enable_flag"] = 0
	err = dao.UpdateReportVersionByReportId(id, m)
	if err != nil {
		logger.Logger().Errorf("fail to delete report version by report id, report id: %d", id)
		return err
	}

	return nil
}

func DownloadReportById(id int64, currentUserId string, isSA bool) ([]byte, int, string, error) {
	reportInfo, err := dao.GetReportByID(id)
	if err != nil {
		logger.Logger().Errorf("fail to get report information, error: %s", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	auth := CheckGroupAuth(currentUserId, reportInfo.GroupId, isSA)
	if !auth {
		return nil, StatusAuthReject, "", nil
	}

	logger.Logger().Debugf("lk-test modelInfo, modelInfo: %+v", reportInfo)
	//tmpDir := "/data/oss-storage/report/" + uuid.New().String()
	//if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
	//	logger.Logger().Errorf("fail to create tmp dir, error: %s", err.Error())
	//	return nil, http.StatusInternalServerError, "", err
	//}

	s3Path := reportInfo.ReportVersion.Source
	if len(s3Path) < 5 {
		logger.Logger().Errorf("s3Path length not correct, s3Path: %s", s3Path)
		return nil, http.StatusInternalServerError, "", fmt.Errorf("s3Path length not correct, s3Path: %s", s3Path)
	}
	s3Array := strings.Split(s3Path[5:], "/")
	logger.Logger().Debugf("lk-test s3Path: %v", s3Path)
	if len(s3Array) < 2 {
		logger.Logger().Errorf("s3 Path is not correct, s3Array: %+v", s3Array)
		return nil, http.StatusInternalServerError, "", errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]

	if bucketName != "mlss-mf" {
		err = fmt.Errorf("bucket name is not 'mlss-mf', bucket name: %s", bucketName)
		return nil, http.StatusInternalServerError, "", err
	}
	// if bucketName == "mlss-mf" {
	// 	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	// 	s3Path = s3Path + "/" + modelInfo.ModelLatestVersion.FileName //单个文件
	// }

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		err = errors.New("Storer client is nil ")
		return nil, http.StatusInternalServerError, "", err
	}

	//if is single file, download
	//bytes := make([]byte, 0)
	//switch strings.Contains(reportInfo.ReportVersion.FileName, ".zip") {
	//case false:
	//	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	//	s3Path = s3Path + "/" + reportInfo.ReportVersion.FileName //single file
	//	// 1. download file
	//	in := grpc_storage.CodeDownloadRequest{
	//		S3Path: s3Path,
	//	}
	//	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
	//		return nil, http.StatusInternalServerError, "", err
	//	}
	//
	//	// 2. read file body
	//	bytes, err = ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
	//		return nil, http.StatusInternalServerError, "", err
	//	}
	//
	//	err = os.RemoveAll(out.HostPath)
	//	if err != nil {
	//		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
	//		return nil, http.StatusInternalServerError, err
	//	}
	//default:
	//	in := grpc_storage.CodeDownloadByDirRequest{
	//		S3Path: s3Path + "/" + reportInfo.ReportVersion.FileName,
	//	}
	//	outs, err := storageClient.Client().DownloadCodeByDir(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code by dir, req s3path: %s,  error: %s", s3Path, err.Error())
	//		return nil, http.StatusInternalServerError, err
	//	}
	//	if outs != nil && len(outs.RespList) > 0 {
	//		for _, resp := range outs.RespList {
	//			logger.Logger().Debugf("download code by dir, hostpath: %s, filename: %s\n", resp.HostPath, resp.FileName)
	//			bytes, err := ioutil.ReadFile(resp.HostPath + "/" + resp.FileName)
	//			if err != nil {
	//				logger.Logger().Errorf("fail to read file, filename: %s", resp.HostPath+"/"+resp.FileName)
	//				return nil, http.StatusInternalServerError, err
	//			}
	//
	//			// write to tmp file
	//			logger.Logger().Debugf("read file %s, content-length: %d\n", resp.HostPath+"/"+resp.FileName, len(bytes))
	//			subPath := strings.Split(resp.FileName, "/")
	//			filename := path.Join(subPath[1:]...) //todo
	//			tmpFile := path.Join(tmpDir, filename)
	//			if err := createDir(tmpFile); err != nil {
	//				logger.Logger().Errorf("fail to create tmp directory of file %q, error: %s", tmpFile, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//			err = ioutil.WriteFile(tmpFile, bytes, os.ModePerm)
	//			// logger.Logger().Debugf("file %q content: %s\n", resp.FileName, string(bytes))
	//			if err != nil {
	//				logger.Logger().Errorf("fail to write data to file %q, error: %s", tmpFile, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//		}
	//
	//		for _, resp := range outs.RespList {
	//			err = os.RemoveAll(resp.HostPath)
	//			if err != nil {
	//				logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", resp.HostPath, err.Error())
	//				return nil, http.StatusInternalServerError, err
	//			}
	//		}
	//
	//		// gzip file
	//		dest := ""
	//		if reportInfo.ReportVersion.FileName != "" {
	//			filenameOnly := getFileNameOnly(reportInfo.ReportVersion.FileName)
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", filenameOnly))
	//		} else {
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", s3Array[len(s3Array)-1]))
	//		}
	//		source := []string{tmpDir}
	//
	//		if err := archiver.Archive(source, dest); err != nil {
	//			logger.Logger().Errorf("fail to archive file, error: %s, source: %s, dest: %s\n", err.Error(), source, dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//
	//		// read from zip file
	//		bytes, err = ioutil.ReadFile(dest)
	//		if err != nil {
	//			logger.Logger().Errorf("fail to read zip file , file: %s", dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//		logger.Logger().Debugf("zip file %q content: %s\n", dest, string(bytes))
	//		if err := os.Remove(dest); err != nil {
	//			logger.Logger().Errorf("fail to remove zip file , file: %s", dest)
	//			return nil, http.StatusInternalServerError, err
	//		}
	//	}
	//}

	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	s3Path = s3Path + "/" + reportInfo.ReportVersion.FileName
	// 1. download file
	in := grpc_storage.CodeDownloadRequest{
		S3Path: s3Path,
	}
	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	if err != nil {
		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	// 2. read file body
	bytes, err := ioutil.ReadFile(out.HostPath + "/" + out.FileName)
	if err != nil {
		logger.Logger().Errorf("fail to read file, filename: %s", out.HostPath+"/"+out.FileName)
		return nil, http.StatusInternalServerError, "", err
	}

	err = os.RemoveAll(out.HostPath)
	if err != nil {
		logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", out.HostPath, err.Error())
		return nil, http.StatusInternalServerError, "", err
	}

	//// remove host tmp file
	//err = os.RemoveAll(tmpDir)
	//if err != nil {
	//	logger.Logger().Errorf("clear folder failed, dir name: %s, error: %s", tmpDir, err.Error())
	//	return nil, http.StatusInternalServerError, "", err
	//}

	return bytes, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(reportInfo.ReportVersion.FileName)), nil
}

func ListReportVersionPushEventsByReportVersionID(reportVersionId, currentPage, pageSize int64, user string, isSA bool) (*models.ReportVersionPushEventResp, error) {
	defaultParams := report.NewListReportVersionPushEventsByReportVersionIDParams()
	if currentPage <= 0 {
		currentPage = *defaultParams.CurrentPage
	}
	if pageSize <= 0 {
		pageSize = *defaultParams.PageSize
	}
	reportVersionInfo, err := dao.GetReportVersionByReportVersionId(reportVersionId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report version information, reportVersionId: %d, error: %s\n", reportVersionId, err)
			return nil, err
		}
		return nil, nil
	}
	reportInfo, err := dao.GetReportBaseByID(reportVersionInfo.ReportId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report information, reportId: %d, error: %s\n", reportVersionInfo.ReportId, err)
			return nil, err
		}
		return nil, nil
	}
	err = PermissionCheck(user, reportInfo.UserId, nil, isSA)
	if err != nil {
		logger.Logger().Errorf("Permission Check Error:" + err.Error())
		return nil, err
	}

	eventsAndUser, err := dao.ListReportPushEventByReportVersionId(reportVersionId, (currentPage-1)*pageSize, pageSize)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report version event list, reportVersionId: %d, error: %s\n", reportVersionId, err)
			return nil, err
		}
		return nil, nil
	}
	logger.Logger().Debugf("ListReportVersionPushEventsByReportVersionID eventsAndUser: %+v\n", eventsAndUser)
	count, err := dao.CountReportPushEventByReportVersionId(reportVersionId)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to count report version event, reportVersionId: %d, error: %s\n", reportVersionId, err)
			return nil, err
		}
		return nil, nil
	}

	resp := reportVersionPushEventResp(count, reportVersionInfo, reportInfo, eventsAndUser)

	return resp, nil
}

func reportVersionPushEventResp(count int64, reportVersion *gormmodels.ReportVersionBase, report *gormmodels.ReportBase,
	eventAndUser []*gormmodels.PushEvent) *models.ReportVersionPushEventResp {
	reportInfo := models.ReportBase{
		CreationTimestamp:     strfmt.DateTime(report.CreationTimestamp),
		EnableFlag:            &report.EnableFlag,
		GroupID:               report.GroupId,
		ID:                    report.Id,
		ModelID:               report.ModelId,
		ModelVersionID:        report.ModelVersionId,
		ReportLatestVersionID: report.ReportLatestVersionId,
		ReportName:            report.ReportName,
		UpdateTimestamp:       strfmt.DateTime(report.UpdateTimestamp),
		UserID:                report.UserId,
	}
	reportVersionInfo := models.ReportVersionBase{
		CreationTimestamp: strfmt.DateTime(reportVersion.CreationTimestamp),
		EnableFlag:        &reportVersion.EnableFlag,
		FileName:          reportVersion.FileName,
		Filepath:          reportVersion.FilePath,
		ID:                reportVersion.Id,
		PushID:            reportVersion.PushId,
		ReportID:          reportVersion.ReportId,
		ReportName:        reportVersion.ReportName,
		Source:            reportVersion.Source,
		UpdateTimestamp:   strfmt.DateTime(reportVersion.UpdateTimestamp),
		Version:           reportVersion.Version,
	}
	events := make([]*models.ReportVersionPushEventRespBase, 0)
	if len(eventAndUser) > 0 {
		for idx := range eventAndUser {
			logger.Logger().Debugf("reportVersionPushEventResp eventsAndUser[%d]: %+v\n", idx, *eventAndUser[idx])
			if eventAndUser[idx] != nil {
				event := models.Event{
					CreationTimestamp: (*strfmt.DateTime)(&eventAndUser[idx].CreationTimestamp),
					Dcn:               &eventAndUser[idx].Dcn,
					EnableFlag:        &eventAndUser[idx].EnableFlag,
					FileID:            &eventAndUser[idx].FileId,
					FileName:          &eventAndUser[idx].FileName,
					FileType:          &eventAndUser[idx].FileType,
					FpsFileID:         &eventAndUser[idx].FpsFileId,
					HashValue:         &eventAndUser[idx].FpsHashValue,
					ID:                &eventAndUser[idx].Id,
					Idc:               &eventAndUser[idx].Idc,
					Params:            &eventAndUser[idx].Params,
					RmbRespFileName:   eventAndUser[idx].RmbRespFileName,
					RmbS3path:         eventAndUser[idx].RmbS3path,
					Status:            &eventAndUser[idx].Status,
					UpdateTimestamp:   (*strfmt.DateTime)(&eventAndUser[idx].UpdateTimestamp),
					Version:           eventAndUser[idx].Version,
				}
				var enableFlag int8 = 0
				if eventAndUser[idx].User.EnableFlag {
					enableFlag = 1
				}
				user := models.UserInfo{
					EnableFlag: enableFlag,
					Gid:        eventAndUser[idx].User.Gid,
					GUIDCheck:  int8(eventAndUser[idx].User.GuidCheck),
					ID:         eventAndUser[idx].User.ID,
					Name:       eventAndUser[idx].User.Name,
					Remarks:    eventAndUser[idx].User.Remarks,
					UID:        eventAndUser[idx].User.UID,
					UserType:   eventAndUser[idx].User.Type,
				}

				reportVersionPushEventRespBase := models.ReportVersionPushEventRespBase{
					Event: &event,
					User:  &user,
				}
				events = append(events, &reportVersionPushEventRespBase)
			}
		}
	} else {
		logger.Logger().Debugln("reportVersionPushEventResp eventsAndUser is nil")
	}

	resp := models.ReportVersionPushEventResp{
		Events:        events,
		Report:        &reportInfo,
		ReportVersion: &reportVersionInfo,
		Total:         count,
	}
	return &resp
}

func UploadReport(fileName string, file io.ReadCloser) (*models.UploadReportResponse, int, error) {
	uid := uuid.New().String()
	hostPath := fmt.Sprintf("/data/oss-storage/model/%v/", uid)
	err := os.MkdirAll(hostPath, os.ModePerm)
	if err != nil {
		logger.Logger().Error("Mkdir err, ", err.Error())
		return nil, StatusError, errors.New("Mkdir err, " + err.Error())
	}

	filePath := hostPath + fileName
	logger.Logger().Infof("lk-test UploadModel filepath: %s", filePath)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Logger().Error("File read err, ", err.Error())
		return nil, StatusError, err
	}
	//err = ioutil.WriteFile(filePath, fileBytes, os.ModePerm)
	err = WriteToFile(filePath, config.Max_RW, fileBytes)
	if err != nil {
		logger.Logger().Error("File write err, ", err.Error())
		return nil, StatusError, err
	}

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return nil, StatusError, err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		return nil, StatusError, errors.New("Storer client is nil")
	}
	req := &grpc_storage.CodeUploadRequest{
		Bucket:   "mlss-mf",
		FileName: fileName,
		HostPath: hostPath,
	}
	//var res *grpc_storage.CodeUploadResponse
	//if strings.HasSuffix(fileName, ".zip") {
	//	res, err = storageClient.Client().UploadModelZip(context.TODO(), req)
	//	if err != nil {
	//		logger.Logger().Error("Storer client upload report err, ", err.Error())
	//		return nil, StatusError, err
	//	}
	//} else {
	//	res, err = storageClient.Client().UploadCode(context.TODO(), req)
	//	if err != nil {
	//		logger.Logger().Error("Storer client upload report err, ", err.Error())
	//		return nil, StatusError, err
	//	}
	//}

	res, err := storageClient.Client().UploadCode(context.TODO(), req)
	if err != nil {
		logger.Logger().Error("Storer client upload report err, ", err.Error())
		return nil, StatusError, err
	}

	// clear folder
	err = os.RemoveAll(hostPath)
	if err != nil {
		logger.Logger().Error("clear folder err, " + err.Error())
		return nil, StatusError, err
	}
	resp := models.UploadReportResponse{
		FileName: fileName,
		Filepath: filePath,
		S3Path:   res.S3Path,
	}
	return &resp, StatusSuccess, nil
}

func checkReportPushEventParams(params *models.PushReportRequest) error {
	if params == nil {
		return errors.New("params can't be nil")
	}
	if params.FactoryName == nil || *params.FactoryName == "" {
		return errors.New("params factorr_name can't be nil or empty string")
	}

	return nil
}

//func PushReportByReportID(reportId int64, pushReportEventParams *models.PushReportRequest, user string,
//	isSA bool) (*models.Event, error) {
//	// check parmas
//	if err := checkReportPushEventParams(pushReportEventParams); err != nil {
//		return nil, err
//	}
//
//	reportInDb, err := dao.GetReportByID(reportId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get report information, reportId: %d, err: %s\n", reportId, err.Error())
//			return nil, err
//		}
//		logger.Logger().Infof("record not exist, reportId: %d, err: %s\n", reportId, err.Error())
//		return nil, err
//	}
//
//	//Permission Check
//	err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
//	if err != nil {
//		logger.Logger().Errorf("Permission Check Error:" + err.Error())
//		return nil, err
//	}
//
//	params, err := json.Marshal(pushReportEventParams)
//	if err != nil {
//		logger.Logger().Errorf("fail to marshal report event params, params: %+v, err: %s\n", *pushReportEventParams, err.Error())
//		return nil, err
//	}
//
//	eventVersion := ""
//	latestPushEvent, err := dao.GetEventById(reportInDb.ReportVersion.PushId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get latest push event information, reportVersionId: %d, eventId: %d,  err: %s\n",
//				reportInDb.ReportLatestVersionId, reportInDb.ReportVersion.PushId, err.Error())
//			return nil, err
//		}
//		eventVersion = "v1"
//	} else {
//		if latestPushEvent.Version != "" {
//			versionLatestBytes := []byte(latestPushEvent.Version)
//			versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
//			if err != nil {
//				logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", latestPushEvent.Version, err.Error())
//				return nil, err
//			}
//			eventVersion = fmt.Sprintf("v%d", versionLatestInt64+1)
//		} else {
//			eventVersion = "v1"
//		}
//	}
//
//	now := time.Now()
//	event := gormmodels.PushEventBase{
//		FileType:          "DATA",
//		FileId:            reportInDb.ReportVersion.Id,
//		FileName:          reportInDb.ReportVersion.FileName,
//		FpsFileId:         "",
//		FpsHashValue:      "",
//		Params:            string(params),
//		Idc:               reportInDb.Group.RmbIdc,
//		Dcn:               reportInDb.Group.RmbDcn,
//		Status:            config.RMB_PushEvent_Status_Creating,
//		CreationTimestamp: now,
//		UpdateTimestamp:   now,
//		EnableFlag:        1,
//		RmbRespFileName:   "",
//		RmbS3path:         "",
//		UserId:            reportInDb.UserId,
//		Version:           eventVersion,
//	}
//	if err := dao.AddEvent(&event); err != nil {
//		logger.Logger().Errorf("fail to add event, event: %+v, err: %s\n", event, err.Error())
//		return nil, err
//	}
//	go pushReport(reportInDb.ReportVersion.Source, reportInDb.ReportVersion.FileName,
//		reportInDb.Group.RmbDcn, reportInDb.Group.ServiceId, event.Id, pushReportEventParams)
//
//	//update report version push id
//	m := make(map[string]interface{})
//	m["push_id"] = event.Id
//	err = dao.UpdateReportVersionById(reportInDb.ReportVersion.Id, m)
//	if err != nil {
//		logger.Logger().Errorf("fail to update report version's push_id, reportVersion id: %d, m: %+v, err: %s\n",
//			reportInDb.ReportVersion.Id, m, err.Error())
//		return nil, err
//	}
//	eventInfo := models.Event{
//		CreationTimestamp: (*strfmt.DateTime)(&event.CreationTimestamp),
//		Dcn:               &event.Dcn,
//		EnableFlag:        &event.EnableFlag,
//		FileID:            &event.FileId,
//		FileName:          &event.FileName,
//		FileType:          &event.FileType,
//		FpsFileID:         &event.FpsFileId,
//		HashValue:         &event.FpsHashValue,
//		ID:                &event.Id,
//		Idc:               &event.Idc,
//		Params:            &event.Params,
//		RmbS3path:         event.RmbS3path,
//		RmbRespFileName:   event.FileName,
//		Status:            &event.Status,
//		UpdateTimestamp:   (*strfmt.DateTime)(&event.UpdateTimestamp),
//		Version:           eventVersion,
//	}
//
//	return &eventInfo, nil
//}

//func PushReportByReportVersionID(reportVersionId int64, pushReportEventParams *models.PushReportRequest, user string, isSA bool) (*models.Event, error) {
//	// check parmas
//	if err := checkReportPushEventParams(pushReportEventParams); err != nil {
//		return nil, err
//	}
//
//	reportVersionInDb, err := dao.GetReportVersionByReportVersionId(reportVersionId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get report version information, reportVersionId: %d, err: %s\n", reportVersionId, err.Error())
//			return nil, err
//		}
//		logger.Logger().Infof("record not exist, reportVersionId: %d, err: %s\n", reportVersionId, err.Error())
//		return nil, nil
//	}
//	reportInDb, err := dao.GetReportBaseByID(reportVersionInDb.ReportId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get report information, reportId: %d, err: %s\n", reportVersionId, err.Error())
//			return nil, err
//		}
//		logger.Logger().Infof("record not exist, reportId: %d, err: %s\n", reportVersionId, err.Error())
//		return nil, nil
//	}
//	//Permission Check
//	err = PermissionCheck(user, reportInDb.UserId, nil, isSA)
//	if err != nil {
//		logger.Logger().Errorf("Permission Check Error:" + err.Error())
//		return nil, err
//	}
//
//	groupInDb, err := groupDao.GetGroupByGroupId(reportInDb.GroupId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get group information, GroupId: %d, err: %s\n", reportInDb.GroupId, err.Error())
//			return nil, err
//		}
//		logger.Logger().Infof("record not exist, GroupId: %d, err: %s\n", reportInDb.GroupId, err.Error())
//		return nil, nil
//	}
//
//	params, err := json.Marshal(pushReportEventParams)
//	if err != nil {
//		logger.Logger().Errorf("fail to marshal report event params, params: %+v, err: %s\n", *pushReportEventParams, err.Error())
//		return nil, err
//	}
//
//	eventVersion := ""
//	latestPushEvent, err := dao.GetEventById(reportVersionInDb.PushId)
//	if err != nil {
//		if err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to get latest push event information, reportVersionId: %d, eventId: %d,  err: %s\n",
//				reportVersionId, reportVersionInDb.PushId, err.Error())
//			return nil, err
//		}
//		eventVersion = "v1"
//	} else {
//		if latestPushEvent.Version != "" {
//			versionLatestBytes := []byte(latestPushEvent.Version)
//			versionLatestInt64, err := strconv.ParseInt(string(versionLatestBytes[1:]), 10, 64)
//			if err != nil {
//				logger.Logger().Errorf("fail to parse version, version: %s, err: %s\n", latestPushEvent.Version, err.Error())
//				return nil, err
//			}
//			eventVersion = fmt.Sprintf("v%d", versionLatestInt64+1)
//		} else {
//			eventVersion = "v1"
//		}
//	}
//
//	now := time.Now()
//	event := gormmodels.PushEventBase{
//		FileType:          "DATA",
//		FileId:            reportVersionInDb.Id,
//		FileName:          reportVersionInDb.FileName,
//		FpsFileId:         "",
//		FpsHashValue:      "",
//		Params:            string(params),
//		Idc:               groupInDb.RmbIdc,
//		Dcn:               groupInDb.RmbDcn,
//		Status:            config.RMB_PushEvent_Status_Creating,
//		CreationTimestamp: now,
//		UpdateTimestamp:   now,
//		EnableFlag:        1,
//		RmbRespFileName:   "",
//		RmbS3path:         "",
//		UserId:            reportInDb.UserId,
//		Version:           eventVersion,
//	}
//
//	if err := dao.AddEvent(&event); err != nil {
//		logger.Logger().Errorf("fail to add event, event: %+v, err: %s\n", event, err.Error())
//		return nil, err
//	}
//
//	go pushReport(reportVersionInDb.Source, reportVersionInDb.FileName, groupInDb.RmbDcn, groupInDb.ServiceId, event.Id, pushReportEventParams)
//
//	//update report version push id
//	m := make(map[string]interface{})
//	m["push_id"] = event.Id
//	err = dao.UpdateReportVersionById(reportVersionInDb.Id, m)
//	if err != nil {
//		logger.Logger().Errorf("fail to update report version's push_id, reportVersion id: %d, m: %+v, err: %s\n",
//			reportVersionInDb.Id, m, err.Error())
//		return nil, err
//	}
//	eventInfo := models.Event{
//		CreationTimestamp: (*strfmt.DateTime)(&event.CreationTimestamp),
//		Dcn:               &event.Dcn,
//		EnableFlag:        &event.EnableFlag,
//		FileID:            &event.FileId,
//		FileName:          &event.FileName,
//		FileType:          &event.FileType,
//		FpsFileID:         &event.FpsFileId,
//		HashValue:         &event.FpsHashValue,
//		ID:                &event.Id,
//		Idc:               &event.Idc,
//		Params:            &event.Params,
//		RmbS3path:         event.RmbS3path,
//		RmbRespFileName:   event.FileName,
//		Status:            &event.Status,
//		UpdateTimestamp:   (*strfmt.DateTime)(&event.UpdateTimestamp),
//		Version:           eventVersion,
//	}
//
//	return &eventInfo, nil
//}

//func pushReport(s3path, filename, rmbDcn string, serviceId, eventId int64, pushReportEventParams *models.PushReportRequest) error {
//	// 1.download file from minio
//	filePathOfDownload, err := downloadFileFromMinio(s3path, filename)
//	if err != nil {
//		// update pushevent
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["update_timestamp"] = time.Now()
//		if err := dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		logger.Logger().Errorf("fail to download report file, s3path: %s, fileName: %s, err: %s\n",
//			s3path, filename, err.Error())
//		return err
//	}
//
//	// 2. upload file to fps
//	// java -jar fps-client-app-1.4.0-fat.jar -t upload -i fpsip-0.0.0.0 -P 7008 -u 账号  -p 密码  --savetype month --savecnt 1 --file test1.txt
//	fpsConfig := config.GetFpsConfig()
//	fpsJdkPath := os.Getenv("FPS_JDK_PATH")
//	cmd := exec.Command(path.Join(os.Getenv("JAVA_HOME"), "/bin/java"), "-jar", path.Join(fpsJdkPath, "fps-client-app-1.4.0-fat.jar"),
//		"-t", "upload", "-i", fpsConfig.Ip, "-P", strconv.FormatInt(fpsConfig.Port, 10), "-u", fpsConfig.User, "-p",
//		fpsConfig.Pwd, "--savetype", fpsConfig.SaveType, "--savecnt", strconv.FormatInt(int64(fpsConfig.SaveCnt), 10), "--file", filePathOfDownload)
//	fpsBytes, err := cmd.Output()
//	if err != nil {
//		logger.Logger().Errorf("fail to exec fps command, cmd: %v, err: %s\n",
//			cmd.Args, err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//
//	fpsResp := parseFPSResp(fpsBytes)
//	//filePath := path.Dir(filePathOfDownload)
//	fileName := path.Base(filePathOfDownload)
//
//	// 3. send rmb event
//	// sh ./bin/wtss-rmbsender.sh  --targetDcn="FA0" --serviceId="11200257" --messageType="sync"
//	// --configPath="./conf/" --msgBody="{\"serviceType\":\"uploadFile\",
//	// \"request\":{\"factoryName\":\"report_factory_name_test\",\"fileType\":\"DATA\",
//	// \"fileName\":\"main.go\",\"fileId\":\"25000000163074214650200325537659000000016279773465023418\",
//	// \"hashValue\":\"cdb7d6b231303426f77b49b6b5b6f312\"}}"
//	rmbConfig := config.GetRmbConfig()
//	if rmbDcn == "" {
//		logger.Logger().Errorf("rmb dcn can't be empty, rmbDcn: %s\n", rmbDcn)
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		err = fmt.Errorf("rmb dcn can't be empty, rmbDcn: %s", rmbDcn)
//		return err
//	}
//	rmbJdkPath := os.Getenv("RMB_JDK_PATH")
//	rmbJdkPath = path.Join(rmbJdkPath, "bin/wtss-rmbsender.sh")
//	if serviceId <= 0 {
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		logger.Logger().Errorf("rmb serviceId can't be empty, report ServiceID: %d\n", serviceId)
//		err = fmt.Errorf("rmb serviceId can't be empty, serviceId: %d", serviceId)
//		return err
//	}
//	serviceIdStr := strconv.FormatInt(serviceId, 10)
//
//	msgBodyStr, err := generateRmbReportReq(pushReportEventParams, fileName, fpsResp.FileId, fpsResp.Hash)
//	if err != nil {
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		logger.Logger().Errorf("fail to generate rmb request body, err: %s\n", err.Error())
//		return err
//	}
//
//	status := config.RMB_PushEvent_Status_Failed
//	command := rmbJdkPath + " " + " --targetDcn='" + rmbDcn + "' --serviceId='" + serviceIdStr +
//		"' --messageType='" + rmbConfig.MessageType + "' --msgBody='" + msgBodyStr + "' --configPath='" + rmbConfig.ConfigPath + "'"
//	cmd = exec.Command("/bin/sh", "-c", command)
//	logger.Logger().Debugf("debug rmb rmbJdkPath: %s, rmb dcn: %s, serviceIdStr:%s,"+
//		" messageType: %s, msgBody: %s, configPathFlag: %s, cmd.Params: %v\n",
//		rmbJdkPath, rmbDcn, serviceIdStr, rmbConfig.MessageType, msgBodyStr, rmbConfig.ConfigPath, cmd.Args)
//	rmbBytes, err := cmd.CombinedOutput()
//	logger.Logger().Debugf("rmbBytes: %s\n", string(rmbBytes))
//	if err != nil {
//		logger.Logger().Errorf("fail to exec rmb command, command: %s, err: %s\n",
//			command, err.Error())
//	} else {
//		flagBytes := []byte("content='{\"code\":")
//		for _, subRmbBytes := range bytes.Split(rmbBytes, []byte("\n")) {
//			//logger.Logger().Debugf("%s\n", string(subRmbBytes))
//			if bytes.Contains(subRmbBytes, flagBytes) {
//				if bytes.Contains(subRmbBytes, []byte("content='{\"code\": 2000")) {
//					logger.Logger().Info("Get RMB Response Success, Reponse Code is 2000, res is "+
//						": ", string(subRmbBytes))
//					status = config.RMB_PushEvent_Status_Success
//				} else {
//					logger.Logger().Errorf("fail to send rmb request, response is: ", string(subRmbBytes))
//					status = config.RMB_PushEvent_Status_Failed
//				}
//				logger.Logger().Info("Get RMB Response: ", string(subRmbBytes))
//				break
//			}
//		}
//	}
//
//	//event := gormmodels.PushEventBase{}
//	rmbTmpDir := "/data/oss-storage/rmb-resp/" + uuid.New().String()
//	rmbRespFileName := fmt.Sprintf("rmb_resp_%s.txt", uuid.New().String())
//	if err := os.MkdirAll(rmbTmpDir, os.ModePerm); err != nil {
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		logger.Logger().Errorf("fail to create tmp directory of file %q, error: %s", rmbTmpDir, err.Error())
//		return err
//	}
//	f, err := os.Create(path.Join(rmbTmpDir, rmbRespFileName))
//	if err != nil {
//		f.Close()
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		logger.Logger().Errorf("fail to create file %q, error: %s", path.Join(rmbTmpDir, rmbRespFileName), err.Error())
//		return err
//	}
//	if _, err := f.Write(rmbBytes); err != nil {
//		f.Close()
//		logger.Logger().Errorf("fail to write data to file %q, error: %s", path.Join(rmbTmpDir, rmbRespFileName), err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//	f.Close()
//
//	storageClient, err := client.NewStorage()
//	if err != nil {
//		logger.Logger().Error("Storer client create err, ", err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//	if storageClient == nil {
//		logger.Logger().Error("Storer client is nil ")
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return errors.New("Storer client is nil")
//	}
//	rmbReq := &grpc_storage.CodeUploadRequest{
//		Bucket:   "mlss-mf",
//		FileName: rmbRespFileName,
//		HostPath: rmbTmpDir,
//	}
//	rmbResp, err := storageClient.Client().UploadCode(context.TODO(), rmbReq)
//	if err != nil {
//		logger.Logger().Errorf("fail to upload rmb resp file, rmb request: %+v, err: %s", err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//	if err := os.RemoveAll(rmbTmpDir); err != nil {
//		logger.Logger().Errorf("fail to remove tmp file path: %s,  err: %s\n", rmbTmpDir, err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//	if err := os.RemoveAll(filePathOfDownload); err != nil {
//		logger.Logger().Errorf("fail to remove download file path: %s,  err: %s\n", filePathOfDownload, err.Error())
//		m := make(map[string]interface{})
//		m["status"] = config.RMB_PushEvent_Status_Failed
//		m["fps_file_id"] = fpsResp.FileId
//		m["fps_hash_value"] = fpsResp.Hash
//		m["update_timestamp"] = time.Now()
//		if err = dao.UpdateEventById(eventId, m); err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//			logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//			return err
//		}
//		return err
//	}
//
//	// update pushevent
//	m := make(map[string]interface{})
//	m["status"] = status
//	m["rmb_resp_file_name"] = rmbRespFileName
//	m["rmb_s3path"] = rmbResp.S3Path
//	m["fps_file_id"] = fpsResp.FileId
//	m["fps_hash_value"] = fpsResp.Hash
//	m["update_timestamp"] = time.Now()
//	err = dao.UpdateEventById(eventId, m)
//	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
//		logger.Logger().Errorf("fail to add event, eventId: %+v, err: %s\n", eventId, err.Error())
//		return err
//	}
//
//	return nil
//}

func parseFPSResp(bytes []byte) *FPSResp {
	subStr := strings.Split(string(bytes), "\n")
	fileid := ""
	hash := ""
	for k := range subStr {
		if strings.HasPrefix(subStr[k], "fileid:") {
			logger.Logger().Debugf("subStr[%d]: %q\n", k, subStr[k])
			fileInfoSubStr := strings.Split(subStr[k], ",")
			for k2 := range fileInfoSubStr {
				if strings.HasPrefix(fileInfoSubStr[k2], "fileid:") {
					fileid = strings.TrimPrefix(fileInfoSubStr[k2], "fileid:")
					fileid = strings.TrimSpace(fileid)
				}
				if strings.HasPrefix(fileInfoSubStr[k2], "hash:") {
					hash = strings.TrimPrefix(fileInfoSubStr[k2], "hash:")
					hash = strings.TrimSpace(hash)
				}
			}
		}
	}
	fps := FPSResp{
		FileId: fileid,
		Hash:   hash,
	}
	logger.Logger().Debugf("parse fps resp result: %+v\n", fps)
	return &fps
}

type FPSResp struct {
	FileId string `json:"fileid"`
	Hash   string `json:"hash"`
}

func downloadFileFromMinio(s3Path, fileName string) (string, error) {
	tmpDir := "/data/oss-storage/minio/" + uuid.New().String()
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		logger.Logger().Errorf("fail to create tmp dir, error: %s", err.Error())
		return "", err
	}
	if len(s3Path) < 5 {
		logger.Logger().Errorf("s3Path length not correct, s3Path: %s", s3Path)
		return "", fmt.Errorf("s3Path length not correct, s3Path: %s", s3Path)
	}

	s3Array := strings.Split(s3Path[5:], "/")
	logger.Logger().Debugf("lk-test s3Path: %v", s3Path)
	if len(s3Array) < 2 {
		logger.Logger().Errorf("s3 Path is not correct, s3Array: %+v", s3Array)
		return "", errors.New("s3 Path is not correct")
	}
	bucketName := s3Array[0]

	if bucketName != "mlss-mf" {
		err := fmt.Errorf("bucket name is not 'mlss-mf', bucket name: %s", bucketName)
		return "", err
	}
	// if bucketName == "mlss-mf" {
	// 	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	// 	s3Path = s3Path + "/" + modelInfo.ModelLatestVersion.FileName //单个文件
	// }

	storageClient, err := client.NewStorage()
	if err != nil {
		logger.Logger().Error("Storer client create err, ", err.Error())
		return "", err
	}
	if storageClient == nil {
		logger.Logger().Error("Storer client is nil ")
		err = fmt.Errorf("Storer client is nil ")
		return "", err
	}

	//if is single file, download
	//switch strings.Contains(fileName, ".zip") {
	//case false:
	//	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	//	s3Path = s3Path + "/" + fileName //single file
	//	// 1. download file
	//	in := grpc_storage.CodeDownloadRequest{
	//		S3Path: s3Path,
	//	}
	//	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
	//		return "", err
	//	}
	//
	//	return path.Join(out.HostPath, out.FileName), nil
	//default:
	//	in := grpc_storage.CodeDownloadByDirRequest{
	//		S3Path: s3Path + "/" + fileName,
	//	}
	//	outs, err := storageClient.Client().DownloadCodeByDir(context.TODO(), &in)
	//	if err != nil {
	//		logger.Logger().Errorf("fail to download code by dir, error: %s", err.Error())
	//		return "", err
	//	}
	//	if outs != nil && len(outs.RespList) > 0 {
	//		for _, resp := range outs.RespList {
	//			logger.Logger().Debugf("download code by dir, hostpath: %s, filename: %s\n", resp.HostPath, resp.FileName)
	//			bytes, err := ioutil.ReadFile(resp.HostPath + "/" + resp.FileName)
	//			if err != nil {
	//				logger.Logger().Errorf("fail to read file, filename: %s", resp.HostPath+"/"+resp.FileName)
	//				return "", err
	//			}
	//
	//			// write to tmp file
	//			logger.Logger().Debugf("read file %s, content-length: %d\n", resp.HostPath+"/"+resp.FileName, len(bytes))
	//			subPath := strings.Split(resp.FileName, "/")
	//			filename := path.Join(subPath[1:]...)
	//			tmpFile := path.Join(tmpDir, filename)
	//			if err := createDir(tmpFile); err != nil {
	//				logger.Logger().Errorf("fail to create tmp directory of file %q, error: %s", tmpFile, err.Error())
	//				return "", err
	//			}
	//			err = ioutil.WriteFile(tmpFile, bytes, os.ModePerm)
	//			// logger.Logger().Debugf("file %q content: %s\n", resp.FileName, string(bytes))
	//			if err != nil {
	//				logger.Logger().Errorf("fail to write data to file %q, error: %s", tmpFile, err.Error())
	//				return "", err
	//			}
	//		}
	//		for _, resp := range outs.RespList {
	//			err = os.RemoveAll(resp.HostPath)
	//			if err != nil {
	//				logger.Logger().Errorf("clear local folder failed from minio, dir name: %s, error: %s", resp.HostPath, err.Error())
	//				return "", err
	//			}
	//		}
	//
	//		// gzip file
	//		dest := ""
	//		if fileName != "" {
	//			filenameOnly := getFileNameOnly(fileName)
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", filenameOnly))
	//		} else {
	//			dest = path.Join(tmpDir, fmt.Sprintf("%s.zip", s3Array[len(s3Array)-1]))
	//		}
	//		source := []string{tmpDir}
	//
	//		if err := archiver.Archive(source, dest); err != nil {
	//			logger.Logger().Errorf("fail to archive file, error: %s, source: %s, dest: %s\n", err.Error(), source, dest)
	//			return "", err
	//		}
	//
	//		return dest, nil
	//	}
	//}

	logger.Logger().Debugf("lk-test s3Path in bucket mlss-mf: %v", s3Path)
	s3Path = s3Path + "/" + fileName //single file
	// 1. download file
	in := grpc_storage.CodeDownloadRequest{
		S3Path: s3Path,
	}
	out, err := storageClient.Client().DownloadCode(context.TODO(), &in)
	if err != nil {
		logger.Logger().Errorf("fail to download code, error: %s", err.Error())
		return "", err
	}

	return path.Join(out.HostPath, out.FileName), nil

	return "", nil
}

type RmbReq struct {
	Sender  string `json:"sender"`
	Topic   string `json:"topic"`
	MsgName string `json:"msgName"`
	MsgBody string `json:"msgBody"`
}

func generateRmbModelReq(params *models.PushModelRequest, fileName, fileId, hash string) (string, error) {
	rmb := RmbModelReqBody{
		ServiceType: "uploadFile",
		Request: ModelReq{
			FactoryName: *params.FactoryName,
			FileType:    "MODEL",
			FileName:    fileName,
			FileId:      fileId,
			HashValue:   hash,
			ModelType:   *params.ModelType,
			ModelUsage:  *params.ModelUsage,
		},
	}

	by, err := json.Marshal(&rmb)
	return string(by), err
}

type RmbModelReqBody struct {
	ServiceType string   `json:"serviceType"`
	Request     ModelReq `json:"request"`
}
type ModelReq struct {
	FactoryName string `json:"factoryName"`
	FileType    string `json:"fileType"`
	FileName    string `json:"fileName"`
	FileId      string `json:"fileId"`
	HashValue   string `json:"hashValue"`
	ModelType   string `json:"modelType"` // if fileType is MODEL, can't be empty
	ModelUsage  string `json:"modelUsage"`
}

type RmbReportReqBody struct {
	ServiceType string    `json:"serviceType"`
	Request     ReportReq `json:"request"`
}
type ReportReq struct {
	FactoryName string `json:"factoryName"`
	FileType    string `json:"fileType"`
	FileName    string `json:"fileName"`
	FileId      string `json:"fileId"`
	HashValue   string `json:"hashValue"`
}

func generateRmbReportReq(params *models.PushReportRequest, fileName, fileId, hash string) (string, error) {
	rmb := RmbReportReqBody{
		ServiceType: "uploadFile",
		Request: ReportReq{
			FactoryName: *params.FactoryName,
			FileType:    "DATA",
			FileName:    fileName,
			FileId:      fileId,
			HashValue:   hash,
		},
	}
	bytes, err := json.Marshal(&rmb)
	if err != nil {
		return "", err
	}
	bodyStr := strings.ReplaceAll(string(bytes), `\`, "")

	logger.Logger().Debugf("generateRmbReportReq bodyStr: %s\n", bodyStr)
	return bodyStr, err
}

func GetPushEventByID(id int64) (*models.Event, error) {
	eventInDb, err := dao.GetEventById(id)
	var event models.Event
	logger.Logger().Infof("GetPushEventByID report: %+v\n", eventInDb)
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			logger.Logger().Errorf("fail to get report by model name and model version,  error: %s", err.Error())
			return nil, err
		}
	}
	err = copier.Copy(&event, &eventInDb)
	if err != nil {
		logger.Logger().Errorf("copier event model err,  error: %s", err.Error())
		return nil, err
	}
	return &event, nil

}
