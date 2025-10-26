package events

import "github.com/gTahidi/wapi.go/pkg/components"

// ReactionMessageEvent represents an event that occurs when a reaction is added to a message.
type ReactionMessageEvent struct {
	BaseMessageEvent `json:",inline"`
	Reaction         components.ReactionMessage
}

// NewReactionMessageEvent creates a new ReactionMessageEvent instance.
func NewReactionMessageEvent(baseMessageEvent BaseMessageEvent, reaction components.ReactionMessage) *ReactionMessageEvent {
	return &ReactionMessageEvent{
		BaseMessageEvent: baseMessageEvent,
		Reaction:         reaction,
	}
}
