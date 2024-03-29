// Code generated by go-swagger; DO NOT EDIT.

package experiment_runs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// GetExperimentRunsHistoryOKCode is the HTTP code returned for type GetExperimentRunsHistoryOK
const GetExperimentRunsHistoryOKCode int = 200

/*GetExperimentRunsHistoryOK OK

swagger:response getExperimentRunsHistoryOK
*/
type GetExperimentRunsHistoryOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.ProphecisExperiments `json:"body,omitempty"`
}

// NewGetExperimentRunsHistoryOK creates GetExperimentRunsHistoryOK with default headers values
func NewGetExperimentRunsHistoryOK() *GetExperimentRunsHistoryOK {

	return &GetExperimentRunsHistoryOK{}
}

// WithPayload adds the payload to the get experiment runs history o k response
func (o *GetExperimentRunsHistoryOK) WithPayload(payload *restmodels.ProphecisExperiments) *GetExperimentRunsHistoryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get experiment runs history o k response
func (o *GetExperimentRunsHistoryOK) SetPayload(payload *restmodels.ProphecisExperiments) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetExperimentRunsHistoryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetExperimentRunsHistoryUnauthorizedCode is the HTTP code returned for type GetExperimentRunsHistoryUnauthorized
const GetExperimentRunsHistoryUnauthorizedCode int = 401

/*GetExperimentRunsHistoryUnauthorized Unauthorized

swagger:response getExperimentRunsHistoryUnauthorized
*/
type GetExperimentRunsHistoryUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewGetExperimentRunsHistoryUnauthorized creates GetExperimentRunsHistoryUnauthorized with default headers values
func NewGetExperimentRunsHistoryUnauthorized() *GetExperimentRunsHistoryUnauthorized {

	return &GetExperimentRunsHistoryUnauthorized{}
}

// WithPayload adds the payload to the get experiment runs history unauthorized response
func (o *GetExperimentRunsHistoryUnauthorized) WithPayload(payload *restmodels.Error) *GetExperimentRunsHistoryUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get experiment runs history unauthorized response
func (o *GetExperimentRunsHistoryUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetExperimentRunsHistoryUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetExperimentRunsHistoryNotFoundCode is the HTTP code returned for type GetExperimentRunsHistoryNotFound
const GetExperimentRunsHistoryNotFoundCode int = 404

/*GetExperimentRunsHistoryNotFound Get ExperimentRun's History Failed

swagger:response getExperimentRunsHistoryNotFound
*/
type GetExperimentRunsHistoryNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewGetExperimentRunsHistoryNotFound creates GetExperimentRunsHistoryNotFound with default headers values
func NewGetExperimentRunsHistoryNotFound() *GetExperimentRunsHistoryNotFound {

	return &GetExperimentRunsHistoryNotFound{}
}

// WithPayload adds the payload to the get experiment runs history not found response
func (o *GetExperimentRunsHistoryNotFound) WithPayload(payload *restmodels.Error) *GetExperimentRunsHistoryNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get experiment runs history not found response
func (o *GetExperimentRunsHistoryNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetExperimentRunsHistoryNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
