// Code generated by go-swagger; DO NOT EDIT.

package model_deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetNamespacedServiceRunParams creates a new GetNamespacedServiceRunParams object
// no default values defined in spec.
func NewGetNamespacedServiceRunParams() GetNamespacedServiceRunParams {

	return GetNamespacedServiceRunParams{}
}

// GetNamespacedServiceRunParams contains all the bound params for the get namespaced service run operation
// typically these are obtained from a http.Request
//
// swagger:parameters getNamespacedServiceRun
type GetNamespacedServiceRunParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ID string
	/*
	  Required: true
	  In: query
	*/
	Namespace string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetNamespacedServiceRunParams() beforehand.
func (o *GetNamespacedServiceRunParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	qNamespace, qhkNamespace, _ := qs.GetOK("namespace")
	if err := o.bindNamespace(qNamespace, qhkNamespace, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from path.
func (o *GetNamespacedServiceRunParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ID = raw

	return nil
}

// bindNamespace binds and validates parameter Namespace from query.
func (o *GetNamespacedServiceRunParams) bindNamespace(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("namespace", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("namespace", "query", raw); err != nil {
		return err
	}

	o.Namespace = raw

	return nil
}
