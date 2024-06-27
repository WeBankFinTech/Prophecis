package third

import "time"

type UserRole struct {
	UserRoleID string    `gorm:"column:user_role_id; primaryKey" json:"user_role_id"`
	UserID     string    `json:"user_id"`
	RoleName   string    `json:"role_name"`
	RoleType   string    `json:"role_type"`
	GroupID    string    `json:"group_id"`
	CreateUser string    `json:"create_user"`
	CreateTime time.Time `json:"create_time"`
	UpdateUser string    `json:"update_user"`
	UpdateTime time.Time `json:"update_time"`
}

func (UserRole) TableName() string {
	return "user_role_v2"
}
