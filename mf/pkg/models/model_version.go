// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ModelVersion model version
// swagger:model ModelVersion
type ModelVersion struct {

	// Timestamp recorded when this model_version was created.
	CreationTimestamp string `json:"creation_timestamp,omitempty"`

	// enable flag
	EnableFlag int8 `json:"enable_flag,omitempty"`

	// file name
	FileName string `json:"file_name,omitempty"`

	// Running status of this model_version.
	// Required: true
	Filepath *string `json:"filepath"`

	// Id for the ModelVersion.
	// Read Only: true
	ID int64 `json:"id,omitempty"`

	// Timestamp recorded when metadata for this model_version was last updated.
	LatestFlag int64 `json:"latest_flag,omitempty"`

	// Description of this model_version.
	ModelID int64 `json:"model_id,omitempty"`

	// params
	Params string `json:"params,omitempty"`

	// push id
	PushID int64 `json:"push_id,omitempty"`

	// Timestamp pushed.
	PushTimestamp string `json:"push_timestamp,omitempty"`

	// User that created this model_version.
	// Required: true
	Source *string `json:"source"`

	// The flag that whether generated by model operator.
	TrainingFlag int8 `json:"training_flag,omitempty"`

	// The id of training
	TrainingID string `json:"training_id,omitempty"`

	// Model’s version number.
	// Required: true
	Version *string `json:"version"`
}

// Validate validates this model version
func (m *ModelVersion) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFilepath(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSource(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelVersion) validateFilepath(formats strfmt.Registry) error {

	if err := validate.Required("filepath", "body", m.Filepath); err != nil {
		return err
	}

	return nil
}

func (m *ModelVersion) validateSource(formats strfmt.Registry) error {

	if err := validate.Required("source", "body", m.Source); err != nil {
		return err
	}

	return nil
}

func (m *ModelVersion) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelVersion) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelVersion) UnmarshalBinary(b []byte) error {
	var res ModelVersion
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}