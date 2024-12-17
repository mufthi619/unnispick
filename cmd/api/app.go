package api

import (
	"Unnispick/internal/config"
	"Unnispick/internal/domain/delivery/router"
	"Unnispick/pkg/databases"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	cfg    *config.Config
	echo   *echo.Echo
	router *router.Router
	db     databases.DB
	logger *zap.Logger
}

func NewApp(
	cfg *config.Config,
	echo *echo.Echo,
	router *router.Router,
	db databases.DB,
	logger *zap.Logger,
) *App {
	return &App{
		cfg:    cfg,
		echo:   echo,
		router: router,
		db:     db,
		logger: logger,
	}
}

func (a *App) Start() error {
	// Setup routes
	a.router.Setup()

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
		if err := a.echo.Start(addr); err != nil {
			a.logger.Error("shutting down the server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.Timeout.Write)
	defer cancel()

	// Close database connection
	if err := a.db.Close(); err != nil {
		a.logger.Error("failed to close database connection", zap.Error(err))
	}

	// Shutdown server
	if err := a.echo.Shutdown(ctx); err != nil {
		a.logger.Error("failed to shutdown server", zap.Error(err))
	}

	return nil
}
