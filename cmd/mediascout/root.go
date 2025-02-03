package mediascout

import (
	"context"
	"fmt"

	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/NawafSwe/media-scout-service/pkg/worker"
	"github.com/jmoiron/sqlx"
)

// RunHTTPServer run http server.
func RunHTTPServer(ctx context.Context, db *sqlx.DB, cfg config.Config) error {
	w, err := worker.NewHTTPWorker(cfg, db, "media_scout.http_srv")
	if err != nil {
		return fmt.Errorf("failed to create http server: %w", err)
	}
	if err := w.Run(ctx); err != nil {
		return fmt.Errorf("failed to run http worker: %w", err)
	}
	return nil
}
