package product

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
)

func (u *productUsecase) Delete(ctx context.Context, command *models.DeleteProductCommand) error {
	//TODO implement me
	panic("implement me")
}
