package usecase

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
)

func (u *productUsecase) GetByID(ctx context.Context, query *models.GetProductByIdQuery) (*models.Product, error) {
	return u.productRepository.GetProductById(ctx, query.ProductID)
}
