package events

// UserIdUpdateEvent is emitted when Meta reports a user-id (BSUID) change on the
// `user_id_update` webhook — e.g. a contact's business-scoped user id is created
// or changes. Consumers use it to bridge/merge identity across the transition.
// All fields are best-effort from Meta's payload; any may be empty depending on
// the transition shape.
type UserIdUpdateEvent struct {
	BaseSystemEvent
	// BusinessAccountId is the WABA (webhook entry) id that scopes this change —
	// BSUID continuity is scoped to a WABA/portfolio, never global, so consumers
	// need it to resolve the affected org/contact.
	BusinessAccountId string `json:"business_account_id,omitempty"`
	// The phone identifier of the affected contact when Meta includes it.
	WaId string `json:"wa_id,omitempty"`
	// Previous and new BSUID values across the transition.
	OldUserId string `json:"old_user_id,omitempty"`
	NewUserId string `json:"new_user_id,omitempty"`
	// Parent BSUID when Meta scopes the id to a parent business.
	ParentUserId string `json:"parent_user_id,omitempty"`
}

func NewUserIdUpdateEvent(baseEvent BaseSystemEvent, businessAccountId, waId, oldUserId, newUserId, parentUserId string) *UserIdUpdateEvent {
	return &UserIdUpdateEvent{
		BaseSystemEvent:   baseEvent,
		BusinessAccountId: businessAccountId,
		WaId:              waId,
		OldUserId:         oldUserId,
		NewUserId:         newUserId,
		ParentUserId:      parentUserId,
	}
}

func (e UserIdUpdateEvent) GetEventType() string {
	return string(UserIdUpdateEventType)
}
