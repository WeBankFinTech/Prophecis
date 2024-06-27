package repo

import (
	datasource "webank/DI/pkg/datasource/mysql"
	gormmodels "webank/DI/pkg/v2/model/third"
	"webank/DI/pkg/v2/repo"
)

const TableNameModelVersion = "t_model_version_v2"
const EnableFlagTrue = 1

var ModelVersionRepo ModelVersionRepoIF

type ModelVersionRepoIF interface {
	GetById(modelVersionId string) (*gormmodels.ModelVersion, error)
}

type ModelVersionRepoImpl struct {
	repo.BaseRepoImpl
}

func init() {
	ModelVersionRepo = &ModelVersionRepoImpl{
		repo.BaseRepoImpl{
			TableName: TableNameModelVersion,
		},
	}
}

func (m *ModelVersionRepoImpl) GetById(mvId string) (*gormmodels.ModelVersion, error) {
	db := datasource.GetDB().Debug()
	var mv gormmodels.ModelVersion

	err := db.Table(m.TableName).
		Where("model_version_id = ?", mvId).
		Where("enable_flag = ?", EnableFlagTrue).
		Preload("Model").
		Find(&mv).Error
	return &mv, err
}
