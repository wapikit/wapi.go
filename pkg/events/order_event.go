package events

import "github.com/gTahidi/wapi.go/pkg/components"

// OrderEvent represents an event related to an order.
type OrderEvent struct {
	BaseMessageEvent `json:",inline"`
	Order            components.Order `json:"order"`
}

// NewOrderEvent creates a new OrderEvent instance.
func NewOrderEvent(baseMessageEvent BaseMessageEvent, order components.Order) *OrderEvent {
	return &OrderEvent{
		BaseMessageEvent: baseMessageEvent,
		Order:            order,
	}
}
