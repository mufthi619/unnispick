package entity

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	Brand struct {
		ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
		BrandName string     `json:"brand_name" validate:"required,min=1,max=255" gorm:"column:brand_name;type:varchar(255);not null"`
		CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
		UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index;type:timestamp with time zone"`
		Products  []Product  `json:"products,omitempty" gorm:"foreignKey:BrandID"`
	}

	BrandRepository interface {
		Create(ctx context.Context, brand *Brand) error
		GetByID(ctx context.Context, id uuid.UUID) (*Brand, error)
		GetAllWithFilter(ctx context.Context, filter BrandFilterRepository) (brands []Brand, count int64, err error)
		Update(ctx context.Context, brand *Brand) error
		Delete(ctx context.Context, id uuid.UUID) error
		ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
		GetByName(ctx context.Context, name string) (*Brand, error)
		ExistsByName(ctx context.Context, name string) (bool, error)
	}

	BrandService interface {
		Create(ctx context.Context, req CreateBrandRequest) (*BrandResponse, error)
		GetByID(ctx context.Context, id uuid.UUID) (*BrandResponse, error)
		GetAll(ctx context.Context, filter BrandFilterRequest) ([]BrandResponse, int64, error)
		Update(ctx context.Context, id uuid.UUID, req UpdateBrandRequest) (*BrandResponse, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	BrandFilterRequest struct {
		Search  string `query:"search"`
		Page    int    `query:"page"`
		PerPage int    `query:"per_page"`
	}

	BrandFilterRepository struct {
		Search string
		Limit  int
		Offset int
	}

	CreateBrandRequest struct {
		BrandName string `json:"brand_name" validate:"required,min=1,max=255"`
	}

	UpdateBrandRequest struct {
		BrandName string `json:"brand_name" validate:"required,min=1,max=255"`
	}

	BrandResponse struct {
		ID        uuid.UUID `json:"id"`
		BrandName string    `json:"brand_name"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
	}
)

func (*Brand) TableName() string {
	return "brands"
}

func (req BrandFilterRequest) ToBrandFilterRepo() BrandFilterRepository {
	return BrandFilterRepository{
		Search: req.Search,
		Limit:  req.PerPage,
		Offset: (req.Page - 1) * req.PerPage,
	}
}

func (req *CreateBrandRequest) ToBrandEntity() *Brand {
	return &Brand{
		BrandName: req.BrandName,
	}
}

func (b *Brand) UpdateFromRequest(req UpdateBrandRequest) {
	if req.BrandName != "" {
		b.BrandName = req.BrandName
	}
}

func (b *Brand) ToResponseDTO() *BrandResponse {
	return &BrandResponse{
		ID:        b.ID,
		BrandName: b.BrandName,
	}
}

func (req *CreateBrandRequest) Validate() error {
	if req.BrandName == "" {
		return ErrEmptyBrandName
	}
	return nil
}

func (req *UpdateBrandRequest) Validate() error {
	if req.BrandName == "" {
		return ErrEmptyBrandName
	}
	return nil
}
