// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// GroupStorage group storage
// swagger:model GroupStorage
type GroupStorage struct {

	// Whether the data is valid
	EnableFlag int64 `json:"enableFlag,omitempty"`

	// the groupId
	GroupID int64 `json:"groupId,omitempty"`

	// the type of group
	GroupType string `json:"groupType,omitempty"`

	// the id of GroupStorage
	ID int64 `json:"id,omitempty"`

	// the path of storage
	Path string `json:"path,omitempty"`

	// the permissions of storage
	Permissions string `json:"permissions,omitempty"`

	// the remarks UserGroup
	Remarks string `json:"remarks,omitempty"`

	// the storageId
	StorageID int64 `json:"storageId,omitempty"`

	// the type of storage
	Type string `json:"type,omitempty"`
}

// Validate validates this group storage
func (m *GroupStorage) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GroupStorage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GroupStorage) UnmarshalBinary(b []byte) error {
	var res GroupStorage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
