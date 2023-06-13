package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisLock_UnLockCtx(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "zeze.com:6379",
		Password: "qwer1234",
	})
	lock := NewRedisLock(client, WithKey("ttt"), WithSeconds(100))
	err := lock.TryLockCtx(context.Background())
	if err != nil {
		t.Logf("lock failed: %s", err)
	}
	err = lock.UnLockCtx(context.Background())
	if err != nil {
		panic(err)
	}

}

func Test_t(t *testing.T) {
	check := assert.New(t)
	client := redis.NewClient(&redis.Options{
		Addr:     "zeze.com:6379",
		Password: "qwer1234",
	})
	redisLock := NewRedisLock(client, WithKey("ttt"), WithSeconds(1000), WithBlocked(true), WithTimeout(10))
	redisLock.UnLock()
	err := redisLock.Lock(context.Background())
	check.Nil(err)
	err = redisLock.Lock(context.Background())
	check.Equal(ErrLockBlockedTimeOut, err)
}

func TestRedisLock_UnLock(t *testing.T) {
	check := assert.New(t)
	client := redis.NewClient(&redis.Options{
		Addr:     "zeze.com:6379",
		Password: "qwer1234",
	})
	redisLock := NewRedisLock(client, WithKey("ttt"), WithSeconds(1000), WithBlocked(true))
	err := redisLock.UnLock()
	check.Nil(err)
	err = redisLock.UnLock()
	check.Nil(err)

}
