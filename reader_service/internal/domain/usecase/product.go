package usecase

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/models"
)

type IProductUsecase interface {
	Create(ctx context.Context, command *models.CreateProductCommand) error
	Update(ctx context.Context, command *models.UpdateProductCommand) error
	Delete(ctx context.Context, command *models.DeleteProductCommand) error

	GetProductById(ctx context.Context, query *models.GetProductByIdQuery) (*models.Product, error)
	Search(ctx context.Context, query *models.SearchProductQuery) (*models.ProductsList, error)
}
