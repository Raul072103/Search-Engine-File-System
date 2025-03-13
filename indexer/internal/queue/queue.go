package events

import (
	"MyFileExporer/common/models"
	"errors"
	"sync"
)

const (
	InsertEvent = "INSERT"
	UpdateEvent = "UPDATE"
	DeleteEvent = "DELETE"
)

var (
	ErrQueueEmpty = errors.New("there are no events in the queue at the moment")
)

type DBEvent struct {
	Type string
	File models.File
}

type InMemoryQueue struct {
	events []DBEvent
	lock   sync.Mutex
}

func (q *InMemoryQueue) Push(event DBEvent) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.events = append(q.events, event)
}

func (q *InMemoryQueue) Pop() (DBEvent, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.events) == 0 {
		return DBEvent{}, ErrQueueEmpty
	}
	event := q.events[0]
	q.events = q.events[1:]
	return event, nil
}
