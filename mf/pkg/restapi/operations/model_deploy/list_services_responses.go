// Code generated by go-swagger; DO NOT EDIT.

package model_deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-mf/pkg/models"
)

// ListServicesOKCode is the HTTP code returned for type ListServicesOK
const ListServicesOKCode int = 200

/*ListServicesOK OK

swagger:response listServicesOK
*/
type ListServicesOK struct {

	/*
	  In: Body
	*/
	Payload models.Services `json:"body,omitempty"`
}

// NewListServicesOK creates ListServicesOK with default headers values
func NewListServicesOK() *ListServicesOK {

	return &ListServicesOK{}
}

// WithPayload adds the payload to the list services o k response
func (o *ListServicesOK) WithPayload(payload models.Services) *ListServicesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list services o k response
func (o *ListServicesOK) SetPayload(payload models.Services) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListServicesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Services{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ListServicesUnauthorizedCode is the HTTP code returned for type ListServicesUnauthorized
const ListServicesUnauthorizedCode int = 401

/*ListServicesUnauthorized Unauthorized

swagger:response listServicesUnauthorized
*/
type ListServicesUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListServicesUnauthorized creates ListServicesUnauthorized with default headers values
func NewListServicesUnauthorized() *ListServicesUnauthorized {

	return &ListServicesUnauthorized{}
}

// WithPayload adds the payload to the list services unauthorized response
func (o *ListServicesUnauthorized) WithPayload(payload *models.Error) *ListServicesUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list services unauthorized response
func (o *ListServicesUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListServicesUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListServicesNotFoundCode is the HTTP code returned for type ListServicesNotFound
const ListServicesNotFoundCode int = 404

/*ListServicesNotFound The Models cannot be found

swagger:response listServicesNotFound
*/
type ListServicesNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListServicesNotFound creates ListServicesNotFound with default headers values
func NewListServicesNotFound() *ListServicesNotFound {

	return &ListServicesNotFound{}
}

// WithPayload adds the payload to the list services not found response
func (o *ListServicesNotFound) WithPayload(payload *models.Error) *ListServicesNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list services not found response
func (o *ListServicesNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListServicesNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
