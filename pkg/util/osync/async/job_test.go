package async

import (
	"fmt"
	"testing"
	"time"
)

func testConsumerPrint(data interface{}) error {
	fmt.Println("开始处理:", data)
	time.Sleep(1 * time.Second)
	if data == 1 {
		// 测试 运行中panic
		panic("test-do panic")
	}
	fmt.Println("处理完成:", data)
	return nil
}

func TestAsyncJob_Start(t *testing.T) {
	// 处理过程,超时(保证超时后的数据全部正常处理完在停止)
	asyncJob := NewAsyncJob()
	err := asyncJob.
		SetJobName("updateJob").
		SetParallel(2).
		SetConsumer(testConsumerPrint).
		Start()
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}
	for i := 0; i < 20; i++ {
		_ = asyncJob.Push(i)
	}
	t.Log("push 完成")
	if err = asyncJob.Stop(); err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}

func TestAsyncJob_Start2(t *testing.T) {
	// 全部处理完成
	asyncJob := NewAsyncJob()
	err := asyncJob.
		SetJobName("updateJob").
		SetParallel(2).
		SetConsumer(testConsumerPrint).
		Start()
	if err != nil {
		t.Errorf("Start() error = %v", err)
	}
	for i := 0; i < 3; i++ {
		_ = asyncJob.Push(i)
	}
	t.Log("push 完成")
	if err = asyncJob.Stop(); err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}
