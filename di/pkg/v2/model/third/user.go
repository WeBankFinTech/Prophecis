package third

import "time"

type UserConfiguration struct {
	Token string `json:"token"`
	// group id
	Gid string `json:"gid,omitempty"`

	// is Guid
	IsGUID int64  `json:"isGuid,omitempty"`
	UID    string `json:"uid"`
}

type UserV2 struct {

	// the create time
	CreateTime time.Time `json:"create_time"`

	// the create user name
	CreateUser string `json:"create_user"`

	// is  enable
	EnableFlag int64 `json:"enable_flag"`

	// mlss cluster name
	MlssClusterName string `json:"mlss_cluster_name"`

	// the remark
	Remarks string `json:"remarks"`

	// token
	//Token string `json:"token,omitempty"`
	UserConfiguration string `json:"user_configuration"`
	// uid
	//UID string `json:"uid,omitempty"`

	// the update time
	UpdateTime time.Time `json:"update_time"`
	// update user
	UpdateUser string `json:"update_user"`

	// user cn name
	UserCnName string `json:"user_cn_name"`

	// the id of user
	UserID string `json:"user_id"`

	// the username
	UserName string `json:"user_name"`

	// the user type  system or other
	UserType string `json:"user_type"`
}
