// Code generated by go-swagger; DO NOT EDIT.

package model_storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-mf/pkg/models"
)

// PostModelOKCode is the HTTP code returned for type PostModelOK
const PostModelOKCode int = 200

/*PostModelOK OK

swagger:response postModelOK
*/
type PostModelOK struct {

	/*
	  In: Body
	*/
	Payload *models.PostModelResp `json:"body,omitempty"`
}

// NewPostModelOK creates PostModelOK with default headers values
func NewPostModelOK() *PostModelOK {

	return &PostModelOK{}
}

// WithPayload adds the payload to the post model o k response
func (o *PostModelOK) WithPayload(payload *models.PostModelResp) *PostModelOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post model o k response
func (o *PostModelOK) SetPayload(payload *models.PostModelResp) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostModelOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostModelUnauthorizedCode is the HTTP code returned for type PostModelUnauthorized
const PostModelUnauthorizedCode int = 401

/*PostModelUnauthorized Unauthorized

swagger:response postModelUnauthorized
*/
type PostModelUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostModelUnauthorized creates PostModelUnauthorized with default headers values
func NewPostModelUnauthorized() *PostModelUnauthorized {

	return &PostModelUnauthorized{}
}

// WithPayload adds the payload to the post model unauthorized response
func (o *PostModelUnauthorized) WithPayload(payload *models.Error) *PostModelUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post model unauthorized response
func (o *PostModelUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostModelUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostModelNotFoundCode is the HTTP code returned for type PostModelNotFound
const PostModelNotFoundCode int = 404

/*PostModelNotFound Model create failed

swagger:response postModelNotFound
*/
type PostModelNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPostModelNotFound creates PostModelNotFound with default headers values
func NewPostModelNotFound() *PostModelNotFound {

	return &PostModelNotFound{}
}

// WithPayload adds the payload to the post model not found response
func (o *PostModelNotFound) WithPayload(payload *models.Error) *PostModelNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post model not found response
func (o *PostModelNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostModelNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
