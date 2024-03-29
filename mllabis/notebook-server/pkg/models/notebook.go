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

// Notebook notebook
// swagger:model Notebook
type Notebook struct {

	// CPU used
	CPU string `json:"cpu,omitempty"`

	// Notebook data volumns mounted
	DataVolume []*MountInfo `json:"dataVolume"`

	// User driver memory in BDAP Yarn Cluster
	DriverMemory string `json:"driverMemory,omitempty"`

	// User excutor core setting in BDAP Yarn Cluster
	ExecutorCores string `json:"executorCores,omitempty"`

	// User excutor memory setting in BDAP Yarn Cluster
	ExecutorMemory string `json:"executorMemory,omitempty"`

	// User excutors setting in BDAP Yarn Cluster
	Executors string `json:"executors,omitempty"`

	// GPU used
	Gpu string `json:"gpu,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// Docker Image used
	Image string `json:"image,omitempty"`

	// Memory used
	Memory string `json:"memory,omitempty"`

	// Notebook name
	Name string `json:"name,omitempty"`

	// Namesapce where notebook has been created
	Namespace string `json:"namespace,omitempty"`

	// Notebook pods
	Pods string `json:"pods,omitempty"`

	// Proxy User of Notebook.
	ProxyUser string `json:"proxyUser,omitempty"`

	// BDAP Yarn Queue Setting
	Queue string `json:"queue,omitempty"`

	// Notebook services
	Service string `json:"service,omitempty"`

	// The number of spark session
	SparkSessionNum int64 `json:"sparkSessionNum,omitempty"`

	// Image short name
	SrtImage string `json:"srtImage,omitempty"`

	// Pod status
	Status string `json:"status,omitempty"`

	// Notebook create time
	Uptime string `json:"uptime,omitempty"`

	// Notebook Owner
	User string `json:"user,omitempty"`

	// Notebook volumns mounted
	Volumns string `json:"volumns,omitempty"`

	// Notebook workspace volumns mounted
	WorkspaceVolume *MountInfo `json:"workspaceVolume,omitempty"`
}

// Validate validates this notebook
func (m *Notebook) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDataVolume(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWorkspaceVolume(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Notebook) validateDataVolume(formats strfmt.Registry) error {

	if swag.IsZero(m.DataVolume) { // not required
		return nil
	}

	for i := 0; i < len(m.DataVolume); i++ {
		if swag.IsZero(m.DataVolume[i]) { // not required
			continue
		}

		if m.DataVolume[i] != nil {
			if err := m.DataVolume[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("dataVolume" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Notebook) validateWorkspaceVolume(formats strfmt.Registry) error {

	if swag.IsZero(m.WorkspaceVolume) { // not required
		return nil
	}

	if m.WorkspaceVolume != nil {
		if err := m.WorkspaceVolume.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("workspaceVolume")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Notebook) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Notebook) UnmarshalBinary(b []byte) error {
	var res Notebook
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
