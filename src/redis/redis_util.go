package redis

import (
	"status-server/constants"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

//添加有序set
func Zadd(key string, score float64, member interface{}) {
	if Model == constants.RedisModel {
		ClusterClient.ZAdd(key, redis.Z{Score: score, Member: member})
	} else {
		Client.ZAdd(key, redis.Z{Score: score, Member: member})
	}
}

//查询有序成员分
func Zscore(key string, member string) (float64, error) {
	var r *redis.FloatCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ZScore(key, member)
	} else {
		r = Client.ZScore(key, member)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return 0, r.Err()
	}
}

//查询加权区间成员数量
func ZCount(key, min, max string) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ZCount(key, min, max)
	} else {
		r = Client.ZCount(key, min, max)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return 0, r.Err()
	}
}

//移除有序集合中的一个或多个成员
func Zrem(key string, member ...interface{}) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ZRem(key, member)
	} else {
		r = Client.ZRem(key, member)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return 0, r.Err()
	}
}

//移除有序集中，指定排名(rank)区间内的所有成员。
func ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ZRemRangeByRank(key, start, stop)
	} else {
		r = Client.ZRemRangeByRank(key, start, stop)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return 0, r.Err()
	}
}

//命令用于移除有序集中，指定分数（score）区间内的所有成员。
func ZRemRangeByScore(key, min, max string) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ZRemRangeByScore(key, min, max)
	} else {
		r = Client.ZRemRangeByScore(key, min, max)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return 0, r.Err()
	}
}

//将脚本 script 添加到脚本缓存中，但并不立即执行这个脚本。
func ScriptLoad(script string) (string, error) {
	var r *redis.StringCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ScriptLoad(script)
	} else {
		r = Client.ScriptLoad(script)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return "", r.Err()
	}
}

//命令根据给定的 sha1 校验码，执行缓存在服务器中的脚本。
func EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	var r *redis.Cmd
	if Model == constants.RedisModel {
		r = ClusterClient.EvalSha(sha1, keys, args...)
	} else {
		r = Client.EvalSha(sha1, keys, args...)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return nil, r.Err()
	}
}

//校验指定的脚本是否已经被保存在缓存当中
func ScriptExists(sha1 string) ([]bool, error) {
	var r *redis.BoolSliceCmd
	if Model == constants.RedisModel {
		r = ClusterClient.ScriptExists(sha1)
	} else {
		r = Client.ScriptExists(sha1)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return nil, r.Err()
	}
}

//获取管道
func Pipeline() redis.Pipeliner {
	var pipeLiner redis.Pipeliner
	if Model == constants.RedisModel {
		pipeLiner = ClusterClient.Pipeline()
	} else {
		pipeLiner = Client.Pipeline()
	}
	return pipeLiner
}

//获取有序集合数据 redis start 0开始
func Zrange(key string, start, stop int64, withScores bool) (interface{}, error) {
	var zsc *redis.ZSliceCmd
	var ssc *redis.StringSliceCmd
	if Model == constants.RedisModel {
		//带权重查询
		if withScores {
			zsc = ClusterClient.ZRangeWithScores(key, start, stop)
		} else {
			ssc = ClusterClient.ZRange(key, start, stop)
		}
	} else {
		if withScores {
			zsc = Client.ZRangeWithScores(key, start, stop)
		} else {
			ssc = Client.ZRange(key, start, stop)
		}
	}

	//带权重查询
	if withScores {
		//查询异常
		if zsc != nil && zsc.Err() != nil {
			return nil, zsc.Err()
		} else if zsc == nil {
			return nil, nil
		} else {
			return zsc.Val(), nil
		}
	} else {
		//查询异常
		if ssc != nil && ssc.Err() != nil {
			return nil, ssc.Err()
		} else if ssc == nil {
			return nil, nil
		} else {
			return ssc.Val(), nil
		}
	}

}

//通过管道设置set字段，并且设置过期时间
func HsetExpire(key, field string, value interface{}, expiration time.Duration) error {
	pipeLiner := Pipeline()
	defer pipeLiner.Close()

	pipeLiner.HSet(key, field, value)
	pipeLiner.Expire(key, expiration)
	_, e := pipeLiner.Exec()

	return e
}

//删除hset指定field
func HDelFiled(key string, fields []string) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.HDel(key, fields...)
	} else {
		r = Client.HDel(key, fields...)
	}
	return r.Result()
}

