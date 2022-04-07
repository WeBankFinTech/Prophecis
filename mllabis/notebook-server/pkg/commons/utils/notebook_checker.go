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
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"strings"
	config "webank/AIDE/notebook-server/pkg/commons/config"
	ccClient "webank/AIDE/notebook-server/pkg/commons/controlcenter/client"
	"webank/AIDE/notebook-server/pkg/commons/logger"
	"webank/AIDE/notebook-server/pkg/models"
	"webank/AIDE/notebook-server/pkg/restapi/operations"
)

const (
	LimitsCpu            = "limits.cpu"
	LimitsMemory         = "limits.memory"
	RequestsNvidiaComGpu = "requests.nvidia.com/gpu"
	//BDAP                 = "bdap"
	//BDAPSAFE             = "bdapsafe"
	//BDP                  = "bdp"
	Standard             = "Standard"
	Custom               = "Custom"
	New                  = "New"
	Existing             = "Existing"
	Non                  = "Non"
)

func CheckPatchNotebook(notebook *models.PatchNotebookRequest, namespace string, name string) error {
	if notebook == nil {
		return errors.New("notebook can not be null")
	}

	maxSparkSession := viper.GetInt64(config.MaxSparkSessionNum)
	if notebook.SparkSessionNum < 0 {
		notebook.SparkSessionNum = 0
	}
	if notebook.SparkSessionNum > maxSparkSession {
		err := fmt.Errorf("sparkSessionNum '%d' in notebook's param is not more than "+
			"max spark session number '%d'", notebook.SparkSessionNum, maxSparkSession)
		return err
	}

	ncClient := GetNBCClient()

	//get notebook crd
	var (
		selfCPU    string
		selfGPU    string
		selfMemory string
	)
	crd, res, err := ncClient.GetNoteBookCRD(namespace, name)
	if err != nil {
		if res.StatusCode != 404 {
			return err
		}
	}
	status, res, err := ncClient.GetNoteBookStatus(namespace, name)
	if err != nil {
		if res.StatusCode != 404 {
			return err
		}
	}
	logger.Logger().Infof("CheckPatchNotebook-GetNoteBookCRD res: %v, status:%v", crd, status)
	if status == "Ready" {
		selfCPU = crd.Limits["cpu"]
		selfGPU = crd.Limits["nvidia.com/gpu"]
		selfMemory = crd.Limits["memory"]
	}

	nodeList, nodeErr := ncClient.listNamespaceNodes(namespace)
	if nil != nodeErr {
		logger.Logger().Errorf("CheckNotebook namespace: %s, err: %v\n", namespace, nodeErr)
		return nodeErr
	}

	if len(nodeList) < 1 {
		return errors.New("namespace requires binding node first, please contact the administrator")
	}

	RQList, err := ncClient.listNamespacedResourceQuota(namespace)
	if len(RQList) == 0 {
		return errors.New("resourceQuota is not existed, please contact the administrator")
	}
	//V1ResourceQuota v1ResourceQuota = items.get(0);
	rq := RQList[0]
	//Map<String, String> hard = v1ResourceQuota.getStatus().getHard();
	//Map<String, String> used = v1ResourceQuota.getStatus().getUsed();
	hard := rq.Status.Hard
	used := rq.Status.Used

	//4.addNotebook时,CPU不能为null或0,不能超配额
	cpu := notebook.CPU
	if cpu <= 0 {
		return errors.New("cpu can not be null or 0")
	}

	//String hardLimitsCPU = hard.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_CPU);
	//String usedLimitsCPU = used.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_CPU);
	hardLimitsCPU := hard[LimitsCpu]
	usedLimitsCPU := used[LimitsCpu]
	if "" == hardLimitsCPU || "" == usedLimitsCPU {
		return errors.New("need to set the RQ to a reasonable size first, please contact the administrator")
	}

	//if (CheckUtils.CheckAvailableCPU(cpu, hardLimitsCPU, usedLimitsCPU) < 0)
	//throw new BDPException("cpu can not greater than the limit");
	availableCPU, err := CheckAvailableCPU(cpu, hardLimitsCPU, usedLimitsCPU, selfCPU) //todo
	if err != nil {
		return err
	}
	if *availableCPU < 0 {
		return errors.New("cpu can not greater than the limit")
	}

	extraResources := notebook.ExtraResources
	if extraResources != "" {
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
			availableGPU, err := CheckAvailableGPU(gpu, hardLimitsGPU, usedLimitsGPU, selfGPU) //todo
			if err != nil {
				return err
			}
			if *availableGPU < 0 {
				return errors.New("gpu can not greater than the limit")
			}
		}
	}

	////        6.addNotebook时,memory不能为null或0,不能超配额

	memory := models.Memory{
		MemoryAmount: &notebook.MemoryAmount,
		MemoryUnit:   *(notebook.MemoryUnit),
	}
	if notebook.MemoryAmount == 0 {
		return errors.New("memory can not be 0")
	}
	hardLimitsMemory := hard[LimitsMemory]
	usedLimitsMemory := used[LimitsMemory]
	mem, err := CheckAvailableMemory(memory, hardLimitsMemory, usedLimitsMemory, selfMemory) //todo
	if *mem < 0 {
		return errors.New("memory can not greater than the limit")
	}

	//split := strings.Split(namespace, "-")
	//if BDAP != split[2] && BDP != split[2] {
	//	return errors.New("cluster name must be BDAP or BDP")
	//}
	////        9.addNotebook时,imageType != null,可选Standard或Custom;
	logger.Logger().Infof("notebook.Image: %v, %v", notebook.ImageName, notebook.ImageType)
	if notebook.ImageName == "" || (notebook.ImageType != Custom && notebook.ImageType != Standard) {
		return errors.New("imageType must be Standard or Custom")
	}

	return nil
}

