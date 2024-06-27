// Code generated by go-swagger; DO NOT EDIT.

package experiment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewCreateExperimentByUploadParams creates a new CreateExperimentByUploadParams object
// with the default values initialized.
func NewCreateExperimentByUploadParams() CreateExperimentByUploadParams {

	var (
		// initialize parameters with default values

		clusterTypeDefault = string("BDAP")

		originSystemDefault = string("MLSS")
	)

	return CreateExperimentByUploadParams{
		ClusterType: &clusterTypeDefault,

		OriginSystem: &originSystemDefault,
	}
}

// CreateExperimentByUploadParams contains all the bound params for the create experiment by upload operation
// typically these are obtained from a http.Request
//
// swagger:parameters CreateExperimentByUpload
type CreateExperimentByUploadParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*实验所属的集群类型，默认为"BDAP"；可选枚举值为 "BDAP" "BDP"
	  In: formData
	  Default: "BDAP"
	*/
	ClusterType *string
	/*创建实验的描述
	  In: formData
	*/
	Description *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的id）
	  In: formData
	*/
	DssFlowID *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的名称）
	  In: formData
	*/
	DssFlowName *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的版本）
	  In: formData
	*/
	DssFlowVersion *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的项目id）
	  In: formData
	*/
	DssProjectID *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的项目名）
	  In: formData
	*/
	DssProjectName *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的工作空间id）
	  In: formData
	*/
	DssWorkspaceID *string
	/*创建DSS类型的实验的dss相关信息（dss工作流的工作空间名称）
	  In: formData
	*/
	DssWorkspaceName *string
	/*上传的表示实验工作流的zip文件，zip文件里应该包含一个名为flow.json的json文件
	  Required: true
	  In: formData
	*/
	File io.ReadCloser
	/*上传文件的文件名
	  In: formData
	*/
	FileName *string
	/*创建实验所属的项目组
	  Required: true
	  Min Length: 1
	  In: formData
	*/
	GroupName string
	/*创建实验的名称
	  Required: true
	  Min Length: 1
	  In: formData
	*/
	Name string
	/*表示实验在哪个系统创建的, 默认为"MLSS"；可选的枚举值为 "WTSS" "DSS" "MLSS"
	  In: formData
	  Default: "MLSS"
	*/
	OriginSystem *string
	/*创建实验的标签
	  In: formData
	*/
	Tags []string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateExperimentByUploadParams() beforehand.