//设置Hset字段值
func Hset(key, field string, value interface{}) error {
	var r *redis.BoolCmd
	if Model == constants.RedisModel {
		r = ClusterClient.HSet(key, field, value)
	} else {
		r = Client.HSet(key, field, value)
	}
	return r.Err()
}

//设置hget字段值
func Hget(key, field string) (string, error) {
	var r *redis.StringCmd
	if Model == constants.RedisModel {
		r = ClusterClient.HGet(key, field)
	} else {
		r = Client.HGet(key, field)
	}
	return r.Val(), r.Err()
}

//取得key剩过期时间，单位为秒
func TTL(key string) (time.Duration, error) {
	var r *redis.DurationCmd
	if Model == constants.RedisModel {
		r = ClusterClient.TTL(key)
	} else {
		r = Client.TTL(key)
	}
	return r.Val(), r.Err()
}

//取得key剩过期时间，单位为毫秒
func PTTL(key string) (time.Duration, error) {
	var r *redis.DurationCmd
	if Model == constants.RedisModel {
		r = ClusterClient.PTTL(key)
	} else {
		r = Client.PTTL(key)
	}
	return r.Val(), r.Err()
}

//如果不存在，则就进行保存。否则返回失败
func SetNx(key, requestId string, duration time.Duration) (bool, error) {
	var r *redis.BoolCmd
	if Model == constants.RedisModel {
		r = ClusterClient.SetNX(key, requestId, duration)
	} else {
		r = Client.SetNX(key, requestId, duration)
	}
	return r.Val(), r.Err()
}

func Get(key string) (string, error) {
	var r *redis.StringCmd
	if Model == constants.RedisModel {
		r = ClusterClient.Get(key)
	} else {
		r = Client.Get(key)
	}
	return r.Val(), r.Err()
}

//如果不存在，则就进行保存。否则返回失败
func Set(key string, value interface{}, duration time.Duration) (string, error) {
	var r *redis.StatusCmd
	if Model == constants.RedisModel {
		r = ClusterClient.Set(key, value, duration)
	} else {
		r = Client.Set(key, value, duration)
	}
	return r.Val(), r.Err()
}

func ReleseDistributedLoock(key, requestId string) (interface{}, error) {
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	var r *redis.Cmd
	if Model == constants.RedisModel {
		r = ClusterClient.Eval(script, []string{key}, requestId)
	} else {
		r = Client.Eval(script, []string{key}, requestId)
	}
	return r.Val(), r.Err()
}

func Del(key string) (int64, error) {
	var r *redis.IntCmd
	if Model == constants.RedisModel {
		r = ClusterClient.Del(key)
	} else {
		r = Client.Del(key)
	}
	return r.Val(), r.Err()
}

//查询redis值
func GetsPipeline(keys []string) ([]string, error) {
	pipeline := Pipeline()
	defer pipeline.Close()
	//批量进行redis查询
	stringCmds := make([]*redis.StringCmd, len(keys))
	result := make([]string, len(keys))
	for idx, v := range keys {
		stringCmds[idx] = pipeline.Get(v)
	}
	//执行redis查询
	_, e := pipeline.Exec()
	if e != nil && e.Error() != "redis: nil" {
		return nil, e
	}

	for idx, v := range stringCmds {
		if len(v.Val()) > 0 {
			result[idx] = v.Val()
		}
	}

	return result, nil
}

//查询redis值
func HGetsPipeline(keys, fileds []string) ([]string, error) {
	if len(keys) != len(fileds) {
		return nil, errors.New("keys and fileds lenght no equal")
	}
	pipeline := Pipeline()
	defer pipeline.Close()
	//批量进行redis查询
	stringCmds := make([]*redis.StringCmd, len(keys))
	result := make([]string, len(keys))
	for idx, v := range keys {
		stringCmds[idx] = pipeline.HGet(v, fileds[idx])
	}
	//执行redis查询
	_, e := pipeline.Exec()
	if e != nil && e.Error() != "redis: nil" {
		return nil, e
	}

	for idx, v := range stringCmds {
		if len(v.Val()) > 0 {
			result[idx] = v.Val()
		}
	}

	return result, nil
}

//命令根据给定的 sha1 校验码，执行缓存在服务器中的脚本。
func Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	var r *redis.Cmd
	if Model == constants.RedisModel {
		r = ClusterClient.Eval(script, keys, args...)
	} else {
		r = Client.Eval(script, keys, args...)
	}
	if r != nil && r.Err() == nil {
		return r.Val(), nil
	} else {
		return nil, r.Err()
	}
}
