package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	stringsUtil "strings"
	"time"
	cfg "webank/AIDE/notebook-server/pkg/commons/config"
	"webank/AIDE/notebook-server/pkg/commons/constants"
	ccClient "webank/AIDE/notebook-server/pkg/commons/controlcenter/client"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/commons/utils"
	"webank/AIDE/notebook-server/pkg/es"
	mw "webank/AIDE/notebook-server/pkg/middleware"
	"webank/AIDE/notebook-server/pkg/models"
	"webank/AIDE/notebook-server/pkg/mongo"
	"webank/AIDE/notebook-server/pkg/restapi/operations"

	"gopkg.in/mgo.v2/bson"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var PlatformNamespace = viper.GetString(cfg.PlatformNamespace)

func getSuperadmin(r *http.Request) string {
	return r.Header.Get(mw.CcAuthSuperadmin)
}

func getSuperadminBool(r *http.Request) bool {
	superAdminStr := r.Header.Get(mw.CcAuthSuperadmin)
	if "true" == superAdminStr {
		return true
	}
	return false
}

// Create a Namespaced Notebook
func generateNotebookInMongo(req *models.NewNotebookRequest, user string) (*utils.NotebookInMongo, error) {
	if req == nil {
		return nil, errors.New("notebook request param can't be nil")
	}
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	now := time.Now().Format(constants.TimeLayoutStr)
	cpu := fmt.Sprintf("%0.1f", *req.CPU)
	workspaceVolume := utils.MountInfoInMongo{
		SubPath: req.WorkspaceVolume.SubPath,
	}
	if req.WorkspaceVolume.MountType != nil {
		workspaceVolume.MountType = *req.WorkspaceVolume.MountType
	}
	if req.WorkspaceVolume.LocalPath != nil {
		workspaceVolume.LocalPath = *req.WorkspaceVolume.LocalPath
	}
	if req.WorkspaceVolume.Size != nil {
		workspaceVolume.Size = int(*req.WorkspaceVolume.Size)
	}
	if req.WorkspaceVolume.MountPath != nil {
		workspaceVolume.MountPath = *req.WorkspaceVolume.MountPath
	}
	if req.WorkspaceVolume.AccessMode != nil {
		workspaceVolume.AccessMode = *req.WorkspaceVolume.AccessMode
	}

	dataVolume := utils.MountInfoInMongo{}
	if req.DataVolume != nil {
		dataVolume.SubPath = req.DataVolume.SubPath
		if req.DataVolume.MountType != nil {
			dataVolume.MountType = *req.DataVolume.MountType
		}
		if req.DataVolume.LocalPath != nil {
			dataVolume.LocalPath = *req.DataVolume.LocalPath
		}
		if req.DataVolume.Size != nil {
			dataVolume.Size = int(*req.DataVolume.Size)
		}
		if req.DataVolume.MountPath != nil {
			dataVolume.MountPath = *req.DataVolume.MountPath
		}
		if req.DataVolume.AccessMode != nil {
			dataVolume.AccessMode = *req.DataVolume.AccessMode
		}
	}

	nb := utils.NotebookInMongo{
		Id:          uid.String(),
		Namespace:   *req.Namespace,
		Name:        *req.Name,
		Status:      constants.NoteBookStatusStoping,
		NotebookMsg: "",
		EnableFlag:  1,
		User:        user,
		ProxyUser:   req.ProxyUser,
		CreateTime:  now,
		Image: utils.ImageInMongo{
			ImageType: *req.Image.ImageType,
			ImageName: *req.Image.ImageName,
		},
		Cpu: cpu,
		Memory: utils.MemoryInMongo{
			MemoryAmount: *req.Memory.MemoryAmount,
			MemoryUnit:   req.Memory.MemoryUnit,
		},
		WorkspaceVolume: workspaceVolume,
		DataVolume:      []utils.MountInfoInMongo{dataVolume},
		ExtraResources:  req.ExtraResources,
		Queue:           req.Queue,
		ExecutorCores:   req.ExecutorCores,
		Executors:       req.Executors,
		DriverMemory:    req.DriverMemory,
		SparkSessionNum: int(req.SparkSessionNum),
	}
	return &nb, nil
}
func postNamespacedNotebook(params operations.PostNamespacedNotebookParams) middleware.Responder {
	// 1.校验mongo notebook是否存在
	query := bson.M{"namespace": params.Notebook.Namespace, "name": params.Notebook.Name, "enable_flag": 1}
	count, err := mongo.CountNoteBook(query)
	if err != nil {
		logger.Logger().Errorf("fail to check notebook existed in mongo: %s", err.Error())
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
	if count > 0 {
		msg := "Notebook在该命名空间下已存在，请更换Notebook名称进行创建"
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusInternalServerError,
				Message: msg,
			})
			w.Write(payload)
		})
	}
	// 2. 校验notebook 在k8s平台是否存在
	logger.Logger().Infof("create time now: %v", time.Now().Format(time.RFC1123))
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	ncClient := utils.GetNBCClient()
	existed, err := ncClient.CheckNotebookExisted(*params.Notebook.Namespace, *params.Notebook.Name)
	if err != nil {
		logger.Logger().Errorf("Check Notebook existed error: ", err.Error())
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
	if existed {
		msg := "Notebook在该命名空间下已存在，请更换Notebook名称进行创建"
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusInternalServerError,
				Message: msg,
			})
			w.Write(payload)
		})
	}

	//Get User Message
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	containerUser := currentUserId
	var gid, uid, token *string
	if params.Notebook.ProxyUser != "" {
		containerUser = params.Notebook.ProxyUser
		gid, uid, token, err = client.GetProxyUserGUIDFromUser(params.HTTPRequest.Header.Get(ccClient.CcAuthToken),
			currentUserId, containerUser)
	} else {
		gid, uid, token, err = client.GetGUIDFromUserId(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId)
	}
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

	// Check Namespace
	logger.Logger().Infof("create notebook in namespace %s", *params.Notebook.Namespace)
	err = client.UserNamespaceCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *params.Notebook.Namespace)
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

	// 3.check notebook request
	//err = client.PostNotebookRequestCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), params.Notebook)
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

	// 4.sync to mongo
	var notebookInMongo *utils.NotebookInMongo
	notebookInMongo, err = generateNotebookInMongo(params.Notebook, currentUserId)
	notebookInMongo.Status = constants.NoteBookStatusCreating
	err = mongo.InsertNoteBook(*notebookInMongo)
	if err != nil {
		logger.Logger().Errorf("fail to insert notebook to mongo, notebook: %+v, err: %v", *notebookInMongo, err)
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

	// 5.创建k8s notebook
	nbClient := utils.GetNBCClient()
	maxRetryNum := 5
	sparkBytes := make([]byte, 0)
	for count := 0; count < maxRetryNum; count++ {
		sparkBytes, err = nbClient.CreateNotebook(params.Notebook, currentUserId, *gid, *uid, *token)
		if err != nil {
			if !stringsUtil.Contains(err.Error(), "provided port is already allocated") {
				logger.Logger().Errorf("CreateNotebook failed: %s, Code: %s", err.Error())
				var error = models.Error{
					Code:    http.StatusInternalServerError,
					Error:   "create notebook failed",
					Message: err.Error(),
				}
				return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
					payload, _ := json.Marshal(error)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(payload)
				})
			}

			//端口被占用，开始重试
			logger.Logger().Infof("[%d]provided port is already allocated, start retry ...\n", count)
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}

	// 6.创建k8s notebook configmap
	err = nbClient.CreateYarnResourceConfigMap(params.Notebook, containerUser, sparkBytes)
	if err != nil {
		logger.Logger().Errorf("CreateNotebook configmap failed: %s, Code: %s", err.Error())
		var error = models.Error{
			Code:    http.StatusInternalServerError,
			Error:   "create notebook configmap failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(payload)
		})
	}
	//return operations.NewPostNamespacedNotebookOK()

	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
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
func deleteNotebookById(params operations.DeleteNotebookByIDParams) middleware.Responder {
	// 1.校验notebook是否存在
	m := bson.M{"id": params.ID, "enable_flag": 1}
	count, err := mongo.CountNoteBook(m)
	if err != nil {
		logger.Logger().Errorf("fail to count notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "error",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	if count <= 0 {
		logger.Logger().Errorf("nootbook not exists")
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "nootbook not exist",
				Code:    http.StatusForbidden,
				Message: fmt.Sprintf("notebook not exist in id %q", params.ID),
			})
			w.Write(payload)
		})
	}

	nbInMongo, err := mongo.GetNotebookById(params.ID)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 2. 校验当前用户是否具有该namespace权限
	namespace := nbInMongo.Namespace
	name := nbInMongo.Name
	currentUserId := params.HTTPRequest.Header.Get(mw.CcAuthUser)
	logger.Logger().Infof("delete Namespace(%s) Notebook %s.", namespace, name)
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
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

	// 3.get list from k8s
	ncClient := utils.GetNBCClient()
	listFromK8s, err := ncClient.GetNotebooks(namespace, "", "")
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "getNamespacedNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
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
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "parse resultMap Unmarshal failed.",
		})
	}
	//if current user is GU, he only access his own notebook
	if role == "GU" {
		for _, k := range res {
			if k.Name == name && nbInMongo.User != currentUserId {
				logger.Logger().Errorf("GU only access his own notebook.")
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusNotFound,
					Error:   "you can't access this notebook",
					Message: "GU only access his own notebook.",
				})
			}
		}
	}

	// 4.删除k8s notebook
	err = mongo.DeleteNotebookById(params.ID)
	if err != nil {
		logger.Logger().Errorf("fail to delete notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 5. delete nb from mongo
	err = ncClient.DeleteNotebook(name, namespace)
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "delete notebook failed",
		})
	}

	//return operations.NewDeleteNamespacedNotebookOK()
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
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
	// 1.参数校验
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

	// 2 namespace权限校验
	namespace := params.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	logger.Logger().Infof("Get Notebooks in namespace %s.", namespace)
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
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

	var role string
	err = models.GetResultData(body, &role)
	if err != nil {
		logger.Logger().Errorf("parse resultMap Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "parse resultMap Unmarshal failed.",
		})
	}
	// 3. 获取notebook 列表
	notebooksInMongo, err := mongo.ListNootBook(params.Namespace, role, currentUserId)
	if err != nil {
		logger.Logger().Errorf("fail to get notebooks from mongo, namespace: %s, role: %s,"+
			" currentUserId: %s, err: %v", params.Namespace, role, currentUserId, err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "fail to get notebooks.",
		})
	}
	nsMap := map[string]struct{}{}
	for _, v := range notebooksInMongo {
		if _, ok := nsMap[v.Namespace]; !ok {
			nsMap[v.Namespace] = struct{}{}
		}
	}

	// 4.获取k8s notebook状态
	ncClient := utils.GetNBCClient()
	notebookStatusMap := make(map[string]string)
	if len(nsMap) > 0 {
		for ns := range nsMap {
			listFromK8s, err := ncClient.ListNotebooks(ns)
			if err != nil {
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: fmt.Sprintf("fail to get notebook by ns %q", ns),
				})
			}
			notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
			if err != nil {
				logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "parseNotebookFromK8s Unmarshal failed.",
				})
			}
			configFromK8s, err := ncClient.GetNBConfigMaps(ns)
			if err != nil {
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: fmt.Sprintf("fail to get notebook configmap by ns %q", ns),
				})
			}
			res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
			for _, nb := range res {
				if nb != nil {
					key := fmt.Sprintf("%s#$#%s", nb.Namespace, nb.Name)
					notebookStatusMap[key] = nb.Status
				}
			}
		}
	}
	logger.Logger().Debugf("notebook status map: %+v", notebookStatusMap)

	//add clusterName to filter notebook
	clusterName := params.ClusterName
	logger.Logger().Debugf("namespace_userId clusterName: %v", *clusterName)
	notebooksInMongo = getNoteBookInMongoByClusterName(*clusterName, notebooksInMongo)
	// 5. 生成response 数据列表
	nbs := generateNotebooks(notebooksInMongo, notebookStatusMap)

	//pagination
	//string to int
	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(nbs, page, size)
	pages := utils.GetPages(len(nbs), size)
	logger.Logger().Debugf("utils.GetPages(len(res), size): %v", pages)

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(nbs),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

