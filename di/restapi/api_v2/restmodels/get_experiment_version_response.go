// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetExperimentVersionResponse get experiment version response
//
// swagger:model GetExperimentVersionResponse
type GetExperimentVersionResponse struct {

	// 实验的id
	ID string `json:"id,omitempty"`
}

// Validate validates this get experiment version response
func (m *GetExperimentVersionResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetExperimentVersionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetExperimentVersionResponse) UnmarshalBinary(b []byte) error {
	var res GetExperimentVersionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
