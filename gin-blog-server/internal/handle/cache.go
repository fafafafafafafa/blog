package handle

import (
	"context"
	g "gin-blog/internal/global"

	"github.com/redis/go-redis/v9"
)

// redis context
var rctx = context.Background()

// Config

// 将博客配置缓存到 Redis 中
func addConfigCache(rdb *redis.Client, config map[string]string) error {
	return rdb.HMSet(rctx, g.CONFIG, config).Err()
}

// 从 Redis 中获取博客配置缓存
// rdb.HGetAll 如果不存在 key, 不会返回 redis.Nil 错误, 而是返回空 map
func getConfigCache(rdb *redis.Client) (cache map[string]string, err error) {
	return rdb.HGetAll(rctx, g.CONFIG).Result()
}

// 删除 Redis 中博客配置缓存
func removeConfigCache(rdb *redis.Client) error {
	return rdb.Del(rctx, g.CONFIG).Err()
}
