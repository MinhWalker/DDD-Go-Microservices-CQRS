package server

import (
	"context"
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"github.com/minhwalker/cqrs-microservices/core/pkg/constants"
	kafkaClient "github.com/minhwalker/cqrs-microservices/core/pkg/kafka"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	stackSize = 1 << 10 // 1 KB
)

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.log.Infof("established new kafka controller connection: %s", controllerURI)

	productCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductCreate.ReplicationFactor,
	}

	productCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductCreated.ReplicationFactor,
	}

	productUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductUpdate.ReplicationFactor,
	}

	productUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductUpdated.ReplicationFactor,
	}

	productDeleteTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductDelete.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductDelete.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductDelete.ReplicationFactor,
	}

	productDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ProductDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ProductDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ProductDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		productCreateTopic,
		productUpdateTopic,
		productCreatedTopic,
		productUpdatedTopic,
		productDeleteTopic,
		productDeletedTopic,
	); err != nil {
		s.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{productCreateTopic, productUpdateTopic, productCreatedTopic, productUpdatedTopic, productDeleteTopic, productDeletedTopic})
}

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.ProductCreate.TopicName,
		s.cfg.KafkaTopics.ProductUpdate.TopicName,
		s.cfg.KafkaTopics.ProductDelete.TopicName,
	}
}

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck(constants.Postgres, healthcheck.AsyncWithContext(ctx, func() error {
		return s.pgConn.Ping(ctx)
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
		_, err := s.kafkaConn.Brokers()
		if err != nil {
			return err
		}
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	go func() {
		s.log.Infof("Writer microservice Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := &fasthttp.Server{
		Handler: s.metricsHandler,
	}

	go func() {
		addr := fmt.Sprintf(":%s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.ListenAndServe(addr); err != nil {
			s.log.Printf("metricsServer.ListenAndServe: %v", err)
			cancel()
		}
	}()

	s.log.Printf("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
}

func (s *server) metricsHandler(ctx *fasthttp.RequestCtx) {
	// Your middleware/recovery logic here
	// Note: FastHTTP doesn't have a built-in middleware system like Echo
	// You need to implement the middleware logic manually if needed

	// Serve Prometheus metrics using the fasthttpadaptor and promhttp packages
	if string(ctx.Path()) == s.cfg.Probes.PrometheusPath {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	}
}
