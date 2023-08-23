package server

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minhwalker/cqrs-microservices/pkg/interceptors"
	kafkaClient "github.com/minhwalker/cqrs-microservices/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/pkg/postgres"
	"github.com/minhwalker/cqrs-microservices/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/writer_service/config"
	kafkaConsumer "github.com/minhwalker/cqrs-microservices/writer_service/internal/delivery/kafka"
	product3 "github.com/minhwalker/cqrs-microservices/writer_service/internal/metrics"
	product2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/repository/product"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/services/product"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log       logger.Logger
	cfg       *config.Config
	v         *validator.Validate
	kafkaConn *kafka.Conn
	ps        *product.ProductService
	im        interceptors.InterceptorManager
	pgConn    *pgxpool.Pool
	metrics   *product3.WriterServiceMetrics
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)
	s.metrics = product3.NewWriterServiceMetrics(s.cfg)

	pgxConn, err := postgres.NewPgxConn(s.cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn")
	}
	s.pgConn = pgxConn
	s.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())
	defer pgxConn.Close()

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck

	productRepo := product2.NewProductRepository(s.log, s.cfg, pgxConn)
	s.ps = product.NewProductService(s.log, s.cfg, productRepo, kafkaProducer)
	productMessageProcessor := kafkaConsumer.NewProductMessageProcessor(s.log, s.cfg, s.v, s.ps, s.metrics)

	s.log.Info("Starting Writer Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), kafkaConsumer.PoolSize, productMessageProcessor.ProcessMessages)

	closeGrpcServer, grpcServer, err := s.newWriterGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	s.runHealthCheck(ctx)
	s.runMetrics(cancel)

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	<-ctx.Done()
	grpcServer.GracefulStop()

	return nil
}
