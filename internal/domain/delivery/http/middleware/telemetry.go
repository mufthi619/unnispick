package middleware

import (
	"Unnispick/internal/infra/metrics"
	"Unnispick/internal/infra/tracing"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"time"
)

type TelemetryMiddleware struct {
	logger     *zap.Logger
	tracer     *tracing.Tracer
	metrics    *metrics.Metrics
	propagator propagation.TextMapPropagator
}

func NewTelemetryMiddleware(
	logger *zap.Logger,
	tracer *tracing.Tracer,
	metrics *metrics.Metrics,
) *TelemetryMiddleware {
	return &TelemetryMiddleware{
		logger:     logger,
		tracer:     tracer,
		metrics:    metrics,
		propagator: propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
	}
}

func (m *TelemetryMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			startTime := time.Now()

			// Extract trace context from incoming request headers
			ctx := m.propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))
			req = req.WithContext(ctx)
			c.SetRequest(req)

			// Start a new span for this request
			spanCtx, span := m.tracer.Start(ctx, "http.request")
			defer span.End()

			// Set basic span attributes
			span.SetAttributes(
				attribute.String("http.method", req.Method),
				attribute.String("http.path", req.URL.Path),
				attribute.String("http.user_agent", req.UserAgent()),
				attribute.String("http.host", req.Host),
				attribute.String("http.scheme", req.URL.Scheme),
			)

			// Store the span context in the echo.Context
			c.Set("span_context", spanCtx)

			// Call the next handler
			err := next(c)

			// Record response status and duration
			status := c.Response().Status
			duration := time.Since(startTime)

			// Update span with response information
			span.SetAttributes(
				attribute.Int("http.status_code", status),
				attribute.Float64("http.request_duration_ms", float64(duration.Milliseconds())),
			)

			// If there was an error, record it in the span
			if err != nil {
				m.logger.Error("request error",
					zap.Error(err),
					zap.String("method", req.Method),
					zap.String("path", req.URL.Path),
					zap.Int("status", status),
					zap.Duration("duration", duration),
				)
				span.RecordError(err)
			}

			// Record metrics
			m.metrics.RecordRequest(ctx, req.Method, req.URL.Path, status, duration)

			// Add trace context to response headers
			m.propagator.Inject(spanCtx, propagation.HeaderCarrier(c.Response().Header()))

			return err
		}
	}
}

func (m *TelemetryMiddleware) ExtractTraceID(c echo.Context) string {
	if spanCtx, ok := c.Get("span_context").(string); ok {
		return spanCtx
	}
	return ""
}

func (m *TelemetryMiddleware) MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			err := next(c)
			duration := time.Since(startTime)

			// Record basic request metrics
			m.metrics.RecordRequest(
				c.Request().Context(),
				c.Request().Method,
				c.Request().URL.Path,
				c.Response().Status,
				duration,
			)

			return err
		}
	}
}

func (m *TelemetryMiddleware) TracingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := m.propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

			spanCtx, span := m.tracer.Start(ctx, "http.request")
			defer span.End()

			// Add trace context to response headers
			m.propagator.Inject(spanCtx, propagation.HeaderCarrier(c.Response().Header()))

			// Update request context with span context
			c.SetRequest(req.WithContext(spanCtx))

			return next(c)
		}
	}
}
