// Code generated by go-swagger; DO NOT EDIT.

package experiments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// ExportExperimentDssOKCode is the HTTP code returned for type ExportExperimentDssOK
const ExportExperimentDssOKCode int = 200

/*ExportExperimentDssOK Export Experiment(Dss) Response definition

swagger:response exportExperimentDssOK
*/
type ExportExperimentDssOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.ProphecisExportExperimentDssResponse `json:"body,omitempty"`
}

// NewExportExperimentDssOK creates ExportExperimentDssOK with default headers values
func NewExportExperimentDssOK() *ExportExperimentDssOK {

	return &ExportExperimentDssOK{}
}

// WithPayload adds the payload to the export experiment dss o k response
func (o *ExportExperimentDssOK) WithPayload(payload *restmodels.ProphecisExportExperimentDssResponse) *ExportExperimentDssOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export experiment dss o k response
func (o *ExportExperimentDssOK) SetPayload(payload *restmodels.ProphecisExportExperimentDssResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportExperimentDssOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExportExperimentDssUnauthorizedCode is the HTTP code returned for type ExportExperimentDssUnauthorized
const ExportExperimentDssUnauthorizedCode int = 401

/*ExportExperimentDssUnauthorized Unauthorized

swagger:response exportExperimentDssUnauthorized
*/
type ExportExperimentDssUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewExportExperimentDssUnauthorized creates ExportExperimentDssUnauthorized with default headers values
func NewExportExperimentDssUnauthorized() *ExportExperimentDssUnauthorized {

	return &ExportExperimentDssUnauthorized{}
}

// WithPayload adds the payload to the export experiment dss unauthorized response
func (o *ExportExperimentDssUnauthorized) WithPayload(payload *restmodels.Error) *ExportExperimentDssUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export experiment dss unauthorized response
func (o *ExportExperimentDssUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportExperimentDssUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ExportExperimentDssNotFoundCode is the HTTP code returned for type ExportExperimentDssNotFound
const ExportExperimentDssNotFoundCode int = 404

/*ExportExperimentDssNotFound The Experiment cannot be found

swagger:response exportExperimentDssNotFound
*/
type ExportExperimentDssNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewExportExperimentDssNotFound creates ExportExperimentDssNotFound with default headers values
func NewExportExperimentDssNotFound() *ExportExperimentDssNotFound {

	return &ExportExperimentDssNotFound{}
}

// WithPayload adds the payload to the export experiment dss not found response
func (o *ExportExperimentDssNotFound) WithPayload(payload *restmodels.Error) *ExportExperimentDssNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the export experiment dss not found response
func (o *ExportExperimentDssNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ExportExperimentDssNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