func generateNotebooks(nbs []utils.NotebookInMongo, m map[string]string) []*models.Notebook {
	result := make([]*models.Notebook, 0)
	if len(nbs) > 0 {
		for _, v := range nbs {
			logger.Logger().Debugf("mongo generateNotebooks nootbook: %+v\n", v)
			gpu := ""
			if v.ExtraResources != "" {
				m := make(map[string]string)
				json.Unmarshal([]byte(v.ExtraResources), &m)
				if gpuStr, ok := m["nvidia.com/gpu"]; ok {
					gpu = gpuStr
				}
				logger.Logger().Debugf("generateNotebooks gpu: %v", gpu)
			}

			dataVolumes := make([]*models.MountInfo, 0)
			if len(v.DataVolume) > 0 {
				for _, vol := range v.DataVolume {
					size := (int64)(vol.Size)
					mountInfo := models.MountInfo{
						AccessMode: &vol.AccessMode,
						LocalPath:  &vol.LocalPath,
						MountPath:  &vol.MountPath,
						MountType:  &vol.MountType,
						Size:       &size,
						SubPath:    vol.SubPath,
					}
					dataVolumes = append(dataVolumes, &mountInfo)
				}
			}
			var size int64 = int64(v.WorkspaceVolume.Size)
			workspaceVolume := models.MountInfo{
				AccessMode: &v.WorkspaceVolume.AccessMode,
				LocalPath:  &v.WorkspaceVolume.LocalPath,
				MountPath:  &v.WorkspaceVolume.MountPath,
				MountType:  &v.WorkspaceVolume.MountType,
				Size:       &size,
				SubPath:    v.WorkspaceVolume.SubPath,
			}
			status := v.Status
			key := fmt.Sprintf("%s#$#%s", v.Namespace, v.Name)
			if v, ok := m[key]; ok {
				status = v
			}

			nb := models.Notebook{
				ID:              v.Id,
				CPU:             v.Cpu,
				DataVolume:      dataVolumes,
				DriverMemory:    v.DriverMemory,
				ExecutorCores:   v.ExecutorCores,
				ExecutorMemory:  v.ExecutorMemory,
				Executors:       v.Executors,
				Gpu:             gpu,
				Image:           v.Image.ImageName,
				Memory:          fmt.Sprintf("%0.1f%s", v.Memory.MemoryAmount, v.Memory.MemoryUnit),
				Name:            v.Name,
				Namespace:       v.Namespace,
				Pods:            "",
				ProxyUser:       v.ProxyUser,
				Queue:           v.Queue,
				Service:         "",
				SparkSessionNum: int64(v.SparkSessionNum),
				SrtImage:        "",
				Status:          status,
				Uptime:          v.CreateTime,
				User:            v.User,
				Volumns:         "",
				WorkspaceVolume: &workspaceVolume,
			}
			result = append(result, &nb)
		}
	}
	return result
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
	/**
	todo
	Step
	*/

	// 1.校验参数
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

	// 2.按照相应权限获取k8s notebook列表
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
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	var listFromK8s interface{}
	ncClient := utils.GetNBCClient()
	// TODO check if user is SA
	if superadmin == TRUE && namespace == "null" && userId == "null" {
		logger.Logger().Debugf("getNamespacedUserNotebooks debug is SA: %v", superadmin)
		listFromK8s, err = ncClient.GetNotebooks("", "", workDir)
	} else {
		err := client.CheckNamespaceUser(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace, params.User)
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

		//listFromK8s, err = ncClient.GetNotebooks(namespace, userId, workDir)
		listFromK8s, err = ncClient.ListNotebooks(namespace)
	}
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "getNamespacedUserNotebooks failed",
		})
	}

	//parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
	logger.Logger().Debugf("getNamespacedUserNotebooks debug notebookFromK8s.length: %v", len(notebookFromK8s.Items))
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
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
	////add clusterName to filter notebook
	//clusterName := params.ClusterName
	//logger.Logger().Debugf("getNamespacedUserNotebooks debug namespace: %v, user: %v, clusterName: %v", namespace, userId, *clusterName)
	//res = getNoteBookByClusterName(*clusterName, res)
	//logger.Logger().Debugf("getNamespacedUserNotebooks debug res.length: %v", len(res))

	notebookStatusMap := make(map[string]string)
	for _, nb := range res {
		if nb != nil {
			key := fmt.Sprintf("%s#$#%s", nb.Namespace, nb.Name)
			notebookStatusMap[key] = nb.Status
		}
	}
	logger.Logger().Debugf("notebook status map: %+v", notebookStatusMap)

	// 3.获取mongo notebook列表
	logger.Logger().Debugf("mongo.ListNamespaceUserNotebook namespace: %s, userId: %v", params.Namespace, userId)
	notebooksInMongo, err := mongo.ListNamespaceUserNotebook(params.Namespace, userId)
	if err != nil {
		logger.Logger().Errorf("fail to list namespace user notebooks from mongo, namespace: %s,"+
			" userId: %s, err: %v", params.Namespace, userId, err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "fail to list notebooks.",
		})
	}

	//add clusterName to filter notebook
	clusterName := params.ClusterName
	logger.Logger().Debugf("namespace_userId clusterName: %v", *clusterName)
	notebooksInMongo = getNoteBookInMongoByClusterName(*clusterName, notebooksInMongo)
	// 4.生成resp
	nbs := generateNotebooks(notebooksInMongo, notebookStatusMap)

	//pagination
	//string to int
	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(nbs, page, size)
	pages := utils.GetPages(len(res), size)

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(nbs),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})

}

