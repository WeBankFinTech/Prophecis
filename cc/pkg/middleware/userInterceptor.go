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
	"mlss-controlcenter-go/pkg/logger"
)

func UserInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		//logger.Logger().Infof("UserInterceptor c.Request.RequestURI: %v", c.Request.RequestURI)

		//check path
		if common.CheckPattern(constants.UserInterceptor, c.Request.RequestURI) {

			logger.Logger().Infof("check user in UserInterceptor")

			if c.GetHeader("token") != "123" && c.GetHeader("token") != "" {
				c.AbortWithStatusJSON(401, gin.H{"msg": "token is invalid"})
				return
			}

			if c.GetHeader("user") != "buoy" && c.GetHeader("pwd") != "123" {
				c.AbortWithStatusJSON(401, gin.H{"msg": "user not match pwd"})
				return
			}
		}
		//next interceptor
		//logger.Logger().Infof("UserInterceptor Next!")
		c.Next()

	}
}
