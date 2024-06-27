// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v2/restmodels"
)

// UploadExperimentFlowJSONOKCode is the HTTP code returned for type UploadExperimentFlowJSONOK
const UploadExperimentFlowJSONOKCode int = 200

/*UploadExperimentFlowJSONOK OK

swagger:response uploadExperimentFlowJsonOK
*/
type UploadExperimentFlowJSONOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.UploadExperimentFlowJSONResponse `json:"body,omitempty"`
}

// NewUploadExperimentFlowJSONOK creates UploadExperimentFlowJSONOK with default headers values
func NewUploadExperimentFlowJSONOK() *UploadExperimentFlowJSONOK {

	return &UploadExperimentFlowJSONOK{}
}

// WithPayload adds the payload to the upload experiment flow Json o k response
func (o *UploadExperimentFlowJSONOK) WithPayload(payload *restmodels.UploadExperimentFlowJSONResponse) *UploadExperimentFlowJSONOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload experiment flow Json o k response
func (o *UploadExperimentFlowJSONOK) SetPayload(payload *restmodels.UploadExperimentFlowJSONResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadExperimentFlowJSONOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UploadExperimentFlowJSONUnauthorizedCode is the HTTP code returned for type UploadExperimentFlowJSONUnauthorized
const UploadExperimentFlowJSONUnauthorizedCode int = 401

/*UploadExperimentFlowJSONUnauthorized Unauthorized

swagger:response uploadExperimentFlowJsonUnauthorized
*/
type UploadExperimentFlowJSONUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUploadExperimentFlowJSONUnauthorized creates UploadExperimentFlowJSONUnauthorized with default headers values
func NewUploadExperimentFlowJSONUnauthorized() *UploadExperimentFlowJSONUnauthorized {

	return &UploadExperimentFlowJSONUnauthorized{}
}

// WithPayload adds the payload to the upload experiment flow Json unauthorized response
func (o *UploadExperimentFlowJSONUnauthorized) WithPayload(payload *restmodels.Error) *UploadExperimentFlowJSONUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload experiment flow Json unauthorized response
func (o *UploadExperimentFlowJSONUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadExperimentFlowJSONUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UploadExperimentFlowJSONNotFoundCode is the HTTP code returned for type UploadExperimentFlowJSONNotFound
const UploadExperimentFlowJSONNotFoundCode int = 404

/*UploadExperimentFlowJSONNotFound The Models cannot be found

swagger:response uploadExperimentFlowJsonNotFound
*/
type UploadExperimentFlowJSONNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewUploadExperimentFlowJSONNotFound creates UploadExperimentFlowJSONNotFound with default headers values
func NewUploadExperimentFlowJSONNotFound() *UploadExperimentFlowJSONNotFound {

	return &UploadExperimentFlowJSONNotFound{}
}

// WithPayload adds the payload to the upload experiment flow Json not found response
func (o *UploadExperimentFlowJSONNotFound) WithPayload(payload *restmodels.Error) *UploadExperimentFlowJSONNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the upload experiment flow Json not found response
func (o *UploadExperimentFlowJSONNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UploadExperimentFlowJSONNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}