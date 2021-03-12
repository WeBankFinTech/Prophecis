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

// PostModelReader is a Reader for the PostModel structure.
type PostModelReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostModelReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewPostModelCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewPostModelBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewPostModelUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewPostModelCreated creates a PostModelCreated with default headers values
func NewPostModelCreated() *PostModelCreated {
	return &PostModelCreated{}
}

/*PostModelCreated handles this case with default header values.

Deep learning model successfully accepted.
*/
type PostModelCreated struct {
	/*Location header containing the model id.
	 */
	Location string

	Payload *restmodels.BasicNewModel
}

func (o *PostModelCreated) Error() string {
	return fmt.Sprintf("[POST /v1/models][%d] postModelCreated  %+v", 201, o.Payload)
}

func (o *PostModelCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Location
	o.Location = response.GetHeader("Location")

	o.Payload = new(restmodels.BasicNewModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostModelBadRequest creates a PostModelBadRequest with default headers values
func NewPostModelBadRequest() *PostModelBadRequest {
	return &PostModelBadRequest{}
}

/*PostModelBadRequest handles this case with default header values.

Error in the the model_definition or manifest.
*/
type PostModelBadRequest struct {
	Payload *restmodels.Error
}

func (o *PostModelBadRequest) Error() string {
	return fmt.Sprintf("[POST /v1/models][%d] postModelBadRequest  %+v", 400, o.Payload)
}

func (o *PostModelBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostModelUnauthorized creates a PostModelUnauthorized with default headers values
func NewPostModelUnauthorized() *PostModelUnauthorized {
	return &PostModelUnauthorized{}
}

/*PostModelUnauthorized handles this case with default header values.

Unauthorized
*/
type PostModelUnauthorized struct {
	Payload *restmodels.Error
}

func (o *PostModelUnauthorized) Error() string {
	return fmt.Sprintf("[POST /v1/models][%d] postModelUnauthorized  %+v", 401, o.Payload)
}

func (o *PostModelUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
