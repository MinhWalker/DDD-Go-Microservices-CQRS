package mappers

import (
	"github.com/minhwalker/cqrs-microservices/core/models"
	"github.com/minhwalker/cqrs-microservices/core/proto/kafka"
	models2 "github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/dto/proto/product_writer"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToGrpcMessage(product *models2.Product) *kafkaMessages.Product {
	return &kafkaMessages.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func ProductFromGrpcMessage(product *kafkaMessages.Product) (*models2.Product, error) {

	proUUID, err := uuid.FromString(product.GetProductID())
	if err != nil {
		return nil, err
	}

	return &models2.Product{
		models.Product{
			ProductID:   proUUID,
			Name:        product.GetName(),
			Description: product.GetDescription(),
			Price:       product.GetPrice(),
			CreatedAt:   product.GetCreatedAt().AsTime(),
			UpdatedAt:   product.GetUpdatedAt().AsTime(),
		},
	}, nil
}

func WriterProductToGrpc(product *models2.Product) *writerService.Product {
	return &writerService.Product{
		ProductID:   product.ProductID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}
