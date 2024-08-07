// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetExperimentVersionHandlerFunc turns a function with the right signature into a get experiment version handler
type GetExperimentVersionHandlerFunc func(GetExperimentVersionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetExperimentVersionHandlerFunc) Handle(params GetExperimentVersionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetExperimentVersionHandler interface for that can handle valid get experiment version params
type GetExperimentVersionHandler interface {
	Handle(GetExperimentVersionParams, interface{}) middleware.Responder
}

// NewGetExperimentVersion creates a new http.Handler for the get experiment version operation
func NewGetExperimentVersion(ctx *middleware.Context, handler GetExperimentVersionHandler) *GetExperimentVersion {
	return &GetExperimentVersion{Context: ctx, Handler: handler}
}

/*GetExperimentVersion swagger:route GET /di/v2/experiment/{exp_id}/version/{version_name} Experiment getExperimentVersion

GetExperimentVersion

获取某个版本的实验信息（如果传入的版本为最新版本，同 GetExperiment 接口效果一样）

*/
type GetExperimentVersion struct {
	Context *middleware.Context
	Handler GetExperimentVersionHandler
}

func (o *GetExperimentVersion) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetExperimentVersionParams()

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
