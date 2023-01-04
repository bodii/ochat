package service

import (
	"context"
	"ochat/bootstrap"

	"github.com/go-redis/redis/v8"
)

var REDIS_CTX context.Context = bootstrap.RedisContext

func NewRedis() *redis.Client {
	return bootstrap.RedisClient
}
