package mutex

import (
	"sync"
	"testing"
)

// 该文件为mutex的功能测试文件
const (
	numGoroutines = 100
	numIncrements = 100
)

func TestMutexV1(t *testing.T) {

	var counter int32
	m := NewMutexV1()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				m.Lock()
				counter++
				m.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * numIncrements)
	if counter != expected {
		t.Errorf("expected counter to be %d, but got %d", expected, counter)
	}
}

func TestMutexV3(t *testing.T) {

	var counter int32
	m := NewMutexV3()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				m.Lock()
				counter++
				m.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * numIncrements)
	if counter != expected {
		t.Errorf("expected counter to be %d, but got %d", expected, counter)
	}
}

func TestMutexV2(t *testing.T) {

	var counter int32
	m := NewMutexV2()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				m.Lock()
				counter++
				m.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * numIncrements)
	if counter != expected {
		t.Errorf("expected counter to be %d, but got %d", expected, counter)
	}
}

func TestSpinLock(t *testing.T) {

	var counter int32
	m := NewSpinLock()
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				m.Lock()
				counter++
				m.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * numIncrements)
	if counter != expected {
		t.Errorf("expected counter to be %d, but got %d", expected, counter)
	}
}

func TestSyncMutex(t *testing.T) {

	var counter int32
	m := sync.Mutex{}
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIncrements; j++ {
				m.Lock()
				counter++
				m.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := int32(numGoroutines * numIncrements)
	if counter != expected {
		t.Errorf("expected counter to be %d, but got %d", expected, counter)
	}
}
