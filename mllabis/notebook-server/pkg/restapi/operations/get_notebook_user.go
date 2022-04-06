// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetNotebookUserHandlerFunc turns a function with the right signature into a get notebook user handler
type GetNotebookUserHandlerFunc func(GetNotebookUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetNotebookUserHandlerFunc) Handle(params GetNotebookUserParams) middleware.Responder {
	return fn(params)
}

// GetNotebookUserHandler interface for that can handle valid get notebook user params
type GetNotebookUserHandler interface {
	Handle(GetNotebookUserParams) middleware.Responder
}

// NewGetNotebookUser creates a new http.Handler for the get notebook user operation
func NewGetNotebookUser(ctx *middleware.Context, handler GetNotebookUserHandler) *GetNotebookUser {
	return &GetNotebookUser{Context: ctx, Handler: handler}
}

/*GetNotebookUser swagger:route GET /aide/v1/notebook/user/{Namespace}/{Name} getNotebookUser

Get user of notebook.

Get user of notebook.

*/
type GetNotebookUser struct {
	Context *middleware.Context
	Handler GetNotebookUserHandler
}

func (o *GetNotebookUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetNotebookUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
