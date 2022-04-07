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
package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	"mlss-controlcenter-go/pkg/authcache"
	"mlss-controlcenter-go/pkg/client"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/auths"
	"mlss-controlcenter-go/pkg/service"
	"net/http"
	"strconv"
	"strings"
)

func UserStorageCheck(params auths.UserStorageCheckParams) middleware.Responder {
	username := params.Username
	path := params.Path

	logger.Logger().Debugf("UserStorageCheck username: %v, path: %v", username, path)

	userStorageCheckErr := service.UserStorageCheck(username, path)

	logger.Logger().Debugf("UserStorageCheck userStorageCheckErr: %v", userStorageCheckErr)

	if nil != userStorageCheckErr {
		return ResponderFunc(http.StatusBadRequest, "failed to check user and path", userStorageCheckErr.Error())
	}

	return GetResult([]byte{}, nil)
}

func UserNamespaceCheck(params auths.UserNamespaceCheckParams) middleware.Responder {
	response, err := service.UserNamespaceCheck(params.Username, params.Namespace)

	if nil != err {
		logger.Logger().Debugf("UserNamespaceCheck with error: %v", err.Error())
		return ResponderFunc(http.StatusBadRequest, "failed to check user and namespace", err.Error())
	}

	marshal, marshalErr := json.Marshal(response)

	return GetResult(marshal, marshalErr)
}

func UserStoragePathCheck(params auths.UserStoragePathCheckParams) middleware.Responder {
	err := service.UserStoragePathCheck(params.Username, params.Namespace, params.Path)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check user and namespace", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func AdminUserCheck(params auths.AdminUserCheckParams) middleware.Responder {
	err := service.AdminUserCheck(params.AdminUsername, params.Username)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check admin and user", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func CheckNamespace(params auths.CheckNamespaceParams) middleware.Responder {
	namespace := params.Namespace

	sessionUser := GetSessionByContext(params.HTTPRequest)

	logger.Logger().Infoln("session user is : ", sessionUser)
	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}
	logger.Logger().Debugf("namespace check : ", sessionUser)

	role, err := service.CheckNamespace(namespace, sessionUser.UserName)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check admin and user", err.Error())
	}

	marshal, marshalErr := json.Marshal(role)

	return GetResult(marshal, marshalErr)
}

func CheckNamespaceUser(params auths.CheckNamespaceUserParams) middleware.Responder {
	username := params.Username
	namespace := params.Namespace

	sessionUser := GetSessionByContext(params.HTTPRequest)

	if nil == sessionUser || "" == sessionUser.UserName {
		logger.Logger().Debugf("session user is : ", sessionUser)
		if sessionUser != nil {
			logger.Logger().Debugf("session user is : ", sessionUser.UserName)
		}
		return FailedToGetUserByToken()
	}
	logger.Logger().Debugf("namespace check : ", sessionUser)
	err := service.CheckNamespaceUser(namespace, sessionUser.UserName, username)

	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to check namespace and user", err.Error())
	}

	return GetResult([]byte{}, nil)
}

func CheckUserGetNamespace(params auths.CheckUserGetNamespaceParams) middleware.Responder {
	username := params.Username
	logger.Logger().Debugf("CheckUserGetNamespace username: %v", username)

	sessionUser := GetSessionByContext(params.HTTPRequest)

	if nil == sessionUser || "" == sessionUser.UserName {
		return FailedToGetUserByToken()
	}

	userNoteBook, err := service.CheckUserGetNamespace(sessionUser.UserName, username)
	if nil != err {
		return ResponderFunc(http.StatusBadRequest, "failed to get userNoteBook by username", err.Error())
	}

	marshal, marshalErr := json.Marshal(userNoteBook)

	return GetResult(marshal, marshalErr)
}

