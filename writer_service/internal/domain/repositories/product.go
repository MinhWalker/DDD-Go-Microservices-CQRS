package repositories

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	uuid "github.com/satori/go.uuid"
)

type IProductRepositoryWriter interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
}
