package worker

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/NawafSwe/media-scout-service/pkg/clients/itunes"
	"github.com/NawafSwe/media-scout-service/pkg/internal/business"
	"github.com/NawafSwe/media-scout-service/pkg/internal/repository/mediadb"
	"github.com/NawafSwe/media-scout-service/pkg/internal/repository/mediafetcher"
	"github.com/NawafSwe/media-scout-service/pkg/logging"
	"github.com/NawafSwe/media-scout-service/pkg/transport"
	kithttptransport "github.com/NawafSwe/media-scout-service/pkg/transport/http"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// HTTPWorker represents http worker.
type HTTPWorker struct {
	cfg     config.Config
	Name    string
	db      *sqlx.DB
	lgr     logging.Logger
	port    int
	router  *mux.Router
	srv     *http.Server
	signals chan os.Signal
}

// NewHTTPWorker function creates http worker.
func NewHTTPWorker(cfg config.Config, db *sqlx.DB, name string) (*HTTPWorker, error) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	lgr := slog.New(handler)
	lgrWithAttrs := lgr.With("service", name)
	return &HTTPWorker{
		cfg:     cfg,
		Name:    name,
		lgr:     lgrWithAttrs,
		db:      db,
		port:    cfg.HTTP.Port,
		router:  mux.NewRouter(),
		signals: make(chan os.Signal, 1),
	}, nil
}

func (h *HTTPWorker) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", h.port))
	if err != nil {
		h.lgr.ErrorContext(ctx, "failed to listen to port", "port", h.port)
		return fmt.Errorf("failed to listen on port %d", h.port)
	}
	h.registerHandlers()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", h.port),
		Handler: enableCORS(h.router),
	}
	h.srv = &srv

	go func() {
		if err := srv.Serve(lis); err != nil {
			h.lgr.ErrorContext(ctx, "failed to serve http server", "port", h.port)
		}
	}()
	h.lgr.InfoContext(ctx, "running server", "port", h.port)
	signal.Notify(h.signals, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	<-h.signals
	h.lgr.InfoContext(ctx, "graceful shutdown started")
	ctx, cancel := context.WithTimeout(ctx, h.cfg.HTTP.GracefulShutdown)
	defer cancel()
	if err := h.srv.Shutdown(ctx); err != nil {
		h.lgr.ErrorContext(ctx, "failed to stop server in graceful shutdown", "error", err.Error())
		return fmt.Errorf("failed to stop server in graceful shutdown: %v", err)
	}
	h.lgr.InfoContext(ctx, "stopped server gracefully.")
	return nil
}

func (h *HTTPWorker) SIGINT() {
	h.signals <- syscall.SIGINT
}

func (h *HTTPWorker) registerHandlers() {
	r := h.router.PathPrefix("").Subrouter()
	r.HandleFunc("/health", h.healthHandler).Methods(http.MethodGet)
	v1APIs := r.PathPrefix("/api/v1").Subrouter()

	v1APIs.Handle("/media/search", makeSearchMediaHandler(h.db, h.lgr)).Methods(http.MethodGet)
}

func (h *HTTPWorker) healthHandler(r http.ResponseWriter, _ *http.Request) {
	if err := h.db.DB.Ping(); err != nil {
		r.WriteHeader(http.StatusServiceUnavailable)
		_, _ = r.Write([]byte(fmt.Sprintf("%s is unavailable due to db unavailability %s", config.ServiceName, err.Error())))
		return
	}
	r.WriteHeader(http.StatusOK)
	_, _ = r.Write([]byte(fmt.Sprintf("%s is healthy", config.ServiceName)))
}

// makeSearchMediaHandler function to return http handler for search media.
func makeSearchMediaHandler(db *sqlx.DB, lgr logging.Logger, middlewares ...endpoint.Middleware) http.Handler {
	itunesClient := itunes.NewClient()
	mediaFetcher := mediafetcher.NewMediaFetcher(itunesClient)
	mediaDBRepo := mediadb.NewMediaRepository(db)
	handler := business.NewSearchMediaHandler(mediaDBRepo, mediaFetcher, lgr)
	ep := transport.MakeSearchMediaEndpoint(handler)
	// applying middlewares, if any.
	if middlewares != nil {
		for _, m := range middlewares {
			ep = m(ep)
		}
	}
	return kithttp.NewServer(ep, kithttptransport.DecodeSearchMediaRequest, kithttptransport.EncodeSearchMediaResponse)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
