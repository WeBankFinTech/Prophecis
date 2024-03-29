// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// PostAlertRequest post alert request
// swagger:model PostAlertRequest
type PostAlertRequest struct {

	// alert list
	AlertList []*PostAlertSubRequest `json:"alertList"`
}

// Validate validates this post alert request
func (m *PostAlertRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlertList(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostAlertRequest) validateAlertList(formats strfmt.Registry) error {

	if swag.IsZero(m.AlertList) { // not required
		return nil
	}

	for i := 0; i < len(m.AlertList); i++ {
		if swag.IsZero(m.AlertList[i]) { // not required
			continue
		}

		if m.AlertList[i] != nil {
			if err := m.AlertList[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("alertList" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostAlertRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostAlertRequest) UnmarshalBinary(b []byte) error {
	var res PostAlertRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
