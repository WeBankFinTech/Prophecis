package models

import "time"

//id                    bigint auto_increment primar
//enable_flag           tinyint(1) default 1 not nul
//exp_id                bigint               null,
//exp_name              varchar(255)         null,
//exp_desc              varchar(255)         null,
//dss_project_id        bigint               null,
//dss_workspace_id      bigint               null,
//dss_flow_id           bigint               null,
//dss_flow_last_version varchar(255)         null,
//dss_bml_last_version  varchar(255)         null,
//exp_create_time       timestamp            null,
//exp_modify_time       timestamp            null,
//exp_create_user_id    bigint               null,
//exp_modify_user_id    bigint               null

type Experiment struct {
	BaseModel
	ExpDesc                    *string         `json:"exp_desc"`
	ExpName                    *string         `json:"exp_name"`
	GroupName                  string          `json:"group_name"`
	MlflowExpID                *string         `json:"mlflow_exp_id" gorm:"column:mlflow_exp_id"`
	DssProjectID               int64           `json:"dss_project_id"`
	DssProjectName             string          `json:"dss_project_name"`
	DssProjectVersionID        int64           `json:"dss_project_version"`
	DssWorkspaceID             int64           `json:"dss_workspace_id"`
	DssFlowID                  int64           `json:"dss_flow_id"`
	DssFlowLastVersion         *string         `json:"dss_flow_last_version"`
	DssBmlLastVersion          *string         `json:"dss_bml_last_version"`
	ExpCreateTime              time.Time       `json:"exp_create_time"`
	ExpModifyTime              time.Time       `json:"exp_modify_time"`
	ExpCreateUserID            int64           `json:"exp_create_user_id"`
	ExpModifyUserID            int64           `json:"exp_modify_user_id"`
	CreateType                 *string         `json:"create_type"`
	DssDssProjectId            int64           `json:"dss_dss_project_id"`
	DssDssProjectName          *string         `json:"dss_dss_project_name"`
	DssDssFlowVersion          *string         `json:"dss_dss_flow_version"`
	DssDssFlowName             *string         `json:"dss_dss_flow_name"`
	DssDssFlowId               int64           `json:"dss_dss_flow_id"`
	DssDssFlowProjectVersionId int64           `json:"dss_dss_flow_project_version_id"`
	TagList                    []ExperimentTag `gorm:"ForeignKey:exp_id;AssociationForeignKey:exp_id" json:"tag_list"`
	CreateUser                 User            `gorm:"ForeignKey:ExpCreateUserID;AssociationForeignKey:id" json:"user"`
	ModifyUser                 User            `gorm:"ForeignKey:ExpModifyUserID;AssociationForeignKey:id" json:"user"`
}

type ExperimentTag struct {
	BaseModel
	TagID              int64     `gorm:"column:id;" json:"tag_id,omitempty"`
	ExpID              int64     `json:"exp_id,omitempty"`
	ExpTag             *string   `json:"exp_tag"`
	ExpTagCreateTime   time.Time `json:"exp_tag_create_time"`
	ExpTagCreateUserID int64     `json:"exp_tag_create_user_id"`
}
