package repo

import (
	"fmt"
	"testing"
	"webank/DI/commons/models"
	"webank/DI/pkg/datasource/mysql"
)

func TestAddExperimentRun(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()
	experimentRun := models.ExperimentRun{
		DssExecID: "buoy-aha",
	}

	err := ExperimentRunRepo.Add(&experimentRun, db)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
	} else {
		fmt.Printf("result: %v", experimentRun)
	}
}

func TestDelete(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	err := ExperimentRunRepo.Delete(2, db)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
	}
}

func TestUpdate(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	run := models.ExperimentRun{
		DssExecID: "aha-1",
		BaseModel: models.BaseModel{
			ID: 1,
		},
	}

	err := ExperimentRunRepo.Update(&run, db)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
	}
}

func TestGet(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	run, err := ExperimentRunRepo.Get(2, db)
	if err != nil {
		fmt.Printf("err: %v", err.Error())
	} else {
		fmt.Printf("result: %v", run)
	}
}

func TestList(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	list, _ := ExperimentRunRepo.GetAllByOffset(0, 10, "", db)
	//if err != nil {
	//	fmt.Printf("err: %v", err.Error())
	//} else {
	fmt.Printf("result: %v", list)
	//}
}

func TestGetAll(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	list, _ := ExperimentRunRepo.GetAll(db)
	fmt.Printf("result: %v", list)

}

func TestCount(t *testing.T) {
	datasource.InitDS()
	db := datasource.GetDB()

	list, _ := ExperimentRunRepo.Count("",db)
	fmt.Printf("result: %v", list)
}
