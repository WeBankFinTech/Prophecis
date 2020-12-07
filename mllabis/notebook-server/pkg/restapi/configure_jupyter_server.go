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
	"github.com/dre1080/recover"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	mw "webank/AIDE/notebook-server/pkg/middleware"
	"webank/AIDE/notebook-server/pkg/restapi/operations"
)

//go:generate swagger generate server --target ../../notebook-server --name JupyterServer --spec ../../../../../../../swagger.yaml

func configureFlags(api *operations.JupyterServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.JupyterServerAPI) http.Handler {
	logger.Logger().Info("configure_jupyter_server.go,configureAPI start")
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.DeleteNamespacedNotebookHandler = operations.DeleteNamespacedNotebookHandlerFunc(func(params operations.DeleteNamespacedNotebookParams) middleware.Responder {
		return deleteNamespacedNotebook(params)
	})

	api.GetNamespacedNotebooksHandler = operations.GetNamespacedNotebooksHandlerFunc(func(params operations.GetNamespacedNotebooksParams) middleware.Responder {
		return getNamespacedNotebooks(params)

	})

	api.GetNamespacedUserNotebooksHandler = operations.GetNamespacedUserNotebooksHandlerFunc(func(params operations.GetNamespacedUserNotebooksParams) middleware.Responder {
		return getNamespacedUserNotebooks(params)
	})

	api.GetUserNotebooksHandler = operations.GetUserNotebooksHandlerFunc(func(params operations.GetUserNotebooksParams) middleware.Responder {
		return getUserNotebooks(params)
	})

	logger.Logger().Info("configure_jupyter_server.go, api.PostNamespacedNotebookHandler == nil")
	api.PostNamespacedNotebookHandler = operations.PostNamespacedNotebookHandlerFunc(func(params operations.PostNamespacedNotebookParams) middleware.Responder {
		return postNamespacedNotebook(params)
	})

	api.GetDashboardsHandler = operations.GetDashboardsHandlerFunc(func(params operations.GetDashboardsParams) middleware.Responder {
		return getDashboards(params)
	})

	api.PatchNamespacedNotebookHandler = operations.PatchNamespacedNotebookHandlerFunc(func(params operations.PatchNamespacedNotebookParams) middleware.Responder {
		return patchNamespacedNotebook(params)
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

	return logging.Handle(recovery(auth(handler)))
}
