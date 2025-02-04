package mediafetcher_test

import (
	"context"
	"errors"
	"github.com/NawafSwe/media-scout-service/pkg/clients/itunes"
	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/NawafSwe/media-scout-service/pkg/internal/repository/mediafetcher"
	"github.com/NawafSwe/media-scout-service/pkg/internal/repository/mediafetcher/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchMediaByTerm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMocksearcherClient(ctrl)
	fetcher := mediafetcher.NewMediaFetcher(mockClient)

	tests := []struct {
		name           string
		term           string
		limit          int
		mockSetup      func()
		expectedError  string
		expectedResult business.MediaResult
	}{
		{
			name:  "successful fetch",
			term:  "test",
			limit: 1,
			mockSetup: func() {
				mockClient.EXPECT().Search(gomock.Any(), "test", 1).Return(itunes.SearchResponse{
					ResultCount: 1,
					Results: []itunes.Media{
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
			expectedResult: business.MediaResult{
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
			},
		},
		{
			name:  "fetch error",
			term:  "test",
			limit: 1,
			mockSetup: func() {
				mockClient.EXPECT().Search(gomock.Any(), "test", 1).Return(itunes.SearchResponse{}, errors.New("search error"))
			},
			expectedError:  "failed to fetch media by term: search error",
			expectedResult: business.MediaResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := fetcher.FetchMediaByTerm(context.Background(), tt.term, tt.limit)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
