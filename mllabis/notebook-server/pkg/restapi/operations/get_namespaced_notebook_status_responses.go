// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "webank/AIDE/notebook-server/pkg/models"
)

// GetNamespacedNotebookStatusOKCode is the HTTP code returned for type GetNamespacedNotebookStatusOK
const GetNamespacedNotebookStatusOKCode int = 200

/*GetNamespacedNotebookStatusOK OK

swagger:response getNamespacedNotebookStatusOK
*/
type GetNamespacedNotebookStatusOK struct {

	/*
	  In: Body
	*/
	Payload *models.GetNotebookStatusResponse `json:"body,omitempty"`
}

// NewGetNamespacedNotebookStatusOK creates GetNamespacedNotebookStatusOK with default headers values
func NewGetNamespacedNotebookStatusOK() *GetNamespacedNotebookStatusOK {

	return &GetNamespacedNotebookStatusOK{}
}

// WithPayload adds the payload to the get namespaced notebook status o k response
func (o *GetNamespacedNotebookStatusOK) WithPayload(payload *models.GetNotebookStatusResponse) *GetNamespacedNotebookStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebook status o k response
func (o *GetNamespacedNotebookStatusOK) SetPayload(payload *models.GetNotebookStatusResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebookStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNamespacedNotebookStatusUnauthorizedCode is the HTTP code returned for type GetNamespacedNotebookStatusUnauthorized
const GetNamespacedNotebookStatusUnauthorizedCode int = 401

/*GetNamespacedNotebookStatusUnauthorized Unauthorized

swagger:response getNamespacedNotebookStatusUnauthorized
*/
type GetNamespacedNotebookStatusUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNamespacedNotebookStatusUnauthorized creates GetNamespacedNotebookStatusUnauthorized with default headers values
func NewGetNamespacedNotebookStatusUnauthorized() *GetNamespacedNotebookStatusUnauthorized {

	return &GetNamespacedNotebookStatusUnauthorized{}
}

// WithPayload adds the payload to the get namespaced notebook status unauthorized response
func (o *GetNamespacedNotebookStatusUnauthorized) WithPayload(payload *models.Error) *GetNamespacedNotebookStatusUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebook status unauthorized response
func (o *GetNamespacedNotebookStatusUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebookStatusUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNamespacedNotebookStatusNotFoundCode is the HTTP code returned for type GetNamespacedNotebookStatusNotFound
const GetNamespacedNotebookStatusNotFoundCode int = 404

/*GetNamespacedNotebookStatusNotFound The Notebook cannot be found

swagger:response getNamespacedNotebookStatusNotFound
*/
type GetNamespacedNotebookStatusNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNamespacedNotebookStatusNotFound creates GetNamespacedNotebookStatusNotFound with default headers values
func NewGetNamespacedNotebookStatusNotFound() *GetNamespacedNotebookStatusNotFound {

	return &GetNamespacedNotebookStatusNotFound{}
}

// WithPayload adds the payload to the get namespaced notebook status not found response
func (o *GetNamespacedNotebookStatusNotFound) WithPayload(payload *models.Error) *GetNamespacedNotebookStatusNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebook status not found response
func (o *GetNamespacedNotebookStatusNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebookStatusNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
