package redlock

import (
	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func NewClient(c *redis.Client) *Client {
	return &Client{Client: c}
}

func (c *Client) NewMutex(key string) *Mutex {
	return New(c, key)
}
