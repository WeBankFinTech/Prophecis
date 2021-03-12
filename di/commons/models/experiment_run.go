package models

import "time"

//id                       bigint auto_increment primary key,
//enable_flag              tinyint(1) default 1,

//exp_exec_id              bigint       null,
//exp_id                   bigint       null,
//dss_exec_id              varchar(255) null,
//dss_flow_id              bigint       null,
//dss_flow_last_version_id bigint       null,
//exp_exec_type            varchar(255) null,
//exp_exec_status          varchar(255) null,
//exp_run_create_time      bigint       null,
//exp_run_end_time         bigint       null,
//exp_run_create_user_id   varchar(255) null,
//exp_run_modify_user_id   varchar(255) null

type ExperimentRun struct {
	BaseModel
	//ExpExecId                 int64  `json:"exp_exec_id"`
	ExpID              int64     `json:"exp_id"`
	DssExecID          string    `json:"dss_exec_id"`
	DssTaskID          int64    `json:"dss_task_id"`
	DssFlowID          int64     `json:"dss_flow_id"`
	DssFlowLastVersion string    `json:"dss_flow_last_version_version"`
	ExpExecType        string    `json:"exp_exec_type"`
	ExpExecStatus      string    `json:"exp_exec_status"`
	ExpRunCreateTime   time.Time `json:"exp_run_create_time"`
	ExpRunEndTime      time.Time `json:"exp_run_end_time"`
	ExpRunCreateUserID int64    `json:"exp_run_create_user_id"`
	ExpRunModifyUserID int64    `json:"exp_run_modify_user_id"`
	Experiment Experiment `gorm:"ForeignKey:ExpID" json:"experiment"`
	User User `gorm:"ForeignKey:ExpRunCreateUserID;AssociationForeignKey:id" json:"user"`
}

const (
	RUN_STATUS_WAITFORRETRY = "WaitForRetry"
	RUN_STATUS_INITED       = "Inited"
	RUN_STATUS_SCHEDULED    = "Scheduled"
	RUN_STATUS_RUNNING      = "Running"
	RUN_STATUS_SUCCEED      = "Succeed"
	RUN_STATUS_FAILED       = "Failed"
	RUN_STATUS_CANCELLED    = "Cancelled"
	RUN_STATUS_TIMEOUT      = "Timeout"
)
