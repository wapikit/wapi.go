package events

// MessageReadEvent represents an event indicating that a message has been read.
type MessageReadEvent struct {
	BaseSystemEvent `json:",inline"`
	MessageId       string `json:"messageId"`
	SentTo          string `json:"sentTo"`

	// Pricing carries Meta's `messages.statuses[].pricing` block when
	// present. Most pricing payloads arrive earlier (on `sent` or
	// `delivered`); included here for parity since Meta sometimes
	// retroactively populates pricing on later status updates.
	Pricing *MessagePricingInfo `json:"pricing,omitempty"`

	// Recipient BSUID identity (identity rollout). Present alongside SentTo
	// (recipient_id) in the dual-identifier shape; empty otherwise.
	RecipientUserId       string `json:"recipientUserId,omitempty"`
	RecipientParentUserId string `json:"recipientParentUserId,omitempty"`
}

// NewMessageReadEvent creates a new instance of MessageReadEvent.
func NewMessageReadEvent(baseSystemEvent BaseSystemEvent, messageId, sendTo string) *MessageReadEvent {
	return &MessageReadEvent{
		BaseSystemEvent: baseSystemEvent,
		MessageId:       messageId,
		SentTo:          sendTo,
	}
}
