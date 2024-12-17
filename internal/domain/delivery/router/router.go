package router

import (
	"Unnispick/internal/domain/delivery/http/handler"
	"Unnispick/internal/domain/delivery/http/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	e               *echo.Echo
	brandHandler    *handler.BrandHandler
	productHandler  *handler.ProductHandler
	telemetryMiddle *middleware.TelemetryMiddleware
}

func NewRouter(
	e *echo.Echo,
	brandHandler *handler.BrandHandler,
	productHandler *handler.ProductHandler,
	telemetryMiddle *middleware.TelemetryMiddleware,
) *Router {
	return &Router{
		e:               e,
		brandHandler:    brandHandler,
		productHandler:  productHandler,
		telemetryMiddle: telemetryMiddle,
	}
}

func (r *Router) Setup() {
	// Middleware
	r.e.Use(echoMiddleware.Logger())
	r.e.Use(echoMiddleware.Recover())
	r.e.Use(echoMiddleware.CORS())
	r.e.Use(r.telemetryMiddle.Middleware())

	// Health Check
	r.e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
		})
	})

	// API v1 group
	v1 := r.e.Group("/api/v1")

	// Brand routes
	brands := v1.Group("/brands")
	brands.POST("", r.brandHandler.Create)
	brands.GET("", r.brandHandler.GetAll)
	brands.GET("/:id", r.brandHandler.GetByID)
	brands.PUT("/:id", r.brandHandler.Update)
	brands.DELETE("/:id", r.brandHandler.Delete)

	// Product routes
	products := v1.Group("/products")
	products.POST("", r.productHandler.Create)
	products.GET("", r.productHandler.GetAll)
	products.GET("/:id", r.productHandler.GetByID)
	products.PUT("/:id", r.productHandler.Update)
	products.DELETE("/:id", r.productHandler.Delete)

	// When we add Swagger, we'll add it here
	// r.e.GET("/swagger/*", echoSwagger.WrapHandler)
}
