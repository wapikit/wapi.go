package events

type MessageTemplateCategoryEnum string

const (
	MessageTemplateCategoryEnumMarketing      MessageTemplateCategoryEnum = "MARKETING"
	MessageTemplateCategoryEnumOtp            MessageTemplateCategoryEnum = "OTP"
	MessageTemplateCategoryEnumTransactional  MessageTemplateCategoryEnum = "TRANSACTIONAL"
	MessageTemplateCategoryEnumAuthentication MessageTemplateCategoryEnum = "AUTHENTICATION"
	MessageTemplateCategoryEnumUtility        MessageTemplateCategoryEnum = "UTILITY"
)

type TemplateCategoryUpdateEvent struct {
	BaseBusinessAccountEvent
	MessageTemplateId       int64
	MessageTemplateName     string
	MessageTemplateLanguage string
	PreviousCategory        MessageTemplateCategoryEnum
	NewCategory             MessageTemplateCategoryEnum
	// CorrectCategory is the category Meta determined the template belongs to,
	// when it sends one.
	CorrectCategory MessageTemplateCategoryEnum
}

func NewMessageTemplateCategoryUpdateEvent(baseEvent *BaseBusinessAccountEvent, messageTemplateId int64, messageTemplateName string, messageTemplateLanguage string, previousCategory MessageTemplateCategoryEnum, newCategory MessageTemplateCategoryEnum) *TemplateCategoryUpdateEvent {
	return &TemplateCategoryUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		MessageTemplateId:        messageTemplateId,
		MessageTemplateName:      messageTemplateName,
		MessageTemplateLanguage:  messageTemplateLanguage,
		PreviousCategory:         previousCategory,
		NewCategory:              newCategory,
	}
}
