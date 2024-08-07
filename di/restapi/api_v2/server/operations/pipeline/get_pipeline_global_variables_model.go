// Code generated by go-swagger; DO NOT EDIT.

package pipeline

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetPipelineGlobalVariablesModelHandlerFunc turns a function with the right signature into a get pipeline global variables model handler
type GetPipelineGlobalVariablesModelHandlerFunc func(GetPipelineGlobalVariablesModelParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPipelineGlobalVariablesModelHandlerFunc) Handle(params GetPipelineGlobalVariablesModelParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetPipelineGlobalVariablesModelHandler interface for that can handle valid get pipeline global variables model params
type GetPipelineGlobalVariablesModelHandler interface {
	Handle(GetPipelineGlobalVariablesModelParams, interface{}) middleware.Responder
}

// NewGetPipelineGlobalVariablesModel creates a new http.Handler for the get pipeline global variables model operation
func NewGetPipelineGlobalVariablesModel(ctx *middleware.Context, handler GetPipelineGlobalVariablesModelHandler) *GetPipelineGlobalVariablesModel {
	return &GetPipelineGlobalVariablesModel{Context: ctx, Handler: handler}
}

/*GetPipelineGlobalVariablesModel swagger:route GET /di/v2/pipeline/global_variables/model Pipeline getPipelineGlobalVariablesModel

GetPipelineGlobalVariablesModel

获取可视化编辑/工作流中的全局变量的模型变量数据结构

*/
type GetPipelineGlobalVariablesModel struct {
	Context *middleware.Context
	Handler GetPipelineGlobalVariablesModelHandler
}

func (o *GetPipelineGlobalVariablesModel) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPipelineGlobalVariablesModelParams()

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
