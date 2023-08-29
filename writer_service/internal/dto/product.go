package dto

import (
	"github.com/minhwalker/cqrs-microservices/writer_service/internal/domain/models"
	uuid "github.com/satori/go.uuid"
)

func NewCreateProductCommand(productID uuid.UUID, name string, description string, price float64) *models.CreateProductCommand {
	return &models.CreateProductCommand{ProductID: productID, Name: name, Description: description, Price: price}
}

func NewUpdateProductCommand(productID uuid.UUID, name string, description string, price float64) *models.UpdateProductCommand {
	return &models.UpdateProductCommand{ProductID: productID, Name: name, Description: description, Price: price}
}

func NewDeleteProductCommand(productID uuid.UUID) *models.DeleteProductCommand {
	return &models.DeleteProductCommand{ProductID: productID}
}

func NewGetProductByIdQuery(productID uuid.UUID) *models.GetProductByIdQuery {
	return &models.GetProductByIdQuery{ProductID: productID}
}
