package repo

import (
	datasource "webank/DI/pkg/datasource/mysql"
	gormmodels "webank/DI/pkg/v2/model/third"
	"webank/DI/pkg/v2/repo"
)

var UserRoleRepo UserRoleIF

const TableNameUserRole = "t_user_role_v2"

type UserRoleIF interface {
	GetGroupIdsByUserId(userId string) ([]string, error)
	IsSA(userId string) bool
}

type UserRoleRepoImpl struct {
	repo.BaseRepoImpl
}

func init() {
	UserRoleRepo = &UserRoleRepoImpl{
		repo.BaseRepoImpl{
			TableName: TableNameUserRole,
		},
	}
}

func (u *UserRoleRepoImpl) GetGroupIdsByUserId(userId string) ([]string, error) {
	var userRoles []gormmodels.UserRole
	db := datasource.GetDB()
	err := db.Table(u.TableName).Find(&userRoles,
		"user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, user := range userRoles {
		ids = append(ids, user.GroupID)
	}
	return ids, err
}

func (u *UserRoleRepoImpl) IsSA(userId string) bool {
	var userRole gormmodels.UserRole
	db := datasource.GetDB()
	_ = db.Table(u.TableName).Find(&userRole,
		"user_id = ? and (role_name = 'ADMIN' or role_name = 'PLATFORM_OPS') ", userId).Error
	isSa := false
	if userRole.UserID != "" {
		isSa = true
	}
	return isSa
}
