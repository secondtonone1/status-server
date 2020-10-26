package config

import "github.com/micro/go-micro/v2/registry"

var (
	config *Config
)

type BaseConf struct {
	GRPCAddr    string `toml:"grpc_addr"`
	ServiceName string `toml:"service_name"`
	RootDir     string `toml:"root_dir"`
	VarDir      string `toml:"var_dir"`
}

//日志配置
type LogConf struct {
	Project         string `toml:"project"`
	Name            string `toml:"name"`
	LogDir          string `toml:"log_dir"`
	LogLevel        string `toml:"log_level"`
	Extname         string `toml:"extname"`
	MaxSize         int    `toml:"max_size"`
	MaxNum          int    `toml:"max_num"`
	MaxDay          int    `toml:"max_day"`
	RotateSeconds   int    `toml:"rotate_seconds"`
	NotPrintLogTime bool   `toml:"not_print_log_time"`
}

//redis连接配置
type RedisConf struct {
	LuaPath string `toml:"lua_path"`
	//1单机模式 2代表集群模式。默认为1
	RedisModel              int    `toml:"redis_model"`
	SingleRedisHost         string `toml:"single_redis_host"`
	SingleRedisDb           int    `toml:"single_redis_db"`
	SingleRedisPoolSize     int    `toml:"single_redis_PoolSize"`
	SingleRedisMinIdleConns int    `toml:"single_redis_MinIdleConns"`
	SingleRedisPassword     string `toml:"single_redis_password"`

	ClusterRedisHost         []string `toml:"cluster_redis_host"`
	ClusterRedisPoolSize     int      `toml:"cluster_redis_PoolSize"`
	ClusterRedisMinIdleConns int      `toml:"cluster_redis_MinIdleConns"`
	ClusterRedisPassword     string   `toml:"cluster_reis_password"`
}

type Config struct {
	Base           BaseConf       `toml:"base"`
	LogConf        LogConf        `toml:"log_conf"`
	RedisConf      RedisConf      `toml:"redis_conf"`
	RegisterCenter RegisterCenter `toml:"register_center"`
}

type RegisterCenter struct {
	register registry.Registry
	Address  []string `toml:"address"`
	Timeout  int64    `toml:"timeout"`
}

func (r RegisterCenter) GetRegisterCenter() registry.Registry {
	return GetConf().RegisterCenter.register
}

func GetConf() *Config {
	return config
}

func GetBaseConfig() BaseConf {
	return config.Base
}
