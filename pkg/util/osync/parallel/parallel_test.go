package parallel

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func sleep1second(second int) {
	print("start: ", second)
	if second > 10 {
		second = 10
	}
	time.Sleep(time.Duration(second) * time.Second)
	fmt.Println("second:", second)
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	errGroup := WithContext(ctx)
	now := time.Now()
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(1)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(2)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(3)
		return nil
	})
	if err := errGroup.Wait(); err != nil {
		t.Error(err)
	}
	fmt.Println("总用时:", time.Now().Sub(now).String())
}

func TestWithContext1(t *testing.T) {
	ctx := context.Background()
	errGroup := WithContext(ctx)
	now := time.Now()
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(1)
		return fmt.Errorf("test: %d", 1)
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(2)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(4)
		return nil
	})
	if err := errGroup.Wait(); err == nil {
		t.Error(err)
	} else {
		t.Log(err)
	}
	fmt.Println("总用时:", time.Now().Sub(now).String())
}

func TestWithCancel(t *testing.T) {
	ctx := context.Background()
	errGroup := WithCancel(ctx)
	now := time.Now()
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(1)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(2)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(3)
		return nil
	})
	if err := errGroup.Wait(); err != nil {
		t.Error(err)
	}
	fmt.Println("总用时:", time.Now().Sub(now).String())
}

func TestWithCancel1(t *testing.T) {
	ctx := context.Background()
	errGroup := WithCancel(ctx)
	now := time.Now()
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(1)
		return nil
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(0)
		return fmt.Errorf("test: %d", 1)
	})
	errGroup.Go(func(ctx context.Context) error {
		sleep1second(1)
		return nil
	})
	if err := errGroup.Wait(); err == nil {
		t.Error("err")
	} else {
		t.Log("some:", err)
	}

	fmt.Println("总用时:", time.Now().Sub(now).String())
}

func TestSetMaxGo(t *testing.T) {
	ctx := context.Background()
	errGroup := WithCancel(ctx)
	now := time.Now()
	errGroup.SetMaxGo(100)
	for i := 0; i < 200; i++ {
		a := i
		errGroup.Go(func(ctx context.Context) error {
			sleep1second(a)
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		t.Error("err", err)
	}

	fmt.Println("总用时:", time.Now().Sub(now).String())
}
