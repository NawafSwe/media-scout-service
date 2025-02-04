package mediafetcher

import (
	"context"
	"fmt"

	"github.com/NawafSwe/media-scout-service/pkg/clients/itunes"
	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/samber/lo"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type (
	searcherClient interface {
		Search(ctx context.Context, term string, limit int) (itunes.SearchResponse, error)
	}
)

// MediaFetcher is a service that fetches media from iTunes.
type MediaFetcher struct {
	client searcherClient
}

// NewMediaFetcher creates a new instance of MediaFetcherService.
func NewMediaFetcher(client searcherClient) *MediaFetcher {
	return &MediaFetcher{client: client}
}

func (s *MediaFetcher) FetchMediaByTerm(ctx context.Context, term string, limit int) (business.MediaResult, error) {
	response, err := s.client.Search(ctx, term, limit)
	if err != nil {
		return business.MediaResult{}, fmt.Errorf("failed to fetch media by term: %w", err)
	}

	mediaResult := business.MediaResult{
		SearchTerm:  term,
		ResultCount: response.ResultCount,
		Media: lo.Map(response.Results, func(m itunes.Media, _ int) business.Media {
			return business.Media{
				WrapperType:            m.WrapperType,
				Kind:                   m.Kind,
				ArtistID:               m.ArtistID,
				CollectionID:           m.CollectionID,
				TrackID:                m.TrackID,
				ArtistName:             m.ArtistName,
				CollectionName:         m.CollectionName,
				TrackName:              m.TrackName,
				ArtistViewURL:          m.ArtistViewURL,
				CollectionViewURL:      m.CollectionViewURL,
				FeedURL:                m.FeedURL,
				TrackViewURL:           m.TrackViewURL,
				ArtworkURL30:           m.ArtworkURL30,
				ArtworkURL60:           m.ArtworkURL60,
				ArtworkURL100:          m.ArtworkURL100,
				ReleaseDate:            m.ReleaseDate,
				CollectionExplicitness: m.CollectionExplicitness,
				TrackExplicitness:      m.TrackExplicitness,
				TrackCount:             m.TrackCount,
				TrackTimeMillis:        m.TrackTimeMillis,
				Country:                m.Country,
				Currency:               m.Currency,
				PrimaryGenreName:       m.PrimaryGenreName,
				ContentAdvisoryRating:  m.ContentAdvisoryRating,
				ArtworkURL600:          m.ArtworkURL600,
				GenreIDs:               m.GenreIDs,
				Genres:                 m.Genres,
			}
		}),
	}

	return mediaResult, nil
}
