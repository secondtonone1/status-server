package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"status-server/config"
	"status-server/constants"
	"time"
)

var (
	//定义变量
	ClusterClient *redis.ClusterClient
	Client        *redis.Client
	Model         int
	IdleTimeout   = time.Second * 10
	PoolTimeout   = time.Second * 60
	MaxConnAge    = time.Second * 120
	PoolSize      = 512
	MinIdleConn   = 512
)

/**
初始化redis连接配置
*/
func InitRedis(cfg config.RedisConf) {

	//判断redis使用模式  1单机模式 2代表集群模式。默认为1
	if cfg.RedisModel == 0 {
		cfg.RedisModel = 1
	}
	Model = cfg.RedisModel

	if Model == constants.RedisModel {
		host := cfg.ClusterRedisHost
		if cfg.ClusterRedisPoolSize <= 0 {
			cfg.ClusterRedisPoolSize = PoolSize
		}
		if cfg.ClusterRedisMinIdleConns <= 0 {
			cfg.ClusterRedisMinIdleConns = MinIdleConn
		}

		ClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        host,
			PoolSize:     cfg.ClusterRedisPoolSize,
			MinIdleConns: cfg.SingleRedisMinIdleConns,
			Password:     cfg.ClusterRedisPassword,
			IdleTimeout:  IdleTimeout,
			PoolTimeout:  PoolTimeout,
			MaxConnAge:   MaxConnAge,
		})
		statusCmd := ClusterClient.Ping()
		if statusCmd.Err() != nil {
			panic("redis初始化连接失败！" + fmt.Sprint(statusCmd.Err()))
		}
	} else {
		if cfg.SingleRedisDb < 0 || cfg.SingleRedisDb > 15 {
			cfg.SingleRedisDb = 0
		}
		if cfg.SingleRedisPoolSize <= 0 {
			cfg.SingleRedisPoolSize = PoolSize
		}
		if cfg.SingleRedisMinIdleConns <= 0 {
			cfg.SingleRedisMinIdleConns = MinIdleConn
		}
		Client = redis.NewClient(&redis.Options{
			Addr:         cfg.SingleRedisHost,
			Password:     cfg.SingleRedisPassword,
			DB:           cfg.SingleRedisDb,
			PoolSize:     cfg.SingleRedisPoolSize,
			MinIdleConns: cfg.SingleRedisMinIdleConns,
			IdleTimeout:  IdleTimeout,
			PoolTimeout:  PoolTimeout,
			MaxConnAge:   MaxConnAge,
		})
		statusCmd := Client.Ping()
		if statusCmd.Err() != nil {
			panic("redis初始化连接失败！" + fmt.Sprint(statusCmd.Err()))
		}
	}

}

func Close() {
	if Model == constants.RedisModel {
		ClusterClient.Close()
	} else {
		Client.Close()
	}
}