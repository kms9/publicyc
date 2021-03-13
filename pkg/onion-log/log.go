package onion_log

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/onion-log/formatter"
	"github.com/kms9/publicyc/pkg/onion-log/hook"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

type Log struct {
	Level  logrus.Level
	Logger *logrus.Logger
}

// New 实例化 Log
//     level: debug info warn errs panic
//     goEnv: development stage production
func New(level string, goEnv string, hooks ...logrus.Hook) *Log {
	parseLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parseLevel = logrus.InfoLevel
		fmt.Printf("err: log level is err:(%s) set level info", err)
	}

	l := logrus.New()
	l.Formatter = &formatter.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	}
	l.Level = parseLevel
	l.Out = os.Stdout

	if len(hooks) > 0 {
		for _, v := range hooks {
			l.AddHook(v)
		}
	}

	// 记录行号 // 并发环境会有性能问题; 只在非生产环境使用
	if goEnv != "production" {
		l.AddHook(hook.NewLineHook(logrus.ErrorLevel, logrus.PanicLevel))
	}

	return &Log{
		Level:  parseLevel,
		Logger: l,
	}
}

// With 携带基础信息
func (l *Log) With(baseContentInfo *logger.BaseContentInfo) *logger.Logger {
	if baseContentInfo == nil {
		baseContentInfo = &logger.BaseContentInfo{}
	}

	// if baseContentInfo.TraceID == "" {
	//	baseContentInfo.InitTraceID()
	// }

	return &logger.Logger{
		Content: baseContentInfo,
		Entry: l.Logger.WithFields(logrus.Fields{
			"uid":     baseContentInfo.UID,
			"traceId": baseContentInfo.TraceID,
			"appFrom": baseContentInfo.From,
			"spanId":  baseContentInfo.SpanID,
			"type":    "content",
		}),
	}
}

// Debug 调试信息输出
func (l *Log) Debug(msg ...interface{}) {
	l.Logger.Debug(msg...)
}

// Debugf 格式化输出
func (l *Log) Debugf(format string, msg ...interface{}) {
	l.Logger.Debugf(format, msg...)
}

// Info 普通信息输出
func (l *Log) Info(msg ...interface{}) {
	l.Logger.Info(msg...)
}

// Infof 格式化输出
func (l *Log) Infof(format string, msg ...interface{}) {
	l.Logger.Infof(format, msg...)
}

// Warn 警告信息输出
func (l *Log) Warn(msg ...interface{}) {
	l.Logger.Warning(msg...)
}

// Warnf 警告信息输出
func (l *Log) Warnf(format string, msg ...interface{}) {
	l.Logger.Warningf(format, msg...)
}

// Error 执行错误信息输出
func (l *Log) Error(msg ...interface{}) {
	l.Logger.Error(msg...)
}

// Errorf 执行错误信息输出
func (l *Log) Errorf(format string, msg ...interface{}) {
	l.Logger.Errorf(format, msg...)
}

// Panic 启动错误信息、意外退出错误信息输出
func (l *Log) Panic(msg ...interface{}) {
	l.Logger.Panic(msg...)
}

// Panicf 启动错误信息、意外退出错误信息输出
func (l *Log) Panicf(format string, msg ...interface{}) {
	l.Logger.Panicf(format, msg...)
}

// Notice 钉钉通知调用通知
func (l *Log) Notice(msg *hook.DingMsg) {
	str, _ := json.Marshal(msg)
	l.Logger.WithFields(logrus.Fields{
		"notice": "ding",
	}).Info(string(str))
}
