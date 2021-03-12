package repo

import (
	"gorm.io/gorm"
	"webank/DI/commons/models"
)

const TableNameExperimentTag = "t_experiment_tag"

var ExperimentTagRepo ExperimentTagRepoIF

type ExperimentTagRepoIF interface {
	//-------------base-----------
	Add(model *models.ExperimentTag, db *gorm.DB) error
	BatchAdd(model []*models.ExperimentTag, db *gorm.DB) error
	BatchDelete(id int64, db *gorm.DB) error
	Update(model *models.ExperimentTag, db *gorm.DB) error
	Delete(id int64, db *gorm.DB) error
	Get(id int64, db *gorm.DB) (*models.ExperimentTag, error)
	GetAllByOffset(offSet int64, size int64, db *gorm.DB) ([]*models.ExperimentTag, error)
	GetAllByUserOffset(offSet int64, size int64, user string, db *gorm.DB) ([]*models.ExperimentTag, error)
	GetAll(db *gorm.DB) ([]*models.ExperimentTag, error)
	Count(db *gorm.DB) (int64, error)
	CountByUser(user string, db *gorm.DB) (int64, error)
	GetByExpId(expId int64, db *gorm.DB) ([]*models.ExperimentTag, error)
}

type ExperimentTagRepoImpl struct {
	BaseRepoImpl
}

func init() {
	ExperimentTagRepo = &ExperimentTagRepoImpl{
		BaseRepoImpl{TableNameExperimentTag},
	}
}


func (i ExperimentTagRepoImpl) Add(model *models.ExperimentTag, db *gorm.DB) error {
	return db.Table(i.TableName).Create(model).Error
}

func (i ExperimentTagRepoImpl) BatchAdd(models []*models.ExperimentTag, db *gorm.DB) error {
	return db.Table(i.TableName).Create(&models).Error
}

func (i ExperimentTagRepoImpl) Update(model *models.ExperimentTag, db *gorm.DB) error {
	return db.Table(i.TableName).Model(model).Where("id = ?", model.ID).Updates(model).Error
}

func (i ExperimentTagRepoImpl) Delete(id int64, db *gorm.DB) error {
	return db.Table(i.TableName).Where("id = ?", id).Update("enable_flag", 0).Error
}

func (i ExperimentTagRepoImpl) BatchDelete(expId int64, db *gorm.DB) error {
	return db.Table(i.TableName).Where("exp_id = ?", expId).Update("enable_flag", 0).Error
}

func (i ExperimentTagRepoImpl) Get(id int64, db *gorm.DB) (*models.ExperimentTag, error) {
	var result models.ExperimentTag
	err := db.Table(i.TableName).Find(&result, "id = ? AND enable_flag = ?", id, 1).Error
	return &result, err
}

func (i ExperimentTagRepoImpl) GetByExpId(expId int64, db *gorm.DB) ([]*models.ExperimentTag, error) {
	var result []*models.ExperimentTag
	err := db.Where("exp_id = ?", expId).Find(&result).Error
	return result, err
}

func (i ExperimentTagRepoImpl) GetAllByOffset(offSet int64, size int64, db *gorm.DB) ([]*models.ExperimentTag, error) {
	var result []*models.ExperimentTag
	err := db.Offset(int(offSet)).Limit(int(size)).Table(i.TableName).Find(&result, "enable_flag = ?", 1).Error
	return result, err
}

func (i ExperimentTagRepoImpl) GetAllByUserOffset(offSet int64, size int64, user string, db *gorm.DB) ([]*models.ExperimentTag, error) {
	var result []*models.ExperimentTag
	err := db.Offset(int(offSet)).Limit(int(size)).Table(i.TableName).
		Joins("LEFT JOIN t_user on "+i.TableName+".exp_run_create_user_id = t_user.id").
		Find(&result, "t_user.name = ? and enable_flag = ?", user, 1).Error
	return result, err
}

func (i ExperimentTagRepoImpl) GetAll(db *gorm.DB) ([]*models.ExperimentTag, error) {
	var result []*models.ExperimentTag
	err := db.Table(i.TableName).Find(&result, "enable_flag = ?", 1).Error
	return result, err
}

func (i ExperimentTagRepoImpl) Count(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Table(i.TableName).Where("enable_flag = ?", 1).Count(&count).Error
	return count, err
}

func (i ExperimentTagRepoImpl) CountByUser(user string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Table(i.TableName).Joins("LEFT JOIN t_user on "+i.TableName+".exp_run_create_user_id = t_user.id").
		Where("t_user.name = ? and enable_flag = ?", user, 1, 1).Count(&count).Error
	return count, err
}


