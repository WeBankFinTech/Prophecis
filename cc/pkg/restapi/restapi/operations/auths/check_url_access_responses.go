// Code generated by go-swagger; DO NOT EDIT.

package auths

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "mlss-controlcenter-go/pkg/models"
)

// CheckURLAccessOKCode is the HTTP code returned for type CheckURLAccessOK
const CheckURLAccessOKCode int = 200

/*CheckURLAccessOK auth by namespace and notebook.

swagger:response checkUrlAccessOK
*/
type CheckURLAccessOK struct {

	/*
	  In: Body
	*/
	Payload *models.Result `json:"body,omitempty"`
}

// NewCheckURLAccessOK creates CheckURLAccessOK with default headers values
func NewCheckURLAccessOK() *CheckURLAccessOK {

	return &CheckURLAccessOK{}
}

// WithPayload adds the payload to the check Url access o k response
func (o *CheckURLAccessOK) WithPayload(payload *models.Result) *CheckURLAccessOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check Url access o k response
func (o *CheckURLAccessOK) SetPayload(payload *models.Result) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckURLAccessOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckURLAccessUnauthorizedCode is the HTTP code returned for type CheckURLAccessUnauthorized
const CheckURLAccessUnauthorizedCode int = 401

/*CheckURLAccessUnauthorized Unauthorized

swagger:response checkUrlAccessUnauthorized
*/
type CheckURLAccessUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckURLAccessUnauthorized creates CheckURLAccessUnauthorized with default headers values
func NewCheckURLAccessUnauthorized() *CheckURLAccessUnauthorized {

	return &CheckURLAccessUnauthorized{}
}

// WithPayload adds the payload to the check Url access unauthorized response
func (o *CheckURLAccessUnauthorized) WithPayload(payload *models.Error) *CheckURLAccessUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check Url access unauthorized response
func (o *CheckURLAccessUnauthorized) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckURLAccessUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CheckURLAccessNotFoundCode is the HTTP code returned for type CheckURLAccessNotFound
const CheckURLAccessNotFoundCode int = 404

/*CheckURLAccessNotFound url to add namespace not found.

swagger:response checkUrlAccessNotFound
*/
type CheckURLAccessNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCheckURLAccessNotFound creates CheckURLAccessNotFound with default headers values
func NewCheckURLAccessNotFound() *CheckURLAccessNotFound {

	return &CheckURLAccessNotFound{}
}

// WithPayload adds the payload to the check Url access not found response
func (o *CheckURLAccessNotFound) WithPayload(payload *models.Error) *CheckURLAccessNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the check Url access not found response
func (o *CheckURLAccessNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CheckURLAccessNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
