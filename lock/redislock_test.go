package lock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestRedisLock_UnLockCtx(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "zeze.com:6379",
		Password: "qwer1234",
	})
	lock := NewRedisLock(client)
	_, err := lock.TryLockCtx(context.Background())
	if err != nil {
		panic(err)
	}
	client.Set(context.Background(), "xx", "", redis.KeepTTL)
}

func Test_t(t *testing.T) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "zeze.com:6379",
		Password: "qwer1234",
	})
	_ = rdb.FlushDB(ctx).Err()

	fmt.Printf("# INCR BY\n")
	for _, change := range []int{+1, +5, 0} {
		num, err := incrBy.Run(ctx, rdb, []string{"my_counter"}, change).Int()
		if err != nil {
			panic(err)
		}
		fmt.Printf("incr by %d: %d\n", change, num)
	}

	fmt.Printf("\n# SUM\n")
	sum, err := sum.Run(ctx, rdb, []string{"my_sum"}, 1, 2, 3).Int()
	if err != nil {
		panic(err)
	}
	fmt.Printf("sum is: %d\n", sum)
}

var incrBy = redis.NewScript(`
local key = KEYS[1]
local change = ARGV[1]

local value = redis.call("GET", key)
if not value then
  value = 0
end

value = value + change
redis.call("SET", key, value)

return value
`)

var sum = redis.NewScript(`
local key = KEYS[1]

local sum = redis.call("GET", key)
if not sum then
  sum = 0
end

local num_arg = #ARGV
for i = 1, num_arg do
  sum = sum + ARGV[i]
end

redis.call("SET", key, sum)

return sum
`)
