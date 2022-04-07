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
package middleware

import (
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
	"mlss-controlcenter-go/pkg/logger"

	//"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func IpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		//url := location.Get(c)
		//logger.Logger().Infof("AuthRequired url.Host: %v, url.Path: %v, url.Scheme: %v, url.RequestURI: %v", url.Host, url.Path, url.Scheme,url.RequestURI())
		//logger.Logger().Infof("AuthRequired url: %v", url)
		//logger.Logger().Infof("IpInterceptor c.Request.RequestURI: %v", c.Request.RequestURI)

		//check path
		if common.CheckPattern(constants.IpInterceptor, c.Request.RequestURI) {
			logger.Logger().Infof("check ip addr in IpInterceptor")
		}
		//next interceptor
		c.Next()

	}
}
