package otel

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TracingConfig holds the configuration for OpenTelemetry tracing
type TracingConfig struct {
	ServiceName        string
	ServiceVersion     string
	Environment        string
	OTLPEndpoint       string
	SamplingRate       float64
	EnableConsole      bool
	MaxExportBatchSize int
	MaxQueueSize       int
}

// DefaultTracingConfig returns a default tracing configuration
func DefaultTracingConfig(serviceName string) *TracingConfig {
	return &TracingConfig{
		ServiceName:        serviceName,
		ServiceVersion:     getEnv("SERVICE_VERSION", "1.0.0"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		OTLPEndpoint:       getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		SamplingRate:       0.1, // Sample 10% of traces by default
		EnableConsole:      getEnv("ENVIRONMENT", "development") == "development",
		MaxExportBatchSize: 512,
		MaxQueueSize:       2048,
	}
}

// InitTracing initializes OpenTelemetry tracing
func InitTracing(serviceName string) (func(), error) {
	return InitTracingWithConfig(DefaultTracingConfig(serviceName))
}

// InitTracingWithConfig initializes OpenTelemetry tracing with custom configuration
func InitTracingWithConfig(config *TracingConfig) (func(), error) {
	ctx := context.Background()

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.ServiceVersion),
			semconv.DeploymentEnvironment(config.Environment),
			attribute.String("service.instance.id", getInstanceID()),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
		resource.WithContainer(),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP trace exporter
	traceExporter, err := createTraceExporter(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Create trace provider with sampling
	sampler := createSampler(config.SamplingRate)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithMaxExportBatchSize(config.MaxExportBatchSize),
			sdktrace.WithMaxQueueSize(config.MaxQueueSize),
			sdktrace.WithBatchTimeout(5*time.Second),
		),
		sdktrace.WithSampler(sampler),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global propagator for context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return shutdown function
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

// createTraceExporter creates an OTLP trace exporter
func createTraceExporter(ctx context.Context, config *TracingConfig) (sdktrace.SpanExporter, error) {
	// Create gRPC connection to OTLP collector
	conn, err := grpc.DialContext(ctx, config.OTLPEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Create OTLP trace exporter
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	return exporter, nil
}

// createSampler creates a sampler based on the sampling rate
func createSampler(samplingRate float64) sdktrace.Sampler {
	// Always sample errors and slow requests
	return sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(samplingRate),
		sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
		sdktrace.WithRemoteParentNotSampled(sdktrace.TraceIDRatioBased(samplingRate)),
	)
}

// Tracer returns a tracer for the given name
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("go-pro").Start(ctx, name, opts...)
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// SetSpanStatus sets the status of the current span
func SetSpanStatus(ctx context.Context, code trace.StatusCode, description string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(code, description)
}

// GetTraceID returns the trace ID from the context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().TraceID().String()
}

// GetSpanID returns the span ID from the context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	return span.SpanContext().SpanID().String()
}

// InstrumentHTTPHandler wraps an HTTP handler with tracing
func InstrumentHTTPHandler(handler http.Handler, operationName string) http.Handler {
	return otelhttp.NewHandler(handler, operationName,
		otelhttp.WithSpanOptions(
			trace.WithAttributes(
				semconv.HTTPMethod("GET"),
			),
		),
	)
}

// InstrumentHTTPClient wraps an HTTP client with tracing
func InstrumentHTTPClient(client *http.Client) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	client.Transport = otelhttp.NewTransport(client.Transport)
	return client
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInstanceID() string {
	// Try to get instance ID from environment
	if id := os.Getenv("INSTANCE_ID"); id != "" {
		return id
	}

	// Try to get hostname
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}

	return "unknown"
}

// Middleware for HTTP servers
type TracingMiddleware struct {
	serviceName string
}

// NewTracingMiddleware creates a new tracing middleware
func NewTracingMiddleware(serviceName string) *TracingMiddleware {
	return &TracingMiddleware{
		serviceName: serviceName,
	}
}

// Handler wraps an HTTP handler with tracing
func (m *TracingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Start span
		ctx, span := StartSpan(ctx, r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(r.Method),
				semconv.HTTPTarget(r.URL.Path),
				semconv.HTTPScheme(r.URL.Scheme),
				semconv.HTTPHost(r.Host),
				semconv.HTTPUserAgent(r.UserAgent()),
			),
		)
		defer span.End()

		// Create response writer wrapper to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		// Call next handler
		next.ServeHTTP(rw, r.WithContext(ctx))

		// Add response attributes
		span.SetAttributes(
			semconv.HTTPStatusCode(rw.statusCode),
		)

		// Set span status based on HTTP status code
		if rw.statusCode >= 400 {
			span.SetStatus(trace.StatusError, fmt.Sprintf("HTTP %d", rw.statusCode))
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Example usage:
//
// func main() {
//     // Initialize tracing
//     shutdown, err := otel.InitTracing("my-service")
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer shutdown()
//
//     // Create HTTP handler with tracing
//     mux := http.NewServeMux()
//     mux.HandleFunc("/api/users", handleUsers)
//
//     // Wrap with tracing middleware
//     middleware := otel.NewTracingMiddleware("my-service")
//     handler := middleware.Handler(mux)
//
//     // Start server
//     http.ListenAndServe(":8080", handler)
// }
//
// func handleUsers(w http.ResponseWriter, r *http.Request) {
//     ctx := r.Context()
//
//     // Add custom span attributes
//     otel.AddSpanAttributes(ctx,
//         attribute.String("user.id", "123"),
//         attribute.String("action", "list"),
//     )
//
//     // Create child span for database operation
//     ctx, span := otel.StartSpan(ctx, "database.query")
//     defer span.End()
//
//     // Simulate database query
//     users, err := fetchUsers(ctx)
//     if err != nil {
//         otel.RecordError(ctx, err)
//         otel.SetSpanStatus(ctx, trace.StatusError, err.Error())
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//
//     // Add event
//     otel.AddSpanEvent(ctx, "users.fetched",
//         attribute.Int("count", len(users)),
//     )
//
//     json.NewEncoder(w).Encode(users)
// }
