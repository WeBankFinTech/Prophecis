package service

import (
	"mlss-controlcenter-go/pkg/models"
	"mlss-controlcenter-go/pkg/repo"
)

func GetRolePermissionsByRoleId(roleId int64) ([]*models.RolePermission, error) {
	return repo.GetRolePermissionsByRoleId(roleId)
}