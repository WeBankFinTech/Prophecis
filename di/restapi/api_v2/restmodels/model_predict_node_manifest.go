// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// ModelPredictNodeManifest model predict node manifest
// swagger:model ModelPredictNodeManifest
type ModelPredictNodeManifest struct {

	// compute resources
	ComputeResources *ModelPredictNodeManifestComputeResources `json:"compute_resources,omitempty"`

	// 节点输入信息
	Input *ModelPredictNodeInput `json:"input,omitempty"`

	// 节点的元数据信息
	MetaData *MetaData `json:"meta_data,omitempty"`

	// 节点的输出信息
	Output *Output `json:"output,omitempty"`

	// 节点的运行环境信息
	RunEnvironment *RunEnvironment `json:"run_environment,omitempty"`
}

// Validate validates this model predict node manifest
func (m *ModelPredictNodeManifest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComputeResources(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInput(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetaData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOutput(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRunEnvironment(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelPredictNodeManifest) validateComputeResources(formats strfmt.Registry) error {

	if swag.IsZero(m.ComputeResources) { // not required
		return nil
	}

	if m.ComputeResources != nil {
		if err := m.ComputeResources.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("compute_resources")
			}
			return err
		}
	}

	return nil
}

func (m *ModelPredictNodeManifest) validateInput(formats strfmt.Registry) error {

	if swag.IsZero(m.Input) { // not required
		return nil
	}

	if m.Input != nil {
		if err := m.Input.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("input")
			}
			return err
		}
	}

	return nil
}

func (m *ModelPredictNodeManifest) validateMetaData(formats strfmt.Registry) error {

	if swag.IsZero(m.MetaData) { // not required
		return nil
	}

	if m.MetaData != nil {
		if err := m.MetaData.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("meta_data")
			}
			return err
		}
	}

	return nil
}

func (m *ModelPredictNodeManifest) validateOutput(formats strfmt.Registry) error {

	if swag.IsZero(m.Output) { // not required
		return nil
	}

	if m.Output != nil {
		if err := m.Output.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("output")
			}
			return err
		}
	}

	return nil
}

func (m *ModelPredictNodeManifest) validateRunEnvironment(formats strfmt.Registry) error {

	if swag.IsZero(m.RunEnvironment) { // not required
		return nil
	}

	if m.RunEnvironment != nil {
		if err := m.RunEnvironment.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("run_environment")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelPredictNodeManifest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelPredictNodeManifest) UnmarshalBinary(b []byte) error {
	var res ModelPredictNodeManifest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ModelPredictNodeManifestComputeResources 节点的计算资源信息
// swagger:model ModelPredictNodeManifestComputeResources
type ModelPredictNodeManifestComputeResources struct {

	// worker
	Worker *ComputeResource `json:"worker,omitempty"`
}

// Validate validates this model predict node manifest compute resources
func (m *ModelPredictNodeManifestComputeResources) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateWorker(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModelPredictNodeManifestComputeResources) validateWorker(formats strfmt.Registry) error {

	if swag.IsZero(m.Worker) { // not required
		return nil
	}

	if m.Worker != nil {
		if err := m.Worker.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("compute_resources" + "." + "worker")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ModelPredictNodeManifestComputeResources) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModelPredictNodeManifestComputeResources) UnmarshalBinary(b []byte) error {
	var res ModelPredictNodeManifestComputeResources
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
