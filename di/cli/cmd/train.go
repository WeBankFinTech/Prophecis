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
	"github.com/go-openapi/runtime"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"webank/DI/restapi/api_v1/client/models"

	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin"
	"github.com/urfave/cli"
	"path/filepath"
	dlaasClient "webank/DI/restapi/api_v1/client"
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

const maxRetryTimeForGetState = 99
const maxRetryTimeForGetLog = 99

//var alright = true

// Run is the handler for the model-deploy CLI command.
func (cmd *TrainCmd) Run(cliContext *cli.Context) error {
	iMonitoring := true
	cmd.config = cmd.context.PluginConfig()
	args := cliContext.Args()
	params := models.NewPostModelParams().WithTimeout(defaultOpTimeout)

	_, manifestFile := filepath.Split(args[0])
	tmpYamlUuid := uuid.New()
	tmpYamlFile, err := os.Create(fmt.Sprintf("%v.yml",tmpYamlUuid))
	if err != nil {
		cmd.ui.Failed("Create tmpYaml file failed, ",err)
		os.Exit(1)
	}
	params.WithManifest(openManifestFile(cmd.ui, args[0]))
	if len(args) >= 3 {
		commandArgs, err := commandConvertParams(args)
		if err != nil {
			cmd.ui.Failed(err.Error())
			os.Exit(1)
		}
		manifestBytes, err := ioutil.ReadFile(manifestFile)
		if err != nil {
			cmd.ui.Failed("Read manifest file failed, ",err)
			os.Exit(1)
		}
		manifestStr := string(manifestBytes)
		for key, value := range commandArgs {
			manifestStr = strings.ReplaceAll(manifestStr, fmt.Sprintf("[#%v]",key),value)
		}
		_, err = tmpYamlFile.Write([]byte(manifestStr))
		if err != nil {
			cmd.ui.Failed("Write tmpYaml file failed, ",err)
			os.Exit(1)
		}
		params.WithManifest(openManifestFile(cmd.ui, fmt.Sprintf("%v.yml",tmpYamlUuid)))
	}
	if len(args) >= 2 {
		if strings.Contains(args[1], ".zip") {
			modelDefinitionFile := args[1]
			cmd.ui.Say("Deploying model with manifest '%s' and model file '%s'...", terminal.EntityNameColor(manifestFile), terminal.EntityNameColor(modelDefinitionFile))
			params.WithModelDefinition(openModelDefinitionFile(cmd.ui, modelDefinitionFile))
		} else {
			modelDir := args[1]
			cmd.ui.Say("Deploying model with manifest '%s' and model files in '%s'...", terminal.EntityNameColor(manifestFile), terminal.EntityNameColor(modelDir))
			f, err2 := zipit(modelDir + "/")
			if err2 != nil {
				cmd.ui.Failed("Unexpected error when compressing model directory: %v", err2)
				os.Exit(1)
			}
			// reopen the file (I tried to leave the file open but it did not work)
			zip, err := os.Open(f.Name())
			if err != nil {
				fmt.Println("Error open temporary ZIP file: ", err)
			}
			defer os.Remove(zip.Name())

			params.WithModelDefinition(zip)
		}
	}

	c, err := NewDlaaSClient()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}
	response, err := c.Models.PostModel(params, BasicAuth())
	if err != nil {
		dlassUrl := os.Getenv("DLAAS_URL")
		if dlassUrl == "" {
			cmd.ui.Failed("DLass url is not set ")
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
				} else {
					s = fmt.Sprintf("Bad request: %s", err.Error())
				}
			}
		}
		if s == "" {
			apiError := err.(*runtime.APIError)
			clientResp := apiError.Response.(runtime.ClientResponse)
			s = fmt.Sprintf("%v (code: %v) Message: %v",apiError.OperationName, clientResp.Code(), clientResp.Message())
		}
		responseError(s, err, cmd.ui)
		err = tmpYamlFile.Close()
		if err != nil {
			cmd.ui.Failed("Remove tmpYaml file failed, ",err.Error())
			os.Exit(1)
		}
		err = os.Remove(fmt.Sprintf("%v.yml",tmpYamlUuid))
		if err != nil {
			cmd.ui.Failed("Remove tmpYaml file failed, ",err.Error())
			os.Exit(1)
		}
		return nil
	}

	err = tmpYamlFile.Close()
	if err != nil {
		cmd.ui.Failed("Remove tmpYaml file failed, ",err)
		os.Exit(1)
	}
	err = os.Remove(fmt.Sprintf("%v.yml",tmpYamlUuid))
	if err != nil {
		cmd.ui.Failed("Remove tmpYaml file failed, ",err)
		os.Exit(1)
	}

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
			//cmd.ui.Failed(err.Error())
			cmd.ui.Failed(err.Error())
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
				cmd.ui.Failed("exceed max retry times for get state")
				break
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
				log.Debugf("get jobStatus break, currentGetStatusFailedTimes: %s", status)
				break
			}
			time.Sleep(60 * time.Second)
		}
	}()

	//获取 job log 从 es
	//当job status 结束时,继续获取 job log ,当 log返回为空是,结束monitoring
	_, err := getLogs(modelID, c, cmd)
	if err != nil {
		//cmd.ui.Failed(err.Error())
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
		cmd.ui.Say("getMonitorStatus failed")
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

	cmd.ui.Say("debug training-logs: doing web socket")
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
		"MLSS-UserID":    []string{username},
		"MLSS-Passwd":    []string{password},
		"MLSS-Auth-Type": []string{authType},
	}
	cmd.ui.Say("call: websocket.Dial: %s\n", urlForNewClient.String())

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
	return "", nil
}

func doWebsocketForNewClient(ws *websocket.Conn, cmd *TrainCmd) error {
	cmd.ui.Say("doWebsocket now")
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
			fmt.Printf("line: %v: %s\n", line, msg)
			currentLine = line
		}

		line++
	}
	return nil
}

func commandConvertParams(commands cli.Args) (map[string]string, error) {
	convertParamsMap := map[string]string{}
	if len(commands) > 2{
		cmdArgs := commands[2:]
		if strings.ToUpper(cmdArgs[0]) != "-P" {
			return nil, errors.New("Command convert failed, argument must start with -P")
		}
		if len(cmdArgs[1]) <= 0 {
			return nil, errors.New("Command convert failed, key and value are required")
		}
		paramsStr := cmdArgs[1]
		params := strings.Split(paramsStr, ",")
		if len(params) <= 0 {
			return nil, errors.New("Params convert failed, key and value format error")
		}
		for _, v := range params {
			if !strings.Contains(v, "=") {
				return nil, errors.New("Params convert failed, key and value format error")
			}
			optionIndex := strings.Index(v, "=")
			key := v[:optionIndex]
			value := v[optionIndex + 1:]
			convertParamsMap[key] = value
		}
	}
	return convertParamsMap, nil
}