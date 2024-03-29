// Code generated by go-swagger; DO NOT EDIT.

package auths

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// CheckGroupByUserOKCode is the HTTP code returned for type CheckGroupByUserOK
const CheckGroupByUserOKCode int = 200

/*CheckGroupByUserOK auth group.

swagger:response checkGroupByUserOK
*/
type CheckGroupByUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.Result `json:"body,omitempty"`
}

// NewCheckGroupByUserOK creates CheckGroupByUserOK with default headers values
func NewCheckGroupByUserOK() *CheckGroupByUserOK {

	return &CheckGroupByUserOK{}
}

// WithPayload adds the payload to the check group by user o k response
func (o *CheckGroupByUserOK) WithPayload(payload *models.Result) *CheckGroupByUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check group by user o k response
func (o *CheckGroupByUserOK) SetPayload(payload *models.Result) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckGroupByUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckGroupByUserUnauthorizedCode is the HTTP code returned for type CheckGroupByUserUnauthorized
const CheckGroupByUserUnauthorizedCode int = 401

/*CheckGroupByUserUnauthorized Unauthorized

swagger:response checkGroupByUserUnauthorized
*/
type CheckGroupByUserUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckGroupByUserUnauthorized creates CheckGroupByUserUnauthorized with default headers values
func NewCheckGroupByUserUnauthorized() *CheckGroupByUserUnauthorized {

	return &CheckGroupByUserUnauthorized{}
}

// WithPayload adds the payload to the check group by user unauthorized response
func (o *CheckGroupByUserUnauthorized) WithPayload(payload *models.Error) *CheckGroupByUserUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check group by user unauthorized response
func (o *CheckGroupByUserUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckGroupByUserUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckGroupByUserNotFoundCode is the HTTP code returned for type CheckGroupByUserNotFound
const CheckGroupByUserNotFoundCode int = 404

/*CheckGroupByUserNotFound not found.

swagger:response checkGroupByUserNotFound
*/
type CheckGroupByUserNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckGroupByUserNotFound creates CheckGroupByUserNotFound with default headers values
func NewCheckGroupByUserNotFound() *CheckGroupByUserNotFound {

	return &CheckGroupByUserNotFound{}
}

// WithPayload adds the payload to the check group by user not found response
func (o *CheckGroupByUserNotFound) WithPayload(payload *models.Error) *CheckGroupByUserNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check group by user not found response
func (o *CheckGroupByUserNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckGroupByUserNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
