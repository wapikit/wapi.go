package request_client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	API_VERSION      = "v20.0"
	BASE_URL         = "graph.facebook.com"
	REQUEST_PROTOCOL = "https"
)

type WhatsappApiType string

const (
	WhatsappApiTypeMessaging WhatsappApiType = "messaging"
	WhatsappApiTypeBusiness  WhatsappApiType = "business"
)

// RequestClient represents a client for making requests to a cloud API.
type RequestClient struct {
	apiVersion     string
	baseUrl        string
	apiAccessToken string
}

func (client *RequestClient) BaseUrl() string {
	return client.baseUrl
}

func (client *RequestClient) ApiVersion() string {
	return client.apiVersion
}

// NewRequestClient creates a new instance of RequestClient.
func NewRequestClient(apiAccessToken string) *RequestClient {
	return &RequestClient{
		apiVersion:     API_VERSION,
		baseUrl:        BASE_URL,
		apiAccessToken: apiAccessToken,
	}
}

// RequestCloudApiParams represents the parameters for making a request to the cloud API.
type RequestCloudApiParams struct {
	Body       string
	Path       string
	Method     string
	QueryParam map[string]string
}

func (requestClientInstance *RequestClient) request(params RequestCloudApiParams) (string, error) {
	queryParamString := ""
	if len(params.QueryParam) > 0 {
		queryParamString = "?"
		for key, value := range params.QueryParam {
			// if first query param, don't add "&"
			if queryParamString != "?" {
				queryParamString += "&"
				queryParamString += strings.Join([]string{queryParamString, key, "=", value}, "")
			} else {
				queryParamString += strings.Join([]string{key, "=", value}, "")
			}
		}
	}

	requestPath := strings.Join(
		[]string{REQUEST_PROTOCOL, "://", requestClientInstance.baseUrl, "/", requestClientInstance.apiVersion, "/", params.Path, queryParamString}, "")

	httpRequest, err := http.NewRequest(params.Method,
		requestPath,
		strings.NewReader(params.Body))
	if err != nil {
		return "", err
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", requestClientInstance.apiAccessToken))
	httpClient := &http.Client{}
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		fmt.Println("Error while requesting cloud api", err)
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
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

func (request *ApiRequest) AddField(field ApiRequestQueryParamField) *ApiRequestQueryParamField {
	// * NOTE:  when we say we need to add a field to the request, it means we need to add a query param to the request with key "fields"
	// * note that if there need to be multiple fields in a single request then the list of fields should be command separated
	// * for example: fields=field1,field2,field3
	// * also note that if there filters in a field then they should be called like a function in the param string, for ex: fields=field1.filter1(value1).filter2(value2),field2.filter1(value1)
	request.Fields = append(request.Fields, field)
	return &field
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

// Execute executes the request and returns the response.
func (request *ApiRequest) Execute() (string, error) {
	// check if there are any fields in the request
	var queryParam = map[string]string{}
	if len(request.Fields) > 0 {
		fieldsString := ""
		for _, field := range request.Fields {
			newFieldString := ""
			if fieldsString != "" {
				newFieldString = ","
			}
			filterString := ""
			for key, value := range field.Filters {
				filterString += strings.Join([]string{".", key, "(", value, ")"}, "")
			}
			newFieldString += strings.Join([]string{field.Name, filterString}, "")
			fieldsString += newFieldString
		}

		queryParam["fields"] = fieldsString
	}

	if len(request.QueryParams) > 0 {
		for key, value := range request.QueryParams {
			queryParam[key] = value
		}
	}

	response, err := request.Requester.request(RequestCloudApiParams{
		Path:       request.Path,
		Body:       request.Body,
		Method:     request.Method,
		QueryParam: queryParam,
	})

	if err != nil {
		fmt.Println("Error while executing business api request", err)
		return "", nil
	}

	return response, err
}

// RequestMultipart allows sending an arbitrary body with a custom Content-Type.
// This is needed for file uploads (multipart/form-data).
func (rc *RequestClient) RequestMultipart(
	method string,
	path string,
	body io.Reader,
	contentType string,
) (string, error) {
	// 1. Build the final URL
	requestPath := strings.Join([]string{
		REQUEST_PROTOCOL, "://",
		rc.baseUrl, "/",
		rc.apiVersion, "/",
		path,
	}, "")

	// 2. Create the HTTP request with the custom body & Content-Type
	httpRequest, err := http.NewRequest(method, requestPath, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rc.apiAccessToken))
	httpRequest.Header.Set("Content-Type", contentType)

	// 3. Send the request using a standard HTTP client
	httpClient := &http.Client{}
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(respBody), nil
}
