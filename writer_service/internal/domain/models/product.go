package models

import (
	"github.com/minhwalker/cqrs-microservices/core/models"
	uuid "github.com/satori/go.uuid"
)

type Product struct {
	models.Product
}

// Grpc
type CreateProductCommand struct {
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

type UpdateProductCommand struct {
	ProductID   uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

type GetProductByIdQuery struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

// HTTP
type CreateProductRequest struct {
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	ProductID   uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}
