// Code generated by go-swagger; DO NOT EDIT.

package namespaces

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// ListNamespaceByRoleNameAndUserNameOKCode is the HTTP code returned for type ListNamespaceByRoleNameAndUserNameOK
const ListNamespaceByRoleNameAndUserNameOKCode int = 200

/*ListNamespaceByRoleNameAndUserNameOK get namespace by role name and user name.

swagger:response listNamespaceByRoleNameAndUserNameOK
*/
type ListNamespaceByRoleNameAndUserNameOK struct {

	/*
	  In: Body
	*/
	Payload models.ListNamespaceByRolenameAndUserNameResp `json:"body,omitempty"`
}

// NewListNamespaceByRoleNameAndUserNameOK creates ListNamespaceByRoleNameAndUserNameOK with default headers values
func NewListNamespaceByRoleNameAndUserNameOK() *ListNamespaceByRoleNameAndUserNameOK {

	return &ListNamespaceByRoleNameAndUserNameOK{}
}

// WithPayload adds the payload to the list namespace by role name and user name o k response
func (o *ListNamespaceByRoleNameAndUserNameOK) WithPayload(payload models.ListNamespaceByRolenameAndUserNameResp) *ListNamespaceByRoleNameAndUserNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list namespace by role name and user name o k response
func (o *ListNamespaceByRoleNameAndUserNameOK) SetPayload(payload models.ListNamespaceByRolenameAndUserNameResp) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListNamespaceByRoleNameAndUserNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.ListNamespaceByRolenameAndUserNameResp{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ListNamespaceByRoleNameAndUserNameUnauthorizedCode is the HTTP code returned for type ListNamespaceByRoleNameAndUserNameUnauthorized
const ListNamespaceByRoleNameAndUserNameUnauthorizedCode int = 401

/*ListNamespaceByRoleNameAndUserNameUnauthorized Unauthorized

swagger:response listNamespaceByRoleNameAndUserNameUnauthorized
*/
type ListNamespaceByRoleNameAndUserNameUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListNamespaceByRoleNameAndUserNameUnauthorized creates ListNamespaceByRoleNameAndUserNameUnauthorized with default headers values
func NewListNamespaceByRoleNameAndUserNameUnauthorized() *ListNamespaceByRoleNameAndUserNameUnauthorized {

	return &ListNamespaceByRoleNameAndUserNameUnauthorized{}
}

// WithPayload adds the payload to the list namespace by role name and user name unauthorized response
func (o *ListNamespaceByRoleNameAndUserNameUnauthorized) WithPayload(payload *models.Error) *ListNamespaceByRoleNameAndUserNameUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list namespace by role name and user name unauthorized response
func (o *ListNamespaceByRoleNameAndUserNameUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListNamespaceByRoleNameAndUserNameUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListNamespaceByRoleNameAndUserNameNotFoundCode is the HTTP code returned for type ListNamespaceByRoleNameAndUserNameNotFound
const ListNamespaceByRoleNameAndUserNameNotFoundCode int = 404

/*ListNamespaceByRoleNameAndUserNameNotFound Model with the given ID not found.

swagger:response listNamespaceByRoleNameAndUserNameNotFound
*/
type ListNamespaceByRoleNameAndUserNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListNamespaceByRoleNameAndUserNameNotFound creates ListNamespaceByRoleNameAndUserNameNotFound with default headers values
func NewListNamespaceByRoleNameAndUserNameNotFound() *ListNamespaceByRoleNameAndUserNameNotFound {

	return &ListNamespaceByRoleNameAndUserNameNotFound{}
}

// WithPayload adds the payload to the list namespace by role name and user name not found response
func (o *ListNamespaceByRoleNameAndUserNameNotFound) WithPayload(payload *models.Error) *ListNamespaceByRoleNameAndUserNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list namespace by role name and user name not found response
func (o *ListNamespaceByRoleNameAndUserNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListNamespaceByRoleNameAndUserNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
