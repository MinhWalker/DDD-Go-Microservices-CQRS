package product_usecase

import (
	kafkaClient "github.com/minhwalker/cqrs-microservices/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	commands2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/handler/product_usecase/commands"
	queries2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/handler/product_usecase/queries"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/repository/product"
)

type ProductService struct {
	Commands *commands2.ProductCommands
	Queries  *queries2.ProductQueries
}

func NewProductService(log logger.Logger, cfg *config.Config, pgRepo repository.RepositoryWriter, kafkaProducer kafkaClient.Producer) *ProductService {

	updateProductHandler := commands2.NewUpdateProductHandler(log, cfg, pgRepo, kafkaProducer)
	createProductHandler := commands2.NewCreateProductHandler(log, cfg, pgRepo, kafkaProducer)
	deleteProductHandler := commands2.NewDeleteProductHandler(log, cfg, pgRepo, kafkaProducer)

	getProductByIdHandler := queries2.NewGetProductByIdHandler(log, cfg, pgRepo)

	productCommands := commands2.NewProductCommands(createProductHandler, updateProductHandler, deleteProductHandler)
	productQueries := queries2.NewProductQueries(getProductByIdHandler)

	return &ProductService{Commands: productCommands, Queries: productQueries}
}
