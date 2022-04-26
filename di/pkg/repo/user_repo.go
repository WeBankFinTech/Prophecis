package repo

import (
	"gorm.io/gorm"
	"webank/DI/commons/models"
)

const TableNameUser = "t_user"

var UserRepo UserRepoIF

type UserRepoIF interface {
	Get(username *string, db *gorm.DB) (*models.User, error)
}

type UserRepoImpl struct {
	BaseRepoImpl
}

func (i UserRepoImpl) Get(username *string, db *gorm.DB) (*models.User, error) {
	var user models.User
	err := db.Table(i.TableName).Find(&user, "name = ? AND enable_flag = ?", username, 1).Error
	return &user, err
}

func init() {
	UserRepo = &UserRepoImpl{
		BaseRepoImpl{TableNameUser},
	}
}

