package http_test

import (
	"context"
	"github.com/NawafSwe/media-scout-service/pkg/internal/transport"
	kithttp "github.com/NawafSwe/media-scout-service/pkg/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeSearchMediaRequest(t *testing.T) {
	tests := []struct {
		name          string
		queryParams   string
		expectedError string
		expectedTerm  string
		expectedLimit int
	}{
		{
			name:          "valid request",
			queryParams:   "term=test&limit=10",
			expectedError: "",
			expectedTerm:  "test",
			expectedLimit: 10,
		},
		{
			name:          "missing term",
			queryParams:   "limit=10",
			expectedError: "term shouldn't be empty",
			expectedTerm:  "",
			expectedLimit: 0,
		},
		{
			name:          "invalid limit",
			queryParams:   "term=test&limit=invalid",
			expectedError: "",
			expectedTerm:  "test",
			expectedLimit: 20, // default limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?"+tt.queryParams, nil)
			result, err := kithttp.DecodeSearchMediaRequest(context.Background(), req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, transport.SearchMediaRequest{Term: tt.expectedTerm, Limit: tt.expectedLimit}, result)
			}
		})
	}
}

func TestEncodeSearchMediaResponse(t *testing.T) {
	tests := []struct {
		name           string
		response       any
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid response",
			response: transport.SearchMediaResponse{
				ID:          1,
				SearchTerm:  "test",
				ResultCount: 1,
				Media:       []transport.Media{},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"search_term":"test","result_count":1,"media":[]}`,
		},
		{
			name:           "invalid response type",
			response:       "invalid response",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"failed to parse search media response, got invalid response"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			err := kithttp.EncodeSearchMediaResponse(context.Background(), recorder, tt.response)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
