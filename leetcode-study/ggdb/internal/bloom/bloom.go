package bloom

import (
	"hash/fnv"
	syncx "leetcode/ggdb/internal/sync"
	"sync"
)

// BloomFilter 结构定义
type BloomFilter struct {
	size      uint
	hashFuncs uint
	bitset    []bool
	mutex     sync.Locker
}

// NewBloomFilter 创建一个新的布隆过滤器
func NewBloomFilter(size uint, hashFuncs uint) *BloomFilter {
	return &BloomFilter{
		size:      size,
		hashFuncs: hashFuncs,
		bitset:    make([]bool, size),
		mutex:     syncx.NewSpinLock(),
	}
}

// Add 添加一个元素到布隆过滤器
func (bf *BloomFilter) Add(key []byte) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := uint(0); i < bf.hashFuncs; i++ {
		index := bf.hash(key, i) % bf.size
		bf.bitset[index] = true
	}
}

// Contains 检查元素是否可能存在于布隆过滤器
func (bf *BloomFilter) Contains(key []byte) bool {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := uint(0); i < bf.hashFuncs; i++ {
		index := bf.hash(key, i) % bf.size
		if !bf.bitset[index] {
			return false
		}
	}

	return true
}

// hash 计算哈希值
func (bf *BloomFilter) hash(key []byte, seed uint) uint {
	hash := fnv.New64a()
	hash.Write(key)
	hash.Write([]byte{byte(seed)})
	return uint(hash.Sum64())
}
