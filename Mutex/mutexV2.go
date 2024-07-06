package mutex

import (
	"sync/atomic"
	"syscall"
	"unsafe"
)

// 我们在之前的基准测试中已经发现，当协程数量不断增加的时候，sync.Mutex的性能将远大于自旋锁，这是因为CPU资源是有限的。
// 所以我们后续的Mutex的实现，将基于互斥锁实现。就是阻塞，唤醒的方式

type MutexV2 uint32

func NewMutexV2() *MutexV2 {
	return new(MutexV2)
}
func (m *MutexV2) Lock() {
	for {
		if atomic.CompareAndSwapUint32((*uint32)(m), 0, 1) {
			return
		}
		// 没能上锁，就得阻塞
		syscall.Syscall(syscall.SYS_FUTEX, uintptr(unsafe.Pointer((*uint32)(m))), 0, 0)
	}
}

func (m *MutexV2) Unlock() {
	atomic.StoreUint32((*uint32)(m), 0)
	syscall.Syscall(syscall.SYS_FUTEX, uintptr(unsafe.Pointer((*uint32)(m))), 1, 0)
}
