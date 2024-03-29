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

// GetNotebookStatusResponse get notebook status response
// swagger:model GetNotebookStatusResponse
type GetNotebookStatusResponse struct {

	// The information of notebook containers status
	ContainersStatusInfo []*ContainerStatusInfo `json:"containers_status_info"`

	// Notebook create time
	CreateTime string `json:"create_time,omitempty"`

	// Docker Image used
	Image string `json:"image,omitempty"`

	// Notebook name
	Name string `json:"name,omitempty"`

	// Namesapce where notebook has been created
	Namespace string `json:"namespace,omitempty"`

	// Notebook status
	NotebookStatus string `json:"notebook_status,omitempty"`

	// The information that trigger notebook status, consist of Container State
	NotebookStatusInfo string `json:"notebook_status_info,omitempty"`

	// Rsource that Notebook limted
	ResourceLimit *Resource `json:"resource_limit,omitempty"`

	// Notebook request resource
	ResourceReq *Resource `json:"resource_req,omitempty"`
}

// Validate validates this get notebook status response
func (m *GetNotebookStatusResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContainersStatusInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceLimit(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceReq(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GetNotebookStatusResponse) validateContainersStatusInfo(formats strfmt.Registry) error {

	if swag.IsZero(m.ContainersStatusInfo) { // not required
		return nil
	}

	for i := 0; i < len(m.ContainersStatusInfo); i++ {
		if swag.IsZero(m.ContainersStatusInfo[i]) { // not required
			continue
		}

		if m.ContainersStatusInfo[i] != nil {
			if err := m.ContainersStatusInfo[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("containers_status_info" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *GetNotebookStatusResponse) validateResourceLimit(formats strfmt.Registry) error {

	if swag.IsZero(m.ResourceLimit) { // not required
		return nil
	}

	if m.ResourceLimit != nil {
		if err := m.ResourceLimit.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("resource_limit")
			}
			return err
		}
	}

	return nil
}

func (m *GetNotebookStatusResponse) validateResourceReq(formats strfmt.Registry) error {

	if swag.IsZero(m.ResourceReq) { // not required
		return nil
	}

	if m.ResourceReq != nil {
		if err := m.ResourceReq.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("resource_req")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GetNotebookStatusResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetNotebookStatusResponse) UnmarshalBinary(b []byte) error {
	var res GetNotebookStatusResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
