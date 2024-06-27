package rest_impl

import (
	"context"
	"encoding/json"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	commonModels "webank/DI/commons/models"
	"webank/DI/pkg/v2/util"
	"webank/DI/restapi/api_v2/restmodels"
	"webank/DI/restapi/api_v2/server/operations/experiment"
	"webank/DI/restapi/service_v2"
)

func StringPtr(input string) *string {
	return &input
}

func GetUserID(r *http.Request) string {
	return r.Header.Get(config.CcAuthUser)
}

func GetClusterType(r *http.Request) string {
	return r.Header.Get(config.CcCluster)
}

type Result struct {
	Code      string          `json:"code"`
	Msg       string          `json:"msg"`
	RequestId string          `json:"requestId"`
	Result    json.RawMessage `json:"result"`
}

// 如果 err 为空，就表示没有错误，就返回 200，success，正常result数据
// 如果 err 不为空，就表示有错误，返回非200， msg根据具体错误返回具体的， result为空
func HttpResponseHandle2(ctx context.Context, err error, returnResult interface{}) middleware.Responder {
	result := Result{}
	if v, ok := ctx.Value("X-Request-Id").(string); ok {
		result.RequestId = v
	}
	var statusCode int
	if err == nil {
		statusCode = http.StatusOK
		result.Code = strconv.Itoa(statusCode)
		result.Msg = "success"
		marshal, err := json.Marshal(returnResult)
		if err != nil {
			log.WithError(err).Errorf("json.Marshal failed")
			statusCode = http.StatusInternalServerError
			result.Code = strconv.Itoa(statusCode)
			result.Msg = service_v2.InternalError{}.Error()
		} else {
			result.Result = marshal
		}
	} else {
		switch err.(type) {
		case service_v2.InputParamsError:
			statusCode = http.StatusBadRequest
		case service_v2.AccessError:
			statusCode = http.StatusForbidden
		default:
			statusCode = http.StatusInternalServerError
			err = service_v2.InternalError{}
		}
		result.Code = strconv.Itoa(statusCode)
		result.Msg = err.Error()
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(statusCode)
		response, err := json.Marshal(result)
		_, err = w.Write(response)
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	})
}

// HttpResponseHandle copy from restapi/api_v1/server/rest_impl/models_impl.go httpResponseHandle
// todo(gaoyuanhe): 该函数逻辑有点混乱，需要改造
func HttpResponseHandle(status int, err error, operation string, resultMsg []byte) middleware.Responder {
	result := commonModels.Result{
		Code:   "200",
		Msg:    "success",
		Result: json.RawMessage(resultMsg),
	}

	if resultMsg == nil {
		jsonStr, err := json.Marshal(operation + " " + "success")
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
		result.Result = jsonStr
	}

	if status != http.StatusOK {
		result.Code = strconv.Itoa(status)
		result.Msg = "Error"
		jsonStr, _ := json.Marshal(operation + " " + " Error")
		result.Result = jsonStr
		if err != nil {
			logger.GetLogger().Error("Operation: " + operation + ";" + "Error: " + err.Error())
			jsonStr, _ := json.Marshal(operation + " " + " error: " + err.Error())
			result.Result = jsonStr
			// 前端应该是获取Msg字段显示错误信息，如果只是简单的给Error，用户完全不知道是什么错误
			result.Msg = err.Error()
		}
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		w.WriteHeader(status)
		response, err := json.Marshal(result)
		_, err = w.Write(response)
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	})
}

func validateCreateExperiment(params experiment.CreateExperimentParams) error {
	// todo(gaoyuanhe): source_system有默认值就是指针？但是没有传值又会是空指针？
	if params.Experiment.SourceSystem == nil || *params.Experiment.SourceSystem == restmodels.CreateExperimentRequestSourceSystemMLSS {
		if params.Experiment.GroupID == "" {
			return service_v2.InputParamsError{Field: "GroupID", Reason: "来源类型为MLSS的时候，GroupID不能为空"}
		}
	}
	return nil
}

func GenerateContext(r *http.Request) context.Context {
	var reqID string
	if reqID = r.Header.Get("X-Request-Id"); reqID == "" {
		// 如果调用方的header中没有，那就自己生成一个
		reqID = util.CreateUUID()
		r.Header.Set("X-Request-Id", reqID)
	}
	ctx := context.Background()
	return context.WithValue(ctx, "X-Request-Id", reqID)
}
