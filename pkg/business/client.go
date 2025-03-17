package business

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wapikit/wapi.go/internal"
	"github.com/wapikit/wapi.go/internal/request_client"
	"github.com/wapikit/wapi.go/manager"
)

// BusinessClient is responsible for managing business account related operations.
type BusinessClient struct {
	BusinessAccountId string `json:"businessAccountId" validate:"required"`
	AccessToken       string `json:"accessToken" validate:"required"`
	PhoneNumber       *manager.PhoneNumberManager
	Template          *manager.TemplateManager
	requester         *request_client.RequestClient
	Catalog           *manager.CatalogManager
}

// BusinessClientConfig holds the configuration for BusinessClient.
type BusinessClientConfig struct {
	BusinessAccountId string `json:"businessAccountId" validate:"required"`
	AccessToken       string `json:"accessToken" validate:"required"`
	Requester         *request_client.RequestClient
}

// NewBusinessClient creates a new instance of BusinessClient.
func NewBusinessClient(config *BusinessClientConfig) *BusinessClient {
	return &BusinessClient{
		BusinessAccountId: config.BusinessAccountId,
		AccessToken:       config.AccessToken,
		PhoneNumber: manager.NewPhoneNumberManager(&manager.PhoneNumberManagerConfig{
			BusinessAccountId: config.BusinessAccountId,
			ApiAccessToken:    config.AccessToken,
			Requester:         config.Requester,
		}),
		Template: manager.NewTemplateManager(&manager.TemplateManagerConfig{
			BusinessAccountId: config.BusinessAccountId,
			ApiAccessToken:    config.AccessToken,
			Requester:         config.Requester,
		}),
		Catalog: manager.NewCatalogManager(&manager.CatalogManagerConfig{
			BusinessAccountId: config.BusinessAccountId,
			Requester:         config.Requester,
		}),
		requester: config.Requester,
	}
}

// GetBusinessId returns the business account ID.
func (bc *BusinessClient) GetBusinessId() string {
	return bc.BusinessAccountId
}

// SetBusinessId sets the business account ID.
func (bc *BusinessClient) SetBusinessId(id string) {
	bc.BusinessAccountId = id
}

// WhatsappBusinessAccount represents a WhatsApp Business Account.
type WhatsappBusinessAccount struct {
	BusinessVerificationStatus string `json:"business_verification_status,omitempty"`
	Country                    string `json:"country,omitempty"`
	Currency                   string `json:"currency,omitempty"`
	IsTemplateAnalyticsEnabled string `json:"is_enabled_for_insights,omitempty"`
	MessageTemplateNamespace   string `json:"message_template_namespace,omitempty"`
	Name                       string `json:"name,omitempty"`
	OwnershipType              string `json:"ownership_type,omitempty"`
	PrimaryFundingId           string `json:"primary_funding_id,omitempty"`
	PurchaseOrderNumber        string `json:"purchase_order_number,omitempty"`
	TimezoneId                 string `json:"timezone_id,omitempty"`
}

type FetchBusinessAccountResponse struct {
	Id                       string `json:"id" validate:"required"`
	Name                     string `json:"name" validate:"required"`
	TimezoneId               string `json:"timezone_id" validate:"required"`
	MessageTemplateNamespace string `json:"message_template_namespace" validate:"required"`
}

// This method fetches the business account details.
func (client *BusinessClient) Fetch() FetchBusinessAccountResponse {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		// return wapi.go custom error here
		fmt.Println("Error while fetching business account", err)
	}
	var responseToReturn FetchBusinessAccountResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn
}

type AnalyticsRequestGranularityType string

const (
	AnalyticsRequestGranularityTypeHalfHour AnalyticsRequestGranularityType = "HALF_HOUR"
	AnalyticsRequestGranularityTypeDay      AnalyticsRequestGranularityType = "DAY"
	AnalyticsRequestGranularityTypeMonth    AnalyticsRequestGranularityType = "MONTH"
)

type WhatsAppBusinessAccountAnalyticsProductType int

const (
	WhatsAppBusinessAccountAnalyticsProductTypeNotificationMessages    WhatsAppBusinessAccountAnalyticsProductType = 0
	WhatsAppBusinessAccountAnalyticsProductTypeCustomerSupportMessages WhatsAppBusinessAccountAnalyticsProductType = 2
)

type AccountAnalyticsOptions struct {
	Start        time.Time                       `json:"start" validate:"required"`
	End          time.Time                       `json:"end" validate:"required"`
	Granularity  AnalyticsRequestGranularityType `json:"granularity" validate:"required"`
	PhoneNumbers []string                        `json:"phone_numbers,omitempty"`

	// * NOT SUPPORTED AS OF NOW
	// ProductTypes []WhatsAppBusinessAccountAnalyticsProductType `json:"product_types,omitempty"`
	CountryCodes []string `json:"country_codes,omitempty"`
}

