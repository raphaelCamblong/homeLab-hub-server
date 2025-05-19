package streaming

import (
	"context"
	"sync"
)

type EventType string

const (
	ActionEvent EventType = "action"
	JobEvent    EventType = "job"
	StepEvent   EventType = "step"
)

type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

type Subscription struct {
	Channel chan Event
	Type    EventType
}

type StreamHub struct {
	mu            sync.RWMutex
	subscriptions map[EventType]map[*Subscription]struct{}
}

func NewStreamHub() *StreamHub {
	return &StreamHub{
		subscriptions: make(map[EventType]map[*Subscription]struct{}),
	}
}

// Subscribe creates a new subscription for the given event type
func (h *StreamHub) Subscribe(ctx context.Context, eventType EventType) (<-chan Event, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan Event, 1)
	sub := &Subscription{
		Channel: ch,
		Type:    eventType,
	}

	if h.subscriptions[eventType] == nil {
		h.subscriptions[eventType] = make(map[*Subscription]struct{})
	}

	h.subscriptions[eventType][sub] = struct{}{}

	go func() {
		<-ctx.Done()
		h.unsubscribe(sub)
	}()

	return ch, nil
}

// Publish sends an event to all subscribers of the event type
func (h *StreamHub) Publish(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subs := h.subscriptions[event.Type]
	if subs == nil {
		return
	}

	// Send event to all subscribers
	for sub := range subs {
		select {
		case sub.Channel <- event:
			// Event sent successfully
		default:
			// Channel is full or closed, remove subscription
			go h.unsubscribe(sub)
		}
	}
}

// CloseAllSubscriptions closes all subscriptions for a given event type
func (h *StreamHub) CloseAllSubscriptions(eventType EventType) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs := h.subscriptions[eventType]; subs != nil {
		for sub := range subs {
			close(sub.Channel)
		}
		delete(h.subscriptions, eventType)
	}
}

func (h *StreamHub) unsubscribe(sub *Subscription) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if subs := h.subscriptions[sub.Type]; subs != nil {
		delete(subs, sub)
		if len(subs) == 0 {
			delete(h.subscriptions, sub.Type)
		}
	}
	close(sub.Channel)
}
