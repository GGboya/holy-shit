package mutex

import (
	"runtime"
	"sync/atomic"
)

type MutexV4 struct {
	state int32
	ch    chan struct{}
	sema  uint32
}

func (m *MutexV4) Lock() {
	// Fast path: 幸运case，能够直接获取到锁
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}
	awoke := false
	backoff := 1
	for {
		old := m.state
		new := old | mutexLocked // 新状态加锁
		if old&mutexLocked != 0 {

			for i := 0; i < backoff; i++ {
				if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
					atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
					awoke = true
				}
				runtime.Gosched()
			}
			if backoff < maxBackoff {
				backoff <<= 1
			}

			new = old + 1<<mutexWaiterShift //等待者数量加一
		}
		if awoke {
			// goroutine是被唤醒的，
			// 新状态清除唤醒标志
			new &^= mutexWoken

		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) { //设置新状态
			if old&mutexLocked == 0 { // 锁原状态未加锁
				break
			}
			// 阻塞
			<-m.ch
			awoke = true
		}
	}
}

func (m *MutexV4) Unlock() {
	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked) //去掉锁标志
	if (new+mutexLocked)&mutexLocked == 0 {        //本来就没有加锁
		panic("sync: unlock of unlocked mutex")
	}

	old := new
	for {
		if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken) != 0 { // 没有等待者，或者有唤醒的waiter，或者锁原来已加锁
			return
		}
		new = (old - 1<<mutexWaiterShift) | mutexWoken // 新状态，准备唤醒goroutine，并设置唤醒标志
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			m.ch <- struct{}{}

			return
		}
		old = m.state
	}
}

func NewMutexV4() *MutexV4 {
	return &MutexV4{
		state: 0,
		ch:    make(chan struct{}, 1),
	}
}