func CheckCurrentUserNamespacedNotebook(params auths.CheckCurrentUserNamespacedNotebookParams) middleware.Responder {
	namespace := params.Namespace
	notebook := params.Notebook
	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook namespace", namespace)
	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook notebook", notebook)
	headerToken := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)
	url := params.HTTPRequest.Header.Get(constants.X_ORIGINAL_URI)
	logger.Logger().Infof("URL:%v", url)

	createUser, proxyUser, err := service.GetNotebookProxyUserAndUser(namespace, notebook)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "failed to get proxy user", err.Error())
	}
	if strings.Contains(url, "/files/") && proxyUser != "" {
		logger.Logger().Infof("file donwload fotbit createUser:%v, proxyUser:%v", createUser, proxyUser)
		return ResponderFunc(http.StatusForbidden, "failed to get proxy user", "forbid")
	}
	authType := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_AUTH_TYPE)


	var sessionUser *models.SessionUser
	if authType == constants.TypeSYSTEM{
		user := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_USERID)
		logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook get user from header: %v", user)
		sessionUser = &models.SessionUser{UserName:user}
	}else{
		//get, b := ipCache.Get(realIp)
		tokenCache := authcache.TokenCache
		get, b := tokenCache.Get(headerToken)

		logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser b: %v", b)
		if b {
			session, err := common.FromInterfaceToSessionUser(get)
			if nil != err {
				return ResponderFunc(http.StatusInternalServerError, "failed to get sessionUser", err.Error())
			}
			sessionUser = session
		}

	}
	//logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser %+v", *sessionUser)

	if nil == sessionUser {
		return ResponderFunc(http.StatusForbidden, "failed to CheckCurrentUserNamespacedNotebook", "token is empty or invalid")
	}

	username := sessionUser.UserName
	user, err := service.GetUserByUserName(username)
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetUserByUserName", err.Error())
	}

	logger.Logger().Debugf("CheckCurrentUserNamespacedNotebook sessionUser.IsSuperAdmin: %v", sessionUser.IsSuperadmin)

	namespaceByName, err := service.GetNamespaceByName(namespace) //todo
	if err != nil {
		return ResponderFunc(http.StatusForbidden, "failed to GetNamespaceByName", err.Error())
	}
	if namespace != namespaceByName.Namespace {
		return ResponderFunc(http.StatusForbidden, "failed to access notebook", "namespace is not exist in db")
	}

	//var checkErr error
	checkPer := false
	if !sessionUser.IsSuperadmin {
		var flag = false
		groupIds, err := service.GetGroupIdListByUserIdRoleId(user.ID, constants.GA_ID)
		if err != nil {
			return ResponderFunc(http.StatusInternalServerError, "failed to GetGroupIdListByUserIdRoleId", err.Error())
		}
		for _, groupId := range groupIds {
			gn, err := service.GetGroupNamespaceByGroupIdAndNamespaceId(groupId, namespaceByName.ID)
			if err == nil {
				if gn.GroupID == groupId && gn.NamespaceID == namespaceByName.ID {
					flag = true
					checkPer = true
					break
				}
			}
		}
		if !flag {
			checkPer = createUser == username
			logger.Logger().Debugf("createUser : %v, proxyUser: %v, username:%v", createUser, proxyUser, username)
		}

	}

	if checkPer != true && !sessionUser.IsSuperadmin {
		return ResponderFunc(http.StatusForbidden, "access notebook failed", "user can access to edit this notebook")
	}

	var nbAddress string

	var notebookBuffer bytes.Buffer
	notebookBuffer.WriteString("/notebook/")
	notebookBuffer.WriteString(namespace)
	notebookBuffer.WriteString("/")
	notebookBuffer.WriteString(notebook)

	//if constants.BDP == strings.Split(namespace, "-")[2] {
	//	logger.Logger().Debugf("gatewayProperties.getBdpAddress(): %v", common.GetAppConfig().Core.Gateway.BdpAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdpAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//} else if constants.BDAP == strings.Split(namespace, "-")[2] && !(constants.SAFE == strings.Split(namespace, "-")[6]) {
	//	logger.Logger().Debugf("gatewayProperties.BdapAddress(): %v", common.GetAppConfig().Core.Gateway.BdapAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdapAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//} else if constants.BDAP == strings.Split(namespace, "-")[2] && constants.SAFE == strings.Split(namespace, "-")[6] {
	//	logger.Logger().Debugf("gatewayProperties.BdapsafeAddress(): %v", common.GetAppConfig().Core.Gateway.BdapsafeAddress)
	//	//notebookBuffer.WriteString(common.GetAppConfig().Core.Gateway.BdapsafeAddress)
	//	notebookBuffer.WriteString("/notebook/")
	//	notebookBuffer.WriteString(namespace)
	//	notebookBuffer.WriteString("/")
	//	notebookBuffer.WriteString(notebook)
	//}

	nbAddress = notebookBuffer.String()

	var userNotebookAddressResponse = models.UserNotebookAddressResponse{
		NotebookAddress: nbAddress,
	}

	marshal, marshalErr := json.Marshal(userNotebookAddressResponse)

	return GetResult(marshal, marshalErr)
}

