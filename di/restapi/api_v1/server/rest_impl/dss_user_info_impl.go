package rest_impl

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"net/http"
	"webank/DI/pkg/client/dss"
	"webank/DI/restapi/api_v1/server/operations/dss_user_info"
)

func DssGetUserInfo(params dss_user_info.GetDssUserInfoParams) middleware.Responder  {
	operation := "DssGetUserInfo"
	dssClient := dss.GetDSSClient()
	cookies := []*http.Cookie{}
	test,_ := ioutil.ReadAll(params.HTTPRequest.Body)
	print(string(test))
	ticketIdCookies := http.Cookie{
		Name: "bdp-user-ticket-id",
		Value: params.ID,
	}
	workspaceIdCookies := http.Cookie{
		Name: "workspaceId",
		Value: params.WorkspaceID,
	}
	cookies = append(cookies, &ticketIdCookies, &workspaceIdCookies)
	res, err := dssClient.GetUserInfo(cookies)
	if err != nil{
		log.Println("get dss user info error, ",err.Error(), "ticket id is ", params.ID)
		return httpResponseHandle(http.StatusInternalServerError, err, operation,[]byte(err.Error()))
	}
	marshal, err := json.Marshal(&res)
	if err != nil{
		log.Println("get dss user info error, ",err)
		return httpResponseHandle(http.StatusInternalServerError, err, operation, marshal)
	}
	return httpResponseHandle(http.StatusOK, err, operation, marshal)
}
