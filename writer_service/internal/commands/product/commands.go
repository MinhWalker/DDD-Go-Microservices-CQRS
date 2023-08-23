package product

import (
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/repositories"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/usecase"
)

type productCommands struct {
	productRepository repositories.IProductRepositoryWriter
}

func NewProductCommands(productRepository repositories.IProductRepositoryWriter) usecase.IProductUsecase {
	return &productCommands{productRepository: productRepository}
}

func (u *productCommands) Create(ctx context.Context, command *models.CreateProductCommand) error {
	//TODO implement me
	panic("implement me")
}

func (u *productCommands) Update(ctx context.Context, command *models.UpdateProductCommand) error {
	//TODO implement me
	panic("implement me")
}

//func (c *deleteProductHandler) Handle(ctx context.Context, command *dto.DeleteProductCommand) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
//	defer span.Finish()
//
//	if err := c.pgRepo.DeleteProduct(ctx, command.ProductID); err != nil {
//		return err
//	}
//
//	msg := &kafkaMessages.ProductDeleted{ProductID: command.ProductID.String()}
//	msgBytes, err := proto.Marshal(msg)
//	if err != nil {
//		return err
//	}
//
//	message := kafkaClient.Message{
//		Topic:   c.cfg.KafkaTopics.ProductDeleted.TopicName,
//		Value:   msgBytes,
//		Time:    time.Now().UTC(),
//		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
//	}
//
//	return c.kafkaProducer.PublishMessage(ctx, message)
//}

//func (c *createProductHandler) Handle(ctx context.Context, command *dto.CreateProductCommand) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
//	defer span.Finish()
//
//	productDto := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price}
//
//	product, err := c.pgRepo.CreateProduct(ctx, productDto)
//	if err != nil {
//		return err
//	}
//
//	msg := &kafkaMessages.ProductCreated{Product: mappers.ProductToGrpcMessage(product)}
//	msgBytes, err := proto.Marshal(msg)
//	if err != nil {
//		return err
//	}
//
//	message := kafkaClient.Message{
//		Topic:   c.cfg.KafkaTopics.ProductCreated.TopicName,
//		Value:   msgBytes,
//		Time:    time.Now().UTC(),
//		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
//	}
//
//	return c.kafkaProducer.PublishMessage(ctx, message)
//}

//func (c *updateProductHandler) Handle(ctx context.Context, command *dto.UpdateProductCommand) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "updateProductHandler.Handle")
//	defer span.Finish()
//
//	productDto := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price}
//
//	product, err := c.pgRepo.UpdateProduct(ctx, productDto)
//	if err != nil {
//		return err
//	}
//
//	msg := &kafkaMessages.ProductUpdated{Product: mappers.ProductToGrpcMessage(product)}
//	msgBytes, err := proto.Marshal(msg)
//	if err != nil {
//		return err
//	}
//
//	message := kafkaClient.Message{
//		Topic:   c.cfg.KafkaTopics.ProductUpdated.TopicName,
//		Value:   msgBytes,
//		Time:    time.Now().UTC(),
//		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
//	}
//
//	return c.kafkaProducer.PublishMessage(ctx, message)
//}
