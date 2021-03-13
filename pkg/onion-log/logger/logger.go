package logger

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/onion-log/hook"
)

// BaseContentInfo 外部输入的content log数据结构
type BaseContentInfo struct {
	UID     string `json:"uid"`
	TraceID string `json:"traceId"`
	From    string `json:"appFrom"`
	SpanID  string `json:"spanId"`
}

// // InitTraceID 初始化traceId
// func (b *BaseContentInfo) InitTraceID() {
//	id, _ := uuid.NewV4()
//	b.TraceID = id.String()
// }

// Logger 打印
type Logger struct {
	Content *BaseContentInfo
	Entry   *logrus.Entry
}

// WithFields 增加字段信息:一直携带此信息 知道更换content
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.Entry = l.Entry.WithFields(fields)
	return l
}

// Debug 调试信息输出
func (l *Logger) Debug(msg ...interface{}) {
	l.Entry.Debug(msg...)
}

// Debugf 格式化输出
func (l *Logger) Debugf(format string, msg ...interface{}) {
	l.Entry.Debugf(format, msg...)
}

// Info 普通信息输出
func (l *Logger) Info(msg ...interface{}) {
	l.Entry.Info(msg...)
}

// Infof 格式化输出
func (l *Logger) Infof(format string, msg ...interface{}) {
	l.Entry.Infof(format, msg...)
}

// Warn 警告信息输出
func (l *Logger) Warn(msg ...interface{}) {
	l.Entry.Warning(msg...)
}

// Warnf 警告信息输出
func (l *Logger) Warnf(format string, msg ...interface{}) {
	l.Entry.Warningf(format, msg...)
}

// Error 执行错误信息输出
func (l *Logger) Error(msg ...interface{}) {
	l.Entry.Error(msg...)
}

// Errorf 执行错误信息输出
func (l *Logger) Errorf(format string, msg ...interface{}) {
	l.Entry.Errorf(format, msg...)
}

// Panic 启动错误信息、意外退出错误信息输出
func (l *Logger) Panic(msg ...interface{}) {
	l.Entry.Panic(msg...)
}

// Panicf 启动错误信息、意外退出错误信息输出
func (l *Logger) Panicf(format string, msg ...interface{}) {
	l.Entry.Panicf(format, msg...)
}

// Notice 钉钉通知调用通知 暂时移除
func (l *Logger) Notice(msg *hook.DingMsg) {
	str, _ := json.Marshal(msg)
	l.Entry.WithFields(logrus.Fields{
		"notice": "ding",
	}).Info(string(str))
}
