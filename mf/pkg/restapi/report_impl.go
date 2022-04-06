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

package restapi

import (
	"encoding/json"
	"mlss-mf/pkg/common"
	cc "mlss-mf/pkg/common/controlcenter"
	"mlss-mf/pkg/restapi/operations/report"
	"mlss-mf/pkg/service"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func GetReportByModelNameAndModelVersion(params report.GetReportByModelNameAndModelVersionParams) middleware.Responder {
	return returnResponse(http.StatusOK, nil, "GetReportByModelNameAndModelVersion", []byte("not implement"))
}

func GetReportByID(params report.GetReportByIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	resp, err := service.GetReportByID(params.ReportID, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "GetReportByID", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "GetReportByID", nil)
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetReportByID", nil)
	}
	return returnResponse(http.StatusOK, nil, "GetReportByID", respJson)
}

func CreateReport(params report.CreateReportParams) middleware.Responder {
	resp, err := service.CreateReport(params)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "CreateReport", nil)
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "CreateReport", nil)
	}

	return returnResponse(http.StatusOK, nil, "CreateReport", respJson)
}

func ListReports(params report.ListReportsParams) middleware.Responder {
	currentUserName := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	queryStr := ""
	if params.QueryStr != nil {
		queryStr = *(params.QueryStr)
	}
	var currentPage, pageSize int64 = 0, 0
	if params.CurrentPage != nil {
		currentPage = *params.CurrentPage
	}
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	resp, err := service.ListReports(currentPage, pageSize, currentUserName, queryStr, isSA)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "ListReports", nil)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "ListReports", nil)
	}
	return returnResponse(http.StatusOK, nil, "ListReports", respJson)
}

func ListReportVersionsByReportID(params report.ListReportVersionsByReportIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	reportVersions, err := service.ListReportVersionsByReportID(params, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "ListReportVersionsByReportID", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "ListReportVersionsByReportID", nil)
	}
	respJson, err := json.Marshal(reportVersions)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "ListReportVersionsByReportID", nil)
	}
	return returnResponse(http.StatusOK, nil, "ListReportVersionsByReportID", respJson)
}

func DeleteReportByID(params report.DeleteReportByIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	err := service.DeleteReportByID(params.ReportID, user, isSA)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "DeleteReportByID", nil)
	}

	err = service.DeleteReportVersionByReportID(params.ReportID, user, isSA)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "DeleteReportByID", nil)
	}

	return returnResponse(http.StatusOK, nil, "DeleteReportByID", nil)
}

func DownloadReportByID(params report.DownloadReportByIDParams) middleware.Responder {
	currentUserID := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := false
	if params.HTTPRequest.Header.Get(cc.CcSuperadmin) == "true" {
		isSA = true
	}
	bytes, statusCode, fileName, err := service.DownloadReportById(params.ReportID, currentUserID, isSA)
	if err != nil {
		return returnResponse(statusCode, err, "DownloadReportByID", nil)
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.Header().Add("File-Name", fileName)
		w.WriteHeader(statusCode)
		w.Write(bytes)
	})
}

func UploadReport(params report.UploadReportParams) middleware.Responder {
	resp, statusCode, err := service.UploadReport(params.FileName, params.File)
	if err != nil {
		return returnResponse(statusCode, err, "UploadReport", nil)
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "UploadReport", nil)
	}
	return returnResponse(statusCode, nil, "UploadReport", respJson)
}

func ListReportVersionPushEventsByReportVersionID(params report.ListReportVersionPushEventsByReportVersionIDParams) middleware.Responder {
	user := params.HTTPRequest.Header.Get(common.UserIDHeader)
	isSA := isSuperAdmin(params.HTTPRequest)
	events, err := service.ListReportVersionPushEventsByReportVersionID(params.ReportVersionID, *params.CurrentPage, *params.PageSize, user, isSA)
	if err != nil {
		if err == service.PermissionDeniedError {
			return returnResponse(http.StatusForbidden, err, "ListReportVersionPushEventsByReportVersionID", nil)
		}
		return returnResponse(http.StatusInternalServerError, err, "ListReportVersionPushEventsByReportVersionID", nil)
	}

	respJson, err := json.Marshal(events)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "ListReportVersionPushEventsByReportVersionID", nil)
	}
	return returnResponse(http.StatusOK, nil, "ListReportVersionPushEventsByReportVersionID", respJson)
}

func GetPushEventById(params report.GetPushEventByIDParams) middleware.Responder {
	resp, err := service.GetPushEventByID(params.EventID)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetReportByID", nil)
	}
	respJson, err := json.Marshal(resp)
	if err != nil {
		return returnResponse(http.StatusInternalServerError, err, "GetReportByID", nil)
	}
	return returnResponse(http.StatusOK, nil, "GetPushEventById", respJson)
}
