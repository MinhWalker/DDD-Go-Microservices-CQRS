package queries

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"

	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	dto "github.com/minhwalker/cqrs-microservices/reader_service/internal/dto/product"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *dto.GetProductByIdQuery) (*models.Product, error)
}

type getProductByIdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.RepositoryReader
	redisRepo repository.CacheRepository
	pgRepo    repository.RepositoryReader
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.RepositoryReader, redisRepo repository.CacheRepository, pgRepo repository.RepositoryReader) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *dto.GetProductByIdQuery) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	if product, err := q.redisRepo.GetProduct(ctx, query.ProductID.String()); err == nil && product != nil {
		return product, nil
	}

	product, err := q.pgRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, err
	}

	q.redisRepo.PutProduct(ctx, product.ProductID.String(), product)
	return product, nil
}
