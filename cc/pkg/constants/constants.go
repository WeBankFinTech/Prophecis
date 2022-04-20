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
package constants

const UserInterceptor = "UserInterceptor"
const IpInterceptor = "IpInterceptor"
const AuthInterceptor = "AuthInterceptor"
const USER_GROUP = "个人用户组"

const (
	AUTH_HEADER_APIKEY    = "MLSS-APPID"
	AUTH_HEADER_TIMESTAMP = "MLSS-APPTimestamp"
	AUTH_HEADER_AUTH_TYPE = "MLSS-Auth-Type"
	AUTH_HEADER_SIGNATURE = "MLSS-APPSignature"
	AUTH_HEADER_TICKET    = "MLSS-Ticket"
	AUTH_HEADER_UIURL     = "MLSS-UIURL"
	AUTH_HEADER_TOKEN     = "MLSS-Token"
	AUTH_HEADER_USERID    = "MLSS-UserID"
	AUTH_HEADER_COOKIE    = "MLSS-Cookie-key"
	AUTH_HEADER_PWD       = "MLSS-Passwd"
	AUTH_REAL_PATH        = "MLSS-RealPath"
	AUTH_REAL_METHOD      = "MLSS-RealMethod"
	SAMPLE_REQUEST_USER   = "MLSS-Sample-User"
	X_ORIGINAL_URI        = "X-Original-URI"
)

const (
	AUTH_RSP_DEPTCODE   = "MLSS-DeptCode"
	AUTH_RSP_USERID     = "MLSS-UserID"
	AUTH_RSP_ORGCODE    = "MLSS-OrgCode"
	AUTH_RSP_LOGINID    = "MLSS-LoginID"
	AUTH_RSP_SUPERADMIN = "MLSS-Superadmin"
	AUTH_HEADER_REAL_IP = "MLSS-RealIP"
	X_REAL_IP           = "X-Real-IP"
)

const (
	NameLengthLimit = 128
	GUIdMinValue    = 0
	GUIdMaxValue    = 90000

	TypePRIVATE = "PRIVATE"
	TypeSYSTEM  = "SYSTEM"
	TypeUser    = "USER"

	NodeSelectorKey = "scheduler.alpha.kubernetes.io/node-selector"
	NVIDIAGPU       = "NVIDIAGPU"
	LB_GPU_MODEL    = "lb-gpu-model"
	LB_BUS_TYPE     = "lb-bus-type"

	GA     = "GA"
	GU     = "GU"
	GA_ID  = 1
	GU_ID  = 2
	USER   = "USER"
	MLSS   = "mlss"
	ADMIN  = "ADMIN"
	VIEWER = "VIEWER"
)

const (
	PrivateGroupNameRegx = "^gp(-private)(-[_a-z0-9]+)$"
	SystemGroupNameRegx  = "^gp(-[a-z0-9]+){2}$"
	//NamespaceNameRegx    = "^ns(-[a-z0-9]+){5}(-sf|-dmz|-csf|-cdmz|-safe|-sfsafe)$"
	NamespaceNameRegx    = "^ns(-[a-z0-9]+){2}"
	SystemNamespaceName  = "^(default)|(kube.*)|(monitoring)|(rook.*)|(system.*)|(admin.*)" +
		"|(DEFAULT)|(KUBE.*)|(MONITORING)|(ROOK.*)|(SYSTEM.*)|(ADMIN.*)$"
	UsernameRegx = "^[a-zA-Z0-9]+|[a-zA-Z0-9]+[_][c]|[v][_][a-zA-Z]+|[v][_][a-zA-Z]+[_][c]$"
)

const (
	NBApiGroup   = "kubeflow.org"
	NBApiVersion = "v1alpha1"
	NBApiPlural  = "notebooks"
)

const (
	GROUP_SUB_NAME = "gp-private-"
	PRO_USER_PATH  = "/data/prophecis/prophecis-data"
)
