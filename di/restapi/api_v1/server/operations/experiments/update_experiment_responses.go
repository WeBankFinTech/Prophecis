// Code generated by go-swagger; DO NOT EDIT.

package experiments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// UpdateExperimentOKCode is the HTTP code returned for type UpdateExperimentOK
const UpdateExperimentOKCode int = 200

/*UpdateExperimentOK OK

swagger:response updateExperimentOK
*/
type UpdateExperimentOK struct {
}

// NewUpdateExperimentOK creates UpdateExperimentOK with default headers values
func NewUpdateExperimentOK() *UpdateExperimentOK {

	return &UpdateExperimentOK{}
}

// WriteResponse to the client
func (o *UpdateExperimentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// UpdateExperimentUnauthorizedCode is the HTTP code returned for type UpdateExperimentUnauthorized
const UpdateExperimentUnauthorizedCode int = 401

/*UpdateExperimentUnauthorized Unauthorized

swagger:response updateExperimentUnauthorized
*/
type UpdateExperimentUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUpdateExperimentUnauthorized creates UpdateExperimentUnauthorized with default headers values
func NewUpdateExperimentUnauthorized() *UpdateExperimentUnauthorized {

	return &UpdateExperimentUnauthorized{}
}

// WithPayload adds the payload to the update experiment unauthorized response
func (o *UpdateExperimentUnauthorized) WithPayload(payload *restmodels.Error) *UpdateExperimentUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update experiment unauthorized response
func (o *UpdateExperimentUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateExperimentUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateExperimentNotFoundCode is the HTTP code returned for type UpdateExperimentNotFound
const UpdateExperimentNotFoundCode int = 404

/*UpdateExperimentNotFound Model create failed

swagger:response updateExperimentNotFound
*/
type UpdateExperimentNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUpdateExperimentNotFound creates UpdateExperimentNotFound with default headers values
func NewUpdateExperimentNotFound() *UpdateExperimentNotFound {

	return &UpdateExperimentNotFound{}
}

// WithPayload adds the payload to the update experiment not found response
func (o *UpdateExperimentNotFound) WithPayload(payload *restmodels.Error) *UpdateExperimentNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update experiment not found response
func (o *UpdateExperimentNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateExperimentNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
