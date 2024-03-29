// Code generated by go-swagger; DO NOT EDIT.

package experiments

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

// NewImportExperimentDssParams creates a new ImportExperimentDssParams object
// with the default values initialized.
func NewImportExperimentDssParams() ImportExperimentDssParams {

	var (
		// initialize parameters with default values

		experimentIDDefault = int64(0)
	)

	return ImportExperimentDssParams{
		ExperimentID: &experimentIDDefault,
	}
}

// ImportExperimentDssParams contains all the bound params for the import experiment dss operation
// typically these are obtained from a http.Request
//
// swagger:parameters importExperimentDss
type ImportExperimentDssParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: formData
	*/
	Desc *string
	/*if experimentId is 0 or not provided, create experiment in db; if exists, get experiment by experimentId, update it.
	  In: formData
	  Default: 0
	*/
	ExperimentID *int64
	/*
	  Required: true
	  In: formData
	*/
	ResourceID string
	/*
	  In: formData
	*/
	Version *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewImportExperimentDssParams() beforehand.
func (o *ImportExperimentDssParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	fdDesc, fdhkDesc, _ := fds.GetOK("desc")
	if err := o.bindDesc(fdDesc, fdhkDesc, route.Formats); err != nil {
		res = append(res, err)
	}

	fdExperimentID, fdhkExperimentID, _ := fds.GetOK("experimentId")
	if err := o.bindExperimentID(fdExperimentID, fdhkExperimentID, route.Formats); err != nil {
		res = append(res, err)
	}

	fdResourceID, fdhkResourceID, _ := fds.GetOK("resourceId")
	if err := o.bindResourceID(fdResourceID, fdhkResourceID, route.Formats); err != nil {
		res = append(res, err)
	}

	fdVersion, fdhkVersion, _ := fds.GetOK("version")
	if err := o.bindVersion(fdVersion, fdhkVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDesc binds and validates parameter Desc from formData.
func (o *ImportExperimentDssParams) bindDesc(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Desc = &raw

	return nil
}

// bindExperimentID binds and validates parameter ExperimentID from formData.
func (o *ImportExperimentDssParams) bindExperimentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewImportExperimentDssParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("experimentId", "formData", "int64", raw)
	}
	o.ExperimentID = &value

	return nil
}

// bindResourceID binds and validates parameter ResourceID from formData.
func (o *ImportExperimentDssParams) bindResourceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("resourceId", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("resourceId", "formData", raw); err != nil {
		return err
	}

	o.ResourceID = raw

	return nil
}

// bindVersion binds and validates parameter Version from formData.
func (o *ImportExperimentDssParams) bindVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Version = &raw

	return nil
}
