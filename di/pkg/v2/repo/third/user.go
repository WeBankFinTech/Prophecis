package repo

import (
	datasource "webank/DI/pkg/datasource/mysql"
	gormmodels "webank/DI/pkg/v2/model/third"
	"webank/DI/pkg/v2/repo"
)

const TableNameUser = "t_user_v2"

var UserRepo UserIF

type UserIF interface {
	GetUser(userName string, clusterType string) (*gormmodels.UserV2, error)
}

type UserRepoImpl struct {
	repo.BaseRepoImpl
}

func init() {
	UserRepo = &UserRepoImpl{
		repo.BaseRepoImpl{
			TableName: TableNameUser,
		},
	}
}

func (u *UserRepoImpl) GetUser(userName string, clusterType string) (*gormmodels.UserV2, error) {
	var user gormmodels.UserV2
	db := datasource.GetDB()
	db = db.Table(u.TableName)
	err := db.Where("user_name = ? AND mlss_cluster_name = ?", userName, clusterType).First(&user).Error
	return &user, err
}
