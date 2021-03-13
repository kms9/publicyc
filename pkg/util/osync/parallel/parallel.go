// Package parallel 参考 https://pkg.go.dev/golang.org/x/sync/errgroup
package parallel

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// ErrCtxCancel 级联取消
var ErrCtxCancel = fmt.Errorf("context cancel")

// Group 并行方法
type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error

	ctx    context.Context
	cancel func()

	// 并发控制
	workerOnce sync.Once
	ch         chan func(ctx context.Context) error
	chs        []func(ctx context.Context) error
}

// WithCancel 级联取消,只要有一处报错即返回
func WithCancel(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{ctx: ctx, cancel: cancel}
}

// WithContext 不会因一处取消全部取消
func WithContext(ctx context.Context) *Group {
	return &Group{ctx: ctx}
}

// SetMaxGo 设置最大并发数
func (g *Group) SetMaxGo(n int) {
	if n <= 0 {
		panic("parallel: SetGoMax n must > 0")
	}
	g.workerOnce.Do(func() {
		g.ch = make(chan func(context.Context) error, n)
		for i := 0; i < n; i++ {
			go func() {
				for f := range g.ch {
					g.do(f)
				}
			}()
		}
	})
}

// Go 执行并行
func (g *Group) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	if g.ch != nil {
		select {
		case g.ch <- f:
		default:
			g.chs = append(g.chs, f)
		}
		// 如果是并发控制的代码,go代码由其他地方执行
		return
	}
	go g.do(f)
}

func (g *Group) do(f func(ctx context.Context) error) {
	ctx := g.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	var err error
	defer func() {
		// 防止goroutine内panic
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = f(ctx)
}

// Wait 并行阻塞代码
func (g *Group) Wait() error {
	// 并发控制
	if g.ch != nil {
		// 持续将数组中的方法提前出来,去执行
		for _, f := range g.chs {
			g.ch <- f
		}
	}
	g.wg.Wait()
	// 关闭并发channel
	if g.ch != nil {
		close(g.ch)
	}
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}
