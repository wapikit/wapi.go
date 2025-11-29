package events

// SmbMessageEchoesEvent represents an event for SMB message echoes
type SmbMessageEchoesEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	MessagingProduct         string        `json:"messaging_product"`
	DisplayPhoneNumber       string        `json:"display_phone_number"`
	PhoneNumberId            string        `json:"phone_number_id"`
	MessageEchoes            []MessageEcho `json:"message_echoes"`
}

type MessageEcho struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

// NewSmbMessageEchoesEvent creates a new SmbMessageEchoesEvent instance
func NewSmbMessageEchoesEvent(
	baseEvent *BaseBusinessAccountEvent,
	messagingProduct string,
	displayPhoneNumber string,
	phoneNumberId string,
	messageEchoes []MessageEcho,
) *SmbMessageEchoesEvent {
	return &SmbMessageEchoesEvent{
		BaseBusinessAccountEvent: *baseEvent,
		MessagingProduct:         messagingProduct,
		DisplayPhoneNumber:       displayPhoneNumber,
		PhoneNumberId:            phoneNumberId,
		MessageEchoes:            messageEchoes,
	}
}

func (e *SmbMessageEchoesEvent) GetEventType() string {
	return string(SmbMessageEchoesEventType)
}
