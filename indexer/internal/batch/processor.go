package batch

import (
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/database"
	"context"
	"errors"
	"time"
)

type Processor interface {
	Run(ctx context.Context) error
	PopEvent() (queue.DBEvent, error)
	HandleEvent(event queue.DBEvent) error
	ExecuteInsertEvent(event queue.DBEvent) error
	ExecuteUpdateEvent(event queue.DBEvent) error
	ExecuteDeleteEvent(event queue.DBEvent) error
}

type processor struct {
	DBRepo     database.Repo
	EventQueue *queue.InMemoryQueue
}

const (
	ProcessorSleepTime = time.Second * 5
)

var (
	ErrExpectedInsertEvent = errors.New("expected the event to of type INSERT")
	ErrExpectedUpdateEvent = errors.New("expected the event to of type UPDATE")
	ErrExpectedDeleteEvent = errors.New("expected the event to of type DELETE")
	ErrUnknownEvent        = errors.New("unknown event sent for processing")
)

func NewProcessor(dbRepo database.Repo, eventQueue *queue.InMemoryQueue) Processor {
	return &processor{
		DBRepo:     dbRepo,
		EventQueue: eventQueue,
	}
}

func (p *processor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if p.EventQueue.Length() == 0 {
				time.Sleep(ProcessorSleepTime)
			}

			if p.EventQueue.Length() > 0 {
				dbEvent, err := p.PopEvent()
				if err != nil {
					return err
				}

				err = p.HandleEvent(dbEvent)
				if err != nil {
					return err
				}
			}
		}

	}
}

func (p *processor) HandleEvent(event queue.DBEvent) error {
	if event.IsInsert() {
		return p.ExecuteInsertEvent(event)
	}

	if event.IsUpdate() {
		return p.ExecuteUpdateEvent(event)
	}

	if event.IsDelete() {
		return p.ExecuteDeleteEvent(event)
	}

	return ErrUnknownEvent
}

func (p *processor) PopEvent() (queue.DBEvent, error) {
	return p.EventQueue.Pop()
}

func (p *processor) ExecuteInsertEvent(event queue.DBEvent) error {
	if !event.IsInsert() {
		return ErrExpectedInsertEvent
	}

	err := p.DBRepo.InsertFile(context.Background(), &event.File)
	return err
}

func (p *processor) ExecuteUpdateEvent(event queue.DBEvent) error {
	if !event.IsUpdate() {
		return ErrExpectedUpdateEvent
	}

	err := p.DBRepo.UpdateFile(context.Background(), &event.File)
	return err
}

func (p *processor) ExecuteDeleteEvent(event queue.DBEvent) error {
	if !event.IsDelete() {
		return ErrExpectedDeleteEvent
	}

	err := p.DBRepo.DeleteFile(context.Background(), &event.File)
	return err
}
