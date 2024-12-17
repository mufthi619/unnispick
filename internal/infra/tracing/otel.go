package tracing

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
}

type Tracer struct {
	tracer trace.Tracer
	logger *zap.Logger
}

func NewTracer(logger *zap.Logger) *Tracer {
	return &Tracer{
		tracer: otel.Tracer("ecommerce"),
		logger: logger,
	}
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

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdkTrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithResource(res),
		sdkTrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tracerProvider.Shutdown, nil
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

func (t *Tracer) StartFromEcho(c echo.Context, name string) (context.Context, trace.Span) {
	ctx := c.Request().Context()
	span := trace.SpanFromContext(ctx)

	attrs := []attribute.KeyValue{
		attribute.String("http.method", c.Request().Method),
		attribute.String("http.path", c.Request().URL.Path),
		attribute.String("http.user_agent", c.Request().UserAgent()),
	}
	span.SetAttributes(attrs...)

	return t.Start(ctx, name)
}

func (t *Tracer) End(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		t.logger.Error("Error in span",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
			zap.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	span.End()
}
