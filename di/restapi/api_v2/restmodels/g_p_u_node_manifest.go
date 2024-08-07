// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// GPUNodeManifest g p u node manifest
// swagger:model GPUNodeManifest
type GPUNodeManifest struct {

	// [调用接口获取数据] 参数设置-任务执行设置-执行代码设置 为codeFile时，调用 codeUpload 接口返回的s3Path
	CodePath string `json:"code_path,omitempty"`

	// 参数设置-任务执行设置-执行代码设置（storagePath对应训练代码子目录；codeFile对应手动上传）
	// Enum: [storagePath codeFile]
	CodeSelector string `json:"code_selector,omitempty"`

	// 资源设置-计算资源设置-计算任务CPU
	Cpus int64 `json:"cpus,omitempty"`

	// 资源设置-目录设置
	DataStores *GPUNodeDataStores `json:"data_stores,omitempty"`

	// 参数设置-任务执行设置-高级参数设置
	DependResources *DependResources `json:"depend_resources,omitempty"`

	// 基本信息-节点描述
	Description string `json:"description,omitempty"`

	// evaluation metrics
	EvaluationMetrics *GPUNodeManifestEvaluationMetrics `json:"evaluation_metrics,omitempty"`

	// 该节点所属工作流所属的MLSS实验的id
	ExpID string `json:"exp_id,omitempty"`

	// 该节点所属工作流所属的MLSS实验的name
	ExpName string `json:"exp_name,omitempty"`

	// [调用接口获取数据] 参数设置-任务执行设置-执行代码设置 为codeFile时，上传文件的文件名
	FileName string `json:"fileName,omitempty"`

	// 资源设置-镜像设置 和 参数设置-执行入口 相关信息
	Framework *FrameWork `json:"framework,omitempty"`

	// 资源设置-计算资源设置-计算任务GPU
	Gpus int64 `json:"gpus,omitempty"`

	// 参数设置-任务执行设置-高级参数设置-忽略失败状态（1表示是，0表示否）
	IgnoreFail int64 `json:"ignore_fail,omitempty"`

	// 资源设置-计算资源设置-任务类型（Local表示单机; dist-tf表示分布式）
	// Enum: [Local dist-tf]
	JobType string `json:"job_type,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-计算任务数
	Learners int64 `json:"learners,omitempty"`

	// 资源设置-计算资源设置-计算任务内存
	Memory string `json:"memory,omitempty"`

	// 基本信息-节点名
	Name string `json:"name,omitempty"`

	// [下拉接口获取数据] 资源设置-计算资源设置-命名空间
	Namespace string `json:"namespace,omitempty"`

	// [下拉接口获取数据] 基本信息-代理用户设置开启-代理用户设置
	ProxyUser string `json:"proxy_user,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-参数服务器CPU
	PsCPU int64 `json:"ps_cpu,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-参数服务器镜像
	PsImage int64 `json:"ps_image,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-参数服务器镜像类型
	// Enum: [Standard Custom]
	PsImageType string `json:"ps_imageType,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-参数服务器内存
	PsMemory string `json:"ps_memory,omitempty"`

	// 当job_type为dist_tf时，资源设置-计算资源设置-参数服务器数
	Pss int64 `json:"pss,omitempty"`

	// 参数设置-任务执行设置-高级参数设置-python环境配置(其value来自哪里todo)
	PythonHdfsPath string `json:"python_hdfs_path,omitempty"`

	// 目前固定为 1.0
	Version string `json:"version,omitempty"`
}

// Validate validates this g p u node manifest
func (m *GPUNodeManifest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCodeSelector(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDataStores(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDependResources(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEvaluationMetrics(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFramework(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateJobType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePsImageType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var gPUNodeManifestTypeCodeSelectorPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["storagePath","codeFile"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		gPUNodeManifestTypeCodeSelectorPropEnum = append(gPUNodeManifestTypeCodeSelectorPropEnum, v)
	}
}

const (

	// GPUNodeManifestCodeSelectorStoragePath captures enum value "storagePath"
	GPUNodeManifestCodeSelectorStoragePath string = "storagePath"

	// GPUNodeManifestCodeSelectorCodeFile captures enum value "codeFile"
	GPUNodeManifestCodeSelectorCodeFile string = "codeFile"
)

// prop value enum
func (m *GPUNodeManifest) validateCodeSelectorEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, gPUNodeManifestTypeCodeSelectorPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *GPUNodeManifest) validateCodeSelector(formats strfmt.Registry) error {

	if swag.IsZero(m.CodeSelector) { // not required
		return nil
	}

	// value enum
	if err := m.validateCodeSelectorEnum("code_selector", "body", m.CodeSelector); err != nil {
		return err
	}

	return nil
}

func (m *GPUNodeManifest) validateDataStores(formats strfmt.Registry) error {

	if swag.IsZero(m.DataStores) { // not required
		return nil
	}

	if m.DataStores != nil {
		if err := m.DataStores.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data_stores")
			}
			return err
		}
	}

	return nil
}

func (m *GPUNodeManifest) validateDependResources(formats strfmt.Registry) error {

	if swag.IsZero(m.DependResources) { // not required
		return nil
	}

	if m.DependResources != nil {
		if err := m.DependResources.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("depend_resources")
			}
			return err
		}
	}

	return nil
}

func (m *GPUNodeManifest) validateEvaluationMetrics(formats strfmt.Registry) error {

	if swag.IsZero(m.EvaluationMetrics) { // not required
		return nil
	}

	if m.EvaluationMetrics != nil {
		if err := m.EvaluationMetrics.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("evaluation_metrics")
			}
			return err
		}
	}

	return nil
}

func (m *GPUNodeManifest) validateFramework(formats strfmt.Registry) error {

	if swag.IsZero(m.Framework) { // not required
		return nil
	}

	if m.Framework != nil {
		if err := m.Framework.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("framework")
			}
			return err
		}
	}

	return nil
}

var gPUNodeManifestTypeJobTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Local","dist-tf"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		gPUNodeManifestTypeJobTypePropEnum = append(gPUNodeManifestTypeJobTypePropEnum, v)
	}
}

const (

	// GPUNodeManifestJobTypeLocal captures enum value "Local"
	GPUNodeManifestJobTypeLocal string = "Local"

	// GPUNodeManifestJobTypeDistTf captures enum value "dist-tf"
	GPUNodeManifestJobTypeDistTf string = "dist-tf"
)

// prop value enum
func (m *GPUNodeManifest) validateJobTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, gPUNodeManifestTypeJobTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *GPUNodeManifest) validateJobType(formats strfmt.Registry) error {

	if swag.IsZero(m.JobType) { // not required
		return nil
	}

	// value enum
	if err := m.validateJobTypeEnum("job_type", "body", m.JobType); err != nil {
		return err
	}

	return nil
}

var gPUNodeManifestTypePsImageTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Standard","Custom"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		gPUNodeManifestTypePsImageTypePropEnum = append(gPUNodeManifestTypePsImageTypePropEnum, v)
	}
}

const (

	// GPUNodeManifestPsImageTypeStandard captures enum value "Standard"
	GPUNodeManifestPsImageTypeStandard string = "Standard"

	// GPUNodeManifestPsImageTypeCustom captures enum value "Custom"
	GPUNodeManifestPsImageTypeCustom string = "Custom"
)

// prop value enum
func (m *GPUNodeManifest) validatePsImageTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, gPUNodeManifestTypePsImageTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *GPUNodeManifest) validatePsImageType(formats strfmt.Registry) error {

	if swag.IsZero(m.PsImageType) { // not required
		return nil
	}

	// value enum
	if err := m.validatePsImageTypeEnum("ps_imageType", "body", m.PsImageType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GPUNodeManifest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GPUNodeManifest) UnmarshalBinary(b []byte) error {
	var res GPUNodeManifest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// GPUNodeManifestEvaluationMetrics g p u node manifest evaluation metrics
// swagger:model GPUNodeManifestEvaluationMetrics
type GPUNodeManifestEvaluationMetrics struct {

	// 目前固定为fluent-bit
	Type string `json:"type,omitempty"`
}

// Validate validates this g p u node manifest evaluation metrics
func (m *GPUNodeManifestEvaluationMetrics) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GPUNodeManifestEvaluationMetrics) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GPUNodeManifestEvaluationMetrics) UnmarshalBinary(b []byte) error {
	var res GPUNodeManifestEvaluationMetrics
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
