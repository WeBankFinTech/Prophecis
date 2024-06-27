package models

import "time"

type Experiment struct {
	ID                                string    `json:"id"`
	Name                              string    `json:"name"`
	Description                       string    `json:"description"`
	GroupID                           string    `json:"group_id"`
	GroupName                         string    `json:"group_name"`
	VersionName                       string    `json:"version_name"`
	SourceSystem                      string    `json:"source_system"`
	ClusterType                       string    `json:"cluster_type"`
	CreateUser                        string    `json:"create_user"`
	CreateTime                        time.Time `json:"create_time"`
	UpdateUser                        string    `json:"update_user"`
	UpdateTime                        time.Time `json:"update_time"`
	LatestExecuteTime                 time.Time `json:"latest_execute_time"`
	Status                            string    `json:"status"` // 表示实验的状态（归档，活跃等）
	Tags                              string    `json:"tags"`   // 逗号隔开的tag list string
	KfpPipelineId                     string    `json:"kfp_pipeline_id"`
	KfpPipelineVersionId              string    `json:"kfp_pipeline_version_id"`
	CanDeployAsOfflineService         bool      `json:"can_deploy_as_offline_service"`
	CanDeployAsOfflineServiceOperator string    `json:"can_deploy_as_offline_service_operator"`

	DSSWorkspaceID   int64  `json:"dss_workspace_id"`
	DSSWorkspaceName string `json:"dss_workspace_name"`
	DSSProjectID     int64  `json:"dss_project_id"`
	DSSProjectName   string `json:"dss_project_name"`
	DSSFlowID        int64  `json:"dss_flow_id"`
	DSSFlowName      string `json:"dss_flow_name"`
	DSSFlowVersion   string `json:"dss_flow_version"`
	DSSPublished     bool   `json:"dss_published"`
}

type ExperimentStatus string

const ExperimentStatusArchive ExperimentStatus = "Archive"
const ExperimentStatusActive ExperimentStatus = "Active"
const ExperimentStatusDeleted ExperimentStatus = "Deleted"
