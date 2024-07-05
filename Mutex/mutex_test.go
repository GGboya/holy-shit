package mutex

import (
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	const numGoroutines = 1000
	const numIncrements = 1000

	var counter int32
	m := NewMutex()
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
	const numGoroutines = 1000
	const numIncrements = 1000

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
