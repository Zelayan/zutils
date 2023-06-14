package redis

import (
	"github.com/Zelayan/zutils/lock"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedLock(t *testing.T) {
	check := assert.New(t)
	lock, err := initRedLockClient()
	check.Nil(err)
	err = lock.TryLock()
	check.Nil(err)
}

func TestRedLock_UnLock(t *testing.T) {
	check := assert.New(t)
	lock, err := initRedLockClient()
	check.Nil(err)
	err = lock.UnLock()
	check.Nil(err)
}

func initRedLockClient() (lock.Lock, error) {
	var clients []*redis.Client
	clients = append(clients, redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "qwer1234",
	}))
	clients = append(clients, redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "qwer1234",
	}))
	clients = append(clients, redis.NewClient(&redis.Options{
		Addr:     "localhost:6381",
		Password: "qwer1234",
	}))
	clients = append(clients, redis.NewClient(&redis.Options{
		Addr:     "localhost:6382",
		Password: "qwer1234",
	}))
	lock, err := NewRedLock(clients, WithKey("ch"))
	return lock, err
}
