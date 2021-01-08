package httpgateway

import (
	"fmt"
	"time"
)

var defaultLogger = Logger{}
var CustomerLogger LogInterface = nil

func GetLogger() LogInterface {
	if CustomerLogger == nil {
		return defaultLogger
	} else {
		return CustomerLogger
	}
}


// 定义log的接口
type LogInterface interface {
	Debug(format string)
	Info(format string)
	Warn(format string)
	Error(format string)
}

type Logger struct {

}

func (l Logger) Debug(format string) {
	format = time.Now().Format("2006-01-02 15:04:05.000") + "    " + format
	fmt.Println(format)
}

func (l Logger) Info(format string) {
	format = time.Now().Format("2006-01-02 15:04:05.000") + "    " + format
	fmt.Println(format)
}

func (l Logger) Warn(format string) {
	format = time.Now().Format("2006-01-02 15:04:05.000") + "    " + format
	fmt.Println(format)
}

func (l Logger) Error(format string) {
	format = time.Now().Format("2006-01-02 15:04:05.000") + "    " + format
	fmt.Println(format)
}


