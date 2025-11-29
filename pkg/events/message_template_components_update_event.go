package events

// MessageTemplateComponentsUpdateEvent represents an event for message template components update
type MessageTemplateComponentsUpdateEvent struct {
	BaseBusinessAccountEvent `json:",inline"`
	MessageTemplateId        int64                   `json:"message_template_id"`
	MessageTemplateName      string                  `json:"message_template_name"`
	MessageTemplateLanguage  string                  `json:"message_template_language"`
	MessageTemplateElement   string                  `json:"message_template_element"` // Template body text
	MessageTemplateTitle     string                  `json:"message_template_title,omitempty"`
	MessageTemplateFooter    string                  `json:"message_template_footer,omitempty"`
	MessageTemplateButtons   []MessageTemplateButton `json:"message_template_buttons,omitempty"`
}

type MessageTemplateButton struct {
	Type        string `json:"type"`
	Text        string `json:"text"`
	Url         string `json:"url,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// NewMessageTemplateComponentsUpdateEvent creates a new MessageTemplateComponentsUpdateEvent instance
func NewMessageTemplateComponentsUpdateEvent(
	baseEvent *BaseBusinessAccountEvent,
	templateId int64,
	templateName string,
	templateLanguage string,
	element string,
	title string,
	footer string,
	buttons []MessageTemplateButton,
) *MessageTemplateComponentsUpdateEvent {
	return &MessageTemplateComponentsUpdateEvent{
		BaseBusinessAccountEvent: *baseEvent,
		MessageTemplateId:        templateId,
		MessageTemplateName:      templateName,
		MessageTemplateLanguage:  templateLanguage,
		MessageTemplateElement:   element,
		MessageTemplateTitle:     title,
		MessageTemplateFooter:    footer,
		MessageTemplateButtons:   buttons,
	}
}

func (e *MessageTemplateComponentsUpdateEvent) GetEventType() string {
	return string(MessageTemplateComponentsUpdateEventType)
}
