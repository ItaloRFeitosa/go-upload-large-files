package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/italorfeitosa/go-upload-large-files/internal/cleaner"
	"github.com/italorfeitosa/go-upload-large-files/internal/handler"
	"github.com/italorfeitosa/go-upload-large-files/pkg/storage"
	"github.com/italorfeitosa/go-upload-large-files/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	setupLogger()
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Application panicked", "error", r)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownOtel, err := telemetry.SetupOtel(ctx)
	if err != nil {
		slog.Error("Failed to create S3 client", "error", err)
		panic(err)
	}
	defer func() {
		if err := shutdownOtel(ctx); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", "error", err)
		}
	}()

	s3Client, err := storage.CreateS3Client()
	if err != nil {
		slog.Error("Failed to create S3 client", "error", err)
		panic(err)
	}

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "uploads"
		slog.Warn("S3 bucket name is not set, using default bucket", "bucket", bucket)
	}

	server := newServer(s3Client, bucket)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "error", err)
			stop()
		}
	}()

	teardownCleaner := cleaner.Start(ctx, s3Client, bucket)
	defer func() {
		if err := teardownCleaner(); err != nil {
			slog.Error("Failed to shutdown cleaner", "error", err)
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		slog.Error("Error during server shutdown", "error", err)
		panic(err)
	}

	slog.Info("Server stopped")
}

func setupLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default log level
	}

	logFormat := os.Getenv("LOG_FORMAT")
	var handler slog.Handler
	switch logFormat {
	case "JSON":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))
}

func newServer(s3Client *s3.Client, bucket string) *http.Server {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	mux := http.NewServeMux()

	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	handleFunc("/health", handler.HealthCheck)
	handleFunc("/upload", handler.Upload(s3Client, bucket))

	handler := otelhttp.NewHandler(mux, "/")

	slog.Info("Starting server", "port", port)

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}
