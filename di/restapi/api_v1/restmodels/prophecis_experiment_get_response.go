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

// ProphecisExperimentGetResponse prophecis experiment get response
// swagger:model ProphecisExperimentGetResponse
type ProphecisExperimentGetResponse struct {

	// flow json
	// Required: true
	FlowJSON *string `json:"flow_json"`

	// prophecis experiment
	// Required: true
	ProphecisExperiment *ProphecisExperiment `json:"prophecis_experiment"`
}

// Validate validates this prophecis experiment get response
func (m *ProphecisExperimentGetResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFlowJSON(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProphecisExperiment(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProphecisExperimentGetResponse) validateFlowJSON(formats strfmt.Registry) error {

	if err := validate.Required("flow_json", "body", m.FlowJSON); err != nil {
		return err
	}

	return nil
}

func (m *ProphecisExperimentGetResponse) validateProphecisExperiment(formats strfmt.Registry) error {

	if err := validate.Required("prophecis_experiment", "body", m.ProphecisExperiment); err != nil {
		return err
	}

	if m.ProphecisExperiment != nil {
		if err := m.ProphecisExperiment.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("prophecis_experiment")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentGetResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentGetResponse) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentGetResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
