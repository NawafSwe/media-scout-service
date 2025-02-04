package business_test

import (
	"context"
	"errors"
	"github.com/NawafSwe/media-scout-service/pkg/internal/business/mock"
	"testing"

	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAndInsertMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockmediaRepository(ctrl)
	mockFetcher := mock.NewMockmediaFetcher(ctrl)
	mockLogger := mock.NewMocklogger(ctrl)

	handler := business.NewSearchMediaHandler(mockRepo, mockFetcher, mockLogger)

	tests := []struct {
		name           string
		term           string
		limit          int
		mockSetup      func()
		expectedError  string
		expectedResult business.MediaResult
	}{
		{
			name:  "successful fetch and insert",
			term:  "test",
			limit: 1,
			mockSetup: func() {
				mockFetcher.EXPECT().FetchMediaByTerm(gomock.Any(), "test", 1).Return(business.MediaResult{
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
				mockRepo.EXPECT().InsertMedia(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			expectedError: "",
			expectedResult: business.MediaResult{
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
			},
		},
		{
			name:  "fetch error",
			term:  "test",
			limit: 1,
			mockSetup: func() {
				mockFetcher.EXPECT().FetchMediaByTerm(gomock.Any(), "test", 1).Return(business.MediaResult{}, errors.New("fetch error"))
				mockLogger.EXPECT().ErrorContext(gomock.Any(), "failed to fetch media", "error", errors.New("fetch error"))
			},
			expectedError:  "failed to fetch media: fetch error",
			expectedResult: business.MediaResult{},
		},
		{
			name:  "insert error",
			term:  "test",
			limit: 1,
			mockSetup: func() {
				mockFetcher.EXPECT().FetchMediaByTerm(gomock.Any(), "test", 1).Return(business.MediaResult{
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
				mockRepo.EXPECT().InsertMedia(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("insert error"))
				mockLogger.EXPECT().ErrorContext(gomock.Any(), "failed to insert media", "error", errors.New("insert error"))
			},
			expectedError: "",
			expectedResult: business.MediaResult{
				ID:          0,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := handler.FetchAndInsertMedia(context.Background(), tt.term, tt.limit)

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
