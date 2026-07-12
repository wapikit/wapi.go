package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Business username management (identity/username rollout). Only the three
// verified username endpoints on the business phone number id are implemented —
// GET/POST/DELETE /<BUSINESS_PHONE_NUMBER_ID>/username. Suggestion/reserved
// endpoints are intentionally NOT invented here.

// BusinessUsernameResponse is the GET /<phone_number_id>/username result.
type BusinessUsernameResponse struct {
	Id       string            `json:"id,omitempty"`
	Username string            `json:"username,omitempty"`
	Error    *MessageSendError `json:"error,omitempty"`
}

// UsernameMutationResponse is the POST/DELETE username result.
type UsernameMutationResponse struct {
	Success bool              `json:"success,omitempty"`
	Error   *MessageSendError `json:"error,omitempty"`
}

// GetUsername fetches the current business username for the phone number id.
// Username is empty when none is set.
func (manager *PhoneNumberManager) GetUsername(phoneNumberId string) (*BusinessUsernameResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumberId, "username"}, "/"), http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var out BusinessUsernameResponse
	if err := json.Unmarshal([]byte(response), &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling username response: %v", err)
	}
	if out.Error != nil {
		return &out, fmt.Errorf("error getting username: %s", out.Error.Message)
	}
	return &out, nil
}

// SetUsername sets (or changes) the business username for the phone number id.
func (manager *PhoneNumberManager) SetUsername(phoneNumberId, username string) (*UsernameMutationResponse, error) {
	body, err := json.Marshal(map[string]string{"username": username})
	if err != nil {
		return nil, fmt.Errorf("error marshalling username body: %v", err)
	}
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumberId, "username"}, "/"), http.MethodPost)
	apiRequest.SetBody(string(body))
	return manager.executeUsernameMutation(apiRequest.Execute())
}

// DeleteUsername removes the business username for the phone number id.
func (manager *PhoneNumberManager) DeleteUsername(phoneNumberId string) (*UsernameMutationResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumberId, "username"}, "/"), http.MethodDelete)
	return manager.executeUsernameMutation(apiRequest.Execute())
}

func (manager *PhoneNumberManager) executeUsernameMutation(response string, execErr error) (*UsernameMutationResponse, error) {
	if execErr != nil {
		return nil, execErr
	}
	var out UsernameMutationResponse
	if err := json.Unmarshal([]byte(response), &out); err != nil {
		return nil, fmt.Errorf("error unmarshalling username mutation response: %v", err)
	}
	if out.Error != nil {
		return &out, fmt.Errorf("error mutating username: %s", out.Error.Message)
	}
	return &out, nil
}
