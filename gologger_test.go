package gologger

import (
	"testing"
)

func TestGoLogger(t *testing.T) {
	LogConf.TimeZone = "Asia/Jakarta"
	LogConf.TimeFormat = "2006-01-02T15:04:05-0700"
	LogConf.CreateLogFile = true
	LogConf.Path = "./logs"
	LogConf.DebugMode = true
	LogConf.NestedLocationLevel = 1
	LogConf.LogFuncName = true
	LogConf.NestedFuncLevel = 1
	LogConf.RuntimeCallerSkip = 2

	Warning("Log warning")
	Info("Log information")
	Info("Log information")
	Info("Log information")
	Info("Log information")
	Info("Log information")
	Info("Log information")
	Error("Log error")
	PrintJSONIndent("Log Config", LogConf)
}
