package redlock

import (
	"fmt"
	redis "github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
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
