package server

import (
	"context"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	v1 "github.com/minhwalker/cqrs-microservices/reader_service/internal/delivery/http/v1"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"

	"github.com/minhwalker/cqrs-microservices/core/pkg/interceptors"
	kafkaClient "github.com/minhwalker/cqrs-microservices/core/pkg/kafka"
	"github.com/minhwalker/cqrs-microservices/core/pkg/logger"
	"github.com/minhwalker/cqrs-microservices/core/pkg/mongodb"
	"github.com/minhwalker/cqrs-microservices/core/pkg/postgres"
	redisClient "github.com/minhwalker/cqrs-microservices/core/pkg/redis"
	"github.com/minhwalker/cqrs-microservices/core/pkg/tracing"
	"github.com/minhwalker/cqrs-microservices/reader_service/config"
	readerKafka "github.com/minhwalker/cqrs-microservices/reader_service/internal/delivery/kafka"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/domain/usecase"
	"github.com/minhwalker/cqrs-microservices/reader_service/internal/metrics"
	product2 "github.com/minhwalker/cqrs-microservices/reader_service/internal/repositories/product"
	usecase2 "github.com/minhwalker/cqrs-microservices/reader_service/internal/usecase"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	log         logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	kafkaConn   *kafka.Conn
	im          interceptors.InterceptorManager
	mongoClient *mongo.Client
	redisClient redis.UniversalClient
	pgConn      *pgxpool.Pool
	ps          usecase.IProductUsecase
	metrics     *metrics.ReaderServiceMetrics
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)
	s.metrics = metrics.NewReaderServiceMetrics(s.cfg)

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, s.cfg.Mongo)
	if err != nil {
		return errors.Wrap(err, "NewMongoDBConn")
	}
	s.mongoClient = mongoDBConn
	defer mongoDBConn.Disconnect(ctx) // nolint: errcheck

	pgxConn, err := postgres.NewPgxConn(s.cfg.Postgresql)
	if err != nil {
		return errors.Wrap(err, "postgresql.NewPgxConn")
	}
	s.pgConn = pgxConn
	s.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())
	defer pgxConn.Close()

	s.log.Infof("Mongo connected: %v", mongoDBConn.NumberSessionsInProgress())

	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close() // nolint: errcheck
	s.log.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	productRepo := product2.NewProductRepository(s.log, s.cfg, pgxConn)
	mongoRepo := product2.NewMongoRepository(s.log, s.cfg, s.mongoClient)
	redisRepo := product2.NewRedisRepository(s.log, s.cfg, s.redisClient)

	s.ps = usecase2.NewProductUsecase(s.log, s.cfg, mongoRepo, redisRepo, productRepo)

	readerMessageProcessor := readerKafka.NewReaderMessageProcessor(s.log, s.cfg, s.v, s.ps, s.metrics)

	s.log.Info("Starting Reader Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), readerKafka.PoolSize, readerMessageProcessor.ProcessMessages)

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	router := fasthttprouter.New()

	productHandlers := v1.NewReaderHttpService(router, s.log, s.cfg, s.v, s.ps, s.metrics)
	productHandlers.MapRoutes()

	port := s.cfg.Http.Port
	addr := fmt.Sprintf("%s", port)

	fmt.Printf("Server is listening on %s\n", addr)
	if err := fasthttp.ListenAndServe(addr, router.Handler); err != nil {
		s.log.Fatalf("Error starting server: %s", err)
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

	closeGrpcServer, grpcServer, err := s.newReaderGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}