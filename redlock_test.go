package redlock

import (
	redis "github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var count int

const number = 100000

func Test_Redlock_lock_unlock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()
	assert.NoError(t, err)
	//	fmt.Println(pong, err)

	var wg sync.WaitGroup
	wg.Add(2)
	defer func() {
		wg.Wait()
		assert.Equal(t, count, number*2)
	}()

	go func() {
		defer wg.Done()
		m := NewClient(client).NewMutex("123")
		for i := 0; i < number; i++ {
			m.Lock()
			count++
			m.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		m := NewClient(client).NewMutex("123")
		for i := 0; i < number; i++ {
			m.Lock()
			count++
			m.Unlock()
		}
	}()
}
