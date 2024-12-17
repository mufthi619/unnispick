//go:build wireinject
// +build wireinject

package api

import (
	"Unnispick/internal/config"
	"Unnispick/internal/domain/delivery/http/handler"
	"Unnispick/internal/domain/delivery/http/middleware"
	"Unnispick/internal/domain/delivery/router"
	"Unnispick/internal/domain/repository"
	"Unnispick/internal/domain/service"
	"Unnispick/internal/infra/metrics"
	"Unnispick/internal/infra/tracing"
	"Unnispick/pkg/databases"
	"Unnispick/pkg/databases/postgres"
	"Unnispick/pkg/logger"
	"Unnispick/pkg/validator"
	"context"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var configSet = wire.NewSet(
	config.Load,
)

func provideContext() context.Context {
	return context.Background()
}

func provideEcho() *echo.Echo {
	return echo.New()
}

func provideLoggerConfig(cfg *config.Config) logger.Config {
	return logger.Config{
		Environment: cfg.Logger.Environment,
		LogLevel:    cfg.Logger.Level,
	}
}

func provideZapLogger(l *logger.Logger) *zap.Logger {
	return l.WithContext(context.Background())
}

func provideDB(conn *postgres.Database) *gorm.DB {
	return conn.DB()
}

func provideDatabaseOptions(cfg *config.Config) postgres.Options {
	return postgres.Options{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		Database:     cfg.Database.Name,
		SSLMode:      cfg.Database.SSLMode,
		MaxOpenConns: cfg.Database.Pool.MaxOpen,
		MaxIdleConns: cfg.Database.Pool.MaxIdle,
		MaxLifetime:  cfg.Database.Pool.MaxLifetime,
	}
}

var infraSet = wire.NewSet(
	provideContext,
	provideEcho,
	provideDB,
	provideDatabaseOptions,
	provideLoggerConfig,
	logger.NewLogger,
	provideZapLogger,
	postgres.NewConnection,
	wire.Bind(new(databases.DB), new(*postgres.Database)),
	tracing.NewTracer,
	metrics.NewMetrics,
	validator.NewValidator,
)

var repositorySet = wire.NewSet(
	repository.NewBrandRepository,
	repository.NewProductRepository,
)

var serviceSet = wire.NewSet(
	service.NewBrandService,
	service.NewProductService,
)

var handlerSet = wire.NewSet(
	handler.NewBrandHandler,
	handler.NewProductHandler,
)

var middlewareSet = wire.NewSet(
	middleware.NewTelemetryMiddleware,
)

var routerSet = wire.NewSet(
	router.NewRouter,
)

func InitializeApp() (*App, error) {
	wire.Build(
		NewApp,
		configSet,
		infraSet,
		repositorySet,
		serviceSet,
		handlerSet,
		middlewareSet,
		routerSet,
	)
	return nil, nil
}
