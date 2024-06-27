// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v2/restmodels"
)

// CopyExperimentOKCode is the HTTP code returned for type CopyExperimentOK
const CopyExperimentOKCode int = 200

/*CopyExperimentOK OK

swagger:response copyExperimentOK
*/
type CopyExperimentOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.CreateExperimentResponse `json:"body,omitempty"`
}

// NewCopyExperimentOK creates CopyExperimentOK with default headers values
func NewCopyExperimentOK() *CopyExperimentOK {

	return &CopyExperimentOK{}
}

// WithPayload adds the payload to the copy experiment o k response
func (o *CopyExperimentOK) WithPayload(payload *restmodels.CreateExperimentResponse) *CopyExperimentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the copy experiment o k response
func (o *CopyExperimentOK) SetPayload(payload *restmodels.CreateExperimentResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CopyExperimentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CopyExperimentUnauthorizedCode is the HTTP code returned for type CopyExperimentUnauthorized
const CopyExperimentUnauthorizedCode int = 401

/*CopyExperimentUnauthorized Unauthorized

swagger:response copyExperimentUnauthorized
*/
type CopyExperimentUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCopyExperimentUnauthorized creates CopyExperimentUnauthorized with default headers values
func NewCopyExperimentUnauthorized() *CopyExperimentUnauthorized {

	return &CopyExperimentUnauthorized{}
}

// WithPayload adds the payload to the copy experiment unauthorized response
func (o *CopyExperimentUnauthorized) WithPayload(payload *restmodels.Error) *CopyExperimentUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the copy experiment unauthorized response
func (o *CopyExperimentUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CopyExperimentUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CopyExperimentNotFoundCode is the HTTP code returned for type CopyExperimentNotFound
const CopyExperimentNotFoundCode int = 404

/*CopyExperimentNotFound The Models cannot be found

swagger:response copyExperimentNotFound
*/
type CopyExperimentNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewCopyExperimentNotFound creates CopyExperimentNotFound with default headers values
func NewCopyExperimentNotFound() *CopyExperimentNotFound {

	return &CopyExperimentNotFound{}
}

// WithPayload adds the payload to the copy experiment not found response
func (o *CopyExperimentNotFound) WithPayload(payload *restmodels.Error) *CopyExperimentNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the copy experiment not found response
func (o *CopyExperimentNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CopyExperimentNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}