// Get the list of Notebooks belongs to a User
func getUserNotebooks(params operations.GetUserNotebooksParams) middleware.Responder {
	// 1.校验参数
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

	// 2.校验权限
	//currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	currentUserId := params.HTTPRequest.Header.Get(mw.CcAuthUser)
	userId := params.User
	isSA := getSuperadminBool(params.HTTPRequest)
	logger.Logger().Infof("Current User %s Get Notebooks of user %s.", currentUserId,userId)
	if !isSA && currentUserId != userId {
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusForbidden)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusForbidden,
				Message: "permission denied",
			})
			w.Write(payload)
		})
	}

	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	clusterName := params.ClusterName
	logger.Logger().Infof("getUserNotebooks for clusterName: %v", clusterName)
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

	// 3.获取notebook列表
	var notebooksInMongo []utils.NotebookInMongo
	if userNotebookVOFromCC.Role == "GA" {
		notebooksInMongo, err = mongo.ListUserNotebook("")
	} else {
		notebooksInMongo, err = mongo.ListUserNotebook(params.User)
	}

	strings := userNotebookVOFromCC.NamespaceList
	var gaNotebooksInMongo []utils.NotebookInMongo
	if userNotebookVOFromCC.Role == "GA" {
		for _, v := range notebooksInMongo {
			if v.User == userId {
				gaNotebooksInMongo = append(gaNotebooksInMongo, v)
				continue
			}
			for _, nv := range strings {
				if v.Namespace == nv {
					gaNotebooksInMongo = append(gaNotebooksInMongo, v)
				}
			}
		}
		notebooksInMongo = gaNotebooksInMongo
	}



	if err != nil {
		logger.Logger().Errorf("fail to list user notebooks from mong0,"+
			" userId: %s, err: %v", userId, err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Error:   err.Error(),
			Message: "fail to list notebooks.",
		})
	}
	nsMap := map[string]struct{}{}
	for _, v := range notebooksInMongo {
		if _, ok := nsMap[v.Namespace]; !ok {
			nsMap[v.Namespace] = struct{}{}
		}
	}

	// 4.获取k8s notebook状态
	ncClient := utils.GetNBCClient()
	notebookStatusMap := make(map[string]string)
	if len(nsMap) > 0 {
		for ns := range nsMap {
			listFromK8s, err := ncClient.ListNotebooks(ns)
			if err != nil {
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: fmt.Sprintf("fail to get notebook by ns %q", ns),
				})
			}
			notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
			if err != nil {
				logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: "parseNotebookFromK8s Unmarshal failed.",
				})
			}
			configFromK8s, err := ncClient.GetNBConfigMaps(ns)
			if err != nil {
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusInternalServerError,
					Error:   err.Error(),
					Message: fmt.Sprintf("fail to get notebook configmap by ns %q", ns),
				})
			}
			res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
			for _, nb := range res {
				if nb != nil {
					key := fmt.Sprintf("%s#$#%s", nb.Namespace, nb.Name)
					notebookStatusMap[key] = nb.Status
				}
			}
		}
	}

	logger.Logger().Infof("notebook status map: %+v", notebookStatusMap)

	//add clusterName to filter notebook
	logger.Logger().Infof("namespace_userId clusterName: %v", *clusterName)
	notebooksInMongo = getNoteBookInMongoByClusterName(*clusterName, notebooksInMongo)
	// 5.生成resp
	nbs := generateNotebooks(notebooksInMongo, notebookStatusMap)

	//pagination
	//string to int
	page, err := strconv.Atoi(*params.Page)
	size, err := strconv.Atoi(*params.Size)

	subList, err := utils.GetSubListByPageAndSize(nbs, page, size)
	pages := utils.GetPages(len(nbs), size)

	list := models.PageListVO{
		List:  subList,
		Pages: pages,
		Total: len(nbs),
	}

	marshal, err := json.Marshal(list)
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  json.RawMessage(marshal),
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})

}

