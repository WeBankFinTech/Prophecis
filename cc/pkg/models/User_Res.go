/*
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package models

import (
	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

type UserRes struct {
	// Whether the data is valid
	EnableFlag int64 `json:"enableFlag,omitempty"`

	// the gid of user
	Gid int64 `json:"gid,omitempty"`

	// true or false
	GUIDCheck bool `json:"guidCheck,omitempty"`

	// the id of user
	ID int64 `json:"id,omitempty"`

	// the username
	Name string `json:"name,omitempty"`

	// the uid remarks user
	Remarks string `json:"remarks,omitempty"`

	// the type of user, user or system
	Type string `json:"type,omitempty"`

	// the uid of user
	UID int64 `json:"uid,omitempty"`

	// the client_token of user
	Token string `json:"token,omitempty"`
}

// Validate validates this user
func (m *UserRes) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UserRes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserRes) UnmarshalBinary(b []byte) error {
	var res UserRes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
