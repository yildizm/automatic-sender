package redis

import (
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
    "golang.org/x/net/context"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis(redisURL string) {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        panic(err)
    }

    rdb = redis.NewClient(opt)
}

func CacheMessageID(id int, messageID string, sentAt time.Time) error {
    key := fmt.Sprintf("message:%d", id)
    value := fmt.Sprintf("messageID:%s, sentAt:%s", messageID, sentAt.Format(time.RFC3339))
    return rdb.Set(ctx, key, value, 0).Err()
}
