package events

type SecurityEvent struct {
	BaseBusinessAccountEvent
	DisplayPhoneNumber string
	// Event is the security change, for example "PIN_CHANGED".
	Event     string
	Requester string
}

func NewSecurity() *SecurityEvent {
	return &SecurityEvent{}
}