func ParseNotebooksFromK8s(listFromK8s interface{}) (*models.NotebookFromK8s, error) {
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

func ParseNotebookFromK8s(nbFromK8s interface{}) (*models.K8sNotebook, error) {
	if nbFromK8s == nil {
		return nil, nil
	}
	logger.Logger().Debugf("ParseNotebookFromK8s nbFromK8s： %+v", nbFromK8s)
	var notebookFromK8s models.K8sNotebook
	bytes, err := json.Marshal(nbFromK8s)
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

	if superadmin == TRUE {
		role = "SA"
		listFromK8s, err = ncClient.GetNotebooks("", "", "")
	} else {
		client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
		if client == nil {
			return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
				Code:    500,
				Error:   "CC Client init Error",
				Message: "getUserNotebooks failed",
			})
		}
		nsResFromCC, err := client.CheckUserGetNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId)

		logger.Logger().Infof("getDashboards CheckUserGetNamespace nsResFromCC: %v", string(nsResFromCC))
		if err != nil {
			return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
				Code:    500,
				Error:   "CC Client init Error",
				Message: "getUserNotebooks failed",
			})
		}
		err = models.GetResultData(nsResFromCC, &userNotebookVOFromCC)

		if userNotebookVOFromCC.Role == "GU" {
			listFromK8s, err = ncClient.GetNotebooks("", currentUserId, "")
		} else {
			role = "GA"
			listFromK8s, err = ncClient.GetNotebooks("", "", "")
		}

	}
	if err != nil {
		return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "getUserNotebooks failed",
		})
	}
	logger.Logger().Infof("getDashboards listFromK8s:%v", listFromK8s)
	notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s Unmarshal failed.",
		})
	}

	configFromK8s, err := ncClient.GetNBConfigMaps("")
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s get notebook config map failed. %v", err.Error())
		return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "parseNotebookFromK8s get notebook config map failed.",
		})
	}
	res := models.ParseToNotebookRes(*notebookFromK8s, *configFromK8s)
	logger.Logger().Infof("getDashboards first res:%v", res)

	var finalResult []*models.Notebook
	if role == "GA" {
		namespaceList := userNotebookVOFromCC.NamespaceList
		logger.Logger().Infof("getDashboards namespaceList:%v", namespaceList)
		for _, v := range res {
			if v.User == currentUserId {
				finalResult = append(finalResult, v)
				continue
			}
			logger.Logger().Infof("getDashboards v.Namespace: %v", v.Namespace)
			for _, nv := range namespaceList {
				if v.Namespace == nv {
					finalResult = append(finalResult, v)
				}
			}
		}
		res = finalResult
	}
	logger.Logger().Infof("getDashboards finalResult res:%v", res)

	//res = getNoteBookByClusterName(*clusterName, res)
	//nbTotal:current user's notebooks; nbRunning: number of current user's notebooks is Ready for status; gpuCount:for gpu number of nbRunning
	nbTotal := len(res)
	nbRunning := 0
	gpuCount := 0
	for _, nb := range res {
		if nb != nil && "Ready" == nb.Status {
			if "" != nb.Gpu {
				gpu, err := strconv.Atoi(nb.Gpu)
				if nil != err {
					logger.Logger().Errorf("strconv.Atoi(nb.Gpu) failed. %v", err.Error())
					return operations.NewGetDashboardsNotFound().WithPayload(&models.Error{
						Code:    http.StatusNotFound,
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
	marshal, _ := json.Marshal(mp)
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
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

//FIXME MLSS Change: v_1.5.1 add logic for path when create notebook
func getNoteBookByClusterName(clusterName string, res []*models.Notebook) []*models.Notebook {
	return res
}

func getNoteBookInMongoByClusterName(clusterName string, res []utils.NotebookInMongo) []utils.NotebookInMongo {
	//clusterResult := make([]utils.NotebookInMongo, 0)
	//clusterResult = res
	return res
}

//FIXME MLSS Change: v_1.5.1 add logic for path when create notebook
func patchNamespacedNotebook(params operations.PatchNamespacedNotebookParams) middleware.Responder {
	logger.Logger().Infof("patch time now: %v", time.Now().Format(time.RFC1123))
	// 1.校验notebook在mongo是否存在
	m := bson.M{"id": params.ID, "enable_flag": 1}
	count, err := mongo.CountNoteBook(m)
	if err != nil {
		logger.Logger().Errorf("fail to count notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "error",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	if count <= 0 {
		logger.Logger().Errorf("nootbook not exists")
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "nootbook not exist",
				Code:    http.StatusForbidden,
				Message: fmt.Sprintf("notebook not exist in id %q", params.ID),
			})
			w.Write(payload)
		})
	}

	// 2. 获取notebook信息
	nbInMongo, err := mongo.GetNotebookById(params.ID)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 3.校验并更新参数
	if params.Notebook.CPU <= 0 {
		cpu, err := strconv.ParseFloat(nbInMongo.Cpu, 64)
		if err != nil {
			logger.Logger().Errorf("fail to parse cpu, cpu: %v, err: %v", nbInMongo.Cpu, err)
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusInternalServerError)
				payload, _ := json.Marshal(models.Error{
					Error:   "failed",
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				w.Write(payload)
			})
		}
		params.Notebook.CPU = cpu
	}
	if params.Notebook.MemoryAmount <= 0 {
		memoryUnit := nbInMongo.Memory.MemoryUnit
		params.Notebook.MemoryAmount = nbInMongo.Memory.MemoryAmount
		params.Notebook.MemoryUnit = &memoryUnit
	}
	if params.Notebook.ImageName == "" {
		params.Notebook.ImageName = nbInMongo.Image.ImageName
	}
	if params.Notebook.ImageType == "" {
		imageType := nbInMongo.Image.ImageType
		if imageType == "" {
			imageType = utils.Standard
		}
		params.Notebook.ImageType = imageType
	}
	namespace := nbInMongo.Namespace
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)

	logger.Logger().Infof("Patch Notebooks in namespace %s.", namespace)
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	// TODO check if user is SA or user has GA/GU access to namespace
	// 5. check namespace
	err = client.UserNamespaceCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, namespace)
	if err != nil {
		logger.Logger().Errorf("UserNamespaceCheck failed: %s, token: %s, "+
			"userName: %s, namespace: %s", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken),
			currentUserId, namespace)
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

	// 6.获取k8s notebook状态信息并校验
	// notebook status check
	ncClient := utils.GetNBCClient()
	nbInterface, err := ncClient.GetNotebookByNamespaceAndName(namespace, nbInMongo.Name)
	if err != nil {
		if !stringsUtil.Contains(err.Error(), "not found") {
			logger.Logger().Errorf("fail to get nootbook namespace: %s, name: %s, err: %v", namespace, nbInMongo.Name, err)
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusInternalServerError)
				payload, _ := json.Marshal(models.Error{
					Error:   "error",
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				w.Write(payload)
			})
		}
	}
	notebook, err := ParseNotebookFromK8s(nbInterface)
	if err != nil {
		logger.Logger().Errorf("fail to parse notebook notebook interface: %v, err: %v", nbInterface, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "error",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	if notebook != nil {
		logger.Logger().Debugf("patchNamespacedNotebook nootebook: %+v", notebook)
	}
	status := models.ParseNotebookStatus(notebook)
	if status == "" {
		status = nbInMongo.Status
	}
	switch status {
	case constants.NoteBookStatusTerminated, constants.NoteBookStatusCreating, constants.NoteBookStatusStoping:
		err = errors.New("when notebook status is Terminated、Creating、Stopping or Stopped, can't update notebook resource setting")
		logger.Logger().Errorf("patchNamespacedNotebook notebook status: %s, err: %v", status, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "error",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	default:
		logger.Logger().Debugf("patchNamespacedNotebook notebook status: %s", status)
	}

	// 4.校验notebook资源信息
	if status != constants.NoteBookStatusStoped {
		if err := utils.CheckPatchNotebook(params.Notebook, namespace, nbInMongo.Name); err != nil {
			logger.Logger().Errorf("fail to check notebook params, notebook request: %+v, err: %v", *params.Notebook, err)
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusInternalServerError)
				payload, _ := json.Marshal(models.Error{
					Error:   "failed",
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				w.Write(payload)
			})
		}
	}

	// 7.更新notebook信息到mongo
	//sparkSessionStr := strconv.FormatInt(params.Notebook.SparkSessionNum, 10)
	updateMap := bson.M{
		"driverMemory":    params.Notebook.DriverMemory,
		"executorCores":   params.Notebook.ExecutorCores,
		"executorMemory":  params.Notebook.ExecutorMemory,
		"executors":       params.Notebook.Executors,
		"queue":           params.Notebook.Queue,
		"sparkSessionNum": params.Notebook.SparkSessionNum,
	}
	if params.Notebook.CPU > 0 {
		updateMap["cpu"] = fmt.Sprintf("%0.1f", params.Notebook.CPU)
	}
	if params.Notebook.MemoryAmount > 0 && params.Notebook.MemoryUnit != nil && *params.Notebook.MemoryUnit != "" {
		updateMap["memory"] = utils.MemoryInMongo{
			MemoryAmount: params.Notebook.MemoryAmount,
			MemoryUnit:   *(params.Notebook.MemoryUnit),
		}
	}
	if params.Notebook.ExtraResources != "" {
		updateMap["extraResources"] = params.Notebook.ExtraResources
	}
	if params.Notebook.ImageName != "" {
		img := utils.ImageInMongo{
			ImageType: params.Notebook.ImageType,
			ImageName: params.Notebook.ImageName,
		}
		updateMap["image"] = img
	}
	logger.Logger().Debugf("patchNamespacedNotebook mongo updateMap: %+v", updateMap)
	if err := mongo.UpdateNotebookById(params.ID, updateMap); err != nil {
		logger.Logger().Errorf("fail to update notebook by id %q, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	if status == constants.NoteBookStatusStoped {
		var result = models.Result{
			Code:    fmt.Sprintf("%d", http.StatusOK),
			Message: "success",
			Result:  nil,
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(result)
			w.Write(payload)
			//w.WriteHeader(200)
		})
	}

	// 8.更新k8s notebook service与configmap信息
	//update spark session service
	sparkSessionMap, sparkService, portArray, err := ncClient.PatchSparkSessionService(params.Notebook,
		nbInMongo.Name, nbInMongo.Namespace)
	logger.Logger().Debugf("portArray: %v", portArray)
	defer func(portArray []int) {
		if len(portArray) > 0 {
			for _, port := range portArray {
				if _, ok := utils.Ports.M[port]; ok {
					delete(utils.Ports.M, port)
				}
			}
			utils.Ports.Mutex.Unlock()
		}
	}(portArray)
	if err != nil {
		logger.Logger().Errorf("fail to patch spark session service, err: %v\n", err)
		var error = models.Error{
			Code:    http.StatusInternalServerError,
			Error:   "update notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(payload)
		})
	}
	//logger.Logger().Debugf("patchNamespacedNotebook sparkService.Spec:%+v\n", *sparkService.Spec)
	if sparkService == nil {
		logger.Logger().Debugf("patchNamespacedNotebook sparkService is nil, sparkSessionMap: %v\n", sparkSessionMap)
	} else {
		logger.Logger().Debugf("patchNamespacedNotebook sparkService is not nil, sparkService: %+v, "+
			"sparkSessionMap: %v\n", *sparkService, sparkSessionMap)
	}

	err = ncClient.PatchYarnSettingConfigMap(params.Notebook, sparkSessionMap, nbInMongo.Name, nbInMongo.Namespace)
	if err != nil {
		logger.Logger().Errorf(" fail to update yarn setting configMap, Code: %s", err.Error())
		var error = models.Error{
			Code:    http.StatusInternalServerError,
			Error:   "update notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(payload)
		})
	}

	// 9. 更新k8s notebook crd
	// update notebook crd
	err = ncClient.PatchNoteBookCRD(namespace, nbInMongo.Name, sparkService, params.Notebook)
	if err != nil {
		logger.Logger().Errorf("fail to update notebook CRD, namespace: %s, notebook name: %s,"+
			"  Code: %s", namespace, nbInMongo.Name, err.Error())
		var error = models.Error{
			Code:    http.StatusInternalServerError,
			Error:   "update notebook failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(payload)
		})
	}

	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

