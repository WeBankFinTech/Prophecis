package repo

import (
	log "github.com/sirupsen/logrus"
	"time"
	datasource "webank/DI/pkg/datasource/mysql"
	models "webank/DI/pkg/v2/model"
)

const TableNameExperimentRun = "t_experiment_run_v2"

var ExperimentRunRepo ExperimentRunRepoIF

type ExperimentRunRepoIF interface {
	Add(experimentRun *models.ExperimentRun) error
	Get(expRunId string) (*models.ExperimentRun, error)
	UpdateByMap(expRunId string, updateMap map[string]interface{}) error
	GetAllByOffset(expId *string, expVersionName *string, status []string,
		createTimeSt *string, createTimeEd *string,
		size int, offset int) ([]*models.ExperimentRun, int64, error)
}

type ExperimentRunRepoImpl struct {
	BaseRepoImpl
}

func init() {
	ExperimentRunRepo = &ExperimentRunRepoImpl{
		BaseRepoImpl{
			TableName: TableNameExperimentRun,
		},
	}
}

func (r *ExperimentRunRepoImpl) GetAllByOffset(expId *string, expVersionName *string, status []string,
	createTimeSt *string, createTimeEd *string,
	size int, offset int) ([]*models.ExperimentRun, int64, error) {
	db := datasource.GetDB()

	db = db.Where("deleted = 0")
	if expId != nil {
		db = db.Where("exp_id = ?", *expId)
	}
	// todo(gaoyuanhe); 应该在前面就处理传入为null的情况
	if expVersionName != nil && *expVersionName != "null" {
		db = db.Where("exp_version_name = ?", *expVersionName)
	}
	if len(status) != 0 {
		db = db.Where("execute_status in (?)", status)
	}
	if createTimeSt != nil {
		createTimeStParse, err := time.Parse(time.RFC3339, *createTimeSt)
		if err != nil {
			log.Errorf("time.Parse(%s) failed: %+v", *createTimeSt, err)
		} else {
			db = db.Where("create_time >= ?", createTimeStParse)
		}
	}
	if createTimeEd != nil {
		createTimeEdParse, err := time.Parse(time.RFC3339, *createTimeEd)
		if err != nil {
			log.Errorf("time.Parse(%s) failed: %+v", *createTimeEd, err)
		} else {
			db = db.Where("create_time <= ?", createTimeEdParse)
		}
	}
	var expRuns []*models.ExperimentRun
	var total int64
	err := db.Table(r.TableName).Count(&total).Order("create_time DESC").Offset(offset).Limit(size).Find(&expRuns).Error
	return expRuns, total, err
}

func (r *ExperimentRunRepoImpl) Add(exp *models.ExperimentRun) error {
	db := datasource.GetDB()
	return db.Table(r.TableName).Create(exp).Error
}

func (r *ExperimentRunRepoImpl) Get(expRunId string) (*models.ExperimentRun, error) {
	db := datasource.GetDB()
	var experimentRun models.ExperimentRun
	err := db.Table(r.TableName).Where("id=?", expRunId).Find(&experimentRun).Error
	return &experimentRun, err
}

func (r *ExperimentRunRepoImpl) UpdateByMap(expRunId string, updateMap map[string]interface{}) error {
	db := datasource.GetDB()
	return db.Table(r.TableName).Where("id = ?", expRunId).Updates(updateMap).Error
}
