package manager

import (
	"encoding/json"
	"net/http"
	"strings"

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

type WhatsappBusinessAccountPhoneNumberCodeVerificationStatus string

const (
	WhatsappBusinessAccountPhoneNumberCodeVerificationStatusVerified    WhatsappBusinessAccountPhoneNumberCodeVerificationStatus = "VERIFIED"
	WhatsappBusinessAccountPhoneNumberCodeVerificationStatusNotVerified WhatsappBusinessAccountPhoneNumberCodeVerificationStatus = "NOT_VERIFIED"
	WhatsappBusinessAccountPhoneNumberCodeVerificationStatusExpired     WhatsappBusinessAccountPhoneNumberCodeVerificationStatus = "EXPIRED"
)

type WhatsappBusinessAccountPhoneNumberMessagingLimitTier string

const (
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTier50        WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_50"
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTier250       WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_250"
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTier1K        WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_1K"
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTier10K       WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_10K"
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTier100K      WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_100K"
	WhatsappBusinessAccountPhoneNumberMessagingLimitTierTierUnlimited WhatsappBusinessAccountPhoneNumberMessagingLimitTier = "TIER_UNLIMITED"
)

type WhatsappBusinessAccountPhoneNumberNameStatus string

const (
	WhatsappBusinessAccountPhoneNumberNameStatusApproved               WhatsappBusinessAccountPhoneNumberNameStatus = "APPROVED"
	WhatsappBusinessAccountPhoneNumberNameStatusAvailableWithoutReview WhatsappBusinessAccountPhoneNumberNameStatus = "AVAILABLE_WITHOUT_REVIEW"
	WhatsappBusinessAccountPhoneNumberNameStatusDeclined               WhatsappBusinessAccountPhoneNumberNameStatus = "DECLINED"
	WhatsappBusinessAccountPhoneNumberNameStatusExpired                WhatsappBusinessAccountPhoneNumberNameStatus = "EXPIRED"
	WhatsappBusinessAccountPhoneNumberNameStatusNonExists              WhatsappBusinessAccountPhoneNumberNameStatus = "NON_EXISTS"
	WhatsappBusinessAccountPhoneNumberNameStatusNone                   WhatsappBusinessAccountPhoneNumberNameStatus = "NONE"
	WhatsappBusinessAccountPhoneNumberNameStatusPendingReview          WhatsappBusinessAccountPhoneNumberNameStatus = "PENDING_REVIEW"
)

type WhatsappBusinessAccountPhoneNumberPlatformType string

const (
	WhatsappBusinessAccountPhoneNumberPlatformTypeCloudApi      WhatsappBusinessAccountPhoneNumberPlatformType = "CLOUD_API"
	WhatsappBusinessAccountPhoneNumberPlatformTypeNotApplicable WhatsappBusinessAccountPhoneNumberPlatformType = "NOT_APPLICABLE"
	WhatsappBusinessAccountPhoneNumberPlatformTypeOnPremise     WhatsappBusinessAccountPhoneNumberPlatformType = "ON_PREMISE"
)

type WhatsappBusinessAccountPhoneNumberStatus string

const (
	WhatsappBusinessAccountPhoneNumberStatusPending      WhatsappBusinessAccountPhoneNumberStatus = "PENDING"
	WhatsappBusinessAccountPhoneNumberStatusDeleted      WhatsappBusinessAccountPhoneNumberStatus = "DELETED"
	WhatsappBusinessAccountPhoneNumberStatusMigrated     WhatsappBusinessAccountPhoneNumberStatus = "MIGRATED"
	WhatsappBusinessAccountPhoneNumberStatusBanned       WhatsappBusinessAccountPhoneNumberStatus = "BANNED"
	WhatsappBusinessAccountPhoneNumberStatusRestricted   WhatsappBusinessAccountPhoneNumberStatus = "RESTRICTED"
	WhatsappBusinessAccountPhoneNumberStatusRateLimited  WhatsappBusinessAccountPhoneNumberStatus = "RATE_LIMITED"
	WhatsappBusinessAccountPhoneNumberStatusFlagged      WhatsappBusinessAccountPhoneNumberStatus = "FLAGGED"
	WhatsappBusinessAccountPhoneNumberStatusConnected    WhatsappBusinessAccountPhoneNumberStatus = "CONNECTED"
	WhatsappBusinessAccountPhoneNumberStatusDisconnected WhatsappBusinessAccountPhoneNumberStatus = "DISCONNECTED"
	WhatsappBusinessAccountPhoneNumberStatusUnknown      WhatsappBusinessAccountPhoneNumberStatus = "UNKNOWN"
	WhatsappBusinessAccountPhoneNumberStatusUnverified   WhatsappBusinessAccountPhoneNumberStatus = "UNVERIFIED"
)

type WabaHealthStatus struct {
	CanSendMessage string `json:"can_send_message"`
	Entities       []struct {
		EntityType        string `json:"entity_type"`
		ID                string `json:"id"`
		CanSendMessage    string `json:"can_send_message,omitempty"`
		CanReceiveCallSIP string `json:"can_receive_call_sip,omitempty"`
		Errors            []struct {
			ErrorCode        int    `json:"error_code"`
			ErrorDescription string `json:"error_description"`
			PossibleSolution string `json:"possible_solution"`
		} `json:"errors,omitempty"`
	} `json:"entities"`
}

type WabaQualityScore struct {
	Score string `json:"score"`
}

type WabaThroughput struct {
	Level string `json:"level"`
}

type WabaOfficialBusinessAccount struct {
	OBAStatus string `json:"oba_status"`
}

// WhatsappBusinessAccountPhoneNumber represents a WhatsApp Business Account phone number.
type WhatsappBusinessAccountPhoneNumber struct {
	VerifiedName                          string                                                   `json:"verified_name,omitempty"`
	DisplayPhoneNumber                    string                                                   `json:"display_phone_number,omitempty"`
	Id                                    string                                                   `json:"id,omitempty"`
	QualityRating                         string                                                   `json:"quality_rating,omitempty"`
	CodeVerificationStatus                WhatsappBusinessAccountPhoneNumberCodeVerificationStatus `json:"code_verification_status,omitempty"`
	Status                                WhatsappBusinessAccountPhoneNumberStatus                 `json:"status,omitempty"` // CONNECTED
	PlatformType                          WhatsappBusinessAccountPhoneNumberPlatformType           `json:"platform_type,omitempty"`
	CountryDialCode                       string                                                   `json:"country_dial_code,omitempty"`
	SearchVisibility                      string                                                   `json:"search_visibility,omitempty"`
	Certificate                           string                                                   `json:"certificate,omitempty"`
	HealthStatus                          *WabaHealthStatus                                        `json:"health_status,omitempty"`
	IsOfficialBusinessAccount             bool                                                     `json:"is_official_business_account,omitempty"`
	IsOnBizApp                            bool                                                     `json:"is_on_biz_app,omitempty"`
	IsPinEnabled                          bool                                                     `json:"is_pin_enabled,omitempty"`
	IsPreverifiedNumber                   bool                                                     `json:"is_preverified_number,omitempty"`
	LastOnboardedTime                     string                                                   `json:"last_onboarded_time,omitempty"`
	MessagingLimitTier                    WhatsappBusinessAccountPhoneNumberMessagingLimitTier     `json:"messaging_limit_tier,omitempty"`
	NameStatus                            WhatsappBusinessAccountPhoneNumberNameStatus             `json:"name_status,omitempty"`
	NewCertificate                        string                                                   `json:"new_certificate,omitempty"`
	NewDisplayName                        string                                                   `json:"new_display_name,omitempty"`
	NewNameStatus                         WhatsappBusinessAccountPhoneNumberNameStatus             `json:"new_name_status,omitempty"`
	EligibilityForApiBusinessGlobalSearch string                                                   `json:"eligibility_for_api_business_global_search,omitempty"`
	AccountMode                           string                                                   `json:"account_mode,omitempty"`
	Throughput                            *WabaThroughput                                          `json:"throughput,omitempty"`
	OfficialBusinessAccount               *WabaOfficialBusinessAccount                             `json:"official_business_account,omitempty"`
	QualityScore                          *WabaQualityScore                                        `json:"quality_score,omitempty"`
}

// WhatsappBusinessAccountPhoneNumberEdge represents a list of WhatsApp Business Account phone numbers.
type WhatsappBusinessAccountPhoneNumberEdge struct {
	Data    []WhatsappBusinessAccountPhoneNumber `json:"data,omitempty"`
	Paging  PaginationDetails                    `json:"paging,omitempty"`
	Summary string                               `json:"summary,omitempty"`
}

// FetchAll fetches all phone numbers based on the provided filters.
func (manager *PhoneNumberManager) FetchAll(getSandBoxNumbers bool) (*WhatsappBusinessAccountPhoneNumberEdge, error) {
	apiRequest := manager.requester.NewApiRequest(strings.Join([]string{manager.businessAccountId, "/", "phone_numbers"}, ""), http.MethodGet)

	apiRequest.AddQueryParam("fields", "id,account_mode,certificate,code_verification_status,conversational_automation,display_phone_number,health_status,eligibility_for_api_business_global_search,is_official_business_account,is_on_biz_app,is_pin_enabled,is_preverified_number,last_onboarded_time,messaging_limit_tier,name_status,new_certificate,new_display_name,new_name_status,official_business_account,platform_type,quality_score,search_visibility,status,throughput,verified_name")
	apiRequest.AddQueryParam("filtering", `[{"field":"account_mode","operator":"EQUAL","value":"LIVE"}]`)
	response, err := apiRequest.Execute()

	if err != nil {
		return nil, err
	}

	var responseToReturn WhatsappBusinessAccountPhoneNumberEdge
	err = json.Unmarshal([]byte(response), &responseToReturn)
	if err != nil {
		return nil, err
	}

	return &responseToReturn, nil
}

// Fetch fetches a phone number by its ID.
func (manager *PhoneNumberManager) Fetch(phoneNumberId string) (*WhatsappBusinessAccountPhoneNumber, error) {
	apiRequest := manager.requester.NewApiRequest(phoneNumberId, http.MethodGet)
	apiRequest.AddQueryParam("fields", "id,account_mode,certificate,code_verification_status,conversational_automation,display_phone_number,health_status,eligibility_for_api_business_global_search,is_official_business_account,is_on_biz_app,is_pin_enabled,is_preverified_number,last_onboarded_time,messaging_limit_tier,name_status,new_certificate,new_display_name,new_name_status,official_business_account,platform_type,quality_score,search_visibility,status,throughput,verified_name")
	response, err := apiRequest.Execute()

	if err != nil {
		return nil, err
	}

	var responseToReturn WhatsappBusinessAccountPhoneNumber
	err = json.Unmarshal([]byte(response), &responseToReturn)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return CreatePhoneNumberResponse{}, err
	}
	responseToReturn := CreatePhoneNumberResponse{}
	err = json.Unmarshal([]byte(response), &responseToReturn)
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
