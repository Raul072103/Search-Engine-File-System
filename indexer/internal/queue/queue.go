package queue

import (
	"MyFileExporer/common/models"
	"errors"
	"sync"
)

const (
	InsertEvent     = "INSERT"
	UpdateEvent     = "UPDATE"
	DeleteEvent     = "DELETE"
	RecursiveDelete = "RECURSIVE_DELETE"
)

var (
	ErrQueueEmpty = errors.New("there are no queue in the queue at the moment")
)

type DBEvent struct {
	Type string
	File models.File
}

func (d *DBEvent) IsInsert() bool {
	return d.Type == InsertEvent
}

func (d *DBEvent) IsUpdate() bool {
	return d.Type == UpdateEvent
}

func (d *DBEvent) IsDelete() bool {
	return d.Type == DeleteEvent
}

func (d *DBEvent) IsRecursiveDelete() bool {
	return d.Type == RecursiveDelete
}

type InMemoryQueue struct {
	events []DBEvent
	lock   sync.Mutex
}

func NewQueue() *InMemoryQueue {
	return &InMemoryQueue{
		events: make([]DBEvent, 0),
		lock:   sync.Mutex{},
	}
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

func (q *InMemoryQueue) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.events)
}
