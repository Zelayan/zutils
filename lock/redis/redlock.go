package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zelayan/zutils/lock"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var ErrLocked = errors.New("red locked failed")

type RedLock struct {
	ClientOptions
	successStores []*redis.Client // 锁成功的 redis 节点
	stores        []*redis.Client
}

func (r *RedLock) TryLock() error {
	return r.TryLockCtx(context.Background())
}

func (r *RedLock) UnLock() error {
	return r.UnLockCtx(context.Background())
}

func (r *RedLock) TryLockCtx(ctx context.Context) error {
	var wg sync.WaitGroup
	lockSuccessClient := make(chan *redis.Client, len(r.stores))
	wg.Add(len(r.stores))
	for _, client := range r.stores {
		go func(redis *redis.Client) {
			defer wg.Done()
			nx := redis.SetNX(ctx, r.key, r.token, time.Duration(r.seconds)*time.Second)
			if nx.Err() != nil {
				fmt.Println(nx.Err())
				return
			}
			if nx.Val() {
				fmt.Println("add lock success")
				lockSuccessClient <- client
			} else {
				fmt.Println("add lock failed")
			}
		}(client)
	}
	wg.Wait()
	close(lockSuccessClient)
	if len(lockSuccessClient) < len(lockSuccessClient)/2+1 {
		fmt.Println("red lock failed, recovery")
		var wgc sync.WaitGroup
		wgc.Add(len(lockSuccessClient))
		// 如果没有加锁成功，需要把已经加锁的给释放掉
		for client := range lockSuccessClient {
			go func(redis *redis.Client) {
				defer wgc.Done()
				redis.Eval(ctx, delScript, []string{r.key}, []string{r.token})
			}(client)
		}
		wgc.Wait()
		return ErrLocked
	}
	// TODO watchDog
	return nil

}

func (r *RedLock) UnLockCtx(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(len(r.stores))
	for _, store := range r.stores {
		go func(client *redis.Client) {
			defer wg.Done()
			client.Eval(ctx, delScript, []string{r.key}, []string{r.token})
		}(store)
	}
	wg.Wait()
	return nil
}

func (r RedLock) Lock(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewRedLock(clients []*redis.Client, opts ...Options) (lock.Lock, error) {
	if len(clients) < 3 {
		return nil, errors.New("red lock len must > 3")
	}
	l := &RedLock{
		stores: clients,
	}
	for _, opt := range opts {
		opt(&l.ClientOptions)
	}
	return l, nil
}
