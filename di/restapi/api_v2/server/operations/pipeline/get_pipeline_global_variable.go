// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetPipelineGlobalVariableHandlerFunc turns a function with the right signature into a get pipeline global variable handler
type GetPipelineGlobalVariableHandlerFunc func(GetPipelineGlobalVariableParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPipelineGlobalVariableHandlerFunc) Handle(params GetPipelineGlobalVariableParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetPipelineGlobalVariableHandler interface for that can handle valid get pipeline global variable params
type GetPipelineGlobalVariableHandler interface {
	Handle(GetPipelineGlobalVariableParams, interface{}) middleware.Responder
}

// NewGetPipelineGlobalVariable creates a new http.Handler for the get pipeline global variable operation
func NewGetPipelineGlobalVariable(ctx *middleware.Context, handler GetPipelineGlobalVariableHandler) *GetPipelineGlobalVariable {
	return &GetPipelineGlobalVariable{Context: ctx, Handler: handler}
}

/*GetPipelineGlobalVariable swagger:route GET /di/v2/pipeline/global_variable Pipeline getPipelineGlobalVariable

GetPipelineGlobalVariable

获取可视化编辑/工作流中的全局变量的数据结构

*/
type GetPipelineGlobalVariable struct {
	Context *middleware.Context
	Handler GetPipelineGlobalVariableHandler
}

func (o *GetPipelineGlobalVariable) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPipelineGlobalVariableParams()

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
