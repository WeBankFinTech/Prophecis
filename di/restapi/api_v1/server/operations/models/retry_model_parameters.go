// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewRetryModelParams creates a new RetryModelParams object
// with the default values initialized.
func NewRetryModelParams() RetryModelParams {

	var (
		// initialize parameters with default values

		versionDefault = string("2017-02-13")
	)

	return RetryModelParams{
		Version: &versionDefault,
	}
}

// RetryModelParams contains all the bound params for the retry model operation
// typically these are obtained from a http.Request
//
// swagger:parameters retryModel
type RetryModelParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The id of the model.
	  Required: true
	  In: path
	*/
	ModelID string
	/*The release date of the version of the API you want to use. Specify dates in YYYY-MM-DD format.
	  In: query
	  Default: "2017-02-13"
	*/
	Version *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewRetryModelParams() beforehand.
func (o *RetryModelParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rModelID, rhkModelID, _ := route.Params.GetOK("model_id")
	if err := o.bindModelID(rModelID, rhkModelID, route.Formats); err != nil {
		res = append(res, err)
	}

	qVersion, qhkVersion, _ := qs.GetOK("version")
	if err := o.bindVersion(qVersion, qhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindModelID binds and validates parameter ModelID from path.
func (o *RetryModelParams) bindModelID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ModelID = raw

	return nil
}

// bindVersion binds and validates parameter Version from query.
func (o *RetryModelParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewRetryModelParams()
		return nil
	}

	o.Version = &raw

	return nil
}
