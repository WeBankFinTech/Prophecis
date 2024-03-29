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

// Event event
// swagger:model Event
type Event struct {

	// create timestamp
	// Required: true
	// Format: date-time
	CreationTimestamp *strfmt.DateTime `json:"creation_timestamp"`

	// DCN of RMB
	// Required: true
	Dcn *string `json:"dcn"`

	// 1 normal; 0 freeze
	// Required: true
	EnableFlag *int8 `json:"enable_flag"`

	// report id or model id
	// Required: true
	FileID *int64 `json:"file_id"`

	// file name
	// Required: true
	FileName *string `json:"file_name"`

	// file type eg:MODEL, DATA
	// Required: true
	FileType *string `json:"file_type"`

	// file id from FPS
	// Required: true
	FpsFileID *string `json:"fps_file_id"`

	// hash from FPS
	// Required: true
	HashValue *string `json:"hash_value"`

	// report id
	// Required: true
	ID *int64 `json:"id"`

	// IDC of RMB
	// Required: true
	Idc *string `json:"idc"`

	// params
	// Required: true
	Params *string `json:"params"`

	// rmb resp file name
	RmbRespFileName string `json:"rmb_resp_file_name,omitempty"`

	// rmb s3path
	RmbS3path string `json:"rmb_s3path,omitempty"`

	// hash from FPS
	// Required: true
	Status *string `json:"status"`

	// report pushed timestamp
	// Required: true
	// Format: date-time
	UpdateTimestamp *strfmt.DateTime `json:"update_timestamp"`

	// version
	Version string `json:"version,omitempty"`
}

// Validate validates this event
func (m *Event) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreationTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDcn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnableFlag(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFileID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFileName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFileType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFpsFileID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHashValue(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIdc(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParams(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Event) validateCreationTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("creation_timestamp", "body", m.CreationTimestamp); err != nil {
		return err
	}

	if err := validate.FormatOf("creation_timestamp", "body", "date-time", m.CreationTimestamp.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateDcn(formats strfmt.Registry) error {

	if err := validate.Required("dcn", "body", m.Dcn); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateEnableFlag(formats strfmt.Registry) error {

	if err := validate.Required("enable_flag", "body", m.EnableFlag); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateFileID(formats strfmt.Registry) error {

	if err := validate.Required("file_id", "body", m.FileID); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateFileName(formats strfmt.Registry) error {

	if err := validate.Required("file_name", "body", m.FileName); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateFileType(formats strfmt.Registry) error {

	if err := validate.Required("file_type", "body", m.FileType); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateFpsFileID(formats strfmt.Registry) error {

	if err := validate.Required("fps_file_id", "body", m.FpsFileID); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateHashValue(formats strfmt.Registry) error {

	if err := validate.Required("hash_value", "body", m.HashValue); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateIdc(formats strfmt.Registry) error {

	if err := validate.Required("idc", "body", m.Idc); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateParams(formats strfmt.Registry) error {

	if err := validate.Required("params", "body", m.Params); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

func (m *Event) validateUpdateTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("update_timestamp", "body", m.UpdateTimestamp); err != nil {
		return err
	}

	if err := validate.FormatOf("update_timestamp", "body", "date-time", m.UpdateTimestamp.String(), formats); err != nil {
		return err
	}

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
