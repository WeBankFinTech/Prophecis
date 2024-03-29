// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "webank/AIDE/notebook-server/pkg/models"
)

// NewPatchNamespacedNotebookParams creates a new PatchNamespacedNotebookParams object
// no default values defined in spec.
func NewPatchNamespacedNotebookParams() PatchNamespacedNotebookParams {

	return PatchNamespacedNotebookParams{}
}

// PatchNamespacedNotebookParams contains all the bound params for the patch namespaced notebook operation
// typically these are obtained from a http.Request
//
// swagger:parameters PatchNamespacedNotebook
type PatchNamespacedNotebookParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ID string
	/*The Patch Notebook Request
	  Required: true
	  In: body
	*/
	Notebook *models.PatchNotebookRequest
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPatchNamespacedNotebookParams() beforehand.
func (o *PatchNamespacedNotebookParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.PatchNotebookRequest
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("notebook", "body"))
			} else {
				res = append(res, errors.NewParseError("notebook", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Notebook = &body
			}
		}
	} else {
		res = append(res, errors.Required("notebook", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from path.
func (o *PatchNamespacedNotebookParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ID = raw

	return nil
}
