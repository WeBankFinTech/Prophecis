package repo

import (
	"gorm.io/gorm"
	"webank/DI/commons/models"
)

const TableNameExperiment = "t_experiment"

var ExperimentRepo ExperimentRepoIF

type ExperimentRepoIF interface {
	//-------------base-----------
	Add(model *models.Experiment, db *gorm.DB) error
	Update(model *models.Experiment, db *gorm.DB) error
	UpdateByMap(id int64, m map[string]interface{}, db *gorm.DB) error
	Delete(id int64, db *gorm.DB) error
	Get(id int64, db *gorm.DB) (*models.Experiment, error)
	GetAllByOffset(offSet int64, size int64, query string, db *gorm.DB) ([]*models.Experiment, error)
	GetAllByUserOffset(offSet int64, size int64, query string, user string, db *gorm.DB) ([]*models.Experiment, error)
	GetAll(db *gorm.DB) ([]*models.Experiment, error)
	Count(query string, db *gorm.DB) (int64, error)
	CountByUser(query string, user string, db *gorm.DB) (int64, error)
	//-------------base-----------
}
type ExperimentRepoImpl struct {
	BaseRepoImpl
}

func init() {
	ExperimentRepo = &ExperimentRepoImpl{
		BaseRepoImpl{
			TableName: TableNameExperiment,
		},
	}
}

func (r *ExperimentRepoImpl) Add(exp *models.Experiment, db *gorm.DB) error {
	e := db.Table(r.TableName).Save(exp).Error
	return e
}

func (r *ExperimentRepoImpl) Update(exp *models.Experiment, db *gorm.DB) error {
	return db.Table(r.TableName).Model(exp).Where("id = ?", exp.ID).Updates(exp).Error
	//return db.Model(exp).Updates(exp).Error
}


func (r *ExperimentRepoImpl) UpdateByMap(id int64,m map[string]interface{}, db *gorm.DB) error {
	return db.Table(r.TableName).Where("id = ?", id).Updates(m).Error
}

func (r *ExperimentRepoImpl) Delete(id int64, db *gorm.DB) error {
	return db.Table(r.TableName).Where("id = ?", id).Update("enable_flag", 0).Error
}

func (r *ExperimentRepoImpl) Get(groupId int64, db *gorm.DB) (*models.Experiment, error) {
	var exp models.Experiment
	err := db.Table(r.TableName).Preload("CreateUser").
		Preload("ModifyUser").Preload("TagList", "enable_flag = ?", 1).
		Find(&exp, "id = ? AND enable_flag = ?", groupId, 1).Error
	return &exp, err
}

func (r *ExperimentRepoImpl) GetAllByOffset(offSet, size int64, query string, db *gorm.DB) ([]*models.Experiment, error) {
	var exps []*models.Experiment
	err := db.Where( "CONCAT(t_user.name,"+r.TableName+".exp_name,"+r.TableName+".exp_desc)  LIKE ? " +
		"AND " + r.TableName + ".enable_flag = ? ","%"+query+"%",1).
		Joins("LEFT JOIN t_user " +
			" ON "+r.TableName+".exp_create_user_id = t_user.id ").
		Order("exp_create_time DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		Preload("CreateUser").
		Preload("TagList", "enable_flag = ?", 1).
		Preload("ModifyUser").
		Find(&exps).Error
	return exps, err
}

func (r *ExperimentRepoImpl) GetAllByUserOffset(offSet, size int64, query, user string, db *gorm.DB) ([]*models.Experiment, error) {
	var exps []*models.Experiment
	err := db.Where("CONCAT(t_user.name,"+r.TableName+".exp_name,"+r.TableName+".exp_desc) LIKE ? " +
		"AND t_user.name = ? AND "+r.TableName+".enable_flag = ? " ,"%"+query+"%",user,1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_create_user_id = t_user.id").
		Order("exp_create_time DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		Preload("CreateUser").
		Preload("TagList", "enable_flag = ?", 1).
		Preload("ModifyUser").
		Find(&exps).Error
	return exps, err
}

func (r *ExperimentRepoImpl) GetAll(db *gorm.DB) ([]*models.Experiment, error) {
	var exps []*models.Experiment
	err := db.Order("exp_create_time desc").Table(r.TableName).Preload("TagList").Preload("User").Find(&exps, "enable_flag = ?", 1).Error
	return exps, err
}

func (r *ExperimentRepoImpl) Count(query string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.Experiment{}).Where("CONCAT(t_user.name,exp_name,exp_desc) LIKE ? " +
		"AND " + r.TableName+".enable_flag = ?","%"+query+"%",1).
		Joins("LEFT JOIN t_user " +
			" ON "+r.TableName+".exp_create_user_id = t_user.id ").
		Count(&count).Error
	return count, err
}
func (r *ExperimentRepoImpl) CountByUser(query string, user string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.Experiment{}).Where("CONCAT(t_user.name,"+r.TableName+".exp_name,"+r.TableName+".exp_desc) LIKE ? " +
		"AND t_user.name = ? " +
		"AND "+r.TableName+".enable_flag = ?","%"+query+"%",user,1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_create_user_id = t_user.id").
		Count(&count).Error
	return count, err
}
