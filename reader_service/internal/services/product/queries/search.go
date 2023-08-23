package queries

import (
	"context"

	"github.com/minhwalker/cqrs-microservices/models"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"
)

type SearchProductHandler interface {
	Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error)
}

type searchProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo product.Repository
	redisRepo product.CacheRepository
	pgRepo    product.Repository
}

func NewSearchProductHandler(log logger.Logger, cfg *config.Config, mongoRepo product.Repository, redisRepo product.CacheRepository, pgRepo product.Repository) *searchProductHandler {
	return &searchProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (s *searchProductHandler) Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error) {
	return s.pgRepo.Search(ctx, query.Text, query.Pagination)
}
