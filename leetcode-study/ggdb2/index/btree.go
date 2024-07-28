package index

import (
	"leetcode/ggdb/data"
	"sync"

	syncx "leetcode/ggdb/internal/sync"

	"github.com/google/btree"
)

// BTree 主要封装了 google 的 btree 库
type BTree struct {
	tree *btree.BTree // 基于B树的索引
	lock sync.Locker  // 确保对B树的并发访问的线程安全
}

// NewBTree 新建一个B树结构体
func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),       // 设置结点的最大键数，可以根据需要调整
		lock: syncx.NewSpinLock(), // 对B树的访问进行并发控制
	}
}

func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) bool {
	it := &Item{key: key, pos: pos}
	bt.lock.Lock()
	defer bt.lock.Unlock()
	bt.tree.ReplaceOrInsert(it)
	return true
}

func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := &Item{key: key}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos
}

func (bt *BTree) Delete(key []byte) bool {
	it := &Item{key: key}
	bt.lock.Lock()
	defer bt.lock.Unlock()
	oldItem := bt.tree.Delete(it)
	if oldItem == nil {
		return false
	}
	return true
}

func (bt *BTree) Iterator() [][]byte {
	// 获取b树中所有的key
	var keys [][]byte
	bt.tree.Ascend(func(item btree.Item) bool {
		myItem := item.(*Item)
		keys = append(keys, myItem.key)
		return true
	})
	return keys
}
