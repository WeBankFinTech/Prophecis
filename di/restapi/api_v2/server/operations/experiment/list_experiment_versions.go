// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ListExperimentVersionsHandlerFunc turns a function with the right signature into a list experiment versions handler
type ListExperimentVersionsHandlerFunc func(ListExperimentVersionsParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn ListExperimentVersionsHandlerFunc) Handle(params ListExperimentVersionsParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// ListExperimentVersionsHandler interface for that can handle valid list experiment versions params
type ListExperimentVersionsHandler interface {
	Handle(ListExperimentVersionsParams, interface{}) middleware.Responder
}

// NewListExperimentVersions creates a new http.Handler for the list experiment versions operation
func NewListExperimentVersions(ctx *middleware.Context, handler ListExperimentVersionsHandler) *ListExperimentVersions {
	return &ListExperimentVersions{Context: ctx, Handler: handler}
}

/*ListExperimentVersions swagger:route GET /di/v2/experiment/{exp_id}/versions Experiment listExperimentVersions

ListExperimentVersions

获取实验的所有版本信息

*/
type ListExperimentVersions struct {
	Context *middleware.Context
	Handler ListExperimentVersionsHandler
}

func (o *ListExperimentVersions) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListExperimentVersionsParams()

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
