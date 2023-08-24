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

type CreateProductCmdHandler interface {
	Handle(ctx context.Context, command *dto.CreateProductCommand) error
}

type createProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.RepositoryReader
	redisRepo repository.CacheRepository
	pgRepo    repository.RepositoryReader
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.RepositoryReader, redisRepo repository.CacheRepository, pgRepo repository.RepositoryReader) CreateProductCmdHandler {
	return &createProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (c *createProductHandler) Handle(ctx context.Context, command *dto.CreateProductCommand) error {
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