package third

import "time"

type ProcessLineVersion struct {
	//1
	ProcesslineVersionId string `json:"processline_version_id"`
	//2
	ProcesslineId string `json:"processline_id"`
	//3
	SourceType string `json:"source_type"`
	//4
	VersionNum string `json:"version_num"`
	//5
	ProcesslineStorageConfiguration string `json:"processline_storage_configuration"`
	//6.
	Options string `json:"options"`
	//7.
	ProcesslineVersionStatus string `json:"processline_version_status"`
	//8.
	ProcesslineVersionDesc string `json:"process_line_version_desc"`
	//9
	IsEncrypted int64 `json:"is_encrypted"`
	//10
	Remarks string `json:"remarks"`
	//11
	EnableFlag int64 `json:"enable_flag"`
	//12
	CreateUser string `json:"create_user"`
	//13
	CreateTime time.Time `json:"create_time"`
	//14.
	UpdateUser string `json:"update_user"`
	//15.
	UpdateTime time.Time `json:"update_time"`
	//16
	ProcesslineName string `json:"processline_name"`
	//17
	//ProcesslineVersion string `json:"processline_version"`
	//18
	GroupId string `json:"group_id"`
	//19
	VariableInfo string `json:"variable_info"`
	//20
	ConfigInfo string `json:"config_info"`

	FpsFileId string `json:"fps_file_id"`

	FpsFileHash string `json:"fps_file_hash"`

	//StoragePath string `json:"storage_path"`

	Request string `json:"request"`

	Response string `json:"response"`
}
