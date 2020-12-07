// Code generated by go-swagger; DO NOT EDIT.

package resources

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// GetNodeByNameOKCode is the HTTP code returned for type GetNodeByNameOK
const GetNodeByNameOKCode int = 200

/*GetNodeByNameOK Detailed namespace and namespace information.

swagger:response getNodeByNameOK
*/
type GetNodeByNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.LabelsResponse `json:"body,omitempty"`
}

// NewGetNodeByNameOK creates GetNodeByNameOK with default headers values
func NewGetNodeByNameOK() *GetNodeByNameOK {

	return &GetNodeByNameOK{}
}

// WithPayload adds the payload to the get node by name o k response
func (o *GetNodeByNameOK) WithPayload(payload *models.LabelsResponse) *GetNodeByNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get node by name o k response
func (o *GetNodeByNameOK) SetPayload(payload *models.LabelsResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNodeByNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNodeByNameUnauthorizedCode is the HTTP code returned for type GetNodeByNameUnauthorized
const GetNodeByNameUnauthorizedCode int = 401

/*GetNodeByNameUnauthorized Unauthorized

swagger:response getNodeByNameUnauthorized
*/
type GetNodeByNameUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNodeByNameUnauthorized creates GetNodeByNameUnauthorized with default headers values
func NewGetNodeByNameUnauthorized() *GetNodeByNameUnauthorized {

	return &GetNodeByNameUnauthorized{}
}

// WithPayload adds the payload to the get node by name unauthorized response
func (o *GetNodeByNameUnauthorized) WithPayload(payload *models.Error) *GetNodeByNameUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get node by name unauthorized response
func (o *GetNodeByNameUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNodeByNameUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetNodeByNameNotFoundCode is the HTTP code returned for type GetNodeByNameNotFound
const GetNodeByNameNotFoundCode int = 404

/*GetNodeByNameNotFound url to add namespace not found.

swagger:response getNodeByNameNotFound
*/
type GetNodeByNameNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetNodeByNameNotFound creates GetNodeByNameNotFound with default headers values
func NewGetNodeByNameNotFound() *GetNodeByNameNotFound {

	return &GetNodeByNameNotFound{}
}

// WithPayload adds the payload to the get node by name not found response
func (o *GetNodeByNameNotFound) WithPayload(payload *models.Error) *GetNodeByNameNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get node by name not found response
func (o *GetNodeByNameNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNodeByNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
