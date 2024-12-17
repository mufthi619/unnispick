package service

import (
	"Unnispick/internal/domain/entity"
	"Unnispick/internal/infra/tracing"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type productService struct {
	repo      entity.ProductRepository
	brandRepo entity.BrandRepository
	logger    *zap.Logger
	tracer    *tracing.Tracer
}

func NewProductService(repo entity.ProductRepository, brandRepo entity.BrandRepository, logger *zap.Logger, tracer *tracing.Tracer) entity.ProductService {
	return &productService{
		repo:      repo,
		brandRepo: brandRepo,
		logger:    logger,
		tracer:    tracer,
	}
}

func (s *productService) Create(ctx context.Context, req entity.CreateProductRequest) (*entity.ProductResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.product.Create")
	defer span.End()

	// Check if brand exists
	exists, err := s.brandRepo.ExistsByID(ctx, req.BrandID)
	if err != nil {
		s.logger.Error("failed to check brand existence", zap.Error(err))
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("brand with ID %s not found", req.BrandID)
	}

	// Check if product name already exists
	exists, err = s.repo.ExistsByName(ctx, req.ProductName)
	if err != nil {
		s.logger.Error("failed to check product existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("product with name %s already exists", req.ProductName)
	}

	product := req.ToProductEntity()
	if err := s.repo.Create(ctx, product); err != nil {
		s.logger.Error("failed to create product", zap.Error(err))
		return nil, err
	}

	return s.toResponse(product), nil
}

func (s *productService) GetByID(ctx context.Context, id uuid.UUID) (*entity.ProductResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.product.GetByID")
	defer span.End()

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err))
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	return s.toResponse(product), nil
}

func (s *productService) GetAll(ctx context.Context, filter entity.ProductFilterRequest) ([]entity.ProductResponse, int64, error) {
	ctx, span := s.tracer.Start(ctx, "service.product.GetAll")
	defer span.End()

	products, count, err := s.repo.GetAllWithFilter(ctx, filter.ToProductFilterRepo())
	if err != nil {
		s.logger.Error("failed to get products", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]entity.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *s.toResponse(&product)
	}

	return responses, count, nil
}

func (s *productService) Update(ctx context.Context, id uuid.UUID, req entity.UpdateProductRequest) (*entity.ProductResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.product.Update")
	defer span.End()

	// Check if product exists
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err))
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Check if brand exists if brand ID is being updated
	if req.BrandID != uuid.Nil && req.BrandID != product.BrandID {
		exists, err := s.brandRepo.ExistsByID(ctx, req.BrandID)
		if err != nil {
			s.logger.Error("failed to check brand existence", zap.Error(err))
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("brand with ID %s not found", req.BrandID)
		}
	}

	// Check if new name conflicts with existing product
	if req.ProductName != product.ProductName {
		exists, err := s.repo.ExistsByName(ctx, req.ProductName)
		if err != nil {
			s.logger.Error("failed to check product existence", zap.Error(err))
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("product with name %s already exists", req.ProductName)
		}
	}

	product.UpdateFromRequest(req)
	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error("failed to update product", zap.Error(err))
		return nil, err
	}

	return s.toResponse(product), nil
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "service.product.Delete")
	defer span.End()

	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to check product existence", zap.Error(err))
		return err
	}
	if !exists {
		return fmt.Errorf("product not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete product", zap.Error(err))
		return err
	}

	return nil
}

func (s *productService) toResponse(product *entity.Product) *entity.ProductResponse {
	response := &entity.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		Price:       product.Price,
		Quantity:    product.Quantity,
		BrandID:     product.BrandID,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}

	if product.Brand != nil {
		response.Brand = &entity.BrandResponse{
			ID:        product.Brand.ID,
			BrandName: product.Brand.BrandName,
			CreatedAt: product.Brand.CreatedAt.Format(time.RFC3339),
			UpdatedAt: product.Brand.UpdatedAt.Format(time.RFC3339),
		}
	}

	return response
}
