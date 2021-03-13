package hook

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

// DingHook
type NotifyHook struct {
	Type string `json:"type"`
	Ding *Ding
}

// Fire 对数据执行相关操作
func (n *NotifyHook) Fire(entry *logrus.Entry) error {
	if v, exists := entry.Data["notice"]; exists {
		switch v {
		case "ding":
			if err := n.ding(entry); err != nil {
				return err
			}
		case "mail":
			// 可以单独配置其他的通知方式
		default:
			return errors.New("notify type is unkown")
		}
	}
	return nil
}

// Levels 哪个级别通知
func (n *NotifyHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

// ding 钉钉通知
func (n *NotifyHook) ding(entry *logrus.Entry) error {
	if n.Ding == nil {
		return errors.New("is not have dingding conf")
	}
	str := entry.Message
	// 执行钉钉通知
	go func() {
		var err error
		defer func() {
			if err != nil {
				fmt.Printf("Failed to fire hook: %s", err.Error())
			}
			if err := recover(); err != nil {
				fmt.Printf("Failed to fire hook: %v", err)
			}
		}()
		msg := &DingMsg{}
		if err = json.Unmarshal([]byte(str), msg); err != nil {
			return
		}
		if err = n.Ding.Send(msg); err != nil {
			return
		}
	}()
	return nil
}
