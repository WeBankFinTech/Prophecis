/*
 * Copyright 2017-2018 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
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

package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"webank/DI/restapi/api_v1/client/models"

	"github.com/go-openapi/runtime"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"

	"path/filepath"
	dlaasClient "webank/DI/restapi/api_v1/client"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin"
	"github.com/urfave/cli"
)

// TrainCmd is the struct to deploy a model.
type TrainCmd struct {
	ui      terminal.UI
	config  plugin.PluginConfig
	context plugin.PluginContext
}

// NewTrainCmd is used to deploy a model's.
func NewTrainCmd(ui terminal.UI, context plugin.PluginContext) *TrainCmd {
	return &TrainCmd{
		ui:      ui,
		context: context,
	}
}

var jobStatus = "not_started"
var currentLine = 0

const maxRetryTimeForGetState = 10
const maxRetryTimeForGetLog = 10

func (cmd *TrainCmd) buildTempFile() (*os.File, string) {
	filename := fmt.Sprintf("%v.yml", uuid.New())
	tmpYamlFile, err := os.Create(fmt.Sprintf("%v", filename))
	if err != nil {
		cmd.ui.Failed("Create tmpYaml file failed, ", err)
		os.Exit(1)
	}
	return tmpYamlFile, filename
}

// Run is the handler for the model-deploy CLI command.
func (cmd *TrainCmd) Run(cliContext *cli.Context) error {
	iMonitoring := true
	cmd.config = cmd.context.PluginConfig()

	args := cliContext.Args()
	if len(args) == 0 {
		cmd.ui.Failed("Incorrect number of arguments.")
	}
	params := models.NewPostModelParams().WithTimeout(defaultOpTimeout)
	_, manifestFileName := filepath.Split(args[0])
	manifestFile := args[0]
	modelDefinitionFile := ""
	modelDefinitionExisted := 0

	//If Command contains model definition file
	if len(args)>1 && !strings.HasPrefix(args[1],"-"){
		modelDefinitionExisted = 1
		if strings.Contains(args[1], ".zip") {
			modelDefinitionFile = args[1]
			cmd.ui.Say("Deploying model with manifest '%s' and model file '%s'...", terminal.EntityNameColor(manifestFileName), terminal.EntityNameColor(modelDefinitionFile))
			params.WithModelDefinition(openModelDefinitionFile(cmd.ui, modelDefinitionFile))
		} else {
			modelDefinitionFile := args[1]
			cmd.ui.Say("Deploying model with manifest '%s' and model files in '%s'...", terminal.EntityNameColor(manifestFileName), terminal.EntityNameColor(modelDefinitionFile))
			f, err := zipit(modelDefinitionFile + "/")
			if err != nil {
				cmd.ui.Failed("Unexpected error when compressing model directory: %v", err.Error())
				os.Exit(1)
			}
			// reopen the file (I tried to leave the file open but it did not work)
			zip, err := os.Open(f.Name())
			if err != nil {
				fmt.Println("Error open temporary ZIP file: ", err.Error())
			}
			defer os.Remove(zip.Name())
			params.WithModelDefinition(zip)
		}
	}


	// Params process
	tmpFilePath := ""
	if (len(args) - modelDefinitionExisted)  > 1  {
		argsIndex := 1 + modelDefinitionExisted
		//if not has prefix -
		// Params check
		if !strings.HasPrefix(args[argsIndex], "-P") {
			cmd.ui.Failed("Args is not start with -P error: %v", args[argsIndex])
			os.Exit(1)
		} else {
			// Param Num Check, <-P xxxx> is Right
			if (len(args)-modelDefinitionExisted) < 3  {
				cmd.ui.Failed("Params number is not correct: %v", args)
				os.Exit(1)
			}

			//Convert -P params to map
			commandArgs, err := commandConvertParams(args[argsIndex+1])
			if err != nil {
				cmd.ui.Failed("Params convert error", err.Error())
				os.Exit(1)
			}

			//Todo: Control File Size
			//Read and Replace file
			manifestBytes, err := ioutil.ReadFile(manifestFile)
			if err != nil {
				cmd.ui.Failed("Params convert read manifest file failed, ", err.Error())
				os.Exit(1)
			}
			manifestStr := string(manifestBytes)
			for key, value := range commandArgs {
				cmd.ui.Say("Transform Key ", key, value)
				manifestStr = strings.ReplaceAll(manifestStr, "[#"+ key+"]", value)
			}

			//Save tmp file
			tmpYamlFile, tmpFilePath := cmd.buildTempFile()
			_, err = tmpYamlFile.Write([]byte(manifestStr))
			if err != nil {
				cmd.ui.Failed("Write tmpYaml file failed, ", err)
				os.Exit(1)
			}
			manifestFile = tmpFilePath
			defer tmpYamlFile.Close()
		}
	}


	manifestF := openManifestFile(cmd.ui, manifestFile)
	params.WithManifest(&manifestF)
	c, err := NewDlaaSClient()
	if err != nil {
		cmd.ui.Failed("Open Mainfest file error: ", err.Error())
	}

	response, err := c.Models.PostModel(params, BasicAuth())
	if err != nil {
		dlassUrl := os.Getenv("DLAAS_URL")
		if dlassUrl == "" {
			cmd.ui.Failed("Env MLSS URL is not set ")
			os.Exit(1)
		}
		s := ""
		switch err.(type) {
		case *models.PostModelUnauthorized:
			s = badUsernameOrPWD
		case *models.PostModelBadRequest:
			if resp, ok := err.(*models.PostModelBadRequest); ok {
				//FIXME MLSS Change: change resp.Payload.Description to resp.Payload.Msg
				if resp.Payload != nil { // we may not have a payload
					s = fmt.Sprintf("Error: %s. %s", resp.Payload.Error, resp.Payload.Msg)
					cmd.ui.Failed(s)
					cmd.fileRemove(tmpFilePath)
					os.Exit(1)
				} else {
					s = fmt.Sprintf("Bad request: %s", err.Error())
					cmd.ui.Failed(s)
					cmd.fileRemove(tmpFilePath)
					os.Exit(1)
				}
			}
		}
		if s == "" {
			apiError := err.(*runtime.APIError)
			clientResp := apiError.Response.(runtime.ClientResponse)
			s = fmt.Sprintf("%v (code: %v) Message: %v", apiError.OperationName, clientResp.Code(), clientResp.Message())
			cmd.ui.Failed(s)
			os.Exit(1)
		}
		responseError(s, err, cmd.ui)
		if err != nil {
			cmd.ui.Failed("Remove tmpYaml file failed, ", err.Error())
			cmd.fileRemove(tmpFilePath)
			os.Exit(1)
		}
		if err != nil {
			cmd.ui.Failed("Remove tmpYaml file failed, ", err.Error())
			cmd.fileRemove(tmpFilePath)
			os.Exit(1)
		}
		return nil
	}

	cmd.fileRemove(tmpFilePath)

	id := LocationToID(response.Location)
	cmd.ui.Say("Model ID: %s", terminal.EntityNameColor(id))

	//FIXME MLSS Change: MLSS 1.7 monitoring job
	if iMonitoring {
		cmd.ui.Say("starting monitoring")
		//获取 job status
		//获取 job log 从 es
		//当job status 结束时,继续获取 job log ,当 log返回为空是,结束monitoring
		err := monitorAndGetLogs(id, c, cmd)
		if err != nil {
			cmd.ui.Failed("monitoring job error:", err.Error())
			os.Exit(1)
		}
	}
	wg.Wait()
	if jobStatus == "COMPLETED" {
		cmd.ui.Ok()
		os.Exit(0)
	} else {
		cmd.ui.Failed("FAILED")
		os.Exit(1)
	}
	return nil
}


func (cmd *TrainCmd) fileRemove(filePath string) {
	if filePath != "" {
		err := os.Remove(fmt.Sprintf("%v.yml", filePath))
		if err != nil {
			cmd.ui.Failed("Remove tmpYaml file failed, ", err)
			os.Exit(1)
		}
	}
}

var wg sync.WaitGroup

func monitorAndGetLogs(modelID string, c *dlaasClient.Dlaas, cmd *TrainCmd) error {
	wg.Add(1)
	//获取 job status
	go func() {
		var currentRetryTime = 0
		defer wg.Done()
		//const max_get_status_failed_times = 99
		//var currentGetStatusFailedTimes = 0
		for {
			if currentRetryTime > maxRetryTimeForGetState {
				cmd.ui.Failed("Exceed max retry times for get state.")
				// return errors.New("Get Status Error, Exceed max retry tiems.")
				os.Exit(1)
			}
			status, err := getMonitorStatus(modelID, c, cmd)
			if err != nil {
				//cmd.ui.Failed(err.Error())
				cmd.ui.Say(err.Error())
				//currentGetStatusFailedTimes++
				time.Sleep(60 * time.Second)
				currentRetryTime++
				continue
			}

			cmd.ui.Say("current %s jobStatus: %s", modelID, status)
			jobStatus = status

			//if jobStatus == "COMPLETED" || jobStatus == "FAILED" || currentGetStatusFailedTimes >= max_get_status_failed_times {
			if jobStatus == "COMPLETED" || jobStatus == "FAILED" {
				log.Debugf("Get jobStatus break, currentGetStatusFailedTimes: %s", status)
				cmd.ui.Say("Get jobStatus break, currentGetStatusFailedTimes: %s", status)
				break
			}
			time.Sleep(60 * time.Second) // todo why retry ?
		}
	}()

	//获取 job log 从 es
	//当job status 结束时,继续获取 job log ,当 log返回为空是,结束monitoring
	cmd.ui.Say("Get Log Start: ")
	_, err := getLogs(modelID, c, cmd)
	if err != nil {
		cmd.ui.Failed("Job get log error:", err.Error())
		return err
	}
	_, err = getMonitorStatus(modelID, c, cmd)
	if err != nil {
		cmd.ui.Failed("Job Get Status error:", err.Error())
		return err
	}
	return nil
}

func getMonitorStatus(modelID string, c *dlaasClient.Dlaas, cmd *TrainCmd) (string, error) {
	//获取 job status
	//获取 job log 从 es
	//当job status 结束时,继续获取 job log ,当 log返回为空是,结束monitoring

	params := models.NewGetModelParams().
		WithModelID(modelID).
		WithTimeout(defaultOpTimeout)

	modelInfo, err := c.Models.GetModel(params, BasicAuth())

	if err != nil {
		//var s string
		//switch err.(type) {
		//case *models.GetModelUnauthorized:
		//	s = badUsernameOrPWD
		//case *models.GetModelNotFound:
		//	s = "Model not found."
		//}
		//responseError(s, err, cmd.ui)
		cmd.ui.Failed("getMonitorStatus failed")
		return "", err
	}

	m := modelInfo.Payload
	status := m.Training.TrainingStatus.Status

	return status, nil
}

func getLogs(modelID string, c *dlaasClient.Dlaas, cmd *TrainCmd) (string, error) {
	wg.Add(1)
	//获取 job status
	//获取 job log 从 es
	//当job status 结束时,继续获取 job log ,当 log返回为空是,结束monitoring

	// cmd.ui.Say("debug training-logs: doing web socket")
	dlaasURL := os.Getenv("DLAAS_URL")
	if dlaasURL == "" {
		return "", errors.New("Dlaas url is null")
	}
	if !strings.HasSuffix(dlaasURL, dlaasClient.DefaultBasePath) {
		dlaasURL = dlaasURL + dlaasClient.DefaultBasePath
	}
	u, _ := url.Parse(dlaasURL)
	host := u.Host
	cmd.ui.Say("debug training-logs: host: %s", host)

	// what's the best way to do this?
	//var streamHost string
	//var wsprotocol string
	//if strings.HasPrefix(host, "gateway.") || strings.HasPrefix(host, "gateway-") {
	//	streamHost = strings.Replace(host, "gateway", "stream", 1) + ":443"
	//} else {
	//	streamHost = host
	//}
	//if strings.HasPrefix(streamHost, "stream.") || strings.HasPrefix(streamHost, "stream-") {
	//	wsprotocol = "wss"
	//} else {
	//wsprotocol = "ws"
	//}

	//var logsOrMetrics string
	//logsOrMetrics = "logs"

	//var fullpath = fmt.Sprintf("%s/v1/models/%s/%s", u.Path, modelID, logsOrMetrics)
	//rawQuery := fmt.Sprintf("version=2017-03-17&mlss-token-for-logs=dc4e9956-594c-48f4-af3d-36f2fff1a428&mlss-userid=%v", os.Getenv("MLSS_AUTH_USER"))
	//rawQuery := fmt.Sprintf("version=2017-03-17&mlss-userid=%v", os.Getenv("MLSS_AUTH_USER"))
	//rawQuery := fmt.Sprintf("version=2017-03-17")
	//trainingLogURL := url.URL{Scheme: wsprotocol, Host: streamHost, Path: fullpath, RawQuery: rawQuery}

	//FIXME wtss
	//headers := http.Header{
	//	"Accept-Encoding": []string{"gzip, deflate"},
	//	"Cache-Control":   []string{"no-cache"},
	//	"MLSS-UserID":     []string{username},
	//	"MLSS-Passwd":     []string{password},
	//	"MLSS-Auth-Type":  []string{authType},
	//	"Authorization":   []string{"Basic dGVkOndlbGNvbWUx"},
	//	//"MLSS-Token":        []string{"dc4e9956-594c-48f4-af3d-36f2fff1a428"},
	//	//"X-Watson-Userinfo": []string{"bluemix-instance-id=test-user"},
	//
	//	//"Upgrade": []string{"websocket"}
	//	//"Connection": []string{"Upgrade"},
	//	//"Sec-WebSocket-Key": []string{challengeKey},
	//	//"Sec-WebSocket-Version": []string{"13"},
	//	//"Host": []string{host},
	//	//"Origin": []string{fmt.Sprintf("http://%v", os.Getenv("PUBLIC_IP"))},
	//}

	//cmd.ui.Say("call: websocket.Dial: %s\n", trainingLogURL.String())

	//wsConfig, err := websocket.NewConfig(trainingLogURL.String(), origin)
	//wsConfig, err := websocket.NewConfig(trainingLogURL.String(), fmt.Sprintf("http://%v", os.Getenv("PUBLIC_IP")))
	//wsConfig.Header = headers
	//wsConfig.Protocol = []string{wsprotocol}
	var path = fmt.Sprintf("/di/v1/models/%v/logs", modelID)
	var query = fmt.Sprintf("follow=true&version=2017-02-13")

	urlForNewClient := url.URL{
		Scheme:     "ws",
		Host:       host,
		Path:       path,
		RawQuery:   query,
		ForceQuery: true,
	}
	cmd.ui.Say("connecting to %s", u.String())

	headersForNewClient := http.Header{
		//"Accept-Encoding": []string{"gzip, deflate"},
		//"Cache-Control":   []string{"no-cache"},
		//"MLSS-Token": []string{token},
		"MLSS-UserID":    []string{username},
		"MLSS-Passwd":    []string{password},
		"MLSS-Auth-Type": []string{authType},
		//"Authorization":     []string{"Basic dGVkOndlbGNvbWUx"},
		//"X-Watson-Userinfo": []string{"bluemix-instance-id=test-user"},

		//"Upgrade": []string{"websocket"},
		//"Connection": []string{"Upgrade"},
		//"Sec-WebSocket-Key": []string{challengeKey},
		//"Sec-WebSocket-Version": []string{"13"},
		//"Host": []string{host},
		//"Origin": []string{fmt.Sprintf("http://%v", host)},
	}
	cmd.ui.Say("call: websocket.Dial: %s\n", urlForNewClient.String())

	//const max_reconnection_times = 3
	//var currentTimes = 0
	go func() {
		var currentRetryTimes = 0
		defer wg.Done()
		for {
			if currentRetryTimes > maxRetryTimeForGetLog {
				cmd.ui.Failed("exceed max retry times for get log")
				break
			}
			wsClient, resp, err := websocket.DefaultDialer.Dial(urlForNewClient.String(), headersForNewClient)
			if err != nil {
				cmd.ui.Say("Cannot configure websocket: %s", err)
				return
			}
			if resp != nil {
				bytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					cmd.ui.Failed("dial:", err)
				}
				cmd.ui.Say("body: %s\n", string(bytes))
			}
			defer wsClient.Close()

			//err := doWebsocket(wsConfig, cmd)
			err = doWebsocketForNewClient(wsClient, cmd)

			if err != nil {
				cmd.ui.Failed("doWebsocket, error: %s", err.Error())
				currentRetryTimes++
				time.Sleep(10 * time.Second)
			} else {
				cmd.ui.Say("getting logs finished")
				break
			}
		}
	}()
	//fmt.Println("sleep!")
	//time.Sleep(999 * time.Second)
	return "", nil
}

//func doWebsocketForNewClient(conn *websocket.Conn, cmd *TrainCmd) interface{} {
//
//}

//func doWebsocket(config *websocket.Config, cmd *TrainCmd) error {
//	cmd.ui.Say("doWebsocket now")
//	ws, err := websocket.DialConfig(config)
//	if err != nil {
//		cmd.ui.Say("Cannot connect to remote websocket API: %s", err)
//		return fmt.Errorf("Cannot connect to remote websocket API: %s", err)
//	}
//	defer ws.Close()
//
//	var msg = make([]byte, 1024*8)
//	var line = 1
//	for {
//		nread, err := ws.Read(msg)
//
//		//test
//		//if line == 50 && alright == true{
//		//	alright = false
//		//	return fmt.Errorf("test err in line: %v", line)
//		//}
//
//		if err != nil {
//			if err == io.EOF || strings.Contains(err.Error(), "close 1000 (normal)") {
//				break
//			}
//
//			cmd.ui.Say("Reading training logs failed: %s", err.Error())
//			return fmt.Errorf("Reading training logs failed: %s", err.Error())
//		}
//
//		if line > currentLine {
//			fmt.Printf("line: %v: %s\n", line, msg[:nread])
//			currentLine = line
//		}
//
//		line++
//	}
//	return nil
//}

func doWebsocketForNewClient(ws *websocket.Conn, cmd *TrainCmd) error {
	cmd.ui.Say("doWebsocket now")
	//ws, err := websocket.DialConfig(config)
	//if err != nil {
	//	cmd.ui.Say("Cannot connect to remote websocket API: %s", err)
	//	return fmt.Errorf("Cannot connect to remote websocket API: %s", err)
	//}
	//defer ws.Close()

	//var msg = make([]byte, 1024*8)
	var line = 1
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "close 1000 (normal)") || strings.Contains(err.Error(), "close 1006 (abnormal closure)") {
				break
			}

			cmd.ui.Say("Reading training logs failed: %s", err.Error())
			return fmt.Errorf("Reading training logs failed: %s", err.Error())
		}

		if line > currentLine {
			cmd.ui.Say(string(msg))
			currentLine = line
		}

		line++
	}
	return nil
}

func commandConvertParams(args string) (map[string]string, error) {
	if args == "" {
		return nil, errors.New("Command convert failed, key and value are required")
	}

	convertParamsMap := map[string]string{}
	params := strings.Split(args, ",")
	if len(params) <= 0 {
		return nil, errors.New("Params convert failed, key and value format error")
	}

	for _, v := range params {
		if !strings.Contains(v, "=") {
			return nil, errors.New("Params convert failed, key and value format error")
		}
		optionIndex := strings.Index(v, "=")
		key := v[:optionIndex]
		value := v[optionIndex+1:]
		convertParamsMap[key] = value
	}
	return convertParamsMap, nil
}