type AnalyticsDataPoint struct {
	Start     int `json:"start,omitempty"`
	End       int `json:"end,omitempty"`
	Sent      int `json:"sent,omitempty"`
	Delivered int `json:"delivered,omitempty"`
}

type WhatsappBusinessAccountAnalyticsResponse struct {
	PhoneNumbers []string             `json:"phone_numbers,omitempty"`
	Granularity  string               `json:"granularity,omitempty"`
	DataPoints   []AnalyticsDataPoint `json:"data_points,omitempty"`
}

// FetchAnalytics fetches the analytics for the business account.
func (client *BusinessClient) FetchAnalytics(options AccountAnalyticsOptions) (WhatsappBusinessAccountAnalyticsResponse, error) {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	analyticsField := apiRequest.AddField(request_client.ApiRequestQueryParamField{
		Name:    "analytics",
		Filters: map[string]string{},
	})
	analyticsField.AddFilter("start", fmt.Sprint(options.Start.Unix()))
	analyticsField.AddFilter("end", fmt.Sprint(options.End.Unix()))
	analyticsField.AddFilter("granularity", string(options.Granularity))

	if len(options.PhoneNumbers) > 0 {
		// get specific phone numbers
		analyticsField.AddFilter("phone_numbers", strings.Join(options.PhoneNumbers, ","))
	} else {
		// get all phone numbers
		analyticsField.AddFilter("phone_numbers", "[]")
	}

	if len(options.CountryCodes) > 0 {
		analyticsField.AddFilter("country_codes", strings.Join(options.CountryCodes, ","))
	} else {
		// get all country codes
		analyticsField.AddFilter("country_codes", "[]")
	}
	response, err := apiRequest.Execute()
	if err != nil {
		// return wapi.go custom error here
		fmt.Println("Error while fetching business account", err)
	}
	var responseToReturn WhatsappBusinessAccountAnalyticsResponse
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn, nil
}

type ConversationCategoryType string

const (
	ConversationCategoryTypeAuthentication ConversationCategoryType = "AUTHENTICATION"
	ConversationCategoryTypeMarketing      ConversationCategoryType = "MARKETING"
	ConversationCategoryTypeService        ConversationCategoryType = "SERVICE"
	ConversationCategoryTypeUtility        ConversationCategoryType = "UTILITY"
)

type ConversationType string

const (
	ConversationTypeFreeEntry ConversationType = "FREE_ENTRY"
	ConversationTypeFreeTier  ConversationType = "FREE_TIER"
	ConversationTypeRegular   ConversationType = "REGULAR"
)

type ConversationDirection string

const (
	ConversationDirectionBusinessInitiated ConversationDirection = "BUSINESS_INITIATED"
	ConversationDirectionUserInitiated     ConversationDirection = "USER_INITIATED"
)

type ConversationDimensionType string

const (
	ConversationDimensionTypeConversationCategory  ConversationDimensionType = "CONVERSATION_CATEGORY"
	ConversationDimensionTypeConversationDirection ConversationDimensionType = "CONVERSATION_DIRECTION"
	ConversationDimensionTypeConversationType      ConversationDimensionType = "CONVERSATION_TYPE"
	ConversationDimensionTypeCountry               ConversationDimensionType = "COUNTRY"
	ConversationDimensionTypePhone                 ConversationDimensionType = "PHONE"
)

type ConversationAnalyticsGranularityType string

const (
	ConversationAnalyticsGranularityTypeHalfHour ConversationAnalyticsGranularityType = "HALF_HOUR"
	ConversationAnalyticsGranularityTypeDay      ConversationAnalyticsGranularityType = "DAILY"
	ConversationAnalyticsGranularityTypeMonth    ConversationAnalyticsGranularityType = "MONTHLY"
)

type ConversationAnalyticsOptions struct {
	Start        time.Time                            `json:"start" validate:"required"`
	End          time.Time                            `json:"end" validate:"required"`
	Granularity  ConversationAnalyticsGranularityType `json:"granularity" validate:"required"`
	PhoneNumbers []string                             `json:"phone_numbers,omitempty"`

	ConversationCategory  []ConversationCategoryType  `json:"conversation_category,omitempty"`
	ConversationTypes     []ConversationCategoryType  `json:"conversation_types,omitempty"`
	ConversationDirection []ConversationDirection     `json:"conversation_direction,omitempty"`
	Dimensions            []ConversationDimensionType `json:"dimensions,omitempty"`
}

