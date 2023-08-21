### Golang CQRS Kafka gRPC Postgresql MongoDB Redis microservices example 👋

#### 👨‍💻 Full list what has been used:
[Kafka](https://github.com/segmentio/kafka-go) as messages broker<br/>
[gRPC](https://github.com/grpc/grpc-go) Go implementation of gRPC<br/>
[PostgreSQL](https://github.com/jackc/pgx) as database<br/>
[Jaeger](https://www.jaegertracing.io/) open source, end-to-end distributed [tracing](https://opentracing.io/)<br/>
[Prometheus](https://prometheus.io/) monitoring and alerting<br/>
[Grafana](https://grafana.com/) for to compose observability dashboards with everything from Prometheus<br/>
[MongoDB](https://github.com/mongodb/mongo-go-driver) Web and API based SMTP testing<br/>
[Redis](https://github.com/go-redis/redis) Type-safe Redis client for Golang<br/>
[swag](https://github.com/swaggo/swag) Swagger for Go<br/>
[Echo](https://github.com/labstack/echo) web framework<br/>

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### Swagger UI:

http://localhost:5001/swagger/index.html


For local development 🙌👨‍💻🚀:

```
make migrate_up // run sql migrations
make mongo // run mongo init scripts
make swagger // generate swagger documentation
make local or docker_dev // for run docker compose files
```

### Project struct:

#### Overview

![Untitled-2023-08-22-0136.svg](..%2F..%2F..%2FDownloads%2Fdiagram%2FUntitled-2023-08-22-0136.svg)

Api_gateway_service
```
api_gateway_service
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config.yaml
└── internal
    ├── client
    │   └── reader_service.go
    ├── dto
    │   ├── create_product.go
    │   ├── product_list_response.go
    │   ├── product_response.go
    │   └── update_product.go
    ├── metrics
    │   └── metrics.go
    ├── middlewares
    │   └── middlewares.go
    ├── products
    │   ├── commands
    │   │   ├── commands.go
    │   │   ├── create_product.go
    │   │   ├── delete_product.go
    │   │   └── update_product.go
    │   ├── delivery
    │   │   └── http
    │   │       └── v1
    │   │           ├── handlers.go
    │   │           └── routes.go
    │   ├── delivery.go
    │   ├── queries
    │   │   ├── get_by_id.go
    │   │   ├── queries.go
    │   │   └── search_product.go
    │   └── service
    │       └── service.go
    └── server
        ├── http.go
        ├── server.go
        └── utils.go
```

reader_service
```
reader_service
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config.yaml
├── internal
│   ├── metrics
│   │   └── metrics.go
│   ├── models
│   │   └── product.go
│   ├── product
│   │   ├── commands
│   │   │   ├── commands.go
│   │   │   ├── create_product.go
│   │   │   ├── delete_product.go
│   │   │   └── update_product.go
│   │   ├── delivery
│   │   │   ├── grpc
│   │   │   │   └── grpc_service.go
│   │   │   └── kafka
│   │   │       ├── consumer_group.go
│   │   │       ├── create_product_consumer.go
│   │   │       ├── delete_product_consumer.go
│   │   │       ├── update_product_consumer.go
│   │   │       └── utils.go
│   │   ├── queries
│   │   │   ├── get_by_id.go
│   │   │   ├── queries.go
│   │   │   └── search.go
│   │   ├── repository
│   │   │   ├── mongo_repository.go
│   │   │   ├── redis_repository.go
│   │   │   └── repository.go
│   │   └── service
│   │       └── service.go
│   └── server
│       ├── grpc_server.go
│       ├── server.go
│       └── utils.go
└── proto
    └── product_reader
        ├── product_reader_grpc.pb.go
        ├── product_reader_messages.pb.go
        ├── product_reader_messages.proto
        ├── product_reader.pb.go
        └── product_reader.proto
```

writer_service
```
writer_service
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config.yaml
├── internal
│   ├── metrics
│   │   └── metrics.go
│   ├── models
│   │   └── product.go
│   ├── product
│   │   ├── commands
│   │   │   ├── commands.go
│   │   │   ├── create_product.go
│   │   │   ├── delete_product.go
│   │   │   └── update_product.go
│   │   ├── delivery
│   │   │   ├── grpc
│   │   │   │   └── grpc_service.go
│   │   │   └── kafka
│   │   │       ├── consumer_group.go
│   │   │       ├── create_product_consumer.go
│   │   │       ├── delete_product_consumer.go
│   │   │       ├── update_product_consumer.go
│   │   │       └── utils.go
│   │   ├── queries
│   │   │   ├── get_product_by_id.go
│   │   │   └── queries.go
│   │   ├── repository
│   │   │   ├── pg_repository.go
│   │   │   ├── repository.go
│   │   │   └── sql_queries.go
│   │   └── service
│   │       └── service.go
│   └── server
│       ├── grpc_server.go
│       ├── server.go
│       └── utils.go
├── mappers
│   └── product_mapper.go
└── proto
    └── product_writer
        ├── product_writer_grpc.pb.go
        ├── product_writer_messages.pb.go
        ├── product_writer_messages.proto
        ├── product_writer.pb.go
        └── product_writer.proto
```