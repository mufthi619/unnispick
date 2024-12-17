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

type BrandHandler struct {
	service  entity.BrandService
	logger   *zap.Logger
	tracer   *tracing.Tracer
	metrics  *metrics.Metrics
	validate *validator.Validator
}

func NewBrandHandler(
	service entity.BrandService,
	logger *zap.Logger,
	tracer *tracing.Tracer,
	metrics *metrics.Metrics,
	validate *validator.Validator,
) *BrandHandler {
	return &BrandHandler{
		service:  service,
		logger:   logger,
		tracer:   tracer,
		metrics:  metrics,
		validate: validate,
	}
}

// Create
// @Summary Create a new brand
// @Description Create a new brand with the provided information
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body entity.CreateBrandRequest true "Brand creation request"
// @Success 201 {object} response_formatter.Response{data=entity.BrandResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /brands [post]
func (h *BrandHandler) Create(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.brand.Create")
	defer span.End()

	var req entity.CreateBrandRequest
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

	brand, err := h.service.Create(ctx, req)
	if err != nil {
		h.logger.Error("failed to create brand", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response_formatter.Error(
			http.StatusInternalServerError,
			"Failed to create brand",
			[]string{err.Error()},
		))
	}

	h.metrics.RecordBrandCreated(ctx)
	return c.JSON(http.StatusCreated, response_formatter.Created(brand, "Brand created successfully"))
}

// GetAll
// @Summary Get all brands with pagination
// @Description Get a list of all brands with pagination support
// @Tags brands
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Param search query string false "Search term for brand name"
// @Success 200 {object} response_formatter.Response{data=[]entity.BrandResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /brands [get]
func (h *BrandHandler) GetAll(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.brand.GetAll")
	defer span.End()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	search := c.QueryParam("search")

	page, perPage = response_formatter.ValidatePagination(page, perPage)
	filter := entity.BrandFilterRequest{
		Search:  search,
		Page:    page,
		PerPage: perPage,
	}

	brands, total, err := h.service.GetAll(ctx, filter)
	if err != nil {
		h.logger.Error("failed to get brands", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, response_formatter.Error(
			http.StatusInternalServerError,
			"Failed to get brands",
			[]string{err.Error()},
		))
	}

	return c.JSON(http.StatusOK, response_formatter.WithPagination(
		brands,
		"Brands retrieved successfully",
		page,
		perPage,
		total,
	))
}

// GetByID
// @Summary Get a brand by ID
// @Description Get detailed information about a brand by its ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} response_formatter.Response{data=entity.BrandResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /brands/{id} [get]
func (h *BrandHandler) GetByID(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.brand.GetByID")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid brand ID",
			[]string{err.Error()},
		))
	}

	brand, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get brand", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "brand not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to get brand",
			[]string{err.Error()},
		))
	}

	return c.JSON(http.StatusOK, response_formatter.Success(brand, "Brand retrieved successfully"))
}

// Update
// @Summary Update a brand
// @Description Update a brand's information by its ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param brand body entity.UpdateBrandRequest true "Brand update request"
// @Success 200 {object} response_formatter.Response{data=entity.BrandResponse}
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /brands/{id} [put]
func (h *BrandHandler) Update(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.brand.Update")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid brand ID",
			[]string{err.Error()},
		))
	}

	var req entity.UpdateBrandRequest
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

	brand, err := h.service.Update(ctx, id, req)
	if err != nil {
		h.logger.Error("failed to update brand", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "brand not found" {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to update brand",
			[]string{err.Error()},
		))
	}

	return c.JSON(http.StatusOK, response_formatter.Success(brand, "Brand updated successfully"))
}

// Delete
// @Summary Delete a brand
// @Description Delete a brand by its ID
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} response_formatter.Response
// @Failure 400 {object} response_formatter.Response
// @Failure 404 {object} response_formatter.Response
// @Failure 500 {object} response_formatter.Response
// @Router /brands/{id} [delete]
func (h *BrandHandler) Delete(c echo.Context) error {
	ctx, span := h.tracer.StartFromEcho(c, "handler.brand.Delete")
	defer span.End()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response_formatter.Error(
			http.StatusBadRequest,
			"Invalid brand ID",
			[]string{err.Error()},
		))
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error("failed to delete brand", zap.Error(err))
		statusCode := http.StatusInternalServerError
		if err.Error() == "brand not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "cannot delete brand: still has associated products" {
			statusCode = http.StatusBadRequest
		}
		return c.JSON(statusCode, response_formatter.Error(
			statusCode,
			"Failed to delete brand",
			[]string{err.Error()},
		))
	}

	h.metrics.RecordBrandDeleted(ctx)
	return c.JSON(http.StatusOK, response_formatter.Success(nil, "Brand deleted successfully"))
}