func CheckNotebook(params operations.PostNamespacedNotebookParams) error {
	//logr := logger.LocLogger(logrus.StandardLogger().WithField("utils", "CheckNotebook"))
	notebook := params.Notebook
	if notebook == nil {
		return errors.New("notebook can not be null")
	}
	namespace := notebook.Namespace
	//err := CheckNamespace(namespace)
	//if err != nil {
	//	return err
	//}

	maxSparkSession := viper.GetInt64(config.MaxSparkSessionNum)
	if params.Notebook.SparkSessionNum < 0 {
		params.Notebook.SparkSessionNum = 0
	}
	if params.Notebook.SparkSessionNum > maxSparkSession {
		err := fmt.Errorf("sparkSessionNum '%d' in notebook's param is not more than "+
			"max spark session number '%d'", params.Notebook.SparkSessionNum, maxSparkSession)
		return err
	}

	ncClient := GetNBCClient()

	nodeList, nodeErr := ncClient.listNamespaceNodes(*namespace)
	if nil != nodeErr {
		logger.Logger().Errorf("CheckNotebook namespace: %s, err: %v\n", *namespace, nodeErr)
		return nodeErr
	}

	if len(nodeList) < 1 {
		return errors.New("namespace requires binding node first, please contact the administrator")
	}

	RQList, err := ncClient.listNamespacedResourceQuota(*namespace)
	if len(RQList) == 0 {
		return errors.New("resourceQuota is not existed, please contact the administrator")
	}
	//V1ResourceQuota v1ResourceQuota = items.get(0);
	rq := RQList[0]
	//Map<String, String> hard = v1ResourceQuota.getStatus().getHard();
	//Map<String, String> used = v1ResourceQuota.getStatus().getUsed();
	hard := rq.Status.Hard
	used := rq.Status.Used

	//4.addNotebook时,CPU不能为null或0,不能超配额
	//Double cpu = notebook.getCpu();
	cpu := notebook.CPU
	//if (cpu == null || cpu.equals(0))
	if cpu == nil || *cpu == 0 {
		//throw new BDPException("cpu can not be null or 0");
		return errors.New("cpu can not be null or 0")
	}

	//String hardLimitsCPU = hard.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_CPU);
	//String usedLimitsCPU = used.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_CPU);
	hardLimitsCPU := hard[LimitsCpu]
	usedLimitsCPU := used[LimitsCpu]

	if "" == hardLimitsCPU || "" == usedLimitsCPU {
		return errors.New("need to set the RQ to a reasonable size first, please contact the administrator")
	}

	//if (CheckUtils.CheckAvailableCPU(cpu, hardLimitsCPU, usedLimitsCPU) < 0)
	//throw new BDPException("cpu can not greater than the limit");
	availableCPU, err := CheckAvailableCPU(*cpu, hardLimitsCPU, usedLimitsCPU, "")
	if err != nil {
		return err
	}
	if *availableCPU < 0 {
		return errors.New("cpu can not greater than the limit")
	}

	extraResources := notebook.ExtraResources
	if extraResources != "" {
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
			availableGPU, err := CheckAvailableGPU(gpu, hardLimitsGPU, usedLimitsGPU, "")
			if err != nil {
				return err
			}
			if *availableGPU < 0 {
				return errors.New("gpu can not greater than the limit")
			}
		}
	}

	////        6.addNotebook时,memory不能为null或0,不能超配额
	//Memory memory = notebook.getMemory();
	memory := notebook.Memory
	//if (memory == null || memory.getMemoryAmount() == null || memory.getMemoryAmount().equals(0))
	//throw new BDPException("memory can not be null or 0");
	if memory == nil || memory.MemoryAmount == nil || *memory.MemoryAmount == 0 {
		return errors.New("memory can not be null or 0")
	}
	//String hardLimitsMemory = hard.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_MEMORY);
	//String usedLimitsMemory = used.get(BDPConstants.Notebook.NotebookRequestRQ.LIMITS_MEMORY);
	hardLimitsMemory := hard[LimitsMemory]
	usedLimitsMemory := used[LimitsMemory]
	//if (CheckUtils.CheckAvailableMemory(memory, hardLimitsMemory, usedLimitsMemory) < 0)
	//throw new BDPException("memory can not greater than the limit");
	mem, err := CheckAvailableMemory(*memory, hardLimitsMemory, usedLimitsMemory, "")
	if *mem < 0 {
		return errors.New("memory can not greater than the limit")
	}
	////        7.addNotebook时,name != null, name.len <= 255, reg.match(a-z A-Z 0-9 "-");
	//String name = notebook.getName();
	//if (name == null || name.length() > 255 || !name.matches("^[a-zA-Z0-9-]+$"))
	//throw new BDPException("notebook's name not available");
	name := notebook.Name
	matched, err := regexp.MatchString("^[a-zA-Z0-9-]+$", *name)
	if err != nil {
		return errors.New(err.Error())
	}
	if name == nil || len(*name) > 255 || !matched {
		return errors.New("notebook's name not available")
	}
	////        8.addNotebook时,namespace != null,ns in BDAP
	//String[] split = namespace.split("-");
	//if (!BDPConstants.ClusterName.BDAP.equals(split[2]))
	//throw new BDPException("cluster name must be BDAP");
	//split := strings.Split(*namespace, "-")
	//if BDAP != split[2] && BDP != split[2] {
	//	return errors.New("cluster name must be BDAP or BDP")
	//}
	////        9.addNotebook时,imageType != null,可选Standard或Custom;
	//if (notebook.getImage() == null || !EnumUtils.isValidEnum(BDPConstants.imageType.class, notebook.getImage().getImageType()))
	//throw new BDPException("imageType must be Standard or Custom");
	logger.Logger().Infof("notebook.Image: %v, %v", *notebook.Image, *notebook.Image.ImageType)
	if notebook.Image == nil || (*notebook.Image.ImageType != Custom && *notebook.Image.ImageType != Standard) {
		return errors.New("imageType must be Standard or Custom")
	}
	////check workspaceVolume
	//if notebook.WorkspaceVolume == nil {
	//	return errors.New("workspaceVolume can not be null")
	//}
	//FIXME MLSS Change: added check notebook workPath
	client := ccClient.GetCcClient(viper.GetString(config.CCAddress))
	currentUserId := params.HTTPRequest.Header.Get("X-DLaaS-UserID")
	//namespaceForNoteBook := params.Notebook.Namespace
	if notebook.WorkspaceVolume != nil && notebook.ProxyUser != "" {
		if nil == notebook.WorkspaceVolume.LocalPath || *notebook.WorkspaceVolume.LocalPath == "" {
			return errors.New("localPath can not be null")
		}
		if nil == notebook.WorkspaceVolume.MountPath || *notebook.WorkspaceVolume.MountPath == "" {
			return errors.New("mountPath can not be null")
		}
		//if (!EnumUtils.isValidEnum(BDPConstants.Notebook.MountType.class, notebook.getWorkspaceVolume().getMountType()))
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
	//FIXME MLSS Change: added check notebook dataPath
	if notebook.DataVolume != nil && notebook.ProxyUser == "" {
		logger.Logger().Debugf("CheckNotebook params.Notebook.DataVolume: %v", notebook.DataVolume)
		//for _, dataVolume := range notebook.DataVolume {
		dataPath := notebook.DataVolume.LocalPath
		logger.Logger().Debugf("UserStoragePathCheck dataPath: %v", dataPath)
		if dataPath != nil && *dataPath != "" {
			errForDataPath := client.UserStoragePathCheck(params.HTTPRequest.Header.Get(ccClient.CcAuthToken), currentUserId, *dataPath)
			if errForDataPath != nil {
				logger.Logger().Errorf("check dataPath failed: %v", errForDataPath.Error())
				return fmt.Errorf(errForDataPath.Error()+">>>dataPath: %s", *dataPath)
			}
		}
		//}
	}

	return nil
}

