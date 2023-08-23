package commands

import (
	"context"
	dto "github.com/minhwalker/cqrs-microservices/writer_service/internal/dto/product"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/repository/product"
	"github.com/opentracing/opentracing-go"
	"time"

	"github.com/minhwalker/cqrs-microservices/models"
	kafkaClient "github.com/minhwalker/cqrs-microservices/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/pkg/tracing"
	kafkaMessages "github.com/minhwalker/cqrs-microservices/proto/kafka"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	"github.com/minhwalker/cqrs-microservices/writer_service/mappers"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type CreateProductCmdHandler interface {
	Handle(ctx context.Context, command *dto.CreateProductCommand) error
}

type createProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        repository.RepositoryWriter
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, pgRepo repository.RepositoryWriter, kafkaProducer kafkaClient.Producer) *createProductHandler {
	return &createProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *createProductHandler) Handle(ctx context.Context, command *dto.CreateProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "createProductHandler.Handle")
	defer span.Finish()

	productDto := &models.Product{ProductID: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price}

	product, err := c.pgRepo.CreateProduct(ctx, productDto)
	if err != nil {
		return err
	}

	msg := &kafkaMessages.ProductCreated{Product: mappers.ProductToGrpcMessage(product)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   c.cfg.KafkaTopics.ProductCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return c.kafkaProducer.PublishMessage(ctx, message)
}
