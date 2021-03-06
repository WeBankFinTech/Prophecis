// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "webank/AIDE/notebook-server/pkg/models"
)

// GetNamespacedNotebooksOKCode is the HTTP code returned for type GetNamespacedNotebooksOK
const GetNamespacedNotebooksOKCode int = 200

/*GetNamespacedNotebooksOK OK

swagger:response getNamespacedNotebooksOK
*/
type GetNamespacedNotebooksOK struct {

	/*
	  In: Body
	*/
	Payload *models.GetNotebooksResponse `json:"body,omitempty"`
}

// NewGetNamespacedNotebooksOK creates GetNamespacedNotebooksOK with default headers values
func NewGetNamespacedNotebooksOK() *GetNamespacedNotebooksOK {

	return &GetNamespacedNotebooksOK{}
}

// WithPayload adds the payload to the get namespaced notebooks o k response
func (o *GetNamespacedNotebooksOK) WithPayload(payload *models.GetNotebooksResponse) *GetNamespacedNotebooksOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebooks o k response
func (o *GetNamespacedNotebooksOK) SetPayload(payload *models.GetNotebooksResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebooksOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNamespacedNotebooksUnauthorizedCode is the HTTP code returned for type GetNamespacedNotebooksUnauthorized
const GetNamespacedNotebooksUnauthorizedCode int = 401

/*GetNamespacedNotebooksUnauthorized Unauthorized

swagger:response getNamespacedNotebooksUnauthorized
*/
type GetNamespacedNotebooksUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNamespacedNotebooksUnauthorized creates GetNamespacedNotebooksUnauthorized with default headers values
func NewGetNamespacedNotebooksUnauthorized() *GetNamespacedNotebooksUnauthorized {

	return &GetNamespacedNotebooksUnauthorized{}
}

// WithPayload adds the payload to the get namespaced notebooks unauthorized response
func (o *GetNamespacedNotebooksUnauthorized) WithPayload(payload *models.Error) *GetNamespacedNotebooksUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebooks unauthorized response
func (o *GetNamespacedNotebooksUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebooksUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNamespacedNotebooksNotFoundCode is the HTTP code returned for type GetNamespacedNotebooksNotFound
const GetNamespacedNotebooksNotFoundCode int = 404

/*GetNamespacedNotebooksNotFound The Notebook cannot be found

swagger:response getNamespacedNotebooksNotFound
*/
type GetNamespacedNotebooksNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNamespacedNotebooksNotFound creates GetNamespacedNotebooksNotFound with default headers values
func NewGetNamespacedNotebooksNotFound() *GetNamespacedNotebooksNotFound {

	return &GetNamespacedNotebooksNotFound{}
}

// WithPayload adds the payload to the get namespaced notebooks not found response
func (o *GetNamespacedNotebooksNotFound) WithPayload(payload *models.Error) *GetNamespacedNotebooksNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get namespaced notebooks not found response
func (o *GetNamespacedNotebooksNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNamespacedNotebooksNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
