// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// PatchExperimentRequest patch experiment request
// swagger:model PatchExperimentRequest
type PatchExperimentRequest struct {

	// 修改的实验的发布配置
	DeploySetting *DeploySetting `json:"deploy_setting,omitempty"`

	// 修改实验的描述.
	Description string `json:"description,omitempty"`

	// 修改的实验的工作流
	FlowJSON string `json:"flow_json,omitempty"`

	// 修改实验所属的项目组ID.(理论上传group_id就够了，但是group的相关接口暂时不支持byId查询，所以暂时需要传group_id和group_name)
	GroupID string `json:"group_id,omitempty"`

	// 修改实验所属的项目组名字.(理论上传group_id就够了，但是group的相关接口暂时不支持byId查询，所以暂时需要传group_id和group_name)
	GroupName string `json:"group_name,omitempty"`

	// 修改实验的名称.
	Name string `json:"name,omitempty"`

	// 修改的实验的标签
	Tags []string `json:"tags"`
}

// Validate validates this patch experiment request
func (m *PatchExperimentRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDeploySetting(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PatchExperimentRequest) validateDeploySetting(formats strfmt.Registry) error {

	if swag.IsZero(m.DeploySetting) { // not required
		return nil
	}

	if m.DeploySetting != nil {
		if err := m.DeploySetting.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("deploy_setting")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PatchExperimentRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PatchExperimentRequest) UnmarshalBinary(b []byte) error {
	var res PatchExperimentRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
