// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Fields TFOS Archives & PyFile Fields
// swagger:model Fields
type Fields struct {

	// HDFS Path.
	Hdfs string `json:"hdfs,omitempty"`

	// BML Resource ID.
	ResourceID string `json:"resource_id,omitempty"`

	// BML Resource Version.
	Version string `json:"version,omitempty"`
}

// Validate validates this fields
func (m *Fields) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Fields) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Fields) UnmarshalBinary(b []byte) error {
	var res Fields
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
