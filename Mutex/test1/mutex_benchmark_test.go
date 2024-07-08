package mutex

import (
	"mutex"
	"sync"
	"testing"
)

/*
goos: linux
goarch: amd64
pkg: mutex
cpu: AMD Ryzen 5 4600H with Radeon Graphics
BenchmarkGoschedMutex-12        83871055                14.58 ns/op
BenchmarkMutex-12                8004612               136.6 ns/op
BenchmarkSpinLock-12            89880213                14.13 ns/op
BenchmarkSyncMutex-12           29304279                47.39 ns/op

BenchmarkGoschedMutex-12        29515287                34.33 ns/op
BenchmarkMutex-12                8752268               126.0 ns/op
BenchmarkSpinLock-12            84994533                13.84 ns/op
BenchmarkSyncMutex-12           24279477                49.00 ns/op

BenchmarkGoschedMutex-12        84296689                14.43 ns/op
BenchmarkMutex-12                8598596               124.2 ns/op
BenchmarkSpinLock-12            91560927                13.88 ns/op
BenchmarkSyncMutex-12           24894915                44.99 ns/op

可以看到我们自己目前实现的锁，跟许多开源的锁相比，性能还差很多。

sync.Mutex作为官方库，什么时候会性能优先呢？前三种锁都是自旋锁，所以在协程竞争不激烈的时候，速度会很快
大概到了10w的级别，sync.Mutex的性能将会远超其他的锁。
cpu: AMD Ryzen 5 4600H with Radeon Graphics
BenchmarkOriginMuteHighConcurrency-12               1066           1123749 ns/op
BenchmarkSpinLockHighConcurrency-12                   72          13957688 ns/op
BenchmarkSyncMutexHighConcurrency-12             1650859               621.3 ns/op
PASS
ok      mutex   70.937s

因此本项目手撕Mutex，将会一边阅读源码，把源码的思想应用进来。
*/

func BenchmarkMutexV1(b *testing.B) {
	mutex := mutex.NewMutexV1()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV2(b *testing.B) {
	mutex := mutex.NewMutexV2()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV3(b *testing.B) {
	mutex := mutex.NewMutexV3()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})
}

func BenchmarkMutexV4(b *testing.B) {
	mutex := mutex.NewMutexV4()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.Lock()
			mutex.Unlock()
		}
	})

}

func BenchmarkSpinLock(b *testing.B) {
	spin := mutex.NewSpinLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}

func BenchmarkSyncMutex(b *testing.B) {
	spin := sync.Mutex{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			spin.Lock()
			//nolint:staticcheck
			spin.Unlock()
		}
	})
}
