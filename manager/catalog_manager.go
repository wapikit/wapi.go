package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"

	"github.com/gTahidi/wapi.go/internal"
	"github.com/gTahidi/wapi.go/internal/request_client"
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

type ProductFeed struct {
	Id       string `json:"id"`
	FileName string `json:"file_name"`
	Name     string `json:"name"`
}

// FeedUpload represents a single upload attempt and its status/diagnostics.
type FeedUpload struct {
	Id           string   `json:"id,omitempty"`
	Status       string   `json:"status,omitempty"`
	Errors       []string `json:"errors,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
	CreatedTime  string   `json:"created_time,omitempty"`
	LastUpdated  string   `json:"last_updated_time,omitempty"`
}

// FeedUploadResponse wraps basic responses for CSV upload operations.
type FeedUploadResponse struct {
	Id     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

// FeedUploadSession represents an upload session listing entry (start/end time).
type FeedUploadSession struct {
	Id        string `json:"id"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// FeedUploadErrorSample represents a row sample in an error entry.
type FeedUploadErrorSample struct {
	RowNumber  int    `json:"row_number"`
	RetailerId string `json:"retailer_id"`
	Id         string `json:"id"`
}

// FeedUploadError represents an individual ingestion error/warning.
type FeedUploadError struct {
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Severity    string `json:"severity"` // fatal | warning
	Samples     struct {
		Data []FeedUploadErrorSample `json:"data"`
	} `json:"samples"`
}

// FeedUploadErrorReport represents the generated error report metadata.
type FeedUploadErrorReport struct {
	ReportStatus string `json:"report_status,omitempty"`
	FileHandle   string `json:"file_handle,omitempty"`
}

type FeedUploadErrorReportResponse struct {
	ErrorReport FeedUploadErrorReport `json:"error_report,omitempty"`
	Id          string                `json:"id,omitempty"`
}

type ProductError struct {
	ErrorType     string `json:"error_type"`
	ErrorPriority string `json:"error_priority"`
	Title         string `json:"title"`
	Description   string `json:"description"`
}

// ProductItem represents a product in a catalog.

type ProductItem struct {
	Id                                                    string         `json:"id"`
	AdditionalImageCdnUrls                                []ImageCdnUrl  `json:"additional_image_cdn_urls,omitempty"`
	AdditionalImageUrls                                   []string       `json:"additional_image_urls,omitempty"`
	AdditionalVariantAttributes                           []KeyValue     `json:"additional_variant_attributes,omitempty"`
	AgeGroup                                              string         `json:"age_group,omitempty"`
	Applinks                                              *Applinks      `json:"applinks,omitempty"`
	Availability                                          string         `json:"availability,omitempty"`
	Brand                                                 string         `json:"brand,omitempty"`
	CapabilityToReviewStatus                              []KeyValue     `json:"capability_to_review_status,omitempty"`
	Category                                              string         `json:"category,omitempty"`
	CategorySpecificFields                                interface{}    `json:"category_specific_fields,omitempty"`
	Color                                                 string         `json:"color,omitempty"`
	CommerceInsights                                      string         `json:"commerce_insights,omitempty"`
	Condition                                             string         `json:"condition,omitempty"`
	Currency                                              string         `json:"currency,omitempty"`
	CustomData                                            []KeyValue     `json:"custom_data,omitempty"`
	CustomLabel0                                          string         `json:"custom_label_0,omitempty"`
	CustomLabel1                                          string         `json:"custom_label_1,omitempty"`
	CustomLabel2                                          string         `json:"custom_label_2,omitempty"`
	CustomLabel3                                          string         `json:"custom_label_3,omitempty"`
	CustomLabel4                                          string         `json:"custom_label_4,omitempty"`
	CustomNumber0                                         string         `json:"custom_number_0,omitempty"`
	CustomNumber1                                         string         `json:"custom_number_1,omitempty"`
	CustomNumber2                                         string         `json:"custom_number_2,omitempty"`
	CustomNumber3                                         string         `json:"custom_number_3,omitempty"`
	CustomNumber4                                         string         `json:"custom_number_4,omitempty"`
	Description                                           string         `json:"description,omitempty"`
	Errors                                                []ProductError `json:"errors,omitempty"`
	ExpirationDate                                        string         `json:"expiration_date,omitempty"`
	FbProductCategory                                     string         `json:"fb_product_category,omitempty"`
	Gender                                                string         `json:"gender,omitempty"`
	Gtin                                                  string         `json:"gtin,omitempty"`
	ImageCdnUrls                                          []ImageCdnUrl  `json:"image_cdn_urls,omitempty"`
	ImageFetchStatus                                      string         `json:"image_fetch_status,omitempty"`
	ImageUrl                                              string         `json:"image_url,omitempty"`
	Images                                                []string       `json:"images,omitempty"`
	ImporterAddress                                       string         `json:"importer_address,omitempty"`
	ImporterName                                          string         `json:"importer_name,omitempty"`
	InvalidationErrors                                    []string       `json:"invalidation_errors,omitempty"`
	Inventory                                             int            `json:"inventory,omitempty"`
	ManufacturerInfo                                      string         `json:"manufacturer_info,omitempty"`
	ManufacturerPartNumber                                string         `json:"manufacturer_part_number,omitempty"`
	MarkedForProductLaunch                                string         `json:"marked_for_product_launch,omitempty"`
	Material                                              string         `json:"material,omitempty"`
	MobileLink                                            string         `json:"mobile_link,omitempty"`
	Name                                                  string         `json:"name,omitempty"`
	OrderingIndex                                         int            `json:"ordering_index,omitempty"`
	OriginCountry                                         string         `json:"origin_country,omitempty"`
	ParentProductID                                       string         `json:"parent_product_id,omitempty"`
	Pattern                                               string         `json:"pattern,omitempty"`
	PostConversionSignalBasedEnforcementAppealEligibility bool           `json:"post_conversion_signal_based_enforcement_appeal_eligibility,omitempty"`
	Price                                                 string         `json:"price,omitempty"`
	ProductFeed                                           *ProductFeed   `json:"product_feed,omitempty"`
	ProductGroup                                          *ProductGroup  `json:"product_group,omitempty"`
	ProductLocalInfo                                      string         `json:"product_local_info,omitempty"`
	ProductType                                           string         `json:"product_type,omitempty"`
	QuantityToSellOnFacebook                              int            `json:"quantity_to_sell_on_facebook,omitempty"`
	RetailerId                                            string         `json:"retailer_id,omitempty"`
	RetailerProductGroupID                                string         `json:"retailer_product_group_id,omitempty"`
	ReviewRejectionReasons                                []string       `json:"review_rejection_reasons,omitempty"`
	ReviewStatus                                          string         `json:"review_status,omitempty"`
	SalePrice                                             string         `json:"sale_price,omitempty"`
	SalePriceEndDate                                      string         `json:"sale_price_end_date,omitempty"`
	SalePriceStartDate                                    string         `json:"sale_price_start_date,omitempty"`
	ShippingWeightUnit                                    string         `json:"shipping_weight_unit,omitempty"`
	ShippingWeightValue                                   float64        `json:"shipping_weight_value,omitempty"`
	ShortDescription                                      string         `json:"short_description,omitempty"`
	Size                                                  string         `json:"size,omitempty"`
	StartDate                                             string         `json:"start_date,omitempty"`
	Tags                                                  []string       `json:"tags,omitempty"`
	Url                                                   string         `json:"url,omitempty"`
	VendorId                                              string         `json:"vendor_id,omitempty"`
	VideoFetchStatus                                      string         `json:"video_fetch_status,omitempty"`
	Visibility                                            string         `json:"visibility,omitempty"`
	WaComplianceCategory                                  string         `json:"wa_compliance_category,omitempty"`
}

// ProductSetMetadata represents metadata for a product set.
type ProductSetMetadata struct {
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
		"price",
		"additional_image_cdn_urls",
		"additional_image_urls",
		"condition",
		"additional_variant_attributes",
		"age_group",
		"availability",
		"brand",
		"category",
		"category_specific_fields",
		"color",
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
		"inventory",
		"material",
		"name",
		"ordering_index",
		"origin_country",
		"parent_product_id",
		"pattern",
		"product_local_info",
		"product_type",
		"quantity_to_sell_on_facebook",
		"retailer_id",
		"retailer_product_group_id",
		"sale_price",
		"sale_price_end_date",
		"sale_price_start_date",
		"short_description",
		"size",
		"start_date",
		"tags",
		"url",
		"vendor_id",
		"visibility",
		"wa_compliance_category",
	}

	for _, field := range fields {
		apiRequest.AddField(request_client.ApiRequestQueryParamField{
			Name:    field,
			Filters: map[string]string{},
		})
	}

	// ! TODO: proper pagination must be implemented here
	apiRequest.AddQueryParam("limit", "1000")

	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		// Return immediately on execution error; do not attempt to decode
		return CreateProductCatalogOptions{}, fmt.Errorf("create catalog request failed: %w", err)
	}
	var responseToReturn CreateProductCatalogOptions
	if err := json.Unmarshal([]byte(response), &responseToReturn); err != nil {
		// Return zero-value options with decoding context for callers
		return CreateProductCatalogOptions{}, fmt.Errorf("decode create catalog response failed: %w", err)
	}
	return responseToReturn, nil
}

