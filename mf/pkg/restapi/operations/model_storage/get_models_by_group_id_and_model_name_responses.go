// Code generated by go-swagger; DO NOT EDIT.

package model_storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-mf/pkg/models"
)

// GetModelsByGroupIDAndModelNameOKCode is the HTTP code returned for type GetModelsByGroupIDAndModelNameOK
const GetModelsByGroupIDAndModelNameOKCode int = 200

/*GetModelsByGroupIDAndModelNameOK OK

swagger:response getModelsByGroupIdAndModelNameOK
*/
type GetModelsByGroupIDAndModelNameOK struct {

	/*
	  In: Body
	*/
	Payload models.Models `json:"body,omitempty"`
}

// NewGetModelsByGroupIDAndModelNameOK creates GetModelsByGroupIDAndModelNameOK with default headers values
func NewGetModelsByGroupIDAndModelNameOK() *GetModelsByGroupIDAndModelNameOK {

	return &GetModelsByGroupIDAndModelNameOK{}
}

// WithPayload adds the payload to the get models by group Id and model name o k response
func (o *GetModelsByGroupIDAndModelNameOK) WithPayload(payload models.Models) *GetModelsByGroupIDAndModelNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get models by group Id and model name o k response
func (o *GetModelsByGroupIDAndModelNameOK) SetPayload(payload models.Models) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetModelsByGroupIDAndModelNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.Models{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetModelsByGroupIDAndModelNameUnauthorizedCode is the HTTP code returned for type GetModelsByGroupIDAndModelNameUnauthorized
const GetModelsByGroupIDAndModelNameUnauthorizedCode int = 401

/*GetModelsByGroupIDAndModelNameUnauthorized Unauthorized

swagger:response getModelsByGroupIdAndModelNameUnauthorized
*/
type GetModelsByGroupIDAndModelNameUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetModelsByGroupIDAndModelNameUnauthorized creates GetModelsByGroupIDAndModelNameUnauthorized with default headers values
func NewGetModelsByGroupIDAndModelNameUnauthorized() *GetModelsByGroupIDAndModelNameUnauthorized {

	return &GetModelsByGroupIDAndModelNameUnauthorized{}
}

// WithPayload adds the payload to the get models by group Id and model name unauthorized response
func (o *GetModelsByGroupIDAndModelNameUnauthorized) WithPayload(payload *models.Error) *GetModelsByGroupIDAndModelNameUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get models by group Id and model name unauthorized response
func (o *GetModelsByGroupIDAndModelNameUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetModelsByGroupIDAndModelNameUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetModelsByGroupIDAndModelNameNotFoundCode is the HTTP code returned for type GetModelsByGroupIDAndModelNameNotFound
const GetModelsByGroupIDAndModelNameNotFoundCode int = 404

/*GetModelsByGroupIDAndModelNameNotFound The Models cannot be found

swagger:response getModelsByGroupIdAndModelNameNotFound
*/
type GetModelsByGroupIDAndModelNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetModelsByGroupIDAndModelNameNotFound creates GetModelsByGroupIDAndModelNameNotFound with default headers values
func NewGetModelsByGroupIDAndModelNameNotFound() *GetModelsByGroupIDAndModelNameNotFound {

	return &GetModelsByGroupIDAndModelNameNotFound{}
}

// WithPayload adds the payload to the get models by group Id and model name not found response
func (o *GetModelsByGroupIDAndModelNameNotFound) WithPayload(payload *models.Error) *GetModelsByGroupIDAndModelNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get models by group Id and model name not found response
func (o *GetModelsByGroupIDAndModelNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetModelsByGroupIDAndModelNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
