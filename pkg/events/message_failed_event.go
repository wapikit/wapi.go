package events

type MessageFailedEvent struct {
	BaseSystemEvent `json:",inline"`
	MessageId       string `json:"messageId"`
	SentTo          string `json:"sentTo"`
	FailReason      string `json:"failReason"`
	ErrorCode       int    `json:"errorCode"`
	ErrorMessage    string `json:"errorMessage"`
}

func NewMessageFailedEvent(baseSystemEvent BaseSystemEvent, messageId, sendTo, failReason string, errCode int, errorMessage string) *MessageFailedEvent {
	return &MessageFailedEvent{
		BaseSystemEvent: baseSystemEvent,
		MessageId:       messageId,
		SentTo:          sendTo,
		FailReason:      failReason,
		ErrorCode:       errCode,
		ErrorMessage:    errorMessage,
	}

}
