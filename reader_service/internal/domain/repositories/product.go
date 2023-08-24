package repositories

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/models"
	uuid "github.com/satori/go.uuid"
)

type IProductRepositoryReader interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error

	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error)
}

type IProductCacheRepository interface {
	PutProduct(ctx context.Context, key string, product *models.Product)
	GetProduct(ctx context.Context, key string) (*models.Product, error)
	DelProduct(ctx context.Context, key string)
	DelAllProducts(ctx context.Context)
}
