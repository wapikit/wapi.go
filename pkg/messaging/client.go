package messaging

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/manager"
)

// MessagingClient represents a WhatsApp client.
type MessagingClient struct {
	Media             manager.MediaManager
	Message           manager.MessageManager
	PhoneNumberId     string
	ApiAccessToken    string
	BusinessAccountId string
	Requester         *request_client.RequestClient
}

// GetPhoneNumberId returns the phone number ID associated with the client.
func (client *MessagingClient) GetPhoneNumberId() string {
	return client.PhoneNumberId
}

// SetPhoneNumberId sets the phone number ID for the client.
func (client *MessagingClient) SetPhoneNumberId(phoneNumberId string) {
	client.PhoneNumberId = phoneNumberId
}

func (client *MessagingClient) GetApiAccessToken() string {
	return client.ApiAccessToken
}

func (client *MessagingClient) SetApiAccessToken(apiAccessToken string) {
	client.ApiAccessToken = apiAccessToken
}

func (client *MessagingClient) GetBusinessAccountId() string {
	return client.BusinessAccountId
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

// Register registers the phone number for use with the WhatsApp Cloud API.
// pin is the six-digit two-step verification PIN (set or reused by this call).
func (client *MessagingClient) Register(pin string) (RegisterResponse, error) {
	body, err := json.Marshal(map[string]string{
		"messaging_product": "whatsapp",
		"pin":               pin,
	})
	if err != nil {
		return RegisterResponse{}, err
	}
	apiRequest := client.Requester.NewApiRequest(strings.Join([]string{client.PhoneNumberId, "register"}, "/"), http.MethodPost)
	apiRequest.SetBody(string(body))
	return executeRegisterRequest(apiRequest)
}

// Deregister deregisters the phone number from the WhatsApp Cloud API.
func (client *MessagingClient) Deregister() (RegisterResponse, error) {
	apiRequest := client.Requester.NewApiRequest(strings.Join([]string{client.PhoneNumberId, "deregister"}, "/"), http.MethodPost)
	return executeRegisterRequest(apiRequest)
}

func executeRegisterRequest(apiRequest *request_client.ApiRequest) (RegisterResponse, error) {
	response, err := apiRequest.Execute()
	if err != nil {
		return RegisterResponse{}, err
	}
	var registerResponse RegisterResponse
	if err := json.Unmarshal([]byte(response), &registerResponse); err != nil {
		return RegisterResponse{}, err
	}
	return registerResponse, nil
}
