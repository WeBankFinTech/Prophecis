// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// Event event
// swagger:model Event
type Event struct {

	// The number of times this event has occurred
	Count int32 `json:"count,omitempty"`

	// The time at which the event was first recorded
	FirstTime string `json:"first_time,omitempty"`

	// Type of this event (Normal, Warning), new types could be added in the future
	Kind string `json:"kind,omitempty"`

	// The time at which the most recent occurrence of this event was recorded
	LastTime string `json:"last_time,omitempty"`

	// Description of the status of this event
	Message string `json:"message,omitempty"`

	// event name
	Name string `json:"name,omitempty"`

	// event namespace
	Namespace string `json:"namespace,omitempty"`
}

// Validate validates this event
func (m *Event) Validate(formats strfmt.Registry) error {
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
