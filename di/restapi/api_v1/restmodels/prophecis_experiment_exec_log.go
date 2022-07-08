// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ProphecisExperimentExecLog prophecis experiment exec log
// swagger:model ProphecisExperimentExecLog
type ProphecisExperimentExecLog struct {

	// Experiment ID.
	ExecID int64 `json:"exec_id,omitempty"`

	// Experiment Tag Name.
	Log string `json:"log,omitempty"`
}

// Validate validates this prophecis experiment exec log
func (m *ProphecisExperimentExecLog) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentExecLog) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentExecLog) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentExecLog
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}