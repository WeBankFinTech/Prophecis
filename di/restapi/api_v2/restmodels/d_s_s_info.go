// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// DSSInfo d s s info
// swagger:model DSSInfo
type DSSInfo struct {

	// flow id
	FlowID int64 `json:"flow_id,omitempty"`

	// flow name
	FlowName string `json:"flow_name,omitempty"`

	// flow version
	FlowVersion string `json:"flow_version,omitempty"`

	// project id
	ProjectID int64 `json:"project_id,omitempty"`

	// project name
	ProjectName string `json:"project_name,omitempty"`

	// workspace id
	WorkspaceID int64 `json:"workspace_id,omitempty"`

	// workspace name,choose next value, "bdapWorkspace" "bdapWorkspace0x"
	WorkspaceName string `json:"workspace_name,omitempty"`
}

// Validate validates this d s s info
func (m *DSSInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DSSInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DSSInfo) UnmarshalBinary(b []byte) error {
	var res DSSInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
