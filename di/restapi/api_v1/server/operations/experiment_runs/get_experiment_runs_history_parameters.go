// Code generated by go-swagger; DO NOT EDIT.

package experiment_runs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetExperimentRunsHistoryParams creates a new GetExperimentRunsHistoryParams object
// no default values defined in spec.
func NewGetExperimentRunsHistoryParams() GetExperimentRunsHistoryParams {

	return GetExperimentRunsHistoryParams{}
}

// GetExperimentRunsHistoryParams contains all the bound params for the get experiment runs history operation
// typically these are obtained from a http.Request
//
// swagger:parameters getExperimentRunsHistory
type GetExperimentRunsHistoryParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ExpID int64
	/*
	  In: query
	*/
	Page *int64
	/*
	  In: query
	*/
	QueryStr *string
	/*
	  In: query
	*/
	Size *int64
	/*
	  In: query
	*/
	Username *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetExperimentRunsHistoryParams() beforehand.
func (o *GetExperimentRunsHistoryParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rExpID, rhkExpID, _ := route.Params.GetOK("exp_id")
	if err := o.bindExpID(rExpID, rhkExpID, route.Formats); err != nil {
		res = append(res, err)
	}

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qQueryStr, qhkQueryStr, _ := qs.GetOK("query_str")
	if err := o.bindQueryStr(qQueryStr, qhkQueryStr, route.Formats); err != nil {
		res = append(res, err)
	}

	qSize, qhkSize, _ := qs.GetOK("size")
	if err := o.bindSize(qSize, qhkSize, route.Formats); err != nil {
		res = append(res, err)
	}

	qUsername, qhkUsername, _ := qs.GetOK("username")
	if err := o.bindUsername(qUsername, qhkUsername, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindExpID binds and validates parameter ExpID from path.
func (o *GetExperimentRunsHistoryParams) bindExpID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("exp_id", "path", "int64", raw)
	}
	o.ExpID = value

	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *GetExperimentRunsHistoryParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = &value

	return nil
}

// bindQueryStr binds and validates parameter QueryStr from query.
func (o *GetExperimentRunsHistoryParams) bindQueryStr(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.QueryStr = &raw

	return nil
}

// bindSize binds and validates parameter Size from query.
func (o *GetExperimentRunsHistoryParams) bindSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("size", "query", "int64", raw)
	}
	o.Size = &value

	return nil
}

// bindUsername binds and validates parameter Username from query.
func (o *GetExperimentRunsHistoryParams) bindUsername(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Username = &raw

	return nil
}
