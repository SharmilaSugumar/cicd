package eventbus

import (
	"context"
	"sync"
)

type EventType string

const (
	EventJobStatusChanged EventType = "JobStatusChanged"
	EventWorkerHeartbeat  EventType = "WorkerHeartbeat"
	EventPipelineUpdated  EventType = "PipelineUpdated"
)

type Event struct {
	Type    EventType
	Payload interface{}
}

type EventHandler func(context.Context, Event)

type EventBus interface {
	Publish(ctx context.Context, event Event)
	Subscribe(eventType EventType, handler EventHandler)
}

type internalEventBus struct {
	subscribers map[EventType][]EventHandler
	mu          sync.RWMutex
}

func NewInternalEventBus() EventBus {
	return &internalEventBus{
		subscribers: make(map[EventType][]EventHandler),
	}
}

func (b *internalEventBus) Publish(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers := b.subscribers[event.Type]
	b.mu.RUnlock()

	for _, h := range handlers {
		go h(ctx, event) // Async execution
	}
}

func (b *internalEventBus) Subscribe(eventType EventType, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[eventType] = append(b.subscribers[eventType], handler)
}
