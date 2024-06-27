package third

import "time"

type ModelVersion struct {
	ModelVersionBase
	Model Model `gorm:"ForeignKey:ModelID;AssociationForeignKey:model_id" json:"model"`
}

type ModelVersionBase struct {
	ModelVersionID            string    `gorm:"column:model_version_id; primaryKey" json:"model_version_id"`
	ModelID                   string    `json:"model_id"`
	VersionNum                string    `json:"version_num"`
	SourceType                string    `json:"source_type"`
	FrameworkType             string    `json:"framework_type"`
	FormatType                string    `json:"format_type"`
	AlgorithmType             string    `json:"algorithm_type"`
	ImageID                   string    `json:"image_id"`
	ModelStorageConfiguration string    `json:"model_storage_configuration"`
	ModelVersionStatus        string    `json:"model_version_status"`
	ModelVersionDesc          string    `json:"model_version_desc"`
	Options                   string    `json:"options"`
	IsEncrypted               int8      `json:"is_encrypted"`
	FpsFileID                 string    `json:"fps_file_id"`
	FpsFileHash               string    `json:"fps_file_hash"`
	IsPsz                     bool      `json:"is_psz"`
	Remarks                   string    `json:"remarks"`
	EnableFlag                int8      `json:"enable_flag"`
	CreateUser                string    `json:"create_user"`
	CreateTime                time.Time `json:"create_time"`
	UpdateUser                string    `json:"update_user"`
	UpdateTime                time.Time `json:"update_time"`
}

type Model struct {
	ModelBase
	//CreateUserBean     UserBase          `gorm:"ForeignKey:CreateUser;AssociationForeignKey:user_name" json:"create_user_bean"`
	//UpdateUserBean     UserBase          `gorm:"ForeignKey:UpdateUser;AssociationForeignKey:user_name" json:"update_user_bean"`
	//ModelLatestVersion ModelVersionBase  `gorm:"ForeignKey:LatestVersionID;AssociationForeignKey:model_version_id" json:"model_latest_version"`
	//ModelPermission    []ModelPermission `gorm:"ForeignKey:ModelID;AssociationForeignKey:model_id" json:"model_permission"`
}

type ModelBase struct {
	ModelID            string    `gorm:"column:model_id; primaryKey" json:"model_id"`
	ModelName          string    `gorm:"column:model_name; unique" json:"model_name"`
	LatestVersionId    string    `json:"latest_version_id"`
	PlatformType       string    `json:"platform_type"`
	ModelUsage         string    `json:"model_usage"`
	ModelStatus        int8      `json:"model_status"`
	ModelDoc           string    `json:"model_doc"`
	ModelDesc          string    `json:"model_desc"`
	ModelOrigin        string    `json:"model_origin"`
	ModelAccessibility string    `json:"model_accessibility"`
	IsGpu              bool      `json:"is_gpu"`
	Remarks            string    `json:"remarks"`
	EnableFlag         int8      `json:"enable_flag"`
	CreateUser         string    `json:"create_user"`
	CreateTime         time.Time `json:"create_time"`
	UpdateUser         string    `json:"update_user"`
	UpdateTime         time.Time `json:"update_time"`
}

type UserBase struct {
	UserID            string    `gorm:"column:user_id; primaryKey" json:"user_id"`
	UserType          string    `json:"user_type"`
	MlssClusterName   string    `json:"mlss_cluster_name"`
	UserName          string    `json:"user_name"`
	UserCnName        string    `json:"user_cn_name"`
	UserConfiguration string    `json:"user_configuration"`
	Remarks           string    `json:"remarks"`
	EnableFlag        int8      `json:"enable_flag"`
	UpdateUser        string    `json:"update_user"`
	CreateUser        string    `json:"create_user"`
	UpdateTime        time.Time `json:"update_time"`
	CreateTime        time.Time `json:"create_time"`
}

type ModelPermission struct {
	ModelID    string    `gorm:"column:model_id; primaryKey" json:"model_id"`
	GroupID    string    `gorm:"column:group_id; primaryKey" json:"group_id"`
	Permission string    `json:"permission"`
	CreateUser string    `json:"create_user"`
	CreateTime time.Time `json:"create_time"`
	UpdateUser string    `json:"update_user"`
	UpdateTime time.Time `json:"update_time"`
	Group      GroupBase `gorm:"ForeignKey:GroupID;AssociationForeignKey:group_id" json:"group"`
}

type GroupBase struct {
	GroupID             string    `gorm:"PRIMARY_KEY" json:"group_id"`
	GroupName           string    `json:"group_name"`
	MlssClusterName     string    `json:"mlss_cluster_name"`
	IsCrossDepartmental int8      `json:"is_cross_departmental"`
	Remarks             string    `json:"remarks"`
	GroupStatus         int8      `json:"group_status"`
	EnableFlag          int8      `json:"enable_flag"`
	CreateUser          string    `json:"create_user"`
	CreateTime          time.Time `json:"create_time"`
	UpdateUser          string    `json:"update_user"`
	UpdateTime          time.Time `json:"update_time"`
}

func (ModelBase) TableName() string {
	return "t_model_v2"
}

func (Model) TableName() string {
	return "t_model_v2"
}
