// Code generated by go-swagger; DO NOT EDIT.

package restmodels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// EventSenderNodeContent event sender node content
// swagger:model EventSenderNodeContent
type EventSenderNodeContent struct {

	// 现在固定为FA0?
	Dcn string `json:"dcn,omitempty"`

	// 信号发送节点所属的工作流所属的实验的id
	ExpID string `json:"exp_id,omitempty"`

	// 信号发送节点所属的工作流所属的实验的名称
	ExpName string `json:"exp_name,omitempty"`

	// 属性信息-项目组的id
	GroupID string `json:"group_id,omitempty"`

	// 属性信息-项目组的名称
	GroupName string `json:"group_name,omitempty"`

	// 属性信息-msg.body
	MsgBody string `json:"msg_body,omitempty"`

	// 属性信息-msg.messageType
	MsgType string `json:"msg_type,omitempty"`

	// 信号发送节点的名称
	Name string `json:"name,omitempty"`

	// 现在固定为 11200257 ?
	ServiceID string `json:"service_id,omitempty"`
}

// Validate validates this event sender node content
func (m *EventSenderNodeContent) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *EventSenderNodeContent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EventSenderNodeContent) UnmarshalBinary(b []byte) error {
	var res EventSenderNodeContent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
