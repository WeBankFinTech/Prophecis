// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ModelParameters model parameters
// swagger:model ModelParameters
type ModelParameters struct {

	// params name
	Name string `json:"name,omitempty"`

	// params type, STRING
	Type string `json:"type,omitempty"`

	// params value
	Value string `json:"value,omitempty"`
}

// Validate validates this model parameters
func (m *ModelParameters) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModelParameters) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelParameters) UnmarshalBinary(b []byte) error {
	var res ModelParameters
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
