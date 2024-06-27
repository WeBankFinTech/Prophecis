package repo

import (
	datasource "webank/DI/pkg/datasource/mysql"
	gormmodels "webank/DI/pkg/v2/model/third"
	"webank/DI/pkg/v2/repo"
)

const TableNameProcesslineVersion = "t_processline_version_v2"

var ProcesslineVersionRepo ProcessLineVersionRepoIF

type ProcessLineVersionRepoIF interface {
	GetById(modelVersionId string) (*gormmodels.ProcessLineVersion, error)
}

type ProcessLineVersionRepoImpl struct {
	repo.BaseRepoImpl
}

func init() {
	ProcesslineVersionRepo = &ProcessLineVersionRepoImpl{
		repo.BaseRepoImpl{
			TableName: TableNameProcesslineVersion,
		},
	}
}

func (p *ProcessLineVersionRepoImpl) GetById(processLineVersionId string) (*gormmodels.ProcessLineVersion, error) {
	tx := datasource.GetDB().Debug()
	processlineVersion := gormmodels.ProcessLineVersion{}
	err := tx.Table(p.TableName).Where("processline_version_id =? and enable_flag =1", processLineVersionId).Find(&processlineVersion).Error
	return &processlineVersion, err
}
