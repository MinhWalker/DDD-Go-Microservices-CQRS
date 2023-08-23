package queries

import (
	"context"

	"github.com/minhwalker/cqrs-microservices/models"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"

	"github.com/opentracing/opentracing-go"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error)
}

type getProductByIdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo product.Repository
	redisRepo product.CacheRepository
	pgRepo    product.Repository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, mongoRepo product.Repository, redisRepo product.CacheRepository, pgRepo product.Repository) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error) {
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
