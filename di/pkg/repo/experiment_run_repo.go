package repo

import (
	"gorm.io/gorm"
	"webank/DI/commons/models"
)

const TableNameExperimentRun = "t_experiment_run"

var ExperimentRunRepo ExperimentRunRepoIF

type ExperimentRunRepoIF interface {
	//-------------base-----------
	Add(exp *models.ExperimentRun, db *gorm.DB) error
	Update(model *models.ExperimentRun, db *gorm.DB) error
	Delete(id int64, db *gorm.DB) error
	Get(id int64, db *gorm.DB) (*models.ExperimentRun, error)
	GetAllByOffset(offSet int64, size int64, query string, db *gorm.DB) ([]*models.ExperimentRun, error)
	GetAllByUserOffset(offSet int64, size int64, query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error)
	GetAll(db *gorm.DB) ([]*models.ExperimentRun, error)
	GetAllByUser(query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error)
	Count(query string, db *gorm.DB) (int64, error)
	CountExpId(expId int64, query string, db *gorm.DB) (int64, error)
	CountByUser(query string, user string, db *gorm.DB) (int64, error)
	CountByExpIdAndUser(expId int64, query string, user string, db *gorm.DB) (int64, error)
	UpdateStatus(expRunId int64, status string, db *gorm.DB) error
	GetByUserAndExpIdAndOffset(expId int64, page int64, size int64, query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error)
	GetByExpIdAndOffset(expId int64, page int64, size int64, query string, db *gorm.DB) ([]*models.ExperimentRun, error)
	GetByExecIdAndTaskId(taskId int64, db *gorm.DB) (int64, error)

	//-------------base-----------
}
type ExperimentRunRepoImpl struct {
	BaseRepoImpl
}

func (r *ExperimentRunRepoImpl) UpdateStatus(expRunId int64, status string, db *gorm.DB) error {
	return db.Table(r.TableName).Where("id = ?", expRunId).Update("exp_exec_status", status).Error
}

func init() {
	ExperimentRunRepo = &ExperimentRunRepoImpl{
		BaseRepoImpl{
			TableName: TableNameExperimentRun,
		},
	}
}

func (r *ExperimentRunRepoImpl) Add(model *models.ExperimentRun, db *gorm.DB) error {
	e := db.Table(r.TableName).Create(model).Error
	return e
}
func (r *ExperimentRunRepoImpl) Update(model *models.ExperimentRun, db *gorm.DB) error {
	return db.Table(r.TableName).Model(model).Where("id = ?", model.ID).Updates(model).Error
}

func (r *ExperimentRunRepoImpl) Delete(id int64, db *gorm.DB) error {
	return db.Table(r.TableName).Where("id = ?", id).Update("enable_flag", 0).Error
}

func (r *ExperimentRunRepoImpl) Get(id int64, db *gorm.DB) (*models.ExperimentRun, error) {
	var result models.ExperimentRun
	err := db.Table(r.TableName).Preload("Experiment").Preload("User").Preload("User").Find(&result, "id = ? AND enable_flag = ?", id, 1).Error
	return &result, err
}

func (r *ExperimentRunRepoImpl) GetByExecIdAndTaskId(taskId int64, db *gorm.DB) (int64, error) {
	var result models.ExperimentRun
	err := db.Where("dss_task_id = ?", taskId).Find(&result).Error
	return result.ID, err
}

func (r *ExperimentRunRepoImpl) GetAllByOffset(offSet int64, size int64, query string, db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Where( "CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND "+r.TableName + ".enable_flag = ? ","%"+query+"%",1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Order("exp_run_create_time DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		Preload("User").
		Preload("Experiment").
		Find(&result).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) GetAllByUserOffset(offSet int64, size int64, query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND t_user.name = ? " +
		"AND "+r.TableName + ".enable_flag = ?","%"+query+"%",user,1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Order("exp_run_create_time DESC").
		Offset(int(offSet)).
		Limit(int(size)).
		Preload("User").
		Preload("Experiment").
		Find(&result).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) GetAll(db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Order("exp_run_create_time desc").Table(r.TableName).Preload("Experiment").Preload("User").Find(&result, "enable_flag = ?", 1).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) GetAllByUser(query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"t_user.name = ? AND "+r.TableName+"."+"enable_flag = ?","%"+query+"%", user, 1).
		Order("exp_run_create_time desc").
		Preload("Experiment").
		Preload("User").
		Joins("LEFT JOIN t_user ON "+r.TableName+".exp_run_create_user_id = t_user.id").
		Find(&result).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) GetByUserAndExpIdAndOffset(expId int64, page int64, size int64, query string, user string, db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND t_user.name = ? " +
		"AND "+r.TableName + ".enable_flag = ? " +
		"AND "+r.TableName + ".exp_id = ?","%"+query+"%",user,1,expId).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Preload("Experiment").
		Preload("User").
		Offset(int(page)).
		Limit(int(size)).
		Order("exp_run_create_time DESC").
		Find(&result).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) GetByExpIdAndOffset(expId int64, page int64, size int64, query string, db *gorm.DB) ([]*models.ExperimentRun, error) {
	var result []*models.ExperimentRun
	err := db.Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND "+r.TableName + ".enable_flag = ? " +
		"AND "+r.TableName + ".exp_id = ?","%"+query+"%",1,expId).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Preload("Experiment").
		Preload("User").
		Offset(int(page)).
		Limit(int(size)).
		Order("exp_run_create_time DESC").
		Find(&result).Error
	return result, err
}

func (r *ExperimentRunRepoImpl) Count(query string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.ExperimentRun{}).Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND " + r.TableName+".enable_flag = ? ","%"+query+"%",1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Count(&count).Error
	return count, err
}

func (r *ExperimentRunRepoImpl) CountExpId(expId int64, query string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.ExperimentRun{}).Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND " + r.TableName+".enable_flag = ? " +
		"AND " + r.TableName+".exp_id = ?","%"+query+"%",1, expId).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Count(&count).Error
	return count, err
}

func (r *ExperimentRunRepoImpl) CountByUser(query string, user string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.ExperimentRun{}).Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND t_user.name = ? " +
		"AND "+r.TableName + ".enable_flag = ?" ,"%"+query+"%",user,1).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Count(&count).Error
	return count, err
}

func (r *ExperimentRunRepoImpl) CountByExpIdAndUser(expId int64, query string, user string, db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&models.ExperimentRun{}).Where("CONCAT(t_experiment.exp_name,t_user.name,"+r.TableName+".exp_exec_type) LIKE ? " +
		"AND t_user.name = ? " +
		"AND "+r.TableName + ".enable_flag = ? " +
		"AND "+r.TableName + ".exp_id = ? " ,"%"+query+"%",user,1, expId).
		Joins("LEFT JOIN t_user " +
			"ON "+r.TableName+".exp_run_create_user_id = t_user.id " +
			"LEFT JOIN t_experiment " +
			"ON t_experiment.id = "+r.TableName+".exp_id ").
		Count(&count).Error
	return count, err
}
