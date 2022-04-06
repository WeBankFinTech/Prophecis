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

// PostModelRequest post model request
// swagger:model PostModelRequest
type PostModelRequest struct {

	// File's child path
	ChildPath string `json:"child_path,omitempty"`

	// File's name
	FileName string `json:"file_name,omitempty"`

	// directory of file storage
	Filepath string `json:"filepath,omitempty"`

	// group id
	// Required: true
	GroupID *int64 `json:"group_id"`

	// model name
	// Required: true
	ModelName *string `json:"model_name"`

	// Type for the model
	// Required: true
	ModelType *string `json:"model_type"`

	// model version , eg v1, v2 ...
	ModelVersion string `json:"model_version,omitempty"`

	// File's root path
	RootPath string `json:"root_path,omitempty"`

	// S3Path
	S3Path string `json:"s3_path,omitempty"`

	// The flag that whether generated by model operator.
	TrainingFlag int8 `json:"training_flag,omitempty"`

	// The id of training
	TrainingID string `json:"training_id,omitempty"`
}

// Validate validates this post model request
func (m *PostModelRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGroupID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModelName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModelType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostModelRequest) validateGroupID(formats strfmt.Registry) error {

	if err := validate.Required("group_id", "body", m.GroupID); err != nil {
		return err
	}

	return nil
}

func (m *PostModelRequest) validateModelName(formats strfmt.Registry) error {

	if err := validate.Required("model_name", "body", m.ModelName); err != nil {
		return err
	}

	return nil
}

func (m *PostModelRequest) validateModelType(formats strfmt.Registry) error {

	if err := validate.Required("model_type", "body", m.ModelType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostModelRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostModelRequest) UnmarshalBinary(b []byte) error {
	var res PostModelRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}