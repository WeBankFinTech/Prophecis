// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// PostModelResp post model resp
// swagger:model PostModelResp
type PostModelResp struct {

	// model id
	ModelID int64 `json:"model_id,omitempty"`

	// model version
	ModelVersion string `json:"model_version,omitempty"`

	// model version id
	ModelVersionID int64 `json:"model_version_id,omitempty"`
}

// Validate validates this post model resp
func (m *PostModelResp) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PostModelResp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostModelResp) UnmarshalBinary(b []byte) error {
	var res PostModelResp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
