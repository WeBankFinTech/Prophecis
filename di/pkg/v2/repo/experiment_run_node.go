package repo

import (
	datasource "webank/DI/pkg/datasource/mysql"
	models "webank/DI/pkg/v2/model"
)

const TableNameExperimentRunNode = "t_experiment_run_v2"

var ExperimentRunNodeRepo ExperimentRunNodeRepoIF

type ExperimentRunNodeRepoIF interface {
	Add(experimentRunNode *models.ExperimentRunNode) error
	UpdateByMap1(expRunId string, key string, updateMap map[string]interface{}) error
	GetAllByExpRunID(expRunId string) ([]*models.ExperimentRunNode, error)
}

type ExperimentRunNodeRepoImpl struct {
	BaseRepoImpl
}

func init() {
	ExperimentRunNodeRepo = &ExperimentRunNodeRepoImpl{
		BaseRepoImpl{
			TableName: TableNameExperimentRunNode,
		},
	}
}

func (r *ExperimentRunNodeRepoImpl) GetAllByExpRunID(expRunId string) ([]*models.ExperimentRunNode, error) {
	db := datasource.GetDB()
	var expRunNodes []*models.ExperimentRunNode
	err := db.Table(r.TableName).Where("experiment_run_id = ?", expRunId).Find(&expRunNodes).Error
	return expRunNodes, err
}

func (r *ExperimentRunNodeRepoImpl) UpdateByMap1(expRunId string, key string, updateMap map[string]interface{}) error {
	db := datasource.GetDB()
	return db.Table(r.TableName).Where("experiment_run_id = ? AND key = ?", expRunId, key).Updates(updateMap).Error
}

func (r *ExperimentRunNodeRepoImpl) Add(experimentRunNode *models.ExperimentRunNode) error {
	db := datasource.GetDB()
	return db.Table(r.TableName).Create(experimentRunNode).Error
}
