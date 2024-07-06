package mutex

import (
	"runtime"
	"sync/atomic"
)

type MutexV1 struct {
	lock uint32 // 互斥锁标识，0表示未锁定，1表示锁定
}

func (m *MutexV1) Lock() {
	for !atomic.CompareAndSwapUint32(&m.lock, 0, 1) {
		// 如果不能从0变为1，则一直等待这把锁
		// 如果进入到了循环，表示可以加锁
		runtime.Gosched() // 让出cpu
	}
}

func (m *MutexV1) Unlock() {
	atomic.StoreUint32(&m.lock, 0) // 解锁
}

func NewMutexV1() *MutexV1 {
	return &MutexV1{}
}
