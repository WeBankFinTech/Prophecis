// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"webank/DI/restapi/api_v1/restmodels"
)

// PatchModelReader is a Reader for the PatchModel structure.
type PatchModelReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchModelReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 202:
		result := NewPatchModelAccepted()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPatchModelBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewPatchModelUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewPatchModelNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPatchModelAccepted creates a PatchModelAccepted with default headers values
func NewPatchModelAccepted() *PatchModelAccepted {
	return &PatchModelAccepted{}
}

/*PatchModelAccepted handles this case with default header values.

Training successfully halted.
*/
type PatchModelAccepted struct {
	Payload *restmodels.BasicModel
}

func (o *PatchModelAccepted) Error() string {
	return fmt.Sprintf("[PATCH /v1/models/{model_id}][%d] patchModelAccepted  %+v", 202, o.Payload)
}

func (o *PatchModelAccepted) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.BasicModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchModelBadRequest creates a PatchModelBadRequest with default headers values
func NewPatchModelBadRequest() *PatchModelBadRequest {
	return &PatchModelBadRequest{}
}

/*PatchModelBadRequest handles this case with default header values.

Incorrect status specified.
*/
type PatchModelBadRequest struct {
	Payload *restmodels.Error
}

func (o *PatchModelBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /v1/models/{model_id}][%d] patchModelBadRequest  %+v", 400, o.Payload)
}

func (o *PatchModelBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchModelUnauthorized creates a PatchModelUnauthorized with default headers values
func NewPatchModelUnauthorized() *PatchModelUnauthorized {
	return &PatchModelUnauthorized{}
}

/*PatchModelUnauthorized handles this case with default header values.

Unauthorized
*/
type PatchModelUnauthorized struct {
	Payload *restmodels.Error
}

func (o *PatchModelUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /v1/models/{model_id}][%d] patchModelUnauthorized  %+v", 401, o.Payload)
}

func (o *PatchModelUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchModelNotFound creates a PatchModelNotFound with default headers values
func NewPatchModelNotFound() *PatchModelNotFound {
	return &PatchModelNotFound{}
}

/*PatchModelNotFound handles this case with default header values.

Model with the given ID not found.
*/
type PatchModelNotFound struct {
	Payload *restmodels.Error
}

func (o *PatchModelNotFound) Error() string {
	return fmt.Sprintf("[PATCH /v1/models/{model_id}][%d] patchModelNotFound  %+v", 404, o.Payload)
}

func (o *PatchModelNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
