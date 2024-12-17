package metrics

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"time"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

type Metrics struct {
	meter           metric.Meter
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	databaseCalls   metric.Int64Counter
	databaseErrors  metric.Int64Counter
	productCreated  metric.Int64Counter
	productUpdated  metric.Int64Counter
	productDeleted  metric.Int64Counter
	brandCreated    metric.Int64Counter
	brandDeleted    metric.Int64Counter
}

func InitProvider(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := sdkMetric.NewMeterProvider(
		sdkMetric.WithResource(res),
		sdkMetric.WithReader(sdkMetric.NewPeriodicReader(metricExporter)),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

func NewMetrics(ctx context.Context) (*Metrics, error) {
	meter := otel.Meter("ecommerce")

	requestCounter, err := meter.Int64Counter("http.request.total",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request counter: %w", err)
	}

	requestDuration, err := meter.Float64Histogram("http.request.duration",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request duration: %w", err)
	}

	databaseCalls, err := meter.Int64Counter("db.call.total",
		metric.WithDescription("Total number of database calls"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database calls counter: %w", err)
	}

	databaseErrors, err := meter.Int64Counter("db.error.total",
		metric.WithDescription("Total number of database errors"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database errors counter: %w", err)
	}

	productCreated, err := meter.Int64Counter("product.created.total",
		metric.WithDescription("Total number of products created"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product created counter: %w", err)
	}

	productUpdated, err := meter.Int64Counter("product.updated.total",
		metric.WithDescription("Total number of products updated"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product updated counter: %w", err)
	}

	productDeleted, err := meter.Int64Counter("product.deleted.total",
		metric.WithDescription("Total number of products deleted"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product deleted counter: %w", err)
	}

	brandCreated, err := meter.Int64Counter("brand.created.total",
		metric.WithDescription("Total number of brands created"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create brand created counter: %w", err)
	}

	brandDeleted, err := meter.Int64Counter("brand.deleted.total",
		metric.WithDescription("Total number of brands deleted"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create brand deleted counter: %w", err)
	}

	return &Metrics{
		meter:           meter,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
		databaseCalls:   databaseCalls,
		databaseErrors:  databaseErrors,
		productCreated:  productCreated,
		productUpdated:  productUpdated,
		productDeleted:  productDeleted,
		brandCreated:    brandCreated,
		brandDeleted:    brandDeleted,
	}, nil
}

func (m *Metrics) RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.path", path),
		attribute.Int("http.status_code", statusCode),
	}

	m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.requestDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
}

func (m *Metrics) RecordDatabaseCall(ctx context.Context, operation string) {
	m.databaseCalls.Add(ctx, 1, metric.WithAttributes(attribute.String("db.operation", operation)))
}

func (m *Metrics) RecordDatabaseError(ctx context.Context, operation string) {
	m.databaseErrors.Add(ctx, 1, metric.WithAttributes(attribute.String("db.operation", operation)))
}

func (m *Metrics) RecordProductCreated(ctx context.Context) {
	m.productCreated.Add(ctx, 1)
}

func (m *Metrics) RecordProductUpdated(ctx context.Context) {
	m.productUpdated.Add(ctx, 1)
}

func (m *Metrics) RecordProductDeleted(ctx context.Context) {
	m.productDeleted.Add(ctx, 1)
}

func (m *Metrics) RecordBrandCreated(ctx context.Context) {
	m.brandCreated.Add(ctx, 1)
}

func (m *Metrics) RecordBrandDeleted(ctx context.Context) {
	m.brandDeleted.Add(ctx, 1)
}
