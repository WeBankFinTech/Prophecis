package repo

import (
	"fmt"
	"github.com/prometheus/common/log"
	"testing"
	"webank/DI/pkg/datasource/mysql"
)

func TestCreateExp(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	//s := "1"
	//exp := models.Experiment{
	//	//ExpID: 1,
	//	ExpDesc: &s,
	//	ExpName: &s,
	//}
	//err := Add(&exp, db)
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//s := "12355577"
	//name := "asd555555555555555"
	//bool := false
	//exp := models.Experiment{
	//
	//	ExpID:   1,
	//	ExpDesc: &s,
	//	ExpName: &name,
	//	BaseModel: models.BaseModel{
	//		ID:         4,
	//		EnableFlag: &bool,
	//	},
	//}
	//err := UpdateExperiment(&exp, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//err := DeleteExperimentById(1, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//}

	//experiment, err := GetExperimentByExperimentId(3, db.Debug())
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//} else {
	//	fmt.Printf("experiment: %v", experiment)
	//}

	//offset := GetAllExperimentsByOffset(0, 99, db.Debug())
	//fmt.Printf("experiment: %v", offset)

	//experiment := GetAllExperiments(db.Debug())
	//fmt.Printf("experiment: %v", experiment)
	user := "alexwu"

	experiment,err := ExperimentRepo.CountByUser("",user, db.Debug())
	if err != nil{
		log.Error(err.Error())
	}
	fmt.Printf("experiment: %v", experiment)
}


func TestGetAllExperiments(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()
	experiments, err := ExperimentRepo.GetAllByOffset(0, 100, "", db.Debug())
	if err != nil{
		println(err.Error())
	}
	println(len(experiments))
}

func TestGetExperiments(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()
	experiments, err := ExperimentRepo.Get(9,db)
	if err != nil{
		println(err.Error())
	}
	println(experiments)
}