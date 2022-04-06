package logger

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

var mylogger *logrus.Logger

func Logger() *logrus.Logger {
	return mylogger
}

func init() {
	rl, _ := rotatelogs.New("logs/access_log.%Y%m%d%H%M", rotatelogs.WithRotationTime(24*time.Hour))

	log := logrus.New()
	log.SetReportCaller(true)
	log.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	mw := io.MultiWriter(os.Stdout, rl)
	//log.SetFormatter(&logrus.JSONFormatter{})
	level := strings.ToUpper(viper.GetString("loglevel"))
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Error("init logger level err, ",err)
		return
	}
	log.SetLevel(logLevel)
	log.SetOutput(mw)
	mylogger = log
}
