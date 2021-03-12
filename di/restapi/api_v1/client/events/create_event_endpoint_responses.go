// Code generated by go-swagger; DO NOT EDIT.

package events

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"webank/DI/restapi/api_v1/restmodels"
)

// CreateEventEndpointReader is a Reader for the CreateEventEndpoint structure.
type CreateEventEndpointReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *CreateEventEndpointReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewCreateEventEndpointOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 401:
		result := NewCreateEventEndpointUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewCreateEventEndpointNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 410:
		result := NewCreateEventEndpointGone()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateEventEndpointOK creates a CreateEventEndpointOK with default headers values
func NewCreateEventEndpointOK(writer io.Writer) *CreateEventEndpointOK {
	return &CreateEventEndpointOK{
		Payload: writer,
	}
}

/*CreateEventEndpointOK handles this case with default header values.

Event updated successfully
*/
type CreateEventEndpointOK struct {
	Payload io.Writer
}

func (o *CreateEventEndpointOK) Error() string {
	return fmt.Sprintf("[POST /v1/models/{model_id}/events/{event_type}/{endpoint_id}][%d] createEventEndpointOK  %+v", 200, o.Payload)
}

func (o *CreateEventEndpointOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateEventEndpointUnauthorized creates a CreateEventEndpointUnauthorized with default headers values
func NewCreateEventEndpointUnauthorized() *CreateEventEndpointUnauthorized {
	return &CreateEventEndpointUnauthorized{}
}

/*CreateEventEndpointUnauthorized handles this case with default header values.

Unauthorized
*/
type CreateEventEndpointUnauthorized struct {
	Payload *restmodels.Error
}

func (o *CreateEventEndpointUnauthorized) Error() string {
	return fmt.Sprintf("[POST /v1/models/{model_id}/events/{event_type}/{endpoint_id}][%d] createEventEndpointUnauthorized  %+v", 401, o.Payload)
}

func (o *CreateEventEndpointUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateEventEndpointNotFound creates a CreateEventEndpointNotFound with default headers values
func NewCreateEventEndpointNotFound() *CreateEventEndpointNotFound {
	return &CreateEventEndpointNotFound{}
}

/*CreateEventEndpointNotFound handles this case with default header values.

The model or event type cannot be found.
*/
type CreateEventEndpointNotFound struct {
	Payload *restmodels.Error
}

func (o *CreateEventEndpointNotFound) Error() string {
	return fmt.Sprintf("[POST /v1/models/{model_id}/events/{event_type}/{endpoint_id}][%d] createEventEndpointNotFound  %+v", 404, o.Payload)
}

func (o *CreateEventEndpointNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateEventEndpointGone creates a CreateEventEndpointGone with default headers values
func NewCreateEventEndpointGone() *CreateEventEndpointGone {
	return &CreateEventEndpointGone{}
}

/*CreateEventEndpointGone handles this case with default header values.

If the trained model storage time has expired and it has been deleted. It only gets deleted if it not stored on an external data store.
*/
type CreateEventEndpointGone struct {
	Payload *restmodels.Error
}

func (o *CreateEventEndpointGone) Error() string {
	return fmt.Sprintf("[POST /v1/models/{model_id}/events/{event_type}/{endpoint_id}][%d] createEventEndpointGone  %+v", 410, o.Payload)
}

func (o *CreateEventEndpointGone) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
