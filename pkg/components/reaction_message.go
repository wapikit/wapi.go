package components

import (
	"encoding/json"
	"fmt"

	"github.com/gTahidi/wapi.go/internal"
)

// ReactionMessage represents a reaction to a message.
type ReactionMessage struct {
	MessageId string `json:"message_id" validate:"required"` // The ID of the message to react to.
	Emoji     string `json:"emoji" validate:"required"`      // The emoji representing the reaction.
}

// ReactionMessageParams is an alias for ReactionMessage.
type ReactionMessageParams = ReactionMessage

// ReactionMessageApiPayload represents the API payload for a reaction message.
type ReactionMessageApiPayload struct {
	BaseMessagePayload
	Reaction ReactionMessage `json:"reaction" validate:"required"` // The reaction message.
}

// NewReactionMessage creates a new ReactionMessage instance.
func NewReactionMessage(params ReactionMessageParams) (*ReactionMessage, error) {
	if err := internal.GetValidator().Struct(params); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	return &ReactionMessage{
		MessageId: params.MessageId,
		Emoji:     params.Emoji,
	}, nil
}

// ToJson converts the ReactionMessage to JSON.
func (reaction *ReactionMessage) ToJson(configs ApiCompatibleJsonConverterConfigs) ([]byte, error) {
	if err := internal.GetValidator().Struct(configs); err != nil {
		return nil, fmt.Errorf("error validating configs: %v", err)
	}

	jsonData := ReactionMessageApiPayload{
		BaseMessagePayload: NewBaseMessagePayload(configs.SendToPhoneNumber, MessageTypeReaction),
		Reaction:           *reaction,
	}

	if configs.ReplyToMessageId != "" {
		jsonData.Context = &Context{
			MessageId: configs.ReplyToMessageId,
		}
	}

	jsonToReturn, err := json.Marshal(jsonData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling json: %v", err)
	}

	return jsonToReturn, nil
}
