package onion_log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

// GinLogMiddle ogin 日志中间件
func GinLogMiddle(appName string, goEnv string) gin.HandlerFunc {
	log := New("info", goEnv)

	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientVersion := c.GetHeader("setup-version")
		clientType := c.GetHeader("setup-type")
		bodyData := ""

		c.Request.Header.Set("span-id", uuid.NewV4().String()) // 设置spanId
		baseContent := GetBaseByContext(c)
		from := c.GetHeader("app-from")
		if from != "" {
			from = fmt.Sprintf("%s;%s", from, appName)
		} else {
			from = appName
		}
		c.Request.Header.Set("app-from", from)

		body, err := c.GetRawData()
		if err == nil && len(body) > 0 {
			bodyData = string(body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		log.With(baseContent).Entry.WithFields(logrus.Fields{
			"type":          "api",
			"method":        method,
			"path":          path,
			"query":         query,
			"body":          bodyData,
			"inout":         "in",
			"clientVersion": clientVersion,
			"clientType":    clientType,
			"appFrom":       from,
		}).Info("in")

		c.Next()

		latencyTime := time.Now().Sub(startTime) // 执行时间
		statusCode := c.Writer.Status()

		log.With(baseContent).Entry.WithFields(logrus.Fields{
			"type":          "api",
			"method":        method,
			"path":          path,
			"query":         query,
			"body":          bodyData,
			"inout":         "out",
			"clientVersion": clientVersion,
			"clientType":    clientType,
			"status":        statusCode,
			"latency":       latencyTime.Milliseconds(),
			"appFrom":       from,
		}).Info("out")
	}
}

// GetBaseByContext 通过gin context 获取基础内容
func GetBaseByContext(c *gin.Context) *logger.BaseContentInfo {
	if c == nil {
		return &logger.BaseContentInfo{}
	}
	baseContentInfo := &logger.BaseContentInfo{
		UID:     c.GetHeader("uid"),
		TraceID: c.GetHeader("trace-id"),
		From:    c.GetHeader("app-from"),
		SpanID:  c.GetHeader("span-id"),
	}

	//if logger.BaseContentInfo.TraceID == "" {
	//	logger.BaseContentInfo.InitTraceID()
	//	c.Request.Header.Add("traceId", logger.BaseContentInfo.TraceID)
	//}

	return baseContentInfo
}
