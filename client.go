package redlock

import (
	redis "github.com/go-redis/redis/v7"
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
