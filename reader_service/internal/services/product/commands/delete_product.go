package commands

import (
	"context"

	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"

	"github.com/opentracing/opentracing-go"
)

type DeleteProductCmdHandler interface {
	Handle(ctx context.Context, command *DeleteProductCommand) error
}

type deleteProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo product.Repository
	redisRepo product.CacheRepository
	pgRepo    product.Repository
}

func NewDeleteProductCmdHandler(log logger.Logger, cfg *config.Config, mongoRepo product.Repository, redisRepo product.CacheRepository, pgRepo product.Repository) *deleteProductCmdHandler {
	return &deleteProductCmdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *deleteProductCmdHandler) Handle(ctx context.Context, command *DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductCmdHandler.Handle")
	defer span.Finish()

	if err := c.pgRepo.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	c.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}
