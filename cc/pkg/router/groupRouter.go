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
	"mlss-controlcenter-go/pkg/controller"
)

type GroupRouter struct {
	//Root *gin.RouterGroup
}

func init() {
	//RouterList = append(RouterList, GroupRouter{})
	InitGroupRouter()
}

func InitGroupRouter() {
	//router := GroupRouter{}
	//router.Init(Root)
	group := Root.Group("/v1/group")
	group.GET("", controller.GetGroup)
}

//func (r GroupRouter) Init(*gin.RouterGroup) {
//	group := Root.Group("/group")
//	group.GET("/", controller.GetGroup)
//}
