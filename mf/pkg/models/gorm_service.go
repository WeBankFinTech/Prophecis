package models

import (
	"time"
)

type GormService struct{
	ModelService
	ModelGroup
	GormModel
	User
	ModelVersionService
	Status string `json:"status"`
	Image string `json:"image"`
}

type GormModel struct {
	Id int64 `json:"id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"id"`
	CreationTimestamp time.Time `json:"creation_timestamp,omitempty" gorm:"column:creation_timestamp"`
	UpdateTimestamp time.Time `json:"update_timestamp,omitempty" gorm:"column:creation_timestamp"`
	ModelName string `json:"model_name,omitempty" gorm:"column:model_name"`
	ModelType string `json:"model_type,omitempty" gorm:"column:model_type"`
	Reamrk string `json:"reamrk,omitempty" gorm:"column:reamrk"`
}

type GormModelVersion struct {
	Id int64 `json:"id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"id"`
	Version string `json:"version,omitempty" gorm:"column:version"`
	LatestFlag int64 `json:"latest_flag,omitempty" gorm:"column:latest_flag"`
	Source string `json:"source,omitempty" gorm:"column:source"`
	Params string `json:"params,omitempty" gorm:"column:params"`
	ModelID int64 `json:"model_id,omitempty" gorm:"column:model_id"`
}

type User struct {
	Id int64 `json:"id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"id"`
	GroupName string `json:"group_name,omitempty" gorm:"column:name"`
}

type ModelGroup struct {
	Id int64 `json:"group_id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"group_id"`
	Username string `json:"user_name,omitempty" gorm:"column:name"`
}

type ModelService struct {
	Id int64 `json:"service_id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"id"`
	ServiceName string `json:"service_name,omitempty" gorm:"column:service_name"`
	Type string `json:"type,omitempty" gorm:"column:type"`
	Cpu float64 `json:"cpu,omitempty" gorm:"column:cpu"`
	Memory float64 `json:"memory,omitempty" gorm:"column:memory"`
	Gpu float64 `json:"gpu,omitempty" gorm:"column:gpu"`
	Namespace string `json:"namespace,omitempty" gorm:"column:namespace"`
	LastUpdatedTimestamp string `json:"last_updated_timestamp,omitempty" gorm:"column:last_updated_timestamp"`
	LogPath string `json:"log_path,omitempty" gorm:"column:log_path"`
	Remark string `json:"remark,omitempty" gorm:"column:remark"`
	EndpointType string `json:"endpoint_type,omitempty" gorm:"column:endpoint_type"`
}

type ModelVersionService struct {
	Id int64 `json:"id,omitempty" gorm:"column:id; PRIMARY_KEY" json:"id"`
	ModelVersionId int64 `json:"modelversion_id,omitempty" gorm:"column:modelversion_id"`
	ModelVersion GormModelVersion `gorm:"ForeignKey:ModelVersionId;AssociationForeignKey:id" json:"modelversion"`
}

type Namespace struct {
	Id int64
	GroupId int64
	Namespace string
}

type GormServiceModelVersion struct {
	Id int64 `gorm:"column:id; PRIMARY_KEY" json:"id"`
	ServiceId int64 `gorm:"column:service_id"`
	ModelVersionId int64 `gorm:"column:modelversion_id"`
	EnableFlag bool `json:"enable_flag"`
}


type PageModel struct {
	// models
	Models    interface{} `json:"models"`
	// page number
	PageNumber int64 `json:"pageNumber,omitempty"`

	// page size
	PageSize int64 `json:"pageSize,omitempty"`

	// total
	Total int64 `json:"total,omitempty"`

	// total page
	TotalPage int64 `json:"totalPage,omitempty"`
}