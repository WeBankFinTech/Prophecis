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
package utils

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"strings"
	"webank/AIDE/notebook-server/pkg/commons/config"
	ccClient "webank/AIDE/notebook-server/pkg/commons/controlcenter/client"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"
	"webank/AIDE/notebook-server/pkg/restapi/operations"
)

const (
	LimitsCpu            = "limits.cpu"
	LimitsMemory         = "limits.memory"
	RequestsNvidiaComGpu = "requests.nvidia.com/gpu"
	Standard             = "Standard"
	Custom               = "Custom"
	New                  = "New"
	Existing             = "Existing"
	Non                  = "Non"
)

func CheckNotebook(params operations.PostNamespacedNotebookParams) error {
	notebook := params.Notebook
	if notebook == nil {
		return errors.New("notebook can not be null")
	}
	namespace := notebook.Namespace
	ncClient := GetNBCClient()

	nodeList, nodeErr := ncClient.listNamespaceNodes(*namespace)
	if nil != nodeErr {
		return nodeErr
	}

	if len(nodeList) < 1 {
		return errors.New("namespace requires binding node first, please contact the administrator")
	}

	RQList, err := ncClient.listNamespacedResourceQuota(*namespace)
	if len(RQList) == 0 {
		return errors.New("resourceQuota is not existed, please contact the administrator")
	}

	rq := RQList[0]
	hard := rq.Status.Hard
	used := rq.Status.Used

	//4.addNotebook时,CPU不能为null或0,不能超配额
	//Double cpu = notebook.getCpu();
	cpu := notebook.CPU
	if cpu == nil || *cpu == 0 {
		return errors.New("cpu can not be null or 0")
	}

	hardLimitsCPU := hard[LimitsCpu]
	usedLimitsCPU := used[LimitsCpu]

	if "" == hardLimitsCPU || "" == usedLimitsCPU {
		return errors.New("need to set the RQ to a reasonable size first, please contact the administrator")
	}

	//if (CheckUtils.CheckAvailableCPU(cpu, hardLimitsCPU, usedLimitsCPU) < 0)
	//throw new BDPException("cpu can not greater than the limit");
	availableCPU, err := CheckAvailableCPU(*cpu, hardLimitsCPU, usedLimitsCPU)
	if err != nil {
		return err
	}
	if *availableCPU < 0 {
		return errors.New("cpu can not greater than the limit")
	}

	extraResources := notebook.ExtraResources
	extraResourcesMap := make(map[string]string)
	err = json.Unmarshal([]byte(extraResources), &extraResourcesMap)
	if err != nil {
		return errors.New("json Unmarshal extraResources failed, " + err.Error())
	}
	logger.Logger().Infof("extraResourcesMap: %v", extraResourcesMap)
	gpuString := extraResourcesMap["nvidia.com/gpu"]
	matched, err := regexp.MatchString("^\\d+$", gpuString)
	if err != nil {
		return err
	}
	if gpuString != "" && !matched {
		return errors.New("gpu must be integer")
	}
	if gpuString != "" && "0" != gpuString {
		gpu, err := strconv.Atoi(gpuString)
		hardLimitsGPU := hard[RequestsNvidiaComGpu]
		usedLimitsGPU := used[RequestsNvidiaComGpu]
		availableGPU, err := CheckAvailableGPU(gpu, hardLimitsGPU, usedLimitsGPU)
		if err != nil {
			return err
		}
		if *availableGPU < 0 {
			return errors.New("gpu can not greater than the limit")
		}
	}

	//6.addNotebook时,memory不能为null或0,不能超配额
	memory := notebook.Memory
	if memory == nil || memory.MemoryAmount == nil || *memory.MemoryAmount == 0 {
		return errors.New("memory can not be null or 0")
	}
	hardLimitsMemory := hard[LimitsMemory]
	usedLimitsMemory := used[LimitsMemory]
	mem, err := CheckAvailableMemory(*memory, hardLimitsMemory, usedLimitsMemory)
	if *mem < 0 {
		return errors.New("memory can not greater than the limit")
	}

	//7.Notebook名称必须为字母、数字或横杆
	name := notebook.Name
	matched, err = regexp.MatchString("^[a-zA-Z0-9-]+$", *name)
	if err != nil {
		return errors.New(err.Error())
	}
	if name == nil || len(*name) > 255 || !matched {
		return errors.New("notebook's name not available")
	}

	//8.addNotebook时,imageType != null,可选Standard或Custom;
	logger.Logger().Infof("notebook.Image: %v, %v", *notebook.Image, *notebook.Image.ImageType)
	if notebook.Image == nil || (*notebook.Image.ImageType != Custom && *notebook.Image.ImageType != Standard) {
		return errors.New("imageType must be Standard or Custom")
	}

	// Get User Message From CC Server
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	currentUserId := params.HTTPRequest.Header.Get("X-DLaaS-UserID")
	//9.WorkSpace Volumne Check
	if notebook.WorkspaceVolume != nil {
		if nil == notebook.WorkspaceVolume.LocalPath || *notebook.WorkspaceVolume.LocalPath == "" {
			return errors.New("localPath can not be null")
		}
		if nil == notebook.WorkspaceVolume.MountPath || *notebook.WorkspaceVolume.MountPath == "" {
			return errors.New("mountPath can not be null")
		}
		if nil == notebook.WorkspaceVolume || (*notebook.WorkspaceVolume.MountType != New && *notebook.WorkspaceVolume.MountType != Existing && *notebook.WorkspaceVolume.MountType != Non) {
			return errors.New("mountType only can be New,Existing,None")
		}
		if notebook.WorkspaceVolume.Size == nil {
			return errors.New("size can not be null")
		}
		workPath := notebook.WorkspaceVolume.LocalPath
		errForWorkPath := client.UserStoragePathCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *workPath)
		if errForWorkPath != nil {
			logger.Logger().Errorf("check workPath failed: %v", errForWorkPath.Error())
			return errors.New(errForWorkPath.Error() + ">>>workPath: %s" + *workPath)
		}
	}
	//10. DataVolume Check
	if notebook.DataVolume != nil {
		logger.Logger().Debugf("CheckNotebook params.Notebook.DataVolume: %v", notebook.DataVolume)
		dataPath := notebook.DataVolume.LocalPath
		errForDataPath := client.UserStoragePathCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *dataPath)
		if errForDataPath != nil {
			logger.Logger().Errorf("check dataPath failed: %v", errForDataPath.Error())
			return errors.New(errForDataPath.Error() + ">>>dataPath: %s" + *dataPath)
		}
	}

	return nil
}

