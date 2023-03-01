package bootstrap

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	_redisInitOnce sync.Once
	RedisClient    *redis.Client
	RedisContext   context.Context
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

	_redisInitOnce.Do(initRedis)

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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log.Println("init redis success!")
}
