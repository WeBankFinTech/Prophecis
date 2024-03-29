// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProphecisExperimentIDResponse prophecis experiment ID response
// swagger:model ProphecisExperimentIDResponse
type ProphecisExperimentIDResponse struct {

	// ProphecisExperiment Name
	ExpName string `json:"exp_name,omitempty"`

	// ProphecisExperiment ID
	// Required: true
	ID *int64 `json:"id"`
}

// Validate validates this prophecis experiment ID response
func (m *ProphecisExperimentIDResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProphecisExperimentIDResponse) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentIDResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentIDResponse) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentIDResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
