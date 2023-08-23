package kafka

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/handler"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/metrics"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

type readerMessageProcessor struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      *handler.ProductService
	metrics *metrics.ReaderServiceMetrics
}

func NewReaderMessageProcessor(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *handler.ProductService, metrics *metrics.ReaderServiceMetrics) *readerMessageProcessor {
	return &readerMessageProcessor{log: log, cfg: cfg, v: v, ps: ps, metrics: metrics}
}

func (s *readerMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.ProductCreated.TopicName:
			s.processProductCreated(ctx, r, m)
		case s.cfg.KafkaTopics.ProductUpdated.TopicName:
			s.processProductUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.ProductDeleted.TopicName:
			s.processProductDeleted(ctx, r, m)
		}
	}
}
