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

// ProphecisExperimentRunRequest prophecis experiment run request
// swagger:model ProphecisExperimentRunRequest
type ProphecisExperimentRunRequest struct {

	// Experiment Run type, Include CLI/UI/DSS/Schedulis.
	// Required: true
	ExpExecType *string `json:"exp_exec_type"`

	// Experiment Run type, Include CLI/UI/DSS/Schedulis.
	FlowJSON string `json:"flow_json,omitempty"`
}

// Validate validates this prophecis experiment run request
func (m *ProphecisExperimentRunRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExpExecType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProphecisExperimentRunRequest) validateExpExecType(formats strfmt.Registry) error {

	if err := validate.Required("exp_exec_type", "body", m.ExpExecType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentRunRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentRunRequest) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentRunRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
