// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"
	"webank/DI/restapi/api_v1/server/rest_impl"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"webank/DI/restapi/api_v1/server/operations"
	"webank/DI/restapi/api_v1/server/operations/experiments"
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
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.BinProducer = runtime.ByteStreamProducer()

	// Applies when the Authorization header is set with the Basic scheme
	//api.BasicAuthAuth = func(user string, pass string) (interface{}, error) {
	//	// We can assume that the basic authentication is handled by the frontend
	//	// gateway. We don't need any checks here anymore.
	//	if user == "" || pass == "" {
	//		return nil, errors.Unauthenticated("user/password missing")
	//	}
	//	return user, nil
	//}

	api.ExperimentsCodeUploadHandler = experiments.CodeUploadHandlerFunc(func(params experiments.CodeUploadParams) middleware.Responder {
		return rest_impl.CodeUpload(params)
	})
	api.ModelsDeleteModelHandler = models.DeleteModelHandlerFunc(func(params models.DeleteModelParams) middleware.Responder {
		return rest_impl.DeleteModel(params)
	})
	api.ModelsDownloadModelDefinitionHandler = models.DownloadModelDefinitionHandlerFunc(func(params models.DownloadModelDefinitionParams) middleware.Responder {
		return rest_impl.DownloadModelDefinition(params)
	})
	api.ModelsDownloadTrainedModelHandler = models.DownloadTrainedModelHandlerFunc(func(params models.DownloadTrainedModelParams) middleware.Responder {
		return rest_impl.DownloadTrainedModel(params)
	})
	api.ModelsExportModelHandler = models.ExportModelHandlerFunc(func(params models.ExportModelParams) middleware.Responder {
		return rest_impl.ExportModel(params)
	})
	api.ModelsGetLogsHandler = models.GetLogsHandlerFunc(func(params models.GetLogsParams) middleware.Responder {
		return rest_impl.GetLogs(params)
	})
	api.ModelsGetModelHandler = models.GetModelHandlerFunc(func(params models.GetModelParams) middleware.Responder {
		return rest_impl.GetModel(params)
	})
	api.ModelsListModelsHandler = models.ListModelsHandlerFunc(func(params models.ListModelsParams) middleware.Responder {
		return rest_impl.ListModels(params)
	})
	api.ModelsPostModelHandler = models.PostModelHandlerFunc(func(params models.PostModelParams) middleware.Responder {
		return rest_impl.PostModel(params)
	})
	api.ModelsPatchModelHandler = models.PatchModelHandlerFunc(func(params models.PatchModelParams) middleware.Responder {
		return rest_impl.PatchModel(params)
	})
	api.TrainingDataGetLoglinesHandler = training_data.GetLoglinesHandlerFunc(func(params training_data.GetLoglinesParams) middleware.Responder {
		return rest_impl.GetLoglines(params)
	})
	api.GetDashboardsHandler = operations.GetDashboardsHandlerFunc(func(params operations.GetDashboardsParams) middleware.Responder {
		return rest_impl.GetDashboards(params)
	})
	api.ServerShutdown = func() {}
	// return setupGlobalMiddleware(api.Serve(setupMiddlewares))
	return api.Serve(setupMiddlewares)
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
