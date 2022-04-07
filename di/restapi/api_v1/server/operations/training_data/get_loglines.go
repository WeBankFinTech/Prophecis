// Code generated by go-swagger; DO NOT EDIT.

package training_data

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetLoglinesHandlerFunc turns a function with the right signature into a get loglines handler
type GetLoglinesHandlerFunc func(GetLoglinesParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLoglinesHandlerFunc) Handle(params GetLoglinesParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// GetLoglinesHandler interface for that can handle valid get loglines params
type GetLoglinesHandler interface {
	Handle(GetLoglinesParams, interface{}) middleware.Responder
}

// NewGetLoglines creates a new http.Handler for the get loglines operation
func NewGetLoglines(ctx *middleware.Context, handler GetLoglinesHandler) *GetLoglines {
	return &GetLoglines{Context: ctx, Handler: handler}
}

/*GetLoglines swagger:route GET /di/v1/logs/{model_id}/loglines TrainingData getLoglines

Get loglines, based on query

*/
type GetLoglines struct {
	Context *middleware.Context
	Handler GetLoglinesHandler
}

func (o *GetLoglines) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetLoglinesParams()

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
