package manager

import (
	"sync"
	"testing"
	"time"

	"github.com/wapikit/wapi.go/pkg/events"
)

func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestEveryHandlerReceivesEveryEvent(t *testing.T) {
	em := NewEventManager()
	var mu sync.Mutex
	counts := map[string]int{}
	for _, name := range []string{"a", "b"} {
		name := name
		em.On(events.TextMessageEventType, func(events.BaseEvent) {
			mu.Lock()
			counts[name]++
			mu.Unlock()
		})
	}

	for i := 0; i < 10; i++ {
		if err := em.Publish(events.TextMessageEventType, events.NewReadyEvent()); err != nil {
			t.Fatalf("Publish: %v", err)
		}
	}

	waitFor(t, "both handlers to see 10 events", func() bool {
		mu.Lock()
		defer mu.Unlock()
		return counts["a"] == 10 && counts["b"] == 10
	})
}

func TestHandlerPanicIsRecoveredAndHandlerKeepsRunning(t *testing.T) {
	em := NewEventManager()
	var mu sync.Mutex
	calls := 0
	em.On(events.TextMessageEventType, func(events.BaseEvent) {
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			var event *events.TextMessageEvent
			_ = event.Text // nil pointer dereference
		}
	})

	_ = em.Publish(events.TextMessageEventType, events.NewReadyEvent())
	_ = em.Publish(events.TextMessageEventType, events.NewReadyEvent())

	waitFor(t, "second event after a panic", func() bool {
		mu.Lock()
		defer mu.Unlock()
		return calls == 2
	})
}

func TestPublishWithoutSubscribersIsNoop(t *testing.T) {
	if err := NewEventManager().Publish(events.TextMessageEventType, events.NewReadyEvent()); err != nil {
		t.Fatalf("Publish: %v", err)
	}
}

func TestPublishReportsFullQueue(t *testing.T) {
	em := NewEventManager()
	if _, err := em.Subscribe(events.TextMessageEventType); err != nil { // nobody reads it
		t.Fatal(err)
	}
	for i := 0; i < subscriberBufferSize; i++ {
		if err := em.Publish(events.TextMessageEventType, events.NewReadyEvent()); err != nil {
			t.Fatalf("publish %d: %v", i, err)
		}
	}
	if err := em.Publish(events.TextMessageEventType, events.NewReadyEvent()); err == nil {
		t.Fatalf("expected an error once the queue is full")
	}
}

func TestUnsubscribeClosesChannels(t *testing.T) {
	em := NewEventManager()
	ch, _ := em.Subscribe(events.TextMessageEventType)
	em.Unsubscribe(events.TextMessageEventType)
	if _, open := <-ch; open {
		t.Fatalf("channel still open after Unsubscribe")
	}
	if err := em.Publish(events.TextMessageEventType, events.NewReadyEvent()); err != nil {
		t.Fatalf("Publish after Unsubscribe: %v", err)
	}
}
