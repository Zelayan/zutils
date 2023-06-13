package lock

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"strconv"
)

var (
	lockScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

type RedisLock struct {
	store   *redis.Client
	seconds uint32
	key     string
	id      string
}

func (r *RedisLock) TryLock() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisLock) UnLock() (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisLock) TryLockCtx(ctx context.Context) (bool, error) {
	script := redis.NewScript(lockScript)
	seconds := strconv.Itoa(int(r.seconds)) // 将 seconds 转换为字符串类型

	eval := script.Run(ctx, r.store, []string{r.key}, []string{r.id, seconds})
	if eval.Err() == redis.Nil {
		return false, nil
	} else if eval.Err() != nil {
		return false, eval.Err()
	} else if eval.Val() == nil {
		return false, nil
	}
	if eval.Val() != "OK" {

	}
	return true, nil
}

func (r *RedisLock) UnLockCtx(ctx context.Context) (bool, error) {
	eval := r.store.Eval(ctx, delScript, []string{r.key}, []string{r.id})
	if eval.Err() != nil {
		return false, eval.Err()
	}
	return true, nil
}

func NewRedisLock(redis *redis.Client) Lock {
	return &RedisLock{
		store:   redis,
		seconds: 1000,
		key:     "12312312321",
		id:      "2222",
	}
}
