// Code generated by go-swagger; DO NOT EDIT.

package experiments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// CodeUploadReader is a Reader for the CodeUpload structure.
type CodeUploadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CodeUploadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCodeUploadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewCodeUploadUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewCodeUploadNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCodeUploadOK creates a CodeUploadOK with default headers values
func NewCodeUploadOK() *CodeUploadOK {
	return &CodeUploadOK{}
}

/*CodeUploadOK handles this case with default header values.

OK
*/
type CodeUploadOK struct {
	Payload *restmodels.CodeUploadResponse
}

func (o *CodeUploadOK) Error() string {
	return fmt.Sprintf("[POST /di/v1/codeUpload][%d] codeUploadOK  %+v", 200, o.Payload)
}

func (o *CodeUploadOK) GetPayload() *restmodels.CodeUploadResponse {
	return o.Payload
}

func (o *CodeUploadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.CodeUploadResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCodeUploadUnauthorized creates a CodeUploadUnauthorized with default headers values
func NewCodeUploadUnauthorized() *CodeUploadUnauthorized {
	return &CodeUploadUnauthorized{}
}

/*CodeUploadUnauthorized handles this case with default header values.

Unauthorized
*/
type CodeUploadUnauthorized struct {
	Payload *restmodels.Error
}

func (o *CodeUploadUnauthorized) Error() string {
	return fmt.Sprintf("[POST /di/v1/codeUpload][%d] codeUploadUnauthorized  %+v", 401, o.Payload)
}

func (o *CodeUploadUnauthorized) GetPayload() *restmodels.Error {
	return o.Payload
}

func (o *CodeUploadUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCodeUploadNotFound creates a CodeUploadNotFound with default headers values
func NewCodeUploadNotFound() *CodeUploadNotFound {
	return &CodeUploadNotFound{}
}

/*CodeUploadNotFound handles this case with default header values.

Code ZIP upload failed
*/
type CodeUploadNotFound struct {
	Payload *restmodels.Error
}

func (o *CodeUploadNotFound) Error() string {
	return fmt.Sprintf("[POST /di/v1/codeUpload][%d] codeUploadNotFound  %+v", 404, o.Payload)
}

func (o *CodeUploadNotFound) GetPayload() *restmodels.Error {
	return o.Payload
}

func (o *CodeUploadNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
