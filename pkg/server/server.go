package server

import "context"

type Option func(c *ServiceInfo)

// ServiceInfo represents service info
type ServiceInfo struct {
	Name    string               `json:"name"`
	Scheme  string               `json:"scheme"`
	Address string               `json:"address"`
}

// Server ...
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name: "",
	}
	return si
}
