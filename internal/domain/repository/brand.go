package repository

import (
	"Unnispick/internal/domain/entity"
	"Unnispick/internal/infra/tracing"
	"Unnispick/pkg/telemetry/tracer"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type brandRepository struct {
	db     *gorm.DB
	tracer *tracing.Tracer
}

func NewBrandRepository(db *gorm.DB, tracer *tracing.Tracer) entity.BrandRepository {
	return &brandRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *brandRepository) Create(ctx context.Context, brand *entity.Brand) error {
	ctx, span := r.tracer.Start(ctx, "repository.brand.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(brand).Error; err != nil {
		tracer.RecordError(span, err)
		return fmt.Errorf("failed to create brand: %w", err)
	}

	return nil
}

func (r *brandRepository) GetAllWithFilter(ctx context.Context, filter entity.BrandFilterRepository) (brands []entity.Brand, count int64, err error) {
	ctx, span := r.tracer.Start(ctx, "repository.brand.GetAllWithFilter")
	defer span.End()

	if filter.Limit < 0 || filter.Offset < 0 {
		return nil, 0, fmt.Errorf("invalid pagination parameters: limit and offset must be non-negative")
	}

	query := r.db.WithContext(ctx).Model(&entity.Brand{})

	// Apply search filter
	if filter.Search != "" {
		query = query.Where("brand_name ILIKE ?", "%"+filter.Search+"%")
	}

	// Count total records
	if err = query.Count(&count).Error; err != nil {
		tracer.RecordError(span, err)
		return nil, 0, fmt.Errorf("failed to count brands: %w", err)
	}

	// Check if offset is beyond total count
	if count > 0 && filter.Offset >= int(count) {
		return []entity.Brand{}, count, nil
	}

	// Get paginated records
	if err = query.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Order("created_at DESC").
		Find(&brands).Error; err != nil {
		tracer.RecordError(span, err)
		return nil, 0, fmt.Errorf("failed to list brands: %w", err)
	}

	return brands, count, nil
}

func (r *brandRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Brand, error) {
	ctx, span := r.tracer.Start(ctx, "repository.brand.GetByID")
	defer span.End()

	var brand entity.Brand
	if err := r.db.WithContext(ctx).First(&brand, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracer.RecordError(span, err)
		return nil, fmt.Errorf("failed to get brand: %w", err)
	}

	return &brand, nil
}

func (r *brandRepository) Update(ctx context.Context, brand *entity.Brand) error {
	ctx, span := r.tracer.Start(ctx, "repository.brand.Update")
	defer span.End()

	result := r.db.WithContext(ctx).Model(brand).Updates(map[string]interface{}{
		"brand_name": brand.BrandName,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		tracer.RecordError(span, result.Error)
		return fmt.Errorf("failed to update brand: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("brand not found")
	}

	return nil
}

func (r *brandRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "repository.brand.Delete")
	defer span.End()

	// Check if brand is used in products
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Count(&count).Error; err != nil {
		tracer.RecordError(span, err)
		return fmt.Errorf("failed to check brand usage: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete brand: still has associated products")
	}

	result := r.db.WithContext(ctx).Delete(&entity.Brand{}, "id = ?", id)
	if result.Error != nil {
		tracer.RecordError(span, result.Error)
		return fmt.Errorf("failed to delete brand: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("brand not found")
	}

	return nil
}

func (r *brandRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "repository.brand.ExistsByID")
	defer span.End()

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&entity.Brand{}).
		Select("1").
		Where("id = ?", id).
		Scan(&exists).Error

	if err != nil {
		tracer.RecordError(span, err)
		return false, fmt.Errorf("failed to check brand existence: %w", err)
	}

	return exists, nil
}

func (r *brandRepository) GetByName(ctx context.Context, name string) (*entity.Brand, error) {
	ctx, span := r.tracer.Start(ctx, "repository.brand.GetByName")
	defer span.End()

	var brand entity.Brand
	if err := r.db.WithContext(ctx).First(&brand, "brand_name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracer.RecordError(span, err)
		return nil, fmt.Errorf("failed to get brand by name: %w", err)
	}

	return &brand, nil
}

func (r *brandRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "repository.brand.ExistsByName")
	defer span.End()

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&entity.Brand{}).
		Select("1").
		Where("brand_name = ?", name).
		Scan(&exists).Error

	if err != nil {
		tracer.RecordError(span, err)
		return false, fmt.Errorf("failed to check brand name existence: %w", err)
	}

	return exists, nil
}
