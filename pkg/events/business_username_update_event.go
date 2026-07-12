package events

// BusinessUsernameUpdateEvent is emitted on Meta's `business_username_updates`
// webhook when the business's own WhatsApp username changes for a phone number
// (set, changed, or removed). Fields are best-effort from Meta's payload.
type BusinessUsernameUpdateEvent struct {
	BaseBusinessAccountEvent
	PhoneNumberId    string `json:"phone_number_id,omitempty"`
	Username         string `json:"username,omitempty"`
	PreviousUsername string `json:"previous_username,omitempty"`
	// Event describes the change (e.g. created/updated/deleted) when Meta
	// provides it; raw string preserved.
	Event string `json:"event,omitempty"`
}

func NewBusinessUsernameUpdateEvent(baseEvent *BaseBusinessAccountEvent, phoneNumberId, username, previousUsername, event string) *BusinessUsernameUpdateEvent {
	return &BusinessUsernameUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		PhoneNumberId:            phoneNumberId,
		Username:                 username,
		PreviousUsername:         previousUsername,
		Event:                    event,
	}
}

func (e BusinessUsernameUpdateEvent) GetEventType() string {
	return string(BusinessUsernameUpdateEventType)
}
