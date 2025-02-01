package manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wapikit/wapi.go/internal"
	"github.com/wapikit/wapi.go/internal/request_client"
)

// PhoneNumberManager is responsible for managing phone numbers for WhatsApp Business API and phone number specific operations.
type PhoneNumberManager struct {
	businessAccountId string
	apiAccessToken    string
	requester         *request_client.RequestClient
}

// PhoneNumberManagerConfig holds the configuration for PhoneNumberManager.
type PhoneNumberManagerConfig struct {
	BusinessAccountId string
	ApiAccessToken    string
	Requester         *request_client.RequestClient
}

// NewPhoneNumberManager creates a new instance of PhoneNumberManager.
func NewPhoneNumberManager(config *PhoneNumberManagerConfig) *PhoneNumberManager {
	return &PhoneNumberManager{
		apiAccessToken:    config.ApiAccessToken,
		businessAccountId: config.BusinessAccountId,
		requester:         config.Requester,
	}
}

// WhatsappBusinessAccountPhoneNumber represents a WhatsApp Business Account phone number.
type WhatsappBusinessAccountPhoneNumber struct {
	VerifiedName       string `json:"verified_name,omitempty"`
	DisplayPhoneNumber string `json:"display_phone_number,omitempty"`
	Id                 string `json:"id,omitempty"`
	QualityRating      string `json:"quality_rating,omitempty"`
	CodeVerification   struct {
		Status string `json:"code_verification_status,omitempty"`
	} `json:"code_verification_status,omitempty"`
	PlatformType string `json:"platform_type,omitempty"`
}

// WhatsappBusinessAccountPhoneNumberEdge represents a list of WhatsApp Business Account phone numbers.
type WhatsappBusinessAccountPhoneNumberEdge struct {
	Data    []WhatsappBusinessAccountPhoneNumber       `json:"data,omitempty"`
	Paging  internal.WhatsAppBusinessApiPaginationMeta `json:"paging,omitempty"`
	Summary string                                     `json:"summary,omitempty"`
}

// FetchAll fetches all phone numbers based on the provided filters.
func (manager *PhoneNumberManager) FetchAll(getSandBoxNumbers bool) (*WhatsappBusinessAccountPhoneNumberEdge, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "/", "phone_numbers"}, ""), http.MethodGet)

	apiRequest.AddQueryParam("filtering", `[{"field":"account_mode","operator":"EQUAL","value":"LIVE"}]`)
	response, err := apiRequest.Execute()

	if err != nil {
		return nil, err
	}

	var responseToReturn WhatsappBusinessAccountPhoneNumberEdge
	json.Unmarshal([]byte(response), &responseToReturn)

	return &responseToReturn, nil
}

// Fetch fetches a phone number by its ID.
func (manager *PhoneNumberManager) Fetch(phoneNumberId string) (*WhatsappBusinessAccountPhoneNumber, error) {
	apiRequest := manager.requester.NewApiRequest(phoneNumberId, http.MethodGet)

	response, err := apiRequest.Execute()

	if err != nil {
		return nil, err
	}

	var responseToReturn WhatsappBusinessAccountPhoneNumber
	json.Unmarshal([]byte(response), &responseToReturn)

	return &responseToReturn, nil
}

type CreatePhoneNumberResponse struct {
	Id string `json:"id,omitempty"`
}

func (manager *PhoneNumberManager) Create(phoneNumber, verifiedName, countryCode string) (CreatePhoneNumberResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "/phone_numbers"}, ""), http.MethodPost)
	apiRequest.AddQueryParam("phone_number", phoneNumber)
	apiRequest.AddQueryParam("cc", countryCode)
	apiRequest.AddQueryParam("verified_name", verifiedName)
	response, err := apiRequest.Execute()
	responseToReturn := CreatePhoneNumberResponse{}
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn, err
}

type VerifyCodeMethod string

const (
	VerifyCodeMethodSms   VerifyCodeMethod = "SMS"
	VerifyCodeMethodVoice VerifyCodeMethod = "VOICE"
)

