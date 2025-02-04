package business

import (
	"context"
	"fmt"
)

// Media represents a single media item with various attributes.
type Media struct {
	WrapperType            string
	Kind                   string
	ArtistID               int
	CollectionID           int
	TrackID                int
	ArtistName             string
	CollectionName         string
	TrackName              string
	ArtistViewURL          string
	CollectionViewURL      string
	FeedURL                string
	TrackViewURL           string
	ArtworkURL30           string
	ArtworkURL60           string
	ArtworkURL100          string
	ReleaseDate            string
	CollectionExplicitness string
	TrackExplicitness      string
	TrackCount             int
	TrackTimeMillis        int
	Country                string
	Currency               string
	PrimaryGenreName       string
	ContentAdvisoryRating  string
	ArtworkURL600          string
	GenreIDs               []string
	Genres                 []string
}

// MediaResult represents the result user searched for.
type MediaResult struct {
	ID          int64
	SearchTerm  string
	Media       []Media
	ResultCount int
}

//go:generate mockgen -source=search_media.go -destination=mock/search_media.go -package=mock
type (
	// mediaRepository defines the interface for media repository operations.
	mediaRepository interface {
		InsertMedia(ctx context.Context, media MediaResult) (int64, error)
	}
	// mediaFetcher defines the interface for fetching media.
	mediaFetcher interface {
		FetchMediaByTerm(ctx context.Context, term string, limit int) (MediaResult, error)
	}
	// logger logging error.
	logger interface {
		ErrorContext(ctx context.Context, msg string, args ...any)
	}
)

type SearchMediaHandler struct {
	repo    mediaRepository
	fetcher mediaFetcher
	lgr     logger
}

// NewSearchMediaHandler creates a new instance of SearchMediaHandler.
func NewSearchMediaHandler(repo mediaRepository, fetcher mediaFetcher, lgr logger) SearchMediaHandler {
	return SearchMediaHandler{repo: repo, fetcher: fetcher, lgr: lgr}
}

// FetchAndInsertMedia fetches media by term, inserts it into the repository, and returns the result.
func (h SearchMediaHandler) FetchAndInsertMedia(ctx context.Context, term string, limit int) (MediaResult, error) {
	// Fetch media by term
	mediaResult, err := h.fetcher.FetchMediaByTerm(ctx, term, limit)
	if err != nil {
		h.lgr.ErrorContext(ctx, "failed to fetch media", "error", err)
		return MediaResult{}, fmt.Errorf("failed to fetch media: %w", err)
	}

	// Insert media into the repository
	id, err := h.repo.InsertMedia(ctx, mediaResult)
	// If we failed to insert to db, it is ok to return to requester the result.
	if err != nil {
		h.lgr.ErrorContext(ctx, "failed to insert media", "error", err)
	}
	mediaResult.ID = id
	return mediaResult, nil
}
