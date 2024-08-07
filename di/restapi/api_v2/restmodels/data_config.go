// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DataConfig data config
// swagger:model DataConfig
type DataConfig struct {

	// 将存储挂载进容器内部的映射路径, 当前为空表示默认的映射
	MappingPath string `json:"mapping_path,omitempty"`

	// 当data_source_type为 MLSSPlatformStorage，需要填充该字段的内容，平台存储根目录
	//
	// 当前输出的存储根目录需要和输入的存储根目录一致
	MlssPlatformParentDir string `json:"mlss_platform_parent_dir,omitempty"`

	// 当data_source_type为 MLSSPlatformStorage，需要填充该字段的内容，平台存储子目录
	MlssPlatformSubDir string `json:"mlss_platform_sub_dir,omitempty"`

	// 数据来源类型
	// Enum: [MLSSPlatformStorage]
	SourceType string `json:"source_type,omitempty"`
}

// Validate validates this data config
func (m *DataConfig) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSourceType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var dataConfigTypeSourceTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["MLSSPlatformStorage"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		dataConfigTypeSourceTypePropEnum = append(dataConfigTypeSourceTypePropEnum, v)
	}
}

const (

	// DataConfigSourceTypeMLSSPlatformStorage captures enum value "MLSSPlatformStorage"
	DataConfigSourceTypeMLSSPlatformStorage string = "MLSSPlatformStorage"
)

// prop value enum
func (m *DataConfig) validateSourceTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, dataConfigTypeSourceTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DataConfig) validateSourceType(formats strfmt.Registry) error {

	if swag.IsZero(m.SourceType) { // not required
		return nil
	}

	// value enum
	if err := m.validateSourceTypeEnum("source_type", "body", m.SourceType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DataConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DataConfig) UnmarshalBinary(b []byte) error {
	var res DataConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
