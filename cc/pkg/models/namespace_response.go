// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// NamespaceResponse namespace response
// swagger:model NamespaceResponse
type NamespaceResponse struct {

	// the annotations of namespace
	Annotations map[string]string `json:"annotations,omitempty"`

	// Whether the data is valid
	EnableFlag int64 `json:"enableFlag,omitempty"`

	// the annotations of namespace
	Hard map[string]interface{} `json:"hard,omitempty"`

	// the id of namespace
	ID int64 `json:"id,omitempty"`

	// the namespace name
	Namespace string `json:"namespace,omitempty"`

	// the namespace platformNamespace
	PlatformNamespace string `json:"platformNamespace,omitempty"`

	// the remarks of namespace
	Remarks string `json:"remarks,omitempty"`
}

// Validate validates this namespace response
func (m *NamespaceResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NamespaceResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NamespaceResponse) UnmarshalBinary(b []byte) error {
	var res NamespaceResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
