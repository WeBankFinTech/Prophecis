package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	"webank/DI/restapi/api_v1/restmodels"
)

/*UpdateEventEndpointOK Event updated successfully

swagger:response updateEventEndpointOK
*/
type UpdateEventEndpointOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewUpdateEventEndpointOK creates UpdateEventEndpointOK with default headers values
func NewUpdateEventEndpointOK() *UpdateEventEndpointOK {
	return &UpdateEventEndpointOK{}
}

// WithPayload adds the payload to the update event endpoint o k response
func (o *UpdateEventEndpointOK) WithPayload(payload io.ReadCloser) *UpdateEventEndpointOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update event endpoint o k response
func (o *UpdateEventEndpointOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateEventEndpointOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*UpdateEventEndpointUnauthorized Unauthorized

swagger:response updateEventEndpointUnauthorized
*/
type UpdateEventEndpointUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUpdateEventEndpointUnauthorized creates UpdateEventEndpointUnauthorized with default headers values
func NewUpdateEventEndpointUnauthorized() *UpdateEventEndpointUnauthorized {
	return &UpdateEventEndpointUnauthorized{}
}

// WithPayload adds the payload to the update event endpoint unauthorized response
func (o *UpdateEventEndpointUnauthorized) WithPayload(payload *restmodels.Error) *UpdateEventEndpointUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update event endpoint unauthorized response
func (o *UpdateEventEndpointUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateEventEndpointUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*UpdateEventEndpointNotFound The model or event type cannot be found.

swagger:response updateEventEndpointNotFound
*/
type UpdateEventEndpointNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUpdateEventEndpointNotFound creates UpdateEventEndpointNotFound with default headers values
func NewUpdateEventEndpointNotFound() *UpdateEventEndpointNotFound {
	return &UpdateEventEndpointNotFound{}
}

// WithPayload adds the payload to the update event endpoint not found response
func (o *UpdateEventEndpointNotFound) WithPayload(payload *restmodels.Error) *UpdateEventEndpointNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update event endpoint not found response
func (o *UpdateEventEndpointNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateEventEndpointNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*UpdateEventEndpointGone If the trained model storage time has expired and it has been deleted. It only gets deleted if it not stored on an external data store.

swagger:response updateEventEndpointGone
*/
type UpdateEventEndpointGone struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUpdateEventEndpointGone creates UpdateEventEndpointGone with default headers values
func NewUpdateEventEndpointGone() *UpdateEventEndpointGone {
	return &UpdateEventEndpointGone{}
}

// WithPayload adds the payload to the update event endpoint gone response
func (o *UpdateEventEndpointGone) WithPayload(payload *restmodels.Error) *UpdateEventEndpointGone {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update event endpoint gone response
func (o *UpdateEventEndpointGone) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateEventEndpointGone) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(410)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
