// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ProphecisExperimentTags prophecis experiment tags
// swagger:model ProphecisExperimentTags
type ProphecisExperimentTags struct {

	// experiment tags
	ExperimentTags []*ProphecisExperimentTag `json:"experiment_tags"`

	// page number
	PageNumber int64 `json:"page_number,omitempty"`

	// page size
	PageSize int64 `json:"page_size,omitempty"`
}

// Validate validates this prophecis experiment tags
func (m *ProphecisExperimentTags) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExperimentTags(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProphecisExperimentTags) validateExperimentTags(formats strfmt.Registry) error {

	if swag.IsZero(m.ExperimentTags) { // not required
		return nil
	}

	for i := 0; i < len(m.ExperimentTags); i++ {
		if swag.IsZero(m.ExperimentTags[i]) { // not required
			continue
		}

		if m.ExperimentTags[i] != nil {
			if err := m.ExperimentTags[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("experiment_tags" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProphecisExperimentTags) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProphecisExperimentTags) UnmarshalBinary(b []byte) error {
	var res ProphecisExperimentTags
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
