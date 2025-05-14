package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/pkg/components"
)

// MessageManager is responsible for managing messages.
type MessageManager struct {
	requester     request_client.RequestClient
	PhoneNumberId string
}

// NewMessageManager creates a new instance of MessageManager.
func NewMessageManager(requester request_client.RequestClient, phoneNumberId string) *MessageManager {
	return &MessageManager{
		requester:     requester,
		PhoneNumberId: phoneNumberId,
	}
}

// MessageSendResponse represents the structured API response for sending a message.
type MessageSendResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Error *MessageSendError `json:"error,omitempty"`
}

// MessageSendError represents the error object in an API response.
type MessageSendError struct {
	Message   string `json:"message"` // Error description.
	Type      string `json:"type"`    // Error type (e.g., OAuthException).
	Code      int    `json:"code"`    // Error code.
	ErrorData struct {
		MessagingProduct string `json:"messaging_product"`
		Details          string `json:"details"`
	} `json:"error_data"` // Additional error details.
	ErrorSubcode int    `json:"error_subcode"`
	FbtraceID    string `json:"fbtrace_id"`
}

// Reply sends a reply message using the provided BaseMessage and returns a structured response.
// If the API response contains an error, it returns that error.
func (mm *MessageManager) Reply(message components.BaseMessage, phoneNumber string, replyTo string) (*MessageSendResponse, error) {
	body, err := message.ToJson(components.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: phoneNumber,
		ReplyToMessageId:  replyTo,
	})
	if err != nil {
		return nil, fmt.Errorf("error converting message to json: %v", err)
	}

	// Build the API request.
	apiRequest := mm.requester.NewApiRequest(strings.Join([]string{mm.PhoneNumberId, "messages"}, "/"), http.MethodPost)
	apiRequest.SetBody(string(body))
	responseStr, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	// Unmarshal the API response.
	var sendResponse MessageSendResponse
	err = json.Unmarshal([]byte(responseStr), &sendResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// If an error object is present in the response, return it.
	if sendResponse.Error != nil {
		return &sendResponse, fmt.Errorf("error sending message: %s", sendResponse.Error.Message)
	}

	return &sendResponse, nil
}

// Send sends a message using the provided BaseMessage and returns a structured response.
// If the API response contains an error, it returns that error.
func (mm *MessageManager) Send(message components.BaseMessage, phoneNumber string) (*MessageSendResponse, error) {
	// Convert the message to JSON.
	body, err := message.ToJson(components.ApiCompatibleJsonConverterConfigs{
		SendToPhoneNumber: phoneNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("error converting message to json: %v", err)
	}

	// Build the API request.
	apiRequest := mm.requester.NewApiRequest(strings.Join([]string{mm.PhoneNumberId, "messages"}, "/"), http.MethodPost)
	apiRequest.SetBody(string(body))
	responseStr, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	// Unmarshal the API response.
	var sendResponse MessageSendResponse
	err = json.Unmarshal([]byte(responseStr), &sendResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// If an error object is present in the response, return it.
	if sendResponse.Error != nil {
		return &sendResponse, fmt.Errorf("error sending message: %s", sendResponse.Error.Message)
	}

	return &sendResponse, nil
}