func (o *CreateExperimentByUploadParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	fdClusterType, fdhkClusterType, _ := fds.GetOK("cluster_type")
	if err := o.bindClusterType(fdClusterType, fdhkClusterType, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDescription, fdhkDescription, _ := fds.GetOK("description")
	if err := o.bindDescription(fdDescription, fdhkDescription, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssFlowID, fdhkDssFlowID, _ := fds.GetOK("dss_flow_id")
	if err := o.bindDssFlowID(fdDssFlowID, fdhkDssFlowID, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssFlowName, fdhkDssFlowName, _ := fds.GetOK("dss_flow_name")
	if err := o.bindDssFlowName(fdDssFlowName, fdhkDssFlowName, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssFlowVersion, fdhkDssFlowVersion, _ := fds.GetOK("dss_flow_version")
	if err := o.bindDssFlowVersion(fdDssFlowVersion, fdhkDssFlowVersion, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssProjectID, fdhkDssProjectID, _ := fds.GetOK("dss_project_id")
	if err := o.bindDssProjectID(fdDssProjectID, fdhkDssProjectID, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssProjectName, fdhkDssProjectName, _ := fds.GetOK("dss_project_name")
	if err := o.bindDssProjectName(fdDssProjectName, fdhkDssProjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssWorkspaceID, fdhkDssWorkspaceID, _ := fds.GetOK("dss_workspace_id")
	if err := o.bindDssWorkspaceID(fdDssWorkspaceID, fdhkDssWorkspaceID, route.Formats); err != nil {
		res = append(res, err)
	}

	fdDssWorkspaceName, fdhkDssWorkspaceName, _ := fds.GetOK("dss_workspace_name")
	if err := o.bindDssWorkspaceName(fdDssWorkspaceName, fdhkDssWorkspaceName, route.Formats); err != nil {
		res = append(res, err)
	}

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

	fdGroupName, fdhkGroupName, _ := fds.GetOK("group_name")
	if err := o.bindGroupName(fdGroupName, fdhkGroupName, route.Formats); err != nil {
		res = append(res, err)
	}

	fdName, fdhkName, _ := fds.GetOK("name")
	if err := o.bindName(fdName, fdhkName, route.Formats); err != nil {
		res = append(res, err)
	}

	fdOriginSystem, fdhkOriginSystem, _ := fds.GetOK("origin_system")
	if err := o.bindOriginSystem(fdOriginSystem, fdhkOriginSystem, route.Formats); err != nil {
		res = append(res, err)
	}

	fdTags, fdhkTags, _ := fds.GetOK("tags")
	if err := o.bindTags(fdTags, fdhkTags, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindClusterType binds and validates parameter ClusterType from formData.
func (o *CreateExperimentByUploadParams) bindClusterType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCreateExperimentByUploadParams()
		return nil
	}

	o.ClusterType = &raw

	if err := o.validateClusterType(formats); err != nil {
		return err
	}

	return nil
}

// validateClusterType carries on validations for parameter ClusterType
func (o *CreateExperimentByUploadParams) validateClusterType(formats strfmt.Registry) error {

	if err := validate.Enum("cluster_type", "formData", *o.ClusterType, []interface{}{"BDAP", "BDP"}); err != nil {
		return err
	}

	return nil
}

// bindDescription binds and validates parameter Description from formData.
func (o *CreateExperimentByUploadParams) bindDescription(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Description = &raw

	return nil
}

// bindDssFlowID binds and validates parameter DssFlowID from formData.
func (o *CreateExperimentByUploadParams) bindDssFlowID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssFlowID = &raw

	return nil
}

// bindDssFlowName binds and validates parameter DssFlowName from formData.
func (o *CreateExperimentByUploadParams) bindDssFlowName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssFlowName = &raw

	return nil
}

// bindDssFlowVersion binds and validates parameter DssFlowVersion from formData.
func (o *CreateExperimentByUploadParams) bindDssFlowVersion(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssFlowVersion = &raw

	return nil
}

// bindDssProjectID binds and validates parameter DssProjectID from formData.
func (o *CreateExperimentByUploadParams) bindDssProjectID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssProjectID = &raw

	return nil
}

// bindDssProjectName binds and validates parameter DssProjectName from formData.
func (o *CreateExperimentByUploadParams) bindDssProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssProjectName = &raw

	return nil
}

// bindDssWorkspaceID binds and validates parameter DssWorkspaceID from formData.
func (o *CreateExperimentByUploadParams) bindDssWorkspaceID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssWorkspaceID = &raw

	return nil
}

// bindDssWorkspaceName binds and validates parameter DssWorkspaceName from formData.
func (o *CreateExperimentByUploadParams) bindDssWorkspaceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.DssWorkspaceName = &raw

	return nil
}

// bindFile binds file parameter File.
//
// The only supported validations on files are MinLength and MaxLength
func (o *CreateExperimentByUploadParams) bindFile(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

// bindFileName binds and validates parameter FileName from formData.
func (o *CreateExperimentByUploadParams) bindFileName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.FileName = &raw

	return nil
}

// bindGroupName binds and validates parameter GroupName from formData.
func (o *CreateExperimentByUploadParams) bindGroupName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("group_name", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("group_name", "formData", raw); err != nil {
		return err
	}

	o.GroupName = raw

	if err := o.validateGroupName(formats); err != nil {
		return err
	}

	return nil
}

// validateGroupName carries on validations for parameter GroupName
func (o *CreateExperimentByUploadParams) validateGroupName(formats strfmt.Registry) error {

	if err := validate.MinLength("group_name", "formData", o.GroupName, 1); err != nil {
		return err
	}

	return nil
}

// bindName binds and validates parameter Name from formData.
func (o *CreateExperimentByUploadParams) bindName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("name", "formData")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("name", "formData", raw); err != nil {
		return err
	}

	o.Name = raw

	if err := o.validateName(formats); err != nil {
		return err
	}

	return nil
}

// validateName carries on validations for parameter Name
func (o *CreateExperimentByUploadParams) validateName(formats strfmt.Registry) error {

	if err := validate.MinLength("name", "formData", o.Name, 1); err != nil {
		return err
	}

	return nil
}

// bindOriginSystem binds and validates parameter OriginSystem from formData.
func (o *CreateExperimentByUploadParams) bindOriginSystem(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewCreateExperimentByUploadParams()
		return nil
	}

	o.OriginSystem = &raw

	if err := o.validateOriginSystem(formats); err != nil {
		return err
	}

	return nil
}

// validateOriginSystem carries on validations for parameter OriginSystem
func (o *CreateExperimentByUploadParams) validateOriginSystem(formats strfmt.Registry) error {

	if err := validate.Enum("origin_system", "formData", *o.OriginSystem, []interface{}{"MLSS", "DSS", "WTSS"}); err != nil {
		return err
	}

	return nil
}

// bindTags binds and validates array parameter Tags from formData.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *CreateExperimentByUploadParams) bindTags(rawData []string, hasKey bool, formats strfmt.Registry) error {

	var qvTags string
	if len(rawData) > 0 {
		qvTags = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	tagsIC := swag.SplitByFormat(qvTags, "")
	if len(tagsIC) == 0 {
		return nil
	}

	var tagsIR []string
	for _, tagsIV := range tagsIC {
		tagsI := tagsIV

		tagsIR = append(tagsIR, tagsI)
	}

	o.Tags = tagsIR

	return nil
}