// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// GPUNodeDataStoresSubDirectory g p u node data stores sub directory
// swagger:model GPUNodeDataStoresSubDirectory
type GPUNodeDataStoresSubDirectory struct {

	// 子目录的路径
	Container string `json:"container,omitempty"`
}

// Validate validates this g p u node data stores sub directory
func (m *GPUNodeDataStoresSubDirectory) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GPUNodeDataStoresSubDirectory) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GPUNodeDataStoresSubDirectory) UnmarshalBinary(b []byte) error {
	var res GPUNodeDataStoresSubDirectory
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
