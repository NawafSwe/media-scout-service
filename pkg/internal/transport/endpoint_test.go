package transport_test

import (
	"context"
	"errors"
	"github.com/NawafSwe/media-scout-service/pkg/internal/transport"
	"github.com/NawafSwe/media-scout-service/pkg/internal/transport/mock"
	"testing"

	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMakeSearchMediaEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := mock.NewMockhandler(ctrl)
	endpoint := transport.MakeSearchMediaEndpoint(mockHandler)

	tests := []struct {
		name             string
		request          any
		mockSetup        func()
		expectedError    string
		expectedResponse any
	}{
		{
			name: "successful fetch",
			request: transport.SearchMediaRequest{
				Term:  "test",
				Limit: 1,
			},
			mockSetup: func() {
				mockHandler.EXPECT().FetchAndInsertMedia(gomock.Any(), "test", 1).Return(business.MediaResult{
					ID:          1,
					SearchTerm:  "test",
					ResultCount: 1,
					Media: []business.Media{
						{
							WrapperType: "track",
							Kind:        "song",
							ArtistID:    123,
							TrackID:     456,
							ArtistName:  "Test Artist",
							TrackName:   "Test Track",
						},
					},
				}, nil)
			},
			expectedError: "",
			expectedResponse: transport.SearchMediaResponse{
				ID:          1,
				SearchTerm:  "test",
				ResultCount: 1,
				Media: []transport.Media{
					{
						WrapperType: "track",
						Kind:        "song",
						ArtistID:    123,
						TrackID:     456,
						ArtistName:  "Test Artist",
						TrackName:   "Test Track",
					},
				},
			},
		},
		{
			name: "fetch error",
			request: transport.SearchMediaRequest{
				Term:  "test",
				Limit: 1,
			},
			mockSetup: func() {
				mockHandler.EXPECT().FetchAndInsertMedia(gomock.Any(), "test", 1).Return(business.MediaResult{}, errors.New("fetch error"))
			},
			expectedError:    "failed to fetch and insert media: fetch error",
			expectedResponse: transport.SearchMediaResponse{},
		},
		{
			name:             "invalid request type",
			request:          "invalid request",
			mockSetup:        func() {},
			expectedError:    "failed to parse search media request",
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := endpoint(context.Background(), tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, response)
			}
		})
	}
}
