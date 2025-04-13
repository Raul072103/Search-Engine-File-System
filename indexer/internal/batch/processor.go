package batch

import (
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/database"
	"context"
	"errors"
	"go.uber.org/zap"
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
	logger     *zap.Logger
}

const (
	ProcessorSleepTime = time.Second * 5
)

var (
	ErrExpectedInsertEvent          = errors.New("expected the event to be of type INSERT")
	ErrExpectedUpdateEvent          = errors.New("expected the event to be of type UPDATE")
	ErrExpectedDeleteEvent          = errors.New("expected the event to be of type DELETE")
	ErrExpectedRecursiveDeleteEvent = errors.New("expected the event to be of type RECURSIVE DELETE")
	ErrUnknownEvent                 = errors.New("unknown event sent for processing")
)

// NewProcessor creates a new instance of the Processor.
func NewProcessor(dbRepo database.Repo, eventQueue *queue.InMemoryQueue, logger *zap.Logger) Processor {
	return &processor{
		DBRepo:     dbRepo,
		EventQueue: eventQueue,
		logger:     logger,
	}
}

// Run is the entry point into the start of the Processor component.
// It makes use of the context passed to it, in order to detect a STOP signal from the
// main goroutine to gracefully shutdown.
func (p *processor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Processor stopped")
			return nil
		default:
			if p.EventQueue.Length() == 0 {
				time.Sleep(ProcessorSleepTime)
			}

			if p.EventQueue.Length() > 0 {
				dbEvent, err := p.PopEvent()
				if err != nil {
					p.logger.Error("Couldn't pop event from events queue", zap.Error(err))
					return err
				}

				err = p.HandleEvent(dbEvent)
				if err != nil {
					// if it is an error, log it
					p.logger.Error("Couldn't handle db event", zap.String("path", dbEvent.File.Path), zap.Error(err))
				}
			}
		}
	}
}

// HandleEvent handles the corresponding event based on the type of the event.
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

	if event.IsRecursiveDelete() {
		return p.ExecuteRecursiveDeleteEvent(event)
	}

	return ErrUnknownEvent
}

// PopEvent pops an element from the InMemoryQueue.
func (p *processor) PopEvent() (queue.DBEvent, error) {
	return p.EventQueue.Pop()
}

// ExecuteInsertEvent takes a INSERT event and applies an INSERT operation on the database.
func (p *processor) ExecuteInsertEvent(event queue.DBEvent) error {
	if !event.IsInsert() {
		return ErrExpectedInsertEvent
	}

	err := p.DBRepo.Files.Insert(context.Background(), &event.File)
	return err
}

// ExecuteUpdateEvent takes a UPDATE event and applies an UPDATE operation on the database.
func (p *processor) ExecuteUpdateEvent(event queue.DBEvent) error {
	if !event.IsUpdate() {
		return ErrExpectedUpdateEvent
	}

	err := p.DBRepo.Files.Update(context.Background(), &event.File)
	return err
}

// ExecuteDeleteEvent takes a DELETE event and applies a DELETE operation on the database.
func (p *processor) ExecuteDeleteEvent(event queue.DBEvent) error {
	if !event.IsDelete() {
		return ErrExpectedDeleteEvent
	}

	err := p.DBRepo.Files.Delete(context.Background(), &event.File)
	return err
}

// ExecuteRecursiveDeleteEvent takes a RECURSIVE DELETE operation and applies it on the database.
func (p *processor) ExecuteRecursiveDeleteEvent(event queue.DBEvent) error {
	if !event.IsRecursiveDelete() {
		return ErrExpectedRecursiveDeleteEvent
	}

	err := p.DBRepo.Files.DeleteAllUnderDirectory(context.Background(), &event.File)
	return err
}
