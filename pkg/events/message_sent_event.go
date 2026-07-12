package events

// MessageSentEvent represents an event indicating that a message has been sent.
type MessageSentEvent struct {
	BaseSystemEvent `json:",inline"`
	MessageId       string `json:"messageId"`
	SentTo          string `json:"sentTo"`

	// Pricing carries Meta's `messages.statuses[].pricing` block when
	// present. Nil for status events Meta didn't tag with a pricing
	// payload. Consumers that bill on the upstream charge should check
	// `Pricing != nil && Pricing.Billable` before committing.
	Pricing *MessagePricingInfo `json:"pricing,omitempty"`

	// Recipient BSUID identity (identity rollout). Present alongside SentTo
	// (recipient_id) in the dual-identifier shape; empty otherwise.
	RecipientUserId       string `json:"recipientUserId,omitempty"`
	RecipientParentUserId string `json:"recipientParentUserId,omitempty"`
}

// NewMessageSentEvent creates a new instance of MessageSentEvent.
func NewMessageSentEvent(baseSystemEvent BaseSystemEvent, messageId, sendTo string) *MessageSentEvent {
	return &MessageSentEvent{
		BaseSystemEvent: baseSystemEvent,
		MessageId:       messageId,
		SentTo:          sendTo,
	}
}
