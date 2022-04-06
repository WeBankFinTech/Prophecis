// Code generated by go-swagger; DO NOT EDIT.

package model_storage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewUploadModelParams creates a new UploadModelParams object
// no default values defined in spec.
func NewUploadModelParams() UploadModelParams {

	return UploadModelParams{}
}

// UploadModelParams contains all the bound params for the upload model operation
// typically these are obtained from a http.Request
//
// swagger:parameters uploadModel
type UploadModelParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: formData
	*/
	File io.ReadCloser
	/*
	  Required: true
	  In: formData
	*/
	FileName string
	/*
	  Required: true
	  In: formData
	*/
	ModelType string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUploadModelParams() beforehand.
func (o *UploadModelParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}
	fds := runtime.Values(r.Form)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "file", err))
	} else if err := o.bindFile(file, fileHeader); err != nil {
		// Required: true
		res = append(res, err)
	} else {
		o.File = &runtime.File{Data: file, Header: fileHeader}
	}

	fdFileName, fdhkFileName, _ := fds.GetOK("fileName")
	if err := o.bindFileName(fdFileName, fdhkFileName, route.Formats); err != nil {
		res = append(res, err)
	}

	fdModelType, fdhkModelType, _ := fds.GetOK("modelType")
	if err := o.bindModelType(fdModelType, fdhkModelType, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindFile binds file parameter File.
//
// The only supported validations on files are MinLength and MaxLength
func (o *UploadModelParams) bindFile(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

// bindFileName binds and validates parameter FileName from formData.
func (o *UploadModelParams) bindFileName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("fileName", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("fileName", "formData", raw); err != nil {
		return err
	}

	o.FileName = raw

	return nil
}

// bindModelType binds and validates parameter ModelType from formData.
func (o *UploadModelParams) bindModelType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("modelType", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("modelType", "formData", raw); err != nil {
		return err
	}

	o.ModelType = raw

	return nil
}