func CheckNamespace(Namespace *string) error {
	return nil
}

func CheckAvailableCPU(cpu float64, hardLimitsCPU string, usedLimitsCPU string, selfCPU string) (*float64, error) {
	hard, e := GetCPUValueWithUnits(hardLimitsCPU)
	if e != nil {
		return nil, e
	}
	used, e := GetCPUValueWithUnits(usedLimitsCPU)
	if e != nil {
		return nil, e
	}
	self, e := GetCPUValueWithUnits(selfCPU)
	if e != nil {
		return nil, e
	}
	logger.Logger().Infof("CheckAvailableCPU, hard: %v, used: %v, cpu: %v, self: %v ", *hard, *used, cpu*1000, *self)
	f := *hard - (*used - *self) - cpu*1000
	return &f, nil
}

func CheckAvailableGPU(gpu int, hardLimitsGPU string, usedLimitsGPU string, selfGPU string) (*int, error) {
	hard, e := strconv.Atoi(hardLimitsGPU)
	if e != nil {
		return nil, e
	}
	used, e := strconv.Atoi(usedLimitsGPU)
	if e != nil {
		return nil, e
	}
	self := 0
	if selfGPU != "" {
		self, _ = strconv.Atoi(selfGPU)
	}

	logger.Logger().Infof("CheckAvailableCPU, hard: %v, used: %v, gpu: %v, self: %v ", hard, used, gpu, self)
	re := hard - (used - self) - gpu
	return &re, nil
}

