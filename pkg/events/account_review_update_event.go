package events

type AccountReviewUpdateEventEnum string

const (
	AccountReviewUpdateEventEnumApproved AccountReviewUpdateEventEnum = "APPROVED"
	AccountReviewUpdateEventEnumRejected AccountReviewUpdateEventEnum = "REJECTED"
)

type AccountReviewUpdateEvent struct {
	BaseBusinessAccountEvent
	Decision AccountReviewUpdateEventEnum
}

func NewAccountReviewUpdateEvent(baseEvent *BaseBusinessAccountEvent, decision AccountReviewUpdateEventEnum) *AccountReviewUpdateEvent {
	return &AccountReviewUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		Decision:                 decision,
	}
}
