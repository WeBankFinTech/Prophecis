// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetExperimentVersionGlobalVariablesStrHandlerFunc turns a function with the right signature into a get experiment version global variables str handler
type GetExperimentVersionGlobalVariablesStrHandlerFunc func(GetExperimentVersionGlobalVariablesStrParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetExperimentVersionGlobalVariablesStrHandlerFunc) Handle(params GetExperimentVersionGlobalVariablesStrParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetExperimentVersionGlobalVariablesStrHandler interface for that can handle valid get experiment version global variables str params
type GetExperimentVersionGlobalVariablesStrHandler interface {
	Handle(GetExperimentVersionGlobalVariablesStrParams, interface{}) middleware.Responder
}

// NewGetExperimentVersionGlobalVariablesStr creates a new http.Handler for the get experiment version global variables str operation
func NewGetExperimentVersionGlobalVariablesStr(ctx *middleware.Context, handler GetExperimentVersionGlobalVariablesStrHandler) *GetExperimentVersionGlobalVariablesStr {
	return &GetExperimentVersionGlobalVariablesStr{Context: ctx, Handler: handler}
}

/*GetExperimentVersionGlobalVariablesStr swagger:route GET /di/v2/experiment/{exp_id}/version/{version_name}/global_variables_str Experiment getExperimentVersionGlobalVariablesStr

GetExperimentVersionGlobalVariablesStr

获取某个版本的实验的工作流的默认全局变量

*/
type GetExperimentVersionGlobalVariablesStr struct {
	Context *middleware.Context
	Handler GetExperimentVersionGlobalVariablesStrHandler
}

func (o *GetExperimentVersionGlobalVariablesStr) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetExperimentVersionGlobalVariablesStrParams()

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