package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"webank/DI/commons/config"
	"webank/DI/commons/models"
)

type CcClient struct {
	ccUrl      string
	httpClient http.Client
}
type UserResponse struct {
	Name      string `json:"name"`
	GID       int    `json:"gid"`
	UID       int    `json:"uid"`
	Type      string `json:"type"`
	GUIDCheck bool   `json:"guidCheck"`
	Remarks   string `json:"remarks"`
	Token     string `json:"token,omitempty"`
}

type ProxyUserResponse struct {
	Name      string `json:"name"`
	GID       int    `json:"gid"`
	UID       int    `json:"uid"`
	Token     string `json:"token,omitempty"`
	UserId    int `json:"user_id,omitempty"`
	Path      string `json:"path,omitempty"`
}

type ProxyUserCheckResponse struct {
	Code    string `json:"code"`
	Message int    `json:"message"`
	Result  bool   `json:"result"`
}

const CcAuthToken = "MLSS-Token"

// FIXME MLSS Change: get models filter by username and namespace
const CcSuperadmin = "MLSS-Superadmin"

var once sync.Once
var client *CcClient

func GetCcClient(ccUrl string) *CcClient {
	once.Do(func() {
		client = &CcClient{}
		//client.httpClient = http.Client{}
		client.ccUrl = ccUrl

		//FIXME
		pool := x509.NewCertPool()

		caCertPath := config.GetValue(config.CaCertPath)
		log.Debugf("ccClient to get caCertPath: %v", caCertPath)
		//调用ca.crt文件
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		pool.AppendCertsFromPEM(caCrt)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            pool,
			},
		}

		client.httpClient = http.Client{Transport: tr}

	})
	return client
}

// FIXME MLSS Change: get models filter by username and namespace
func (cc *CcClient) AuthAccessCheck(authOptions *map[string]string) (int, string, string, error) {
	apiUrl := cc.ccUrl + "/cc/v1/inter/user"
	//apiUrl := cc.ccUrl + "/cc/v1/sample"

	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		log.Debug("Create AuthAccessCheck Request Failed !")
		log.Debug(err.Error())
		return 500, "", "", err
	}
	for key, value := range *authOptions {
		req.Header.Set(key, value)
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Debug("Do AuthAccessCheck Request Failed !")
		log.Debug(err.Error())
		return 500, "", "", err
	}
	defer resp.Body.Close()
	token := resp.Header.Get(CcAuthToken)
	isSA := resp.Header.Get(CcSuperadmin)
	statusCode := resp.StatusCode
	if statusCode == 401 || (statusCode >= 200 && statusCode < 300 && token == "") {
		log.Debugf("Auth Access Check Request Failed, StatusCode: %v, Token: %v", statusCode, token)
		log.Debugf("Message: %v", resp.Body)
		return 401, "", "", *new(error)
	}
	if statusCode > 300 || token == "" {
		log.Debugf("Auth Access Check Request Failed from cc, StatusCode: %v, Token: %v", statusCode, token)
		log.Debugf("Message: %v", resp.Body)
		return 500, "", "", *new(error)
	}
	return statusCode, token, isSA, nil
}

func (cc *CcClient) UserStorageCheck(token string, username string, path string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/storages?path=%v", username, path)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CcClient) UserNamespaceCheck(token string, userName string, namespace string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/namespaces/%v", userName, namespace)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New("user cannot access the specified namespace")
	}
	return nil
}

func (cc *CcClient) GetProxyUserFromCC(token string, proxyUserId string, userId string) (*ProxyUserResponse, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/proxyUser/%v/user/%v", proxyUserId, userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, err
	}
	var u = new(ProxyUserResponse)
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &u)
	if err != nil {
		return nil,  errors.New("GetResultData failed" + err.Error())
	}
	log.Debugf("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)
	if  (u.GID == 0 || u.UID == 0) {
		return nil, errors.New("error GUIDCheck is on and GID or UID is 0")
	}
	//gid := strconv.Itoa(u.GID)
	//uid := strconv.Itoa(u.UID)

	return u, nil
}

func (cc *CcClient) GetGUIDFromProxyUserId(token string, proxyUserId string) (*string, *string, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/proxyUser/%v", proxyUserId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, nil, err
	}
	var u = new(UserResponse)
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &u)
	if err != nil {
		return nil, nil, errors.New("GetResultData failed")
	}
	log.Debugf("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)
	if u.GUIDCheck && (u.GID == 0 || u.UID == 0) {
		return nil, nil, errors.New("error GUIDCheck is on and GID or UID is 0")
	}
	gid := strconv.Itoa(u.GID)
	uid := strconv.Itoa(u.UID)

	return &gid, &uid, nil
}

