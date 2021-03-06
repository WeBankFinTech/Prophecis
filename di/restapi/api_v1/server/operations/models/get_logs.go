// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetLogsHandlerFunc turns a function with the right signature into a get logs handler
type GetLogsHandlerFunc func(GetLogsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLogsHandlerFunc) Handle(params GetLogsParams) middleware.Responder {
	return fn(params)
}

// GetLogsHandler interface for that can handle valid get logs params
type GetLogsHandler interface {
	Handle(GetLogsParams) middleware.Responder
}

// NewGetLogs creates a new http.Handler for the get logs operation
func NewGetLogs(ctx *middleware.Context, handler GetLogsHandler) *GetLogs {
	return &GetLogs{Context: ctx, Handler: handler}
}

/*GetLogs swagger:route GET /di/v1/models/{model_id}/logs Models getLogs

Get training logs as websocket stream.


Get training logs for the given model as websocket stream. Each message can contain one or more log lines.


*/
type GetLogs struct {
	Context *middleware.Context
	Handler GetLogsHandler
}

func (o *GetLogs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetLogsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
