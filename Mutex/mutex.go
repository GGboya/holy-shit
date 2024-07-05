package mutex

import "sync/atomic"

type Mutex struct {
	lock uint32 // 互斥锁标识，0表示未锁定，1表示锁定
}

func (m *Mutex) Lock() {
	for !atomic.CompareAndSwapUint32(&m.lock, 0, 1) {
		// 如果不能从0变为1，则一直等待这把锁
		// 如果进入到了循环，表示可以加锁
	}
}

func (m *Mutex) Unlock() {
	atomic.StoreUint32(&m.lock, 0) // 解锁
}

func NewMutex() *Mutex {
	return &Mutex{}
}
