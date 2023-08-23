package queries

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	dto "github.com/minhwalker/cqrs-microservices/reader_service/internal/dto/product"

	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/repository/product"
)

type SearchProductHandler interface {
	Handle(ctx context.Context, query *dto.SearchProductQuery) (*models.ProductsList, error)
}

type searchProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.RepositoryReader
	redisRepo repository.CacheRepository
	pgRepo    repository.RepositoryReader
}

func NewSearchProductHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.RepositoryReader, redisRepo repository.CacheRepository, pgRepo repository.RepositoryReader) *searchProductHandler {
	return &searchProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo, pgRepo: pgRepo}
}

func (s *searchProductHandler) Handle(ctx context.Context, query *dto.SearchProductQuery) (*models.ProductsList, error) {
	return s.pgRepo.Search(ctx, query.Text, query.Pagination)
}
