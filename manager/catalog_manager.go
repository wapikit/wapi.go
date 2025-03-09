package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/wapikit/wapi.go/internal"
	"github.com/wapikit/wapi.go/internal/request_client"
)

type CatalogManager struct {
	requester         *request_client.RequestClient
	businessAccountId string
}

type CatalogManagerConfig struct {
	BusinessAccountId string
	Requester         *request_client.RequestClient
}

func NewCatalogManager(config *CatalogManagerConfig) *CatalogManager {
	return &CatalogManager{
		requester:         config.Requester,
		businessAccountId: config.BusinessAccountId,
	}
}

// New helper type for key/value pairs.
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Applinks struct {
	Web WebApplink `json:"web"`
}

type WebApplink struct {
	ShouldFallback bool   `json:"should_fallback"`
	URL            string `json:"url"`
}

type CapabilityReview struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ImageCdnUrl struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductGroup struct {
	ID         string `json:"id"`
	RetailerID string `json:"retailer_id"`
}

// ProductItem represents a product in a catalog.

type ProductItem struct {
	Id                                                    string        `json:"id"`
	AdditionalImageCdnUrls                                []ImageCdnUrl `json:"additional_image_cdn_urls,omitempty"`
	AdditionalImageUrls                                   []string      `json:"additional_image_urls,omitempty"`
	AdditionalVariantAttributes                           []KeyValue    `json:"additional_variant_attributes,omitempty"`
	AgeGroup                                              string        `json:"age_group,omitempty"`
	Applinks                                              *Applinks     `json:"applinks,omitempty"`
	Availability                                          string        `json:"availability,omitempty"`
	Brand                                                 string        `json:"brand,omitempty"`
	CapabilityToReviewStatus                              []KeyValue    `json:"capability_to_review_status,omitempty"`
	Category                                              string        `json:"category,omitempty"`
	CategorySpecificFields                                string        `json:"category_specific_fields,omitempty"`
	Color                                                 string        `json:"color,omitempty"`
	CommerceInsights                                      string        `json:"commerce_insights,omitempty"`
	Condition                                             string        `json:"condition,omitempty"`
	Currency                                              string        `json:"currency,omitempty"`
	CustomData                                            []KeyValue    `json:"custom_data,omitempty"`
	CustomLabel0                                          string        `json:"custom_label_0,omitempty"`
	CustomLabel1                                          string        `json:"custom_label_1,omitempty"`
	CustomLabel2                                          string        `json:"custom_label_2,omitempty"`
	CustomLabel3                                          string        `json:"custom_label_3,omitempty"`
	CustomLabel4                                          string        `json:"custom_label_4,omitempty"`
	CustomNumber0                                         string        `json:"custom_number_0,omitempty"`
	CustomNumber1                                         string        `json:"custom_number_1,omitempty"`
	CustomNumber2                                         string        `json:"custom_number_2,omitempty"`
	CustomNumber3                                         string        `json:"custom_number_3,omitempty"`
	CustomNumber4                                         string        `json:"custom_number_4,omitempty"`
	Description                                           string        `json:"description,omitempty"`
	Errors                                                []string      `json:"errors,omitempty"`
	ExpirationDate                                        string        `json:"expiration_date,omitempty"`
	FbProductCategory                                     string        `json:"fb_product_category,omitempty"`
	Gender                                                string        `json:"gender,omitempty"`
	Gtin                                                  string        `json:"gtin,omitempty"`
	ImageCdnUrls                                          []ImageCdnUrl `json:"image_cdn_urls,omitempty"`
	ImageFetchStatus                                      string        `json:"image_fetch_status,omitempty"`
	ImageUrl                                              string        `json:"image_url,omitempty"`
	Images                                                []string      `json:"images,omitempty"`
	ImporterAddress                                       string        `json:"importer_address,omitempty"`
	ImporterName                                          string        `json:"importer_name,omitempty"`
	InvalidationErrors                                    []string      `json:"invalidation_errors,omitempty"`
	Inventory                                             int           `json:"inventory,omitempty"`
	ManufacturerInfo                                      string        `json:"manufacturer_info,omitempty"`
	ManufacturerPartNumber                                string        `json:"manufacturer_part_number,omitempty"`
	MarkedForProductLaunch                                string        `json:"marked_for_product_launch,omitempty"`
	Material                                              string        `json:"material,omitempty"`
	MobileLink                                            string        `json:"mobile_link,omitempty"`
	Name                                                  string        `json:"name,omitempty"`
	OrderingIndex                                         int           `json:"ordering_index,omitempty"`
	OriginCountry                                         string        `json:"origin_country,omitempty"`
	ParentProductID                                       string        `json:"parent_product_id,omitempty"`
	Pattern                                               string        `json:"pattern,omitempty"`
	PostConversionSignalBasedEnforcementAppealEligibility bool          `json:"post_conversion_signal_based_enforcement_appeal_eligibility,omitempty"`
	Price                                                 string        `json:"price,omitempty"`
	ProductFeed                                           string        `json:"product_feed,omitempty"`
	ProductGroup                                          *ProductGroup `json:"product_group,omitempty"`
	ProductLocalInfo                                      string        `json:"product_local_info,omitempty"`
	ProductType                                           string        `json:"product_type,omitempty"`
	QuantityToSellOnFacebook                              int           `json:"quantity_to_sell_on_facebook,omitempty"`
	RetailerId                                            string        `json:"retailer_id,omitempty"`
	RetailerProductGroupID                                string        `json:"retailer_product_group_id,omitempty"`
	ReviewRejectionReasons                                []string      `json:"review_rejection_reasons,omitempty"`
	ReviewStatus                                          string        `json:"review_status,omitempty"`
	SalePrice                                             string        `json:"sale_price,omitempty"`
	SalePriceEndDate                                      string        `json:"sale_price_end_date,omitempty"`
	SalePriceStartDate                                    string        `json:"sale_price_start_date,omitempty"`
	ShippingWeightUnit                                    string        `json:"shipping_weight_unit,omitempty"`
	ShippingWeightValue                                   float64       `json:"shipping_weight_value,omitempty"`
	ShortDescription                                      string        `json:"short_description,omitempty"`
	Size                                                  string        `json:"size,omitempty"`
	StartDate                                             string        `json:"start_date,omitempty"`
	Tags                                                  []string      `json:"tags,omitempty"`
	Url                                                   string        `json:"url,omitempty"`
	VendorId                                              string        `json:"vendor_id,omitempty"`
	VideoFetchStatus                                      string        `json:"video_fetch_status,omitempty"`
	Visibility                                            string        `json:"visibility,omitempty"`
	WaComplianceCategory                                  string        `json:"wa_compliance_category,omitempty"`
}

