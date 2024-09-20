package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func SetLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Unable to start logger" + err.Error())
	}
}
func I(log ...interface{}) {
	Logger.Info(fmt.Sprint(log...))
}
func D(log ...interface{}) {
	Logger.Debug(fmt.Sprint(log...))
}
func E(log ...interface{}) {
	Logger.Error(fmt.Sprint(log...))
}
