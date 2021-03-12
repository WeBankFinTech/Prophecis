package models

type User struct {
	BaseModel
	Name       *string `json:"name,omitempty"`
	Gid        *int64  `json:"gid"`
	UID        *int64  `json:"uid"`
	Token      *string `json:"token"`
	Type       *string `json:"type"`
	Remarks    *string `json:"remarks"`
	GuidCheck *int64  `json:"guid_check"`
}
