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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
	"webank/AIDE/notebook-server/pkg/commons/config"
	ccClient "webank/AIDE/notebook-server/pkg/commons/controlcenter/client"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/commons/utils"
	mw "webank/AIDE/notebook-server/pkg/middleware"
	"webank/AIDE/notebook-server/pkg/models"
	"webank/AIDE/notebook-server/pkg/restapi/operations"
)

func getSuperadmin(r *http.Request) string {
	return r.Header.Get(mw.CcAuthSuperadmin)
}

// Create a Namespaced Notebook
func postNamespacedNotebook(params operations.PostNamespacedNotebookParams) middleware.Responder {
	logger.Logger().Infof("create time now: %v", time.Now().Format(time.RFC1123))

	//namespace := params.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	logger.Logger().Infof("create notebook in namespace %s", *params.Notebook.Namespace)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	// check namespace
	err := client.UserNamespaceCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *params.Notebook.Namespace)
	if err != nil {
		logger.Logger().Errorf("UserNamespaceCheck failed: %s, token: %s, userName: %s, namespace: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *params.Notebook.Namespace)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	//check notebook request
	//err = client.PostNotebookRequestCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), params.Notebook)
	//FIXME MLSS Change: change Description to Message
	err = utils.CheckNotebook(params)
	if err != nil {
		logger.Logger().Errorf("PostNotebookRequestCheck failed: %s, token: %s, userName: %s, Notebook: %v", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, params.Notebook)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "StatusInternalServerError",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// get guid from cc
	gid, uid, token, err := client.GetGUIDFromUserId(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId)
	logger.Logger().Infof("gid: %s, uid: %s", *gid, *uid)
	if err != nil || gid == nil || uid == nil {
		logger.Logger().Errorf("GetGUIDFromUserId failed: %s, token: %s, userName: %s, Notebook: %v", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, params.Notebook)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	ncClient := utils.GetNBCClient()

	// Create Notebook
	err = ncClient.CreateYarnResourceConfigMap(params.Notebook, currentUserId)
	if err != nil {
		logger.Logger().Errorf("User %v CreateNotebook configmap failed:%v", currentUserId, err.Error())
		var error = models.Error{
			Code:    500,
			Error:   "create notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(500)
			w.Write(payload)
		})
	}
	err = ncClient.CreateNotebook(params.Notebook, currentUserId, *gid, *uid, *token)
	if err != nil {
		logger.Logger().Errorf("CreateNotebook failed: %s, Code: %s", err.Error())
		var error = models.Error{
			Code:    500,
			Error:   "create notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(500)
			w.Write(payload)
		})
	}
	//return operations.NewPostNamespacedNotebookOK()

	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

// Delete a Namespaced Notebook
func deleteNamespacedNotebook(params operations.DeleteNamespacedNotebookParams) middleware.Responder {
	namespace := params.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.CcAuthUser)
	logger.Logger().Infof("delete Namespace(%s) Notebook %s.", namespace, params.Notebook)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	// check current user accessable for namespace
	body, err := client.CheckNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
	if err != nil {
		logger.Logger().Errorf("CheckNamespace failed: %s, token: %s, namespace: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	logger.Logger().Infof("CheckNamespace body: %v", string(body))
	//get list from k8s
	ncClient := utils.GetNBCClient()
	listFromK8s, err := ncClient.GetNotebooks(namespace, "", "")
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "getNamespacedNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := parseNotebookFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}

	configFromK8s, err := ncClient.GetNBConfigMaps(namespace)
	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
	//filter books
	var role string
	err = models.GetResultData(body, &role)
	if err != nil {
		logger.Logger().Errorf("parse resultMap Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parse resultMap Unmarshal failed.",
		})
	}
	//if current user is GU, he only access his own notebook
	if role == "GU" {
		for _, k := range res {
			if k.Name == params.Notebook && k.User != currentUserId {
				logger.Logger().Errorf("GU only access his own notebook.")
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    404,
					Error:   "you can't access this notebook",
					Message: "GU only access his own notebook.",
				})
			}
		}
	}

	// delete nb
	err = ncClient.DeleteNotebook(params.Notebook, namespace)
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "delete notebook failed",
		})
	}

	//return operations.NewDeleteNamespacedNotebookOK()
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