func CheckResource(params auths.CheckResourceParams) middleware.Responder {
	k8sClient := client.GetK8sClient()
	nodeList, nodeErr := k8sClient.ListNamespaceNodes(params.Namespace)
	if nil != nodeErr {
		logger.Logger().Error("List namespace nodes err, ", nodeErr)
		return ResponderFunc(http.StatusInternalServerError, "failed List namespace nodes", nodeErr.Error())
	}
	if len(nodeList) < 1 {
		return ResponderFunc(http.StatusInternalServerError, "failed List namespace nodes", "namespace requires binding node first, please contact the administrator")
	}
	RQList, err := k8sClient.ListNamespacedResourceQuota(params.Namespace)
	if len(RQList) == 0 {
		return ResponderFunc(http.StatusInternalServerError, "failed List namespace resource quota", "resourceQuota is not existed, please contact the administrator")
	}
	if nil != err {
		logger.Logger().Error("List namespace resource quota err, ", err)
		return ResponderFunc(http.StatusInternalServerError, "failed List namespace resource quota", err.Error())
	}
	rq := RQList[0]
	hard := rq.Status.Hard
	used := rq.Status.Used
	hardCpu := hard["limits.cpu"]
	hardGpu := hard["requests.nvidia.com/gpu"]
	hardMemory := hard["limits.memory"]
	usedCpu := used["limits.cpu"]
	usedGpu := used["requests.nvidia.com/gpu"]
	usedMemory := used["limits.memory"]
	cpuFlag := CheckCpu(hardCpu.String(), usedCpu.String(), params.CPU)
	gpuFlag := CheckGpu(hardGpu.String(), usedGpu.String(), params.Gpu)
	memoryFlag := CheckMemory(hardMemory.String(), usedMemory.String(), params.Memory)
	if cpuFlag == false && gpuFlag == false && memoryFlag == false {
		logger.Logger().Infof("CpuFlag: %v, GpuFlag: %v, MemoryFlag: %v", cpuFlag, gpuFlag, memoryFlag)
		return ResponderFunc(http.StatusOK, "", "ok")
	}
	if cpuFlag {
		return ResponderFunc(http.StatusInternalServerError, "cpu can not greater than the limit", "cpu can not greater than the limit")
	}
	if gpuFlag {
		return ResponderFunc(http.StatusInternalServerError, "gpu can not greater than the limit", "gpu can not greater than the limit")
	}
	if memoryFlag {
		return ResponderFunc(http.StatusInternalServerError, "memory can not greater than the limit", "memory can not greater than the limit")
	}
	return ResponderFunc(http.StatusInternalServerError, "resource can not greater than the limit", "resource can not greater than the limit")
}