type WhatsAppConversationAnalyticsNode struct {
	Start                 int    `json:"start" validate:"required"`
	End                   int    `json:"end,omitempty" validate:"required"`
	Conversation          int    `json:"conversation,omitempty"`
	PhoneNumber           string `json:"phone_number,omitempty"`
	Country               string `json:"country,omitempty"`
	ConversationType      string `json:"conversation_type,omitempty"`
	ConversationDirection string `json:"conversation_direction,omitempty"`
	ConversationCategory  string `json:"conversation_category,omitempty"`
	Cost                  int    `json:"cost,omitempty"`
}

type WhatsAppConversationAnalyticsEdge struct {
	Data []struct {
		DataPoints []WhatsAppConversationAnalyticsNode `json:"data_points,omitempty"`
	} `json:"data,omitempty"`
	Paging internal.WhatsAppBusinessApiPaginationMeta `json:"paging,omitempty"`
}

type WhatsAppConversationAnalyticsResponse struct {
	ConversationAnalytics []WhatsAppConversationAnalyticsEdge `json:"conversation_analytics" validate:"required"`
}

// ConversationAnalytics fetches the conversation analytics for the business account.
func (client *BusinessClient) ConversationAnalytics(options ConversationAnalyticsOptions) (*WhatsAppConversationAnalyticsResponse, error) {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	analyticsField := apiRequest.AddField(request_client.ApiRequestQueryParamField{
		Name:    "conversation_analytics",
		Filters: map[string]string{},
	})
	analyticsField.AddFilter("start", fmt.Sprint(options.Start.Unix()))
	analyticsField.AddFilter("end", fmt.Sprint(options.End.Unix()))
	analyticsField.AddFilter("granularity", string(options.Granularity))

	if len(options.PhoneNumbers) > 0 {
		// get specific phone numbers
		analyticsField.AddFilter("phone_numbers", strings.Join(options.PhoneNumbers, ","))
	} else {
		// get all phone numbers
		analyticsField.AddFilter("phone_numbers", "[]")
	}

	if len(options.ConversationCategory) > 0 {
		categoryStrings := make([]string, len(options.ConversationCategory))
		for i, category := range options.ConversationCategory {
			categoryStrings[i] = string(category)
		}
		analyticsField.AddFilter("conversation_category", strings.Join(categoryStrings, ","))
	} else {
		analyticsField.AddFilter("conversation_category", "[]") // Empty slice
	}

	if len(options.ConversationTypes) > 0 {
		typeStrings := make([]string, len(options.ConversationTypes))
		for i, ctype := range options.ConversationTypes {
			typeStrings[i] = string(ctype)
		}
		analyticsField.AddFilter("conversation_types", strings.Join(typeStrings, ","))
	} else {
		analyticsField.AddFilter("conversation_types", "[]") // Empty slice
	}

	if len(options.ConversationDirection) > 0 {
		directionStrings := make([]string, len(options.ConversationDirection))
		for i, direction := range options.ConversationDirection {
			directionStrings[i] = string(direction)
		}
		analyticsField.AddFilter("conversation_direction", strings.Join(directionStrings, ","))
	} else {
		analyticsField.AddFilter("conversation_direction", "[]") // Empty slice
	}

	if len(options.Dimensions) > 0 {
		dimensionsStrings := make([]string, len(options.Dimensions))
		for i, dim := range options.Dimensions {
			dimensionsStrings[i] = string(dim)
		}
		analyticsField.AddFilter("dimensions", strings.Join(dimensionsStrings, ","))
	} else {
		// get all country codes
		analyticsField.AddFilter("dimensions", "[]")
	}

	response, err := apiRequest.Execute()
	if err != nil {
		// return wapi.go custom error here
		fmt.Println("Error while fetching business account", err)
	}
	var responseToReturn WhatsAppConversationAnalyticsResponse
	json.Unmarshal([]byte(response), &responseToReturn)

	fmt.Println("Response to return is", responseToReturn)

	return &responseToReturn, nil
}

type BusinessRole string

const (
	BusinessRoleManage               BusinessRole = "MANAGE"
	BusinessRoleDevelop              BusinessRole = "DEVELOP"
	BusinessRoleManageTemplates      BusinessRole = "MANAGE_TEMPLATES"
	BusinessRoleManagePhone          BusinessRole = "MANAGE_PHONE"
	BusinessRoleViewCost             BusinessRole = "VIEW_COST"
	BusinessRoleManageExtensions     BusinessRole = "MANAGE_EXTENSIONS"
	BusinessRoleViewPhoneAssets      BusinessRole = "VIEW_PHONE_ASSETS"
	BusinessRoleManagePhoneAssets    BusinessRole = "MANAGE_PHONE_ASSETS"
	BusinessRoleViewTemplates        BusinessRole = "VIEW_TEMPLATES"
	BusinessRoleMessaging            BusinessRole = "MESSAGING"
	BusinessRoleManageBusinessPhones BusinessRole = "MANAGE_BUSINESS_PHONES"
)
