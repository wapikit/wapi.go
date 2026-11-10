package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

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
	Id                string `json:"id" validate:"required"`
	Name              string `json:"name" validate:"required"`
	TimezoneId        string `json:"timezone_id" validate:"required"`
	Currency          string `json:"currency" validate:"required"`
	OwnerBusinessInfo struct {
		Name string `json:"name" validate:"required"`
		Id   string `json:"id" validate:"required"`
	} `json:"owner_business_info" validate:"required"`
	MarketingMessagesOnboardingStatus string `json:"marketing_messages_onboarding_status" validate:"required"`
	MessageTemplateNamespace          string `json:"message_template_namespace" validate:"required"`
	AccountReviewStatus               string `json:"account_review_status" validate:"required"`
	BusinessVerificationStatus        string `json:"business_verification_status" validate:"required"`
	Country                           string `json:"country" validate:"required"`
	OwnershipType                     string `json:"ownership_type" validate:"required"`
	PrimaryBusinessLocation           string `json:"primary_business_location" validate:"required"`
}

// This method fetches the business account details.
func (client *BusinessClient) Fetch() (*FetchBusinessAccountResponse, error) {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	fields := []string{
		"id",
		"name",
		"timezone_id",
		"message_template_namespace",
		"account_review_status",
		"business_verification_status",
		"country",
		"ownership_type",
		"primary_business_location",
		"currency",
		"owner_business_info",
		"marketing_messages_onboarding_status",
	}

	for _, field := range fields {
		apiRequest.AddField(request_client.ApiRequestQueryParamField{
			Name:    field,
			Filters: map[string]string{},
		})
	}

	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var responseToReturn FetchBusinessAccountResponse
	if err := json.Unmarshal([]byte(response), &responseToReturn); err != nil {
		// A field whose JSON type differs from the struct (Meta changing a
		// field's shape) still leaves every other field populated. Earlier
		// releases ignored that case and returned the partial result; keep doing
		// so, and only fail on a body that is not JSON at all.
		var typeErr *json.UnmarshalTypeError
		if !errors.As(err, &typeErr) {
			return nil, fmt.Errorf("error unmarshalling business account response: %w", err)
		}
	}
	return &responseToReturn, nil
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

// FetchAnalytics fetches message analytics (sent and delivered counts) for
// the business account.
func (client *BusinessClient) FetchAnalytics(options AccountAnalyticsOptions) (WhatsappBusinessAccountAnalyticsResponse, error) {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	analyticsField := apiRequest.AddField(request_client.ApiRequestQueryParamField{
		Name:    "analytics",
		Filters: map[string]string{},
	})
	analyticsField.AddFilter("start", fmt.Sprint(options.Start.Unix()))
	analyticsField.AddFilter("end", fmt.Sprint(options.End.Unix()))
	analyticsField.AddFilter("granularity", string(options.Granularity))
	// An empty list means every phone number / country.
	analyticsField.AddFilter("phone_numbers", graphListParam(options.PhoneNumbers))
	analyticsField.AddFilter("country_codes", graphListParam(options.CountryCodes))

	response, err := apiRequest.Execute()
	if err != nil {
		return WhatsappBusinessAccountAnalyticsResponse{}, err
	}
	// Meta nests the result under the requested field name.
	var envelope struct {
		Analytics WhatsappBusinessAccountAnalyticsResponse `json:"analytics"`
	}
	if err := json.Unmarshal([]byte(response), &envelope); err != nil {
		return WhatsappBusinessAccountAnalyticsResponse{}, fmt.Errorf("error unmarshalling analytics response: %w", err)
	}
	return envelope.Analytics, nil
}

// graphListParam renders values as a JSON array, the form Graph API field
// filters such as phone_numbers(["1","2"]) expect.
func graphListParam[T ~string](values []T) string {
	if len(values) == 0 {
		return "[]"
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(encoded)
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
	// Cost is rounded to a whole number; ConversationAnalyticsDataPoint.Cost
	// keeps the exact value.
	Cost int `json:"cost,omitempty"`
}

type WhatsAppConversationAnalyticsEdge struct {
	Data []struct {
		DataPoints []WhatsAppConversationAnalyticsNode `json:"data_points,omitempty"`
	} `json:"data,omitempty"`
	Paging manager.PaginationDetails `json:"paging,omitempty"`
}

type WhatsAppConversationAnalyticsResponse struct {
	ConversationAnalytics []WhatsAppConversationAnalyticsEdge `json:"conversation_analytics" validate:"required"`
}

// ConversationAnalyticsDataPoint is one bucket of conversation analytics.
type ConversationAnalyticsDataPoint struct {
	Start                 int64   `json:"start"`
	End                   int64   `json:"end"`
	Conversation          int     `json:"conversation,omitempty"`
	PhoneNumber           string  `json:"phone_number,omitempty"`
	Country               string  `json:"country,omitempty"`
	ConversationType      string  `json:"conversation_type,omitempty"`
	ConversationDirection string  `json:"conversation_direction,omitempty"`
	ConversationCategory  string  `json:"conversation_category,omitempty"`
	Cost                  float64 `json:"cost,omitempty"`
}

// ConversationAnalyticsResult is the conversation_analytics object Meta
// returns: one or more groups of data points, plus paging.
type ConversationAnalyticsResult struct {
	Data []struct {
		DataPoints []ConversationAnalyticsDataPoint `json:"data_points,omitempty"`
	} `json:"data,omitempty"`
	Paging manager.PaginationDetails `json:"paging,omitempty"`
}

// FetchConversationAnalytics fetches conversation counts and cost for the
// business account.
func (client *BusinessClient) FetchConversationAnalytics(options ConversationAnalyticsOptions) (*ConversationAnalyticsResult, error) {
	apiRequest := client.requester.NewApiRequest(client.BusinessAccountId, http.MethodGet)
	analyticsField := apiRequest.AddField(request_client.ApiRequestQueryParamField{
		Name:    "conversation_analytics",
		Filters: map[string]string{},
	})
	analyticsField.AddFilter("start", fmt.Sprint(options.Start.Unix()))
	analyticsField.AddFilter("end", fmt.Sprint(options.End.Unix()))
	analyticsField.AddFilter("granularity", string(options.Granularity))
	// An empty list means every phone number.
	analyticsField.AddFilter("phone_numbers", graphListParam(options.PhoneNumbers))
	// The remaining filters are left out when empty, which Meta treats as "all".
	if len(options.ConversationCategory) > 0 {
		analyticsField.AddFilter("conversation_categories", graphListParam(options.ConversationCategory))
	}
	if len(options.ConversationTypes) > 0 {
		analyticsField.AddFilter("conversation_types", graphListParam(options.ConversationTypes))
	}
	if len(options.ConversationDirection) > 0 {
		analyticsField.AddFilter("conversation_directions", graphListParam(options.ConversationDirection))
	}
	if len(options.Dimensions) > 0 {
		analyticsField.AddFilter("dimensions", graphListParam(options.Dimensions))
	}

	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var envelope struct {
		ConversationAnalytics ConversationAnalyticsResult `json:"conversation_analytics"`
	}
	if err := json.Unmarshal([]byte(response), &envelope); err != nil {
		return nil, fmt.Errorf("error unmarshalling conversation analytics response: %w", err)
	}
	return &envelope.ConversationAnalytics, nil
}

// ConversationAnalytics fetches the conversation analytics for the business account.
//
// Deprecated: use FetchConversationAnalytics, which keeps fractional costs.
// This method returns the same data with Cost rounded to a whole number.
func (client *BusinessClient) ConversationAnalytics(options ConversationAnalyticsOptions) (*WhatsAppConversationAnalyticsResponse, error) {
	result, err := client.FetchConversationAnalytics(options)
	if err != nil {
		return nil, err
	}
	edge := WhatsAppConversationAnalyticsEdge{Paging: result.Paging}
	for _, group := range result.Data {
		nodes := make([]WhatsAppConversationAnalyticsNode, len(group.DataPoints))
		for i, point := range group.DataPoints {
			nodes[i] = WhatsAppConversationAnalyticsNode{
				Start:                 int(point.Start),
				End:                   int(point.End),
				Conversation:          point.Conversation,
				PhoneNumber:           point.PhoneNumber,
				Country:               point.Country,
				ConversationType:      point.ConversationType,
				ConversationDirection: point.ConversationDirection,
				ConversationCategory:  point.ConversationCategory,
				Cost:                  int(math.Round(point.Cost)),
			}
		}
		edge.Data = append(edge.Data, struct {
			DataPoints []WhatsAppConversationAnalyticsNode `json:"data_points,omitempty"`
		}{DataPoints: nodes})
	}
	return &WhatsAppConversationAnalyticsResponse{ConversationAnalytics: []WhatsAppConversationAnalyticsEdge{edge}}, nil
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
