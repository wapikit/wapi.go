package events

type BusinessCapabilityUpdateEvent struct {
	BaseBusinessAccountEvent
	MaxDailyConversationPerPhone int64
	MaxPhoneNumbersPerBusiness   int64
}

func NewBusinessCapabilityUpdateEvent(baseEvent *BaseBusinessAccountEvent, maxDailyConversationPerPhone int64, maxPhoneNumbersPerBusiness int64) *BusinessCapabilityUpdateEvent {
	return &BusinessCapabilityUpdateEvent{
		BaseBusinessAccountEvent:     *baseEvent,
		MaxDailyConversationPerPhone: maxDailyConversationPerPhone,
		MaxPhoneNumbersPerBusiness:   maxPhoneNumbersPerBusiness,
	}
}
