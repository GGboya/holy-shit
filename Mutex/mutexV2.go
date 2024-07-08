package mutex

import (
	"sync/atomic"
)

type MutexV2 struct {
	state int32
	ch    chan struct{}
}

func NewMutexV2() *MutexV2 {
	return &MutexV2{ch: make(chan struct{}, 1)} // 缓冲通道
}

func (m *MutexV2) Lock() {
	// 尝试获取锁
	if atomic.AddInt32(&m.state, 1) == 1 {
		return
	}
	// 等待锁释放
	<-m.ch
}

func (m *MutexV2) Unlock() {
	if atomic.AddInt32(&m.state, -1) == 0 {
		return
	}
	// 尝试唤醒一个等待者
	m.ch <- struct{}{}
}
