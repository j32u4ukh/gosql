package utils

import (
	"fmt"

	"github.com/j32u4ukh/glog"
)

var logger *glog.Logger

func InitLogger(folder string, loggerName string, level glog.LogLevel) {
	logger = glog.GetLogger(folder, loggerName, level, false)
}

func SetOptions(options ...glog.Option) {
	if logger != nil {
		logger.SetOptions(options...)
	}
}

func Debug(message string, a ...any) {
	if logger != nil {
		logger.Debug(message, a...)
	} else {
		fmt.Printf("[Debug] %s\n", fmt.Sprintf(message, a...))
	}
}

func Info(message string, a ...any) {
	if logger != nil {
		logger.Info(message, a...)
	} else {
		fmt.Printf("[Info] %s\n", fmt.Sprintf(message, a...))
	}
}

func Warn(message string, a ...any) {
	if logger != nil {
		logger.Warn(message, a...)
	} else {
		fmt.Printf("[Warn] %s\n", fmt.Sprintf(message, a...))
	}
}

func Error(message string, a ...any) {
	if logger != nil {
		logger.Debug(message, a...)
	} else {
		fmt.Printf("[Error] %s\n", fmt.Sprintf(message, a...))
	}
}
