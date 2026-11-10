package events

type AccountReviewUpdateEventEnum string

const (
	AccountReviewUpdateEventEnumApproved AccountReviewUpdateEventEnum = "APPROVED"
	AccountReviewUpdateEventEnumRejected AccountReviewUpdateEventEnum = "REJECTED"
	// AccountReviewUpdateEventEnumPending means the review is still open.
	AccountReviewUpdateEventEnumPending AccountReviewUpdateEventEnum = "PENDING"
	// AccountReviewUpdateEventEnumDeferred means the decision is on hold.
	AccountReviewUpdateEventEnumDeferred AccountReviewUpdateEventEnum = "DEFERRED"
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
