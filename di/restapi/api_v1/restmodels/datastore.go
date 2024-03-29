// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Datastore datastore
// swagger:model Datastore
type Datastore struct {

	// fields
	Fields map[string]string `json:"Fields,omitempty"`

	// connection
	Connection map[string]string `json:"connection,omitempty"`

	// the id of the data store as defined in the manifest.
	DataStoreID string `json:"data_store_id,omitempty"`

	// the type of the data store as defined in the manifest.
	Type string `json:"type,omitempty"`
}

// Validate validates this datastore
func (m *Datastore) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Datastore) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Datastore) UnmarshalBinary(b []byte) error {
	var res Datastore
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
