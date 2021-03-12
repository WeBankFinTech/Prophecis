package main

import (
	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/commons/util"
	"webank/DI/storage/storage"

	"github.com/spf13/viper"
)

func main() {
	config.InitViper()
	logger.Config()

	getLogger := logger.GetLogger()
	uploadContainerPath := viper.GetString(config.UploadContainerPath)
	getLogger.Debugf("uploadContainerPath: %v", uploadContainerPath)
	getLogger.Debugf("start storage service: %v", uploadContainerPath)

	port := viper.GetInt(config.PortKey)
	if port == 0 {
		port = 30005 // TODO don't hardcode
	}
	service := storage.NewService()

	util.HandleOSSignals(func() {
		service.StopTrainer()
	})
	service.Start(port, false)
}
