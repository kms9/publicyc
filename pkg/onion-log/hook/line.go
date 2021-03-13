package hook

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// ContextHook for log the call context
type contextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// NewLineHook 行号
// 根据上面的推断, 我们递归深度可以设置到5即可.
func NewLineHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		Field:  "line",
		Skip:   4,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook *contextHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire implement fire
func (hook *contextHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] %s", findCaller(hook.Skip), entry.Message)
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "sirupsen/logrus") && !strings.HasPrefix(file, "onion-log/logger/logger.go") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 3 {
				file = file[i+1:]
				break
			}
		}
	}

	return file, line
}
