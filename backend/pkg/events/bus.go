package events

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// Event represents a published event.
type Event struct {
	ID        uuid.UUID
	Topic     string
	TenantID  string
	Payload   []byte
	Timestamp int64
}

// HandlerFunc is a function that handles an event.
type HandlerFunc func(context.Context, *Event) error

// Bus is the interface for publishing and subscribing to events.
type Bus interface {
	Publish(ctx context.Context, event *Event) error
	Subscribe(topic string, handler HandlerFunc) (unsubscribe func())
}

// InMemoryBus is a thread-safe in-memory event bus implementation.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]HandlerFunc
}

// NewInMemoryBus creates a new InMemoryBus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]HandlerFunc),
	}
}

// Publish publishes an event to the bus.
func (b *InMemoryBus) Publish(ctx context.Context, event *Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.Topic]
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	// Create a new context with cancellation for each handler
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Execute handlers concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h HandlerFunc) {
			defer wg.Done()
			if err := h(ctx, event); err != nil {
				errChan <- err
			}
		}(handler)
	}

	wg.Wait()
	close(errChan)

	// Check if any handler returned an error
	for err := range errChan {
		if err != nil {
			return err // Return first error encountered
		}
	}

	return nil
}

// Subscribe subscribes to a topic with a handler function.
func (b *InMemoryBus) Subscribe(topic string, handler HandlerFunc) (unsubscribe func()) {
	b.mu.Lock()
	b.handlers[topic] = append(b.handlers[topic], handler)
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		
		handlers := b.handlers[topic]
		for i, h := range handlers {
			if h == handler {
				b.handlers[topic] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}