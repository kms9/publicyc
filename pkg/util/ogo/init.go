package ogo

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
	"github.com/kms9/publicyc/pkg/util"
	"github.com/kms9/publicyc/pkg/util/ostring"
)

var (
	_logger = onion_log.New("info", util.GetGoEnv()).With(&logger.BaseContentInfo{UID: "ogo"})
)

// try 执行
func try(fn func() error) (ret error) {
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(2)
			_logger.Errorf("recover err: %v line: %s:%d", err, file, line)

			if _, ok := err.(error); ok {
				ret = err.(error)
			} else {
				ret = fmt.Errorf("%+v", err)
			}
		ret = errors.Wrap(ret, fmt.Sprintf("%s:%d", ostring.FunctionName(fn), line))
		}
	}()
	return fn()
}
