package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zelayan/zutils/lock"
	redis "github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var (
	// 基于 lua 脚本实现，如果 token 是自己的就删除对应的 key
	delScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	// deplay
	delayScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("expire", KEYS[1], ARGV[2])
else
	return 0
end
`
)

var (
	ErrLockAcquiredByOthers = errors.New("set key failed")
	ErrLockBlockedTimeOut   = errors.New("lock failed, blocked waiting timeout")
)

type RedisLock struct {
	store *redis.Client
	ClientOptions
	watchDog chan struct{}
}

func (r *RedisLock) Lock(ctx context.Context) error {
	// 先尝试加锁
	err := r.TryLockCtx(ctx)
	if err == nil {
		return nil
	}
	if !r.blocked {
		return err
	}

	if !errors.Is(err, ErrLockAcquiredByOthers) {
		return err
	}

	// 阻塞加锁
	err = r.blockedLock(ctx)
	return err
}

func (r *RedisLock) TryLock() error {
	return r.TryLockCtx(context.Background())
}

func (r *RedisLock) UnLock() error {
	return r.UnLockCtx(context.Background())
}

func (r *RedisLock) TryLockCtx(ctx context.Context) error {
	nx := r.store.SetNX(ctx, r.key, r.token, time.Duration(r.seconds)*time.Second)
	if nx.Err() != nil {
		return nx.Err()
	}
	if !nx.Val() {
		return ErrLockAcquiredByOthers
	}
	go r.startWatchDog()
	return nil
}

func (r *RedisLock) UnLockCtx(ctx context.Context) error {
	eval := r.store.Eval(ctx, delScript, []string{r.key}, []string{r.token})
	close(r.watchDog)
	if eval.Err() != nil {
		return eval.Err()
	}
	return nil
}

func (r *RedisLock) blockedLock(ctx context.Context) error {
	after := time.After(time.Duration(r.timeout) * time.Second)
	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("lock failed, ctx timeout, err:%w", ctx.Err())
		case <-ticker.C:
			err := r.TryLockCtx(ctx)
			if !errors.Is(err, ErrLockAcquiredByOthers) {
				return err
			}
		case <-after:
			return ErrLockBlockedTimeOut
		}
	}
}

func (r *RedisLock) startWatchDog() {
	if r.ttl == 0 {
		r.ttl = 10
	}
	r.watchDog = make(chan struct{})
	ticker := time.NewTicker(time.Duration(r.ttl/3) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			timeout, cancelFunc := context.WithTimeout(context.Background(), time.Duration(r.ttl/3*2)*time.Second)
			eval := r.store.Eval(timeout, delayScript, []string{r.key}, []string{r.token, strconv.Itoa(int(r.ttl))})
			cancelFunc()
			fmt.Println("watch dog")
			if eval.Err() != nil {
				return
			}
		case <-r.watchDog:
			return
		}
	}
}

func NewRedisLock(redis *redis.Client, Ops ...Options) lock.Lock {
	redisLock := RedisLock{
		store:    redis,
		watchDog: make(chan struct{}),
	}
	for _, opt := range Ops {
		opt(&redisLock.ClientOptions)
	}
	fmt.Println(redisLock.key)
	return &redisLock
}
