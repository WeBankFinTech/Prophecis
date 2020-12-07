// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// UpdateGroupStorageOKCode is the HTTP code returned for type UpdateGroupStorageOK
const UpdateGroupStorageOKCode int = 200

/*UpdateGroupStorageOK Detailed group and group information.

swagger:response updateGroupStorageOK
*/
type UpdateGroupStorageOK struct {

	/*
	  In: Body
	*/
	Payload *models.GroupStorage `json:"body,omitempty"`
}

// NewUpdateGroupStorageOK creates UpdateGroupStorageOK with default headers values
func NewUpdateGroupStorageOK() *UpdateGroupStorageOK {

	return &UpdateGroupStorageOK{}
}

// WithPayload adds the payload to the update group storage o k response
func (o *UpdateGroupStorageOK) WithPayload(payload *models.GroupStorage) *UpdateGroupStorageOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update group storage o k response
func (o *UpdateGroupStorageOK) SetPayload(payload *models.GroupStorage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateGroupStorageOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateGroupStorageUnauthorizedCode is the HTTP code returned for type UpdateGroupStorageUnauthorized
const UpdateGroupStorageUnauthorizedCode int = 401

/*UpdateGroupStorageUnauthorized Unauthorized

swagger:response updateGroupStorageUnauthorized
*/
type UpdateGroupStorageUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUpdateGroupStorageUnauthorized creates UpdateGroupStorageUnauthorized with default headers values
func NewUpdateGroupStorageUnauthorized() *UpdateGroupStorageUnauthorized {

	return &UpdateGroupStorageUnauthorized{}
}

// WithPayload adds the payload to the update group storage unauthorized response
func (o *UpdateGroupStorageUnauthorized) WithPayload(payload *models.Error) *UpdateGroupStorageUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update group storage unauthorized response
func (o *UpdateGroupStorageUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateGroupStorageUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateGroupStorageNotFoundCode is the HTTP code returned for type UpdateGroupStorageNotFound
const UpdateGroupStorageNotFoundCode int = 404

/*UpdateGroupStorageNotFound url to get allGroups not found.

swagger:response updateGroupStorageNotFound
*/
type UpdateGroupStorageNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewUpdateGroupStorageNotFound creates UpdateGroupStorageNotFound with default headers values
func NewUpdateGroupStorageNotFound() *UpdateGroupStorageNotFound {

	return &UpdateGroupStorageNotFound{}
}

// WithPayload adds the payload to the update group storage not found response
func (o *UpdateGroupStorageNotFound) WithPayload(payload *models.Error) *UpdateGroupStorageNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update group storage not found response
func (o *UpdateGroupStorageNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateGroupStorageNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
