// Code generated by go-swagger; DO NOT EDIT.

package training_data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	restmodels "webank/DI/restapi/api_v1/restmodels"
)

// GetLoglinesReader is a Reader for the GetLoglines structure.
type GetLoglinesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetLoglinesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetLoglinesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetLoglinesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewGetLoglinesOK creates a GetLoglinesOK with default headers values
func NewGetLoglinesOK() *GetLoglinesOK {
	return &GetLoglinesOK{}
}

/*GetLoglinesOK handles this case with default header values.

(streaming responses)
*/
type GetLoglinesOK struct {
	Payload *restmodels.V1LogLinesList
}

func (o *GetLoglinesOK) Error() string {
	return fmt.Sprintf("[GET /di/v1/logs/{model_id}/loglines][%d] getLoglinesOK  %+v", 200, o.Payload)
}

func (o *GetLoglinesOK) GetPayload() *restmodels.V1LogLinesList {
	return o.Payload
}

func (o *GetLoglinesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.V1LogLinesList)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetLoglinesUnauthorized creates a GetLoglinesUnauthorized with default headers values
func NewGetLoglinesUnauthorized() *GetLoglinesUnauthorized {
	return &GetLoglinesUnauthorized{}
}

/*GetLoglinesUnauthorized handles this case with default header values.

Unauthorized
*/
type GetLoglinesUnauthorized struct {
	Payload *restmodels.Error
}

func (o *GetLoglinesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /di/v1/logs/{model_id}/loglines][%d] getLoglinesUnauthorized  %+v", 401, o.Payload)
}

func (o *GetLoglinesUnauthorized) GetPayload() *restmodels.Error {
	return o.Payload
}

func (o *GetLoglinesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(restmodels.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
