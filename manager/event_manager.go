package manager

import (
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/wapikit/wapi.go/pkg/events"
)

// subscriberBufferSize is how many events each subscriber can have queued
// before Publish starts dropping events for it.
const subscriberBufferSize = 100

// ChannelEvent represents an event that can be published and subscribed to.
type ChannelEvent struct {
	Type events.EventType // Type is the type of the event.
	Data events.BaseEvent // Data is the data associated with the event.
}

// EventManager is responsible for managing events and their subscribers.
//
// Every subscriber gets its own buffered channel, so each subscriber receives
// every event published for its type. An EventManager must not be copied
// after first use; pass it around as a pointer.
type EventManager struct {
	subscribers  map[events.EventType][]chan ChannelEvent // subscribers maps an event type to one channel per subscriber.
	sync.RWMutex                                          // RWMutex is used to synchronize access to the subscribers map.
}

// NewEventManager creates a new instance of EventManger.
func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[events.EventType][]chan ChannelEvent),
	}
}

// Subscribe adds a new subscriber to the specified event type and returns its
// channel. Each call creates a new channel that receives every event of that
// type. The channel is closed by Unsubscribe.
func (em *EventManager) Subscribe(eventName events.EventType) (chan ChannelEvent, error) {
	em.Lock()
	defer em.Unlock()
	if em.subscribers == nil {
		em.subscribers = make(map[events.EventType][]chan ChannelEvent)
	}
	ch := make(chan ChannelEvent, subscriberBufferSize)
	em.subscribers[eventName] = append(em.subscribers[eventName], ch)
	return ch, nil
}

// Unsubscribe removes every subscriber of the specified event type and closes
// their channels, which stops the goroutines started by On for that type.
func (em *EventManager) Unsubscribe(id events.EventType) {
	em.Lock()
	defer em.Unlock()
	for _, ch := range em.subscribers[id] {
		close(ch)
	}
	delete(em.subscribers, id)
}

// Publish delivers an event to every subscriber of its type without blocking.
// If a subscriber's queue is full the event is dropped for that subscriber
// and an error is returned.
func (em *EventManager) Publish(event events.EventType, data events.BaseEvent) error {
	em.RLock()
	defer em.RUnlock()

	dropped := 0
	for _, ch := range em.subscribers[event] {
		select {
		case ch <- ChannelEvent{Type: event, Data: data}:
		default:
			dropped++
		}
	}
	if dropped > 0 {
		slog.Warn("wapi: event queue full, event dropped", "event_type", string(event), "subscribers_dropped", dropped)
		return fmt.Errorf("event queue full for type: %s", event)
	}
	return nil
}

// On registers a handler function for the specified event type.
// The handler is called, one event at a time, for every event published for
// that type. Calling On several times for one type registers several
// handlers, and each of them receives every event. A panic inside the handler
// is recovered and logged so it cannot crash the process.
// It returns the event type that the handler is registered for.
func (em *EventManager) On(eventName events.EventType, handler func(events.BaseEvent)) events.EventType {
	ch, _ := em.Subscribe(eventName)
	go func() {
		for event := range ch {
			callHandler(eventName, handler, event.Data)
		}
	}()
	return eventName
}

func callHandler(eventName events.EventType, handler func(events.BaseEvent), data events.BaseEvent) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("wapi: event handler panicked",
				"event_type", string(eventName),
				"panic", fmt.Sprint(r),
				"stack", string(debug.Stack()))
		}
	}()
	handler(data)
}
