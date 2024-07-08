package mutex

import (
	"mutex"
	"sync"
	"testing"
)

var t = 50000

func BenchmarkMutexV1(b *testing.B) {
	mutex := mutex.NewMutexV1()
	b.SetParallelism(t)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV2(b *testing.B) {
	mutex := mutex.NewMutexV2()
	b.SetParallelism(t)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkSpinLock(b *testing.B) {
	spin := mutex.NewSpinLock()
	b.SetParallelism(t)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			spin.Unlock()
		}
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	spin := sync.Mutex{}
	b.SetParallelism(t)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			spin.Unlock()
		}
	})
}
