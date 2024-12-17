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

type brandService struct {
	repo   entity.BrandRepository
	logger *zap.Logger
	tracer *tracing.Tracer
}

func NewBrandService(repo entity.BrandRepository, logger *zap.Logger, tracer *tracing.Tracer) entity.BrandService {
	return &brandService{
		repo:   repo,
		logger: logger,
		tracer: tracer,
	}
}

func (s *brandService) Create(ctx context.Context, req entity.CreateBrandRequest) (*entity.BrandResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.brand.Create")
	defer span.End()

	if err := req.Validate(); err != nil {
		s.logger.Error("failed to validate brand request", zap.Error(err))
		return nil, err
	}

	exists, err := s.repo.ExistsByName(ctx, req.BrandName)
	if err != nil {
		s.logger.Error("failed to check brand existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("brand with name %s already exists", req.BrandName)
	}

	brand := req.ToBrandEntity()
	if err := s.repo.Create(ctx, brand); err != nil {
		s.logger.Error("failed to create brand", zap.Error(err))
		return nil, err
	}

	return s.toResponse(brand), nil
}

func (s *brandService) GetByID(ctx context.Context, id uuid.UUID) (*entity.BrandResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.brand.GetByID")
	defer span.End()

	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get brand", zap.Error(err))
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	return s.toResponse(brand), nil
}

func (s *brandService) GetAll(ctx context.Context, filter entity.BrandFilterRequest) ([]entity.BrandResponse, int64, error) {
	ctx, span := s.tracer.Start(ctx, "service.brand.GetAll")
	defer span.End()

	brands, count, err := s.repo.GetAllWithFilter(ctx, filter.ToBrandFilterRepo())
	if err != nil {
		s.logger.Error("failed to get brands", zap.Error(err))
		return nil, 0, err
	}

	responses := make([]entity.BrandResponse, len(brands))
	for i, brand := range brands {
		responses[i] = *s.toResponse(&brand)
	}

	return responses, count, nil
}

func (s *brandService) Update(ctx context.Context, id uuid.UUID, req entity.UpdateBrandRequest) (*entity.BrandResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.brand.Update")
	defer span.End()

	if err := req.Validate(); err != nil {
		s.logger.Error("failed to validate brand request", zap.Error(err))
		return nil, err
	}

	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get brand", zap.Error(err))
		return nil, err
	}
	if brand == nil {
		return nil, fmt.Errorf("brand not found")
	}

	if brand.BrandName != req.BrandName {
		exists, err := s.repo.ExistsByName(ctx, req.BrandName)
		if err != nil {
			s.logger.Error("failed to check brand existence", zap.Error(err))
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("brand with name %s already exists", req.BrandName)
		}
	}

	brand.UpdateFromRequest(req)
	if err := s.repo.Update(ctx, brand); err != nil {
		s.logger.Error("failed to update brand", zap.Error(err))
		return nil, err
	}

	return s.toResponse(brand), nil
}

func (s *brandService) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "service.brand.Delete")
	defer span.End()

	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to check brand existence", zap.Error(err))
		return err
	}
	if !exists {
		return fmt.Errorf("brand not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete brand", zap.Error(err))
		return err
	}

	return nil
}

func (s *brandService) toResponse(brand *entity.Brand) *entity.BrandResponse {
	return &entity.BrandResponse{
		ID:        brand.ID,
		BrandName: brand.BrandName,
		CreatedAt: brand.CreatedAt.Format(time.RFC3339),
		UpdatedAt: brand.UpdatedAt.Format(time.RFC3339),
	}
}
