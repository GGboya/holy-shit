package mutex

type MutexV1 struct {
	ch chan struct{}
}

func NewMutexV1() *MutexV2 {
	return &MutexV2{ch: make(chan struct{}, 1)} // 缓冲通道
}

func (m *MutexV1) Lock() {
	// 尝试获取锁

	<-m.ch
}

func (m *MutexV1) Unlock() {

	m.ch <- struct{}{}
}
