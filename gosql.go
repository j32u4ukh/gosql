package gosql

import (
	"fmt"

	"github.com/j32u4ukh/glog"
)

var gs *GoSql

func init() {
	gs = &GoSql{}
}

type GoSql struct {
	*glog.Logger
}

func InitLogger(folder string, loggerName string, level glog.LogLevel) {
	gs.Logger = glog.GetLogger(folder, loggerName, level, false)
}

func SetOptions(options ...glog.Option) {
	if gs.Logger != nil {
		gs.Logger.SetOptions(options...)
	}
}

func Debug(message string, a ...any) {
	if gs.Logger != nil {
		gs.Logger.Debug(message, a...)
	} else {
		fmt.Printf("[Debug] %s\n", fmt.Sprintf(message, a...))
	}
}

func Info(message string, a ...any) {
	if gs.Logger != nil {
		gs.Logger.Info(message, a...)
	} else {
		fmt.Printf("[Info] %s\n", fmt.Sprintf(message, a...))
	}
}

func Warn(message string, a ...any) {
	if gs.Logger != nil {
		gs.Logger.Warn(message, a...)
	} else {
		fmt.Printf("[Warn] %s\n", fmt.Sprintf(message, a...))
	}
}

func Error(message string, a ...any) {
	if gs.Logger != nil {
		gs.Logger.Debug(message, a...)
	} else {
		fmt.Printf("[Error] %s\n", fmt.Sprintf(message, a...))
	}
}
