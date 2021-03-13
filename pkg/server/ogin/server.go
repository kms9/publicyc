package ogin

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kms9/publicyc/pkg/server"
)

// ModName 模块名称
const ModName = "server.gin"

// Server 服务
type Server struct {
	*gin.Engine
	Server   *http.Server
	config   *Config
	listener net.Listener
}

func newServer(config *Config) *Server {
	listener, err := net.Listen("tcp", config.Address())
	if err != nil {
		config._logger.Panicf("new gin server err: listen err %s", err)
	}
	config.Port = listener.Addr().(*net.TCPAddr).Port

	gin.SetMode(config.Mode)
	return &Server{
		Engine:   gin.New(),
		config:   config,
		listener: listener,
	}
}

// Serve implements server.Server interface.
func (s *Server) Serve() error {
	for _, route := range s.Engine.Routes() {
		s.config._logger.Logger.WithFields(map[string]interface{}{
			"method": route.Method,
			"path":   route.Path,
		}).Infof("add route")
	}
	s.Server = &http.Server{
		Addr:    s.config.Address(),
		Handler: s,
	}

	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

// Stop implements server.Server interface
// it will terminate gin server immediately
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop implements server.Server interface
// it will stop gin server gracefully
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Info returns server info, used by governor and consumer balancer
func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.listener.Addr().String()),
	)
	info.Name = fmt.Sprintf("%s.%s%s", ModName, s.config.Name, s.config.Address())
	return &info
}
