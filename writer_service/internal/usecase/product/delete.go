package usecase

import (
	"context"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	kafkaMessages "github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

func (u *productUsecase) Delete(ctx context.Context, command *models.DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductHandler.Handle")
	defer span.Finish()

	if err := u.productRepository.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	msg := &kafkaMessages.ProductDeleted{ProductID: command.ProductID.String()}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   u.cfg.KafkaTopics.ProductDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	return u.kafkaProducer.PublishMessage(ctx, message)
}
