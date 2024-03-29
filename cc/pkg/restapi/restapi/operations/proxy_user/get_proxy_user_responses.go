// Code generated by go-swagger; DO NOT EDIT.

package proxy_user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// GetProxyUserOKCode is the HTTP code returned for type GetProxyUserOK
const GetProxyUserOKCode int = 200

/*GetProxyUserOK Proxy User.

swagger:response getProxyUserOK
*/
type GetProxyUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.ProxyUser `json:"body,omitempty"`
}

// NewGetProxyUserOK creates GetProxyUserOK with default headers values
func NewGetProxyUserOK() *GetProxyUserOK {

	return &GetProxyUserOK{}
}

// WithPayload adds the payload to the get proxy user o k response
func (o *GetProxyUserOK) WithPayload(payload *models.ProxyUser) *GetProxyUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get proxy user o k response
func (o *GetProxyUserOK) SetPayload(payload *models.ProxyUser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProxyUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProxyUserUnauthorizedCode is the HTTP code returned for type GetProxyUserUnauthorized
const GetProxyUserUnauthorizedCode int = 401

/*GetProxyUserUnauthorized Unauthorized

swagger:response getProxyUserUnauthorized
*/
type GetProxyUserUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProxyUserUnauthorized creates GetProxyUserUnauthorized with default headers values
func NewGetProxyUserUnauthorized() *GetProxyUserUnauthorized {

	return &GetProxyUserUnauthorized{}
}

// WithPayload adds the payload to the get proxy user unauthorized response
func (o *GetProxyUserUnauthorized) WithPayload(payload *models.Error) *GetProxyUserUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get proxy user unauthorized response
func (o *GetProxyUserUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProxyUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProxyUserNotFoundCode is the HTTP code returned for type GetProxyUserNotFound
const GetProxyUserNotFoundCode int = 404

/*GetProxyUserNotFound url to add keyPair not found.

swagger:response getProxyUserNotFound
*/
type GetProxyUserNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProxyUserNotFound creates GetProxyUserNotFound with default headers values
func NewGetProxyUserNotFound() *GetProxyUserNotFound {

	return &GetProxyUserNotFound{}
}

// WithPayload adds the payload to the get proxy user not found response
func (o *GetProxyUserNotFound) WithPayload(payload *models.Error) *GetProxyUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get proxy user not found response
func (o *GetProxyUserNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProxyUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
