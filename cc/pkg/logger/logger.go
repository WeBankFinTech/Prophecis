/*
 * Copyright 2020 WeBank
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package logger

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
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
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(mw)
	mylogger = log
}

//func InitLogger2() {
//
//	rl, _ := rotatelogs.New("logs/access_log.%Y%m%d%H%M", rotatelogs.WithRotationTime(time.Minute))
//
//	log := logrus.New()
//	log.SetReportCaller(true)
//	log.Formatter = &logrus.TextFormatter{
//		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
//			filename := path.Base(f.File)
//			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
//		},
//	}
//
//	mw := io.MultiWriter(os.Stdout, rl)
//	log.SetFormatter(&logrus.JSONFormatter{})
//	log.SetLevel(logrus.DebugLevel)
//	log.SetOutput(mw)
//
//}
