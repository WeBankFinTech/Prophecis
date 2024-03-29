// Code generated by go-swagger; DO NOT EDIT.

package auths

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// CheckUserByGroupOKCode is the HTTP code returned for type CheckUserByGroupOK
const CheckUserByGroupOKCode int = 200

/*CheckUserByGroupOK auth group.

swagger:response checkUserByGroupOK
*/
type CheckUserByGroupOK struct {

	/*
	  In: Body
	*/
	Payload *models.Result `json:"body,omitempty"`
}

// NewCheckUserByGroupOK creates CheckUserByGroupOK with default headers values
func NewCheckUserByGroupOK() *CheckUserByGroupOK {

	return &CheckUserByGroupOK{}
}

// WithPayload adds the payload to the check user by group o k response
func (o *CheckUserByGroupOK) WithPayload(payload *models.Result) *CheckUserByGroupOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check user by group o k response
func (o *CheckUserByGroupOK) SetPayload(payload *models.Result) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckUserByGroupOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckUserByGroupUnauthorizedCode is the HTTP code returned for type CheckUserByGroupUnauthorized
const CheckUserByGroupUnauthorizedCode int = 401

/*CheckUserByGroupUnauthorized Unauthorized

swagger:response checkUserByGroupUnauthorized
*/
type CheckUserByGroupUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckUserByGroupUnauthorized creates CheckUserByGroupUnauthorized with default headers values
func NewCheckUserByGroupUnauthorized() *CheckUserByGroupUnauthorized {

	return &CheckUserByGroupUnauthorized{}
}

// WithPayload adds the payload to the check user by group unauthorized response
func (o *CheckUserByGroupUnauthorized) WithPayload(payload *models.Error) *CheckUserByGroupUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check user by group unauthorized response
func (o *CheckUserByGroupUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckUserByGroupUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckUserByGroupNotFoundCode is the HTTP code returned for type CheckUserByGroupNotFound
const CheckUserByGroupNotFoundCode int = 404

/*CheckUserByGroupNotFound not found.

swagger:response checkUserByGroupNotFound
*/
type CheckUserByGroupNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckUserByGroupNotFound creates CheckUserByGroupNotFound with default headers values
func NewCheckUserByGroupNotFound() *CheckUserByGroupNotFound {

	return &CheckUserByGroupNotFound{}
}

// WithPayload adds the payload to the check user by group not found response
func (o *CheckUserByGroupNotFound) WithPayload(payload *models.Error) *CheckUserByGroupNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check user by group not found response
func (o *CheckUserByGroupNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckUserByGroupNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
