package ogo

// SerialUntilError 创建一个迭代器
func SerialUntilError(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := try(fn); err != nil {
				return err
			}
		}
		return nil
	}
}
