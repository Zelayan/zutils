package bloom

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"hash/fnv"
	"testing"
)

func TestBloomFilter_Add(t *testing.T) {
	check := assert.New(t)
	bloom := NewBoomFilter(1000, 3)
	bloom.Add([]byte("hello"))
	bloom.Add([]byte("hello2"))
	bloom.Add([]byte("hello3"))

	check.True(bloom.Contains([]byte("hello")))
	check.True(bloom.Contains([]byte("hello2")))
	check.False(bloom.Contains([]byte("hello4")))
}

func Test_fnc(t *testing.T) {
	data := []byte("example")

	hashFunc := fnv.New64a()
	_, err := hashFunc.Write(data)
	if err != nil {
		return
	}
	hashValue := hashFunc.Sum64()

	fmt.Printf("Hash value: %d\n", hashValue%10)
}