// ListProductFeeds lists product feeds for a given catalog.
func (cm *CatalogManager) ListProductFeeds(catalogId string) ([]ProductFeed, error) {
	apiPath := strings.Join([]string{catalogId, "product_feeds"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []ProductFeed `json:"data"`
	}
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// CreateProductFeed creates a product feed for CSV ingestion. We are using metas base accepted format, for more info check the meta docs
// meta whatsapp catalogs: https://developers.facebook.com/docs/commerce-platform/catalog/fields
// name: Human-readable feed name
// fileFormat: e.g., "CSV"
// fileName: default file name reference (optional)
func (cm *CatalogManager) CreateProductFeed(catalogId, name, fileFormat, fileName string) (*ProductFeed, error) {
	apiPath := strings.Join([]string{catalogId, "product_feeds"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodPost)
	body := map[string]string{
		"name":        name,
		"file_format": fileFormat,
	}
	if fileName != "" {
		body["file_name"] = fileName
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create feed body: %w", err)
	}
	apiRequest.SetBody(string(payload))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res ProductFeed
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// UploadFeedCSV uploads a CSV file to a product feed using multipart/form-data.
func (cm *CatalogManager) UploadFeedCSV(feedId string, file io.Reader, filename, mimeType string, updateOnly bool) (*FeedUploadResponse, error) {
	// Prepare multipart body with update_only and a single file part
	bodyBuf := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyBuf)
	// update_only as string field
	if err := writer.WriteField("update_only", func() string {
		if updateOnly {
			return "true"
		}
		return "false"
	}()); err != nil {
		return nil, fmt.Errorf("failed to write update_only: %w", err)
	}
	// file part
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filename)))
	partHeader.Set("Content-Type", mimeType)
	filePart, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart part: %w", err)
	}
	if _, err := io.Copy(filePart, file); err != nil {
		return nil, fmt.Errorf("failed to copy csv into part: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	apiPath := strings.Join([]string{feedId, "uploads"}, "/")
	contentType := writer.FormDataContentType()
	responseBody, err := cm.requester.RequestMultipart(http.MethodPost, apiPath, bodyBuf, contentType)
	if err != nil {
		return nil, fmt.Errorf("error uploading CSV: %w", err)
	}

	var res FeedUploadResponse
	if err := json.Unmarshal([]byte(responseBody), &res); err != nil {
		return nil, fmt.Errorf("failed to parse upload response: %w", err)
	}
	return &res, nil
}

// UploadFeedCSVFromURL triggers a feed ingestion from a hosted CSV URL.
func (cm *CatalogManager) UploadFeedCSVFromURL(feedId, csvURL string, updateOnly bool) (*FeedUploadResponse, error) {
	apiPath := strings.Join([]string{feedId, "uploads"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodPost)
	body := map[string]interface{}{
		// Meta docs show using 'url' for hosted feed uploads
		"url":         csvURL,
		"update_only": updateOnly,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hosted feed body: %w", err)
	}
	apiRequest.SetBody(string(payload))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res FeedUploadResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ListFeedUploads lists uploads for a product feed.
func (cm *CatalogManager) ListFeedUploads(feedId string) ([]FeedUploadSession, error) {
	apiPath := strings.Join([]string{feedId, "uploads"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []FeedUploadSession `json:"data"`
	}
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// GetFeedUploadStatus fetches a single upload’s status/diagnostics.
func (cm *CatalogManager) GetFeedUploadStatus(uploadId string) (*FeedUploadErrorReportResponse, error) {
	apiRequest := cm.requester.NewApiRequest(uploadId, http.MethodGet)
	// include error_report field for convenience
	apiRequest.AddField(request_client.ApiRequestQueryParamField{Name: "error_report", Filters: map[string]string{}})
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res FeedUploadErrorReportResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetFeedUploadErrors fetches a sampling of errors/warnings for an upload session.
func (cm *CatalogManager) GetFeedUploadErrors(uploadId string) ([]FeedUploadError, error) {
	apiPath := strings.Join([]string{uploadId, "errors"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodGet)
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res struct {
		Data []FeedUploadError `json:"data"`
	}
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return res.Data, nil
}

// RequestFeedUploadErrorReport triggers generation of a full error report.
func (cm *CatalogManager) RequestFeedUploadErrorReport(uploadId string) (bool, error) {
	apiPath := strings.Join([]string{uploadId, "error_report"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodPost)
	response, err := apiRequest.Execute()
	if err != nil {
		return false, err
	}
	var res struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return false, err
	}
	return res.Success, nil
}

// GetFeedUploadErrorReport fetches the error_report field of an upload session.
func (cm *CatalogManager) GetFeedUploadErrorReport(uploadId string) (*FeedUploadErrorReportResponse, error) {
	apiRequest := cm.requester.NewApiRequest(uploadId, http.MethodGet)
	apiRequest.AddField(request_client.ApiRequestQueryParamField{Name: "error_report", Filters: map[string]string{}})
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res FeedUploadErrorReportResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ProductFeedSchedule represents schedule config for Google Sheets or hosted feeds.
type ProductFeedSchedule struct {
	Url      string `json:"url"`
	Interval string `json:"interval"` // e.g., HOURLY, DAILY; Meta specific values
	Hour     *int   `json:"hour,omitempty"`
}

// CreateScheduledProductFeed creates a scheduled feed that fetches from a URL (Google Sheets supported via shareable link).
// When updateOnly is true, the feed behaves in update-only mode.
// ingestionSourceType: "PRIMARY_FEED" or "SUPPLEMENTARY_FEED" (optional)
// primaryFeedIds required for Supplementary feeds.
func (cm *CatalogManager) CreateScheduledProductFeed(
	catalogId string,
	name string,
	schedule ProductFeedSchedule,
	updateOnly bool,
	ingestionSourceType string,
	primaryFeedIds []string,
) (*ProductFeed, error) {
	apiPath := strings.Join([]string{catalogId, "product_feeds"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodPost)
	body := map[string]interface{}{
		"name":        name,
		"schedule":    schedule,
		"update_only": updateOnly,
	}
	if ingestionSourceType != "" {
		body["ingestion_source_type"] = ingestionSourceType
	}
	if len(primaryFeedIds) > 0 {
		body["primary_feed_ids"] = primaryFeedIds
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scheduled feed body: %w", err)
	}
	apiRequest.SetBody(string(payload))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res ProductFeed
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// UpsertProductItem updates or creates a product item using Meta’s format.
// fields should include at least retailer_id, name, price, currency, image_url, availability, etc.
func (cm *CatalogManager) UpsertProductItem(catalogId string, fields map[string]interface{}) (*ProductItem, error) {
	apiPath := strings.Join([]string{catalogId, "products"}, "/")
	apiRequest := cm.requester.NewApiRequest(apiPath, http.MethodPost)
	payload, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall product fields: %w", err)

	}
	apiRequest.SetBody(string(payload))
	response, err := apiRequest.Execute()
	if err != nil {
		return nil, err
	}
	var res ProductItem
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// BatchUpsertProductItems performs multiple upserts sequentially.
// Returns the successfully upserted items and a map of index->error for failures.
func (cm *CatalogManager) BatchUpsertProductItems(catalogId string, items []map[string]interface{}) ([]ProductItem, map[int]error) {
	var results []ProductItem
	errs := make(map[int]error)
	for i, fields := range items {
		item, err := cm.UpsertProductItem(catalogId, fields)
		if err != nil {
			errs[i] = err
			continue
		}
		results = append(results, *item)
	}
	return results, errs
}

// UpdateProductImages updates image_url and additional_image_urls for a retailer_id.
func (cm *CatalogManager) UpdateProductImages(catalogId, retailerId, imageURL string, additionalImageURLs []string) (*ProductItem, error) {
	fields := map[string]interface{}{
		"retailer_id":           retailerId,
		"image_url":             imageURL,
		"additional_image_urls": additionalImageURLs,
	}
	return cm.UpsertProductItem(catalogId, fields)
}
