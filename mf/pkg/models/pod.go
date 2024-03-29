// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Pod pod
// swagger:model Pod
type Pod struct {

	// containers in this pod.
	Containers []*Container `json:"containers"`

	// The message why pod's status that it is.
	Message string `json:"message,omitempty"`

	// pod's name
	Name string `json:"name,omitempty"`

	// pod's namespaces
	Namespace string `json:"namespace,omitempty"`

	// The reason why pod's status that it is.
	Reason string `json:"reason,omitempty"`

	// pod's status
	Status string `json:"status,omitempty"`
}

// Validate validates this pod
func (m *Pod) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContainers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Pod) validateContainers(formats strfmt.Registry) error {

	if swag.IsZero(m.Containers) { // not required
		return nil
	}

	for i := 0; i < len(m.Containers); i++ {
		if swag.IsZero(m.Containers[i]) { // not required
			continue
		}

		if m.Containers[i] != nil {
			if err := m.Containers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("containers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Pod) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Pod) UnmarshalBinary(b []byte) error {
	var res Pod
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
