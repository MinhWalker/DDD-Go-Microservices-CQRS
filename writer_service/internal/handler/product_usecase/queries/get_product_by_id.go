package queries

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/models"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	dto "github.com/minhwalker/cqrs-microservices/writer_service/internal/dto/product"
	repository "github.com/minhwalker/cqrs-microservices/writer_service/internal/repository/product"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *dto.GetProductByIdQuery) (*models.Product, error)
}

type getProductByIdHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo repository.RepositoryWriter
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, pgRepo repository.RepositoryWriter) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *dto.GetProductByIdQuery) (*models.Product, error) {
	return q.pgRepo.GetProductById(ctx, query.ProductID)
}
