package mediadb

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

// Media represents a single media item with various attributes.
type Media struct {
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

type Medias []Media

// Scan implements the sql.Scanner interface for Medias.
//
// for more details see https://stackoverflow.com/questions/41375563/unsupported-scan-storing-driver-value-type-uint8-into-type-string
func (m *Medias) Scan(src any) error {
	if src == nil {
		*m = Medias{} // Make sure it's initialized if NULL was stored.
		return nil
	}
	v, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("invalid data received, expected []byte got %T", src)
	}
	if err := json.Unmarshal(v, m); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from bytes: %w", err)
	}
	return nil
}

// Value implements the driver.Valuer interface for Medias.
func (m *Medias) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// MediaResult represents the result user searched for.
type MediaResult struct {
	ID          int64     `db:"id"`
	SearchTerm  string    `db:"search_term"`
	Media       Medias    `db:"returned_result"` // JSONB field
	ResultCount int       `db:"result_count"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// MediaRepositoryImpl is the implementation of MediaRepository.
type MediaRepositoryImpl struct {
	db *sqlx.DB
}

// NewMediaRepository creates a new instance of MediaRepository.
func NewMediaRepository(db *sqlx.DB) *MediaRepositoryImpl {
	return &MediaRepositoryImpl{db: db}
}

// mapBusinessToDBModel maps a business.MediaResult to a MediaResult.
func mapBusinessToDBModel(media business.MediaResult) MediaResult {
	return MediaResult{
		ID:         media.ID,
		SearchTerm: media.SearchTerm,
		Media: lo.Map(media.Media, func(m business.Media, _ int) Media {
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
		}),
		ResultCount: media.ResultCount,
	}
}

// InsertMedia inserts a new media result into the database.
func (repo *MediaRepositoryImpl) InsertMedia(ctx context.Context, media business.MediaResult) (int64, error) {
	dbMedia := mapBusinessToDBModel(media)
	mediaResult, err := json.Marshal(dbMedia.Media)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal media result: %w", err)
	}
	query := `
		INSERT INTO media_result (search_term, returned_result, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64

	if err = repo.db.QueryRowContext(ctx, query, dbMedia.SearchTerm, mediaResult, time.Now().UTC(), time.Now().UTC()).Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to insert media to db: %w", err)
	}

	return id, nil
}
