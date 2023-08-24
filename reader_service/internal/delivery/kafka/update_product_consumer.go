package kafka

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/dto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func (s *readerMessageProcessor) processProductUpdated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.UpdateProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processProductUpdated")
	defer span.Finish()

	msg := &kafkaMessages.ProductUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	p := msg.GetProduct()
	command := dto.NewUpdateProductCommand(p.GetProductID(), p.GetName(), p.GetDescription(), p.GetPrice(), p.GetUpdatedAt().AsTime())
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
