package events

// HistoryEvent represents an event for chat history synchronization
type HistoryEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	MessagingProduct         string         `json:"messaging_product"`
	DisplayPhoneNumber       string         `json:"display_phone_number"`
	PhoneNumberId            string         `json:"phone_number_id"`
	HistoryChunks            []HistoryChunk `json:"history_chunks"`
}

type HistoryChunk struct {
	Phase      int             `json:"phase"`
	ChunkOrder int             `json:"chunk_order"`
	Progress   int             `json:"progress"`
	Threads    []HistoryThread `json:"threads"`
}

type HistoryThread struct {
	Id       string           `json:"id"` // WhatsApp user phone number
	Messages []HistoryMessage `json:"messages"`
}

type HistoryMessage struct {
	From          string `json:"from"`
	To            string `json:"to,omitempty"` // Only included if SMB message echo
	Id            string `json:"id"`
	Timestamp     string `json:"timestamp"`
	Type          string `json:"type"`
	MessageStatus string `json:"message_status"` // From history_context.status
}

// NewHistoryEvent creates a new HistoryEvent instance
func NewHistoryEvent(
	baseEvent *BaseBusinessAccountEvent,
	messagingProduct string,
	displayPhoneNumber string,
	phoneNumberId string,
	historyChunks []HistoryChunk,
) *HistoryEvent {
	return &HistoryEvent{
		BaseBusinessAccountEvent: *baseEvent,
		MessagingProduct:         messagingProduct,
		DisplayPhoneNumber:       displayPhoneNumber,
		PhoneNumberId:            phoneNumberId,
		HistoryChunks:            historyChunks,
	}
}

func (e *HistoryEvent) GetEventType() string {
	return string(HistoryEventType)
}
