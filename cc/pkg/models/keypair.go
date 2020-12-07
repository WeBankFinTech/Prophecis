// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Keypair keypair
// swagger:model Keypair
type Keypair struct {

	// the apiKey
	APIKey string `json:"apiKey,omitempty"`

	// the enableFlag
	EnableFlag int64 `json:"enableFlag,omitempty"`

	// the id
	ID int64 `json:"id,omitempty"`

	// the name
	Name string `json:"name,omitempty"`

	// the remarks
	Remarks string `json:"remarks,omitempty"`

	// the secretKey
	SecretKey string `json:"secretKey,omitempty"`
}

// Validate validates this keypair
func (m *Keypair) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Keypair) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Keypair) UnmarshalBinary(b []byte) error {
	var res Keypair
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
