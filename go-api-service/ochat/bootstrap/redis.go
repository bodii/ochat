package bootstrap

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
)

var (
	redis_init_once sync.Once
	RedisClient     *redis.Client
	RedisContext    context.Context
)

// redisList config struct type
type redisListConfT struct {
	Servs []redisConfT `toml:"redis-server"`
}

// redis config struct type
type redisConfT struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	Auth string `toml:"auth"`
	Db   int    `toml:"db"`
}

func RedisOnceInit() *redis.Client {

	redis_init_once.Do(initRedis)

	return RedisClient
}

// read  cache.yaml config and set var
func loadRedisListConfig() redisListConfT {
	return readTomlConfig[redisListConfT]("redis.toml")
}

func initRedis() {
	// loading redis config info
	redisList := loadRedisListConfig()

	// fmt.Printf("%#v\n", redisList)
	db01 := redisList.Servs[0]
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", db01.Host, db01.Port),
		Password: db01.Auth,
		DB:       db01.Db,
	})

	RedisContext := context.Background()
	ping, err := RedisClient.Ping(RedisContext).Result()
	if err != nil || ping != "PONG" {
		RedisClient.Close()
		panic(err)
	}

	slog.Info("init redis success!")
}
