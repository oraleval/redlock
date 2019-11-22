package redlock

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func Test_New(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}
