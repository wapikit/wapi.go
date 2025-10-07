package events

// MessageUndeliveredEvent represents an event related to an undelivered message.
type MessageUndeliveredEvent struct {
	BaseSystemEvent `json:",inline"`
	MessageId       string `json:"messageId"`
	SentTo          string `json:"sentTo"`
	Reason          string `json:"reason"`
	ErrorCode       int    `json:"errorCode"`
	ErrorMessage    string `json:"errorMessage"`
}

// NewMessageUndeliveredEvent creates a new instance of MessageUndeliveredEvent.
func NewMessageUndeliveredEvent(baseSystemEvent BaseSystemEvent, messageId, sendTo, reason string, errorCode int, errorMessage string) *MessageUndeliveredEvent {
	return &MessageUndeliveredEvent{
		BaseSystemEvent: baseSystemEvent,
		MessageId:       messageId,
		SentTo:          sendTo,
		Reason:          reason,
		ErrorCode:       errorCode,
		ErrorMessage:    errorMessage,
	}
}
