package usecase

import (
	"context"
	kafkaClient "github.com/minhwalker/cqrs-microservices/core/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/repositories"
)

type IProductUsecase interface {
	Create(ctx context.Context, command *models.CreateProductCommand) error
	Delete(ctx context.Context, command *models.DeleteProductCommand) error
	Update(ctx context.Context, command *models.UpdateProductCommand) error

	GetByID(ctx context.Context, query *models.GetProductByIdQuery) (*models.Product, error)
}

type productUsecase struct {
	log               logger.Logger
	cfg               *config.Config
	productRepository repositories.IProductRepositoryWriter
	kafkaProducer     kafkaClient.Producer
}

func NewProductUsecase(log logger.Logger, cfg *config.Config, productRepository repositories.IProductRepositoryWriter, kafkaProducer kafkaClient.Producer) IProductUsecase {
	return &productUsecase{log: log, cfg: cfg, productRepository: productRepository, kafkaProducer: kafkaProducer}
}