// ProductSetMetadata represents metadata for a product set.
type ProductSetMetadata struct {
	// Define minimal metadata fields as needed.
	// For example:
	UpdateTime string `json:"update_time,omitempty"`
}

// ProductSet represents a product set within a catalog.
type ProductSet struct {
	Id              string             `json:"id"`
	AutoCreationUrl string             `json:"auto_creation_url,omitempty"`
	Filter          string             `json:"filter,omitempty"`
	LatestMetadata  ProductSetMetadata `json:"latest_metadata,omitempty"`
	LiveMetadata    ProductSetMetadata `json:"live_metadata,omitempty"`
	Name            string             `json:"name"`
	// Omit full ProductCatalog to avoid cyclic properties; use catalog Id instead.
	ProductCount uint32 `json:"product_count"`
	RetailerId   string `json:"retailer_id,omitempty"`
}

// Update Catalog to use []ProductItem for Products.
type Catalog struct {
	// Basic fields
	Id           string `json:"id"`
	Name         string `json:"name"`
	ProductCount int    `json:"product_count"`
	Vertical     string `json:"vertical,omitempty"`

	// Edge relationships (represented as slices for brevity)
	Agencies                      *[]interface{} `json:"agencies,omitempty"`
	AssignedUsers                 *[]interface{} `json:"assigned_users,omitempty"`
	AutomotiveModels              *[]interface{} `json:"automotive_models,omitempty"`
	Categories                    *[]interface{} `json:"categories,omitempty"`
	CheckBatchRequestStatus       *[]interface{} `json:"check_batch_request_status,omitempty"`
	CollaborativeAdsShareSettings *[]interface{} `json:"collaborative_ads_share_settings,omitempty"`
	DataSources                   *interface{}   `json:"data_sources,omitempty"`
	Destinations                  *[]interface{} `json:"destinations,omitempty"`
	Diagnostics                   *[]interface{} `json:"diagnostics,omitempty"`
	EventStats                    *[]interface{} `json:"event_stats,omitempty"`
	ExternalEventSources          *[]interface{} `json:"external_event_sources,omitempty"`
	Flights                       *[]interface{} `json:"flights,omitempty"`
	HomeListings                  *[]interface{} `json:"home_listings,omitempty"`
	HotelRoomsBatch               *[]interface{} `json:"hotel_rooms_batch,omitempty"`
	Hotels                        *[]interface{} `json:"hotels,omitempty"`
	PricingVariablesBatch         *[]interface{} `json:"pricing_variables_batch,omitempty"`
	ProductGroups                 *interface{}   `json:"product_groups,omitempty"`
	ProductSets                   *interface{}   `json:"product_sets,omitempty"`
	ProductSetsBatch              *[]interface{} `json:"product_sets_batch,omitempty"`
	Products                      struct {
		Data []ProductItem `json:"data"`
	} `json:"products,omitempty"`
	VehicleOffers *[]interface{} `json:"vehicle_offers,omitempty"`
	Vehicles      *[]interface{} `json:"vehicles,omitempty"`
}

