package transport

import (
	"context"
	"fmt"

	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/go-kit/kit/endpoint"
	"github.com/samber/lo"
)

//go:generate mockgen -source=endpoint.go -destination=mock/endpoint.go -package=mock
type handler interface {
	FetchAndInsertMedia(ctx context.Context, term string, limit int) (business.MediaResult, error)
}
type (
	// SearchMediaRequest represents the received request to search for media.
	SearchMediaRequest struct {
		Term  string
		Limit int
	}

	// Media represents a single media item with various attributes.
	Media struct {
		WrapperType            string   `json:"wrapperType"`
		Kind                   string   `json:"kind"`
		ArtistID               int      `json:"artistId"`
		CollectionID           int      `json:"collectionId"`
		TrackID                int      `json:"trackId"`
		ArtistName             string   `json:"artistName"`
		CollectionName         string   `json:"collectionName"`
		TrackName              string   `json:"trackName"`
		ArtistViewURL          string   `json:"artistViewUrl"`
		CollectionViewURL      string   `json:"collectionViewUrl"`
		FeedURL                string   `json:"feedUrl"`
		TrackViewURL           string   `json:"trackViewUrl"`
		ArtworkURL30           string   `json:"artworkUrl30"`
		ArtworkURL60           string   `json:"artworkUrl60"`
		ArtworkURL100          string   `json:"artworkUrl100"`
		ReleaseDate            string   `json:"releaseDate"`
		CollectionExplicitness string   `json:"collectionExplicitness"`
		TrackExplicitness      string   `json:"trackExplicitness"`
		TrackCount             int      `json:"trackCount"`
		TrackTimeMillis        int      `json:"trackTimeMillis"`
		Country                string   `json:"country"`
		Currency               string   `json:"currency"`
		PrimaryGenreName       string   `json:"primaryGenreName"`
		ContentAdvisoryRating  string   `json:"contentAdvisoryRating"`
		ArtworkURL600          string   `json:"artworkUrl600"`
		GenreIDs               []string `json:"genreIds"`
		Genres                 []string `json:"genres"`
	}

	// SearchMediaResponse represents the media user searched for.
	SearchMediaResponse struct {
		ID          int64   `json:"id"`
		SearchTerm  string  `json:"search_term"`
		ResultCount int     `json:"result_count"`
		Media       []Media `json:"media"`
	}
)

// MakeSearchMediaEndpoint function to make search media endpoint call.
func MakeSearchMediaEndpoint(handler handler) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		body, ok := request.(SearchMediaRequest)
		if !ok {
			return nil, fmt.Errorf("failed to parse search media request")
		}

		res, err := handler.FetchAndInsertMedia(ctx, body.Term, body.Limit)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch and insert media: %w", err)
		}

		media := lo.Map(res.Media, func(m business.Media, _ int) Media {
			return Media{
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
		})

		return SearchMediaResponse{
			ID:          res.ID,
			SearchTerm:  res.SearchTerm,
			ResultCount: res.ResultCount,
			Media:       media,
		}, nil
	}
}
