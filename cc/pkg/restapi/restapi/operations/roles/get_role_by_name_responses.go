// Code generated by go-swagger; DO NOT EDIT.

package roles

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// GetRoleByNameOKCode is the HTTP code returned for type GetRoleByNameOK
const GetRoleByNameOKCode int = 200

/*GetRoleByNameOK Detailed Role and Role information.

swagger:response getRoleByNameOK
*/
type GetRoleByNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.Role `json:"body,omitempty"`
}

// NewGetRoleByNameOK creates GetRoleByNameOK with default headers values
func NewGetRoleByNameOK() *GetRoleByNameOK {

	return &GetRoleByNameOK{}
}

// WithPayload adds the payload to the get role by name o k response
func (o *GetRoleByNameOK) WithPayload(payload *models.Role) *GetRoleByNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get role by name o k response
func (o *GetRoleByNameOK) SetPayload(payload *models.Role) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRoleByNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRoleByNameUnauthorizedCode is the HTTP code returned for type GetRoleByNameUnauthorized
const GetRoleByNameUnauthorizedCode int = 401

/*GetRoleByNameUnauthorized Unauthorized

swagger:response getRoleByNameUnauthorized
*/
type GetRoleByNameUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRoleByNameUnauthorized creates GetRoleByNameUnauthorized with default headers values
func NewGetRoleByNameUnauthorized() *GetRoleByNameUnauthorized {

	return &GetRoleByNameUnauthorized{}
}

// WithPayload adds the payload to the get role by name unauthorized response
func (o *GetRoleByNameUnauthorized) WithPayload(payload *models.Error) *GetRoleByNameUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get role by name unauthorized response
func (o *GetRoleByNameUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRoleByNameUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRoleByNameNotFoundCode is the HTTP code returned for type GetRoleByNameNotFound
const GetRoleByNameNotFoundCode int = 404

/*GetRoleByNameNotFound url to add namespace not found.

swagger:response getRoleByNameNotFound
*/
type GetRoleByNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRoleByNameNotFound creates GetRoleByNameNotFound with default headers values
func NewGetRoleByNameNotFound() *GetRoleByNameNotFound {

	return &GetRoleByNameNotFound{}
}

// WithPayload adds the payload to the get role by name not found response
func (o *GetRoleByNameNotFound) WithPayload(payload *models.Error) *GetRoleByNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get role by name not found response
func (o *GetRoleByNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRoleByNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
