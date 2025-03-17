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

	"github.com/wapikit/wapi.go/internal/request_client"
)

// MediaManager is responsible for managing media related operations.
type MediaManager struct {
	requester request_client.RequestClient
}

// NewMediaManager creates a new instance of MediaManager.
func NewMediaManager(requester request_client.RequestClient) *MediaManager {
	return &MediaManager{
		requester: requester,
	}
}

type MediaMetadata struct {
	MessagingProduct string `json:"messaging_product"`
	Url              string `json:"url"`
	MimeType         string `json:"mime_type"`
	Sha256           string `json:"sha256"`
	FileSize         int    `json:"file_size"`
	ID               string `json:"id"`
}

func (mm *MediaManager) GetMediaUrlById(id string) (string, error) {
	// Build GET request to: e.g. "<MEDIA_ID>" (the request client automatically prefixes the base URL and version)
	apiRequest := mm.requester.NewApiRequest(id, http.MethodGet)

	// Execute the request and get the raw JSON response
	rawResponse, err := apiRequest.Execute()
	if err != nil {
		return "", err
	}

	// Parse into a struct
	var res MediaMetadata
	if err := json.Unmarshal([]byte(rawResponse), &res); err != nil {
		return "", fmt.Errorf("failed to parse media metadata: %w", err)
	}

	if res.Url == "" {
		return "", fmt.Errorf("no media url found in response: %s", rawResponse)
	}

	return res.Url, nil
}

type DeleteSuccessResponse struct {
	Success bool `json:"success"`
}

func (mm *MediaManager) DeleteMedia(id string) (string, error) {
	// The path becomes "media/<MEDIA_ID>"
	apiRequest := mm.requester.NewApiRequest(strings.Join([]string{"media", id}, "/"), http.MethodDelete)

	rawResponse, err := apiRequest.Execute()
	if err != nil {
		return "", err
	}

	// Parse the JSON
	var res DeleteSuccessResponse
	if err := json.Unmarshal([]byte(rawResponse), &res); err != nil {
		return "", fmt.Errorf("failed to parse delete response: %w", err)
	}

	if !res.Success {
		return "", fmt.Errorf("media deletion failed or returned success=false: %s", rawResponse)
	}

	return "media deleted successfully", nil
}

// UploadMedia uploads a media file to WhatsApp's Cloud API.
func (mm *MediaManager) UploadMedia(phoneNumberId string, file io.Reader, filename, mimeType string) (string, error) {
	// 1. Build the multipart form in memory
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", fmt.Errorf("failed to write field: %w", err)
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filename)))
	partHeader.Set("Content-Type", mimeType)

	filePart, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", fmt.Errorf("failed to create multipart part: %w", err)
	}
	if _, err := io.Copy(filePart, file); err != nil {
		return "", fmt.Errorf("failed to copy file into part: %w", err)
	}

	// Close the writer to finalize the multipart data
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	apiPath := strings.Join([]string{phoneNumberId, "media"}, "/")

	contentType := writer.FormDataContentType()

	responseBody, err := mm.requester.RequestMultipart(http.MethodPost, apiPath, body, contentType)
	if err != nil {
		return "", fmt.Errorf("error uploading media: %w", err)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(responseBody), &result); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}
	if result.ID == "" {
		// Possibly an error or partial success
		return "", fmt.Errorf("no media id in response: %s", responseBody)
	}

	return result.ID, nil
}
