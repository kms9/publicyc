package oredis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	config *Config
	pool   *redis.Pool
}

// ping 测试链接
func (r *Redis) ping() bool {
	if result, err := redis.String(r.Do("PING")); err != nil || result != "PONG" {
		return false
	}
	return true
}

// GetConfig 查询配置
func (r *Redis) GetConfig() *Config {
	return r.config
}

// GetConfig 获取格式化的名称
func (r *Redis) GetAddrDB() string {
	return fmt.Sprintf("%s %s/%d", r.config.Name, r.config.Addr, r.config.DB)
}

// Do 执行命令
func (r *Redis) Do(command string, args ...interface{}) (interface{}, error) {
	var pool = r.pool.Get()
	slime := time.Now()
	defer func() {
		if err := pool.Close(); err != nil {
			r.config._logger.Errorf("%s: %s/%d pool.Close err:%s", r.config.Name, r.config.Addr, r.config.DB, err)
		}
		// 打印慢查询
		if t := time.Now().Sub(slime); t > r.config.SlowThreshold {
			r.config._logger.Warnf("redis: %s do command over %s use: %s command: %s %s", r.config.Name, r.config.SlowThreshold.String(), t.String(), command, args)
		}
	}()
	return pool.Do(command, args...)
}

// Lock redis 加锁
func (r *Redis) Lock(key string, value string, expire int64) bool {
	if result, err := r.Do("SET", key, value, "EX", expire, "NX"); err != nil || result == nil {
		return false
	}
	return true
}

// RmLock redis 移除锁
//     value 为空则直接删除;否则按照value值删除
func (r *Redis) RmLock(key string, value string) (bool, error) {
	if value == "" {
		if _, err := r.Do("DEL", key); err != nil {
			return false, err
		}
		return true, nil
	}
	result, err := redis.Int64(r.Do("EVAL", "if redis.call(\"get\", KEYS[1]) == ARGV[1] then return redis.call(\"del\", KEYS[1]) else return 0 end", 1, key, value))
	if err != nil {
		return false, err
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

// GetConnect 获取连接&close方法
func (r *Redis) GetConnect() (redis.Conn, func() error) {
	conn := r.pool.Get()
	return conn, func() error {
		return conn.Close()
	}
}
