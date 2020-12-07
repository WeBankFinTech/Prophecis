// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// AddGroupHandlerFunc turns a function with the right signature into a add group handler
type AddGroupHandlerFunc func(AddGroupParams) middleware.Responder

// Handle executing the request and returning a response
func (fn AddGroupHandlerFunc) Handle(params AddGroupParams) middleware.Responder {
	return fn(params)
}

// AddGroupHandler interface for that can handle valid add group params
type AddGroupHandler interface {
	Handle(AddGroupParams) middleware.Responder
}

// NewAddGroup creates a new http.Handler for the add group operation
func NewAddGroup(ctx *middleware.Context, handler AddGroupHandler) *AddGroup {
	return &AddGroup{Context: ctx, Handler: handler}
}

/*AddGroup swagger:route POST /cc/v1/groups Groups addGroup

Returns a list of groups.

Optional extended description in Markdown.

*/
type AddGroup struct {
	Context *middleware.Context
	Handler AddGroupHandler
}

func (o *AddGroup) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddGroupParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
