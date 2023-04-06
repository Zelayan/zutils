package zmap

// ChMap 基于channel实现的map
type ChMap struct {
	data  map[string]interface{}
	mutex chan struct{}
}

func NewChMap() *ChMap {
	return &ChMap{
		data:  make(map[string]interface{}),
		mutex: make(chan struct{}, 1),
	}
}

func (c *ChMap) Load(key string) (any, bool) {
	c.mutex <- struct{}{}
	defer func() { <-c.mutex }()
	value, ok := c.data[key]
	return value, ok
}

func (c *ChMap) Store(key string, value any) {
	c.mutex <- struct{}{}
	defer func() { <-c.mutex }()
	c.data[key] = value
}

func (c *ChMap) Delete(key string) {
	c.mutex <- struct{}{}
	defer func() { <-c.mutex }()
	delete(c.data, key)
}
