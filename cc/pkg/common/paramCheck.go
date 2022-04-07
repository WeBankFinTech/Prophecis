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
package common

import (
	"errors"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/restapi/restapi/operations/resources"
	"regexp"
	"strconv"
	"strings"
)

func CheckForAddUser(user models.UserRequest) error {
	if user.Name == "" {
		return errors.New("username can't be empty")
	}

	matched, _ := regexp.MatchString(constants.UsernameRegx, user.Name)

	if !matched {
		return errors.New("the name is not valid, must be one of [a-zA-Z0-9]+|[a-zA-Z0-9]+[_][c]|[v][_][a-zA-Z]+|[v][_][a-zA-Z]+[_][c]")
	}

	if user.Type != constants.TypeUser && user.Type != constants.TypeSYSTEM {
		return errors.New("user's type must be SYSTEM or USER")
	}

	if user.UID < constants.GUIdMinValue || user.UID > constants.GUIdMaxValue {
		return errors.New("uid is not valid, must be (0,90000)")
	}

	if user.Gid < constants.GUIdMinValue || user.Gid > constants.GUIdMaxValue {
		return errors.New("gid is not valid, must be (0,90000)")
	}

	if "" != user.GUIDCheck {
		if user.GUIDCheck != "0" && user.GUIDCheck != "1" {
			return errors.New("GUIDCheck must be 0 or 1")
		}
	} else {
		user.GUIDCheck = "1"
	}

	if user.GUIDCheck == "1" {
		if user.UID == constants.GUIdMinValue && user.Gid == constants.GUIdMinValue {
			return errors.New("gid,uid can't be 0 at the same time, where guidCheck is true")
		}
	}

	return nil
}

func CheckForAddingUser(user models.User) error {
	if user.Name == "" {
		return errors.New("username can't be empty")
	}

	matched, _ := regexp.MatchString(constants.UsernameRegx, user.Name)

	if !matched {
		return errors.New("the name is not valid, must be one of [a-zA-Z0-9]+|[a-zA-Z0-9]+[_][c]|[v][_][a-zA-Z]+|[v][_][a-zA-Z]+[_][c]")
	}

	if user.Type != constants.TypeUser && user.Type != constants.TypeSYSTEM {
		return errors.New("user's type must be SYSTEM or USER")
	}

	if user.UID < constants.GUIdMinValue || user.UID > constants.GUIdMaxValue {
		return errors.New("uid is not valid, must be (0,90000)")
	}

	if user.Gid < constants.GUIdMinValue || user.Gid > constants.GUIdMaxValue {
		return errors.New("gid is not valid, must be (0,90000)")
	}

	if "" != user.GUIDCheck {
		if user.GUIDCheck != "0" && user.GUIDCheck != "1" {
			return errors.New("GUIDCheck must be 0 or 1")
		}
	} else {
		user.GUIDCheck = "1"
	}

	if user.GUIDCheck == "1" {
		if user.UID == constants.GUIdMinValue && user.Gid == constants.GUIdMinValue {
			return errors.New("gid,uid can't be 0 at the same time, where guidCheck is true")
		}
	}

	if user.ID == 0 {
		//return ResponderFunc(http.StatusBadRequest, "failed to update user", "user id is invalid")
		return errors.New("user id is invalid")
	}

	return nil
}

func CheckForUpdateGroup(group models.Group) error {
	if group.GroupType != constants.TypePRIVATE && group.GroupType != constants.TypeSYSTEM {
		return errors.New("group type must be PRIVATE or SYSTEM")
	}

	if group.DepartmentName == "" {
		return errors.New("department name can't be empty")
	}

	return nil
}

func CheckForGroup(group models.Group) error {
	if "" == group.Name {
		return errors.New("group name can't be empty")
	}

	if len(group.Name) > constants.NameLengthLimit {
		return errors.New("the length of group name is not up to 128")
	}

	if group.GroupType != constants.TypePRIVATE && group.GroupType != constants.TypeSYSTEM {
		return errors.New("group type must be PRIVATE or SYSTEM")
	}

	if group.GroupType == constants.TypePRIVATE {
		matched, _ := regexp.MatchString(constants.PrivateGroupNameRegx, group.Name)
		if !matched {
			return errors.New("the private group naming convention is gp(-private)(-[_a-z0-9]+)")
		}
	}

	if group.GroupType == constants.TypeSYSTEM {
		matched, _ := regexp.MatchString(constants.SystemGroupNameRegx, group.Name)
		if !matched {
			return errors.New("the system group naming convention is gp-[department]-[application]")
		}
	}

	if group.DepartmentName == "" {
		return errors.New("department name can't be empty")
	}

	return nil
}

