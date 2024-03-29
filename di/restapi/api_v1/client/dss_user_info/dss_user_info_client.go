// Code generated by go-swagger; DO NOT EDIT.

package dss_user_info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new dss user info API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for dss user info API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
GetDssUserInfo gets d s s user info

Get DSS UserInfo
*/
func (a *Client) GetDssUserInfo(params *GetDssUserInfoParams, authInfo runtime.ClientAuthInfoWriter) (*GetDssUserInfoOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetDssUserInfoParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getDssUserInfo",
		Method:             "GET",
		PathPattern:        "/di/v1/dssUserInfo",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetDssUserInfoReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetDssUserInfoOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for getDssUserInfo: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
