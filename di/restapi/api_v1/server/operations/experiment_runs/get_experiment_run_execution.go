// Code generated by go-swagger; DO NOT EDIT.

package experiment_runs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetExperimentRunExecutionHandlerFunc turns a function with the right signature into a get experiment run execution handler
type GetExperimentRunExecutionHandlerFunc func(GetExperimentRunExecutionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetExperimentRunExecutionHandlerFunc) Handle(params GetExperimentRunExecutionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetExperimentRunExecutionHandler interface for that can handle valid get experiment run execution params
type GetExperimentRunExecutionHandler interface {
	Handle(GetExperimentRunExecutionParams, interface{}) middleware.Responder
}

// NewGetExperimentRunExecution creates a new http.Handler for the get experiment run execution operation
func NewGetExperimentRunExecution(ctx *middleware.Context, handler GetExperimentRunExecutionHandler) *GetExperimentRunExecution {
	return &GetExperimentRunExecution{Context: ctx, Handler: handler}
}

/*GetExperimentRunExecution swagger:route GET /di/v1/experimentRun/{exec_id}/execution ExperimentRuns getExperimentRunExecution

Linkis Execution Method.

Get Flow Execution Msg From /api/entrance/${exec_id}/execution

*/
type GetExperimentRunExecution struct {
	Context *middleware.Context
	Handler GetExperimentRunExecutionHandler
}

func (o *GetExperimentRunExecution) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetExperimentRunExecutionParams()

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