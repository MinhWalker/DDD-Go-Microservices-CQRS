package usecase

import (
	"context"
	"time"

	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	kafkaMessages "github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	models2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/mappers"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

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
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   u.cfg.KafkaTopics.ProductUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return u.kafkaProducer.PublishMessage(ctx, message)
}
