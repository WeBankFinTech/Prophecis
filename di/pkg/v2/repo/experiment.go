package repo

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	datasource "webank/DI/pkg/datasource/mysql"
	models "webank/DI/pkg/v2/model"
)

const TableNameExperiment = "t_experiment_v2"

var ExperimentRepo ExperimentRepoIF

type ExperimentRepoIF interface {
	Add(experiment *models.Experiment) error
	GetByVersion(expId, expVersionName string) (*models.Experiment, error)
	GetLatestVersion(expId string) (*models.Experiment, error)
	UpdateByMap(expId string, versionName string, m map[string]interface{}) error
	UpdateByMapIfEmpty(expId string, versionName string, m map[string]interface{}) error
	ListVersionNames(expId string) ([]string, error)

	GetAllByOffset(currentUserName string, expId *string, name *string, tag *string, status []string,
		groupIDs []string, groupID *string, groupName *string, versionName *string, sourceSystem *string,
		clusterType string, createUser *string, createTimeSt *string, createTimeEd *string,
		updateTimeSt *string, updateTimeEd *string, dssProjectName, dssFlowName *string,
		size int, offset int) ([]*models.Experiment, int64, error)
	GetDSSFieldValues(fieldName string) ([]string, error)
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

func (r *ExperimentRepoImpl) GetDSSFieldValues(fieldName string) ([]string, error) {
	db := datasource.GetDB()
	// 目前的gorm版本没有Distinct函数
	var distinctValues []string
	db = db.Where("source_system = ?", "DSS")
	err := db.Table(r.TableName).Model(&models.Experiment{}).Select(fieldName).Group(fieldName).Pluck(fieldName, &distinctValues).Error
	// 下面的语句 distinctValues 是空的
	// err := db.Raw(fmt.Sprintf("SELECT DISTINCT %s FROM %s", fieldName, r.TableName)).Scan(&distinctValues).Error
	return distinctValues, err
}

func (r *ExperimentRepoImpl) ListVersionNames(expId string) ([]string, error) {
	db := datasource.GetDB()
	versionNameList := make([]string, 0)
	err := db.Table(r.TableName).Where("id=? and status in (?)", expId, []string{"Archive", "Active"}).Order("create_time DESC").Pluck("version_name", &versionNameList).Error
	if err != nil {
		return nil, err
	}
	return versionNameList, nil
}

func (r *ExperimentRepoImpl) GetAllByOffset(currentUserName string, expId *string, name *string, inputTag *string,
	status []string, groupIDs []string, groupID *string, groupName *string, versionName *string, sourceSystem *string,
	clusterType string, createUser *string, createTimeSt *string, createTimeEd *string,
	updateTimeSt *string, updateTimeEd *string, dssProjectName, dssFlowName *string,
	size int, offset int) ([]*models.Experiment, int64, error) {
	db := datasource.GetDB()
	db = db.Where("status != \"Deleted\"")
	// %% 是%的转义表达
	if name != nil {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", *name))
	}
	// from ChatGPT4
	if inputTag != nil {
		tags := strings.Split(*inputTag, ",")
		for _, tag := range tags {
			db = db.Where("FIND_IN_SET(?, tags) > 0", tag)
		}
	}
	if expId != nil {
		db = db.Where("id = ?", *expId)
	}
	if len(status) != 0 {
		db = db.Where("status in (?)", status)
	}
	if versionName != nil {
		db = db.Where("version_name = ?", *versionName)
	}
	if sourceSystem != nil {
		db = db.Where("source_system = ?", *sourceSystem)
	}
	if clusterType != "" {
		db = db.Where("cluster_type = ?", clusterType)
	}
	if len(groupIDs) != 0 {
		db = db.Where("group_id in (?) OR create_user = ?", groupIDs, currentUserName)
	}
	if groupID != nil {
		db = db.Where("group_id = ?", *groupID)
	}
	if groupName != nil {
		db = db.Where("group_name = ?", *groupName)
	}
	if createUser != nil {
		db = db.Where("create_user = ?", createUser)
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
	if updateTimeSt != nil {
		updateTimeStParse, err := time.Parse(time.RFC3339, *updateTimeSt)
		if err != nil {
			log.Errorf("time.Parse(%s) failed: %+v", *updateTimeSt, err)
		} else {
			db = db.Where("update_time >= ?", updateTimeStParse)
		}
	}
	if updateTimeEd != nil {
		updateTimeEdParse, err := time.Parse(time.RFC3339, *updateTimeEd)
		if err != nil {
			log.Errorf("time.Parse(%s) failed: %+v", *updateTimeEd, err)
		} else {
			db = db.Where("update_time <= ?", updateTimeEdParse)
		}
	}
	if dssFlowName != nil {
		db = db.Where("dss_flow_name = ?", *dssFlowName)
	}
	if dssProjectName != nil {
		db = db.Where("dss_project_name = ?", *dssProjectName)
	}
	var exps []*models.Experiment
	var total int64
	err := db.Table(r.TableName).Count(&total).Order("create_time DESC").Offset(offset).Limit(size).Find(&exps).Error
	return exps, total, err
}

func (r *ExperimentRepoImpl) Add(exp *models.Experiment) error {
	db := datasource.GetDB()
	return db.Table(r.TableName).Create(exp).Error
}

func (r *ExperimentRepoImpl) GetByVersion(expId, expVersionName string) (*models.Experiment, error) {
	var exp models.Experiment
	db := datasource.GetDB()
	err := db.Table(r.TableName).Find(&exp, "id = ? AND version_name = ?", expId, expVersionName).Error
	return &exp, err
}

// todo(gaoyuanhe): 这样实现，如果expId不存在，会返回一个空的Experiment，不会返回nil, 该函数的调用者需要检查返回
func (r *ExperimentRepoImpl) GetLatestVersion(expId string) (*models.Experiment, error) {
	db := datasource.GetDB()
	var exp models.Experiment
	err := db.Table(r.TableName).Find(&exp, "id = ? AND status = ?", expId, string(models.ExperimentStatusActive)).Error
	return &exp, err
}

func (r *ExperimentRepoImpl) UpdateByMap(expId string, versionName string, m map[string]interface{}) error {
	db := datasource.GetDB()
	db = db.Table(r.TableName)
	if versionName != "" {
		db = db.Where("version_name = ?", versionName)
	}
	return db.Where("id = ?", expId).Updates(m).Error
}

func (r *ExperimentRepoImpl) UpdateByMapIfEmpty(expId string, versionName string, m map[string]interface{}) error {
	db := datasource.GetDB()
	db = db.Table(r.TableName)
	if versionName != "" {
		db = db.Where("version_name = ?", versionName)
	}
	// 过滤得到那些updateMap key值本来为空的数据, 要除去 update_user, update_time,因为他们不是用来过滤的
	for k, _ := range m {
		if k == "update_user" || k == "update_time" {
			continue
		}
		db = db.Where(fmt.Sprintf("%s = ?", k), "")
	}
	return db.Where("id = ?", expId).Updates(m).Error
}
