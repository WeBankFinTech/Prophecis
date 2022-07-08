// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// CreateEventEndpointOKCode is the HTTP code returned for type CreateEventEndpointOK
const CreateEventEndpointOKCode int = 200

/*CreateEventEndpointOK Event updated successfully

swagger:response createEventEndpointOK
*/
type CreateEventEndpointOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewCreateEventEndpointOK creates CreateEventEndpointOK with default headers values
func NewCreateEventEndpointOK() *CreateEventEndpointOK {

	return &CreateEventEndpointOK{}
}

// WithPayload adds the payload to the create event endpoint o k response
func (o *CreateEventEndpointOK) WithPayload(payload io.ReadCloser) *CreateEventEndpointOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create event endpoint o k response
func (o *CreateEventEndpointOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEventEndpointOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// CreateEventEndpointUnauthorizedCode is the HTTP code returned for type CreateEventEndpointUnauthorized
const CreateEventEndpointUnauthorizedCode int = 401

/*CreateEventEndpointUnauthorized Unauthorized

swagger:response createEventEndpointUnauthorized
*/
type CreateEventEndpointUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCreateEventEndpointUnauthorized creates CreateEventEndpointUnauthorized with default headers values
func NewCreateEventEndpointUnauthorized() *CreateEventEndpointUnauthorized {

	return &CreateEventEndpointUnauthorized{}
}

// WithPayload adds the payload to the create event endpoint unauthorized response
func (o *CreateEventEndpointUnauthorized) WithPayload(payload *restmodels.Error) *CreateEventEndpointUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create event endpoint unauthorized response
func (o *CreateEventEndpointUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEventEndpointUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEventEndpointNotFoundCode is the HTTP code returned for type CreateEventEndpointNotFound
const CreateEventEndpointNotFoundCode int = 404

/*CreateEventEndpointNotFound The model or event type cannot be found.

swagger:response createEventEndpointNotFound
*/
type CreateEventEndpointNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCreateEventEndpointNotFound creates CreateEventEndpointNotFound with default headers values
func NewCreateEventEndpointNotFound() *CreateEventEndpointNotFound {

	return &CreateEventEndpointNotFound{}
}

// WithPayload adds the payload to the create event endpoint not found response
func (o *CreateEventEndpointNotFound) WithPayload(payload *restmodels.Error) *CreateEventEndpointNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create event endpoint not found response
func (o *CreateEventEndpointNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEventEndpointNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateEventEndpointGoneCode is the HTTP code returned for type CreateEventEndpointGone
const CreateEventEndpointGoneCode int = 410

/*CreateEventEndpointGone If the trained model storage time has expired and it has been deleted. It only gets deleted if it not stored on an external data store.

swagger:response createEventEndpointGone
*/
type CreateEventEndpointGone struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCreateEventEndpointGone creates CreateEventEndpointGone with default headers values
func NewCreateEventEndpointGone() *CreateEventEndpointGone {

	return &CreateEventEndpointGone{}
}

// WithPayload adds the payload to the create event endpoint gone response
func (o *CreateEventEndpointGone) WithPayload(payload *restmodels.Error) *CreateEventEndpointGone {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create event endpoint gone response
func (o *CreateEventEndpointGone) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateEventEndpointGone) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(410)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
