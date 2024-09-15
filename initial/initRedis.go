package initial

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func init() {
	loadEnv()
}

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:" + os.Getenv("RedisPort"),
	})
}
