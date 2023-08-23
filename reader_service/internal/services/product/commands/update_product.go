package commands

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/repository"

	"github.com/minhwalker/cqrs-microservices/models"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
)

type UpdateProductCmdHandler interface {
	Handle(ctx context.Context, command *UpdateProductCommand) error
}

type updateProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.Repository
	redisRepo repository.CacheRepository
	pgRepo    repository.Repository
}

func NewUpdateProductCmdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository, pgRepo repository.Repository) *updateProductCmdHandler {
	return &updateProductCmdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *updateProductCmdHandler) Handle(ctx context.Context, command *UpdateProductCommand) error {
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
