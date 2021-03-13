package onion_log

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/onion-log/hook"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
	"github.com/kms9/publicyc/pkg/util"
)

var log *Log
var _logger *logger.Logger

func init() {
	log = New("info", util.GetGoEnv())
	_logger = log.With(&logger.BaseContentInfo{UID: "yc", SpanID: uuid.NewV4().String()})
}

func DefaultLog() *logger.Logger {
	return _logger
}

func DefaultLogger() *Log {
	return log
}

// Debug 调试信息输出
func Debug(msg ...interface{}) {
	_logger.Entry.Debug(msg...)
}

// Debugf 格式化输出
func Debugf(format string, msg ...interface{}) {
	_logger.Entry.Debugf(format, msg...)
}

// Info 普通信息输出
func Info(msg ...interface{}) {
	_logger.Entry.Info(msg...)
}

// Infof 格式化输出
func Infof(format string, msg ...interface{}) {
	_logger.Entry.Infof(format, msg...)
}

// Warn 警告信息输出
func Warn(msg ...interface{}) {
	_logger.Entry.Warning(msg...)
}

// Warnf 警告信息输出
func Warnf(format string, msg ...interface{}) {
	_logger.Entry.Warningf(format, msg...)
}

// Error 执行错误信息输出
func Error(msg ...interface{}) {
	_logger.Entry.Error(msg...)
}

// Errorf 执行错误信息输出
func Errorf(format string, msg ...interface{}) {
	_logger.Entry.Errorf(format, msg...)
}

// Panic 启动错误信息、意外退出错误信息输出
func Panic(msg ...interface{}) {
	_logger.Entry.Panic(msg...)
}

// Panicf 启动错误信息、意外退出错误信息输出
func Panicf(format string, msg ...interface{}) {
	_logger.Entry.Panicf(format, msg...)
}

// Notice 钉钉通知调用通知 暂时移除
func Notice(msg *hook.DingMsg) {
	str, _ := json.Marshal(msg)
	_logger.Entry.WithFields(logrus.Fields{
		"notice": "ding",
	}).Info(string(str))
}