func GetNotebookUser(params operations.GetNotebookUserParams) middleware.Responder {
	k8sNotebook, res, err := utils.GetNotebook(params.Namespace, params.Name)
	userIdLabel := PlatformNamespace + "-UserId"
	proxyUserIdLabel := PlatformNamespace + "-ProxyUserId"
	resp := models.GetNotebookUserResponse{
		User:      "",
		ProxyUser: "",
	}
	if err != nil {
		logger.Logger().Error("GetNotebook From K8s Error:%v", err.Error())
		if res != nil {
			if res.StatusCode == http.StatusNotFound {
				logger.Logger().Errorf("GetNotebook: Fail to get notebook by namespace %q and name %q",
					params.Namespace, params.Name)
				return buildResponse(http.StatusNotFound, nil, errors.New("notebook not found"))
			}
		}
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	notebook, err := parseNotebookFromK8s(k8sNotebook)
	if err != nil {
		logger.Logger().Error("parseNotebookFromK8s From K8s Error:%v", err.Error())
		return buildResponse(http.StatusInternalServerError, nil, err)
	}
	resp.User = notebook.Metadata.Labels[userIdLabel]
	resp.ProxyUser = notebook.Metadata.Labels[proxyUserIdLabel]

	return buildResponse(200, resp, nil)
}

func generateNootbookStatusResp(nb utils.NotebookInMongo) (*models.GetNotebookStatusResponse, error) {
	gpu := ""
	if nb.ExtraResources != "" {
		m := map[string]string{}
		err := json.Unmarshal([]byte(nb.ExtraResources), &m)
		if err != nil {
			return nil, err
		}
		if gpuStr, ok := m["nvidia.com/gpu"]; ok {
			gpu = gpuStr
		}
	}
	resp := models.GetNotebookStatusResponse{
		ContainersStatusInfo: nil,
		CreateTime:           nb.CreateTime,
		Image:                nb.Image.ImageName,
		Name:                 nb.Name,
		Namespace:            nb.Namespace,
		NotebookStatus:       nb.Status,
		NotebookStatusInfo:   "",
		ResourceReq: &models.Resource{
			CPU:    nb.Cpu,
			Gpu:    gpu,
			Memory: fmt.Sprintf("%0.1f%s", nb.Memory.MemoryAmount, nb.Memory.MemoryUnit),
		},
		ResourceLimit: &models.Resource{
			CPU:    nb.Cpu,
			Gpu:    gpu,
			Memory: fmt.Sprintf("%0.1f%s", nb.Memory.MemoryAmount, nb.Memory.MemoryUnit),
		},
	}
	return &resp, nil
}

func GetNamespacedNotebookStatus(params operations.GetNamespacedNotebookStatusParams) middleware.Responder {
	// 1.获取notebook信息
	nbInMgo, err := mongo.GetNotebookById(params.ID)
	if err != nil {
		logger.Logger().Error("fail to get notebook from mongo, id:%v,  Error:%v", params.ID, err.Error())
		return buildResponse(http.StatusNotFound, nil, errors.New("notebook not found"))
	}
	ns := nbInMgo.Namespace
	name := nbInMgo.Name

	// 2.获取k8s notebook信息
	k8sNotebook, res, err := utils.GetNotebook(ns, name)
	if err != nil {
		logger.Logger().Errorf("GetNamespacedNotebookStatus From K8s Error:%v", err.Error())
		if res != nil {
			if res.StatusCode == http.StatusNotFound {
				logger.Logger().Errorf("GetNamespacedNotebookStatus: Fail to get notebook by namespace %q and name %q",
					ns, name)
				resp, err := generateNootbookStatusResp(nbInMgo)
				if err != nil {
					return buildResponse(http.StatusInternalServerError, nil, err)
				}
				return buildResponse(http.StatusOK, *resp, nil)
			}
			logger.Logger().Debugf("GetNamespacedNotebookStatus From K8s res.StatusCode:%v", res.StatusCode)
		}
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	nb, err := parseNotebookFromK8s(k8sNotebook)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s From K8s Error:%v", err.Error())
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	var timeLayoutStr = "2006-01-02 15:04:05"
	createTime := nb.Metadata.CreationTimestamp.Format(timeLayoutStr)

	status := "NotReady"
	notebookStatusInfo := ""
	if nil != nb.Status.ContainerState.Running {
		status = "Ready"
	} else if nil != nb.Status.ContainerState.Waiting {
		logger.Logger().Debugf("GetNamespacedNotebookStatus for state Waiting: %v", nb.Status.ContainerState.Waiting.Message)
		status = "Waiting"
		notebookStatusInfo = fmt.Sprintf("reason: %s, message: %s",
			nb.Status.ContainerState.Waiting.Reason, nb.Status.ContainerState.Waiting.Message)
	} else if nil != nb.Status.ContainerState.Terminated {
		logger.Logger().Debugf("GetNamespacedNotebookStatus for state terminater: %v", nb.Status.ContainerState.Terminated.Message)
		status = "Terminated"
		notebookStatusInfo = fmt.Sprintf("reason: %s, message: %s",
			nb.Status.ContainerState.Terminated.Reason, nb.Status.ContainerState.Terminated.Message)
	} else {
		logger.Logger().Debugf("GetNamespacedNotebookStatus for state terminater: %v", nb.Status.ContainerState)
		status = "Waiting"
	}

	//get container info from pod
	podName := fmt.Sprintf("%s-0", nb.Metadata.Name)
	pod, _, err := utils.GetPod(ns, podName)
	if err != nil {
		if !stringsUtil.Contains(err.Error(), "not found") {
			logger.Logger().Errorf("GetNamespacedNotebookStatus: Fail to get pod info, namespace: %s, pod name: %s, "+
				"error: %v\n", ns, podName, err)
			return buildResponse(http.StatusInternalServerError, nil, err)
		}
		logger.Logger().Infof("%q not found in namespace %q\n", podName, ns)
	}
	containerStatusInfoList := []*models.ContainerStatusInfo{}
	if pod.Status != nil && len(pod.Status.ContainerStatuses) > 0 {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			containerName := containerStatus.Name
			containerContainerID := containerStatus.ContainerID

			startedTime := ""
			finishedTime := ""
			containerStatusType := "NotReady"
			message := ""
			if containerStatus.State.Running != nil {
				containerStatusType = "Running"
				startedTime = containerStatus.State.Running.StartedAt.Format(timeLayoutStr)
			} else if containerStatus.State.Terminated != nil {
				containerStatusType = "Terminated"
				startedTime = containerStatus.State.Terminated.StartedAt.Format(timeLayoutStr)
				finishedTime = containerStatus.State.Terminated.FinishedAt.Format(timeLayoutStr)
				message = fmt.Sprintf("containerId: %s, exitCode: %v, reason: %s, message: %s",
					containerStatus.State.Terminated.ContainerID, containerStatus.State.Terminated.ExitCode,
					containerStatus.State.Terminated.Reason, containerStatus.State.Terminated.Message)
			} else if containerStatus.State.Waiting != nil {
				containerStatusType = "Waiting"
				message = fmt.Sprintf("reason: %s, message: %s", containerStatus.State.Waiting.Reason,
					containerStatus.State.Waiting.Message)
			}

			lastStartedTime := ""
			lastFinishedTime := ""
			lastContainerStateType := "NotReady"
			lastMessage := ""
			if containerStatus.LastState.Running != nil {
				lastContainerStateType = "Running"
				lastStartedTime = containerStatus.LastState.Running.StartedAt.Format(timeLayoutStr)
			} else if containerStatus.LastState.Terminated != nil {
				lastContainerStateType = "Terminated"
				lastStartedTime = containerStatus.LastState.Terminated.StartedAt.Format(timeLayoutStr)
				lastFinishedTime = containerStatus.LastState.Terminated.FinishedAt.Format(timeLayoutStr)
				lastMessage = fmt.Sprintf("containerId: %s, exitCode: %v, reason: %s, message: %s",
					containerStatus.LastState.Terminated.ContainerID, containerStatus.LastState.Terminated.ExitCode,
					containerStatus.LastState.Terminated.Reason, containerStatus.LastState.Terminated.Message)
			} else if containerStatus.LastState.Waiting != nil {
				lastContainerStateType = "Waiting"
				lastMessage = fmt.Sprintf("reason: %s, message: %s", containerStatus.LastState.Waiting.Reason,
					containerStatus.LastState.Waiting.Message)
			}

			containerStatusInfo := models.ContainerStatusInfo{
				ContainerID:            containerContainerID,
				ContainerStatusType:    containerStatusType,
				ContainerName:          containerName,
				FineshedTime:           finishedTime,
				LastContainerStateType: lastContainerStateType,
				LastFinishedTime:       lastFinishedTime,
				LastMessage:            lastMessage,
				LastStartedTime:        lastStartedTime,
				Message:                message,
				Namespace:              pod.Metadata.Namespace,
				StartedTime:            startedTime,
			}
			containerStatusInfoList = append(containerStatusInfoList, &containerStatusInfo)
		}
	}

	var reqCpu, reqMemory, reqGpu, limitCpu, limitMemory, limitGpu string
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Requests["cpu"]; ok {
		reqCpu = v
	}
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Requests["memory"]; ok {
		reqMemory = v
	}
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Requests["nvidia.com/gpu"]; ok {
		reqGpu = v
	}
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Limits["cpu"]; ok {
		limitCpu = v
	}
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Limits["memory"]; ok {
		limitMemory = v
	}
	if v, ok := nb.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"]; ok {
		limitGpu = v
	}
	//logger.Logger().Debugf("GetNamespacedNotebookStatus Containers[0].Resources: %+v\n", *nb.Spec.Template.Spec.Containers[0].Resources)

	resp := models.GetNotebookStatusResponse{
		ContainersStatusInfo: containerStatusInfoList,
		CreateTime:           createTime,
		Image:                nb.Spec.Template.Spec.Containers[0].Image,
		Name:                 nb.Metadata.Name,
		Namespace:            nb.Metadata.Namespace,
		NotebookStatus:       status,
		NotebookStatusInfo:   notebookStatusInfo,
		ResourceReq: &models.Resource{
			CPU:    reqCpu,
			Gpu:    reqGpu,
			Memory: reqMemory,
		},
		ResourceLimit: &models.Resource{
			CPU:    limitCpu,
			Gpu:    limitGpu,
			Memory: limitMemory,
		},
	}

	return buildResponse(200, resp, nil)
}

