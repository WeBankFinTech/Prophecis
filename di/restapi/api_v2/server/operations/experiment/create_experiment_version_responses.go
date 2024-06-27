// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v2/restmodels"
)

// CreateExperimentVersionOKCode is the HTTP code returned for type CreateExperimentVersionOK
const CreateExperimentVersionOKCode int = 200

/*CreateExperimentVersionOK OK

swagger:response createExperimentVersionOK
*/
type CreateExperimentVersionOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.CreateExperimentVersionResponse `json:"body,omitempty"`
}

// NewCreateExperimentVersionOK creates CreateExperimentVersionOK with default headers values
func NewCreateExperimentVersionOK() *CreateExperimentVersionOK {

	return &CreateExperimentVersionOK{}
}

// WithPayload adds the payload to the create experiment version o k response
func (o *CreateExperimentVersionOK) WithPayload(payload *restmodels.CreateExperimentVersionResponse) *CreateExperimentVersionOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create experiment version o k response
func (o *CreateExperimentVersionOK) SetPayload(payload *restmodels.CreateExperimentVersionResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateExperimentVersionOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateExperimentVersionUnauthorizedCode is the HTTP code returned for type CreateExperimentVersionUnauthorized
const CreateExperimentVersionUnauthorizedCode int = 401

/*CreateExperimentVersionUnauthorized Unauthorized

swagger:response createExperimentVersionUnauthorized
*/
type CreateExperimentVersionUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCreateExperimentVersionUnauthorized creates CreateExperimentVersionUnauthorized with default headers values
func NewCreateExperimentVersionUnauthorized() *CreateExperimentVersionUnauthorized {

	return &CreateExperimentVersionUnauthorized{}
}

// WithPayload adds the payload to the create experiment version unauthorized response
func (o *CreateExperimentVersionUnauthorized) WithPayload(payload *restmodels.Error) *CreateExperimentVersionUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create experiment version unauthorized response
func (o *CreateExperimentVersionUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateExperimentVersionUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateExperimentVersionNotFoundCode is the HTTP code returned for type CreateExperimentVersionNotFound
const CreateExperimentVersionNotFoundCode int = 404

/*CreateExperimentVersionNotFound The Models cannot be found

swagger:response createExperimentVersionNotFound
*/
type CreateExperimentVersionNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCreateExperimentVersionNotFound creates CreateExperimentVersionNotFound with default headers values
func NewCreateExperimentVersionNotFound() *CreateExperimentVersionNotFound {

	return &CreateExperimentVersionNotFound{}
}

// WithPayload adds the payload to the create experiment version not found response
func (o *CreateExperimentVersionNotFound) WithPayload(payload *restmodels.Error) *CreateExperimentVersionNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create experiment version not found response
func (o *CreateExperimentVersionNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateExperimentVersionNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
