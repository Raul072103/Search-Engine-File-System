package queue

import (
	"MyFileExporer/common/models"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func testQueue(queue *InMemoryQueue, dbevent DBEvent, wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()
	sleepTime := rand.Intn(3000)

	queue.Push(dbevent)

	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	_, err := queue.Pop()
	if err != nil {
		t.Errorf("Failed to pop event: %v", err)
		return
	}
}

func TestInMemoryQueue(t *testing.T) {
	// Create the InMemoryQueue for testing
	queue := &InMemoryQueue{}

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Simulate creating and testing multiple queue concurrently
	for i := 0; i < 10; i++ {
		// Each goroutine gets a different event
		dbevent := DBEvent{
			Type: "update",
			File: models.File{ID: int64(i)},
		}

		// Increment the wait group counter
		wg.Add(1)

		// Launch a goroutine for each event
		go testQueue(queue, dbevent, &wg, t)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
