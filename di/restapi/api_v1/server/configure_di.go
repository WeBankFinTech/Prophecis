// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"
	"webank/DI/commons/service"
	"webank/DI/restapi/api_v1/server/rest_impl"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"webank/DI/restapi/api_v1/server/operations"
	"webank/DI/restapi/api_v1/server/operations/dss_user_info"
	"webank/DI/restapi/api_v1/server/operations/experiment_runs"
	"webank/DI/restapi/api_v1/server/operations/experiments"
	"webank/DI/restapi/api_v1/server/operations/linkis_job"
	"webank/DI/restapi/api_v1/server/operations/models"
	"webank/DI/restapi/api_v1/server/operations/training_data"
)

//go:generate swagger generate server --target ..\..\api_v1 --name Di --spec ..\swagger\swagger.yml --model-package restmodels --server-package server --exclude-main

func configureFlags(api *operations.DiAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DiAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.BinProducer = runtime.ByteStreamProducer()

	// Important: we don't really autenticate here anymore. It is done at the gateway level.
	api.WatsonAuthTokenAuth = func(token string) (interface{}, error) {
		if token == "" {
			return nil, errors.Unauthenticated("token")
		}
		return &service.User{}, nil
	}

	api.WatsonAuthTokenQueryAuth = func(token string) (interface{}, error) {
		return api.WatsonAuthTokenAuth(token)
	}

	// Important: we don't really autenticate here anymore. It is done at the gateway level.
	api.BasicAuthTokenAuth = func(token string) (interface{}, error) {
		return api.WatsonAuthTokenAuth(token)
	}
	api.ModelsDeleteModelHandler = models.DeleteModelHandlerFunc(func(params models.DeleteModelParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteModel(params)
	})
	api.ModelsDownloadModelDefinitionHandler = models.DownloadModelDefinitionHandlerFunc(func(params models.DownloadModelDefinitionParams, principal interface{}) middleware.Responder {
		return rest_impl.DownloadModelDefinition(params)
	})
	api.ModelsDownloadTrainedModelHandler = models.DownloadTrainedModelHandlerFunc(func(params models.DownloadTrainedModelParams, principal interface{}) middleware.Responder {
		return rest_impl.DownloadTrainedModel(params)
	})
	api.ModelsExportModelHandler = models.ExportModelHandlerFunc(func(params models.ExportModelParams, i interface{}) middleware.Responder {
		return rest_impl.ExportModel(params)
	})
	api.ModelsGetLogsHandler = models.GetLogsHandlerFunc(func(params models.GetLogsParams) middleware.Responder {
		return rest_impl.GetLogs(params)
	})
	api.ModelsGetModelHandler = models.GetModelHandlerFunc(func(params models.GetModelParams, principal interface{}) middleware.Responder {
		return rest_impl.GetModel(params)
	})
	api.ModelsRetryModelHandler = models.RetryModelHandlerFunc(func(params models.RetryModelParams, principal interface{}) middleware.Responder {
		return rest_impl.RetryModel(params)
	})
	api.ModelsListModelsHandler = models.ListModelsHandlerFunc(func(params models.ListModelsParams, principal interface{}) middleware.Responder {
		return rest_impl.ListModels(params)
	})
	api.ModelsPostModelHandler = models.PostModelHandlerFunc(func(params models.PostModelParams, principal interface{}) middleware.Responder {
		return rest_impl.PostModel(params)
	})
	api.ModelsPatchModelHandler = models.PatchModelHandlerFunc(func(params models.PatchModelParams, principal interface{}) middleware.Responder {
		return rest_impl.PatchModel(params)
	})
	api.TrainingDataGetLoglinesHandler = training_data.GetLoglinesHandlerFunc(func(params training_data.GetLoglinesParams, principal interface{}) middleware.Responder {
		return rest_impl.GetLoglines(params)
	})
	api.GetDashboardsHandler = operations.GetDashboardsHandlerFunc(func(params operations.GetDashboardsParams, principal interface{}) middleware.Responder {
		return rest_impl.GetDashboards(params)
	})
	api.ExperimentsDeleteExperimentHandler = experiments.DeleteExperimentHandlerFunc(func(params experiments.DeleteExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteExperiment(params)
	})
	api.ExperimentRunsDeleteExperimentRunHandler = experiment_runs.DeleteExperimentRunHandlerFunc(func(params experiment_runs.DeleteExperimentRunParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteExperimentRun(params)
	})
	api.ExperimentsDeleteExperimentTagHandler = experiments.DeleteExperimentTagHandlerFunc(func(params experiments.DeleteExperimentTagParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteExperimentTag(params)
	})
	api.ExperimentsExportExperimentHandler = experiments.ExportExperimentHandlerFunc(func(params experiments.ExportExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.ExportExperiment(params)
	})
	api.ExperimentsGetExperimentHandler = experiments.GetExperimentHandlerFunc(func(params experiments.GetExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperiment(params)
	})
	api.ExperimentRunsGetExperimentRunHandler = experiment_runs.GetExperimentRunHandlerFunc(func(params experiment_runs.GetExperimentRunParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperimentRun(params)
	})
	api.ExperimentsImportExperimentHandler = experiments.ImportExperimentHandlerFunc(func(params experiments.ImportExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.ImportExperiment(params)
	})
	api.ExperimentRunsListExperimentRunsHandler = experiment_runs.ListExperimentRunsHandlerFunc(func(params experiment_runs.ListExperimentRunsParams, principal interface{}) middleware.Responder {
		return rest_impl.ListExperimentRuns(params)
	})
	api.ExperimentsListExperimentsHandler = experiments.ListExperimentsHandlerFunc(func(params experiments.ListExperimentsParams, principal interface{}) middleware.Responder {
		return rest_impl.ListExperiments(params)
	})
	api.ExperimentsUpdateExperimentHandler = experiments.UpdateExperimentHandlerFunc(func(params experiments.UpdateExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.UpdateExperiment(params)
	})
	api.ExperimentsCreateExperimentHandler = experiments.CreateExperimentHandlerFunc(func(params experiments.CreateExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperiment(params)
	})
	api.ExperimentsCopyExperimentHandler = experiments.CopyExperimentHandlerFunc(func(params experiments.CopyExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.CopyExperiment(params)
	})
	api.ExperimentRunsCreateExperimentRunHandler = experiment_runs.CreateExperimentRunHandlerFunc(func(params experiment_runs.CreateExperimentRunParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperimentRun(params)
	})
	api.ExperimentRunsGetExperimentRunStatusHandler = experiment_runs.GetExperimentRunStatusHandlerFunc(func(params experiment_runs.GetExperimentRunStatusParams, principal interface{}) middleware.Responder {
		return rest_impl.StatusExperimentRun(params)
	})
	api.ExperimentRunsKillExperimentRunHandler = experiment_runs.KillExperimentRunHandlerFunc(func(params experiment_runs.KillExperimentRunParams, principal interface{}) middleware.Responder {
		return rest_impl.KillExperimentRun(params)
	})
	api.ModelsGetJobLogByLineHandler = models.GetJobLogByLineHandlerFunc(func(params models.GetJobLogByLineParams, principal interface{}) middleware.Responder {
		return rest_impl.GetJobLogByLine(params)
	})
	api.ExperimentRunsGetExperimentRunLogHandler = experiment_runs.GetExperimentRunLogHandlerFunc(func(params experiment_runs.GetExperimentRunLogParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunLog(params)
	})
	api.ExperimentsCodeUploadHandler = experiments.CodeUploadHandlerFunc(func(params experiments.CodeUploadParams, principal interface{}) middleware.Responder {
		return rest_impl.CodeUpload(params)
	})
	api.ExperimentsUpdateExperimentInfoHandler = experiments.UpdateExperimentInfoHandlerFunc(func(params experiments.UpdateExperimentInfoParams, principal interface{}) middleware.Responder {
		return rest_impl.UpdateExperimentInfo(params)
	})
	api.LinkisJobGetLinkisJobStatusHandler = linkis_job.GetLinkisJobStatusHandlerFunc(func(params linkis_job.GetLinkisJobStatusParams, principal interface{}) middleware.Responder {
		return rest_impl.GetLinkisJobStatus(params)
	})
	api.LinkisJobGetLinkisJobLogHandler = linkis_job.GetLinkisJobLogHandlerFunc(func(params linkis_job.GetLinkisJobLogParams, principal interface{}) middleware.Responder {
		return rest_impl.GetLinkisJobLog(params)
	})
	api.ExperimentRunsGetExperimentRunExecutionHandler = experiment_runs.GetExperimentRunExecutionHandlerFunc(func(params experiment_runs.GetExperimentRunExecutionParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunExecution(params)
	})
	api.ExperimentsCreateExperimentTagHandler = experiments.CreateExperimentTagHandlerFunc(func(params experiments.CreateExperimentTagParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperimentTag(params)
	})
	api.ExperimentsDeleteExperimentTagHandler = experiments.DeleteExperimentTagHandlerFunc(func(params experiments.DeleteExperimentTagParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteExperimentTag(params)
	})
	api.ExperimentsExportExperimentDssHandler = experiments.ExportExperimentDssHandlerFunc(func(params experiments.ExportExperimentDssParams, principal interface{}) middleware.Responder {
		return rest_impl.ExportExperimentDss(params)
	})
	api.ExperimentsImportExperimentDssHandler = experiments.ImportExperimentDssHandlerFunc(func(params experiments.ImportExperimentDssParams, principal interface{}) middleware.Responder {
		return rest_impl.ImportExperimentDss(params)
	})
	api.DssUserInfoGetDssUserInfoHandler = dss_user_info.GetDssUserInfoHandlerFunc(func(params dss_user_info.GetDssUserInfoParams, i interface{}) middleware.Responder {
		return rest_impl.DssGetUserInfo(params)
	})
	api.ExperimentRunsGetExperimentRunsHistoryHandler = experiment_runs.GetExperimentRunsHistoryHandlerFunc(func(params experiment_runs.GetExperimentRunsHistoryParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunsHistory(params)
	})
	api.ModelsKillTrainingModelHandler = models.KillTrainingModelHandlerFunc(func(params models.KillTrainingModelParams, i interface{}) middleware.Responder {
		return rest_impl.KillTrainingModel(params)
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
	return handler
}
