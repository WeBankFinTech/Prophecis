// Code generated by go-swagger; DO NOT EDIT.

package training_data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetLoglinesParams creates a new GetLoglinesParams object
// with the default values initialized.
func NewGetLoglinesParams() GetLoglinesParams {

	var (
		// initialize parameters with default values

		searchTypeDefault = string("TERM")
		sinceTimeDefault  = string("")
		versionDefault    = string("2017-10-01")
	)

	return GetLoglinesParams{
		SearchType: &searchTypeDefault,

		SinceTime: &sinceTimeDefault,

		Version: &versionDefault,
	}
}

// GetLoglinesParams contains all the bound params for the get loglines operation
// typically these are obtained from a http.Request
//
// swagger:parameters getLoglines
type GetLoglinesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*The id of the model.
	  Required: true
	  In: path
	*/
	ModelID string
	/*Number of lines to output.
	  In: query
	*/
	Pagesize *int32
	/*If positive, line number from start, if negative line counting from end
	  In: query
	*/
	Pos *int64
	/*
	  In: query
	  Default: "TERM"
	*/
	SearchType *string
	/*Time from which to show logs. If this value precedes the time a pod was started, only logs since the pod start will be returned. If this value is in the future, no logs will be returned. If this value is a raw integer, it represents the time that the metric occured: representing the number of milliseconds since midnight January 1, 1970. If this value is a negative integer, it represents the number of lines to count backwards. If this value is empty, the logs since the beginning of the job will be returned
	  In: query
	  Default: ""
	*/
	SinceTime *string
	/*The release date of the version of the API you want to use. Specify dates in YYYY-MM-DD format.
	  In: query
	  Default: "2017-10-01"
	*/
	Version *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetLoglinesParams() beforehand.
func (o *GetLoglinesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rModelID, rhkModelID, _ := route.Params.GetOK("model_id")
	if err := o.bindModelID(rModelID, rhkModelID, route.Formats); err != nil {
		res = append(res, err)
	}

	qPagesize, qhkPagesize, _ := qs.GetOK("pagesize")
	if err := o.bindPagesize(qPagesize, qhkPagesize, route.Formats); err != nil {
		res = append(res, err)
	}

	qPos, qhkPos, _ := qs.GetOK("pos")
	if err := o.bindPos(qPos, qhkPos, route.Formats); err != nil {
		res = append(res, err)
	}

	qSearchType, qhkSearchType, _ := qs.GetOK("searchType")
	if err := o.bindSearchType(qSearchType, qhkSearchType, route.Formats); err != nil {
		res = append(res, err)
	}

	qSinceTime, qhkSinceTime, _ := qs.GetOK("since_time")
	if err := o.bindSinceTime(qSinceTime, qhkSinceTime, route.Formats); err != nil {
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
func (o *GetLoglinesParams) bindModelID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ModelID = raw

	return nil
}

// bindPagesize binds and validates parameter Pagesize from query.
func (o *GetLoglinesParams) bindPagesize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("pagesize", "query", "int32", raw)
	}
	o.Pagesize = &value

	return nil
}

// bindPos binds and validates parameter Pos from query.
func (o *GetLoglinesParams) bindPos(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
		return errors.InvalidType("pos", "query", "int64", raw)
	}
	o.Pos = &value

	return nil
}

// bindSearchType binds and validates parameter SearchType from query.
func (o *GetLoglinesParams) bindSearchType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetLoglinesParams()
		return nil
	}

	o.SearchType = &raw

	if err := o.validateSearchType(formats); err != nil {
		return err
	}

	return nil
}

// validateSearchType carries on validations for parameter SearchType
func (o *GetLoglinesParams) validateSearchType(formats strfmt.Registry) error {

	if err := validate.Enum("searchType", "query", *o.SearchType, []interface{}{"TERM", "NESTED", "MATCH"}); err != nil {
		return err
	}

	return nil
}

// bindSinceTime binds and validates parameter SinceTime from query.
func (o *GetLoglinesParams) bindSinceTime(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetLoglinesParams()
		return nil
	}

	o.SinceTime = &raw

	return nil
}

// bindVersion binds and validates parameter Version from query.
func (o *GetLoglinesParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetLoglinesParams()
		return nil
	}

	o.Version = &raw

	return nil
}
