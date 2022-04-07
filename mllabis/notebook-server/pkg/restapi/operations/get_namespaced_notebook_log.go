// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetNamespacedNotebookLogHandlerFunc turns a function with the right signature into a get namespaced notebook log handler
type GetNamespacedNotebookLogHandlerFunc func(GetNamespacedNotebookLogParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetNamespacedNotebookLogHandlerFunc) Handle(params GetNamespacedNotebookLogParams) middleware.Responder {
	return fn(params)
}

// GetNamespacedNotebookLogHandler interface for that can handle valid get namespaced notebook log params
type GetNamespacedNotebookLogHandler interface {
	Handle(GetNamespacedNotebookLogParams) middleware.Responder
}

// NewGetNamespacedNotebookLog creates a new http.Handler for the get namespaced notebook log operation
func NewGetNamespacedNotebookLog(ctx *middleware.Context, handler GetNamespacedNotebookLogHandler) *GetNamespacedNotebookLog {
	return &GetNamespacedNotebookLog{Context: ctx, Handler: handler}
}

/*GetNamespacedNotebookLog swagger:route GET /aide/v1/namespaces/{namespace}/notebooks/{notebook_name}/log getNamespacedNotebookLog

Get log of notebook in the given namespace

Get notebook log.

*/
type GetNamespacedNotebookLog struct {
	Context *middleware.Context
	Handler GetNamespacedNotebookLogHandler
}

func (o *GetNamespacedNotebookLog) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetNamespacedNotebookLogParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
