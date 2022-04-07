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
	"mlss-controlcenter-go/pkg/logger"
	"mlss-controlcenter-go/pkg/models"
	"os"
	"regexp"
	"strings"
)

func CheckURLPattern(pathPattern string, requestPath string) bool {

	pathPattern = strings.Replace(pathPattern, "**", "(.*)", -1)
	match, _ := regexp.MatchString(pathPattern, requestPath)

	return match
}

func CheckPattern(interceptorName string, requestPath string) bool {
	var interceptor models.InterceptorConfig

	prefix := strings.HasPrefix(requestPath, GetAppConfig().Server.Servlet.ContextPath)

	if prefix {
		requestPath = strings.TrimPrefix(requestPath, GetAppConfig().Server.Servlet.ContextPath)
	}

	config := GetAppConfig()
	for _, v := range config.Core.Interceptor.Configs {
		logger.Logger().Debugf("inr name: %v", v.Name)

		if v.Name == interceptorName {
			interceptor = v
			break
		}
	}

	if interceptor.Name == "" {
		panic(interceptorName + " not exists")
	}

	if !CheckAddList(interceptor.Add, requestPath) {
		return false
	}
	if CheckExcludeList(interceptor.Exclude, requestPath) {
		return false
	}

	return true
}

func CheckAddList(pathPattern []string, requestPath string) bool {
	for _, v := range pathPattern {
		if CheckURLPattern(v, requestPath) {
			return true
		}
	}
	return false
}

func CheckAdd(pathPattern string, requestPath string) bool {
	return false
}

func CheckExcludeList(pathPattern []string, requestPath string) bool {
	for _, v := range pathPattern {
		if CheckURLPattern(v, requestPath) {
			return true
		}
	}
	return false
}

func CheckExclude(pathPattern string, requestPath string) bool {
	return false

}

func PathExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true, nil
		}
		return false, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
