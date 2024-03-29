// Code generated by go-swagger; DO NOT EDIT.

package model_deploy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// CreateNamespacedServiceRunHandlerFunc turns a function with the right signature into a create namespaced service run handler
type CreateNamespacedServiceRunHandlerFunc func(CreateNamespacedServiceRunParams) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateNamespacedServiceRunHandlerFunc) Handle(params CreateNamespacedServiceRunParams) middleware.Responder {
	return fn(params)
}

// CreateNamespacedServiceRunHandler interface for that can handle valid create namespaced service run params
type CreateNamespacedServiceRunHandler interface {
	Handle(CreateNamespacedServiceRunParams) middleware.Responder
}

// NewCreateNamespacedServiceRun creates a new http.Handler for the create namespaced service run operation
func NewCreateNamespacedServiceRun(ctx *middleware.Context, handler CreateNamespacedServiceRunHandler) *CreateNamespacedServiceRun {
	return &CreateNamespacedServiceRun{Context: ctx, Handler: handler}
}

/*CreateNamespacedServiceRun swagger:route POST /mf/v1/serviceRun modelDeploy createNamespacedServiceRun

Run a Service CRD Object.

Run Service.

*/
type CreateNamespacedServiceRun struct {
	Context *middleware.Context
	Handler CreateNamespacedServiceRunHandler
}

func (o *CreateNamespacedServiceRun) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewCreateNamespacedServiceRunParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
