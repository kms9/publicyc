package ogo

import (
	"fmt"
	"runtime"

	onion_log "github.com/kms9/publicyc/pkg/onion-log"
)

// Recover 安全执行func 不panic
func Recover(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("panic recovered: %s\n%s", r, buf)
		}
	}()
	return fn()
}

// Go 安全执行func,避免野生goroutine
func Go(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			onion_log.Warnf("panic recovered: %s\n%s", r, buf)
		}
	}()
	fn()
}
