// Code generated by go-swagger; DO NOT EDIT.

package experiments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ImportExperimentDssHandlerFunc turns a function with the right signature into a import experiment dss handler
type ImportExperimentDssHandlerFunc func(ImportExperimentDssParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ImportExperimentDssHandlerFunc) Handle(params ImportExperimentDssParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ImportExperimentDssHandler interface for that can handle valid import experiment dss params
type ImportExperimentDssHandler interface {
	Handle(ImportExperimentDssParams, interface{}) middleware.Responder
}

// NewImportExperimentDss creates a new http.Handler for the import experiment dss operation
func NewImportExperimentDss(ctx *middleware.Context, handler ImportExperimentDssHandler) *ImportExperimentDss {
	return &ImportExperimentDss{Context: ctx, Handler: handler}
}

/*ImportExperimentDss swagger:route POST /di/v1/experiment/importdss Experiments importExperimentDss

Import Experiment(DSS)

import Experiment(DSS)

*/
type ImportExperimentDss struct {
	Context *middleware.Context
	Handler ImportExperimentDssHandler
}

func (o *ImportExperimentDss) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewImportExperimentDssParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
