package onion_log

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/onion-log/hook"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

func BenchmarkTestLogrus(b *testing.B) {
	b.ReportAllocs()
	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}
	l.Level = logrus.InfoLevel
	l.Out = os.Stdout
	en := l.WithFields(map[string]interface{}{"traceId": "39b99823-9f9b-4afc-a9a2-2ae33d90a690", "uid": ""})
	l.Info("test")
	for i := 0; i < b.N; i++ {
		en.Info("test:1")
	}
}

func BenchmarkTestLog_Info(b *testing.B) {
	b.ReportAllocs()
	log := New("info", "development")

	logWith := log.With(&logger.BaseContentInfo{
		UID:     "",
		TraceID: "39b99823-9f9b-4afc-a9a2-2ae33d90a690",
	})
	for i := 0; i < b.N; i++ {
		logWith.Info("test:1")
	}
}

func BenchmarkTestLog_Infof2(b *testing.B) {
	log := New("info", "development")

	for i := 0; i < b.N; i++ {
		log.With(&logger.BaseContentInfo{
			UID:     "",
			TraceID: "39b99823-9f9b-4afc-a9a2-2ae33d90a690",
		}).Info("test:1")
	}
}

/**
测试环境:
{"level":"warn","logType":"content","msg":"[onion-log/logger/log_test.go:73] 这是测试warn","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"warn","logType":"content","msg":"[onion-log/logger/log_test.go:74] 这是测试 warnf 1","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"info","logType":"content","msg":"[onion-log/logger/log_test.go:75] 这是测试 info","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"info","logType":"content","msg":"[onion-log/logger/log_test.go:76] 这是测试 infof","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"errs","logType":"content","msg":"[onion-log/logger/log_test.go:77] 这是测试error","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"errs","logType":"content","msg":"[onion-log/logger/log_test.go:78] 这是测试 errorf 1","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"debug","logType":"content","msg":"[onion-log/logger/log_test.go:79] 这是测试 debug","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"debug","logType":"content","msg":"[onion-log/logger/log_test.go:80] 这是测试 debugf","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"notice","logType":"content","msg":"[onion-log/logger/log_test.go:81] {\"msgtype\":\"text\",\"text\":{\"content\":\"test-3214321432-342143214321-43214231\"}}","notice":"ding","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"panic","logType":"content","msg":"[onion-log/logger/log_test.go:85] 这是测试 panic","time":"2020-05-12 15:11:05","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
*/
func TestLog_development_Info(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
	l := New("debug", "development").With(&logger.BaseContentInfo{
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

/**
{"level":"warn","logType":"content","msg":"这是测试warn","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"warn","logType":"content","msg":"这是测试 warnf 1","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"info","logType":"content","msg":"这是测试 info","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"info","logType":"content","msg":"这是测试 infof","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"errs","logType":"content","msg":"这是测试error","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"errs","logType":"content","msg":"这是测试 errorf 1","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"debug","logType":"content","msg":"这是测试 debug","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"debug","logType":"content","msg":"这是测试 debugf","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"notice","logType":"content","msg":"{\"msgtype\":\"text\",\"text\":{\"content\":\"test-3214321432-342143214321-43214231\"}}","notice":"ding","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
{"level":"panic","logType":"content","msg":"这是测试 panic","time":"2020-05-12 15:09:54","traceId":"39b99823-9f9b-4afc-a9a2-2ae33d90a690","uid":"5eb5200ec23e490001cf405a"}
*/
func TestLog_production_Info(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
		}
	}()
	l := New("debug", "production").With(&logger.BaseContentInfo{
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

// 钉钉通知测试
func TestLog_Notice(t *testing.T) {
	l := New("debug", "development", &hook.NotifyHook{
		Type: "ding",
		Ding: &hook.Ding{
			"SEC2da2d500ed1b2ed09251906b24db6c588d2ff6d83b4aef4e1b602f9296c60afc",
			"f136ce14e6746141c3646ca097923b822fc9387497b3ff61202f90230600ae57",
		},
	}).With(&logger.BaseContentInfo{
		TraceID: "8x8x8f8safdsafdsaa-12312321312312",
	})
	l.Notice(&hook.DingMsg{
		Msgtype: "text",
		Text:    &hook.DingMsgText{Content: "test-3214321432-342143214321-43214231"},
	})
	time.Sleep(3 * time.Second)
}
