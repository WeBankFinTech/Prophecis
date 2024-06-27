// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetDSSJumpMLSSUserParams creates a new GetDSSJumpMLSSUserParams object
// no default values defined in spec.
func NewGetDSSJumpMLSSUserParams() GetDSSJumpMLSSUserParams {

	return GetDSSJumpMLSSUserParams{}
}

// GetDSSJumpMLSSUserParams contains all the bound params for the get d s s jump m l s s user operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetDSSJumpMLSSUser
type GetDSSJumpMLSSUserParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*dss的url后带的dss_user_ticket_id内容
	  Required: true
	  In: query
	*/
	DssUserTicketID string
	/*dss的url后带的dss_workspace_id内容
	  Required: true
	  In: query
	*/
	DssWorkspaceID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetDSSJumpMLSSUserParams() beforehand.
func (o *GetDSSJumpMLSSUserParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDssUserTicketID, qhkDssUserTicketID, _ := qs.GetOK("dss_user_ticket_id")
	if err := o.bindDssUserTicketID(qDssUserTicketID, qhkDssUserTicketID, route.Formats); err != nil {
		res = append(res, err)
	}

	qDssWorkspaceID, qhkDssWorkspaceID, _ := qs.GetOK("dss_workspace_id")
	if err := o.bindDssWorkspaceID(qDssWorkspaceID, qhkDssWorkspaceID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindDssUserTicketID binds and validates parameter DssUserTicketID from query.
func (o *GetDSSJumpMLSSUserParams) bindDssUserTicketID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("dss_user_ticket_id", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("dss_user_ticket_id", "query", raw); err != nil {
		return err
	}

	o.DssUserTicketID = raw

	return nil
}

// bindDssWorkspaceID binds and validates parameter DssWorkspaceID from query.
func (o *GetDSSJumpMLSSUserParams) bindDssWorkspaceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("dss_workspace_id", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("dss_workspace_id", "query", raw); err != nil {
		return err
	}

	o.DssWorkspaceID = raw

	return nil
}