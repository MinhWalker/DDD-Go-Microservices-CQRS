package commands

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	dto "github.com/minhwalker/cqrs-microservices/reader_service/internal/dto/product"

	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"

	uuid "github.com/satori/go.uuid"
)

type UpdateProductCmdHandler interface {
	Handle(ctx context.Context, command *dto.UpdateProductCommand) error
}

type updateProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.RepositoryReader
	redisRepo repository.CacheRepository
	pgRepo    repository.RepositoryReader
}

func NewUpdateProductCmdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.RepositoryReader, redisRepo repository.CacheRepository, pgRepo repository.RepositoryReader) *updateProductCmdHandler {
	return &updateProductCmdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *updateProductCmdHandler) Handle(ctx context.Context, command *dto.UpdateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductCmdHandler.Handle")
	defer span.Finish()

	parsedUUID, err := uuid.FromString(command.ProductID)
	if err != nil {
		return err
	}

	product := &models.Product{
		ProductID:   parsedUUID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		UpdatedAt:   command.UpdatedAt,
	}

	updated, err := c.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	c.redisRepo.PutProduct(ctx, updated.ProductID.String(), updated)
	return nil
}
