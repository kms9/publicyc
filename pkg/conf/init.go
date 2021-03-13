package conf

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	_logger = logrus.New()
)

type Config struct {
	init sync.Once

	*viper.Viper
}

func (c *Config) UnmarshalKey(k string, result interface{}) error {
	k = strings.ToLower(k)
	return c.Viper.UnmarshalKey(k, result)
}

func (c *Config) Get(k string) interface{} {
	k = strings.ToLower(k)
	return c.Viper.Get(k)
}
