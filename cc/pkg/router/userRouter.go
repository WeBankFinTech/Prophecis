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

type UserRouter struct {
	//Root *gin.RouterGroup
}

func init() {
	//RouterList = append(RouterList, UserRouter{})
	InitUserRouter()
}

func InitUserRouter() {
	//router := UserRouter{}
	//router.Init(Root)
	//return router
	//user := Root.Group("/v1/user")
	//user.GET("", controller.GetUser)
	//user.GET("/:userId", controller.GetUserByUserId)
	//user.GET("/:groupId", controller.GetUsersByGroupId)
	//user.GET("/users", controller.GetAllUsers)
	//user.POST("", controller.AddUser)
	//user.POST("/join", controller.JoinGroup)
}

//func (r UserRouter) Init(*gin.RouterGroup) {
//	user := Root.Group("/v1/user")
//	user.GET("/", controller.GetUser)
//	user.POST("/", controller.AddUser)
//	user.POST("/join", controller.JoinGroup)
//}
