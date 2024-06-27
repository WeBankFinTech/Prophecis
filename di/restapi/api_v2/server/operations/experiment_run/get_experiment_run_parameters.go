// Code generated by go-swagger; DO NOT EDIT.

package experiment_run

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetExperimentRunParams creates a new GetExperimentRunParams object
// no default values defined in spec.
func NewGetExperimentRunParams() GetExperimentRunParams {

	return GetExperimentRunParams{}
}

// GetExperimentRunParams contains all the bound params for the get experiment run operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetExperimentRun
type GetExperimentRunParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*实验执行的Id
	  Required: true
	  In: path
	*/
	RunID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetExperimentRunParams() beforehand.
func (o *GetExperimentRunParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rRunID, rhkRunID, _ := route.Params.GetOK("run_id")
	if err := o.bindRunID(rRunID, rhkRunID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindRunID binds and validates parameter RunID from path.
func (o *GetExperimentRunParams) bindRunID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.RunID = raw

	return nil
}
