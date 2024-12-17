package postgres

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"time"
)

func WithTelemetry() gorm.Plugin {
	return &telemetryPlugin{}
}

type telemetryPlugin struct{}

func (p *telemetryPlugin) Name() string {
	return "telemetryPlugin"
}

func (p *telemetryPlugin) Initialize(db *gorm.DB) error {
	// Add callbacks for Create, Query, Update, Delete operations
	if err := db.Callback().Create().Before("gorm:create").Register("telemetry:before_create", before("db.create")); err != nil {
		return err
	}
	if err := db.Callback().Query().Before("gorm:query").Register("telemetry:before_query", before("db.query")); err != nil {
		return err
	}
	if err := db.Callback().Update().Before("gorm:update").Register("telemetry:before_update", before("db.update")); err != nil {
		return err
	}
	if err := db.Callback().Delete().Before("gorm:delete").Register("telemetry:before_delete", before("db.delete")); err != nil {
		return err
	}

	// Add callbacks for after operations to end spans
	if err := db.Callback().Create().After("gorm:create").Register("telemetry:after_create", after()); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("telemetry:after_query", after()); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("telemetry:after_update", after()); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("telemetry:after_delete", after()); err != nil {
		return err
	}

	return nil
}

func before(operation string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		tracer := db.Statement.Context.Value("tracer").(trace.Tracer)
		ctx, span := tracer.Start(db.Statement.Context, operation)
		db.Statement.Context = ctx
		db.InstanceSet("telemetry:span", span)
		db.InstanceSet("telemetry:start_time", time.Now())
	}
}

func after() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		spanValue, exists := db.InstanceGet("telemetry:span")
		if !exists {
			return
		}

		span := spanValue.(trace.Span)
		defer span.End()

		if db.Error != nil {
			span.RecordError(db.Error)
			span.SetStatus(codes.Error, db.Error.Error())
		}

		if db.Statement != nil && db.Statement.SQL.String() != "" {
			span.SetAttributes(
				attribute.String("db.statement", db.Statement.SQL.String()),
				attribute.Int64("db.rows_affected", db.RowsAffected),
			)
		}
	}
}
