package bloom

import (
	"hash"
	"hash/fnv"
)

type Bloom interface {
	Add(key []byte)
	Contains(key []byte) bool
	Reserve()
}

type BloomFilter struct {
	bitArray []bool
	hashFunc []hash.Hash64
}

func NewBoomFilter(bitArrayLen int, hashLen int) *BloomFilter {
	hashFunc := make([]hash.Hash64, hashLen)

	for i := 0; i < hashLen; i++ {
		hashFunc[i] = fnv.New64a()
	}

	bl := &BloomFilter{
		bitArray: make([]bool, bitArrayLen),
		hashFunc: hashFunc,
	}
	return bl
}

// Add 通过三次 hash 计算其下标
func (b *BloomFilter) Add(key []byte) {
	for _, hs := range b.hashFunc {
		sum64 := b.hash(key, hs)
		u := sum64 % (uint64)(len(b.bitArray))
		b.bitArray[u] = true
	}
}

// Contains 通过三次 hash 计算其下标，如果是 true 则表示包含
func (b *BloomFilter) Contains(key []byte) bool {
	for _, hs := range b.hashFunc {
		sum64 := b.hash(key, hs)
		u := sum64 % (uint64)(len(b.bitArray))
		if ok := b.bitArray[u]; !ok {
			return false
		}
	}
	return true
}

// hash 计算 hash 值
func (b *BloomFilter) hash(key []byte, hs hash.Hash64) uint64 {
	hs.Reset()
	_, err := hs.Write(key)
	if err != nil {
		return 0
	}
	sum64 := hs.Sum64()
	return sum64
}
