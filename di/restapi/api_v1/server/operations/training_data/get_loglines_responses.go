// Code generated by go-swagger; DO NOT EDIT.

package training_data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// GetLoglinesOKCode is the HTTP code returned for type GetLoglinesOK
const GetLoglinesOKCode int = 200

/*GetLoglinesOK (streaming responses)

swagger:response getLoglinesOK
*/
type GetLoglinesOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.V1LogLinesList `json:"body,omitempty"`
}

// NewGetLoglinesOK creates GetLoglinesOK with default headers values
func NewGetLoglinesOK() *GetLoglinesOK {

	return &GetLoglinesOK{}
}

// WithPayload adds the payload to the get loglines o k response
func (o *GetLoglinesOK) WithPayload(payload *restmodels.V1LogLinesList) *GetLoglinesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get loglines o k response
func (o *GetLoglinesOK) SetPayload(payload *restmodels.V1LogLinesList) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLoglinesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetLoglinesUnauthorizedCode is the HTTP code returned for type GetLoglinesUnauthorized
const GetLoglinesUnauthorizedCode int = 401

/*GetLoglinesUnauthorized Unauthorized

swagger:response getLoglinesUnauthorized
*/
type GetLoglinesUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewGetLoglinesUnauthorized creates GetLoglinesUnauthorized with default headers values
func NewGetLoglinesUnauthorized() *GetLoglinesUnauthorized {

	return &GetLoglinesUnauthorized{}
}

// WithPayload adds the payload to the get loglines unauthorized response
func (o *GetLoglinesUnauthorized) WithPayload(payload *restmodels.Error) *GetLoglinesUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get loglines unauthorized response
func (o *GetLoglinesUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetLoglinesUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
