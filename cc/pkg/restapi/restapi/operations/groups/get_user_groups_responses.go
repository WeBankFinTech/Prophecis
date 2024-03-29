// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// GetUserGroupsOKCode is the HTTP code returned for type GetUserGroupsOK
const GetUserGroupsOKCode int = 200

/*GetUserGroupsOK Detailed userGroup and userGroup information.

swagger:response getUserGroupsOK
*/
type GetUserGroupsOK struct {

	/*
	  In: Body
	*/
	Payload models.UserGroups `json:"body,omitempty"`
}

// NewGetUserGroupsOK creates GetUserGroupsOK with default headers values
func NewGetUserGroupsOK() *GetUserGroupsOK {

	return &GetUserGroupsOK{}
}

// WithPayload adds the payload to the get user groups o k response
func (o *GetUserGroupsOK) WithPayload(payload models.UserGroups) *GetUserGroupsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user groups o k response
func (o *GetUserGroupsOK) SetPayload(payload models.UserGroups) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserGroupsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.UserGroups{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetUserGroupsUnauthorizedCode is the HTTP code returned for type GetUserGroupsUnauthorized
const GetUserGroupsUnauthorizedCode int = 401

/*GetUserGroupsUnauthorized Unauthorized

swagger:response getUserGroupsUnauthorized
*/
type GetUserGroupsUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUserGroupsUnauthorized creates GetUserGroupsUnauthorized with default headers values
func NewGetUserGroupsUnauthorized() *GetUserGroupsUnauthorized {

	return &GetUserGroupsUnauthorized{}
}

// WithPayload adds the payload to the get user groups unauthorized response
func (o *GetUserGroupsUnauthorized) WithPayload(payload *models.Error) *GetUserGroupsUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user groups unauthorized response
func (o *GetUserGroupsUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserGroupsUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUserGroupsNotFoundCode is the HTTP code returned for type GetUserGroupsNotFound
const GetUserGroupsNotFoundCode int = 404

/*GetUserGroupsNotFound Model with the given ID not found.

swagger:response getUserGroupsNotFound
*/
type GetUserGroupsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUserGroupsNotFound creates GetUserGroupsNotFound with default headers values
func NewGetUserGroupsNotFound() *GetUserGroupsNotFound {

	return &GetUserGroupsNotFound{}
}

// WithPayload adds the payload to the get user groups not found response
func (o *GetUserGroupsNotFound) WithPayload(payload *models.Error) *GetUserGroupsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user groups not found response
func (o *GetUserGroupsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserGroupsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
