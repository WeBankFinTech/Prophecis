// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	restmodels "webank/DI/restapi/api_v2/restmodels"
)

// GetPipelineModelPredictNodeManifestOKCode is the HTTP code returned for type GetPipelineModelPredictNodeManifestOK
const GetPipelineModelPredictNodeManifestOKCode int = 200

/*GetPipelineModelPredictNodeManifestOK OK

swagger:response getPipelineModelPredictNodeManifestOK
*/
type GetPipelineModelPredictNodeManifestOK struct {

	/*
	  In: Body
	*/
	Payload *restmodels.ModelPredictNodeManifest `json:"body,omitempty"`
}

// NewGetPipelineModelPredictNodeManifestOK creates GetPipelineModelPredictNodeManifestOK with default headers values
func NewGetPipelineModelPredictNodeManifestOK() *GetPipelineModelPredictNodeManifestOK {

	return &GetPipelineModelPredictNodeManifestOK{}
}

// WithPayload adds the payload to the get pipeline model predict node manifest o k response
func (o *GetPipelineModelPredictNodeManifestOK) WithPayload(payload *restmodels.ModelPredictNodeManifest) *GetPipelineModelPredictNodeManifestOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pipeline model predict node manifest o k response
func (o *GetPipelineModelPredictNodeManifestOK) SetPayload(payload *restmodels.ModelPredictNodeManifest) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPipelineModelPredictNodeManifestOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPipelineModelPredictNodeManifestUnauthorizedCode is the HTTP code returned for type GetPipelineModelPredictNodeManifestUnauthorized
const GetPipelineModelPredictNodeManifestUnauthorizedCode int = 401

/*GetPipelineModelPredictNodeManifestUnauthorized Unauthorized

swagger:response getPipelineModelPredictNodeManifestUnauthorized
*/
type GetPipelineModelPredictNodeManifestUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewGetPipelineModelPredictNodeManifestUnauthorized creates GetPipelineModelPredictNodeManifestUnauthorized with default headers values
func NewGetPipelineModelPredictNodeManifestUnauthorized() *GetPipelineModelPredictNodeManifestUnauthorized {

	return &GetPipelineModelPredictNodeManifestUnauthorized{}
}

// WithPayload adds the payload to the get pipeline model predict node manifest unauthorized response
func (o *GetPipelineModelPredictNodeManifestUnauthorized) WithPayload(payload *restmodels.Error) *GetPipelineModelPredictNodeManifestUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pipeline model predict node manifest unauthorized response
func (o *GetPipelineModelPredictNodeManifestUnauthorized) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPipelineModelPredictNodeManifestUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPipelineModelPredictNodeManifestNotFoundCode is the HTTP code returned for type GetPipelineModelPredictNodeManifestNotFound
const GetPipelineModelPredictNodeManifestNotFoundCode int = 404

/*GetPipelineModelPredictNodeManifestNotFound Model create failed

swagger:response getPipelineModelPredictNodeManifestNotFound
*/
type GetPipelineModelPredictNodeManifestNotFound struct {

	/*
	  In: Body
	*/
	Payload *restmodels.Error `json:"body,omitempty"`
}

// NewGetPipelineModelPredictNodeManifestNotFound creates GetPipelineModelPredictNodeManifestNotFound with default headers values
func NewGetPipelineModelPredictNodeManifestNotFound() *GetPipelineModelPredictNodeManifestNotFound {

	return &GetPipelineModelPredictNodeManifestNotFound{}
}

// WithPayload adds the payload to the get pipeline model predict node manifest not found response
func (o *GetPipelineModelPredictNodeManifestNotFound) WithPayload(payload *restmodels.Error) *GetPipelineModelPredictNodeManifestNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get pipeline model predict node manifest not found response
func (o *GetPipelineModelPredictNodeManifestNotFound) SetPayload(payload *restmodels.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPipelineModelPredictNodeManifestNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
