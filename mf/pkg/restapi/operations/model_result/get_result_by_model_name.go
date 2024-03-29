// Code generated by go-swagger; DO NOT EDIT.

package model_result

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetResultByModelNameHandlerFunc turns a function with the right signature into a get result by model name handler
type GetResultByModelNameHandlerFunc func(GetResultByModelNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetResultByModelNameHandlerFunc) Handle(params GetResultByModelNameParams) middleware.Responder {
	return fn(params)
}

// GetResultByModelNameHandler interface for that can handle valid get result by model name params
type GetResultByModelNameHandler interface {
	Handle(GetResultByModelNameParams) middleware.Responder
}

// NewGetResultByModelName creates a new http.Handler for the get result by model name operation
func NewGetResultByModelName(ctx *middleware.Context, handler GetResultByModelNameHandler) *GetResultByModelName {
	return &GetResultByModelName{Context: ctx, Handler: handler}
}

/*GetResultByModelName swagger:route GET /mf/v1/result/{modelName} model_result getResultByModelName

Get report by model name and model version

*/
type GetResultByModelName struct {
	Context *middleware.Context
	Handler GetResultByModelNameHandler
}

func (o *GetResultByModelName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetResultByModelNameParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
