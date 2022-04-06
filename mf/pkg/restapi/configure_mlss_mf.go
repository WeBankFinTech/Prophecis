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

// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	mw "mlss-mf/pkg/middleware"
	"mlss-mf/pkg/restapi/operations"
	"mlss-mf/pkg/restapi/operations/container"
	"mlss-mf/pkg/restapi/operations/image"
	"mlss-mf/pkg/restapi/operations/model_deploy"
	"mlss-mf/pkg/restapi/operations/model_result"
	"mlss-mf/pkg/restapi/operations/model_storage"
	"mlss-mf/pkg/restapi/operations/report"
	"mlss-mf/pkg/restapi/operations/rmb"
	"net/http"
	"strings"

	"github.com/dre1080/recover"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ..\..\pkg --name MlssMf --spec ..\swagger\swagger.yaml

func configureFlags(api *operations.MlssMfAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MlssMfAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf
	api.JSONConsumer = runtime.JSONConsumer()

	api.ModelDeployListServicesHandler = model_deploy.ListServicesHandlerFunc(func(params model_deploy.ListServicesParams) middleware.Responder {
		return ServiceList(params)
	})
	api.ModelDeployGetServiceHandler = model_deploy.GetServiceHandlerFunc(func(params model_deploy.GetServiceParams) middleware.Responder {
		return ServiceGet(params)
	})
	api.ModelDeployDeleteServiceHandler = model_deploy.DeleteServiceHandlerFunc(func(params model_deploy.DeleteServiceParams) middleware.Responder {
		return ServiceDelete(params)
	})
	api.ModelDeployUpdateServiceHandler = model_deploy.UpdateServiceHandlerFunc(func(params model_deploy.UpdateServiceParams) middleware.Responder {
		return ServiceUpdate(params)
	})
	api.ModelDeployStopNamespacedServiceHandler = model_deploy.StopNamespacedServiceHandlerFunc(func(params model_deploy.StopNamespacedServiceParams) middleware.Responder {
		return ServiceStop(params)
	})
	api.ModelDeployCreateNamespacedServiceRunHandler = model_deploy.CreateNamespacedServiceRunHandlerFunc(func(params model_deploy.CreateNamespacedServiceRunParams) middleware.Responder {
		return ServiceRun(params)
	})
	api.ModelDeployPostServiceHandler = model_deploy.PostServiceHandlerFunc(func(params model_deploy.PostServiceParams) middleware.Responder {
		return ServicePost(params)
	})
	api.ModelDeployServiceDashboradHandler = model_deploy.ServiceDashboradHandlerFunc(func(params model_deploy.ServiceDashboradParams) middleware.Responder {
		return ServiceDashborad(params)
	})
	api.ModelStorageGetModelHandler = model_storage.GetModelHandlerFunc(func(params model_storage.GetModelParams) middleware.Responder {
		return GetModel(params)
	})
	api.ModelStoragePostModelHandler = model_storage.PostModelHandlerFunc(func(params model_storage.PostModelParams) middleware.Responder {
		return PostModel(params)
	})
	api.ModelStorageDeleteModelHandler = model_storage.DeleteModelHandlerFunc(func(params model_storage.DeleteModelParams) middleware.Responder {
		return DeleteModel(params)
	})
	api.ModelStorageUpdateModelHandler = model_storage.UpdateModelHandlerFunc(func(params model_storage.UpdateModelParams) middleware.Responder {
		return UpdateModel(params)
	})
	api.ModelStorageGetModelsHandler = model_storage.GetModelsHandlerFunc(func(params model_storage.GetModelsParams) middleware.Responder {
		return ListModels(params)
	})
	api.ModelStorageGetModelsByClusterHandler = model_storage.GetModelsByClusterHandlerFunc(func(params model_storage.GetModelsByClusterParams) middleware.Responder {
		return ListModelsByCluster(params)
	})
	api.ModelStorageListModelsByGroupIDHandler = model_storage.ListModelsByGroupIDHandlerFunc(func(params model_storage.ListModelsByGroupIDParams) middleware.Responder {
		return ListModelsByGroup(params)
	})
	api.ModelStorageUploadModelHandler = model_storage.UploadModelHandlerFunc(func(params model_storage.UploadModelParams) middleware.Responder {
		return UploadModel(params)
	})
	api.ModelStorageGetModelVersionHandler = model_storage.GetModelVersionHandlerFunc(func(params model_storage.GetModelVersionParams) middleware.Responder {
		return ListModelVersions(params)
	})
	api.ModelStorageExportModelHandler = model_storage.ExportModelHandlerFunc(func(params model_storage.ExportModelParams) middleware.Responder {
		return ExportModel(params)
	})
	api.ModelStorageGetModelsByGroupIDAndModelNameHandler = model_storage.GetModelsByGroupIDAndModelNameHandlerFunc(
		func(params model_storage.GetModelsByGroupIDAndModelNameParams) middleware.Responder {
			return GetModelsByGroupIdAndModelName(params)
		})
	api.ImageCreateImageHandler = image.CreateImageHandlerFunc(func(params image.CreateImageParams) middleware.Responder {
		return CreateImage(params)
	})
	api.ImageDeleteImageHandler = image.DeleteImageHandlerFunc(func(params image.DeleteImageParams) middleware.Responder {
		return DeleteImage(params)
	})
	api.ImageListImageHandler = image.ListImageHandlerFunc(func(params image.ListImageParams) middleware.Responder {
		return ListImage(params)
	})
	api.ImageListImageByModelVersionIDHandler = image.ListImageByModelVersionIDHandlerFunc(func(params image.ListImageByModelVersionIDParams) middleware.Responder {
		return ListImageByModelVersionId(params)
	})
	api.ImageGetImageHandler = image.GetImageHandlerFunc(func(params image.GetImageParams) middleware.Responder {
		return GetImage(params)
	})
	api.ImageUpdateImageHandler = image.UpdateImageHandlerFunc(func(params image.UpdateImageParams) middleware.Responder {
		return UpdateImage(params)
	})
	api.ContainerListContainerHandler = container.ListContainerHandlerFunc(func(params container.ListContainerParams) middleware.Responder {
		return ListContainer(params)
	})
	api.ContainerGetNamespacedServiceContainerLogHandler = container.GetNamespacedServiceContainerLogHandlerFunc(
		func(params container.GetNamespacedServiceContainerLogParams) middleware.Responder {
			return GetNamespacedServiceContainerLog(params)
		})
	api.ReportGetReportByModelNameAndModelVersionHandler = report.GetReportByModelNameAndModelVersionHandlerFunc(
		func(params report.GetReportByModelNameAndModelVersionParams) middleware.Responder {
			return GetReportByModelNameAndModelVersion(params)
		})
	api.ReportGetReportByIDHandler = report.GetReportByIDHandlerFunc(
		func(params report.GetReportByIDParams) middleware.Responder {
			return GetReportByID(params)
		})
	api.ReportCreateReportHandler = report.CreateReportHandlerFunc(
		func(params report.CreateReportParams) middleware.Responder {
			return CreateReport(params)
		})
	api.ReportListReportsHandler = report.ListReportsHandlerFunc(
		func(params report.ListReportsParams) middleware.Responder {
			return ListReports(params)
		})
	api.ReportDeleteReportByIDHandler = report.DeleteReportByIDHandlerFunc(
		func(params report.DeleteReportByIDParams) middleware.Responder {
			return DeleteReportByID(params)
		})
	//api.ReportPushReportByReportIDHandler = report.PushReportByReportIDHandlerFunc(
	//	func(params report.PushReportByReportIDParams) middleware.Responder {
	//		return PushReportByReportID(params)
	//	})
	//api.ReportPushReportByReportVersionIDHandler = report.PushReportByReportVersionIDHandlerFunc(
	//	func(params report.PushReportByReportVersionIDParams) middleware.Responder {
	//		return PushReportByReportVersionID(params)
	//	})
	api.ReportListReportVersionPushEventsByReportVersionIDHandler = report.ListReportVersionPushEventsByReportVersionIDHandlerFunc(
		func(lrpebri report.ListReportVersionPushEventsByReportVersionIDParams) middleware.Responder {
			return ListReportVersionPushEventsByReportVersionID(lrpebri)
		})
	api.ReportDownloadReportByIDHandler = report.DownloadReportByIDHandlerFunc(
		func(params report.DownloadReportByIDParams) middleware.Responder {
			return DownloadReportByID(params)
		})
	api.ReportUploadReportHandler = report.UploadReportHandlerFunc(
		func(urp report.UploadReportParams) middleware.Responder {
			return UploadReport(urp)
		})
	api.ReportListReportVersionsByReportIDHandler = report.ListReportVersionsByReportIDHandlerFunc(
		func(lrvbri report.ListReportVersionsByReportIDParams) middleware.Responder {
			return ListReportVersionsByReportID(lrvbri)
		})
	api.ModelResultGetResultByModelNameHandler = model_result.GetResultByModelNameHandlerFunc(
		func(params model_result.GetResultByModelNameParams) middleware.Responder {
			return GetResultByModelName(params)
		})
	api.ModelResultDeleteResultByIDHandler = model_result.DeleteResultByIDHandlerFunc(
		func(params model_result.DeleteResultByIDParams) middleware.Responder {
			return DeleteResultByID(params)
		})
	api.ModelResultCreateResultHandler = model_result.CreateResultHandlerFunc(
		func(params model_result.CreateResultParams) middleware.Responder {
			return CreateResult(params)
		})
	api.ModelResultUpdateResultByIDHandler = model_result.UpdateResultByIDHandlerFunc(
		func(params model_result.UpdateResultByIDParams) middleware.Responder {
			return UpdateResultByID(params)
		})
	api.ModelStorageDownloadModelByIDHandler = model_storage.DownloadModelByIDHandlerFunc(
		func(params model_storage.DownloadModelByIDParams) middleware.Responder {
			return DownloadModelByID(params)
		})
	api.ModelStorageDownloadModelVersionByIDHandler = model_storage.DownloadModelVersionByIDHandlerFunc(
		func(params model_storage.DownloadModelVersionByIDParams) middleware.Responder {
			return DownloadModelVersionByID(params)
		})
	api.ModelStoragePushModelByModelIDHandler = model_storage.PushModelByModelIDHandlerFunc(
		func(pmbmi model_storage.PushModelByModelIDParams) middleware.Responder {
			return PushModelByModelId(pmbmi)
		})
	api.ModelStoragePushModelByModelVersionIDHandler = model_storage.PushModelByModelVersionIDHandlerFunc(
		func(pmbmi model_storage.PushModelByModelVersionIDParams) middleware.Responder {
			return PushModelByModelVersionId(pmbmi)
		})

	api.ModelStorageListModelVersionPushEventsByModelVersionIDHandler = model_storage.ListModelVersionPushEventsByModelVersionIDHandlerFunc(
		func(lmpebmi model_storage.ListModelVersionPushEventsByModelVersionIDParams) middleware.Responder {
			return ListModelVersionPushEventsByModelVersionID(lmpebmi)
		})
	api.RmbDownloadRmbLogByEventIDHandler = rmb.DownloadRmbLogByEventIDHandlerFunc(
		func(params rmb.DownloadRmbLogByEventIDParams) middleware.Responder {
			return DownloadRmbLogByEventID(params)
		})
	api.ModelStorageGetModelVersionByNameAndVersionHandler = model_storage.GetModelVersionByNameAndVersionHandlerFunc(
		func(params model_storage.GetModelVersionByNameAndVersionParams) middleware.Responder {
			return GetModelVersionByNameAndGroupID(params)
		})
	api.ReportGetPushEventByIDHandler = report.GetPushEventByIDHandlerFunc(
		func(params report.GetPushEventByIDParams) middleware.Responder {
			return GetPushEventById(params)
		})
	api.ServerShutdown = func() {}
	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	logging := mw.NewLoggingMiddleware("rest-api")
	recovery := recover.New(&recover.Options{
		Log: log.Print,
	})

	auth := mw.NewAuthMiddleware(&mw.AuthOptions{
		ExcludedURLs: []string{},
	})

	return logging.Handle(recovery(auth(uiMiddleware(handler))))
	//return uiMiddleware(handler)
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
