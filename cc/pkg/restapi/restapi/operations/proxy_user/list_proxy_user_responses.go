// Code generated by go-swagger; DO NOT EDIT.

package proxy_user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// ListProxyUserOKCode is the HTTP code returned for type ListProxyUserOK
const ListProxyUserOKCode int = 200

/*ListProxyUserOK keyPair interceptor.

swagger:response listProxyUserOK
*/
type ListProxyUserOK struct {

	/*
	  In: Body
	*/
	Payload models.ProxyUsers `json:"body,omitempty"`
}

// NewListProxyUserOK creates ListProxyUserOK with default headers values
func NewListProxyUserOK() *ListProxyUserOK {

	return &ListProxyUserOK{}
}

// WithPayload adds the payload to the list proxy user o k response
func (o *ListProxyUserOK) WithPayload(payload models.ProxyUsers) *ListProxyUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list proxy user o k response
func (o *ListProxyUserOK) SetPayload(payload models.ProxyUsers) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListProxyUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = models.ProxyUsers{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// ListProxyUserUnauthorizedCode is the HTTP code returned for type ListProxyUserUnauthorized
const ListProxyUserUnauthorizedCode int = 401

/*ListProxyUserUnauthorized Unauthorized

swagger:response listProxyUserUnauthorized
*/
type ListProxyUserUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListProxyUserUnauthorized creates ListProxyUserUnauthorized with default headers values
func NewListProxyUserUnauthorized() *ListProxyUserUnauthorized {

	return &ListProxyUserUnauthorized{}
}

// WithPayload adds the payload to the list proxy user unauthorized response
func (o *ListProxyUserUnauthorized) WithPayload(payload *models.Error) *ListProxyUserUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list proxy user unauthorized response
func (o *ListProxyUserUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListProxyUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListProxyUserNotFoundCode is the HTTP code returned for type ListProxyUserNotFound
const ListProxyUserNotFoundCode int = 404

/*ListProxyUserNotFound url to add keyPair not found.

swagger:response listProxyUserNotFound
*/
type ListProxyUserNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewListProxyUserNotFound creates ListProxyUserNotFound with default headers values
func NewListProxyUserNotFound() *ListProxyUserNotFound {

	return &ListProxyUserNotFound{}
}

// WithPayload adds the payload to the list proxy user not found response
func (o *ListProxyUserNotFound) WithPayload(payload *models.Error) *ListProxyUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list proxy user not found response
func (o *ListProxyUserNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListProxyUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
