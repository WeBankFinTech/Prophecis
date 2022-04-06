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
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var logger *logrus.Logger

func init() {
	log := logrus.New()

	log.Formatter = &logrus.JSONFormatter{}

	rl, err := rotatelogs.New("logs/access_log.%Y%m%d%H%M", rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, rl)

	log.SetOutput(mw)
	log.SetLevel(logrus.DebugLevel)
	logger = log
}

func Logger() *logrus.Logger {
	return logger
}