// List Namespaced Notebook
func getNamespacedNotebooks(params operations.GetNamespacedNotebooksParams) middleware.Responder {
	err := checkPageAndSize(*params.Page, *params.Size)
	if err != nil {
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "StatusInternalServerError",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	namespace := params.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	logger.Logger().Infof("Get Notebooks in namespace %s.", namespace)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	// TODO check if user is SA or user has GA/GU access to namespace
	body, err := client.CheckNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
	if err != nil {
		logger.Logger().Errorf("CheckNamespace failed: %s, token: %s, namespace: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	logger.Logger().Infof("CheckNamespace body: %v", string(body))

	ncClient := utils.GetNBCClient()

	listFromK8s, err := ncClient.GetNotebooks(namespace, "", "")
	configFromK8s, err := ncClient.GetNBConfigMaps(namespace)

	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "getNamespacedNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := parseNotebookFromK8s(listFromK8s)

	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}

	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
	var finalResult []*models.Notebook
	//todo filter books
	//roleJsonStr := string(body)
	//var resultMap map[string]string
	var role string
	err = models.GetResultData(body, &role)
	if err != nil {
		logger.Logger().Errorf("parse resultMap Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parse resultMap Unmarshal failed.",
		})
	}
	// if current user is GU, only return his own nb
	if role == "GU" {
		for _, v := range res {
			if v.User == currentUserId {
				finalResult = append(finalResult, v)
			}
		}
		res = finalResult
	}

	if res == nil {
		logger.Logger().Infof("res == nil")
		res = make([]*models.Notebook, 0)
	}
	//add clusterName to filter notebook
	//clusterName := params.ClusterName
	//logger.Logger().Debugf("namespace_userId clusterName: %v", clusterName)
	//res = getNoteBookByClusterName(*clusterName, res)

	//pagination
	//string to int
	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(res, page, size)
	pages := utils.GetPages(len(res), size)
	logger.Logger().Debugf("utils.GetPages(len(res), size): %v", pages)

	//int to string
	//pagesStr := strconv.Itoa(pages)
	//totalStr := strconv.Itoa(len(res))
	//return operations2.NewGetNamespacedUserNotebooksOK().WithPayload(&models.GetNotebooksResponse{
	//	Notebooks: subList,
	//	Pages:     &pagesStr,
	//	Total:     &totalStr,
	//})

	//response := &models.GetNotebooksResponse{
	//	Notebooks: subList,
	//	Pages:     &pagesStr,
	//	Total:     &totalStr,
	//}

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(res),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})

}

func checkPageAndSize(page string, size string) error {
	_, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	_, err = strconv.Atoi(size)
	if err != nil {
		return err
	}
	return nil
}

// Get the list of Notebooks in the given Namespace belongs to a User
func getNamespacedUserNotebooks(params operations.GetNamespacedUserNotebooksParams) middleware.Responder {
	var err error
	err = checkPageAndSize(*params.Page, *params.Size)
	if err != nil {
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "StatusInternalServerError",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	namespace := params.Namespace
	//currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	superadmin := getSuperadmin(params.HTTPRequest)
	const TRUE = "true"
	userId := params.User
	workDir := ""
	if params.WorkDir != nil {
		workDir = *params.WorkDir
	}
	logger.Logger().Infof("Get user %s's Notebooks in namespace %s and workDir is %s.", userId, namespace, workDir)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	var listFromK8s interface{}
	ncClient := utils.GetNBCClient()
	// TODO check if user is SA
	if superadmin == TRUE && namespace == "null" && userId == "null" {
		logger.Logger().Debugf("getNamespacedUserNotebooks debug is SA: %v", superadmin)
		listFromK8s, err = ncClient.GetNotebooks("", "", workDir)
	} else {
		var err error
		if !(namespace == "null" && userId == "null") {
			err = client.CheckNamespaceUser(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace, params.User)
		}
		if err != nil {
			logger.Logger().Errorf("CheckNamespaceUser failed: %s, token: %s, namespace: %s, user: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace, params.User)
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusForbidden)
				payload, _ := json.Marshal(models.Error{
					Error:   "forbidden",
					Code:    http.StatusForbidden,
					Message: err.Error(),
				})
				w.Write(payload)
			})
		}

		listFromK8s, err = ncClient.GetNotebooks(namespace, userId, workDir)
	}
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "getNamespacedUserNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := parseNotebookFromK8s(listFromK8s)
	logger.Logger().Debugf("getNamespacedUserNotebooks debug notebookFromK8s.length: %v", len(notebookFromK8s.Items))
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}
	configFromK8s, err := ncClient.GetNBConfigMaps(namespace)
	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
	logger.Logger().Debugf("getNamespacedUserNotebooks debug res.length: %v", len(res))
	if res == nil {
		logger.Logger().Infof("res == nil")
		res = make([]*models.Notebook, 0)
	}

	logger.Logger().Debugf("getNamespacedUserNotebooks debug res.length: %v", len(res))
	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(res, page, size)
	pages := utils.GetPages(len(res), size)

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(res),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)

	})

}

