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

// this register function is for one time registration of the phone number to enable the usage with WhatsApp Cloud API
func (client *MessagingClient) Register(pin string) (RegisterResponse, error) {
	apiRequest := client.Requester.NewApiRequest(strings.Join([]string{client.PhoneNumberId, "resgiter"}, "/"), http.MethodPost)
	apiRequest.AddQueryParam("messaging_product", "WHATSAPP")
	apiRequest.AddQueryParam("pin", pin)
	response, err := apiRequest.Execute()
	if err != nil {
		return RegisterResponse{}, err
	}
	var registerResponse RegisterResponse
	json.Unmarshal([]byte(response), &registerResponse)
	return registerResponse, nil
}

func (client *MessagingClient) Deregister() (RegisterResponse, error) {
	apiRequest := client.Requester.NewApiRequest(strings.Join([]string{client.PhoneNumberId, "deregister"}, "/"), http.MethodPost)
	response, err := apiRequest.Execute()
	if err != nil {
		return RegisterResponse{}, err
	}
	var registerResponse RegisterResponse
	json.Unmarshal([]byte(response), &registerResponse)
	return registerResponse, nil
}
