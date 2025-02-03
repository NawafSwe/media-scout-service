package logging

import (
	"context"
)

// Logger interface for logging.
//
//go:generate mockgen -source=logger.go -destination=mock/logger.go -package=mock
type Logger interface {
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
}
