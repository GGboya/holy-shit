package mutex

import (
	"sync"
	"testing"
)

func BenchmarkMutexV1_100(b *testing.B) {
	mutex := NewMutexV1()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV2_100(b *testing.B) {
	mutex := NewMutexV2()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkSpinLock_100(b *testing.B) {
	spin := NewSpinLock()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}

func BenchmarkSyncMutex_100(b *testing.B) {
	spin := sync.Mutex{}
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}
