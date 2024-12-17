// @title Unnispick K-Style API
// @version 1.0
// @description This is the API server for Unnispick Korean beauty and fashion products
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email admin@k-stylehub.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

package docs

import (
	"github.com/swaggo/swag"
)

// Brand represents a brand in the system
// swagger:model Brand
type Brand struct {
	// The unique identifier for the brand
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// The name of the brand
	// example: SANGCLI
	// required: true
	// min length: 1
	// max length: 255
	BrandName string `json:"brand_name"`

	// When the brand was created
	// example: 2024-01-01T00:00:00Z
	CreatedAt string `json:"created_at"`

	// When the brand was last updated
	// example: 2024-01-01T00:00:00Z
	UpdatedAt string `json:"updated_at"`
}

// Product represents a product in the system
// swagger:model Product
type Product struct {
	// The unique identifier for the product
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// The name of the product
	// example: SANGCLI - Oat Barrier Sparkling Spa Bath Barm
	// required: true
	// min length: 1
	// max length: 255
	ProductName string `json:"product_name"`

	// The price of the product
	// example: 289000
	// required: true
	// minimum: 0.01
	Price float64 `json:"price"`

	// The quantity in stock
	// example: 50
	// required: true
	// minimum: 0
	Quantity int `json:"quantity"`

	// The ID of the brand this product belongs to
	// example: 123e4567-e89b-12d3-a456-426614174000
	// required: true
	BrandID string `json:"brand_id"`

	// The brand information
	Brand *Brand `json:"brand,omitempty"`

	// When the product was created
	// example: 2024-01-01T00:00:00Z
	CreatedAt string `json:"created_at"`

	// When the product was last updated
	// example: 2024-01-01T00:00:00Z
	UpdatedAt string `json:"updated_at"`
}

// CreateBrandRequest represents the request body for creating a brand
// swagger:model CreateBrandRequest
type CreateBrandRequest struct {
	// The name of the brand
	// example: SANGCLI
	// required: true
	// min length: 1
	// max length: 255
	BrandName string `json:"brand_name"`
}

// UpdateBrandRequest represents the request body for updating a brand
// swagger:model UpdateBrandRequest
type UpdateBrandRequest struct {
	// The name of the brand
	// example: SANGCLI
	// required: true
	// min length: 1
	// max length: 255
	BrandName string `json:"brand_name"`
}

// CreateProductRequest represents the request body for creating a product
// swagger:model CreateProductRequest
type CreateProductRequest struct {
	// The name of the product
	// example: SANGCLI - Oat Barrier Sparkling Spa Bath Barm
	// required: true
	// min length: 1
	// max length: 255
	ProductName string `json:"product_name"`

	// The price of the product
	// example: 159000
	// required: true
	// minimum: 0.01
	Price float64 `json:"price"`

	// The quantity in stock
	// example: 100
	// required: true
	// minimum: 0
	Quantity int `json:"quantity"`

	// The ID of the brand this product belongs to
	// example: 123e4567-e89b-12d3-a456-426614174000
	// required: true
	BrandID string `json:"brand_id"`
}

// Example responses for documentation
const exampleResponses = `
    "responses": {
        "200": {
            "description": "Successful response example",
            "schema": {
                "type": "object",
                "properties": {
                    "data": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "id": {
                                    "type": "string",
                                    "example": "123e4567-e89b-12d3-a456-426614174000"
                                },
                                "product_name": {
                                    "type": "string",
                                    "example": "SANGCLI - Oat Barrier Sparkling Spa Bath Barm"
                                },
                                "price": {
                                    "type": "number",
                                    "example": 249000
                                },
                                "quantity": {
                                    "type": "integer",
                                    "example": 75
                                },
                                "brand": {
                                    "type": "object",
                                    "properties": {
                                        "id": {
                                            "type": "string",
                                            "example": "123e4567-e89b-12d3-a456-426614174000"
                                        },
                                        "brand_name": {
                                            "type": "string",
                                            "example": "SANGCLI"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "meta": {
                        "type": "object",
                        "properties": {
                            "page": {
                                "type": "integer",
                                "example": 1
                            },
                            "per_page": {
                                "type": "integer",
                                "example": 10
                            },
                            "total": {
                                "type": "integer",
                                "example": 100
                            }
                        }
                    }
                }
            }
        }
    }
`

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}