// Get the list of Notebooks belongs to a User
func getUserNotebooks(params operations.GetUserNotebooksParams) middleware.Responder {
	var err error

	err = checkPageAndSize(*params.Page, *params.Size)
	if err != nil {
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "StatusInternalServerError",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	//currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	userId := params.User
	logger.Logger().Infof("Get Notebooks of user %s.", userId)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	clusterName := params.ClusterName
	logger.Logger().Debugf("getUserNotebooks for clusterName: %v", clusterName)
	//check currentUser accessibility for user
	nsResFromCC, err := client.CheckUserGetNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), userId)
	if err != nil {
		logger.Logger().Errorf("CheckUserGetNamespace failed: %s, token: %s, user: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), userId)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	logger.Logger().Infof("CheckUserGetNamespace nsResFromCC: %v, length: %v", string(nsResFromCC), len(nsResFromCC))

	//Get User Role
	var userNotebookVOFromCC models.UserNotebookVO
	err = models.GetResultData(nsResFromCC, &userNotebookVOFromCC)
	if err != nil {
		logger.Logger().Errorf("parse userNotebookVOFromCC Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parse userNotebookVOFromCC Unmarshal failed.",
		})
	}

	// get list from k8s
	ncClient := utils.GetNBCClient()
	var listFromK8s interface{}
	if userNotebookVOFromCC.Role == "GA" {
		listFromK8s, err = ncClient.GetNotebooks("", "", "")
	} else {
		listFromK8s, err = ncClient.GetNotebooks("", userId, "")
	}
	//listFromK8s, err := ncClient.GetNotebooks("", userId, "")
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "getUserNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := parseNotebookFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}
	configFromK8s, err := ncClient.GetNBConfigMaps("")
	if err != nil {
		logger.Logger().Errorf("Failed to get notebook configmap %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    500,
			Error:   err.Error(),
			Message: "Failed to get notebook configmap" + err.Error(),
		})
	}
	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)

	//filter notebook
	strings := userNotebookVOFromCC.NamespaceList
	var finalResult []*models.Notebook

	if userNotebookVOFromCC.Role == "GA" {
		for _, v := range res {
			if v.User == userId {
				finalResult = append(finalResult, v)
				continue
			}
			for _, nv := range strings {
				if v.Namespace == nv {
					finalResult = append(finalResult, v)
				}
			}
		}
		res = finalResult
	}
	if res == nil {
		logger.Logger().Infof("res == nil")
		res = make([]*models.Notebook, 0)
	} else {
		//res = getNoteBookByClusterName(*clusterName, res)
	}

	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(res, page, size)
	pages := utils.GetPages(len(res), size)

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(res),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
	})

}

