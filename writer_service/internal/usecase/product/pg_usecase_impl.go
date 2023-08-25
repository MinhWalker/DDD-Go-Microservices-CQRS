package usecase

import (
	"context"

	"github.com/minhwalker/cqrs-microservices/core/models"
	kafkaClient "github.com/minhwalker/cqrs-microservices/core/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/core/pkg/utils"
	kafkaMessages "github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	models2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/repositories"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/usecase"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/mappers"

	"github.com/opentracing/opentracing-go"
)

type productUsecase struct {
	log               logger.Logger
	cfg               *config.Config
	productRepository repositories.IProductRepositoryWriter
	kafkaProducer     kafkaClient.Producer
}

func NewProductUsecase(log logger.Logger, cfg *config.Config, productRepository repositories.IProductRepositoryWriter, kafkaProducer kafkaClient.Producer) usecase.IProductUsecase {
	return &productUsecase{log: log, cfg: cfg, productRepository: productRepository, kafkaProducer: kafkaProducer}
}

func (u *productUsecase) Create(ctx context.Context, command *models2.CreateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
	defer span.Finish()

	productDto := &models2.Product{Product: models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
	}}

	product, err := u.productRepository.CreateProduct(ctx, productDto)
	if err != nil {
		return err
	}

	msg := &kafkaMessages.ProductCreated{Product: mappers.ProductToGrpcMessage(product)}

	return utils.PublishKafkaMessage(ctx, span.Context(), u.kafkaProducer, msg, u.cfg.KafkaTopics.ProductCreated.TopicName)
}

func (u *productUsecase) Delete(ctx context.Context, command *models2.DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
	defer span.Finish()

	if err := u.productRepository.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	msg := &kafkaMessages.ProductDeleted{ProductID: command.ProductID.String()}

	return utils.PublishKafkaMessage(ctx, span.Context(), u.kafkaProducer, msg, u.cfg.KafkaTopics.ProductDeleted.TopicName)
}

func (u *productUsecase) Update(ctx context.Context, command *models2.UpdateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductHandler.Handle")
	defer span.Finish()

	productDto := &models2.Product{Product: models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
	}}

	product, err := u.productRepository.UpdateProduct(ctx, productDto)
	if err != nil {
		return err
	}

	msg := &kafkaMessages.ProductUpdated{Product: mappers.ProductToGrpcMessage(product)}

	return utils.PublishKafkaMessage(ctx, span.Context(), u.kafkaProducer, msg, u.cfg.KafkaTopics.ProductUpdated.TopicName)
}

func (u *productUsecase) GetByID(ctx context.Context, query *models2.GetProductByIdQuery) (*models2.Product, error) {
	return u.productRepository.GetProductById(ctx, query.ProductID)
}
