package events

import "github.com/gTahidi/wapi.go/pkg/components"

// LocationMessageEvent represents an event that contains a location message.
type LocationMessageEvent struct {
	BaseMessageEvent `json:",inline"`
	Location         components.LocationMessage `json:"location"`
}

// NewLocationMessageEvent creates a new LocationMessageEvent instance.
func NewLocationMessageEvent(baseMessageEvent BaseMessageEvent, location components.LocationMessage) *LocationMessageEvent {
	return &LocationMessageEvent{
		BaseMessageEvent: baseMessageEvent,
		Location:         location,
	}
}
