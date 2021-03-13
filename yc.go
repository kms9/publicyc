package yc

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kms9/publicyc/pkg/conf"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/server"
	"github.com/kms9/publicyc/pkg/util"
	"github.com/kms9/publicyc/pkg/util/ogo"
	"github.com/kms9/publicyc/pkg/worker"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	initOnce    sync.Once // 应用属性初始化
	startupOnce sync.Once // 必备func启动
	stopOnce    sync.Once
	smu         *sync.RWMutex
	servers     []server.Server
	jobs        []worker.Worker
	_logger     *onion_log.Log
}

// New 实例化一个应用
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Start(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

// initialize TODO:初始化字段信息
func (a *Application) initialize() {
	a.initOnce.Do(func() {
		a.smu = &sync.RWMutex{}
		a.servers = make([]server.Server, 0)
		a._logger = onion_log.New("info", util.GetGoEnv())
	})
}

// start TODO:服务治理相关信息
func (a *Application) startFunc() (err error) {
	a.startupOnce.Do(func() {
		err = ogo.SerialUntilError(
			a.loadConfig,
			a.initLog,
		)()
	})
	return
}

// Start 服务启动
func (a *Application) Start(fns ...func() error) error {
	// 初始化app的字段
	a.initialize()

	// 初始化框架集成的方法
	if err := a.startFunc(); err != nil {
		return err
	}

	// 初始化应用传入的方法
	return ogo.SerialUntilError(fns...)()
}

// Serve 注入服务
func (a *Application) Serve(s ...server.Server) error {
	a.smu.Lock()
	defer a.smu.Unlock()
	a.servers = append(a.servers, s...)
	return nil
}

// Job 注入Job
func (a *Application) Job(jobs ...worker.Worker) error {
	a.jobs = append(a.jobs, jobs...)
	return nil
}

// loadConfig 初始化加载配置文件
func (a *Application) loadConfig() error {
	var configAddr = "config"
	err := conf.NewConfig(configAddr)
	return err
}

// initLog 初始化加载日志
func (a *Application) initLog() error {
	a._logger = onion_log.UseConfig("logger").Build()
	return nil
}

// waitSignals wait signal
func (a *Application) waitSignals() error {
	a._logger.Info("init listen signal")
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	<-sig
	return a.Stop()
}

// Stop 停止应用
func (a *Application) Stop() (err error) {
	a.stopOnce.Do(func() {
		//stop servers
		a.smu.RLock()
		for _, s := range a.servers {
			func(s server.Server) {
				_ = s.Stop()
				a._logger.Info("exit server: ", s.Info().Name)
			}(s)
		}
		a.smu.RUnlock()

		// stop jobs
		for _, j := range a.jobs {
			func(j worker.Worker) {
				_ = j.Stop()
				a._logger.Info("exit server: ", j.Name())
			}(j)
		}
	})
	return
}

// startServers 启动服务
func (a *Application) startServers() error {
	var eg errgroup.Group

	// start multi servers
	for _, s := range a.servers {
		s := s
		eg.Go(func() (err error) {
			a._logger.Info("start server: ", s.Info().Name)
			err = s.Serve()
			return
		})
	}
	return eg.Wait()
}

// startJobs 启动job
func (a *Application) startJobs() error {
	var eg errgroup.Group

	// start multi jobs
	for _, job := range a.jobs {
		j := job
		eg.Go(func() (err error) {
			a._logger.Info("start job: ", j.Name())
			err = j.Start()
			return
		})
	}
	return eg.Wait()
}

// Run 开始启动
func (a *Application) Run(servers ...server.Server) error {
	a.smu.Lock()
	a.servers = append(a.servers, servers...)
	a.smu.Unlock()

	var eg errgroup.Group

	eg.Go(func() error {
		return a.waitSignals()
	})
	eg.Go(func() error {
		return a.startJobs()
	})
	eg.Go(func() error {
		return a.startServers()
	})

	return eg.Wait()
}
