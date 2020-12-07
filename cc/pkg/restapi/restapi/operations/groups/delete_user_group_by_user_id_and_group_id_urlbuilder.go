// Code generated by go-swagger; DO NOT EDIT.

package groups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/swag"
)

// DeleteUserGroupByUserIDAndGroupIDURL generates an URL for the delete user group by user Id and group Id operation
type DeleteUserGroupByUserIDAndGroupIDURL struct {
	GroupID int64
	UserID  int64

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteUserGroupByUserIDAndGroupIDURL) WithBasePath(bp string) *DeleteUserGroupByUserIDAndGroupIDURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteUserGroupByUserIDAndGroupIDURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *DeleteUserGroupByUserIDAndGroupIDURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/cc/v1/groups/userGroup/user/{userId}/group/{groupId}"

	groupID := swag.FormatInt64(o.GroupID)
	if groupID != "" {
		_path = strings.Replace(_path, "{groupId}", groupID, -1)
	} else {
		return nil, errors.New("groupId is required on DeleteUserGroupByUserIDAndGroupIDURL")
	}

	userID := swag.FormatInt64(o.UserID)
	if userID != "" {
		_path = strings.Replace(_path, "{userId}", userID, -1)
	} else {
		return nil, errors.New("userId is required on DeleteUserGroupByUserIDAndGroupIDURL")
	}

	_basePath := o._basePath
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *DeleteUserGroupByUserIDAndGroupIDURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *DeleteUserGroupByUserIDAndGroupIDURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *DeleteUserGroupByUserIDAndGroupIDURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on DeleteUserGroupByUserIDAndGroupIDURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on DeleteUserGroupByUserIDAndGroupIDURL")
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
func (o *DeleteUserGroupByUserIDAndGroupIDURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