type RequestVerificationCodeResponse struct {
	Success bool `json:"success,omitempty"`
}

func (manager *PhoneNumberManager) RequestVerificationCode(phoneNumberId string, codeMethod VerifyCodeMethod, languageCode string) (RequestVerificationCodeResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumberId, "request_code"}, "/"), http.MethodPost)
	apiRequest.AddQueryParam("code_method", string(codeMethod))
	apiRequest.AddQueryParam("language", languageCode)
	response, err := apiRequest.Execute()
	responseToReturn := RequestVerificationCodeResponse{}
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn, err
}

type VerifyCodeResponse struct {
	Success bool `json:"success,omitempty"`
}

func (manager *PhoneNumberManager) VerifyCode(phoneNumberId, verificationCode string) (VerifyCodeResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumberId, "verify_code"}, "/"), http.MethodPost)
	apiRequest.AddQueryParam("code", verificationCode)
	response, err := apiRequest.Execute()
	responseToReturn := VerifyCodeResponse{}
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn, err
}

// GenerateQrCodeResponse represents the response of generating a QR code.
type GenerateQrCodeResponse struct {
	Code             string `json:"code,omitempty"`
	PrefilledMessage string `json:"prefilled_message,omitempty"`
	DeepLinkUrl      string `json:"deep_link_url,omitempty"`
	QrImageUrl       string `json:"qr_image_url,omitempty"`
}

// GenerateQrCode generates a QR code for the specified phone number with the given prefilled message.
func (manager *PhoneNumberManager) GenerateQrCode(phoneNumber string, prefilledMessage string) (*GenerateQrCodeResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumber, "/message_qrdls"}, ""), http.MethodPost)
	jsonBody, err := json.Marshal(map[string]string{
		"prefilled_message": prefilledMessage,
		"generate_qr_image": "PNG",
	})
	if err != nil {
		return nil, err
	}
	apiRequest.SetBody(string(jsonBody))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var responseToReturn GenerateQrCodeResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// GetAllQrCodesResponse represents the response of getting all QR codes for a phone number.
type GetAllQrCodesResponse struct {
	Data []GenerateQrCodeResponse `json:"data,omitempty"`
}

// GetAllQrCodes gets all QR codes for the specified phone number.
func (manager *PhoneNumberManager) GetAllQrCodes(phoneNumber string) (*GetAllQrCodesResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumber, "/message_qrdls"}, ""), http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	var responseToReturn GetAllQrCodesResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// GetQrCodeById gets a QR code by its ID for the specified phone number.
func (manager *PhoneNumberManager) GetQrCodeById(phoneNumber, id string) (*GetAllQrCodesResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumber, "/message_qrdls", "/", id}, ""), http.MethodDelete)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var responseToReturn GetAllQrCodesResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// DeleteQrCodeResponse represents the response of deleting a QR code.
type DeleteQrCodeResponse struct {
	Success bool `json:"success,omitempty"`
}

// DeleteQrCode deletes a QR code by its ID for the specified phone number.
func (manager *PhoneNumberManager) DeleteQrCode(phoneNumber, id string) (*DeleteQrCodeResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumber, "/message_qrdls", "/", id}, ""), http.MethodDelete)
	response, err := apiRequest.Execute()

	if err != nil {
		return nil, err
	}
	var responseToReturn DeleteQrCodeResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}

// UpdateQrCode updates a QR code by its ID for the specified phone number with the given prefilled message.
func (manager *PhoneNumberManager) UpdateQrCode(phoneNumber, id, prefilledMessage string) (*GenerateQrCodeResponse, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{phoneNumber, "/message_qrdls"}, ""), http.MethodPost)
	jsonBody, err := json.Marshal(map[string]string{
		"prefilled_message": prefilledMessage,
		"code":              id,
	})
	if err != nil {
		return nil, err
	}
	apiRequest.SetBody(string(jsonBody))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

	var responseToReturn GenerateQrCodeResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return &responseToReturn, nil
}
