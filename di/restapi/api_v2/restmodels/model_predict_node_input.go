// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ModelPredictNodeInput model predict node input
// swagger:model ModelPredictNodeInput
type ModelPredictNodeInput struct {

	// 数据输入
	Data *DataConfig `json:"data,omitempty"`

	// 模型输入
	Model *ModelInput `json:"model,omitempty"`
}

// Validate validates this model predict node input
func (m *ModelPredictNodeInput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModel(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelPredictNodeInput) validateData(formats strfmt.Registry) error {

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

func (m *ModelPredictNodeInput) validateModel(formats strfmt.Registry) error {

	if swag.IsZero(m.Model) { // not required
		return nil
	}

	if m.Model != nil {
		if err := m.Model.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("model")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelPredictNodeInput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelPredictNodeInput) UnmarshalBinary(b []byte) error {
	var res ModelPredictNodeInput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
