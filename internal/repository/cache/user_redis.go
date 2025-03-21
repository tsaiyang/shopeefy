package cache

import "github.com/redis/go-redis/v9"

type redisUserCache struct {
	client redis.Cmdable
}

var _ UserCache = (*redisUserCache)(nil)

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &redisUserCache{client: client}
}
