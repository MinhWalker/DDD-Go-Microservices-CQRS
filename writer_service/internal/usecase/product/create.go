package product

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
)

func (u *productUsecase) Create(ctx context.Context, command *models.CreateProductCommand) error {
	_, err := u.productRepository.CreateProduct(ctx, command)
	return err
}
