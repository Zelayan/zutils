package redis

type ClientOptions struct {
	seconds uint32 // 锁的有效期
	key     string // 锁的 key
	token   string // 谁的锁
	blocked bool   // 是否阻塞获取锁
	timeout uint32 // 阻塞的超时时间
	ttl     uint32 // 看门狗的续约时间
}

type Options func(c *ClientOptions)

func WithKey(key string) Options {
	return func(c *ClientOptions) {
		c.key = key
	}
}

func WithSeconds(seconds uint32) Options {
	return func(c *ClientOptions) {
		c.seconds = seconds
	}
}

func WithToken(token string) Options {
	return func(c *ClientOptions) {
		c.token = token
	}
}

func WithBlocked(blocked bool) Options {
	return func(c *ClientOptions) {
		c.blocked = blocked
	}
}

func WithTimeout(timeout uint32) Options {
	return func(c *ClientOptions) {
		c.timeout = timeout
	}
}

func WithTTL(ttl uint32) Options {
	return func(c *ClientOptions) {
		c.ttl = ttl
	}
}