func CheckCpu(hardCpu, usedCpu, reqCpu string) bool {
	cpu, err := strconv.ParseFloat(reqCpu, 64)
	if err != nil {
		logger.Logger().Error("Request cpu convert float err, ", err)
		return false
	}
	hard, err := GetCPUValueWithUnits(hardCpu)
	if err != nil {
		logger.Logger().Error("Get cpu value with units err, ", err)
		return false
	}
	used, err := GetCPUValueWithUnits(usedCpu)
	if err != nil {
		logger.Logger().Error("Get cpu value with units err, ", err)
		return false
	}
	logger.Logger().Infof("CheckCPU, hard: %v, used: %v, cpu: %v, ", *hard, *used, cpu*1000)
	return *hard-*used-cpu*1000 < 0
}

func CheckGpu(hardGpu, usedGpu, reqGpu string) bool {
	gpu, err := strconv.Atoi(reqGpu)
	if err != nil {
		logger.Logger().Error("Request gpu atoi err, ", err)
		return false
	}
	hard, err := strconv.Atoi(hardGpu)
	if err != nil {
		logger.Logger().Error("Gpu used atoi err, ", err)
		return false
	}
	used, err := strconv.Atoi(usedGpu)
	if err != nil {
		logger.Logger().Error("Gpu used atoi err, ", err)
		return false
	}
	logger.Logger().Infof("CheckGPU, hard: %v, used: %v, gpu: %v, ", hard, used, gpu*1000)
	return hard-used-gpu < 0
}

func CheckMemory(hardMemory, usedMemory, reqMemory string) bool {
	memoryFloat, err := strconv.ParseFloat(reqMemory, 64)
	if err != nil {
		logger.Logger().Error("Request memory convert float err, ", err)
		return false
	}
	memory := memoryFloat * 1024 * 1024 * 1024
	hard, err := GetMemoryValueWithUnits(hardMemory)
	if err != nil {
		logger.Logger().Error("Get memory hard value with units err, ", err)
		return false
	}
	used, err := GetMemoryValueWithUnits(usedMemory)
	if err != nil {
		logger.Logger().Error("Get memory used value with units err, ", err)
		return false
	}
	logger.Logger().Infof("CheckMemory, hard: %v, used: %v, memory: %v, ", hard, used, memory)
	return *hard-*used-memory < 0
}

func GetMemoryValueWithUnits(withUnits string) (*float64, error) {
	f, e := strconv.ParseFloat(withUnits, 64)
	if e == nil {
		return &f, nil
	}
	if strings.HasSuffix(withUnits, "Mi") {
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-2]), 64)
		if e != nil {
			return nil, errors.New("memory not a number with Mi")
		} else {
			i := float * 1024 * 1024
			return &i, nil
		}
	}
	if strings.HasSuffix(withUnits, "Gi") {
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-2]), 64)
		if e != nil {
			return nil, errors.New("memory not a number with Gi")
		} else {
			i := float * 1024 * 1024 * 1024
			return &i, nil
		}
	}

	if strings.HasSuffix(withUnits, "k") {
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-1]), 64)
		if e != nil {
			return nil, errors.New("memory not a number with k")
		} else {
			i := float * 1024
			return &i, nil
		}
	}
	return nil, errors.New("GetMemoryValueWithUnits failed, str: " + withUnits)
}

func GetCPUValueWithUnits(withUnits string) (*float64, error) {
	f, e := strconv.ParseFloat(withUnits, 64)
	if e == nil {
		i := f * 1000
		return &i, nil
	}
	if strings.HasSuffix(withUnits, "m") {
		runes := []rune(withUnits)
		float, e := strconv.ParseFloat(string(runes[0:len(runes)-1]), 64)
		if e != nil {
			return nil, errors.New("cpu not a number with m")
		} else {
			i := float
			return &i, nil
		}
	}
	return nil, errors.New("GetCPUValueWithUnits failed, str: " + withUnits)
}