func GetNamespacedNotebookLog(params operations.GetNamespacedNotebookLogParams) middleware.Responder {
	// 1.参数校验
	if err := checkGetNamespacedNotebookLogParams(&params); err != nil {
		logger.Logger().Errorf("GetNamespacedNotebookLog: fail to check param, error: %s", err.Error())
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	// 2.获取notebook
	k8sNotebook, res, err := utils.GetNotebook(params.Namespace, params.NotebookName)
	if err != nil {
		logger.Logger().Error("GetNamespacedNotebookLog From K8s Error:%v", err.Error())
		if res != nil {
			if res.StatusCode == http.StatusNotFound {
				logger.Logger().Errorf("GetNamespacedNotebookLog: Fail to get notebook by namespace %q and name %q",
					params.Namespace, params.NotebookName)
				return buildResponse(http.StatusNotFound, nil, errors.New("notebook not found"))
			}
		}
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	nb, err := parseNotebookFromK8s(k8sNotebook)
	if err != nil {
		logger.Logger().Error("parseNotebookFromK8s From K8s Error:%v", err.Error())
		return buildResponse(http.StatusInternalServerError, nil, err)
	}

	// 3.从es查询日志
	podName := fmt.Sprintf("%s-0", nb.Metadata.Name)
	esResult, err := es.LogList(es.Log{
		EsIndex: es.K8S_LOG_CONTAINER_ES_INDEX,
		EsType:  es.K8S_LOG_CONTAINER_ES_TYPE,
		PodName: podName,
		Ns:      params.Namespace,
		From:    int((*params.CurrentPage - 1) * (*params.PageSize)),
		Size:    int(*params.PageSize),
		Asc:     *params.Asc,
	})

	if err != nil {
		logger.Logger().Errorf("fail to get log info list, error: %s", err.Error())
		return buildResponse(http.StatusInternalServerError, nil, err)
	}
	resp, err := generateLogList(esResult)

	return buildResponse(http.StatusOK, resp, err)
}

func parseNotebookFromK8s(notebook interface{}) (*models.K8sNotebook, error) {
	var notebookFromK8s models.K8sNotebook
	bytes, err := json.Marshal(notebook)
	if err != nil {
		logger.Logger().Errorf("Marshal failed. %v", err.Error())
		return nil, err
	}
	err = json.Unmarshal(bytes, &notebookFromK8s)
	if err != nil {
		logger.Logger().Errorf("parse []models.K8sNotebook Unmarshal failed. %v", err.Error())
		return nil, err
	}
	// logger.Logger().Debugf("parseNotebookFromK8s notebook: %+v\n", notebookFromK8s)
	return &notebookFromK8s, nil
}

func buildResponse(statusCode int, result interface{}, err error) middleware.Responder {
	if err != nil {
		if statusCode == http.StatusNotFound {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusNotFound)
				payload, _ := json.Marshal(models.Error{
					Error:   "Not found",
					Code:    http.StatusNotFound,
					Message: err.Error(),
				})
				w.Write(payload)
			})

		} else {
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
	}

	response, err := json.Marshal(result)
	if err != nil {
		if statusCode == http.StatusNotFound {
			return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
				w.WriteHeader(http.StatusNotFound)
				payload, _ := json.Marshal(models.Error{
					Error:   "Not found",
					Code:    http.StatusNotFound,
					Message: err.Error(),
				})
				w.Write(payload)
			})

		} else {
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
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(statusCode)
		payload, _ := json.Marshal(models.Result{
			Code:    fmt.Sprintf("%d", http.StatusOK),
			Message: "success",
			Result:  response,
		})
		_, _ = w.Write(payload)
	})
}

