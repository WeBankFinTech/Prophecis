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

type GroupBase struct {
	ID             int64  `gorm:"column:id; PRIMARY_KEY" json:"id"`
	EnableFlag     int8   `json:"enable_flag"`
	Name           string `json:"name"`
	GroupType      string `json:"group_type"`
	SubsystemId    int64  `json:"subsystem_id"`
	SubsystemName  string `json:"subsystem_name"`
	Remarks        string `json:"remarks"`
	DepartmentId   string `json:"department_id"`
	DepartmentName string `json:"department_name"`
	ClusterName    string `json:"cluster_name"`
	RmbIdc         string `json:"rmb_idc"`
	RmbDcn         string `json:"rmb_dcn"`
	ServiceId      int64  `json:"service_id"`
}
