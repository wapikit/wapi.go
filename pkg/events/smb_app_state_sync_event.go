package events

// SmbAppStateSyncEvent represents an event for SMB app state synchronization
type SmbAppStateSyncEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	MessagingProduct         string          `json:"messaging_product"`
	DisplayPhoneNumber       string          `json:"display_phone_number"`
	PhoneNumberId            string          `json:"phone_number_id"`
	StateSync                []StateSyncItem `json:"state_sync"`
}

type StateSyncItem struct {
	Type      string      `json:"type"` // e.g., "contact"
	Contact   ContactSync `json:"contact,omitempty"`
	Action    string      `json:"action"` // e.g., "add", "remove", "update"
	Timestamp string      `json:"timestamp"`
}

type ContactSync struct {
	FullName    string `json:"full_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// NewSmbAppStateSyncEvent creates a new SmbAppStateSyncEvent instance
func NewSmbAppStateSyncEvent(
	baseEvent *BaseBusinessAccountEvent,
	messagingProduct string,
	displayPhoneNumber string,
	phoneNumberId string,
	stateSync []StateSyncItem,
) *SmbAppStateSyncEvent {
	return &SmbAppStateSyncEvent{
		BaseBusinessAccountEvent: *baseEvent,
		MessagingProduct:         messagingProduct,
		DisplayPhoneNumber:       displayPhoneNumber,
		PhoneNumberId:            phoneNumberId,
		StateSync:                stateSync,
	}
}

func (e *SmbAppStateSyncEvent) GetEventType() string {
	return string(SmbAppStateSyncEventType)
}
