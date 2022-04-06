// Code generated by go-swagger; DO NOT EDIT.

package model_deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-mf/pkg/models"
)

// CreateNamespacedServiceRunOKCode is the HTTP code returned for type CreateNamespacedServiceRunOK
const CreateNamespacedServiceRunOKCode int = 200

/*CreateNamespacedServiceRunOK OK

swagger:response createNamespacedServiceRunOK
*/
type CreateNamespacedServiceRunOK struct {
}

// NewCreateNamespacedServiceRunOK creates CreateNamespacedServiceRunOK with default headers values
func NewCreateNamespacedServiceRunOK() *CreateNamespacedServiceRunOK {

	return &CreateNamespacedServiceRunOK{}
}

// WriteResponse to the client
func (o *CreateNamespacedServiceRunOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// CreateNamespacedServiceRunUnauthorizedCode is the HTTP code returned for type CreateNamespacedServiceRunUnauthorized
const CreateNamespacedServiceRunUnauthorizedCode int = 401

/*CreateNamespacedServiceRunUnauthorized Unauthorized

swagger:response createNamespacedServiceRunUnauthorized
*/
type CreateNamespacedServiceRunUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateNamespacedServiceRunUnauthorized creates CreateNamespacedServiceRunUnauthorized with default headers values
func NewCreateNamespacedServiceRunUnauthorized() *CreateNamespacedServiceRunUnauthorized {

	return &CreateNamespacedServiceRunUnauthorized{}
}

// WithPayload adds the payload to the create namespaced service run unauthorized response
func (o *CreateNamespacedServiceRunUnauthorized) WithPayload(payload *models.Error) *CreateNamespacedServiceRunUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create namespaced service run unauthorized response
func (o *CreateNamespacedServiceRunUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNamespacedServiceRunUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateNamespacedServiceRunNotFoundCode is the HTTP code returned for type CreateNamespacedServiceRunNotFound
const CreateNamespacedServiceRunNotFoundCode int = 404

/*CreateNamespacedServiceRunNotFound Model create failed

swagger:response createNamespacedServiceRunNotFound
*/
type CreateNamespacedServiceRunNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateNamespacedServiceRunNotFound creates CreateNamespacedServiceRunNotFound with default headers values
func NewCreateNamespacedServiceRunNotFound() *CreateNamespacedServiceRunNotFound {

	return &CreateNamespacedServiceRunNotFound{}
}

// WithPayload adds the payload to the create namespaced service run not found response
func (o *CreateNamespacedServiceRunNotFound) WithPayload(payload *models.Error) *CreateNamespacedServiceRunNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create namespaced service run not found response
func (o *CreateNamespacedServiceRunNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNamespacedServiceRunNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
