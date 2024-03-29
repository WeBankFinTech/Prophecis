// Code generated by go-swagger; DO NOT EDIT.

package report

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-mf/pkg/models"
)

// GetPushEventByIDOKCode is the HTTP code returned for type GetPushEventByIDOK
const GetPushEventByIDOKCode int = 200

/*GetPushEventByIDOK OK

swagger:response getPushEventByIdOK
*/
type GetPushEventByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Event `json:"body,omitempty"`
}

// NewGetPushEventByIDOK creates GetPushEventByIDOK with default headers values
func NewGetPushEventByIDOK() *GetPushEventByIDOK {

	return &GetPushEventByIDOK{}
}

// WithPayload adds the payload to the get push event by Id o k response
func (o *GetPushEventByIDOK) WithPayload(payload *models.Event) *GetPushEventByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get push event by Id o k response
func (o *GetPushEventByIDOK) SetPayload(payload *models.Event) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPushEventByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPushEventByIDUnauthorizedCode is the HTTP code returned for type GetPushEventByIDUnauthorized
const GetPushEventByIDUnauthorizedCode int = 401

/*GetPushEventByIDUnauthorized Unauthorized

swagger:response getPushEventByIdUnauthorized
*/
type GetPushEventByIDUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetPushEventByIDUnauthorized creates GetPushEventByIDUnauthorized with default headers values
func NewGetPushEventByIDUnauthorized() *GetPushEventByIDUnauthorized {

	return &GetPushEventByIDUnauthorized{}
}

// WithPayload adds the payload to the get push event by Id unauthorized response
func (o *GetPushEventByIDUnauthorized) WithPayload(payload *models.Error) *GetPushEventByIDUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get push event by Id unauthorized response
func (o *GetPushEventByIDUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPushEventByIDUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPushEventByIDNotFoundCode is the HTTP code returned for type GetPushEventByIDNotFound
const GetPushEventByIDNotFoundCode int = 404

/*GetPushEventByIDNotFound Report get fail

swagger:response getPushEventByIdNotFound
*/
type GetPushEventByIDNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetPushEventByIDNotFound creates GetPushEventByIDNotFound with default headers values
func NewGetPushEventByIDNotFound() *GetPushEventByIDNotFound {

	return &GetPushEventByIDNotFound{}
}

// WithPayload adds the payload to the get push event by Id not found response
func (o *GetPushEventByIDNotFound) WithPayload(payload *models.Error) *GetPushEventByIDNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get push event by Id not found response
func (o *GetPushEventByIDNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPushEventByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