func CheckForStorage(storage models.Storage) error {
	path := storage.Path

	if !strings.HasPrefix(path, "/") {
		return errors.New("path must be start with /")
	}

	if len(strings.Split(path[1:], "/")) < 4 {
		return errors.New("the storage path is at least level 4")
	}

	if storage.Type != strings.Split(path[1:], "/")[1] {
		return errors.New("type must be the same as in path")
	}
	return nil
}

func CheckForNamespace(namespaceRequest models.NamespaceRequest) error {
	namespace := namespaceRequest.Namespace

	return CheckNS(namespace)
}

func CheckNS(namespace string) error {
	if "" == namespace {
		return errors.New("namespace can't be empty")
	}

	if len(namespace) > constants.NameLengthLimit {
		return errors.New("namespace len can't more than 128")
	}

	matchedNamespace, _ := regexp.MatchString(constants.NamespaceNameRegx, namespace)

	if !matchedNamespace {
		//return errors.New("the namespace naming convention is ns-[IDC]-[cluster]-[department]-[application]-{[option]-}[network]")
		return errors.New("the namespace naming convention is ns-[department]-[application]")
	}

	matchedSystemNamespaceName, _ := regexp.MatchString(constants.SystemNamespaceName, namespace)

	if matchedSystemNamespaceName {
		return errors.New("the system namespace can't be created")
	}

	//split := strings.Split(namespace, "-")
	//netWorkName := split[len(split)-1]
	//clusterName := split[2]

	//if constants.BDP == clusterName {
	//	if strings.Contains(namespace, constants.SAFE) {
	//		return errors.New("namespace is unavailable")
	//	}
	//}

	//if netWorkName == constants.SAFE && clusterName != constants.BDAP {
	//	return errors.New("bdap must be in safe netWork")
	//}
	return nil
}

func CheckForRQParams(rq models.ResourcesQuota) error {

	if "" == rq.Namespace || "" == rq.CPU || "" == rq.Gpu || nil == rq.Memory {
		return errors.New("namespace and cpu and gpu and memory can't be empty")
	}

	if "" == rq.Memory["memoryAmount"] || "" == rq.Memory["memoryUnit"] {
		return errors.New("memoryAmount and memoryUnit can't be empty")
	}

	if "Mi" != rq.Memory["memoryUnit"] {
		return errors.New("memoryUnit must be Mi")
	}

	if strings.Contains(rq.CPU, ".") {
		replace := strings.Replace(rq.CPU, ".", "a", 1)
		logger.Logger().Debugf("format cpu: %v", replace)
		if len(strings.Split(replace, "a")[1]) == 0 || len(strings.Split(replace, "a")[1]) > 1 {
			return errors.New("cpu can only have one decimal place at most")
		}
	}

	parCPU, cpuErr := StringToFloat64(rq.CPU)

	if nil != cpuErr {
		return cpuErr
	}

	if parCPU < 0 {
		return errors.New("cpu is invalid data")
	}

	parGPU, gpuErr := StringToFloat64(rq.Gpu)

	if nil != gpuErr {
		return gpuErr
	}

	if parGPU < 0 {
		return errors.New("gpu is invalid data")
	}

	parMem, parErr := StringToInt64(rq.Memory["memoryAmount"])
	if nil != parErr {
		return parErr
	}

	if parMem < 0 {
		return errors.New("memoryAmount is invalid data")
	}

	return nil
}

func CheckForLabelsRequest(params resources.UpdateLabelsParams) error {
	ip := params.ResourcesQuota.IPAddress
	namespace := params.ResourcesQuota.Namespace

	if "" == ip || namespace == "" {
		return errors.New("ip and namespace can't be empty")
	}
	return nil
}

func CheckForAddLabels(params resources.AddLabelsParams) error {
	ip := params.LabelsRequest.IPAddress
	namespace := params.LabelsRequest.Namespace

	if "" == ip || namespace == "" {
		return errors.New("ip and namespace can't be empty")
	}
	return nil
}

func CheckPageParams(page int64, size int64) error {
	if page < 0 || size < 0 {
		return errors.New("page and size must be greater than 0")
	}
	return nil
}

func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func StringToFloat64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
