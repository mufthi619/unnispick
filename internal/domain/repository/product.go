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

type productRepository struct {
	db     *gorm.DB
	tracer *tracing.Tracer
}

func NewProductRepository(db *gorm.DB, tracer *tracing.Tracer) entity.ProductRepository {
	return &productRepository{
		db:     db,
		tracer: tracer,
	}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	ctx, span := r.tracer.Start(ctx, "repository.product.Create")
	defer span.End()

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		tracer.RecordError(span, err)
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	ctx, span := r.tracer.Start(ctx, "repository.product.GetByID")
	defer span.End()

	var product entity.Product
	if err := r.db.WithContext(ctx).
		Preload("Brand").
		First(&product, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracer.RecordError(span, err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *productRepository) GetAllWithFilter(ctx context.Context, filter entity.ProductFilterRepository) (products []entity.Product, count int64, err error) {
	if filter.Limit < 0 || filter.Offset < 0 {
		return nil, 0, fmt.Errorf("invalid pagination parameters: limit and offset must be non-negative")
	}

	query := r.db.WithContext(ctx).Model(&entity.Product{})

	// Apply filters
	if filter.BrandID != uuid.Nil {
		query = query.Where("brand_id = ?", filter.BrandID)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.MinQty > 0 {
		query = query.Where("quantity >= ?", filter.MinQty)
	}
	if filter.MaxQty > 0 {
		query = query.Where("quantity <= ?", filter.MaxQty)
	}

	// Count total records
	if err = query.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Check if offset is beyond total count
	if count > 0 && filter.Offset >= int(count) {
		return []entity.Product{}, count, nil
	}

	// Get paginated records
	if err = query.
		Preload("Brand").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, count, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	ctx, span := r.tracer.Start(ctx, "repository.product.Update")
	defer span.End()

	result := r.db.WithContext(ctx).Model(product).Updates(map[string]interface{}{
		"product_name": product.ProductName,
		"price":        product.Price,
		"quantity":     product.Quantity,
		"brand_id":     product.BrandID,
		"updated_at":   time.Now(),
	})

	if result.Error != nil {
		tracer.RecordError(span, result.Error)
		return fmt.Errorf("failed to update product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "repository.product.Delete")
	defer span.End()

	result := r.db.WithContext(ctx).Delete(&entity.Product{}, "id = ?", id)
	if result.Error != nil {
		tracer.RecordError(span, result.Error)
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *productRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "repository.product.ExistsByID")
	defer span.End()

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Select("1").
		Where("id = ?", id).
		Scan(&exists).Error

	if err != nil {
		tracer.RecordError(span, err)
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}

	return exists, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) (*entity.Product, error) {
	ctx, span := r.tracer.Start(ctx, "repository.product.GetByName")
	defer span.End()

	var product entity.Product
	if err := r.db.WithContext(ctx).
		Preload("Brand").
		First(&product, "product_name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		tracer.RecordError(span, err)
		return nil, fmt.Errorf("failed to get product by name: %w", err)
	}

	return &product, nil
}

func (r *productRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "repository.product.ExistsByName")
	defer span.End()

	var exists bool
	err := r.db.WithContext(ctx).
		Model(&entity.Product{}).
		Select("1").
		Where("product_name = ?", name).
		Scan(&exists).Error

	if err != nil {
		tracer.RecordError(span, err)
		return false, fmt.Errorf("failed to check product name existence: %w", err)
	}

	return exists, nil
}
