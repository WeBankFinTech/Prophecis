// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// StopNotebookByIDHandlerFunc turns a function with the right signature into a stop notebook by Id handler
type StopNotebookByIDHandlerFunc func(StopNotebookByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn StopNotebookByIDHandlerFunc) Handle(params StopNotebookByIDParams) middleware.Responder {
	return fn(params)
}

// StopNotebookByIDHandler interface for that can handle valid stop notebook by Id params
type StopNotebookByIDHandler interface {
	Handle(StopNotebookByIDParams) middleware.Responder
}

// NewStopNotebookByID creates a new http.Handler for the stop notebook by Id operation
func NewStopNotebookByID(ctx *middleware.Context, handler StopNotebookByIDHandler) *StopNotebookByID {
	return &StopNotebookByID{Context: ctx, Handler: handler}
}

/*StopNotebookByID swagger:route DELETE /aide/v1/notebooks/{id}/stop stopNotebookById

stop a Notebook, request information from mongo

stop Notebooks.

*/
type StopNotebookByID struct {
	Context *middleware.Context
	Handler StopNotebookByIDHandler
}

func (o *StopNotebookByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewStopNotebookByIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
