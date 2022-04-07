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

type Group struct {
	BaseModel
	Name           string `json:"name"`
	GroupType      string `json:"group_type"`
	SubsystemId    int64  `json:"subsystem_id"`
	SubsystemName  string `json:"subsystem_name"`
	Remarks        string `json:"remarks"`
	DepartmentId   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	ClusterName    string `json:"cluster_name"`

	// add by version 1.16.0
	RmbIdc    string `json:"rmb_idc"`
	RmbDcn    string `json:"rmb_dcn"`
	ServiceId int64  `json:"service_id"`
}

type UserGroup struct {
	BaseModel
	GroupID int64  `json:"group_id"`
	RoleID  int64  `json:"role_id"`
	Remarks string `json:"remarks"`
	UserID  int64  `json:"user_id"`
	Group   Group  `gorm:"ForeignKey:GroupID;AssociationForeignKey:id" json:"group"`
	Role    Role   `gorm:"ForeignKey:RoleID;AssociationForeignKey:id" json:"role"`
	User    User   `gorm:"ForeignKey:UserID;AssociationForeignKey:id" json:"user"`

	// add by version 1.16.0
	RmbIdc    string `json:"rmb_idc"`
	RmbDcn    string `json:"rmb_dcn"`
	ServiceId int64  `json:"service_id"`
}

type ImageGroup struct {
	BaseModel
	Name string `json:"name"`
}

type ResultGroup struct {
	BaseModel
	GroupID int64  `json:"group_id"`
	RoleID  int64  `json:"role_id"`
	Remarks string `json:"remarks"`
	UserID  int64  `json:"user_id"`
}

func (ImageGroup) TableName() string {
	return "group"
}
