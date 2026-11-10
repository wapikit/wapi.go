package request_client

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	API_VERSION      = "v20.0"
	BASE_URL         = "graph.facebook.com"
	REQUEST_PROTOCOL = "https"
)

// DefaultResponseHeaderTimeout bounds how long the default HTTP client waits
// for Meta to start responding once a request has been fully written. It does
// not limit how long writing the request body takes, so large media uploads
// are not cut off.
const DefaultResponseHeaderTimeout = 60 * time.Second

// defaultHTTPClient is shared by every RequestClient that is not given its own
// *http.Client, so connections are pooled across clients. It has no overall
// Timeout (that would also cap slow uploads); instead each phase that can hang
// — dial, TLS handshake, waiting for response headers — is bounded.
var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: DefaultResponseHeaderTimeout,
	},
}

type WhatsappApiType string

const (
	WhatsappApiTypeMessaging WhatsappApiType = "messaging"
	WhatsappApiTypeBusiness  WhatsappApiType = "business"
)

// RequestClient represents a client for making requests to a cloud API.
//
// It is safe to copy and to use from multiple goroutines: all fields are
// read-only after construction and the *http.Client is shared.
type RequestClient struct {
	apiVersion     string
	baseUrl        string // host only, e.g. "graph.facebook.com" (kept for BaseUrl()).
	endpoint       string // scheme + host, e.g. "https://graph.facebook.com".
	apiAccessToken string
	httpClient     *http.Client
}

func (client *RequestClient) BaseUrl() string {
	return client.baseUrl
}

func (client *RequestClient) ApiVersion() string {
	return client.apiVersion
}

// Options customises a RequestClient. Every field is optional; the zero value
// gives the historical defaults.
type Options struct {
	// HTTPClient is used for every request. Defaults to a shared client with
	// dial/TLS/response-header timeouts (see DefaultResponseHeaderTimeout).
	HTTPClient *http.Client
	// ApiVersion is the Graph API version, e.g. "v21.0". Defaults to API_VERSION.
	ApiVersion string
	// BaseURL is the scheme and host of the Graph API, e.g.
	// "https://graph.facebook.com". Mainly useful to point at an httptest
	// server in tests. Defaults to REQUEST_PROTOCOL://BASE_URL.
	BaseURL string
}

// NewRequestClient creates a new instance of RequestClient.
func NewRequestClient(apiAccessToken string) *RequestClient {
	return NewRequestClientWithOptions(apiAccessToken, Options{})
}

// NewRequestClientWithOptions creates a RequestClient with custom options.
func NewRequestClientWithOptions(apiAccessToken string, opts Options) *RequestClient {
	client := &RequestClient{
		apiVersion:     API_VERSION,
		baseUrl:        BASE_URL,
		endpoint:       REQUEST_PROTOCOL + "://" + BASE_URL,
		apiAccessToken: apiAccessToken,
		httpClient:     defaultHTTPClient,
	}
	if opts.HTTPClient != nil {
		client.httpClient = opts.HTTPClient
	}
	if opts.ApiVersion != "" {
		client.apiVersion = opts.ApiVersion
	}
	if opts.BaseURL != "" {
		client.endpoint = strings.TrimRight(opts.BaseURL, "/")
		if u, err := url.Parse(client.endpoint); err == nil && u.Host != "" {
			client.baseUrl = u.Host
		}
	}
	return client
}

// RequestCloudApiParams represents the parameters for making a request to the cloud API.
type RequestCloudApiParams struct {
	Body       string
	Path       string
	Method     string
	QueryParam map[string]string
}

// buildURL joins the endpoint, API version, path and (URL-encoded) query
// parameters. Query parameters are emitted in sorted key order.
func (client *RequestClient) buildURL(path string, query map[string]string) string {
	requestURL := client.endpoint + "/" + client.apiVersion + "/" + path
	if len(query) > 0 {
		values := url.Values{}
		for key, value := range query {
			values.Set(key, value)
		}
		requestURL += "?" + values.Encode()
	}
	return requestURL
}

func (client *RequestClient) do(ctx context.Context, method, requestURL string, body io.Reader, contentType string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	httpRequest, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", contentType)
	httpRequest.Header.Set("Authorization", "Bearer "+client.apiAccessToken)

	httpClient := client.httpClient
	if httpClient == nil { // zero-value RequestClient
		httpClient = defaultHTTPClient
	}
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	// Surface non-2xx Graph API responses as a typed error (the body is still
	// returned so callers that already parse Meta's error envelope keep working).
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return string(respBody), newGraphAPIError(response.StatusCode, string(respBody))
	}
	return string(respBody), nil
}

