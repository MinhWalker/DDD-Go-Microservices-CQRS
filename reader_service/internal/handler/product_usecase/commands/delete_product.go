package commands

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	dto "github.com/minhwalker/cqrs-microservices/reader_service/internal/dto/product"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"

	"github.com/minhwalker/cqrs-microservices/reader_service/config"
)

type DeleteProductCmdHandler interface {
	Handle(ctx context.Context, command *dto.DeleteProductCommand) error
}

type deleteProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.RepositoryReader
	redisRepo repository.CacheRepository
	pgRepo    repository.RepositoryReader
}

func NewDeleteProductCmdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.RepositoryReader, redisRepo repository.CacheRepository, pgRepo repository.RepositoryReader) *deleteProductCmdHandler {
	return &deleteProductCmdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *deleteProductCmdHandler) Handle(ctx context.Context, command *dto.DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductCmdHandler.Handle")
	defer span.Finish()

	if err := c.pgRepo.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	c.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}
