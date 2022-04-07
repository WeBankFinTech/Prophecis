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
	"mlss-controlcenter-go/pkg/controller"
	mw "mlss-controlcenter-go/pkg/restapi/restapi/middleware"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/proxy_user"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"mlss-controlcenter-go/pkg/restapi/restapi/operations"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/alerts"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/auths"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/groups"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/inters"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/keys"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/logins"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/logouts"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/namespaces"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/resources"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/roles"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/samples"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/storages"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/users"
)

//go:generate swagger generate server --target ..\..\restapi --name MlssCc --spec ..\swagger\swagger.yaml --model-package ../models

func configureFlags(api *operations.MlssCcAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MlssCcAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GroupsAddGroupHandler = groups.AddGroupHandlerFunc(func(params groups.AddGroupParams) middleware.Responder {
		return controller.AddGroup(params)
	})
	api.GroupsGetGroupByUsernameHandler = groups.GetGroupByUsernameHandlerFunc(func(params groups.GetGroupByUsernameParams) middleware.Responder {
		return controller.GetGroupByUsername(params)
	})

	api.KeysAddKeyHandler = keys.AddKeyHandlerFunc(func(params keys.AddKeyParams) middleware.Responder {
		return controller.AddKey(params)
	})
	api.NamespacesAddNamespaceHandler = namespaces.AddNamespaceHandlerFunc(func(params namespaces.AddNamespaceParams) middleware.Responder {
		return controller.AddNamespace(params)
	})
	api.GroupsAddNamespaceToGroupHandler = groups.AddNamespaceToGroupHandlerFunc(func(params groups.AddNamespaceToGroupParams) middleware.Responder {
		return controller.AddNamespaceToGroup(params)
	})
	api.RolesAddRoleHandler = roles.AddRoleHandlerFunc(func(params roles.AddRoleParams) middleware.Responder {
		return controller.AddRole(params)
	})
	api.StoragesAddStorageHandler = storages.AddStorageHandlerFunc(func(params storages.AddStorageParams) middleware.Responder {
		return controller.AddStorage(params)
	})
	api.GroupsAddStorageToGroupHandler = groups.AddStorageToGroupHandlerFunc(func(params groups.AddStorageToGroupParams) middleware.Responder {
		return controller.AddStorageToGroup(params)
	})
	api.UsersAddUserHandler = users.AddUserHandlerFunc(func(params users.AddUserParams) middleware.Responder {
		return controller.AddUser(params)
	})
	api.GroupsAddUserToGroupHandler = groups.AddUserToGroupHandlerFunc(func(params groups.AddUserToGroupParams) middleware.Responder {
		return controller.AddUserToGroup(params)
	})
	api.AuthsAdminUserCheckHandler = auths.AdminUserCheckHandlerFunc(func(params auths.AdminUserCheckParams) middleware.Responder {
		return controller.AdminUserCheck(params)
	})
	api.IntersAuthInterceptorHandler = inters.AuthInterceptorHandlerFunc(func(params inters.AuthInterceptorParams) middleware.Responder {
		return controller.AuthInterceptor(params)
	})
	api.AuthsCheckCurrentUserNamespacedNotebookHandler = auths.CheckCurrentUserNamespacedNotebookHandlerFunc(func(params auths.CheckCurrentUserNamespacedNotebookParams) middleware.Responder {
		return controller.CheckCurrentUserNamespacedNotebook(params)
	})
	api.AuthsCheckNamespaceHandler = auths.CheckNamespaceHandlerFunc(func(params auths.CheckNamespaceParams) middleware.Responder {
		return controller.CheckNamespace(params)
	})
	api.AuthsCheckNamespaceUserHandler = auths.CheckNamespaceUserHandlerFunc(func(params auths.CheckNamespaceUserParams) middleware.Responder {
		return controller.CheckNamespaceUser(params)
	})
	api.AuthsCheckUserGetNamespaceHandler = auths.CheckUserGetNamespaceHandlerFunc(func(params auths.CheckUserGetNamespaceParams) middleware.Responder {
		return controller.CheckUserGetNamespace(params)
	})
	api.KeysDeleteByNameHandler = keys.DeleteByNameHandlerFunc(func(params keys.DeleteByNameParams) middleware.Responder {
		return controller.DeleteByName(params)
	})
	api.GroupsDeleteGroupByIDHandler = groups.DeleteGroupByIDHandlerFunc(func(params groups.DeleteGroupByIDParams) middleware.Responder {
		return controller.DeleteGroupById(params)
	})
	api.GroupsDeleteGroupByNameHandler = groups.DeleteGroupByNameHandlerFunc(func(params groups.DeleteGroupByNameParams) middleware.Responder {
		return controller.DeleteGroupByName(params)
	})
	api.GroupsDeleteGroupNamespaceHandler = groups.DeleteGroupNamespaceHandlerFunc(func(params groups.DeleteGroupNamespaceParams) middleware.Responder {
		return controller.DeleteGroupNamespace(params)
	})
	api.NamespacesDeleteNamespaceByNameHandler = namespaces.DeleteNamespaceByNameHandlerFunc(func(params namespaces.DeleteNamespaceByNameParams) middleware.Responder {
		return controller.DeleteNamespaceByName(params)
	})
	api.StoragesDeleteStorageByIDHandler = storages.DeleteStorageByIDHandlerFunc(func(params storages.DeleteStorageByIDParams) middleware.Responder {
		return controller.DeleteStorageByID(params)
	})
	api.StoragesDeleteStorageByPathHandler = storages.DeleteStorageByPathHandlerFunc(func(params storages.DeleteStorageByPathParams) middleware.Responder {
		return controller.DeleteStorageByPath(params)
	})
	api.GroupsDeleteStorageFromGroupHandler = groups.DeleteStorageFromGroupHandlerFunc(func(params groups.DeleteStorageFromGroupParams) middleware.Responder {
		return controller.DeleteStorageFromGroup(params)
	})
	api.UsersDeleteUserByIDHandler = users.DeleteUserByIDHandlerFunc(func(params users.DeleteUserByIDParams) middleware.Responder {
		return controller.DeleteUserById(params)
	})
	api.UsersDeleteUserByNameHandler = users.DeleteUserByNameHandlerFunc(func(params users.DeleteUserByNameParams) middleware.Responder {
		return controller.DeleteUserByName(params)
	})
	api.GroupsDeleteUserFromGroupHandler = groups.DeleteUserFromGroupHandlerFunc(func(params groups.DeleteUserFromGroupParams) middleware.Responder {
		return controller.DeleteUserFromGroup(params)
	})
	api.GroupsGetAllGroupNamespaceByNamespaceIDHandler = groups.GetAllGroupNamespaceByNamespaceIDHandlerFunc(func(params groups.GetAllGroupNamespaceByNamespaceIDParams) middleware.Responder {
		return controller.GetAllGroupNamespaceByNamespaceId(params)
	})
	api.GroupsGetAllGroupStorageByStorageIDHandler = groups.GetAllGroupStorageByStorageIDHandlerFunc(func(params groups.GetAllGroupStorageByStorageIDParams) middleware.Responder {
		return controller.GetAllGroupStorageByStorageId(params)
	})
	api.GroupsGetAllGroupsHandler = groups.GetAllGroupsHandlerFunc(func(params groups.GetAllGroupsParams) middleware.Responder {
		return controller.GetAllGroups(params)
	})
	api.NamespacesGetAllNamespacesHandler = namespaces.GetAllNamespacesHandlerFunc(func(params namespaces.GetAllNamespacesParams) middleware.Responder {
		return controller.GetAllNamespaces(params)
	})
	api.StoragesGetAllStorageHandler = storages.GetAllStorageHandlerFunc(func(params storages.GetAllStorageParams) middleware.Responder {
		return controller.GetAllStorage(params)
	})
	api.GroupsGetAllUserGroupByUserIDHandler = groups.GetAllUserGroupByUserIDHandlerFunc(func(params groups.GetAllUserGroupByUserIDParams) middleware.Responder {
		return controller.GetAllUserGroupByUserId(params)
	})
	api.UsersGetAllUsersHandler = users.GetAllUsersHandlerFunc(func(params users.GetAllUsersParams) middleware.Responder {
		return controller.GetAllUsers(params)
	})
	api.KeysGetByNameHandler = keys.GetByNameHandlerFunc(func(params keys.GetByNameParams) middleware.Responder {
		return controller.GetByName(params)
	})
	api.GroupsGetCurrentUserNamespaceWithRoleHandler = groups.GetCurrentUserNamespaceWithRoleHandlerFunc(func(params groups.GetCurrentUserNamespaceWithRoleParams) middleware.Responder {
		return controller.GetCurrentUserNamespaceWithRole(params)
	})
	api.GroupsGetGroupByGroupIDHandler = groups.GetGroupByGroupIDHandlerFunc(func(params groups.GetGroupByGroupIDParams) middleware.Responder {
		return controller.GetGroupByGroupId(params)
	})
	api.GroupsGetGroupByNameHandler = groups.GetGroupByNameHandlerFunc(func(params groups.GetGroupByNameParams) middleware.Responder {
		return controller.GetGroupByName(params)
	})
	api.ResourcesGetLabelsOfNodeHandler = resources.GetLabelsOfNodeHandlerFunc(func(params resources.GetLabelsOfNodeParams) middleware.Responder {
		return controller.GetLabelsOfNode(params)
	})
	api.UsersGetMyUsersHandler = users.GetMyUsersHandlerFunc(func(params users.GetMyUsersParams) middleware.Responder {
		return controller.GetMyUsers(params)
	})
	api.NamespacesGetNamespaceByNameHandler = namespaces.GetNamespaceByNameHandlerFunc(func(params namespaces.GetNamespaceByNameParams) middleware.Responder {
		return controller.GetNamespaceByName(params)
	})
	api.ResourcesGetNodeByNameHandler = resources.GetNodeByNameHandlerFunc(func(params resources.GetNodeByNameParams) middleware.Responder {
		return controller.GetNodeByName(params)
	})
	api.RolesGetRoleByIDHandler = roles.GetRoleByIDHandlerFunc(func(params roles.GetRoleByIDParams) middleware.Responder {
		return controller.GetRoleById(params)
	})
	api.RolesGetRoleByNameHandler = roles.GetRoleByNameHandlerFunc(func(params roles.GetRoleByNameParams) middleware.Responder {
		return controller.GetRoleByName(params)
	})
	api.RolesGetRolesHandler = roles.GetRolesHandlerFunc(func(params roles.GetRolesParams) middleware.Responder {
		return controller.GetRoles(params)
	})
	api.UsersGetSAByNameHandler = users.GetSAByNameHandlerFunc(func(params users.GetSAByNameParams) middleware.Responder {
		return controller.GetSAByName(params)
	})
	api.StoragesGetStorageByPathHandler = storages.GetStorageByPathHandlerFunc(func(params storages.GetStorageByPathParams) middleware.Responder {
		return controller.GetStorageByPath(params)
	})
	api.UsersGetUserByUserIDHandler = users.GetUserByUserIDHandlerFunc(func(params users.GetUserByUserIDParams) middleware.Responder {
		return controller.GetUserByUserId(params)
	})
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams) middleware.Responder {
		return controller.GetUserByUserName(params)
	})
	api.IntersIPInterceptorHandler = inters.IPInterceptorHandlerFunc(func(params inters.IPInterceptorParams) middleware.Responder {
		return controller.IpInterceptor(params)
	})
	api.LogoutsLogoutHandler = logouts.LogoutHandlerFunc(func(params logouts.LogoutParams) middleware.Responder {
		return controller.Logout(params)
	})
	api.AlertsReceiveTaskAlertHandler = alerts.ReceiveTaskAlertHandlerFunc(func(params alerts.ReceiveTaskAlertParams) middleware.Responder {
		return controller.ReceiveTaskAlert(params)
	})
	api.AlertsPostAlertHandler = alerts.PostAlertHandlerFunc(func(params alerts.PostAlertParams) middleware.Responder {
		return middleware.NotImplemented("operation AlertsPostAlert has not yet been implemented")
	})
	api.SamplesSampleGetHandler = samples.SampleGetHandlerFunc(func(params samples.SampleGetParams) middleware.Responder {
		return controller.SampleGet(params)
	})
	api.SamplesSamplePostHandler = samples.SamplePostHandlerFunc(func(params samples.SamplePostParams) middleware.Responder {
		return controller.SamplePost(params)
	})
	api.ResourcesSetNamespaceRQHandler = resources.SetNamespaceRQHandlerFunc(func(params resources.SetNamespaceRQParams) middleware.Responder {
		return controller.SetNamespaceRQ(params)
	})
	api.LoginsUMLoginHandler = logins.UMLoginHandlerFunc(func(params logins.UMLoginParams) middleware.Responder {
		return controller.UMLogin(params)
	})
	api.GroupsUpdateGroupHandler = groups.UpdateGroupHandlerFunc(func(params groups.UpdateGroupParams) middleware.Responder {
		return controller.UpdateGroup(params)
	})
	api.GroupsUpdateGroupStorageHandler = groups.UpdateGroupStorageHandlerFunc(func(params groups.UpdateGroupStorageParams) middleware.Responder {
		return controller.UpdateGroupStorage(params)
	})
	api.GroupsGetNamespaceByGroupIDAndNamespaceHandler = groups.GetNamespaceByGroupIDAndNamespaceHandlerFunc(func(params groups.GetNamespaceByGroupIDAndNamespaceParams) middleware.Responder {
		return controller.GetNamespaceByGroupIDAndNamespace(params)
	})
	api.GroupsGetNamespacesByGroupIDHandler = groups.GetNamespacesByGroupIDHandlerFunc(func(params groups.GetNamespacesByGroupIDParams) middleware.Responder {
		return controller.GetNamespacesByGroupId(params)
	})
	api.GroupsGetUserGroupByUserIDAndGroupIDHandler = groups.GetUserGroupByUserIDAndGroupIDHandlerFunc(func(params groups.GetUserGroupByUserIDAndGroupIDParams) middleware.Responder {
		return controller.GetUserGroupByUserIdAndGroupId(params)
	})
	api.GroupsDeleteUserGroupByUserIDAndGroupIDHandler = groups.DeleteUserGroupByUserIDAndGroupIDHandlerFunc(func(params groups.DeleteUserGroupByUserIDAndGroupIDParams) middleware.Responder {
		return controller.DeleteUserGroupByUserIdAndGroupId(params)
	})
	api.ResourcesUpdateLabelsHandler = resources.UpdateLabelsHandlerFunc(func(params resources.UpdateLabelsParams) middleware.Responder {
		return controller.UpdateLabels(params)
	})
	api.ResourcesAddLabelsHandler = resources.AddLabelsHandlerFunc(func(params resources.AddLabelsParams) middleware.Responder {
		return controller.AddLabels(params)
	})
	api.ResourcesRemoveNodeLabelHandler = resources.RemoveNodeLabelHandlerFunc(func(params resources.RemoveNodeLabelParams) middleware.Responder {
		return controller.RemoveNodeLabel(params)
	})
	api.NamespacesUpdateNamespaceHandler = namespaces.UpdateNamespaceHandlerFunc(func(params namespaces.UpdateNamespaceParams) middleware.Responder {
		return controller.UpdateNamespace(params)
	})
	api.NamespacesGetMyNamespaceHandler = namespaces.GetMyNamespaceHandlerFunc(func(params namespaces.GetMyNamespaceParams) middleware.Responder {
		return controller.GetMyNamespace(params)
	})
	api.NamespacesListNamespaceByRoleNameAndUserNameHandler = namespaces.ListNamespaceByRoleNameAndUserNameHandlerFunc(
		func(params namespaces.ListNamespaceByRoleNameAndUserNameParams) middleware.Responder {
			return controller.ListNamespaceByRoleAndUserName(params)
		})
	api.RolesUpdateRoleHandler = roles.UpdateRoleHandlerFunc(func(params roles.UpdateRoleParams) middleware.Responder {
		return controller.UpdateRole(params)
	})
	api.StoragesUpdateStorageHandler = storages.UpdateStorageHandlerFunc(func(params storages.UpdateStorageParams) middleware.Responder {
		return controller.UpdateStorage(params)
	})
	api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder {
		return controller.UpdateUser(params)
	})
	api.GroupsUpdateUserGroupHandler = groups.UpdateUserGroupHandlerFunc(func(params groups.UpdateUserGroupParams) middleware.Responder {
		return controller.UpdateUserGroup(params)
	})
	api.GroupsGetCurrentUserStoragePathHandler = groups.GetCurrentUserStoragePathHandlerFunc(func(params groups.GetCurrentUserStoragePathParams) middleware.Responder {
		return controller.GetCurrentUserStoragePath(params)
	})
	api.IntersUserInterceptorHandler = inters.UserInterceptorHandlerFunc(func(params inters.UserInterceptorParams) middleware.Responder {
		return controller.UserInterceptor(params)
	})
	api.AuthsUserNamespaceCheckHandler = auths.UserNamespaceCheckHandlerFunc(func(params auths.UserNamespaceCheckParams) middleware.Responder {
		return controller.UserNamespaceCheck(params)
	})
	api.AuthsUserStorageCheckHandler = auths.UserStorageCheckHandlerFunc(func(params auths.UserStorageCheckParams) middleware.Responder {
		return controller.UserStorageCheck(params)
	})
	api.AuthsUserStoragePathCheckHandler = auths.UserStoragePathCheckHandlerFunc(func(params auths.UserStoragePathCheckParams) middleware.Responder {
		return controller.UserStoragePathCheck(params)
	})
	api.LoginsLDAPLoginHandler = logins.LDAPLoginHandlerFunc(func(params logins.LDAPLoginParams) middleware.Responder {
		return controller.LDAPLogin(params)
	})
	api.LoginsGetRsaPubKeyHandler = logins.GetRsaPubKeyHandlerFunc(func(params logins.GetRsaPubKeyParams) middleware.Responder {
		return controller.GetRsaPubKey(params)
	})
	api.AuthsCheckGroupByUserHandler = auths.CheckGroupByUserHandlerFunc(func(params auths.CheckGroupByUserParams) middleware.Responder {
		return controller.CheckGroupByUser(params)
	})
	api.GroupsGetUserGroupsHandler = groups.GetUserGroupsHandlerFunc(func(params groups.GetUserGroupsParams) middleware.Responder {
		return controller.GetUserGroups(params)
	})
	api.ProxyUserAddProxyUserHandler = proxy_user.AddProxyUserHandlerFunc(func(params proxy_user.AddProxyUserParams) middleware.Responder {
		return controller.AddProxyUser(params)
	})
	api.ProxyUserDeleteProxyUserHandler = proxy_user.DeleteProxyUserHandlerFunc(func(params proxy_user.DeleteProxyUserParams) middleware.Responder {
		return controller.DeleteProxyUser(params)
	})
	api.ProxyUserUpdateProxyUserHandler = proxy_user.UpdateProxyUserHandlerFunc(func(params proxy_user.UpdateProxyUserParams) middleware.Responder {
		return controller.UpdateProxyUser(params)
	})
	api.ProxyUserGetProxyUserHandler = proxy_user.GetProxyUserHandlerFunc(func(params proxy_user.GetProxyUserParams) middleware.Responder {
		return controller.GetProxyUser(params)
	})
	api.ProxyUserListProxyUserHandler = proxy_user.ListProxyUserHandlerFunc(func(params proxy_user.ListProxyUserParams) middleware.Responder {
		return controller.ProxyUserList(params)
	})
	api.ProxyUserGetProxyUserByUserIDHandler = proxy_user.GetProxyUserByUserIDHandlerFunc(func(params proxy_user.GetProxyUserByUserIDParams) middleware.Responder {
		return controller.GetProxyUserByUser(params)
	})
	api.AuthsCheckResourceHandler = auths.CheckResourceHandlerFunc(func(params auths.CheckResourceParams) middleware.Responder {
		return controller.CheckResource(params)
	})
	api.AuthsCheckURLAccessHandler = auths.CheckURLAccessHandlerFunc(func(params auths.CheckURLAccessParams) middleware.Responder {
		return controller.CheckURLAccess(params)
	})
	api.UsersGetUserTokenHandler = users.GetUserTokenHandlerFunc(func(params users.GetUserTokenParams) middleware.Responder {
		return controller.GetUserToken(params)
	})
	api.ProxyUserProxyUserCheckHandler = proxy_user.ProxyUserCheckHandlerFunc(func(params proxy_user.ProxyUserCheckParams) middleware.Responder {
		return controller.ProxyUserCheck(params)
	})
	api.AuthsCheckMLFlowResourceAccessHandler = auths.CheckMLFlowResourceAccessHandlerFunc(
		func(params auths.CheckMLFlowResourceAccessParams) middleware.Responder{
		return 	controller.CheckMLFlowResource(params)
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
	cxt := mw.NewContextMiddleWare(handler)

	return cxt
}
