package manager

import (
	"net/http"
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

// GetMediaUrlById retrieves the media URL by its ID.
func (mm *MediaManager) GetMediaUrlById(id string) (string, error) {
	apiRequest := mm.requester.NewApiRequest(strings.Join([]string{"media", id}, "/"), http.MethodGet)
	response, err := apiRequest.Execute()
	return response, err
}

// GetMediaIdByUrl retrieves the media ID by its URL.
func (mm *MediaManager) DeleteMedia(id string) (string, error) {
	apiRequest := mm.requester.NewApiRequest(strings.Join([]string{"media", id}, "/"), http.MethodDelete)
	response, err := apiRequest.Execute()
	return response, err
}
