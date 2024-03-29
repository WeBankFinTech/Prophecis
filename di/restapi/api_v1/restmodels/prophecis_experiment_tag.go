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

// ProphecisExperimentTag prophecis experiment tag
// swagger:model ProphecisExperimentTag
type ProphecisExperimentTag struct {

	// Flag of Delete.
	EnableFlag string `json:"enable_flag,omitempty"`

	// Experiment ID.
	ExpID int64 `json:"exp_id,omitempty"`

	// Experiment Tag Name.
	// Required: true
	ExpTag *string `json:"exp_tag"`

	// Timestamp recorded when this Experiment Tag  was created.
	// Format: date-time
	ExpTagCreateTime strfmt.DateTime `json:"exp_tag_create_time,omitempty"`

	// Creator recorded when this Experiment Tag  was created.
	ExpTagCreateUserID int64 `json:"exp_tag_create_user_id,omitempty"`

	// Experiment Tag ID.
	ID int64 `json:"id,omitempty"`
}

// Validate validates this prophecis experiment tag
func (m *ProphecisExperimentTag) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExpTag(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExpTagCreateTime(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProphecisExperimentTag) validateExpTag(formats strfmt.Registry) error {

	if err := validate.Required("exp_tag", "body", m.ExpTag); err != nil {
		return err
	}

	return nil
}

func (m *ProphecisExperimentTag) validateExpTagCreateTime(formats strfmt.Registry) error {

	if swag.IsZero(m.ExpTagCreateTime) { // not required
		return nil
	}

	if err := validate.FormatOf("exp_tag_create_time", "body", "date-time", m.ExpTagCreateTime.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentTag) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentTag) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentTag
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
