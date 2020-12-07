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
	"github.com/gin-gonic/gin"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/constants"
)

func AuthInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		//t := time.Now()
		// Set example variable
		//c.Set("pwd", "666")

		//logger.Logger().Infof("AuthRequired RequestURI: %v", c.Request.RequestURI)

		// before request

		// after request
		//latency := time.Since(t)
		//log.Print(latency)

		// access the status we are sending
		//status := c.Writer.Status()
		//log.Println(status)

		//pwd, _ := c.Get("pwd")
		//logger.Logger().Infof("AuthRequired pwd: %v", pwd)

		if common.CheckPattern(constants.AuthInterceptor, c.Request.RequestURI) {
			logger.Logger().Infof("check Auth in AuthInterceptor")
		}
		c.Next()
	}
}
