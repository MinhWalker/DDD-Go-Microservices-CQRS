package utils

import (
	"context"
	"time"

	kafkaClient "github.com/minhwalker/cqrs-microservices/core/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"

	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

func PublishKafkaMessage(ctx context.Context, span opentracing.SpanContext, kafkaProducer kafkaClient.Producer, msg proto.Message, topicName string) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Topic:   topicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span),
	}

	return kafkaProducer.PublishMessage(ctx, message)
}
