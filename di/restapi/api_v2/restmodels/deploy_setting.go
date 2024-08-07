// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// DeploySetting deploy setting
// swagger:model DeploySetting
type DeploySetting struct {

	// 是否可发布为离线服务
	CanDeployAsOfflineService *bool `json:"can_deploy_as_offline_service,omitempty"`
}

// Validate validates this deploy setting
func (m *DeploySetting) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeploySetting) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeploySetting) UnmarshalBinary(b []byte) error {
	var res DeploySetting
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
