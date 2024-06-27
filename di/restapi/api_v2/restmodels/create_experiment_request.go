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

// CreateExperimentRequest create experiment request
// swagger:model CreateExperimentRequest
type CreateExperimentRequest struct {

	// 实验所属的集群类型，默认为"BDAP"；可选枚举值为 "BDAP" "BDP"
	// Enum: [BDAP BDP]
	ClusterType *string `json:"cluster_type,omitempty"`

	// 创建实验的描述
	Description string `json:"description,omitempty"`

	// DSS相关的信息，如果实验的origin_system是DSS，则需要传入该字段
	DssInfo *DSSInfo `json:"dss_info,omitempty"`

	// 创建的实验的工作流
	FlowJSON string `json:"flow_json,omitempty"`

	// 创建的实验的工作流上传的id（来自UploadExperimentFlowJson接口的返回）
	FlowJSONUploadID string `json:"flow_json_upload_id,omitempty"`

	// 创建实验所属的项目组ID(理论上传group_id就够了，但是group的相关接口暂时不支持byId查询，所以暂时需要传group_id和group_name)
	GroupID string `json:"group_id,omitempty"`

	// 创建实验所属的项目组名字(理论上传group_id就够了，但是group的相关接口暂时不支持byId查询，所以暂时需要传group_id和group_name)
	GroupName string `json:"group_name,omitempty"`

	// 创建实验的名称
	// Required: true
	// Min Length: 1
	Name *string `json:"name"`

	// 表示实验在哪个系统创建的, 默认为"MLSS"；可选的枚举值为 "WTSS" "DSS" "MLSS"
	// Enum: [MLSS DSS WTSS]
	SourceSystem *string `json:"source_system,omitempty"`

	// 创建实验的标签
	Tags []string `json:"tags"`
}

// Validate validates this create experiment request
func (m *CreateExperimentRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClusterType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDssInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSourceSystem(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var createExperimentRequestTypeClusterTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["BDAP","BDP"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createExperimentRequestTypeClusterTypePropEnum = append(createExperimentRequestTypeClusterTypePropEnum, v)
	}
}

const (

	// CreateExperimentRequestClusterTypeBDAP captures enum value "BDAP"
	CreateExperimentRequestClusterTypeBDAP string = "BDAP"

	// CreateExperimentRequestClusterTypeBDP captures enum value "BDP"
	CreateExperimentRequestClusterTypeBDP string = "BDP"
)

// prop value enum
func (m *CreateExperimentRequest) validateClusterTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, createExperimentRequestTypeClusterTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CreateExperimentRequest) validateClusterType(formats strfmt.Registry) error {

	if swag.IsZero(m.ClusterType) { // not required
		return nil
	}

	// value enum
	if err := m.validateClusterTypeEnum("cluster_type", "body", *m.ClusterType); err != nil {
		return err
	}

	return nil
}

func (m *CreateExperimentRequest) validateDssInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.DssInfo) { // not required
		return nil
	}

	if m.DssInfo != nil {
		if err := m.DssInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("dss_info")
			}
			return err
		}
	}

	return nil
}

func (m *CreateExperimentRequest) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	if err := validate.MinLength("name", "body", string(*m.Name), 1); err != nil {
		return err
	}

	return nil
}

var createExperimentRequestTypeSourceSystemPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["MLSS","DSS","WTSS"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		createExperimentRequestTypeSourceSystemPropEnum = append(createExperimentRequestTypeSourceSystemPropEnum, v)
	}
}

const (

	// CreateExperimentRequestSourceSystemMLSS captures enum value "MLSS"
	CreateExperimentRequestSourceSystemMLSS string = "MLSS"

	// CreateExperimentRequestSourceSystemDSS captures enum value "DSS"
	CreateExperimentRequestSourceSystemDSS string = "DSS"

	// CreateExperimentRequestSourceSystemWTSS captures enum value "WTSS"
	CreateExperimentRequestSourceSystemWTSS string = "WTSS"
)

// prop value enum
func (m *CreateExperimentRequest) validateSourceSystemEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, createExperimentRequestTypeSourceSystemPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *CreateExperimentRequest) validateSourceSystem(formats strfmt.Registry) error {

	if swag.IsZero(m.SourceSystem) { // not required
		return nil
	}

	// value enum
	if err := m.validateSourceSystemEnum("source_system", "body", *m.SourceSystem); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CreateExperimentRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateExperimentRequest) UnmarshalBinary(b []byte) error {
	var res CreateExperimentRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