func CheckAvailableCPU(cpu float64, hardLimitsCPU string, usedLimitsCPU string) (*float64, error) {
	hard, e := GetCPUValueWithUnits(hardLimitsCPU)
	if e != nil {
		return nil, e
	}
	used, e := GetCPUValueWithUnits(usedLimitsCPU)
	if e != nil {
		return nil, e
	}
	logger.Logger().Infof("CheckAvailableCPU, hard: %v, used: %v, cpu: %v, ", *hard, *used, cpu*1000)
	f := *hard - *used - cpu*1000
	return &f, nil
}

func CheckAvailableGPU(gpu int, hardLimitsGPU string, usedLimitsGPU string) (*int, error) {
	hard, e := strconv.Atoi(hardLimitsGPU)
	if e != nil {
		return nil, e
	}
	used, e := strconv.Atoi(usedLimitsGPU)
	if e != nil {
		return nil, e
	}
	logger.Logger().Infof("CheckAvailableCPU, hard: %v, used: %v, gpu: %v, ", hard, used, gpu)
	re := hard - used - gpu
	return &re, nil
}

func CheckAvailableMemory(memory models.Memory, hardLimitsMemory string, usedLimitsMemory string) (*float64, error) {
	logger.Logger().Infof("CheckAvailableMemory08")
	mem := *memory.MemoryAmount * 1024 * 1024 * 1024
	hard, e := GetMemoryValueWithUnits(hardLimitsMemory)

	if e != nil {
		return nil, e
	}
	used, e := GetMemoryValueWithUnits(usedLimitsMemory)
	logger.Logger().Debugf("CheckAvailableMemory, usedLimitsMemory:%v,"+
		" hardLimitsMemory: %v, men:%v , hard: %v", usedLimitsMemory, hardLimitsMemory, mem, *hard)
	logger.Logger().Infof("CheckAvailableCPU, hard: %v, used: %v, mem: %v, ", *hard, *used, mem)
	re := *hard - *used - mem
	return &re, nil
}

func GetMemoryValueWithUnits(withUnits string) (*float64, error) {
	f, e := strconv.ParseFloat(withUnits, 64)
	if e == nil {
		return &f, nil
	}
	if strings.HasSuffix(withUnits, "Mi") {
		//return strconv.ParseFloat(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024;
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
		//return Double.parseDouble(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024 * 1024;
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
		//return Double.parseDouble(withUnits.substring(0, withUnits.length()-2)) * 1024 * 1024 * 1024;
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
