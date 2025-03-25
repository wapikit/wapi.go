package events

type PhoneNumberNameUpdateEvent struct {
	BaseBusinessAccountEvent
	PhoneNumber     string
	Name            string
	Decision        string
	RejectionReason *string
}

func NewPhoneNumberNameUpdateEvent(baseEvent *BaseBusinessAccountEvent, phoneNumber string, name string, decision string, reason *string) *PhoneNumberNameUpdateEvent {
	return &PhoneNumberNameUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		PhoneNumber:              phoneNumber,
		Name:                     name,
		Decision:                 decision,
		RejectionReason:          reason,
	}
}
