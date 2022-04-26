package dss1_0

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestBMLDownload(t *testing.T) {
	dssClient := DSSClient{}
	dssClient.Address = ""
	dssClient.TokenCode = "QML-AUTH"
	downloadData := DowloadData{
		ResourceId: "173784bc-ac6f-45ec-83cf-8f89687678c9",
		Version:    "v000001",
	}

	fileByte,err := dssClient.Download(&downloadData,"alexwu")
	if err != nil {
		fmt.Errorf("")
	}

	reader := bytes.NewReader(fileByte)
	data := ioutil.NopCloser(reader)


	folderUid := uuid.New().String()
	flowFolder := "D://importFolder/" + folderUid
	//ossStorageFolder := "/data/oss-storage"
	//flowFileFolder := "flow"
	//flowDefYaml := "flow-definition.yaml"
	var flowFileName string
	uid := uuid.New().String()
	flowFileName = fmt.Sprintf("%v.zip", uid)
	flowFolderFile, err := os.Open(flowFolder)
	defer flowFolderFile.Close()
	if err != nil && os.IsNotExist(err) {
		//create file
		err := os.MkdirAll(flowFolder, os.ModePerm)
		if err != nil {
			log.Println("create folder flowFolder failed, ", err)
			//return nil, err
		}
	}
	//save zip file
	importZipFilePath := fmt.Sprintf("%v/%v", flowFolder, flowFileName)
	log.Println("import zip file, path:", importZipFilePath)
	//check file is not zip
	if !strings.HasSuffix(importZipFilePath, ".zip") {
		log.Println("import zip file is not zip file, path:", importZipFilePath)
		//return nil, errors.New("import zip file is not zip file, path:" + importZipFilePath)
	}
	//un's zip file
	unZipPathOfFlowZip := strings.TrimRight(importZipFilePath, ".zip")
	log.Println("un zip, path:", unZipPathOfFlowZip)
	//create flow.zip file and write
	flowZipFile, err := os.Create(importZipFilePath)
	if err != nil {
		log.Println("create flow.zip failed, err:", err)
		//return nil, err
	}
	//read data
	fileBytes, err := ioutil.ReadAll(data)
	_, err = flowZipFile.Write(fileBytes)
	if err != nil {
		log.Println("flow.zip write failed, ", err)
		//return nil, err
	}
	//un zip file
	err = archiver.Unarchive(importZipFilePath, unZipPathOfFlowZip)
	if err != nil {
		log.Println("un archive file failed, ", err)
		//return nil, err
	}
}