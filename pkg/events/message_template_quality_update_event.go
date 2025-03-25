package events

type MessageTemplateQualityUpdateQualityScoreEnum string

const (
	MessageTemplateQualityUpdateQualityScoreEnumGreen   MessageTemplateQualityUpdateQualityScoreEnum = "GREEN"
	MessageTemplateQualityUpdateQualityScoreEnumRed     MessageTemplateQualityUpdateQualityScoreEnum = "RED"
	MessageTemplateQualityUpdateQualityScoreEnumUnknown MessageTemplateQualityUpdateQualityScoreEnum = "UNKNOWN"
	MessageTemplateQualityUpdateQualityScoreEnumYellow  MessageTemplateQualityUpdateQualityScoreEnum = "YELLOW"
)

type MessageTemplateQualityUpdateEvent struct {
	BaseBusinessAccountEvent
	PreviousQualityScore    MessageTemplateQualityUpdateQualityScoreEnum
	NewQualityScore         MessageTemplateQualityUpdateQualityScoreEnum
	MessageTemplateId       int64
	MessageTemplateName     string
	MessageTemplateLanguage string
}

func NewMessageTemplateQualityUpdateEvent(baseEvent *BaseBusinessAccountEvent, previousQualityScore MessageTemplateQualityUpdateQualityScoreEnum, newQualityScore MessageTemplateQualityUpdateQualityScoreEnum, messageTemplateId int64, messageTemplateName string, messageTemplateLanguage string) *MessageTemplateQualityUpdateEvent {
	return &MessageTemplateQualityUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		PreviousQualityScore:     previousQualityScore,
		NewQualityScore:          newQualityScore,
		MessageTemplateId:        messageTemplateId,
		MessageTemplateName:      messageTemplateName,
		MessageTemplateLanguage:  messageTemplateLanguage,
	}
}