func parseNotebookFromK8s(listFromK8s interface{}) (*models.NotebookFromK8s, error) {
	var notebookFromK8s models.NotebookFromK8s
	bytes, err := json.Marshal(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("Marshal failed. %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(bytes, &notebookFromK8s)
	if err != nil {
		logger.Logger().Errorf("parse []models.K8sNotebook Unmarshal failed. %v", err.Error())
		return nil, err
	}
	return &notebookFromK8s, nil
}

func getDashboards(params operations.GetDashboardsParams) middleware.Responder {
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	logger.Logger().Infof("Get Dashboards of user %s.", currentUserId)
	superadmin := getSuperadmin(params.HTTPRequest)
	const TRUE = "true"
	var err error
	var listFromK8s interface{}
	ncClient := utils.GetNBCClient()
	role := "GU"
	var userNotebookVOFromCC models.UserNotebookVO
	//clusterName := params.ClusterName

	if superadmin == TRUE {
		role = "SA"
		listFromK8s, err = ncClient.GetNotebooks("", "", "")
	} else {
		client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
		if client == nil{
			return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
				Code:    500,
				Error:   "CC Client init Error",
				Message: "getUserNotebooks failed",
			})
		}
		nsResFromCC, err := client.
			CheckUserGetNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId)
		if err != nil{
			return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
				Code:    500,
				Error:   "CC Client init Error",
				Message: "getUserNotebooks failed",
			})
		}
		err = models.GetResultData(nsResFromCC, &userNotebookVOFromCC)

		if userNotebookVOFromCC.Role == "GU"{
			listFromK8s, err = ncClient.GetNotebooks("", currentUserId, "")
		}else{
			role = "GA"
			listFromK8s, err = ncClient.GetNotebooks("", "", "")
		}

	}
	if err != nil {
		return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "getUserNotebooks failed",
		})
	}
	notebookFromK8s, err := parseNotebookFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
			Code:    404,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}

	configFromK8s, err := ncClient.GetNBConfigMaps("")
	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)

	var finalResult []*models.Notebook
	if role == "GA"{
		namespaceList := userNotebookVOFromCC.NamespaceList
		for _, v := range res {
			if v.User == currentUserId {
				finalResult = append(finalResult, v)
				continue
			}
			for _, nv := range namespaceList {
				if v.Namespace == nv {
					finalResult = append(finalResult, v)
				}
			}
		}
		res = finalResult
	}

	//res = getNoteBookByClusterName(*clusterName, res)

	//nbTotal:current user's notebooks; nbRunning: number of current user's notebooks is Ready for status; gpuCount:for gpu number of nbRunning
	nbTotal := len(res)
	nbRunning := 0
	gpuCount := 0
	for _, nb := range res {
		if "Ready" == nb.Status {
			if "" != nb.Gpu {
				gpu, err := strconv.Atoi(nb.Gpu)
				if nil != err {
					logger.Logger().Errorf("strconv.Atoi(nb.Gpu) failed. %v", err.Error())
					return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
						Code:    404,
						Error:   err.Error(),
						Message: "strconv.Atoi(nb.Gpu) failed." + nb.Gpu,
					})
				}
				gpuCount += gpu
			}
			nbRunning += 1
		}
	}
	mp := map[string]int{"nbTotal": nbTotal, "gpuCount": gpuCount, "nbRunning": nbRunning}
	marshal, err := json.Marshal(mp)
	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	logger.Logger().Infof("Get Dashboards of result %s.", json.RawMessage(marshal))

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

func patchNamespacedNotebook(params operations.PatchNamespacedNotebookParams) middleware.Responder {
	logger.Logger().Infof("patch time now: %v", time.Now().Format(time.RFC1123))
	namespace := params.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	logger.Logger().Infof("Patch Notebooks in namespace %s.", namespace)
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	// check namespace
	err := client.UserNamespaceCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *params.Notebook.Namespace)
	if err != nil {
		logger.Logger().Errorf("UserNamespaceCheck failed: %s, token: %s, userName: %s, namespace: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *params.Notebook.Namespace)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	ncClient := utils.GetNBCClient()

	// Create Notebook
	err = ncClient.PatchYarnSettingConfigMap(params.Notebook)
	if err != nil {
		logger.Logger().Errorf("CreateNotebook failed: %s, Code: %s", err.Error())
		var error = models.Error{
			Code:    500,
			Error:   "create notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(500)
			w.Write(payload)
		})
	}
	//return operations.NewPostNamespacedNotebookOK()

	var result = models.Result{
		Code:    "200",
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})

}
