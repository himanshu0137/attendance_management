package events

import "sync"

// Event is a cross-module message published on the event bus.
type Event struct {
	Name string
	Data any
}

// Bus is the cross-module publication interface.
type Bus interface {
	Publish(e Event)
	Subscribe(name string, fn func(e Event))
}

// InMemoryBus is a thread-safe in-memory publish/subscribe bus. Events with no
// subscribers are dropped silently; modules running under -modules must tolerate
// missing peer subscribers without panicking.
type InMemoryBus struct {
	sync.RWMutex
	subscribers map[string][]func(e Event)
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		subscribers: map[string][]func(e Event){},
	}
}

func (b *InMemoryBus) Publish(e Event) {
	b.RLock()
	n := len(b.subscribers[e.Name])
	b.RUnlock()
	if n == 0 {
		return
	}
	b.Lock()
	for _, fn := range b.subscribers[e.Name] {
		fn(e)
	}
	b.Unlock()
}

func (b *InMemoryBus) Subscribe(name string, fn func(e Event)) {
	b.Lock()
	b.subscribers[name] = append(b.subscribers[name], fn)
	b.Unlock()
}
