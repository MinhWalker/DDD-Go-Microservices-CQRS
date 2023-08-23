package commands

import (
	"context"

	"github.com/minhwalker/cqrs-microservices/models"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
)

type CreateProductCmdHandler interface {
	Handle(ctx context.Context, command *CreateProductCommand) error
}

type createProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo product.Repository
	redisRepo product.CacheRepository
	pgRepo    product.Repository
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, mongoRepo product.Repository, redisRepo product.CacheRepository, pgRepo product.Repository) CreateProductCmdHandler {
	return &createProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *createProductHandler) Handle(ctx context.Context, command *CreateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
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
		CreatedAt:   command.CreatedAt,
		UpdatedAt:   command.UpdatedAt,
	}

	created, err := c.pgRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	c.redisRepo.PutProduct(ctx, created.ProductID.String(), created)
	return nil
}
