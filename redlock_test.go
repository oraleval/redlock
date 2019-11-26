package redlock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

var count int

const number = 100000

// 测试错误的情况
func Test_Redlock_Fail(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "不存在", //一个不存在的密码
		DB:       0,
	})

	err := NewClient(client).NewMutex("aaa").Lock()
	assert.Error(t, err)
}

type retries struct {
	to       time.Duration
	tryCount int
}

func Test_Redlock_Withcontext(t *testing.T) {
	log.SetFlags(log.Lmicroseconds)
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
		// 测试的时候自己配置
	})

	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		m := NewClient(client).NewMutex("aaa").WithContext(ctx)
		log.Printf("A:start lock\n")

		err := m.Lock()

		log.Printf("A:lock return:%v\n", err)
		time.Sleep(time.Second * 5)

		m.Unlock()

	}()

	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		m := NewClient(client).NewMutex("aaa").WithContext(ctx)
		log.Printf("B:start lock\n")

		err := m.Lock()

		log.Printf("B:lock return:%v\n", err)
		time.Sleep(time.Second * 5)

		m.Unlock()

	}()
}

func Test_Redlock_maxRetries(t *testing.T) {
	tests := []retries{
		{time.Second, Retries},
		{3 * time.Second, Retries},
		{10 * time.Second, int((10 * time.Second) / RetryInterval)},
	}

	for _, v := range tests {
		assert.Equal(t, maxRetries(v.to), v.tryCount)
	}
}

func Test_Redlock_lock_unlock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	defer func() {
		wg.Wait()
		assert.Equal(t, count, number*2)
	}()

	go func() {
		defer wg.Done()
		m := NewClient(client).NewMutex("testkey")
		for i := 0; i < number; i++ {
			func() {
				defer func() {
					err := m.Unlock()
					assert.NoError(t, err)
					if err != nil {
						fmt.Printf("unlock fail:%s\n", err)
						return
					}
				}()

				err := m.Lock()
				assert.NoError(t, err)
				if err != nil {
					fmt.Printf("lock fail:%s\n", err)
					return
				}

				count++
			}()
		}
	}()

	go func() {
		defer wg.Done()
		m := NewClient(client).NewMutex("testkey")
		for i := 0; i < number; i++ {
			func() {
				defer func() {
					err = m.Unlock()
					assert.NoError(t, err)
					if err != nil {
						fmt.Printf("unlock fail:%s\n", err)
						return
					}
				}()

				err := m.Lock()
				assert.NoError(t, err)
				if err != nil {
					fmt.Printf("lock fail:%s\n", err)
					return
				}

				count++
			}()
		}
	}()
}
