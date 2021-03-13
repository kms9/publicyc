package async

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"

	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/util/ogo"
)

//  NewAsyncJob ...
func NewAsyncJob() *Job {
	return &Job{
		ch:                  make(chan interface{}, 1),
		chs:                 list.New(),
		stopSignal:          make(chan int),
		stop:                false,
		bufferDone:          make(chan int),
		allDone:             make(chan int),
		StopTimeOut:         8 * time.Second,
		StopParallelTimeout: 4 * time.Second,
		SleepDuration:       1 * time.Second,
	}
}

// Job 更新Job->独立协程 适合长时间大量创建goroutine的
type Job struct {
	StopTimeOut         time.Duration // stop后,清理链表中的超时时间
	StopParallelTimeout time.Duration // Stop后,处理channel的超时时间
	SleepDuration       time.Duration // 如果队列中没有数据 sleep一段时间

	name       string
	wg         sync.WaitGroup // 消费并发控制
	parallel   int            // 并发数
	bufferDone chan int       // 缓冲链表

	allDone    chan int // 缓冲列表中是否全部处理完成
	stopSignal chan int // 停止信号
	stop       bool     // 是否关闭

	workerOnce   sync.Once                    // 只启动一次consumer
	consumerFunc func(data interface{}) error // 消费方法
	ch           chan interface{}             // 并发处理数据
	chs          *list.List                   // 缓冲链表,用于存储来不及处理的数据
}

// SetParallel 设置并发数
func (c *Job) SetParallel(parallel int) *Job {
	if parallel <= 0 {
		parallel = 1
	}
	c.ch = make(chan interface{}, parallel)
	c.parallel = parallel
	return c
}

// SetJobName 设置消费job名称
func (c *Job) SetJobName(name string) *Job {
	c.name = name
	return c
}

// SetConsumer 加载消费者
func (c *Job) SetConsumer(f func(data interface{}) error) *Job {
	c.consumerFunc = f
	return c
}

// SetTimeout 设置超时时间
func (c *Job) SetTimeout(stopTimeOut, stopParallelTimeout time.Duration) *Job {
	c.StopTimeOut = stopTimeOut
	c.StopParallelTimeout = stopParallelTimeout
	return c
}

// SetSleep 设置没有数据sleep的时间
func (c *Job) SetSleep(sleepDuration time.Duration) *Job {
	c.SleepDuration = sleepDuration
	return c
}

// Name ...
func (c *Job) Name() string {
	return c.name
}

// Start 开始启动协程
func (c *Job) Start() error {
	go func() {
		if err := ogo.Recover(func() error {
			go c.consumer()
			c.runBufferList()
			return nil
		}); err != nil {
			onion_log.Errorf("Job: Start err: %v", err)
			_ = c.Start() // 重新拉起处理协程
		}
	}()
	return nil
}

// Stop 全部停止
func (c *Job) Stop() error {
	c.stop = true
	if c.chs.Len() != 0 {
		// 设置超时&通知结束
		timeoutCtx, _ := context.WithTimeout(context.Background(), c.StopTimeOut)
		// 通知准备结束了,如果结束则发送allDone通知
		select {
		case <-timeoutCtx.Done():
			onion_log.Warnf("%s timeout", c.name)
			c.stopSignal <- 1
		case <-c.bufferDone:
			onion_log.Infof("%s bufferDone", c.name)
			close(c.ch)
		}
	} else {
		c.stopSignal <- 1
	}

	// 设置channel中处理逻辑:超时
	channelTimeout, _ := context.WithTimeout(context.Background(), c.StopParallelTimeout)
	select {
	case <-channelTimeout.Done():
		onion_log.Warnf("%s queue timeout:%s", c.name, c.StopParallelTimeout)
	case <-c.allDone:
		onion_log.Infof("%s queue allDone", c.name)
	}
	return nil
}

// Push 推送数据
func (c *Job) Push(data interface{}) error {
	if c.stop {
		return errors.New("Job Stopped")
	}
	c.chs.PushBack(data)
	return nil
}

// do 执行相关操作
func (c *Job) do(data interface{}) {
	if err := ogo.Recover(func() error {
		return c.consumerFunc(data)
	}); err != nil {
		onion_log.Warnf("Job run: err: data:%s %v", data, err)
	}
}

// consumer 并发处理排队数据
func (c *Job) consumer() {
	if c.consumerFunc == nil {
		panic("consumerFunc is nil")
	}
	c.workerOnce.Do(func() {
		for i := 0; i < c.parallel; i++ {
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				for data := range c.ch {
					c.do(data)
				}
			}()
		}
		c.wg.Wait()
		c.allDone <- 1
	})
}

// runBufferList 开始运行,处理排队数据
func (c *Job) runBufferList() {
	for {
		select {
		case <-c.stopSignal:
			close(c.ch)
			return
		default:
			element := c.chs.Front()
			// 获取不到元素
			if element == nil {
				if c.stop {
					c.bufferDone <- 1
					return
				}
				time.Sleep(c.SleepDuration)
				continue
			}

			c.ch <- element.Value
			c.chs.Remove(element)
		}
	}
}
