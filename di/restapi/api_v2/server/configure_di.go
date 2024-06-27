// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"
	"webank/DI/commons/service"
	"webank/DI/restapi/api_v2/server/operations/experiment"
	"webank/DI/restapi/api_v2/server/operations/experiment_run"
	"webank/DI/restapi/api_v2/server/rest_impl"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	mw "webank/DI/restapi/middleware"

	"webank/DI/restapi/api_v2/server/operations"
)

//go:generate swagger generate server --target ..\..\api_v2 --name Di --spec ..\swagger\swagger.yml --model-package restmodels --server-package server --exclude-main

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

	api.JSONProducer = runtime.JSONProducer()

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

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	rest_impl.UpdateExperimentRunExecuteStatus()

	api.ExperimentGetDSSJumpMLSSUserHandler = experiment.GetDSSJumpMLSSUserHandlerFunc(func(params experiment.GetDSSJumpMLSSUserParams, i interface{}) middleware.Responder {
		return rest_impl.GetDSSJumpMLSSUser(params)
	})

	api.ExperimentListExperimentVersionNamesHandler = experiment.ListExperimentVersionNamesHandlerFunc(func(params experiment.ListExperimentVersionNamesParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperimentVersionNames(params)
	})

	api.ExperimentListExperimentSourceSystemsHandler = experiment.ListExperimentSourceSystemsHandlerFunc(func(params experiment.ListExperimentSourceSystemsParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperimentSourceSystems(params)
	})

	api.ExperimentCreateExperimentHandler = experiment.CreateExperimentHandlerFunc(func(params experiment.CreateExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperiment(params)
	})

	api.ExperimentCreateExperimentByUploadHandler = experiment.CreateExperimentByUploadHandlerFunc(func(params experiment.CreateExperimentByUploadParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperimentByUpload(params)
	})

	api.ExperimentListExperimentsHandler = experiment.ListExperimentsHandlerFunc(func(params experiment.ListExperimentsParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperiments(params)
	})

	api.ExperimentListExperimentDSSProjectNamesHandler = experiment.ListExperimentDSSProjectNamesHandlerFunc(func(params experiment.ListExperimentDSSProjectNamesParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperimentDssProjectNames(params)
	})

	api.ExperimentListExperimentDSSFlowNamesHandler = experiment.ListExperimentDSSFlowNamesHandlerFunc(func(params experiment.ListExperimentDSSFlowNamesParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperimentDssFlowNames(params)
	})

	api.ExperimentPatchExperimentVersionHandler = experiment.PatchExperimentVersionHandlerFunc(func(params experiment.PatchExperimentVersionParams, i interface{}) middleware.Responder {
		return rest_impl.PatchExperimentVersion(params)
	})

	api.ExperimentCreateExperimentVersionHandler = experiment.CreateExperimentVersionHandlerFunc(func(params experiment.CreateExperimentVersionParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperimentVersion(params)
	})

	api.ExperimentGetExperimentVersionFlowJSONHandler = experiment.GetExperimentVersionFlowJSONHandlerFunc(func(params experiment.GetExperimentVersionFlowJSONParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentVersionFlowJson(params)
	})

	api.ExperimentGetExperimentVersionGlobalVariablesStrHandler = experiment.GetExperimentVersionGlobalVariablesStrHandlerFunc(func(params experiment.GetExperimentVersionGlobalVariablesStrParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentVersionGlobalVariablesStr(params)
	})

	api.ExperimentGetExperimentVersionHandler = experiment.GetExperimentVersionHandlerFunc(func(params experiment.GetExperimentVersionParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperimentVersion(params)
	})

	api.ExperimentListExperimentVersionsHandler = experiment.ListExperimentVersionsHandlerFunc(func(params experiment.ListExperimentVersionsParams, principal interface{}) middleware.Responder {
		return rest_impl.ListExperimentVersions(params)
	})

	api.ExperimentPatchExperimentHandler = experiment.PatchExperimentHandlerFunc(func(params experiment.PatchExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.PatchExperiment(params)
	})

	api.ExperimentCopyExperimentHandler = experiment.CopyExperimentHandlerFunc(func(params experiment.CopyExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.CopyExperiment(params)
	})

	api.ExperimentExportExperimentFlowJSONHandler = experiment.ExportExperimentFlowJSONHandlerFunc(func(params experiment.ExportExperimentFlowJSONParams, principal interface{}) middleware.Responder {
		return rest_impl.ExportExperimentFlowJson(params)
	})

	api.ExperimentCodeUploadHandler = experiment.CodeUploadHandlerFunc(func(params experiment.CodeUploadParams, i interface{}) middleware.Responder {
		return rest_impl.CodeUpload(params)
	})

	api.ExperimentUploadExperimentFlowJSONHandler = experiment.UploadExperimentFlowJSONHandlerFunc(func(params experiment.UploadExperimentFlowJSONParams, i interface{}) middleware.Responder {
		return rest_impl.UploadExperimentFlowJson(params)
	})

	api.ExperimentRunGetExperimentRunNodeExecutionsHandler = experiment_run.GetExperimentRunNodeExecutionsHandlerFunc(func(params experiment_run.GetExperimentRunNodeExecutionsParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunNodeExecutions(params)
	})

	api.ExperimentRunGetExperimentRunNodeLogsHandler = experiment_run.GetExperimentRunNodeLogsHandlerFunc(func(params experiment_run.GetExperimentRunNodeLogsParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunNodeLogs(params)
	})

	api.ExperimentRunDeleteExperimentRunHandler = experiment_run.DeleteExperimentRunHandlerFunc(func(params experiment_run.DeleteExperimentRunParams, i interface{}) middleware.Responder {
		return rest_impl.DeleteExperimentRun(params)
	})

	api.ExperimentRunGetExperimentRunHandler = experiment_run.GetExperimentRunHandlerFunc(func(params experiment_run.GetExperimentRunParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRun(params)
	})

	api.ExperimentRunGetExperimentRunStatusHandler = experiment_run.GetExperimentRunStatusHandlerFunc(func(params experiment_run.GetExperimentRunStatusParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunStatus(params)
	})

	api.ExperimentRunGetExperimentRunFlowJSONHandler = experiment_run.GetExperimentRunFlowJSONHandlerFunc(func(params experiment_run.GetExperimentRunFlowJSONParams, i interface{}) middleware.Responder {
		return rest_impl.GetExperimentRunFlowJson(params)
	})

	api.ExperimentDeleteExperimentHandler = experiment.DeleteExperimentHandlerFunc(func(params experiment.DeleteExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.DeleteExperiment(params)
	})

	api.ExperimentGetExperimentHandler = experiment.GetExperimentHandlerFunc(func(params experiment.GetExperimentParams, principal interface{}) middleware.Responder {
		return rest_impl.GetExperiment(params)
	})

	api.ExperimentRunCreateExperimentRunHandler = experiment_run.CreateExperimentRunHandlerFunc(func(params experiment_run.CreateExperimentRunParams, principal interface{}) middleware.Responder {
		return rest_impl.CreateExperimentRun(params)
	})

	api.ExperimentRunListExperimentRunsHandler = experiment_run.ListExperimentRunsHandlerFunc(func(params experiment_run.ListExperimentRunsParams, i interface{}) middleware.Responder {
		return rest_impl.ListExperimentRuns(params)
	})

	api.ExperimentRunKillExperimentRunHandler = experiment_run.KillExperimentRunHandlerFunc(func(params experiment_run.KillExperimentRunParams, i interface{}) middleware.Responder {
		return rest_impl.KillExperimentRun(params)
	})

	api.ExperimentRunRetryExperimentRunHandler = experiment_run.RetryExperimentRunHandlerFunc(func(params experiment_run.RetryExperimentRunParams, i interface{}) middleware.Responder {
		return rest_impl.RetryExperimentRun(params)
	})

	api.ExperimentDeleteExperimentVersionHandler = experiment.DeleteExperimentVersionHandlerFunc(func(params experiment.DeleteExperimentVersionParams, i interface{}) middleware.Responder {
		return rest_impl.DeleteExperimentVersion(params)
	})

	api.ExperimentBatchDeleteExperimentsHandler = experiment.BatchDeleteExperimentsHandlerFunc(func(params experiment.BatchDeleteExperimentsParams, i interface{}) middleware.Responder {
		return rest_impl.BatchDeleteExperiments(params)
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
	return logging.Handle(handler)
}
