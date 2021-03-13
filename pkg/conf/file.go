package conf

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/kms9/publicyc/pkg/util"
)

var config = &Config{Viper: viper.New()}

// NewConfig path 文件夹地址
func NewConfig(path string) (err error) {
	config.init.Do(func() {
		_logger.SetFormatter(&logrus.JSONFormatter{})

		config.SetConfigType("yaml")
		_ = config.BindEnv("GO_ENV")
		absolutePath, err1 := filepath.Abs(path)
		if err1 != nil {
			_logger.Errorf("conf get folder absolutePath err %s", err1)
		}
		config.AddConfigPath(absolutePath)
		// 读取默认配置文件
		config.SetConfigName("default")
		if err := config.ReadInConfig(); err != nil || len(config.AllKeys()) == 0 {
			_logger.Errorf("Fatal errs conf file: %s", err)
		}

		// 读取环境变量的配置文件
		config.SetConfigName(util.GetGoEnv())
		if err := config.MergeInConfig(); err != nil || len(config.AllKeys()) == 0 {
			_logger.Panicf("GoEnv Fatal errs conf file: %s", err)
		}

		// 监听配置文件变化
		config.WatchConfig()
		config.OnConfigChange(func(e fsnotify.Event) {
			_logger.Info("Config file changed: ", e.Name)
		})
	})

	return err
}

// Detail 查询配置文件详情
func Detail() *Config {
	return config
}
