package handler

import (
	"Unnispick/internal/domain/entity"
	"Unnispick/internal/infra/metrics"
	"Unnispick/internal/infra/tracing"
	"Unnispick/pkg/validator"
	"Unnispick/utils/response_formatter"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	service  entity.ProductService
	logger   *zap.Logger
	tracer   *tracing.Tracer
	metrics  *metrics.Metrics
	validate *validator.Validator
}

func NewProductHandler(
	service entity.ProductService,
	logger *zap.Logger,
	tracer *tracing.Tracer,
	metrics *metrics.Metrics,
	validate *validator.Validator,
) *ProductHandler {
	return &ProductHandler{
		service:  service,
		logger:   logger,
		tracer:   tracer,
		metrics:  metrics,
		validate: validate,
	}
}

// Create
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body entity.CreateProductRequest true "Product creation request"
// @Success 201 {object} response_formatter.Response{data=entity.ProductResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /products [post]
func (h *ProductHandler) Create(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.product.Create")
	defer span.End()

	var req entity.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid request body",
			[]string{err.Error()},
		))
	}

	if err := h.validate.Validate(ctx, req); err != nil {
		validationErrors := h.validate.ExtractValidationErrors(err)
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, ve.Message)
		}
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Validation failed",
			errorMessages,
		))
	}

	product, err := h.service.Create(ctx, req)
	if err != nil {
		h.logger.Error("failed to create product", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response_formatter.Error(
			http.StatusInternalServerError,
			"Failed to create product",
			[]string{err.Error()},
		))
	}

	h.metrics.RecordProductCreated(ctx)
	return c.JSON(http.StatusCreated, response_formatter.Created(product, "Product created successfully"))
}

// GetAll
// @Summary Get all products with pagination and filters
// @Description Get a list of all products with pagination and filtering support
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Param brand_id query string false "Filter by brand ID"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param min_qty query int false "Minimum quantity filter"
// @Param max_qty query int false "Maximum quantity filter"
// @Success 200 {object} response_formatter.Response{data=[]entity.ProductResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /products [get]
func (h *ProductHandler) GetAll(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.product.GetAll")
	defer span.End()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	page, perPage = response_formatter.ValidatePagination(page, perPage)

	filter := entity.ProductFilterRequest{
		Page:    perPage,
		PerPage: response_formatter.CalculateOffset(page, perPage),
	}

	// Parse optional filters
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		if id, err := uuid.Parse(brandID); err == nil {
			filter.BrandID = id
		}
	}
	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filter.MinPrice = price
		}
	}
	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filter.MaxPrice = price
		}
	}
	if minQty := c.QueryParam("min_qty"); minQty != "" {
		if qty, err := strconv.Atoi(minQty); err == nil {
			filter.MinQty = qty
		}
	}
	if maxQty := c.QueryParam("max_qty"); maxQty != "" {
		if qty, err := strconv.Atoi(maxQty); err == nil {
			filter.MaxQty = qty
		}
	}

	products, total, err := h.service.GetAll(ctx, filter)
	if err != nil {
		h.logger.Error("failed to get products", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response_formatter.Error(
			http.StatusInternalServerError,
			"Failed to get products",
			[]string{err.Error()},
		))
	}

	return c.JSON(http.StatusOK, response_formatter.WithPagination(
		products,
		"Products retrieved successfully",
		page,
		perPage,
		total,
	))
}

// GetByID
// @Summary Get a product by ID
// @Description Get detailed information about a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response_formatter.Response{data=entity.ProductResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.product.GetByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid product ID",
			[]string{err.Error()},
		))
	}

	product, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get product", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to get product",
			[]string{err.Error()},
		))
	}

	return c.JSON(http.StatusOK, response_formatter.Success(product, "Product retrieved successfully"))
}

// Update
// @Summary Update a product
// @Description Update a product's information by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body entity.UpdateProductRequest true "Product update request"
// @Success 200 {object} response_formatter.Response{data=entity.ProductResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.product.Update")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid product ID",
			[]string{err.Error()},
		))
	}

	var req entity.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid request body",
			[]string{err.Error()},
		))
	}

	if err := h.validate.Validate(ctx, req); err != nil {
		validationErrors := h.validate.ExtractValidationErrors(err)
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, ve.Message)
		}
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Validation failed",
			errorMessages,
		))
	}

	product, err := h.service.Update(ctx, id, req)
	if err != nil {
		h.logger.Error("failed to update product", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to update product",
			[]string{err.Error()},
		))
	}

	h.metrics.RecordProductUpdated(ctx)
	return c.JSON(http.StatusOK, response_formatter.Success(product, "Product updated successfully"))
}

// Delete
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response_formatter.Response
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.product.Delete")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid product ID",
			[]string{err.Error()},
		))
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error("failed to delete product", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "product not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to delete product",
			[]string{err.Error()},
		))
	}

	h.metrics.RecordProductDeleted(ctx)
	return c.JSON(http.StatusOK, response_formatter.Success(nil, "Product deleted successfully"))
}
