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

package router

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"mlss-controlcenter-go/pkg/common"
	"mlss-controlcenter-go/pkg/middleware"
)

var route = GetDefaultRoute()
var Root = GetRoot()

//var RouterList []IRouter

func GetDefaultRoute() *gin.Engine {
	return gin.Default()
}

func GetRoot() *gin.RouterGroup {
	return route.Group(common.GetAppConfig().Server.Servlet.ContextPath, middleware.UserInterceptor(), middleware.IpInterceptor(), middleware.AuthInterceptor())
}

type IRouter interface {
	Init(*gin.RouterGroup)
}

func GetRoute() *gin.Engine {
	//location middleware
	route.Use(location.Default())

	//InitRouter(UserRouter{}, GroupRouter{})
	//InitAllRouter(RouterList...)

	return route
}

func InitAllRouter(routers ...IRouter) {
	for _, v := range routers {
		v.Init(Root)
	}
}

func testNotebook(c *gin.Context) {
	c.String(200, "aha")
}