func (cc *CcClient) GetGUIDFromUserId(token string, userId string) (*string, *string, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/users/name/%v", userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, nil, err
	}
	var u = new(UserResponse)
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &u)
	if err != nil {
		return nil, nil, errors.New("GetResultData failed")
	}
	log.Debugf("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)
	if u.GUIDCheck && (u.GID == 0 || u.UID == 0) {
		return nil, nil, errors.New("error GUIDCheck is on and GID or UID is 0")
	}
	gid := strconv.Itoa(u.GID)
	uid := strconv.Itoa(u.UID)

	return &gid, &uid, nil
}


func (cc *CcClient) GetUserByName(token string, userId string) (*UserResponse, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/users/name/%v", userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return nil, err
	}
	var u = new(UserResponse)
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &u)
	if err != nil {
		return nil, errors.New("GetResultData failed")
	}
	log.Debugf("response=%+v, GID=%v, UID=%v", u, u.GID, u.UID)

	return u, nil
}

func (cc *CcClient) ProxyUserCheck(token string, userId string, proxyUser string) (bool, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/proxyUserCheck/%v/%v", proxyUser, userId)
	body, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return false, err
	}
	var auth bool
	//json.Unmarshal(body, &u)
	err = models.GetResultData(body, &auth)

	if err != nil {
		log.Error("Unmarshal Error", err)
		return false, errors.New("Check Proxy User failed")
	}

	return auth, nil
}

func (cc *CcClient) AdminUserCheck(token string, adminUserId string, userId string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/users/%v", adminUserId, userId)
	_, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CcClient) UserStoragePathCheck(token string, username string, namespace string, path string) (result []byte, error error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/users/%v/namespaces/%v/storages?path=%v", username, namespace, path)
	res, err := cc.accessCheckFromCc(token, apiUrl)
	if err != nil {
		error = err
		return res, err
	}
	result = res
	return res, err
}

func (cc *CcClient) accessCheckFromCc(token string, apiUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, doErr := client.httpClient.Do(req)
	statusCode := resp.StatusCode
	body, bodyErr := ioutil.ReadAll(resp.Body)
	if doErr != nil {
		return body, doErr
	}
	defer resp.Body.Close()
	//statusCode := resp.StatusCode
	if statusCode != 200 {
		return body, errors.New("error from cc ")
	}
	//body, err := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		return body, bodyErr
	}
	return body, nil
}

func (cc *CcClient) AlertTraining(alertDTO models.AlertDTO) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/alerts/training")
	bytes, err := json.Marshal(alertDTO)
	if err != nil {
		return nil, err
	}
	//init request
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}
	return body, nil
}

func (cc *CcClient) CheckNamespace(token string, namespace string) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/namespaces/%v", namespace)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}

	defer resp.Body.Close()
	return body, nil
}

func (cc *CcClient) CheckNamespaceUser(token string, namespace string, user string) error {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/auth/access/admin/namespaces/%v/users/%v", namespace, user)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return errors.New(string(body))
	}
	return nil
}

func (cc *CcClient) GetCurrentUserNamespaceWithRole(token string, roleId string) ([]byte, error) {
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/groups/users/roles/%v/namespaces", roleId)
	log.Debugf("get apiUrl: %v", apiUrl)
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set(CcAuthToken, token)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	statusCode := resp.StatusCode
	//	if statusCode == 401 || statusCode == 500 {
	if statusCode != 200 {
		//		return *new(error)
		return nil, errors.New(string(body))
	}

	defer resp.Body.Close()
	return body, nil
}


type GetGroupFromUserRequest struct {

}

type GetGroupFromUserResponse struct {
	ID         int64  `gorm:"column:id; PRIMARY_KEY" json:"id"`
	Name       string `json:"name"`
	EnableFlag int8   `json:"enable_flag"`
	UserId     int64  `json:"user_id"`
	RoleId     int64  `json:"role_id"`
	GroupId    int64  `json:"group_id"`
	Remarks    string `json:"remarks"`
}



func (cc *CcClient) GetGroupFromUser(username string) (*GetGroupFromUserResponse,error){
	apiUrl := cc.ccUrl + fmt.Sprintf("/cc/v1/group/username/%v",  username)

	req,err := http.NewRequest("GET",apiUrl,nil)
	if err != nil {
		return nil,err
	}

	res,err := cc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	group := GetGroupFromUserResponse{}
	err = json.Unmarshal(resByte,&group)
	if err != nil {
		return nil, err
	}

	return &group,nil


}