package itunes

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
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

// SearchResponse represents the response from the iTunes search API.
type SearchResponse struct {
	ResultCount int     `json:"resultCount"`
	Results     []Media `json:"results"`
}

// Client represents the iTunes API client.
type Client struct {
	httpClient http.Client
}

// NewClient creates a new iTunes API client.
func NewClient(tracer *trace.TracerProvider) *Client {
	otelTransport := otelhttp.NewTransport(nil, otelhttp.WithTracerProvider(tracer))
	return &Client{
		httpClient: http.Client{
			Transport: otelTransport,
		},
	}
}

// Search fetches media items from the iTunes API based on the search term.
func (c *Client) Search(ctx context.Context, term string, limit int) (SearchResponse, error) {
	url := fmt.Sprintf("https://itunes.apple.com/search?term=%s&limit=%d", term, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("failed to fetch data from iTunes API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return SearchResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return searchResponse, nil
}
