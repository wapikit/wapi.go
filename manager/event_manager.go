package manager

import (
	"fmt"
	"sync"

	"github.com/wapikit/wapi.go/pkg/events"
)

// ChannelEvent represents an event that can be published and subscribed to.
type ChannelEvent struct {
	Type events.EventType // Type is the type of the event.
	Data events.BaseEvent // Data is the data associated with the event.
}

// EventManager is responsible for managing events and their subscribers.
type EventManager struct {
	subscribers  map[events.EventType]chan ChannelEvent // subscribers is a map of event types to channels of ChannelEvent.
	sync.RWMutex                                        // RWMutex is used to synchronize access to the subscribers map.
}

// NewEventManager creates a new instance of EventManger.
func NewEventManager() *EventManager {
	return &EventManager{
		subscribers: make(map[events.EventType]chan ChannelEvent),
	}
}

// Subscribe adds a new subscriber to the specified event type.
// The subscriber will be notified when the event is published.
func (em *EventManager) Subscribe(eventName events.EventType) (chan ChannelEvent, error) {
	em.Lock()
	defer em.Unlock()
	if ch, ok := em.subscribers[eventName]; ok {
		return ch, nil
	}
	em.subscribers[eventName] = make(chan ChannelEvent, 100)
	return em.subscribers[eventName], nil
}

// Unsubscribe removes a subscriber from the specified event type.
func (em *EventManager) Unsubscribe(id events.EventType) {
	em.Lock()
	defer em.Unlock()
	delete(em.subscribers, id)
}

// Publish publishes an event to the event system and notifies all the subscribers.
func (em *EventManager) Publish(event events.EventType, data events.BaseEvent) error {
	em.Lock()
	defer em.Unlock()

	if ch, ok := em.subscribers[event]; ok {
		select {
		case ch <- ChannelEvent{
			Type: event,
			Data: data,
		}:
		default:
			return fmt.Errorf("event queue full for type: %s", event)
		}
	}
	return nil
}

// On registers a handler function for the specified event type.
// The handler function will be called whenever the event is published.
// It returns the event type that the handler is registered for.
func (em *EventManager) On(eventName events.EventType, handler func(events.BaseEvent)) events.EventType {
	ch, _ := em.Subscribe(eventName)
	go func() {
		for {
			select {
			case event := <-ch:
				handler(event.Data)
			}
		}
	}()
	return eventName
}
