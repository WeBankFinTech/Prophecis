// Code generated by go-swagger; DO NOT EDIT.

package model_storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateModelHandlerFunc turns a function with the right signature into a update model handler
type UpdateModelHandlerFunc func(UpdateModelParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateModelHandlerFunc) Handle(params UpdateModelParams) middleware.Responder {
	return fn(params)
}

// UpdateModelHandler interface for that can handle valid update model params
type UpdateModelHandler interface {
	Handle(UpdateModelParams) middleware.Responder
}

// NewUpdateModel creates a new http.Handler for the update model operation
func NewUpdateModel(ctx *middleware.Context, handler UpdateModelHandler) *UpdateModel {
	return &UpdateModel{Context: ctx, Handler: handler}
}

/*UpdateModel swagger:route PUT /mf/v1/model/{modelID} modelStorage updateModel

Update a Model in the given Model ID

Update Model.

*/
type UpdateModel struct {
	Context *middleware.Context
	Handler UpdateModelHandler
}

func (o *UpdateModel) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateModelParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
