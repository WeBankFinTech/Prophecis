// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Service service
// swagger:model Service
type Service struct {

	// Serivce CPU resource.
	CPU float64 `json:"cpu,omitempty"`

	// creation timestamp of service.
	CreationTimestamp string `json:"creation_timestamp,omitempty"`

	// Endpoint type
	EndpointType string `json:"endpoint_type,omitempty"`

	// Serivce GPU resource.
	Gpu string `json:"gpu,omitempty"`

	// the group id of this service.
	GroupID int64 `json:"group_id,omitempty"`

	// Id for the Service.
	ID int64 `json:"id,omitempty"`

	// last update timestamp of this service.
	LastUpdatedTimestamp string `json:"last_updated_timestamp,omitempty"`

	// The location of log path.
	LogPath string `json:"log_path,omitempty"`

	// Service memory resource.
	Memory float64 `json:"memory,omitempty"`

	// models of servier.
	Modelversion ServiceModelVersions `json:"modelversion,omitempty"`

	// Service NameSpace.
	Namespace string `json:"namespace,omitempty"`

	// The location path of the model material.
	Remark string `json:"remark,omitempty"`

	// Service Name.
	ServiceName string `json:"service_name,omitempty"`

	// Service Type, Include Single, ABTest, Graph.
	Type string `json:"type,omitempty"`

	// the owner id of this service.
	UserID int64 `json:"user_id,omitempty"`
}

// Validate validates this service
func (m *Service) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateModelversion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Service) validateModelversion(formats strfmt.Registry) error {

	if swag.IsZero(m.Modelversion) { // not required
		return nil
	}

	if err := m.Modelversion.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("modelversion")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Service) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Service) UnmarshalBinary(b []byte) error {
	var res Service
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
