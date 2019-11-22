package redlock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/satori/go.uuid"
	"time"
)

var ErrCanced = errors.New("canced")

type Mutex struct {
	ctx      context.Context
	key      string
	tryCount int
	*Client
}

// key name 不能为空
func New(c *Client, key string) *Mutex {
	if key == "" {
		return nil
	}

	return &Mutex{key: key, Client: c, tryCount: 16}
}

// value使用uuidv4
// TODO:有没有更好的算法
func getValue() string {
	u, err := uuid.NewV4()
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return u.String()
}

// 设置context
func (m *Mutex) WithContext(ctx context.Context) *Mutex {
	m.ctx = ctx
	return m
}

// 锁
func (m *Mutex) Lock(d ...time.Duration) error {
	/*
		to := 3 * time.Second
		if len(d) > 0 {
			to = d[0]
		}

		value := getValue()

		if err == nil {
			return nil
		}

		tk := time.NewTicker(d)
		for i := 0; i < m.tryCount; i++ {

			set, err := m.SetNx(m.key, value, to).Result()
			// 锁住直接返回
			if err == nil {
				return nil
			}

			// 锁失败的,说明已经被锁住了，就轮训等待
			select {
			case <-m.ctx.Done():
				return ErrCanced
			case <-tk.C:
			}
		}

	*/
	return nil
}

// 解锁
func (m *Mutex) Unlock() error {
	return nil
}

var deleteScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

var touchScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`)