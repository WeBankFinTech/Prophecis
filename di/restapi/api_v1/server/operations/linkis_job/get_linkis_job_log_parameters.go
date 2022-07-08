// Code generated by go-swagger; DO NOT EDIT.

package linkis_job

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetLinkisJobLogParams creates a new GetLinkisJobLogParams object
// no default values defined in spec.
func NewGetLinkisJobLogParams() GetLinkisJobLogParams {

	return GetLinkisJobLogParams{}
}

// GetLinkisJobLogParams contains all the bound params for the get linkis job log operation
// typically these are obtained from a http.Request
//
// swagger:parameters getLinkisJobLog
type GetLinkisJobLogParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	LinkisTaskID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetLinkisJobLogParams() beforehand.
func (o *GetLinkisJobLogParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rLinkisTaskID, rhkLinkisTaskID, _ := route.Params.GetOK("linkis_task_id")
	if err := o.bindLinkisTaskID(rLinkisTaskID, rhkLinkisTaskID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindLinkisTaskID binds and validates parameter LinkisTaskID from path.
func (o *GetLinkisJobLogParams) bindLinkisTaskID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("linkis_task_id", "path", "int64", raw)
	}
	o.LinkisTaskID = value

	return nil
}
