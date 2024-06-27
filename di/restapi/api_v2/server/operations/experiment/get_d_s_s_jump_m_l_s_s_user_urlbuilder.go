// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
)

// GetDSSJumpMLSSUserURL generates an URL for the get d s s jump m l s s user operation
type GetDSSJumpMLSSUserURL struct {
	DssUserTicketID string
	DssWorkspaceID  string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetDSSJumpMLSSUserURL) WithBasePath(bp string) *GetDSSJumpMLSSUserURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetDSSJumpMLSSUserURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetDSSJumpMLSSUserURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/di/v2/experiment/dss_jump_mlss_user"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	dssUserTicketIDQ := o.DssUserTicketID
	if dssUserTicketIDQ != "" {
		qs.Set("dss_user_ticket_id", dssUserTicketIDQ)
	}

	dssWorkspaceIDQ := o.DssWorkspaceID
	if dssWorkspaceIDQ != "" {
		qs.Set("dss_workspace_id", dssWorkspaceIDQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetDSSJumpMLSSUserURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetDSSJumpMLSSUserURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetDSSJumpMLSSUserURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetDSSJumpMLSSUserURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetDSSJumpMLSSUserURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetDSSJumpMLSSUserURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}