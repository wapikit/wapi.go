package components

// MessageType represents the type of message.
type MessageType string

// Constants for different message types.
const (
	MessageTypeLocation    MessageType = "location"
	MessageTypeAudio       MessageType = "audio"
	MessageTypeVideo       MessageType = "video"
	MessageTypeDocument    MessageType = "document"
	MessageTypeText        MessageType = "text"
	MessageTypeContact     MessageType = "contacts"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeTemplate    MessageType = "template"
	MessageTypeReaction    MessageType = "reaction"
	MessageTypeSticker     MessageType = "sticker"
	MessageTypeImage       MessageType = "image"
)

// ApiCompatibleJsonConverterConfigs represents the configuration for converting to JSON.
type ApiCompatibleJsonConverterConfigs struct {
	ReplyToMessageId  string
	SendToPhoneNumber string
	// SendToRecipient carries a BSUID (or parent BSUID) send target. When set,
	// the payload serializes `recipient` instead of `to` (Meta BSUID send
	// contract). Empty for phone sends — existing behavior is unchanged.
	SendToRecipient string
}

// Context represents the context of the message.
type Context struct {
	MessageId string `json:"message_id,omitempty"`
}

// BaseMessagePayload represents the base payload to send messages.
type BaseMessagePayload struct {
	Context *Context `json:"context,omitempty"`
	// To is the recipient phone number for phone-targeted sends. Omitted when a
	// BSUID `Recipient` is used (exactly one of To/Recipient is set).
	To string `json:"to,omitempty"`
	// Recipient is the BSUID / parent-BSUID send target (identity rollout).
	// Omitted for phone sends. Serialized only when set.
	Recipient        string      `json:"recipient,omitempty"`
	Type             MessageType `json:"type"`
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
}

// NewBaseMessagePayload creates a new instance of BaseMessagePayload from the
// send configs, choosing phone (`to`) vs BSUID (`recipient`) serialization.
// SendToRecipient (BSUID) takes precedence when both are set — callers should
// set exactly one; setting both is only meaningful for Meta precedence tests.
func NewBaseMessagePayload(configs ApiCompatibleJsonConverterConfigs, messageType MessageType) BaseMessagePayload {
	payload := BaseMessagePayload{
		Type:             messageType,
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
	}
	if configs.SendToRecipient != "" {
		payload.Recipient = configs.SendToRecipient
	}
	if configs.SendToPhoneNumber != "" {
		payload.To = configs.SendToPhoneNumber
	}
	return payload
}

// BaseMessage is an interface for sending messages.
type BaseMessage interface {
	ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error)
}
