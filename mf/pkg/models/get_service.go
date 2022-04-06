// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// GetService get service
// swagger:model GetService
type GetService struct {

	// Serivce CPU resource.
	CPU float64 `json:"cpu,omitempty"`

	// creation timestamp of service.
	CreationTimestamp string `json:"creation_timestamp,omitempty"`

	// Service's Endpoint Type
	EndpointType string `json:"endpoint_type,omitempty"`

	// Serivce GPU resource.
	Gpu string `json:"gpu,omitempty"`

	// the group id of this service.
	GroupID int64 `json:"group_id,omitempty"`

	// Id for the Service.
	ID int64 `json:"id,omitempty"`

	// The service's image
	Image string `json:"image,omitempty"`

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

	// models of servier.
	Parameters ModelParametersList `json:"parameters,omitempty"`

	// The location path of the model material.
	Remark string `json:"remark,omitempty"`

	// Service Name.
	ServiceName string `json:"service_name,omitempty"`

	// Service's Status
	Status string `json:"status,omitempty"`

	// Service Type, Include Single, ABTest, Graph.
	Type string `json:"type,omitempty"`

	// the owner id of this service.
	UserID int64 `json:"user_id,omitempty"`
}

// Validate validates this get service
func (m *GetService) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateModelversion(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParameters(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GetService) validateModelversion(formats strfmt.Registry) error {

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

func (m *GetService) validateParameters(formats strfmt.Registry) error {

	if swag.IsZero(m.Parameters) { // not required
		return nil
	}

	if err := m.Parameters.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("parameters")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GetService) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetService) UnmarshalBinary(b []byte) error {
	var res GetService
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
