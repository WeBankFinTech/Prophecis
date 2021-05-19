// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetLogsParams creates a new GetLogsParams object
// with the default values initialized.
func NewGetLogsParams() *GetLogsParams {
	var (
		followDefault    = bool(false)
		sinceTimeDefault = string("")
		versionDefault   = string("2017-02-13")
	)
	return &GetLogsParams{
		Follow:    &followDefault,
		SinceTime: &sinceTimeDefault,
		Version:   versionDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewGetLogsParamsWithTimeout creates a new GetLogsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetLogsParamsWithTimeout(timeout time.Duration) *GetLogsParams {
	var (
		followDefault    = bool(false)
		sinceTimeDefault = string("")
		versionDefault   = string("2017-02-13")
	)
	return &GetLogsParams{
		Follow:    &followDefault,
		SinceTime: &sinceTimeDefault,
		Version:   versionDefault,

		timeout: timeout,
	}
}

// NewGetLogsParamsWithContext creates a new GetLogsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetLogsParamsWithContext(ctx context.Context) *GetLogsParams {
	var (
		followDefault    = bool(false)
		sinceTimeDefault = string("")
		versionDefault   = string("2017-02-13")
	)
	return &GetLogsParams{
		Follow:    &followDefault,
		SinceTime: &sinceTimeDefault,
		Version:   versionDefault,

		Context: ctx,
	}
}

// NewGetLogsParamsWithHTTPClient creates a new GetLogsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetLogsParamsWithHTTPClient(client *http.Client) *GetLogsParams {
	var (
		followDefault    = bool(false)
		sinceTimeDefault = string("")
		versionDefault   = string("2017-02-13")
	)
	return &GetLogsParams{
		Follow:     &followDefault,
		SinceTime:  &sinceTimeDefault,
		Version:    versionDefault,
		HTTPClient: client,
	}
}

/*GetLogsParams contains all the parameters to send to the API endpoint
for the get logs operation typically these are written to a http.Request
*/
type GetLogsParams struct {

	/*Follow
	  Follow the log stream if true. Default false.

	*/
	Follow *bool
	/*ModelID
	  The id of the model.

	*/
	ModelID string
	/*SinceTime
	  An RFC3339 timestamp from which to show logs. If this value precedes the time a pod was started, only logs since the pod start will be returned.  If this value is in the future, no logs will be returned.

	*/
	SinceTime *string
	/*Version
	  The release date of the version of the API you want to use. Specify dates in YYYY-MM-DD format.

	*/
	Version string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get logs params
func (o *GetLogsParams) WithTimeout(timeout time.Duration) *GetLogsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get logs params
func (o *GetLogsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get logs params
func (o *GetLogsParams) WithContext(ctx context.Context) *GetLogsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get logs params
func (o *GetLogsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get logs params
func (o *GetLogsParams) WithHTTPClient(client *http.Client) *GetLogsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get logs params
func (o *GetLogsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFollow adds the follow to the get logs params
func (o *GetLogsParams) WithFollow(follow *bool) *GetLogsParams {
	o.SetFollow(follow)
	return o
}

// SetFollow adds the follow to the get logs params
func (o *GetLogsParams) SetFollow(follow *bool) {
	o.Follow = follow
}

// WithModelID adds the modelID to the get logs params
func (o *GetLogsParams) WithModelID(modelID string) *GetLogsParams {
	o.SetModelID(modelID)
	return o
}

// SetModelID adds the modelId to the get logs params
func (o *GetLogsParams) SetModelID(modelID string) {
	o.ModelID = modelID
}

// WithSinceTime adds the sinceTime to the get logs params
func (o *GetLogsParams) WithSinceTime(sinceTime *string) *GetLogsParams {
	o.SetSinceTime(sinceTime)
	return o
}

// SetSinceTime adds the sinceTime to the get logs params
func (o *GetLogsParams) SetSinceTime(sinceTime *string) {
	o.SinceTime = sinceTime
}

// WithVersion adds the version to the get logs params
func (o *GetLogsParams) WithVersion(version string) *GetLogsParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the get logs params
func (o *GetLogsParams) SetVersion(version string) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *GetLogsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Follow != nil {

		// query param follow
		var qrFollow bool
		if o.Follow != nil {
			qrFollow = *o.Follow
		}
		qFollow := swag.FormatBool(qrFollow)
		if qFollow != "" {
			if err := r.SetQueryParam("follow", qFollow); err != nil {
				return err
			}
		}

	}

	// path param model_id
	if err := r.SetPathParam("model_id", o.ModelID); err != nil {
		return err
	}

	if o.SinceTime != nil {

		// query param since_time
		var qrSinceTime string
		if o.SinceTime != nil {
			qrSinceTime = *o.SinceTime
		}
		qSinceTime := qrSinceTime
		if qSinceTime != "" {
			if err := r.SetQueryParam("since_time", qSinceTime); err != nil {
				return err
			}
		}

	}

	// query param version
	qrVersion := o.Version
	qVersion := qrVersion
	if qVersion != "" {
		if err := r.SetQueryParam("version", qVersion); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}