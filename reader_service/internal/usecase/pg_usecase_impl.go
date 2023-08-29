package usecase

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	models2 "github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/repositories"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/usecase"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
)

type productUsecase struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repositories.IProductRepositoryReader
	redisRepo repositories.IProductCacheRepository
	pgRepo    repositories.IProductRepositoryReader
}

func NewProductUsecase(
	log logger.Logger,
	cfg *config.Config,
	mongoRepo repositories.IProductRepositoryReader,
	redisRepo repositories.IProductCacheRepository,
	pgRepo repositories.IProductRepositoryReader,
) usecase.IProductUsecase {
	return &productUsecase{
		log:       log,
		cfg:       cfg,
		mongoRepo: mongoRepo,
		redisRepo: redisRepo,
		pgRepo:    pgRepo,
	}
}

func (p *productUsecase) Create(ctx context.Context, command *models2.CreateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
	defer span.Finish()

	parsedUUID, err := uuid.FromString(command.ProductID)
	if err != nil {
		return err
	}

	product := &models2.Product{
		Product: models.Product{
			ProductID:   parsedUUID,
			Name:        command.Name,
			Description: command.Description,
			Price:       command.Price,
			CreatedAt:   command.CreatedAt,
			UpdatedAt:   command.UpdatedAt,
		},
	}

	created, err := p.pgRepo.CreateProduct(ctx, product)
	if err != nil {
		return err
	}

	p.redisRepo.PutProduct(ctx, created.ProductID.String(), created)
	return nil
}

func (p *productUsecase) Update(ctx context.Context, command *models2.UpdateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductCmdHandler.Handle")
	defer span.Finish()

	parsedUUID, err := uuid.FromString(command.ProductID)
	if err != nil {
		return err
	}

	product := &models2.Product{
		Product: models.Product{
			ProductID:   parsedUUID,
			Name:        command.Name,
			Description: command.Description,
			Price:       command.Price,
			UpdatedAt:   command.UpdatedAt,
		},
	}

	updated, err := p.pgRepo.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	p.redisRepo.PutProduct(ctx, updated.ProductID.String(), updated)
	return nil
}

func (p *productUsecase) Delete(ctx context.Context, command *models2.DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductCmdHandler.Handle")
	defer span.Finish()

	if err := p.pgRepo.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	p.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}

func (p *productUsecase) GetProductById(ctx context.Context, query *models2.GetProductByIdQuery) (*models2.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	if product, err := p.redisRepo.GetProduct(ctx, query.ProductID.String()); err == nil && product != nil {
		return product, nil
	}

	product, err := p.pgRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, err
	}

	p.redisRepo.PutProduct(ctx, product.ProductID.String(), product)
	return product, nil
}

func (p *productUsecase) Search(ctx context.Context, query *models2.SearchProductQuery) (*models2.ProductsList, error) {
	return p.pgRepo.Search(ctx, query.Text, query.Pagination)
}
