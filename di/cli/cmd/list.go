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
	"github.com/IBM-Bluemix/bluemix-cli-sdk/bluemix/terminal"
	"github.com/IBM-Bluemix/bluemix-cli-sdk/plugin"
	"github.com/urfave/cli"
	"webank/DI/restapi/api_v1/client/models"
	"webank/DI/restapi/api_v1/restmodels"
)

// ListCmd is the struct to get all models.
type ListCmd struct {
	ui      terminal.UI
	config  plugin.PluginConfig
	context plugin.PluginContext
}

// NewListCmd is used to get all models.
func NewListCmd(ui terminal.UI, context plugin.PluginContext) *ListCmd {
	return &ListCmd{
		ui:      ui,
		context: context,
	}
}

// Run is the handler for the model-list CLI command.
func (cmd *ListCmd) Run(cliContext *cli.Context) error {
	cmd.config = cmd.context.PluginConfig()
	cmd.ui.Say("Getting all models ...")

	// FIXME MLSS Change: get models filter by username and namespace
	userid := cliContext.String("userid")
	namespace := cliContext.String("namespace")
	page := cliContext.String("page")
	size := cliContext.String("size")
	lflog().Debugf("userid: %s, namespace: %s, page: %s, size: %s", userid, namespace, page, size)
	c, err := NewDlaaSClient()
	if err != nil {
		lflog().Debugf("NewDlaaSClient failed: %s", err.Error())
		cmd.ui.Failed(err.Error())
	}
	lflog().Debugf("Calling ListModels")
	params := models.NewListModelsParams().WithTimeout(defaultOpTimeout)
	params.Userid = &userid
	params.Namespace = &namespace
	if page != "" {
		params.Page = &page
	}
	if size != "" {
		params.Size = &size
	}
	modelz, err := c.Models.ListModels(params, BasicAuth())
	if err != nil {
		lflog().WithError(err).Debugf("ListModels failed")
		var s string
		switch err.(type) {
		case *models.ListModelsUnauthorized:
			s = badUsernameOrPWD
		}
		responseError(s, err, cmd.ui)
		return nil
	}
	lflog().Debugf("Constructing table")
	table := cmd.ui.Table([]string{"ID", "Name", "Framework", "Training status", "Submitted", "Completed"})
	// FIXME MLSS Change: get models filter by username and namespace
	Models := modelz.Payload.Models
	reverse(Models)
	for _, v := range Models {
		ts := v.Training.TrainingStatus
		table.Add(v.ModelID, v.Name, v.Framework.Name+":"+v.Framework.Version, ts.Status, formatTimestamp(ts.Submitted), formatTimestamp(ts.Completed))
	}
	table.Print()
	cmd.ui.Say("\n%d records found,  current page: %v, page size: %v, total page: %d, total records: %d",
		len(modelz.Payload.Models), *params.Page, *params.Size, modelz.Payload.Pages, modelz.Payload.Total)
	return nil
}

// FIXME MLSS Change: get models filter by username and namespace
func reverse(models []*restmodels.Model) {
	if models == nil || len(models) < 2 {
		return
	}
	for i, j := 0, len(models)-1; i < j; i, j = i+1, j-1 {
		models[i], models[j] = models[j], models[i]
	}
}
