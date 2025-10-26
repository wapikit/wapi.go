package events

import "github.com/gTahidi/wapi.go/pkg/components"

// ContactsMessageEvent represents an event that occurs when a message with contacts is received.
type ContactsMessageEvent struct {
	BaseMessageEvent `json:",inline"`
	Contacts         components.ContactMessage `json:"contacts"`
}

// NewContactsMessageEvent creates a new ContactsMessageEvent instance.
func NewContactsMessageEvent(baseMessageEvent BaseMessageEvent, contacts components.ContactMessage) *ContactsMessageEvent {
	return &ContactsMessageEvent{
		BaseMessageEvent: baseMessageEvent,
		Contacts:         contacts,
	}
}
