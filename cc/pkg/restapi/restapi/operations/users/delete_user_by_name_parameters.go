// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDeleteUserByNameParams creates a new DeleteUserByNameParams object
// no default values defined in spec.
func NewDeleteUserByNameParams() DeleteUserByNameParams {

	return DeleteUserByNameParams{}
}

// DeleteUserByNameParams contains all the bound params for the delete user by name operation
// typically these are obtained from a http.Request
//
// swagger:parameters DeleteUserByName
type DeleteUserByNameParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*the id of user
	  Required: true
	  In: path
	*/
	UserName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteUserByNameParams() beforehand.
func (o *DeleteUserByNameParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rUserName, rhkUserName, _ := route.Params.GetOK("userName")
	if err := o.bindUserName(rUserName, rhkUserName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindUserName binds and validates parameter UserName from path.
func (o *DeleteUserByNameParams) bindUserName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.UserName = raw

	return nil
}
