package ogin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
)

var log *onion_log.Log

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "    %s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// ErrorTrace 异常错误捕获
func ErrorTrace(log *onion_log.Log) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				log.With(onion_log.GetBaseByContext(c)).WithFields(map[string]interface{}{"err": err, "stack": string(stack), "brokenPipe": brokenPipe}).Error("recovery")

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Abort()
					return
				}

				c.JSON(http.StatusInternalServerError, ErrMsg{
					Msg:   "Sorry，服务器累瘫了 开发人员正拼命抢修\n待会再来试试~",
					Debug: fmt.Sprint(err),
				})
			}
		}()
		c.Next()
	}
}

// LogMiddle ogin 日志中间件
func LogMiddle(log *onion_log.Log, name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		clientVersion := c.GetHeader("client-version")
		clientType := c.GetHeader("client-type")
		bodyData := ""

		c.Request.Header.Set("span-id", uuid.NewV4().String()) // 设置spanId

		baseContent := onion_log.GetBaseByContext(c)
		from := c.GetHeader("app-from")
		if from != "" {
			from = fmt.Sprintf("%s;%s", from, name)
		} else {
			from = name
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
