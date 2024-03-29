// Code generated by go-swagger; DO NOT EDIT.

package report

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListReportVersionPushEventsByReportVersionIDParams creates a new ListReportVersionPushEventsByReportVersionIDParams object
// with the default values initialized.
func NewListReportVersionPushEventsByReportVersionIDParams() ListReportVersionPushEventsByReportVersionIDParams {

	var (
		// initialize parameters with default values

		currentPageDefault = int64(1)
		pageSizeDefault    = int64(10)
	)

	return ListReportVersionPushEventsByReportVersionIDParams{
		CurrentPage: &currentPageDefault,

		PageSize: &pageSizeDefault,
	}
}

// ListReportVersionPushEventsByReportVersionIDParams contains all the bound params for the list report version push events by report version Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters listReportVersionPushEventsByReportVersionId
type ListReportVersionPushEventsByReportVersionIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	  Default: 1
	*/
	CurrentPage *int64
	/*
	  In: query
	  Default: 10
	*/
	PageSize *int64
	/*
	  Required: true
	  In: path
	*/
	ReportVersionID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewListReportVersionPushEventsByReportVersionIDParams() beforehand.
func (o *ListReportVersionPushEventsByReportVersionIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qCurrentPage, qhkCurrentPage, _ := qs.GetOK("currentPage")
	if err := o.bindCurrentPage(qCurrentPage, qhkCurrentPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageSize, qhkPageSize, _ := qs.GetOK("pageSize")
	if err := o.bindPageSize(qPageSize, qhkPageSize, route.Formats); err != nil {
		res = append(res, err)
	}

	rReportVersionID, rhkReportVersionID, _ := route.Params.GetOK("reportVersionId")
	if err := o.bindReportVersionID(rReportVersionID, rhkReportVersionID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindCurrentPage binds and validates parameter CurrentPage from query.
func (o *ListReportVersionPushEventsByReportVersionIDParams) bindCurrentPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewListReportVersionPushEventsByReportVersionIDParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("currentPage", "query", "int64", raw)
	}
	o.CurrentPage = &value

	return nil
}

// bindPageSize binds and validates parameter PageSize from query.
func (o *ListReportVersionPushEventsByReportVersionIDParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewListReportVersionPushEventsByReportVersionIDParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("pageSize", "query", "int64", raw)
	}
	o.PageSize = &value

	return nil
}

// bindReportVersionID binds and validates parameter ReportVersionID from path.
func (o *ListReportVersionPushEventsByReportVersionIDParams) bindReportVersionID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("reportVersionId", "path", "int64", raw)
	}
	o.ReportVersionID = value

	return nil
}
