package rest_impl

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	logr "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"webank/DI/pkg/client/dss"
	"webank/DI/restapi/api_v1/server/operations/dss_user_info"
)

func DssGetUserInfo(params dss_user_info.GetDssUserInfoParams) middleware.Responder {
	operation := "DssGetUserInfo"
	dssClient := dss.GetDSSClient()
	cookies := []*http.Cookie{}
	test, _ := ioutil.ReadAll(params.HTTPRequest.Body)
	print(string(test))
	ticket, err := url.QueryUnescape(params.ID)
	if err != nil {
		logr.Error(operation + " error: " + err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte(err.Error()))
	}
	ticketIdCookies := http.Cookie{
		Name:  "linkis_user_session_ticket_id_v1",
		Value: ticket,
	}
	workspaceIdCookies := http.Cookie{
		Name:  "workspaceId",
		Value: params.WorkspaceID,
	}
	cookies = append(cookies, &ticketIdCookies, &workspaceIdCookies)

	log.Infof("DssGetUserInfo cookies: %+v", cookies)

	res, err := dssClient.GetUserInfo(cookies)
	if err != nil {
		log.Errorf("get dss user info error, ", err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, []byte("DSS获取用户信息失败"))
	}

	marshal, err := json.Marshal(&res)
	if err != nil {
		log.Errorf("get dss user info error, ", err.Error())
		return httpResponseHandle(http.StatusInternalServerError, err, operation, marshal)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}
