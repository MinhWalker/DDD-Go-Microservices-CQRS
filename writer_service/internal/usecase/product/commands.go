package product

import (
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/repositories"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/usecase"
)

type productUsecase struct {
	productRepository repositories.IProductRepositoryWriter
}

func New(productRepository repositories.IProductRepositoryWriter) usecase.IProductUsecase {
	return &productUsecase{productRepository: productRepository}
}
