package events

import (
	"net/http"
	"strings"

	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/pkg/components"
)

type MessageContext struct {
	From string `json:"from"`
}

type BaseEvent interface {
	GetEventType() string
}

type BaseMessageEventInterface interface {
	BaseEvent
	Reply() (string, error)
	React() (string, error)
}

type BaseSystemEventInterface interface {
	BaseEvent
}

type BaseBusinessAccountEventInterface interface {
	BaseEvent
}

type BaseMessageEvent struct {
	BusinessAccountId string `json:"business_account_id"`
	requester         request_client.RequestClient
	MessageId         string         `json:"message_id"`
	Context           MessageContext `json:"context"`
	Timestamp         string         `json:"timestamp"`
	IsForwarded       bool           `json:"is_forwarded"`
	PhoneNumber       string         `json:"phone_number"`
}

func NewBaseMessageEvent(businessAccountId, messageId, phoneNumber, timestamp, from string, isForwarded bool, requester request_client.RequestClient) BaseMessageEvent {
	return BaseMessageEvent{
		MessageId: messageId,
		Context: MessageContext{
			From: from,
		},
		requester:   requester,
		Timestamp:   timestamp,
		IsForwarded: isForwarded,
		PhoneNumber: phoneNumber,
	}
}

func (bme BaseMessageEvent) GetEventType() string {
	return "message"
}

// Reply to the message
func (baseMessageEvent *BaseMessageEvent) Reply(Message components.BaseMessage) (string, error) {
	body, err := Message.ToJson(components.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: baseMessageEvent.Context.From,
		ReplyToMessageId:  baseMessageEvent.MessageId,
	})

	if err != nil {
		return "", err
	}

	apiRequest := baseMessageEvent.requester.NewApiRequest(strings.Join([]string{baseMessageEvent.PhoneNumber, "messages"}, "/"), http.MethodPost)
	apiRequest.SetBody(string(body))
	apiRequest.Execute()

	return "", nil

}

// React to the message
func (baseMessageEvent *BaseMessageEvent) React(emoji string) (string, error) {
	reactionMessage, err := components.NewReactionMessage(components.ReactionMessageParams{
		Emoji:     emoji,
		MessageId: baseMessageEvent.MessageId,
	})
	if err != nil {
		return "", err
	}
	baseMessageEvent.Reply(reactionMessage)
	return "", nil
}

// BaseMediaMessageEvent represents a base media message event which contains media information.
type BaseMediaMessageEvent struct {
	BaseMessageEvent `json:",inline"`
	MediaId          string `json:"media_id"`
	MimeType         string `json:"mime_type"`
	Sha256           string `json:"sha256"`
}

type BaseSystemEvent struct {
	Timestamp string `json:"timestamp"`
}

func (bme BaseSystemEvent) GetEventType() string {
	return "system"
}

type BaseBusinessAccountEvent struct {
	Timestamp string `json:"timestamp"`
}

func (bme BaseBusinessAccountEvent) GetEventType() string {
	return "business_account"
}
