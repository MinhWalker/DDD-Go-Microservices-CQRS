package services

import (
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/services/product/commands"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/services/product/queries"
	"github.com/minhwalker/cqrs-microservices/repository"
)

type ProductService struct {
	Commands *commands.ProductCommands
	Queries  *queries.ProductQueries
}

func NewProductService(
	log logger.Logger,
	cfg *config.Config,
	mongoRepo repository.Repository,
	redisRepo repository.CacheRepository,
	pgRepo repository.Repository,
) *ProductService {

	createProductHandler := commands.NewCreateProductHandler(log, cfg, mongoRepo, redisRepo, pgRepo)
	deleteProductCmdHandler := commands.NewDeleteProductCmdHandler(log, cfg, mongoRepo, redisRepo, pgRepo)
	updateProductCmdHandler := commands.NewUpdateProductCmdHandler(log, cfg, mongoRepo, redisRepo, pgRepo)

	getProductByIdHandler := queries.NewGetProductByIdHandler(log, cfg, mongoRepo, redisRepo, pgRepo)
	searchProductHandler := queries.NewSearchProductHandler(log, cfg, mongoRepo, redisRepo, pgRepo)

	productCommands := commands.NewProductCommands(createProductHandler, updateProductCmdHandler, deleteProductCmdHandler)
	productQueries := queries.NewProductQueries(getProductByIdHandler, searchProductHandler)

	return &ProductService{Commands: productCommands, Queries: productQueries}
}
