package ohttp

import "fmt"

var (
	// ErrNoURL 未传url参数
	ErrNoURL = fmt.Errorf("ms url is ''")

	// ErrNoMethod 方法未找到
	ErrNoMethod = fmt.Errorf("ms method not found")

	// ErrProjectIsNull 项目名称不存在
	ErrProjectIsNull = fmt.Errorf("ms project not found")

	// ErrAPIIsNull 项目的api不存在
	ErrAPIIsNull = fmt.Errorf("ms project's api not found")
)
