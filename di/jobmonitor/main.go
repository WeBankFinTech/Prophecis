package main

import (
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"webank/DI/commons/config"
	"webank/DI/commons/logger"
	"webank/DI/jobmonitor/jobmonitor"
)

func main() {
	config.InitViper()
	viper.Set(config.LogLevelKey, "debug")
	logger.Config()
	logr := logger.GetLogger()

	signals := make(chan os.Signal)
	signal.Notify(signals)
	logr.Infof("jm server up!")

	monitor, err := jobmonitor.NewJobMonitor(nil)
	if err != nil {
		logr.Errorf("NewJobMonitor failed, %v", err.Error())
	}
	monitor.Start()

	//block goroutine/clean up
	select {
	case s := <-signals:
		logr.Infof("", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			println("someone wanna exit.")
			//TODO ExitFunc
			Cleanup()
		}
	}
}

func Cleanup() {
	os.Exit(0)
}
