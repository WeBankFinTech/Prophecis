// Code generated by go-swagger; DO NOT EDIT.

package storages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetStorageByPathHandlerFunc turns a function with the right signature into a get storage by path handler
type GetStorageByPathHandlerFunc func(GetStorageByPathParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetStorageByPathHandlerFunc) Handle(params GetStorageByPathParams) middleware.Responder {
	return fn(params)
}

// GetStorageByPathHandler interface for that can handle valid get storage by path params
type GetStorageByPathHandler interface {
	Handle(GetStorageByPathParams) middleware.Responder
}

// NewGetStorageByPath creates a new http.Handler for the get storage by path operation
func NewGetStorageByPath(ctx *middleware.Context, handler GetStorageByPathHandler) *GetStorageByPath {
	return &GetStorageByPath{Context: ctx, Handler: handler}
}

/*GetStorageByPath swagger:route GET /cc/v1/storages/path Storages getStorageByPath

Returns a Storage.

Optional extended description in Markdown.

*/
type GetStorageByPath struct {
	Context *middleware.Context
	Handler GetStorageByPathHandler
}

func (o *GetStorageByPath) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetStorageByPathParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