type CatalogFetchResponseEdge struct {
	Data   []Catalog                                  `json:"data"`
	Paging internal.WhatsAppBusinessApiPaginationMeta `json:"paging"`
}

func (cm *CatalogManager) GetAllCatalogs() (*CatalogFetchResponseEdge, error) {
	apiPath := strings.Join([]string{cm.businessAccountId, "product_catalogs"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodGet)

	fields := []string{
		"id",
		"name",
		"vertical",
		"product_count",
		"data_sources",
		"product_groups",
		"product_sets",
		"products",
	}

	for _, field := range fields {
		apiRequest.AddField(request_client.ApiRequestQueryParamField{
			Name:    field,
			Filters: map[string]string{},
		})
	}

	response, err := apiRequest.Execute()

	if err != nil {
		fmt.Println("GetAllCatalogs error", err)
		return nil, err
	}
	fmt.Println("GetAllCatalogs response", response)
	var result CatalogFetchResponseEdge
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCatalogProducts retrieves the list of products for a given catalog.
func (cm *CatalogManager) GetCatalogProducts(catalogId string) ([]ProductItem, error) {
	apiPath := strings.Join([]string{catalogId, "products"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodGet)

	fields := []string{
		"id",
		"additional_image_cdn_urls",
		"additional_image_urls",
		"additional_variant_attributes",
		"age_group",
		"applinks",
		"availability",
		"brand",
		"capability_to_review_status",
		"category",
		"category_specific_fields",
		"color",
		"commerce_insights",
		"condition",
		"currency",
		"custom_data",
		"custom_label_0",
		"custom_label_1",
		"custom_label_2",
		"custom_label_3",
		"custom_label_4",
		"custom_number_0",
		"custom_number_1",
		"custom_number_2",
		"custom_number_3",
		"custom_number_4",
		"description",
		"errors",
		"expiration_date",
		"fb_product_category",
		"gender",
		"gtin",
		"image_cdn_urls",
		"image_fetch_status",
		"image_url",
		"images",
		"importer_address",
		"importer_name",
		"invalidation_errors",
		"inventory",
		"manufacturer_info",
		"manufacturer_part_number",
		"marked_for_product_launch",
		"material",
		"mobile_link",
		"name",
		"ordering_index",
		"origin_country",
		"parent_product_id",
		"pattern",
		"post_conversion_signal_based_enforcement_appeal_eligibility",
		"price",
		"product_feed",
		"product_group",
		"product_local_info",
		"product_type",
		"quantity_to_sell_on_facebook",
		"retailer_id",
		"retailer_product_group_id",
		"review_rejection_reasons",
		"review_status",
		"sale_price",
		"sale_price_end_date",
		"sale_price_start_date",
		"shipping_weight_unit",
		"shipping_weight_value",
		"short_description",
		"size",
		"start_date",
		"tags",
		"url",
		"vendor_id",
		"video_fetch_status",
		"visibility",
		"wa_compliance_category",
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

	fmt.Println("GetCatalogProducts response", response)

	// Temporary response structure.
	var res struct {
		Data    []ProductItem   `json:"data"`
		Paging  json.RawMessage `json:"paging"` // using RawMessage to ignore paging details
		Summary struct {
			TotalCount int `json:"total_count"`
		} `json:"summary"`
	}
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

type CreateProductCatalogOptions struct {
	Success string `json:"success,omitempty"`
}

func (cm *CatalogManager) CreateNewProductCatalog() (CreateProductCatalogOptions, error) {
	apiRequest := cm.requester.NewApiRequest(strings.Join([]string{cm.businessAccountId, "product_catalogs"}, "/"), http.MethodPost)
	response, err := apiRequest.Execute()
	var responseToReturn CreateProductCatalogOptions
	json.Unmarshal([]byte(response), &responseToReturn)
	return responseToReturn, err
}
