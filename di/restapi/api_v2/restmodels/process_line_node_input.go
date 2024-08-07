// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ProcessLineNodeInput process line node input
// swagger:model ProcessLineNodeInput
type ProcessLineNodeInput struct {

	// 数据输入
	Data *DataConfig `json:"data,omitempty"`

	// 模型输入
	Processline *ProcessLineInput `json:"processline,omitempty"`
}

// Validate validates this process line node input
func (m *ProcessLineNodeInput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProcessline(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProcessLineNodeInput) validateData(formats strfmt.Registry) error {

	if swag.IsZero(m.Data) { // not required
		return nil
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

func (m *ProcessLineNodeInput) validateProcessline(formats strfmt.Registry) error {

	if swag.IsZero(m.Processline) { // not required
		return nil
	}

	if m.Processline != nil {
		if err := m.Processline.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("processline")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProcessLineNodeInput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProcessLineNodeInput) UnmarshalBinary(b []byte) error {
	var res ProcessLineNodeInput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
