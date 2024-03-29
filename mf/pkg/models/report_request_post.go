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

// ReportRequestPost report request post
// swagger:model ReportRequestPost
type ReportRequestPost struct {

	// File's child path
	ChildPath string `json:"child_path,omitempty"`

	// File's name
	FileName string `json:"file_name,omitempty"`

	// group id
	// Required: true
	GroupID *int64 `json:"group_id"`

	// model name
	// Required: true
	ModelName *string `json:"model_name"`

	// model version , eg v1, v2 ...
	// Required: true
	ModelVersion *string `json:"model_version"`

	// report name
	// Required: true
	ReportName *string `json:"report_name"`

	// if report_id + report_version not exist in report_version table, create report_version
	ReportVersion string `json:"report_version,omitempty"`

	// File's root path
	RootPath string `json:"root_path,omitempty"`

	// s3 path
	S3Path string `json:"s3_path,omitempty"`
}

// Validate validates this report request post
func (m *ReportRequestPost) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGroupID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModelName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModelVersion(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReportName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ReportRequestPost) validateGroupID(formats strfmt.Registry) error {

	if err := validate.Required("group_id", "body", m.GroupID); err != nil {
		return err
	}

	return nil
}

func (m *ReportRequestPost) validateModelName(formats strfmt.Registry) error {

	if err := validate.Required("model_name", "body", m.ModelName); err != nil {
		return err
	}

	return nil
}

func (m *ReportRequestPost) validateModelVersion(formats strfmt.Registry) error {

	if err := validate.Required("model_version", "body", m.ModelVersion); err != nil {
		return err
	}

	return nil
}

func (m *ReportRequestPost) validateReportName(formats strfmt.Registry) error {

	if err := validate.Required("report_name", "body", m.ReportName); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ReportRequestPost) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ReportRequestPost) UnmarshalBinary(b []byte) error {
	var res ReportRequestPost
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
