package redlock

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"time"
)

var (
	LockTime      = 3 * time.Second               //默认lock时间
	RetryInterval = 100 * time.Millisecond        //默认间隔时间
	Retries       = int(LockTime / RetryInterval) //默认重试次数
)

var ErrCanced = errors.New("canced")

// Mutext定义
type Mutex struct {
	ctx      context.Context
	key      string
	value    string
	tryCount int
	*Client
}

// key name 不能为空
func New(c *Client, key string) *Mutex {
	if key == "" {
		return nil
	}

	return &Mutex{key: key, Client: c, tryCount: Retries, ctx: context.Background()}
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

// 取最大重试次数
func maxRetries(timeout time.Duration) (retries int) {
	retries = Retries
	if r := int(timeout / RetryInterval); r > retries {
		retries = r
	}
	return
}

// 锁
func (m *Mutex) Lock(d ...time.Duration) error {
	to := LockTime
	if len(d) > 0 {
		to = d[0]
		m.tryCount = maxRetries(to)
	}

	value := getValue()
	if m.value == "" {
		m.value = value
	}

	tk := time.NewTicker(RetryInterval)
	for i := 0; i < m.tryCount; i++ {

		b, err := m.SetNX(m.key, m.value, to).Result()
		// 出错直接返回
		if err != nil {
			return err
		}
		// 锁住直接返回
		if b {
			return nil
		}

		// 锁失败的,说明已经被锁住了，就轮训等待
		// 加入context是为了方便外层直接取消
		select {
		case <-m.ctx.Done():
			return ErrCanced
		case <-tk.C:
		}
	}

	return nil
}

// 解锁
func (m *Mutex) Unlock() error {
	return deleteScript.Run(m.Client, []string{m.key}, m.value).Err()
}

// 删除脚本
var deleteScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

// 租约续租script
// TODO 调用
var touchScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("pexpire", KEYS[1], ARGV[2])
	else
		return 0
	end
`)
