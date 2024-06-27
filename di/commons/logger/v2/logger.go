package v2

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	RequestIdLogFieldKey string = "X-Request-Id"
	OperationLogFieldKey string = "Operation"
)

func GetLogger(context context.Context) *log.Entry {
	formatter := new(log.TextFormatter)

	l := &log.Logger{
		Out:       os.Stdout,
		Formatter: formatter,
		Hooks:     make(log.LevelHooks),
		Level:     log.DebugLevel,
		ExitFunc:  os.Exit,
		// ReportCaller 会把file(调用的具体哪一行）和func（调用的函数）同时打开
		ReportCaller: true,
	}
	entry := log.NewEntry(l)
	if v, ok := context.Value(RequestIdLogFieldKey).(string); ok {
		entry = entry.WithField(RequestIdLogFieldKey, v)
	}
	if v, ok := context.Value(OperationLogFieldKey).(string); ok {
		entry = entry.WithField(OperationLogFieldKey, v)
	}
	return entry
}
