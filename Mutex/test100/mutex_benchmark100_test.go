package mutex

import (
	"mutex"
	"sync"
	"testing"
)

func BenchmarkMutexV1_100(b *testing.B) {
	mutex := mutex.NewMutexV1()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV2_100(b *testing.B) {
	mutex := mutex.NewMutexV2()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV3_100(b *testing.B) {
	mutex := mutex.NewMutexV3()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV4_100(b *testing.B) {
	mutex := mutex.NewMutexV4()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkSpinLock_100(b *testing.B) {
	spin := mutex.NewSpinLock()
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
