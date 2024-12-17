package entity

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type (
	Product struct {
		ID          uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
		ProductName string     `json:"product_name" validate:"required,min=1,max=255" gorm:"column:product_name;type:varchar(255);not null"`
		Price       float64    `json:"price" validate:"required,gt=0" gorm:"type:decimal(15,2);not null;check:price > 0"`
		Quantity    int        `json:"quantity" validate:"required,gte=0" gorm:"type:integer;not null;check:quantity >= 0"`
		BrandID     uuid.UUID  `json:"brand_id" validate:"required,uuid" gorm:"column:brand_id;type:uuid;not null"`
		Brand       *Brand     `json:"brand,omitempty" gorm:"foreignKey:BrandID"`
		CreatedAt   time.Time  `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
		DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index;type:timestamp with time zone"`
	}

	ProductRepository interface {
		Create(ctx context.Context, product *Product) error
		GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
		GetAllWithFilter(ctx context.Context, filter ProductFilterRepository) (products []Product, count int64, err error)
		Update(ctx context.Context, product *Product) error
		Delete(ctx context.Context, id uuid.UUID) error
		ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
		GetByName(ctx context.Context, name string) (*Product, error)
		ExistsByName(ctx context.Context, name string) (bool, error)
	}

	ProductService interface {
		Create(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
		GetByID(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
		GetAll(ctx context.Context, filter ProductFilterRequest) ([]ProductResponse, int64, error)
		Update(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	ProductFilterRequest struct {
		BrandID  uuid.UUID `query:"brand_id"`
		MinPrice float64   `query:"min_price"`
		MaxPrice float64   `query:"max_price"`
		MinQty   int       `query:"min_qty"`
		MaxQty   int       `query:"max_qty"`
		Page     int       `query:"page"`
		PerPage  int       `query:"per_page"`
	}

	ProductFilterRepository struct {
		BrandID  uuid.UUID
		MinPrice float64
		MaxPrice float64
		MinQty   int
		MaxQty   int
		Limit    int
		Offset   int
	}

	CreateProductRequest struct {
		ProductName string    `json:"product_name" validate:"required,min=1,max=255"`
		Price       float64   `json:"price" validate:"required,gt=0"`
		Quantity    int       `json:"quantity" validate:"required,gte=0"`
		BrandID     uuid.UUID `json:"brand_id" validate:"required,uuid"`
	}

	UpdateProductRequest struct {
		ProductName string    `json:"product_name" validate:"required,min=1,max=255"`
		Price       float64   `json:"price" validate:"required,gt=0"`
		Quantity    int       `json:"quantity" validate:"required,gte=0"`
		BrandID     uuid.UUID `json:"brand_id" validate:"required,uuid"`
	}

	ProductResponse struct {
		ID          uuid.UUID      `json:"id"`
		ProductName string         `json:"product_name"`
		Price       float64        `json:"price"`
		Quantity    int            `json:"quantity"`
		BrandID     uuid.UUID      `json:"brand_id"`
		Brand       *BrandResponse `json:"brand,omitempty"`
		CreatedAt   string         `json:"created_at"`
		UpdatedAt   string         `json:"updated_at"`
	}
)

func (*Product) TableName() string {
	return "products"
}

func (req ProductFilterRequest) ToProductFilterRepo() ProductFilterRepository {
	return ProductFilterRepository{
		BrandID:  req.BrandID,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		MinQty:   req.MinQty,
		MaxQty:   req.MaxQty,
		Limit:    req.PerPage,
		Offset:   (req.Page - 1) * req.PerPage,
	}
}

func (req *CreateProductRequest) ToProductEntity() *Product {
	return &Product{
		ProductName: req.ProductName,
		Price:       req.Price,
		Quantity:    req.Quantity,
		BrandID:     req.BrandID,
	}
}

func (p *Product) UpdateFromRequest(req UpdateProductRequest) {
	if req.ProductName != "" {
		p.ProductName = req.ProductName
	}
	if req.Price > 0 {
		p.Price = req.Price
	}
	if req.Quantity >= 0 {
		p.Quantity = req.Quantity
	}
	if req.BrandID != uuid.Nil {
		p.BrandID = req.BrandID
	}
}

func (p *Product) ToResponseDTO() *ProductResponse {
	response := &ProductResponse{
		ID:          p.ID,
		ProductName: p.ProductName,
		Price:       p.Price,
		Quantity:    p.Quantity,
		BrandID:     p.BrandID,
	}

	if p.Brand != nil {
		response.Brand = p.Brand.ToResponseDTO()
	}

	return response
}
