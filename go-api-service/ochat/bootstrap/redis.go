package bootstrap

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	redis_init_once sync.Once
	RedisClient     *redis.Client
	RedisContext    context.Context
)

// redisList config struct type
type redisListConfT struct {
	Client01 redisConfT `yaml:"client_01"`
}

// redis config struct type
type redisConfT struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Auth string `yaml:"auth"`
	Db   int    `yaml:"db"`
}

func RedisOnceInit() *redis.Client {

	redis_init_once.Do(initRedis)

	return RedisClient
}

// read  cache.yaml config and set var
func loadRedisListConfig() redisListConfT {
	return readYamlConfig[redisListConfT]("redis.yaml")
}

func initRedis() {
	// loading redis config info
	redisList := loadRedisListConfig()

	// fmt.Printf("%#v\n", redisList)
	db01 := redisList.Client01
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

	log.Println("init redis success!")
}