func (requestClientInstance *RequestClient) request(ctx context.Context, params RequestCloudApiParams) (string, error) {
	requestURL := requestClientInstance.buildURL(params.Path, params.QueryParam)
	return requestClientInstance.do(ctx, params.Method, requestURL, strings.NewReader(params.Body), "application/json")
}

func (client *RequestClient) NewApiRequest(path, method string) *ApiRequest {
	return &ApiRequest{
		Path:        path,
		Fields:      []ApiRequestQueryParamField{},
		Requester:   client,
		Method:      method,
		QueryParams: map[string]string{},
	}
}

type ApiRequestQueryParamField struct {
	Name    string
	Filters map[string]string
}

func (field *ApiRequestQueryParamField) AddFilter(key, value string) {
	if field.Filters == nil {
		field.Filters = map[string]string{}
	}
	field.Filters[key] = value
}

type ApiRequest struct {
	Path        string
	Method      string
	Body        string
	Fields      []ApiRequestQueryParamField
	QueryParams map[string]string
	Requester   *RequestClient
}

// AddField adds a field to the request's "fields" query parameter and returns
// a pointer to the stored field, so filters added through it are sent.
//
// Fields are comma separated, and filters are rendered as chained calls, e.g.
// fields=field1.filter1(value1).filter2(value2),field2
func (request *ApiRequest) AddField(field ApiRequestQueryParamField) *ApiRequestQueryParamField {
	if field.Filters == nil {
		field.Filters = map[string]string{}
	}
	request.Fields = append(request.Fields, field)
	return &request.Fields[len(request.Fields)-1]
}

// AddQueryParam adds a query parameter to the request.
func (request *ApiRequest) AddQueryParam(key, value string) {
	request.QueryParams[key] = value
}

// SetMethod sets the method for the request.
func (request *ApiRequest) SetMethod(method string) {
	request.Method = method
}

// SetBody sets the body for the request.
func (request *ApiRequest) SetBody(body string) {
	request.Body = body
}

// fieldsParam renders Fields as the Graph API "fields" query value. Filters
// within a field are emitted in sorted key order so the URL is deterministic.
func (request *ApiRequest) fieldsParam() string {
	parts := make([]string, 0, len(request.Fields))
	for _, field := range request.Fields {
		keys := make([]string, 0, len(field.Filters))
		for key := range field.Filters {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var b strings.Builder
		b.WriteString(field.Name)
		for _, key := range keys {
			b.WriteString("." + key + "(" + field.Filters[key] + ")")
		}
		parts = append(parts, b.String())
	}
	return strings.Join(parts, ",")
}

// Execute executes the request and returns the response body.
//
// A non-2xx status yields a *GraphAPIError while still returning the body, so
// callers can both detect the failure and inspect Meta's error envelope.
func (request *ApiRequest) Execute() (string, error) {
	return request.ExecuteContext(context.Background())
}

// ExecuteContext is like Execute but carries ctx for cancellation and deadlines.
func (request *ApiRequest) ExecuteContext(ctx context.Context) (string, error) {
	queryParam := map[string]string{}
	if len(request.Fields) > 0 {
		queryParam["fields"] = request.fieldsParam()
	}
	for key, value := range request.QueryParams {
		queryParam[key] = value
	}

	return request.Requester.request(ctx, RequestCloudApiParams{
		Path:       request.Path,
		Body:       request.Body,
		Method:     request.Method,
		QueryParam: queryParam,
	})
}

// RequestMultipart allows sending an arbitrary body with a custom Content-Type.
// This is needed for file uploads (multipart/form-data).
func (rc *RequestClient) RequestMultipart(
	method string,
	path string,
	body io.Reader,
	contentType string,
) (string, error) {
	return rc.RequestMultipartContext(context.Background(), method, path, body, contentType)
}

// RequestMultipartContext is like RequestMultipart but carries ctx.
func (rc *RequestClient) RequestMultipartContext(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	contentType string,
) (string, error) {
	response, err := rc.do(ctx, method, rc.buildURL(path, nil), body, contentType)
	var graphErr *GraphAPIError
	if err != nil && !asGraphAPIError(err, &graphErr) {
		return response, fmt.Errorf("failed to execute request: %w", err)
	}
	return response, err
}
