package main

import (
	"context"
	"github.com/NawafSwe/media-scout-service/cmd/config"
	"github.com/NawafSwe/media-scout-service/cmd/mediascout"
	"github.com/NawafSwe/media-scout-service/pkg/db"
	"go.opentelemetry.io/otel"
	otelgrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelgrpctrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	cfg, err := config.NewConfig(".", ".env")
	if err != nil {
		log.Fatalf("err loading config, err: %v", err)
	}
	tp, err := initTracer(cfg)
	if err != nil {
		log.Fatalf("err initializing trace, err: %v", err)
	}
	dbConn, err := db.NewDBConn(cfg.DB, tp)
	if err != nil {
		log.Fatalf("err creating db conn, err: %v", err)
	}
	if err := mediascout.RunHTTPServer(context.Background(), tp, dbConn, cfg); err != nil {
		log.Fatalf("failed to run http server: %v", err)
	}
}

func initTracer(cfg config.Config) (*trace.TracerProvider, error) {
	ctx := context.Background()
	// Create the gRPC connection
	grpcConn, err := grpc.NewClient(cfg.General.Tracing.ReceiverEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// trace exporter
	tpExporter, err := otelgrpctrace.New(ctx, otelgrpctrace.WithGRPCConn(grpcConn))
	if err != nil {
		return nil, err
	}

	// Create the OTLP gRPC exporter for logs
	exporter, err := otelgrpc.New(ctx, otelgrpc.WithGRPCConn(grpcConn))
	if err != nil {
		return nil, err
	}

	// Create the resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.General.ServiceName),
			semconv.ServiceNamespaceKey.String(cfg.General.ServiceName),
			semconv.ServiceVersionKey.String(cfg.General.AppVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.General.AppEnvironment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create the tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(tpExporter),
		trace.WithResource(res),
	)

	// Create the log provider
	logProcessor := otellog.NewBatchProcessor(exporter)
	logProvider := otellog.NewLoggerProvider(otellog.WithProcessor(logProcessor))
	global.SetLoggerProvider(logProvider)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	return tp, nil
}
