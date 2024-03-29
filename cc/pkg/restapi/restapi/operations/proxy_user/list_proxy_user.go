// Code generated by go-swagger; DO NOT EDIT.

package proxy_user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// ListProxyUserHandlerFunc turns a function with the right signature into a list proxy user handler
type ListProxyUserHandlerFunc func(ListProxyUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListProxyUserHandlerFunc) Handle(params ListProxyUserParams) middleware.Responder {
	return fn(params)
}

// ListProxyUserHandler interface for that can handle valid list proxy user params
type ListProxyUserHandler interface {
	Handle(ListProxyUserParams) middleware.Responder
}

// NewListProxyUser creates a new http.Handler for the list proxy user operation
func NewListProxyUser(ctx *middleware.Context, handler ListProxyUserHandler) *ListProxyUser {
	return &ListProxyUser{Context: ctx, Handler: handler}
}

/*ListProxyUser swagger:route GET /cc/v1/proxyUsers/{user} proxyUser listProxyUser

list proxy user.

Optional extended description in Markdown.

*/
type ListProxyUser struct {
	Context *middleware.Context
	Handler ListProxyUserHandler
}

func (o *ListProxyUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListProxyUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
