package kafka

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/dto"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *productMessageProcessor) processUpdateProduct(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "productMessageProcessor.processUpdateProduct")
	defer span.Finish()

	msg := &kafkaMessages.ProductUpdate{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	proUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := dto.NewUpdateProductCommand(proUUID, msg.GetName(), msg.GetDescription(), msg.GetPrice())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.ps.Update(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateProduct.Handle", err)
		s.metrics.ErrorKafkaMessages.Inc()
		return
	}

	s.commitMessage(ctx, r, m)
}