func CheckURLAccess(params auths.CheckURLAccessParams) middleware.Responder {
	token := authcache.TokenCache
	headerToken := params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN)
	url := params.HTTPRequest.Header.Get(constants.X_ORIGINAL_URI)
	logger.Logger().Infof("CheckURLAccess token: %v, url:%v ", headerToken, url)
	userIF, getBool := token.Get(headerToken)
	if !getBool && userIF == nil {
		return ResponderFunc(http.StatusForbidden, "Get User From Token Error", "User Not in Cache")
	}
	user, err := common.FromInterfaceToSessionUser(userIF)
	if err != nil {
		return ResponderFunc(http.StatusInternalServerError, "Get User From Token Error", err.Error())
	}

	urls := strings.Split(url, "/")
	ns := ""
	mfService := ""
	if len(urls) < 4 {
		return ResponderFunc(http.StatusInternalServerError, "Path is not correct", "Path")
	} else {
		ns = urls[2]
		mfService = urls[3]
		logger.Logger().Infof("ns:%v,mfservice:%v ", ns, mfService)
	}

	permission := false
	//check
	mfCreateUser, err := service.GetMFServiceUser(ns, mfService)
	if err != nil {
		logger.Logger().Errorf("get MFUser error:%v", err.Error())
		return ResponderFunc(http.StatusForbidden, "Check URL False", "GetMFServiceUser Error:"+err.Error())
	}
	if *mfCreateUser == user.UserName {
		logger.Logger().Debugf("mf create user:%v", *mfCreateUser)
		permission = true
		return getSuccessRes(params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN), nil)
	}

	//check namespace auth
	role, err := service.CheckNamespace(ns, user.UserName)
	if role == "GA" || role == "SA" {
		permission = true
	}
	logger.Logger().Infof("get MFUser role:%v", role)

	if !permission {
		return ResponderFunc(http.StatusForbidden, "Check URL False", "User Do Not Have Service Permission")
	}

	return getSuccessRes(params.HTTPRequest.Header.Get(constants.AUTH_HEADER_TOKEN), nil)
}

func ResponderFunc(code int, error string, msg string) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(code)
		payload, _ := json.Marshal(models.Error{
			Error: error,
			Code:  int32(code),
			Msg:   msg,
		})
		w.Write(payload)
	})
}

func CheckMLFlowResource(params auths.CheckMLFlowResourceAccessParams) middleware.Responder {
	//todo remove to caddy
	var r http.Request
	err := copier.Copy(&r, params.HTTPRequest)
	if err != nil {
		logger.Logger().Error("MLFlow copy request error: ", err.Error())
		ResponderFunc(http.StatusInternalServerError, err.Error(), "copy request error")
	}

	r.Header.Set(constants.AUTH_REAL_PATH, getOriginURI(params.HTTPRequest))
	r.Header.Set(constants.AUTH_REAL_METHOD, getOriginMethod(params.HTTPRequest))
	//r. = getOriginURI(params.HTTPRequest)

	pBool, err := checkMLFlowPermissionAuth(&r)
	if err != nil {
		logger.Logger().Error("MLFlow copy request error: ", err.Error())
		ResponderFunc(http.StatusInternalServerError, err.Error(), "checkPermissionAuth error")
	}
	if !pBool {
		logger.Logger().Error("")
		return ResponderFunc(http.StatusForbidden, service.PermissionDeniedError.Error(), "permission denied")
	}

	requestUri := params.RequestURI
	isSA := isSuperAdmin(params.HTTPRequest)
	user := getUserID(params.HTTPRequest)

	b, err := service.CheckMLFlowResource(requestUri, user, isSA)

	if err != nil {
		if err == service.PermissionDeniedError {
			return ResponderFunc(http.StatusForbidden, err.Error(), "permission denied")
		}
		return ResponderFunc(http.StatusInternalServerError, err.Error(), "permission denied")
	}

	//TODO
	if b {
		return ResponderFunc(http.StatusForbidden, err.Error(), "permission denied")
	}

	return GetResult([]byte{}, nil)
}

func getOriginURI(r *http.Request) string {
	return r.Header.Get("X-Original-URI")
}

func getOriginMethod(r *http.Request) string {
	return r.Header.Get("X-Original-METHOD")
}