func generateLogList(r map[string]interface{}) (models.GetNotebookLogResponse, error) {
	rst := models.GetNotebookLogResponse{}
	if len(r) == 0 || r["hits"] == nil {
		logger.Logger().Errorf("es hits is nil, r: %+v", r)
		err := fmt.Errorf("es hits is nil, r: %+v", r)
		return rst, err
	}

	switch r["hits"].(type) {
	case map[string]interface{}:
		if len(r["hits"].(map[string]interface{})) == 0 {
			logger.Logger().Info("es hits type is map[string]interface{}, but map is empty")
			return rst, nil
		}
	default:
		logger.Logger().Errorf("hits type is undefined, type by  reflect is %v", reflect.TypeOf(r["hits"]).Kind().String())
		err := fmt.Errorf("hits type is undefined, type by  reflect is %v", reflect.TypeOf(r["hits"]).Kind().String())
		return rst, err
	}

	notebookLogList := make([]*models.NotebookLog, 0)
	var total int64 = 0
	if hitsTotalInter, ok := r["hits"].(map[string]interface{})["total"]; ok {
		var hitsTotal int64 = 0
		switch hitsTotalVal := hitsTotalInter.(type) {
		case int:
			hitsTotal = int64(hitsTotalVal)
		case int64:
			hitsTotal = hitsTotalVal
		case float32:
			hitsTotal = int64(hitsTotalVal)
		case float64:
			hitsTotal = int64(hitsTotalVal)
		case string:
			hitsTotalInt64, err := strconv.ParseInt(hitsTotalVal, 10, 64)
			if err != nil {
				logger.Logger().Infof("hits total type is string, but parse error: %v", err)
			}
			hitsTotal = hitsTotalInt64
		case map[string]interface{}:
			if valueInter, exist := hitsTotalVal["value"]; exist {
				if totalFloat64, existFloat64 := valueInter.(float64); existFloat64 {
					hitsTotal = int64(totalFloat64)
				}
			}
		default:
			logger.Logger().Infof("hits total type undefined, %v", reflect.TypeOf(hitsTotalInter).Kind().String())
		}
		total = hitsTotal
	}
	if _, ok := r["hits"].(map[string]interface{})["hits"]; ok {
		if hits, ok2 := r["hits"].(map[string]interface{})["hits"].([]interface{}); ok2 {
			for _, hit := range hits {
				var log interface{} = ""
				switch logV := hit.(type) {
				case map[string]interface{}:
					if source, ok := logV["_source"]; ok {
						sourceSub, sourceOK := source.(map[string]interface{})
						if !sourceOK && len(sourceSub) == 0 {
							continue
						}
						if v, ok2 := sourceSub["log"]; ok2 {
							log = v
						}
					}
					notebookLog := models.NotebookLog{
						Log: log.(string),
					}
					notebookLogList = append(notebookLogList, &notebookLog)
				default:
					logger.Logger().Infof("hit type undefined, hit type: %v", reflect.TypeOf(logV).Kind().String())
					continue
				}
			}
		}
	}

	rst.LogList = notebookLogList
	rst.Total = &total

	return rst, nil
}

func checkGetNamespacedNotebookLogParams(param *operations.GetNamespacedNotebookLogParams) error {
	if param.Namespace == "" || param.NotebookName == "" {
		return fmt.Errorf("namespace %q or notebook %q can not be empty string", param.Namespace, param.NotebookName)
	}
	if param.CurrentPage == nil || *param.CurrentPage == 0 {
		var currentPage int64 = 1
		param.CurrentPage = &currentPage
	}
	if param.PageSize == nil || *param.PageSize == 0 {
		var pageSize int64 = 10
		param.CurrentPage = &pageSize
	}

	asc := false
	if *(param.Asc) {
		asc = true
	}
	param.Asc = &asc

	return nil
}