func CheckAvailableMemory(memory models.Memory, hardLimitsMemory string, usedLimitsMemory string, selfMemory string) (*float64, error) {
	logger.Logger().Infof("CheckAvailableMemory hardLimitsMemory: %v, usedLimitsMemory: %v, "+
		"memoryAmount: %v, memoryUnit: %v", hardLimitsMemory, usedLimitsMemory, *(memory.MemoryAmount), memory.MemoryUnit)
	//	mem := *memory.MemoryAmount * 1024 * 1024 * 1024 //todo
	var mem float64 = 0
	switch memory.MemoryUnit {
	case "Gi":
		mem = *memory.MemoryAmount * 1024 * 1024 * 1024
	case "Mi":
		mem = *memory.MemoryAmount * 1024 * 1024
	case "k", "K", "ki":
		mem = *memory.MemoryAmount * 1024
	default:
		return nil, fmt.Errorf("memory unit error, unit: %v", memory.MemoryUnit)
	}

	hard, e := GetMemoryValueWithUnits(hardLimitsMemory)

	if e != nil {
		return nil, e
	}
	logger.Logger().Infof("CheckAvailableMemory00")
	used, e := GetMemoryValueWithUnits(usedLimitsMemory)
	//if e != nil {
	//	return nil, e
	//}
	self, e := GetMemoryValueWithUnits(selfMemory)
	if e != nil {
		return nil, e
	}
	logger.Logger().Infof("CheckAvailableMemory0")
	logger.Logger().Infof("CheckAvailableMemory1 usedLimitsMemory:%v", usedLimitsMemory)
	logger.Logger().Infof("CheckAvailableMemory2 hardLimitsMemory:%v", hardLimitsMemory)
	logger.Logger().Infof("CheckAvailableMemory2 selfMemory:%v", selfMemory)
	logger.Logger().Infof("CheckAvailableMemory2 selfMemory:%v", self)

	re := *hard - (*used - *self) - mem
	return &re, nil
}

func GetMemoryValueWithUnits(withUnits string) (*float64, error) {
	if withUnits == "" {
		r := float64(0)
		return &r, nil
	}
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
	//if (NumberUtils.isParsable(withUnits))
	//return Double.parseDouble(withUnits) * 1000;
	//if (withUnits.endsWith("m"))
	//return Double.parseDouble(withUnits.substring(0, withUnits.length()-1));
	//return nil, errors.New("GetCPUValueWithUnits failed, str: " + withUnits);
	if withUnits == "" {
		r := float64(0)
		return &r, nil
	}
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
