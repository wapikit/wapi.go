package events

// MessageDeliveredEvent represents an event related to an undelivered message.
type MessageDeliveredEvent struct {
	BaseSystemEvent `json:",inline"`
	MessageId       string `json:"messageId"`
	SentTo          string `json:"sentTo"`

	// Pricing carries Meta's `messages.statuses[].pricing` block when
	// present. For CBP-priced sends Meta typically populates this on
	// `delivered` rather than `sent`. Nil otherwise.
	Pricing *MessagePricingInfo `json:"pricing,omitempty"`

	// Recipient BSUID identity (identity rollout). Present alongside SentTo
	// (recipient_id) in the dual-identifier shape; empty otherwise.
	RecipientUserId       string `json:"recipientUserId,omitempty"`
	RecipientParentUserId string `json:"recipientParentUserId,omitempty"`
}

// MessageDeliveredEvent creates a new instance of MessageUndeliveredEvent.
func NewMessageDeliveredEvent(baseSystemEvent BaseSystemEvent, messageId, sendTo string) *MessageDeliveredEvent {
	return &MessageDeliveredEvent{
		BaseSystemEvent: baseSystemEvent,
		MessageId:       messageId,
		SentTo:          sendTo,
	}
}
