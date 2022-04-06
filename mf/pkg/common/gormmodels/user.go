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

package gormmodels

type User struct {
	BaseModel
	Name      string `json:"name"`
	Gid       int64  `json:"gid"`
	UID       int64  `json:"uid"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	Remarks   string `json:"remarks"`
	GuidCheck int64  `json:"guid_check"`
}

type ImageUser struct {
	BaseModel
	Name string `json:"name"`
}

func (ImageUser) TableName() string {
	return "user"
}

func (User) TableName() string {
	return "user"
}
