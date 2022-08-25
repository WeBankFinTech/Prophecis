// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Event event
// swagger:model Event
type Event struct {

	// endpoints
	Endpoints *EndpointList `json:"endpoints,omitempty"`

	// the type of event (i.e. update, metrics, logs, all...)
	Type string `json:"type,omitempty"`
}

// Validate validates this event
func (m *Event) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEndpoints(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Event) validateEndpoints(formats strfmt.Registry) error {

	if swag.IsZero(m.Endpoints) { // not required
		return nil
	}

	if m.Endpoints != nil {
		if err := m.Endpoints.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("endpoints")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Event) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Event) UnmarshalBinary(b []byte) error {
	var res Event
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
