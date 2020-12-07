// Code generated by go-swagger; DO NOT EDIT.

package storages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteStorageByPathHandlerFunc turns a function with the right signature into a delete storage by path handler
type DeleteStorageByPathHandlerFunc func(DeleteStorageByPathParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteStorageByPathHandlerFunc) Handle(params DeleteStorageByPathParams) middleware.Responder {
	return fn(params)
}

// DeleteStorageByPathHandler interface for that can handle valid delete storage by path params
type DeleteStorageByPathHandler interface {
	Handle(DeleteStorageByPathParams) middleware.Responder
}

// NewDeleteStorageByPath creates a new http.Handler for the delete storage by path operation
func NewDeleteStorageByPath(ctx *middleware.Context, handler DeleteStorageByPathHandler) *DeleteStorageByPath {
	return &DeleteStorageByPath{Context: ctx, Handler: handler}
}

/*DeleteStorageByPath swagger:route DELETE /cc/v1/storages/path Storages deleteStorageByPath

Returns a Storage.

Optional extended description in Markdown.

*/
type DeleteStorageByPath struct {
	Context *middleware.Context
	Handler DeleteStorageByPathHandler
}

func (o *DeleteStorageByPath) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteStorageByPathParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
