package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NawafSwe/media-scout-service/pkg/internal/transport"
	"net/http"
	"strconv"
)

// DecodeSearchMediaRequest function decodes search media request.
func DecodeSearchMediaRequest(_ context.Context, r *http.Request) (any, error) {
	term := r.URL.Query().Get("term")
	defaultLimit := 20
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = defaultLimit
	}
	if term == "" {
		return nil, errors.New("term shouldn't be empty")
	}

	return transport.SearchMediaRequest{Term: term, Limit: limit}, nil
}

// EncodeSearchMediaResponse function to encode media search response back.
func EncodeSearchMediaResponse(_ context.Context, w http.ResponseWriter, response any) error {
	_, ok := response.(transport.SearchMediaResponse)
	if !ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(w).Encode(map[string]any{
			"errors": fmt.Errorf("failed to parse search media response, got %v", response).Error(),
		})
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
