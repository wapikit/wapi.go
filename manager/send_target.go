package manager

import "github.com/wapikit/wapi.go/pkg/components"

// MessageTargetType is how a send is addressed. Phone targets serialize as
// Meta's `to`; BSUID and parent-BSUID targets serialize as `recipient` (Meta
// business-scoped user id send contract).
type MessageTargetType string

const (
	MessageTargetTypePhone       MessageTargetType = "phone"
	MessageTargetTypeBSUID       MessageTargetType = "bsuid"
	MessageTargetTypeParentBSUID MessageTargetType = "parent_bsuid"
)

// MessageTarget identifies who a message is sent to, independent of message
// type. Use the constructors below; the zero value is not a valid target.
type MessageTarget struct {
	Type  MessageTargetType
	Value string
}

// NewPhoneTarget targets a recipient by phone number (serialized as `to`). This
// is the historical default and keeps existing behavior byte-for-byte.
func NewPhoneTarget(phoneNumber string) MessageTarget {
	return MessageTarget{Type: MessageTargetTypePhone, Value: phoneNumber}
}

// NewBSUIDTarget targets a recipient by business-scoped user id (serialized as
// `recipient`).
func NewBSUIDTarget(bsuid string) MessageTarget {
	return MessageTarget{Type: MessageTargetTypeBSUID, Value: bsuid}
}

// NewParentBSUIDTarget targets a recipient by parent business-scoped user id
// (serialized as `recipient`).
func NewParentBSUIDTarget(parentBSUID string) MessageTarget {
	return MessageTarget{Type: MessageTargetTypeParentBSUID, Value: parentBSUID}
}

// configs builds the send configs for this target: phone → `to`, BSUID/parent →
// `recipient`. replyTo is optional (empty for a non-reply send).
func (t MessageTarget) configs(replyTo string) components.ApiCompatibleJsonConverterConfigs {
	configs := components.ApiCompatibleJsonConverterConfigs{ReplyToMessageId: replyTo}
	switch t.Type {
	case MessageTargetTypePhone:
		configs.SendToPhoneNumber = t.Value
	default: // bsuid, parent_bsuid
		configs.SendToRecipient = t.Value
	}
	return configs
}
