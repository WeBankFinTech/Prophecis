package models

import "time"

type ExperimentRun struct {
	ExpRunID       string    `gorm:"column:id" json:"id"`
	Name           string    `json:"name"`
	CreateUser     string    `json:"create_user"`
	CreateTime     time.Time `json:"create_time"`
	ExpID          string    `json:"exp_id"`
	ExpVersionName string    `json:"exp_version_name"`
	KfpRunId       string    `json:"kfp_run_id"`
	EndTime        time.Time `json:"end_time"`
	ExecuteStatus  string    `json:"execute_status"` // 表示实验执行的状态（初始化，运行，成功，失败等）
	FlowJson       string    `json:"flow_json"`
}

type ExperimentRunNode struct {
	ID              string        `json:"id"`
	ExperimentRunID string        `json:"experiment_run_id"`
	Key             string        `json:"key"`      // 节点的名字, 对应flowjson的node.key
	Title           string        `json:"title"`    // 节点的名字, 对应flowjson的node.title
	Type            string        `json:"type"`     // 节点的类型，比如gpu节点，模型预测节点
	State           string        `json:"state"`    // 节点的状态
	ProxyID         string        `json:"proxy_id"` // todo: 代理id，当节点是个代理的时候，它会去第三方创建一个任务，服务，第三方的那个任务，服务的id就是代理id
	ExperimentRun   ExperimentRun `gorm:"ForeignKey:ExperimentRunID;AssociationForeignKey:id" json:"experiment_run"`
}
