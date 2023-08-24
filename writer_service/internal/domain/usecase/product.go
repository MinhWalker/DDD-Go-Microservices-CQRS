package usecase

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
)

type IProductUsecase interface {
	Create(ctx context.Context, command *models.CreateProductCommand) error
	Delete(ctx context.Context, command *models.DeleteProductCommand) error
	Update(ctx context.Context, command *models.UpdateProductCommand) error

	GetByID(ctx context.Context, query *models.GetProductByIdQuery) (*models.Product, error)
}
