package events

type MessageTemplateStatusUpdateEventEnum string

const (
	MessageTemplateStatusUpdateEventEnumApproved   MessageTemplateStatusUpdateEventEnum = "APPROVED"
	MessageTemplateStatusUpdateEventEnumDisabled   MessageTemplateStatusUpdateEventEnum = "DISABLED"
	MessageTemplateStatusUpdateEventEnumInAppeal   MessageTemplateStatusUpdateEventEnum = "IN_APPEAL"
	MessageTemplateStatusUpdateEventEnumPending    MessageTemplateStatusUpdateEventEnum = "PENDING"
	MessageTemplateStatusUpdateEventEnumReinstated MessageTemplateStatusUpdateEventEnum = "REINSTATED"
	MessageTemplateStatusUpdateEventEnumRejected   MessageTemplateStatusUpdateEventEnum = "REJECTED"
	MessageTemplateStatusUpdateEventEnumFlagged    MessageTemplateStatusUpdateEventEnum = "FLAGGED"
)

type MessageTemplateStatusUpdateReason string

const (
	MessageTemplateStatusUpdateReasonAbusiveContent    MessageTemplateStatusUpdateReason = "ABUSIVE_CONTENT"
	MessageTemplateStatusUpdateReasonIncorrectCategory MessageTemplateStatusUpdateReason = "INCORRECT_CATEGORY"
	MessageTemplateStatusUpdateReasonInvalidFormat     MessageTemplateStatusUpdateReason = "INVALID_FORMAT"
	MessageTemplateStatusUpdateReasonNone              MessageTemplateStatusUpdateReason = "NONE"
	MessageTemplateStatusUpdateReasonScam              MessageTemplateStatusUpdateReason = "SCAM"
)

type MessageTemplateStatusUpdateEvent struct {
	BaseBusinessAccountEvent
	Event                   MessageTemplateStatusUpdateEventEnum
	MessageTemplateId       int64
	MessageTemplateName     string
	MessageTemplateLanguage string
	Reason                  MessageTemplateStatusUpdateReason
}

func NewMessageTemplateStatusUpdateEvent(baseEvent *BaseBusinessAccountEvent, event MessageTemplateStatusUpdateEventEnum, messageTemplateId int64, messageTemplateName string, messageTemplateLanguage string, reason MessageTemplateStatusUpdateReason) *MessageTemplateStatusUpdateEvent {
	return &MessageTemplateStatusUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		Event:                    event,
		MessageTemplateId:        messageTemplateId,
		MessageTemplateName:      messageTemplateName,
		MessageTemplateLanguage:  messageTemplateLanguage,
		Reason:                   reason,
	}
}
