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
	// MessageTemplateCategory is the template's category, when Meta sends it.
	MessageTemplateCategory MessageTemplateCategoryEnum
	// DisableInfo is set when the template is disabled.
	DisableInfo *TemplateDisableInfo
	// OtherInfo is set when the template is locked or unlocked.
	OtherInfo *TemplateOtherInfo
	// RejectionInfo is set when the template is rejected with INVALID_FORMAT.
	RejectionInfo *TemplateRejectionInfo
}

// TemplateDisableInfo says when a template was disabled.
type TemplateDisableInfo struct {
	DisableDate string
}

// TemplateOtherInfo carries the lock/unlock details of a template, such as
// Title "FIRST_PAUSE" or "UNPAUSE".
type TemplateOtherInfo struct {
	Title       string
	Description string
}

// TemplateRejectionInfo explains an INVALID_FORMAT rejection.
type TemplateRejectionInfo struct {
	Reason         string
	Recommendation string
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
