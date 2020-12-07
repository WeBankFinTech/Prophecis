// Code generated by go-swagger; DO NOT EDIT.

package samples

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// SampleGetOKCode is the HTTP code returned for type SampleGetOK
const SampleGetOKCode int = 200

/*SampleGetOK Detailed role and role information.

swagger:response sampleGetOK
*/
type SampleGetOK struct {

	/*
	  In: Body
	*/
	Payload *models.Result `json:"body,omitempty"`
}

// NewSampleGetOK creates SampleGetOK with default headers values
func NewSampleGetOK() *SampleGetOK {

	return &SampleGetOK{}
}

// WithPayload adds the payload to the sample get o k response
func (o *SampleGetOK) WithPayload(payload *models.Result) *SampleGetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the sample get o k response
func (o *SampleGetOK) SetPayload(payload *models.Result) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SampleGetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SampleGetUnauthorizedCode is the HTTP code returned for type SampleGetUnauthorized
const SampleGetUnauthorizedCode int = 401

/*SampleGetUnauthorized Unauthorized

swagger:response sampleGetUnauthorized
*/
type SampleGetUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSampleGetUnauthorized creates SampleGetUnauthorized with default headers values
func NewSampleGetUnauthorized() *SampleGetUnauthorized {

	return &SampleGetUnauthorized{}
}

// WithPayload adds the payload to the sample get unauthorized response
func (o *SampleGetUnauthorized) WithPayload(payload *models.Error) *SampleGetUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the sample get unauthorized response
func (o *SampleGetUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SampleGetUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// SampleGetNotFoundCode is the HTTP code returned for type SampleGetNotFound
const SampleGetNotFoundCode int = 404

/*SampleGetNotFound url to add role not found.

swagger:response sampleGetNotFound
*/
type SampleGetNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSampleGetNotFound creates SampleGetNotFound with default headers values
func NewSampleGetNotFound() *SampleGetNotFound {

	return &SampleGetNotFound{}
}

// WithPayload adds the payload to the sample get not found response
func (o *SampleGetNotFound) WithPayload(payload *models.Error) *SampleGetNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the sample get not found response
func (o *SampleGetNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SampleGetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
