package onion_log

import (
	"github.com/kms9/publicyc/pkg/onion-log/hook"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

func Example() {
	log := New("info", "developemnt")

	l := log.With(&logger.BaseContentInfo{
		TraceID: "39b99823-9f9b-4afc-a9a2-2ae33d90a690",
		UID:     "5eb5200ec23e490001cf405a",
	})

	l.Warn("这是测试", "warn")
	l.Warnf("这是测试 %s %d", "warnf", 1)
	l.Info("这是测试 info")
	l.Infof("这是测试 %s", "infof")
	l.Error("这是测试", "errs")
	l.Errorf("这是测试 %s %d", "errorf", 1)
	l.Debug("这是测试 debug")
	l.Debugf("这是测试 %s", "debugf")
	l.Notice(&hook.DingMsg{
		Msgtype: "text",
		Text:    &hook.DingMsgText{Content: "test-3214321432-342143214321-43214231"},
	}) // 不配置不会真实通知
	l.Panic("这是测试 panic")
}
