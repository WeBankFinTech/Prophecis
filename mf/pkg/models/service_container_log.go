// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ServiceContainerLog service container log
// swagger:model ServiceContainerLog
type ServiceContainerLog struct {

	// log
	Log string `json:"log,omitempty"`
}

// Validate validates this service container log
func (m *ServiceContainerLog) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ServiceContainerLog) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServiceContainerLog) UnmarshalBinary(b []byte) error {
	var res ServiceContainerLog
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