func StopNamespacedNotebook(params operations.StopNotebookByIDParams) middleware.Responder {
	// 1.获取notebook
	nbInMongo, err := mongo.GetNotebookById(params.ID)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	namespace := nbInMongo.Namespace
	name := nbInMongo.Name
	currentUserId := params.HTTPRequest.Header.Get(mw.CcAuthUser)
	logger.Logger().Infof("delete Namespace(%s) Notebook %s.", namespace, name)
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	// 2. check current user accessable for namespace
	body, err := client.CheckNamespace(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
	if err != nil {
		logger.Logger().Errorf("CheckNamespace failed: %s, token: %s, namespace: %s", err.Error(),
			params.HTTPRequest.Header.Get(ccClient.CcAuthToken), namespace)
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

	// 3.get list from k8s
	ncClient := utils.GetNBCClient()
	listFromK8s, err := ncClient.GetNotebooks(namespace, "", "")
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "getNamespacedNotebooks failed",
		})
	}

	// 4.parse obj from k8s to models.NotebookFromK8s
	notebookFromK8s, err := ParseNotebooksFromK8s(listFromK8s)
	if err != nil {
		logger.Logger().Errorf("parseNotebookFromK8s Unmarshal failed. %v", err.Error())
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
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
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "parse resultMap Unmarshal failed.",
		})
	}
	//if current user is GU, he only access his own notebook
	if role == "GU" {
		for _, k := range res {
			if k.Name == name && nbInMongo.User != currentUserId {
				logger.Logger().Errorf("GU only access his own notebook.")
				return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
					Code:    http.StatusNotFound,
					Error:   "you can't access this notebook",
					Message: "GU only access his own notebook.",
				})
			}
		}
	}

	// 5.更新mongo notebook 状态信息
	updateMap := bson.M{
		"status": constants.NoteBookStatusStoping,
	}
	err = mongo.UpdateNotebookById(params.ID, updateMap)
	if err != nil {
		logger.Logger().Errorf("fail to stop notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 6. delete nb
	err = ncClient.DeleteNotebook(name, namespace)
	if err != nil {
		return operations.NewPostNamespacedNotebookNotFound().WithPayload(&models.Error{
			Code:    http.StatusNotFound,
			Error:   err.Error(),
			Message: "stop notebook failed",
		})
	}

	// 7.再次更新mongo notebook 状态信息
	stopedMap := bson.M{
		"status": constants.NoteBookStatusTerminating,
	}
	err = mongo.UpdateNotebookById(params.ID, stopedMap)
	if err != nil {
		logger.Logger().Errorf("fail to update notebook status to %q by id from mongo, "+
			"id: %s, err: %v", constants.NoteBookStatusStoped, params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 7.1 update status
	go updateNotebookStatus(params.ID, nbInMongo)

	//return operations.NewDeleteNamespacedNotebookOK()
	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

func updateNotebookStatus(id string, nbInMongo utils.NotebookInMongo) {
	var (
		status string
		err    error
	)
	maxRetryNum := 60 //The maximum time for waiting for pod is 5 minutes. Each cycle takes 5 seconds and 60 cycles
	ncClient := utils.GetNBCClient()
	for count := 0; count < maxRetryNum; count++ {
		status, err = ncClient.GetPodStatus(nbInMongo.Namespace, nbInMongo.Name)
		if status == "" || err != nil {
			status = constants.NoteBookStatusStoped
			break
		}

		if count == maxRetryNum-1 {
			status = constants.NoteBookStatusError
			break
		}

		time.Sleep(time.Second * 5)
	}

	stopedMap := bson.M{

		"status": status,
	}

	err = mongo.UpdateNotebookById(id, stopedMap)
	if err != nil {
		logger.Logger().Errorf("updateNotebookStatus fail to update notebook status to %q by id from mongo, "+
			"id: %s, err: %v", constants.NoteBookStatusStoped, id, err)
	}
}

func StartNamespacedNotebook(params operations.StartNotebookByIDParams) middleware.Responder {
	// 1.获取mongo notebook 信息
	nbInMongo, err := mongo.GetNotebookById(params.ID)
	if err != nil {
		logger.Logger().Errorf("fail to get notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	//query := bson.M{"namespace": nbInMongo.Namespace, "name": nbInMongo.Name, "enable_flag": 1}
	//count, err := mongo.CountNoteBook(query)
	//if err != nil {
	//	logger.Logger().Errorf("fail to check notebook existed in mongo: %s", err.Error())
	//	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		payload, _ := json.Marshal(models.Error{
	//			Error:   "StatusInternalServerError",
	//			Code:    http.StatusInternalServerError,
	//			Message: err.Error(),
	//		})
	//		w.Write(payload)
	//	})
	//}
	//if count > 0 {
	//	msg := "Notebook在该命名空间下已存在，无法创建Notebook"
	//	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
	//		w.WriteHeader(http.StatusInternalServerError)
	//		payload, _ := json.Marshal(models.Error{
	//			Error:   "forbidden",
	//			Code:    http.StatusInternalServerError,
	//			Message: msg,
	//		})
	//		w.Write(payload)
	//	})
	//}

	// 2. 校验k8s notebook
	logger.Logger().Infof("create time now: %v", time.Now().Format(time.RFC1123))
	client := ccClient.GetCcClient(viper.GetString(cfg.CCAddress))
	ncClient := utils.GetNBCClient()
	existed, err := ncClient.CheckNotebookExisted(nbInMongo.Namespace, nbInMongo.Name)
	if err != nil {
		logger.Logger().Errorf("Check Notebook existed error: ", err.Error())
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
	if existed {
		msg := "Notebook在该命名空间下已存在，无法创建Notebook"
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "forbidden",
				Code:    http.StatusInternalServerError,
				Message: msg,
			})
			w.Write(payload)
		})
	}

	// 3.校验用户权限
	//Get User Message
	currentUserId := params.HTTPRequest.Header.Get(mw.UserIDHeader)
	containerUser := currentUserId
	var gid, uid, token *string
	if nbInMongo.ProxyUser != "" {
		containerUser = nbInMongo.ProxyUser
		gid, uid, token, err = client.GetProxyUserGUIDFromUser(params.HTTPRequest.Header.Get(ccClient.CcAuthToken),
			currentUserId, containerUser)
	} else {
		gid, uid, token, err = client.GetGUIDFromUserId(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId)
	}
	logger.Logger().Infof("gid: %s, uid: %s", *gid, *uid)
	if err != nil || gid == nil || uid == nil {
		logger.Logger().Errorf("GetGUIDFromUserId failed: %s, token: %s, userName: %s, id: %v",
			err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, params.ID)
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

	// Check Namespace
	logger.Logger().Infof("create notebook in namespace %s", nbInMongo.Namespace)
	err = client.UserNamespaceCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, nbInMongo.Namespace)
	if err != nil {
		logger.Logger().Errorf("UserNamespaceCheck failed: %s, token: %s, userName: %s, namespace: %s",
			err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, nbInMongo.Namespace)
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

	// 4.生成notebook  request
	nbReq, err := generateNootbookRequest(nbInMongo)
	if err != nil {
		logger.Logger().Errorf("fail to generate notebook request, nbInMongo: %+v, err: %v", nbInMongo, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "error",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	logger.Logger().Debugf("notebook start request: %+v ", *nbReq)
	if nbReq.DataVolume != nil {
		logger.Logger().Debugf("notebook start request dataVolume localpath: %v, SubPath: %v, MountPath: %v", *nbReq.DataVolume.LocalPath,
			nbReq.DataVolume.SubPath, *nbReq.DataVolume.MountPath)
	}
	if nbReq.WorkspaceVolume != nil {
		logger.Logger().Debugf("notebook start request workVolume localpath: %v, SubPath: %v, MountPath: %v", *nbReq.WorkspaceVolume.LocalPath,
			nbReq.WorkspaceVolume.SubPath, *nbReq.WorkspaceVolume.MountPath)
	}

	if nbReq.Image.ImageType == nil || *nbReq.Image.ImageType == "" {
		standard := "Standard"
		nbReq.Image.ImageType = &standard
	}
	req := operations.PostNamespacedNotebookParams{
		HTTPRequest: params.HTTPRequest,
		Namespace:   nbInMongo.Namespace,
		Notebook:    nbReq,
	}

	// 5.check notebook request
	//err = client.PostNotebookRequestCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), params.Notebook)
	err = utils.CheckNotebook(req)
	if err != nil {
		logger.Logger().Errorf("PostNotebookRequestCheck failed: %s, token: %s, userName: %s,"+
			" id: %v", err.Error(), params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, params.ID)
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

	// 6.更新mongo notebook 状态
	updateMap := bson.M{
		"status": constants.NoteBookStatusCreating,
	}
	err = mongo.UpdateNotebookById(params.ID, updateMap)
	if err != nil {
		logger.Logger().Errorf("fail to start notebook by id from mongo, "+
			"id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusForbidden,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}
	if nbInMongo.Namespace == "" || nbInMongo.Name == "" {
		err := fmt.Errorf("namespace %q or name %q can't be empty", nbInMongo.Namespace, nbInMongo.Name)
		logger.Logger().Errorf("StartNamespacedNotebook params error, id: %s, err: %v", params.ID, err)
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			w.WriteHeader(http.StatusInternalServerError)
			payload, _ := json.Marshal(models.Error{
				Error:   "failed",
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			w.Write(payload)
		})
	}

	// 7. 创建k8s notebook
	nbClient := utils.GetNBCClient()
	maxRetryNum := 5
	sparkBytes := make([]byte, 0)
	for count := 0; count < maxRetryNum; count++ {
		sparkBytes, err = nbClient.CreateNotebook(req.Notebook, currentUserId, *gid, *uid, *token)
		if err != nil {
			if !stringsUtil.Contains(err.Error(), "provided port is already allocated") {
				logger.Logger().Errorf("CreateNotebook failed: %s, Code: %s", err.Error())
				var error = models.Error{
					Code:    http.StatusInternalServerError,
					Error:   "create notebook failed",
					Message: err.Error(),
				}
				return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
					payload, _ := json.Marshal(error)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(payload)
				})
			}

			//端口被占用，开始重试
			logger.Logger().Infof("[%d]provided port is already allocated, start retry ...\n", count)
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}

	// 8.创建k8s config map
	err = nbClient.CreateYarnResourceConfigMap(req.Notebook, containerUser, sparkBytes)
	if err != nil {
		logger.Logger().Errorf("CreateNotebook configmap failed: %s, Code: %s", err.Error())
		var error = models.Error{
			Code:    http.StatusInternalServerError,
			Error:   "create notebook configmap failed",
			Message: err.Error(),
		}
		return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
			payload, _ := json.Marshal(error)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(payload)
		})
	}

	var result = models.Result{
		Code:    fmt.Sprintf("%d", http.StatusOK),
		Message: "success",
		Result:  nil,
	}
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		payload, _ := json.Marshal(result)
		w.Write(payload)
		//w.WriteHeader(200)
	})
}

func generateNootbookRequest(nb utils.NotebookInMongo) (*models.NewNotebookRequest, error) {
	cpu, err := strconv.ParseFloat(nb.Cpu, 64)
	if err != nil {
		return nil, err
	}

	image := models.Image{
		ImageType: &nb.Image.ImageType,
		ImageName: &nb.Image.ImageName,
	}
	memory := models.Memory{
		MemoryAmount: &nb.Memory.MemoryAmount,
		MemoryUnit:   nb.Memory.MemoryUnit,
	}

	workspaceVolumeSize := int64(nb.WorkspaceVolume.Size)
	workspaceVolume := models.MountInfo{
		AccessMode: &nb.WorkspaceVolume.AccessMode,
		LocalPath:  &nb.WorkspaceVolume.LocalPath,
		MountPath:  &nb.WorkspaceVolume.MountPath,
		MountType:  &nb.WorkspaceVolume.MountType,
		Size:       &workspaceVolumeSize,
		SubPath:    nb.WorkspaceVolume.SubPath,
	}

	req := models.NewNotebookRequest{
		CPU:             &cpu,
		DataVolume:      nil,
		DriverMemory:    nb.DriverMemory,
		ExecutorCores:   nb.ExecutorCores,
		ExecutorMemory:  nb.ExecutorMemory,
		Executors:       nb.Executors,
		ExtraResources:  nb.ExtraResources,
		Image:           &image,
		Memory:          &memory,
		Name:            &nb.Name,
		Namespace:       &nb.Namespace,
		ProxyUser:       nb.ProxyUser,
		Queue:           nb.Queue,
		SparkSessionNum: int64(nb.SparkSessionNum),
		WorkspaceVolume: &workspaceVolume,
	}

	if len(nb.DataVolume) > 0 {
		size := int64(nb.DataVolume[0].Size)
		dataVolumeObj := models.MountInfo{
			MountType:  &nb.DataVolume[0].MountType,
			LocalPath:  &nb.DataVolume[0].LocalPath,
			SubPath:    nb.DataVolume[0].SubPath,
			Size:       &size,
			MountPath:  &nb.DataVolume[0].MountPath,
			AccessMode: &nb.DataVolume[0].AccessMode,
		}
		req.DataVolume = &dataVolumeObj
	}
	return &req, nil
}
