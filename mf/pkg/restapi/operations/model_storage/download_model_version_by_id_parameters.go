// Code generated by go-swagger; DO NOT EDIT.

package model_storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDownloadModelVersionByIDParams creates a new DownloadModelVersionByIDParams object
// no default values defined in spec.
func NewDownloadModelVersionByIDParams() DownloadModelVersionByIDParams {

	return DownloadModelVersionByIDParams{}
}

// DownloadModelVersionByIDParams contains all the bound params for the download model version by Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters DownloadModelVersionById
type DownloadModelVersionByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ModelVersionID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDownloadModelVersionByIDParams() beforehand.
func (o *DownloadModelVersionByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rModelVersionID, rhkModelVersionID, _ := route.Params.GetOK("modelVersionId")
	if err := o.bindModelVersionID(rModelVersionID, rhkModelVersionID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindModelVersionID binds and validates parameter ModelVersionID from path.
func (o *DownloadModelVersionByIDParams) bindModelVersionID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("modelVersionId", "path", "int64", raw)
	}
	o.ModelVersionID = value

	return nil
}
